package postgres

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/internal/repository/postgres/sqlcgen"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/pagination"
)

// uniqueViolation is the PostgreSQL SQLSTATE for a unique_violation. We map the
// users_email unique index hit to a domain-level EMAIL_TAKEN conflict.
const uniqueViolation = "23505"

// UserRepo is the sqlc-backed implementation of repository.UserRepository.
type UserRepo struct {
	q    *sqlcgen.Queries
	pool *pgxpool.Pool
}

// compile-time assertion: UserRepo satisfies the domain interface.
var _ repository.UserRepository = (*UserRepo)(nil)

// NewUserRepo builds a UserRepo over a pgx pool.
func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{q: sqlcgen.New(pool), pool: pool}
}

// Create inserts a user. A duplicate email collapses to apperror.Conflict so
// callers never have to inspect SQLSTATE.
func (r *UserRepo) Create(ctx context.Context, u *entity.User) error {
	row, err := r.q.CreateUser(ctx, sqlcgen.CreateUserParams{
		Email:           u.Email,
		Phone:           u.Phone,
		PasswordHash:    u.PasswordHash,
		Role:            string(u.Role),
		FirstName:       u.FirstName,
		LastName:        u.LastName,
		AvatarUrl:       u.AvatarURL,
		IsActive:        u.IsActive,
		IsEmailVerified: u.IsEmailVerified,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolation {
			return apperror.Conflict("EMAIL_TAKEN", err)
		}
		return apperror.Internal(err)
	}
	// Reflect DB-assigned fields (id, timestamps, defaults) back onto the entity.
	*u = *userFromRow(row)
	return nil
}

// FindByID returns the user or apperror.NotFound. Soft-deleted rows are excluded
// by the query.
func (r *UserRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	row, err := r.q.GetUserByID(ctx, toPgUUID(id))
	if err != nil {
		return nil, mapUserErr(err)
	}
	return userFromRow(row), nil
}

// FindByEmail returns the user or apperror.NotFound.
func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	row, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, mapUserErr(err)
	}
	return userFromRow(row), nil
}

// Update persists the mutable user fields (the query itself filters
// deleted_at IS NULL).
func (r *UserRepo) Update(ctx context.Context, u *entity.User) error {
	row, err := r.q.UpdateUser(ctx, sqlcgen.UpdateUserParams{
		ID:                     toPgUUID(u.ID),
		Phone:                  u.Phone,
		PasswordHash:           u.PasswordHash,
		FirstName:              u.FirstName,
		LastName:               u.LastName,
		AvatarUrl:              u.AvatarURL,
		IsActive:               u.IsActive,
		IsEmailVerified:        u.IsEmailVerified,
		EmailVerificationToken: u.EmailVerificationToken,
		PasswordResetToken:     u.PasswordResetToken,
		PasswordResetExpiresAt: toPgTimestamptz(u.PasswordResetExpiresAt),
	})
	if err != nil {
		return mapUserErr(err)
	}
	*u = *userFromRow(row)
	return nil
}

// SoftDelete sets deleted_at; the row is never physically removed
// (Non-Negotiable #5).
func (r *UserRepo) SoftDelete(ctx context.Context, id uuid.UUID) error {
	if err := r.q.SoftDeleteUser(ctx, toPgUUID(id)); err != nil {
		return apperror.Internal(err)
	}
	return nil
}

// SetLastLogin stamps the last successful login time.
func (r *UserRepo) SetLastLogin(ctx context.Context, id uuid.UUID, t time.Time) error {
	err := r.q.SetLastLogin(ctx, sqlcgen.SetLastLoginParams{
		ID:          toPgUUID(id),
		LastLoginAt: pgTimestamptzValue(t),
	})
	if err != nil {
		return apperror.Internal(err)
	}
	return nil
}

// mapUserErr translates a query error: no rows → NotFound, anything else →
// Internal.
func mapUserErr(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return apperror.NotFound("user not found", err)
	}
	return apperror.Internal(err)
}

// userFromRow maps a sqlc AuthUser to the domain entity.
func userFromRow(row sqlcgen.AuthUser) *entity.User {
	return &entity.User{
		ID:                     fromPgUUID(row.ID),
		Email:                  row.Email,
		Phone:                  row.Phone,
		PasswordHash:           row.PasswordHash,
		Role:                   entity.Role(row.Role),
		FirstName:              row.FirstName,
		LastName:               row.LastName,
		AvatarURL:              row.AvatarUrl,
		IsActive:               row.IsActive,
		IsEmailVerified:        row.IsEmailVerified,
		EmailVerificationToken: row.EmailVerificationToken,
		PasswordResetToken:     row.PasswordResetToken,
		PasswordResetExpiresAt: fromPgTimestamptz(row.PasswordResetExpiresAt),
		LastLoginAt:            fromPgTimestamptz(row.LastLoginAt),
		CreatedAt:              fromPgTimestamptzValue(row.CreatedAt),
		UpdatedAt:              fromPgTimestamptzValue(row.UpdatedAt),
		DeletedAt:              fromPgTimestamptz(row.DeletedAt),
	}
}

// List returns users matching filter (admin use), newest first.
func (r *UserRepo) List(ctx context.Context, filter repository.AdminUserFilter, cursor *pagination.Cursor, limit int) ([]*entity.User, *pagination.Cursor, error) {
	limit = pagination.NormalizeLimit(limit)

	var (
		where []string
		args  []any
	)
	where = append(where, "deleted_at IS NULL")

	if filter.Role != nil && *filter.Role != "" {
		args = append(args, *filter.Role)
		where = append(where, "role = $"+itoa(len(args)))
	}
	if filter.IsActive != nil {
		args = append(args, *filter.IsActive)
		where = append(where, "is_active = $"+itoa(len(args)))
	}
	if filter.Email != nil && *filter.Email != "" {
		args = append(args, "%"+strings.ToLower(*filter.Email)+"%")
		where = append(where, "LOWER(email) LIKE $"+itoa(len(args)))
	}
	if cursor != nil {
		args = append(args, pgTimestamptzValue(cursor.CreatedAt), toPgUUID(cursor.ID))
		where = append(where, "(created_at, id) < ($"+itoa(len(args)-1)+", $"+itoa(len(args))+")")
	}

	args = append(args, int32(limit+1))
	query := `
SELECT id, email, phone, password_hash, role, first_name, last_name, avatar_url,
       is_active, is_email_verified, email_verification_token,
       password_reset_token, password_reset_expires_at, last_login_at,
       created_at, updated_at, deleted_at
FROM auth.users
WHERE ` + strings.Join(where, " AND ") + `
ORDER BY created_at DESC, id DESC
LIMIT $` + itoa(len(args))

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, nil, apperror.Internal(err)
	}
	defer rows.Close()

	users := make([]*entity.User, 0, limit+1)
	for rows.Next() {
		var row sqlcgen.AuthUser
		if err := rows.Scan(
			&row.ID,
			&row.Email,
			&row.Phone,
			&row.PasswordHash,
			&row.Role,
			&row.FirstName,
			&row.LastName,
			&row.AvatarUrl,
			&row.IsActive,
			&row.IsEmailVerified,
			&row.EmailVerificationToken,
			&row.PasswordResetToken,
			&row.PasswordResetExpiresAt,
			&row.LastLoginAt,
			&row.CreatedAt,
			&row.UpdatedAt,
			&row.DeletedAt,
		); err != nil {
			return nil, nil, apperror.Internal(err)
		}
		users = append(users, userFromRow(row))
	}
	if err := rows.Err(); err != nil {
		return nil, nil, apperror.Internal(err)
	}

	var next *pagination.Cursor
	if len(users) > limit {
		last := users[limit-1]
		next = &pagination.Cursor{CreatedAt: last.CreatedAt, ID: last.ID}
		users = users[:limit]
	}
	return users, next, nil
}

// CountByRole returns the number of active (non-deleted, non-suspended) users
// with the given role. Used by the admin guard that prevents last-admin demotion.
func (r *UserRepo) CountByRole(ctx context.Context, role entity.Role) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM auth.users WHERE role = $1 AND is_active = true AND deleted_at IS NULL`,
		string(role),
	).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

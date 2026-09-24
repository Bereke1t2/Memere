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

// UserRepo is the postgres-backed implementation of repository.UserRepository.
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
	approvalStatus := string(u.ApprovalStatus)
	if approvalStatus == "" {
		if u.Role == entity.RoleStudent {
			approvalStatus = string(entity.ApprovalStatusPending)
		} else {
			approvalStatus = string(entity.ApprovalStatusApproved)
		}
	}

	query := `
INSERT INTO auth.users (
    email, phone, password_hash, role, approval_status, first_name, last_name, avatar_url,
    is_active, is_email_verified
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
RETURNING id, email, phone, password_hash, role, approval_status, first_name, last_name, avatar_url,
          is_active, is_email_verified, email_verification_token, password_reset_token,
          password_reset_expires_at, last_login_at, created_at, updated_at, deleted_at;`

	var scanned entity.User
	var (
		id                     uuid.UUID
		roleStr                string
		apprStr                string
		phone                  *string
		avatar                 *string
		emailVerificationToken *string
		passwordResetToken     *string
		passwordResetExpiresAt *time.Time
		lastLoginAt            *time.Time
		deletedAt              *time.Time
	)

	err := r.pool.QueryRow(ctx, query,
		u.Email,
		u.Phone,
		u.PasswordHash,
		string(u.Role),
		approvalStatus,
		u.FirstName,
		u.LastName,
		u.AvatarURL,
		u.IsActive,
		u.IsEmailVerified,
	).Scan(
		&id,
		&scanned.Email,
		&phone,
		&scanned.PasswordHash,
		&roleStr,
		&apprStr,
		&scanned.FirstName,
		&scanned.LastName,
		&avatar,
		&scanned.IsActive,
		&scanned.IsEmailVerified,
		&emailVerificationToken,
		&passwordResetToken,
		&passwordResetExpiresAt,
		&lastLoginAt,
		&scanned.CreatedAt,
		&scanned.UpdatedAt,
		&deletedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolation {
			return apperror.Conflict("EMAIL_TAKEN", err)
		}
		return apperror.Internal(err)
	}

	scanned.ID = id
	scanned.Role = entity.Role(roleStr)
	scanned.ApprovalStatus = entity.ApprovalStatus(apprStr)
	scanned.Phone = phone
	scanned.AvatarURL = avatar
	scanned.EmailVerificationToken = emailVerificationToken
	scanned.PasswordResetToken = passwordResetToken
	scanned.PasswordResetExpiresAt = passwordResetExpiresAt
	scanned.LastLoginAt = lastLoginAt
	scanned.DeletedAt = deletedAt

	*u = scanned
	return nil
}

// FindByID returns the user or apperror.NotFound. Soft-deleted rows are excluded
// by the query.
func (r *UserRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	query := `
SELECT id, email, phone, password_hash, role, approval_status, first_name, last_name, avatar_url,
       is_active, is_email_verified, email_verification_token, password_reset_token,
       password_reset_expires_at, last_login_at, created_at, updated_at, deleted_at
FROM auth.users
WHERE id = $1 AND deleted_at IS NULL;`

	return r.scanUser(r.pool.QueryRow(ctx, query, id))
}

// FindByEmail returns the user or apperror.NotFound.
func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `
SELECT id, email, phone, password_hash, role, approval_status, first_name, last_name, avatar_url,
       is_active, is_email_verified, email_verification_token, password_reset_token,
       password_reset_expires_at, last_login_at, created_at, updated_at, deleted_at
FROM auth.users
WHERE email = $1 AND deleted_at IS NULL;`

	return r.scanUser(r.pool.QueryRow(ctx, query, email))
}

// Update persists the mutable user fields (the query itself filters deleted_at IS NULL).
func (r *UserRepo) Update(ctx context.Context, u *entity.User) error {
	approvalStatus := string(u.ApprovalStatus)
	if approvalStatus == "" {
		approvalStatus = string(entity.ApprovalStatusPending)
	}

	query := `
UPDATE auth.users
SET phone = $2,
    password_hash = $3,
    role = $4,
    approval_status = $5,
    first_name = $6,
    last_name = $7,
    avatar_url = $8,
    is_active = $9,
    is_email_verified = $10,
    email_verification_token = $11,
    password_reset_token = $12,
    password_reset_expires_at = $13
WHERE id = $1 AND deleted_at IS NULL
RETURNING id, email, phone, password_hash, role, approval_status, first_name, last_name, avatar_url,
          is_active, is_email_verified, email_verification_token, password_reset_token,
          password_reset_expires_at, last_login_at, created_at, updated_at, deleted_at;`

	scanned, err := r.scanUser(r.pool.QueryRow(ctx, query,
		u.ID,
		u.Phone,
		u.PasswordHash,
		string(u.Role),
		approvalStatus,
		u.FirstName,
		u.LastName,
		u.AvatarURL,
		u.IsActive,
		u.IsEmailVerified,
		u.EmailVerificationToken,
		u.PasswordResetToken,
		u.PasswordResetExpiresAt,
	))
	if err != nil {
		return err
	}
	*u = *scanned
	return nil
}

// SoftDelete sets deleted_at; the row is never physically removed
// (Non-Negotiable #5).
func (r *UserRepo) SoftDelete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE auth.users SET deleted_at = now() WHERE id = $1 AND deleted_at IS NULL`, id)
	if err != nil {
		return apperror.Internal(err)
	}
	return nil
}

// SetLastLogin stamps the last successful login time.
func (r *UserRepo) SetLastLogin(ctx context.Context, id uuid.UUID, t time.Time) error {
	_, err := r.pool.Exec(ctx, `UPDATE auth.users SET last_login_at = $2 WHERE id = $1 AND deleted_at IS NULL`, id, t)
	if err != nil {
		return apperror.Internal(err)
	}
	return nil
}

// scanUser scans a single row into an entity.User.
func (r *UserRepo) scanUser(row pgx.Row) (*entity.User, error) {
	var (
		u                      entity.User
		roleStr                string
		apprStr                string
		phone                  *string
		avatar                 *string
		emailVerificationToken *string
		passwordResetToken     *string
		passwordResetExpiresAt *time.Time
		lastLoginAt            *time.Time
		deletedAt              *time.Time
	)

	err := row.Scan(
		&u.ID,
		&u.Email,
		&phone,
		&u.PasswordHash,
		&roleStr,
		&apprStr,
		&u.FirstName,
		&u.LastName,
		&avatar,
		&u.IsActive,
		&u.IsEmailVerified,
		&emailVerificationToken,
		&passwordResetToken,
		&passwordResetExpiresAt,
		&lastLoginAt,
		&u.CreatedAt,
		&u.UpdatedAt,
		&deletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NotFound("user not found", err)
		}
		return nil, apperror.Internal(err)
	}

	u.Role = entity.Role(roleStr)
	u.ApprovalStatus = entity.ApprovalStatus(apprStr)
	u.Phone = phone
	u.AvatarURL = avatar
	u.EmailVerificationToken = emailVerificationToken
	u.PasswordResetToken = passwordResetToken
	u.PasswordResetExpiresAt = passwordResetExpiresAt
	u.LastLoginAt = lastLoginAt
	u.DeletedAt = deletedAt
	return &u, nil
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
	if filter.ApprovalStatus != nil && *filter.ApprovalStatus != "" {
		args = append(args, *filter.ApprovalStatus)
		where = append(where, "approval_status = $"+itoa(len(args)))
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
		args = append(args, cursor.CreatedAt, cursor.ID)
		where = append(where, "(created_at, id) < ($"+itoa(len(args)-1)+", $"+itoa(len(args))+")")
	}

	args = append(args, int32(limit+1))
	query := `
SELECT id, email, phone, password_hash, role, approval_status, first_name, last_name, avatar_url,
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
		var (
			u                      entity.User
			roleStr                string
			apprStr                string
			phone                  *string
			avatar                 *string
			emailVerificationToken *string
			passwordResetToken     *string
			passwordResetExpiresAt *time.Time
			lastLoginAt            *time.Time
			deletedAt              *time.Time
		)
		if err := rows.Scan(
			&u.ID,
			&u.Email,
			&phone,
			&u.PasswordHash,
			&roleStr,
			&apprStr,
			&u.FirstName,
			&u.LastName,
			&avatar,
			&u.IsActive,
			&u.IsEmailVerified,
			&emailVerificationToken,
			&passwordResetToken,
			&passwordResetExpiresAt,
			&lastLoginAt,
			&u.CreatedAt,
			&u.UpdatedAt,
			&deletedAt,
		); err != nil {
			return nil, nil, apperror.Internal(err)
		}
		u.Role = entity.Role(roleStr)
		u.ApprovalStatus = entity.ApprovalStatus(apprStr)
		u.Phone = phone
		u.AvatarURL = avatar
		u.EmailVerificationToken = emailVerificationToken
		u.PasswordResetToken = passwordResetToken
		u.PasswordResetExpiresAt = passwordResetExpiresAt
		u.LastLoginAt = lastLoginAt
		u.DeletedAt = deletedAt
		users = append(users, &u)
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

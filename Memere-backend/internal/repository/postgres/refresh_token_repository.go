package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/internal/repository/postgres/sqlcgen"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// RefreshTokenRepo is the sqlc-backed implementation of
// repository.RefreshTokenRepository. It stores only the SHA-256 hash of each
// refresh token; the raw token is never persisted (spec §7.1).
type RefreshTokenRepo struct {
	q *sqlcgen.Queries
}

var _ repository.RefreshTokenRepository = (*RefreshTokenRepo)(nil)

// NewRefreshTokenRepo builds a RefreshTokenRepo over a pgx pool.
func NewRefreshTokenRepo(pool *pgxpool.Pool) *RefreshTokenRepo {
	return &RefreshTokenRepo{q: sqlcgen.New(pool)}
}

// Create persists a hashed refresh token and reflects DB-assigned fields back
// onto the entity.
func (r *RefreshTokenRepo) Create(ctx context.Context, t *entity.RefreshToken) error {
	row, err := r.q.CreateRefreshToken(ctx, sqlcgen.CreateRefreshTokenParams{
		UserID:     toPgUUID(t.UserID),
		TokenHash:  t.TokenHash,
		DeviceInfo: t.DeviceInfo,
		ExpiresAt:  pgTimestamptzValue(t.ExpiresAt),
	})
	if err != nil {
		return apperror.Internal(err)
	}
	*t = *refreshTokenFromRow(row)
	return nil
}

// FindByHash looks up a token by its SHA-256 hash, returning apperror.NotFound
// when absent.
func (r *RefreshTokenRepo) FindByHash(ctx context.Context, hash string) (*entity.RefreshToken, error) {
	row, err := r.q.GetRefreshTokenByHash(ctx, hash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NotFound("refresh token not found", err)
		}
		return nil, apperror.Internal(err)
	}
	return refreshTokenFromRow(row), nil
}

// Revoke marks a single token revoked (idempotent: already-revoked rows are
// untouched by the query's WHERE clause).
func (r *RefreshTokenRepo) Revoke(ctx context.Context, id uuid.UUID) error {
	if err := r.q.RevokeRefreshToken(ctx, toPgUUID(id)); err != nil {
		return apperror.Internal(err)
	}
	return nil
}

// RevokeAllForUser revokes every active token for a user (used on password
// change / full logout).
func (r *RefreshTokenRepo) RevokeAllForUser(ctx context.Context, userID uuid.UUID) error {
	if err := r.q.RevokeAllRefreshTokensForUser(ctx, toPgUUID(userID)); err != nil {
		return apperror.Internal(err)
	}
	return nil
}

// DeleteExpired physically removes tokens past their expiry — these carry no
// user data, so the soft-delete rule does not apply (cleanup job).
func (r *RefreshTokenRepo) DeleteExpired(ctx context.Context) error {
	if err := r.q.DeleteExpiredRefreshTokens(ctx); err != nil {
		return apperror.Internal(err)
	}
	return nil
}

// refreshTokenFromRow maps a sqlc AuthRefreshToken to the domain entity.
func refreshTokenFromRow(row sqlcgen.AuthRefreshToken) *entity.RefreshToken {
	return &entity.RefreshToken{
		ID:         fromPgUUID(row.ID),
		UserID:     fromPgUUID(row.UserID),
		TokenHash:  row.TokenHash,
		DeviceInfo: row.DeviceInfo,
		ExpiresAt:  fromPgTimestamptzValue(row.ExpiresAt),
		RevokedAt:  fromPgTimestamptz(row.RevokedAt),
		CreatedAt:  fromPgTimestamptzValue(row.CreatedAt),
	}
}

// Package redis holds Redis-backed implementations of the domain repository
// interfaces.
package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	goredis "github.com/redis/go-redis/v9"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// SessionRepo implements repository.SessionRepository over Redis.
type SessionRepo struct {
	client *goredis.Client
}

var _ repository.SessionRepository = (*SessionRepo)(nil)

// NewSessionRepo builds a SessionRepo over a go-redis client.
func NewSessionRepo(client *goredis.Client) *SessionRepo {
	return &SessionRepo{client: client}
}

// sessionKey is the per-user session key. The stored value is the SHA-256 hash
// of the active refresh token — never the raw token.
func sessionKey(userID uuid.UUID) string {
	return fmt.Sprintf("session:%s", userID.String())
}

// SetSession stores the token hash with the given TTL, replacing any existing
// session for the user.
func (r *SessionRepo) SetSession(ctx context.Context, userID uuid.UUID, tokenHash string, ttl time.Duration) error {
	if err := r.client.Set(ctx, sessionKey(userID), tokenHash, ttl).Err(); err != nil {
		return apperror.Internal(err)
	}
	return nil
}

// GetSession returns the stored hash, or "" when no session exists (a Redis
// miss is a normal, non-error outcome — the caller falls back to Postgres).
func (r *SessionRepo) GetSession(ctx context.Context, userID uuid.UUID) (string, error) {
	hash, err := r.client.Get(ctx, sessionKey(userID)).Result()
	if errors.Is(err, goredis.Nil) {
		return "", nil
	}
	if err != nil {
		return "", apperror.Internal(err)
	}
	return hash, nil
}

// DeleteSession removes the user's session. Deleting a missing key is a no-op.
func (r *SessionRepo) DeleteSession(ctx context.Context, userID uuid.UUID) error {
	if err := r.client.Del(ctx, sessionKey(userID)).Err(); err != nil {
		return apperror.Internal(err)
	}
	return nil
}

// --- Account lockout ---

func loginFailKey(userID uuid.UUID) string {
	return fmt.Sprintf("auth:fail:%s", userID.String())
}

// IncrLoginFailure increments the per-account failed-login counter and starts
// the lockout TTL on the first failure. Returns the new count.
func (r *SessionRepo) IncrLoginFailure(ctx context.Context, userID uuid.UUID, lockoutTTL time.Duration) (int64, error) {
	key := loginFailKey(userID)
	count, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, apperror.Internal(err)
	}
	if count == 1 {
		_ = r.client.Expire(ctx, key, lockoutTTL).Err()
	}
	return count, nil
}

// ClearLoginFailures resets the counter after a successful login.
func (r *SessionRepo) ClearLoginFailures(ctx context.Context, userID uuid.UUID) error {
	if err := r.client.Del(ctx, loginFailKey(userID)).Err(); err != nil {
		return apperror.Internal(err)
	}
	return nil
}

// IsLockedOut returns true when the stored failure count >= maxFailures.
func (r *SessionRepo) IsLockedOut(ctx context.Context, userID uuid.UUID, maxFailures int) (bool, error) {
	val, err := r.client.Get(ctx, loginFailKey(userID)).Int64()
	if errors.Is(err, goredis.Nil) {
		return false, nil
	}
	if err != nil {
		return false, apperror.Internal(err)
	}
	return val >= int64(maxFailures), nil
}

// --- Access-token JTI denylist ---

func denyKey(jti string) string { return "auth:deny:" + jti }

// DenyToken adds the JWT ID to the revocation set. The TTL should match the
// token's remaining lifetime so the key auto-expires when the token does.
func (r *SessionRepo) DenyToken(ctx context.Context, jti string, ttl time.Duration) error {
	if err := r.client.Set(ctx, denyKey(jti), "1", ttl).Err(); err != nil {
		return apperror.Internal(err)
	}
	return nil
}

// IsTokenDenied returns true when the JTI is on the denylist.
func (r *SessionRepo) IsTokenDenied(ctx context.Context, jti string) (bool, error) {
	err := r.client.Get(ctx, denyKey(jti)).Err()
	if errors.Is(err, goredis.Nil) {
		return false, nil
	}
	if err != nil {
		return false, apperror.Internal(err)
	}
	return true, nil
}

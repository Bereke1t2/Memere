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

package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// SessionRepository is the fast-path store for a user's current refresh-token
// hash, keyed by user id (key format `session:{user_id}`). It mirrors the
// authoritative record in Postgres so the common refresh path can validate a
// token without a database round trip (spec §7.1). A miss here is not
// authoritative — the caller falls back to Postgres.
type SessionRepository interface {
	// SetSession stores the refresh-token hash for the user with the given TTL,
	// overwriting any previous session (single active session per user).
	SetSession(ctx context.Context, userID uuid.UUID, tokenHash string, ttl time.Duration) error
	// GetSession returns the stored hash, or ("", nil) when no session exists.
	GetSession(ctx context.Context, userID uuid.UUID) (string, error)
	// DeleteSession removes the user's session (logout).
	DeleteSession(ctx context.Context, userID uuid.UUID) error
}

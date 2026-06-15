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

	// --- Account lockout (Phase 6 Skill 2, §7.3) ---

	// IncrLoginFailure increments the per-account failed-login counter and sets
	// its TTL to lockoutTTL on the first increment. Returns the new count.
	IncrLoginFailure(ctx context.Context, userID uuid.UUID, lockoutTTL time.Duration) (int64, error)
	// ClearLoginFailures resets the counter after a successful login.
	ClearLoginFailures(ctx context.Context, userID uuid.UUID) error
	// IsLockedOut returns true when the account is currently in lockout state.
	IsLockedOut(ctx context.Context, userID uuid.UUID, maxFailures int) (bool, error)

	// --- Access-token JTI denylist (Phase 6 Skill 2, §7.3) ---

	// DenyToken adds a JWT ID to the revocation set until the given TTL expires.
	DenyToken(ctx context.Context, jti string, ttl time.Duration) error
	// IsTokenDenied returns true when the JTI is on the denylist.
	IsTokenDenied(ctx context.Context, jti string) (bool, error)
}

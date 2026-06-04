package entity

import (
	"time"

	"github.com/google/uuid"
)

// RefreshToken is a persisted, hashed refresh token (spec §7.1). The raw token
// is never stored — only its hash. RevokedAt being non-nil means the token can
// no longer be used to mint access tokens.
type RefreshToken struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	TokenHash  string
	DeviceInfo *string
	ExpiresAt  time.Time
	RevokedAt  *time.Time
	CreatedAt  time.Time
}

// IsActive reports whether the token is neither revoked nor expired at now.
func (t RefreshToken) IsActive(now time.Time) bool {
	return t.RevokedAt == nil && now.Before(t.ExpiresAt)
}

// Package jwt issues and verifies the access/refresh JSON Web Tokens used by the
// auth vertical (spec §7.1). Access tokens are short-lived bearer credentials;
// refresh tokens are long-lived and additionally tracked server-side by a hash
// of the raw token — the raw refresh token is never persisted.
package jwt

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// Token types carried in the custom TokenType claim. Verification rejects a
// token presented for the wrong purpose (an access token cannot be refreshed,
// and a refresh token cannot authorize a request).
const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

// Claims is the JWT payload: the standard registered claims plus the user
// identity and token purpose the rest of the system needs.
type Claims struct {
	gojwt.RegisteredClaims
	UserID    uuid.UUID `json:"uid"`
	Role      string    `json:"role"`
	TokenType string    `json:"typ"`
}

// Manager signs and verifies tokens with a single HMAC secret. It is safe for
// concurrent use.
type Manager struct {
	secret     []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
	issuer     string
	// now is injectable for tests; defaults to time.Now.
	now func() time.Time
}

// NewManager builds a Manager. secret must be non-empty in production.
func NewManager(secret string, accessTTL, refreshTTL time.Duration, issuer string) *Manager {
	return &Manager{
		secret:     []byte(secret),
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
		issuer:     issuer,
		now:        time.Now,
	}
}

// AccessTTL exposes the configured access-token lifetime (used to populate the
// expires_in field returned to clients).
func (m *Manager) AccessTTL() time.Duration { return m.accessTTL }

// RefreshTTL exposes the configured refresh-token lifetime.
func (m *Manager) RefreshTTL() time.Duration { return m.refreshTTL }

// GenerateAccessToken mints a short-lived access token for the user.
func (m *Manager) GenerateAccessToken(u *entity.User) (string, error) {
	return m.generate(u, TokenTypeAccess, m.accessTTL)
}

// GenerateRefreshToken mints a long-lived refresh token for the user.
func (m *Manager) GenerateRefreshToken(u *entity.User) (string, error) {
	return m.generate(u, TokenTypeRefresh, m.refreshTTL)
}

func (m *Manager) generate(u *entity.User, tokenType string, ttl time.Duration) (string, error) {
	now := m.now()
	claims := Claims{
		RegisteredClaims: gojwt.RegisteredClaims{
			// A random JWT ID makes every token unique, so refresh-token
			// rotation always yields a distinct token even when two are minted
			// within the same second (HS256 over identical claims is otherwise
			// byte-identical).
			ID:        uuid.NewString(),
			Issuer:    m.issuer,
			Subject:   u.ID.String(),
			IssuedAt:  gojwt.NewNumericDate(now),
			NotBefore: gojwt.NewNumericDate(now),
			ExpiresAt: gojwt.NewNumericDate(now.Add(ttl)),
		},
		UserID:    u.ID,
		Role:      string(u.Role),
		TokenType: tokenType,
	}
	tok := gojwt.NewWithClaims(gojwt.SigningMethodHS256, claims)
	return tok.SignedString(m.secret)
}

// Verify validates a token's signature, expiry, issuer, and that its TokenType
// matches expectedType. Every failure mode collapses to apperror.Unauthorized
// so callers cannot distinguish "expired" from "tampered" from "wrong type"
// (spec §7.3 — opaque verification failures).
func (m *Manager) Verify(token, expectedType string) (*Claims, error) {
	claims := &Claims{}
	parsed, err := gojwt.ParseWithClaims(token, claims, func(t *gojwt.Token) (any, error) {
		if _, ok := t.Method.(*gojwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return m.secret, nil
	},
		gojwt.WithValidMethods([]string{"HS256"}),
		gojwt.WithIssuer(m.issuer),
		gojwt.WithTimeFunc(m.now),
	)
	if err != nil || !parsed.Valid {
		return nil, apperror.Unauthorized("invalid token", err)
	}
	if claims.TokenType != expectedType {
		return nil, apperror.Unauthorized("invalid token", nil)
	}
	return claims, nil
}

// HashToken returns the SHA-256 hex digest of a raw token. Refresh tokens are
// stored and looked up by this digest so the raw value never touches the
// database or cache (spec §7.3 — raw refresh tokens are never persisted).
func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

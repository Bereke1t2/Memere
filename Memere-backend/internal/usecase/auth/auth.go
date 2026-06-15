// Package auth holds the authentication usecases: registration, login, token
// refresh (with rotation), and logout. It is pure orchestration over the domain
// repository interfaces plus the jwt/password packages — it imports no
// infrastructure and no delivery code (dependency rule). HTTP handlers in
// Skill 5 adapt these methods to the wire.
package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	infmetrics "github.com/Bereke1t2/Memere/memere-backend/internal/infrastructure/metrics"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/jwt"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/password"
)

// minPasswordLen is the minimum accepted password length at registration.
const minPasswordLen = 8

// LockoutConfig holds the account-lockout policy (Phase 6 §7.3).
type LockoutConfig struct {
	// MaxFailures is the number of consecutive bad-password attempts that trigger
	// a lockout. Zero disables the feature.
	MaxFailures int
	// LockoutTTL is how long a locked account stays locked.
	LockoutTTL time.Duration
}

// AuthTokens is the credential pair returned by login and refresh. ExpiresIn is
// the access-token lifetime in seconds (the refresh token's own expiry is long
// and not surfaced here).
type AuthTokens struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int
}

// RegisterInput is the registration request after transport decoding.
type RegisterInput struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
	Phone     *string
	// Role is optional; empty defaults to student. Self-assigning admin is
	// rejected (only an existing admin may mint admins, in a later phase).
	Role entity.Role
}

// LoginInput is the login request after transport decoding.
type LoginInput struct {
	Email    string
	Password string
	// DeviceInfo is an optional free-form client descriptor stored alongside the
	// refresh token for auditing.
	DeviceInfo *string
}

// Service implements the auth usecases over the domain repository interfaces.
type Service struct {
	users    repository.UserRepository
	tokens   repository.RefreshTokenRepository
	sessions repository.SessionRepository
	jwt      *jwt.Manager
	lockout  LockoutConfig
	// now is injectable for tests; defaults to time.Now.
	now func() time.Time
}

// NewService wires the auth usecase with its dependencies.
func NewService(
	users repository.UserRepository,
	tokens repository.RefreshTokenRepository,
	sessions repository.SessionRepository,
	jwtMgr *jwt.Manager,
) *Service {
	return &Service{
		users:    users,
		tokens:   tokens,
		sessions: sessions,
		jwt:      jwtMgr,
		lockout:  LockoutConfig{MaxFailures: 10, LockoutTTL: 30 * time.Minute},
		now:      time.Now,
	}
}

// WithLockout overrides the default account-lockout policy. Call from main after
// NewService to wire config-driven settings.
func (s *Service) WithLockout(cfg LockoutConfig) *Service {
	s.lockout = cfg
	return s
}

// Register validates the input, ensures the email is free, hashes the password,
// and creates an unverified user with a stored email-verification token (the
// email itself is sent in a later phase). The returned user is sanitized.
func (s *Service) Register(ctx context.Context, in RegisterInput) (*entity.User, error) {
	email := normalizeEmail(in.Email)
	if details := validateRegister(in, email); len(details) > 0 {
		return nil, apperror.Validation(details, nil)
	}

	role := in.Role
	if role == "" {
		role = entity.RoleStudent
	}
	// Self-service registration may never create a privileged account.
	if role == entity.RoleAdmin {
		return nil, apperror.Forbidden("cannot self-assign admin role", nil)
	}
	if !role.Valid() {
		return nil, apperror.Validation(map[string]any{"role": "unknown role"}, nil)
	}

	// Pre-check keeps the common case returning a clean EMAIL_TAKEN; the unique
	// index in Create remains the authoritative guard against the race.
	switch _, err := s.users.FindByEmail(ctx, email); {
	case err == nil:
		return nil, apperror.Conflict("EMAIL_TAKEN", nil)
	case !apperror.IsNotFound(err):
		return nil, err
	}

	hash, err := password.Hash(in.Password)
	if err != nil {
		return nil, apperror.Internal(err)
	}
	verifyToken := randomToken()

	u := &entity.User{
		Email:                  email,
		Phone:                  in.Phone,
		PasswordHash:           hash,
		Role:                   role,
		FirstName:              strings.TrimSpace(in.FirstName),
		LastName:               strings.TrimSpace(in.LastName),
		IsActive:               true,
		IsEmailVerified:        false,
		EmailVerificationToken: &verifyToken,
	}
	if err := s.users.Create(ctx, u); err != nil {
		return nil, err
	}

	sanitized := u.Sanitized()
	return &sanitized, nil
}

// Login verifies the password and, on success, issues an access+refresh pair,
// persists the refresh-token hash in Postgres and Redis, and stamps
// last_login_at. A missing email and a wrong password return the *same* opaque
// INVALID_CREDENTIALS error so the endpoint does not reveal which addresses are
// registered (spec §7.3). After s.lockout.MaxFailures consecutive failures for
// an account the endpoint enters a lockout period and still returns opaque
// INVALID_CREDENTIALS (§7.3 — avoid revealing lockout status to attackers).
func (s *Service) Login(ctx context.Context, in LoginInput) (*AuthTokens, *entity.User, error) {
	email := normalizeEmail(in.Email)

	u, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		if apperror.IsNotFound(err) {
			return nil, nil, invalidCredentials()
		}
		return nil, nil, err
	}
	if !u.IsActive {
		return nil, nil, apperror.Forbidden("account is disabled", nil)
	}

	// Account lockout check (feature enabled when MaxFailures > 0).
	if s.lockout.MaxFailures > 0 {
		locked, err := s.sessions.IsLockedOut(ctx, u.ID, s.lockout.MaxFailures)
		if err != nil {
			return nil, nil, err
		}
		if locked {
			return nil, nil, invalidCredentials()
		}
	}

	if err := password.Compare(u.PasswordHash, in.Password); err != nil {
		// Record the failure and emit a metric on lockout.
		if s.lockout.MaxFailures > 0 {
			count, incrErr := s.sessions.IncrLoginFailure(ctx, u.ID, s.lockout.LockoutTTL)
			if incrErr == nil && count >= int64(s.lockout.MaxFailures) {
				infmetrics.ObserveSuspiciousLogin("lockout")
			}
		}
		return nil, nil, invalidCredentials()
	}

	// Success — clear the failure counter.
	if s.lockout.MaxFailures > 0 {
		_ = s.sessions.ClearLoginFailures(ctx, u.ID)
	}

	tokens, err := s.issueTokens(ctx, u, in.DeviceInfo)
	if err != nil {
		return nil, nil, err
	}

	if err := s.users.SetLastLogin(ctx, u.ID, s.now()); err != nil {
		return nil, nil, err
	}

	sanitized := u.Sanitized()
	return tokens, &sanitized, nil
}

// RevokeAccessToken adds the token's JTI to the denylist so it is rejected by
// RequireAuth before its natural expiry. Used by SuspendUser and logout flows.
// The TTL is set to the token's remaining lifetime so the key auto-expires.
func (s *Service) RevokeAccessToken(ctx context.Context, tokenStr string) error {
	claims, err := s.jwt.Verify(tokenStr, jwt.TokenTypeAccess)
	if err != nil {
		return nil // already expired or invalid — nothing to do
	}
	remaining := claims.ExpiresAt.Time.Sub(s.now())
	if remaining <= 0 {
		return nil
	}
	return s.sessions.DenyToken(ctx, claims.ID, remaining)
}

// Refresh validates a refresh token (signature, type, and presence in the
// session store / Postgres), then issues a new access token and *rotates* the
// refresh token: the presented token is revoked and a fresh one is minted. This
// bounds the damage of a leaked refresh token (spec §7.3, token rotation).
func (s *Service) Refresh(ctx context.Context, refreshToken string) (*AuthTokens, error) {
	claims, err := s.jwt.Verify(refreshToken, jwt.TokenTypeRefresh)
	if err != nil {
		return nil, err
	}
	hash := jwt.HashToken(refreshToken)

	// Source of truth: the stored token must exist and still be active.
	stored, err := s.tokens.FindByHash(ctx, hash)
	if err != nil {
		if apperror.IsNotFound(err) {
			return nil, apperror.Unauthorized("invalid token", err)
		}
		return nil, err
	}
	if !stored.IsActive(s.now()) || stored.UserID != claims.UserID {
		return nil, apperror.Unauthorized("invalid token", nil)
	}

	u, err := s.users.FindByID(ctx, claims.UserID)
	if err != nil {
		if apperror.IsNotFound(err) {
			return nil, apperror.Unauthorized("invalid token", err)
		}
		return nil, err
	}

	// Rotate: revoke the presented token before minting its replacement so a
	// replay of the old token fails even if issuing the new one races.
	if err := s.tokens.Revoke(ctx, stored.ID); err != nil {
		return nil, err
	}
	return s.issueTokens(ctx, u, stored.DeviceInfo)
}

// Logout revokes the presented refresh token in Postgres and clears the Redis
// session. It is idempotent: an unknown or already-revoked token still
// succeeds, so a client can always "log out".
func (s *Service) Logout(ctx context.Context, userID uuid.UUID, refreshToken string) error {
	if refreshToken != "" {
		hash := jwt.HashToken(refreshToken)
		stored, err := s.tokens.FindByHash(ctx, hash)
		switch {
		case err == nil:
			if err := s.tokens.Revoke(ctx, stored.ID); err != nil {
				return err
			}
		case !apperror.IsNotFound(err):
			return err
		}
	}
	return s.sessions.DeleteSession(ctx, userID)
}

// issueTokens mints an access+refresh pair, persists the refresh-token hash in
// Postgres (authoritative) and Redis (fast path), and returns the raw tokens.
func (s *Service) issueTokens(ctx context.Context, u *entity.User, deviceInfo *string) (*AuthTokens, error) {
	access, err := s.jwt.GenerateAccessToken(u)
	if err != nil {
		return nil, apperror.Internal(err)
	}
	refresh, err := s.jwt.GenerateRefreshToken(u)
	if err != nil {
		return nil, apperror.Internal(err)
	}

	refreshHash := jwt.HashToken(refresh)
	refreshTTL := s.jwt.RefreshTTL()
	rt := &entity.RefreshToken{
		UserID:     u.ID,
		TokenHash:  refreshHash,
		DeviceInfo: deviceInfo,
		ExpiresAt:  s.now().Add(refreshTTL),
	}
	if err := s.tokens.Create(ctx, rt); err != nil {
		return nil, err
	}
	if err := s.sessions.SetSession(ctx, u.ID, refreshHash, refreshTTL); err != nil {
		return nil, err
	}

	return &AuthTokens{
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresIn:    int(s.jwt.AccessTTL().Seconds()),
	}, nil
}

// invalidCredentials is the single opaque error used for both "no such email"
// and "wrong password" so the two are indistinguishable to a caller.
func invalidCredentials() error {
	return apperror.New(http.StatusUnauthorized, "INVALID_CREDENTIALS", "invalid email or password", nil)
}

// validateRegister returns a field→message map of validation problems; empty
// means valid.
func validateRegister(in RegisterInput, email string) map[string]any {
	details := map[string]any{}
	if _, err := mail.ParseAddress(email); err != nil || email == "" {
		details["email"] = "must be a valid email address"
	}
	if len(in.Password) < minPasswordLen {
		details["password"] = "must be at least 8 characters"
	}
	if strings.TrimSpace(in.FirstName) == "" {
		details["first_name"] = "is required"
	}
	if strings.TrimSpace(in.LastName) == "" {
		details["last_name"] = "is required"
	}
	return details
}

// normalizeEmail lowercases and trims an email so lookups and the unique index
// agree on casing.
func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// randomToken returns a 32-byte hex token suitable for email verification /
// password reset links.
func randomToken() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

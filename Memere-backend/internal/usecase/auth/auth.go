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
	"fmt"
	"log/slog"
	"math/big"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/service"
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
	Email      string
	Password   string
	// DeviceInfo is an optional free-form client descriptor stored alongside the
	// refresh token for auditing.
	DeviceInfo *string
	// DeviceID is a unique client device identifier to distinguish sessions.
	DeviceID   *string
	// Force when true supersedes any existing session on another device.
	Force      bool
}

// Service implements the auth usecases over the domain repository interfaces.
type Service struct {
	users     repository.UserRepository
	tokens    repository.RefreshTokenRepository
	sessions  repository.SessionRepository
	jwt       *jwt.Manager
	email     service.EmailSender
	publicURL string
	lockout   LockoutConfig
	// now is injectable for tests; defaults to time.Now.
	now func() time.Time
	// syncEmail sends verification emails synchronously when true (used in unit tests).
	syncEmail bool
}

// NewService wires the auth usecase with its dependencies.
func NewService(
	users repository.UserRepository,
	tokens repository.RefreshTokenRepository,
	sessions repository.SessionRepository,
	jwtMgr *jwt.Manager,
	emailSender ...service.EmailSender,
) *Service {
	var es service.EmailSender
	if len(emailSender) > 0 {
		es = emailSender[0]
	}
	return &Service{
		users:    users,
		tokens:   tokens,
		sessions: sessions,
		jwt:      jwtMgr,
		email:    es,
		lockout:  LockoutConfig{MaxFailures: 10, LockoutTTL: 30 * time.Minute},
		now:      time.Now,
	}
}

// WithEmailSender attaches an email sender to the service.
func (s *Service) WithEmailSender(es service.EmailSender) *Service {
	s.email = es
	return s
}

// WithSyncEmail sets email dispatch to synchronous mode (useful for unit tests).
func (s *Service) WithSyncEmail() *Service {
	s.syncEmail = true
	return s
}

// WithPublicURL sets the public base URL for verification links.
func (s *Service) WithPublicURL(publicURL string) *Service {
	s.publicURL = strings.TrimRight(publicURL, "/")
	return s
}

// WithLockout overrides the default account-lockout policy. Call from main after
// NewService to wire config-driven settings.
func (s *Service) WithLockout(cfg LockoutConfig) *Service {
	s.lockout = cfg
	return s
}

// Register validates the input, ensures the email is free, hashes the password,
// and creates an unverified user with a stored email-verification token.
// If an email sender is wired, it delivers a verification email. The returned user is sanitized.
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
	verifyToken := randomVerificationCode()

	approvalStatus := entity.ApprovalStatusPending
	if role != entity.RoleStudent {
		approvalStatus = entity.ApprovalStatusApproved
	}

	u := &entity.User{
		Email:                  email,
		Phone:                  in.Phone,
		PasswordHash:           hash,
		Role:                   role,
		ApprovalStatus:         approvalStatus,
		FirstName:              strings.TrimSpace(in.FirstName),
		LastName:               strings.TrimSpace(in.LastName),
		IsActive:               true,
		IsEmailVerified:        false,
		EmailVerificationToken: &verifyToken,
	}
	if err := s.users.Create(ctx, u); err != nil {
		return nil, err
	}

	s.sendVerificationEmail(ctx, u, verifyToken)

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

	// Check for active session on another device (Single-Device Policy).
	// Only students are restricted to a single device session. Admin and teacher
	// accounts can log in concurrently from multiple/different devices.
	if u.Role == entity.RoleStudent {
		activeSession, err := s.sessions.GetSessionData(ctx, u.ID)
		if err != nil {
			return nil, nil, err
		}
		if activeSession != nil && activeSession.TokenHash != "" {
			reqDeviceID := ""
			if in.DeviceID != nil {
				reqDeviceID = strings.TrimSpace(*in.DeviceID)
			}
			isSameDevice := reqDeviceID != "" && activeSession.DeviceID != "" && reqDeviceID == activeSession.DeviceID
			if !isSameDevice && !in.Force {
				return nil, nil, apperror.New(
					http.StatusConflict,
					"ACTIVE_SESSION_EXISTS",
					"This account is currently active on another device. Please log out from that device first before logging in here.",
					nil,
				)
			}
			// Same device or forced takeover: revoke old refresh token so old session is cleaned up.
			oldToken, err := s.tokens.FindByHash(ctx, activeSession.TokenHash)
			if err == nil && oldToken != nil {
				_ = s.tokens.Revoke(ctx, oldToken.ID)
			}
		}
	}

	tokens, err := s.issueTokens(ctx, u, in.DeviceInfo, in.DeviceID)
	if err != nil {
		return nil, nil, err
	}

	if err := s.users.SetLastLogin(ctx, u.ID, s.now()); err != nil {
		return nil, nil, err
	}

	sanitized := u.Sanitized()
	return tokens, &sanitized, nil
}

// VerifyEmail validates the token and marks the user's email verified.
// If an email is supplied, the verification is scoped to that email for collision safety.
func (s *Service) VerifyEmail(ctx context.Context, token string, email ...string) (*entity.User, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return nil, apperror.BadRequest("verification token is required", nil)
	}

	var u *entity.User
	var err error
	if len(email) > 0 && strings.TrimSpace(email[0]) != "" {
		normalized := normalizeEmail(email[0])
		u, err = s.users.FindByEmail(ctx, normalized)
		if err != nil {
			if apperror.IsNotFound(err) {
				return nil, apperror.New(http.StatusBadRequest, "INVALID_VERIFICATION_TOKEN", "invalid or expired email verification token", nil)
			}
			return nil, err
		}
		if u.IsEmailVerified {
			sanitized := u.Sanitized()
			return &sanitized, nil
		}
		if u.EmailVerificationToken == nil || *u.EmailVerificationToken != token {
			return nil, apperror.New(http.StatusBadRequest, "INVALID_VERIFICATION_TOKEN", "invalid or expired email verification token", nil)
		}
	} else {
		u, err = s.users.FindByEmailVerificationToken(ctx, token)
		if err != nil {
			if apperror.IsNotFound(err) {
				return nil, apperror.New(http.StatusBadRequest, "INVALID_VERIFICATION_TOKEN", "invalid or expired email verification token", nil)
			}
			return nil, err
		}
		if u.IsEmailVerified {
			sanitized := u.Sanitized()
			return &sanitized, nil
		}
	}

	u.IsEmailVerified = true
	u.EmailVerificationToken = nil
	if err := s.users.Update(ctx, u); err != nil {
		return nil, err
	}

	s.sendWelcomeEmail(ctx, u)

	sanitized := u.Sanitized()
	return &sanitized, nil
}

// ResendVerificationEmail issues a fresh verification token and emails the user.
func (s *Service) ResendVerificationEmail(ctx context.Context, email string) error {
	email = normalizeEmail(email)
	if email == "" {
		return apperror.BadRequest("email is required", nil)
	}

	u, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		if apperror.IsNotFound(err) {
			// Do not leak email existence
			return nil
		}
		return err
	}

	if u.IsEmailVerified {
		return nil
	}

	token := randomVerificationCode()
	u.EmailVerificationToken = &token
	if err := s.users.Update(ctx, u); err != nil {
		return err
	}

	s.sendVerificationEmail(ctx, u, token)
	return nil
}

func (s *Service) sendVerificationEmail(_ context.Context, u *entity.User, token string) {
	if s.email == nil {
		return
	}
	subject := "Verify your Mirkuz account email"
	body := fmt.Sprintf(
		`<h2>Welcome to Mirkuz, %s!</h2><p>Please use this verification code to verify your email address:</p><p style="font-size: 32px; font-weight: bold; letter-spacing: 6px; color: #10b981;">%s</p><p>If you did not register this account, please disregard this email.</p>`,
		u.FirstName, token,
	)
	if s.syncEmail {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := s.email.Send(ctx, u.Email, subject, body); err != nil {
			slog.Error("failed to send email verification", "err", err, "email", u.Email)
		}
		return
	}
	// Fire-and-forget with a detached context so the email send is not
	// cancelled when the HTTP client disconnects (e.g. mobile timeout).
	go func() {
		emailCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := s.email.Send(emailCtx, u.Email, subject, body); err != nil {
			slog.Error("failed to send email verification", "err", err, "email", u.Email)
		}
	}()
}

func (s *Service) sendWelcomeEmail(_ context.Context, u *entity.User) {
	if s.email == nil {
		return
	}
	subject := "Welcome to Mirkuz — Your Account is Verified! 🎓"
	body := s.buildWelcomeEmailHTML(u)
	if s.syncEmail {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := s.email.Send(ctx, u.Email, subject, body); err != nil {
			slog.Error("failed to send welcome email", "err", err, "email", u.Email)
		}
		return
	}
	// Fire-and-forget with a detached context so the email send is not
	// cancelled when the HTTP client disconnects (e.g. mobile timeout).
	go func() {
		emailCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := s.email.Send(emailCtx, u.Email, subject, body); err != nil {
			slog.Error("failed to send welcome email", "err", err, "email", u.Email)
		}
	}()
}

func (s *Service) buildWelcomeEmailHTML(u *entity.User) string {
	name := u.FirstName
	if name == "" {
		name = "there"
	}

	appURL := s.publicURL
	if appURL == "" {
		appURL = "https://mirkuz.app"
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Welcome to Mirkuz</title>
</head>
<body style="margin: 0; padding: 0; background-color: #f1f5f9; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif; -webkit-font-smoothing: antialiased;">
  <table role="presentation" border="0" cellpadding="0" cellspacing="0" width="100%%" style="background-color: #f1f5f9; padding: 32px 16px;">
    <tr>
      <td align="center">
        <!-- Main Card Container -->
        <table role="presentation" border="0" cellpadding="0" cellspacing="0" width="100%%" style="max-width: 600px; background-color: #ffffff; border-radius: 16px; overflow: hidden; box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05), 0 2px 4px -2px rgba(0, 0, 0, 0.05); border: 1px solid #e2e8f0;">
          
          <!-- Brand Header -->
          <tr>
            <td style="background-color: #0f172a; padding: 32px 40px; text-align: center;">
              <h1 style="margin: 0; font-size: 26px; font-weight: 800; letter-spacing: 1.5px; color: #ffffff;">
                🎓 MIRKUZ
              </h1>
              <p style="margin: 6px 0 0 0; font-size: 13px; color: #94a3b8; font-weight: 400; letter-spacing: 0.5px;">
                Your Academic Success Partner
              </p>
            </td>
          </tr>

          <!-- Hero Section -->
          <tr>
            <td style="padding: 40px 40px 24px 40px; text-align: center;">
              <!-- Verification Success Badge -->
              <div style="display: inline-block; background-color: #dcfce7; border: 1px solid #86efac; border-radius: 9999px; padding: 6px 16px; margin-bottom: 20px;">
                <span style="color: #166534; font-size: 13px; font-weight: 700; letter-spacing: 0.3px;">
                  ✓ Email Verified Successfully
                </span>
              </div>
              
              <h2 style="margin: 0 0 12px 0; color: #0f172a; font-size: 24px; font-weight: 700; line-height: 1.3;">
                Welcome aboard, %s! 🚀
              </h2>
              <p style="margin: 0; color: #475569; font-size: 15px; line-height: 1.6;">
                Your account is officially active. You have joined a dedicated community of students preparing for academic excellence and exam mastery.
              </p>
            </td>
          </tr>

          <!-- Value / Features Section -->
          <tr>
            <td style="padding: 0 40px 24px 40px;">
              <table role="presentation" border="0" cellpadding="0" cellspacing="0" width="100%%" style="border-top: 1px solid #f1f5f9; padding-top: 24px;">
                <tr>
                  <td style="padding-bottom: 16px;">
                    <h3 style="margin: 0; font-size: 16px; font-weight: 700; color: #0f172a;">
                      Here is what you can do right now:
                    </h3>
                  </td>
                </tr>

                <!-- Feature 1: Video Courses -->
                <tr>
                  <td style="padding-bottom: 16px;">
                    <table role="presentation" border="0" cellpadding="0" cellspacing="0" width="100%%" style="background-color: #f8fafc; border: 1px solid #e2e8f0; border-radius: 12px; padding: 16px;">
                      <tr>
                        <td width="42" valign="top" style="font-size: 26px; line-height: 1;">📚</td>
                        <td style="padding-left: 12px;">
                          <div style="font-size: 15px; font-weight: 700; color: #1e293b; margin-bottom: 4px;">Curated Video Courses</div>
                          <div style="font-size: 13px; color: #64748b; line-height: 1.5;">Master difficult subjects chapter-by-chapter with high-yield video lessons taught by top educators.</div>
                        </td>
                      </tr>
                    </table>
                  </td>
                </tr>

                <!-- Feature 2: Timed Mock Exams -->
                <tr>
                  <td style="padding-bottom: 16px;">
                    <table role="presentation" border="0" cellpadding="0" cellspacing="0" width="100%%" style="background-color: #f8fafc; border: 1px solid #e2e8f0; border-radius: 12px; padding: 16px;">
                      <tr>
                        <td width="42" valign="top" style="font-size: 26px; line-height: 1;">📝</td>
                        <td style="padding-left: 12px;">
                          <div style="font-size: 15px; font-weight: 700; color: #1e293b; margin-bottom: 4px;">Realistic Mock Exams</div>
                          <div style="font-size: 13px; color: #64748b; line-height: 1.5;">Practice under realistic exam conditions with timed question sets and instant step-by-step explanations.</div>
                        </td>
                      </tr>
                    </table>
                  </td>
                </tr>

                <!-- Feature 3: Offline Study Materials -->
                <tr>
                  <td style="padding-bottom: 16px;">
                    <table role="presentation" border="0" cellpadding="0" cellspacing="0" width="100%%" style="background-color: #f8fafc; border: 1px solid #e2e8f0; border-radius: 12px; padding: 16px;">
                      <tr>
                        <td width="42" valign="top" style="font-size: 26px; line-height: 1;">📖</td>
                        <td style="padding-left: 12px;">
                          <div style="font-size: 15px; font-weight: 700; color: #1e293b; margin-bottom: 4px;">Summaries & Offline Notes</div>
                          <div style="font-size: 13px; color: #64748b; line-height: 1.5;">Download essential formulas, key concepts, and chapter notes to study anywhere, anytime.</div>
                        </td>
                      </tr>
                    </table>
                  </td>
                </tr>

                <!-- Feature 4: Analytics & Streaks -->
                <tr>
                  <td style="padding-bottom: 8px;">
                    <table role="presentation" border="0" cellpadding="0" cellspacing="0" width="100%%" style="background-color: #f8fafc; border: 1px solid #e2e8f0; border-radius: 12px; padding: 16px;">
                      <tr>
                        <td width="42" valign="top" style="font-size: 26px; line-height: 1;">📊</td>
                        <td style="padding-left: 12px;">
                          <div style="font-size: 15px; font-weight: 700; color: #1e293b; margin-bottom: 4px;">Track Progress & Daily Streaks</div>
                          <div style="font-size: 13px; color: #64748b; line-height: 1.5;">See your strengths, target your weak areas, and build a lasting daily study habit that pays off.</div>
                        </td>
                      </tr>
                    </table>
                  </td>
                </tr>

              </table>
            </td>
          </tr>

          <!-- CTA Button -->
          <tr>
            <td style="padding: 12px 40px 32px 40px; text-align: center;">
              <a href="%s" target="_blank" style="display: inline-block; background-color: #10b981; color: #ffffff; font-size: 16px; font-weight: 700; text-decoration: none; padding: 14px 36px; border-radius: 10px; box-shadow: 0 4px 6px -1px rgba(16, 185, 129, 0.3);">
                Start Learning Now &rarr;
              </a>
            </td>
          </tr>

          <!-- Pro Tip Callout -->
          <tr>
            <td style="padding: 0 40px 36px 40px;">
              <table role="presentation" border="0" cellpadding="0" cellspacing="0" width="100%%" style="background-color: #eff6ff; border-left: 4px solid #3b82f6; border-radius: 4px 8px 8px 4px; padding: 16px 20px;">
                <tr>
                  <td>
                    <div style="font-size: 13px; color: #1e40af; line-height: 1.5;">
                      <strong>💡 Quick Tip:</strong> Studying just 25 minutes a day with active recall and practice tests improves exam retention by over 80%%. Open the app today and complete your first quiz!
                    </div>
                  </td>
                </tr>
              </table>
            </td>
          </tr>

          <!-- Footer -->
          <tr>
            <td style="background-color: #f8fafc; border-top: 1px solid #e2e8f0; padding: 28px 40px; text-align: center;">
              <p style="margin: 0 0 8px 0; font-size: 13px; color: #64748b;">
                Need help or have questions? Contact us anytime at <a href="mailto:support@mirkuz.app" style="color: #10b981; text-decoration: none; font-weight: 600;">support@mirkuz.app</a>
              </p>
              <p style="margin: 0; font-size: 12px; color: #94a3b8;">
                &copy; %d Mirkuz Inc. All rights reserved.
              </p>
            </td>
          </tr>

        </table>
      </td>
    </tr>
  </table>
</body>
</html>`,
		name, appURL, time.Now().Year(),
	)
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

	var devID *string
	if u.Role == entity.RoleStudent {
		sessionData, err := s.sessions.GetSessionData(ctx, u.ID)
		if err == nil && sessionData != nil && sessionData.DeviceID != "" {
			devID = &sessionData.DeviceID
		}
	}

	// Rotate: revoke the presented token before minting its replacement so a
	// replay of the old token fails even if issuing the new one races.
	if err := s.tokens.Revoke(ctx, stored.ID); err != nil {
		return nil, err
	}
	return s.issueTokens(ctx, u, stored.DeviceInfo, devID)
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

// DeleteAccount soft-deletes the user record and revokes all sessions and refresh tokens.
func (s *Service) DeleteAccount(ctx context.Context, userID uuid.UUID) error {
	if err := s.users.SoftDelete(ctx, userID); err != nil {
		return err
	}
	_ = s.tokens.RevokeAllForUser(ctx, userID)
	_ = s.sessions.DeleteSession(ctx, userID)
	return nil
}

// issueTokens mints an access+refresh pair, persists the refresh-token hash in
// Postgres (authoritative) and Redis (fast path), and returns the raw tokens.
func (s *Service) issueTokens(ctx context.Context, u *entity.User, deviceInfo *string, deviceID *string) (*AuthTokens, error) {
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
	if u.Role == entity.RoleStudent {
		devID := ""
		if deviceID != nil {
			devID = *deviceID
		}
		data := repository.SessionData{
			DeviceID:  devID,
			TokenHash: refreshHash,
		}
		if err := s.sessions.SetSessionData(ctx, u.ID, data, refreshTTL); err != nil {
			return nil, err
		}
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

// randomVerificationCode returns a cryptographically secure 6-digit numeric OTP ("000000"-"999999").
func randomVerificationCode() string {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		b := make([]byte, 4)
		_, _ = rand.Read(b)
		val := (uint32(b[0])<<24 | uint32(b[1])<<16 | uint32(b[2])<<8 | uint32(b[3])) % 1000000
		return fmt.Sprintf("%06d", val)
	}
	return fmt.Sprintf("%06d", n.Int64())
}

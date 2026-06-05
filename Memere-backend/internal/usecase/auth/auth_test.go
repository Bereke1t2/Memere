package auth

import (
	"context"
	"testing"
	"time"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/jwt"
)

func newTestService() (*Service, *fakeUserRepo, *fakeTokenRepo, *fakeSessionRepo) {
	users := newFakeUserRepo()
	tokens := newFakeTokenRepo()
	sessions := newFakeSessionRepo()
	mgr := jwt.NewManager("test-secret", 15*time.Minute, 720*time.Hour, "memere-test")
	return NewService(users, tokens, sessions, mgr), users, tokens, sessions
}

func validRegisterInput() RegisterInput {
	return RegisterInput{
		Email:     "Student@Example.com",
		Password:  "correct horse battery",
		FirstName: "Abel",
		LastName:  "Tesfaye",
	}
}

func TestRegister_HappyPath(t *testing.T) {
	svc, users, _, _ := newTestService()

	u, err := svc.Register(context.Background(), validRegisterInput())
	if err != nil {
		t.Fatalf("Register: unexpected error: %v", err)
	}
	if u.Email != "student@example.com" {
		t.Errorf("email not normalized: got %q", u.Email)
	}
	if u.Role != entity.RoleStudent {
		t.Errorf("expected default student role, got %q", u.Role)
	}
	if u.PasswordHash != "" {
		t.Error("returned user must be sanitized (password hash present)")
	}
	if u.EmailVerificationToken != nil {
		t.Error("returned user must be sanitized (verification token present)")
	}
	// The stored user keeps its hash and verification token.
	stored, _ := users.FindByEmail(context.Background(), "student@example.com")
	if stored.PasswordHash == "" {
		t.Error("stored user should retain password hash")
	}
	if stored.EmailVerificationToken == nil {
		t.Error("stored user should retain verification token")
	}
	if stored.IsEmailVerified {
		t.Error("new user should be unverified")
	}
}

func TestRegister_DuplicateEmail(t *testing.T) {
	svc, _, _, _ := newTestService()
	ctx := context.Background()
	if _, err := svc.Register(ctx, validRegisterInput()); err != nil {
		t.Fatalf("first register failed: %v", err)
	}
	_, err := svc.Register(ctx, validRegisterInput())
	if !apperror.IsCode(err, "CONFLICT") {
		t.Fatalf("expected CONFLICT (EMAIL_TAKEN), got %v", err)
	}
}

func TestRegister_RejectsSelfAssignedAdmin(t *testing.T) {
	svc, _, _, _ := newTestService()
	in := validRegisterInput()
	in.Role = entity.RoleAdmin
	_, err := svc.Register(context.Background(), in)
	if !apperror.IsCode(err, "FORBIDDEN") {
		t.Fatalf("expected FORBIDDEN, got %v", err)
	}
}

func TestRegister_ValidationErrors(t *testing.T) {
	svc, _, _, _ := newTestService()
	in := RegisterInput{Email: "not-an-email", Password: "short"}
	_, err := svc.Register(context.Background(), in)
	if !apperror.IsCode(err, "VALIDATION_ERROR") {
		t.Fatalf("expected VALIDATION_ERROR, got %v", err)
	}
}

func TestLogin_Success(t *testing.T) {
	svc, _, tokens, sessions := newTestService()
	ctx := context.Background()
	if _, err := svc.Register(ctx, validRegisterInput()); err != nil {
		t.Fatalf("register: %v", err)
	}

	tok, u, err := svc.Login(ctx, LoginInput{Email: "student@example.com", Password: "correct horse battery"})
	if err != nil {
		t.Fatalf("Login: unexpected error: %v", err)
	}
	if tok.AccessToken == "" || tok.RefreshToken == "" {
		t.Fatal("login should return both tokens")
	}
	if tok.ExpiresIn != int((15 * time.Minute).Seconds()) {
		t.Errorf("ExpiresIn = %d, want %d", tok.ExpiresIn, int((15 * time.Minute).Seconds()))
	}
	if u.PasswordHash != "" {
		t.Error("login must return a sanitized user")
	}
	// The refresh-token hash, not the raw token, is what gets stored.
	if _, err := tokens.FindByHash(ctx, jwt.HashToken(tok.RefreshToken)); err != nil {
		t.Errorf("refresh hash not persisted in PG: %v", err)
	}
	if sess, _ := sessions.GetSession(ctx, u.ID); sess != jwt.HashToken(tok.RefreshToken) {
		t.Error("session store should hold the refresh-token hash")
	}
	// last_login_at stamped.
	if u.LastLoginAt == nil {
		// sanitized copy is taken before SetLastLogin mutates store; just check store
		stored, _ := svc.users.FindByID(ctx, u.ID)
		if stored.LastLoginAt == nil {
			t.Error("last_login_at not set")
		}
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	svc, _, _, _ := newTestService()
	ctx := context.Background()
	if _, err := svc.Register(ctx, validRegisterInput()); err != nil {
		t.Fatalf("register: %v", err)
	}
	_, _, err := svc.Login(ctx, LoginInput{Email: "student@example.com", Password: "wrong"})
	if !apperror.IsCode(err, "INVALID_CREDENTIALS") {
		t.Fatalf("expected INVALID_CREDENTIALS, got %v", err)
	}
}

func TestLogin_UnknownEmail_SameError(t *testing.T) {
	svc, _, _, _ := newTestService()
	_, _, err := svc.Login(context.Background(), LoginInput{Email: "nobody@example.com", Password: "whatever"})
	if !apperror.IsCode(err, "INVALID_CREDENTIALS") {
		t.Fatalf("unknown email must return INVALID_CREDENTIALS (no user enumeration), got %v", err)
	}
}

func TestRefresh_Success_RotatesToken(t *testing.T) {
	svc, _, tokens, _ := newTestService()
	ctx := context.Background()
	if _, err := svc.Register(ctx, validRegisterInput()); err != nil {
		t.Fatalf("register: %v", err)
	}
	first, _, err := svc.Login(ctx, LoginInput{Email: "student@example.com", Password: "correct horse battery"})
	if err != nil {
		t.Fatalf("login: %v", err)
	}

	next, err := svc.Refresh(ctx, first.RefreshToken)
	if err != nil {
		t.Fatalf("Refresh: unexpected error: %v", err)
	}
	if next.RefreshToken == first.RefreshToken {
		t.Error("refresh must rotate the refresh token")
	}
	if next.AccessToken == "" {
		t.Error("refresh must mint a new access token")
	}
	// Old token revoked; replay must fail.
	old, _ := tokens.FindByHash(ctx, jwt.HashToken(first.RefreshToken))
	if old.RevokedAt == nil {
		t.Error("old refresh token should be revoked after rotation")
	}
	if _, err := svc.Refresh(ctx, first.RefreshToken); !apperror.IsCode(err, "UNAUTHORIZED") {
		t.Fatalf("replay of rotated token should be UNAUTHORIZED, got %v", err)
	}
}

func TestRefresh_RevokedToken(t *testing.T) {
	svc, _, tokens, _ := newTestService()
	ctx := context.Background()
	if _, err := svc.Register(ctx, validRegisterInput()); err != nil {
		t.Fatalf("register: %v", err)
	}
	tok, _, err := svc.Login(ctx, LoginInput{Email: "student@example.com", Password: "correct horse battery"})
	if err != nil {
		t.Fatalf("login: %v", err)
	}
	stored, _ := tokens.FindByHash(ctx, jwt.HashToken(tok.RefreshToken))
	if err := tokens.Revoke(ctx, stored.ID); err != nil {
		t.Fatalf("revoke: %v", err)
	}
	if _, err := svc.Refresh(ctx, tok.RefreshToken); !apperror.IsCode(err, "UNAUTHORIZED") {
		t.Fatalf("revoked token must be UNAUTHORIZED, got %v", err)
	}
}

func TestRefresh_AccessTokenRejected(t *testing.T) {
	svc, _, _, _ := newTestService()
	ctx := context.Background()
	if _, err := svc.Register(ctx, validRegisterInput()); err != nil {
		t.Fatalf("register: %v", err)
	}
	tok, _, err := svc.Login(ctx, LoginInput{Email: "student@example.com", Password: "correct horse battery"})
	if err != nil {
		t.Fatalf("login: %v", err)
	}
	// Presenting an access token to refresh must be rejected on the type check.
	if _, err := svc.Refresh(ctx, tok.AccessToken); !apperror.IsCode(err, "UNAUTHORIZED") {
		t.Fatalf("access token at refresh must be UNAUTHORIZED, got %v", err)
	}
}

func TestLogout_RevokesAndClearsSession(t *testing.T) {
	svc, _, tokens, sessions := newTestService()
	ctx := context.Background()
	if _, err := svc.Register(ctx, validRegisterInput()); err != nil {
		t.Fatalf("register: %v", err)
	}
	tok, u, err := svc.Login(ctx, LoginInput{Email: "student@example.com", Password: "correct horse battery"})
	if err != nil {
		t.Fatalf("login: %v", err)
	}

	if err := svc.Logout(ctx, u.ID, tok.RefreshToken); err != nil {
		t.Fatalf("Logout: unexpected error: %v", err)
	}
	stored, _ := tokens.FindByHash(ctx, jwt.HashToken(tok.RefreshToken))
	if stored.RevokedAt == nil {
		t.Error("logout should revoke the refresh token in PG")
	}
	if sess, _ := sessions.GetSession(ctx, u.ID); sess != "" {
		t.Error("logout should clear the Redis session")
	}
	// A logged-out refresh token can no longer be refreshed.
	if _, err := svc.Refresh(ctx, tok.RefreshToken); !apperror.IsCode(err, "UNAUTHORIZED") {
		t.Fatalf("refresh after logout must be UNAUTHORIZED, got %v", err)
	}
}

func TestLogout_Idempotent(t *testing.T) {
	svc, _, _, _ := newTestService()
	ctx := context.Background()
	// Logging out with no/unknown token still succeeds.
	if err := svc.Logout(ctx, [16]byte{}, ""); err != nil {
		t.Fatalf("logout with empty token should succeed, got %v", err)
	}
}

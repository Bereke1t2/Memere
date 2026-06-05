package jwt

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

func testUser() *entity.User {
	return &entity.User{ID: uuid.New(), Role: entity.RoleStudent}
}

func TestAccessToken_RoundTrip(t *testing.T) {
	u := testUser()
	mgr := NewManager("secret", 15*time.Minute, 720*time.Hour, "memere")

	tok, err := mgr.GenerateAccessToken(u)
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	claims, err := mgr.Verify(tok, TokenTypeAccess)
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if claims.UserID != u.ID {
		t.Errorf("UserID = %v, want %v", claims.UserID, u.ID)
	}
	if claims.Role != string(entity.RoleStudent) {
		t.Errorf("Role = %q, want student", claims.Role)
	}
	if claims.TokenType != TokenTypeAccess {
		t.Errorf("TokenType = %q, want access", claims.TokenType)
	}
}

func TestVerify_WrongTypeRejected(t *testing.T) {
	u := testUser()
	mgr := NewManager("secret", 15*time.Minute, 720*time.Hour, "memere")
	access, _ := mgr.GenerateAccessToken(u)

	// Asking for a refresh token but presenting an access token must fail.
	if _, err := mgr.Verify(access, TokenTypeRefresh); !apperror.IsCode(err, "UNAUTHORIZED") {
		t.Fatalf("expected UNAUTHORIZED, got %v", err)
	}
}

func TestVerify_ExpiredRejected(t *testing.T) {
	u := testUser()
	mgr := NewManager("secret", time.Hour, time.Hour, "memere")
	// Issue a token "in the past" so it is already expired at the real now.
	mgr.now = func() time.Time { return time.Now().Add(-2 * time.Hour) }
	tok, _ := mgr.GenerateAccessToken(u)

	mgr.now = time.Now // verify at real time
	if _, err := mgr.Verify(tok, TokenTypeAccess); !apperror.IsCode(err, "UNAUTHORIZED") {
		t.Fatalf("expired token should be UNAUTHORIZED, got %v", err)
	}
}

func TestVerify_TamperedSignatureRejected(t *testing.T) {
	u := testUser()
	mgr := NewManager("secret", 15*time.Minute, 720*time.Hour, "memere")
	tok, _ := mgr.GenerateAccessToken(u)

	// Flip a character in the signature segment.
	parts := strings.Split(tok, ".")
	if len(parts) != 3 {
		t.Fatalf("unexpected token shape: %q", tok)
	}
	sig := []byte(parts[2])
	if sig[0] == 'A' {
		sig[0] = 'B'
	} else {
		sig[0] = 'A'
	}
	tampered := parts[0] + "." + parts[1] + "." + string(sig)
	if _, err := mgr.Verify(tampered, TokenTypeAccess); !apperror.IsCode(err, "UNAUTHORIZED") {
		t.Fatalf("tampered token should be UNAUTHORIZED, got %v", err)
	}
}

func TestVerify_WrongSecretRejected(t *testing.T) {
	u := testUser()
	signer := NewManager("secret-a", 15*time.Minute, 720*time.Hour, "memere")
	tok, _ := signer.GenerateAccessToken(u)

	verifier := NewManager("secret-b", 15*time.Minute, 720*time.Hour, "memere")
	if _, err := verifier.Verify(tok, TokenTypeAccess); !apperror.IsCode(err, "UNAUTHORIZED") {
		t.Fatalf("token signed with a different secret should be UNAUTHORIZED, got %v", err)
	}
}

func TestHashToken_DeterministicAndOpaque(t *testing.T) {
	raw := "a-refresh-token"
	h1 := HashToken(raw)
	h2 := HashToken(raw)
	if h1 != h2 {
		t.Error("hash must be deterministic")
	}
	if strings.Contains(h1, raw) {
		t.Error("hash must not contain the raw token")
	}
	if HashToken("other") == h1 {
		t.Error("distinct inputs must hash differently")
	}
}

package password

import (
	"strings"
	"testing"
)

func TestHash_DiffersFromPlaintext(t *testing.T) {
	plain := "correct horse battery staple"
	hash, err := Hash(plain)
	if err != nil {
		t.Fatalf("Hash: %v", err)
	}
	if hash == plain {
		t.Fatal("hash must not equal the plaintext")
	}
	if strings.Contains(hash, plain) {
		t.Fatal("hash must not contain the plaintext")
	}
}

func TestHash_Salted(t *testing.T) {
	// bcrypt salts each hash, so the same input hashes differently each time.
	h1, _ := Hash("samepass")
	h2, _ := Hash("samepass")
	if h1 == h2 {
		t.Error("two hashes of the same password should differ (salt)")
	}
}

func TestCompare_Success(t *testing.T) {
	plain := "s3cret-passw0rd"
	hash, err := Hash(plain)
	if err != nil {
		t.Fatalf("Hash: %v", err)
	}
	if err := Compare(hash, plain); err != nil {
		t.Errorf("Compare should succeed on the right password: %v", err)
	}
}

func TestCompare_Failure(t *testing.T) {
	hash, _ := Hash("the-right-one")
	if err := Compare(hash, "the-wrong-one"); err == nil {
		t.Error("Compare should fail on the wrong password")
	}
}

package storage

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CredentialStore persists the single admin Drive account's OAuth refresh token.
// The plaintext token never leaves the backend: it is written once by the admin
// connect flow and read only by the Drive store when minting access tokens.
type CredentialStore interface {
	// Save stores (or replaces) the refresh token for the connected account.
	Save(ctx context.Context, refreshToken, googleEmail string) error
	// Load returns the stored refresh token; ok is false when none is connected.
	Load(ctx context.Context) (refreshToken string, ok bool, err error)
}

// PgCredentialStore is a Postgres-backed CredentialStore. It encrypts the
// refresh token with AES-256-GCM under a key derived from the app's JWT secret,
// so the ciphertext in storage.google_credentials cannot be replayed against
// Google from a database dump alone.
type PgCredentialStore struct {
	pool *pgxpool.Pool
	aead cipher.AEAD
}

// NewPgCredentialStore builds a credential store. secret keys the at-rest
// encryption (the app JWT secret is used in wiring); it must be non-empty.
func NewPgCredentialStore(pool *pgxpool.Pool, secret string) (*PgCredentialStore, error) {
	aead, err := deriveAEAD(secret)
	if err != nil {
		return nil, err
	}
	return &PgCredentialStore{pool: pool, aead: aead}, nil
}

var _ CredentialStore = (*PgCredentialStore)(nil)

// deriveAEAD builds an AES-256-GCM AEAD from a 32-byte key derived from secret
// via a domain-separated SHA-256 (so it is distinct from the raw JWT secret).
func deriveAEAD(secret string) (cipher.AEAD, error) {
	if secret == "" {
		return nil, errors.New("storage: credential encryption needs a non-empty secret")
	}
	key := sha256.Sum256([]byte("memere-gdrive-cred-v1:" + secret))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, fmt.Errorf("storage: cred cipher: %w", err)
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("storage: cred gcm: %w", err)
	}
	return aead, nil
}

// seal returns nonce||ciphertext so open can recover the nonce.
func (s *PgCredentialStore) seal(plaintext string) ([]byte, error) {
	nonce := make([]byte, s.aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("storage: cred nonce: %w", err)
	}
	return s.aead.Seal(nonce, nonce, []byte(plaintext), nil), nil
}

func (s *PgCredentialStore) open(blob []byte) (string, error) {
	ns := s.aead.NonceSize()
	if len(blob) < ns {
		return "", errors.New("storage: cred ciphertext too short")
	}
	nonce, ct := blob[:ns], blob[ns:]
	pt, err := s.aead.Open(nil, nonce, ct, nil)
	if err != nil {
		return "", fmt.Errorf("storage: cred open: %w", err)
	}
	return string(pt), nil
}

func (s *PgCredentialStore) Save(ctx context.Context, refreshToken, googleEmail string) error {
	blob, err := s.seal(refreshToken)
	if err != nil {
		return err
	}
	const q = `
INSERT INTO storage.google_credentials (id, refresh_token, google_email, connected_at, updated_at)
VALUES (1, $1, $2, now(), now())
ON CONFLICT (id) DO UPDATE
   SET refresh_token = EXCLUDED.refresh_token,
       google_email  = EXCLUDED.google_email,
       updated_at    = now()`
	if _, err := s.pool.Exec(ctx, q, blob, googleEmail); err != nil {
		return fmt.Errorf("storage: save credentials: %w", err)
	}
	return nil
}

func (s *PgCredentialStore) Load(ctx context.Context) (string, bool, error) {
	const q = `SELECT refresh_token FROM storage.google_credentials WHERE id = 1`
	var blob []byte
	err := s.pool.QueryRow(ctx, q).Scan(&blob)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", false, nil
	}
	if err != nil {
		return "", false, fmt.Errorf("storage: load credentials: %w", err)
	}
	tok, err := s.open(blob)
	if err != nil {
		return "", false, err
	}
	return tok, true, nil
}

package storage

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/service"
)

// devStore embeds the ObjectStore interface (nil) and overrides only PresignGet,
// the one method the dev signer path delegates to. Any other call would panic,
// surfacing an unexpected dependency.
type devStore struct {
	service.ObjectStore
	lastKey string
	lastTTL time.Duration
}

func (f *devStore) PresignGet(_ context.Context, key string, ttl time.Duration) (string, error) {
	f.lastKey, f.lastTTL = key, ttl
	return "https://minio.test/" + key + "?ttl=" + strconv.Itoa(int(ttl.Seconds())), nil
}

func newPEMKey(t *testing.T) string {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	der := x509.MarshalPKCS1PrivateKey(key)
	block := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: der}
	return string(pem.EncodeToMemory(block))
}

func TestCDNSigner_DevFallsBackToPresignGet(t *testing.T) {
	fb := &devStore{}
	s, err := NewCDNSigner(CDNConfig{}, fb, nil)
	if err != nil {
		t.Fatalf("NewCDNSigner: %v", err)
	}
	got, err := s.SignGet(context.Background(), "hls/abc/master.m3u8", 90*time.Minute)
	if err != nil {
		t.Fatalf("SignGet: %v", err)
	}
	if fb.lastKey != "hls/abc/master.m3u8" {
		t.Fatalf("fallback signed key %q", fb.lastKey)
	}
	if fb.lastTTL != 90*time.Minute {
		t.Fatalf("fallback ttl %v, want 90m", fb.lastTTL)
	}
	if got == "" {
		t.Fatalf("empty URL")
	}
}

func TestCDNSigner_CloudFrontSignsWithExpiry(t *testing.T) {
	fixed := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	clock := func() time.Time { return fixed }

	s, err := NewCDNSigner(CDNConfig{
		Domain:        "d123.cloudfront.net",
		KeyPairID:     "APKAEXAMPLE",
		PrivateKeyPEM: newPEMKey(t),
	}, nil, clock)
	if err != nil {
		t.Fatalf("NewCDNSigner: %v", err)
	}

	ttl := 2 * time.Hour
	signed, err := s.SignGet(context.Background(), "hls/vid/master.m3u8", ttl)
	if err != nil {
		t.Fatalf("SignGet: %v", err)
	}
	u, err := url.Parse(signed)
	if err != nil {
		t.Fatalf("parse signed URL: %v", err)
	}
	if u.Host != "d123.cloudfront.net" {
		t.Fatalf("host = %q", u.Host)
	}
	q := u.Query()
	// Canned policy: Expires == now+ttl, plus Signature + Key-Pair-Id.
	wantExpires := strconv.FormatInt(fixed.Add(ttl).Unix(), 10)
	if q.Get("Expires") != wantExpires {
		t.Fatalf("Expires = %q, want %q", q.Get("Expires"), wantExpires)
	}
	if q.Get("Signature") == "" {
		t.Fatalf("missing Signature")
	}
	if q.Get("Key-Pair-Id") != "APKAEXAMPLE" {
		t.Fatalf("Key-Pair-Id = %q", q.Get("Key-Pair-Id"))
	}
}

func TestCDNSigner_KeyChangesSignature(t *testing.T) {
	clock := func() time.Time { return time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC) }
	s, err := NewCDNSigner(CDNConfig{
		Domain:        "d123.cloudfront.net",
		KeyPairID:     "APKAEXAMPLE",
		PrivateKeyPEM: newPEMKey(t),
	}, nil, clock)
	if err != nil {
		t.Fatalf("NewCDNSigner: %v", err)
	}
	a, err := s.SignGet(context.Background(), "hls/one/master.m3u8", time.Hour)
	if err != nil {
		t.Fatalf("sign a: %v", err)
	}
	b, err := s.SignGet(context.Background(), "hls/two/master.m3u8", time.Hour)
	if err != nil {
		t.Fatalf("sign b: %v", err)
	}
	sigA := mustQuery(t, a, "Signature")
	sigB := mustQuery(t, b, "Signature")
	if sigA == sigB {
		t.Fatalf("different keys produced identical signatures")
	}
}

func TestNewCDNSigner_ConfiguredButNoKeyPairIsError(t *testing.T) {
	_, err := NewCDNSigner(CDNConfig{Domain: "d123.cloudfront.net"}, nil, nil)
	if err == nil {
		t.Fatalf("expected error when domain set without key-pair id")
	}
}

func TestNewCDNSigner_DevWithoutFallbackIsError(t *testing.T) {
	_, err := NewCDNSigner(CDNConfig{}, nil, nil)
	if err == nil {
		t.Fatalf("expected error: dev signer needs a fallback store")
	}
}

func mustQuery(t *testing.T, raw, key string) string {
	t.Helper()
	u, err := url.Parse(raw)
	if err != nil {
		t.Fatalf("parse %q: %v", raw, err)
	}
	return u.Query().Get(key)
}

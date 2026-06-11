package storage

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"
	"time"

	cfsign "github.com/aws/aws-sdk-go-v2/feature/cloudfront/sign"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/service"
)

// URLSigner abstracts "give me a short-lived GET URL for this object key". The
// video delivery usecase depends on this port; it never imports an SDK or knows
// whether the URL is a CloudFront-signed URL or an S3/MinIO pre-signed GET.
type URLSigner interface {
	SignGet(ctx context.Context, key string, ttl time.Duration) (string, error)
}

// CDNSigner issues time-limited GET URLs for stored objects. When a CloudFront
// distribution is configured (domain + key-pair id + private key) it returns a
// CloudFront *signed URL* with a canned policy (Expires = now+ttl) so neither
// the manifest nor — via a wildcard policy at the distribution — the segments are
// publicly reachable. With no distribution configured it falls back to an
// S3/MinIO pre-signed GET, the development path.
//
// HLS segment-access strategy (spec §7.3 — hotlinking mitigation):
//
//	An HLS master/variant playlist references *segment* (.ts) URLs by relative
//	path; signing only the master does NOT protect the segments. The chosen
//	strategy is environment-split:
//
//	  Production (CloudFront): the bucket is private (origin-access only) and the
//	  distribution is configured to require signed access for the hls/* path.
//	  Issue a CloudFront *signed cookie* (or a wildcard signed URL) scoped to
//	  hls/{video_id}/* so the player's subsequent segment GETs carry the same
//	  short-lived credential. The signer here mints the master URL; the segment
//	  authorization rides on the distribution's signed-cookie/wildcard policy. No
//	  object is ever publicly readable.
//
//	  Development (MinIO presign): CloudFront cookies aren't available, so dev
//	  relies on a bucket policy that grants read on the hls/ prefix to the dev
//	  credentials only (not anonymous public). This is a deliberate dev-only
//	  relaxation; it MUST NOT ship to production, where the CloudFront signed
//	  path above is mandatory.
//
// The clock is injected so the canned-policy Expires is deterministic in tests;
// production passes time.Now.
type CDNSigner struct {
	domain    string // CloudFront domain, e.g. dxxxx.cloudfront.net (no scheme)
	keyPairID string
	privKey   *rsa.PrivateKey
	fallback  service.ObjectStore // S3/MinIO presign for dev
	now       func() time.Time
}

var _ URLSigner = (*CDNSigner)(nil)
var _ service.URLSigner = (*CDNSigner)(nil)

// CDNConfig carries the signing settings mapped from config.StorageConfig.
type CDNConfig struct {
	Domain        string // CDN_DOMAIN; empty disables CloudFront signing (dev)
	KeyPairID     string // CDN_KEY_PAIR_ID
	PrivateKeyPEM string // CDN_PRIVATE_KEY_PEM (PEM contents)
}

// NewCDNSigner builds a signer. When cfg.Domain (and the key material) is empty
// it returns a dev signer that delegates entirely to fallback. A configured
// domain with malformed or missing key material is a hard error so production
// never silently degrades to unsigned URLs.
func NewCDNSigner(cfg CDNConfig, fallback service.ObjectStore, now func() time.Time) (*CDNSigner, error) {
	if now == nil {
		now = time.Now
	}
	s := &CDNSigner{
		domain:    strings.TrimSpace(cfg.Domain),
		keyPairID: strings.TrimSpace(cfg.KeyPairID),
		fallback:  fallback,
		now:       now,
	}
	if s.domain == "" {
		// Development: no CloudFront, presign against the object store directly.
		if fallback == nil {
			return nil, errors.New("storage: CDN signer needs a fallback ObjectStore when no CDN is configured")
		}
		return s, nil
	}
	if s.keyPairID == "" {
		return nil, errors.New("storage: CDN_DOMAIN set but CDN_KEY_PAIR_ID is empty")
	}
	key, err := parseRSAPrivateKey(cfg.PrivateKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("storage: CDN private key: %w", err)
	}
	s.privKey = key
	return s, nil
}

// SignGet returns a GET URL valid for ttl. CloudFront-signed when configured,
// otherwise an object-store pre-signed GET.
func (c *CDNSigner) SignGet(ctx context.Context, key string, ttl time.Duration) (string, error) {
	if c.domain == "" || c.privKey == nil {
		return c.fallback.PresignGet(ctx, key, ttl) // dev path
	}
	raw := "https://" + c.domain + "/" + ltrimSlash(key)
	signer := cfsign.NewURLSigner(c.keyPairID, c.privKey)
	signed, err := signer.Sign(raw, c.now().Add(ttl))
	if err != nil {
		return "", fmt.Errorf("storage: cloudfront sign %q: %w", key, err)
	}
	return signed, nil
}

// ltrimSlash drops a single leading slash so the key joins cleanly after the
// domain without producing a double slash.
func ltrimSlash(key string) string {
	return strings.TrimPrefix(key, "/")
}

// parseRSAPrivateKey accepts a PEM-encoded RSA private key in PKCS#1 or PKCS#8.
func parseRSAPrivateKey(pemStr string) (*rsa.PrivateKey, error) {
	pemStr = strings.TrimSpace(pemStr)
	if pemStr == "" {
		return nil, errors.New("empty PEM")
	}
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, errors.New("not valid PEM")
	}
	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return key, nil
	}
	parsed, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse key: %w", err)
	}
	rsaKey, ok := parsed.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("PEM is not an RSA private key")
	}
	return rsaKey, nil
}

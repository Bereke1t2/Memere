package service

import (
	"context"
	"time"
)

// URLSigner issues short-lived, access-controlled GET URLs for stored objects.
// It is the delivery-side counterpart to ObjectStore: the video usecase depends
// on this port to mint streaming/download URLs, so it never imports a CDN or S3
// SDK. The infrastructure ring provides the implementation (CloudFront signed
// URLs in production, S3/MinIO pre-signed GET in dev).
type URLSigner interface {
	// SignGet returns a URL that grants GET access to key for at most ttl.
	SignGet(ctx context.Context, key string, ttl time.Duration) (string, error)
}

// DownloadTokens is the single-use download-token gate (spec §8.3 step 2). A
// signed URL is replayable within its window; this port lets the usecase make
// the *download* flow one-shot by issuing a token bound to (videoID, userID) and
// atomically consuming it on first redeem. Backed by Redis in the infra ring.
type DownloadTokens interface {
	// Issue binds token -> (videoID, userID) for ttl.
	Issue(ctx context.Context, token, videoID, userID string, ttl time.Duration) error
	// Consume atomically fetches and deletes the token. ok is false (no error)
	// for an unknown, expired, or already-consumed token.
	Consume(ctx context.Context, token string) (videoID, userID string, ok bool, err error)
}

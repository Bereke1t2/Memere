// Package service holds domain-level ports: interfaces the usecases depend on
// whose implementations live in the infrastructure ring. ObjectStore is the
// storage port — S3, MinIO or GCS satisfy it, so business logic never imports a
// storage SDK (clean-architecture dependency rule).
package service

import (
	"context"
	"io"
	"time"
)

// ObjectStore abstracts the object/blob store behind the video pipeline.
type ObjectStore interface {
	// PresignPut returns a URL the client uses to upload one object directly
	// (HTTP PUT), pinned to contentType, valid for ttl. The server never proxies
	// the bytes (spec §8.1).
	PresignPut(ctx context.Context, key, contentType string, ttl time.Duration) (string, error)
	// PresignGet returns a time-limited URL to download/stream a single object.
	PresignGet(ctx context.Context, key string, ttl time.Duration) (string, error)
	// Put uploads bytes server-side (used by the transcode worker for outputs).
	Put(ctx context.Context, key, contentType string, body io.Reader) error
	// Get streams an object (used by the worker to fetch the source).
	Get(ctx context.Context, key string) (io.ReadCloser, error)
	// Exists reports whether an object is present.
	Exists(ctx context.Context, key string) (bool, error)
	// Delete removes an object. Deleting a missing key is not an error.
	Delete(ctx context.Context, key string) error
}

// SizedUploader is an OPTIONAL ObjectStore capability: some stores (e.g. Google
// Drive) can stream an upload of known length straight to the backend without
// buffering it in memory. Callers type-assert an ObjectStore to this interface
// and fall back to Put when it is not implemented. Used by the proxied video
// upload so large files never have to be read fully into memory.
type SizedUploader interface {
	PutSized(ctx context.Context, key, contentType string, r io.Reader, size int64) error
}

// PrefixDeleter is an OPTIONAL ObjectStore capability: permanently delete every
// object under a key prefix (e.g. "hls/<videoId>/"). This is how a replaced or
// deleted video's HLS segments are purged — the individual .ts keys are never
// recorded in the database, so they are only reachable by prefix. Callers
// type-assert an ObjectStore to this interface and skip prefix cleanup when it
// is unimplemented (single-file backends like Google Drive store no segments, so
// deleting their known keys already reclaims everything). Deleting a prefix with
// no objects under it is not an error.
type PrefixDeleter interface {
	DeletePrefix(ctx context.Context, prefix string) error
}

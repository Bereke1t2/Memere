// Package storage holds object-store implementations of the domain ObjectStore
// port. S3Store talks to AWS S3 in production, MinIO in local dev (custom
// endpoint + path-style addressing), and Backblaze B2 (its S3-compatible API:
// custom endpoint + application-key credentials, checksums disabled). It is the
// only place the AWS SDK is imported; usecases depend on service.ObjectStore,
// never on this package.
package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	smithy "github.com/aws/smithy-go"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/service"
)

// Config carries the S3/MinIO/B2 connection settings. main.go maps the typed
// config.StorageConfig onto this so infrastructure stays free of the config
// package.
type Config struct {
	Endpoint        string // empty for AWS; e.g. http://localhost:9000 (MinIO) or https://s3.us-west-004.backblazeb2.com (B2)
	Region          string
	Bucket          string
	AccessKeyID     string // B2: the Application Key ID
	SecretAccessKey string // B2: the Application Key
	UsePathStyle    bool   // true for MinIO; false for AWS and B2 (virtual-hosted)

	// DisableChecksums turns off the AWS SDK's default per-request integrity
	// checksums (CRC32 upload trailers, added by default since aws-sdk-go-v2
	// v1.32). Backblaze B2 — and some other S3-compatible stores — reject those
	// headers/trailers, so uploads (and pre-signed PUTs) fail unless this is set.
	// Leave false for AWS S3 and MinIO.
	DisableChecksums bool
}

// S3Store implements service.ObjectStore over an S3-compatible API.
type S3Store struct {
	client  *s3.Client
	presign *s3.PresignClient
	bucket  string
}

var _ service.ObjectStore = (*S3Store)(nil)
var _ service.PrefixDeleter = (*S3Store)(nil)

// NewS3Store builds an S3Store. A bucket is required; static credentials are used
// when supplied (MinIO / B2 / explicit AWS keys), otherwise the default AWS chain.
func NewS3Store(ctx context.Context, cfg Config) (*S3Store, error) {
	if cfg.Bucket == "" {
		return nil, fmt.Errorf("storage: bucket is required")
	}

	opts := []func(*awsconfig.LoadOptions) error{awsconfig.WithRegion(cfg.Region)}
	if cfg.AccessKeyID != "" {
		opts = append(opts, awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		))
	}
	awsCfg, err := awsconfig.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("storage: load aws config: %w", err)
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		if cfg.Endpoint != "" { // MinIO / B2 / custom endpoint
			o.BaseEndpoint = aws.String(cfg.Endpoint)
		}
		o.UsePathStyle = cfg.UsePathStyle
		if cfg.DisableChecksums {
			// Only compute/validate a checksum when the operation requires one, so
			// the SDK stops adding CRC32 upload trailers that Backblaze B2 rejects.
			o.RequestChecksumCalculation = aws.RequestChecksumCalculationWhenRequired
			o.ResponseChecksumValidation = aws.ResponseChecksumValidationWhenRequired
		}
	})
	return &S3Store{
		client:  client,
		presign: s3.NewPresignClient(client),
		bucket:  cfg.Bucket,
	}, nil
}

// PresignPut returns a pre-signed PUT URL pinned to contentType so the client
// cannot upload a different media type than the upload usecase authorized.
func (s *S3Store) PresignPut(ctx context.Context, key, contentType string, ttl time.Duration) (string, error) {
	out, err := s.presign.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      &s.bucket,
		Key:         &key,
		ContentType: &contentType,
	}, s3.WithPresignExpires(ttl))
	if err != nil {
		return "", fmt.Errorf("storage: presign put %q: %w", key, err)
	}
	return out.URL, nil
}

// PresignGet returns a time-limited GET URL for a single object.
func (s *S3Store) PresignGet(ctx context.Context, key string, ttl time.Duration) (string, error) {
	out, err := s.presign.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: &s.bucket,
		Key:    &key,
	}, s3.WithPresignExpires(ttl))
	if err != nil {
		return "", fmt.Errorf("storage: presign get %q: %w", key, err)
	}
	return out.URL, nil
}

// Put uploads bytes server-side (transcode outputs).
func (s *S3Store) Put(ctx context.Context, key, contentType string, body io.Reader) error {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      &s.bucket,
		Key:         &key,
		Body:        body,
		ContentType: &contentType,
	})
	if err != nil {
		return fmt.Errorf("storage: put %q: %w", key, err)
	}
	return nil
}

// Get streams an object; the caller closes the returned reader.
func (s *S3Store) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	out, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &s.bucket,
		Key:    &key,
	})
	if err != nil {
		return nil, fmt.Errorf("storage: get %q: %w", key, err)
	}
	return out.Body, nil
}

// Exists reports object presence via HeadObject, mapping a 404/NotFound to
// (false, nil) rather than an error.
func (s *S3Store) Exists(ctx context.Context, key string) (bool, error) {
	_, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: &s.bucket,
		Key:    &key,
	})
	if err == nil {
		return true, nil
	}
	var notFound *types.NotFound
	if errors.As(err, &notFound) {
		return false, nil
	}
	var apiErr smithy.APIError
	if errors.As(err, &apiErr) {
		switch apiErr.ErrorCode() {
		case "NotFound", "NoSuchKey":
			return false, nil
		}
	}
	return false, fmt.Errorf("storage: head %q: %w", key, err)
}

// Delete removes an object. S3 DeleteObject is idempotent, so deleting a missing
// key is not an error.
func (s *S3Store) Delete(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &s.bucket,
		Key:    &key,
	})
	if err != nil {
		return fmt.Errorf("storage: delete %q: %w", key, err)
	}
	return nil
}

// DeletePrefix permanently deletes every object under prefix. It lists the keys
// (paginated) and deletes them one-by-one with DeleteObject rather than the
// batch DeleteObjects API: the batch call sends a Content-MD5/checksum the same
// way uploads did, which Backblaze B2 rejects (the reason DisableChecksums
// exists). A prefix with no objects is a no-op. Used to purge a video's HLS
// segments (hls/<id>/) and thumbnails (thumbnails/<id>/) on replace/delete.
func (s *S3Store) DeletePrefix(ctx context.Context, prefix string) error {
	paginator := s3.NewListObjectsV2Paginator(s.client, &s3.ListObjectsV2Input{
		Bucket: &s.bucket,
		Prefix: &prefix,
	})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("storage: list %q: %w", prefix, err)
		}
		for _, obj := range page.Contents {
			if obj.Key == nil {
				continue
			}
			if _, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
				Bucket: &s.bucket,
				Key:    obj.Key,
			}); err != nil {
				return fmt.Errorf("storage: delete %q: %w", *obj.Key, err)
			}
		}
	}
	return nil
}

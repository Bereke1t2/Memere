# Phase 3 · Skill 1 — Video Data Layer & Object Storage Abstraction

> **Prerequisite:** Phases 1 & 2 complete and green. Read
> [`docs/skill.md`](../../skill.md) §2 **Non-Negotiables** again — Phase 3 is
> governed by *"pre-signed CDN URLs only, no public S3"* and *"filter every query
> by the authenticated user_id"*.
>
> **Spec references:** `memere_Design_Specification.md` §4.2.4 (videos table),
> §8.1 (upload & processing pipeline), §8.2 (HLS adaptive bitrate), §3.3
> (object storage / CDN), README "Video Streaming" & "Video Endpoints".

---

## Goal

Build the **data + storage foundation** for video: the `videos` table Go layer
(entity, repo interface, sqlc) that Phase 1 left stubbed, and a clean
**object-storage abstraction** (`ObjectStore`) with an AWS S3 implementation plus
pre-signed URL generation. No HTTP, no transcoding yet — just the typed data layer
and the storage port the rest of Phase 3 depends on.

The whole point of the abstraction: the engines and handlers depend on an
**interface**, so S3 can be swapped for MinIO (local dev) or GCS later without
touching business logic (clean-architecture dependency rule).

---

## Key decisions for Phase 3 (monolith-friendly, swappable)

| Concern | Phase 3 choice | Why |
|---|---|---|
| Object storage | `ObjectStore` interface; **AWS S3** impl (`aws-sdk-go-v2`) | Spec-mandated S3; interface keeps domain clean. |
| Local dev storage | **MinIO** (S3-compatible) via docker-compose | No AWS account needed to run locally. |
| CDN signing | Helper for **CloudFront signed URLs**; fall back to S3 pre-signed in dev | Spec wants CDN delivery; S3 presign is the dev stand-in. |
| Upload | **Direct-to-S3 pre-signed PUT** (client uploads, server never proxies bytes) | Spec §8.1 step 2–3; keeps the API stateless/cheap. |
| Keys | Deterministic key scheme (below) | Predictable, collision-free, debuggable. |

### S3 key scheme (use everywhere)

```
originals/{course_id}/{lesson_id}/{video_id}/source{ext}
hls/{video_id}/master.m3u8
hls/{video_id}/{rendition}/playlist.m3u8        # rendition ∈ {480p,720p,1080p}
hls/{video_id}/{rendition}/seg_%05d.ts
thumbnails/{video_id}/thumb.jpg
```

---

## Tasks

### 1.1 — Config additions (`config/config.go`)

Extend the typed `Config` from Phase 1 with a `Storage` group. Add the keys to
`.env.example`.

```go
// config/config.go  (add to the existing Config struct)
type StorageConfig struct {
    Provider        string        `envconfig:"STORAGE_PROVIDER" default:"s3"` // s3 | minio
    Endpoint        string        `envconfig:"S3_ENDPOINT"`                    // empty for AWS; set for MinIO
    Region          string        `envconfig:"AWS_REGION" default:"af-south-1"`
    Bucket          string        `envconfig:"AWS_S3_BUCKET" required:"true"`
    AccessKeyID     string        `envconfig:"AWS_ACCESS_KEY_ID"`
    SecretAccessKey string        `envconfig:"AWS_SECRET_ACCESS_KEY"`
    UsePathStyle    bool          `envconfig:"S3_USE_PATH_STYLE" default:"false"` // true for MinIO
    UploadURLTTL    time.Duration `envconfig:"UPLOAD_URL_TTL" default:"15m"`
    StreamURLTTL    time.Duration `envconfig:"STREAM_URL_TTL" default:"2h"`   // spec §8.3: 2h
    DownloadURLTTL  time.Duration `envconfig:"DOWNLOAD_URL_TTL" default:"2h"`

    CDNDomain        string `envconfig:"CDN_DOMAIN"`         // e.g. dxxxx.cloudfront.net
    CDNKeyPairID     string `envconfig:"CDN_KEY_PAIR_ID"`    // CloudFront signing
    CDNPrivateKeyPEM string `envconfig:"CDN_PRIVATE_KEY_PEM"` // PEM (or path)
}
```

### 1.2 — Migrate the `videos` table Go-ready (`migrations/`)

Phase 1 migration `0003` created `courses.videos` as schema-only. Reconcile it
against spec §4.2.4 and, if anything is missing, add an **additive** migration
`0011_videos_enhancements` (continue numbering after Phase 2's `0010`). Ensure
these columns exist:

```sql
-- migrations/0011_videos_enhancements.up.sql  (only ALTERs for missing pieces)
ALTER TABLE courses.videos
    ADD COLUMN IF NOT EXISTS processing_status TEXT NOT NULL DEFAULT 'pending',
    ADD COLUMN IF NOT EXISTS original_file_key  TEXT,
    ADD COLUMN IF NOT EXISTS hls_master_key     TEXT,
    ADD COLUMN IF NOT EXISTS resolution_480p_key  TEXT,
    ADD COLUMN IF NOT EXISTS resolution_720p_key  TEXT,
    ADD COLUMN IF NOT EXISTS resolution_1080p_key TEXT,
    ADD COLUMN IF NOT EXISTS thumbnail_key      TEXT,
    ADD COLUMN IF NOT EXISTS duration_seconds   INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS file_size_bytes    BIGINT  NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS processing_error   TEXT,
    ADD COLUMN IF NOT EXISTS processed_at       TIMESTAMPTZ;

-- enforce one video per lesson (spec §4.2.4: lesson_id UNIQUE)
CREATE UNIQUE INDEX IF NOT EXISTS videos_lesson_id_uniq
    ON courses.videos (lesson_id) WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS videos_processing_status_idx
    ON courses.videos (processing_status);
```

`processing_status` enum values (use a typed string in Go; CHECK constraint in
SQL is fine): `pending` / `processing` / `ready` / `failed`.

`down.sql` drops the added columns/indexes in reverse.

### 1.3 — Domain entity (`internal/domain/entity/video.go`)

Pure Go, **no db/json tags** (same rule as every phase).

```go
package entity

import (
    "time"
    "github.com/google/uuid"
)

type VideoStatus string

const (
    VideoPending    VideoStatus = "pending"
    VideoProcessing VideoStatus = "processing"
    VideoReady      VideoStatus = "ready"
    VideoFailed     VideoStatus = "failed"
)

func (s VideoStatus) Valid() bool {
    switch s {
    case VideoPending, VideoProcessing, VideoReady, VideoFailed:
        return true
    }
    return false
}

// CanTransitionTo encodes the pipeline state machine (see Skill 3).
func (s VideoStatus) CanTransitionTo(next VideoStatus) bool {
    switch s {
    case VideoPending:    return next == VideoProcessing || next == VideoFailed
    case VideoProcessing: return next == VideoReady || next == VideoFailed
    case VideoFailed:     return next == VideoProcessing // allow retry
    default:              return false                    // ready is terminal
    }
}

type Video struct {
    ID              uuid.UUID
    LessonID        uuid.UUID
    CourseID        uuid.UUID // denormalized for fast authz/lookups
    Status          VideoStatus
    OriginalFileKey *string
    HLSMasterKey    *string
    Res480pKey      *string
    Res720pKey      *string
    Res1080pKey     *string
    ThumbnailKey    *string
    DurationSeconds int
    FileSizeBytes   int64
    ProcessingError *string
    ProcessedAt     *time.Time
    CreatedAt       time.Time
    UpdatedAt       time.Time
    DeletedAt       *time.Time
}
```

### 1.4 — Object storage port (`internal/domain/service/object_store.go`)

This interface lives in the **domain** layer (it's a port). Implementations live
in infrastructure.

```go
package service

import (
    "context"
    "io"
    "time"
)

// ObjectStore is the storage port. S3/MinIO/GCS implement it.
type ObjectStore interface {
    // PresignPut returns a URL the client uses to upload directly (PUT).
    PresignPut(ctx context.Context, key, contentType string, ttl time.Duration) (string, error)
    // PresignGet returns a time-limited download/stream URL for a single object.
    PresignGet(ctx context.Context, key string, ttl time.Duration) (string, error)
    // Put uploads bytes server-side (used by the transcode worker for outputs).
    Put(ctx context.Context, key, contentType string, body io.Reader) error
    // Get streams an object (used by the worker to fetch the source).
    Get(ctx context.Context, key string) (io.ReadCloser, error)
    Exists(ctx context.Context, key string) (bool, error)
    Delete(ctx context.Context, key string) error
}
```

### 1.5 — S3 implementation (`internal/infrastructure/storage/s3_store.go`)

Use `aws-sdk-go-v2` (`s3`, `s3/types`, `feature/s3/manager`, presign client).
Support MinIO via custom endpoint + path-style.

```go
package storage

import (
    "context"
    "io"
    "time"

    "github.com/aws/aws-sdk-go-v2/aws"
    awsconfig "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/credentials"
    "github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Store struct {
    client  *s3.Client
    presign *s3.PresignClient
    bucket  string
}

func NewS3Store(ctx context.Context, cfg AppStorageCfg) (*S3Store, error) {
    awsCfg, err := awsconfig.LoadDefaultConfig(ctx,
        awsconfig.WithRegion(cfg.Region),
        awsconfig.WithCredentialsProvider(
            credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
        ),
    )
    if err != nil {
        return nil, err
    }
    client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
        if cfg.Endpoint != "" { // MinIO / custom
            o.BaseEndpoint = aws.String(cfg.Endpoint)
        }
        o.UsePathStyle = cfg.UsePathStyle
    })
    return &S3Store{client: client, presign: s3.NewPresignClient(client), bucket: cfg.Bucket}, nil
}

func (s *S3Store) PresignPut(ctx context.Context, key, contentType string, ttl time.Duration) (string, error) {
    out, err := s.presign.PresignPutObject(ctx, &s3.PutObjectInput{
        Bucket:      &s.bucket,
        Key:         &key,
        ContentType: &contentType,
    }, s3.WithPresignExpires(ttl))
    if err != nil {
        return "", err
    }
    return out.URL, nil
}

// PresignGet, Put, Get, Exists, Delete: implement analogously.
// var _ service.ObjectStore = (*S3Store)(nil)
```

> ⚠️ The presigned PUT must pin `ContentType` (and ideally a max size via a POST
> policy) so a client can't upload arbitrary huge/foreign files. Note this; the
> upload usecase (Skill 2) sets the allowed content type.

### 1.6 — Repository interface + sqlc + impl

**Interface** (`internal/domain/repository/video_repository.go`):

```go
type VideoRepository interface {
    Create(ctx context.Context, v *entity.Video) error
    GetByID(ctx context.Context, id uuid.UUID) (*entity.Video, error)
    GetByLessonID(ctx context.Context, lessonID uuid.UUID) (*entity.Video, error)
    Update(ctx context.Context, v *entity.Video) error
    // UpdateStatus is a guarded transition: only flips if current status matches.
    UpdateStatus(ctx context.Context, id uuid.UUID, from, to entity.VideoStatus) (bool, error)
    ListByStatus(ctx context.Context, status entity.VideoStatus, limit int) ([]*entity.Video, error)
    SoftDelete(ctx context.Context, id uuid.UUID) error
}
```

**sqlc queries** (`db/query/videos.sql`) — note the guarded `UpdateStatus`:

```sql
-- name: CreateVideo :one
INSERT INTO courses.videos (id, lesson_id, course_id, processing_status, original_file_key, file_size_bytes)
VALUES (gen_random_uuid(), $1, $2, 'pending', $3, $4)
RETURNING *;

-- name: GetVideoByID :one
SELECT * FROM courses.videos WHERE id = $1 AND deleted_at IS NULL;

-- name: GetVideoByLessonID :one
SELECT * FROM courses.videos WHERE lesson_id = $1 AND deleted_at IS NULL;

-- name: UpdateVideoStatusGuarded :execrows
UPDATE courses.videos
SET processing_status = @to_status, updated_at = now()
WHERE id = @id AND processing_status = @from_status AND deleted_at IS NULL;

-- name: SetVideoReady :exec
UPDATE courses.videos
SET processing_status = 'ready',
    hls_master_key = $2, resolution_480p_key = $3, resolution_720p_key = $4,
    resolution_1080p_key = $5, thumbnail_key = $6, duration_seconds = $7,
    processed_at = now(), updated_at = now()
WHERE id = $1;

-- name: ListVideosByStatus :many
SELECT * FROM courses.videos
WHERE processing_status = $1 AND deleted_at IS NULL
ORDER BY created_at ASC LIMIT $2;
```

Implement `internal/repository/postgres/video_repository.go` over the generated
code (map sqlc model ↔ entity; `pgx.ErrNoRows` → `apperror.NotFound`;
`UpdateStatus` returns `rowsAffected > 0`). Add the compile-time assertion:
`var _ repository.VideoRepository = (*VideoRepo)(nil)`.

### 1.7 — Local dev: MinIO + Makefile

Add MinIO to `docker-compose.yml` and a bucket-bootstrap step.

```yaml
# docker-compose.yml  (add service)
  minio:
    image: minio/minio:latest
    command: server /data --console-address ":9001"
    environment:
      MINIO_ROOT_USER: ${AWS_ACCESS_KEY_ID:-minioadmin}
      MINIO_ROOT_PASSWORD: ${AWS_SECRET_ACCESS_KEY:-minioadmin}
    ports: ["9000:9000", "9001:9001"]
    volumes: ["minio_data:/data"]
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:9000/minio/health/live"]
      interval: 10s
      timeout: 5s
      retries: 5
# volumes: minio_data:
```

```make
# Makefile
minio-bucket: ## create the dev bucket in MinIO
	docker run --rm --network host --entrypoint sh minio/mc -c \
	 'mc alias set local http://localhost:9000 $(AWS_ACCESS_KEY_ID) $(AWS_SECRET_ACCESS_KEY) && \
	  mc mb -p local/$(AWS_S3_BUCKET) && mc anonymous set none local/$(AWS_S3_BUCKET)'
```

For dev `.env`: `STORAGE_PROVIDER=minio`, `S3_ENDPOINT=http://localhost:9000`,
`S3_USE_PATH_STYLE=true`, `AWS_S3_BUCKET=memere-media`, keys = `minioadmin`.

---

## Definition of Done

- [ ] `make migrate-up` applies `0011`; `make migrate-down` reverses it; the
      `videos.lesson_id` unique (partial) index exists.
- [ ] `make sqlc && go build ./...` clean; `golangci-lint run` clean.
- [ ] `entity.Video` has no db/json tags (`grep -rn 'db:"\|json:"'
      internal/domain/entity/video.go` is empty).
- [ ] `S3Store` satisfies `service.ObjectStore` (compile-time assertion present).
- [ ] `make up && make minio-bucket` creates the dev bucket; a tiny manual
      `PresignPut` → `curl -X PUT` upload → `Exists` returns true.
- [ ] `VideoRepository` impl satisfies the interface; guarded `UpdateStatus`
      returns false when the `from` status doesn't match (unit test).

## Verification commands

```bash
make migrate-up && make sqlc && go build ./... && golangci-lint run
grep -rn 'db:"\|json:"' internal/domain/entity/video.go && echo FAIL || echo OK
make up && make minio-bucket
go test ./internal/repository/postgres/... -run Video -v
```

## Hand-off to Skill 2

The video data layer and storage port exist. Skill 2 builds the **upload flow**:
the pre-signed upload-URL usecase, the create-video record, the upload-confirm
endpoint, and enqueuing a transcode job onto the message-queue abstraction.

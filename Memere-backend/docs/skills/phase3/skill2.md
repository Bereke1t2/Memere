# Phase 3 · Skill 2 — Upload Flow & Job Queue Abstraction

> **Prerequisite:** Phase 3 Skill 1 done (video entity/repo/sqlc, `ObjectStore`
> port + S3/MinIO impl, MinIO running). Read [`docs/skill.md`](../../skill.md) §2.
>
> **Spec references:** `memere_Design_Specification.md` §8.1 (steps 1–4: request
> pre-signed URL → client uploads to S3 → S3 event → SQS message), §3.3 (message
> queue), README "Video Endpoints" (`POST /videos/upload-url`).

---

## Goal

Implement the **upload half** of the pipeline (spec §8.1 steps 1–4) as usecases +
a **message-queue abstraction** — no HTTP yet (that's Skill 5), no transcoding yet
(Skill 3). A teacher requests a pre-signed upload URL, uploads directly to S3,
confirms the upload, and the system enqueues a transcode job.

We replace the spec's "S3 event → SQS" auto-trigger with an explicit
**confirm-upload** call that enqueues the job. This is simpler and fully testable
in the monolith; the `JobQueue` interface keeps the door open to swap in real SQS
+ S3 event notifications later without touching usecases.

---

## Tasks

### 2.1 — Job queue port (`internal/domain/service/job_queue.go`)

```go
package service

import "context"

type TranscodeJob struct {
    VideoID        string `json:"video_id"`
    CourseID       string `json:"course_id"`
    LessonID       string `json:"lesson_id"`
    SourceKey      string `json:"source_key"`
    AttemptCount   int    `json:"attempt_count"`
}

// JobQueue is the async-work port. In-process channel impl now; SQS/RabbitMQ later.
type JobQueue interface {
    EnqueueTranscode(ctx context.Context, job TranscodeJob) error
}
```

### 2.2 — In-process queue impl (`internal/infrastructure/messaging/inproc_queue.go`)

A buffered-channel implementation for the monolith. The transcode worker
(Skill 3) consumes from it.

```go
package messaging

import (
    "context"
    "github.com/.../internal/domain/service"
)

type InProcQueue struct {
    transcode chan service.TranscodeJob
}

func NewInProcQueue(buffer int) *InProcQueue {
    return &InProcQueue{transcode: make(chan service.TranscodeJob, buffer)}
}

func (q *InProcQueue) EnqueueTranscode(ctx context.Context, job service.TranscodeJob) error {
    select {
    case q.transcode <- job:
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}

// Transcode exposes the receive side to the worker (not part of the port).
func (q *InProcQueue) Transcode() <-chan service.TranscodeJob { return q.transcode }

// var _ service.JobQueue = (*InProcQueue)(nil)
```

> Design note: jobs in a channel are lost on restart. Acceptable for Phase 3
> (the sweeper in 2.5 re-enqueues stuck `pending`/`processing` videos on boot).
> Phase 6 swaps this for a durable SQS consumer behind the same `JobQueue` port.

### 2.3 — Upload usecases (`internal/usecase/video/`)

A `Service` holding `videoRepo`, `lessonRepo` (Phase 1), `courseRepo` (for
ownership), `ObjectStore`, `JobQueue`, and `Config`. All methods take `ctx` +
`actor` and **enforce course ownership** (teacher owns the lesson's course, or
admin) — never trust client IDs (IDOR rule).

```go
package video

import (
    "context"
    "fmt"
    "path"
    "github.com/google/uuid"
)

type RequestUploadInput struct {
    LessonID    uuid.UUID
    FileName    string
    ContentType string // must be video/*
    SizeBytes   int64
}

type RequestUploadResult struct {
    VideoID   uuid.UUID `json:"video_id"`
    UploadURL string    `json:"upload_url"`
    SourceKey string    `json:"source_key"`
    ExpiresIn int       `json:"expires_in"`
}

func (s *Service) RequestUpload(ctx context.Context, actor Actor, in RequestUploadInput) (*RequestUploadResult, error) {
    lesson, err := s.lessonRepo.GetByID(ctx, in.LessonID)
    if err != nil {
        return nil, err // NotFound mapped in repo
    }
    if err := s.assertCourseOwner(ctx, actor, lesson.CourseID); err != nil {
        return nil, err // Forbidden
    }
    if !isAllowedVideoType(in.ContentType) {
        return nil, apperror.BadRequest("UNSUPPORTED_MEDIA_TYPE", "only video/mp4|quicktime|webm allowed")
    }
    if in.SizeBytes <= 0 || in.SizeBytes > s.cfg.MaxUploadBytes {
        return nil, apperror.BadRequest("FILE_TOO_LARGE", "exceeds max upload size")
    }

    videoID := uuid.New()
    ext := path.Ext(in.FileName)
    key := fmt.Sprintf("originals/%s/%s/%s/source%s", lesson.CourseID, in.LessonID, videoID, ext)

    v := &entity.Video{ID: videoID, LessonID: in.LessonID, CourseID: lesson.CourseID,
        Status: entity.VideoPending, OriginalFileKey: &key, FileSizeBytes: in.SizeBytes}
    if err := s.videoRepo.Create(ctx, v); err != nil {
        return nil, err // Conflict if lesson already has a video (unique index)
    }

    url, err := s.store.PresignPut(ctx, key, in.ContentType, s.cfg.UploadURLTTL)
    if err != nil {
        return nil, apperror.Internal(err)
    }
    return &RequestUploadResult{VideoID: videoID, UploadURL: url, SourceKey: key,
        ExpiresIn: int(s.cfg.UploadURLTTL.Seconds())}, nil
}
```

- `ConfirmUpload(ctx, actor, videoID)`:
  - ownership check; load video; verify the object actually exists in storage
    (`store.Exists(sourceKey)`) — reject if the client never uploaded.
  - guarded transition `pending → processing` (`videoRepo.UpdateStatus`).
  - `jobQueue.EnqueueTranscode(...)`. If enqueue fails, roll status back to
    `pending` so it can be retried.
- `GetVideoStatus(ctx, actor, videoID)` — ownership or enrolled-student read;
  returns status + (when ready) duration/thumbnail info. **Never** returns raw
  keys to the client; streaming URLs come from Skill 4.
- `RetryProcessing(ctx, actor, videoID)` — owner/admin; `failed → processing`,
  re-enqueue (bounded attempts).

Helpers: `isAllowedVideoType` (whitelist `video/mp4`, `video/quicktime`,
`video/webm`), `assertCourseOwner` (reuse Phase 1 course-ownership logic — extract
it into a shared helper if not already).

### 2.4 — Config additions

Add `MaxUploadBytes` (`MAX_UPLOAD_BYTES`, default e.g. `2147483648` = 2 GiB),
`TRANSCODE_QUEUE_BUFFER` (default 64), `TRANSCODE_MAX_ATTEMPTS` (default 3) to the
typed config + `.env.example`.

### 2.5 — Boot reconciler (stuck-job recovery)

Because the in-proc queue is non-durable, add a small startup step
(`video.Service.RequeueStuck(ctx)`): on boot, find videos in `processing` (and
optionally `pending` with an existing source object) via
`videoRepo.ListByStatus` and re-enqueue them. Wired in `main.go` (Skill 5).

### 2.6 — Unit tests (mocked store + queue + repos)

Cover: request-upload sets `pending` + returns a presigned URL with the right key;
non-owner → `FORBIDDEN`; disallowed content type / oversize → `400`; duplicate
video per lesson → `CONFLICT`; confirm without an uploaded object → error; confirm
happy path → `processing` + job enqueued; enqueue failure rolls back to `pending`;
retry only from `failed` and respects max attempts.

---

## Definition of Done

- [ ] `go build ./...` clean; `golangci-lint run` clean.
- [ ] `go test ./internal/usecase/video/...` passes (all cases above).
- [ ] `InProcQueue` satisfies `service.JobQueue`; `S3Store` still satisfies
      `ObjectStore` (assertions compile).
- [ ] Presigned PUT pins content type; oversize/foreign types rejected before any
      DB write.
- [ ] `ConfirmUpload` verifies the object exists before transitioning.
- [ ] Status transitions are guarded (`pending→processing` only once); enqueue
      failure does not leave a video stuck in `processing`.
- [ ] No raw storage keys are returned to clients from `GetVideoStatus`.

## Verification commands

```bash
go build ./... && golangci-lint run
go test ./internal/usecase/video/... -v
# Manual presign round-trip against MinIO (after wiring in Skill 5, or via a small main):
#   request-upload -> curl -X PUT --upload-file sample.mp4 "<upload_url>" -> confirm-upload
```

## Hand-off to Skill 3

Uploads land in storage and a transcode job is enqueued. Skill 3 builds the
**transcode worker**: consume the job, pull the source, run FFmpeg to produce HLS
renditions + thumbnail, upload outputs, and flip the video to `ready` (or
`failed` with retry).

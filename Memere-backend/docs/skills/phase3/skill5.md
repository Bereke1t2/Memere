# Phase 3 · Skill 5 — HTTP Delivery & Wiring (Video API + Worker)

> **Prerequisite:** Phase 3 Skills 1–4 done (data layer, upload flow, transcode
> worker, secure delivery — all tested in isolation). Read
> [`docs/skill.md`](../../skill.md) §2 and reuse the Phase 1 delivery layer
> (apperror envelope, middleware, router, constructor wiring).
>
> **Spec references:** `memere_Design_Specification.md` §8 (all), §5.4 (API
> conventions), README "Video Endpoints", §7.2 (RBAC).

---

## Goal

Expose the video pipeline over the existing Gin API and wire the **in-proc queue +
transcode worker + signer** into the app lifecycle. After this skill, Phase 3 is
**fully runnable end-to-end** against MinIO + FFmpeg: a teacher uploads a video, it
transcodes to HLS, and a student gets a short-lived, access-controlled stream URL.

Extend the Phase 1/2 delivery layer — do not rebuild it. Reuse `apperror` mapping,
the middleware stack, `actorFromContext`, and the constructor-wiring style.

---

## API surface for Phase 3 (from README + spec §8)

Base path `/api/v1`.

**Upload / management** (`video_handler.go`)
| Method | Path | Auth | Notes |
|---|---|---|---|
| POST | `/lessons/:id/videos/upload-url` | bearer + owner/admin | request pre-signed PUT; creates `pending` video |
| POST | `/videos/:id/confirm` | bearer + owner/admin | verify object exists → `processing` + enqueue |
| GET | `/videos/:id/status` | bearer (owner/admin or allowed student) | processing state + meta (no raw keys) |
| POST | `/videos/:id/retry` | bearer + owner/admin | `failed → processing`, re-enqueue |

**Delivery**
| Method | Path | Auth | Notes |
|---|---|---|---|
| GET | `/videos/:id/stream` | bearer (allowed) | signed HLS master URL (2h) |
| GET | `/videos/:id/download-url` | bearer (allowed) | signed URL + single-use token |
| GET | `/videos/download/:token` | bearer (allowed) | consume token → 302 redirect to signed manifest |

(Payment/enrollment, progress, notifications remain later phases.)

---

## Tasks

### 5.1 — DTOs (`internal/delivery/http/dto/video.go`)

Request/response structs with `json` tags. **Never expose raw storage keys or
`processing_error` internals to students** (owner/admin may see a sanitized error).

```go
package dto

type RequestUploadRequest struct {
    FileName    string `json:"file_name" binding:"required"`
    ContentType string `json:"content_type" binding:"required"`
    SizeBytes   int64  `json:"size_bytes" binding:"required,gt=0"`
}

type RequestUploadResponse struct {
    VideoID   string `json:"video_id"`
    UploadURL string `json:"upload_url"`
    ExpiresIn int    `json:"expires_in"`
    // note: SourceKey intentionally omitted from the client response
}

type VideoStatusResponse struct {
    VideoID         string `json:"video_id"`
    Status          string `json:"status"`           // pending|processing|ready|failed
    DurationSeconds int    `json:"duration_seconds"`
    Error           string `json:"error,omitempty"`  // populated only for owner/admin
}

type StreamResponse struct {
    MasterURL    string `json:"master_url"`
    ExpiresIn    int    `json:"expires_in"`
    ThumbnailURL string `json:"thumbnail_url,omitempty"`
    DurationSec  int    `json:"duration_seconds"`
}

type DownloadURLResponse struct {
    DownloadURL string `json:"download_url"`
    Token       string `json:"token"`
    ExpiresIn   int    `json:"expires_in"`
}
```

### 5.2 — Handlers (`internal/delivery/http/video_handler.go`)

Bind+validate, pull `actor` from context, call the Phase 3 usecases, map to DTOs,
set status codes. Sketch:

```go
func (h *VideoHandler) RequestUpload(c *gin.Context) {
    actor, _ := actorFromContext(c)
    lessonID, err := uuid.Parse(c.Param("id"))
    if err != nil { respondError(c, apperror.BadRequest("INVALID_ID", "bad lesson id")); return }

    var req dto.RequestUploadRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        respondError(c, apperror.Validation(bindingDetails(err))); return
    }
    out, err := h.uc.RequestUpload(c, actor, video.RequestUploadInput{
        LessonID: lessonID, FileName: req.FileName,
        ContentType: req.ContentType, SizeBytes: req.SizeBytes,
    })
    if err != nil { respondError(c, err); return }
    respondJSON(c, http.StatusCreated, dto.RequestUploadResponse{
        VideoID: out.VideoID.String(), UploadURL: out.UploadURL, ExpiresIn: out.ExpiresIn,
    })
}

func (h *VideoHandler) Stream(c *gin.Context) {
    actor, _ := actorFromContext(c)
    id, err := uuid.Parse(c.Param("id"))
    if err != nil { respondError(c, apperror.BadRequest("INVALID_ID", "bad video id")); return }
    res, err := h.uc.GetStreamURL(c, actor, id)
    if err != nil { respondError(c, err); return }
    respondJSON(c, http.StatusOK, dto.StreamResponse{
        MasterURL: res.MasterURL, ExpiresIn: res.ExpiresIn,
        ThumbnailURL: res.ThumbnailURL, DurationSec: res.DurationSec,
    })
}

func (h *VideoHandler) ConsumeDownload(c *gin.Context) {
    // _, _ = actorFromContext(c) // still require auth at the route
    token := c.Param("token")
    url, err := h.uc.ConsumeDownloadToken(c, token)
    if err != nil { respondError(c, err); return } // 404/410 on miss
    c.Redirect(http.StatusFound, url)
}
```

`Confirm`, `Status`, `Retry`, `DownloadURL` follow the same pattern.

### 5.3 — Routes & RBAC (`router.go`)

Register in the existing `/api/v1` group:

```go
v := r.Group("/api/v1")
{
    vids := v.Group("/videos")
    vids.Use(mw.RequireAuth())
    {
        vids.GET("/:id/status",       h.Video.Status)
        vids.GET("/:id/stream",       h.Video.Stream)
        vids.GET("/:id/download-url", h.Video.DownloadURL)
        vids.GET("/download/:token",  h.Video.ConsumeDownload)

        owner := vids.Group("")
        owner.Use(mw.RequireRole(entity.RoleTeacher, entity.RoleAdmin))
        owner.POST("/:id/confirm", h.Video.Confirm)
        owner.POST("/:id/retry",   h.Video.Retry)
    }
    v.POST("/lessons/:id/videos/upload-url",
        mw.RequireAuth(), mw.RequireRole(entity.RoleTeacher, entity.RoleAdmin),
        h.Video.RequestUpload)
}
```

Ownership of the specific course/lesson is enforced **in the usecase** (not
middleware) — middleware only gates by role.

### 5.4 — Wire the queue + worker + signer into `main.go`

Extend the Phase 1/2 wiring:

```go
// cmd/api/main.go (additions)
store, err := storage.NewS3Store(ctx, cfg.StorageCfg())          // S3/MinIO
signer := storage.NewCDNSigner(cfg, store)                        // CloudFront or presign fallback
queue := messaging.NewInProcQueue(cfg.Storage.QueueBuffer)
videoRepo := postgres.NewVideoRepo(pool, q)                       // q = sqlc Queries
videoUC := video.NewService(videoRepo, lessonRepo, courseRepo, store, signer, queue, cfg)

coder := transcode.NewFFmpeg()
tw := worker.NewTranscodeWorker(queue.Transcode(), store, videoRepo, coder, queue, cfg.WorkerCfg())

// boot reconciler: re-enqueue stuck videos (Skill 2)
_ = videoUC.RequeueStuck(ctx)

go tw.Run(ctx) // ctx cancelled on SIGINT/SIGTERM -> worker drains & stops

// handlers + routes
videoHandler := httpx.NewVideoHandler(videoUC)
// register in NewRouter(...)
```

Ensure the worker stops on the same shutdown signal as the HTTP server (it already
selects on `ctx.Done()`).

### 5.5 — End-to-end smoke test (`scripts/smoke_phase3.sh`)

Against a running stack (`make up && make minio-bucket && make migrate-up &&
make run`), with ffmpeg installed and a tiny `sample.mp4` committed under
`testdata/`:

1. login teacher (Phase 1); create course + section + lesson (`type=video`).
2. `POST /lessons/:id/videos/upload-url` → get `upload_url` + `video_id`.
3. `curl -X PUT --upload-file testdata/sample.mp4 -H "Content-Type: video/mp4"
   "<upload_url>"` → 200 (direct to MinIO).
4. `POST /videos/:id/confirm` → `202/200`; status becomes `processing`.
5. poll `GET /videos/:id/status` until `ready` (worker ran FFmpeg) — assert
   `duration_seconds > 0` and **no raw key** in the response.
6. as the teacher, `GET /videos/:id/stream` → a signed `master_url`; `curl -I`
   it → 200, and confirm it carries an expiry/signature.
7. as a **student** with no access to a **paid** course → `GET .../stream` →
   `403`; mark the lesson `is_free_preview` → student now gets a URL.
8. `GET /videos/:id/download-url` → token; `GET /videos/download/:token` →
   302; calling the same token again → `404/410` (single-use).

### 5.6 — Docs

- Update `README.md` "Video Streaming"/"Getting Started" if commands changed
  (note `make minio-bucket`, ffmpeg dependency).
- Optionally extend the `api/` OpenAPI stub with the Phase 3 endpoints.

---

## Definition of Done

- [ ] `go build ./...` clean; `golangci-lint run` clean; `go test ./...` passes.
- [ ] `make up && make minio-bucket && make migrate-up && make run` serves the
      Phase 3 endpoints; the worker starts and stops with the app.
- [ ] `scripts/smoke_phase3.sh` passes every assertion (upload → transcode →
      ready → signed stream → access control → single-use download).
- [ ] No video response exposes a raw storage key; `processing_error` is shown
      only to owner/admin.
- [ ] Stream/download URLs are signed and expire (≈2h); download token is
      single-use.
- [ ] Access control gates paid content to owner/admin (with `TODO(phase4)`
      enrollment hook) and allows free/preview to authenticated students.
- [ ] Transcode worker concurrency bounded; failures retried within budget then
      terminal `failed`.

## Verification commands

```bash
make up && make minio-bucket && make migrate-up
make run &
bash scripts/smoke_phase3.sh
go test ./... && golangci-lint run
grep -rn 'original_file_key\|hls_master_key\|OriginalFileKey\|HLSMasterKey' \
  internal/delivery/http/dto && echo "FAIL: raw key in a DTO" || echo "OK"
```

---

## 🎉 Phase 3 complete — what now

When this Definition of Done passes, the video pipeline is live: direct-to-S3
pre-signed upload, FFmpeg→HLS transcoding via a bounded async worker, a guarded
processing state machine, and short-lived, access-controlled, single-use-download
delivery — all on the Phase 1/2 foundation and honoring every Non-Negotiable.

**To proceed to Phase 4:**
1. Report "Phase 3 complete" with a one-paragraph summary + the smoke-test result.
2. Ask Claude to author the **Phase 4** skill files into `docs/skills/phase4/` —
   **Payments & Enrollments**: Chapa/Telebirr/Stripe behind a provider
   abstraction, idempotency keys, webhook verification/dedup, the enrollment flow,
   and **replacing every `TODO(phase4)` enrollment hook** left in Phases 2–3
   (quiz/exam taking + paid video access) with real enrollment checks (spec §10).
3. Do **not** scaffold Phase 4 code before its skills are written and reviewed.

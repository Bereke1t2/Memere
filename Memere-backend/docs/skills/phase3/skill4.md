# Phase 3 · Skill 4 — Secure Streaming & Offline Download URLs

> **Prerequisite:** Phase 3 Skills 1–3 done (videos reach `ready` with an HLS
> ladder in storage). Read [`docs/skill.md`](../../skill.md) §2 — this skill is
> the direct expression of *"pre-signed CDN URLs only, no public S3"* and
> *"filter every query by the authenticated user_id"*.
>
> **Spec references:** `memere_Design_Specification.md` §8.3 (offline download
> strategy — 2-hour signed URL, single-use, 30-day expiry), §7.3 (video
> hotlinking mitigation), README "Video Endpoints" (`/stream`, `/download-url`).

---

## Goal

Issue **time-limited, access-controlled** URLs so students can stream (HLS) and
download videos — and **no one** can hotlink or share durable links. Built as
usecases (HTTP in Skill 5). This is the security-critical end of the pipeline:
every URL is short-lived, scoped to one authenticated user, and only granted when
access rules pass.

---

## Access rules (who may get a URL)

A streaming/download URL is issued **only if** the caller is:

1. the **course teacher** (owns it) or an **admin**, OR
2. a **student with valid access** to the course — for Phase 3 that means the
   lesson is `is_free_preview = true` **or** the course `is_free = true`. Paid
   enrollment checks arrive in **Phase 4**; until then, gate paid content to
   owner/admin and free/preview content to any authenticated student. Mark the
   enrollment hook with `TODO(phase4)`.

Always resolve the video through the lesson → course chain server-side and check
authorization against the **authenticated** user — never trust a client-supplied
course/enrollment claim (IDOR).

---

## Tasks

### 4.1 — CDN signing helper (`internal/infrastructure/storage/cdn_signer.go`)

Prefer **CloudFront signed URLs** when a distribution is configured; fall back to
S3/MinIO pre-signed GET in dev.

```go
package storage

import (
    "context"
    "time"
)

// URLSigner abstracts "give me a short-lived GET URL for this key".
type URLSigner interface {
    SignGet(ctx context.Context, key string, ttl time.Duration) (string, error)
}

// CDNSigner signs via CloudFront if configured, else delegates to the ObjectStore.
type CDNSigner struct {
    cdnDomain  string
    keyPairID  string
    privKey    *rsa.PrivateKey
    fallback   service.ObjectStore // S3/MinIO presign for dev
}

func (c *CDNSigner) SignGet(ctx context.Context, key string, ttl time.Duration) (string, error) {
    if c.cdnDomain == "" || c.privKey == nil {
        return c.fallback.PresignGet(ctx, key, ttl) // dev path
    }
    // Build https://{cdnDomain}/{key} and sign with a canned policy (Expires = now+ttl)
    // using aws cloudfront sign (github.com/aws/aws-sdk-go-v2/feature/cloudfront/sign).
    // Return the signed URL.
    ...
}
```

> ⚠️ `now+ttl` requires a real clock. Inject a `Clock` (Phase 2 already
> introduced one) so this is testable; production uses wall-clock.

### 4.2 — Single-use download tokens (Redis)

Spec §8.3 step 2 calls the download URL "single-use". S3/CloudFront signed URLs
are time-limited but **replayable** within the window. To honor single-use, add a
server-side gate for the *download* flow:

```go
// internal/repository/redis/download_token.go
type DownloadTokenStore struct{ rdb *redis.Client }

// Issue stores a one-time token -> videoID+userID with the URL TTL.
func (s *DownloadTokenStore) Issue(ctx context.Context, token, videoID, userID string, ttl time.Duration) error {
    return s.rdb.Set(ctx, "dl:"+token, videoID+"|"+userID, ttl).Err()
}

// Consume atomically fetches+deletes (GETDEL) — second use returns miss.
func (s *DownloadTokenStore) Consume(ctx context.Context, token string) (videoID, userID string, ok bool, err error) {
    val, err := s.rdb.GetDel(ctx, "dl:"+token).Result()
    if err == redis.Nil { return "", "", false, nil }
    if err != nil { return "", "", false, err }
    // split val -> videoID, userID
    ...
}
```

The download endpoint (Skill 5) returns a manifest URL plus this token; the client
calls a `confirm`/redirect route that `Consume`s the token before redirecting to
the signed object. (Streaming is *not* single-use — HLS needs repeated segment
fetches within the window.)

### 4.3 — Delivery usecases (`internal/usecase/video/delivery.go`)

Extend the Phase 3 `video.Service` (or a sibling `DeliveryService`):

```go
type StreamResult struct {
    MasterURL      string `json:"master_url"`
    ExpiresIn      int    `json:"expires_in"`
    DurationSec    int    `json:"duration_seconds"`
    ThumbnailURL   string `json:"thumbnail_url,omitempty"`
}

func (s *Service) GetStreamURL(ctx context.Context, actor Actor, videoID uuid.UUID) (*StreamResult, error) {
    v, err := s.videoRepo.GetByID(ctx, videoID)
    if err != nil { return nil, err }
    if v.Status != entity.VideoReady || v.HLSMasterKey == nil {
        return nil, apperror.Conflict("VIDEO_NOT_READY", "video is still processing")
    }
    if err := s.assertCanWatch(ctx, actor, v); err != nil { // owner/admin/free/preview
        return nil, err // Forbidden
    }
    master, err := s.signer.SignGet(ctx, *v.HLSMasterKey, s.cfg.StreamURLTTL)
    if err != nil { return nil, apperror.Internal(err) }
    res := &StreamResult{MasterURL: master, ExpiresIn: int(s.cfg.StreamURLTTL.Seconds()), DurationSec: v.DurationSeconds}
    if v.ThumbnailKey != nil {
        res.ThumbnailURL, _ = s.signer.SignGet(ctx, *v.ThumbnailKey, s.cfg.StreamURLTTL)
    }
    return res, nil
}
```

- `GetDownloadURL(ctx, actor, videoID)`:
  - same readiness + access checks.
  - mint a single-use token (4.2), store it, return
    `{ download_url: signed-master, token, expires_in }`. (Or return a
    server-relative `confirm` URL that carries the token — Skill 5 decides the
    exact shape.)
  - record the intent for the offline-expiry rule (client enforces the 30-day
    re-download per §8.3 step 7; server just bounds the signed window to 2h).
- `ConsumeDownloadToken(ctx, token) (signedManifestURL, error)` — `Consume` the
  token; if valid, re-sign and return the manifest URL; second attempt → 410/404.

> **HLS caveat:** signed master/variant playlists reference **segment** URLs.
> Segments must also be access-controlled. With CloudFront, sign at the
> distribution (cookie or wildcard signed URL for `hls/{video_id}/*`). With S3
> presign in dev, either (a) rewrite the playlist to embed presigned segment URLs,
> or (b) accept that dev uses a permissive bucket policy scoped to the prefix.
> **Document the choice**; production must use CloudFront signed cookies/URLs so
> segments aren't public.

### 4.4 — Access-control helper

`assertCanWatch(ctx, actor, video)`:
- load lesson + course (denormalized `course_id` on the video makes this one
  fetch); owner/admin → allow; else if `course.is_free || lesson.is_free_preview`
  → allow for authenticated students; else → `apperror.Forbidden`
  with `TODO(phase4): replace with enrollment check`.

### 4.5 — Config

Confirm `StreamURLTTL` (2h) and `DownloadURLTTL` (2h) from Skill 1 are wired. Add
`CDN_*` keys to `.env.example` (already declared in Skill 1 config).

### 4.6 — Tests

- owner gets a stream URL; unauthenticated → `401` (enforced at HTTP layer);
  student blocked from paid non-preview content → `403`; student allowed for free
  course / preview lesson → URL issued.
- `VIDEO_NOT_READY` when status ≠ `ready`.
- signed URL carries an expiry ≈ TTL (fake clock); changing the key changes the
  signature.
- download token is single-use: first `Consume` ok, second → miss.
- HLS segment access: with the chosen strategy, a segment URL outside the signed
  scope/window is rejected (document + test what you can).

---

## Definition of Done

- [ ] `go build ./...` clean; `golangci-lint run` clean.
- [ ] `go test ./internal/usecase/video/... ./internal/repository/redis/...`
      passes (delivery + token tests).
- [ ] No endpoint ever returns a raw S3 key or a non-expiring URL.
- [ ] Stream/download URLs expire in ≈2h (from config, not hardcoded).
- [ ] Download tokens are single-use (Redis `GETDEL`); replay fails.
- [ ] Access control resolves video→lesson→course server-side and checks the
      authenticated user (no IDOR); paid content gated to owner/admin for now
      with a `TODO(phase4)` enrollment hook.
- [ ] The HLS segment-access strategy is documented and the production path uses
      signed CloudFront (not public S3).

## Verification commands

```bash
go build ./... && golangci-lint run
go test ./internal/usecase/video/... ./internal/repository/redis/... -run 'Stream|Download|Token' -v
grep -rn 'PresignGet\|SignGet' internal/usecase internal/delivery && echo "review: URLs only via signer"
```

## Hand-off to Skill 5

Secure delivery logic is ready. Skill 5 exposes the full video API over HTTP
(upload-url, confirm, status, stream, download-url, retry), wires the transcode
worker + in-proc queue + signer into `main.go`, and runs the Phase 3 end-to-end
smoke test against MinIO + FFmpeg.

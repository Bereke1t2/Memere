# Phase 3 · Skill 3 — Transcode Worker (FFmpeg → HLS) & Pipeline State Machine

> **Prerequisite:** Phase 3 Skills 1–2 done (storage port, video repo, job queue,
> upload flow enqueues jobs). Read [`docs/skill.md`](../../skill.md) §2. FFmpeg
> must be installed where the worker runs.
>
> **Spec references:** `memere_Design_Specification.md` §8.1 (steps 5–9:
> transcode to HLS 480/720/1080, generate thumbnail, upload segments + manifest,
> set status=ready, invalidate CDN, notify), §8.2 (adaptive bitrate ladder table).

---

## Goal

Build the **transcode worker** that consumes a `TranscodeJob`, downloads the
source from object storage, runs **FFmpeg** to produce an HLS adaptive-bitrate
ladder (480p/720p/1080p) plus a thumbnail, uploads the outputs, and transitions
the video `processing → ready` (or `failed` with bounded retry). This is the
async backbone of §8.1.

The worker depends only on the `JobQueue` (consume), `ObjectStore` (get source /
put outputs), `VideoRepository` (state), and an `FFmpeg` runner abstraction — so
the actual transcoder can later be swapped for AWS MediaConvert behind the same
boundary.

---

## The HLS ladder (spec §8.2) — encode exactly these

| Rendition | Resolution | Video bitrate | Audio | Use case |
|---|---|---|---|---|
| 480p | 854×480 | 800 kbps | 96 kbps | 3G |
| 720p | 1280×720 | 1500 kbps | 128 kbps | 4G / good WiFi |
| 1080p | 1920×1080 | 3000 kbps | 128 kbps | strong WiFi |

(360p is in the spec's player table but the §4.2.4 video table only stores
480/720/1080 keys — encode 480/720/1080. Note this reconciliation in a comment.)

Output a **master `.m3u8`** referencing all three variant playlists, 6-second
segments, and a thumbnail at the 5-second mark.

---

## Pipeline state machine (enforced via `VideoStatus.CanTransitionTo`)

```
pending ──confirm-upload──▶ processing ──success──▶ ready (terminal)
   ▲                            │
   │                            └──error──▶ failed ──retry──▶ processing
   └── (boot reconciler re-enqueues stuck processing)
```

Use the **guarded** `videoRepo.UpdateStatus(from, to)` for every flip so a retry
and the worker can't double-process (returns false if `from` didn't match).

---

## Tasks

### 3.1 — FFmpeg runner abstraction (`internal/infrastructure/transcode/ffmpeg.go`)

```go
package transcode

import (
    "context"
    "os/exec"
)

type Renditions struct {
    MasterPath    string            // local path to master.m3u8
    VariantDirs   map[string]string // "480p" -> local dir
    ThumbnailPath string
    DurationSec   int
}

type Transcoder interface {
    // ToHLS reads a local source file, writes HLS outputs into workDir, returns paths.
    ToHLS(ctx context.Context, srcPath, workDir string) (*Renditions, error)
    Probe(ctx context.Context, srcPath string) (durationSec int, err error)
}

type FFmpeg struct {
    bin       string // "ffmpeg"
    ffprobe   string // "ffprobe"
}

func NewFFmpeg() *FFmpeg { return &FFmpeg{bin: "ffmpeg", ffprobe: "ffprobe"} }
```

Sketch of the `ToHLS` command construction (one master + three variants, single
ffmpeg invocation using `-var_stream_map`):

```go
func (f *FFmpeg) ToHLS(ctx context.Context, src, workDir string) (*Renditions, error) {
    // -filter_complex split into 3 scaled outputs; map each to a bitrate;
    // HLS muxer with master playlist + per-variant playlists.
    args := []string{
        "-y", "-i", src,
        "-filter_complex",
        "[0:v]split=3[v1][v2][v3];" +
            "[v1]scale=w=854:h=480[v480];" +
            "[v2]scale=w=1280:h=720[v720];" +
            "[v3]scale=w=1920:h=1080[v1080]",
        // 480p
        "-map", "[v480]", "-map", "0:a?", "-c:v:0", "libx264", "-b:v:0", "800k", "-c:a", "aac", "-b:a", "96k",
        // 720p
        "-map", "[v720]", "-map", "0:a?", "-c:v:1", "libx264", "-b:v:1", "1500k",
        // 1080p
        "-map", "[v1080]", "-map", "0:a?", "-c:v:2", "libx264", "-b:v:2", "3000k",
        "-f", "hls",
        "-hls_time", "6",
        "-hls_playlist_type", "vod",
        "-hls_segment_filename", workDir + "/%v/seg_%05d.ts",
        "-master_pl_name", "master.m3u8",
        "-var_stream_map", "v:0,a:0,name:480p v:1,a:1,name:720p v:2,a:2,name:1080p",
        workDir + "/%v/playlist.m3u8",
    }
    cmd := exec.CommandContext(ctx, f.bin, args...)
    // capture stderr for diagnostics (do NOT log full stderr at info level)
    if err := cmd.Run(); err != nil {
        return nil, err
    }
    dur, _ := f.Probe(ctx, src)
    // generate thumbnail at 5s:
    //   ffmpeg -y -ss 5 -i src -frames:v 1 workDir/thumb.jpg
    return &Renditions{ /* fill paths */ DurationSec: dur}, nil
}
```

> The exact ffmpeg flags may need tuning for your test media. Keep them in one
> place; the worker treats the runner as a black box.

### 3.2 — Transcode worker (`internal/worker/transcode_worker.go`)

```go
package worker

import (
    "context"
    "fmt"
    "os"
    "path/filepath"
)

type TranscodeWorker struct {
    jobs       <-chan service.TranscodeJob
    store      service.ObjectStore
    videos     repository.VideoRepository
    coder      transcode.Transcoder
    queue      service.JobQueue // for retry re-enqueue
    cfg        WorkerCfg        // MaxAttempts, WorkDir, concurrency
    notify     NotifyPort       // optional; no-op in Phase 3 (Phase 5 wires FCM)
}

func (w *TranscodeWorker) Run(ctx context.Context) {
    sem := make(chan struct{}, w.cfg.Concurrency) // bound parallel ffmpeg
    for {
        select {
        case <-ctx.Done():
            return
        case job := <-w.jobs:
            sem <- struct{}{}
            go func(j service.TranscodeJob) {
                defer func() { <-sem }()
                if err := w.process(ctx, j); err != nil {
                    w.handleFailure(ctx, j, err)
                }
            }(job)
        }
    }
}

func (w *TranscodeWorker) process(ctx context.Context, j service.TranscodeJob) error {
    // 1. claim: guarded transition is already 'processing' from confirm-upload;
    //    re-assert it's still processing (idempotency vs retries).
    vid := uuid.MustParse(j.VideoID)

    // 2. temp workspace
    work := filepath.Join(w.cfg.WorkDir, j.VideoID)
    if err := os.MkdirAll(work, 0o750); err != nil { return err }
    defer os.RemoveAll(work)

    // 3. download source
    srcPath := filepath.Join(work, "source")
    if err := w.downloadTo(ctx, j.SourceKey, srcPath); err != nil { return err }

    // 4. transcode
    out, err := w.coder.ToHLS(ctx, srcPath, work)
    if err != nil { return fmt.Errorf("ffmpeg: %w", err) }

    // 5. upload outputs under hls/{video_id}/...
    keys, err := w.uploadOutputs(ctx, j.VideoID, out)
    if err != nil { return err }

    // 6. flip to ready (sets all keys + duration) — terminal
    return w.videos.SetReady(ctx, vid, keys, out.DurationSec)
}

func (w *TranscodeWorker) handleFailure(ctx context.Context, j service.TranscodeJob, cause error) {
    vid := uuid.MustParse(j.VideoID)
    if j.AttemptCount+1 < w.cfg.MaxAttempts {
        if ok, _ := w.videos.UpdateStatus(ctx, vid, entity.VideoProcessing, entity.VideoFailed); ok {
            // brief backoff, then retry via failed->processing + re-enqueue
            _, _ = w.videos.UpdateStatus(ctx, vid, entity.VideoFailed, entity.VideoProcessing)
            j.AttemptCount++
            _ = w.queue.EnqueueTranscode(ctx, j)
        }
        return
    }
    // exhausted: mark failed + record error (truncated, no secrets)
    _, _ = w.videos.UpdateStatus(ctx, vid, entity.VideoProcessing, entity.VideoFailed)
    _ = w.videos.SetError(ctx, vid, truncate(cause.Error(), 500))
}
```

- `downloadTo` streams `store.Get(sourceKey)` to a local file.
- `uploadOutputs` walks the HLS dirs and `store.Put`s every `.m3u8`/`.ts` +
  thumbnail under the Skill 1 key scheme; returns the key set for `SetReady`.
- **CDN invalidation** (§8.1 step 7) and **teacher notification** (step 8): expose
  hooks; in Phase 3 the notify hook is a no-op/log (FCM lands in Phase 5). If a
  CloudFront distribution is configured, call `CreateInvalidation` for
  `hls/{video_id}/*`; otherwise skip in dev.

### 3.3 — Repo additions

Add to `VideoRepository` + sqlc (from Skill 1's `SetVideoReady`): a `SetReady`
adapter that fills all rendition keys + duration + `processed_at`, and a `SetError`
that stores `processing_error`. Both already have SQL in Skill 1 — wire the Go
methods here.

### 3.4 — Worker config

`WorkerCfg`: `Concurrency` (`TRANSCODE_CONCURRENCY`, default = min(2, NumCPU)),
`WorkDir` (`TRANSCODE_WORKDIR`, default `os.TempDir()/memere-transcode`),
`MaxAttempts` (from Skill 2). Add to `.env.example`.

### 3.5 — Tests

- With a **fake `Transcoder`** (returns canned `Renditions`) and **mock store +
  repo**: happy path flips `processing → ready` with all keys set; transcoder
  error within attempt budget → re-enqueued (`failed→processing`); error past
  budget → terminal `failed` + error recorded; guarded transitions prevent
  double-ready.
- One **integration test gated by a build tag** (`//go:build ffmpeg`) that runs
  real FFmpeg on a tiny bundled clip and asserts a valid `master.m3u8` + 3
  variant playlists are produced. Skipped in normal CI unless ffmpeg present.

---

## Definition of Done

- [ ] `go build ./...` clean; `golangci-lint run` clean.
- [ ] `go test ./internal/worker/... -run Transcode` passes with the fake
      transcoder (no ffmpeg needed for unit tests).
- [ ] `go test -tags ffmpeg ./internal/infrastructure/transcode/...` produces a
      valid HLS master + 480/720/1080 variant playlists from a sample clip
      (run locally where ffmpeg is installed).
- [ ] Worker concurrency is bounded; one bad job never blocks the pool.
- [ ] All status flips go through the guarded transition; no double-`ready`.
- [ ] Retry is bounded by `MaxAttempts`; final failure records a truncated error
      (no tokens/keys in the message).
- [ ] Outputs land under `hls/{video_id}/...` per the Skill 1 key scheme.

## Verification commands

```bash
go build ./... && golangci-lint run
go test ./internal/worker/... -run Transcode -v
# Local, with ffmpeg installed:
go test -tags ffmpeg ./internal/infrastructure/transcode/... -v
```

## Hand-off to Skill 4

Videos now reach `ready` with an HLS ladder in storage. Skill 4 builds **secure
delivery**: pre-signed streaming URL (master manifest), the offline download-URL
flow (2-hour expiry), and the access-control rules that gate who can stream.

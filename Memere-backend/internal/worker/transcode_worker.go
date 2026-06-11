package worker

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/service"
	"github.com/Bereke1t2/Memere/memere-backend/internal/infrastructure/transcode"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// NotifyPort is the (optional) hook fired when a video finishes transcoding.
// Phase 3 wires the no-op LogNotifier; Phase 5 swaps in FCM/email.
type NotifyPort interface {
	VideoReady(ctx context.Context, videoID, courseID string)
}

// LogNotifier is the default NotifyPort: it just logs readiness.
type LogNotifier struct{}

func (LogNotifier) VideoReady(_ context.Context, videoID, _ string) {
	log.Printf("transcode: video %s ready", videoID)
}

// WorkerCfg tunes the transcode worker. Concurrency<=0 falls back to
// min(2, NumCPU); WorkDir="" falls back to os.TempDir()/memere-transcode.
type WorkerCfg struct {
	Concurrency int
	WorkDir     string
	MaxAttempts int
}

// TranscodeWorker consumes TranscodeJobs, runs the transcoder, uploads the HLS
// outputs, and flips the video to ready (or failed with bounded retry). It
// depends only on ports so the ffmpeg runner or the in-proc queue can be
// swapped without touching this logic.
type TranscodeWorker struct {
	jobs   <-chan service.TranscodeJob
	store  service.ObjectStore
	videos repository.VideoRepository
	coder  transcode.Transcoder
	queue  service.JobQueue
	notify NotifyPort
	cfg    WorkerCfg
}

// NewTranscodeWorker wires the worker, applying defaults for concurrency and
// work directory. A nil notify becomes the logging no-op.
func NewTranscodeWorker(
	jobs <-chan service.TranscodeJob,
	store service.ObjectStore,
	videos repository.VideoRepository,
	coder transcode.Transcoder,
	queue service.JobQueue,
	notify NotifyPort,
	cfg WorkerCfg,
) *TranscodeWorker {
	if cfg.Concurrency <= 0 {
		cfg.Concurrency = min(2, runtime.NumCPU())
	}
	if cfg.WorkDir == "" {
		cfg.WorkDir = filepath.Join(os.TempDir(), "memere-transcode")
	}
	if cfg.MaxAttempts <= 0 {
		cfg.MaxAttempts = 3
	}
	if notify == nil {
		notify = LogNotifier{}
	}
	return &TranscodeWorker{
		jobs: jobs, store: store, videos: videos, coder: coder,
		queue: queue, notify: notify, cfg: cfg,
	}
}

// Run consumes jobs until ctx is cancelled, bounding parallel transcodes to
// cfg.Concurrency. Each job runs in its own goroutine so one slow/bad job never
// blocks the pool; a panic in processing is contained and logged.
func (w *TranscodeWorker) Run(ctx context.Context) {
	log.Printf("transcode worker started: concurrency=%d workdir=%s", w.cfg.Concurrency, w.cfg.WorkDir)
	sem := make(chan struct{}, w.cfg.Concurrency)
	for {
		select {
		case <-ctx.Done():
			log.Println("transcode worker stopped")
			return
		case job := <-w.jobs:
			sem <- struct{}{}
			go func(j service.TranscodeJob) {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("transcode: panic on video %s: %v", j.VideoID, r)
						w.handleFailure(ctx, j, fmt.Errorf("panic: %v", r))
					}
					<-sem
				}()
				if err := w.process(ctx, j); err != nil {
					w.handleFailure(ctx, j, err)
				}
			}(job)
		}
	}
}

// process runs the full transcode for one job. Any returned error is handled by
// handleFailure (retry or terminal failure).
func (w *TranscodeWorker) process(ctx context.Context, j service.TranscodeJob) error {
	vid, err := uuid.Parse(j.VideoID)
	if err != nil {
		return fmt.Errorf("bad video id %q: %w", j.VideoID, err)
	}

	work := filepath.Join(w.cfg.WorkDir, j.VideoID)
	if err := os.MkdirAll(work, 0o750); err != nil {
		return fmt.Errorf("mkdir workspace: %w", err)
	}
	defer os.RemoveAll(work)

	srcPath := filepath.Join(work, "source")
	if err := w.downloadTo(ctx, j.SourceKey, srcPath); err != nil {
		return fmt.Errorf("download source: %w", err)
	}

	out, err := w.coder.ToHLS(ctx, srcPath, work)
	if err != nil {
		return fmt.Errorf("transcode: %w", err)
	}

	keys, err := w.uploadOutputs(ctx, j.VideoID, out)
	if err != nil {
		return fmt.Errorf("upload outputs: %w", err)
	}

	v := &entity.Video{
		ID:              vid,
		HLSMasterKey:    &keys.master,
		Res480pKey:      &keys.v480,
		Res720pKey:      &keys.v720,
		Res1080pKey:     &keys.v1080,
		ThumbnailKey:    &keys.thumb,
		DurationSeconds: out.DurationSec,
	}
	if err := w.videos.SetReady(ctx, v); err != nil {
		// A guarded SetReady that matches no processing row means another worker
		// already finalized this video — not a failure.
		if apperror.IsNotFound(err) {
			return nil
		}
		return fmt.Errorf("set ready: %w", err)
	}

	w.notify.VideoReady(ctx, j.VideoID, j.CourseID)
	// CDN invalidation (spec §8.1 step 7) would fire here for hls/{video_id}/*
	// when a CloudFront distribution is configured; skipped in dev.
	return nil
}

// handleFailure either re-enqueues the job (when attempts remain) or marks the
// video terminally failed with a truncated, secret-free error message.
func (w *TranscodeWorker) handleFailure(ctx context.Context, j service.TranscodeJob, cause error) {
	vid, err := uuid.Parse(j.VideoID)
	if err != nil {
		log.Printf("transcode: cannot fail unparseable video id %q", j.VideoID)
		return
	}

	if j.AttemptCount < w.cfg.MaxAttempts {
		// Attempts remain: bounce processing->failed->processing to record the
		// attempt, then re-enqueue with an incremented count. The guard ensures
		// we only retry a row still owned by this attempt.
		if ok, _ := w.videos.UpdateStatus(ctx, vid, entity.VideoProcessing, entity.VideoFailed); !ok {
			return // already moved on; don't double-handle
		}
		if ok, _ := w.videos.UpdateStatus(ctx, vid, entity.VideoFailed, entity.VideoProcessing); !ok {
			return
		}
		j.AttemptCount++
		if err := w.queue.EnqueueTranscode(ctx, j); err != nil {
			// Could not re-enqueue: leave it failed so the boot reconciler or a
			// manual retry can pick it up rather than stranding it in processing.
			_, _ = w.videos.UpdateStatus(ctx, vid, entity.VideoProcessing, entity.VideoFailed)
			_ = w.videos.SetError(ctx, vid, truncate(err.Error(), errMsgMax))
		}
		return
	}

	// Budget exhausted: terminal failure.
	if ok, _ := w.videos.UpdateStatus(ctx, vid, entity.VideoProcessing, entity.VideoFailed); !ok {
		return
	}
	_ = w.videos.SetError(ctx, vid, truncate(cause.Error(), errMsgMax))
	log.Printf("transcode: video %s failed after %d attempts", j.VideoID, j.AttemptCount+1)
}

// downloadTo streams the source object to a local file.
func (w *TranscodeWorker) downloadTo(ctx context.Context, key, dst string) error {
	rc, err := w.store.Get(ctx, key)
	if err != nil {
		return err
	}
	defer rc.Close()

	f, err := os.Create(dst) //nolint:gosec // dst is built from the controlled WorkDir
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := io.Copy(f, rc); err != nil {
		return err
	}
	return nil
}

// outputKeys is the set of storage keys for a finished transcode.
type outputKeys struct {
	master, v480, v720, v1080, thumb string
}

// uploadOutputs pushes the master playlist, every variant's playlist + segments,
// and the thumbnail to object storage under the Skill 1 key scheme, returning the
// keys SetReady records.
func (w *TranscodeWorker) uploadOutputs(ctx context.Context, videoID string, out *transcode.Renditions) (outputKeys, error) {
	keys := outputKeys{
		master: fmt.Sprintf("hls/%s/master.m3u8", videoID),
		v480:   fmt.Sprintf("hls/%s/480p/playlist.m3u8", videoID),
		v720:   fmt.Sprintf("hls/%s/720p/playlist.m3u8", videoID),
		v1080:  fmt.Sprintf("hls/%s/1080p/playlist.m3u8", videoID),
		thumb:  fmt.Sprintf("thumbnails/%s/thumb.jpg", videoID),
	}

	if err := w.putFile(ctx, out.MasterPath, keys.master); err != nil {
		return keys, err
	}
	for name, dir := range out.VariantDirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			return keys, err
		}
		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			local := filepath.Join(dir, e.Name())
			key := fmt.Sprintf("hls/%s/%s/%s", videoID, name, e.Name())
			if err := w.putFile(ctx, local, key); err != nil {
				return keys, err
			}
		}
	}
	if err := w.putFile(ctx, out.ThumbnailPath, keys.thumb); err != nil {
		return keys, err
	}
	return keys, nil
}

// putFile uploads a single local file under key with a content type inferred
// from its extension.
func (w *TranscodeWorker) putFile(ctx context.Context, localPath, key string) error {
	f, err := os.Open(localPath) //nolint:gosec // localPath is under the controlled WorkDir
	if err != nil {
		return err
	}
	defer f.Close()
	return w.store.Put(ctx, key, contentTypeFor(localPath), f)
}

func contentTypeFor(p string) string {
	switch filepath.Ext(p) {
	case ".m3u8":
		return "application/vnd.apple.mpegurl"
	case ".ts":
		return "video/mp2t"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	default:
		return "application/octet-stream"
	}
}

// errMsgMax bounds the stored processing_error so a noisy ffmpeg dump can't bloat
// the row (and to keep any incidental detail short).
const errMsgMax = 500

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}

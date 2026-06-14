package worker

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/service"
	"github.com/Bereke1t2/Memere/memere-backend/internal/infrastructure/transcode"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// --- fake transcoder ---------------------------------------------------------

// fakeCoder writes minimal-but-real HLS-shaped files into workDir so the worker's
// uploadOutputs walk has actual files to push. failWith makes ToHLS error.
type fakeCoder struct {
	failWith error
	duration int
}

func (c *fakeCoder) Probe(context.Context, string) (int, error) { return c.duration, nil }

func (c *fakeCoder) ToHLS(_ context.Context, _, workDir string) (*transcode.Renditions, error) {
	if c.failWith != nil {
		return nil, c.failWith
	}
	master := filepath.Join(workDir, "master.m3u8")
	if err := os.WriteFile(master, []byte("#EXTM3U\n"), 0o600); err != nil {
		return nil, err
	}
	dirs := map[string]string{}
	for _, name := range []string{"480p", "720p", "1080p"} {
		d := filepath.Join(workDir, name)
		_ = os.MkdirAll(d, 0o750)
		_ = os.WriteFile(filepath.Join(d, "playlist.m3u8"), []byte("#EXTM3U\n"), 0o600)
		_ = os.WriteFile(filepath.Join(d, "seg_00000.ts"), []byte("seg"), 0o600)
		dirs[name] = d
	}
	thumb := filepath.Join(workDir, "thumb.jpg")
	_ = os.WriteFile(thumb, []byte("jpg"), 0o600)
	return &transcode.Renditions{MasterPath: master, VariantDirs: dirs, ThumbnailPath: thumb, DurationSec: c.duration}, nil
}

// --- fake object store -------------------------------------------------------

type fakeObjectStore struct {
	service.ObjectStore
	mu  sync.Mutex
	put map[string]bool
}

func newFakeObjectStore() *fakeObjectStore { return &fakeObjectStore{put: map[string]bool{}} }

func (s *fakeObjectStore) Get(_ context.Context, _ string) (io.ReadCloser, error) {
	return io.NopCloser(strings.NewReader("source-bytes")), nil
}

func (s *fakeObjectStore) Put(_ context.Context, key, _ string, body io.Reader) error {
	_, _ = io.Copy(io.Discard, body)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.put[key] = true
	return nil
}

func (s *fakeObjectStore) putKeys() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]string, 0, len(s.put))
	for k := range s.put {
		out = append(out, k)
	}
	return out
}

// --- fake video repo ---------------------------------------------------------

type fakeVideoStore struct {
	mu  sync.Mutex
	row map[uuid.UUID]*entity.Video
}

func newFakeVideoStore() *fakeVideoStore { return &fakeVideoStore{row: map[uuid.UUID]*entity.Video{}} }

func (r *fakeVideoStore) seed(id uuid.UUID, status entity.VideoStatus) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.row[id] = &entity.Video{ID: id, Status: status}
}

func (r *fakeVideoStore) get(id uuid.UUID) *entity.Video {
	r.mu.Lock()
	defer r.mu.Unlock()
	if v, ok := r.row[id]; ok {
		cp := *v
		return &cp
	}
	return nil
}

func (r *fakeVideoStore) Create(context.Context, *entity.Video) error { return nil }
func (r *fakeVideoStore) GetByID(_ context.Context, id uuid.UUID) (*entity.Video, error) {
	if v := r.get(id); v != nil {
		return v, nil
	}
	return nil, apperror.NotFound("video not found", nil)
}
func (r *fakeVideoStore) GetByLessonID(context.Context, uuid.UUID) (*entity.Video, error) {
	return nil, apperror.NotFound("video not found", nil)
}

// SetReady is guarded on the processing status, mirroring the real SQL: a video
// not currently processing yields NotFound so the worker treats it as
// already-finalized (no double-ready).
func (r *fakeVideoStore) SetReady(_ context.Context, v *entity.Video) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cur, ok := r.row[v.ID]
	if !ok || cur.Status != entity.VideoProcessing {
		return apperror.NotFound("video not found", nil)
	}
	cur.Status = entity.VideoReady
	cur.HLSMasterKey = v.HLSMasterKey
	cur.Res480pKey = v.Res480pKey
	cur.Res720pKey = v.Res720pKey
	cur.Res1080pKey = v.Res1080pKey
	cur.ThumbnailKey = v.ThumbnailKey
	cur.DurationSeconds = v.DurationSeconds
	return nil
}

func (r *fakeVideoStore) SetError(_ context.Context, id uuid.UUID, msg string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if v, ok := r.row[id]; ok {
		v.ProcessingError = &msg
	}
	return nil
}

func (r *fakeVideoStore) UpdateStatus(_ context.Context, id uuid.UUID, from, to entity.VideoStatus) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	v, ok := r.row[id]
	if !ok || v.Status != from {
		return false, nil
	}
	v.Status = to
	return true, nil
}

func (r *fakeVideoStore) ListByStatus(context.Context, entity.VideoStatus, int) ([]*entity.Video, error) {
	return nil, nil
}
func (r *fakeVideoStore) SoftDelete(context.Context, uuid.UUID) error { return nil }

// --- fake queue --------------------------------------------------------------

type fakeJobQueue struct {
	mu   sync.Mutex
	jobs []service.TranscodeJob
}

func (q *fakeJobQueue) EnqueueTranscode(_ context.Context, job service.TranscodeJob) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.jobs = append(q.jobs, job)
	return nil
}

func (q *fakeJobQueue) EnqueueNotification(context.Context, service.NotificationJob) error {
	return nil
}

func (q *fakeJobQueue) enqueued() []service.TranscodeJob {
	q.mu.Lock()
	defer q.mu.Unlock()
	return append([]service.TranscodeJob(nil), q.jobs...)
}

// --- helpers -----------------------------------------------------------------

type rig struct {
	w      *TranscodeWorker
	store  *fakeObjectStore
	videos *fakeVideoStore
	queue  *fakeJobQueue
	coder  *fakeCoder
}

func newRig(t *testing.T, coder *fakeCoder, maxAttempts int) *rig {
	t.Helper()
	store := newFakeObjectStore()
	videos := newFakeVideoStore()
	queue := &fakeJobQueue{}
	w := NewTranscodeWorker(nil, store, videos, coder, queue, nil, WorkerCfg{
		Concurrency: 2, WorkDir: t.TempDir(), MaxAttempts: maxAttempts,
	})
	return &rig{w: w, store: store, videos: videos, queue: queue, coder: coder}
}

func job(id uuid.UUID, attempt int) service.TranscodeJob {
	return service.TranscodeJob{VideoID: id.String(), CourseID: uuid.New().String(), SourceKey: "originals/x/source.mp4", AttemptCount: attempt}
}

// --- tests -------------------------------------------------------------------

func TestTranscodeProcess_HappyPathReady(t *testing.T) {
	r := newRig(t, &fakeCoder{duration: 42}, 3)
	id := uuid.New()
	r.videos.seed(id, entity.VideoProcessing)

	if err := r.w.process(context.Background(), job(id, 1)); err != nil {
		t.Fatalf("process: %v", err)
	}
	v := r.videos.get(id)
	if v.Status != entity.VideoReady {
		t.Fatalf("status = %q, want ready", v.Status)
	}
	if v.HLSMasterKey == nil || *v.HLSMasterKey != "hls/"+id.String()+"/master.m3u8" {
		t.Fatalf("master key not set: %+v", v.HLSMasterKey)
	}
	if v.Res480pKey == nil || v.Res720pKey == nil || v.Res1080pKey == nil || v.ThumbnailKey == nil {
		t.Fatalf("rendition/thumbnail keys missing: %+v", v)
	}
	if v.DurationSeconds != 42 {
		t.Fatalf("duration = %d, want 42", v.DurationSeconds)
	}
	// Outputs landed under the Skill 1 key scheme.
	keys := strings.Join(r.store.putKeys(), "\n")
	for _, want := range []string{
		"hls/" + id.String() + "/master.m3u8",
		"hls/" + id.String() + "/480p/playlist.m3u8",
		"hls/" + id.String() + "/480p/seg_00000.ts",
		"thumbnails/" + id.String() + "/thumb.jpg",
	} {
		if !strings.Contains(keys, want) {
			t.Fatalf("missing uploaded key %q in:\n%s", want, keys)
		}
	}
}

func TestTranscodeFailure_RetriesWithinBudget(t *testing.T) {
	r := newRig(t, &fakeCoder{failWith: errors.New("boom")}, 3)
	id := uuid.New()
	r.videos.seed(id, entity.VideoProcessing)

	j := job(id, 1)
	err := r.w.process(context.Background(), j)
	if err == nil {
		t.Fatalf("expected process error")
	}
	r.w.handleFailure(context.Background(), j, err)

	if got := r.videos.get(id).Status; got != entity.VideoProcessing {
		t.Fatalf("status = %q, want processing (bounced for retry)", got)
	}
	en := r.queue.enqueued()
	if len(en) != 1 || en[0].AttemptCount != 2 {
		t.Fatalf("want 1 re-enqueue at attempt 2, got %+v", en)
	}
}

func TestTranscodeFailure_TerminalAfterBudget(t *testing.T) {
	r := newRig(t, &fakeCoder{failWith: errors.New("boom")}, 3)
	id := uuid.New()
	r.videos.seed(id, entity.VideoProcessing)

	j := job(id, 3) // attempt 3 == MaxAttempts -> exhausted
	err := r.w.process(context.Background(), j)
	r.w.handleFailure(context.Background(), j, err)

	v := r.videos.get(id)
	if v.Status != entity.VideoFailed {
		t.Fatalf("status = %q, want failed", v.Status)
	}
	if v.ProcessingError == nil || *v.ProcessingError == "" {
		t.Fatalf("expected a recorded processing error")
	}
	if len(r.queue.enqueued()) != 0 {
		t.Fatalf("exhausted job must not be re-enqueued")
	}
}

func TestTranscodeProcess_GuardedNoDoubleReady(t *testing.T) {
	r := newRig(t, &fakeCoder{duration: 10}, 3)
	id := uuid.New()
	// Already ready (e.g. a duplicate/stale job): SetReady must not overwrite.
	r.videos.seed(id, entity.VideoReady)
	first := "ORIGINAL"
	r.videos.row[id].ThumbnailKey = &first

	if err := r.w.process(context.Background(), job(id, 1)); err != nil {
		t.Fatalf("process should treat already-ready as success, got %v", err)
	}
	v := r.videos.get(id)
	if v.Status != entity.VideoReady || v.ThumbnailKey == nil || *v.ThumbnailKey != "ORIGINAL" {
		t.Fatalf("a finalized video was overwritten: %+v", v)
	}
}

func TestTranscodeError_MessageTruncated(t *testing.T) {
	long := strings.Repeat("x", 5000)
	r := newRig(t, &fakeCoder{failWith: errors.New(long)}, 1)
	id := uuid.New()
	r.videos.seed(id, entity.VideoProcessing)

	j := job(id, 1) // MaxAttempts=1 -> immediate terminal
	err := r.w.process(context.Background(), j)
	r.w.handleFailure(context.Background(), j, err)

	v := r.videos.get(id)
	if v.ProcessingError == nil || len(*v.ProcessingError) > errMsgMax {
		t.Fatalf("error not truncated to <=%d: len=%v", errMsgMax, v.ProcessingError)
	}
}

func TestTranscodeWorker_RunProcessesAndBadJobDoesNotBlock(t *testing.T) {
	// A worker driven through Run: one good job reaches ready even when a bad
	// job (unparseable id) is interleaved, proving the pool isn't blocked.
	store := newFakeObjectStore()
	videos := newFakeVideoStore()
	queue := &fakeJobQueue{}
	jobs := make(chan service.TranscodeJob, 4)
	w := NewTranscodeWorker(jobs, store, videos, &fakeCoder{duration: 5}, queue, nil, WorkerCfg{
		Concurrency: 2, WorkDir: t.TempDir(), MaxAttempts: 1,
	})

	good := uuid.New()
	videos.seed(good, entity.VideoProcessing)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go w.Run(ctx)

	jobs <- service.TranscodeJob{VideoID: "not-a-uuid", SourceKey: "k"} // bad, must not block
	jobs <- job(good, 1)

	deadline := time.After(3 * time.Second)
	for {
		if videos.get(good).Status == entity.VideoReady {
			return
		}
		select {
		case <-deadline:
			t.Fatalf("good job did not reach ready; pool likely blocked")
		case <-time.After(10 * time.Millisecond):
		}
	}
}

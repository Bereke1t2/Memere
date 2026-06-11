package video

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/service"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// fakeVideoRepo is an in-memory VideoRepository. Create enforces one live video
// per lesson (mirroring the partial unique index) so the duplicate→CONFLICT path
// is exercised without a database; UpdateStatus implements the guarded
// transition (flip only when current == from).
type fakeVideoRepo struct {
	mu   sync.Mutex
	byID map[uuid.UUID]*entity.Video
}

func newFakeVideoRepo() *fakeVideoRepo {
	return &fakeVideoRepo{byID: map[uuid.UUID]*entity.Video{}}
}

func (f *fakeVideoRepo) Create(_ context.Context, v *entity.Video) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, ex := range f.byID {
		if ex.DeletedAt == nil && ex.LessonID == v.LessonID {
			return apperror.Conflict("VIDEO_EXISTS", nil)
		}
	}
	if v.ID == uuid.Nil {
		v.ID = uuid.New()
	}
	if v.Status == "" {
		v.Status = entity.VideoPending
	}
	now := time.Now()
	v.CreatedAt, v.UpdatedAt = now, now
	cp := *v
	f.byID[v.ID] = &cp
	return nil
}

func (f *fakeVideoRepo) GetByID(_ context.Context, id uuid.UUID) (*entity.Video, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if v, ok := f.byID[id]; ok && v.DeletedAt == nil {
		cp := *v
		return &cp, nil
	}
	return nil, apperror.NotFound("video not found", nil)
}

func (f *fakeVideoRepo) GetByLessonID(_ context.Context, lessonID uuid.UUID) (*entity.Video, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, v := range f.byID {
		if v.DeletedAt == nil && v.LessonID == lessonID {
			cp := *v
			return &cp, nil
		}
	}
	return nil, apperror.NotFound("video not found", nil)
}

func (f *fakeVideoRepo) SetReady(_ context.Context, v *entity.Video) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	cur, ok := f.byID[v.ID]
	if !ok {
		return apperror.NotFound("video not found", nil)
	}
	cur.Status = entity.VideoReady
	cur.HLSMasterKey = v.HLSMasterKey
	cur.ThumbnailKey = v.ThumbnailKey
	cur.DurationSeconds = v.DurationSeconds
	cp := *cur
	*v = cp
	return nil
}

func (f *fakeVideoRepo) SetError(_ context.Context, id uuid.UUID, msg string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if v, ok := f.byID[id]; ok {
		v.ProcessingError = &msg
	}
	return nil
}

func (f *fakeVideoRepo) UpdateStatus(_ context.Context, id uuid.UUID, from, to entity.VideoStatus) (bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	v, ok := f.byID[id]
	if !ok || v.DeletedAt != nil {
		return false, nil
	}
	if v.Status != from {
		return false, nil
	}
	v.Status = to
	v.UpdatedAt = time.Now()
	return true, nil
}

func (f *fakeVideoRepo) ListByStatus(_ context.Context, status entity.VideoStatus, limit int) ([]*entity.Video, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	var out []*entity.Video
	for _, v := range f.byID {
		if v.DeletedAt == nil && v.Status == status {
			cp := *v
			out = append(out, &cp)
			if limit > 0 && len(out) >= limit {
				break
			}
		}
	}
	return out, nil
}

func (f *fakeVideoRepo) SoftDelete(_ context.Context, id uuid.UUID) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if v, ok := f.byID[id]; ok {
		now := time.Now()
		v.DeletedAt = &now
	}
	return nil
}

// statusOf reads a video's current status straight from the store (test helper).
func (f *fakeVideoRepo) statusOf(id uuid.UUID) entity.VideoStatus {
	f.mu.Lock()
	defer f.mu.Unlock()
	if v, ok := f.byID[id]; ok {
		return v.Status
	}
	return ""
}

// fakeLessonRepo embeds the interface (nil) and overrides only FindByID, the one
// method the video usecase calls. Any other call would panic, surfacing an
// unexpected dependency.
type fakeLessonRepo struct {
	repository.LessonRepository
	byID map[uuid.UUID]*entity.Lesson
}

func (f *fakeLessonRepo) FindByID(_ context.Context, id uuid.UUID) (*entity.Lesson, error) {
	if l, ok := f.byID[id]; ok {
		cp := *l
		return &cp, nil
	}
	return nil, apperror.NotFound("lesson not found", nil)
}

// fakeCourseRepo embeds the interface and overrides only FindByID.
type fakeCourseRepo struct {
	repository.CourseRepository
	byID map[uuid.UUID]*entity.Course
}

func (f *fakeCourseRepo) FindByID(_ context.Context, id uuid.UUID) (*entity.Course, error) {
	if c, ok := f.byID[id]; ok {
		cp := *c
		return &cp, nil
	}
	return nil, apperror.NotFound("course not found", nil)
}

// fakeStore embeds ObjectStore and overrides PresignPut + Exists, the two methods
// the upload usecase uses.
type fakeStore struct {
	service.ObjectStore
	mu         sync.Mutex
	presigned  []string        // keys a PUT URL was minted for
	objects    map[string]bool // keys that "exist" in storage
	presignErr error
}

func newFakeStore() *fakeStore {
	return &fakeStore{objects: map[string]bool{}}
}

func (f *fakeStore) PresignPut(_ context.Context, key, _ string, _ time.Duration) (string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.presignErr != nil {
		return "", f.presignErr
	}
	f.presigned = append(f.presigned, key)
	return "https://minio.test/" + key + "?sig=x", nil
}

func (f *fakeStore) Exists(_ context.Context, key string) (bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.objects[key], nil
}

// markUploaded simulates the client's direct PUT landing the object.
func (f *fakeStore) markUploaded(key string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.objects[key] = true
}

// fakeQueue records enqueued jobs and can be made to fail.
type fakeQueue struct {
	mu      sync.Mutex
	jobs    []service.TranscodeJob
	failErr error
}

func (f *fakeQueue) EnqueueTranscode(_ context.Context, job service.TranscodeJob) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.failErr != nil {
		return f.failErr
	}
	f.jobs = append(f.jobs, job)
	return nil
}

func (f *fakeQueue) count() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.jobs)
}

var errEnqueue = errors.New("queue unavailable")

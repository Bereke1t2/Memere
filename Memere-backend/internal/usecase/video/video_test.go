package video

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

type harness struct {
	svc     *Service
	videos  *fakeVideoRepo
	courses *fakeCourseRepo
	store   *fakeStore
	queue   *fakeQueue

	teacherID uuid.UUID
	courseID  uuid.UUID
	lessonID  uuid.UUID
}

func newHarness() *harness {
	teacherID := uuid.New()
	courseID := uuid.New()
	lessonID := uuid.New()

	courses := &fakeCourseRepo{byID: map[uuid.UUID]*entity.Course{
		courseID: {ID: courseID, TeacherID: teacherID, IsPublished: false},
	}}
	lessons := &fakeLessonRepo{byID: map[uuid.UUID]*entity.Lesson{
		lessonID: {ID: lessonID, CourseID: courseID, Type: entity.LessonTypeVideo},
	}}
	videos := newFakeVideoRepo()
	store := newFakeStore()
	queue := &fakeQueue{}

	cfg := Config{UploadURLTTL: 15 * time.Minute, MaxUploadBytes: 100, MaxAttempts: 3}
	svc := NewService(videos, lessons, courses, store, queue, cfg)

	return &harness{
		svc: svc, videos: videos, courses: courses, store: store, queue: queue,
		teacherID: teacherID, courseID: courseID, lessonID: lessonID,
	}
}

func (h *harness) owner() *Actor { return &Actor{UserID: h.teacherID, Role: entity.RoleTeacher} }
func otherTeacher() *Actor       { return &Actor{UserID: uuid.New(), Role: entity.RoleTeacher} }
func admin() *Actor              { return &Actor{UserID: uuid.New(), Role: entity.RoleAdmin} }
func student() *Actor            { return &Actor{UserID: uuid.New(), Role: entity.RoleStudent} }

func validInput(lessonID uuid.UUID) RequestUploadInput {
	return RequestUploadInput{LessonID: lessonID, FileName: "lecture.mp4", ContentType: "video/mp4", SizeBytes: 50}
}

func mustCode(t *testing.T, err error, code string) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected error with code %s, got nil", code)
	}
	if !apperror.IsCode(err, code) {
		t.Fatalf("err = %v, want code %s", err, code)
	}
}

func TestRequestUpload_HappyPath(t *testing.T) {
	h := newHarness()
	ctx := context.Background()

	res, err := h.svc.RequestUpload(ctx, h.owner(), validInput(h.lessonID))
	if err != nil {
		t.Fatalf("RequestUpload: %v", err)
	}

	wantPrefix := "originals/" + h.courseID.String() + "/" + h.lessonID.String() + "/" + res.VideoID.String() + "/source.mp4"
	if res.SourceKey != wantPrefix {
		t.Fatalf("SourceKey = %q, want %q", res.SourceKey, wantPrefix)
	}
	if !strings.Contains(res.UploadURL, res.SourceKey) {
		t.Fatalf("UploadURL %q does not contain key %q", res.UploadURL, res.SourceKey)
	}
	if res.ExpiresIn != int((15 * time.Minute).Seconds()) {
		t.Fatalf("ExpiresIn = %d", res.ExpiresIn)
	}
	if got := h.videos.statusOf(res.VideoID); got != entity.VideoPending {
		t.Fatalf("status = %q, want pending", got)
	}
	if len(h.store.presigned) != 1 || h.store.presigned[0] != res.SourceKey {
		t.Fatalf("expected one presign for the source key, got %v", h.store.presigned)
	}
}

func TestRequestUpload_AdminMayUpload(t *testing.T) {
	h := newHarness()
	if _, err := h.svc.RequestUpload(context.Background(), admin(), validInput(h.lessonID)); err != nil {
		t.Fatalf("admin RequestUpload: %v", err)
	}
}

func TestRequestUpload_NonOwnerForbidden(t *testing.T) {
	h := newHarness()
	_, err := h.svc.RequestUpload(context.Background(), otherTeacher(), validInput(h.lessonID))
	mustCode(t, err, "FORBIDDEN")
}

func TestRequestUpload_StudentForbidden(t *testing.T) {
	h := newHarness()
	_, err := h.svc.RequestUpload(context.Background(), student(), validInput(h.lessonID))
	mustCode(t, err, "FORBIDDEN")
}

func TestRequestUpload_UnsupportedTypeRejectedBeforeDBWrite(t *testing.T) {
	h := newHarness()
	in := validInput(h.lessonID)
	in.ContentType = "application/pdf"
	_, err := h.svc.RequestUpload(context.Background(), h.owner(), in)
	mustCode(t, err, "BAD_REQUEST")
	if len(h.videos.byID) != 0 {
		t.Fatalf("no video should be created on a rejected type, found %d", len(h.videos.byID))
	}
}

func TestRequestUpload_OversizeRejectedBeforeDBWrite(t *testing.T) {
	h := newHarness()
	in := validInput(h.lessonID)
	in.SizeBytes = h.svc.cfg.MaxUploadBytes + 1
	_, err := h.svc.RequestUpload(context.Background(), h.owner(), in)
	mustCode(t, err, "BAD_REQUEST")
	if len(h.videos.byID) != 0 {
		t.Fatalf("no video should be created on oversize, found %d", len(h.videos.byID))
	}
}

func TestRequestUpload_DuplicatePerLessonConflict(t *testing.T) {
	h := newHarness()
	ctx := context.Background()
	if _, err := h.svc.RequestUpload(ctx, h.owner(), validInput(h.lessonID)); err != nil {
		t.Fatalf("first RequestUpload: %v", err)
	}
	_, err := h.svc.RequestUpload(ctx, h.owner(), validInput(h.lessonID))
	mustCode(t, err, "CONFLICT")
}

func TestConfirmUpload_RejectsWhenObjectMissing(t *testing.T) {
	h := newHarness()
	ctx := context.Background()
	res, err := h.svc.RequestUpload(ctx, h.owner(), validInput(h.lessonID))
	if err != nil {
		t.Fatalf("RequestUpload: %v", err)
	}
	// Client never PUT the file → store.Exists is false.
	_, err = h.svc.ConfirmUpload(ctx, h.owner(), res.VideoID)
	mustCode(t, err, "BAD_REQUEST")
	if got := h.videos.statusOf(res.VideoID); got != entity.VideoPending {
		t.Fatalf("status = %q, want pending (unchanged)", got)
	}
	if h.queue.count() != 0 {
		t.Fatalf("no job should be enqueued, got %d", h.queue.count())
	}
}

func TestConfirmUpload_HappyPathEnqueues(t *testing.T) {
	h := newHarness()
	ctx := context.Background()
	res, err := h.svc.RequestUpload(ctx, h.owner(), validInput(h.lessonID))
	if err != nil {
		t.Fatalf("RequestUpload: %v", err)
	}
	h.store.markUploaded(res.SourceKey)

	view, err := h.svc.ConfirmUpload(ctx, h.owner(), res.VideoID)
	if err != nil {
		t.Fatalf("ConfirmUpload: %v", err)
	}
	if view.Status != entity.VideoProcessing {
		t.Fatalf("status = %q, want processing", view.Status)
	}
	if h.queue.count() != 1 {
		t.Fatalf("want 1 job enqueued, got %d", h.queue.count())
	}
	job := h.queue.jobs[0]
	if job.VideoID != res.VideoID.String() || job.SourceKey != res.SourceKey || job.AttemptCount != 1 {
		t.Fatalf("unexpected job %+v", job)
	}
}

func TestConfirmUpload_DoubleConfirmConflict(t *testing.T) {
	h := newHarness()
	ctx := context.Background()
	res, _ := h.svc.RequestUpload(ctx, h.owner(), validInput(h.lessonID))
	h.store.markUploaded(res.SourceKey)
	if _, err := h.svc.ConfirmUpload(ctx, h.owner(), res.VideoID); err != nil {
		t.Fatalf("first confirm: %v", err)
	}
	_, err := h.svc.ConfirmUpload(ctx, h.owner(), res.VideoID)
	mustCode(t, err, "CONFLICT")
	if h.queue.count() != 1 {
		t.Fatalf("second confirm must not enqueue again, got %d", h.queue.count())
	}
}

func TestConfirmUpload_EnqueueFailureRollsBack(t *testing.T) {
	h := newHarness()
	ctx := context.Background()
	res, _ := h.svc.RequestUpload(ctx, h.owner(), validInput(h.lessonID))
	h.store.markUploaded(res.SourceKey)
	h.queue.failErr = errEnqueue

	_, err := h.svc.ConfirmUpload(ctx, h.owner(), res.VideoID)
	if err == nil {
		t.Fatalf("expected enqueue failure to surface")
	}
	if got := h.videos.statusOf(res.VideoID); got != entity.VideoPending {
		t.Fatalf("status = %q, want rolled back to pending", got)
	}
	if h.queue.count() != 0 {
		t.Fatalf("no job should remain enqueued, got %d", h.queue.count())
	}
}

// seedVideo inserts a video directly in a given status for the harness lesson.
func (h *harness) seedVideo(status entity.VideoStatus, procErr *string) uuid.UUID {
	key := "originals/x/y/z/source.mp4"
	v := &entity.Video{
		ID: uuid.New(), LessonID: h.lessonID, CourseID: h.courseID,
		Status: status, OriginalFileKey: &key, ProcessingError: procErr,
	}
	_ = h.videos.Create(context.Background(), v)
	return v.ID
}

func TestRetryProcessing_FromFailedReenqueues(t *testing.T) {
	h := newHarness()
	id := h.seedVideo(entity.VideoFailed, nil)

	view, err := h.svc.RetryProcessing(context.Background(), h.owner(), id)
	if err != nil {
		t.Fatalf("RetryProcessing: %v", err)
	}
	if view.Status != entity.VideoProcessing {
		t.Fatalf("status = %q, want processing", view.Status)
	}
	if h.queue.count() != 1 {
		t.Fatalf("want 1 job enqueued, got %d", h.queue.count())
	}
}

func TestRetryProcessing_OnlyFromFailed(t *testing.T) {
	h := newHarness()
	id := h.seedVideo(entity.VideoProcessing, nil) // not failed

	_, err := h.svc.RetryProcessing(context.Background(), h.owner(), id)
	mustCode(t, err, "CONFLICT")
	if h.queue.count() != 0 {
		t.Fatalf("no job should be enqueued, got %d", h.queue.count())
	}
}

func TestRetryProcessing_NonOwnerForbidden(t *testing.T) {
	h := newHarness()
	id := h.seedVideo(entity.VideoFailed, nil)
	_, err := h.svc.RetryProcessing(context.Background(), otherTeacher(), id)
	mustCode(t, err, "FORBIDDEN")
}

func TestGetVideoStatus_OwnerSeesErrorDetail(t *testing.T) {
	h := newHarness()
	msg := "ffmpeg exited 1"
	id := h.seedVideo(entity.VideoFailed, &msg)

	view, err := h.svc.GetVideoStatus(context.Background(), h.owner(), id)
	if err != nil {
		t.Fatalf("GetVideoStatus: %v", err)
	}
	if view.Status != entity.VideoFailed || view.ProcessingError == nil || *view.ProcessingError != msg {
		t.Fatalf("owner should see error detail, got %+v", view)
	}
}

func TestGetVideoStatus_StudentOnUnpublishedNotFound(t *testing.T) {
	h := newHarness() // course is unpublished by default
	id := h.seedVideo(entity.VideoProcessing, nil)
	_, err := h.svc.GetVideoStatus(context.Background(), student(), id)
	mustCode(t, err, "RESOURCE_NOT_FOUND")
}

func TestGetVideoStatus_StudentOnPublishedHidesErrorDetail(t *testing.T) {
	h := newHarness()
	h.courses.byID[h.courseID].IsPublished = true
	msg := "internal detail"
	id := h.seedVideo(entity.VideoFailed, &msg)

	view, err := h.svc.GetVideoStatus(context.Background(), student(), id)
	if err != nil {
		t.Fatalf("GetVideoStatus: %v", err)
	}
	if view.Status != entity.VideoFailed {
		t.Fatalf("status = %q, want failed", view.Status)
	}
	if view.ProcessingError != nil {
		t.Fatalf("student must not see processing error, got %q", *view.ProcessingError)
	}
}

func TestRequeueStuck_ReenqueuesProcessing(t *testing.T) {
	h := newHarness()
	// Two processing videos on distinct lessons + one pending (ignored).
	h.videos.byID[uuid.New()] = &entity.Video{ID: uuid.New(), LessonID: uuid.New(), CourseID: h.courseID, Status: entity.VideoProcessing, OriginalFileKey: ptr("a")}
	h.videos.byID[uuid.New()] = &entity.Video{ID: uuid.New(), LessonID: uuid.New(), CourseID: h.courseID, Status: entity.VideoProcessing, OriginalFileKey: ptr("b")}
	h.videos.byID[uuid.New()] = &entity.Video{ID: uuid.New(), LessonID: uuid.New(), CourseID: h.courseID, Status: entity.VideoPending, OriginalFileKey: ptr("c")}

	n, err := h.svc.RequeueStuck(context.Background())
	if err != nil {
		t.Fatalf("RequeueStuck: %v", err)
	}
	if n != 2 || h.queue.count() != 2 {
		t.Fatalf("want 2 re-enqueued, got n=%d queue=%d", n, h.queue.count())
	}
}

func ptr(s string) *string { return &s }

package video

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
)

// readyVideo seeds a ready video with an HLS master + thumbnail for the harness
// course/lesson and returns its id.
func (h *harness) readyVideo(t *testing.T) uuid.UUID {
	t.Helper()
	id := uuid.New()
	master := "hls/" + id.String() + "/master.m3u8"
	thumb := "thumbnails/" + id.String() + "/thumb.jpg"
	v := &entity.Video{
		ID:              id,
		LessonID:        h.lessonID,
		CourseID:        h.courseID,
		Status:          entity.VideoReady,
		HLSMasterKey:    &master,
		ThumbnailKey:    &thumb,
		DurationSeconds: 360,
	}
	if err := h.videos.Create(context.Background(), v); err != nil {
		t.Fatalf("seed ready video: %v", err)
	}
	return id
}

// setCourseFree / setLessonPreview flip access flags on the harness fixtures.
func (h *harness) setCourseFree(free bool)    { h.courses.byID[h.courseID].IsFree = free }
func (h *harness) setLessonPreview(prev bool) { h.lessons.byID[h.lessonID].IsFreePreview = prev }

func TestGetStreamURL_OwnerGetsSignedMaster(t *testing.T) {
	h := newHarness()
	id := h.readyVideo(t)

	res, err := h.svc.GetStreamURL(context.Background(), h.owner(), id)
	if err != nil {
		t.Fatalf("GetStreamURL: %v", err)
	}
	wantKey := "hls/" + id.String() + "/master.m3u8"
	if !strings.Contains(res.MasterURL, wantKey) {
		t.Fatalf("MasterURL %q does not reference master key %q", res.MasterURL, wantKey)
	}
	// No raw storage key may leak in the response struct itself.
	if res.MasterURL == wantKey {
		t.Fatalf("MasterURL is a raw key, not a signed URL")
	}
	if res.ExpiresIn != int((2 * time.Hour).Seconds()) {
		t.Fatalf("ExpiresIn = %d, want %d", res.ExpiresIn, int((2 * time.Hour).Seconds()))
	}
	if res.DurationSec != 360 {
		t.Fatalf("DurationSec = %d, want 360", res.DurationSec)
	}
	if res.ThumbnailURL == "" {
		t.Fatalf("expected a signed thumbnail URL")
	}
	if h.signer.lastTTL() != 2*time.Hour {
		t.Fatalf("signed with ttl %v, want 2h (from config, not hardcoded)", h.signer.lastTTL())
	}
}

func TestGetStreamURL_AdminAllowed(t *testing.T) {
	h := newHarness()
	id := h.readyVideo(t)
	if _, err := h.svc.GetStreamURL(context.Background(), admin(), id); err != nil {
		t.Fatalf("admin GetStreamURL: %v", err)
	}
}

func TestGetStreamURL_StudentBlockedFromPaidContent(t *testing.T) {
	h := newHarness()
	id := h.readyVideo(t) // course not free, lesson not preview
	_, err := h.svc.GetStreamURL(context.Background(), student(), id)
	mustCode(t, err, "FORBIDDEN")
}

func TestGetStreamURL_StudentAllowedForFreeCourse(t *testing.T) {
	h := newHarness()
	h.setCourseFree(true)
	id := h.readyVideo(t)
	if _, err := h.svc.GetStreamURL(context.Background(), student(), id); err != nil {
		t.Fatalf("student should stream a free course: %v", err)
	}
}

func TestGetStreamURL_StudentAllowedForPreviewLesson(t *testing.T) {
	h := newHarness()
	h.setLessonPreview(true)
	id := h.readyVideo(t)
	if _, err := h.svc.GetStreamURL(context.Background(), student(), id); err != nil {
		t.Fatalf("student should stream a free-preview lesson: %v", err)
	}
}

func TestGetStreamURL_NotReady(t *testing.T) {
	h := newHarness()
	// Seed a video stuck in processing (no master key).
	id := uuid.New()
	src := "originals/x/source.mp4"
	if err := h.videos.Create(context.Background(), &entity.Video{
		ID: id, LessonID: h.lessonID, CourseID: h.courseID,
		Status: entity.VideoProcessing, OriginalFileKey: &src,
	}); err != nil {
		t.Fatalf("seed: %v", err)
	}
	_, err := h.svc.GetStreamURL(context.Background(), h.owner(), id)
	mustCode(t, err, "CONFLICT")
}

func TestGetStreamURL_SignerKeyChangesSignature(t *testing.T) {
	h := newHarness()
	id1 := h.readyVideo(t)
	res1, err := h.svc.GetStreamURL(context.Background(), h.owner(), id1)
	if err != nil {
		t.Fatalf("stream 1: %v", err)
	}
	// A second video (under its own lesson in the same course) has a different
	// master key -> a different signed URL.
	lessonID2 := uuid.New()
	h.lessons.byID[lessonID2] = &entity.Lesson{ID: lessonID2, CourseID: h.courseID, Type: entity.LessonTypeVideo}
	id2 := uuid.New()
	master := "hls/" + id2.String() + "/master.m3u8"
	if err := h.videos.Create(context.Background(), &entity.Video{
		ID: id2, LessonID: lessonID2, CourseID: h.courseID,
		Status: entity.VideoReady, HLSMasterKey: &master,
	}); err != nil {
		t.Fatalf("seed 2: %v", err)
	}
	res2, err := h.svc.GetStreamURL(context.Background(), h.owner(), id2)
	if err != nil {
		t.Fatalf("stream 2: %v", err)
	}
	if res1.MasterURL == res2.MasterURL {
		t.Fatalf("different keys produced identical signed URLs")
	}
}

func TestGetDownloadURL_IssuesSingleUseToken(t *testing.T) {
	h := newHarness()
	h.setCourseFree(true)
	id := h.readyVideo(t)

	res, err := h.svc.GetDownloadURL(context.Background(), student(), id)
	if err != nil {
		t.Fatalf("GetDownloadURL: %v", err)
	}
	if res.Token == "" {
		t.Fatalf("expected a download token")
	}
	if !strings.Contains(res.DownloadURL, id.String()) {
		t.Fatalf("download URL %q does not reference the video", res.DownloadURL)
	}
	if res.ExpiresIn != int((2 * time.Hour).Seconds()) {
		t.Fatalf("ExpiresIn = %d", res.ExpiresIn)
	}
	if h.tokens.issued != 1 {
		t.Fatalf("issued %d tokens, want 1", h.tokens.issued)
	}
}

func TestGetDownloadURL_PaidBlockedForStudent(t *testing.T) {
	h := newHarness()
	id := h.readyVideo(t) // paid course, non-preview lesson
	_, err := h.svc.GetDownloadURL(context.Background(), student(), id)
	mustCode(t, err, "FORBIDDEN")
}

func TestConsumeDownloadToken_SingleUse(t *testing.T) {
	h := newHarness()
	h.setCourseFree(true)
	id := h.readyVideo(t)

	res, err := h.svc.GetDownloadURL(context.Background(), student(), id)
	if err != nil {
		t.Fatalf("GetDownloadURL: %v", err)
	}

	// First redemption succeeds and returns a signed manifest URL.
	url, err := h.svc.ConsumeDownloadToken(context.Background(), res.Token)
	if err != nil {
		t.Fatalf("first Consume: %v", err)
	}
	if !strings.Contains(url, id.String()) {
		t.Fatalf("consumed URL %q does not reference the video", url)
	}

	// Second redemption misses -> NotFound (token already consumed).
	_, err = h.svc.ConsumeDownloadToken(context.Background(), res.Token)
	mustCode(t, err, "RESOURCE_NOT_FOUND")
}

func TestConsumeDownloadToken_UnknownToken(t *testing.T) {
	h := newHarness()
	_, err := h.svc.ConsumeDownloadToken(context.Background(), "deadbeef")
	mustCode(t, err, "RESOURCE_NOT_FOUND")
}

func TestConsumeDownloadToken_RevokedAccessRejected(t *testing.T) {
	h := newHarness()
	h.setCourseFree(true)
	id := h.readyVideo(t)
	res, err := h.svc.GetDownloadURL(context.Background(), student(), id)
	if err != nil {
		t.Fatalf("GetDownloadURL: %v", err)
	}
	// Access is revoked (course no longer free) before redemption.
	h.setCourseFree(false)
	_, err = h.svc.ConsumeDownloadToken(context.Background(), res.Token)
	mustCode(t, err, "FORBIDDEN")
}

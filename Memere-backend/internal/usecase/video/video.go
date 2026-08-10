package video

import (
	"context"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/service"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/access"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// Config carries the upload-flow tunables the usecase needs. main.go maps the
// typed config.VideoConfig / config.StorageConfig onto this so the usecase stays
// free of the config package (dependency rule).
type Config struct {
	UploadURLTTL   time.Duration // pre-signed PUT lifetime
	MaxUploadBytes int64         // reject larger requests before any DB write
	MaxAttempts    int           // transcode retry budget carried on the job
	StreamURLTTL   time.Duration // signed HLS master lifetime (spec §8.3: ~2h)
	DownloadURLTTL time.Duration // signed download lifetime + token TTL (~2h)
}

// LessonAccess decides whether the caller may watch a lesson's video.
// *access.Service satisfies it; routing through it means paid content now
// unlocks for genuinely enrolled/subscribed students, not just owner/admin.
type LessonAccess interface {
	CanAccessLesson(ctx context.Context, actor access.Actor, lesson *entity.Lesson) (bool, error)
}

// Service orchestrates the video upload pipeline over domain ports.
type Service struct {
	videos  repository.VideoRepository
	lessons repository.LessonRepository
	courses repository.CourseRepository
	store   service.ObjectStore
	queue   service.JobQueue
	signer  service.URLSigner
	tokens  service.DownloadTokens
	access  LessonAccess
	cfg     Config
}

// NewService wires the video usecase with its dependencies. signer and tokens
// power the secure-delivery flow (Skill 4); they may be nil in setups that only
// exercise the upload pipeline, in which case the delivery methods return an
// internal error rather than panicking.
func NewService(
	videos repository.VideoRepository,
	lessons repository.LessonRepository,
	courses repository.CourseRepository,
	store service.ObjectStore,
	queue service.JobQueue,
	signer service.URLSigner,
	tokens service.DownloadTokens,
	accessSvc LessonAccess,
	cfg Config,
) *Service {
	return &Service{
		videos:  videos,
		lessons: lessons,
		courses: courses,
		store:   store,
		queue:   queue,
		signer:  signer,
		tokens:  tokens,
		access:  accessSvc,
		cfg:     cfg,
	}
}

// allowedVideoTypes is the whitelist of source container types a teacher may
// upload. Pinning it lets the pre-signed PUT fix Content-Type so a client cannot
// smuggle a foreign or executable payload into the originals/ prefix.
var allowedVideoTypes = map[string]bool{
	"video/mp4":       true,
	"video/quicktime": true,
	"video/webm":      true,
}

func isAllowedVideoType(ct string) bool {
	// Tolerate parameters like "video/mp4; codecs=...".
	if i := strings.IndexByte(ct, ';'); i >= 0 {
		ct = ct[:i]
	}
	return allowedVideoTypes[strings.ToLower(strings.TrimSpace(ct))]
}

// RequestUploadInput is the decoded request for a pre-signed upload URL.
type RequestUploadInput struct {
	LessonID    uuid.UUID
	FileName    string
	ContentType string // must be video/mp4|quicktime|webm
	SizeBytes   int64
}

// RequestUploadResult is returned to the teacher: the new video id, the URL to
// PUT the bytes to, the storage key (opaque to the client), and the URL TTL.
type RequestUploadResult struct {
	VideoID   uuid.UUID
	UploadURL string
	SourceKey string
	ExpiresIn int // seconds
}

// StatusView is the client-safe projection of a video's processing state. It
// deliberately carries NO storage keys — streaming URLs are minted separately
// (Skill 4) and originals are never exposed.
type StatusView struct {
	VideoID         uuid.UUID
	Status          entity.VideoStatus
	DurationSeconds int
	HasThumbnail    bool
	// ProcessingError is populated only for the owning teacher/admin so a
	// student never sees internal failure detail.
	ProcessingError *string
}

// RequestUpload validates the request, creates a pending video row, and returns
// a pre-signed PUT URL the client uploads directly to (the server never proxies
// the bytes, spec §8.1). Ownership is enforced via the lesson's course; size and
// content-type are rejected before any DB write.
func (s *Service) RequestUpload(ctx context.Context, actor *Actor, in RequestUploadInput) (*RequestUploadResult, error) {
	if !actor.isTeacherOrAdmin() {
		return nil, apperror.Forbidden("only teachers may upload videos", nil)
	}
	if !isAllowedVideoType(in.ContentType) {
		return nil, apperror.BadRequest("only video/mp4, video/quicktime or video/webm allowed", nil)
	}
	if in.SizeBytes <= 0 || in.SizeBytes > s.cfg.MaxUploadBytes {
		return nil, apperror.BadRequest("file size is zero or exceeds the maximum upload size", nil)
	}

	lesson, err := s.lessons.FindByID(ctx, in.LessonID)
	if err != nil {
		return nil, err // NotFound mapped in repo
	}
	if err := s.assertCourseOwner(ctx, actor, lesson.CourseID); err != nil {
		return nil, err // Forbidden
	}

	// If an existing video record exists for this lesson (e.g. from an aborted upload or replacement),
	// soft-delete it so a new upload can be requested cleanly without unique constraint conflicts.
	if existing, err := s.videos.GetByLessonID(ctx, in.LessonID); err == nil && existing != nil {
		_ = s.videos.SoftDelete(ctx, existing.ID)
	}

	videoID := uuid.New()
	ext := strings.ToLower(path.Ext(in.FileName))
	key := fmt.Sprintf("originals/%s/%s/%s/source%s", lesson.CourseID, in.LessonID, videoID, ext)

	v := &entity.Video{
		ID:              videoID,
		LessonID:        in.LessonID,
		CourseID:        lesson.CourseID,
		Status:          entity.VideoPending,
		OriginalFileKey: &key,
		FileSizeBytes:   in.SizeBytes,
	}
	if err := s.videos.Create(ctx, v); err != nil {
		return nil, err // Conflict (VIDEO_EXISTS) if the lesson already has one
	}

	url, err := s.store.PresignPut(ctx, key, in.ContentType, s.cfg.UploadURLTTL)
	if err != nil {
		return nil, apperror.Internal(err)
	}
	return &RequestUploadResult{
		VideoID:   videoID,
		UploadURL: url,
		SourceKey: key,
		ExpiresIn: int(s.cfg.UploadURLTTL.Seconds()),
	}, nil
}

// ConfirmUpload is called after the client finishes the direct PUT. It verifies
// the object actually landed in storage, flips the video pending→processing with
// a guarded transition (so a double-confirm is a no-op), and enqueues a
// transcode job. If the enqueue fails the status is rolled back to pending so
// the upload can be confirmed again rather than being stuck in processing.
func (s *Service) ConfirmUpload(ctx context.Context, actor *Actor, videoID uuid.UUID) (*StatusView, error) {
	v, err := s.videos.GetByID(ctx, videoID)
	if err != nil {
		return nil, err
	}
	if err := s.assertCourseOwner(ctx, actor, v.CourseID); err != nil {
		return nil, err
	}
	if v.OriginalFileKey == nil {
		return nil, apperror.BadRequest("no source object for this video", nil)
	}

	exists, err := s.store.Exists(ctx, *v.OriginalFileKey)
	if err != nil {
		return nil, apperror.Internal(err)
	}
	if !exists {
		return nil, apperror.BadRequest("upload not found in storage; client must PUT the file first", nil)
	}

	ok, err := s.videos.UpdateStatus(ctx, videoID, entity.VideoPending, entity.VideoProcessing)
	if err != nil {
		return nil, err
	}
	if !ok {
		// Already confirmed (or otherwise not pending) — not an error to repeat,
		// but report current state so the caller can react.
		return nil, apperror.Conflict("video is not pending; it may already be processing", nil)
	}

	if err := s.enqueue(ctx, v, 1); err != nil {
		// Roll back so the upload isn't stranded in processing with no worker
		// ever picking it up. Best-effort: a failed rollback is logged upstream.
		_, _ = s.videos.UpdateStatus(ctx, videoID, entity.VideoProcessing, entity.VideoPending)
		return nil, apperror.Internal(err)
	}

	v.Status = entity.VideoProcessing
	return s.viewFor(actor, v, true), nil
}

// GetVideoStatus returns the processing status. The owning teacher/admin may
// always read it (with error detail); any other authenticated user may read a
// published course's video status (without error detail) so a student can poll
// readiness. Raw storage keys are never returned.
func (s *Service) GetVideoStatus(ctx context.Context, actor *Actor, videoID uuid.UUID) (*StatusView, error) {
	v, err := s.videos.GetByID(ctx, videoID)
	if err != nil {
		return nil, err
	}
	course, err := s.courses.FindByID(ctx, v.CourseID)
	if err != nil {
		return nil, err
	}
	owner := s.isOwner(actor, course)
	if !owner {
		if actor == nil || !course.IsPublished {
			// Don't reveal existence of unpublished content to non-owners.
			return nil, apperror.NotFound("video not found", nil)
		}
	}
	return s.viewFor(actor, v, owner), nil
}

// RetryProcessing lets an owner/admin re-run a failed transcode: failed→
// processing (guarded) then re-enqueue. The automated retry budget
// (cfg.MaxAttempts) is enforced by the worker as it re-enqueues on failure; a
// manual retry starts a fresh attempt.
func (s *Service) RetryProcessing(ctx context.Context, actor *Actor, videoID uuid.UUID) (*StatusView, error) {
	v, err := s.videos.GetByID(ctx, videoID)
	if err != nil {
		return nil, err
	}
	if err := s.assertCourseOwner(ctx, actor, v.CourseID); err != nil {
		return nil, err
	}

	ok, err := s.videos.UpdateStatus(ctx, videoID, entity.VideoFailed, entity.VideoProcessing)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, apperror.Conflict("retry is only allowed from the failed state", nil)
	}

	if err := s.enqueue(ctx, v, 1); err != nil {
		_, _ = s.videos.UpdateStatus(ctx, videoID, entity.VideoProcessing, entity.VideoFailed)
		return nil, apperror.Internal(err)
	}

	v.Status = entity.VideoProcessing
	return s.viewFor(actor, v, true), nil
}

// RequeueStuck re-enqueues videos left in processing on boot. The in-process
// queue is non-durable, so a restart drops in-flight jobs; this reconciler finds
// the orphaned rows and hands them back to a worker. It is wired into startup
// (Skill 5). Returns the number of jobs re-enqueued.
func (s *Service) RequeueStuck(ctx context.Context) (int, error) {
	const batch = 100
	stuck, err := s.videos.ListByStatus(ctx, entity.VideoProcessing, batch)
	if err != nil {
		return 0, err
	}
	n := 0
	for _, v := range stuck {
		if v.OriginalFileKey == nil {
			continue
		}
		if err := s.enqueue(ctx, v, 1); err != nil {
			return n, err
		}
		n++
	}
	return n, nil
}

// enqueue builds and submits the transcode job for v.
func (s *Service) enqueue(ctx context.Context, v *entity.Video, attempt int) error {
	var src string
	if v.OriginalFileKey != nil {
		src = *v.OriginalFileKey
	}
	return s.queue.EnqueueTranscode(ctx, service.TranscodeJob{
		VideoID:      v.ID.String(),
		CourseID:     v.CourseID.String(),
		LessonID:     v.LessonID.String(),
		SourceKey:    src,
		AttemptCount: attempt,
	})
}

// viewFor projects a video into the client-safe StatusView. ProcessingError is
// included only when the caller owns the content.
func (s *Service) viewFor(_ *Actor, v *entity.Video, owner bool) *StatusView {
	view := &StatusView{
		VideoID:         v.ID,
		Status:          v.Status,
		DurationSeconds: v.DurationSeconds,
		HasThumbnail:    v.ThumbnailKey != nil,
	}
	if owner {
		view.ProcessingError = v.ProcessingError
	}
	return view
}

// assertCourseOwner loads the course and permits an admin always, the owning
// teacher otherwise, and rejects everyone else with FORBIDDEN.
func (s *Service) assertCourseOwner(ctx context.Context, actor *Actor, courseID uuid.UUID) error {
	course, err := s.courses.FindByID(ctx, courseID)
	if err != nil {
		return err
	}
	if s.isOwner(actor, course) {
		return nil
	}
	return apperror.Forbidden("you do not own this course", nil)
}

func (s *Service) isOwner(actor *Actor, course *entity.Course) bool {
	if actor.isAdmin() {
		return true
	}
	return actor != nil && actor.Role == entity.RoleTeacher && course.TeacherID == actor.UserID
}

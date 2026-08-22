package video

import (
	"context"
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/service"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/media"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// UploadVideoDirectInput carries a fully-encoded video streamed through the
// backend (no client-side pre-signed PUT). Body must yield exactly SizeBytes.
type UploadVideoDirectInput struct {
	LessonID    uuid.UUID
	FileName    string
	ContentType string // video/mp4|quicktime|webm
	SizeBytes   int64
	Body        io.Reader
}

// UploadVideoDirect stores a single, already-playable video file (e.g. MP4)
// straight through the backend to the object store and marks it ready WITHOUT
// transcoding. It is the delivery path for stores that serve byte-range playback
// directly — Google Drive via the /media proxy, or S3/MinIO natively — where the
// HLS transcode pipeline is either unavailable (Drive cannot run ffmpeg on the
// object) or simply not wanted.
//
// The object key is recorded in HLSMasterKey so the existing delivery methods
// (GetStreamURL / GetDownloadURL / ConsumeDownloadToken) serve it unchanged: the
// URL signer's dev path delegates to the store's PresignGet, which for Drive
// returns a short-lived signed /media URL the client streams with HTTP Range.
// Ownership, content-type and size are validated exactly as in RequestUpload;
// the bytes are validated to exist because we wrote them.
func (s *Service) UploadVideoDirect(ctx context.Context, actor *Actor, in UploadVideoDirectInput) (*StatusView, error) {
	if actor == nil || !actor.isTeacherOrAdmin() {
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

	// Replace any prior video for this lesson so a re-upload doesn't collide with
	// the one-live-video-per-lesson unique index. Purge its storage first so the
	// old file doesn't linger in the bucket, then soft-delete the row.
	if existing, gerr := s.videos.GetByLessonID(ctx, in.LessonID); gerr == nil && existing != nil {
		media.PurgeVideo(ctx, s.store, existing)
		_ = s.videos.SoftDelete(ctx, existing.ID)
	}

	videoID := uuid.New()
	ext := strings.ToLower(path.Ext(in.FileName))
	if ext == "" {
		ext = ".mp4"
	}
	key := fmt.Sprintf("videos/%s/%s/%s/video%s", lesson.CourseID, in.LessonID, videoID, ext)

	// Stream to the store. Prefer the SizedUploader capability so a large file is
	// never buffered fully in memory; fall back to Put where it is unsupported.
	if su, ok := s.store.(service.SizedUploader); ok {
		if err := su.PutSized(ctx, key, in.ContentType, in.Body, in.SizeBytes); err != nil {
			return nil, apperror.Internal(err)
		}
	} else if err := s.store.Put(ctx, key, in.ContentType, in.Body); err != nil {
		return nil, apperror.Internal(err)
	}

	// Persist the row, then walk the guarded transitions the transcode path uses:
	// pending -> processing -> ready (SetReady is guarded on the processing row).
	v := &entity.Video{
		ID:              videoID,
		LessonID:        in.LessonID,
		CourseID:        lesson.CourseID,
		Status:          entity.VideoPending,
		OriginalFileKey: &key,
		FileSizeBytes:   in.SizeBytes,
	}
	if err := s.videos.Create(ctx, v); err != nil {
		_ = s.store.Delete(ctx, key) // don't leak the uploaded object on a DB failure
		return nil, err
	}
	if ok, err := s.videos.UpdateStatus(ctx, videoID, entity.VideoPending, entity.VideoProcessing); err != nil {
		return nil, err
	} else if !ok {
		return nil, apperror.Conflict("video is not in a state that can be finalized", nil)
	}

	ready := &entity.Video{ID: videoID, HLSMasterKey: &key}
	if err := s.videos.SetReady(ctx, ready); err != nil {
		// Roll the row back so it isn't stranded in processing (a direct upload has
		// no transcode worker that would ever pick it up).
		_, _ = s.videos.UpdateStatus(ctx, videoID, entity.VideoProcessing, entity.VideoFailed)
		return nil, apperror.Internal(err)
	}

	v.Status = entity.VideoReady
	v.HLSMasterKey = &key
	return s.viewFor(actor, v, true), nil
}

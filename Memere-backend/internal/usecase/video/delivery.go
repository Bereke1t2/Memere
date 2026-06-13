package video

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/access"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// StreamResult is the client-safe response for a streaming request. It carries
// only short-lived signed URLs — never a raw storage key. ThumbnailURL is
// omitted when the video has no thumbnail yet.
type StreamResult struct {
	MasterURL    string `json:"master_url"`
	ExpiresIn    int    `json:"expires_in"` // seconds until the signed URL expires
	DurationSec  int    `json:"duration_seconds"`
	ThumbnailURL string `json:"thumbnail_url,omitempty"`
}

// DownloadResult is the client-safe response for a download request. The signed
// URL is replayable within its window, so the single-use Token gates the actual
// redemption: the client redeems it once via ConsumeDownloadToken.
type DownloadResult struct {
	DownloadURL string `json:"download_url"`
	Token       string `json:"token"`
	ExpiresIn   int    `json:"expires_in"`
}

// GetStreamURL issues a signed URL for a ready video's HLS master manifest. The
// video must be ready with a master key, and the caller must pass assertCanWatch
// (owner/admin, or a student with free/preview access). No storage key is ever
// returned. Streaming is not single-use: HLS fetches many segments within the
// window, so no token is minted here.
func (s *Service) GetStreamURL(ctx context.Context, actor *Actor, videoID uuid.UUID) (*StreamResult, error) {
	v, err := s.requireReady(ctx, videoID)
	if err != nil {
		return nil, err
	}
	if err := s.assertCanWatch(ctx, actor, v); err != nil {
		return nil, err
	}
	if s.signer == nil {
		return nil, apperror.Internal(fmt.Errorf("video: URL signer not configured"))
	}

	master, err := s.signer.SignGet(ctx, *v.HLSMasterKey, s.cfg.StreamURLTTL)
	if err != nil {
		return nil, apperror.Internal(err)
	}
	res := &StreamResult{
		MasterURL:   master,
		ExpiresIn:   int(s.cfg.StreamURLTTL.Seconds()),
		DurationSec: v.DurationSeconds,
	}
	if v.ThumbnailKey != nil {
		// A thumbnail failure must not block playback; best-effort.
		if thumb, terr := s.signer.SignGet(ctx, *v.ThumbnailKey, s.cfg.StreamURLTTL); terr == nil {
			res.ThumbnailURL = thumb
		}
	}
	return res, nil
}

// GetDownloadURL mints a single-use download token for a ready video the caller
// may watch, alongside a signed URL for the master manifest. The token (not the
// URL) is the one-shot credential: the client redeems it once via
// ConsumeDownloadToken before downloading. The server only bounds the signed
// window to DownloadURLTTL (~2h); the 30-day offline re-download rule (spec §8.3
// step 7) is enforced client-side.
func (s *Service) GetDownloadURL(ctx context.Context, actor *Actor, videoID uuid.UUID) (*DownloadResult, error) {
	v, err := s.requireReady(ctx, videoID)
	if err != nil {
		return nil, err
	}
	if err := s.assertCanWatch(ctx, actor, v); err != nil {
		return nil, err
	}
	if s.signer == nil || s.tokens == nil {
		return nil, apperror.Internal(fmt.Errorf("video: download not configured"))
	}

	signed, err := s.signer.SignGet(ctx, *v.HLSMasterKey, s.cfg.DownloadURLTTL)
	if err != nil {
		return nil, apperror.Internal(err)
	}
	token, err := newToken()
	if err != nil {
		return nil, apperror.Internal(err)
	}
	if err := s.tokens.Issue(ctx, token, v.ID.String(), actor.UserID.String(), s.cfg.DownloadURLTTL); err != nil {
		return nil, err
	}
	return &DownloadResult{
		DownloadURL: signed,
		Token:       token,
		ExpiresIn:   int(s.cfg.DownloadURLTTL.Seconds()),
	}, nil
}

// ConsumeDownloadToken redeems a single-use download token: it atomically
// consumes the token, re-validates that the bound video is still ready and that
// the bound user still has access, then returns a freshly signed manifest URL. A
// second redemption (or an unknown/expired token) returns NotFound — the token
// is already gone from the store.
func (s *Service) ConsumeDownloadToken(ctx context.Context, token string) (string, error) {
	if s.tokens == nil || s.signer == nil {
		return "", apperror.Internal(fmt.Errorf("video: download not configured"))
	}
	videoIDStr, userIDStr, ok, err := s.tokens.Consume(ctx, token)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", apperror.NotFound("download link is invalid or already used", nil)
	}
	videoID, err := uuid.Parse(videoIDStr)
	if err != nil {
		return "", apperror.NotFound("download link is invalid or already used", nil)
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return "", apperror.NotFound("download link is invalid or already used", nil)
	}

	v, err := s.requireReady(ctx, videoID)
	if err != nil {
		return "", err
	}
	// Re-resolve access for the bound user (role is irrelevant to free/preview
	// gating and owner check uses the id; treat the redeemer as a student so a
	// revoked/paid path can't be replayed). assertCanWatch resolves course/lesson
	// server-side from the video.
	actor := &Actor{UserID: userID, Role: entity.RoleStudent}
	if err := s.assertCanWatch(ctx, actor, v); err != nil {
		return "", err
	}
	signed, err := s.signer.SignGet(ctx, *v.HLSMasterKey, s.cfg.DownloadURLTTL)
	if err != nil {
		return "", apperror.Internal(err)
	}
	return signed, nil
}

// requireReady loads a video and asserts it is ready with an HLS master key,
// returning VIDEO_NOT_READY otherwise so a caller cannot get a URL for content
// that is still processing or failed.
func (s *Service) requireReady(ctx context.Context, videoID uuid.UUID) (*entity.Video, error) {
	v, err := s.videos.GetByID(ctx, videoID)
	if err != nil {
		return nil, err
	}
	if v.Status != entity.VideoReady || v.HLSMasterKey == nil {
		return nil, apperror.Conflict("video is still processing", nil)
	}
	return v, nil
}

// assertCanWatch enforces video access, resolving the lesson server-side (never
// trusting a client-supplied course id — IDOR) and delegating to the shared
// access.Service. Owner/admin, a free course, an active enrollment, an active
// subscription, or a free-preview lesson all pass; anything else is FORBIDDEN.
// Paid content now unlocks for genuinely enrolled/subscribed students, not just
// owner/admin.
func (s *Service) assertCanWatch(ctx context.Context, actor *Actor, v *entity.Video) error {
	lesson, err := s.lessons.FindByID(ctx, v.LessonID)
	if err != nil {
		return err
	}
	var act access.Actor
	if actor != nil {
		act = access.Actor{UserID: actor.UserID, Role: actor.Role}
	}
	ok, err := s.access.CanAccessLesson(ctx, act, lesson)
	if err != nil {
		return err
	}
	if !ok {
		return apperror.Forbidden("this video requires purchasing the course", nil)
	}
	return nil
}

// newToken mints an unguessable 256-bit hex download token.
func newToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

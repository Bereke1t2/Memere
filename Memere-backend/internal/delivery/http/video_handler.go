package http

import (
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/Bereke1t2/Memere/memere-backend/internal/delivery/http/dto"
	"github.com/Bereke1t2/Memere/memere-backend/internal/delivery/middleware"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/video"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// VideoHandler adapts the Phase 3 video usecase to HTTP. Role gating happens in
// middleware; resource ownership and access control stay in the usecase (IDOR
// rule), so handlers only bind/validate, resolve the actor, and map results.
type VideoHandler struct {
	svc *video.Service
}

// NewVideoHandler builds a VideoHandler.
func NewVideoHandler(svc *video.Service) *VideoHandler {
	return &VideoHandler{svc: svc}
}

// videoActor adapts the context's *course.Actor to the video usecase's *Actor.
// Both carry the same authenticated identity (UserID + Role); the conversion
// keeps the usecase packages decoupled. Returns nil for an anonymous request.
func videoActor(c *gin.Context) *video.Actor {
	a, ok := middleware.ActorFromContext(c)
	if !ok || a == nil {
		return nil
	}
	return &video.Actor{UserID: a.UserID, Role: a.Role}
}

// UploadDirect handles POST /lessons/:id/videos/upload (teacher/admin). The
// client streams the full, already-playable video (e.g. MP4) as
// multipart/form-data under the field "file"; the backend stores it straight to
// the object store and marks it ready WITHOUT transcoding. This is the delivery
// path for stores that serve byte-range playback directly — Google Drive via the
// /media proxy, or S3/MinIO natively. Size/type/ownership are validated in the
// usecase. The global body limit is waived for this route (see middleware.BodyLimit);
// the usecase enforces MaxUploadBytes.
func (h *VideoHandler) UploadDirect(c *gin.Context) {
	lessonID, err := parseUUIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		respondError(c, apperror.BadRequest("missing 'file' form field", err))
		return
	}
	defer file.Close()

	// Prefer the multipart part's declared type; fall back to the file extension
	// when the client sent a generic/empty type (common for octet-stream uploads).
	contentType := header.Header.Get("Content-Type")
	if contentType == "" || contentType == "application/octet-stream" {
		if byExt := mime.TypeByExtension(strings.ToLower(filepath.Ext(header.Filename))); byExt != "" {
			if i := strings.IndexByte(byExt, ';'); i >= 0 {
				byExt = byExt[:i]
			}
			contentType = strings.TrimSpace(byExt)
		}
	}

	view, err := h.svc.UploadVideoDirect(c.Request.Context(), videoActor(c), video.UploadVideoDirectInput{
		LessonID:    lessonID,
		FileName:    header.Filename,
		ContentType: contentType,
		SizeBytes:   header.Size,
		Body:        file,
	})
	if err != nil {
		respondError(c, err)
		return
	}
	respondJSON(c, http.StatusCreated, dto.NewVideoStatusResponse(view))
}

// RequestUpload handles POST /lessons/:id/videos/upload-url (teacher/admin):
// creates a pending video and returns a pre-signed PUT URL.
func (h *VideoHandler) RequestUpload(c *gin.Context) {
	lessonID, err := parseUUIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	var req dto.RequestUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, apperror.BadRequest("invalid request body", err))
		return
	}
	out, err := h.svc.RequestUpload(c.Request.Context(), videoActor(c), video.RequestUploadInput{
		LessonID:    lessonID,
		FileName:    req.FileName,
		ContentType: req.ContentType,
		SizeBytes:   req.SizeBytes,
	})
	if err != nil {
		respondError(c, err)
		return
	}
	respondJSON(c, http.StatusCreated, dto.RequestUploadResponse{
		VideoID:   out.VideoID.String(),
		UploadURL: out.UploadURL,
		ExpiresIn: out.ExpiresIn,
	})
}

// Confirm handles POST /videos/:id/confirm (teacher/admin): verifies the object
// landed in storage, flips pending→processing, and enqueues a transcode job.
func (h *VideoHandler) Confirm(c *gin.Context) {
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	view, err := h.svc.ConfirmUpload(c.Request.Context(), videoActor(c), id)
	if err != nil {
		respondError(c, err)
		return
	}
	respondJSON(c, http.StatusAccepted, dto.NewVideoStatusResponse(view))
}

// Status handles GET /videos/:id/status: processing state + metadata, no raw
// keys. The owning teacher/admin additionally sees the processing error.
func (h *VideoHandler) Status(c *gin.Context) {
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	view, err := h.svc.GetVideoStatus(c.Request.Context(), videoActor(c), id)
	if err != nil {
		respondError(c, err)
		return
	}
	respondJSON(c, http.StatusOK, dto.NewVideoStatusResponse(view))
}

// Retry handles POST /videos/:id/retry (teacher/admin): failed→processing and
// re-enqueue.
func (h *VideoHandler) Retry(c *gin.Context) {
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	view, err := h.svc.RetryProcessing(c.Request.Context(), videoActor(c), id)
	if err != nil {
		respondError(c, err)
		return
	}
	respondJSON(c, http.StatusAccepted, dto.NewVideoStatusResponse(view))
}

// Stream handles GET /videos/:id/stream: a short-lived signed HLS master URL for
// a caller allowed to watch (owner/admin, or a student with free/preview access).
func (h *VideoHandler) Stream(c *gin.Context) {
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	res, err := h.svc.GetStreamURL(c.Request.Context(), videoActor(c), id)
	if err != nil {
		respondError(c, err)
		return
	}
	respondJSON(c, http.StatusOK, dto.StreamResponse{
		MasterURL:    res.MasterURL,
		ExpiresIn:    res.ExpiresIn,
		ThumbnailURL: res.ThumbnailURL,
		DurationSec:  res.DurationSec,
	})
}

// DownloadURL handles GET /videos/:id/download-url: a signed URL plus a
// single-use token the client redeems via ConsumeDownload.
func (h *VideoHandler) DownloadURL(c *gin.Context) {
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	res, err := h.svc.GetDownloadURL(c.Request.Context(), videoActor(c), id)
	if err != nil {
		respondError(c, err)
		return
	}
	respondJSON(c, http.StatusOK, dto.DownloadURLResponse{
		DownloadURL: res.DownloadURL,
		Token:       res.Token,
		ExpiresIn:   res.ExpiresIn,
	})
}

// ConsumeDownload handles GET /videos/download/:token: consumes the single-use
// token and 302-redirects to a freshly signed manifest URL. A second use (or an
// unknown/expired token) yields 404. The route is OptionalAuth: the 256-bit
// random token is itself the credential (bound to the redeemer at issue time and
// re-validated on redemption), so it can't be feasibly probed by an anonymous
// caller, and guest downloads must be redeemable without a bearer.
func (h *VideoHandler) ConsumeDownload(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		respondError(c, apperror.BadRequest("missing download token", nil))
		return
	}
	url, err := h.svc.ConsumeDownloadToken(c.Request.Context(), token)
	if err != nil {
		respondError(c, err)
		return
	}
	c.Redirect(http.StatusFound, url)
}

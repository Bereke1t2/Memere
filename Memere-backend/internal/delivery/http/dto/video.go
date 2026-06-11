// Package dto holds the HTTP request/response shapes for the delivery layer.
// Video DTOs deliberately never expose raw storage keys (original/HLS/segment
// object keys) and surface processing_error only to the owning teacher/admin —
// the usecase decides whether to populate it.
package dto

import "github.com/Bereke1t2/Memere/memere-backend/internal/usecase/video"

// RequestUploadRequest is the body for POST /lessons/:id/videos/upload-url.
type RequestUploadRequest struct {
	FileName    string `json:"file_name" binding:"required"`
	ContentType string `json:"content_type" binding:"required"`
	SizeBytes   int64  `json:"size_bytes" binding:"required,gt=0"`
}

// RequestUploadResponse returns the pre-signed PUT URL. The raw object key is
// intentionally omitted — the client uploads to the opaque URL and never sees
// the storage key.
type RequestUploadResponse struct {
	VideoID   string `json:"video_id"`
	UploadURL string `json:"upload_url"`
	ExpiresIn int    `json:"expires_in"`
}

// VideoStatusResponse is the client-safe processing state. Error is populated
// only for the owning teacher/admin (the usecase leaves it nil otherwise).
type VideoStatusResponse struct {
	VideoID         string `json:"video_id"`
	Status          string `json:"status"`
	DurationSeconds int    `json:"duration_seconds"`
	HasThumbnail    bool   `json:"has_thumbnail"`
	Error           string `json:"error,omitempty"`
}

// StreamResponse carries short-lived signed URLs for HLS playback.
type StreamResponse struct {
	MasterURL    string `json:"master_url"`
	ExpiresIn    int    `json:"expires_in"`
	ThumbnailURL string `json:"thumbnail_url,omitempty"`
	DurationSec  int    `json:"duration_seconds"`
}

// DownloadURLResponse carries the signed URL plus the single-use token the
// client redeems via GET /videos/download/:token.
type DownloadURLResponse struct {
	DownloadURL string `json:"download_url"`
	Token       string `json:"token"`
	ExpiresIn   int    `json:"expires_in"`
}

// NewVideoStatusResponse maps a usecase StatusView to the wire shape. The view's
// ProcessingError is non-nil only when the caller owns the content, so the
// student path never receives internal failure detail.
func NewVideoStatusResponse(v *video.StatusView) VideoStatusResponse {
	resp := VideoStatusResponse{
		VideoID:         v.VideoID.String(),
		Status:          string(v.Status),
		DurationSeconds: v.DurationSeconds,
		HasThumbnail:    v.HasThumbnail,
	}
	if v.ProcessingError != nil {
		resp.Error = *v.ProcessingError
	}
	return resp
}

package entity

import (
	"time"

	"github.com/google/uuid"
)

// VideoStatus is the processing state of a lesson video. Values mirror the
// courses.videos CHECK constraint (spec §4.2.4).
type VideoStatus string

const (
	VideoPending    VideoStatus = "pending"
	VideoProcessing VideoStatus = "processing"
	VideoReady      VideoStatus = "ready"
	VideoFailed     VideoStatus = "failed"
)

// Valid reports whether s is one of the known statuses.
func (s VideoStatus) Valid() bool {
	switch s {
	case VideoPending, VideoProcessing, VideoReady, VideoFailed:
		return true
	}
	return false
}

// CanTransitionTo encodes the pipeline state machine the transcode worker
// (Skill 3) enforces: pending -> processing -> ready, with failure reachable
// from the non-terminal states and retry allowed out of failed. ready is
// terminal.
func (s VideoStatus) CanTransitionTo(next VideoStatus) bool {
	switch s {
	case VideoPending:
		return next == VideoProcessing || next == VideoFailed
	case VideoProcessing:
		return next == VideoReady || next == VideoFailed
	case VideoFailed:
		return next == VideoProcessing // allow retry
	default:
		return false // ready is terminal
	}
}

// Video is the domain model for a lesson's video asset. Pure Go, no db/json
// tags (clean-architecture dependency rule). CourseID is denormalized for fast
// authz (video -> course) without a lesson join. The *Key fields hold S3/MinIO
// object keys, never client-visible URLs.
type Video struct {
	ID              uuid.UUID
	LessonID        uuid.UUID
	CourseID        uuid.UUID
	Status          VideoStatus
	OriginalFileKey *string
	HLSMasterKey    *string
	Res480pKey      *string
	Res720pKey      *string
	Res1080pKey     *string
	ThumbnailKey    *string
	DurationSeconds int
	FileSizeBytes   int64
	ProcessingError *string
	ProcessedAt     *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time
}

package service

import "context"

// TranscodeJob is the unit of async work handed to the transcode worker
// (Skill 3): everything the worker needs to fetch the source object, run
// FFmpeg, and write the HLS renditions back. IDs are carried as strings so the
// job marshals cleanly onto a wire queue (SQS/RabbitMQ) when the in-process
// implementation is swapped out in Phase 6.
type TranscodeJob struct {
	VideoID      string `json:"video_id"`
	CourseID     string `json:"course_id"`
	LessonID     string `json:"lesson_id"`
	SourceKey    string `json:"source_key"`
	AttemptCount int    `json:"attempt_count"`
}

// NotificationJob is the unit of async push/email delivery. The in-app row is
// written inline (before enqueue) so the badge count is immediately consistent;
// this job drives the push/email fan-out that must never block the API path.
type NotificationJob struct {
	UserID   string            `json:"user_id"`
	Type     string            `json:"type"`
	Title    string            `json:"title"`
	Body     string            `json:"body"`
	Data     map[string]string `json:"data,omitempty"`
	Channels []Channel         `json:"channels"`
}

// JobQueue is the async-work port. The monolith uses an in-process buffered
// channel (infrastructure/messaging); a durable SQS/RabbitMQ consumer can be
// dropped in behind the same interface without touching usecases.
type JobQueue interface {
	EnqueueTranscode(ctx context.Context, job TranscodeJob) error
	EnqueueNotification(ctx context.Context, job NotificationJob) error
}

// Package messaging holds JobQueue implementations. InProcQueue is the
// monolith's in-process, buffered-channel queue: producers (the upload usecase)
// enqueue transcode jobs and the transcode worker (Skill 3) consumes them in
// the same process. It is the only place coupling to channels lives; usecases
// depend on service.JobQueue, never on this type.
package messaging

import (
	"context"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/service"
)

// InProcQueue is a buffered-channel JobQueue for single-process deployments.
//
// Durability note: jobs live only in the channel and are lost on restart. That
// is acceptable for Phase 3 because the upload usecase's boot reconciler
// (RequeueStuck) re-enqueues videos left in pending/processing on startup.
// Phase 6 replaces this with a durable SQS consumer behind the same port.
type InProcQueue struct {
	transcode    chan service.TranscodeJob
	notification chan service.NotificationJob
}

var _ service.JobQueue = (*InProcQueue)(nil)

// NewInProcQueue builds a queue whose channels buffer up to buffer jobs each.
func NewInProcQueue(buffer int) *InProcQueue {
	if buffer < 0 {
		buffer = 0
	}
	return &InProcQueue{
		transcode:    make(chan service.TranscodeJob, buffer),
		notification: make(chan service.NotificationJob, buffer),
	}
}

// EnqueueTranscode offers the job to the channel, honoring context cancellation
// so a shutting-down caller is never blocked forever on a full buffer.
func (q *InProcQueue) EnqueueTranscode(ctx context.Context, job service.TranscodeJob) error {
	select {
	case q.transcode <- job:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// EnqueueNotification offers a notification delivery job to the channel.
func (q *InProcQueue) EnqueueNotification(ctx context.Context, job service.NotificationJob) error {
	select {
	case q.notification <- job:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Transcode exposes the receive side to the transcode worker.
func (q *InProcQueue) Transcode() <-chan service.TranscodeJob { return q.transcode }

// Notification exposes the receive side to the notification worker.
func (q *InProcQueue) Notification() <-chan service.NotificationJob { return q.notification }

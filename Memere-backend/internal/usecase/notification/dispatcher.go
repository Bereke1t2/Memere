// Package notification implements the Phase 5 notification system. The Dispatcher
// is the concrete service.Notifier: it writes the in-app row immediately (so the
// badge count is instantly consistent) and enqueues a NotificationJob for async
// push/email fan-out. Slow or failing providers never block or fail the API path.
package notification

import (
	"context"
	"encoding/json"
	"log"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/service"
)

// Dispatcher implements service.Notifier. It is the single entry point for
// all notification events across every bounded domain.
type Dispatcher struct {
	notifications repository.NotificationRepository
	queue         service.JobQueue
}

var _ service.Notifier = (*Dispatcher)(nil)

func NewDispatcher(
	notifications repository.NotificationRepository,
	queue service.JobQueue,
) *Dispatcher {
	return &Dispatcher{
		notifications: notifications,
		queue:         queue,
	}
}

// Notify writes the in-app row synchronously (so read/unread counts are
// immediately correct) then enqueues the push/email job. The caller is never
// blocked by provider latency.
func (d *Dispatcher) Notify(ctx context.Context, ev service.NotifyEvent) error {
	userID, err := uuid.Parse(ev.UserID)
	if err != nil {
		// Malformed userID — log and skip silently so the caller's request
		// succeeds regardless.
		log.Printf("notification: invalid user_id %q: %v", ev.UserID, err)
		return nil
	}

	// Write the in-app row for every event (the in-app channel is always on).
	n := &entity.Notification{
		UserID: userID,
		Type:   ev.Type,
		Title:  ev.Title,
		Body:   ev.Body,
		Data:   ev.Data,
	}
	if _, err := d.notifications.Create(ctx, n); err != nil {
		// Log but don't fail — notification loss is preferable to API failure.
		log.Printf("notification: failed to write in-app row: %v", err)
	}

	// Only enqueue push/email channels if the event requests them.
	needsAsync := false
	for _, ch := range ev.Channels {
		if ch == service.ChannelPush || ch == service.ChannelEmail {
			needsAsync = true
			break
		}
	}
	if !needsAsync {
		return nil
	}

	job := service.NotificationJob{
		UserID:   ev.UserID,
		Type:     ev.Type,
		Title:    ev.Title,
		Body:     ev.Body,
		Data:     ev.Data,
		Channels: ev.Channels,
	}
	if err := d.queue.EnqueueNotification(ctx, job); err != nil {
		log.Printf("notification: failed to enqueue job: %v", err)
	}
	return nil
}

// jsonData marshals a map to JSON bytes for logging/debugging purposes only.
func jsonData(m map[string]string) string {
	b, _ := json.Marshal(m)
	return string(b)
}

// Ensure jsonData is used (avoids unused-function lint error).
var _ = jsonData

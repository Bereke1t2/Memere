package worker

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/service"
)

// notificationJobSource exposes the notification channel from the in-process
// queue without coupling the worker to the concrete InProcQueue type.
type notificationJobSource interface {
	Notification() <-chan service.NotificationJob
}

// NotificationWorker consumes NotificationJob from the queue and fans out to
// push/email channels. Provider failure is logged and retried (up to maxRetry);
// it must never propagate to the originating API request.
type NotificationWorker struct {
	source   notificationJobSource
	push     service.PushSender
	email    service.EmailSender
	devices  repository.DeviceTokenRepository
	prefs    repository.PreferenceRepository
	maxRetry int
}

func NewNotificationWorker(
	source notificationJobSource,
	push service.PushSender,
	email service.EmailSender,
	devices repository.DeviceTokenRepository,
	prefs repository.PreferenceRepository,
) *NotificationWorker {
	return &NotificationWorker{
		source:   source,
		push:     push,
		email:    email,
		devices:  devices,
		prefs:    prefs,
		maxRetry: 3,
	}
}

// Run drains the notification channel until ctx is cancelled. One goroutine
// per job; errors are logged, not propagated.
func (w *NotificationWorker) Run(ctx context.Context) {
	for {
		select {
		case job, ok := <-w.source.Notification():
			if !ok {
				return
			}
			go w.process(ctx, job)
		case <-ctx.Done():
			return
		}
	}
}

func (w *NotificationWorker) process(ctx context.Context, job service.NotificationJob) {
	from, err := uuid.Parse(job.UserID)
	if err != nil {
		log.Printf("notification_worker: invalid user_id %q: %v", job.UserID, err)
		return
	}

	prefs, err := w.prefs.Get(ctx, from)
	if err != nil {
		log.Printf("notification_worker: failed to get prefs for %s: %v", job.UserID, err)
		return
	}

	ev := service.NotifyEvent{
		UserID:   job.UserID,
		Type:     job.Type,
		Title:    job.Title,
		Body:     job.Body,
		Data:     job.Data,
		Channels: job.Channels,
	}

	for _, ch := range job.Channels {
		switch ch {
		case service.ChannelPush:
			if !prefs.PushEnabled {
				continue
			}
			tokens, err := w.devices.ListByUser(ctx, from)
			if err != nil {
				log.Printf("notification_worker: list devices: %v", err)
				continue
			}
			if len(tokens) == 0 {
				continue
			}
			fcmTokens := make([]string, len(tokens))
			for i, t := range tokens {
				fcmTokens[i] = t.FCMToken
			}
			w.withRetry(ctx, "push", func() error {
				invalid, err := w.push.Send(ctx, fcmTokens, ev)
				if err != nil {
					return err
				}
				if len(invalid) > 0 {
					_ = w.devices.DeleteInvalid(ctx, invalid)
				}
				return nil
			})

		case service.ChannelEmail:
			if !prefs.EmailEnabled {
				continue
			}
			// Email address must be resolved from the user record; in Phase 5
			// the HTTP delivery layer injects it via job.Data["email"]. If
			// absent, skip gracefully.
			to, ok := job.Data["email"]
			if !ok || to == "" {
				continue
			}
			w.withRetry(ctx, "email", func() error {
				return w.email.Send(ctx, to, job.Title, job.Body)
			})
		}
	}
}

// withRetry calls fn up to w.maxRetry times with exponential backoff.
func (w *NotificationWorker) withRetry(ctx context.Context, label string, fn func() error) {
	delay := 500 * time.Millisecond
	for attempt := 1; attempt <= w.maxRetry; attempt++ {
		if err := fn(); err == nil {
			return
		} else if attempt == w.maxRetry {
			log.Printf("notification_worker: %s failed after %d attempts: %v", label, w.maxRetry, err)
			return
		} else {
			log.Printf("notification_worker: %s attempt %d failed, retrying in %s", label, attempt, delay)
		}
		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return
		}
		delay *= 2
	}
}

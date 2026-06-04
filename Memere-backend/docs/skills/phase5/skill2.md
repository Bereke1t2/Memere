# Phase 5 · Skill 2 — Notification System (FCM, Email, In-App)

> **Prerequisite:** Phase 5 Skill 1 done (progress engine fires notify hooks).
> Read [`docs/skill.md`](../../skill.md) §2.
>
> **Spec references:** `memere_Design_Specification.md` §11 (entire Notifications
> System — §11.1 channels, §11.2 event catalog), §3.3 (message queue), README
> "Push Notifications".

---

## Goal

Build the **notification system** behind a clean `Notifier` port, with three
channels — **Push (FCM)**, **Email (SendGrid)**, and **In-App** (PostgreSQL) — and
wire up **every no-op notify hook** left across Phases 3–5 (video processed,
purchase confirmed, subscription expired, certificate ready, exam graded, etc.).
Delivery is **async** via the job queue so a slow provider never blocks an API
request.

---

## §11.2 event catalog → implement these

| Event | Trigger (already hooked in earlier phases) | Channels |
|---|---|---|
| `welcome_email` | registration (Phase 1) | Email |
| `email_verification` | registration / resend | Email |
| `lesson_published` | teacher publishes lesson | Push + In-App |
| `exam_graded` | attempt graded (Phase 2) | Push + In-App |
| `purchase_confirmed` | payment success (Phase 4) | Email + In-App |
| `video_ready` | transcode done (Phase 3) | Push (to teacher) |
| `certificate_ready` | course 100% (Phase 5 Skill 1) | Push + Email |
| `streak_warning` | no study in 2 days (Phase 5 sweeper) | Push |
| `subscription_expired` | sub sweeper (Phase 4) | Push + Email |
| `announcement` | admin broadcast (Skill 3) | Push + In-App |

---

## Tasks

### 2.1 — Migration (`migrations/0014_notifications`)

```sql
-- migrations/0014_notifications.up.sql
CREATE TABLE IF NOT EXISTS notifications.notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES auth.users(id),
    type TEXT NOT NULL,
    title TEXT NOT NULL,
    body  TEXT NOT NULL,
    data JSONB,                 -- deep-link payload
    read_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS notifications_user_unread_idx
    ON notifications.notifications (user_id, read_at) WHERE read_at IS NULL;

CREATE TABLE IF NOT EXISTS notifications.device_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES auth.users(id),
    fcm_token TEXT NOT NULL,
    platform TEXT NOT NULL,     -- android | ios
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (fcm_token)
);
CREATE TABLE IF NOT EXISTS notifications.preferences (
    user_id UUID PRIMARY KEY REFERENCES auth.users(id),
    push_enabled BOOLEAN NOT NULL DEFAULT true,
    email_enabled BOOLEAN NOT NULL DEFAULT true,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

(Add the `notifications` schema in this migration if Phase 1 didn't create it.)

### 2.2 — The `Notifier` port (`internal/domain/service/notifier.go`)

This is the interface earlier phases referenced as no-op hooks. Define it
concretely now and back the no-ops with the real implementation.

```go
package service

import "context"

type NotifyEvent struct {
    UserID string
    Type   string                 // matches the §11.2 catalog
    Title  string
    Body   string
    Data   map[string]string      // deep-link payload
    Channels []Channel            // Push, Email, InApp
}
type Channel string
const ( ChannelPush Channel="push"; ChannelEmail Channel="email"; ChannelInApp Channel="in_app" )

type Notifier interface {
    Notify(ctx context.Context, ev NotifyEvent) error // enqueues; returns fast
}

// Per-channel senders (infrastructure implements each).
type PushSender  interface{ Send(ctx context.Context, tokens []string, ev NotifyEvent) error }
type EmailSender interface{ Send(ctx context.Context, to, subject, html string) error }
```

Replace the `notifyNoop` placeholders passed into earlier services with the real
`Notifier`. Provide a **typed facade** so call sites stay readable:

```go
// internal/usecase/notification/facade.go
type Hooks struct{ n service.Notifier; users repository.UserRepository }
func (h *Hooks) PurchaseConfirmed(ctx context.Context, p *entity.Payment) {
    _ = h.n.Notify(ctx, service.NotifyEvent{ UserID: p.StudentID.String(),
        Type: "purchase_confirmed", Title: "Payment received",
        Body: "Your enrollment is confirmed.", Channels: []service.Channel{service.ChannelEmail, service.ChannelInApp}})
}
// VideoReady, ExamGraded, CertificateReady, SubscriptionExpired, LessonPublished, ...
```

### 2.3 — Async dispatcher (`internal/usecase/notification/dispatcher.go`)

`Notify` does **not** call providers inline — it:
1. writes the In-App row immediately (so the badge updates), and
2. enqueues a `NotificationJob` on the job queue (Phase 3's `JobQueue`, extended
   with `EnqueueNotification`) for push/email fan-out.

A `NotificationWorker` (`internal/worker/notification_worker.go`) consumes jobs,
checks user `preferences`, resolves FCM device tokens, and calls the channel
senders with **retry + backoff**. Provider failure never affects the API path.

### 2.4 — Channel implementations (`internal/infrastructure/notification/`)

- `fcm_sender.go` — Firebase Cloud Messaging (HTTP v1 or legacy server key from
  `FCM_SERVER_KEY`). Batch tokens; prune invalid tokens (delete `device_tokens`
  rows on `NotRegistered`).
- `sendgrid_sender.go` — SendGrid API (`SENDGRID_API_KEY`); simple HTML templates
  per event type (welcome, receipt, certificate). Keep templates in one place.
- Both must be **swappable** (interfaces) and **stubbable** in tests/dev (a
  `LogSender` that just logs is the dev default when keys are absent).

### 2.5 — In-app + device-token + prefs usecases

- `RegisterDevice(ctx, actor, fcmToken, platform)` / `UnregisterDevice`.
- `ListNotifications(ctx, actor, cursor)` / `MarkRead(ctx, actor, id)` /
  `MarkAllRead` / `UnreadCount`.
- `UpdatePreferences(ctx, actor, push, email)`.
All scoped to the authenticated user (IDOR-safe).

### 2.6 — Wire the real hooks

Replace every `notifyNoop` from Phases 3–4 and the Skill 1 hooks with the
`notification.Hooks` facade. Confirm: video_ready (Phase 3 worker),
purchase_confirmed + subscription_expired (Phase 4), exam_graded (Phase 2 grading
/ sweeper), certificate_ready (Phase 5 Skill 1), lesson_published (course
publish). `welcome_email`/`email_verification` fire from Phase 1 register (wire
here).

### 2.7 — Tests

`Notify` writes the in-app row + enqueues a job (provider not called inline);
worker respects disabled preferences; FCM invalid-token pruning; idempotent-ish
(duplicate in-app rows acceptable but document); dev `LogSender` used when keys
absent; IDOR on notifications list/mark-read blocked.

---

## Definition of Done

- [ ] `make migrate-up` applies `0014`; down reverses.
- [ ] `go build ./...` clean; `golangci-lint run` clean.
- [ ] `go test ./internal/usecase/notification/... ./internal/worker/...` passes.
- [ ] **No `notifyNoop` remains** (`grep -rn 'notifyNoop\|NoopNotifier' internal/`
      is empty); all §11.2 hooks fire the real `Notifier`.
- [ ] Sending is async: a failing/slow provider never blocks or fails the
      originating API request (tested).
- [ ] In-app notifications + unread count + mark-read work and are IDOR-safe.
- [ ] FCM/SendGrid are interfaces with a dev `LogSender` fallback when keys absent.

## Verification commands

```bash
make migrate-up && go build ./... && golangci-lint run
grep -rn 'notifyNoop\|NoopNotifier' internal/ && echo "FAIL: noop remains" || echo OK
go test ./internal/usecase/notification/... ./internal/worker/... -v
```

## Hand-off to Skill 3

Notifications flow on every key event. Skill 3 builds **admin operations**:
user/content management, payment reconciliation, broadcast announcements, and the
platform analytics dashboards.

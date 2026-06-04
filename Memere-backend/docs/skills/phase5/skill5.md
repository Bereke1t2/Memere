# Phase 5 · Skill 5 — HTTP Delivery & Wiring (Progress, Notifications, Admin, Certificates)

> **Prerequisite:** Phase 5 Skills 1–4 done. Read [`docs/skill.md`](../../skill.md)
> §2 and reuse the Phase 1–4 delivery layer.
>
> **Spec references:** `memere_Design_Specification.md` §5.4, §11, README
> "Progress Tracking" / "Admin Features" / "Certificates", §1.2 (DAU/MAU).

---

## Goal

Expose progress, notifications, admin operations, and certificates over the Gin
API; wire the **notification worker + engagement sweeper** into the app; and run
the Phase 5 end-to-end smoke. After this skill the **student-and-admin experience
is feature-complete** for the MVP+growth scope (Phases 1–5).

Extend the existing delivery layer — reuse `apperror` mapping, middleware,
`actorFromContext`, paginated envelope, constructor wiring.

---

## API surface for Phase 5

Base path `/api/v1`.

**Progress** (`progress_handler.go`)
| Method | Path | Auth |
|---|---|---|
| POST | `/lessons/:id/complete` | bearer (access) |
| PUT | `/lessons/:id/video-progress` | bearer (access) |
| GET | `/courses/:id/progress` | bearer (own/owner/admin) |
| GET | `/me/dashboard` | bearer |
| GET | `/me/streak` | bearer |

**Notifications** (`notification_handler.go`)
| Method | Path | Auth |
|---|---|---|
| GET | `/me/notifications` | bearer |
| GET | `/me/notifications/unread-count` | bearer |
| POST | `/me/notifications/:id/read` | bearer |
| POST | `/me/notifications/read-all` | bearer |
| POST | `/me/devices` | bearer (register FCM token) |
| DELETE | `/me/devices/:token` | bearer |
| PUT | `/me/notification-preferences` | bearer |

**Certificates** (`certificate_handler.go`)
| Method | Path | Auth |
|---|---|---|
| POST | `/courses/:id/certificate` | bearer (issue, 100% req) |
| GET | `/me/certificates` | bearer |
| GET | `/certificates/:id/download` | bearer (owner) |
| GET | `/verify/certificates/:serial` | **public** |

**Admin** (`admin_handler.go`) — all `RequireRole(admin)` unless noted
| Method | Path |
|---|---|
| GET | `/admin/users` · `/admin/users/:id` |
| POST | `/admin/users/:id/suspend` · `/admin/users/:id/reactivate` · `/admin/users/:id/role` |
| GET | `/admin/courses` |
| POST | `/admin/courses/:id/unpublish` |
| GET | `/admin/payments` · `/admin/payments/:id` |
| POST | `/admin/payments/reconcile` |
| POST | `/admin/announcements` |
| GET | `/admin/analytics/overview` · `/admin/analytics/revenue` · `/admin/analytics/engagement` |

---

## Tasks

### 5.1 — DTOs

Request/response with json tags per resource. Notable: notifications list uses the
paginated envelope; `DashboardResponse` composes enrolled courses + %s + streak;
admin analytics responses are flat numeric DTOs. **Never** expose audit-log
internals or raw provider payloads to non-admin; `verify` returns only
name/course/date.

### 5.2 — Handlers

Standard pattern (bind → actor → usecase → DTO → status). A couple of notes:

```go
func (h *ProgressHandler) Complete(c *gin.Context) {
    actor, _ := actorFromContext(c)
    lessonID, err := uuid.Parse(c.Param("id"))
    if err != nil { respondError(c, apperror.BadRequest("INVALID_ID","")); return }
    if err := h.uc.MarkLessonComplete(c, actor, lessonID); err != nil { respondError(c, err); return }
    c.Status(http.StatusNoContent)
}

func (h *CertificateHandler) Verify(c *gin.Context) { // public
    v, err := h.uc.VerifyCertificate(c, c.Param("serial"))
    if err != nil { respondError(c, err); return }
    respondJSON(c, http.StatusOK, toPublicCertView(v))
}
```

### 5.3 — Routes & RBAC

Mount under `/api/v1`. Admin group wrapped in `RequireAuth` +
`RequireRole(admin)`. The two **public** routes — `GET /verify/certificates/:serial`
and (from Phase 4) the webhook — sit outside auth. Everything else
`RequireAuth`; ownership in usecases.

### 5.4 — Wire workers + senders into `main.go`

```go
// senders: real if keys present, else LogSender (dev)
push := notif.NewFCMSender(cfg.FCM)        // or LogSender
email := notif.NewSendGridSender(cfg.SendGrid)
notifier := notification.NewDispatcher(notifRepo, queue, clock) // writes in-app + enqueues
hooks := notification.NewHooks(notifier, userRepo)

// re-inject `hooks` into ALL services that had notifyNoop (Phases 2,3,4,5)
notifWorker := worker.NewNotificationWorker(queue.Notifications(), push, email, notifRepo, deviceRepo, prefRepo, cfg)
go notifWorker.Run(ctx)

engSweeper := worker.NewEngagementSweeper(progressRepo, hooks, clock, cfg.EngagementSweep())
go engSweeper.Run(ctx)

certUC := certificate.NewService(certRepo, progressRepo, courseRepo, userRepo, store, signer, pdfRenderer, hooks, cfg)
```

Confirm **all** previously-no-op hooks now reference `hooks`. Every `go ...Run(ctx)`
worker stops on the shared shutdown signal. By now the app starts: HTTP server +
attempt sweeper (P2) + transcode worker (P3) + subscription sweeper (P4) +
notification worker + engagement sweeper (P5).

### 5.5 — End-to-end smoke test (`scripts/smoke_phase5.sh`)

1. student completes all lessons in an enrolled course → `GET /me/dashboard`
   shows 100%; a `certificate_ready` notification appears in
   `GET /me/notifications`; `GET /me/streak` shows current=1.
2. `POST /courses/:id/certificate` → cert issued (idempotent on retry);
   `GET /certificates/:id/download` → signed URL; `GET /verify/certificates/:serial`
   (no auth) → name/course/date.
3. complete a lesson on consecutive simulated days → streak increments; skip a day
   → resets (if the sweeper/clock is testable; otherwise assert streak logic via
   unit tests and smoke the happy path).
4. register an FCM device; trigger an `exam_graded` event (take+submit an exam) →
   an in-app notification row appears; mark-read drops unread-count.
5. as admin: `GET /admin/analytics/overview` → users/revenue/MRR numbers;
   suspend a user → that user can no longer log in; `POST /admin/announcements`
   (segment=all) → recipients get an in-app announcement.
6. as a non-admin hitting any `/admin/*` → `403`.

### 5.6 — Docs

Update README admin/progress/certificate/notification sections + "Getting
Started" (FCM/SendGrid optional in dev). Optionally extend `api/` OpenAPI.

---

## Definition of Done

- [ ] `go build ./...` clean; `golangci-lint run` clean; `go test ./...` passes.
- [ ] `make up && make migrate-up && make run` serves Phase 5 endpoints; all six
      background workers start and stop cleanly with the app.
- [ ] `scripts/smoke_phase5.sh` passes every assertion.
- [ ] Course completion produces a streak bump, a certificate, and a notification.
- [ ] Certificate download is signed + ownership-checked; public verify exposes
      only minimal data.
- [ ] Admin endpoints are admin-gated (403 for others); analytics overview returns
      correct totals/MRR; suspended users cannot authenticate.
- [ ] Every notification hook across Phases 2–5 fires the real notifier (no noop).

## Verification commands

```bash
make up && make migrate-up
make run &
bash scripts/smoke_phase5.sh
go test ./... && golangci-lint run
grep -rn 'notifyNoop\|NoopNotifier\|TODO(phase' internal/ && echo "FAIL: leftovers" || echo "OK"
```

---

## 🎉 Phase 5 complete — what now

The platform is **feature-complete for the MVP + growth scope**: auth, courses,
quizzes, exams, video, payments, subscriptions, progress, streaks, certificates,
notifications, and admin analytics — all on one clean-architecture monolith
honoring every Non-Negotiable.

**To proceed to Phase 6 (final):**
1. Report "Phase 5 complete" with a summary + smoke-test result.
2. Phase 6 is already authored in `docs/skills/phase6/`: **production hardening,
   observability, performance, and the path to the microservice split + k8s**
   (spec §3.1–3.3, §12, §13). It does not add product features — it makes what
   exists production-ready and scalable.

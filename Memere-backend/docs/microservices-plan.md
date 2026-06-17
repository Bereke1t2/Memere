# Memere — Microservices Extraction Plan (Strangler-Fig)

> **Reference:** Design spec §3.1–3.2. Read alongside `docs/scaling.md`.
> **Principle:** Extract a service only when §12.2 triggers fire. The monolith is
> the correct architecture until the metrics say otherwise. The clean-architecture
> boundaries maintained since Phase 1 make each extraction mechanical — not
> architectural — when the time comes.

---

## Why extraction is safe here

Six phases of strict dependency inversion means every domain already lives in its
own `internal/usecase/<domain>/` package that depends only on `internal/domain`
interfaces — never on sibling usecases or infrastructure directly:

```
usecase/notification  →  domain/repository.NotificationRepo (interface)
                      →  domain/service.JobQueue (interface)
                      →  domain/service.PushSender (interface)
                      →  domain/service.EmailSender (interface)
```

To extract `notification` into its own process:
1. Copy `usecase/notification` + its repo implementations into a new module.
2. Replace the in-monolith caller (`hooks.NotifDispatcher`) with an HTTP/gRPC
   client that implements the **same `domain/service.Notifier` interface**.
3. The monolith does not know the difference — its dependency points at an
   interface, not an import path.

This pattern applies to every domain. The extraction order below goes from
lowest coupling (fewest callers) to highest (most callers depend on it).

---

## Pre-work audit — cross-domain DB joins

The following queries span multiple schemas and must become API calls after
the relevant services are extracted. Audit each before extracting that service.

| Join | File | Becomes |
|------|------|---------|
| `enrollments JOIN courses` | `repository/postgres/enrollment.go` | Course Service API call from Enrollment Service |
| `payments JOIN enrollments` | `usecase/payment/webhook.go` | Event: `payment.confirmed` → Enrollment Service saga |
| `progress JOIN lessons JOIN courses` | `repository/postgres/progress.go` | Course Service API call from Progress Service |
| `certificates JOIN courses JOIN users` | `usecase/certificate/` | Course + Auth Service API calls |
| `exam_attempts JOIN exams JOIN users` | `repository/postgres/exam.go` (analytics) | Exam Service exposes analytics read API |

No join crosses more than two schemas today. None block the first two extractions.

---

## Extraction order

### Extraction 1 — Notification Service (port 8086) ← _extract first_

**Coupling:** near-zero. All callers use the `Notifier` interface via `hooks`
(fire-and-forget). No synchronous response needed.

**DB schema ownership:** `notifications`, `device_tokens`, `notification_preferences`
(already isolated in `migrations/0016_notifications.up.sql`).

**Step-by-step:**

```
1. Create repo: github.com/Bereke1t2/Memere/notification-service
2. Copy:
     internal/usecase/notification/   → notification-service/internal/usecase/
     internal/repository/postgres/notification_repo.go  → ...
     internal/repository/postgres/device_token_repo.go  → ...
     internal/repository/postgres/preference_repo.go    → ...
     internal/infrastructure/notification/              → ...
     internal/worker/notification_worker.go             → ...
3. Write a thin HTTP handler: POST /internal/notifications/dispatch
   (body: NotifyPayload — the same struct the hooks currently pass in-proc)
4. Write HTTPNotifier implementing domain/service.Notifier → calls the HTTP endpoint
5. In monolith cmd/api/main.go: swap
     notifDispatcher := notificationuc.NewDispatcher(...)
   for
     notifDispatcher := notificationclient.NewHTTPNotifier("http://notification-svc")
6. Deploy notification-service; update monolith config NOTIFICATION_SVC_URL.
7. Delete notification usecase + repos from monolith (they now live in the service).
```

**Data migration:** none — the DB schema stays in the shared cluster until
DB-per-service is warranted (Tier 3). The service connects to the same
`notifications` schema.

**Rollback:** swap the config back to the in-proc dispatcher and redeploy
monolith. No data migration needed.

---

### Extraction 2 — Video/Transcode Worker (port 8088) ← _already partly split_

**Status:** The worker `Deployment` is already separate from `api`. This extraction
completes the split: give it its own DB schema ownership and HTTP API.

**Coupling:** `video` usecase is called by `lesson` routes (upload-url, confirm,
retry, status, stream) and by the background TranscodeWorker. The worker already
communicates with the monolith only via the `JobQueue` port.

**Step-by-step:**

```
1. Create repo: github.com/Bereke1t2/Memere/media-service
2. Copy:
     internal/usecase/video/
     internal/repository/postgres/video_repo.go
     internal/infrastructure/storage/
     internal/infrastructure/transcode/
     internal/worker/transcode_worker.go
3. Expose HTTP API (subset of current video routes):
     POST /internal/videos/upload-url
     POST /internal/videos/:id/confirm
     GET  /internal/videos/:id/status
     GET  /internal/videos/:id/stream-url
4. Write VideoClient implementing the same interface the monolith's VideoHandler
   already calls (video.Service methods → HTTP calls to media-service).
5. In monolith: wire VideoHandler to VideoClient instead of video.Service.
6. The transcode queue (now SQS per scaling.md §Tier-2) is owned by media-service.
```

**DB schema:** `videos` table moves to media-service's own DB at Tier 3.

---

### Extraction 3 — Payment Service (port 8085)

**Coupling:** moderately coupled. The payment webhook triggers enrollment and
subscription activation — currently done in a single transaction. After extraction,
this becomes a saga (see §Saga pattern below).

**Step-by-step:**

```
1. Extract usecase/payment + usecase/subscription + usecase/coupon + usecase/revenue
   into payment-service.
2. Webhook: payment.confirmed event published to SQS topic.
3. Monolith Course Service (still in monolith at this stage) subscribes and creates
   enrollment + activates subscription.
4. Compensating transaction: if enrollment fails → payment-service issues refund.
5. Expose HTTP API matching §3.2 port 8085 routes.
```

**Schemas owned:** `payments`, `subscriptions`, `coupons`, `revenue` (schema
`payments` from migrations 0006 / 0013 / 0014).

---

### Extraction 4 — Quiz Service (port 8083) & Exam Service (port 8084)

These are extracted together because they share the `access.Service` dependency
(enrollment/subscription check) and the `grading` package.

```
1. Both services call access.Service to verify enrollment.
   After extraction: access.Service becomes an internal HTTP call to Course/Auth
   Service (or a JWT-embedded claim for subscription status).
2. The Redis attempt-state store (AttemptStateRepo) is co-located with quiz/exam —
   it moves to these services.
3. The scoreRankingRepo (Redis sorted set for leaderboards) moves to Exam Service.
```

**Schemas owned:** `quiz_attempts`, `exam_attempts`, `exam_questions` schemas.

---

### Extraction 5 — Course Service (port 8082) & Progress Service (port 8087)

Extracted together because progress is tightly coupled to course structure
(completion % requires knowing lesson count per section per course).

```
1. Course Service owns: courses, sections, lessons, enrollments schemas.
2. Progress Service owns: progress, streaks schemas.
3. Progress Service calls Course Service API for lesson/course metadata.
4. Certificate Service (currently in monolith) moves here or to its own micro-service.
```

---

### Extraction 6 — Auth Service (port 8081) ← _extract last_

**Auth is the hardest** — every other service depends on JWT validation.
Extract last when Auth is a clear bottleneck or needs independent scaling.

```
1. Auth Service exposes JWKS endpoint: GET /.well-known/jwks.json
2. All other services validate JWTs locally using the public key from JWKS
   (no synchronous call to Auth Service per request).
3. The Redis session denylist (JTI blacklist) moves to Auth Service;
   other services call GET /internal/auth/sessions/:jti/valid on the hot path
   OR replicate the denylist to a Redis read replica each service can query locally.
4. Refresh token rotation stays in Auth Service exclusively.
```

---

## Saga pattern — payment→enrollment

The in-transaction fulfillment (`usecase/payment/webhook.go` calls enrollment +
subscription activation in one DB transaction) becomes a choreography saga once
Payment and Course are separate services:

```
Payment Service:
  1. Verify webhook signature
  2. Mark payment COMPLETED in payments DB
  3. Publish event: payment.confirmed { payment_id, user_id, course_id, plan }

Course Service (subscriber):
  4. Receive payment.confirmed
  5. Create enrollment record
  6. Publish: enrollment.created { enrollment_id, user_id, course_id }
  7. On failure: publish enrollment.failed { payment_id, reason }

Payment Service (compensating):
  8. Receive enrollment.failed
  9. Initiate refund via provider
  10. Mark payment REFUNDED
```

The idempotency key on every payment request (Phase 4 non-negotiable) ensures
replayed events are safe at each step.

---

## Cross-cutting concerns after full extraction

| Concern | Approach |
|---------|----------|
| **JWT validation** | Each service validates locally using JWKS public keys cached from Auth Service (5 min TTL). No per-request call to Auth. |
| **Service-to-service auth** | mTLS (cert-manager issues per-service certs) OR signed internal tokens (short-lived JWTs signed with a shared internal key). |
| **Shared libraries** | `pkg/apperror`, `pkg/pagination`, `pkg/jwt`, `pkg/logger` → shared Go module `github.com/Bereke1t2/Memere/shared`. Or copy-vendored into each service (avoids distributed lib versioning). |
| **Distributed tracing** | OpenTelemetry already wired (Phase 6 Skill 1). Each service propagates `traceparent` header. Jaeger/Tempo aggregates cross-service traces. |
| **API Gateway routing** | nginx-ingress routes to each service by path prefix matching §3.2 port table. At Tier 3, replace with Kong or AWS API Gateway for per-service policies. |

---

## What NOT to extract prematurely

The monolith is the right architecture at Tier 0–1. Splitting before the triggers
fire creates:
- Two codebases to maintain with no throughput gain
- Network latency replacing function calls
- Distributed transaction complexity for no scaling benefit

**Decision rule:** extract a service when a §12.2 trigger fires _and_ the service
boundary is clearly beneficial for that specific trigger. Not before.

---

## Target architecture (§3.2 port map)

```
                            ┌─────────────────────┐
  Mobile App (Flutter)  ──▶ │  API Gateway        │  :443 (nginx-ingress + cert-manager)
                            └──────────┬──────────┘
              ┌──────────────┬─────────┼──────────┬──────────────┐
              ▼              ▼         ▼          ▼              ▼
         Auth Svc      Course Svc  Quiz Svc  Exam Svc     Payment Svc
          :8081          :8082      :8083     :8084          :8085
                              ▼                         ▼
                        Progress Svc            Notification Svc
                          :8087                    :8086
                              ▼
                         Media Svc (Video/Transcode)
                           :8088
```

Each service has its own Deployment + HPA in `k8s/deployments/`. The API Gateway
(nginx-ingress today) routes based on path prefix and rewrites to the ClusterIP
service for the owning microservice.

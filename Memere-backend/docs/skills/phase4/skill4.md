# Phase 4 · Skill 4 — Subscriptions, Renewals & Revenue Reporting

> **Prerequisite:** Phase 4 Skills 1–3 done (payments, webhooks, fulfillment,
> enrollment). Read [`docs/skill.md`](../../skill.md) §2.
>
> **Spec references:** `memere_Design_Specification.md` §1.6 (monetization models,
> 70/30 teacher split), §10.1–10.2 (providers/flow), §4.2.7, README "Monetization
> Models".

---

## Goal

Add the **subscription** monetization model (monthly/annual all-access) on top of
the payment flow, an **expiry/renewal sweeper** that deactivates lapsed
subscriptions, and **revenue reporting** (platform + per-teacher earnings with the
70/30 split). Usecases + worker only; HTTP in Skill 5.

---

## Tasks

### 4.1 — Subscription plans (`internal/usecase/subscription/`)

```go
type Plan string
const ( PlanMonthly Plan = "monthly"; PlanAnnual Plan = "annual" )

type PlanSpec struct { Plan Plan; Price decimal.Decimal; Currency string; Period time.Duration }

// Plans are config-driven (price ranges per spec §1.6). Load from config.
func (s *Service) PlanSpec(p Plan) (PlanSpec, error) { ... }
```

- A subscription purchase reuses the **payment flow** from Skill 3 with
  `CourseID = nil` and a `Plan`. On webhook success, Skill 3's fulfillment calls
  `Subscriptions().Activate(payment)` → create/extend a subscription row with
  `current_period_start/end` and `status=active`.
- `Activate` is **idempotent** and **extends** an existing active subscription
  (period end += plan period) rather than creating duplicates.

### 4.2 — Access integration

Update `access.Service.hasActiveSubscription` (stubbed in Skill 2) to check
`subscriptions` for an `active` row with `current_period_end > now`. An active
subscriber gets `FullAccess` to **all** courses — confirm this matches the spec's
"all-access plan" (§1.6). This means subscription holders bypass per-course
enrollment; document that the access engine already short-circuits on it.

### 4.3 — Renewal/expiry sweeper (`internal/worker/subscription_sweeper.go`)

Reuse the worker pattern from Phases 2–3.

```go
func (w *SubscriptionSweeper) tick(ctx context.Context) error {
    // 1. Find active subscriptions with current_period_end < now.
    expired, err := w.subs.ListExpired(ctx, w.clock.Now())
    ...
    for _, sub := range expired {
        // 2. Attempt auto-renew if provider supports recurring AND not canceled:
        //    create a renewal payment; on success Activate extends the period.
        //    (For Phase 4, Chapa/Telebirr may be one-off only -> mark past_due/expired.)
        // 3. Otherwise transition active -> expired (guarded).
        _ = w.subs.UpdateStatusGuarded(ctx, sub.ID, "active", "expired")
        w.notify.SubscriptionExpired(ctx, sub) // no-op hook now; Phase 5 wires it
    }
    return nil
}
```

- Runs on a ticker (`SUBSCRIPTION_SWEEP_INTERVAL`, e.g. hourly). Starts/stops with
  the app (Skill 5 wires it).
- Idempotent + guarded transitions; safe to run repeatedly.
- `CancelSubscription(ctx, actor, subID)` usecase: owner-only; sets `canceled_at`,
  keeps access until `current_period_end`, then the sweeper expires it.

### 4.4 — Revenue reporting (`internal/usecase/revenue/`)

Read-only aggregations over `payments` (completed only):

- `PlatformRevenue(ctx, actor(admin), from, to)` — gross, net, by provider, by
  course/subscription, refunds. Admin-only.
- `TeacherEarnings(ctx, actor, teacherID, from, to)` — sum of completed payments
  for the teacher's courses × **70%** (teacher share; 30% platform fee per §1.6).
  Teacher may view only their own; admin any.
- `CourseSalesStats(ctx, actor, courseID)` — units sold, gross, conversion.
  Owner/admin only.

Keep the split ratio in config (`TEACHER_REVENUE_SHARE=0.70`) so it's not
hardcoded.

```sql
-- name: TeacherGrossEarnings :one
SELECT COALESCE(SUM(p.amount),0)::numeric AS gross
FROM payments.payments p
JOIN courses.courses c ON c.id = p.course_id
WHERE c.teacher_id = $1 AND p.status='completed'
  AND p.paid_at BETWEEN $2 AND $3;
```

### 4.5 — Tests

- Subscription activate idempotent + extends period; access engine grants
  full access to an active subscriber and denies an expired one (fake clock).
- Sweeper expires lapsed subs (guarded, idempotent); canceled sub keeps access
  until period end then expires.
- Revenue: teacher earnings = gross × 0.70; admin-only platform report; teacher
  cannot read another teacher's earnings (403).

---

## Definition of Done

- [ ] `go build ./...` clean; `golangci-lint run` clean.
- [ ] `go test ./internal/usecase/subscription/... ./internal/usecase/revenue/...
      ./internal/worker/...` passes.
- [ ] Active subscription grants all-access via `access.Service`; expiry revokes it
      (fake-clock tests).
- [ ] Subscription sweeper expires lapsed subs idempotently and stops with the app.
- [ ] Teacher earnings apply the configurable 70/30 split; reports are
      role/ownership-scoped.
- [ ] `Activate` never creates duplicate active subscriptions (extends instead).

## Verification commands

```bash
go build ./... && golangci-lint run
go test ./internal/usecase/subscription/... ./internal/usecase/revenue/... ./internal/worker/... -v
grep -rn '0\.7\b\|0\.70\b' internal/usecase/revenue && echo "REVIEW: ratio should be config-driven"
```

## Hand-off to Skill 5

Purchases, subscriptions, and revenue logic are complete. Skill 5 exposes the
payment/enrollment/subscription/revenue API over HTTP, wires the providers +
subscription sweeper into `main.go`, mounts the **webhook endpoints** (raw-body,
signature-verified), and runs the Phase 4 end-to-end smoke test.

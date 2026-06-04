# Phase 5 · Skill 3 — Admin Operations & Platform Analytics

> **Prerequisite:** Phase 5 Skills 1–2 done. Read [`docs/skill.md`](../../skill.md)
> §2.
>
> **Spec references:** `memere_Design_Specification.md` §1.3 (admin persona),
> §1.4 (FA-* admin features), §3.2 (admin responsibilities), §11.2 (announcement),
> README "Teacher & Admin Features", §10 (payment reconciliation).

---

## Goal

Build the **admin/operations** layer: user management, content moderation
(approve/unpublish), payment reconciliation, broadcast announcements, and
platform-wide analytics dashboards. All strictly **admin-gated** (some
teacher-scoped). Usecases only; HTTP in Skill 5.

These are read-heavy aggregations + privileged mutations. The rule: **every admin
action is authorized by role, scoped, and auditable** (write an audit log row for
state-changing admin actions).

---

## Tasks

### 3.1 — Migration (`migrations/0015_admin_audit`)

```sql
-- migrations/0015_admin_audit.up.sql
CREATE TABLE IF NOT EXISTS auth.admin_audit_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    actor_id UUID NOT NULL REFERENCES auth.users(id),
    action TEXT NOT NULL,          -- e.g. user.suspend, course.unpublish, payment.refund
    target_type TEXT NOT NULL,     -- user | course | payment | ...
    target_id UUID,
    details JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS admin_audit_actor_idx ON auth.admin_audit_log (actor_id, created_at);
CREATE INDEX IF NOT EXISTS admin_audit_target_idx ON auth.admin_audit_log (target_type, target_id);
```

### 3.2 — Audit helper (`internal/usecase/admin/audit.go`)

```go
func (s *Service) audit(ctx context.Context, actor Actor, action, targetType string, targetID *uuid.UUID, details map[string]any) {
    _ = s.audit.Insert(ctx, actor.UserID, action, targetType, targetID, details)
}
```

Call it from every state-changing admin usecase. Audit is best-effort (log on
failure) but should never silently swallow.

### 3.3 — User management (`internal/usecase/admin/users.go`)

Admin-only. `ListUsers(filter, cursor)` (by role/email/active), `GetUser`,
`SuspendUser`/`ReactivateUser` (toggle `is_active`; suspended users fail auth —
confirm the Phase 1 login/middleware checks `is_active`, add it if missing),
`ChangeRole(userID, role)` (e.g. promote a student to teacher — **never** allow
self-demotion of the last admin; guard it), `SoftDeleteUser`. Each audited.

### 3.4 — Content moderation (`internal/usecase/admin/content.go`)

`ListAllCourses(filter, cursor)` (including unpublished, any teacher),
`UnpublishCourse(courseID, reason)` (admin override; notify teacher via Skill 2),
`FeatureCourse` (optional), `DeleteCourse` (admin soft-delete any). Reuse Phase 1
course repo; bypass the teacher-ownership restriction **because** the actor is
admin (the access check allows it). Audited.

### 3.5 — Payment reconciliation (`internal/usecase/admin/payments.go`)

`ListPayments(filter, cursor)` (by status/provider/date/student),
`GetPaymentDetail` (incl. webhook events + enrollment linkage),
`ReconcilePending(ctx)` — find `pending` payments older than N minutes and call
`provider.VerifyPayment` to resolve them (reuses Phase 4's guarded fulfillment;
no double-grant). `ManualRefund` delegates to Phase 4's refund (admin). Audited.

### 3.6 — Announcements / broadcast (`internal/usecase/admin/announce.go`)

`Broadcast(ctx, actor(admin), BroadcastInput)` — send an `announcement`
notification (§11.2) to a **segment**: all users, a role, enrolled-in-course,
or active subscribers. Fan-out via the Skill 2 notifier (push + in-app). For large
segments, **page through recipients and enqueue in batches** — never build one
giant in-memory list; `log()` the recipient count. Audited.

### 3.7 — Platform analytics (`internal/usecase/admin/analytics.go`)

Read-only dashboards (admin; some teacher-scoped):

- `Overview(ctx, from, to)` — totals: users (by role), active users (DAU/MAU per
  §1.2 retention goal), courses, enrollments, **MRR** (active subs × plan price),
  gross revenue, completed-payment count, refunds.
- `RevenueBreakdown` — by provider, by course, by plan (composes Phase 4 revenue
  usecase — don't duplicate).
- `EngagementStats` — avg completion %, quiz/exam pass rates (Phase 2 analytics),
  streak distribution (Phase 5 Skill 1).
- `ContentStats` — top courses by enrollment/rating, processing backlog (videos in
  `processing`/`failed` from Phase 3).

Keep these as **bounded** queries (date ranges, LIMITs) with appropriate indexes;
note any expensive aggregation that should later move to a materialized view
(Phase 6).

### 3.8 — Tests

Role gate: non-admin → 403 on every admin usecase; last-admin demotion blocked;
suspend makes a user fail auth; reconcile resolves a pending payment without
double-granting; broadcast batches recipients and calls the notifier per segment;
overview MRR math; every state-changing action writes an audit row.

---

## Definition of Done

- [ ] `make migrate-up` applies `0015`; down reverses.
- [ ] `go build ./...` clean; `golangci-lint run` clean.
- [ ] `go test ./internal/usecase/admin/...` passes (role gates, reconciliation,
      broadcast batching, audit writes).
- [ ] Every admin usecase rejects non-admin actors (403); teacher-scoped reports
      reject other teachers.
- [ ] Every state-changing admin action writes an `admin_audit_log` row.
- [ ] Reconciliation reuses Phase 4's guarded fulfillment (no double-grant).
- [ ] Broadcast pages/batches recipients (no unbounded in-memory list).
- [ ] Suspended/`is_active=false` users cannot authenticate.

## Verification commands

```bash
make migrate-up && go build ./... && golangci-lint run
go test ./internal/usecase/admin/... -v
```

## Hand-off to Skill 4

Admin operations + analytics exist. Skill 4 adds the **inactivity/streak-warning
sweeper** and the **certificate generation** flow (course-completion certificates,
§FR-10) — the last engine pieces before Phase 5's HTTP layer.

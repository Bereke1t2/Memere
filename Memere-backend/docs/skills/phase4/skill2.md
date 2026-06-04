# Phase 4 · Skill 2 — Enrollment Engine & Access Control (replace TODO(phase4))

> **Prerequisite:** Phase 4 Skill 1 done (payment/enrollment/coupon data layer +
> provider port). Read [`docs/skill.md`](../../skill.md) §2.
>
> **Spec references:** `memere_Design_Specification.md` §4.2.7 (enrollments),
> §7.2 (RBAC — students read *enrolled* courses), §10 (monetization), README
> "Course Purchase".

---

## Goal

Build the **enrollment engine** — the authority on "does this student have access
to this course?" — and use it to **replace every `TODO(phase4)` hook** left in
Phases 2–3:

- Phase 2: quiz/exam *taking* was gated on "course published or owner".
- Phase 3: paid video access was gated to owner/admin only.

After this skill, access decisions everywhere go through one shared, tested
`AccessService`. No payment provider calls yet (Skill 3) — enrollment can still be
granted via `free` and (test-only) direct grants.

---

## Tasks

### 2.1 — Access service (`internal/usecase/access/`)

The single source of truth for course access. Everything else depends on this.

```go
package access

type AccessLevel int
const ( NoAccess AccessLevel = iota; PreviewAccess; FullAccess )

type Service struct {
    enroll repository.EnrollmentRepository
    subs   repository.SubscriptionRepository
    course repository.CourseRepository
    clock  Clock
}

// CanAccessCourse is THE function. Owner/admin -> full; active enrollment -> full;
// active subscription -> full; free course -> full; else preview/none.
func (s *Service) CanAccessCourse(ctx context.Context, actor Actor, courseID uuid.UUID) (AccessLevel, error) {
    c, err := s.course.FindByID(ctx, courseID)
    if err != nil { return NoAccess, err }
    if actor.Role == entity.RoleAdmin || (actor.Role==entity.RoleTeacher && c.TeacherID==actor.UserID) {
        return FullAccess, nil
    }
    if c.IsFree { return FullAccess, nil }
    if actor.UserID != uuid.Nil {
        if ok, _ := s.enroll.Exists(ctx, actor.UserID, courseID); ok { return FullAccess, nil }
        if s.hasActiveSubscription(ctx, actor.UserID) { return FullAccess, nil }
    }
    return PreviewAccess, nil // can see preview lessons only
}

// CanAccessLesson layers free_preview on top of course access.
func (s *Service) CanAccessLesson(ctx context.Context, actor Actor, lesson *entity.Lesson) (bool, error) {
    lvl, err := s.CanAccessCourse(ctx, actor, lesson.CourseID)
    if err != nil { return false, err }
    if lvl == FullAccess { return true, nil }
    if lvl == PreviewAccess && lesson.IsFreePreview { return true, nil }
    return false, nil
}
```

Provide `RequireFullAccess(ctx, actor, courseID) error` returning
`apperror.Forbidden("NOT_ENROLLED", ...)` for use as a guard.

### 2.2 — Enrollment usecases (`internal/usecase/enrollment/`)

- `GrantEnrollment(ctx, studentID, courseID, source, expiresAt)` — idempotent
  (`Exists` check first; never double-grant). `source ∈ {purchase, subscription,
  free, coupon}`. Used by the payment fulfillment path (Skill 3) and the free path.
- `EnrollFree(ctx, actor, courseID)` — only if the course `is_free`; grants
  `source=free`.
- `ListMyEnrollments(ctx, actor)`.
- `IsEnrolled(ctx, actor, courseID)` (thin wrapper over `enroll.Exists`).

### 2.3 — Replace the Phase 2 hooks

In the quiz and exam *taking* usecases, replace the
`TODO(phase4): enrollment check` gates with `access.RequireFullAccess(ctx, actor,
quiz.CourseID)` (and the exam's course). Inject the `access.Service` into the quiz
and exam services (constructor change). A student without access → `403
NOT_ENROLLED`; owner/admin/enrolled/free/subscription → allowed.

> Free-preview lessons may contain a quiz; decide whether preview quizzes are
> allowed. **Decision:** quizzes/exams require **FullAccess** (preview is for
> watching sample lessons, not graded attempts). Document it.

### 2.4 — Replace the Phase 3 hooks

In `video.assertCanWatch` (Skill 3/4), replace the owner/admin+free/preview logic
with a call into `access.CanAccessLesson(ctx, actor, lesson)`. Paid content now
unlocks for genuinely enrolled/subscribed students — not just owner/admin. Remove
the `TODO(phase4)` comments.

### 2.5 — Coupon application (read-only here)

A `coupon.Service.Quote(ctx, code, courseID, basePrice) (finalPrice, couponID,
error)` that validates and computes the discount (using `entity.Coupon.Apply`).
**Do not** increment usage here — that happens atomically at payment completion
(Skill 3) so abandoned checkouts don't burn coupon uses.

### 2.6 — Tests

- Access matrix: owner→full; admin→full; free course→full; enrolled→full; active
  sub→full; none→preview; preview lesson visible at preview, paid lesson not.
- Quiz/exam taking blocked for non-enrolled (`NOT_ENROLLED`), allowed for enrolled.
- Video paid lesson blocked for non-enrolled, allowed once enrolled.
- `GrantEnrollment` idempotent (second call no-ops).
- Coupon quote: valid/expired/over-use/scope-mismatch.

---

## Definition of Done

- [ ] `go build ./...` clean; `golangci-lint run` clean.
- [ ] **Zero** `TODO(phase4)` markers remain in the codebase
      (`grep -rn 'TODO(phase4)' internal/` is empty).
- [ ] `go test ./internal/usecase/access/... ./internal/usecase/enrollment/...`
      passes the full access matrix.
- [ ] Quiz/exam taking and paid video access now route through `access.Service`
      (verified by tests that previously could not run).
- [ ] `GrantEnrollment` is idempotent; coupon `Quote` does not mutate usage.

## Verification commands

```bash
go build ./... && golangci-lint run
grep -rn 'TODO(phase4)' internal/ && echo "FAIL: hooks remain" || echo "OK: all hooks replaced"
go test ./internal/usecase/access/... ./internal/usecase/enrollment/... \
        ./internal/usecase/quiz/... ./internal/usecase/exam/... ./internal/usecase/video/... -v
```

## Hand-off to Skill 3

Access control is unified and the carry-over hooks are gone. Skill 3 builds the
**payment flow**: initiate checkout (idempotent), provider redirect, and the
fulfillment path that grants enrollment exactly once.

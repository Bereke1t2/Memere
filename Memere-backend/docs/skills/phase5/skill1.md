# Phase 5 · Skill 1 — Progress Tracking Engine

> **Prerequisite:** Phases 1–4 complete and green. Read
> [`docs/skill.md`](../../skill.md) §2.
>
> **Spec references:** `memere_Design_Specification.md` §4.2.8 (progress table),
> §9.3 (weak areas, trend — feeds here), README "Progress Tracking" & "Progress
> Service", §3.2 (Progress Service responsibilities: completion, streaks,
> analytics aggregation).

---

## Goal

Build the **progress tracking engine**: per-lesson completion, video watch
position, course completion %, and study **streaks** — the data behind the
student dashboard. Usecases + repository only (HTTP in Skill 5). This consumes
the quiz/exam analytics from Phase 2 and the enrollment/access from Phase 4.

---

## Tasks

### 1.1 — Migration (`migrations/0013_progress_enhancements`)

Phase 1's `0007` created `progress.progress`. Add streak + denormalized rollups
(additive; continue after Phase 4's `0012`):

```sql
-- migrations/0013_progress_enhancements.up.sql
ALTER TABLE progress.progress
    ADD COLUMN IF NOT EXISTS video_progress_seconds INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS last_accessed_at TIMESTAMPTZ;

CREATE UNIQUE INDEX IF NOT EXISTS progress_student_lesson_uniq
    ON progress.progress (student_id, lesson_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS progress_student_course_idx
    ON progress.progress (student_id, course_id);

-- streaks: one row per student
CREATE TABLE IF NOT EXISTS progress.study_streaks (
    student_id UUID PRIMARY KEY REFERENCES auth.users(id),
    current_streak INTEGER NOT NULL DEFAULT 0,
    longest_streak INTEGER NOT NULL DEFAULT 0,
    last_study_date DATE,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- per-course completion rollup (denormalized cache; rebuildable)
CREATE TABLE IF NOT EXISTS progress.course_progress (
    student_id UUID NOT NULL REFERENCES auth.users(id),
    course_id  UUID NOT NULL REFERENCES courses.courses(id),
    completed_lessons INTEGER NOT NULL DEFAULT 0,
    total_lessons     INTEGER NOT NULL DEFAULT 0,
    percent_complete  NUMERIC(5,2) NOT NULL DEFAULT 0,
    completed_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (student_id, course_id)
);
```

### 1.2 — Entities (`internal/domain/entity/progress.go`) — no db/json tags

```go
type LessonProgress struct {
    ID uuid.UUID; StudentID uuid.UUID; LessonID uuid.UUID; CourseID uuid.UUID
    IsCompleted bool; CompletedAt *time.Time
    VideoProgressSeconds int; LastAccessedAt *time.Time
    CreatedAt, UpdatedAt time.Time
}
type CourseProgress struct {
    StudentID, CourseID uuid.UUID
    CompletedLessons, TotalLessons int
    PercentComplete decimal.Decimal
    CompletedAt *time.Time
}
type StudyStreak struct {
    StudentID uuid.UUID
    Current, Longest int
    LastStudyDate *time.Time
}
```

### 1.3 — Repository interface + sqlc + impl

`ProgressRepository`: `UpsertLessonProgress`, `GetLessonProgress`,
`ListByCourse(student, course)`, `RecomputeCourseProgress(student, course)`
(counts completed vs total non-deleted lessons → upsert `course_progress`),
`GetCourseProgress`, `GetStreak`, `UpsertStreak`.

```sql
-- name: UpsertLessonProgress :one
INSERT INTO progress.progress (id, student_id, lesson_id, course_id, is_completed,
    completed_at, video_progress_seconds, last_accessed_at)
VALUES (gen_random_uuid(), $1, $2, $3, $4,
        CASE WHEN $4 THEN now() ELSE NULL END, $5, now())
ON CONFLICT (student_id, lesson_id) WHERE deleted_at IS NULL
DO UPDATE SET is_completed = EXCLUDED.is_completed,
    completed_at = COALESCE(progress.progress.completed_at, EXCLUDED.completed_at),
    video_progress_seconds = GREATEST(progress.progress.video_progress_seconds, EXCLUDED.video_progress_seconds),
    last_accessed_at = now(), updated_at = now()
RETURNING *;

-- name: CountCourseCompletion :one
SELECT
  (SELECT COUNT(*) FROM courses.lessons l WHERE l.course_id=$2 AND l.deleted_at IS NULL AND l.is_published) AS total,
  (SELECT COUNT(*) FROM progress.progress p WHERE p.student_id=$1 AND p.course_id=$2 AND p.is_completed AND p.deleted_at IS NULL) AS completed;
```

### 1.4 — Progress usecases (`internal/usecase/progress/`)

- `MarkLessonComplete(ctx, actor, lessonID)` — require access
  (`access.RequireFullAccess` or preview for free lessons); upsert progress
  (`is_completed=true`); `RecomputeCourseProgress`; bump streak (1.5); if course
  hits 100% set `completed_at` and fire `notify.CertificateReady` (no-op hook
  until Skill 2). IDOR-safe (always the authenticated student).
- `UpdateVideoProgress(ctx, actor, lessonID, seconds)` — upsert watch position
  (monotonic via `GREATEST`); marks complete if ≥ ~95% of video duration.
- `GetCourseProgress(ctx, actor, courseID)` — % + per-lesson list (own data, or
  teacher/admin viewing their course's student with authz).
- `GetMyDashboard(ctx, actor)` — enrolled courses with % + current streak +
  recent activity (composes Phase 4 enrollments + this).

### 1.5 — Streak logic (`internal/usecase/progress/streak.go`)

```go
// On any study activity (lesson complete / video progress / attempt submit):
func (s *Service) bumpStreak(ctx context.Context, studentID uuid.UUID) error {
    st, _ := s.repo.GetStreak(ctx, studentID)
    today := s.clock.Now().UTC().Truncate(24*time.Hour)
    switch {
    case st.LastStudyDate == nil: st.Current = 1
    case sameDay(*st.LastStudyDate, today): /* already counted today */ return nil
    case isYesterday(*st.LastStudyDate, today): st.Current++
    default: st.Current = 1 // gap -> reset
    }
    if st.Current > st.Longest { st.Longest = st.Current }
    st.LastStudyDate = &today
    return s.repo.UpsertStreak(ctx, st)
}
```

Expose `bumpStreak` so Phase 2's attempt-submit can call it too (wire in Skill 5).
Use UTC day boundaries; document (Ethiopia is UTC+3 — decide whether to use local
day; **decision:** use a configurable `STREAK_TZ`, default `Africa/Addis_Ababa`).

### 1.6 — "Weak areas" integration

Surface the Phase 2 §9.3 weak-area aggregation in the dashboard: call the
analytics usecase to attach weakest topics per subject. Read-only; no duplication.

### 1.7 — Tests

Mark-complete updates % and streak; idempotent (re-marking same lesson doesn't
double-count); video progress monotonic + auto-complete at threshold; streak:
consecutive days increment, same day no-op, gap resets, longest tracked (fake
clock); dashboard composes enrollments + progress; IDOR blocked.

---

## Definition of Done

- [ ] `make migrate-up` applies `0013`; down reverses; unique
      `(student_id, lesson_id)` index exists.
- [ ] `make sqlc && go build ./...` clean; `golangci-lint run` clean.
- [ ] `go test ./internal/usecase/progress/...` passes (streak + completion +
      idempotency).
- [ ] Completion % matches published-lesson counts; 100% sets `completed_at`.
- [ ] Streak math correct across day boundaries (fake-clock tests), tz-configurable.
- [ ] All progress reads/writes are scoped to the authenticated student (no IDOR).

## Verification commands

```bash
make migrate-up && make sqlc && go build ./... && golangci-lint run
go test ./internal/usecase/progress/... -v
```

## Hand-off to Skill 2

Progress + streaks exist and fire notification *hooks*. Skill 2 builds the
**notification system** (the `Notifier` port, FCM + SendGrid + in-app channels)
and wires up every no-op hook left across Phases 3–5.

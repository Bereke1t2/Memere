# Phase 2 · Skill 4 — Expiry Sweeper & Scoring Analytics

> **Prerequisite:** Phase 2 Skills 1–3 done (both engines grade synchronously;
> lazy expiry on access works). Read [`docs/skill.md`](../../skill.md) §2.
>
> **Spec references:** `memere_Design_Specification.md` §9.3 (scoring &
> analytics — raw score, percentage, subject breakdown, weak areas, percentile
> rank via Redis sorted sets, trend analysis), §9.2 (EXPIRED→GRADED on timer).

---

## Goal

Two things the synchronous path can't do:

1. A **background sweeper** that finds `in_progress` attempts past `expires_at`
   (student abandoned them) and auto-transitions them to `EXPIRED → GRADED` —
   guaranteeing the §9.2 "server timer fires" rule even when no client request
   ever arrives.
2. The **scoring & analytics** layer from §9.3: percentile rank (Redis sorted
   sets), weak-area aggregation by topic, and score trend over attempts.

Still no HTTP for the new analytics endpoints' wiring — that's Skill 5. This skill
delivers the engine + the worker.

---

## Tasks

### 4.1 — Background expiry sweeper (`internal/worker/attempt_sweeper.go`)

- A worker that runs on a ticker (configurable interval, e.g. 30–60s; add
  `SWEEPER_INTERVAL` to config). On each tick:
  - `quizAttemptRepo.ListExpiredInProgress(now)` and
    `examAttemptRepo.ListExpiredInProgress(now)` (queries from Skill 1).
  - For each, call the engine's submit/grade path with the **last saved answers**
    (from Redis; fall back to `answers_snapshot` in PG if Redis is gone),
    transitioning `IN_PROGRESS → EXPIRED → GRADED`.
  - Idempotent: if another path already graded it (state no longer
    `in_progress`), skip — no double grading.
- **Concurrency safety:** grading must be safe if the sweeper and a late client
  submit race. Use a guarded transition (conditional update
  `WHERE status='in_progress'`); whoever flips the row first wins, the other
  no-ops. (A Redis lock `lock:attempt:{id}` is an acceptable belt-and-suspenders;
  the conditional UPDATE is the real guard.)
- Lifecycle: started from `cmd/api/main.go` (Skill 5 wires it) in its own
  goroutine; stops on the same shutdown signal as the HTTP server.
- **Structured logging** of how many attempts it swept per tick; never log answer
  contents.

> Design note: this is an in-process worker for the Phase 1/2 monolith. The spec's
> SQS/worker split is a later phase — keep the sweeper behind a small interface so
> it can later be replaced by a queue consumer without touching the engines.

### 4.2 — Analytics: percentile via Redis sorted sets (§9.3)

- `internal/repository/redis/leaderboard.go` (or `score_ranking.go`):
  - On every graded exam attempt, `ZADD exam:{exam_id}:scores <percentage>
    <student_id>` (best score per student — use `ZADD GT` semantics or compare
    before writing so only the student's best counts; decide and document).
  - `PercentileRank(ctx, examID, studentID) (float64, error)` via `ZRANK` /
    `ZCOUNT` — "score vs all attempts for the same exam" (§9.3).
  - Keep the sorted set as a **cache/derived** structure; PG `exam_attempts`
    remains the source of truth (rebuildable).

### 4.3 — Analytics usecases (`internal/usecase/analytics/`)

Compute the §9.3 metrics from stored attempts:

- `GetAttemptAnalytics(ctx, actor, attemptID)` — raw score, percentage, subject
  breakdown (from the snapshot), **weak areas** = questions answered wrong grouped
  by `topic`/`subject`, plus the percentile rank from Redis. Ownership-checked.
- `GetStudentTrend(ctx, actor, subject)` — score over consecutive attempts for a
  subject (time-series query on `exam_attempts`/`quiz_attempts`), for the
  requesting student (or a teacher viewing their own course's student with proper
  authz).
- `GetExamStats(ctx, actor, examID)` — teacher/admin view: attempt count, average,
  pass rate, distribution. Ownership/role-checked (teacher only for own course's
  exam).

Keep these **read-only** and derived; they must not mutate attempts.

### 4.4 — Weak-area aggregation

- From a graded attempt's per-question results (captured in Skill 2/3 grading
  output and persisted in `answers_snapshot`), aggregate wrong answers by `topic`
  then `subject`. Return a ranked list of weakest topics for the student.
- This feeds the future Progress service (Phase 5) — keep the shape stable and
  documented.

### 4.5 — Config additions

Add to the typed `Config` (Phase 1 Skill 1 pattern): `SWEEPER_INTERVAL`
(duration, default 60s), `SWEEPER_ENABLED` (bool, default true — lets tests
disable it). Document in `.env.example`.

### 4.6 — Tests

- Sweeper: an `in_progress` attempt past expiry gets graded on a tick (fake
  clock + in-memory/real repos); already-graded attempts are skipped; race —
  simultaneous sweeper + client submit results in exactly one grading (guarded
  update test).
- Percentile: known set of scores → correct rank/percentile.
- Weak areas: a crafted attempt → expected ranked topics.
- Trend: multiple attempts → ordered series.

---

## Definition of Done

- [ ] `go build ./...` clean; `golangci-lint run` clean.
- [ ] `go test ./internal/worker/... ./internal/usecase/analytics/...
      ./internal/repository/redis/...` passes.
- [ ] An abandoned (`in_progress`, expired) attempt is auto-graded by the sweeper
      without any client request (tested).
- [ ] Sweeper + late submit cannot double-grade (guarded-transition test green).
- [ ] Percentile rank computed from a Redis sorted set; PG remains source of truth
      (sorted set rebuildable).
- [ ] Subject breakdown + weak-area-by-topic + score trend implemented and tested.
- [ ] Sweeper starts/stops with the app lifecycle (verified in Skill 5 wiring).
- [ ] No analytics path mutates attempt rows; no answer contents logged.

## Verification commands

```bash
go build ./... && golangci-lint run
go test ./internal/worker/... ./internal/usecase/analytics/... \
        ./internal/repository/redis/... -v
```

## Hand-off to Skill 5

Engines + sweeper + analytics are complete and tested in isolation. Skill 5
exposes everything over HTTP: quiz/exam/analytics DTOs (answer-key-free),
handlers, routes, RBAC wiring, the sweeper started in `main.go`, and the Phase 2
end-to-end smoke test.

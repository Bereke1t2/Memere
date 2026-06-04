# Phase 2 · Skill 3 — Mock Exam Engine (State Machine + Server Timer)

> **Prerequisite:** Phase 2 Skills 1–2 done (data layer + quiz engine; Redis
> attempt-state and grading patterns established). Read
> [`docs/skill.md`](../../skill.md) §2.
>
> **Spec references:** `memere_Design_Specification.md` §9.2 (exam attempt state
> machine), §4.2.6 (exams & attempts), §9.3 (scoring & analytics), README
> "Quiz & Exam Endpoints" (start / submit / results).

---

## Goal

Implement the **mock exam engine** as usecases (no HTTP yet — Skill 5). Exams
differ from quizzes: they assemble questions from the bank via `exam_questions`,
are always timed (`duration_minutes`), use **marks** (not percentage pass on
points the same way), and follow the explicit §9.2 state machine with
**server-enforced auto-submit on expiry**.

Reuse everything from Skill 2 (Redis state, grading core, IDOR/ownership
patterns). This skill is mostly the exam-specific orchestration + the rigorous
state machine.

---

## The state machine (§9.2) — implement exactly

```
[*] → NOT_STARTED
NOT_STARTED → IN_PROGRESS   : student clicks Start
IN_PROGRESS → IN_PROGRESS   : save answer (auto-save every 30s)
IN_PROGRESS → SUBMITTED     : student submits manually
IN_PROGRESS → EXPIRED       : server timer fires at duration_minutes
SUBMITTED   → GRADED        : auto-grading completes
EXPIRED     → GRADED        : auto-graded on expiry
GRADED      → [*]
```

- `NOT_STARTED` is implicit (no attempt row). First persisted state is
  `in_progress`.
- Every transition goes through `AttemptStatus.CanTransitionTo` (from Skill 1).
  Illegal transitions return `apperror.Conflict("INVALID_ATTEMPT_STATE")`.
- **Grading is mandatory** after `submitted`/`expired` — never leave an attempt in
  a terminal-but-ungraded state in the synchronous path.

---

## Tasks

### 3.1 — Exam authoring usecases (`internal/usecase/exam/`)

For teachers/admins to build exams:

- `CreateExam(ctx, actor, input)` — actor owns the parent course (or admin);
  validate `duration_minutes > 0`, `total_marks`, `pass_marks <= total_marks`,
  subject, grade.
- `AddExamQuestion(ctx, actor, examID, questionID, marks, orderIndex)` —
  ownership-checked; questions come from the bank created in Skill 2.
  Recompute/validate that Σ `exam_questions.marks` reconciles with
  `exams.total_marks` (either enforce equality or recompute total_marks from the
  parts — **decision: recompute `total_marks` from the sum** and keep it
  consistent in a transaction).
- `ListExamQuestions`, `UpdateExam`, `PublishExam` — ownership-checked.
- `ListExams(ctx, actor, filter, cursor, limit)` — published visibility rules like
  Phase 1 courses (anonymous/student see published only).

### 3.2 — Exam attempt usecases

A `Service` holding exam repo + exam-attempt repo + redis state + clock:

- `StartExam(ctx, actor, examID) (*ExamAttemptClientView, error)`
  - resume an existing `in_progress` attempt (idempotent) with remaining time;
    otherwise create one (`status=in_progress`, `started_at`, `expires_at =
    started_at + duration_minutes`).
  - build the question set (respect `randomize` if the exam sets it; persist the
    order to Redis + `answers_snapshot`/`question_order`).
  - return questions via the **client view** (NO correct answers) + `expires_at`
    + `remaining_seconds` + `total_marks`.
- `SaveExamProgress(ctx, actor, attemptID, answers)` — ownership + `in_progress` +
  within time; write to Redis (auto-save target every 30s).
- `SubmitExam(ctx, actor, attemptID, answers) (*ExamResult, error)`
  - ownership + state check; if past `expires_at` → transition `EXPIRED` then
    grade; else transition `SUBMITTED` then grade.
  - grade server-side using the grading core (reuse Skill 2's `grading.go`,
    generalized to take a marks-per-question map for exams), compute `score`,
    `percentage`, subject breakdown; persist `answers_snapshot`,
    `status=graded`; clear Redis.
- `GetExamResult(ctx, actor, attemptID)` — ownership-checked; terminal states
  only; includes per-subject breakdown + pass/fail vs `pass_marks`.
- `ListMyExamAttempts(ctx, actor, examID)`.

### 3.3 — Shared grading generalization

Refactor Skill 2's grading so both engines share it:

- A `gradeAnswers(questions []GradableQuestion, submitted map[qID]answer)
  → GradeResult` where `GradableQuestion` carries `points`/`marks`, `type`,
  correct option(s), `subject`, `topic`.
- Quiz passes `points`; exam passes per-exam `marks`. Output includes total,
  per-subject breakdown, and per-question correctness for the result view.
- Keep this in a shared internal package (e.g. `internal/usecase/grading/`) so
  neither engine duplicates logic. (If you prefer, leave it in `usecase/quiz` and
  import it — but a neutral package reads cleaner. Pick one and be consistent.)

### 3.4 — Expiry semantics (sync path)

In this skill, expiry is handled **lazily**: any `Start`/`Save`/`Submit`/`Get`
that observes `now > expires_at` on an `in_progress` attempt transitions it to
`expired` and grades it. The **background sweeper** that auto-grades abandoned
attempts (student closed the app) is built in Skill 4 — note the dependency.

### 3.5 — Unit tests (mocked repos + fake clock)

Cover: full state machine (each legal transition + each illegal one rejected);
`StartExam` resumes in-progress; submit before expiry → `submitted`→`graded`;
submit/observe after expiry → `expired`→`graded`; correct answers never in client
view; marks summed to `total_marks`; IDOR on another student's attempt →
`FORBIDDEN`; subject breakdown correct.

---

## Definition of Done

- [ ] `go build ./...` clean; `golangci-lint run` clean.
- [ ] `go test ./internal/usecase/exam/... ./internal/usecase/grading/...` passes.
- [ ] Every §9.2 transition is enforced via `CanTransitionTo`; illegal ones
      return `INVALID_ATTEMPT_STATE` (tested).
- [ ] Exam client view contains no correct-answer data (tested on serialized DTO).
- [ ] `expires_at` computed server-side from `duration_minutes`; lazy expiry on
      access works (fake-clock test).
- [ ] Grading core is shared between quiz and exam (no duplicated grading logic).
- [ ] `total_marks` stays consistent with Σ exam-question marks (transactional).

## Verification commands

```bash
go build ./... && golangci-lint run
go test ./internal/usecase/exam/... ./internal/usecase/grading/... -v
grep -rn 'IsCorrect\|correct_answer' internal/delivery 2>/dev/null \
  && echo "REVIEW: must not appear in client DTOs" || echo "OK"
```

## Hand-off to Skill 4

Both engines work synchronously. Skill 4 adds the **background expiry sweeper**
(auto-submit/auto-grade abandoned attempts) and the **scoring & analytics** layer
from §9.3 (percentile via Redis sorted sets, weak-area aggregation, trend
analysis).

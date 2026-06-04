# Phase 2 · Skill 2 — Quiz Engine (Attempts, Grading, Feedback)

> **Prerequisite:** Phase 2 Skill 1 done (quiz/exam data layer, attempt repos,
> client-vs-grading question queries). Read [`docs/skill.md`](../../skill.md) §2.
>
> **Spec references:** `memere_Design_Specification.md` §9.1 (quiz engine design
> table — answer security, randomization, grading, partial attempts, time
> enforcement, attempt limits), §9.3 (scoring), README "Quiz & Exam Endpoints".

---

## Goal

Implement the **quiz engine** as usecases (no HTTP yet — Skill 5). A student can
start an attempt, receive questions **without answer keys** in a randomized order,
auto-save progress, and submit for **server-side grading** with immediate
feedback. Every rule in §9.1 is enforced here.

This skill is also where the **two non-negotiables that define this phase** live:
*answers never leave the server* and *grading is server-side only*.

---

## §9.1 rules → how we implement each

| Spec rule | Implementation |
|---|---|
| Correct answers NEVER sent to client | Client-facing DTOs are built from `GetQuestionsForClient` (no `is_correct`). Grading uses `GetQuestionsForGrading` server-side only. |
| Question randomization, snapshot stored in Redis | On `StartAttempt`, shuffle question + answer order, persist the order to Redis (`quiz:attempt:{id}`) **and** to `quiz_attempts.question_order` JSONB as the durable copy. Same order for the whole attempt; different per attempt. |
| Grading server-side only, on submission | `SubmitAttempt` recomputes score from the grading query; the client score is never trusted. |
| Partial attempts saved to Redis every 30s | `SaveProgress` writes in-progress answers to Redis with the attempt TTL; flushed to PG on submit/expiry. |
| Time enforcement server-side | `started_at` + `expires_at` (= started_at + time_limit) in PG. The client timer is display-only; the server rejects/auto-submits past `expires_at` (sweeper in Skill 4). |
| Attempt limits configurable | `CountAttempts` vs `quiz.max_attempts` before allowing a new attempt. |

---

## Tasks

### 2.1 — Redis attempt state (`internal/repository/redis/quiz_attempt_state.go`)

A small store for live attempt state (separate from the PG record):

- `SetAttemptSnapshot(ctx, attemptID, snapshot, ttl)` — the randomized
  question/answer order (`quiz:attempt:{id}:order`).
- `SaveAnswers(ctx, attemptID, answers, ttl)` — current in-progress answers
  (`quiz:attempt:{id}:answers`).
- `GetSnapshot`, `GetAnswers`, `DeleteAttemptState`.
- TTL = quiz `time_limit_seconds` + a small grace (e.g. +60s), or a sane default
  (e.g. 2h) for untimed quizzes.

### 2.2 — Randomization helper (`pkg/shuffle` or inside usecase)

- Deterministic-per-attempt shuffle: seed from the attempt ID so the order is
  reproducible if Redis is lost (fallback to `question_order` in PG).
- Shuffle questions (only if `quiz.randomize_questions`) and shuffle answer
  options within each question. Persist the resulting order.

### 2.3 — Quiz engine usecases (`internal/usecase/quiz/`)

A `Service` holding quiz repo + quiz-attempt repo + redis state + clock. Methods
(each takes `ctx` + `actor`):

- `GetQuizForStudent(ctx, actor, quizID) (*QuizClientView, error)`
  - student must be enrolled in the quiz's course (reuse Phase 1 enrollment
    check; in Phase 1 enrollment isn't built yet — for now gate on "course
    published or owner", and add a `TODO(phase4): enrollment check`). Ask if the
    user wants strict enrollment gating now.
  - returns metadata + question **count** only, not the questions (questions come
    on StartAttempt).
- `StartAttempt(ctx, actor, quizID) (*AttemptClientView, error)`
  - reject if an `in_progress` attempt already exists (return it instead, with
    remaining time) — idempotent start.
  - enforce `max_attempts` via `CountAttempts`.
  - create the PG attempt (`status=in_progress`, `attempt_number`, `started_at`,
    `expires_at` if timed); build + persist the randomized snapshot (Redis + PG);
    return questions via the **client view** (no `is_correct`) in snapshot order,
    plus `expires_at`/`remaining_seconds`.
- `SaveProgress(ctx, actor, attemptID, answers) error`
  - assert the attempt belongs to `actor` (IDOR), is `in_progress`, not past
    `expires_at`; write answers to Redis. (Called ~every 30s by the client.)
- `SubmitAttempt(ctx, actor, attemptID, answers) (*AttemptResult, error)`
  - assert ownership + `in_progress` + within time (if past expiry → mark
    `expired` and grade what's there, per §9.2).
  - merge final answers (Redis + payload), **grade server-side** using
    `GetQuestionsForGrading`, compute `score`, `percentage`, `passed`
    (vs `pass_percentage`), persist snapshot + results, set `status=graded`,
    clear Redis state.
  - return the **result with feedback**: per-question correct/incorrect, the
    `explanation`, and (now that the attempt is graded) the correct answer **may**
    be revealed per §4.2.5 "explanation shown after submission". Reveal correct
    answers **only in the post-submission result**, never before.
- `GetAttemptResult(ctx, actor, attemptID) (*AttemptResult, error)` — ownership-
  checked; only for `submitted`/`graded`/`expired`.
- `ListMyAttempts(ctx, actor, quizID)` — student's own attempt history.

### 2.4 — Grading logic (`internal/usecase/quiz/grading.go`)

- `multiple_choice` / `true_false`: exact match on the selected answer's
  `is_correct`.
- `short_answer`: normalized comparison (trim, lowercase, collapse spaces) against
  accepted answer(s). Keep it simple; flag ambiguous ones for manual review if the
  spec later requires it. Document the limitation.
- Score = Σ `points` of correct questions; `percentage = score/total*100`;
  `passed = percentage >= pass_percentage`.
- Capture a **subject/topic breakdown** (§9.3) into the result + `answers_snapshot`
  for later analytics.

### 2.5 — Teacher quiz authoring usecases

(So teachers can create the quizzes students take.)

- `CreateQuiz(ctx, actor, input)` — actor must own the parent course (reuse Phase
  1 course-ownership check); validate `pass_percentage`, `time_limit`, etc.
- `AddQuestion(ctx, actor, quizID, input)` + `AddAnswer` — ownership-checked;
  validate exactly one+ correct answer for choice types.
- `UpdateQuiz` / `PublishQuiz` — ownership-checked.

### 2.6 — Unit tests (mocked repos + fake clock)

Cover: answer key never present in client view; randomization differs per attempt
but is stable within one; `max_attempts` enforced; submit grades correctly for
each question type; submit past expiry → `expired` + graded; IDOR (another user's
attempt) → `FORBIDDEN`; double-start returns the existing attempt; pass/fail
threshold; subject breakdown computed.

---

## Definition of Done

- [ ] `go build ./...` clean; `golangci-lint run` clean.
- [ ] `go test ./internal/usecase/quiz/...` passes, including the answer-leak and
      IDOR tests.
- [ ] Client views are built from `GetQuestionsForClient`; a test asserts the
      serialized client view contains no `is_correct` / correct-answer field.
- [ ] Server-side grading ignores any client-supplied score.
- [ ] `max_attempts` and time-window enforcement proven by tests (fake clock).
- [ ] Redis state set on start, updated on save, cleared on submit.

## Verification commands

```bash
go build ./... && golangci-lint run
go test ./internal/usecase/quiz/... -v
# Assert no answer-key leak anywhere in client-facing code paths:
grep -rn 'IsCorrect' internal/usecase/quiz internal/delivery 2>/dev/null \
  && echo "REVIEW: ensure IsCorrect only in grading path"
```

## Hand-off to Skill 3

The quiz engine works end-to-end in isolation. Skill 3 builds the **exam engine**:
the §9.2 state machine, server-enforced exam timer, exam assembly from the
question bank, and exam scoring — reusing the attempt/grading/Redis patterns
established here.

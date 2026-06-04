# Phase 2 · Skill 1 — Quiz & Exam Data Layer

> **Prerequisite:** Phase 1 is complete and green (auth + course API live, schema
> migrated, sqlc generating). Read [`docs/skill.md`](../../skill.md) and its §2
> **Non-Negotiables** again — Phase 2 is where "answers never leave the server"
> and "timers are server-side" become the central design constraints.
>
> **Spec references:** `memere_Design_Specification.md` §4.2.5 (quizzes,
> questions, answers), §4.2.6 (exams, exam_attempts), §9.1 (quiz engine design),
> §9.2 (exam state machine), the ER diagram in §4 (QUIZ_ATTEMPTS, EXAM_QUESTIONS).

---

## Goal

Build the **data foundation** for the quiz and exam engines: additive migrations
for the tables Phase 1 didn't create, the domain entities, the repository
interfaces, and the sqlc queries. No engine logic, no HTTP yet — just a correct,
type-safe data layer the next skills consume.

Phase 1 created `courses.quizzes / questions / answers` and
`courses.exams / exam_attempts` as **schema-only** (migrations 0004, 0005). This
skill adds what's missing and builds the Go layer over all of it.

---

## What Phase 1 left missing (you must add)

1. **`courses.quiz_attempts`** — referenced by the ER diagram but never created.
2. **`courses.exam_questions`** — the join table letting an exam reuse questions
   (ER diagram `EXAMS ||--o{ EXAM_QUESTIONS`). Exams need their own question set
   distinct from per-lesson quizzes.
3. **Subject/topic tags on questions** — §9.3 needs "subject breakdown" and
   "weak areas grouped by topic". Add `subject` and `topic` columns (or a
   `tags JSONB`) to `courses.questions`.
4. **Attempt-limit column on quizzes** — §9.1 "Attempt limits: configurable per
   quiz". Add `max_attempts INTEGER NULL` (NULL = unlimited) to `courses.quizzes`.
5. **`attempt_number`** on attempts to enforce/limit retries.

---

## Tasks

### 1.1 — Additive migrations (`migrations/`, golang-migrate)

Create new numbered pairs (continue from Phase 1's `0007`):

- `0008_quiz_attempts` — `courses.quiz_attempts`:
  - `id` UUID PK, `quiz_id` FK, `student_id` FK → `auth.users`,
  - `attempt_number` INT, `status` ENUM (`in_progress`/`submitted`/`graded`/`expired`),
  - `started_at`, `submitted_at` (nullable), `expires_at` (nullable — for timed
    quizzes), `score` DECIMAL, `percentage` DECIMAL, `passed` BOOLEAN,
  - `question_order` JSONB (the randomized question/answer order snapshot — see
    Skill 2; persisted here as the durable copy of the Redis snapshot),
  - `answers_snapshot` JSONB (student's submitted answers), audit columns.
  - Indexes: `(student_id, quiz_id)`, `status`, `(quiz_id, student_id, attempt_number)` UNIQUE.
- `0009_exam_questions` — `courses.exam_questions`:
  - `id` UUID PK, `exam_id` FK, `question_id` FK, `order_index` INT,
    `marks` INT (per-question marks within this exam), UNIQUE `(exam_id, question_id)`.
- `0010_quiz_exam_enhancements`:
  - `ALTER courses.quizzes ADD max_attempts INTEGER` (NULL = unlimited).
  - `ALTER courses.questions ADD subject VARCHAR(100)`, `ADD topic VARCHAR(150)`
    (both nullable).
  - Any missing columns the spec §4.2.5/§4.2.6 lists that 0004/0005 skipped —
    reconcile against the spec and add them here.

Every `.up.sql` has a clean reversing `.down.sql`. Keep the `set_updated_at`
trigger on new tables. UUID PKs, audit columns, soft delete where user-facing
(`quiz_attempts`/`exam_attempts` are records — they are **not** soft-deleted, they
transition status; do **not** add `deleted_at` to attempts).

### 1.2 — Domain entities (`internal/domain/entity/`)

Pure Go, **no db/json tags** (same rule as Phase 1). Create:

- `quiz.go` — `Quiz` (id, lessonID*, courseID, title, timeLimitSeconds*,
  passPercentage, randomizeQuestions, maxAttempts*, audit) + `QuizSettings` if
  helpful.
- `question.go` — `Question` (id, quizID*, text, type, points, explanation,
  orderIndex, subject*, topic*) + `QuestionType` typed string
  (`multiple_choice`/`true_false`/`short_answer`) with `Valid()`.
- `answer.go` — `Answer` (id, questionID, text, isCorrect, orderIndex).
  ⚠️ `IsCorrect` lives in the domain/repo layer **only** — it must never reach a
  response DTO (enforced in Skill 5; note it here loudly).
- `quiz_attempt.go` — `QuizAttempt` (all 0008 columns) + `AttemptStatus` typed
  string (`in_progress`/`submitted`/`graded`/`expired`) with a `CanTransitionTo`
  method encoding the §9.2 state machine.
- `exam.go` — `Exam` (id, courseID*, title, subject, grade, durationMinutes,
  totalMarks, passMarks, instructions, isPublished, audit).
- `exam_question.go` — `ExamQuestion` (examID, questionID, orderIndex, marks).
- `exam_attempt.go` — `ExamAttempt` (§4.2.6 columns: startedAt, submittedAt*,
  score, percentage, answersSnapshot JSONB, status) reusing `AttemptStatus`.

Model the **state machine** (§9.2) once as a shared `AttemptStatus` with
transitions: `NOT_STARTED → IN_PROGRESS → {SUBMITTED | EXPIRED} → GRADED`.
(`NOT_STARTED` is implicit — no row yet; the first persisted state is
`in_progress`.)

### 1.3 — Repository interfaces (`internal/domain/repository/`)

Interfaces only. Methods take `ctx` first, return entities + error.

- `quiz_repository.go`: `CreateQuiz`, `GetQuizByID`, `GetQuizWithQuestions`
  (questions **with** answers — repo/usecase internal use), `ListQuizzesByCourse`,
  `UpdateQuiz`, `DeleteQuiz` (soft on the quiz, not attempts),
  `CreateQuestion`, `AddAnswer`, `ListQuestionsByQuiz`.
- `quiz_attempt_repository.go`: `CreateAttempt`, `GetAttemptByID`,
  `GetActiveAttempt(studentID, quizID)`, `CountAttempts(studentID, quizID)`,
  `ListAttemptsByStudent`, `UpdateAttempt` (status/score/snapshot),
  `ListExpiredInProgress(now)` (for the sweeper in Phase 2 Skill 4).
- `exam_repository.go`: `CreateExam`, `GetExamByID`, `GetExamWithQuestions`,
  `ListExams(filter, cursor, limit)`, `UpdateExam`, `DeleteExam`,
  `AddExamQuestion`, `ListExamQuestions`.
- `exam_attempt_repository.go`: mirror of quiz-attempt repo for exams, plus
  `ListExpiredInProgress(now)`.

Reuse Phase 1's `pkg/pagination` cursor for list methods.

### 1.4 — sqlc queries (`db/query/`)

Add query files; regenerate. **Every read excludes soft-deleted parents and never
selects `is_correct` into any DTO-bound struct** (it's fine in repo-internal
grading queries):

- `quizzes.sql`, `questions.sql`, `answers.sql`
- `quiz_attempts.sql` — including `GetActiveAttempt`, `CountAttempts`,
  `ListExpiredInProgress` (`WHERE status='in_progress' AND expires_at < $1`).
- `exams.sql`, `exam_questions.sql`, `exam_attempts.sql`

Add **two distinct** question-fetch queries:
- `GetQuestionsForGrading` — includes `is_correct` (server-side only).
- `GetQuestionsForClient` — **excludes** `is_correct` entirely (used to render the
  attempt to the student). Make the leak structurally impossible at the SQL level.

### 1.5 — Repository implementations (`internal/repository/postgres/`)

Implement all four interfaces over sqlc, mapping sqlc models ↔ entities. Map
`pgx.ErrNoRows` → `apperror.NotFound`, unique violations → `apperror.Conflict`.
Compile-time assertions: `var _ repository.QuizRepository = (*postgres.QuizRepo)(nil)`
for each.

---

## Definition of Done

- [ ] `make migrate-up` applies 0008–0010; `make migrate-down` reverses cleanly.
- [ ] `\d courses.quiz_attempts` and `\d courses.exam_questions` show the new
      tables with correct FKs/indexes.
- [ ] `make sqlc` generates code that `go build ./...` accepts.
- [ ] Domain entities have **no** db/json tags (`grep -rn 'db:"\|json:"'
      internal/domain/entity/` is empty).
- [ ] `GetQuestionsForClient` cannot return `is_correct` (verify the generated
      struct has no such field).
- [ ] All four repository interfaces have compiling implementations with
      compile-time assertions.
- [ ] `AttemptStatus.CanTransitionTo` rejects illegal transitions (unit-tested).

## Verification commands

```bash
make migrate-up
psql "$DB_DSN" -c '\d courses.quiz_attempts; \d courses.exam_questions'
make sqlc && go build ./...
grep -rn 'db:"\|json:"' internal/domain/entity/ && echo "FAIL" || echo "OK clean domain"
go test ./internal/domain/entity/... -run Transition -v
make migrate-down && make migrate-up
```

## Hand-off to Skill 2

The quiz/exam data layer is ready and type-safe. Skill 2 builds the **quiz
engine**: starting an attempt, the randomized snapshot in Redis, auto-save,
server-side grading, attempt limits, and immediate feedback.

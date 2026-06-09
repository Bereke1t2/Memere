-- name: CreateQuizAttempt :one
INSERT INTO courses.quiz_attempts (
    quiz_id, student_id, attempt_number, expires_at, question_order,
    answers_snapshot, status
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: GetQuizAttemptByID :one
-- Filters by student_id to prevent IDOR (Non-Negotiable #7).
SELECT * FROM courses.quiz_attempts
WHERE id = $1 AND student_id = $2;

-- name: GetActiveQuizAttempt :one
-- The student's current in-progress attempt at a quiz, if any. Used to resume an
-- attempt and to block starting a second concurrent one.
SELECT * FROM courses.quiz_attempts
WHERE student_id = $1 AND quiz_id = $2 AND status = 'in_progress'
ORDER BY attempt_number DESC
LIMIT 1;

-- name: ListQuizAttemptsByStudentAndQuiz :many
SELECT * FROM courses.quiz_attempts
WHERE student_id = $1 AND quiz_id = $2
ORDER BY attempt_number DESC;

-- name: CountQuizAttempts :one
SELECT count(*) FROM courses.quiz_attempts
WHERE student_id = $1 AND quiz_id = $2;

-- name: ListExpiredQuizAttempts :many
-- In-progress attempts whose server-side deadline has passed, for the auto-grade
-- sweeper (spec §9.1 time enforcement, Non-Negotiable #2). Only timed quizzes
-- (expires_at NOT NULL) are eligible.
SELECT * FROM courses.quiz_attempts
WHERE status = 'in_progress'
  AND expires_at IS NOT NULL
  AND expires_at < $1
ORDER BY expires_at ASC
LIMIT $2;

-- name: UpdateQuizAttempt :one
-- Auto-save (spec §9.1: state saved every 30s) and status transitions. Persists
-- the latest answers snapshot and status without touching the immutable timer
-- columns (started_at/expires_at).
UPDATE courses.quiz_attempts
SET answers_snapshot = $2,
    status = $3,
    submitted_at = $4
WHERE id = $1
RETURNING *;

-- name: ClaimQuizAttemptForGrading :one
-- Race-safe transition out of in_progress (Non-Negotiable #2 server timer): the
-- WHERE status='in_progress' predicate means only the first writer (a late client
-- submit OR the background sweeper) flips the row; the loser matches no row and
-- no-ops. $2 is the target status ('submitted' or 'expired').
UPDATE courses.quiz_attempts
SET status = $2,
    answers_snapshot = $3,
    submitted_at = $4
WHERE id = $1 AND status = 'in_progress'
RETURNING *;

-- name: GradeQuizAttempt :one
UPDATE courses.quiz_attempts
SET score = $2,
    percentage = $3,
    passed = $4,
    submitted_at = now(),
    status = 'graded'
WHERE id = $1
RETURNING *;

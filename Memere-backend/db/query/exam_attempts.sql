-- name: CreateExamAttempt :one
INSERT INTO courses.exam_attempts (
    exam_id, student_id, answers_snapshot, status
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetExamAttemptByID :one
-- Filters by student_id to prevent IDOR (Non-Negotiable #7).
SELECT * FROM courses.exam_attempts
WHERE id = $1 AND student_id = $2;

-- name: GetActiveExamAttempt :one
-- The student's current in-progress sitting of an exam, if any.
SELECT * FROM courses.exam_attempts
WHERE student_id = $1 AND exam_id = $2 AND status = 'in_progress'
ORDER BY started_at DESC
LIMIT 1;

-- name: ListExamAttemptsByStudent :many
SELECT * FROM courses.exam_attempts
WHERE student_id = $1
ORDER BY started_at DESC;

-- name: FindExpiredExamAttempts :many
-- In-progress attempts whose server-side timer has elapsed, for the auto-submit
-- sweeper (spec §9.2). The deadline is derived from exams.duration_minutes; the
-- client timer is display-only (Non-Negotiable #2).
SELECT a.* FROM courses.exam_attempts a
JOIN courses.exams e ON e.id = a.exam_id
WHERE a.status = 'in_progress'
  AND a.started_at + (e.duration_minutes * interval '1 minute') < $1
ORDER BY a.started_at ASC
LIMIT $2;

-- name: UpdateExamAttempt :one
-- Auto-save and status transitions; the immutable started_at timer column is
-- never touched here.
UPDATE courses.exam_attempts
SET answers_snapshot = $2,
    status = $3,
    submitted_at = $4
WHERE id = $1
RETURNING *;

-- name: ClaimExamAttemptForGrading :one
-- Race-safe transition out of in_progress (spec §9.2): WHERE status='in_progress'
-- ensures only the first writer (late client submit OR sweeper) flips the row;
-- the loser matches no row and no-ops. $2 is 'submitted' or 'expired'.
UPDATE courses.exam_attempts
SET status = $2,
    answers_snapshot = $3,
    submitted_at = $4
WHERE id = $1 AND status = 'in_progress'
RETURNING *;

-- name: GradeExamAttempt :one
UPDATE courses.exam_attempts
SET score = $2,
    percentage = $3,
    submitted_at = now(),
    status = 'graded'
WHERE id = $1
RETURNING *;

-- name: ListGradedExamAttemptsBySubject :many
-- Score trend (§9.3): a student's graded attempts for a subject, oldest first, so
-- the caller can plot score over consecutive attempts. Joins exams for subject.
SELECT a.* FROM courses.exam_attempts a
JOIN courses.exams e ON e.id = a.exam_id
WHERE a.student_id = $1
  AND e.subject = $2
  AND a.status = 'graded'
ORDER BY a.submitted_at ASC, a.started_at ASC;

-- name: GetExamAttemptStats :one
-- Teacher/admin exam stats (§9.3): how many graded attempts, the average
-- percentage, and how many passed (score >= the exam's pass_marks).
SELECT
    count(*)                                                   AS total_attempts,
    coalesce(avg(a.percentage), 0)::float8                     AS avg_percentage,
    count(*) FILTER (WHERE a.score >= e.pass_marks)            AS passed_count
FROM courses.exam_attempts a
JOIN courses.exams e ON e.id = a.exam_id
WHERE a.exam_id = $1 AND a.status = 'graded';

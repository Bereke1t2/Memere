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

-- name: GradeExamAttempt :one
UPDATE courses.exam_attempts
SET score = $2,
    percentage = $3,
    submitted_at = now(),
    status = 'graded'
WHERE id = $1
RETURNING *;

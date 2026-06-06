-- name: CreateQuizAttempt :one
INSERT INTO courses.quiz_attempts (
    quiz_id, student_id, attempt_number, answers_snapshot, status
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetQuizAttemptByID :one
-- Filters by student_id to prevent IDOR (Non-Negotiable #7).
SELECT * FROM courses.quiz_attempts
WHERE id = $1 AND student_id = $2;

-- name: ListQuizAttemptsByStudentAndQuiz :many
SELECT * FROM courses.quiz_attempts
WHERE student_id = $1 AND quiz_id = $2
ORDER BY attempt_number DESC;

-- name: CountQuizAttempts :one
SELECT count(*) FROM courses.quiz_attempts
WHERE student_id = $1 AND quiz_id = $2;

-- name: GradeQuizAttempt :one
UPDATE courses.quiz_attempts
SET score = $2,
    percentage = $3,
    submitted_at = now(),
    status = 'graded'
WHERE id = $1
RETURNING *;

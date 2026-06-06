-- name: CreateQuiz :one
INSERT INTO courses.quizzes (
    lesson_id, course_id, title, time_limit_seconds, pass_percentage,
    randomize_questions, max_attempts
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: GetQuizByID :one
SELECT * FROM courses.quizzes
WHERE id = $1;

-- name: ListQuizzesByCourse :many
SELECT * FROM courses.quizzes
WHERE course_id = $1
ORDER BY created_at DESC, id DESC;

-- name: UpdateQuiz :one
UPDATE courses.quizzes
SET title = $2,
    time_limit_seconds = $3,
    pass_percentage = $4,
    randomize_questions = $5,
    max_attempts = $6
WHERE id = $1
RETURNING *;

-- name: DeleteQuiz :exec
DELETE FROM courses.quizzes
WHERE id = $1;

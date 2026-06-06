-- name: CreateQuestion :one
INSERT INTO courses.questions (
    quiz_id, text, type, points, explanation, order_index
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetQuestionByID :one
SELECT * FROM courses.questions
WHERE id = $1;

-- name: ListQuestionsByQuiz :many
SELECT * FROM courses.questions
WHERE quiz_id = $1
ORDER BY order_index ASC, id ASC;

-- name: UpdateQuestion :one
UPDATE courses.questions
SET text = $2,
    type = $3,
    points = $4,
    explanation = $5,
    order_index = $6
WHERE id = $1
RETURNING *;

-- name: DeleteQuestion :exec
DELETE FROM courses.questions
WHERE id = $1;

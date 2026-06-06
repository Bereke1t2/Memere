-- answers.is_correct is server-only data (Non-Negotiable #1). Delivery DTOs must
-- omit it; only the server-side grading path reads it.

-- name: CreateAnswer :one
INSERT INTO courses.answers (
    question_id, text, is_correct, order_index
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: ListAnswersByQuestion :many
SELECT * FROM courses.answers
WHERE question_id = $1
ORDER BY order_index ASC, id ASC;

-- name: DeleteAnswersByQuestion :exec
DELETE FROM courses.answers
WHERE question_id = $1;

-- answers.is_correct is server-only data (Non-Negotiable #1). The client path
-- uses ListAnswersForClient (no is_correct); only the grading path reads the key.

-- name: CreateAnswer :one
INSERT INTO courses.answers (
    question_id, text, is_correct, order_index
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- ListAnswersByQuestion is server-internal: it includes is_correct for grading
-- and teacher-facing reads. Never map it onto a client DTO.
-- name: ListAnswersByQuestion :many
SELECT * FROM courses.answers
WHERE question_id = $1
ORDER BY order_index ASC, id ASC;

-- ListAnswersForClient renders options to the student. It omits is_correct at
-- the SQL level, so the leak is structurally impossible in the generated struct.
-- name: ListAnswersForClient :many
SELECT id, question_id, text, order_index
FROM courses.answers
WHERE question_id = $1
ORDER BY order_index ASC, id ASC;

-- name: DeleteAnswersByQuestion :exec
DELETE FROM courses.answers
WHERE question_id = $1;

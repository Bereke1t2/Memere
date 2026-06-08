-- name: CreateQuestion :one
INSERT INTO courses.questions (
    quiz_id, text, type, points, explanation, order_index, subject, topic
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: GetQuestionByID :one
SELECT * FROM courses.questions
WHERE id = $1;

-- name: ListQuestionsByQuiz :many
SELECT * FROM courses.questions
WHERE quiz_id = $1
ORDER BY order_index ASC, id ASC;

-- GetQuestionsForClient renders the attempt to the student. It selects ONLY the
-- non-secret columns — there is no is_correct anywhere in the projection, so the
-- generated row struct structurally cannot carry the answer key (Non-Negotiable
-- #1). Answer options for these questions are fetched via ListAnswersForClient.
-- name: GetQuestionsForClient :many
SELECT id, quiz_id, text, type, points, order_index, subject, topic
FROM courses.questions
WHERE quiz_id = $1
ORDER BY order_index ASC, id ASC;

-- GetQuestionsForGrading is the server-side-only counterpart: it returns each
-- question joined to its answer options INCLUDING is_correct, for the grading
-- core. This result must never be mapped onto a client DTO.
-- name: GetQuestionsForGrading :many
SELECT
    q.id          AS question_id,
    q.points      AS points,
    q.type        AS type,
    a.id          AS answer_id,
    a.text        AS answer_text,
    a.is_correct  AS is_correct,
    a.order_index AS answer_order_index
FROM courses.questions q
JOIN courses.answers a ON a.question_id = q.id
WHERE q.quiz_id = $1
ORDER BY q.order_index ASC, q.id ASC, a.order_index ASC, a.id ASC;

-- name: UpdateQuestion :one
UPDATE courses.questions
SET text = $2,
    type = $3,
    points = $4,
    explanation = $5,
    order_index = $6,
    subject = $7,
    topic = $8
WHERE id = $1
RETURNING *;

-- name: DeleteQuestion :exec
DELETE FROM courses.questions
WHERE id = $1;

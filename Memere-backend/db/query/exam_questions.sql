-- name: AddExamQuestion :one
INSERT INTO courses.exam_questions (
    exam_id, question_id, order_index, marks
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: ListExamQuestions :many
SELECT * FROM courses.exam_questions
WHERE exam_id = $1
ORDER BY order_index ASC, id ASC;

-- ListExamQuestionsForClient renders the exam to the student: the per-exam order
-- and marks joined to the question text/type, with NO is_correct anywhere in the
-- projection (Non-Negotiable #1). Answer options come from ListAnswersForClient.
-- name: ListExamQuestionsForClient :many
SELECT
    eq.question_id  AS question_id,
    eq.order_index  AS order_index,
    eq.marks        AS marks,
    q.text          AS text,
    q.type          AS type,
    q.subject       AS subject,
    q.topic         AS topic
FROM courses.exam_questions eq
JOIN courses.questions q ON q.id = eq.question_id
WHERE eq.exam_id = $1
ORDER BY eq.order_index ASC, eq.id ASC;

-- ListExamQuestionsForGrading is server-side-only: per-exam marks joined to each
-- answer option INCLUDING is_correct, for the shared grading core. Never map onto
-- a client DTO.
-- name: ListExamQuestionsForGrading :many
SELECT
    eq.question_id  AS question_id,
    eq.marks        AS marks,
    q.type          AS type,
    q.subject       AS subject,
    q.topic         AS topic,
    a.id            AS answer_id,
    a.text          AS answer_text,
    a.is_correct    AS is_correct,
    a.order_index   AS answer_order_index
FROM courses.exam_questions eq
JOIN courses.questions q ON q.id = eq.question_id
JOIN courses.answers a ON a.question_id = q.id
WHERE eq.exam_id = $1
ORDER BY eq.order_index ASC, eq.id ASC, a.order_index ASC, a.id ASC;

-- name: UpdateExamQuestion :one
-- Re-points an existing link's per-exam order and mark value. Scoped by exam_id
-- as well as question_id so one exam's edit can never touch another's link.
UPDATE courses.exam_questions
SET order_index = $3,
    marks = $4
WHERE exam_id = $1 AND question_id = $2
RETURNING *;

-- name: RemoveExamQuestion :exec
-- Detaches a question from an exam. This is a true DELETE of the JOIN ROW only —
-- the question itself survives in the bank for reuse by other exams, so
-- Non-Negotiable #5 (no hard deletes of user content) is not in play here.
DELETE FROM courses.exam_questions
WHERE exam_id = $1 AND question_id = $2;

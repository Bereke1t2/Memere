-- name: CreateQuiz :one
-- course_id is nullable since 0026: a NULL means a standalone quiz, authored on
-- its own rather than under a course.
INSERT INTO courses.quizzes (
    lesson_id, course_id, title, time_limit_seconds, pass_percentage,
    randomize_questions, max_attempts, is_published
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: GetQuizByID :one
SELECT * FROM courses.quizzes
WHERE id = $1 AND deleted_at IS NULL;

-- name: ListQuizzesByCourse :many
SELECT * FROM courses.quizzes
WHERE course_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC, id DESC;

-- name: ListQuizzes :many
-- Keyset (cursor) pagination ordered by (created_at, id) DESC, for the admin
-- panel's cross-course quiz list. Nullable filter args are ignored when NULL.
-- 'standalone' = true narrows to course-less quizzes, false to course-attached
-- ones. Soft-deleted quizzes are excluded.
SELECT * FROM courses.quizzes
WHERE deleted_at IS NULL
  AND (sqlc.narg('course_id')::uuid IS NULL OR course_id = sqlc.narg('course_id'))
  AND (sqlc.narg('is_published')::bool IS NULL OR is_published = sqlc.narg('is_published'))
  AND (
      sqlc.narg('standalone')::bool IS NULL
      OR (sqlc.narg('standalone')::bool AND course_id IS NULL)
      OR (NOT sqlc.narg('standalone')::bool AND course_id IS NOT NULL)
  )
  AND (
      sqlc.narg('after_created_at')::timestamptz IS NULL
      OR (created_at, id) < (sqlc.narg('after_created_at')::timestamptz, sqlc.narg('after_id')::uuid)
  )
ORDER BY created_at DESC, id DESC
LIMIT sqlc.arg('row_limit');

-- name: CountQuestionsByQuiz :one
-- Question tally for the authoring lists, so the admin panel does not have to
-- fetch every question just to show a count.
SELECT COUNT(*) FROM courses.questions
WHERE quiz_id = $1;

-- name: UpdateQuiz :one
UPDATE courses.quizzes
SET title = $2,
    time_limit_seconds = $3,
    pass_percentage = $4,
    randomize_questions = $5,
    max_attempts = $6,
    lesson_id = $7,
    is_published = $8
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteQuiz :exec
-- Non-Negotiable #5: tombstone, never physically DELETE. Attempts referencing
-- the quiz are preserved as historical records.
UPDATE courses.quizzes
SET deleted_at = now()
WHERE id = $1 AND deleted_at IS NULL;

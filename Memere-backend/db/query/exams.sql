-- name: CreateExam :one
INSERT INTO courses.exams (
    course_id, title, subject, grade, duration_minutes, total_marks,
    pass_marks, instructions, is_published
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: GetExamByID :one
SELECT * FROM courses.exams
WHERE id = $1 AND deleted_at IS NULL;

-- name: ListExams :many
-- Keyset (cursor) pagination ordered by (created_at, id) DESC. Nullable filter
-- args are ignored when NULL. Soft-deleted exams are excluded.
SELECT * FROM courses.exams
WHERE deleted_at IS NULL
  AND (sqlc.narg('subject')::text IS NULL OR subject = sqlc.narg('subject'))
  AND (sqlc.narg('grade')::int IS NULL OR grade = sqlc.narg('grade'))
  AND (sqlc.narg('is_published')::bool IS NULL OR is_published = sqlc.narg('is_published'))
  AND (
      sqlc.narg('after_created_at')::timestamptz IS NULL
      OR (created_at, id) < (sqlc.narg('after_created_at')::timestamptz, sqlc.narg('after_id')::uuid)
  )
ORDER BY created_at DESC, id DESC
LIMIT sqlc.arg('row_limit');

-- name: UpdateExam :one
UPDATE courses.exams
SET title = $2,
    subject = $3,
    grade = $4,
    duration_minutes = $5,
    total_marks = $6,
    pass_marks = $7,
    instructions = $8,
    is_published = $9
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteExam :exec
-- Non-Negotiable #5: tombstone, never physically DELETE.
UPDATE courses.exams
SET deleted_at = now()
WHERE id = $1 AND deleted_at IS NULL;

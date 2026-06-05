-- name: CreateLesson :one
INSERT INTO courses.lessons (
    section_id, course_id, title, type, order_index, is_free_preview,
    duration_seconds, is_published
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: GetLessonByID :one
SELECT * FROM courses.lessons
WHERE id = $1 AND deleted_at IS NULL;

-- name: ListLessonsBySection :many
SELECT * FROM courses.lessons
WHERE section_id = $1 AND deleted_at IS NULL
ORDER BY order_index ASC, created_at ASC;

-- name: ListLessonsByCourse :many
SELECT * FROM courses.lessons
WHERE course_id = $1 AND deleted_at IS NULL
ORDER BY order_index ASC, created_at ASC;

-- name: UpdateLesson :one
UPDATE courses.lessons
SET title = $2,
    type = $3,
    order_index = $4,
    is_free_preview = $5,
    duration_seconds = $6,
    is_published = $7
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteLesson :exec
UPDATE courses.lessons
SET deleted_at = now()
WHERE id = $1 AND deleted_at IS NULL;

-- name: SoftDeleteLessonsBySection :exec
-- Cascade soft-delete: when a section is removed, tombstone its live lessons too.
UPDATE courses.lessons
SET deleted_at = now()
WHERE section_id = $1 AND deleted_at IS NULL;

-- name: MaxLessonOrderIndex :one
-- Highest order_index among a section's live lessons; -1 when none exist.
SELECT COALESCE(MAX(order_index), -1)::int AS max_order_index
FROM courses.lessons
WHERE section_id = $1 AND deleted_at IS NULL;

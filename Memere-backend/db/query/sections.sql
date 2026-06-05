-- name: CreateSection :one
INSERT INTO courses.course_sections (course_id, title, description, order_index, is_published)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetSectionByID :one
SELECT * FROM courses.course_sections
WHERE id = $1 AND deleted_at IS NULL;

-- name: ListSectionsByCourse :many
SELECT * FROM courses.course_sections
WHERE course_id = $1 AND deleted_at IS NULL
ORDER BY order_index ASC, created_at ASC;

-- name: UpdateSection :one
UPDATE courses.course_sections
SET title = $2,
    description = $3,
    order_index = $4,
    is_published = $5
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteSection :exec
UPDATE courses.course_sections
SET deleted_at = now()
WHERE id = $1 AND deleted_at IS NULL;

-- name: MaxSectionOrderIndex :one
-- Highest order_index among a course's live sections; -1 when none exist, so
-- callers can use (max + 1) for the next append position.
SELECT COALESCE(MAX(order_index), -1)::int AS max_order_index
FROM courses.course_sections
WHERE course_id = $1 AND deleted_at IS NULL;

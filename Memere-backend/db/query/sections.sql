-- name: CreateSection :one
INSERT INTO courses.course_sections (course_id, title, description, order_index, is_published)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetSectionByID :one
SELECT * FROM courses.course_sections
WHERE id = $1;

-- name: ListSectionsByCourse :many
SELECT * FROM courses.course_sections
WHERE course_id = $1
ORDER BY order_index ASC, created_at ASC;

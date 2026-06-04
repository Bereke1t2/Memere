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
WHERE id = $1;

-- name: ListLessonsBySection :many
SELECT * FROM courses.lessons
WHERE section_id = $1
ORDER BY order_index ASC, created_at ASC;

-- name: ListLessonsByCourse :many
SELECT * FROM courses.lessons
WHERE course_id = $1
ORDER BY order_index ASC, created_at ASC;

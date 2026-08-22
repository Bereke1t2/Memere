-- name: CreateLesson :one
INSERT INTO courses.lessons (
    section_id, course_id, title, type, order_index, is_free_preview,
    duration_seconds, is_published, content, pdf_url
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
RETURNING *;

-- name: GetLessonByID :one
SELECT l.id, l.section_id, l.course_id, l.title, l.type, l.order_index, l.is_free_preview, l.duration_seconds, l.is_published, l.content, l.pdf_url, l.created_at, l.updated_at, l.deleted_at, v.id AS video_id, q.id AS quiz_id
FROM courses.lessons l
LEFT JOIN courses.videos v ON v.lesson_id = l.id AND v.deleted_at IS NULL
LEFT JOIN courses.quizzes q ON q.lesson_id = l.id AND q.deleted_at IS NULL
WHERE l.id = $1 AND l.deleted_at IS NULL;

-- name: ListLessonsBySection :many
SELECT l.id, l.section_id, l.course_id, l.title, l.type, l.order_index, l.is_free_preview, l.duration_seconds, l.is_published, l.content, l.pdf_url, l.created_at, l.updated_at, l.deleted_at, v.id AS video_id, q.id AS quiz_id
FROM courses.lessons l
LEFT JOIN courses.videos v ON v.lesson_id = l.id AND v.deleted_at IS NULL
LEFT JOIN courses.quizzes q ON q.lesson_id = l.id AND q.deleted_at IS NULL
WHERE l.section_id = $1 AND l.deleted_at IS NULL
ORDER BY l.order_index ASC, l.created_at ASC;

-- name: ListLessonsByCourse :many
SELECT l.id, l.section_id, l.course_id, l.title, l.type, l.order_index, l.is_free_preview, l.duration_seconds, l.is_published, l.content, l.pdf_url, l.created_at, l.updated_at, l.deleted_at, v.id AS video_id, q.id AS quiz_id
FROM courses.lessons l
LEFT JOIN courses.videos v ON v.lesson_id = l.id AND v.deleted_at IS NULL
LEFT JOIN courses.quizzes q ON q.lesson_id = l.id AND q.deleted_at IS NULL
WHERE l.course_id = $1 AND l.deleted_at IS NULL
ORDER BY l.order_index ASC, l.created_at ASC;

-- name: UpdateLesson :one
UPDATE courses.lessons
SET title = $2,
    type = $3,
    order_index = $4,
    is_free_preview = $5,
    duration_seconds = $6,
    is_published = $7,
    content = COALESCE($8, content),
    pdf_url = COALESCE($9, pdf_url)
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

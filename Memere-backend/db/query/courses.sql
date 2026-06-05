-- name: CreateCourse :one
INSERT INTO courses.courses (
    teacher_id, title, slug, description, short_description, subject, grade,
    thumbnail_url, price, currency, is_free, is_published, language, level,
    metadata
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
)
RETURNING *;

-- name: GetCourseByID :one
SELECT * FROM courses.courses
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetCourseBySlug :one
SELECT * FROM courses.courses
WHERE slug = $1 AND deleted_at IS NULL;

-- name: ListCourses :many
-- Keyset (cursor) pagination ordered by (created_at, id) DESC. Nullable filter
-- args are ignored when NULL. after_created_at/after_id is the decoded cursor.
SELECT * FROM courses.courses
WHERE deleted_at IS NULL
  AND (sqlc.narg('subject')::text IS NULL OR subject = sqlc.narg('subject'))
  AND (sqlc.narg('grade')::int IS NULL OR grade = sqlc.narg('grade'))
  AND (sqlc.narg('is_published')::bool IS NULL OR is_published = sqlc.narg('is_published'))
  AND (sqlc.narg('teacher_id')::uuid IS NULL OR teacher_id = sqlc.narg('teacher_id'))
  AND (
      sqlc.narg('after_created_at')::timestamptz IS NULL
      OR (created_at, id) < (sqlc.narg('after_created_at')::timestamptz, sqlc.narg('after_id')::uuid)
  )
ORDER BY created_at DESC, id DESC
LIMIT sqlc.arg('row_limit');

-- name: UpdateCourse :one
UPDATE courses.courses
SET title = $2,
    description = $3,
    short_description = $4,
    subject = $5,
    grade = $6,
    thumbnail_url = $7,
    price = $8,
    currency = $9,
    is_free = $10,
    is_published = $11,
    language = $12,
    level = $13,
    metadata = $14
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteCourse :exec
UPDATE courses.courses
SET deleted_at = now()
WHERE id = $1 AND deleted_at IS NULL;

-- name: RecomputeCourseCounters :exec
-- Recompute the denormalized lesson counters from the live (non-deleted)
-- lessons. Run inside the same transaction as the lesson mutation that
-- triggered it so the course row stays consistent.
UPDATE courses.courses c
SET total_lessons = sub.cnt,
    total_duration_seconds = sub.dur
FROM (
    SELECT COUNT(*)::int AS cnt,
           COALESCE(SUM(duration_seconds), 0)::int AS dur
    FROM courses.lessons
    WHERE course_id = $1 AND deleted_at IS NULL
) AS sub
WHERE c.id = $1 AND c.deleted_at IS NULL;

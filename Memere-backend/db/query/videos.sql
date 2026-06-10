-- name: CreateVideo :one
INSERT INTO courses.videos (id, lesson_id, course_id, processing_status, original_file_key, file_size_bytes)
VALUES (gen_random_uuid(), $1, $2, $3, $4, $5)
RETURNING *;

-- name: GetVideoByID :one
SELECT * FROM courses.videos WHERE id = $1 AND deleted_at IS NULL;

-- name: GetVideoByLessonID :one
SELECT * FROM courses.videos WHERE lesson_id = $1 AND deleted_at IS NULL;

-- UpdateVideoStatusGuarded flips status only when the row still holds from_status;
-- :execrows lets the repo report claimed=false (0 rows) without an error so two
-- workers can race for the same video safely.
-- name: UpdateVideoStatusGuarded :execrows
UPDATE courses.videos
SET processing_status = @to_status, updated_at = now()
WHERE id = @id AND processing_status = @from_status AND deleted_at IS NULL;

-- name: SetVideoReady :one
UPDATE courses.videos
SET processing_status = 'ready',
    hls_master_key = $2, resolution_480p_key = $3, resolution_720p_key = $4,
    resolution_1080p_key = $5, thumbnail_key = $6, duration_seconds = $7,
    processed_at = now(), updated_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: ListVideosByStatus :many
SELECT * FROM courses.videos
WHERE processing_status = $1 AND deleted_at IS NULL
ORDER BY created_at ASC LIMIT $2;

-- name: SoftDeleteVideo :exec
UPDATE courses.videos
SET deleted_at = now(), updated_at = now()
WHERE id = $1 AND deleted_at IS NULL;

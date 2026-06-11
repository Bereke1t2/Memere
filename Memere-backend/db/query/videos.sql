-- CreateVideo inserts a pending video. The id is supplied by the application
-- (not gen_random_uuid()) because the upload usecase must embed it in the
-- storage-key path before the row exists; honoring it keeps the returned id, the
-- key path, and the row consistent.
-- name: CreateVideo :one
INSERT INTO courses.videos (id, lesson_id, course_id, processing_status, original_file_key, file_size_bytes)
VALUES ($1, $2, $3, $4, $5, $6)
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

-- SetVideoReady is guarded on the processing status so a video can only reach
-- ready (atomically, with all its keys) from processing. A second/duplicate
-- worker finds 0 rows (RETURNING empty -> ErrNoRows) and cannot double-ready or
-- overwrite a finalized row.
-- name: SetVideoReady :one
UPDATE courses.videos
SET processing_status = 'ready',
    hls_master_key = $2, resolution_480p_key = $3, resolution_720p_key = $4,
    resolution_1080p_key = $5, thumbnail_key = $6, duration_seconds = $7,
    processed_at = now(), updated_at = now()
WHERE id = $1 AND processing_status = 'processing' AND deleted_at IS NULL
RETURNING *;

-- name: SetVideoError :exec
UPDATE courses.videos
SET processing_error = $2, updated_at = now()
WHERE id = $1 AND deleted_at IS NULL;

-- name: ListVideosByStatus :many
SELECT * FROM courses.videos
WHERE processing_status = $1 AND deleted_at IS NULL
ORDER BY created_at ASC LIMIT $2;

-- name: SoftDeleteVideo :exec
UPDATE courses.videos
SET deleted_at = now(), updated_at = now()
WHERE id = $1 AND deleted_at IS NULL;

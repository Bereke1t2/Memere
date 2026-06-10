-- Phase 3 Skill 1 — make courses.videos Go-ready for the video pipeline.
-- Additive only. The base table (migration 0003) predates the Phase 3 entity,
-- so it lacks: the denormalized course_id (fast video -> course authz without a
-- lesson join), the soft-delete column (Non-Negotiable #5), and the processing
-- bookkeeping columns the transcode worker (Skill 3) writes on success/failure.
ALTER TABLE courses.videos
    ADD COLUMN IF NOT EXISTS course_id        UUID,
    ADD COLUMN IF NOT EXISTS processing_error TEXT,
    ADD COLUMN IF NOT EXISTS processed_at     TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS deleted_at       TIMESTAMPTZ;

-- Backfill the denormalized course_id from each video's lesson (lessons already
-- carry a denormalized course_id), then enforce NOT NULL + the FK.
UPDATE courses.videos v
SET course_id = l.course_id
FROM courses.lessons l
WHERE v.lesson_id = l.id AND v.course_id IS NULL;

ALTER TABLE courses.videos
    ALTER COLUMN course_id SET NOT NULL;

ALTER TABLE courses.videos
    ADD CONSTRAINT videos_course_id_fkey
    FOREIGN KEY (course_id) REFERENCES courses.courses (id) ON DELETE CASCADE;

-- Replace the original table-level UNIQUE(lesson_id) with a partial unique index
-- so a soft-deleted video does not permanently block re-uploading to that lesson
-- (Non-Negotiable #5). The base table's inline UNIQUE is named videos_lesson_id_key.
ALTER TABLE courses.videos DROP CONSTRAINT IF EXISTS videos_lesson_id_key;

CREATE UNIQUE INDEX IF NOT EXISTS videos_lesson_id_uniq
    ON courses.videos (lesson_id) WHERE deleted_at IS NULL;

-- The boot reconciler / transcode worker (Skills 2-3) scan by processing_status;
-- authz and soft-delete reads filter by course_id / deleted_at.
CREATE INDEX IF NOT EXISTS videos_processing_status_idx
    ON courses.videos (processing_status);

CREATE INDEX IF NOT EXISTS videos_course_id_idx
    ON courses.videos (course_id);

CREATE INDEX IF NOT EXISTS videos_deleted_at_idx
    ON courses.videos (deleted_at);

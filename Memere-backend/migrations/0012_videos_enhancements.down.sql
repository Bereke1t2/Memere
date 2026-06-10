DROP INDEX IF EXISTS courses.videos_deleted_at_idx;
DROP INDEX IF EXISTS courses.videos_course_id_idx;
DROP INDEX IF EXISTS courses.videos_processing_status_idx;
DROP INDEX IF EXISTS courses.videos_lesson_id_uniq;

-- Restore the original table-level uniqueness on lesson_id.
ALTER TABLE courses.videos
    ADD CONSTRAINT videos_lesson_id_key UNIQUE (lesson_id);

ALTER TABLE courses.videos DROP CONSTRAINT IF EXISTS videos_course_id_fkey;

ALTER TABLE courses.videos
    DROP COLUMN IF EXISTS deleted_at,
    DROP COLUMN IF EXISTS processed_at,
    DROP COLUMN IF EXISTS processing_error,
    DROP COLUMN IF EXISTS course_id;

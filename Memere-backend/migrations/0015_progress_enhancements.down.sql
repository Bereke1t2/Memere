DROP TABLE IF EXISTS progress.course_progress;
DROP TABLE IF EXISTS progress.study_streaks;
DROP INDEX IF EXISTS progress_student_course_idx;
DROP INDEX IF EXISTS progress_student_lesson_uniq;
ALTER TABLE progress.progress DROP COLUMN IF EXISTS deleted_at;
-- Restore the original full unique constraint.
ALTER TABLE progress.progress
    ADD CONSTRAINT progress_student_id_lesson_id_key UNIQUE (student_id, lesson_id);

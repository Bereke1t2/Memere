DROP INDEX IF EXISTS courses.idx_lessons_deleted_at;
DROP INDEX IF EXISTS courses.idx_course_sections_deleted_at;

ALTER TABLE courses.lessons DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE courses.course_sections DROP COLUMN IF EXISTS deleted_at;

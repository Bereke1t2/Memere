ALTER TABLE progress.study_streaks DROP COLUMN IF EXISTS last_warned_date;
DROP INDEX IF EXISTS courses.certificates_student_idx;
DROP TABLE IF EXISTS courses.certificates;

DROP INDEX IF EXISTS courses.idx_exams_deleted_at;
DROP INDEX IF EXISTS courses.idx_quizzes_deleted_at;
ALTER TABLE courses.exams DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE courses.quizzes DROP COLUMN IF EXISTS deleted_at;

ALTER TABLE courses.questions DROP COLUMN IF EXISTS topic;
ALTER TABLE courses.questions DROP COLUMN IF EXISTS subject;

DROP INDEX IF EXISTS courses.idx_quiz_attempts_status_expires_at;

ALTER TABLE courses.quiz_attempts DROP COLUMN IF EXISTS passed;
ALTER TABLE courses.quiz_attempts DROP COLUMN IF EXISTS expires_at;
ALTER TABLE courses.quiz_attempts DROP COLUMN IF EXISTS question_order;

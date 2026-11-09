-- Reverting re-imposes NOT NULL, which fails while any standalone quiz exists.
-- Standalone quizzes have no course to fall back to, so the rollback deliberately
-- refuses rather than guessing a course_id or destroying authored content; delete
-- or re-parent those quizzes first.
DROP INDEX IF EXISTS courses.idx_quizzes_standalone;

ALTER TABLE courses.quizzes DROP COLUMN IF EXISTS is_published;

ALTER TABLE courses.quizzes ALTER COLUMN course_id SET NOT NULL;

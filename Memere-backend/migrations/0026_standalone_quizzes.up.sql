-- Standalone quizzes (admin-authored, not attached to a course).
--
-- 0004 created quizzes.course_id NOT NULL, which forced every quiz to live under
-- a course. Exams have always allowed a NULL course_id (0005) so an admin can
-- publish a free-standing mock exam; quizzes now match that shape, letting the
-- admin panel author a quiz independently of the course tree.
--
-- Ownership follows the exam rule (see exam.loadOwnedExam): a quiz with a
-- course_id is editable by that course's teacher or an admin, while a quiz with
-- course_id IS NULL is editable by admins only — there is no teacher to scope it
-- to.
ALTER TABLE courses.quizzes ALTER COLUMN course_id DROP NOT NULL;

-- is_published mirrors exams.is_published and exists to gate STANDALONE quizzes
-- only. A course-attached quiz is gated by access to its course exactly as
-- before and ignores this column, so adding it changes no existing behaviour; a
-- course-less quiz has no paywall to inherit, so it stays admin-only until
-- published (see quiz.assertQuizAccess). Defaulting to false means a quiz being
-- authored is never reachable mid-edit.
ALTER TABLE courses.quizzes ADD COLUMN is_published BOOLEAN NOT NULL DEFAULT false;

-- Standalone quizzes are listed on their own page in the admin panel, so the
-- course_id IS NULL partition gets a partial index rather than riding the
-- existing full index on course_id.
CREATE INDEX idx_quizzes_standalone ON courses.quizzes (created_at DESC, id DESC)
    WHERE course_id IS NULL AND deleted_at IS NULL;

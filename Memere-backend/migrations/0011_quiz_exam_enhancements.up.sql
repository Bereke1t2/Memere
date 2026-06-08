-- Phase 2 Skill 1 enhancements (spec §9.1, §9.3). Additive only; reconciles the
-- columns 0009 (quiz_attempts) and 0004 (questions) left out of the engine
-- design.
--
-- quiz_attempts gains the durable copy of the per-attempt randomized order
-- (question_order, the snapshot mirrored from Redis — spec §9.1), the
-- server-side timer deadline for timed quizzes (expires_at — Non-Negotiable #2),
-- and the cached pass/fail outcome (passed).
ALTER TABLE courses.quiz_attempts ADD COLUMN question_order JSONB;
ALTER TABLE courses.quiz_attempts ADD COLUMN expires_at TIMESTAMPTZ;
ALTER TABLE courses.quiz_attempts ADD COLUMN passed BOOLEAN;

-- The sweeper (Phase 2 Skill 4) scans for in-progress attempts past their
-- deadline; index the predicate it filters on.
CREATE INDEX idx_quiz_attempts_status_expires_at
    ON courses.quiz_attempts (status, expires_at);

-- questions gain subject/topic tags so §9.3 can produce the subject breakdown
-- and weak-area-by-topic analytics. Both nullable (legacy rows untagged).
ALTER TABLE courses.questions ADD COLUMN subject VARCHAR(100);
ALTER TABLE courses.questions ADD COLUMN topic VARCHAR(150);

-- quizzes and exams are user-facing content tables, so Non-Negotiable #5
-- (soft delete only) applies — 0004/0005 omitted the tombstone column. Reads
-- filter deleted_at IS NULL; Delete sets the timestamp instead of DELETE.
-- Attempts are deliberately excluded: they transition status, never soft-delete.
ALTER TABLE courses.quizzes ADD COLUMN deleted_at TIMESTAMPTZ;
ALTER TABLE courses.exams ADD COLUMN deleted_at TIMESTAMPTZ;

CREATE INDEX idx_quizzes_deleted_at ON courses.quizzes (deleted_at);
CREATE INDEX idx_exams_deleted_at ON courses.exams (deleted_at);

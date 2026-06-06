-- Quiz attempts (spec §9.1). Additive migration for the Phase 2 quiz engine.
-- answers_snapshot stores the student's submitted answers as JSONB; the
-- correct-answer key is never stored here (Non-Negotiable #1).
CREATE TABLE courses.quiz_attempts (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    quiz_id          UUID NOT NULL REFERENCES courses.quizzes (id) ON DELETE RESTRICT,
    student_id       UUID NOT NULL REFERENCES auth.users (id) ON DELETE RESTRICT,
    attempt_number   INTEGER NOT NULL DEFAULT 1,
    started_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    submitted_at     TIMESTAMPTZ,
    score            DECIMAL(10, 2),
    percentage       DECIMAL(5, 2),
    answers_snapshot JSONB,
    status           VARCHAR(20) NOT NULL DEFAULT 'in_progress'
                         CHECK (status IN ('in_progress', 'submitted', 'graded', 'expired')),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (quiz_id, student_id, attempt_number)
);

CREATE INDEX idx_quiz_attempts_quiz_id ON courses.quiz_attempts (quiz_id);
CREATE INDEX idx_quiz_attempts_student_id ON courses.quiz_attempts (student_id);
CREATE INDEX idx_quiz_attempts_status ON courses.quiz_attempts (status);

CREATE TRIGGER trg_quiz_attempts_updated_at
    BEFORE UPDATE ON courses.quiz_attempts
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Exam schema (spec §4.2.6). Created now; Go layer deferred to Phase 2.
CREATE TABLE courses.exams (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id        UUID REFERENCES courses.courses (id) ON DELETE RESTRICT,
    title            VARCHAR(200) NOT NULL,
    subject          VARCHAR(100) NOT NULL,
    grade            INTEGER NOT NULL,
    duration_minutes INTEGER NOT NULL,
    total_marks      INTEGER NOT NULL,
    pass_marks       INTEGER NOT NULL,
    instructions     TEXT,
    is_published     BOOLEAN NOT NULL DEFAULT false,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_exams_course_id ON courses.exams (course_id);
CREATE INDEX idx_exams_subject ON courses.exams (subject);

CREATE TRIGGER trg_exams_updated_at
    BEFORE UPDATE ON courses.exams
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- exam_attempts.started_at is the server-side timer source of truth (spec §9.2 /
-- Non-Negotiable #2). answers_snapshot is the student's answers at submission.
CREATE TABLE courses.exam_attempts (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    exam_id          UUID NOT NULL REFERENCES courses.exams (id) ON DELETE RESTRICT,
    student_id       UUID NOT NULL REFERENCES auth.users (id) ON DELETE RESTRICT,
    started_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    submitted_at     TIMESTAMPTZ,
    score            DECIMAL(10, 2),
    percentage       DECIMAL(5, 2),
    answers_snapshot JSONB,
    status           VARCHAR(20) NOT NULL DEFAULT 'in_progress'
                         CHECK (status IN ('in_progress', 'submitted', 'graded', 'expired')),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_exam_attempts_exam_id ON courses.exam_attempts (exam_id);
CREATE INDEX idx_exam_attempts_student_id ON courses.exam_attempts (student_id);
CREATE INDEX idx_exam_attempts_status ON courses.exam_attempts (status);

CREATE TRIGGER trg_exam_attempts_updated_at
    BEFORE UPDATE ON courses.exam_attempts
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

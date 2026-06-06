-- Exam question-bank linkage (spec §9.2). Additive migration for Phase 2. Lets a
-- single bank question be reused across multiple exams with a per-exam order and
-- mark value, instead of duplicating question rows per exam.
CREATE TABLE courses.exam_questions (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    exam_id     UUID NOT NULL REFERENCES courses.exams (id) ON DELETE CASCADE,
    question_id UUID NOT NULL REFERENCES courses.questions (id) ON DELETE RESTRICT,
    order_index INTEGER NOT NULL DEFAULT 0,
    marks       INTEGER NOT NULL DEFAULT 1,
    UNIQUE (exam_id, question_id)
);

CREATE INDEX idx_exam_questions_exam_id ON courses.exam_questions (exam_id);
CREATE INDEX idx_exam_questions_question_id ON courses.exam_questions (question_id);

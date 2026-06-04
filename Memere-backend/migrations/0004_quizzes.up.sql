-- Quiz schema (spec §4.2.5). Created now for a coherent schema; the Go layer
-- (entities, repos, queries) is deferred to Phase 2.
CREATE TABLE courses.quizzes (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lesson_id           UUID REFERENCES courses.lessons (id) ON DELETE RESTRICT,
    course_id           UUID NOT NULL REFERENCES courses.courses (id) ON DELETE RESTRICT,
    title               VARCHAR(200) NOT NULL,
    time_limit_seconds  INTEGER,
    pass_percentage     DECIMAL(5, 2) NOT NULL DEFAULT 0,
    randomize_questions BOOLEAN NOT NULL DEFAULT false,
    max_attempts        INTEGER,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_quizzes_lesson_id ON courses.quizzes (lesson_id);
CREATE INDEX idx_quizzes_course_id ON courses.quizzes (course_id);

CREATE TRIGGER trg_quizzes_updated_at
    BEFORE UPDATE ON courses.quizzes
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TABLE courses.questions (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    quiz_id     UUID NOT NULL REFERENCES courses.quizzes (id) ON DELETE CASCADE,
    text        TEXT NOT NULL,
    type        VARCHAR(20) NOT NULL CHECK (type IN ('multiple_choice', 'true_false', 'short_answer')),
    points      INTEGER NOT NULL DEFAULT 1,
    explanation TEXT,
    order_index INTEGER NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_questions_quiz_id ON courses.questions (quiz_id);

CREATE TRIGGER trg_questions_updated_at
    BEFORE UPDATE ON courses.questions
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- answers.is_correct is server-only data and must never be sent to clients
-- (spec §9.1 / Non-Negotiable #1).
CREATE TABLE courses.answers (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    question_id UUID NOT NULL REFERENCES courses.questions (id) ON DELETE CASCADE,
    text        TEXT NOT NULL,
    is_correct  BOOLEAN NOT NULL DEFAULT false,
    order_index INTEGER NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_answers_question_id ON courses.answers (question_id);

CREATE TRIGGER trg_answers_updated_at
    BEFORE UPDATE ON courses.answers
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

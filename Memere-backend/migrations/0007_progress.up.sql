-- Progress tracking (spec §4.2.8). Created now; Go layer deferred to Phase 5.
-- course_id is denormalized for fast course-level aggregation.
CREATE TABLE progress.progress (
    id                    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id            UUID NOT NULL REFERENCES auth.users (id) ON DELETE RESTRICT,
    lesson_id             UUID NOT NULL REFERENCES courses.lessons (id) ON DELETE RESTRICT,
    course_id             UUID NOT NULL REFERENCES courses.courses (id) ON DELETE RESTRICT,
    is_completed          BOOLEAN NOT NULL DEFAULT false,
    completed_at          TIMESTAMPTZ,
    video_progress_seconds INTEGER NOT NULL DEFAULT 0,
    last_accessed_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at            TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (student_id, lesson_id)
);

CREATE INDEX idx_progress_student_id ON progress.progress (student_id);
CREATE INDEX idx_progress_course_id ON progress.progress (course_id);

CREATE TRIGGER trg_progress_updated_at
    BEFORE UPDATE ON progress.progress
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

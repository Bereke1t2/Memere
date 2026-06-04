-- courses.courses (spec §4.2.2)
CREATE TABLE courses.courses (
    id                     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    teacher_id             UUID NOT NULL REFERENCES auth.users (id) ON DELETE RESTRICT,
    title                  VARCHAR(200) NOT NULL,
    slug                   VARCHAR(200) UNIQUE NOT NULL,
    description            TEXT NOT NULL,
    short_description      VARCHAR(500),
    subject                VARCHAR(100) NOT NULL,
    grade                  INTEGER NOT NULL,
    thumbnail_url          TEXT,
    price                  DECIMAL(10, 2) NOT NULL DEFAULT 0,
    currency               VARCHAR(3) NOT NULL DEFAULT 'ETB',
    is_free                BOOLEAN NOT NULL DEFAULT false,
    is_published           BOOLEAN NOT NULL DEFAULT false,
    language               VARCHAR(10) NOT NULL DEFAULT 'en',
    level                  VARCHAR(20) NOT NULL DEFAULT 'beginner'
                               CHECK (level IN ('beginner', 'intermediate', 'advanced')),
    total_duration_seconds INTEGER NOT NULL DEFAULT 0,
    total_lessons          INTEGER NOT NULL DEFAULT 0,
    rating_avg             DECIMAL(3, 2) NOT NULL DEFAULT 0,
    enrollment_count       INTEGER NOT NULL DEFAULT 0,
    metadata               JSONB,
    created_at             TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at             TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at             TIMESTAMPTZ
);

CREATE INDEX idx_courses_teacher_id ON courses.courses (teacher_id);
CREATE INDEX idx_courses_subject ON courses.courses (subject);
CREATE INDEX idx_courses_grade ON courses.courses (grade);
CREATE INDEX idx_courses_is_published ON courses.courses (is_published);
CREATE INDEX idx_courses_price ON courses.courses (price);
CREATE INDEX idx_courses_enrollment_count ON courses.courses (enrollment_count);
CREATE INDEX idx_courses_deleted_at ON courses.courses (deleted_at);
CREATE INDEX idx_courses_created_at ON courses.courses (created_at);

CREATE TRIGGER trg_courses_updated_at
    BEFORE UPDATE ON courses.courses
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- courses.course_sections (spec §4.2.3)
CREATE TABLE courses.course_sections (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id    UUID NOT NULL REFERENCES courses.courses (id) ON DELETE RESTRICT,
    title        VARCHAR(200) NOT NULL,
    description  TEXT,
    order_index  INTEGER NOT NULL DEFAULT 0,
    is_published BOOLEAN NOT NULL DEFAULT false,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_course_sections_course_id ON courses.course_sections (course_id);

CREATE TRIGGER trg_course_sections_updated_at
    BEFORE UPDATE ON courses.course_sections
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- courses.lessons (spec §4.2.3) — course_id denormalized for fast queries.
CREATE TABLE courses.lessons (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    section_id       UUID NOT NULL REFERENCES courses.course_sections (id) ON DELETE RESTRICT,
    course_id        UUID NOT NULL REFERENCES courses.courses (id) ON DELETE RESTRICT,
    title            VARCHAR(200) NOT NULL,
    type             VARCHAR(20) NOT NULL CHECK (type IN ('video', 'note', 'quiz', 'mixed')),
    order_index      INTEGER NOT NULL DEFAULT 0,
    is_free_preview  BOOLEAN NOT NULL DEFAULT false,
    duration_seconds INTEGER NOT NULL DEFAULT 0,
    is_published     BOOLEAN NOT NULL DEFAULT false,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_lessons_section_id ON courses.lessons (section_id);
CREATE INDEX idx_lessons_course_id ON courses.lessons (course_id);

CREATE TRIGGER trg_lessons_updated_at
    BEFORE UPDATE ON courses.lessons
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- courses.videos (spec §4.2.4) — one video per lesson; Go layer deferred to Phase 3.
CREATE TABLE courses.videos (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lesson_id           UUID NOT NULL UNIQUE REFERENCES courses.lessons (id) ON DELETE CASCADE,
    original_file_key   TEXT,
    hls_master_key      TEXT,
    processing_status   VARCHAR(20) NOT NULL DEFAULT 'pending'
                            CHECK (processing_status IN ('pending', 'processing', 'ready', 'failed')),
    duration_seconds    INTEGER NOT NULL DEFAULT 0,
    resolution_480p_key TEXT,
    resolution_720p_key TEXT,
    resolution_1080p_key TEXT,
    thumbnail_key       TEXT,
    file_size_bytes     BIGINT NOT NULL DEFAULT 0,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TRIGGER trg_videos_updated_at
    BEFORE UPDATE ON courses.videos
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

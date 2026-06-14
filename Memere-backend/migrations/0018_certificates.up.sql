-- Phase 5 Skill 4: certificate issuance table and engagement sweeper support.

-- 1. Completion certificates (one per student+course, idempotent issuance).
CREATE TABLE IF NOT EXISTS courses.certificates (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES auth.users(id),
    course_id  UUID NOT NULL REFERENCES courses.courses(id),
    serial     TEXT NOT NULL UNIQUE,     -- public verification code (non-guessable)
    file_key   TEXT,                      -- S3/MinIO key of the rendered PDF
    issued_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (student_id, course_id)        -- one cert per completion
);

CREATE INDEX IF NOT EXISTS certificates_student_idx ON courses.certificates(student_id);

-- 2. Engagement sweeper idempotency: track the last date a streak-warning push
--    was sent so we never double-warn on the same calendar day (spec §11.2).
ALTER TABLE progress.study_streaks
    ADD COLUMN IF NOT EXISTS last_warned_date DATE;

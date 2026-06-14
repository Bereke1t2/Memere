-- Phase 5 Skill 1: progress tracking enhancements.
-- The base progress.progress table was created in 0007 without soft-delete
-- support. This migration adds deleted_at, replaces the full UNIQUE constraint
-- with a partial unique index (required for the ON CONFLICT upsert), adds the
-- course-scoped index, and creates the streak + rollup tables.

-- 1. Soft-delete column (satisfies §2 Non-Negotiable: soft deletes only).
ALTER TABLE progress.progress
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;

-- 2. Drop the old full UNIQUE constraint created by 0007 so the partial index
--    below becomes the sole enforcer. The auto-generated name follows the
--    PostgreSQL convention "{table}_{col1}_{col2}_key".
DO $$ BEGIN
    ALTER TABLE progress.progress
        DROP CONSTRAINT IF EXISTS progress_student_id_lesson_id_key;
EXCEPTION WHEN undefined_object THEN NULL;
END $$;

-- 3. Partial unique index — the ON CONFLICT target for UpsertLessonProgress.
CREATE UNIQUE INDEX IF NOT EXISTS progress_student_lesson_uniq
    ON progress.progress (student_id, lesson_id) WHERE deleted_at IS NULL;

-- 4. Course-scoped index for fast per-student course listing.
CREATE INDEX IF NOT EXISTS progress_student_course_idx
    ON progress.progress (student_id, course_id);

-- 5. Study streaks: one row per student.
CREATE TABLE IF NOT EXISTS progress.study_streaks (
    student_id     UUID PRIMARY KEY REFERENCES auth.users(id),
    current_streak INTEGER NOT NULL DEFAULT 0,
    longest_streak INTEGER NOT NULL DEFAULT 0,
    last_study_date DATE,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- 6. Per-course completion rollup (denormalised cache; rebuildable via
--    RecomputeCourseProgress).
CREATE TABLE IF NOT EXISTS progress.course_progress (
    student_id        UUID NOT NULL REFERENCES auth.users(id),
    course_id         UUID NOT NULL REFERENCES courses.courses(id),
    completed_lessons INTEGER NOT NULL DEFAULT 0,
    total_lessons     INTEGER NOT NULL DEFAULT 0,
    percent_complete  NUMERIC(5,2) NOT NULL DEFAULT 0,
    completed_at      TIMESTAMPTZ,
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (student_id, course_id)
);

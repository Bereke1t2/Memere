-- Soft-delete support for sections and lessons (Non-Negotiable #5). The
-- courses table already carries deleted_at (migration 0003); these two child
-- tables were missing it, so deletes had nowhere to record the tombstone.
ALTER TABLE courses.course_sections ADD COLUMN deleted_at TIMESTAMPTZ;
ALTER TABLE courses.lessons ADD COLUMN deleted_at TIMESTAMPTZ;

CREATE INDEX idx_course_sections_deleted_at ON courses.course_sections (deleted_at);
CREATE INDEX idx_lessons_deleted_at ON courses.lessons (deleted_at);

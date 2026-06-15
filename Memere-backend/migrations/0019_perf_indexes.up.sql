-- Migration 0019: performance indexes for Phase 6 Skill 3 (§12.2 query tuning).
--
-- Composite/partial indexes that cover the hottest filter patterns identified
-- from EXPLAIN ANALYZE on the published-course catalog, payment reconciliation,
-- and in-progress attempt expiry sweeper.

-- Catalog hot path: published courses filtered by subject + grade, ordered by
-- enrollment_count DESC and created_at DESC. The partial index covers only
-- live (non-deleted, published) rows — all catalog reads qualify.
CREATE INDEX IF NOT EXISTS idx_courses_published_subject_grade
    ON courses.courses (subject, grade, enrollment_count DESC, created_at DESC)
    WHERE deleted_at IS NULL AND is_published = true;

-- Payment reconciliation + revenue queries filter by status and created_at or
-- paid_at. The composite index supports both the admin payment list and the
-- revenue roll-up without a seq scan on the full payments table.
CREATE INDEX IF NOT EXISTS idx_payments_status_paid_at
    ON payments.payments (status, paid_at DESC)
    WHERE deleted_at IS NULL;

-- Attempt expiry sweeper (Phase 2 §9.2) scans in_progress attempts past their
-- deadline. A partial index on the status column reduces the scan to the small
-- subset of live attempts rather than the whole table.
CREATE INDEX IF NOT EXISTS idx_exam_attempts_inprogress
    ON courses.exam_attempts (expires_at)
    WHERE status = 'in_progress';

CREATE INDEX IF NOT EXISTS idx_quiz_attempts_inprogress
    ON courses.quiz_attempts (expires_at)
    WHERE status = 'in_progress';

-- Progress dashboard: fetch all progress rows for a student across courses.
-- The existing progress_student_course_idx covers (student_id, course_id) but
-- not ordering by updated_at, which the dashboard query needs.
CREATE INDEX IF NOT EXISTS idx_progress_student_updated
    ON progress.lesson_progress (student_id, updated_at DESC)
    WHERE deleted_at IS NULL;

-- Notification unread count: the existing notifications_user_unread_idx covers
-- (user_id) WHERE is_read = false; extend to include created_at for the list
-- query ordering without a filesort.
CREATE INDEX IF NOT EXISTS idx_notifications_user_created
    ON progress.notifications (user_id, created_at DESC)
    WHERE deleted_at IS NULL;

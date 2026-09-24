DROP TRIGGER IF EXISTS trg_course_access_requests_updated_at ON courses.course_access_requests;
DROP TABLE IF EXISTS courses.course_access_requests;

ALTER TABLE payments.enrollments DROP CONSTRAINT IF EXISTS enrollments_source_check;
ALTER TABLE payments.enrollments ADD CONSTRAINT enrollments_source_check CHECK (source IN ('purchase', 'subscription', 'free', 'coupon'));

DROP INDEX IF EXISTS auth.idx_users_approval_status;
ALTER TABLE auth.users DROP CONSTRAINT IF EXISTS users_approval_status_check;
ALTER TABLE auth.users DROP COLUMN IF EXISTS approval_status;

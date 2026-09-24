-- 0024: Account approval and course access requests

-- 1. Add approval_status to auth.users (defaults to 'approved' for existing rows, then altered to 'pending')
ALTER TABLE auth.users ADD COLUMN IF NOT EXISTS approval_status VARCHAR(20) NOT NULL DEFAULT 'approved';
ALTER TABLE auth.users ADD CONSTRAINT users_approval_status_check CHECK (approval_status IN ('pending', 'approved', 'rejected'));
CREATE INDEX IF NOT EXISTS idx_users_approval_status ON auth.users (approval_status);
ALTER TABLE auth.users ALTER COLUMN approval_status SET DEFAULT 'pending';

-- 2. Extend enrollments.source CHECK constraint to allow 'admin_grant'
ALTER TABLE payments.enrollments DROP CONSTRAINT IF EXISTS enrollments_source_check;
ALTER TABLE payments.enrollments ADD CONSTRAINT enrollments_source_check CHECK (source IN ('purchase', 'subscription', 'free', 'coupon', 'admin_grant'));

-- 3. Create courses.course_access_requests
CREATE TABLE IF NOT EXISTS courses.course_access_requests (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id       UUID NOT NULL REFERENCES auth.users (id) ON DELETE CASCADE,
    course_id        UUID NOT NULL REFERENCES courses.courses (id) ON DELETE CASCADE,
    status           VARCHAR(20) NOT NULL DEFAULT 'pending'
                         CHECK (status IN ('pending', 'approved', 'rejected')),
    reviewed_by      UUID REFERENCES auth.users (id) ON DELETE SET NULL,
    reviewed_at      TIMESTAMPTZ,
    rejection_reason TEXT,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (student_id, course_id)
);

CREATE INDEX IF NOT EXISTS idx_course_access_requests_student ON courses.course_access_requests (student_id);
CREATE INDEX IF NOT EXISTS idx_course_access_requests_course ON courses.course_access_requests (course_id);
CREATE INDEX IF NOT EXISTS idx_course_access_requests_status ON courses.course_access_requests (status);

CREATE TRIGGER trg_course_access_requests_updated_at
    BEFORE UPDATE ON courses.course_access_requests
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

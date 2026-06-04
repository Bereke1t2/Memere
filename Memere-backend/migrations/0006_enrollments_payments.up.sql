-- Enrollment, payment, and coupon schema (spec §4.2.7, §10.3). Created now; Go
-- layer deferred to Phase 4.
CREATE TABLE payments.coupons (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code           VARCHAR(50) UNIQUE NOT NULL,
    discount_type  VARCHAR(20) NOT NULL CHECK (discount_type IN ('percentage', 'fixed_amount')),
    discount_value DECIMAL(10, 2) NOT NULL,
    max_uses       INTEGER,
    used_count     INTEGER NOT NULL DEFAULT 0,
    expires_at     TIMESTAMPTZ,
    applicable_to  VARCHAR(20) NOT NULL DEFAULT 'all'
                       CHECK (applicable_to IN ('all', 'specific_courses', 'subscription_only')),
    course_ids     UUID[],
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TRIGGER trg_coupons_updated_at
    BEFORE UPDATE ON payments.coupons
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TABLE payments.enrollments (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id  UUID NOT NULL REFERENCES auth.users (id) ON DELETE RESTRICT,
    course_id   UUID NOT NULL REFERENCES courses.courses (id) ON DELETE RESTRICT,
    enrolled_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at  TIMESTAMPTZ,
    source      VARCHAR(20) NOT NULL CHECK (source IN ('purchase', 'subscription', 'free', 'coupon')),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (student_id, course_id)
);

CREATE INDEX idx_enrollments_student_id ON payments.enrollments (student_id);
CREATE INDEX idx_enrollments_course_id ON payments.enrollments (course_id);

CREATE TRIGGER trg_enrollments_updated_at
    BEFORE UPDATE ON payments.enrollments
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- idempotency_key enforces the no-double-charge rule (spec §7.3 / Non-Negotiable #4).
CREATE TABLE payments.payments (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id              UUID NOT NULL REFERENCES auth.users (id) ON DELETE RESTRICT,
    course_id               UUID REFERENCES courses.courses (id) ON DELETE RESTRICT,
    amount                  DECIMAL(10, 2) NOT NULL,
    currency                VARCHAR(3) NOT NULL DEFAULT 'ETB',
    provider                VARCHAR(20) NOT NULL CHECK (provider IN ('stripe', 'chapa', 'telebirr')),
    provider_transaction_id VARCHAR(255),
    status                  VARCHAR(20) NOT NULL DEFAULT 'pending'
                                CHECK (status IN ('pending', 'completed', 'failed', 'refunded')),
    coupon_id               UUID REFERENCES payments.coupons (id) ON DELETE RESTRICT,
    idempotency_key         VARCHAR(255) UNIQUE NOT NULL,
    metadata                JSONB,
    paid_at                 TIMESTAMPTZ,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_payments_student_id ON payments.payments (student_id);
CREATE INDEX idx_payments_course_id ON payments.payments (course_id);
CREATE INDEX idx_payments_status ON payments.payments (status);

CREATE TRIGGER trg_payments_updated_at
    BEFORE UPDATE ON payments.payments
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

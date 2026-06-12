-- Phase 4 · Skill 1 — additive payment/enrollment/subscription enhancements.
-- 0006 already created payments.payments / enrollments / coupons with a NOT NULL
-- UNIQUE idempotency_key, so we DO NOT re-add it here. We add the columns the
-- Phase 4 flows need, a partial-unique index that makes (provider,
-- provider_transaction_id) idempotent, the webhook dedup ledger, and the
-- subscriptions table (spec §10.2).

ALTER TABLE payments.payments
    ADD COLUMN IF NOT EXISTS provider_checkout_id TEXT,
    ADD COLUMN IF NOT EXISTS subscription_id      UUID,
    ADD COLUMN IF NOT EXISTS failure_reason       TEXT;

-- A completed payment maps to exactly one provider transaction; the partial
-- unique index makes a re-delivered webhook (same provider txn) a no-op insert,
-- never a second row. NULL provider_transaction_id (pending payments) is exempt.
CREATE UNIQUE INDEX IF NOT EXISTS payments_provider_txn_uniq
    ON payments.payments (provider, provider_transaction_id)
    WHERE provider_transaction_id IS NOT NULL;

-- Webhook dedup ledger (spec §7.3 "webhook deduplication"). The UNIQUE
-- (provider, provider_event_id) is the hard guarantee: InsertIfNew relies on
-- ON CONFLICT DO NOTHING so a re-delivered event is processed at most once.
CREATE TABLE IF NOT EXISTS payments.webhook_events (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    provider          TEXT NOT NULL,
    provider_event_id TEXT NOT NULL,
    event_type        TEXT NOT NULL,
    payload           JSONB NOT NULL,
    processed_at      TIMESTAMPTZ,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (provider, provider_event_id)
);

-- Subscriptions (spec §10.2). An active subscription is all-access via the
-- enrollment/access layer built in Skill 2.
CREATE TABLE IF NOT EXISTS payments.subscriptions (
    id                       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id               UUID NOT NULL REFERENCES auth.users (id) ON DELETE RESTRICT,
    plan                     TEXT NOT NULL CHECK (plan IN ('monthly', 'annual')),
    status                   TEXT NOT NULL CHECK (status IN ('active', 'past_due', 'canceled', 'expired')),
    current_period_start     TIMESTAMPTZ NOT NULL,
    current_period_end       TIMESTAMPTZ NOT NULL,
    provider                 TEXT NOT NULL,
    provider_subscription_id TEXT,
    created_at               TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at               TIMESTAMPTZ NOT NULL DEFAULT now(),
    canceled_at              TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_subscriptions_student ON payments.subscriptions (student_id, status);

CREATE TRIGGER trg_subscriptions_updated_at
    BEFORE UPDATE ON payments.subscriptions
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- subscription_id on a payment references the subscription it settled (nullable;
-- a course purchase has none). Added after subscriptions exists.
ALTER TABLE payments.payments
    ADD CONSTRAINT payments_subscription_id_fkey
    FOREIGN KEY (subscription_id) REFERENCES payments.subscriptions (id) ON DELETE SET NULL;

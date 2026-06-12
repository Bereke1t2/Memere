-- Reverse 0013 in dependency order.
ALTER TABLE payments.payments DROP CONSTRAINT IF EXISTS payments_subscription_id_fkey;

DROP TRIGGER IF EXISTS trg_subscriptions_updated_at ON payments.subscriptions;
DROP INDEX IF EXISTS payments.idx_subscriptions_student;
DROP TABLE IF EXISTS payments.subscriptions;

DROP TABLE IF EXISTS payments.webhook_events;

DROP INDEX IF EXISTS payments.payments_provider_txn_uniq;

ALTER TABLE payments.payments
    DROP COLUMN IF EXISTS failure_reason,
    DROP COLUMN IF EXISTS subscription_id,
    DROP COLUMN IF EXISTS provider_checkout_id;

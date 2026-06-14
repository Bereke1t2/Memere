-- Restore the original provider CHECK (without 'mock'). Any rows written by the
-- mock provider must be removed first or this will fail.
ALTER TABLE payments.payments
    DROP CONSTRAINT IF EXISTS payments_provider_check;

ALTER TABLE payments.payments
    ADD CONSTRAINT payments_provider_check
    CHECK (provider IN ('stripe', 'chapa', 'telebirr'));

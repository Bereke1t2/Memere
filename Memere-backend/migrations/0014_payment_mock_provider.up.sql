-- Phase 4 §5.6: allow the test-only 'mock' payment provider so the end-to-end
-- smoke test can settle a payment offline (no external provider HTTP). The mock
-- provider is only ever registered when PAYMENT_MOCK_ENABLED is set, so this
-- widened CHECK does not expose anything new in production.
ALTER TABLE payments.payments
    DROP CONSTRAINT IF EXISTS payments_provider_check;

ALTER TABLE payments.payments
    ADD CONSTRAINT payments_provider_check
    CHECK (provider IN ('stripe', 'chapa', 'telebirr', 'mock'));

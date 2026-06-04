# Phase 4 · Skill 1 — Payment & Enrollment Data Layer + Provider Port

> **Prerequisite:** Phases 1–3 complete and green. Read
> [`docs/skill.md`](../../skill.md) §2 — Phase 4 is governed by *"all payments use
> idempotency keys; webhooks are deduplicated"* and *"filter every query by the
> authenticated user_id"*.
>
> **Spec references:** `memere_Design_Specification.md` §4.2.7 (enrollments &
> payments), §10.1 (providers), §10.3 (coupon system), §7.3 (payment replay
> mitigation), README "Payment Integration".

---

## Goal

Build the **data + abstraction foundation** for money: the Go layer for
`payments`, `enrollments`, `coupons` (created schema-only in Phase 1 migration
`0006`), plus the **`PaymentProvider` port** that Chapa/Telebirr/Stripe will
implement, and an **idempotency** mechanism. No HTTP, no real provider calls yet.

The cardinal rule of this phase: **money operations must be idempotent and
auditable**. Every payment has a client-supplied idempotency key; every webhook is
deduplicated; enrollment is granted exactly once per successful payment.

---

## Key decisions for Phase 4

| Concern | Choice | Why |
|---|---|---|
| Provider integration | `PaymentProvider` interface; Chapa (primary), Telebirr, Stripe impls | Spec §10.1; abstraction isolates provider churn (a named risk in §14.4). |
| Idempotency | Client `Idempotency-Key` header → `payments.idempotency_key` UNIQUE | Spec §7.3 "no double-charge". |
| Webhook dedup | `webhook_events` table keyed by provider event id | Spec §7.3 "webhook deduplication". |
| Fulfillment | Webhook is the **source of truth** for completion; polling is advisory | Avoids granting access on an unconfirmed payment. |
| Currency | ETB (Chapa/Telebirr), USD/EUR (Stripe) | Spec §10.1. |

---

## Tasks

### 1.1 — Reconcile + extend migrations (`migrations/0012_payments_enhancements`)

Phase 1's `0006` created `payments.payments / enrollments / coupons`. Add what the
flows need (additive; continue numbering after Phase 3's `0011`):

```sql
-- migrations/0012_payments_enhancements.up.sql
ALTER TABLE payments.payments
    ADD COLUMN IF NOT EXISTS idempotency_key TEXT,
    ADD COLUMN IF NOT EXISTS provider_checkout_id TEXT,
    ADD COLUMN IF NOT EXISTS subscription_id UUID,
    ADD COLUMN IF NOT EXISTS failure_reason TEXT;

CREATE UNIQUE INDEX IF NOT EXISTS payments_idempotency_key_uniq
    ON payments.payments (idempotency_key) WHERE idempotency_key IS NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS payments_provider_txn_uniq
    ON payments.payments (provider, provider_transaction_id)
    WHERE provider_transaction_id IS NOT NULL;

-- webhook dedup ledger
CREATE TABLE IF NOT EXISTS payments.webhook_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    provider TEXT NOT NULL,
    provider_event_id TEXT NOT NULL,
    event_type TEXT NOT NULL,
    payload JSONB NOT NULL,
    processed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (provider, provider_event_id)
);

-- subscriptions (spec §10.2 references SUBSCRIPTIONS)
CREATE TABLE IF NOT EXISTS payments.subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES auth.users(id),
    plan TEXT NOT NULL,                 -- monthly | annual
    status TEXT NOT NULL,               -- active | past_due | canceled | expired
    current_period_start TIMESTAMPTZ NOT NULL,
    current_period_end   TIMESTAMPTZ NOT NULL,
    provider TEXT NOT NULL,
    provider_subscription_id TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    canceled_at TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS subscriptions_student_idx ON payments.subscriptions(student_id, status);
```

`down.sql` reverses in order. Keep `set_updated_at` triggers on new tables.

### 1.2 — Domain entities (`internal/domain/entity/`) — no db/json tags

```go
// payment.go
type PaymentProvider string
const ( ProviderChapa PaymentProvider = "chapa"; ProviderTelebirr PaymentProvider = "telebirr"; ProviderStripe PaymentProvider = "stripe" )

type PaymentStatus string
const ( PayPending PaymentStatus="pending"; PayCompleted PaymentStatus="completed"; PayFailed PaymentStatus="failed"; PayRefunded PaymentStatus="refunded" )

func (s PaymentStatus) CanTransitionTo(n PaymentStatus) bool {
    switch s {
    case PayPending:   return n==PayCompleted || n==PayFailed
    case PayCompleted: return n==PayRefunded
    default:           return false // failed/refunded terminal
    }
}

type Payment struct {
    ID uuid.UUID; StudentID uuid.UUID; CourseID *uuid.UUID // nil = subscription
    SubscriptionID *uuid.UUID
    Amount decimal.Decimal; Currency string
    Provider PaymentProvider; ProviderTxnID *string; ProviderCheckoutID *string
    Status PaymentStatus; IdempotencyKey *string; CouponID *uuid.UUID
    FailureReason *string; Metadata map[string]any
    PaidAt *time.Time; CreatedAt, UpdatedAt time.Time
}
```

Also: `enrollment.go` (§4.2.7 — `EnrollmentSource` enum
`purchase|subscription|free|coupon`, `ExpiresAt *time.Time`),
`coupon.go` (§10.3 — `DiscountType` `percentage|fixed_amount`, `MaxUses`,
`UsedCount`, `ExpiresAt`, `ApplicableTo`, `CourseIDs []uuid.UUID`, an
`Apply(price) (discounted, error)` method that validates expiry/uses/scope),
`subscription.go`, and `webhook_event.go`.

> Use a decimal type (`shopspring/decimal`) for money — never `float64`. Add it to
> `go.mod`. Note this in code.

### 1.3 — Provider port (`internal/domain/service/payment_provider.go`)

```go
package service

import "context"

type CheckoutRequest struct {
    PaymentID   string
    Amount      decimal.Decimal
    Currency    string
    CustomerEmail string
    Description string
    ReturnURL   string // where the provider redirects after pay
    CallbackURL string // our webhook endpoint
    IdempotencyKey string
}
type CheckoutResult struct {
    CheckoutID  string // provider's reference
    RedirectURL string // payment_url the client opens (or client_secret)
}

type WebhookEvent struct {
    ProviderEventID string
    Type            string // e.g. "charge.success"
    PaymentRef      string // maps back to our PaymentID or provider_txn
    Status          string // normalized: completed | failed
    RawPayload      []byte
}

type PaymentProvider interface {
    Name() string
    CreateCheckout(ctx context.Context, req CheckoutRequest) (*CheckoutResult, error)
    // VerifyWebhook validates signature + parses into a normalized event.
    VerifyWebhook(ctx context.Context, headers map[string]string, body []byte) (*WebhookEvent, error)
    // VerifyPayment is a server-to-server confirmation (polling fallback).
    VerifyPayment(ctx context.Context, providerRef string) (status string, err error)
}

// Registry resolves a provider by name.
type PaymentProviderRegistry interface {
    Get(p PaymentProvider) (PaymentProvider, error) // by entity.PaymentProvider name
}
```

### 1.4 — Repository interfaces + sqlc + impls

Interfaces (`internal/domain/repository/`): `PaymentRepository`
(`Create`, `GetByID`, `GetByIdempotencyKey`, `GetByProviderTxn`,
`UpdateStatusGuarded(id, from, to, fields)`, `ListByStudent`),
`EnrollmentRepository` (`Create`, `Exists(studentID, courseID)`,
`GetActiveForStudent`, `ListByStudent`), `CouponRepository`
(`GetByCode`, `IncrementUse`), `SubscriptionRepository`,
`WebhookEventRepository` (`InsertIfNew(provider, eventID, type, payload) (bool, error)`,
`MarkProcessed`).

sqlc (`db/query/payments.sql`, `enrollments.sql`, `coupons.sql`,
`subscriptions.sql`, `webhook_events.sql`) — note the **dedup insert**:

```sql
-- name: InsertWebhookIfNew :one
INSERT INTO payments.webhook_events (id, provider, provider_event_id, event_type, payload)
VALUES (gen_random_uuid(), $1, $2, $3, $4)
ON CONFLICT (provider, provider_event_id) DO NOTHING
RETURNING id;        -- no row returned => duplicate, skip processing

-- name: GetPaymentByIdempotencyKey :one
SELECT * FROM payments.payments WHERE idempotency_key = $1;

-- name: UpdatePaymentStatusGuarded :execrows
UPDATE payments.payments
SET status=@to_status, provider_transaction_id=COALESCE(@provider_txn, provider_transaction_id),
    paid_at=CASE WHEN @to_status='completed' THEN now() ELSE paid_at END,
    failure_reason=@failure_reason, updated_at=now()
WHERE id=@id AND status=@from_status;
```

Implement the Postgres repos; compile-time assertions for each interface.

### 1.5 — Idempotency helper (`pkg/idempotency` or inside usecase)

A small helper that, given an `Idempotency-Key`, returns an existing payment if one
already exists for that key (so a retried request returns the **same** checkout,
never a second charge). The UNIQUE index is the hard guarantee; this is the
friendly path.

### 1.6 — Tests

State-machine transitions (legal/illegal); `InsertWebhookIfNew` returns false on
duplicate; coupon `Apply` (percentage, fixed, expired, over-max-use, wrong scope);
idempotency-key lookup returns the existing payment.

---

## Definition of Done

- [ ] `make migrate-up` applies `0012`; down reverses; UNIQUE indexes on
      `idempotency_key` and `(provider, provider_transaction_id)` exist.
- [ ] `make sqlc && go build ./...` clean; `golangci-lint run` clean.
- [ ] Money uses `decimal.Decimal`, never `float64` (grep confirms).
- [ ] Entities have no db/json tags; provider port + registry compile.
- [ ] `InsertWebhookIfNew` is dedup-safe (test); guarded payment transition works.
- [ ] Coupon `Apply` enforces expiry/uses/scope (tests).

## Verification commands

```bash
make migrate-up && make sqlc && go build ./... && golangci-lint run
grep -rn 'float64' internal/domain/entity/payment.go internal/usecase/payment 2>/dev/null && echo "FAIL: float money" || echo OK
go test ./internal/domain/entity/... ./internal/repository/postgres/... -run 'Payment|Coupon|Webhook' -v
```

## Hand-off to Skill 2

Data + ports ready. Skill 2 builds the **enrollment engine** (grant/check access,
coupons) and **replaces every `TODO(phase4)` hook** left in Phases 2–3 with real
enrollment checks.

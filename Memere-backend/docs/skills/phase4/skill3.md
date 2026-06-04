# Phase 4 · Skill 3 — Payment Flow, Webhooks & Fulfillment

> **Prerequisite:** Phase 4 Skills 1–2 done (data + provider port; enrollment +
> access engine). Read [`docs/skill.md`](../../skill.md) §2 — *idempotency keys*
> and *webhook deduplication* are the law of this skill.
>
> **Spec references:** `memere_Design_Specification.md` §10.2 (payment sequence
> diagram), §7.3 (payment replay → idempotency + webhook dedup), §10.3 (coupons),
> README "Payment Endpoints".

---

## Goal

Implement the **payment lifecycle** as usecases (HTTP in Skill 5): initiate an
idempotent checkout, handle the provider webhook (verified + deduplicated), and on
success **grant enrollment exactly once** in the same transaction that completes
the payment. Provider impls are stubbed/sandbox here behind the Skill 1 port.

The §10.2 flow: `initiate → INSERT payment(pending, idempotency_key) → provider
checkout → redirect → webhook(verify) → UPDATE payment(completed) + INSERT
enrollment → notify`.

---

## Tasks

### 3.1 — Provider implementations (`internal/infrastructure/payment/`)

One file per provider implementing `service.PaymentProvider`. **Chapa is primary;
Telebirr and Stripe can be sandbox/stub** for Phase 4 (real keys later), but all
three must satisfy the interface and verify webhook signatures correctly.

```go
// chapa.go
type ChapaProvider struct{ secretKey string; httpc *http.Client; baseURL string }
func (c *ChapaProvider) Name() string { return "chapa" }

func (c *ChapaProvider) CreateCheckout(ctx context.Context, req service.CheckoutRequest) (*service.CheckoutResult, error) {
    // POST {baseURL}/transaction/initialize with tx_ref=req.PaymentID, amount,
    // currency, email, callback_url, return_url; Bearer c.secretKey.
    // Parse checkout_url + reference. Return CheckoutResult.
    ...
}

func (c *ChapaProvider) VerifyWebhook(ctx context.Context, h map[string]string, body []byte) (*service.WebhookEvent, error) {
    // Validate the 'Chapa-Signature' HMAC over body with the webhook secret.
    // On mismatch -> apperror.Unauthorized("BAD_SIGNATURE"). Parse normalized event.
    ...
}
```

Provide a `Registry` (`internal/infrastructure/payment/registry.go`) that maps
`entity.PaymentProvider` → impl, constructed from config.

> **Never** trust webhook *body* without signature verification (§7.3). Treat the
> signing secret like any credential — from config, never logged.

### 3.2 — Initiate payment (`internal/usecase/payment/initiate.go`)

```go
type InitiateInput struct {
    Actor          Actor
    CourseID       *uuid.UUID // nil => subscription purchase
    Plan           string     // for subscription
    Provider       entity.PaymentProvider
    CouponCode     *string
    IdempotencyKey string     // from the Idempotency-Key header (required)
}
type InitiateResult struct {
    PaymentID   string `json:"payment_id"`
    RedirectURL string `json:"redirect_url"`
    Amount      string `json:"amount"`
    Currency    string `json:"currency"`
}

func (s *Service) Initiate(ctx context.Context, in InitiateInput) (*InitiateResult, error) {
    if in.IdempotencyKey == "" {
        return nil, apperror.BadRequest("IDEMPOTENCY_KEY_REQUIRED", "")
    }
    // 1. Idempotent replay: if a payment already exists for this key, return its
    //    existing checkout (re-sign redirect if needed) — never create a second.
    if existing, _ := s.payments.GetByIdempotencyKey(ctx, in.IdempotencyKey); existing != nil {
        return s.resultFrom(existing), nil
    }
    // 2. Resolve price: course price (must not be free) or plan price.
    base, currency := s.priceFor(ctx, in)
    final, couponID := base, (*uuid.UUID)(nil)
    if in.CouponCode != nil {
        q, err := s.coupons.Quote(ctx, *in.CouponCode, in.CourseID, base)
        if err != nil { return nil, err }
        final, couponID = q.Final, &q.CouponID
    }
    // 3. Guard: already enrolled? -> 409 ALREADY_ENROLLED (no charge).
    if in.CourseID != nil {
        if ok, _ := s.enroll.Exists(ctx, in.Actor.UserID, *in.CourseID); ok {
            return nil, apperror.Conflict("ALREADY_ENROLLED", "")
        }
    }
    // 4. INSERT payment(pending) with idempotency_key + couponID.
    p := &entity.Payment{ ID: uuid.New(), StudentID: in.Actor.UserID, CourseID: in.CourseID,
        Amount: final, Currency: currency, Provider: in.Provider, Status: entity.PayPending,
        IdempotencyKey: &in.IdempotencyKey, CouponID: couponID }
    if err := s.payments.Create(ctx, p); err != nil { return nil, err } // UNIQUE key race -> retry-as-replay
    // 5. provider.CreateCheckout(...) ; persist provider_checkout_id.
    prov, _ := s.registry.Get(in.Provider)
    co, err := prov.CreateCheckout(ctx, s.checkoutReq(p, in.Actor))
    if err != nil { _ = s.payments.UpdateStatusGuarded(ctx, p.ID, entity.PayPending, entity.PayFailed, /*reason*/); return nil, err }
    _ = s.payments.SetCheckout(ctx, p.ID, co.CheckoutID)
    return &InitiateResult{PaymentID: p.ID.String(), RedirectURL: co.RedirectURL, Amount: final.String(), Currency: currency}, nil
}
```

Handle the UNIQUE-violation race on `idempotency_key` (two concurrent identical
requests) by re-reading and returning the existing payment.

### 3.3 — Webhook handling + fulfillment (`internal/usecase/payment/webhook.go`)

```go
func (s *Service) HandleWebhook(ctx context.Context, provider entity.PaymentProvider, headers map[string]string, body []byte) error {
    prov, err := s.registry.Get(provider)
    if err != nil { return err }
    ev, err := prov.VerifyWebhook(ctx, headers, body) // signature-checked
    if err != nil { return err }                       // bad sig -> 401

    // 1. DEDUP: insert into webhook_events; if duplicate -> ack & return (no reprocess).
    isNew, err := s.webhooks.InsertIfNew(ctx, string(provider), ev.ProviderEventID, ev.Type, ev.RawPayload)
    if err != nil { return err }
    if !isNew { return nil } // already processed; idempotent ack

    // 2. Resolve our payment by ref (tx_ref == PaymentID, or provider_txn).
    p, err := s.resolvePayment(ctx, ev)
    if err != nil { return err }

    // 3. Fulfill in ONE transaction (guarded):
    switch ev.Status {
    case "completed":
        err = s.tx.Do(ctx, func(q repository.Tx) error {
            ok, err := q.Payments().UpdateStatusGuarded(ctx, p.ID, entity.PayPending, entity.PayCompleted, withTxn(ev))
            if err != nil { return err }
            if !ok { return nil } // already completed -> idempotent no-op
            if p.CouponID != nil { _ = q.Coupons().IncrementUse(ctx, *p.CouponID) } // burn coupon ONLY on success
            if p.CourseID != nil {
                return q.Enrollments().CreateIfNotExists(ctx, p.StudentID, *p.CourseID, entity.SourcePurchase, nil)
            }
            return q.Subscriptions().Activate(ctx, p) // subscription path
        })
    case "failed":
        _, err = s.payments.UpdateStatusGuarded(ctx, p.ID, entity.PayPending, entity.PayFailed, withReason(ev))
    }
    if err != nil { return err }
    _ = s.webhooks.MarkProcessed(ctx, provider, ev.ProviderEventID)
    // 4. fire enrollment-confirmed notification (no-op hook now; Phase 5 wires it)
    s.notify.PurchaseConfirmed(ctx, p)
    return nil
}
```

Critical properties (all from §7.3 / §10.2):
- **Webhook dedup** via `InsertIfNew` — a provider re-delivery never double-grants.
- **Guarded completion** — `pending→completed` flips once; a racing poll/webhook
  no-ops.
- **Coupon usage incremented only on success**, inside the same tx.
- **Enrollment created idempotently** (`CreateIfNotExists`).
- The whole fulfillment is **one DB transaction** — payment, coupon, enrollment
  commit together or not at all. Add a `repository.Tx` / `WithTx` helper if Phase
  1 didn't.

### 3.4 — Status polling (advisory) + refund

- `GetPaymentStatus(ctx, actor, paymentID)` — ownership-checked; returns current
  status. May call `provider.VerifyPayment` as a fallback if still `pending` past
  a grace window, then run the same guarded fulfillment (so polling can complete a
  payment if the webhook was lost). Reuses the guarded/transactional path — **no
  duplicate enrollment**.
- `RefundPayment(ctx, actor(admin), paymentID)` — admin-only; `completed→refunded`;
  optionally revoke enrollment (decision: keep enrollment but mark refunded — log
  it; document). Provider refund call is stubbed if not supported.

### 3.5 — Tests (mocked provider + repos + fake clock)

- Initiate: missing idempotency key → 400; replay with same key → same payment, no
  second insert; already-enrolled → 409; coupon applied to amount.
- Webhook: bad signature → 401; duplicate event → single fulfillment; success →
  payment completed + enrollment created + coupon incremented (all-or-nothing tx);
  failure → payment failed, no enrollment.
- Race: webhook + poll both observe success → exactly one enrollment, one coupon
  increment.
- Refund: admin only; non-admin → 403.

---

## Definition of Done

- [ ] `go build ./...` clean; `golangci-lint run` clean.
- [ ] `go test ./internal/usecase/payment/... ./internal/infrastructure/payment/...`
      passes including the dedup + race tests.
- [ ] All three providers satisfy `service.PaymentProvider`; webhook signature
      verification is enforced (bad sig → 401).
- [ ] Idempotency key required; replays return the same payment (no double-charge).
- [ ] Fulfillment is transactional and idempotent: success grants enrollment
      exactly once and burns the coupon exactly once, even under
      webhook-redelivery or webhook+poll races.
- [ ] No signing secret or raw card data is ever logged.

## Verification commands

```bash
go build ./... && golangci-lint run
go test ./internal/usecase/payment/... ./internal/infrastructure/payment/... -v
grep -rn 'SecretKey\|secretKey\|Signature' internal/ | grep -i 'log\.' && echo "REVIEW: secret near logging" || echo OK
```

## Hand-off to Skill 4

Purchases work end-to-end (mocked providers). Skill 4 adds **subscriptions**
(monthly/annual plans, renewal/expiry sweeper) and the **teacher earnings /
revenue** reporting from §10 + §1.6 (70/30 split).

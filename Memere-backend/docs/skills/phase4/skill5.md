# Phase 4 · Skill 5 — HTTP Delivery & Wiring (Payments, Enrollments, Subscriptions)

> **Prerequisite:** Phase 4 Skills 1–4 done. Read [`docs/skill.md`](../../skill.md)
> §2 and reuse the Phase 1–3 delivery layer.
>
> **Spec references:** `memere_Design_Specification.md` §10.2 (flow), §5.4 (API
> conventions), §7.3 (idempotency, webhook verification), README "Payment
> Endpoints".

---

## Goal

Expose the payment, enrollment, subscription, and revenue functionality over the
Gin API; mount **signature-verified, raw-body webhook endpoints**; and wire
providers + the subscription sweeper into the app. After this skill, Phase 4 is
**fully runnable end-to-end** with sandbox/mocked providers.

The webhook route is special: it needs the **raw request body** (for signature
verification) and must be **outside** the JSON-binding and auth middleware.

---

## API surface for Phase 4

Base path `/api/v1`.

**Payments / purchase** (`payment_handler.go`)
| Method | Path | Auth | Notes |
|---|---|---|---|
| POST | `/payments/initiate` | bearer | requires `Idempotency-Key` header |
| GET | `/payments/:id/status` | bearer (owner) | advisory poll |
| GET | `/payments` | bearer | my payment history |
| POST | `/payments/:id/refund` | bearer + admin | |
| POST | `/webhooks/payments/:provider` | **public, raw body, signature-verified** | provider callback |

**Enrollment**
| Method | Path | Auth |
|---|---|---|
| POST | `/courses/:id/enroll-free` | bearer (free courses only) |
| GET | `/me/enrollments` | bearer |

**Subscriptions**
| Method | Path | Auth |
|---|---|---|
| GET | `/subscription-plans` | public |
| POST | `/subscriptions` | bearer (initiates payment, plan in body) |
| GET | `/me/subscription` | bearer |
| POST | `/subscriptions/:id/cancel` | bearer (owner) |

**Revenue** (`revenue_handler.go`)
| Method | Path | Auth |
|---|---|---|
| GET | `/admin/revenue` | bearer + admin |
| GET | `/me/earnings` | bearer + teacher/admin |
| GET | `/courses/:id/sales` | bearer + owner/admin |

---

## Tasks

### 5.1 — DTOs (`internal/delivery/http/dto/payment.go`)

```go
type InitiatePaymentRequest struct {
    CourseID   *string `json:"course_id"`            // null for subscription
    Plan       string  `json:"plan,omitempty"`       // monthly|annual
    Provider   string  `json:"provider" binding:"required"` // chapa|telebirr|stripe
    CouponCode *string `json:"coupon_code,omitempty"`
}
type InitiatePaymentResponse struct {
    PaymentID   string `json:"payment_id"`
    RedirectURL string `json:"redirect_url"`
    Amount      string `json:"amount"`
    Currency    string `json:"currency"`
}
type PaymentStatusResponse struct {
    PaymentID string `json:"payment_id"`
    Status    string `json:"status"`
    Amount    string `json:"amount"`
    Currency  string `json:"currency"`
    PaidAt    *string `json:"paid_at,omitempty"`
}
```

Plus enrollment, subscription, revenue response DTOs. **Never** serialize
`metadata` blobs containing provider raw payloads to the client.

### 5.2 — Webhook handler (raw body, no auth) — get this right

```go
// payment_handler.go
func (h *PaymentHandler) Webhook(c *gin.Context) {
    provider := entity.PaymentProvider(c.Param("provider"))
    body, err := io.ReadAll(io.LimitReader(c.Request.Body, 1<<20)) // 1MB cap
    if err != nil { respondError(c, apperror.BadRequest("BAD_BODY","")); return }

    headers := map[string]string{}
    for k := range c.Request.Header { headers[k] = c.GetHeader(k) }

    if err := h.uc.HandleWebhook(c, provider, headers, body); err != nil {
        // bad signature -> 401; unknown payment -> 404; else 500.
        respondError(c, err); return
    }
    c.Status(http.StatusOK) // always 200 on successful (idempotent) processing
}
```

Register it **before/outside** the JSON + auth middleware so the raw body is
intact and unauthenticated providers can reach it. Apply only: recovery,
request-id, logger (with body redaction), and a dedicated rate limit.

### 5.3 — Initiate handler (idempotency header)

```go
func (h *PaymentHandler) Initiate(c *gin.Context) {
    actor, _ := actorFromContext(c)
    idem := c.GetHeader("Idempotency-Key")
    if idem == "" { respondError(c, apperror.BadRequest("IDEMPOTENCY_KEY_REQUIRED","")); return }
    var req dto.InitiatePaymentRequest
    if err := c.ShouldBindJSON(&req); err != nil { respondError(c, apperror.Validation(bindingDetails(err))); return }
    out, err := h.uc.Initiate(c, payment.InitiateInput{
        Actor: actor, Provider: entity.PaymentProvider(req.Provider),
        CourseID: parseUUIDPtr(req.CourseID), Plan: req.Plan,
        CouponCode: req.CouponCode, IdempotencyKey: idem,
    })
    if err != nil { respondError(c, err); return }
    respondJSON(c, http.StatusCreated, toInitiateResponse(out))
}
```

### 5.4 — Routes & RBAC

Register all routes above. Webhook group is public + raw; everything else uses
`RequireAuth`, with `RequireRole(admin)` on refund/revenue-admin and
`RequireRole(teacher, admin)` on earnings. Ownership checks stay in usecases.

### 5.5 — Wire providers + sweeper into `main.go`

```go
// build provider registry from config
reg := payment.NewRegistry(
    payment.NewChapaProvider(cfg.Chapa),
    payment.NewTelebirrProvider(cfg.Telebirr),
    payment.NewStripeProvider(cfg.Stripe),
)
payRepo := postgres.NewPaymentRepo(pool, q)
enrollRepo := postgres.NewEnrollmentRepo(pool, q)
subRepo := postgres.NewSubscriptionRepo(pool, q)
whRepo := postgres.NewWebhookRepo(pool, q)
tx := postgres.NewTxRunner(pool, q)

accessSvc := access.NewService(enrollRepo, subRepo, courseRepo, clock)   // replaces Skill 2 stub deps
payUC := payment.NewService(payRepo, enrollRepo, subRepo, whRepo, couponRepo, reg, tx, notifyNoop, clock, cfg)
subSweeper := worker.NewSubscriptionSweeper(subRepo, payUC, notifyNoop, clock, cfg.SubSweep())
go subSweeper.Run(ctx)

// IMPORTANT: re-inject accessSvc into quiz/exam/video services (Skill 2 wiring)
```

Confirm the quiz/exam/video services constructed in earlier phases now receive the
real `access.Service` (the Phase 2/3 `TODO(phase4)` replacements depend on it).

### 5.6 — End-to-end smoke test (`scripts/smoke_phase4.sh`)

With a **mock/sandbox provider** (add a `mock` provider that auto-fires a webhook,
or expose a test-only endpoint to simulate the callback):

1. teacher creates a **paid** course + a video lesson + a quiz (Phases 1–3).
2. student tries to stream the paid video → `403 NOT_ENROLLED`; tries the quiz →
   `403 NOT_ENROLLED` (proves Skill 2 hooks live).
3. student `POST /payments/initiate` (course, provider=mock) **with**
   `Idempotency-Key` → `redirect_url` + `payment_id`.
4. repeat the **same** request with the **same** key → same `payment_id` (no
   double-charge).
5. simulate the provider webhook (success) → payment `completed`, enrollment
   created.
6. re-deliver the **same** webhook → still one enrollment (dedup proven).
7. student now streams the video (200) and takes the quiz (200).
8. apply a coupon on a second purchase → discounted amount; coupon `used_count`
   increments only after success.
9. subscribe (plan=monthly) → on webhook success, student gets all-access; expire
   it via the sweeper (short period) → access revoked.
10. `GET /me/earnings` as teacher → gross × 0.70.

### 5.7 — Docs

Update README "Payment Integration"/"Getting Started"; document the
`Idempotency-Key` requirement and webhook URLs. Optionally extend `api/` OpenAPI.

---

## Definition of Done

- [ ] `go build ./...` clean; `golangci-lint run` clean; `go test ./...` passes.
- [ ] `make up && make migrate-up && make run` serves Phase 4 endpoints; the
      subscription sweeper starts/stops with the app.
- [ ] `scripts/smoke_phase4.sh` passes every assertion, including the
      idempotent-initiate, webhook-dedup, and access-gating steps.
- [ ] The webhook route reads the **raw body**, verifies the signature, is
      unauthenticated, and returns 200 on idempotent (re)processing.
- [ ] Paid content (video + quiz/exam) is correctly gated by enrollment; the
      Phase 2/3 hooks now resolve through the real `access.Service`.
- [ ] No double-charge (idempotency), no double-grant (webhook dedup + guarded tx).
- [ ] No provider secret or raw payload is logged or returned to clients.

## Verification commands

```bash
make up && make migrate-up
make run &
bash scripts/smoke_phase4.sh
go test ./... && golangci-lint run
grep -rn 'TODO(phase4)' internal/ && echo "FAIL: hooks remain" || echo "OK"
```

---

## 🎉 Phase 4 complete — what now

Payments, enrollments, subscriptions, coupons, and revenue reporting are live —
idempotent, webhook-deduplicated, transactionally fulfilled — and every
`TODO(phase4)` access hook from Phases 2–3 is now backed by real enrollment.

**To proceed to Phase 5:**
1. Report "Phase 4 complete" with a summary + smoke-test result.
2. Ask Claude to author **Phase 5** (it is already authored in
   `docs/skills/phase5/`): **Progress tracking, Notifications (FCM/SendGrid), and
   Admin analytics** (spec §11, §4.2.8, §9.3 trends, README dashboards).
3. Phase 5 wires the `notify` hooks left as no-ops in Phases 3–4.

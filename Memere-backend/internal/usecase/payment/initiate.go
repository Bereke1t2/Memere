package payment

import (
	"context"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// InitiateInput is the decoded request to start a checkout. IdempotencyKey comes
// from the client Idempotency-Key header and is required (no double-charge).
// CourseID nil means a subscription purchase (deferred to Skill 4).
type InitiateInput struct {
	Actor          Actor
	CourseID       *uuid.UUID
	Plan           string
	Provider       entity.PaymentProvider
	CouponCode     *string
	IdempotencyKey string
}

// InitiateResult is the client response: the redirect URL to open plus the
// resolved (possibly coupon-discounted) amount.
type InitiateResult struct {
	PaymentID   string `json:"payment_id"`
	RedirectURL string `json:"redirect_url"`
	Amount      string `json:"amount"`
	Currency    string `json:"currency"`
}

// Initiate starts an idempotent checkout (spec §10.2 step 1). A retried request
// with the same Idempotency-Key returns the existing payment's checkout instead
// of charging again. The flow: replay-check → resolve price (+ coupon) →
// already-enrolled guard → create pending payment → provider checkout → persist
// the redirect URL + checkout id on the row.
func (s *Service) Initiate(ctx context.Context, in InitiateInput) (*InitiateResult, error) {
	if in.Actor.UserID == uuid.Nil {
		return nil, apperror.Unauthorized("authentication required", nil)
	}
	if in.IdempotencyKey == "" {
		return nil, apperror.New(400, "IDEMPOTENCY_KEY_REQUIRED",
			"an Idempotency-Key is required to start a payment", nil)
	}
	if !in.Provider.Valid() {
		return nil, apperror.New(400, "UNKNOWN_PROVIDER", "unknown payment provider", nil)
	}
	// Resolve the provider up front so a bad provider fails before any write.
	prov, err := s.registry.Get(in.Provider)
	if err != nil {
		return nil, err
	}

	// 1. Idempotent replay: a payment already created for this key returns its
	//    existing checkout — never a second charge.
	if existing, err := s.payments.GetByIdempotencyKey(ctx, in.IdempotencyKey); err == nil {
		return s.resultFrom(existing), nil
	} else if !apperror.IsNotFound(err) {
		return nil, err
	}

	// 2. Resolve price and (optionally) apply a coupon. Quote is read-only — the
	//    coupon's use count is incremented only when the payment completes.
	base, currency, err := s.priceFor(ctx, in)
	if err != nil {
		return nil, err
	}
	final := base
	var couponID *uuid.UUID
	if in.CouponCode != nil && in.CourseID != nil {
		f, cid, err := s.quoter.Quote(ctx, *in.CouponCode, *in.CourseID, base)
		if err != nil {
			return nil, err
		}
		final, couponID = f, &cid
	}

	// 3. Guard: already enrolled? Never charge for access the student owns.
	if in.CourseID != nil {
		if ok, err := s.enroll.Exists(ctx, in.Actor.UserID, *in.CourseID); err != nil {
			return nil, err
		} else if ok {
			return nil, apperror.New(409, "ALREADY_ENROLLED",
				"you already have access to this course", nil)
		}
	}

	// 4. Create the pending payment, claiming the idempotency key (UNIQUE).
	key := in.IdempotencyKey
	p := &entity.Payment{
		ID:             uuid.New(),
		StudentID:      in.Actor.UserID,
		CourseID:       in.CourseID,
		Amount:         final,
		Currency:       currency,
		Provider:       in.Provider,
		Status:         entity.PayPending,
		IdempotencyKey: &key,
		CouponID:       couponID,
	}

	// 5. Start the hosted checkout (tx_ref == PaymentID). Stash the redirect URL
	//    + checkout id on the row so a replay can answer without re-calling the
	//    provider.
	co, err := prov.CreateCheckout(ctx, s.checkoutReq(p))
	if err != nil {
		return nil, err
	}
	p.ProviderCheckoutID = &co.CheckoutID
	p.Metadata = map[string]any{metaRedirectURL: co.RedirectURL}
	if in.CourseID == nil {
		// Subscription purchase: carry the plan to fulfillment (no plan column).
		p.Metadata[metaPlan] = in.Plan
	}

	if err := s.payments.Create(ctx, p); err != nil {
		// A concurrent identical initiate won the idempotency-key race: re-read and
		// return that payment rather than surfacing the conflict.
		if apperror.IsCode(err, "PAYMENT_DUPLICATE") {
			if existing, gerr := s.payments.GetByIdempotencyKey(ctx, in.IdempotencyKey); gerr == nil {
				return s.resultFrom(existing), nil
			}
		}
		return nil, err
	}

	return &InitiateResult{
		PaymentID:   p.ID.String(),
		RedirectURL: co.RedirectURL,
		Amount:      final.String(),
		Currency:    currency,
	}, nil
}

// resultFrom builds an InitiateResult from an existing payment (the replay path).
func (s *Service) resultFrom(p *entity.Payment) *InitiateResult {
	return &InitiateResult{
		PaymentID:   p.ID.String(),
		RedirectURL: redirectURLOf(p),
		Amount:      p.Amount.String(),
		Currency:    p.Currency,
	}
}

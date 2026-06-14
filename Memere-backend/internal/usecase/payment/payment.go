// Package payment holds the payment-lifecycle usecases (spec §10.2): initiating
// an idempotent checkout, handling the provider webhook (signature-verified and
// deduplicated), and — on success — granting enrollment exactly once inside the
// same transaction that completes the payment. Like the other usecase packages
// it is pure orchestration over domain repositories and ports; the concrete
// Chapa/Telebirr/Stripe clients live in the infrastructure ring behind
// service.PaymentProvider (dependency rule).
//
// Two non-negotiables are enforced here and must never be relaxed: every
// initiate carries a client Idempotency-Key so a retry never double-charges, and
// fulfillment is one guarded transaction so a re-delivered webhook (or a polling
// fallback racing the webhook) never double-grants or double-burns a coupon
// (spec §7.3).
package payment

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/service"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// metaRedirectURL is the Metadata key under which we stash the provider redirect
// URL at create time, so an idempotent replay of initiate can return the same
// checkout without re-calling the provider.
const metaRedirectURL = "redirect_url"

// metaPlan is the Metadata key under which a subscription purchase carries its
// chosen plan from initiate to fulfillment (the Payment entity has no plan
// column). It must match the subscription usecase's key.
const metaPlan = "plan"

// Actor is the authenticated caller's identity. A zero Actor (UserID == Nil) is
// anonymous and may not initiate a payment.
type Actor struct {
	UserID uuid.UUID
	Role   entity.Role
}

// CouponQuoter validates a coupon against a course purchase and returns the
// discounted price plus the coupon id (read-only; usage is incremented only at
// fulfillment). *coupon.Service satisfies it.
type CouponQuoter interface {
	Quote(ctx context.Context, code string, courseID uuid.UUID, basePrice decimal.Decimal) (final decimal.Decimal, couponID uuid.UUID, err error)
}

// PlanPricer resolves a subscription plan's price + currency for a subscription
// checkout. *subscription.Service satisfies it. It may be nil, in which case
// subscription purchases are rejected (deferred wiring).
type PlanPricer interface {
	PriceForPlan(plan string) (decimal.Decimal, string, error)
}

// SubscriptionActivator creates/extends a subscription from a settled payment.
// *subscription.Service satisfies it. It may be nil, in which case a settled
// subscription payment simply completes without activation (deferred wiring).
type SubscriptionActivator interface {
	Activate(ctx context.Context, p *entity.Payment) error
}

// Notifier fires post-fulfillment notifications. It is a no-op hook in Phase 4;
// Phase 5 wires it to the real notification system. PurchaseConfirmed must never
// fail the fulfillment — it runs after the transaction commits.
type Notifier interface {
	PurchaseConfirmed(ctx context.Context, p *entity.Payment)
}

// NoopNotifier is the default Notifier used until Phase 5.
type NoopNotifier struct{}

// PurchaseConfirmed does nothing.
func (NoopNotifier) PurchaseConfirmed(context.Context, *entity.Payment) {}

// Config carries the URLs the usecase hands the provider when creating a
// checkout. main.go maps the typed config onto this so the usecase stays free of
// the config package (dependency rule).
type Config struct {
	// CallbackURL is our webhook endpoint the provider posts settlement to.
	CallbackURL string
	// ReturnURL is where the provider redirects the user after paying.
	ReturnURL string
	// DefaultCurrency is used when a course row carries no currency.
	DefaultCurrency string
}

// Service implements the payment-lifecycle usecases over the domain
// repositories, the provider registry, and the coupon quoter. now is injectable
// so tests drive coupon expiry / timestamps with a fake clock.
type Service struct {
	payments repository.PaymentRepository
	enroll   repository.EnrollmentRepository
	coupons  repository.CouponRepository
	webhooks repository.WebhookEventRepository
	courses  repository.CourseRepository
	registry service.PaymentProviderRegistry
	quoter   CouponQuoter
	pricer   PlanPricer
	activate SubscriptionActivator
	tx       repository.TxManager
	notify   Notifier
	cfg      Config
	now      func() time.Time
}

// NewService wires the payment usecase. notify may be nil (defaults to a no-op);
// clock may be nil (defaults to time.Now). pricer/activator may be nil until the
// subscription usecase is wired (Skill 4), in which case subscription purchases
// are rejected and a settled subscription payment completes without activation.
func NewService(
	payments repository.PaymentRepository,
	enroll repository.EnrollmentRepository,
	coupons repository.CouponRepository,
	webhooks repository.WebhookEventRepository,
	courses repository.CourseRepository,
	registry service.PaymentProviderRegistry,
	quoter CouponQuoter,
	pricer PlanPricer,
	activator SubscriptionActivator,
	tx repository.TxManager,
	notify Notifier,
	cfg Config,
	clock func() time.Time,
) *Service {
	if notify == nil {
		notify = NoopNotifier{}
	}
	if clock == nil {
		clock = time.Now
	}
	if cfg.DefaultCurrency == "" {
		cfg.DefaultCurrency = "ETB"
	}
	return &Service{
		payments: payments,
		enroll:   enroll,
		coupons:  coupons,
		webhooks: webhooks,
		courses:  courses,
		registry: registry,
		quoter:   quoter,
		pricer:   pricer,
		activate: activator,
		tx:       tx,
		notify:   notify,
		cfg:      cfg,
		now:      clock,
	}
}

// PaymentView is the client-safe projection of a payment (no provider secrets).
type PaymentView struct {
	PaymentID string `json:"payment_id"`
	Status    string `json:"status"`
	Amount    string `json:"amount"`
	Currency  string `json:"currency"`
	CourseID  string `json:"course_id,omitempty"`
}

// viewOf projects a payment for a client response.
func viewOf(p *entity.Payment) *PaymentView {
	v := &PaymentView{
		PaymentID: p.ID.String(),
		Status:    string(p.Status),
		Amount:    p.Amount.String(),
		Currency:  p.Currency,
	}
	if p.CourseID != nil {
		v.CourseID = p.CourseID.String()
	}
	return v
}

// priceFor resolves the amount + currency to charge. Course purchases read the
// (server-side) course price and refuse a free course; subscription purchases
// (CourseID == nil) are deferred to Skill 4.
func (s *Service) priceFor(ctx context.Context, in InitiateInput) (decimal.Decimal, string, error) {
	if in.CourseID == nil {
		// Subscription purchase: price the plan via the injected pricer.
		if s.pricer == nil {
			return decimal.Zero, "", apperror.New(400, "SUBSCRIPTION_NOT_AVAILABLE",
				"subscription purchases are not available", nil)
		}
		amount, currency, err := s.pricer.PriceForPlan(in.Plan)
		if err != nil {
			return decimal.Zero, "", err
		}
		if currency == "" {
			currency = s.cfg.DefaultCurrency
		}
		return amount, currency, nil
	}
	c, err := s.courses.FindByID(ctx, *in.CourseID)
	if err != nil {
		return decimal.Zero, "", err
	}
	if c.IsFree {
		return decimal.Zero, "", apperror.New(409, "COURSE_IS_FREE",
			"this course is free; enrol directly instead of paying", nil)
	}
	currency := c.Currency
	if currency == "" {
		currency = s.cfg.DefaultCurrency
	}
	return decimal.NewFromFloat(c.Price), currency, nil
}

// checkoutReq builds the provider-agnostic checkout request from a pending
// payment. tx_ref is the PaymentID, so the webhook maps the settlement straight
// back to our row.
func (s *Service) checkoutReq(p *entity.Payment) service.CheckoutRequest {
	desc := "Course purchase"
	if p.CourseID != nil {
		desc = "Course purchase " + p.CourseID.String()
	}
	key := ""
	if p.IdempotencyKey != nil {
		key = *p.IdempotencyKey
	}
	return service.CheckoutRequest{
		PaymentID:      p.ID.String(),
		Amount:         p.Amount,
		Currency:       p.Currency,
		Description:    desc,
		ReturnURL:      s.cfg.ReturnURL,
		CallbackURL:    s.cfg.CallbackURL,
		IdempotencyKey: key,
	}
}

// redirectURLOf reads the stashed provider redirect URL from a payment's
// metadata (set at create time), used to answer an idempotent replay.
func redirectURLOf(p *entity.Payment) string {
	if p.Metadata == nil {
		return ""
	}
	if u, ok := p.Metadata[metaRedirectURL].(string); ok {
		return u
	}
	return ""
}

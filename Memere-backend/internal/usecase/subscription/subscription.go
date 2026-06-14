// Package subscription holds the all-access subscription usecases (spec §1.6,
// §10.2): config-driven plan pricing, activation (the fulfillment primitive the
// payment flow calls on a successful subscription purchase), and cancellation.
// Like the other usecase packages it is pure orchestration over the domain
// repositories — no HTTP, no SQL.
//
// Activation is idempotent and EXTENDS an existing active subscription rather
// than creating a duplicate, so a renewal or a second purchase stacks the period
// instead of opening a parallel row.
package subscription

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// metaPlan is the payment Metadata key under which the chosen plan is carried
// from initiate to fulfillment (the Payment entity has no plan column).
const metaPlan = "plan"

// Actor is the authenticated caller's identity.
type Actor struct {
	UserID uuid.UUID
	Role   entity.Role
}

// PlanSpec is the price and billing period for a subscription plan. Prices are
// config-driven (spec §1.6) so they are never hardcoded in the usecase.
type PlanSpec struct {
	Plan     entity.SubscriptionPlan
	Price    decimal.Decimal
	Currency string
	Period   time.Duration
}

// Config carries the per-plan specs. main.go maps the typed config onto this so
// the usecase stays free of the config package (dependency rule).
type Config struct {
	Monthly PlanSpec
	Annual  PlanSpec
}

// Service implements the subscription usecases over the subscription repository.
// now is injectable so tests drive period math with a fake clock.
type Service struct {
	subs  repository.SubscriptionRepository
	plans map[entity.SubscriptionPlan]PlanSpec
	now   func() time.Time
}

// NewService wires the subscription service. clock may be nil (defaults to
// time.Now). Plans are taken from cfg, defaulting any unset period to its
// conventional length so a misconfigured period never yields a zero window.
func NewService(subs repository.SubscriptionRepository, cfg Config, clock func() time.Time) *Service {
	if clock == nil {
		clock = time.Now
	}
	monthly := cfg.Monthly
	monthly.Plan = entity.PlanMonthly
	if monthly.Period == 0 {
		monthly.Period = 30 * 24 * time.Hour
	}
	annual := cfg.Annual
	annual.Plan = entity.PlanAnnual
	if annual.Period == 0 {
		annual.Period = 365 * 24 * time.Hour
	}
	return &Service{
		subs: subs,
		plans: map[entity.SubscriptionPlan]PlanSpec{
			entity.PlanMonthly: monthly,
			entity.PlanAnnual:  annual,
		},
		now: clock,
	}
}

// ListPlans returns the configured subscription plans (monthly then annual) for
// the public plan-listing endpoint. Order is stable so the catalogue renders
// predictably.
func (s *Service) ListPlans() []PlanSpec {
	out := make([]PlanSpec, 0, len(s.plans))
	for _, p := range []entity.SubscriptionPlan{entity.PlanMonthly, entity.PlanAnnual} {
		if spec, ok := s.plans[p]; ok {
			out = append(out, spec)
		}
	}
	return out
}

// PlanSpec returns the spec for a plan, or an error for an unknown plan.
func (s *Service) PlanSpec(p entity.SubscriptionPlan) (PlanSpec, error) {
	spec, ok := s.plans[p]
	if !ok {
		return PlanSpec{}, apperror.New(400, "UNKNOWN_PLAN", "unknown subscription plan: "+string(p), nil)
	}
	return spec, nil
}

// PriceForPlan resolves a plan's price + currency. It satisfies the payment
// usecase's PlanPricer port so a subscription checkout can be priced.
func (s *Service) PriceForPlan(plan string) (decimal.Decimal, string, error) {
	spec, err := s.PlanSpec(entity.SubscriptionPlan(plan))
	if err != nil {
		return decimal.Zero, "", err
	}
	return spec.Price, spec.Currency, nil
}

// Activate creates or extends the student's all-access subscription from a
// settled payment. It satisfies the payment usecase's SubscriptionActivator
// port and is called inside the guarded fulfillment transaction, so it runs at
// most once per payment. It is also idempotent against an existing active
// subscription: rather than open a duplicate, it pushes the period end out by the
// plan's length.
func (s *Service) Activate(ctx context.Context, p *entity.Payment) error {
	plan := planFromPayment(p)
	spec, err := s.PlanSpec(plan)
	if err != nil {
		return err
	}
	now := s.now()

	existing, err := s.subs.GetActiveForStudent(ctx, p.StudentID)
	switch {
	case err == nil && existing != nil:
		// Extend from the later of (now, current end) so a renewal stacks and a
		// lapsed-but-not-yet-swept row still extends forward from now.
		base := existing.CurrentPeriodEnd
		if base.Before(now) {
			base = now
		}
		return s.subs.ExtendPeriod(ctx, existing.ID, base.Add(spec.Period))
	case apperror.IsNotFound(err):
		// No active subscription — create one.
	case err != nil:
		return err
	}

	return s.subs.Create(ctx, &entity.Subscription{
		ID:                 uuid.New(),
		StudentID:          p.StudentID,
		Plan:               plan,
		Status:             entity.SubActive,
		CurrentPeriodStart: now,
		CurrentPeriodEnd:   now.Add(spec.Period),
		Provider:           p.Provider,
	})
}

// CancelSubscription flags the caller's active subscription as canceled (it will
// not renew) while keeping access until current_period_end; the sweeper expires
// it when the period lapses. Owner-only (an admin may cancel any).
func (s *Service) CancelSubscription(ctx context.Context, actor Actor, subID uuid.UUID) error {
	if actor.UserID == uuid.Nil {
		return apperror.Unauthorized("authentication required", nil)
	}
	sub, err := s.subs.GetByID(ctx, subID)
	if err != nil {
		return err
	}
	if actor.Role != entity.RoleAdmin && sub.StudentID != actor.UserID {
		return apperror.Forbidden("you do not own this subscription", nil)
	}
	ok, err := s.subs.CancelAtPeriodEnd(ctx, subID)
	if err != nil {
		return err
	}
	if !ok {
		return apperror.New(409, "NOT_CANCELABLE",
			"subscription is not active or already canceled", nil)
	}
	return nil
}

// SweepExpired transitions active subscriptions whose period has lapsed to
// expired and returns the ones it actually expired (for the sweeper to notify
// on). It is idempotent and guarded: ExpireLapsed only flips a row that is still
// active-and-lapsed, so a row renewed between the list and the update is skipped,
// and a repeated tick is a no-op. A per-row error is collected and the sweep
// continues with the rest; the next tick retries the failed row.
func (s *Service) SweepExpired(ctx context.Context, limit int) ([]*entity.Subscription, error) {
	lapsed, err := s.subs.ListExpiring(ctx, limit)
	if err != nil {
		return nil, err
	}
	expired := make([]*entity.Subscription, 0, len(lapsed))
	var firstErr error
	for _, sub := range lapsed {
		ok, err := s.subs.ExpireLapsed(ctx, sub.ID)
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		if ok {
			sub.Status = entity.SubExpired
			expired = append(expired, sub)
		}
	}
	return expired, firstErr
}

// GetMySubscription returns the caller's active subscription, or apperror.NotFound.
func (s *Service) GetMySubscription(ctx context.Context, actor Actor) (*entity.Subscription, error) {
	if actor.UserID == uuid.Nil {
		return nil, apperror.Unauthorized("authentication required", nil)
	}
	return s.subs.GetActiveForStudent(ctx, actor.UserID)
}

// planFromPayment reads the plan a payment was made for from its metadata,
// defaulting to monthly when absent (a defensive fallback; initiate always sets
// it for a subscription purchase).
func planFromPayment(p *entity.Payment) entity.SubscriptionPlan {
	if p.Metadata != nil {
		if v, ok := p.Metadata[metaPlan].(string); ok && v != "" {
			return entity.SubscriptionPlan(v)
		}
	}
	return entity.PlanMonthly
}

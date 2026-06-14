package worker

import (
	"context"
	"log"
	"time"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
)

// SubscriptionExpirer transitions lapsed (active, period-ended) subscriptions to
// expired and returns the ones it expired. *subscription.Service satisfies it via
// SweepExpired. The transition is guarded + idempotent, so the sweeper is safe to
// run repeatedly (spec §10.2).
type SubscriptionExpirer interface {
	SweepExpired(ctx context.Context, limit int) ([]*entity.Subscription, error)
}

// SubscriptionNotifier fires the "your subscription expired" notification. It is a
// no-op hook in Phase 4; Phase 5 wires it to the real notification system.
type SubscriptionNotifier interface {
	SubscriptionExpired(ctx context.Context, sub *entity.Subscription)
}

// NoopSubscriptionNotifier is the default notifier used until Phase 5.
type NoopSubscriptionNotifier struct{}

// SubscriptionExpired does nothing.
func (NoopSubscriptionNotifier) SubscriptionExpired(context.Context, *entity.Subscription) {}

// SubscriptionSweeper expires lapsed subscriptions on a ticker until its context
// is cancelled (wired to the app's shutdown signal in main.go). It mirrors the
// AttemptSweeper pattern: one immediate sweep on start, then every interval.
//
// Phase 4 providers (Chapa/Telebirr) are one-off, so a lapsed subscription simply
// expires; auto-renew for recurring providers is a later phase.
type SubscriptionSweeper struct {
	expirer  SubscriptionExpirer
	notify   SubscriptionNotifier
	interval time.Duration
	limit    int
}

// NewSubscriptionSweeper builds a sweeper that ticks every interval. notify may be
// nil (defaults to a no-op).
func NewSubscriptionSweeper(interval time.Duration, expirer SubscriptionExpirer, notify SubscriptionNotifier) *SubscriptionSweeper {
	if notify == nil {
		notify = NoopSubscriptionNotifier{}
	}
	return &SubscriptionSweeper{
		expirer:  expirer,
		notify:   notify,
		interval: interval,
		limit:    defaultBatchLimit,
	}
}

// Run blocks, sweeping on each tick, until ctx is cancelled. Start it in its own
// goroutine. It performs one immediate sweep on start so a restart promptly clears
// subscriptions that lapsed while the process was down.
func (s *SubscriptionSweeper) Run(ctx context.Context) {
	log.Printf("subscription sweeper started: interval=%s", s.interval)
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	s.sweepOnce(ctx)
	for {
		select {
		case <-ctx.Done():
			log.Println("subscription sweeper stopped")
			return
		case <-ticker.C:
			s.sweepOnce(ctx)
		}
	}
}

// sweepOnce expires one batch of lapsed subscriptions and notifies each. An error
// is logged and the next tick retries; the guarded transition makes repetition
// safe. It logs only counts — never student identifiers (Non-Negotiable #8).
func (s *SubscriptionSweeper) sweepOnce(ctx context.Context) {
	expired, err := s.expirer.SweepExpired(ctx, s.limit)
	if err != nil {
		log.Printf("subscription sweeper: sweep error: %v", err)
	}
	for _, sub := range expired {
		s.notify.SubscriptionExpired(ctx, sub)
	}
	if len(expired) > 0 {
		log.Printf("subscription sweeper: expired %d subscription(s)", len(expired))
	}
}

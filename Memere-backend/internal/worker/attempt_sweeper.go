// Package worker holds in-process background workers for the Phase 1/2 monolith.
// The expiry sweeper guarantees the spec §9.2 "server timer fires" rule for
// attempts a client never returns to submit. It is deliberately kept behind the
// ExpirySweeper interface so a later phase can swap this ticker for an SQS/queue
// consumer without touching the engines.
package worker

import (
	"context"
	"log"
	"time"
)

// ExpirySweeper finalizes abandoned (in-progress, past-deadline) attempts for one
// engine. Both the quiz and exam engines satisfy it via their SweepExpired
// method. now is passed in so the worker (not each engine) owns the clock.
type ExpirySweeper interface {
	SweepExpired(ctx context.Context, now time.Time, limit int) (int, error)
}

// defaultBatchLimit caps how many attempts a single engine sweep grades per tick,
// bounding the work (and DB load) of any one pass; a backlog drains over ticks.
const defaultBatchLimit = 100

// AttemptSweeper runs the registered engine sweepers on a ticker until its
// context is cancelled (wired to the app's shutdown signal in main.go).
type AttemptSweeper struct {
	sweepers []ExpirySweeper
	interval time.Duration
	limit    int
	now      func() time.Time
}

// NewAttemptSweeper builds a sweeper that ticks every interval over the given
// engine sweepers (e.g. the quiz and exam services).
func NewAttemptSweeper(interval time.Duration, sweepers ...ExpirySweeper) *AttemptSweeper {
	return &AttemptSweeper{
		sweepers: sweepers,
		interval: interval,
		limit:    defaultBatchLimit,
		now:      time.Now,
	}
}

// Run blocks, sweeping on each tick, until ctx is cancelled. Start it in its own
// goroutine. It performs one immediate sweep on start so a restart promptly
// clears any attempts that expired while the process was down.
func (s *AttemptSweeper) Run(ctx context.Context) {
	log.Printf("attempt sweeper started: interval=%s", s.interval)
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	s.sweepOnce(ctx)
	for {
		select {
		case <-ctx.Done():
			log.Println("attempt sweeper stopped")
			return
		case <-ticker.C:
			s.sweepOnce(ctx)
		}
	}
}

// sweepOnce runs every engine sweeper once. A failing engine is logged and the
// others still run; the next tick retries. It logs only counts — never answer
// contents (Non-Negotiable #8).
func (s *AttemptSweeper) sweepOnce(ctx context.Context) {
	now := s.now()
	total := 0
	for _, sw := range s.sweepers {
		n, err := sw.SweepExpired(ctx, now, s.limit)
		total += n
		if err != nil {
			log.Printf("attempt sweeper: engine sweep error: %v", err)
		}
	}
	if total > 0 {
		log.Printf("attempt sweeper: graded %d expired attempt(s)", total)
	}
}

// RunOnce performs exactly one sweep and returns. It lets an external scheduler
// (Cloud Scheduler → /internal/jobs/sweep-attempts) drive the sweep when the
// in-process ticker is disabled on a scale-to-zero deployment. Safe to call
// concurrently with Run; the underlying engine sweeps are idempotent.
func (s *AttemptSweeper) RunOnce(ctx context.Context) { s.sweepOnce(ctx) }

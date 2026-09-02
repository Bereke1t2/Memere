package worker

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
)

// StreakNotifier fires the streak-warning push notification.
type StreakNotifier interface {
	StreakWarning(ctx context.Context, studentID uuid.UUID)
}

// noopStreakNotifier is the default when no notifier is wired (dev/test).
type noopStreakNotifier struct{}

func (noopStreakNotifier) StreakWarning(context.Context, uuid.UUID) {}

// EngagementSweeper warns students whose study streak is at risk of breaking.
// It runs daily (configurable via interval) and is idempotent: a student is
// warned at most once per calendar day because ListStreakAtRisk filters on
// last_warned_date and SetLastWarnedDate is called after each notification.
type EngagementSweeper struct {
	progress repository.ProgressRepository
	notify   StreakNotifier
	clock    func() time.Time
	interval time.Duration
	pageSize int
}

// NewEngagementSweeper builds a sweeper. notify may be nil (defaults to noop).
// clock may be nil (defaults to time.Now). interval defaults to 24h when zero.
func NewEngagementSweeper(
	progress repository.ProgressRepository,
	notify StreakNotifier,
	clock func() time.Time,
	interval time.Duration,
) *EngagementSweeper {
	if notify == nil {
		notify = noopStreakNotifier{}
	}
	if clock == nil {
		clock = time.Now
	}
	if interval <= 0 {
		interval = 24 * time.Hour
	}
	return &EngagementSweeper{
		progress: progress,
		notify:   notify,
		clock:    clock,
		interval: interval,
		pageSize: 100,
	}
}

// Run blocks until ctx is cancelled, sweeping once at start and then on each
// tick. Start it in its own goroutine alongside the subscription sweeper.
func (w *EngagementSweeper) Run(ctx context.Context) {
	log.Printf("engagement sweeper started: interval=%s", w.interval)
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	w.tick(ctx)
	for {
		select {
		case <-ctx.Done():
			log.Println("engagement sweeper stopped")
			return
		case <-ticker.C:
			w.tick(ctx)
		}
	}
}

// tick performs one sweep: page through at-risk students and send warnings.
func (w *EngagementSweeper) tick(ctx context.Context) {
	now := w.clock()
	cutoff := now.AddDate(0, 0, -2)
	today := now.Truncate(24 * time.Hour)

	atRisk, err := w.progress.ListStreakAtRisk(ctx, cutoff, today, w.pageSize)
	if err != nil {
		log.Printf("engagement sweeper: list at-risk: %v", err)
		return
	}
	warned := 0
	for _, st := range atRisk {
		w.notify.StreakWarning(ctx, st.StudentID)
		if setErr := w.progress.SetLastWarnedDate(ctx, st.StudentID, today); setErr != nil {
			log.Printf("engagement sweeper: set last warned: %v", setErr)
		}
		warned++
	}
	if warned > 0 {
		log.Printf("engagement sweeper: sent streak warnings to %d student(s)", warned)
	}
}

// RunOnce performs exactly one sweep and returns. It lets an external scheduler
// (Cloud Scheduler → /internal/jobs/sweep-engagement) drive the sweep when the
// in-process ticker is disabled on a scale-to-zero deployment. Idempotent per
// calendar day (last_warned_date), so safe to call repeatedly.
func (w *EngagementSweeper) RunOnce(ctx context.Context) { w.tick(ctx) }

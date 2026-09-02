package http

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SweepRunner performs one pass of a periodic background job. The concrete
// sweepers (attempt, subscription, engagement) satisfy it via their RunOnce
// method. Kept as a small interface so the handler has no dependency on the
// worker package's concrete types.
type SweepRunner interface {
	RunOnce(ctx context.Context)
}

// JobsHandler exposes the periodic sweeps as authenticated HTTP endpoints so an
// external scheduler (Cloud Scheduler) can drive them on a scale-to-zero
// deployment where the in-process tickers are disabled. Each field may be nil;
// a nil runner's route responds 200 with a "disabled" status rather than
// pretending work happened.
//
// The sweeps run synchronously within the request so Cloud Run keeps the CPU
// allocated until they finish (CPU is throttled once the response is sent). They
// are bounded (batch-limited) and idempotent, so a retried or overlapping call
// is safe.
type JobsHandler struct {
	Attempts      SweepRunner
	Subscriptions SweepRunner
	Engagement    SweepRunner
}

// NewJobsHandler builds the handler from the three sweep runners (any may be nil).
func NewJobsHandler(attempts, subscriptions, engagement SweepRunner) *JobsHandler {
	return &JobsHandler{
		Attempts:      attempts,
		Subscriptions: subscriptions,
		Engagement:    engagement,
	}
}

// SweepAttempts finalizes abandoned in-progress attempts past their deadline.
func (h *JobsHandler) SweepAttempts(c *gin.Context) { runSweep(c, h.Attempts) }

// SweepSubscriptions expires lapsed subscriptions and fires their notifications.
func (h *JobsHandler) SweepSubscriptions(c *gin.Context) { runSweep(c, h.Subscriptions) }

// SweepEngagement warns students whose study streak is about to break.
func (h *JobsHandler) SweepEngagement(c *gin.Context) { runSweep(c, h.Engagement) }

// runSweep executes one pass of the runner (using the request context so a
// client/Cloud Run timeout cancels it) and reports the outcome. Per-item errors
// are logged inside the sweep and never surface details to the caller.
func runSweep(c *gin.Context, runner SweepRunner) {
	if runner == nil {
		c.JSON(http.StatusOK, gin.H{"status": "disabled"})
		return
	}
	runner.RunOnce(c.Request.Context())
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

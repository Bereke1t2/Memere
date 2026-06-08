package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// AttemptStateStore holds the live, in-progress state of a quiz/exam attempt
// separately from its durable Postgres record (spec §9.1: randomization snapshot
// in Redis, partial answers saved every 30s). The Postgres row remains the
// source of truth for grading; this store is a fast, expiring scratchpad that is
// cleared once the attempt is submitted/graded.
//
// A miss (no snapshot / no answers) is a normal outcome, not an error: callers
// fall back to the durable question_order copy in Postgres when Redis is lost.
type AttemptStateStore interface {
	// SetSnapshot persists the randomized question/answer order for an attempt.
	SetSnapshot(ctx context.Context, attemptID uuid.UUID, snapshot map[string]any, ttl time.Duration) error
	// GetSnapshot returns the stored order, or (nil, nil) on a miss.
	GetSnapshot(ctx context.Context, attemptID uuid.UUID) (map[string]any, error)
	// SaveAnswers writes the current in-progress answers (auto-save).
	SaveAnswers(ctx context.Context, attemptID uuid.UUID, answers map[string]any, ttl time.Duration) error
	// GetAnswers returns the latest saved answers, or (nil, nil) on a miss.
	GetAnswers(ctx context.Context, attemptID uuid.UUID) (map[string]any, error)
	// DeleteAttemptState clears both the snapshot and answers for an attempt.
	DeleteAttemptState(ctx context.Context, attemptID uuid.UUID) error
}

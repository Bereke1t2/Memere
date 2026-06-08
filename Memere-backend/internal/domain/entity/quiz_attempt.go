package entity

import (
	"time"

	"github.com/google/uuid"
)

// AttemptStatus is the lifecycle state of a quiz or exam attempt (spec §9.2).
type AttemptStatus string

const (
	AttemptInProgress AttemptStatus = "in_progress"
	AttemptSubmitted  AttemptStatus = "submitted"
	AttemptGraded     AttemptStatus = "graded"
	AttemptExpired    AttemptStatus = "expired"
)

// Valid reports whether s is a recognized attempt status.
func (s AttemptStatus) Valid() bool {
	switch s {
	case AttemptInProgress, AttemptSubmitted, AttemptGraded, AttemptExpired:
		return true
	default:
		return false
	}
}

// CanTransitionTo encodes the §9.2 attempt state machine:
//
//	NOT_STARTED → IN_PROGRESS → {SUBMITTED | EXPIRED} → GRADED
//
// NOT_STARTED is implicit (no row yet); the first persisted state is
// in_progress. SUBMITTED is a manual submit, EXPIRED is the server timer firing;
// both auto-grade into GRADED, which is terminal. Any other move is illegal.
func (s AttemptStatus) CanTransitionTo(next AttemptStatus) bool {
	switch s {
	case AttemptInProgress:
		return next == AttemptSubmitted || next == AttemptExpired || next == AttemptGraded
	case AttemptSubmitted, AttemptExpired:
		return next == AttemptGraded
	default:
		return false
	}
}

// QuizAttempt records one student run at a quiz (spec §9.1). Score/Percentage/
// Passed stay nil until grading completes. AnswersSnapshot holds the submitted
// answers as JSONB; the correct-answer key is never stored here. QuestionOrder
// is the durable copy of the per-attempt randomized order (mirrored from Redis).
// ExpiresAt is the server-side timer deadline for timed quizzes (nil = untimed);
// the sweeper auto-grades attempts past it.
type QuizAttempt struct {
	ID              uuid.UUID
	QuizID          uuid.UUID
	StudentID       uuid.UUID
	AttemptNumber   int
	StartedAt       time.Time
	SubmittedAt     *time.Time
	ExpiresAt       *time.Time
	Score           *float64
	Percentage      *float64
	Passed          *bool
	QuestionOrder   map[string]any
	AnswersSnapshot map[string]any
	Status          AttemptStatus
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

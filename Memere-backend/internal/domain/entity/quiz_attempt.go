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

// QuizAttempt records one student run at a quiz (spec §9.1). Score/Percentage
// stay nil until grading completes. AnswersSnapshot holds the submitted answers
// as JSONB; the correct-answer key is never stored here.
type QuizAttempt struct {
	ID              uuid.UUID
	QuizID          uuid.UUID
	StudentID       uuid.UUID
	AttemptNumber   int
	StartedAt       time.Time
	SubmittedAt     *time.Time
	Score           *float64
	Percentage      *float64
	AnswersSnapshot map[string]any
	Status          AttemptStatus
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

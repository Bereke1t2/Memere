package entity

import (
	"time"

	"github.com/google/uuid"
)

// ExamAttempt records one student sitting of an exam (spec §9.2). StartedAt is
// the server-side timer source of truth; the sweeper auto-submits at expiry.
// AnswersSnapshot is the student's answers at submission, never the answer key.
type ExamAttempt struct {
	ID              uuid.UUID
	ExamID          uuid.UUID
	StudentID       uuid.UUID
	StartedAt       time.Time
	SubmittedAt     *time.Time
	Score           *float64
	Percentage      *float64
	AnswersSnapshot map[string]any
	Status          AttemptStatus
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

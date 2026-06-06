package entity

import (
	"time"

	"github.com/google/uuid"
)

// Answer is an option for a Question. IsCorrect is server-only data and must
// never be serialized to a client (spec §9.1 / Non-Negotiable #1) — grading
// happens server-side and DTOs omit this field.
type Answer struct {
	ID         uuid.UUID
	QuestionID uuid.UUID
	Text       string
	IsCorrect  bool
	OrderIndex int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

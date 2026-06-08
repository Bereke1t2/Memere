package entity

import (
	"time"

	"github.com/google/uuid"
)

// Exam is a timed mock exam (spec §4.2.6 / §9.2). DurationMinutes drives the
// server-enforced timer; the client timer is display-only (Non-Negotiable #2).
type Exam struct {
	ID              uuid.UUID
	CourseID        *uuid.UUID
	Title           string
	Subject         string
	Grade           int
	DurationMinutes int
	TotalMarks      int
	PassMarks       int
	Instructions    *string
	IsPublished     bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time
}

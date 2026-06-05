package entity

import (
	"time"

	"github.com/google/uuid"
)

// CourseSection groups lessons within a course (spec §4.2.3).
type CourseSection struct {
	ID          uuid.UUID
	CourseID    uuid.UUID
	Title       string
	Description *string
	OrderIndex  int
	IsPublished bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

package entity

import (
	"time"

	"github.com/google/uuid"
)

// Level is the difficulty level of a course (spec §4.2.2).
type Level string

const (
	LevelBeginner     Level = "beginner"
	LevelIntermediate Level = "intermediate"
	LevelAdvanced     Level = "advanced"
)

// Valid reports whether l is a recognized level.
func (l Level) Valid() bool {
	switch l {
	case LevelBeginner, LevelIntermediate, LevelAdvanced:
		return true
	default:
		return false
	}
}

// Course is a sellable unit of content (spec §4.2.2). Price/RatingAvg map to
// DECIMAL columns; Phase 1 performs no money arithmetic on them, so float64 is
// sufficient for display and storage round-tripping. Metadata is JSONB.
type Course struct {
	ID                   uuid.UUID
	TeacherID            uuid.UUID
	Title                string
	Slug                 string
	Description          string
	ShortDescription     *string
	Subject              string
	Grade                int
	ThumbnailURL         *string
	Price                float64
	Currency             string
	IsFree               bool
	IsPublished          bool
	Language             string
	Level                Level
	TotalDurationSeconds int
	TotalLessons         int
	RatingAvg            float64
	EnrollmentCount      int
	Metadata             map[string]any
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            *time.Time
}

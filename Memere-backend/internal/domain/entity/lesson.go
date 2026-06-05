package entity

import (
	"time"

	"github.com/google/uuid"
)

// LessonType is the kind of content a lesson delivers (spec §4.2.3).
type LessonType string

const (
	LessonTypeVideo LessonType = "video"
	LessonTypeNote  LessonType = "note"
	LessonTypeQuiz  LessonType = "quiz"
	LessonTypeMixed LessonType = "mixed"
)

// Valid reports whether t is a recognized lesson type.
func (t LessonType) Valid() bool {
	switch t {
	case LessonTypeVideo, LessonTypeNote, LessonTypeQuiz, LessonTypeMixed:
		return true
	default:
		return false
	}
}

// Lesson is a single unit of study within a section (spec §4.2.3). CourseID is
// denormalized for fast course-level queries.
type Lesson struct {
	ID              uuid.UUID
	SectionID       uuid.UUID
	CourseID        uuid.UUID
	Title           string
	Type            LessonType
	OrderIndex      int
	IsFreePreview   bool
	DurationSeconds int
	IsPublished     bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time
}

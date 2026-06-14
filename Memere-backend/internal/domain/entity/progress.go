package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// LessonProgress tracks a student's engagement with a single lesson. No
// db/json tags (Clean Architecture: domain entities have zero external deps).
type LessonProgress struct {
	ID                   uuid.UUID
	StudentID            uuid.UUID
	LessonID             uuid.UUID
	CourseID             uuid.UUID
	IsCompleted          bool
	CompletedAt          *time.Time
	VideoProgressSeconds int
	LastAccessedAt       *time.Time
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            *time.Time
}

// CourseProgress is the denormalized per-student completion rollup for a course.
// Rebuilt by RecomputeCourseProgress after each lesson completion.
type CourseProgress struct {
	StudentID        uuid.UUID
	CourseID         uuid.UUID
	CompletedLessons int
	TotalLessons     int
	PercentComplete  decimal.Decimal
	CompletedAt      *time.Time
	UpdatedAt        time.Time
}

// StudyStreak holds the current and longest study streaks for a student. Streaks
// are day-scoped (timezone-configurable, default Africa/Addis_Ababa per spec).
type StudyStreak struct {
	StudentID      uuid.UUID
	CurrentStreak  int
	LongestStreak  int
	LastStudyDate  *time.Time
	UpdatedAt      time.Time
}

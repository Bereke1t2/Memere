package entity

import (
	"time"

	"github.com/google/uuid"
)

// QuestionType enumerates the kinds of quiz/exam questions (spec §4.2.5).
type QuestionType string

const (
	QuestionMultipleChoice QuestionType = "multiple_choice"
	QuestionTrueFalse      QuestionType = "true_false"
	QuestionShortAnswer    QuestionType = "short_answer"
)

// Valid reports whether t is a recognized question type.
func (t QuestionType) Valid() bool {
	switch t {
	case QuestionMultipleChoice, QuestionTrueFalse, QuestionShortAnswer:
		return true
	default:
		return false
	}
}

// Quiz is a per-lesson assessment (spec §4.2.5). TimeLimitSeconds and MaxAttempts
// are nil when unlimited. PassPercentage maps to a DECIMAL column.
type Quiz struct {
	ID                 uuid.UUID
	LessonID           *uuid.UUID
	CourseID           uuid.UUID
	Title              string
	TimeLimitSeconds   *int
	PassPercentage     float64
	RandomizeQuestions bool
	MaxAttempts        *int
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          *time.Time
}

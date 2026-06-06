package entity

import (
	"time"

	"github.com/google/uuid"
)

// Question belongs to a Quiz (spec §4.2.5). Explanation is revealed only after
// grading. OrderIndex fixes the canonical (un-randomized) display order.
type Question struct {
	ID          uuid.UUID
	QuizID      uuid.UUID
	Text        string
	Type        QuestionType
	Points      int
	Explanation *string
	OrderIndex  int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

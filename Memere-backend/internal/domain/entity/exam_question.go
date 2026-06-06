package entity

import "github.com/google/uuid"

// ExamQuestion links an Exam to a Question from the shared question bank, with a
// per-exam order and mark value (spec §9.2). The join lets one question be
// reused across exams without duplicating it.
type ExamQuestion struct {
	ID         uuid.UUID
	ExamID     uuid.UUID
	QuestionID uuid.UUID
	OrderIndex int
	Marks      int
}

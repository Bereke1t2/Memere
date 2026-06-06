package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/pagination"
)

// ExamFilter narrows an exam listing. Nil fields are ignored.
type ExamFilter struct {
	Subject     *string
	Grade       *int
	IsPublished *bool
}

// ExamRepository persists exams and their assembled question set (spec §9.2).
type ExamRepository interface {
	Create(ctx context.Context, e *entity.Exam) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Exam, error)
	// List returns up to limit exams matching filter, ordered by (created_at, id)
	// DESC starting after cursor (nil = first page). The returned cursor is
	// non-nil only when more rows may exist.
	List(ctx context.Context, filter ExamFilter, cursor *pagination.Cursor, limit int) ([]*entity.Exam, *pagination.Cursor, error)
	Update(ctx context.Context, e *entity.Exam) error
	// AddQuestion links a bank question to an exam with order and marks.
	AddQuestion(ctx context.Context, eq *entity.ExamQuestion) error
	// ListQuestions returns the exam's questions in order, for the grading core.
	// Strip answer keys before any client response.
	ListQuestions(ctx context.Context, examID uuid.UUID) ([]*entity.ExamQuestion, error)
}

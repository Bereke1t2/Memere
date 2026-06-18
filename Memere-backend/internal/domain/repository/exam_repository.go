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

// ExamWithQuestions is an exam and its assembled question set, answer keys
// included, for the grading core. Never serialize directly to a client.
type ExamWithQuestions struct {
	Exam      *entity.Exam
	Questions []QuestionWithAnswers
}

// ClientExamQuestion is one exam question rendered to a student: the per-exam
// order and marks plus the question body and options, with no answer key.
type ClientExamQuestion struct {
	QuestionID uuid.UUID
	OrderIndex int
	Marks      int
	Text       string
	Type       entity.QuestionType
	Subject    *string
	Topic      *string
	Answers    []ClientAnswer
}

// ExamRepository persists exams and their assembled question set (spec §9.2).
type ExamRepository interface {
	Create(ctx context.Context, e *entity.Exam) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Exam, error)
	// ListByCourse returns all non-deleted exams for a given course, ordered by
	// created_at DESC. Used by the teacher authoring UI.
	ListByCourse(ctx context.Context, courseID uuid.UUID) ([]*entity.Exam, error)
	// List returns up to limit exams matching filter, ordered by (created_at, id)
	// DESC starting after cursor (nil = first page). The returned cursor is
	// non-nil only when more rows may exist.
	List(ctx context.Context, filter ExamFilter, cursor *pagination.Cursor, limit int) ([]*entity.Exam, *pagination.Cursor, error)
	Update(ctx context.Context, e *entity.Exam) error
	// Delete soft-deletes the exam (Non-Negotiable #5).
	Delete(ctx context.Context, id uuid.UUID) error
	// AddQuestion links a bank question to an exam with order and marks.
	AddQuestion(ctx context.Context, eq *entity.ExamQuestion) error
	// ListQuestions returns the exam's question links in order.
	ListQuestions(ctx context.Context, examID uuid.UUID) ([]*entity.ExamQuestion, error)
	// GetExamWithQuestions loads the exam plus its question set with answer keys,
	// for the grading core. Never serialize directly to a client.
	GetExamWithQuestions(ctx context.Context, examID uuid.UUID) (*ExamWithQuestions, error)
	// GetQuestionsForClient loads the exam's questions and options with the
	// answer key structurally excluded — safe to render to a student.
	GetQuestionsForClient(ctx context.Context, examID uuid.UUID) ([]ClientExamQuestion, error)
}

package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
)

// QuestionRepository persists questions and their answer options (spec §4.2.5).
// Create writes a question and its answers together so the answer key is always
// consistent.
type QuestionRepository interface {
	Create(ctx context.Context, q *entity.Question, answers []*entity.Answer) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Question, error)
	ListByQuiz(ctx context.Context, quizID uuid.UUID) ([]*entity.Question, error)
	ListAnswers(ctx context.Context, questionID uuid.UUID) ([]*entity.Answer, error)
	Update(ctx context.Context, q *entity.Question) error
	Delete(ctx context.Context, id uuid.UUID) error
}

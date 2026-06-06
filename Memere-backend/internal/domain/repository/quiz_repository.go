package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
)

// QuestionWithAnswers is a question together with its answer options.
type QuestionWithAnswers struct {
	Question *entity.Question
	Answers  []*entity.Answer
}

// QuizWithQuestions is a quiz and its ordered question tree. Repositories load
// it with answer keys for server-side grading; callers MUST strip
// Answer.IsCorrect before the data leaves the server (Non-Negotiable #1).
type QuizWithQuestions struct {
	Quiz      *entity.Quiz
	Questions []QuestionWithAnswers
}

// QuizRepository persists quizzes and their question tree (spec §4.2.5).
type QuizRepository interface {
	Create(ctx context.Context, q *entity.Quiz) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Quiz, error)
	ListByCourse(ctx context.Context, courseID uuid.UUID) ([]*entity.Quiz, error)
	Update(ctx context.Context, q *entity.Quiz) error
	Delete(ctx context.Context, id uuid.UUID) error
	// GetQuizWithQuestions loads the full quiz tree including answer keys for
	// grading. Never serialize the result directly to a client.
	GetQuizWithQuestions(ctx context.Context, quizID uuid.UUID) (*QuizWithQuestions, error)
}

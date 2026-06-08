package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
)

// QuestionWithAnswers is a question together with its answer options (answer key
// included). Server-internal only.
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

// ClientAnswer is an answer option safe to render to a student: it carries no
// is_correct field by construction (Non-Negotiable #1).
type ClientAnswer struct {
	ID         uuid.UUID
	QuestionID uuid.UUID
	Text       string
	OrderIndex int
}

// ClientQuestion is a question rendered to a student — its options without any
// answer key. The explanation is also omitted (revealed only after grading).
type ClientQuestion struct {
	ID         uuid.UUID
	QuizID     uuid.UUID
	Text       string
	Type       entity.QuestionType
	Points     int
	OrderIndex int
	Subject    *string
	Topic      *string
	Answers    []ClientAnswer
}

// QuizRepository persists quizzes and their question tree (spec §4.2.5).
type QuizRepository interface {
	Create(ctx context.Context, q *entity.Quiz) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Quiz, error)
	ListByCourse(ctx context.Context, courseID uuid.UUID) ([]*entity.Quiz, error)
	Update(ctx context.Context, q *entity.Quiz) error
	// Delete soft-deletes the quiz (Non-Negotiable #5); attempts are preserved.
	Delete(ctx context.Context, id uuid.UUID) error
	// GetQuizWithQuestions loads the full quiz tree including answer keys for
	// grading. Never serialize the result directly to a client.
	GetQuizWithQuestions(ctx context.Context, quizID uuid.UUID) (*QuizWithQuestions, error)
	// GetQuestionsForClient loads the quiz's questions and options with the
	// answer key structurally excluded — safe to render to a student.
	GetQuestionsForClient(ctx context.Context, quizID uuid.UUID) ([]ClientQuestion, error)
}

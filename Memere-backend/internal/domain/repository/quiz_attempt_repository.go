package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
)

// QuizAttemptRepository persists quiz attempts and their grading results
// (spec §9.1). Every read filters by the authenticated student_id to prevent
// IDOR (Non-Negotiable #7).
type QuizAttemptRepository interface {
	Create(ctx context.Context, a *entity.QuizAttempt) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.QuizAttempt, error)
	ListByStudentAndQuiz(ctx context.Context, studentID, quizID uuid.UUID) ([]*entity.QuizAttempt, error)
	CountByStudentAndQuiz(ctx context.Context, studentID, quizID uuid.UUID) (int, error)
	// Grade persists the final score/percentage and flips status to graded.
	Grade(ctx context.Context, a *entity.QuizAttempt) error
}

package repository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
)

// QuizAttemptRepository persists quiz attempts and their grading results
// (spec §9.1). Every read filters by the authenticated student_id to prevent
// IDOR (Non-Negotiable #7).
type QuizAttemptRepository interface {
	Create(ctx context.Context, a *entity.QuizAttempt) error
	// FindByID scopes the lookup to studentID so a student can never read another
	// student's attempt (IDOR, Non-Negotiable #7).
	FindByID(ctx context.Context, id, studentID uuid.UUID) (*entity.QuizAttempt, error)
	// GetActive returns the student's current in-progress attempt at a quiz, or
	// apperror.NotFound if none — used to resume and to block concurrent starts.
	GetActive(ctx context.Context, studentID, quizID uuid.UUID) (*entity.QuizAttempt, error)
	ListByStudentAndQuiz(ctx context.Context, studentID, quizID uuid.UUID) ([]*entity.QuizAttempt, error)
	CountByStudentAndQuiz(ctx context.Context, studentID, quizID uuid.UUID) (int, error)
	// ListExpired returns in-progress, timed attempts past their deadline, so the
	// sweeper can auto-grade them (spec §9.1 time enforcement, Non-Negotiable #2).
	ListExpired(ctx context.Context, now time.Time, limit int) ([]*entity.QuizAttempt, error)
	// Update persists auto-save state and status transitions (answers snapshot,
	// status, submitted_at) without touching the immutable timer columns.
	Update(ctx context.Context, a *entity.QuizAttempt) error
	// ClaimForGrading atomically transitions the attempt out of in_progress to
	// a.Status (submitted/expired), persisting its answers snapshot, ONLY if it is
	// still in_progress. It reports claimed=false (without error) when another path
	// — a racing client submit or the sweeper — already moved it, so callers can
	// no-op rather than double-grade (race guard, spec §9.2).
	ClaimForGrading(ctx context.Context, a *entity.QuizAttempt) (claimed bool, err error)
	// Grade persists the final score/percentage/passed and flips status to graded.
	Grade(ctx context.Context, a *entity.QuizAttempt) error
}

package repository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
)

// ExamAttemptRepository persists exam sittings (spec §9.2). started_at is the
// server-side timer source of truth; ListExpired backs the auto-submit sweeper.
// Every read filters by the authenticated student_id to prevent IDOR
// (Non-Negotiable #7).
type ExamAttemptRepository interface {
	Create(ctx context.Context, a *entity.ExamAttempt) error
	// FindByID scopes the lookup to studentID to prevent IDOR (Non-Negotiable #7).
	FindByID(ctx context.Context, id, studentID uuid.UUID) (*entity.ExamAttempt, error)
	// GetActive returns the student's current in-progress sitting of an exam, or
	// apperror.NotFound if none.
	GetActive(ctx context.Context, studentID, examID uuid.UUID) (*entity.ExamAttempt, error)
	ListByStudent(ctx context.Context, studentID uuid.UUID) ([]*entity.ExamAttempt, error)
	// ListExpired returns in-progress attempts whose timer (started_at +
	// exams.duration_minutes) has elapsed as of now, for the sweeper (spec §9.2).
	ListExpired(ctx context.Context, now time.Time, limit int) ([]*entity.ExamAttempt, error)
	// Update persists auto-save state and status transitions without touching the
	// immutable started_at timer column.
	Update(ctx context.Context, a *entity.ExamAttempt) error
	Grade(ctx context.Context, a *entity.ExamAttempt) error
}

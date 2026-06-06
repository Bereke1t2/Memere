package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
)

// ExamAttemptRepository persists exam sittings (spec §9.2). started_at is the
// server-side timer source of truth; FindExpired backs the auto-submit sweeper.
type ExamAttemptRepository interface {
	Create(ctx context.Context, a *entity.ExamAttempt) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.ExamAttempt, error)
	ListByStudent(ctx context.Context, studentID uuid.UUID) ([]*entity.ExamAttempt, error)
	// FindExpired returns in-progress attempts whose timer has elapsed, so the
	// sweeper can auto-grade them (spec §9.2).
	FindExpired(ctx context.Context, limit int) ([]*entity.ExamAttempt, error)
	Grade(ctx context.Context, a *entity.ExamAttempt) error
}

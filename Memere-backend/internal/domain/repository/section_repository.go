package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
)

// SectionRepository persists course sections (spec §4.2.3). Every read filters
// deleted_at IS NULL.
type SectionRepository interface {
	Create(ctx context.Context, s *entity.CourseSection) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.CourseSection, error)
	ListByCourse(ctx context.Context, courseID uuid.UUID) ([]*entity.CourseSection, error)
	Update(ctx context.Context, s *entity.CourseSection) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	// MaxOrderIndex returns the highest order_index among a course's live
	// sections, or -1 when the course has none (so callers append at max+1).
	MaxOrderIndex(ctx context.Context, courseID uuid.UUID) (int, error)
}

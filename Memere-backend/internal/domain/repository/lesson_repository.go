package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
)

// LessonRepository persists lessons (spec §4.2.3). Every read filters
// deleted_at IS NULL.
type LessonRepository interface {
	Create(ctx context.Context, l *entity.Lesson) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Lesson, error)
	ListBySection(ctx context.Context, sectionID uuid.UUID) ([]*entity.Lesson, error)
	ListByCourse(ctx context.Context, courseID uuid.UUID) ([]*entity.Lesson, error)
	Update(ctx context.Context, l *entity.Lesson) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	// SoftDeleteBySection tombstones every live lesson in a section, used to
	// cascade a section soft-delete.
	SoftDeleteBySection(ctx context.Context, sectionID uuid.UUID) error
	// MaxOrderIndex returns the highest order_index among a section's live
	// lessons, or -1 when the section has none.
	MaxOrderIndex(ctx context.Context, sectionID uuid.UUID) (int, error)
}

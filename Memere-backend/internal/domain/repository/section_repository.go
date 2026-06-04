package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
)

// SectionRepository persists course sections (spec §4.2.3).
type SectionRepository interface {
	Create(ctx context.Context, s *entity.CourseSection) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.CourseSection, error)
	ListByCourse(ctx context.Context, courseID uuid.UUID) ([]*entity.CourseSection, error)
}

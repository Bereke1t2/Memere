package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
)

// LessonRepository persists lessons (spec §4.2.3).
type LessonRepository interface {
	Create(ctx context.Context, l *entity.Lesson) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Lesson, error)
	ListBySection(ctx context.Context, sectionID uuid.UUID) ([]*entity.Lesson, error)
	ListByCourse(ctx context.Context, courseID uuid.UUID) ([]*entity.Lesson, error)
}

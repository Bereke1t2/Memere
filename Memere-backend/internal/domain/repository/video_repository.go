package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
)

// VideoRepository is the persistence port for lesson videos. Implementations
// filter soft-deleted rows and translate driver errors to apperror.
type VideoRepository interface {
	Create(ctx context.Context, v *entity.Video) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Video, error)
	GetByLessonID(ctx context.Context, lessonID uuid.UUID) (*entity.Video, error)
	// SetReady writes the transcode outputs and flips status to ready.
	SetReady(ctx context.Context, v *entity.Video) error
	// SetError records a (truncated, secret-free) processing error message for a
	// video whose transcode failed. It does not change status; the caller flips
	// to failed via the guarded UpdateStatus.
	SetError(ctx context.Context, id uuid.UUID, msg string) error
	// UpdateStatus is a guarded transition: it only flips when the current
	// status still equals from, returning false (not an error) otherwise so two
	// workers can race for the same row safely.
	UpdateStatus(ctx context.Context, id uuid.UUID, from, to entity.VideoStatus) (bool, error)
	ListByStatus(ctx context.Context, status entity.VideoStatus, limit int) ([]*entity.Video, error)
	SoftDelete(ctx context.Context, id uuid.UUID) error
}

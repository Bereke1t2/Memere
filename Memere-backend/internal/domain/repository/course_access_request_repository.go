package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/pagination"
)

// CourseAccessRequestRepository persists and queries course access requests.
type CourseAccessRequestRepository interface {
	Create(ctx context.Context, req *entity.CourseAccessRequest) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.CourseAccessRequest, error)
	GetByStudentAndCourse(ctx context.Context, studentID, courseID uuid.UUID) (*entity.CourseAccessRequest, error)
	Update(ctx context.Context, req *entity.CourseAccessRequest) error
	ListByTeacher(ctx context.Context, teacherID uuid.UUID, status *string, cursor *pagination.Cursor, limit int) ([]*entity.CourseAccessRequest, *pagination.Cursor, error)
	ListAll(ctx context.Context, status *string, cursor *pagination.Cursor, limit int) ([]*entity.CourseAccessRequest, *pagination.Cursor, error)
	ListByStudent(ctx context.Context, studentID uuid.UUID) ([]*entity.CourseAccessRequest, error)
}

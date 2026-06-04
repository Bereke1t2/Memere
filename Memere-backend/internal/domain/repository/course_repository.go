package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/pagination"
)

// CourseFilter narrows a course listing. Nil fields are ignored.
type CourseFilter struct {
	Subject     *string
	Grade       *int
	IsPublished *bool
	TeacherID   *uuid.UUID
}

// SectionWithLessons is a section together with its ordered lessons.
type SectionWithLessons struct {
	Section *entity.CourseSection
	Lessons []*entity.Lesson
}

// CourseWithContent is a course with its full section/lesson tree, used to
// render a course detail page in one round trip.
type CourseWithContent struct {
	Course   *entity.Course
	Sections []SectionWithLessons
}

// CourseRepository persists courses and their content tree (spec §4.2.2–4.2.3).
// Every read filters deleted_at IS NULL.
type CourseRepository interface {
	Create(ctx context.Context, c *entity.Course) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Course, error)
	FindBySlug(ctx context.Context, slug string) (*entity.Course, error)
	// List returns up to limit courses matching filter, ordered by
	// (created_at, id) DESC starting after cursor (nil = first page). The
	// returned cursor is non-nil only when more rows may exist.
	List(ctx context.Context, filter CourseFilter, cursor *pagination.Cursor, limit int) ([]*entity.Course, *pagination.Cursor, error)
	Update(ctx context.Context, c *entity.Course) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	GetCourseWithSectionsAndLessons(ctx context.Context, courseID uuid.UUID) (*CourseWithContent, error)
}

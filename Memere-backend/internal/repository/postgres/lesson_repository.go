package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/internal/repository/postgres/sqlcgen"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// LessonRepo is the sqlc-backed implementation of repository.LessonRepository.
type LessonRepo struct {
	q *sqlcgen.Queries
}

var _ repository.LessonRepository = (*LessonRepo)(nil)

// NewLessonRepo builds a LessonRepo over a pgx pool.
func NewLessonRepo(pool *pgxpool.Pool) *LessonRepo {
	return &LessonRepo{q: sqlcgen.New(pool)}
}

// Create inserts a lesson and reflects DB-assigned fields back onto the entity.
func (r *LessonRepo) Create(ctx context.Context, l *entity.Lesson) error {
	row, err := queriesFor(ctx, r.q).CreateLesson(ctx, sqlcgen.CreateLessonParams{
		SectionID:       toPgUUID(l.SectionID),
		CourseID:        toPgUUID(l.CourseID),
		Title:           l.Title,
		Type:            string(l.Type),
		OrderIndex:      int32(l.OrderIndex),
		IsFreePreview:   l.IsFreePreview,
		DurationSeconds: int32(l.DurationSeconds),
		IsPublished:     l.IsPublished,
		Content:         l.Content,
		PdfUrl:          l.PdfURL,
	})
	if err != nil {
		return apperror.Internal(err)
	}
	*l = *lessonFromCoursesLesson(row)
	return nil
}

// FindByID returns the lesson or apperror.NotFound (soft-deleted excluded).
func (r *LessonRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.Lesson, error) {
	row, err := queriesFor(ctx, r.q).GetLessonByID(ctx, toPgUUID(id))
	if err != nil {
		return nil, mapLessonErr(err)
	}
	return lessonFromGetRow(row), nil
}

// ListBySection returns a section's live lessons ordered by order_index.
func (r *LessonRepo) ListBySection(ctx context.Context, sectionID uuid.UUID) ([]*entity.Lesson, error) {
	rows, err := queriesFor(ctx, r.q).ListLessonsBySection(ctx, toPgUUID(sectionID))
	if err != nil {
		return nil, apperror.Internal(err)
	}
	lessons := make([]*entity.Lesson, len(rows))
	for i, row := range rows {
		lessons[i] = lessonFromSectionRow(row)
	}
	return lessons, nil
}

// ListByCourse returns all of a course's live lessons ordered by order_index.
func (r *LessonRepo) ListByCourse(ctx context.Context, courseID uuid.UUID) ([]*entity.Lesson, error) {
	rows, err := queriesFor(ctx, r.q).ListLessonsByCourse(ctx, toPgUUID(courseID))
	if err != nil {
		return nil, apperror.Internal(err)
	}
	lessons := make([]*entity.Lesson, len(rows))
	for i, row := range rows {
		lessons[i] = lessonFromCourseRow(row)
	}
	return lessons, nil
}

// Update persists the mutable lesson fields (the query filters deleted_at IS NULL).
func (r *LessonRepo) Update(ctx context.Context, l *entity.Lesson) error {
	row, err := queriesFor(ctx, r.q).UpdateLesson(ctx, sqlcgen.UpdateLessonParams{
		ID:              toPgUUID(l.ID),
		Title:           l.Title,
		Type:            string(l.Type),
		OrderIndex:      int32(l.OrderIndex),
		IsFreePreview:   l.IsFreePreview,
		DurationSeconds: int32(l.DurationSeconds),
		IsPublished:     l.IsPublished,
		Content:         l.Content,
		PdfUrl:          l.PdfURL,
	})
	if err != nil {
		return mapLessonErr(err)
	}
	updated := lessonFromCoursesLesson(row)
	updated.VideoID = l.VideoID
	*l = *updated
	return nil
}

// SoftDelete sets deleted_at (Non-Negotiable #5).
func (r *LessonRepo) SoftDelete(ctx context.Context, id uuid.UUID) error {
	if err := queriesFor(ctx, r.q).SoftDeleteLesson(ctx, toPgUUID(id)); err != nil {
		return apperror.Internal(err)
	}
	return nil
}

// SoftDeleteBySection tombstones every live lesson in a section.
func (r *LessonRepo) SoftDeleteBySection(ctx context.Context, sectionID uuid.UUID) error {
	if err := queriesFor(ctx, r.q).SoftDeleteLessonsBySection(ctx, toPgUUID(sectionID)); err != nil {
		return apperror.Internal(err)
	}
	return nil
}

// MaxOrderIndex returns the highest order_index among a section's live lessons, or -1.
func (r *LessonRepo) MaxOrderIndex(ctx context.Context, sectionID uuid.UUID) (int, error) {
	max, err := queriesFor(ctx, r.q).MaxLessonOrderIndex(ctx, toPgUUID(sectionID))
	if err != nil {
		return 0, apperror.Internal(err)
	}
	return int(max), nil
}

func mapLessonErr(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return apperror.NotFound("lesson not found", err)
	}
	return apperror.Internal(err)
}

func lessonFromCoursesLesson(row sqlcgen.CoursesLesson) *entity.Lesson {
	return &entity.Lesson{
		ID:              fromPgUUID(row.ID),
		SectionID:       fromPgUUID(row.SectionID),
		CourseID:        fromPgUUID(row.CourseID),
		Title:           row.Title,
		Type:            entity.LessonType(row.Type),
		OrderIndex:      int(row.OrderIndex),
		IsFreePreview:   row.IsFreePreview,
		DurationSeconds: int(row.DurationSeconds),
		IsPublished:     row.IsPublished,
		Content:         row.Content,
		PdfURL:          row.PdfUrl,
		CreatedAt:       fromPgTimestamptzValue(row.CreatedAt),
		UpdatedAt:       fromPgTimestamptzValue(row.UpdatedAt),
		DeletedAt:       fromPgTimestamptz(row.DeletedAt),
	}
}

func lessonFromGetRow(row sqlcgen.GetLessonByIDRow) *entity.Lesson {
	var videoID *uuid.UUID
	if row.VideoID.Valid {
		id := fromPgUUID(row.VideoID)
		videoID = &id
	}
	return &entity.Lesson{
		ID:              fromPgUUID(row.ID),
		SectionID:       fromPgUUID(row.SectionID),
		CourseID:        fromPgUUID(row.CourseID),
		Title:           row.Title,
		Type:            entity.LessonType(row.Type),
		OrderIndex:      int(row.OrderIndex),
		IsFreePreview:   row.IsFreePreview,
		DurationSeconds: int(row.DurationSeconds),
		IsPublished:     row.IsPublished,
		VideoID:         videoID,
		Content:         row.Content,
		PdfURL:          row.PdfUrl,
		CreatedAt:       fromPgTimestamptzValue(row.CreatedAt),
		UpdatedAt:       fromPgTimestamptzValue(row.UpdatedAt),
		DeletedAt:       fromPgTimestamptz(row.DeletedAt),
	}
}

func lessonFromSectionRow(row sqlcgen.ListLessonsBySectionRow) *entity.Lesson {
	var videoID *uuid.UUID
	if row.VideoID.Valid {
		id := fromPgUUID(row.VideoID)
		videoID = &id
	}
	return &entity.Lesson{
		ID:              fromPgUUID(row.ID),
		SectionID:       fromPgUUID(row.SectionID),
		CourseID:        fromPgUUID(row.CourseID),
		Title:           row.Title,
		Type:            entity.LessonType(row.Type),
		OrderIndex:      int(row.OrderIndex),
		IsFreePreview:   row.IsFreePreview,
		DurationSeconds: int(row.DurationSeconds),
		IsPublished:     row.IsPublished,
		VideoID:         videoID,
		Content:         row.Content,
		PdfURL:          row.PdfUrl,
		CreatedAt:       fromPgTimestamptzValue(row.CreatedAt),
		UpdatedAt:       fromPgTimestamptzValue(row.UpdatedAt),
		DeletedAt:       fromPgTimestamptz(row.DeletedAt),
	}
}

func lessonFromCourseRow(row sqlcgen.ListLessonsByCourseRow) *entity.Lesson {
	var videoID *uuid.UUID
	if row.VideoID.Valid {
		id := fromPgUUID(row.VideoID)
		videoID = &id
	}
	return &entity.Lesson{
		ID:              fromPgUUID(row.ID),
		SectionID:       fromPgUUID(row.SectionID),
		CourseID:        fromPgUUID(row.CourseID),
		Title:           row.Title,
		Type:            entity.LessonType(row.Type),
		OrderIndex:      int(row.OrderIndex),
		IsFreePreview:   row.IsFreePreview,
		DurationSeconds: int(row.DurationSeconds),
		IsPublished:     row.IsPublished,
		VideoID:         videoID,
		Content:         row.Content,
		PdfURL:          row.PdfUrl,
		CreatedAt:       fromPgTimestamptzValue(row.CreatedAt),
		UpdatedAt:       fromPgTimestamptzValue(row.UpdatedAt),
		DeletedAt:       fromPgTimestamptz(row.DeletedAt),
	}
}

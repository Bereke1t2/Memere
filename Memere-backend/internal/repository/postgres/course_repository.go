package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/internal/repository/postgres/sqlcgen"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/pagination"
)

// CourseRepo is the sqlc-backed implementation of repository.CourseRepository.
type CourseRepo struct {
	q *sqlcgen.Queries
}

var _ repository.CourseRepository = (*CourseRepo)(nil)

// NewCourseRepo builds a CourseRepo over a pgx pool.
func NewCourseRepo(pool *pgxpool.Pool) *CourseRepo {
	return &CourseRepo{q: sqlcgen.New(pool)}
}

// Create inserts a course. A duplicate slug collapses to apperror.Conflict so
// callers never inspect SQLSTATE.
func (r *CourseRepo) Create(ctx context.Context, c *entity.Course) error {
	row, err := queriesFor(ctx, r.q).CreateCourse(ctx, sqlcgen.CreateCourseParams{
		TeacherID:        toPgUUID(c.TeacherID),
		Title:            c.Title,
		Slug:             c.Slug,
		Description:      c.Description,
		ShortDescription: c.ShortDescription,
		Subject:          c.Subject,
		Grade:            int32(c.Grade),
		ThumbnailUrl:     c.ThumbnailURL,
		Price:            toPgNumeric(c.Price),
		Currency:         c.Currency,
		IsFree:           c.IsFree,
		IsPublished:      c.IsPublished,
		Language:         c.Language,
		Level:            string(c.Level),
		Metadata:         toJSONB(c.Metadata),
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolation {
			return apperror.Conflict("SLUG_TAKEN", err)
		}
		return apperror.Internal(err)
	}
	*c = *courseFromRow(row)
	return nil
}

// FindByID returns the course or apperror.NotFound. Soft-deleted rows are
// excluded by the query.
func (r *CourseRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.Course, error) {
	row, err := queriesFor(ctx, r.q).GetCourseByID(ctx, toPgUUID(id))
	if err != nil {
		return nil, mapCourseErr(err)
	}
	return courseFromRow(row), nil
}

// FindBySlug returns the course or apperror.NotFound.
func (r *CourseRepo) FindBySlug(ctx context.Context, slug string) (*entity.Course, error) {
	row, err := queriesFor(ctx, r.q).GetCourseBySlug(ctx, slug)
	if err != nil {
		return nil, mapCourseErr(err)
	}
	return courseFromRow(row), nil
}

// List returns up to limit courses matching filter, ordered by (created_at, id)
// DESC after cursor. It fetches limit+1 rows: the extra row, if present, signals
// a next page and yields the returned cursor.
func (r *CourseRepo) List(ctx context.Context, filter repository.CourseFilter, cursor *pagination.Cursor, limit int) ([]*entity.Course, *pagination.Cursor, error) {
	limit = pagination.NormalizeLimit(limit)

	params := sqlcgen.ListCoursesParams{
		Subject:     filter.Subject,
		IsPublished: filter.IsPublished,
		RowLimit:    int32(limit + 1),
	}
	if filter.Grade != nil {
		g := int32(*filter.Grade)
		params.Grade = &g
	}
	if filter.TeacherID != nil {
		params.TeacherID = toPgUUID(*filter.TeacherID)
	}
	if cursor != nil {
		params.AfterCreatedAt = pgTimestamptzValue(cursor.CreatedAt)
		params.AfterID = toPgUUID(cursor.ID)
	}

	rows, err := queriesFor(ctx, r.q).ListCourses(ctx, params)
	if err != nil {
		return nil, nil, apperror.Internal(err)
	}

	var next *pagination.Cursor
	if len(rows) > limit {
		last := rows[limit-1]
		next = &pagination.Cursor{
			CreatedAt: fromPgTimestamptzValue(last.CreatedAt),
			ID:        fromPgUUID(last.ID),
		}
		rows = rows[:limit]
	}

	courses := make([]*entity.Course, len(rows))
	for i, row := range rows {
		courses[i] = courseFromRow(row)
	}
	return courses, next, nil
}

// Update persists the mutable course fields (the query filters
// deleted_at IS NULL). Counters and audit columns are not touched here.
func (r *CourseRepo) Update(ctx context.Context, c *entity.Course) error {
	row, err := queriesFor(ctx, r.q).UpdateCourse(ctx, sqlcgen.UpdateCourseParams{
		ID:               toPgUUID(c.ID),
		Title:            c.Title,
		Description:      c.Description,
		ShortDescription: c.ShortDescription,
		Subject:          c.Subject,
		Grade:            int32(c.Grade),
		ThumbnailUrl:     c.ThumbnailURL,
		Price:            toPgNumeric(c.Price),
		Currency:         c.Currency,
		IsFree:           c.IsFree,
		IsPublished:      c.IsPublished,
		Language:         c.Language,
		Level:            string(c.Level),
		Metadata:         toJSONB(c.Metadata),
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolation {
			return apperror.Conflict("SLUG_TAKEN", err)
		}
		return mapCourseErr(err)
	}
	*c = *courseFromRow(row)
	return nil
}

// SoftDelete sets deleted_at; the row is never physically removed
// (Non-Negotiable #5).
func (r *CourseRepo) SoftDelete(ctx context.Context, id uuid.UUID) error {
	if err := queriesFor(ctx, r.q).SoftDeleteCourse(ctx, toPgUUID(id)); err != nil {
		return apperror.Internal(err)
	}
	return nil
}

// RecomputeCounters refreshes total_lessons / total_duration_seconds from the
// course's live lessons. Run it inside a WithinTx alongside the lesson change.
func (r *CourseRepo) RecomputeCounters(ctx context.Context, courseID uuid.UUID) error {
	if err := queriesFor(ctx, r.q).RecomputeCourseCounters(ctx, toPgUUID(courseID)); err != nil {
		return apperror.Internal(err)
	}
	return nil
}

// GetCourseWithSectionsAndLessons assembles the nested course view with three
// queries (course, its live sections, its live lessons) stitched in Go. Lessons
// are grouped under their section by section_id, preserving each query's
// order_index ordering.
func (r *CourseRepo) GetCourseWithSectionsAndLessons(ctx context.Context, courseID uuid.UUID) (*repository.CourseWithContent, error) {
	q := queriesFor(ctx, r.q)

	courseRow, err := q.GetCourseByID(ctx, toPgUUID(courseID))
	if err != nil {
		return nil, mapCourseErr(err)
	}

	sectionRows, err := q.ListSectionsByCourse(ctx, toPgUUID(courseID))
	if err != nil {
		return nil, apperror.Internal(err)
	}
	lessonRows, err := q.ListLessonsByCourse(ctx, toPgUUID(courseID))
	if err != nil {
		return nil, apperror.Internal(err)
	}

	// Bucket lessons by section so each section gets its own ordered slice.
	lessonsBySection := make(map[uuid.UUID][]*entity.Lesson, len(sectionRows))
	for _, lr := range lessonRows {
		l := lessonFromRow(lr)
		lessonsBySection[l.SectionID] = append(lessonsBySection[l.SectionID], l)
	}

	sections := make([]repository.SectionWithLessons, len(sectionRows))
	for i, sr := range sectionRows {
		s := sectionFromRow(sr)
		sections[i] = repository.SectionWithLessons{
			Section: s,
			Lessons: lessonsBySection[s.ID],
		}
	}

	return &repository.CourseWithContent{
		Course:   courseFromRow(courseRow),
		Sections: sections,
	}, nil
}

// mapCourseErr translates a query error: no rows → NotFound, else → Internal.
func mapCourseErr(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return apperror.NotFound("course not found", err)
	}
	return apperror.Internal(err)
}

// courseFromRow maps a sqlc CoursesCourse to the domain entity.
func courseFromRow(row sqlcgen.CoursesCourse) *entity.Course {
	return &entity.Course{
		ID:                   fromPgUUID(row.ID),
		TeacherID:            fromPgUUID(row.TeacherID),
		Title:                row.Title,
		Slug:                 row.Slug,
		Description:          row.Description,
		ShortDescription:     row.ShortDescription,
		Subject:              row.Subject,
		Grade:                int(row.Grade),
		ThumbnailURL:         row.ThumbnailUrl,
		Price:                fromPgNumeric(row.Price),
		Currency:             row.Currency,
		IsFree:               row.IsFree,
		IsPublished:          row.IsPublished,
		Language:             row.Language,
		Level:                entity.Level(row.Level),
		TotalDurationSeconds: int(row.TotalDurationSeconds),
		TotalLessons:         int(row.TotalLessons),
		RatingAvg:            fromPgNumeric(row.RatingAvg),
		EnrollmentCount:      int(row.EnrollmentCount),
		Metadata:             fromJSONB(row.Metadata),
		CreatedAt:            fromPgTimestamptzValue(row.CreatedAt),
		UpdatedAt:            fromPgTimestamptzValue(row.UpdatedAt),
		DeletedAt:            fromPgTimestamptz(row.DeletedAt),
	}
}

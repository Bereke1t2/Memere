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

// SectionRepo is the sqlc-backed implementation of repository.SectionRepository.
type SectionRepo struct {
	q *sqlcgen.Queries
}

var _ repository.SectionRepository = (*SectionRepo)(nil)

// NewSectionRepo builds a SectionRepo over a pgx pool.
func NewSectionRepo(pool *pgxpool.Pool) *SectionRepo {
	return &SectionRepo{q: sqlcgen.New(pool)}
}

// Create inserts a section and reflects DB-assigned fields back onto the entity.
func (r *SectionRepo) Create(ctx context.Context, s *entity.CourseSection) error {
	row, err := queriesFor(ctx, r.q).CreateSection(ctx, sqlcgen.CreateSectionParams{
		CourseID:    toPgUUID(s.CourseID),
		Title:       s.Title,
		Description: s.Description,
		OrderIndex:  int32(s.OrderIndex),
		IsPublished: s.IsPublished,
	})
	if err != nil {
		return apperror.Internal(err)
	}
	*s = *sectionFromRow(row)
	return nil
}

// FindByID returns the section or apperror.NotFound (soft-deleted excluded).
func (r *SectionRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.CourseSection, error) {
	row, err := queriesFor(ctx, r.q).GetSectionByID(ctx, toPgUUID(id))
	if err != nil {
		return nil, mapSectionErr(err)
	}
	return sectionFromRow(row), nil
}

// ListByCourse returns a course's live sections ordered by order_index.
func (r *SectionRepo) ListByCourse(ctx context.Context, courseID uuid.UUID) ([]*entity.CourseSection, error) {
	rows, err := queriesFor(ctx, r.q).ListSectionsByCourse(ctx, toPgUUID(courseID))
	if err != nil {
		return nil, apperror.Internal(err)
	}
	sections := make([]*entity.CourseSection, len(rows))
	for i, row := range rows {
		sections[i] = sectionFromRow(row)
	}
	return sections, nil
}

// Update persists the mutable section fields (the query filters deleted_at IS
// NULL).
func (r *SectionRepo) Update(ctx context.Context, s *entity.CourseSection) error {
	row, err := queriesFor(ctx, r.q).UpdateSection(ctx, sqlcgen.UpdateSectionParams{
		ID:          toPgUUID(s.ID),
		Title:       s.Title,
		Description: s.Description,
		OrderIndex:  int32(s.OrderIndex),
		IsPublished: s.IsPublished,
	})
	if err != nil {
		return mapSectionErr(err)
	}
	*s = *sectionFromRow(row)
	return nil
}

// SoftDelete sets deleted_at (Non-Negotiable #5).
func (r *SectionRepo) SoftDelete(ctx context.Context, id uuid.UUID) error {
	if err := queriesFor(ctx, r.q).SoftDeleteSection(ctx, toPgUUID(id)); err != nil {
		return apperror.Internal(err)
	}
	return nil
}

// MaxOrderIndex returns the highest order_index among a course's live sections,
// or -1 when none exist.
func (r *SectionRepo) MaxOrderIndex(ctx context.Context, courseID uuid.UUID) (int, error) {
	max, err := queriesFor(ctx, r.q).MaxSectionOrderIndex(ctx, toPgUUID(courseID))
	if err != nil {
		return 0, apperror.Internal(err)
	}
	return int(max), nil
}

// mapSectionErr translates a query error: no rows → NotFound, else → Internal.
func mapSectionErr(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return apperror.NotFound("section not found", err)
	}
	return apperror.Internal(err)
}

// sectionFromRow maps a sqlc CoursesCourseSection to the domain entity.
func sectionFromRow(row sqlcgen.CoursesCourseSection) *entity.CourseSection {
	return &entity.CourseSection{
		ID:          fromPgUUID(row.ID),
		CourseID:    fromPgUUID(row.CourseID),
		Title:       row.Title,
		Description: row.Description,
		OrderIndex:  int(row.OrderIndex),
		IsPublished: row.IsPublished,
		CreatedAt:   fromPgTimestamptzValue(row.CreatedAt),
		UpdatedAt:   fromPgTimestamptzValue(row.UpdatedAt),
		DeletedAt:   fromPgTimestamptz(row.DeletedAt),
	}
}

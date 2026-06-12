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

// EnrollmentRepo is the sqlc-backed implementation of
// repository.EnrollmentRepository.
type EnrollmentRepo struct {
	q *sqlcgen.Queries
}

var _ repository.EnrollmentRepository = (*EnrollmentRepo)(nil)

// NewEnrollmentRepo builds an EnrollmentRepo over a pgx pool.
func NewEnrollmentRepo(pool *pgxpool.Pool) *EnrollmentRepo {
	return &EnrollmentRepo{q: sqlcgen.New(pool)}
}

// Create grants access exactly once. The CreateEnrollment query is
// ON CONFLICT (student_id, course_id) DO NOTHING, so a re-processed payment
// produces no row (pgx.ErrNoRows): that is the idempotent "already enrolled"
// case, reported as success — never an error — with e left unchanged.
func (r *EnrollmentRepo) Create(ctx context.Context, e *entity.Enrollment) error {
	id := e.ID
	if id == uuid.Nil {
		id = uuid.New()
	}
	row, err := queriesFor(ctx, r.q).CreateEnrollment(ctx, sqlcgen.CreateEnrollmentParams{
		ID:        toPgUUID(id),
		StudentID: toPgUUID(e.StudentID),
		CourseID:  toPgUUID(e.CourseID),
		Source:    string(e.Source),
		ExpiresAt: toPgTimestamptz(e.ExpiresAt),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil // already enrolled — exactly-once grant
		}
		return apperror.Internal(err)
	}
	*e = *enrollmentFromRow(row)
	return nil
}

func (r *EnrollmentRepo) Exists(ctx context.Context, studentID, courseID uuid.UUID) (bool, error) {
	ok, err := queriesFor(ctx, r.q).EnrollmentExists(ctx, sqlcgen.EnrollmentExistsParams{
		StudentID: toPgUUID(studentID),
		CourseID:  toPgUUID(courseID),
	})
	if err != nil {
		return false, apperror.Internal(err)
	}
	return ok, nil
}

func (r *EnrollmentRepo) GetActiveForStudent(ctx context.Context, studentID, courseID uuid.UUID) (*entity.Enrollment, error) {
	row, err := queriesFor(ctx, r.q).GetActiveEnrollment(ctx, sqlcgen.GetActiveEnrollmentParams{
		StudentID: toPgUUID(studentID),
		CourseID:  toPgUUID(courseID),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NotFound("enrollment not found", err)
		}
		return nil, apperror.Internal(err)
	}
	return enrollmentFromRow(row), nil
}

func (r *EnrollmentRepo) ListByStudent(ctx context.Context, studentID uuid.UUID, limit int) ([]*entity.Enrollment, error) {
	rows, err := queriesFor(ctx, r.q).ListEnrollmentsByStudent(ctx, sqlcgen.ListEnrollmentsByStudentParams{
		StudentID: toPgUUID(studentID),
		Limit:     int32(limit),
	})
	if err != nil {
		return nil, apperror.Internal(err)
	}
	out := make([]*entity.Enrollment, len(rows))
	for i, row := range rows {
		out[i] = enrollmentFromRow(row)
	}
	return out, nil
}

func enrollmentFromRow(row sqlcgen.PaymentsEnrollment) *entity.Enrollment {
	return &entity.Enrollment{
		ID:         fromPgUUID(row.ID),
		StudentID:  fromPgUUID(row.StudentID),
		CourseID:   fromPgUUID(row.CourseID),
		Source:     entity.EnrollmentSource(row.Source),
		EnrolledAt: fromPgTimestamptzValue(row.EnrolledAt),
		ExpiresAt:  fromPgTimestamptz(row.ExpiresAt),
		CreatedAt:  fromPgTimestamptzValue(row.CreatedAt),
		UpdatedAt:  fromPgTimestamptzValue(row.UpdatedAt),
	}
}

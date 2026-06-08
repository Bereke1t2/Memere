package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/internal/repository/postgres/sqlcgen"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// ExamAttemptRepo is the sqlc-backed implementation of
// repository.ExamAttemptRepository.
type ExamAttemptRepo struct {
	q *sqlcgen.Queries
}

var _ repository.ExamAttemptRepository = (*ExamAttemptRepo)(nil)

// NewExamAttemptRepo builds an ExamAttemptRepo over a pgx pool.
func NewExamAttemptRepo(pool *pgxpool.Pool) *ExamAttemptRepo {
	return &ExamAttemptRepo{q: sqlcgen.New(pool)}
}

func (r *ExamAttemptRepo) Create(ctx context.Context, a *entity.ExamAttempt) error {
	row, err := queriesFor(ctx, r.q).CreateExamAttempt(ctx, sqlcgen.CreateExamAttemptParams{
		ExamID:          toPgUUID(a.ExamID),
		StudentID:       toPgUUID(a.StudentID),
		AnswersSnapshot: toJSONB(a.AnswersSnapshot),
		Status:          string(a.Status),
	})
	if err != nil {
		return apperror.Internal(err)
	}
	*a = *examAttemptFromRow(row)
	return nil
}

func (r *ExamAttemptRepo) FindByID(ctx context.Context, id, studentID uuid.UUID) (*entity.ExamAttempt, error) {
	row, err := queriesFor(ctx, r.q).GetExamAttemptByID(ctx, sqlcgen.GetExamAttemptByIDParams{
		ID:        toPgUUID(id),
		StudentID: toPgUUID(studentID),
	})
	if err != nil {
		return nil, mapExamAttemptErr(err)
	}
	return examAttemptFromRow(row), nil
}

func (r *ExamAttemptRepo) GetActive(ctx context.Context, studentID, examID uuid.UUID) (*entity.ExamAttempt, error) {
	row, err := queriesFor(ctx, r.q).GetActiveExamAttempt(ctx, sqlcgen.GetActiveExamAttemptParams{
		StudentID: toPgUUID(studentID),
		ExamID:    toPgUUID(examID),
	})
	if err != nil {
		return nil, mapExamAttemptErr(err)
	}
	return examAttemptFromRow(row), nil
}

func (r *ExamAttemptRepo) ListByStudent(ctx context.Context, studentID uuid.UUID) ([]*entity.ExamAttempt, error) {
	rows, err := queriesFor(ctx, r.q).ListExamAttemptsByStudent(ctx, toPgUUID(studentID))
	if err != nil {
		return nil, apperror.Internal(err)
	}
	attempts := make([]*entity.ExamAttempt, len(rows))
	for i, row := range rows {
		attempts[i] = examAttemptFromRow(row)
	}
	return attempts, nil
}

func (r *ExamAttemptRepo) ListExpired(ctx context.Context, now time.Time, limit int) ([]*entity.ExamAttempt, error) {
	rows, err := queriesFor(ctx, r.q).FindExpiredExamAttempts(ctx, sqlcgen.FindExpiredExamAttemptsParams{
		StartedAt: pgTimestamptzValue(now),
		Limit:     int32(limit),
	})
	if err != nil {
		return nil, apperror.Internal(err)
	}
	attempts := make([]*entity.ExamAttempt, len(rows))
	for i, row := range rows {
		attempts[i] = examAttemptFromRow(row)
	}
	return attempts, nil
}

func (r *ExamAttemptRepo) Update(ctx context.Context, a *entity.ExamAttempt) error {
	row, err := queriesFor(ctx, r.q).UpdateExamAttempt(ctx, sqlcgen.UpdateExamAttemptParams{
		ID:              toPgUUID(a.ID),
		AnswersSnapshot: toJSONB(a.AnswersSnapshot),
		Status:          string(a.Status),
		SubmittedAt:     toPgTimestamptz(a.SubmittedAt),
	})
	if err != nil {
		return mapExamAttemptErr(err)
	}
	*a = *examAttemptFromRow(row)
	return nil
}

func (r *ExamAttemptRepo) Grade(ctx context.Context, a *entity.ExamAttempt) error {
	row, err := queriesFor(ctx, r.q).GradeExamAttempt(ctx, sqlcgen.GradeExamAttemptParams{
		ID:         toPgUUID(a.ID),
		Score:      toPgNumericPtr(a.Score),
		Percentage: toPgNumericPtr(a.Percentage),
	})
	if err != nil {
		return mapExamAttemptErr(err)
	}
	*a = *examAttemptFromRow(row)
	return nil
}

func mapExamAttemptErr(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return apperror.NotFound("exam attempt not found", err)
	}
	return apperror.Internal(err)
}

func examAttemptFromRow(row sqlcgen.CoursesExamAttempt) *entity.ExamAttempt {
	return &entity.ExamAttempt{
		ID:              fromPgUUID(row.ID),
		ExamID:          fromPgUUID(row.ExamID),
		StudentID:       fromPgUUID(row.StudentID),
		StartedAt:       fromPgTimestamptzValue(row.StartedAt),
		SubmittedAt:     fromPgTimestamptz(row.SubmittedAt),
		Score:           fromPgNumericPtr(row.Score),
		Percentage:      fromPgNumericPtr(row.Percentage),
		AnswersSnapshot: fromJSONB(row.AnswersSnapshot),
		Status:          entity.AttemptStatus(row.Status),
		CreatedAt:       fromPgTimestamptzValue(row.CreatedAt),
		UpdatedAt:       fromPgTimestamptzValue(row.UpdatedAt),
	}
}

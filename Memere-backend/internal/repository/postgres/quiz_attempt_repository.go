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

// QuizAttemptRepo is the sqlc-backed implementation of
// repository.QuizAttemptRepository.
type QuizAttemptRepo struct {
	q *sqlcgen.Queries
}

var _ repository.QuizAttemptRepository = (*QuizAttemptRepo)(nil)

// NewQuizAttemptRepo builds a QuizAttemptRepo over a pgx pool.
func NewQuizAttemptRepo(pool *pgxpool.Pool) *QuizAttemptRepo {
	return &QuizAttemptRepo{q: sqlcgen.New(pool)}
}

func (r *QuizAttemptRepo) Create(ctx context.Context, a *entity.QuizAttempt) error {
	row, err := queriesFor(ctx, r.q).CreateQuizAttempt(ctx, sqlcgen.CreateQuizAttemptParams{
		QuizID:          toPgUUID(a.QuizID),
		StudentID:       toPgUUID(a.StudentID),
		AttemptNumber:   int32(a.AttemptNumber),
		ExpiresAt:       toPgTimestamptz(a.ExpiresAt),
		QuestionOrder:   toJSONB(a.QuestionOrder),
		AnswersSnapshot: toJSONB(a.AnswersSnapshot),
		Status:          string(a.Status),
	})
	if err != nil {
		return apperror.Internal(err)
	}
	*a = *quizAttemptFromRow(row)
	return nil
}

func (r *QuizAttemptRepo) FindByID(ctx context.Context, id, studentID uuid.UUID) (*entity.QuizAttempt, error) {
	row, err := queriesFor(ctx, r.q).GetQuizAttemptByID(ctx, sqlcgen.GetQuizAttemptByIDParams{
		ID:        toPgUUID(id),
		StudentID: toPgUUID(studentID),
	})
	if err != nil {
		return nil, mapQuizAttemptErr(err)
	}
	return quizAttemptFromRow(row), nil
}

func (r *QuizAttemptRepo) GetActive(ctx context.Context, studentID, quizID uuid.UUID) (*entity.QuizAttempt, error) {
	row, err := queriesFor(ctx, r.q).GetActiveQuizAttempt(ctx, sqlcgen.GetActiveQuizAttemptParams{
		StudentID: toPgUUID(studentID),
		QuizID:    toPgUUID(quizID),
	})
	if err != nil {
		return nil, mapQuizAttemptErr(err)
	}
	return quizAttemptFromRow(row), nil
}

func (r *QuizAttemptRepo) ListByStudentAndQuiz(ctx context.Context, studentID, quizID uuid.UUID) ([]*entity.QuizAttempt, error) {
	rows, err := queriesFor(ctx, r.q).ListQuizAttemptsByStudentAndQuiz(ctx, sqlcgen.ListQuizAttemptsByStudentAndQuizParams{
		StudentID: toPgUUID(studentID),
		QuizID:    toPgUUID(quizID),
	})
	if err != nil {
		return nil, apperror.Internal(err)
	}
	attempts := make([]*entity.QuizAttempt, len(rows))
	for i, row := range rows {
		attempts[i] = quizAttemptFromRow(row)
	}
	return attempts, nil
}

func (r *QuizAttemptRepo) CountByStudentAndQuiz(ctx context.Context, studentID, quizID uuid.UUID) (int, error) {
	n, err := queriesFor(ctx, r.q).CountQuizAttempts(ctx, sqlcgen.CountQuizAttemptsParams{
		StudentID: toPgUUID(studentID),
		QuizID:    toPgUUID(quizID),
	})
	if err != nil {
		return 0, apperror.Internal(err)
	}
	return int(n), nil
}

func (r *QuizAttemptRepo) ListExpired(ctx context.Context, now time.Time, limit int) ([]*entity.QuizAttempt, error) {
	rows, err := queriesFor(ctx, r.q).ListExpiredQuizAttempts(ctx, sqlcgen.ListExpiredQuizAttemptsParams{
		ExpiresAt: pgTimestamptzValue(now),
		Limit:     int32(limit),
	})
	if err != nil {
		return nil, apperror.Internal(err)
	}
	attempts := make([]*entity.QuizAttempt, len(rows))
	for i, row := range rows {
		attempts[i] = quizAttemptFromRow(row)
	}
	return attempts, nil
}

func (r *QuizAttemptRepo) Update(ctx context.Context, a *entity.QuizAttempt) error {
	row, err := queriesFor(ctx, r.q).UpdateQuizAttempt(ctx, sqlcgen.UpdateQuizAttemptParams{
		ID:              toPgUUID(a.ID),
		AnswersSnapshot: toJSONB(a.AnswersSnapshot),
		Status:          string(a.Status),
		SubmittedAt:     toPgTimestamptz(a.SubmittedAt),
	})
	if err != nil {
		return mapQuizAttemptErr(err)
	}
	*a = *quizAttemptFromRow(row)
	return nil
}

// ClaimForGrading runs the conditional UPDATE ... WHERE status='in_progress'.
// pgx.ErrNoRows means the row was already moved by another writer — reported as
// claimed=false, not an error.
func (r *QuizAttemptRepo) ClaimForGrading(ctx context.Context, a *entity.QuizAttempt) (bool, error) {
	row, err := queriesFor(ctx, r.q).ClaimQuizAttemptForGrading(ctx, sqlcgen.ClaimQuizAttemptForGradingParams{
		ID:              toPgUUID(a.ID),
		Status:          string(a.Status),
		AnswersSnapshot: toJSONB(a.AnswersSnapshot),
		SubmittedAt:     toPgTimestamptz(a.SubmittedAt),
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, apperror.Internal(err)
	}
	*a = *quizAttemptFromRow(row)
	return true, nil
}

func (r *QuizAttemptRepo) Grade(ctx context.Context, a *entity.QuizAttempt) error {
	row, err := queriesFor(ctx, r.q).GradeQuizAttempt(ctx, sqlcgen.GradeQuizAttemptParams{
		ID:         toPgUUID(a.ID),
		Score:      toPgNumericPtr(a.Score),
		Percentage: toPgNumericPtr(a.Percentage),
		Passed:     a.Passed,
	})
	if err != nil {
		return mapQuizAttemptErr(err)
	}
	*a = *quizAttemptFromRow(row)
	return nil
}

func (r *QuizAttemptRepo) SumBestScores(ctx context.Context, studentID uuid.UUID) (repository.StudentScoreTotals, error) {
	row, err := queriesFor(ctx, r.q).GetStudentQuizPoints(ctx, toPgUUID(studentID))
	if err != nil {
		return repository.StudentScoreTotals{}, apperror.Internal(err)
	}
	return repository.StudentScoreTotals{
		Points:        row.TotalPoints,
		Count:         row.QuizCount,
		AvgPercentage: row.AvgPercentage,
	}, nil
}

func mapQuizAttemptErr(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return apperror.NotFound("quiz attempt not found", err)
	}
	return apperror.Internal(err)
}

func quizAttemptFromRow(row sqlcgen.CoursesQuizAttempt) *entity.QuizAttempt {
	return &entity.QuizAttempt{
		ID:              fromPgUUID(row.ID),
		QuizID:          fromPgUUID(row.QuizID),
		StudentID:       fromPgUUID(row.StudentID),
		AttemptNumber:   int(row.AttemptNumber),
		StartedAt:       fromPgTimestamptzValue(row.StartedAt),
		SubmittedAt:     fromPgTimestamptz(row.SubmittedAt),
		ExpiresAt:       fromPgTimestamptz(row.ExpiresAt),
		Score:           fromPgNumericPtr(row.Score),
		Percentage:      fromPgNumericPtr(row.Percentage),
		Passed:          row.Passed,
		QuestionOrder:   fromJSONB(row.QuestionOrder),
		AnswersSnapshot: fromJSONB(row.AnswersSnapshot),
		Status:          entity.AttemptStatus(row.Status),
		CreatedAt:       fromPgTimestamptzValue(row.CreatedAt),
		UpdatedAt:       fromPgTimestamptzValue(row.UpdatedAt),
	}
}

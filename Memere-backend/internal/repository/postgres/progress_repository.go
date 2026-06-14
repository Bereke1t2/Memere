package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/internal/repository/postgres/sqlcgen"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// ProgressRepo is the sqlc-backed implementation of repository.ProgressRepository.
type ProgressRepo struct {
	q *sqlcgen.Queries
}

var _ repository.ProgressRepository = (*ProgressRepo)(nil)

func NewProgressRepo(pool *pgxpool.Pool) *ProgressRepo {
	return &ProgressRepo{q: sqlcgen.New(pool)}
}

func (r *ProgressRepo) UpsertLessonProgress(ctx context.Context, p *entity.LessonProgress) (*entity.LessonProgress, error) {
	row, err := queriesFor(ctx, r.q).UpsertLessonProgress(ctx, sqlcgen.UpsertLessonProgressParams{
		StudentID:            toPgUUID(p.StudentID),
		LessonID:             toPgUUID(p.LessonID),
		CourseID:             toPgUUID(p.CourseID),
		IsCompleted:          p.IsCompleted,
		VideoProgressSeconds: int32(p.VideoProgressSeconds),
	})
	if err != nil {
		return nil, apperror.Internal(err)
	}
	return progressRowToEntity(row), nil
}

func (r *ProgressRepo) GetLessonProgress(ctx context.Context, studentID, lessonID uuid.UUID) (*entity.LessonProgress, error) {
	row, err := queriesFor(ctx, r.q).GetLessonProgress(ctx, toPgUUID(studentID), toPgUUID(lessonID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NotFound("lesson progress not found", nil)
		}
		return nil, apperror.Internal(err)
	}
	return progressRowToEntity(row), nil
}

func (r *ProgressRepo) ListByCourse(ctx context.Context, studentID, courseID uuid.UUID) ([]*entity.LessonProgress, error) {
	rows, err := queriesFor(ctx, r.q).ListByCourse(ctx, toPgUUID(studentID), toPgUUID(courseID))
	if err != nil {
		return nil, apperror.Internal(err)
	}
	out := make([]*entity.LessonProgress, 0, len(rows))
	for _, row := range rows {
		out = append(out, progressRowToEntity(row))
	}
	return out, nil
}

func (r *ProgressRepo) RecomputeCourseProgress(ctx context.Context, studentID, courseID uuid.UUID) (*entity.CourseProgress, error) {
	counts, err := queriesFor(ctx, r.q).CountCourseCompletion(ctx, toPgUUID(studentID), toPgUUID(courseID))
	if err != nil {
		return nil, apperror.Internal(err)
	}
	total := int(counts.Total)
	completed := int(counts.Completed)

	var pct decimal.Decimal
	if total > 0 {
		pct = decimal.NewFromInt(int64(completed)).Div(decimal.NewFromInt(int64(total))).Mul(decimal.NewFromInt(100)).Round(2)
	}

	var completedAt pgtype.Timestamptz
	if total > 0 && completed >= total {
		now := time.Now().UTC()
		completedAt = pgTimestamptzValue(now)
	}

	if err := queriesFor(ctx, r.q).UpsertCourseProgress(ctx, sqlcgen.UpsertCourseProgressParams{
		StudentID:        toPgUUID(studentID),
		CourseID:         toPgUUID(courseID),
		CompletedLessons: int32(completed),
		TotalLessons:     int32(total),
		PercentComplete:  toPgDecimal(pct),
		CompletedAt:      completedAt,
	}); err != nil {
		return nil, apperror.Internal(err)
	}

	cp := &entity.CourseProgress{
		StudentID:        studentID,
		CourseID:         courseID,
		CompletedLessons: completed,
		TotalLessons:     total,
		PercentComplete:  pct,
	}
	if completedAt.Valid {
		t := completedAt.Time
		cp.CompletedAt = &t
	}
	return cp, nil
}

func (r *ProgressRepo) GetCourseProgress(ctx context.Context, studentID, courseID uuid.UUID) (*entity.CourseProgress, error) {
	row, err := queriesFor(ctx, r.q).GetCourseProgress(ctx, toPgUUID(studentID), toPgUUID(courseID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NotFound("course progress not found", nil)
		}
		return nil, apperror.Internal(err)
	}
	return courseProgressRowToEntity(row), nil
}

func (r *ProgressRepo) GetStreak(ctx context.Context, studentID uuid.UUID) (*entity.StudyStreak, error) {
	row, err := queriesFor(ctx, r.q).GetStreak(ctx, toPgUUID(studentID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// No streak row yet — return a zero-value streak, never NotFound.
			return &entity.StudyStreak{StudentID: studentID}, nil
		}
		return nil, apperror.Internal(err)
	}
	return streakRowToEntity(row), nil
}

func (r *ProgressRepo) ListStreakAtRisk(ctx context.Context, cutoff time.Time, today time.Time, limit int) ([]*entity.StudyStreak, error) {
	rows, err := queriesFor(ctx, r.q).ListStreakAtRisk(ctx, sqlcgen.ListStreakAtRiskParams{
		Cutoff: pgtype.Date{Time: cutoff, Valid: true},
		Today:  pgtype.Date{Time: today, Valid: true},
		Limit:  int32(limit),
	})
	if err != nil {
		return nil, apperror.Internal(err)
	}
	out := make([]*entity.StudyStreak, 0, len(rows))
	for _, row := range rows {
		out = append(out, streakRowToEntity(row))
	}
	return out, nil
}

func (r *ProgressRepo) SetLastWarnedDate(ctx context.Context, studentID uuid.UUID, date time.Time) error {
	if err := queriesFor(ctx, r.q).SetLastWarnedDate(ctx, toPgUUID(studentID), pgtype.Date{Time: date, Valid: true}); err != nil {
		return apperror.Internal(err)
	}
	return nil
}

func (r *ProgressRepo) UpsertStreak(ctx context.Context, st *entity.StudyStreak) error {
	var pgDate pgtype.Date
	if st.LastStudyDate != nil {
		pgDate = pgtype.Date{Time: *st.LastStudyDate, Valid: true}
	}
	if err := queriesFor(ctx, r.q).UpsertStreak(ctx, sqlcgen.UpsertStreakParams{
		StudentID:     toPgUUID(st.StudentID),
		CurrentStreak: int32(st.CurrentStreak),
		LongestStreak: int32(st.LongestStreak),
		LastStudyDate: pgDate,
	}); err != nil {
		return apperror.Internal(err)
	}
	return nil
}

// --- mappers ---

func progressRowToEntity(r sqlcgen.ProgressProgress) *entity.LessonProgress {
	return &entity.LessonProgress{
		ID:                   fromPgUUID(r.ID),
		StudentID:            fromPgUUID(r.StudentID),
		LessonID:             fromPgUUID(r.LessonID),
		CourseID:             fromPgUUID(r.CourseID),
		IsCompleted:          r.IsCompleted,
		CompletedAt:          fromPgTimestamptz(r.CompletedAt),
		VideoProgressSeconds: int(r.VideoProgressSeconds),
		LastAccessedAt:       fromPgTimestamptz(r.LastAccessedAt),
		CreatedAt:            fromPgTimestamptzValue(r.CreatedAt),
		UpdatedAt:            fromPgTimestamptzValue(r.UpdatedAt),
		DeletedAt:            fromPgTimestamptz(r.DeletedAt),
	}
}

func courseProgressRowToEntity(r sqlcgen.ProgressCourseProgress) *entity.CourseProgress {
	return &entity.CourseProgress{
		StudentID:        fromPgUUID(r.StudentID),
		CourseID:         fromPgUUID(r.CourseID),
		CompletedLessons: int(r.CompletedLessons),
		TotalLessons:     int(r.TotalLessons),
		PercentComplete:  fromPgDecimal(r.PercentComplete),
		CompletedAt:      fromPgTimestamptz(r.CompletedAt),
		UpdatedAt:        fromPgTimestamptzValue(r.UpdatedAt),
	}
}

func streakRowToEntity(r sqlcgen.ProgressStudyStreak) *entity.StudyStreak {
	st := &entity.StudyStreak{
		StudentID:     fromPgUUID(r.StudentID),
		CurrentStreak: int(r.CurrentStreak),
		LongestStreak: int(r.LongestStreak),
		UpdatedAt:     fromPgTimestamptzValue(r.UpdatedAt),
	}
	if r.LastStudyDate.Valid {
		t := r.LastStudyDate.Time
		st.LastStudyDate = &t
	}
	return st
}

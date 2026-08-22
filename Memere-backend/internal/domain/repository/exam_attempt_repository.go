package repository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
)

// ExamAttemptRepository persists exam sittings (spec §9.2). started_at is the
// server-side timer source of truth; ListExpired backs the auto-submit sweeper.
// Every read filters by the authenticated student_id to prevent IDOR
// (Non-Negotiable #7).
type ExamAttemptRepository interface {
	Create(ctx context.Context, a *entity.ExamAttempt) error
	// FindByID scopes the lookup to studentID to prevent IDOR (Non-Negotiable #7).
	FindByID(ctx context.Context, id, studentID uuid.UUID) (*entity.ExamAttempt, error)
	// GetActive returns the student's current in-progress sitting of an exam, or
	// apperror.NotFound if none.
	GetActive(ctx context.Context, studentID, examID uuid.UUID) (*entity.ExamAttempt, error)
	ListByStudent(ctx context.Context, studentID uuid.UUID) ([]*entity.ExamAttempt, error)
	// ListExpired returns in-progress attempts whose timer (started_at +
	// exams.duration_minutes) has elapsed as of now, for the sweeper (spec §9.2).
	ListExpired(ctx context.Context, now time.Time, limit int) ([]*entity.ExamAttempt, error)
	// Update persists auto-save state and status transitions without touching the
	// immutable started_at timer column.
	Update(ctx context.Context, a *entity.ExamAttempt) error
	// ClaimForGrading atomically transitions the attempt out of in_progress to
	// a.Status (submitted/expired) ONLY if it is still in_progress, reporting
	// claimed=false when a racing submit or the sweeper already moved it (race
	// guard, spec §9.2).
	ClaimForGrading(ctx context.Context, a *entity.ExamAttempt) (claimed bool, err error)
	Grade(ctx context.Context, a *entity.ExamAttempt) error
	// ListGradedBySubject returns the student's graded attempts for a subject,
	// oldest first, for the §9.3 score trend.
	ListGradedBySubject(ctx context.Context, studentID uuid.UUID, subject string) ([]*entity.ExamAttempt, error)
	// Stats returns aggregate analytics for an exam (§9.3): graded attempt count,
	// average percentage, and pass count (score >= the exam's pass_marks).
	Stats(ctx context.Context, examID uuid.UUID) (ExamAttemptStats, error)
	// SumBestScores returns the student's cumulative exam points: the best graded
	// score per exam, summed, with the distinct-exam count and average of those
	// best percentages. Best-per-exam so retakes never double-count (§9.3).
	SumBestScores(ctx context.Context, studentID uuid.UUID) (StudentScoreTotals, error)
}

// ExamAttemptStats is the aggregate result behind GetExamStats (§9.3).
type ExamAttemptStats struct {
	TotalAttempts int
	AvgPercentage float64
	PassedCount   int
}

// StudentScoreTotals is a student's cumulative points across quizzes or exams:
// the best graded score per item, summed, plus the number of distinct items
// completed and the average of those best percentages (§9.3). Shared by the
// quiz and exam attempt repositories.
type StudentScoreTotals struct {
	Points        float64
	Count         int64
	AvgPercentage float64
}

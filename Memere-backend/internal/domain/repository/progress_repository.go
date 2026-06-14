package repository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
)

// ProgressRepository is the persistence port for progress tracking. All reads
// and writes are scoped to the authenticated student — callers must supply the
// correct studentID; there is no server-side IDOR check at this layer.
type ProgressRepository interface {
	// UpsertLessonProgress creates or updates the progress row for one lesson.
	// video_progress_seconds advances monotonically (GREATEST); completed_at is
	// set once and never cleared. Returns the updated row.
	UpsertLessonProgress(ctx context.Context, p *entity.LessonProgress) (*entity.LessonProgress, error)

	// GetLessonProgress returns the progress row for one student+lesson pair, or
	// apperror.NotFound if none exists.
	GetLessonProgress(ctx context.Context, studentID, lessonID uuid.UUID) (*entity.LessonProgress, error)

	// ListByCourse returns all lesson progress rows for a student in a course,
	// most-recently-updated first.
	ListByCourse(ctx context.Context, studentID, courseID uuid.UUID) ([]*entity.LessonProgress, error)

	// RecomputeCourseProgress counts published lessons vs completed rows, writes
	// the course_progress rollup (upsert), and returns the updated rollup. Sets
	// completed_at on the first time total == completed (100%).
	RecomputeCourseProgress(ctx context.Context, studentID, courseID uuid.UUID) (*entity.CourseProgress, error)

	// GetCourseProgress returns the cached rollup for one student+course pair, or
	// apperror.NotFound if the student has not started the course.
	GetCourseProgress(ctx context.Context, studentID, courseID uuid.UUID) (*entity.CourseProgress, error)

	// GetStreak returns the study streak for a student, or a zero-value
	// StudyStreak when no row exists yet (never returns NotFound).
	GetStreak(ctx context.Context, studentID uuid.UUID) (*entity.StudyStreak, error)

	// UpsertStreak creates or updates the streak row for a student.
	UpsertStreak(ctx context.Context, st *entity.StudyStreak) error

	// ListStreakAtRisk returns students whose streak is at risk of breaking:
	// current_streak > 0 and last_study_date < cutoff, skipping any student
	// already warned on today's date (idempotency for the engagement sweeper).
	ListStreakAtRisk(ctx context.Context, cutoff time.Time, today time.Time, limit int) ([]*entity.StudyStreak, error)

	// SetLastWarnedDate records the date a streak-warning notification was sent
	// for a student so the engagement sweeper doesn't double-warn within the day.
	SetLastWarnedDate(ctx context.Context, studentID uuid.UUID, date time.Time) error
}

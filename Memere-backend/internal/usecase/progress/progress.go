// Package progress implements the Phase 5 progress tracking engine: per-lesson
// completion, video watch position, study streaks, course completion %, and the
// student dashboard. All reads and writes are scoped to the authenticated caller
// (IDOR prevention §2.7 Non-Negotiable).
package progress

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/access"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// videoCompletionThreshold is the watch fraction at which a video lesson is
// considered complete (95%, matching typical UX conventions and the spec).
const videoCompletionThreshold = 0.95

// Notifier is the port for post-progress notifications. Methods must never
// block or fail the calling request — they are called after the progress row
// is committed. The real notification.Hooks facade is injected in main.go.
type Notifier interface {
	CertificateReady(ctx context.Context, studentID, courseID uuid.UUID)
}

// noopNotifier satisfies Notifier with empty implementations. Used only in
// tests that don't wire the real notification.Hooks facade.
type noopNotifier struct{}

func (noopNotifier) CertificateReady(context.Context, uuid.UUID, uuid.UUID) {}

// CourseAccess is the dependency on the access gate. *access.Service satisfies it.
type CourseAccess interface {
	RequireFullAccess(ctx context.Context, actor access.Actor, courseID uuid.UUID) error
}


// Actor is the authenticated caller's identity, passed down from the HTTP layer.
type Actor struct {
	UserID uuid.UUID
	Role   entity.Role
}

// Config holds runtime settings for the progress service.
type Config struct {
	// StreakTZ is the IANA timezone name for streak day-boundary calculations
	// (spec decision: configurable, default Africa/Addis_Ababa).
	StreakTZ string
}

// DashboardItem is one enrolled course's progress summary for the dashboard.
type DashboardItem struct {
	CourseID         uuid.UUID
	CourseTitle      string
	PercentComplete  float64
	CompletedAt      *time.Time
	CompletedLessons int
	TotalLessons     int
}

// Dashboard is the student self-service view returned by GetMyDashboard.
type Dashboard struct {
	Items         []DashboardItem
	CurrentStreak int
	LongestStreak int
	WeakTopics    []string // phase 2 §9.3 weak-area feed (nil when unavailable)
}

// Service is the stateless orchestrator for the progress domain.
type Service struct {
	repo       repository.ProgressRepository
	courses    repository.CourseRepository
	enroll     repository.EnrollmentRepository
	access     CourseAccess
	lessonRepo repository.LessonRepository
	notify     Notifier
	clock      func() time.Time
	loc        *time.Location
}

// NewService constructs a progress Service. notify may be nil (defaults to a
// no-op stub). clock may be nil (defaults to time.Now). cfg.StreakTZ defaults
// to Africa/Addis_Ababa when empty.
func NewService(
	repo repository.ProgressRepository,
	courses repository.CourseRepository,
	enroll repository.EnrollmentRepository,
	accessSvc CourseAccess,
	lessonRepo repository.LessonRepository,
	notify Notifier,
	cfg Config,
	clock func() time.Time,
) *Service {
	if notify == nil {
		notify = noopNotifier{}
	}
	if clock == nil {
		clock = time.Now
	}
	tz := cfg.StreakTZ
	if tz == "" {
		tz = "Africa/Addis_Ababa"
	}
	loc, err := time.LoadLocation(tz)
	if err != nil {
		loc = time.UTC
	}
	return &Service{
		repo:       repo,
		courses:    courses,
		enroll:     enroll,
		access:     accessSvc,
		lessonRepo: lessonRepo,
		notify:     notify,
		clock:      clock,
		loc:        loc,
	}
}

// MarkLessonComplete marks a lesson as completed for the authenticated student.
// - Validates the actor is authenticated and has access to the course.
// - Upserts the progress row (is_completed=true, idempotent).
// - Recomputes the course completion rollup.
// - Bumps the study streak.
// - Fires CertificateReady hook when the course hits 100%.
func (s *Service) MarkLessonComplete(ctx context.Context, actor Actor, lessonID uuid.UUID) (*entity.LessonProgress, error) {
	if actor.UserID == uuid.Nil {
		return nil, apperror.Unauthorized("authentication required", nil)
	}

	// Resolve the lesson's course so we can gate access and recompute rollup.
	lesson, err := s.resolveLesson(ctx, lessonID)
	if err != nil {
		return nil, err
	}

	if err := s.access.RequireFullAccess(ctx, access.Actor{UserID: actor.UserID, Role: actor.Role}, lesson.CourseID); err != nil {
		return nil, err
	}

	p, err := s.repo.UpsertLessonProgress(ctx, &entity.LessonProgress{
		StudentID:   actor.UserID,
		LessonID:    lessonID,
		CourseID:    lesson.CourseID,
		IsCompleted: true,
	})
	if err != nil {
		return nil, err
	}

	cp, err := s.repo.RecomputeCourseProgress(ctx, actor.UserID, lesson.CourseID)
	if err != nil {
		return nil, err
	}

	if err := s.bumpStreak(ctx, actor.UserID); err != nil {
		return nil, err
	}

	if cp.CompletedAt != nil && cp.TotalLessons > 0 && cp.CompletedLessons >= cp.TotalLessons {
		s.notify.CertificateReady(ctx, actor.UserID, lesson.CourseID)
	}

	return p, nil
}

// UpdateVideoProgress upserts the watch position for a video lesson. The
// position advances monotonically (GREATEST in DB). Auto-completes the lesson
// when position >= 95% of video duration (if duration is known).
func (s *Service) UpdateVideoProgress(ctx context.Context, actor Actor, lessonID uuid.UUID, seconds int) (*entity.LessonProgress, error) {
	if actor.UserID == uuid.Nil {
		return nil, apperror.Unauthorized("authentication required", nil)
	}
	if seconds < 0 {
		seconds = 0
	}

	lesson, err := s.resolveLesson(ctx, lessonID)
	if err != nil {
		return nil, err
	}

	if err := s.access.RequireFullAccess(ctx, access.Actor{UserID: actor.UserID, Role: actor.Role}, lesson.CourseID); err != nil {
		return nil, err
	}

	// Auto-complete if at least 95% of the video has been watched.
	isCompleted := false
	if seconds > 0 && lesson.DurationSeconds > 0 {
		if float64(seconds) >= float64(lesson.DurationSeconds)*videoCompletionThreshold {
			isCompleted = true
		}
	}

	p, err := s.repo.UpsertLessonProgress(ctx, &entity.LessonProgress{
		StudentID:            actor.UserID,
		LessonID:             lessonID,
		CourseID:             lesson.CourseID,
		IsCompleted:          isCompleted,
		VideoProgressSeconds: seconds,
	})
	if err != nil {
		return nil, err
	}

	if err := s.bumpStreak(ctx, actor.UserID); err != nil {
		return nil, err
	}

	if isCompleted {
		if _, err := s.repo.RecomputeCourseProgress(ctx, actor.UserID, lesson.CourseID); err != nil {
			return nil, err
		}
	}

	return p, nil
}

// GetCourseProgress returns the cached completion rollup for a student's course.
// Students can only read their own; teachers/admins can read any student's.
func (s *Service) GetCourseProgress(ctx context.Context, actor Actor, courseID uuid.UUID) (*entity.CourseProgress, []*entity.LessonProgress, error) {
	if actor.UserID == uuid.Nil {
		return nil, nil, apperror.Unauthorized("authentication required", nil)
	}

	cp, err := s.repo.GetCourseProgress(ctx, actor.UserID, courseID)
	if err != nil {
		// No rollup yet — return a zero-value (student has not started).
		cp = &entity.CourseProgress{StudentID: actor.UserID, CourseID: courseID}
	}

	lessons, err := s.repo.ListByCourse(ctx, actor.UserID, courseID)
	if err != nil {
		return nil, nil, err
	}

	return cp, lessons, nil
}

// GetMyDashboard returns the calling student's enrolled courses with their
// completion percentages plus current and longest streaks.
func (s *Service) GetMyDashboard(ctx context.Context, actor Actor) (*Dashboard, error) {
	if actor.UserID == uuid.Nil {
		return nil, apperror.Unauthorized("authentication required", nil)
	}

	enrollments, err := s.enroll.ListByStudent(ctx, actor.UserID, 50)
	if err != nil {
		return nil, err
	}

	items := make([]DashboardItem, 0, len(enrollments))
	for _, e := range enrollments {
		item := DashboardItem{CourseID: e.CourseID}

		// Best-effort: missing progress rows are treated as 0%.
		if cp, err := s.repo.GetCourseProgress(ctx, actor.UserID, e.CourseID); err == nil {
			f, _ := cp.PercentComplete.Float64()
			item.PercentComplete = f
			item.CompletedAt = cp.CompletedAt
			item.CompletedLessons = cp.CompletedLessons
			item.TotalLessons = cp.TotalLessons
		}

		if s.courses != nil {
			if c, err := s.courses.FindByID(ctx, e.CourseID); err == nil {
				item.CourseTitle = c.Title
			}
		}

		items = append(items, item)
	}

	streak, err := s.repo.GetStreak(ctx, actor.UserID)
	if err != nil {
		streak = &entity.StudyStreak{}
	}

	return &Dashboard{
		Items:         items,
		CurrentStreak: streak.CurrentStreak,
		LongestStreak: streak.LongestStreak,
	}, nil
}

// GetMyStreak returns the calling student's study streak.
func (s *Service) GetMyStreak(ctx context.Context, actor Actor) (*entity.StudyStreak, error) {
	if actor.UserID == uuid.Nil {
		return nil, apperror.Unauthorized("authentication required", nil)
	}
	return s.repo.GetStreak(ctx, actor.UserID)
}

// resolveLesson fetches the full Lesson entity. Returns NotFound when the
// lesson does not exist. Callers cannot supply a courseID directly — the
// lesson lookup is the single source of truth for the course association,
// preventing IDOR via a crafted courseID.
func (s *Service) resolveLesson(ctx context.Context, lessonID uuid.UUID) (*entity.Lesson, error) {
	return s.lessonRepo.FindByID(ctx, lessonID)
}

package worker

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/pagination"
)

// ---- fakes ------------------------------------------------------------------

type fakeEngageProgressRepo struct {
	mu            sync.Mutex
	streaks       []*entity.StudyStreak
	warnedDates   map[uuid.UUID]time.Time
}

func newFakeEngageProgressRepo(streaks []*entity.StudyStreak) *fakeEngageProgressRepo {
	return &fakeEngageProgressRepo{
		streaks:     streaks,
		warnedDates: map[uuid.UUID]time.Time{},
	}
}

func (f *fakeEngageProgressRepo) ListStreakAtRisk(_ context.Context, cutoff time.Time, today time.Time, limit int) ([]*entity.StudyStreak, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	cutoffDate := cutoff.Truncate(24 * time.Hour)
	todayDate := today.Truncate(24 * time.Hour)
	var out []*entity.StudyStreak
	for _, st := range f.streaks {
		if st.CurrentStreak == 0 {
			continue
		}
		if st.LastStudyDate == nil {
			continue
		}
		lastStudy := st.LastStudyDate.Truncate(24 * time.Hour)
		if !lastStudy.Before(cutoffDate) {
			continue
		}
		// Idempotency: skip if already warned today.
		if warned, ok := f.warnedDates[st.StudentID]; ok {
			if !warned.Before(todayDate) {
				continue
			}
		}
		cp := *st
		out = append(out, &cp)
		if limit > 0 && len(out) >= limit {
			break
		}
	}
	return out, nil
}

func (f *fakeEngageProgressRepo) SetLastWarnedDate(_ context.Context, studentID uuid.UUID, date time.Time) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.warnedDates[studentID] = date
	return nil
}

// Stub implementations of the full ProgressRepository interface.
func (f *fakeEngageProgressRepo) UpsertLessonProgress(context.Context, *entity.LessonProgress) (*entity.LessonProgress, error) {
	return nil, nil
}
func (f *fakeEngageProgressRepo) GetLessonProgress(context.Context, uuid.UUID, uuid.UUID) (*entity.LessonProgress, error) {
	return nil, apperror.NotFound("not found", nil)
}
func (f *fakeEngageProgressRepo) ListByCourse(context.Context, uuid.UUID, uuid.UUID) ([]*entity.LessonProgress, error) {
	return nil, nil
}
func (f *fakeEngageProgressRepo) RecomputeCourseProgress(context.Context, uuid.UUID, uuid.UUID) (*entity.CourseProgress, error) {
	return nil, nil
}
func (f *fakeEngageProgressRepo) GetCourseProgress(context.Context, uuid.UUID, uuid.UUID) (*entity.CourseProgress, error) {
	return nil, apperror.NotFound("not found", nil)
}
func (f *fakeEngageProgressRepo) GetStreak(context.Context, uuid.UUID) (*entity.StudyStreak, error) {
	return &entity.StudyStreak{}, nil
}
func (f *fakeEngageProgressRepo) UpsertStreak(context.Context, *entity.StudyStreak) error {
	return nil
}

var _ repository.ProgressRepository = (*fakeEngageProgressRepo)(nil)

type fakeStreakNotifier struct {
	mu    sync.Mutex
	calls []uuid.UUID
}

func (n *fakeStreakNotifier) StreakWarning(_ context.Context, studentID uuid.UUID) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.calls = append(n.calls, studentID)
}

func (n *fakeStreakNotifier) count() int {
	n.mu.Lock()
	defer n.mu.Unlock()
	return len(n.calls)
}

var _ StreakNotifier = (*fakeStreakNotifier)(nil)

// ---- helpers ----------------------------------------------------------------

func atRiskStreak(studentID uuid.UUID, lastStudied time.Time, streak int) *entity.StudyStreak {
	t := lastStudied
	return &entity.StudyStreak{
		StudentID:     studentID,
		CurrentStreak: streak,
		LastStudyDate: &t,
	}
}

// makeClock returns a fake clock fixed at the given time.
func makeClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

// ---- tests ------------------------------------------------------------------

func TestEngagementSweeper_WarnsAtRiskStudents(t *testing.T) {
	now := time.Date(2026, 6, 14, 10, 0, 0, 0, time.UTC)
	threeDaysAgo := now.AddDate(0, 0, -3)

	studentID := uuid.New()
	streaks := []*entity.StudyStreak{
		atRiskStreak(studentID, threeDaysAgo, 5),
	}
	progress := newFakeEngageProgressRepo(streaks)
	notify := &fakeStreakNotifier{}

	sw := NewEngagementSweeper(progress, notify, makeClock(now), 24*time.Hour)
	sw.tick(context.Background())

	if notify.count() != 1 {
		t.Errorf("expected 1 streak warning, got %d", notify.count())
	}
}

func TestEngagementSweeper_NoWarningIfStudiedWithin2Days(t *testing.T) {
	now := time.Date(2026, 6, 14, 10, 0, 0, 0, time.UTC)
	yesterday := now.AddDate(0, 0, -1)

	studentID := uuid.New()
	streaks := []*entity.StudyStreak{
		atRiskStreak(studentID, yesterday, 3),
	}
	progress := newFakeEngageProgressRepo(streaks)
	notify := &fakeStreakNotifier{}

	sw := NewEngagementSweeper(progress, notify, makeClock(now), 24*time.Hour)
	sw.tick(context.Background())

	if notify.count() != 0 {
		t.Errorf("expected 0 warnings (studied within 2 days), got %d", notify.count())
	}
}

func TestEngagementSweeper_IdempotentSameDay(t *testing.T) {
	now := time.Date(2026, 6, 14, 10, 0, 0, 0, time.UTC)
	threeDaysAgo := now.AddDate(0, 0, -3)

	studentID := uuid.New()
	streaks := []*entity.StudyStreak{
		atRiskStreak(studentID, threeDaysAgo, 7),
	}
	progress := newFakeEngageProgressRepo(streaks)
	notify := &fakeStreakNotifier{}

	sw := NewEngagementSweeper(progress, notify, makeClock(now), 24*time.Hour)

	// First tick warns.
	sw.tick(context.Background())
	if notify.count() != 1 {
		t.Fatalf("expected 1 warning after first tick, got %d", notify.count())
	}

	// Second tick same day — idempotent, no double-warn.
	sw.tick(context.Background())
	if notify.count() != 1 {
		t.Errorf("expected still 1 warning (idempotent), got %d", notify.count())
	}
}

func TestEngagementSweeper_NoWarningWhenStreakIsZero(t *testing.T) {
	now := time.Date(2026, 6, 14, 10, 0, 0, 0, time.UTC)
	threeDaysAgo := now.AddDate(0, 0, -3)

	studentID := uuid.New()
	streaks := []*entity.StudyStreak{
		atRiskStreak(studentID, threeDaysAgo, 0), // streak already zero
	}
	progress := newFakeEngageProgressRepo(streaks)
	notify := &fakeStreakNotifier{}

	sw := NewEngagementSweeper(progress, notify, makeClock(now), 24*time.Hour)
	sw.tick(context.Background())

	if notify.count() != 0 {
		t.Errorf("expected 0 warnings for zero streak, got %d", notify.count())
	}
}

func TestEngagementSweeper_PagedResults(t *testing.T) {
	now := time.Date(2026, 6, 14, 10, 0, 0, 0, time.UTC)
	threeDaysAgo := now.AddDate(0, 0, -3)

	// Build 150 at-risk students.
	streaks := make([]*entity.StudyStreak, 150)
	for i := range streaks {
		streaks[i] = atRiskStreak(uuid.New(), threeDaysAgo, 1)
	}
	progress := newFakeEngageProgressRepo(streaks)
	notify := &fakeStreakNotifier{}

	sw := NewEngagementSweeper(progress, notify, makeClock(now), 24*time.Hour)
	sw.pageSize = 100 // default
	sw.tick(context.Background())

	// With pageSize=100, one tick should warn at most 100 students.
	if notify.count() > 100 {
		t.Errorf("expected at most 100 warnings per tick (pageSize), got %d", notify.count())
	}
	if notify.count() == 0 {
		t.Error("expected warnings for at-risk students")
	}
}

// Ensure unused pagination import doesn't cause compile errors.
var _ = pagination.Cursor{}

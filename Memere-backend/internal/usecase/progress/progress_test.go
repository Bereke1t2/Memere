package progress

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
)

// newService builds a Service wired to the provided fakes. clock is optional
// (defaults to time.Now via NewService). When clock is non-nil it drives streak
// calculations, enabling fake-clock tests.
func newTestService(
	pRepo *fakeProgressRepo,
	lRepo *fakeLessonRepo,
	cRepo *fakeCourseRepo,
	eRepo *fakeEnrollRepo,
	access CourseAccess,
	notify Notifier,
	clock func() time.Time,
) *Service {
	return NewService(
		pRepo,
		cRepo,
		eRepo,
		access,
		lRepo,
		notify,
		Config{StreakTZ: "UTC"}, // UTC simplifies day-boundary assertions
		clock,
	)
}

// day returns a UTC time representing midnight of the given day in 2024.
func day(d int) time.Time {
	return time.Date(2024, 1, d, 0, 0, 0, 0, time.UTC)
}

// ---- MarkLessonComplete tests -----------------------------------------------

func TestMarkLessonComplete_UpdatesProgressAndStreak(t *testing.T) {
	ctx := context.Background()
	pRepo := newFakeProgressRepo()
	lRepo := newFakeLessonRepo()
	courseID := uuid.New()
	lessonID := uuid.New()
	studentID := uuid.New()

	lRepo.add(&entity.Lesson{ID: lessonID, CourseID: courseID})
	pRepo.setTotalLessons(courseID, 1)

	svc := newTestService(pRepo, lRepo, newFakeCourseRepo(), &fakeEnrollRepo{},
		fakeAccess{}, noopNotifier{}, func() time.Time { return day(1) })

	p, err := svc.MarkLessonComplete(ctx, Actor{UserID: studentID}, lessonID)
	if err != nil {
		t.Fatalf("MarkLessonComplete: %v", err)
	}
	if !p.IsCompleted {
		t.Error("expected IsCompleted=true")
	}

	st, _ := pRepo.GetStreak(ctx, studentID)
	if st.CurrentStreak != 1 {
		t.Errorf("expected streak=1, got %d", st.CurrentStreak)
	}

	cp, _ := pRepo.GetCourseProgress(ctx, studentID, courseID)
	pct, _ := cp.PercentComplete.Float64()
	if pct != 100.0 {
		t.Errorf("expected 100%% course progress, got %.2f", pct)
	}
}

func TestMarkLessonComplete_Idempotent(t *testing.T) {
	ctx := context.Background()
	pRepo := newFakeProgressRepo()
	lRepo := newFakeLessonRepo()
	courseID := uuid.New()
	lessonID := uuid.New()
	studentID := uuid.New()

	lRepo.add(&entity.Lesson{ID: lessonID, CourseID: courseID})
	pRepo.setTotalLessons(courseID, 1)

	svc := newTestService(pRepo, lRepo, newFakeCourseRepo(), &fakeEnrollRepo{},
		fakeAccess{}, noopNotifier{}, func() time.Time { return day(1) })

	if _, err := svc.MarkLessonComplete(ctx, Actor{UserID: studentID}, lessonID); err != nil {
		t.Fatalf("first mark: %v", err)
	}
	p, err := svc.MarkLessonComplete(ctx, Actor{UserID: studentID}, lessonID)
	if err != nil {
		t.Fatalf("second mark: %v", err)
	}
	if !p.IsCompleted {
		t.Error("expected IsCompleted=true after second mark")
	}
}

func TestMarkLessonComplete_FiresCertificateReady_OnCourseCompletion(t *testing.T) {
	ctx := context.Background()
	pRepo := newFakeProgressRepo()
	lRepo := newFakeLessonRepo()
	courseID := uuid.New()
	lessonID := uuid.New()
	studentID := uuid.New()
	notify := &fakeNotifier{}

	lRepo.add(&entity.Lesson{ID: lessonID, CourseID: courseID})
	pRepo.setTotalLessons(courseID, 1)

	svc := newTestService(pRepo, lRepo, newFakeCourseRepo(), &fakeEnrollRepo{},
		fakeAccess{}, notify, func() time.Time { return day(1) })

	if _, err := svc.MarkLessonComplete(ctx, Actor{UserID: studentID}, lessonID); err != nil {
		t.Fatalf("MarkLessonComplete: %v", err)
	}
	if notify.count() != 1 {
		t.Errorf("expected 1 certificate notification, got %d", notify.count())
	}
}

func TestMarkLessonComplete_UnauthenticatedReturnsUnauthorized(t *testing.T) {
	ctx := context.Background()
	svc := newTestService(newFakeProgressRepo(), newFakeLessonRepo(), newFakeCourseRepo(),
		&fakeEnrollRepo{}, fakeAccess{}, noopNotifier{}, nil)

	_, err := svc.MarkLessonComplete(ctx, Actor{UserID: uuid.Nil}, uuid.New())
	if err == nil {
		t.Fatal("expected error for nil UserID")
	}
}

// ---- UpdateVideoProgress tests ----------------------------------------------

func TestUpdateVideoProgress_Monotonic(t *testing.T) {
	ctx := context.Background()
	pRepo := newFakeProgressRepo()
	lRepo := newFakeLessonRepo()
	courseID := uuid.New()
	lessonID := uuid.New()
	studentID := uuid.New()

	lRepo.add(&entity.Lesson{ID: lessonID, CourseID: courseID, DurationSeconds: 200})
	pRepo.setTotalLessons(courseID, 1)

	svc := newTestService(pRepo, lRepo, newFakeCourseRepo(), &fakeEnrollRepo{},
		fakeAccess{}, noopNotifier{}, func() time.Time { return day(1) })

	if _, err := svc.UpdateVideoProgress(ctx, Actor{UserID: studentID}, lessonID, 120); err != nil {
		t.Fatalf("first update: %v", err)
	}
	p, err := svc.UpdateVideoProgress(ctx, Actor{UserID: studentID}, lessonID, 50)
	if err != nil {
		t.Fatalf("second update: %v", err)
	}
	// Monotonic: position must not regress.
	if p.VideoProgressSeconds != 120 {
		t.Errorf("expected VideoProgressSeconds=120 (monotonic), got %d", p.VideoProgressSeconds)
	}
}

func TestUpdateVideoProgress_AutoCompleteAtThreshold(t *testing.T) {
	ctx := context.Background()
	pRepo := newFakeProgressRepo()
	lRepo := newFakeLessonRepo()
	courseID := uuid.New()
	lessonID := uuid.New()
	studentID := uuid.New()

	// Duration 100s; 95s == exactly at 95% threshold.
	lRepo.add(&entity.Lesson{ID: lessonID, CourseID: courseID, DurationSeconds: 100})
	pRepo.setTotalLessons(courseID, 1)

	svc := newTestService(pRepo, lRepo, newFakeCourseRepo(), &fakeEnrollRepo{},
		fakeAccess{}, noopNotifier{}, func() time.Time { return day(1) })

	p, err := svc.UpdateVideoProgress(ctx, Actor{UserID: studentID}, lessonID, 95)
	if err != nil {
		t.Fatalf("UpdateVideoProgress: %v", err)
	}
	if !p.IsCompleted {
		t.Error("expected auto-complete at 95% threshold")
	}
}

func TestUpdateVideoProgress_NoAutoCompleteBelowThreshold(t *testing.T) {
	ctx := context.Background()
	pRepo := newFakeProgressRepo()
	lRepo := newFakeLessonRepo()
	courseID := uuid.New()
	lessonID := uuid.New()
	studentID := uuid.New()

	lRepo.add(&entity.Lesson{ID: lessonID, CourseID: courseID, DurationSeconds: 100})
	pRepo.setTotalLessons(courseID, 1)

	svc := newTestService(pRepo, lRepo, newFakeCourseRepo(), &fakeEnrollRepo{},
		fakeAccess{}, noopNotifier{}, func() time.Time { return day(1) })

	p, err := svc.UpdateVideoProgress(ctx, Actor{UserID: studentID}, lessonID, 50)
	if err != nil {
		t.Fatalf("UpdateVideoProgress: %v", err)
	}
	if p.IsCompleted {
		t.Error("expected no auto-complete at 50% watch")
	}
}

// ---- Streak tests (fake clock) ----------------------------------------------

func TestStreak_ConsecutiveDays(t *testing.T) {
	ctx := context.Background()
	pRepo := newFakeProgressRepo()
	lRepo := newFakeLessonRepo()
	courseID := uuid.New()
	studentID := uuid.New()

	// Two lessons so each day completion makes sense.
	l1, l2 := uuid.New(), uuid.New()
	lRepo.add(&entity.Lesson{ID: l1, CourseID: courseID})
	lRepo.add(&entity.Lesson{ID: l2, CourseID: courseID})
	pRepo.setTotalLessons(courseID, 2)

	now := day(1)
	svc := newTestService(pRepo, lRepo, newFakeCourseRepo(), &fakeEnrollRepo{},
		fakeAccess{}, noopNotifier{}, func() time.Time { return now })

	actor := Actor{UserID: studentID}
	if _, err := svc.MarkLessonComplete(ctx, actor, l1); err != nil {
		t.Fatalf("day 1: %v", err)
	}
	now = day(2)
	if _, err := svc.MarkLessonComplete(ctx, actor, l2); err != nil {
		t.Fatalf("day 2: %v", err)
	}

	st, _ := pRepo.GetStreak(ctx, studentID)
	if st.CurrentStreak != 2 {
		t.Errorf("expected streak=2, got %d", st.CurrentStreak)
	}
	if st.LongestStreak != 2 {
		t.Errorf("expected longest=2, got %d", st.LongestStreak)
	}
}

func TestStreak_SameDayNoChange(t *testing.T) {
	ctx := context.Background()
	pRepo := newFakeProgressRepo()
	lRepo := newFakeLessonRepo()
	courseID := uuid.New()
	studentID := uuid.New()

	l1, l2 := uuid.New(), uuid.New()
	lRepo.add(&entity.Lesson{ID: l1, CourseID: courseID})
	lRepo.add(&entity.Lesson{ID: l2, CourseID: courseID})
	pRepo.setTotalLessons(courseID, 2)

	// Both marks on the same calendar day.
	now := day(1)
	svc := newTestService(pRepo, lRepo, newFakeCourseRepo(), &fakeEnrollRepo{},
		fakeAccess{}, noopNotifier{}, func() time.Time { return now })

	actor := Actor{UserID: studentID}
	if _, err := svc.MarkLessonComplete(ctx, actor, l1); err != nil {
		t.Fatalf("mark 1: %v", err)
	}
	if _, err := svc.MarkLessonComplete(ctx, actor, l2); err != nil {
		t.Fatalf("mark 2: %v", err)
	}

	st, _ := pRepo.GetStreak(ctx, studentID)
	if st.CurrentStreak != 1 {
		t.Errorf("expected streak=1 (same day), got %d", st.CurrentStreak)
	}
}

func TestStreak_GapResetsToOne(t *testing.T) {
	ctx := context.Background()
	pRepo := newFakeProgressRepo()
	lRepo := newFakeLessonRepo()
	courseID := uuid.New()
	studentID := uuid.New()

	l1, l2 := uuid.New(), uuid.New()
	lRepo.add(&entity.Lesson{ID: l1, CourseID: courseID})
	lRepo.add(&entity.Lesson{ID: l2, CourseID: courseID})
	pRepo.setTotalLessons(courseID, 2)

	now := day(1)
	svc := newTestService(pRepo, lRepo, newFakeCourseRepo(), &fakeEnrollRepo{},
		fakeAccess{}, noopNotifier{}, func() time.Time { return now })

	actor := Actor{UserID: studentID}
	if _, err := svc.MarkLessonComplete(ctx, actor, l1); err != nil {
		t.Fatalf("day 1: %v", err)
	}
	// Skip day 2; study on day 3 → gap → reset.
	now = day(3)
	if _, err := svc.MarkLessonComplete(ctx, actor, l2); err != nil {
		t.Fatalf("day 3: %v", err)
	}

	st, _ := pRepo.GetStreak(ctx, studentID)
	if st.CurrentStreak != 1 {
		t.Errorf("expected streak=1 after gap, got %d", st.CurrentStreak)
	}
}

func TestStreak_LongestPreservedAfterReset(t *testing.T) {
	ctx := context.Background()
	pRepo := newFakeProgressRepo()
	lRepo := newFakeLessonRepo()
	courseID := uuid.New()
	studentID := uuid.New()

	lessons := make([]uuid.UUID, 5)
	for i := range lessons {
		lessons[i] = uuid.New()
		lRepo.add(&entity.Lesson{ID: lessons[i], CourseID: courseID})
	}
	pRepo.setTotalLessons(courseID, 5)

	now := day(1)
	svc := newTestService(pRepo, lRepo, newFakeCourseRepo(), &fakeEnrollRepo{},
		fakeAccess{}, noopNotifier{}, func() time.Time { return now })

	actor := Actor{UserID: studentID}
	// Build streak of 3.
	for i := 0; i < 3; i++ {
		now = day(i + 1)
		if _, err := svc.MarkLessonComplete(ctx, actor, lessons[i]); err != nil {
			t.Fatalf("day %d: %v", i+1, err)
		}
	}
	// Gap then rebuild streak of 2.
	now = day(10)
	if _, err := svc.MarkLessonComplete(ctx, actor, lessons[3]); err != nil {
		t.Fatalf("day 10: %v", err)
	}
	now = day(11)
	if _, err := svc.MarkLessonComplete(ctx, actor, lessons[4]); err != nil {
		t.Fatalf("day 11: %v", err)
	}

	st, _ := pRepo.GetStreak(ctx, studentID)
	if st.CurrentStreak != 2 {
		t.Errorf("expected current=2, got %d", st.CurrentStreak)
	}
	if st.LongestStreak != 3 {
		t.Errorf("expected longest=3, got %d", st.LongestStreak)
	}
}

// ---- Dashboard tests --------------------------------------------------------

func TestGetMyDashboard_Composition(t *testing.T) {
	ctx := context.Background()
	pRepo := newFakeProgressRepo()
	lRepo := newFakeLessonRepo()
	cRepo := newFakeCourseRepo()
	eRepo := &fakeEnrollRepo{}

	studentID := uuid.New()
	courseA := uuid.New()
	courseB := uuid.New()
	lessonA := uuid.New()

	// Two enrollments.
	eRepo.add(&entity.Enrollment{ID: uuid.New(), StudentID: studentID, CourseID: courseA})
	eRepo.add(&entity.Enrollment{ID: uuid.New(), StudentID: studentID, CourseID: courseB})

	// Course titles.
	cRepo.add(&entity.Course{ID: courseA, Title: "Math"})
	cRepo.add(&entity.Course{ID: courseB, Title: "Physics"})

	// One lesson completed in course A.
	lRepo.add(&entity.Lesson{ID: lessonA, CourseID: courseA})
	pRepo.setTotalLessons(courseA, 1)
	if _, err := pRepo.UpsertLessonProgress(ctx, &entity.LessonProgress{
		ID: uuid.New(), StudentID: studentID, LessonID: lessonA, CourseID: courseA, IsCompleted: true,
	}); err != nil {
		t.Fatalf("seed: %v", err)
	}
	if _, err := pRepo.RecomputeCourseProgress(ctx, studentID, courseA); err != nil {
		t.Fatalf("recompute: %v", err)
	}

	svc := newTestService(pRepo, lRepo, cRepo, eRepo, fakeAccess{}, noopNotifier{}, nil)

	dash, err := svc.GetMyDashboard(ctx, Actor{UserID: studentID})
	if err != nil {
		t.Fatalf("GetMyDashboard: %v", err)
	}
	if len(dash.Items) != 2 {
		t.Fatalf("expected 2 dashboard items, got %d", len(dash.Items))
	}

	// Find course A item and verify 100%.
	var foundA bool
	for _, item := range dash.Items {
		if item.CourseID == courseA {
			foundA = true
			if item.PercentComplete != 100.0 {
				t.Errorf("course A: expected 100%%, got %.2f", item.PercentComplete)
			}
			if item.CourseTitle != "Math" {
				t.Errorf("course A: expected title 'Math', got %q", item.CourseTitle)
			}
		}
	}
	if !foundA {
		t.Error("course A missing from dashboard")
	}
}

func TestGetMyDashboard_UnauthenticatedReturnsUnauthorized(t *testing.T) {
	ctx := context.Background()
	svc := newTestService(newFakeProgressRepo(), newFakeLessonRepo(), newFakeCourseRepo(),
		&fakeEnrollRepo{}, fakeAccess{}, noopNotifier{}, nil)

	_, err := svc.GetMyDashboard(ctx, Actor{UserID: uuid.Nil})
	if err == nil {
		t.Fatal("expected error for nil UserID")
	}
}

// ---- GetCourseProgress tests ------------------------------------------------

func TestGetCourseProgress_ReturnsZeroWhenNotStarted(t *testing.T) {
	ctx := context.Background()
	pRepo := newFakeProgressRepo()
	studentID := uuid.New()
	courseID := uuid.New()

	svc := newTestService(pRepo, newFakeLessonRepo(), newFakeCourseRepo(), &fakeEnrollRepo{},
		fakeAccess{}, noopNotifier{}, nil)

	cp, lessons, err := svc.GetCourseProgress(ctx, Actor{UserID: studentID}, courseID)
	if err != nil {
		t.Fatalf("GetCourseProgress: %v", err)
	}
	pct, _ := cp.PercentComplete.Float64()
	if pct != 0 {
		t.Errorf("expected 0%% for unstarted course, got %.2f", pct)
	}
	if len(lessons) != 0 {
		t.Errorf("expected 0 lesson rows, got %d", len(lessons))
	}
}

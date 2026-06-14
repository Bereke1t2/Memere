package progress

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/access"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/pagination"
)

// ---- fake progress repo -------------------------------------------------------

type fakeProgressRepo struct {
	mu           sync.Mutex
	progress     map[[2]uuid.UUID]*entity.LessonProgress // (student,lesson)
	streaks      map[uuid.UUID]*entity.StudyStreak
	courseProgress map[[2]uuid.UUID]*entity.CourseProgress // (student,course)
	totalLessons map[uuid.UUID]int                        // courseID → total published
}

func newFakeProgressRepo() *fakeProgressRepo {
	return &fakeProgressRepo{
		progress:       map[[2]uuid.UUID]*entity.LessonProgress{},
		streaks:        map[uuid.UUID]*entity.StudyStreak{},
		courseProgress: map[[2]uuid.UUID]*entity.CourseProgress{},
		totalLessons:   map[uuid.UUID]int{},
	}
}

func (f *fakeProgressRepo) setTotalLessons(courseID uuid.UUID, n int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.totalLessons[courseID] = n
}

func (f *fakeProgressRepo) UpsertLessonProgress(_ context.Context, p *entity.LessonProgress) (*entity.LessonProgress, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	key := [2]uuid.UUID{p.StudentID, p.LessonID}
	if existing, ok := f.progress[key]; ok {
		// Upsert: preserve completed (once set, never cleared), advance video pos.
		if p.IsCompleted {
			existing.IsCompleted = true
		}
		if p.VideoProgressSeconds > existing.VideoProgressSeconds {
			existing.VideoProgressSeconds = p.VideoProgressSeconds
		}
		cp := *existing
		return &cp, nil
	}
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	cp := *p
	f.progress[key] = &cp
	return &cp, nil
}

func (f *fakeProgressRepo) GetLessonProgress(_ context.Context, studentID, lessonID uuid.UUID) (*entity.LessonProgress, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if p, ok := f.progress[[2]uuid.UUID{studentID, lessonID}]; ok {
		cp := *p
		return &cp, nil
	}
	return nil, apperror.NotFound("lesson progress not found", nil)
}

func (f *fakeProgressRepo) ListByCourse(_ context.Context, studentID, courseID uuid.UUID) ([]*entity.LessonProgress, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	var out []*entity.LessonProgress
	for _, p := range f.progress {
		if p.StudentID == studentID && p.CourseID == courseID {
			cp := *p
			out = append(out, &cp)
		}
	}
	return out, nil
}

func (f *fakeProgressRepo) RecomputeCourseProgress(_ context.Context, studentID, courseID uuid.UUID) (*entity.CourseProgress, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	total := f.totalLessons[courseID]
	completed := 0
	for _, p := range f.progress {
		if p.StudentID == studentID && p.CourseID == courseID && p.IsCompleted {
			completed++
		}
	}
	var pct decimal.Decimal
	if total > 0 {
		pct = decimal.NewFromInt(int64(completed)).Div(decimal.NewFromInt(int64(total))).Mul(decimal.NewFromInt(100)).Round(2)
	}
	cp := &entity.CourseProgress{
		StudentID:        studentID,
		CourseID:         courseID,
		CompletedLessons: completed,
		TotalLessons:     total,
		PercentComplete:  pct,
	}
	if total > 0 && completed >= total {
		existing := f.courseProgress[[2]uuid.UUID{studentID, courseID}]
		if existing != nil && existing.CompletedAt != nil {
			// Preserve the original completion timestamp (set once, never cleared).
			cp.CompletedAt = existing.CompletedAt
		} else {
			now := time.Now().UTC()
			cp.CompletedAt = &now
		}
	}
	f.courseProgress[[2]uuid.UUID{studentID, courseID}] = cp
	return cp, nil
}

func (f *fakeProgressRepo) GetCourseProgress(_ context.Context, studentID, courseID uuid.UUID) (*entity.CourseProgress, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if cp, ok := f.courseProgress[[2]uuid.UUID{studentID, courseID}]; ok {
		return cp, nil
	}
	return nil, apperror.NotFound("course progress not found", nil)
}

func (f *fakeProgressRepo) GetStreak(_ context.Context, studentID uuid.UUID) (*entity.StudyStreak, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if st, ok := f.streaks[studentID]; ok {
		cp := *st
		return &cp, nil
	}
	return &entity.StudyStreak{StudentID: studentID}, nil
}

func (f *fakeProgressRepo) UpsertStreak(_ context.Context, st *entity.StudyStreak) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	cp := *st
	f.streaks[st.StudentID] = &cp
	return nil
}

func (f *fakeProgressRepo) ListStreakAtRisk(_ context.Context, _ time.Time, _ time.Time, _ int) ([]*entity.StudyStreak, error) {
	return nil, nil
}

func (f *fakeProgressRepo) SetLastWarnedDate(_ context.Context, _ uuid.UUID, _ time.Time) error {
	return nil
}

// ---- fake lesson repo -------------------------------------------------------

type fakeLessonRepo struct {
	mu    sync.Mutex
	byID  map[uuid.UUID]*entity.Lesson
}

func newFakeLessonRepo() *fakeLessonRepo {
	return &fakeLessonRepo{byID: map[uuid.UUID]*entity.Lesson{}}
}

func (f *fakeLessonRepo) add(l *entity.Lesson) {
	f.mu.Lock()
	defer f.mu.Unlock()
	cp := *l
	f.byID[l.ID] = &cp
}

func (f *fakeLessonRepo) FindByID(_ context.Context, id uuid.UUID) (*entity.Lesson, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if l, ok := f.byID[id]; ok {
		cp := *l
		return &cp, nil
	}
	return nil, apperror.NotFound("lesson not found", nil)
}

func (f *fakeLessonRepo) Create(context.Context, *entity.Lesson) error             { return nil }
func (f *fakeLessonRepo) ListBySection(context.Context, uuid.UUID) ([]*entity.Lesson, error) {
	return nil, nil
}
func (f *fakeLessonRepo) ListByCourse(context.Context, uuid.UUID) ([]*entity.Lesson, error) {
	return nil, nil
}
func (f *fakeLessonRepo) Update(context.Context, *entity.Lesson) error              { return nil }
func (f *fakeLessonRepo) SoftDelete(context.Context, uuid.UUID) error               { return nil }
func (f *fakeLessonRepo) SoftDeleteBySection(context.Context, uuid.UUID) error      { return nil }
func (f *fakeLessonRepo) MaxOrderIndex(context.Context, uuid.UUID) (int, error)     { return -1, nil }

// ---- fake course repo -------------------------------------------------------

type fakeCourseRepo struct {
	mu   sync.Mutex
	byID map[uuid.UUID]*entity.Course
}

func newFakeCourseRepo() *fakeCourseRepo {
	return &fakeCourseRepo{byID: map[uuid.UUID]*entity.Course{}}
}

func (f *fakeCourseRepo) add(c *entity.Course) {
	f.mu.Lock()
	defer f.mu.Unlock()
	cp := *c
	f.byID[c.ID] = &cp
}

func (f *fakeCourseRepo) FindByID(_ context.Context, id uuid.UUID) (*entity.Course, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if c, ok := f.byID[id]; ok {
		cp := *c
		return &cp, nil
	}
	return nil, apperror.NotFound("course not found", nil)
}

func (f *fakeCourseRepo) Create(context.Context, *entity.Course) error { return nil }
func (f *fakeCourseRepo) FindBySlug(context.Context, string) (*entity.Course, error) {
	return nil, apperror.NotFound("course not found", nil)
}
func (f *fakeCourseRepo) List(context.Context, repository.CourseFilter, *pagination.Cursor, int) ([]*entity.Course, *pagination.Cursor, error) {
	return nil, nil, nil
}
func (f *fakeCourseRepo) Update(context.Context, *entity.Course) error       { return nil }
func (f *fakeCourseRepo) SoftDelete(context.Context, uuid.UUID) error        { return nil }
func (f *fakeCourseRepo) RecomputeCounters(context.Context, uuid.UUID) error { return nil }
func (f *fakeCourseRepo) GetCourseWithSectionsAndLessons(context.Context, uuid.UUID) (*repository.CourseWithContent, error) {
	return nil, apperror.NotFound("course not found", nil)
}

// ---- fake enrollment repo ---------------------------------------------------

type fakeEnrollRepo struct {
	mu          sync.Mutex
	enrollments []*entity.Enrollment
}

func (f *fakeEnrollRepo) add(e *entity.Enrollment) {
	f.mu.Lock()
	defer f.mu.Unlock()
	cp := *e
	f.enrollments = append(f.enrollments, &cp)
}

func (f *fakeEnrollRepo) Create(_ context.Context, e *entity.Enrollment) error {
	f.add(e)
	return nil
}

func (f *fakeEnrollRepo) Exists(_ context.Context, studentID, courseID uuid.UUID) (bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, e := range f.enrollments {
		if e.StudentID == studentID && e.CourseID == courseID {
			return true, nil
		}
	}
	return false, nil
}

func (f *fakeEnrollRepo) GetActiveForStudent(_ context.Context, studentID, courseID uuid.UUID) (*entity.Enrollment, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, e := range f.enrollments {
		if e.StudentID == studentID && e.CourseID == courseID {
			cp := *e
			return &cp, nil
		}
	}
	return nil, apperror.NotFound("enrollment not found", nil)
}

func (f *fakeEnrollRepo) ListByStudent(_ context.Context, studentID uuid.UUID, limit int) ([]*entity.Enrollment, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	var out []*entity.Enrollment
	for _, e := range f.enrollments {
		if e.StudentID == studentID {
			cp := *e
			out = append(out, &cp)
			if limit > 0 && len(out) >= limit {
				break
			}
		}
	}
	return out, nil
}

// ---- fake course access gate ------------------------------------------------

type fakeAccess struct{ err error }

func (f fakeAccess) RequireFullAccess(_ context.Context, _ access.Actor, _ uuid.UUID) error {
	return f.err
}

// ---- fake notifier ----------------------------------------------------------

type fakeNotifier struct {
	mu    sync.Mutex
	calls []notifyCall
}

type notifyCall struct {
	studentID uuid.UUID
	courseID  uuid.UUID
}

func (n *fakeNotifier) CertificateReady(_ context.Context, studentID, courseID uuid.UUID) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.calls = append(n.calls, notifyCall{studentID, courseID})
}

func (n *fakeNotifier) count() int {
	n.mu.Lock()
	defer n.mu.Unlock()
	return len(n.calls)
}

package access

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/pagination"
)

// ---- fakes -------------------------------------------------------------------

type fakeCourseRepo struct{ byID map[uuid.UUID]*entity.Course }

func (f *fakeCourseRepo) FindByID(_ context.Context, id uuid.UUID) (*entity.Course, error) {
	if c, ok := f.byID[id]; ok {
		cp := *c
		return &cp, nil
	}
	return nil, apperror.NotFound("course not found", nil)
}

// The access service only calls FindByID; the rest satisfy the interface.
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

type fakeEnrollRepo struct {
	active map[[2]uuid.UUID]*entity.Enrollment // (student,course) -> active enrollment
}

func (f *fakeEnrollRepo) Create(context.Context, *entity.Enrollment) error { return nil }
func (f *fakeEnrollRepo) Exists(_ context.Context, s, c uuid.UUID) (bool, error) {
	_, ok := f.active[[2]uuid.UUID{s, c}]
	return ok, nil
}
func (f *fakeEnrollRepo) GetActiveForStudent(_ context.Context, s, c uuid.UUID) (*entity.Enrollment, error) {
	if e, ok := f.active[[2]uuid.UUID{s, c}]; ok {
		return e, nil
	}
	return nil, apperror.NotFound("enrollment not found", nil)
}
func (f *fakeEnrollRepo) ListByStudent(context.Context, uuid.UUID, int) ([]*entity.Enrollment, error) {
	return nil, nil
}

type fakeSubRepo struct{ active map[uuid.UUID]*entity.Subscription }

func (f *fakeSubRepo) Create(context.Context, *entity.Subscription) error { return nil }
func (f *fakeSubRepo) GetByID(context.Context, uuid.UUID) (*entity.Subscription, error) {
	return nil, apperror.NotFound("subscription not found", nil)
}
func (f *fakeSubRepo) GetActiveForStudent(_ context.Context, s uuid.UUID) (*entity.Subscription, error) {
	if sub, ok := f.active[s]; ok {
		return sub, nil
	}
	return nil, apperror.NotFound("subscription not found", nil)
}
func (f *fakeSubRepo) UpdateStatus(context.Context, uuid.UUID, entity.SubscriptionStatus) error {
	return nil
}
func (f *fakeSubRepo) ListExpiring(context.Context, int) ([]*entity.Subscription, error) {
	return nil, nil
}

// ---- helpers -----------------------------------------------------------------

func newServiceWith(course *entity.Course, enroll *fakeEnrollRepo, sub *fakeSubRepo, now time.Time) *Service {
	cr := &fakeCourseRepo{byID: map[uuid.UUID]*entity.Course{course.ID: course}}
	return NewService(enroll, sub, cr, func() time.Time { return now })
}

// ---- tests -------------------------------------------------------------------

func TestCanAccessCourse_Matrix(t *testing.T) {
	now := time.Date(2026, 6, 12, 0, 0, 0, 0, time.UTC)
	teacher := uuid.New()
	student := uuid.New()
	courseID := uuid.New()

	paid := &entity.Course{ID: courseID, TeacherID: teacher, IsFree: false}

	cases := []struct {
		name   string
		actor  Actor
		course *entity.Course
		enroll *fakeEnrollRepo
		sub    *fakeSubRepo
		want   AccessLevel
	}{
		{
			name:   "owner teacher -> full",
			actor:  Actor{UserID: teacher, Role: entity.RoleTeacher},
			course: paid,
			enroll: &fakeEnrollRepo{active: map[[2]uuid.UUID]*entity.Enrollment{}},
			sub:    &fakeSubRepo{active: map[uuid.UUID]*entity.Subscription{}},
			want:   FullAccess,
		},
		{
			name:   "admin -> full",
			actor:  Actor{UserID: uuid.New(), Role: entity.RoleAdmin},
			course: paid,
			enroll: &fakeEnrollRepo{active: map[[2]uuid.UUID]*entity.Enrollment{}},
			sub:    &fakeSubRepo{active: map[uuid.UUID]*entity.Subscription{}},
			want:   FullAccess,
		},
		{
			name:   "free course -> full for any student",
			actor:  Actor{UserID: student, Role: entity.RoleStudent},
			course: &entity.Course{ID: courseID, TeacherID: teacher, IsFree: true},
			enroll: &fakeEnrollRepo{active: map[[2]uuid.UUID]*entity.Enrollment{}},
			sub:    &fakeSubRepo{active: map[uuid.UUID]*entity.Subscription{}},
			want:   FullAccess,
		},
		{
			name:  "enrolled student -> full",
			actor: Actor{UserID: student, Role: entity.RoleStudent},
			course: paid,
			enroll: &fakeEnrollRepo{active: map[[2]uuid.UUID]*entity.Enrollment{
				{student, courseID}: {StudentID: student, CourseID: courseID, ExpiresAt: nil},
			}},
			sub:  &fakeSubRepo{active: map[uuid.UUID]*entity.Subscription{}},
			want: FullAccess,
		},
		{
			name:   "active subscription -> full",
			actor:  Actor{UserID: student, Role: entity.RoleStudent},
			course: paid,
			enroll: &fakeEnrollRepo{active: map[[2]uuid.UUID]*entity.Enrollment{}},
			sub: &fakeSubRepo{active: map[uuid.UUID]*entity.Subscription{
				student: {StudentID: student, Status: entity.SubActive, CurrentPeriodEnd: now.Add(24 * time.Hour)},
			}},
			want: FullAccess,
		},
		{
			name:   "no entitlement -> preview",
			actor:  Actor{UserID: student, Role: entity.RoleStudent},
			course: paid,
			enroll: &fakeEnrollRepo{active: map[[2]uuid.UUID]*entity.Enrollment{}},
			sub:    &fakeSubRepo{active: map[uuid.UUID]*entity.Subscription{}},
			want:   PreviewAccess,
		},
		{
			name:   "anonymous on paid -> preview",
			actor:  Actor{},
			course: paid,
			enroll: &fakeEnrollRepo{active: map[[2]uuid.UUID]*entity.Enrollment{}},
			sub:    &fakeSubRepo{active: map[uuid.UUID]*entity.Subscription{}},
			want:   PreviewAccess,
		},
		{
			name:   "expired subscription -> preview",
			actor:  Actor{UserID: student, Role: entity.RoleStudent},
			course: paid,
			enroll: &fakeEnrollRepo{active: map[[2]uuid.UUID]*entity.Enrollment{}},
			sub: &fakeSubRepo{active: map[uuid.UUID]*entity.Subscription{
				student: {StudentID: student, Status: entity.SubActive, CurrentPeriodEnd: now.Add(-time.Hour)},
			}},
			want: PreviewAccess,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := newServiceWith(tc.course, tc.enroll, tc.sub, now)
			got, err := s.CanAccessCourse(context.Background(), tc.actor, tc.course.ID)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("CanAccessCourse = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestCanAccessLesson_PreviewGate(t *testing.T) {
	now := time.Date(2026, 6, 12, 0, 0, 0, 0, time.UTC)
	teacher := uuid.New()
	student := uuid.New()
	courseID := uuid.New()
	paid := &entity.Course{ID: courseID, TeacherID: teacher, IsFree: false}

	s := newServiceWith(paid,
		&fakeEnrollRepo{active: map[[2]uuid.UUID]*entity.Enrollment{}},
		&fakeSubRepo{active: map[uuid.UUID]*entity.Subscription{}},
		now)

	actor := Actor{UserID: student, Role: entity.RoleStudent}

	preview := &entity.Lesson{CourseID: courseID, IsFreePreview: true}
	ok, err := s.CanAccessLesson(context.Background(), actor, preview)
	if err != nil || !ok {
		t.Fatalf("preview lesson should be visible at preview access: ok=%v err=%v", ok, err)
	}

	paidLesson := &entity.Lesson{CourseID: courseID, IsFreePreview: false}
	ok, err = s.CanAccessLesson(context.Background(), actor, paidLesson)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatalf("paid lesson must NOT be visible without full access")
	}
}

func TestRequireFullAccess(t *testing.T) {
	now := time.Date(2026, 6, 12, 0, 0, 0, 0, time.UTC)
	teacher := uuid.New()
	student := uuid.New()
	courseID := uuid.New()
	paid := &entity.Course{ID: courseID, TeacherID: teacher, IsFree: false}

	s := newServiceWith(paid,
		&fakeEnrollRepo{active: map[[2]uuid.UUID]*entity.Enrollment{}},
		&fakeSubRepo{active: map[uuid.UUID]*entity.Subscription{}},
		now)

	err := s.RequireFullAccess(context.Background(), Actor{UserID: student, Role: entity.RoleStudent}, courseID)
	if !apperror.IsCode(err, "NOT_ENROLLED") {
		t.Fatalf("non-enrolled student should get NOT_ENROLLED, got %v", err)
	}
}

package enrollment

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

type fakeEnrollRepo struct {
	rows    map[[2]uuid.UUID]*entity.Enrollment
	creates int
}

func newFakeEnrollRepo() *fakeEnrollRepo {
	return &fakeEnrollRepo{rows: map[[2]uuid.UUID]*entity.Enrollment{}}
}

func (f *fakeEnrollRepo) Create(_ context.Context, e *entity.Enrollment) error {
	f.creates++
	f.rows[[2]uuid.UUID{e.StudentID, e.CourseID}] = e
	return nil
}
func (f *fakeEnrollRepo) Exists(_ context.Context, s, c uuid.UUID) (bool, error) {
	_, ok := f.rows[[2]uuid.UUID{s, c}]
	return ok, nil
}
func (f *fakeEnrollRepo) GetActiveForStudent(_ context.Context, s, c uuid.UUID) (*entity.Enrollment, error) {
	if e, ok := f.rows[[2]uuid.UUID{s, c}]; ok {
		return e, nil
	}
	return nil, apperror.NotFound("enrollment not found", nil)
}
func (f *fakeEnrollRepo) ListByStudent(_ context.Context, s uuid.UUID, _ int) ([]*entity.Enrollment, error) {
	var out []*entity.Enrollment
	for k, e := range f.rows {
		if k[0] == s {
			out = append(out, e)
		}
	}
	return out, nil
}

// ---- tests -------------------------------------------------------------------

func newService(enroll *fakeEnrollRepo, courseFree bool) *Service {
	now := time.Date(2026, 6, 12, 0, 0, 0, 0, time.UTC)
	return &Service{
		enroll: enroll,
		course: courseRepoStub{free: courseFree},
		clock:  func() time.Time { return now },
	}
}

// courseRepoStub satisfies repository.CourseRepository; only FindByID is
// exercised by the enrollment engine.
type courseRepoStub struct{ free bool }

func (s courseRepoStub) FindByID(_ context.Context, id uuid.UUID) (*entity.Course, error) {
	return &entity.Course{ID: id, IsFree: s.free}, nil
}
func (s courseRepoStub) Create(context.Context, *entity.Course) error { return nil }
func (s courseRepoStub) FindBySlug(context.Context, string) (*entity.Course, error) {
	return nil, apperror.NotFound("course not found", nil)
}
func (s courseRepoStub) List(context.Context, repository.CourseFilter, *pagination.Cursor, int) ([]*entity.Course, *pagination.Cursor, error) {
	return nil, nil, nil
}
func (s courseRepoStub) Update(context.Context, *entity.Course) error       { return nil }
func (s courseRepoStub) SoftDelete(context.Context, uuid.UUID) error        { return nil }
func (s courseRepoStub) RecomputeCounters(context.Context, uuid.UUID) error { return nil }
func (s courseRepoStub) GetCourseWithSectionsAndLessons(context.Context, uuid.UUID) (*repository.CourseWithContent, error) {
	return nil, apperror.NotFound("course not found", nil)
}

func TestGrantEnrollment_Idempotent(t *testing.T) {
	enroll := newFakeEnrollRepo()
	svc := &Service{enroll: enroll, clock: func() time.Time { return time.Unix(0, 0) }}
	student, course := uuid.New(), uuid.New()

	if err := svc.GrantEnrollment(context.Background(), student, course, entity.SourcePurchase, nil); err != nil {
		t.Fatalf("first grant: %v", err)
	}
	if err := svc.GrantEnrollment(context.Background(), student, course, entity.SourcePurchase, nil); err != nil {
		t.Fatalf("second grant: %v", err)
	}
	if enroll.creates != 1 {
		t.Fatalf("expected exactly 1 Create, got %d (grant not idempotent)", enroll.creates)
	}
}

func TestGrantEnrollment_RejectsBadInput(t *testing.T) {
	enroll := newFakeEnrollRepo()
	svc := &Service{enroll: enroll, clock: func() time.Time { return time.Unix(0, 0) }}

	if err := svc.GrantEnrollment(context.Background(), uuid.Nil, uuid.New(), entity.SourcePurchase, nil); !apperror.IsCode(err, "INVALID_ARGUMENT") {
		t.Fatalf("nil student should be rejected, got %v", err)
	}
	if err := svc.GrantEnrollment(context.Background(), uuid.New(), uuid.New(), entity.EnrollmentSource("bogus"), nil); !apperror.IsCode(err, "INVALID_ARGUMENT") {
		t.Fatalf("bad source should be rejected, got %v", err)
	}
}

func TestEnrollFree_RefusesPaidCourse(t *testing.T) {
	enroll := newFakeEnrollRepo()
	svc := newService(enroll, false) // course is paid
	actor := Actor{UserID: uuid.New(), Role: entity.RoleStudent}

	err := svc.EnrollFree(context.Background(), actor, uuid.New())
	if !apperror.IsCode(err, "COURSE_NOT_FREE") {
		t.Fatalf("EnrollFree on a paid course must fail with COURSE_NOT_FREE, got %v", err)
	}
	if enroll.creates != 0 {
		t.Fatalf("no enrollment should have been created for a paid course")
	}
}

func TestEnrollFree_GrantsFreeCourse(t *testing.T) {
	enroll := newFakeEnrollRepo()
	svc := newService(enroll, true) // course is free
	actor := Actor{UserID: uuid.New(), Role: entity.RoleStudent}

	if err := svc.EnrollFree(context.Background(), actor, uuid.New()); err != nil {
		t.Fatalf("EnrollFree on a free course: %v", err)
	}
	if enroll.creates != 1 {
		t.Fatalf("expected 1 free enrollment, got %d", enroll.creates)
	}
}

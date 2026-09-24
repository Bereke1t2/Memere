package courseaccess

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

type fakeRequestRepo struct {
	mu       sync.Mutex
	byID     map[uuid.UUID]*entity.CourseAccessRequest
	byPair   map[string]*entity.CourseAccessRequest
}

func newFakeRequestRepo() *fakeRequestRepo {
	return &fakeRequestRepo{
		byID:   make(map[uuid.UUID]*entity.CourseAccessRequest),
		byPair: make(map[string]*entity.CourseAccessRequest),
	}
}

func pairKey(s, c uuid.UUID) string {
	return s.String() + ":" + c.String()
}

func (r *fakeRequestRepo) Create(_ context.Context, req *entity.CourseAccessRequest) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if req.ID == uuid.Nil {
		req.ID = uuid.New()
	}
	cp := *req
	r.byID[req.ID] = &cp
	r.byPair[pairKey(req.StudentID, req.CourseID)] = &cp
	return nil
}

func (r *fakeRequestRepo) GetByID(_ context.Context, id uuid.UUID) (*entity.CourseAccessRequest, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if req, ok := r.byID[id]; ok {
		cp := *req
		return &cp, nil
	}
	return nil, apperror.NotFound("not found", nil)
}

func (r *fakeRequestRepo) GetByStudentAndCourse(_ context.Context, studentID, courseID uuid.UUID) (*entity.CourseAccessRequest, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if req, ok := r.byPair[pairKey(studentID, courseID)]; ok {
		cp := *req
		return &cp, nil
	}
	return nil, apperror.NotFound("not found", nil)
}

func (r *fakeRequestRepo) Update(_ context.Context, req *entity.CourseAccessRequest) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *req
	r.byID[req.ID] = &cp
	r.byPair[pairKey(req.StudentID, req.CourseID)] = &cp
	return nil
}

func (r *fakeRequestRepo) ListByTeacher(_ context.Context, teacherID uuid.UUID, status *string, _ *pagination.Cursor, _ int) ([]*entity.CourseAccessRequest, *pagination.Cursor, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var out []*entity.CourseAccessRequest
	for _, req := range r.byID {
		if status != nil && string(req.Status) != *status {
			continue
		}
		cp := *req
		out = append(out, &cp)
	}
	return out, nil, nil
}

func (r *fakeRequestRepo) ListAll(_ context.Context, status *string, _ *pagination.Cursor, _ int) ([]*entity.CourseAccessRequest, *pagination.Cursor, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var out []*entity.CourseAccessRequest
	for _, req := range r.byID {
		if status != nil && string(req.Status) != *status {
			continue
		}
		cp := *req
		out = append(out, &cp)
	}
	return out, nil, nil
}

func (r *fakeRequestRepo) ListByStudent(_ context.Context, studentID uuid.UUID) ([]*entity.CourseAccessRequest, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var out []*entity.CourseAccessRequest
	for _, req := range r.byID {
		if req.StudentID == studentID {
			cp := *req
			out = append(out, &cp)
		}
	}
	return out, nil
}

type fakeEnrollRepo struct {
	mu     sync.Mutex
	byPair map[string]*entity.Enrollment
}

func newFakeEnrollRepo() *fakeEnrollRepo {
	return &fakeEnrollRepo{byPair: make(map[string]*entity.Enrollment)}
}

func (r *fakeEnrollRepo) Create(_ context.Context, e *entity.Enrollment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *e
	r.byPair[pairKey(e.StudentID, e.CourseID)] = &cp
	return nil
}

func (r *fakeEnrollRepo) Exists(_ context.Context, studentID, courseID uuid.UUID) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.byPair[pairKey(studentID, courseID)]
	return ok, nil
}

func (r *fakeEnrollRepo) GetActiveForStudent(_ context.Context, studentID, courseID uuid.UUID) (*entity.Enrollment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if e, ok := r.byPair[pairKey(studentID, courseID)]; ok {
		cp := *e
		return &cp, nil
	}
	return nil, apperror.NotFound("not found", nil)
}

func (r *fakeEnrollRepo) ListByStudent(_ context.Context, studentID uuid.UUID, _ int) ([]*entity.Enrollment, error) {
	return r.ListAllByStudent(context.Background(), studentID)
}

func (r *fakeEnrollRepo) ListAllByStudent(_ context.Context, studentID uuid.UUID) ([]*entity.Enrollment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var out []*entity.Enrollment
	for _, e := range r.byPair {
		if e.StudentID == studentID {
			cp := *e
			out = append(out, &cp)
		}
	}
	return out, nil
}

func (r *fakeEnrollRepo) Delete(_ context.Context, studentID, courseID uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.byPair, pairKey(studentID, courseID))
	return nil
}

type fakeCourseRepo struct {
	courses map[uuid.UUID]*entity.Course
}

func (r *fakeCourseRepo) Create(_ context.Context, _ *entity.Course) error { return nil }
func (r *fakeCourseRepo) FindByID(_ context.Context, id uuid.UUID) (*entity.Course, error) {
	if c, ok := r.courses[id]; ok {
		cp := *c
		return &cp, nil
	}
	return nil, apperror.NotFound("course not found", nil)
}
func (r *fakeCourseRepo) FindBySlug(_ context.Context, _ string) (*entity.Course, error) { return nil, nil }
func (r *fakeCourseRepo) Update(_ context.Context, _ *entity.Course) error { return nil }
func (r *fakeCourseRepo) SoftDelete(_ context.Context, _ uuid.UUID) error { return nil }
func (r *fakeCourseRepo) List(_ context.Context, _ repository.CourseFilter, _ *pagination.Cursor, _ int) ([]*entity.Course, *pagination.Cursor, error) {
	var out []*entity.Course
	for _, c := range r.courses {
		cp := *c
		out = append(out, &cp)
	}
	return out, nil, nil
}
func (r *fakeCourseRepo) RecomputeCounters(_ context.Context, _ uuid.UUID) error { return nil }
func (r *fakeCourseRepo) GetCourseWithSectionsAndLessons(_ context.Context, _ uuid.UUID) (*repository.CourseWithContent, error) {
	return nil, nil
}

type fakeUserRepo struct {
	users map[uuid.UUID]*entity.User
}

func (r *fakeUserRepo) Create(_ context.Context, _ *entity.User) error { return nil }
func (r *fakeUserRepo) FindByID(_ context.Context, id uuid.UUID) (*entity.User, error) {
	if u, ok := r.users[id]; ok {
		cp := *u
		return &cp, nil
	}
	return nil, apperror.NotFound("user not found", nil)
}
func (r *fakeUserRepo) FindByEmail(_ context.Context, _ string) (*entity.User, error) { return nil, nil }
func (r *fakeUserRepo) Update(_ context.Context, u *entity.User) error {
	cp := *u
	r.users[u.ID] = &cp
	return nil
}
func (r *fakeUserRepo) SoftDelete(_ context.Context, _ uuid.UUID) error { return nil }
func (r *fakeUserRepo) SetLastLogin(_ context.Context, _ uuid.UUID, _ time.Time) error { return nil }
func (r *fakeUserRepo) List(_ context.Context, _ repository.AdminUserFilter, _ *pagination.Cursor, _ int) ([]*entity.User, *pagination.Cursor, error) { return nil, nil, nil }
func (r *fakeUserRepo) CountByRole(_ context.Context, _ entity.Role) (int, error) { return 0, nil }

func setupTest(t *testing.T) (*Service, *fakeRequestRepo, *fakeEnrollRepo, uuid.UUID, uuid.UUID, uuid.UUID, uuid.UUID) {
	reqRepo := newFakeRequestRepo()
	enrollRepo := newFakeEnrollRepo()

	teacherID := uuid.New()
	pendingStudentID := uuid.New()
	approvedStudentID := uuid.New()
	paidCourseID := uuid.New()

	courses := map[uuid.UUID]*entity.Course{
		paidCourseID: {
			ID:          paidCourseID,
			TeacherID:   teacherID,
			Title:       "Advanced Mathematics",
			IsFree:      false,
			IsPublished: true,
			Price:       500,
		},
	}
	users := map[uuid.UUID]*entity.User{
		teacherID: {
			ID:             teacherID,
			Role:           entity.RoleTeacher,
			ApprovalStatus: entity.ApprovalStatusApproved,
		},
		pendingStudentID: {
			ID:             pendingStudentID,
			Role:           entity.RoleStudent,
			ApprovalStatus: entity.ApprovalStatusPending,
		},
		approvedStudentID: {
			ID:             approvedStudentID,
			Role:           entity.RoleStudent,
			ApprovalStatus: entity.ApprovalStatusApproved,
		},
	}

	courseRepo := &fakeCourseRepo{courses: courses}
	userRepo := &fakeUserRepo{users: users}

	svc := NewService(reqRepo, enrollRepo, courseRepo, userRepo, nil, time.Now)
	return svc, reqRepo, enrollRepo, teacherID, pendingStudentID, approvedStudentID, paidCourseID
}

func TestRequestAccess_PendingStudent_Forbidden(t *testing.T) {
	svc, _, _, _, pendingID, _, courseID := setupTest(t)

	_, err := svc.RequestAccess(context.Background(), Actor{UserID: pendingID, Role: entity.RoleStudent}, courseID)
	if err == nil {
		t.Fatal("expected error for pending student, got nil")
	}
	if !apperror.IsCode(err, "FORBIDDEN") {
		t.Errorf("expected FORBIDDEN error, got %v", err)
	}
}

func TestRequestAccess_ApprovedStudent_Success(t *testing.T) {
	svc, reqRepo, _, _, _, approvedID, courseID := setupTest(t)

	req, err := svc.RequestAccess(context.Background(), Actor{UserID: approvedID, Role: entity.RoleStudent}, courseID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Status != entity.RequestStatusPending {
		t.Errorf("expected request status pending, got %v", req.Status)
	}

	stored, err := reqRepo.GetByStudentAndCourse(context.Background(), approvedID, courseID)
	if err != nil || stored == nil {
		t.Fatalf("expected stored request, got %v", err)
	}
}

func TestRequestAccess_DuplicatePending_Conflict(t *testing.T) {
	svc, _, _, _, _, approvedID, courseID := setupTest(t)

	_, err := svc.RequestAccess(context.Background(), Actor{UserID: approvedID, Role: entity.RoleStudent}, courseID)
	if err != nil {
		t.Fatalf("first request failed: %v", err)
	}

	_, err = svc.RequestAccess(context.Background(), Actor{UserID: approvedID, Role: entity.RoleStudent}, courseID)
	if err == nil || !apperror.IsCode(err, "CONFLICT") {
		t.Fatalf("expected CONFLICT error on duplicate request, got %v", err)
	}
}

func TestTeacherApproval_Success(t *testing.T) {
	svc, reqRepo, enrollRepo, teacherID, _, approvedID, courseID := setupTest(t)

	req, err := svc.RequestAccess(context.Background(), Actor{UserID: approvedID, Role: entity.RoleStudent}, courseID)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	err = svc.ApproveTeacherRequest(context.Background(), Actor{UserID: teacherID, Role: entity.RoleTeacher}, req.ID)
	if err != nil {
		t.Fatalf("teacher approve failed: %v", err)
	}

	updated, _ := reqRepo.GetByID(context.Background(), req.ID)
	if updated.Status != entity.RequestStatusApproved {
		t.Errorf("expected status approved, got %v", updated.Status)
	}

	enrolled, _ := enrollRepo.Exists(context.Background(), approvedID, courseID)
	if !enrolled {
		t.Error("expected student to be enrolled after teacher approval")
	}
}

func TestTeacherApproval_UnauthorizedTeacher_Forbidden(t *testing.T) {
	svc, _, _, _, _, approvedID, courseID := setupTest(t)
	otherTeacherID := uuid.New()

	req, err := svc.RequestAccess(context.Background(), Actor{UserID: approvedID, Role: entity.RoleStudent}, courseID)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	err = svc.ApproveTeacherRequest(context.Background(), Actor{UserID: otherTeacherID, Role: entity.RoleTeacher}, req.ID)
	if err == nil || !apperror.IsCode(err, "FORBIDDEN") {
		t.Fatalf("expected FORBIDDEN when unassociated teacher approves, got %v", err)
	}
}

func TestGetAccessStatus_States(t *testing.T) {
	svc, _, enrollRepo, teacherID, pendingID, approvedID, courseID := setupTest(t)

	// 1. Pending Student
	st, err := svc.GetAccessStatus(context.Background(), Actor{UserID: pendingID, Role: entity.RoleStudent}, courseID)
	if err != nil || st.Status != "account_pending" {
		t.Errorf("expected status account_pending, got %v (err: %v)", st.Status, err)
	}

	// 2. Approved Student without request
	st, err = svc.GetAccessStatus(context.Background(), Actor{UserID: approvedID, Role: entity.RoleStudent}, courseID)
	if err != nil || st.Status != "not_enrolled" {
		t.Errorf("expected status not_enrolled, got %v", st.Status)
	}

	// 3. Approved Student with pending request
	req, _ := svc.RequestAccess(context.Background(), Actor{UserID: approvedID, Role: entity.RoleStudent}, courseID)
	st, err = svc.GetAccessStatus(context.Background(), Actor{UserID: approvedID, Role: entity.RoleStudent}, courseID)
	if err != nil || st.Status != "request_pending" {
		t.Errorf("expected status request_pending, got %v", st.Status)
	}

	// 4. Approved Student after approval (enrollment exists)
	_ = svc.ApproveTeacherRequest(context.Background(), Actor{UserID: teacherID, Role: entity.RoleTeacher}, req.ID)
	st, err = svc.GetAccessStatus(context.Background(), Actor{UserID: approvedID, Role: entity.RoleStudent}, courseID)
	if err != nil || st.Status != "granted" || !st.HasAccess {
		t.Errorf("expected status granted with has_access=true, got status=%v has_access=%v", st.Status, st.HasAccess)
	}

	_ = enrollRepo
}

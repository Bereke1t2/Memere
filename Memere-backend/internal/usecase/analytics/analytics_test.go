package analytics

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

func f64(v float64) *float64 { return &v }

func snapshot(breakdown map[string][2]int) map[string]any {
	b := map[string]any{}
	for k, v := range breakdown {
		b[k] = map[string]any{"earned": float64(v[0]), "possible": float64(v[1])}
	}
	return map[string]any{"answers": map[string]any{}, "subject_breakdown": b}
}

func TestGetAttemptAnalytics_WeakAreasRankedByMarksLost(t *testing.T) {
	stu := &Actor{UserID: uuid.New(), Role: entity.RoleStudent}
	attemptID := uuid.New()
	examID := uuid.New()

	attempts := &fakeExamAttemptRepo{byID: map[uuid.UUID]*entity.ExamAttempt{
		attemptID: {
			ID:         attemptID,
			ExamID:     examID,
			StudentID:  stu.UserID,
			Status:     entity.AttemptGraded,
			Score:      f64(6),
			Percentage: f64(60),
			AnswersSnapshot: snapshot(map[string][2]int{
				"Algebra":  {1, 5}, // lost 4 — weakest
				"Geometry": {3, 4}, // lost 1
				"Calculus": {2, 2}, // lost 0 — not weak
			}),
		},
	}}
	ranking := &fakeRanking{pct: 75, ok: true}
	svc := NewService(nil, attempts, nil, ranking)

	out, err := svc.GetAttemptAnalytics(context.Background(), stu, attemptID)
	if err != nil {
		t.Fatalf("GetAttemptAnalytics: %v", err)
	}
	if len(out.WeakAreas) != 2 {
		t.Fatalf("want 2 weak areas (Calculus excluded), got %d", len(out.WeakAreas))
	}
	if out.WeakAreas[0].Key != "Algebra" || out.WeakAreas[1].Key != "Geometry" {
		t.Errorf("weak areas mis-ranked: %+v", out.WeakAreas)
	}
	if out.Percentile == nil || *out.Percentile != 75 {
		t.Errorf("percentile = %v, want 75", out.Percentile)
	}
}

func TestGetAttemptAnalytics_OtherStudentForbidden(t *testing.T) {
	owner := uuid.New()
	attemptID := uuid.New()
	attempts := &fakeExamAttemptRepo{byID: map[uuid.UUID]*entity.ExamAttempt{
		attemptID: {ID: attemptID, StudentID: owner, Status: entity.AttemptGraded},
	}}
	svc := NewService(nil, attempts, nil, nil)

	_, err := svc.GetAttemptAnalytics(context.Background(), &Actor{UserID: uuid.New(), Role: entity.RoleStudent}, attemptID)
	if !apperror.IsNotFound(err) {
		t.Fatalf("another student's attempt must not be found (IDOR), got %v", err)
	}
}

func TestGetStudentTrend_OrderedSeries(t *testing.T) {
	stu := &Actor{UserID: uuid.New(), Role: entity.RoleStudent}
	attempts := &fakeExamAttemptRepo{trend: []*entity.ExamAttempt{
		{ID: uuid.New(), Percentage: f64(40)},
		{ID: uuid.New(), Percentage: f64(55)},
		{ID: uuid.New(), Percentage: f64(70)},
	}}
	svc := NewService(nil, attempts, nil, nil)

	pts, err := svc.GetStudentTrend(context.Background(), stu, "Math")
	if err != nil {
		t.Fatalf("GetStudentTrend: %v", err)
	}
	if len(pts) != 3 || pts[0].Percentage != 40 || pts[2].Percentage != 70 {
		t.Errorf("trend series wrong: %+v", pts)
	}
}

func TestGetExamStats_TeacherOwnsCourse(t *testing.T) {
	teacher := &Actor{UserID: uuid.New(), Role: entity.RoleTeacher}
	courseID := uuid.New()
	examID := uuid.New()

	exams := &fakeExamRepo{byID: map[uuid.UUID]*entity.Exam{
		examID: {ID: examID, CourseID: &courseID},
	}}
	courses := &fakeCourseRepo{byID: map[uuid.UUID]*entity.Course{
		courseID: {ID: courseID, TeacherID: teacher.UserID},
	}}
	attempts := &fakeExamAttemptRepo{stats: repository.ExamAttemptStats{TotalAttempts: 4, AvgPercentage: 62.5, PassedCount: 3}}
	svc := NewService(exams, attempts, courses, nil)

	stats, err := svc.GetExamStats(context.Background(), teacher, examID)
	if err != nil {
		t.Fatalf("GetExamStats: %v", err)
	}
	if stats.TotalAttempts != 4 || stats.AvgPercentage != 62.5 || stats.PassRate != 75 {
		t.Errorf("stats wrong: %+v", stats)
	}
}

func TestGetExamStats_StudentForbidden(t *testing.T) {
	svc := NewService(nil, nil, nil, nil)
	_, err := svc.GetExamStats(context.Background(), &Actor{UserID: uuid.New(), Role: entity.RoleStudent}, uuid.New())
	if !apperror.IsCode(err, "FORBIDDEN") {
		t.Fatalf("student must not view exam stats, got %v", err)
	}
}

func TestGetExamStats_OtherTeacherForbidden(t *testing.T) {
	courseID := uuid.New()
	examID := uuid.New()
	exams := &fakeExamRepo{byID: map[uuid.UUID]*entity.Exam{examID: {ID: examID, CourseID: &courseID}}}
	courses := &fakeCourseRepo{byID: map[uuid.UUID]*entity.Course{courseID: {ID: courseID, TeacherID: uuid.New()}}}
	svc := NewService(exams, &fakeExamAttemptRepo{}, courses, nil)

	_, err := svc.GetExamStats(context.Background(), &Actor{UserID: uuid.New(), Role: entity.RoleTeacher}, examID)
	if !apperror.IsCode(err, "FORBIDDEN") {
		t.Fatalf("non-owning teacher must be forbidden, got %v", err)
	}
}

// ---- minimal fakes -----------------------------------------------------------

type fakeRanking struct {
	pct float64
	ok  bool
}

func (f *fakeRanking) RecordExamScore(context.Context, uuid.UUID, uuid.UUID, float64) error {
	return nil
}
func (f *fakeRanking) PercentileRank(context.Context, uuid.UUID, uuid.UUID) (float64, bool, error) {
	return f.pct, f.ok, nil
}
func (f *fakeRanking) GetTopN(context.Context, uuid.UUID, int) ([]repository.LeaderboardEntry, error) {
	return nil, nil
}
func (f *fakeRanking) GetRank(context.Context, uuid.UUID, uuid.UUID) (int64, int64, float64, bool, error) {
	return 0, 0, 0, false, nil
}
func (f *fakeRanking) RebuildFromScores(context.Context, uuid.UUID, map[uuid.UUID]float64) error {
	return nil
}

type fakeExamAttemptRepo struct {
	byID  map[uuid.UUID]*entity.ExamAttempt
	trend []*entity.ExamAttempt
	stats repository.ExamAttemptStats
}

func (f *fakeExamAttemptRepo) FindByID(_ context.Context, id, studentID uuid.UUID) (*entity.ExamAttempt, error) {
	a, ok := f.byID[id]
	if !ok || a.StudentID != studentID {
		return nil, apperror.NotFound("exam attempt not found", nil)
	}
	cp := *a
	return &cp, nil
}
func (f *fakeExamAttemptRepo) ListGradedBySubject(context.Context, uuid.UUID, string) ([]*entity.ExamAttempt, error) {
	return f.trend, nil
}
func (f *fakeExamAttemptRepo) Stats(context.Context, uuid.UUID) (repository.ExamAttemptStats, error) {
	return f.stats, nil
}
func (f *fakeExamAttemptRepo) Create(context.Context, *entity.ExamAttempt) error { return nil }
func (f *fakeExamAttemptRepo) GetActive(context.Context, uuid.UUID, uuid.UUID) (*entity.ExamAttempt, error) {
	return nil, apperror.NotFound("none", nil)
}
func (f *fakeExamAttemptRepo) ListByStudent(context.Context, uuid.UUID) ([]*entity.ExamAttempt, error) {
	return nil, nil
}
func (f *fakeExamAttemptRepo) ListExpired(context.Context, time.Time, int) ([]*entity.ExamAttempt, error) {
	return nil, nil
}
func (f *fakeExamAttemptRepo) Update(context.Context, *entity.ExamAttempt) error { return nil }
func (f *fakeExamAttemptRepo) ClaimForGrading(context.Context, *entity.ExamAttempt) (bool, error) {
	return true, nil
}
func (f *fakeExamAttemptRepo) Grade(context.Context, *entity.ExamAttempt) error { return nil }

type fakeExamRepo struct {
	byID map[uuid.UUID]*entity.Exam
}

func (f *fakeExamRepo) FindByID(_ context.Context, id uuid.UUID) (*entity.Exam, error) {
	if e, ok := f.byID[id]; ok {
		cp := *e
		return &cp, nil
	}
	return nil, apperror.NotFound("exam not found", nil)
}
func (f *fakeExamRepo) Create(context.Context, *entity.Exam) error { return nil }
func (f *fakeExamRepo) List(context.Context, repository.ExamFilter, *pagination.Cursor, int) ([]*entity.Exam, *pagination.Cursor, error) {
	return nil, nil, nil
}
func (f *fakeExamRepo) Update(context.Context, *entity.Exam) error              { return nil }
func (f *fakeExamRepo) Delete(context.Context, uuid.UUID) error                 { return nil }
func (f *fakeExamRepo) AddQuestion(context.Context, *entity.ExamQuestion) error { return nil }
func (f *fakeExamRepo) ListQuestions(context.Context, uuid.UUID) ([]*entity.ExamQuestion, error) {
	return nil, nil
}
func (f *fakeExamRepo) GetExamWithQuestions(context.Context, uuid.UUID) (*repository.ExamWithQuestions, error) {
	return nil, nil
}
func (f *fakeExamRepo) ListByCourse(_ context.Context, _ uuid.UUID) ([]*entity.Exam, error) {
	return nil, nil
}
func (f *fakeExamRepo) GetQuestionsForClient(context.Context, uuid.UUID) ([]repository.ClientExamQuestion, error) {
	return nil, nil
}

type fakeCourseRepo struct {
	byID map[uuid.UUID]*entity.Course
}

func (f *fakeCourseRepo) FindByID(_ context.Context, id uuid.UUID) (*entity.Course, error) {
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
	return nil, apperror.NotFound("not found", nil)
}

var (
	_ repository.ExamAttemptRepository = (*fakeExamAttemptRepo)(nil)
	_ repository.ExamRepository        = (*fakeExamRepo)(nil)
	_ repository.CourseRepository      = (*fakeCourseRepo)(nil)
	_ repository.ScoreRanking          = (*fakeRanking)(nil)
)

package analytics

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

func TestGetExamsOverview_AdminAggregatesTotals(t *testing.T) {
	admin := &Actor{UserID: uuid.New(), Role: entity.RoleAdmin}
	examA, examB := uuid.New(), uuid.New()

	exams := &fakeExamRepo{list: []*entity.Exam{
		{ID: examA, Title: "A"},
		{ID: examB, Title: "B"},
	}}
	attempts := &fakeExamAttemptRepo{
		overview: []repository.ExamOverviewStat{
			{ExamID: examA, TotalAttempts: 4, DistinctStudents: 3, AvgPercentage: 50, PassedCount: 2, InProgressCount: 1},
			{ExamID: examB, TotalAttempts: 6, DistinctStudents: 5, AvgPercentage: 80, PassedCount: 6},
		},
		students: 7,
		daily:    []repository.DailyAttemptCount{{Attempts: 10}},
	}
	svc := NewService(exams, attempts, &fakeCourseRepo{}, nil, nil)

	out, err := svc.GetExamsOverview(context.Background(), admin)
	if err != nil {
		t.Fatalf("GetExamsOverview: %v", err)
	}
	if len(out.Rows) != 2 {
		t.Fatalf("want 2 rows, got %d", len(out.Rows))
	}
	if out.TotalAttempts != 10 || out.InProgressCount != 1 || out.DistinctStudents != 7 {
		t.Errorf("totals wrong: %+v", out)
	}
	// Attempt-weighted mean: (50*4 + 80*6) / 10 = 68; pass rate 8/10 = 80%.
	if out.AvgPercentage != 68 || out.PassRate != 80 {
		t.Errorf("avg=%v passRate=%v, want 68 / 80", out.AvgPercentage, out.PassRate)
	}
	if len(out.Daily) != 1 {
		t.Errorf("daily series missing: %+v", out.Daily)
	}
}

func TestGetExamsOverview_TeacherScopedToOwnCourses(t *testing.T) {
	teacher := &Actor{UserID: uuid.New(), Role: entity.RoleTeacher}
	ownCourse, otherCourse := uuid.New(), uuid.New()
	ownExam := uuid.New()

	exams := &fakeExamRepo{list: []*entity.Exam{
		{ID: ownExam, CourseID: &ownCourse, Title: "mine"},
		{ID: uuid.New(), CourseID: &otherCourse, Title: "someone else's"},
		{ID: uuid.New(), Title: "standalone (admin-only)"},
	}}
	courses := &fakeCourseRepo{list: []*entity.Course{
		{ID: ownCourse, TeacherID: teacher.UserID},
		{ID: otherCourse, TeacherID: uuid.New()},
	}}
	svc := NewService(exams, &fakeExamAttemptRepo{}, courses, nil, nil)

	out, err := svc.GetExamsOverview(context.Background(), teacher)
	if err != nil {
		t.Fatalf("GetExamsOverview: %v", err)
	}
	if len(out.Rows) != 1 || out.Rows[0].Exam.ID != ownExam {
		t.Fatalf("teacher must only see exams under their own courses, got %+v", out.Rows)
	}
}

func TestGetExamsOverview_StudentForbidden(t *testing.T) {
	svc := NewService(nil, nil, nil, nil, nil)
	_, err := svc.GetExamsOverview(context.Background(), &Actor{UserID: uuid.New(), Role: entity.RoleStudent})
	if !apperror.IsCode(err, "FORBIDDEN") {
		t.Fatalf("student must not view the overview, got %v", err)
	}
}

func TestGetExamInsights_DenseDistribution(t *testing.T) {
	admin := &Actor{UserID: uuid.New(), Role: entity.RoleAdmin}
	examID := uuid.New()

	exams := &fakeExamRepo{byID: map[uuid.UUID]*entity.Exam{examID: {ID: examID}}}
	attempts := &fakeExamAttemptRepo{
		overview: []repository.ExamOverviewStat{
			{ExamID: examID, TotalAttempts: 3, DistinctStudents: 2, AvgPercentage: 70, PassedCount: 2},
		},
		buckets: []repository.ScoreBucket{{Bucket: 5, Attempts: 2}, {Bucket: 9, Attempts: 1}},
	}
	svc := NewService(exams, attempts, &fakeCourseRepo{}, nil, nil)

	out, err := svc.GetExamInsights(context.Background(), admin, examID)
	if err != nil {
		t.Fatalf("GetExamInsights: %v", err)
	}
	if out.Stats.TotalAttempts != 3 || out.PassRate < 66 || out.PassRate > 67 {
		t.Errorf("stats wrong: %+v passRate=%v", out.Stats, out.PassRate)
	}
	want := [10]int{5: 2, 9: 1}
	if out.Distribution != want {
		t.Errorf("distribution = %v, want %v", out.Distribution, want)
	}
}

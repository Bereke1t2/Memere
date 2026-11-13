package analytics

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/pagination"
)

// overviewMaxExams caps how many exams the overview aggregates so a huge
// catalog cannot make the dashboard page unbounded.
const overviewMaxExams = 1000

// dailyWindowDays is how far back the attempts-per-day series reaches.
const dailyWindowDays = 30

// ExamOverviewRow pairs an exam with its attempt aggregates for the overview.
type ExamOverviewRow struct {
	Exam  *entity.Exam
	Stats repository.ExamOverviewStat
}

// ExamsOverview is the teacher/admin analytics dashboard payload: one row per
// visible exam plus cross-exam totals and the attempts-per-day series.
type ExamsOverview struct {
	Rows             []ExamOverviewRow
	TotalAttempts    int
	DistinctStudents int64
	AvgPercentage    float64 // attempt-weighted mean across all rows
	PassRate         float64 // 0–100 across all graded attempts
	InProgressCount  int
	Daily            []repository.DailyAttemptCount
}

// ExamInsights is the per-exam analytics detail: the aggregate stats plus a
// dense 10-bucket score histogram and the attempts-per-day series.
type ExamInsights struct {
	Exam         *entity.Exam
	Stats        repository.ExamOverviewStat
	PassRate     float64 // 0–100
	Distribution [10]int // bucket n covers [n*10, n*10+10)%, 9 includes 100
	Daily        []repository.DailyAttemptCount
}

// GetExamsOverview returns per-exam attempt aggregates for every exam the
// actor may see stats for: all exams for an admin, the exams under their own
// courses for a teacher (standalone exams stay admin-only, matching
// GetExamStats). Rows keep the catalog's created_at DESC order.
func (s *Service) GetExamsOverview(ctx context.Context, actor *Actor) (*ExamsOverview, error) {
	if !actor.isTeacherOrAdmin() {
		return nil, apperror.Forbidden("only teachers may view exam statistics", nil)
	}

	exams, err := s.listVisibleExams(ctx, actor)
	if err != nil {
		return nil, err
	}

	out := &ExamsOverview{Rows: make([]ExamOverviewRow, 0, len(exams))}
	if len(exams) == 0 {
		return out, nil
	}

	ids := make([]uuid.UUID, len(exams))
	for i, e := range exams {
		ids[i] = e.ID
	}

	stats, err := s.examAttempts.StatsByExamIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	byExam := make(map[uuid.UUID]repository.ExamOverviewStat, len(stats))
	for _, st := range stats {
		byExam[st.ExamID] = st
	}

	var passed int
	var pctSum float64
	for _, e := range exams {
		st := byExam[e.ID]
		st.ExamID = e.ID // zero row when the exam had no attempts
		out.Rows = append(out.Rows, ExamOverviewRow{Exam: e, Stats: st})
		out.TotalAttempts += st.TotalAttempts
		out.InProgressCount += st.InProgressCount
		passed += st.PassedCount
		pctSum += st.AvgPercentage * float64(st.TotalAttempts)
	}
	if out.TotalAttempts > 0 {
		out.AvgPercentage = pctSum / float64(out.TotalAttempts)
		out.PassRate = float64(passed) / float64(out.TotalAttempts) * 100
	}

	if out.DistinctStudents, err = s.examAttempts.CountDistinctStudents(ctx, ids); err != nil {
		return nil, err
	}
	since := time.Now().AddDate(0, 0, -dailyWindowDays)
	if out.Daily, err = s.examAttempts.DailyAttempts(ctx, ids, since); err != nil {
		return nil, err
	}
	return out, nil
}

// GetExamInsights returns the per-exam analytics detail for the teacher who
// owns the exam's course or an admin, under the same authorization rule as
// GetExamStats.
func (s *Service) GetExamInsights(ctx context.Context, actor *Actor, examID uuid.UUID) (*ExamInsights, error) {
	exam, err := s.authorizeExamStats(ctx, actor, examID)
	if err != nil {
		return nil, err
	}

	stats, err := s.examAttempts.StatsByExamIDs(ctx, []uuid.UUID{examID})
	if err != nil {
		return nil, err
	}
	out := &ExamInsights{Exam: exam, Stats: repository.ExamOverviewStat{ExamID: examID}}
	if len(stats) > 0 {
		out.Stats = stats[0]
	}
	if out.Stats.TotalAttempts > 0 {
		out.PassRate = float64(out.Stats.PassedCount) / float64(out.Stats.TotalAttempts) * 100
	}

	buckets, err := s.examAttempts.ScoreDistribution(ctx, examID)
	if err != nil {
		return nil, err
	}
	for _, b := range buckets {
		if b.Bucket >= 0 && b.Bucket < len(out.Distribution) {
			out.Distribution[b.Bucket] = b.Attempts
		}
	}

	since := time.Now().AddDate(0, 0, -dailyWindowDays)
	if out.Daily, err = s.examAttempts.DailyAttempts(ctx, []uuid.UUID{examID}, since); err != nil {
		return nil, err
	}
	return out, nil
}

// authorizeExamStats loads the exam and asserts the actor may view its cohort
// statistics: any admin, or the teacher who owns its parent course. Standalone
// (course-less) exams are admin-only.
func (s *Service) authorizeExamStats(ctx context.Context, actor *Actor, examID uuid.UUID) (*entity.Exam, error) {
	if !actor.isTeacherOrAdmin() {
		return nil, apperror.Forbidden("only teachers may view exam statistics", nil)
	}
	exam, err := s.exams.FindByID(ctx, examID)
	if err != nil {
		return nil, err
	}
	if actor.Role != entity.RoleAdmin {
		if exam.CourseID == nil {
			return nil, apperror.Forbidden("only an admin may view stats for a standalone exam", nil)
		}
		course, err := s.courses.FindByID(ctx, *exam.CourseID)
		if err != nil {
			return nil, err
		}
		if course.TeacherID != actor.UserID {
			return nil, apperror.Forbidden("you do not own this exam's course", nil)
		}
	}
	return exam, nil
}

// listVisibleExams pages through the exam catalog (up to overviewMaxExams) and,
// for a teacher, keeps only exams under one of their own courses.
func (s *Service) listVisibleExams(ctx context.Context, actor *Actor) ([]*entity.Exam, error) {
	var ownCourses map[uuid.UUID]bool
	if actor.Role != entity.RoleAdmin {
		var err error
		if ownCourses, err = s.listOwnCourseIDs(ctx, actor.UserID); err != nil {
			return nil, err
		}
	}

	var out []*entity.Exam
	var cursor *pagination.Cursor
	for len(out) < overviewMaxExams {
		exams, next, err := s.exams.List(ctx, repository.ExamFilter{}, cursor, pagination.MaxLimit)
		if err != nil {
			return nil, err
		}
		for _, e := range exams {
			if ownCourses != nil && (e.CourseID == nil || !ownCourses[*e.CourseID]) {
				continue
			}
			out = append(out, e)
		}
		if next == nil {
			break
		}
		cursor = next
	}
	return out, nil
}

// listOwnCourseIDs collects the IDs of every course the teacher owns.
func (s *Service) listOwnCourseIDs(ctx context.Context, teacherID uuid.UUID) (map[uuid.UUID]bool, error) {
	ids := make(map[uuid.UUID]bool)
	var cursor *pagination.Cursor
	for {
		courses, next, err := s.courses.List(ctx, repository.CourseFilter{TeacherID: &teacherID}, cursor, pagination.MaxLimit)
		if err != nil {
			return nil, err
		}
		for _, c := range courses {
			ids[c.ID] = true
		}
		if next == nil {
			break
		}
		cursor = next
	}
	return ids, nil
}

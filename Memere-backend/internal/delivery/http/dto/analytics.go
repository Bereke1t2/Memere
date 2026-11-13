package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/analytics"
)

// AnalyticsSubjectScore is a per-subject (or per-topic) tally in an analytics
// response, carrying the subject key alongside the earned/possible marks.
type AnalyticsSubjectScore struct {
	Key      string `json:"key"`
	Earned   int    `json:"earned"`
	Possible int    `json:"possible"`
}

// AttemptAnalyticsResponse is the §9.3 per-attempt analytics view.
type AttemptAnalyticsResponse struct {
	AttemptID        uuid.UUID               `json:"attempt_id"`
	ExamID           uuid.UUID               `json:"exam_id"`
	Score            float64                 `json:"score"`
	Percentage       float64                 `json:"percentage"`
	SubjectBreakdown []AnalyticsSubjectScore `json:"subject_breakdown"`
	WeakAreas        []AnalyticsSubjectScore `json:"weak_areas"`
	Percentile       *float64                `json:"percentile"`
}

func newAnalyticsSubjectScores(in []analytics.SubjectScore) []AnalyticsSubjectScore {
	out := make([]AnalyticsSubjectScore, 0, len(in))
	for _, s := range in {
		out = append(out, AnalyticsSubjectScore{Key: s.Key, Earned: s.Earned, Possible: s.Possible})
	}
	return out
}

// NewAttemptAnalyticsResponse maps the analytics usecase view to its response.
func NewAttemptAnalyticsResponse(a *analytics.AttemptAnalytics) AttemptAnalyticsResponse {
	return AttemptAnalyticsResponse{
		AttemptID:        a.AttemptID,
		ExamID:           a.ExamID,
		Score:            a.Score,
		Percentage:       a.Percentage,
		SubjectBreakdown: newAnalyticsSubjectScores(a.SubjectBreakdown),
		WeakAreas:        newAnalyticsSubjectScores(a.WeakAreas),
		Percentile:       a.Percentile,
	}
}

// TrendPointResponse is one attempt's score in a subject trend.
type TrendPointResponse struct {
	AttemptID  uuid.UUID `json:"attempt_id"`
	ExamID     uuid.UUID `json:"exam_id"`
	Percentage float64   `json:"percentage"`
}

// TrendResponse is the ordered (oldest-first) score series for a subject.
type TrendResponse struct {
	Subject string               `json:"subject"`
	Points  []TrendPointResponse `json:"points"`
}

// NewTrendResponse maps the trend points to their response.
func NewTrendResponse(subject string, points []analytics.TrendPoint) TrendResponse {
	out := make([]TrendPointResponse, 0, len(points))
	for _, p := range points {
		out = append(out, TrendPointResponse{AttemptID: p.AttemptID, ExamID: p.ExamID, Percentage: p.Percentage})
	}
	return TrendResponse{Subject: subject, Points: out}
}

// ExamStatsResponse is the teacher/admin aggregate for an exam.
type ExamStatsResponse struct {
	ExamID        uuid.UUID `json:"exam_id"`
	TotalAttempts int       `json:"total_attempts"`
	AvgPercentage float64   `json:"avg_percentage"`
	PassRate      float64   `json:"pass_rate"`
}

// NewExamStatsResponse maps the exam stats view to its response.
func NewExamStatsResponse(s *analytics.ExamStats) ExamStatsResponse {
	return ExamStatsResponse{
		ExamID:        s.ExamID,
		TotalAttempts: s.TotalAttempts,
		AvgPercentage: s.AvgPercentage,
		PassRate:      s.PassRate,
	}
}

// StudentPointsResponse is a student's cumulative points across all graded
// quizzes and exams (best attempt per item), with per-source subtotals, the
// number of distinct quizzes/exams completed, and an overall average percentage.
type StudentPointsResponse struct {
	TotalPoints   float64 `json:"total_points"`
	QuizPoints    float64 `json:"quiz_points"`
	ExamPoints    float64 `json:"exam_points"`
	AvgPercentage float64 `json:"avg_percentage"`
	QuizCount     int64   `json:"quiz_count"`
	ExamCount     int64   `json:"exam_count"`
}

// NewStudentPointsResponse maps the cumulative points view to its response.
func NewStudentPointsResponse(p *analytics.StudentPoints) StudentPointsResponse {
	return StudentPointsResponse{
		TotalPoints:   p.TotalPoints,
		QuizPoints:    p.QuizPoints,
		ExamPoints:    p.ExamPoints,
		AvgPercentage: p.AvgPercentage,
		QuizCount:     p.QuizCount,
		ExamCount:     p.ExamCount,
	}
}

// ExamOverviewStatResponse is one exam's attempt aggregates in the analytics
// overview and insights responses.
type ExamOverviewStatResponse struct {
	TotalAttempts    int        `json:"total_attempts"`
	DistinctStudents int        `json:"distinct_students"`
	AvgPercentage    float64    `json:"avg_percentage"`
	PassedCount      int        `json:"passed_count"`
	PassRate         float64    `json:"pass_rate"`
	InProgressCount  int        `json:"in_progress_count"`
	LastSubmittedAt  *time.Time `json:"last_submitted_at"`
}

func newExamOverviewStatResponse(st repository.ExamOverviewStat) ExamOverviewStatResponse {
	out := ExamOverviewStatResponse{
		TotalAttempts:    st.TotalAttempts,
		DistinctStudents: st.DistinctStudents,
		AvgPercentage:    st.AvgPercentage,
		PassedCount:      st.PassedCount,
		InProgressCount:  st.InProgressCount,
		LastSubmittedAt:  st.LastSubmittedAt,
	}
	if st.TotalAttempts > 0 {
		out.PassRate = float64(st.PassedCount) / float64(st.TotalAttempts) * 100
	}
	return out
}

// DailyAttemptsResponse is one day's attempt count in a time series.
type DailyAttemptsResponse struct {
	Date     string `json:"date"` // YYYY-MM-DD
	Attempts int    `json:"attempts"`
}

func newDailyAttemptsResponse(in []repository.DailyAttemptCount) []DailyAttemptsResponse {
	out := make([]DailyAttemptsResponse, 0, len(in))
	for _, d := range in {
		out = append(out, DailyAttemptsResponse{Date: d.Day.Format("2006-01-02"), Attempts: d.Attempts})
	}
	return out
}

// ExamOverviewRowResponse pairs an exam with its aggregates.
type ExamOverviewRowResponse struct {
	Exam  ExamResponse             `json:"exam"`
	Stats ExamOverviewStatResponse `json:"stats"`
}

// ExamsOverviewResponse is the teacher/admin analytics dashboard payload.
type ExamsOverviewResponse struct {
	Exams            []ExamOverviewRowResponse `json:"exams"`
	TotalAttempts    int                       `json:"total_attempts"`
	DistinctStudents int64                     `json:"distinct_students"`
	AvgPercentage    float64                   `json:"avg_percentage"`
	PassRate         float64                   `json:"pass_rate"`
	InProgressCount  int                       `json:"in_progress_count"`
	Daily            []DailyAttemptsResponse   `json:"daily"`
}

// NewExamsOverviewResponse maps the overview usecase view to its response.
func NewExamsOverviewResponse(o *analytics.ExamsOverview) ExamsOverviewResponse {
	rows := make([]ExamOverviewRowResponse, 0, len(o.Rows))
	for _, r := range o.Rows {
		rows = append(rows, ExamOverviewRowResponse{
			Exam:  NewExamResponse(r.Exam),
			Stats: newExamOverviewStatResponse(r.Stats),
		})
	}
	return ExamsOverviewResponse{
		Exams:            rows,
		TotalAttempts:    o.TotalAttempts,
		DistinctStudents: o.DistinctStudents,
		AvgPercentage:    o.AvgPercentage,
		PassRate:         o.PassRate,
		InProgressCount:  o.InProgressCount,
		Daily:            newDailyAttemptsResponse(o.Daily),
	}
}

// DistributionBucketResponse is one 10-point histogram bucket of graded scores.
type DistributionBucketResponse struct {
	Range    string `json:"range"` // e.g. "70-79"; the last bucket is "90-100"
	Attempts int    `json:"attempts"`
}

// ExamInsightsResponse is the per-exam analytics detail served on
// GET /exams/:id/stats. It is a superset of the original stats payload, so the
// legacy exam_id / total_attempts / avg_percentage / pass_rate fields stay.
type ExamInsightsResponse struct {
	ExamID           uuid.UUID                    `json:"exam_id"`
	TotalAttempts    int                          `json:"total_attempts"`
	AvgPercentage    float64                      `json:"avg_percentage"`
	PassRate         float64                      `json:"pass_rate"`
	DistinctStudents int                          `json:"distinct_students"`
	PassedCount      int                          `json:"passed_count"`
	InProgressCount  int                          `json:"in_progress_count"`
	LastSubmittedAt  *time.Time                   `json:"last_submitted_at"`
	Distribution     []DistributionBucketResponse `json:"distribution"` // dense: all 10 buckets
	Daily            []DailyAttemptsResponse      `json:"daily"`
}

// NewExamInsightsResponse maps the insights usecase view to its response.
func NewExamInsightsResponse(in *analytics.ExamInsights) ExamInsightsResponse {
	dist := make([]DistributionBucketResponse, 0, len(in.Distribution))
	for i, n := range in.Distribution {
		rangeLabel := fmt.Sprintf("%d-%d", i*10, i*10+9)
		if i == len(in.Distribution)-1 {
			rangeLabel = "90-100"
		}
		dist = append(dist, DistributionBucketResponse{Range: rangeLabel, Attempts: n})
	}
	return ExamInsightsResponse{
		ExamID:           in.Stats.ExamID,
		TotalAttempts:    in.Stats.TotalAttempts,
		AvgPercentage:    in.Stats.AvgPercentage,
		PassRate:         in.PassRate,
		DistinctStudents: in.Stats.DistinctStudents,
		PassedCount:      in.Stats.PassedCount,
		InProgressCount:  in.Stats.InProgressCount,
		LastSubmittedAt:  in.Stats.LastSubmittedAt,
		Distribution:     dist,
		Daily:            newDailyAttemptsResponse(in.Daily),
	}
}

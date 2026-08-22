package dto

import (
	"github.com/google/uuid"

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

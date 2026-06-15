// Package analytics holds the read-only §9.3 scoring-analytics usecases: per-
// attempt analysis (score, subject breakdown, weak areas, percentile), a
// student's score trend for a subject, and teacher/admin exam statistics. These
// usecases are strictly derived — they never mutate attempt rows — and they
// never log answer contents (Non-Negotiable #8).
package analytics

import (
	"context"
	"sort"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// Actor is the authenticated caller's identity.
type Actor struct {
	UserID uuid.UUID
	Role   entity.Role
}

func (a *Actor) isTeacherOrAdmin() bool {
	return a != nil && (a.Role == entity.RoleTeacher || a.Role == entity.RoleAdmin)
}

// Service implements the analytics usecases over the exam repositories and the
// Redis score ranking. It is read-only.
type Service struct {
	exams        repository.ExamRepository
	examAttempts repository.ExamAttemptRepository
	courses      repository.CourseRepository
	ranking      repository.ScoreRanking
}

// NewService wires the analytics usecase.
func NewService(
	exams repository.ExamRepository,
	examAttempts repository.ExamAttemptRepository,
	courses repository.CourseRepository,
	ranking repository.ScoreRanking,
) *Service {
	return &Service{exams: exams, examAttempts: examAttempts, courses: courses, ranking: ranking}
}

// SubjectScore is a per-subject (or per-topic) earned/possible tally.
type SubjectScore struct {
	Key      string
	Earned   int
	Possible int
}

// AttemptAnalytics is the §9.3 per-attempt view: the raw score and percentage,
// the subject breakdown, the weakest topics, and the percentile rank.
type AttemptAnalytics struct {
	AttemptID        uuid.UUID
	ExamID           uuid.UUID
	Score            float64
	Percentage       float64
	SubjectBreakdown []SubjectScore
	WeakAreas        []SubjectScore // ranked weakest-first by lost marks
	Percentile       *float64       // nil when no ranking is recorded yet
}

// TrendPoint is one attempt's score in a student's subject trend.
type TrendPoint struct {
	AttemptID  uuid.UUID
	ExamID     uuid.UUID
	Percentage float64
}

// ExamStats is the teacher/admin aggregate for an exam.
type ExamStats struct {
	ExamID        uuid.UUID
	TotalAttempts int
	AvgPercentage float64
	PassRate      float64 // 0–100
}

// GetAttemptAnalytics returns the §9.3 metrics for one of the actor's own
// attempts. It reads the per-subject breakdown persisted in answers_snapshot,
// derives weak areas (subjects where marks were lost, ranked), and fetches the
// percentile from Redis. It never mutates the attempt.
func (s *Service) GetAttemptAnalytics(ctx context.Context, actor *Actor, attemptID uuid.UUID) (*AttemptAnalytics, error) {
	if actor == nil {
		return nil, apperror.Unauthorized("authentication required", nil)
	}
	attempt, err := s.examAttempts.FindByID(ctx, attemptID, actor.UserID)
	if err != nil {
		return nil, err
	}
	if attempt.Status != entity.AttemptGraded {
		return nil, apperror.Conflict("attempt is not graded yet", nil)
	}

	breakdown := subjectBreakdownFromSnapshot(attempt.AnswersSnapshot)
	weak := weakAreas(breakdown)

	out := &AttemptAnalytics{
		AttemptID:        attempt.ID,
		ExamID:           attempt.ExamID,
		SubjectBreakdown: breakdown,
		WeakAreas:        weak,
	}
	if attempt.Score != nil {
		out.Score = *attempt.Score
	}
	if attempt.Percentage != nil {
		out.Percentage = *attempt.Percentage
	}

	if s.ranking != nil {
		if pct, ok, err := s.ranking.PercentileRank(ctx, attempt.ExamID, actor.UserID); err == nil && ok {
			p := pct
			out.Percentile = &p
		}
	}
	return out, nil
}

// GetStudentTrend returns the actor's graded-attempt scores for a subject, oldest
// first, so the client can plot improvement over time (§9.3 trend analysis).
func (s *Service) GetStudentTrend(ctx context.Context, actor *Actor, subject string) ([]TrendPoint, error) {
	if actor == nil {
		return nil, apperror.Unauthorized("authentication required", nil)
	}
	attempts, err := s.examAttempts.ListGradedBySubject(ctx, actor.UserID, subject)
	if err != nil {
		return nil, err
	}
	points := make([]TrendPoint, 0, len(attempts))
	for _, a := range attempts {
		p := TrendPoint{AttemptID: a.ID, ExamID: a.ExamID}
		if a.Percentage != nil {
			p.Percentage = *a.Percentage
		}
		points = append(points, p)
	}
	return points, nil
}

// GetExamStats returns aggregate stats for an exam, for the teacher who owns its
// course or an admin (§9.3). Students may not see cohort stats.
func (s *Service) GetExamStats(ctx context.Context, actor *Actor, examID uuid.UUID) (*ExamStats, error) {
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

	stats, err := s.examAttempts.Stats(ctx, examID)
	if err != nil {
		return nil, err
	}
	out := &ExamStats{
		ExamID:        examID,
		TotalAttempts: stats.TotalAttempts,
		AvgPercentage: stats.AvgPercentage,
	}
	if stats.TotalAttempts > 0 {
		out.PassRate = float64(stats.PassedCount) / float64(stats.TotalAttempts) * 100
	}
	return out, nil
}

// subjectBreakdownFromSnapshot reads the per-subject earned/possible tally the
// grading core persisted under "subject_breakdown" in answers_snapshot, returned
// sorted by subject key for a stable response.
func subjectBreakdownFromSnapshot(snapshot map[string]any) []SubjectScore {
	raw, ok := snapshot["subject_breakdown"].(map[string]any)
	if !ok {
		return nil
	}
	out := make([]SubjectScore, 0, len(raw))
	for key, v := range raw {
		m, ok := v.(map[string]any)
		if !ok {
			continue
		}
		out = append(out, SubjectScore{
			Key:      key,
			Earned:   toInt(m["earned"]),
			Possible: toInt(m["possible"]),
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Key < out[j].Key })
	return out
}

// weakAreas ranks subjects/topics by marks lost (possible-earned), weakest first,
// keeping only those where the student actually lost marks. Ties break by key for
// determinism.
func weakAreas(breakdown []SubjectScore) []SubjectScore {
	weak := make([]SubjectScore, 0, len(breakdown))
	for _, s := range breakdown {
		if s.Possible-s.Earned > 0 {
			weak = append(weak, s)
		}
	}
	sort.Slice(weak, func(i, j int) bool {
		li := weak[i].Possible - weak[i].Earned
		lj := weak[j].Possible - weak[j].Earned
		if li != lj {
			return li > lj
		}
		return weak[i].Key < weak[j].Key
	})
	return weak
}

// toInt coerces a JSONB-decoded number (float64) or int into an int.
// LeaderboardResult is the response for the GET /exams/:id/leaderboard route.
type LeaderboardResult struct {
	TopN   []LeaderboardRow `json:"top_n"`
	MyRank int64            `json:"my_rank"`   // 0 when caller has no score
	MyScore float64         `json:"my_score"`  // 0 when caller has no score
	Total  int64            `json:"total"`
}

// LeaderboardRow is one entry in the leaderboard.
type LeaderboardRow struct {
	Rank      int64   `json:"rank"`
	StudentID string  `json:"student_id"`
	Score     float64 `json:"score"`
}

// GetLeaderboard returns the top limit students for an exam (any authenticated
// caller may view), plus the caller's own rank when they have a score.
func (s *Service) GetLeaderboard(ctx context.Context, actor *Actor, examID uuid.UUID, limit int) (*LeaderboardResult, error) {
	if actor == nil {
		return nil, apperror.Unauthorized("authentication required", nil)
	}
	entries, err := s.ranking.GetTopN(ctx, examID, limit)
	if err != nil {
		return nil, err
	}
	rows := make([]LeaderboardRow, 0, len(entries))
	for _, e := range entries {
		rows = append(rows, LeaderboardRow{
			Rank:      e.Rank,
			StudentID: e.StudentID.String(),
			Score:     e.Score,
		})
	}
	result := &LeaderboardResult{TopN: rows}
	rank, total, myScore, ok, err := s.ranking.GetRank(ctx, examID, actor.UserID)
	if err != nil {
		return nil, err
	}
	if ok {
		result.MyRank = rank
		result.MyScore = myScore
		result.Total = total
	}
	return result, nil
}

func toInt(v any) int {
	switch n := v.(type) {
	case float64:
		return int(n)
	case int:
		return n
	default:
		return 0
	}
}

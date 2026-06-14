package dto

import (
	"time"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/progress"
)

// VideoProgressRequest is the body of PUT /lessons/:id/video-progress.
type VideoProgressRequest struct {
	PositionSeconds int `json:"position_seconds"`
}

// LessonProgressResponse is the projection of entity.LessonProgress.
type LessonProgressResponse struct {
	LessonID             string     `json:"lesson_id"`
	IsCompleted          bool       `json:"is_completed"`
	CompletedAt          *time.Time `json:"completed_at,omitempty"`
	VideoProgressSeconds int        `json:"video_progress_seconds"`
}

// CourseProgressResponse is the response for GET /courses/:id/progress.
type CourseProgressResponse struct {
	CourseID         string                   `json:"course_id"`
	CompletedLessons int                      `json:"completed_lessons"`
	TotalLessons     int                      `json:"total_lessons"`
	PercentComplete  string                   `json:"percent_complete"`
	CompletedAt      *time.Time               `json:"completed_at,omitempty"`
	Lessons          []LessonProgressResponse `json:"lessons"`
}

// DashboardItemResponse is one enrolled course's progress summary.
type DashboardItemResponse struct {
	CourseID         string     `json:"course_id"`
	CourseTitle      string     `json:"course_title"`
	PercentComplete  float64    `json:"percent_complete"`
	CompletedAt      *time.Time `json:"completed_at,omitempty"`
	CompletedLessons int        `json:"completed_lessons"`
	TotalLessons     int        `json:"total_lessons"`
}

// DashboardResponse is returned by GET /me/dashboard.
type DashboardResponse struct {
	Courses       []DashboardItemResponse `json:"courses"`
	CurrentStreak int                     `json:"current_streak"`
	LongestStreak int                     `json:"longest_streak"`
}

// StreakResponse is returned by GET /me/streak.
type StreakResponse struct {
	CurrentStreak int        `json:"current_streak"`
	LongestStreak int        `json:"longest_streak"`
	LastStudyDate *time.Time `json:"last_study_date,omitempty"`
}

// NewLessonProgressResponse maps a domain entity to the wire DTO.
func NewLessonProgressResponse(p *entity.LessonProgress) LessonProgressResponse {
	return LessonProgressResponse{
		LessonID:             p.LessonID.String(),
		IsCompleted:          p.IsCompleted,
		CompletedAt:          p.CompletedAt,
		VideoProgressSeconds: p.VideoProgressSeconds,
	}
}

// NewCourseProgressResponse builds the course progress DTO.
func NewCourseProgressResponse(cp *entity.CourseProgress, lessons []*entity.LessonProgress) CourseProgressResponse {
	items := make([]LessonProgressResponse, 0, len(lessons))
	for _, l := range lessons {
		items = append(items, NewLessonProgressResponse(l))
	}
	return CourseProgressResponse{
		CourseID:         cp.CourseID.String(),
		CompletedLessons: cp.CompletedLessons,
		TotalLessons:     cp.TotalLessons,
		PercentComplete:  cp.PercentComplete.StringFixed(2),
		CompletedAt:      cp.CompletedAt,
		Lessons:          items,
	}
}

// NewDashboardResponse builds the dashboard DTO.
func NewDashboardResponse(d *progress.Dashboard) DashboardResponse {
	items := make([]DashboardItemResponse, 0, len(d.Items))
	for _, item := range d.Items {
		items = append(items, DashboardItemResponse{
			CourseID:         item.CourseID.String(),
			CourseTitle:      item.CourseTitle,
			PercentComplete:  item.PercentComplete,
			CompletedAt:      item.CompletedAt,
			CompletedLessons: item.CompletedLessons,
			TotalLessons:     item.TotalLessons,
		})
	}
	return DashboardResponse{
		Courses:       items,
		CurrentStreak: d.CurrentStreak,
		LongestStreak: d.LongestStreak,
	}
}

// NewStreakResponse builds the streak DTO.
func NewStreakResponse(st *entity.StudyStreak) StreakResponse {
	return StreakResponse{
		CurrentStreak: st.CurrentStreak,
		LongestStreak: st.LongestStreak,
		LastStudyDate: st.LastStudyDate,
	}
}

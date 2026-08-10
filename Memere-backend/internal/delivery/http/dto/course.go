package dto

import (
	"time"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
)

// CreateCourseRequest is the body of POST /courses.
type CreateCourseRequest struct {
	Title            string         `json:"title"`
	Description      string         `json:"description"`
	ShortDescription *string        `json:"short_description"`
	Subject          string         `json:"subject"`
	Grade            int            `json:"grade"`
	ThumbnailURL     *string        `json:"thumbnail_url"`
	Price            float64        `json:"price"`
	Currency         string         `json:"currency"`
	IsFree           bool           `json:"is_free"`
	Language         string         `json:"language"`
	Level            string         `json:"level"`
	Metadata         map[string]any `json:"metadata"`
}

// UpdateCourseRequest is the body of PUT /courses/:id. Pointer fields are
// optional: a nil field leaves the existing value unchanged (partial update).
type UpdateCourseRequest struct {
	Title            *string        `json:"title"`
	Description      *string        `json:"description"`
	ShortDescription *string        `json:"short_description"`
	Subject          *string        `json:"subject"`
	Grade            *int           `json:"grade"`
	ThumbnailURL     *string        `json:"thumbnail_url"`
	Price            *float64       `json:"price"`
	IsFree           *bool          `json:"is_free"`
	Language         *string        `json:"language"`
	Level            *string        `json:"level"`
	Metadata         map[string]any `json:"metadata"`
}

// CreateSectionRequest is the body of POST /courses/:id/sections.
type CreateSectionRequest struct {
	Title       string  `json:"title"`
	Description *string `json:"description"`
	IsPublished bool    `json:"is_published"`
	OrderIndex  *int    `json:"order_index"`
}

// CreateLessonRequest is the body of POST /sections/:id/lessons.
type CreateLessonRequest struct {
	Title           string  `json:"title"`
	Type            string  `json:"type"`
	IsFreePreview   bool    `json:"is_free_preview"`
	DurationSeconds int     `json:"duration_seconds"`
	IsPublished     bool    `json:"is_published"`
	OrderIndex      *int    `json:"order_index"`
	Content         *string `json:"content,omitempty"`
	PdfURL          *string `json:"pdf_url,omitempty"`
}

// CourseResponse is the flat projection of a course (list views).
type CourseResponse struct {
	ID                   string         `json:"id"`
	TeacherID            string         `json:"teacher_id"`
	Title                string         `json:"title"`
	Slug                 string         `json:"slug"`
	Description          string         `json:"description"`
	ShortDescription     *string        `json:"short_description,omitempty"`
	Subject              string         `json:"subject"`
	Grade                int            `json:"grade"`
	ThumbnailURL         *string        `json:"thumbnail_url,omitempty"`
	Price                float64        `json:"price"`
	Currency             string         `json:"currency"`
	IsFree               bool           `json:"is_free"`
	IsPublished          bool           `json:"is_published"`
	Language             string         `json:"language"`
	Level                string         `json:"level"`
	TotalDurationSeconds int            `json:"total_duration_seconds"`
	TotalLessons         int            `json:"total_lessons"`
	RatingAvg            float64        `json:"rating_avg"`
	EnrollmentCount      int            `json:"enrollment_count"`
	Metadata             map[string]any `json:"metadata,omitempty"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
}

// SectionResponse is the projection of a course section.
type SectionResponse struct {
	ID          string    `json:"id"`
	CourseID    string    `json:"course_id"`
	Title       string    `json:"title"`
	Description *string   `json:"description,omitempty"`
	OrderIndex  int       `json:"order_index"`
	IsPublished bool      `json:"is_published"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// LessonResponse is the projection of a lesson.
type LessonResponse struct {
	ID              string    `json:"id"`
	SectionID       string    `json:"section_id"`
	CourseID        string    `json:"course_id"`
	Title           string    `json:"title"`
	Type            string    `json:"type"`
	OrderIndex      int       `json:"order_index"`
	IsFreePreview   bool      `json:"is_free_preview"`
	DurationSeconds int       `json:"duration_seconds"`
	IsPublished     bool      `json:"is_published"`
	VideoID         *string   `json:"video_id,omitempty"`
	Content         *string   `json:"content,omitempty"`
	PdfURL          *string   `json:"pdf_url,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// SectionDetailResponse is a section with its ordered lessons (nested view).
type SectionDetailResponse struct {
	SectionResponse
	Lessons []LessonResponse `json:"lessons"`
}

// CourseDetailResponse is a course with its full section/lesson tree.
type CourseDetailResponse struct {
	CourseResponse
	Sections []SectionDetailResponse `json:"sections"`
}

// NewCourseResponse maps a domain course to its flat projection.
func NewCourseResponse(c *entity.Course) CourseResponse {
	return CourseResponse{
		ID:                   c.ID.String(),
		TeacherID:            c.TeacherID.String(),
		Title:                c.Title,
		Slug:                 c.Slug,
		Description:          c.Description,
		ShortDescription:     c.ShortDescription,
		Subject:              c.Subject,
		Grade:                c.Grade,
		ThumbnailURL:         c.ThumbnailURL,
		Price:                c.Price,
		Currency:             c.Currency,
		IsFree:               c.IsFree,
		IsPublished:          c.IsPublished,
		Language:             c.Language,
		Level:                string(c.Level),
		TotalDurationSeconds: c.TotalDurationSeconds,
		TotalLessons:         c.TotalLessons,
		RatingAvg:            c.RatingAvg,
		EnrollmentCount:      c.EnrollmentCount,
		Metadata:             c.Metadata,
		CreatedAt:            c.CreatedAt,
		UpdatedAt:            c.UpdatedAt,
	}
}

// NewSectionResponse maps a domain section to its projection.
func NewSectionResponse(s *entity.CourseSection) SectionResponse {
	return SectionResponse{
		ID:          s.ID.String(),
		CourseID:    s.CourseID.String(),
		Title:       s.Title,
		Description: s.Description,
		OrderIndex:  s.OrderIndex,
		IsPublished: s.IsPublished,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}

// NewLessonResponse maps a domain lesson to its projection.
func NewLessonResponse(l *entity.Lesson) LessonResponse {
	var videoID *string
	if l.VideoID != nil {
		str := l.VideoID.String()
		videoID = &str
	}
	return LessonResponse{
		ID:              l.ID.String(),
		SectionID:       l.SectionID.String(),
		CourseID:        l.CourseID.String(),
		Title:           l.Title,
		Type:            string(l.Type),
		OrderIndex:      l.OrderIndex,
		IsFreePreview:   l.IsFreePreview,
		DurationSeconds: l.DurationSeconds,
		IsPublished:     l.IsPublished,
		VideoID:         videoID,
		Content:         l.Content,
		PdfURL:          l.PdfURL,
		CreatedAt:       l.CreatedAt,
		UpdatedAt:       l.UpdatedAt,
	}
}

// NewCourseListResponse maps a slice of courses to flat projections.
func NewCourseListResponse(courses []*entity.Course) []CourseResponse {
	out := make([]CourseResponse, len(courses))
	for i, c := range courses {
		out[i] = NewCourseResponse(c)
	}
	return out
}

// NewSectionListResponse maps a slice of sections to projections.
func NewSectionListResponse(sections []*entity.CourseSection) []SectionResponse {
	out := make([]SectionResponse, len(sections))
	for i, s := range sections {
		out[i] = NewSectionResponse(s)
	}
	return out
}

// NewLessonListResponse maps a slice of lessons to projections.
func NewLessonListResponse(lessons []*entity.Lesson) []LessonResponse {
	out := make([]LessonResponse, len(lessons))
	for i, l := range lessons {
		out[i] = NewLessonResponse(l)
	}
	return out
}

// NewCourseDetailResponse maps the nested course view to its wire shape.
func NewCourseDetailResponse(content *repository.CourseWithContent) CourseDetailResponse {
	sections := make([]SectionDetailResponse, len(content.Sections))
	for i, sec := range content.Sections {
		sections[i] = SectionDetailResponse{
			SectionResponse: NewSectionResponse(sec.Section),
			Lessons:         NewLessonListResponse(sec.Lessons),
		}
	}
	return CourseDetailResponse{
		CourseResponse: NewCourseResponse(content.Course),
		Sections:       sections,
	}
}

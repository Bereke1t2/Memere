package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/exam"
)

// CreateExamRequest is the body of POST /courses/:id/exams. The course id comes
// from the path.
type CreateExamRequest struct {
	Title           string  `json:"title"`
	Subject         string  `json:"subject"`
	Grade           int     `json:"grade"`
	DurationMinutes int     `json:"duration_minutes"`
	PassMarks       int     `json:"pass_marks"`
	Instructions    *string `json:"instructions"`
}

// AddExamQuestionRequest is the body of POST /exams/:id/questions. It references
// an existing question by id and assigns its per-exam marks and order.
type AddExamQuestionRequest struct {
	QuestionID uuid.UUID `json:"question_id"`
	Marks      int       `json:"marks"`
	OrderIndex int       `json:"order_index"`
}

// SubmitExamRequest is the body of POST /exam-attempts/:id/submit.
type SubmitExamRequest struct {
	Answers map[string]any `json:"answers"`
}

// ExamResponse is the authoring-side view of an exam.
type ExamResponse struct {
	ID              uuid.UUID  `json:"id"`
	CourseID        *uuid.UUID `json:"course_id"`
	Title           string     `json:"title"`
	Subject         string     `json:"subject"`
	Grade           int        `json:"grade"`
	DurationMinutes int        `json:"duration_minutes"`
	TotalMarks      int        `json:"total_marks"`
	PassMarks       int        `json:"pass_marks"`
	Instructions    *string    `json:"instructions"`
	IsPublished     bool       `json:"is_published"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// NewExamResponse maps an exam entity to its response.
func NewExamResponse(e *entity.Exam) ExamResponse {
	return ExamResponse{
		ID:              e.ID,
		CourseID:        e.CourseID,
		Title:           e.Title,
		Subject:         e.Subject,
		Grade:           e.Grade,
		DurationMinutes: e.DurationMinutes,
		TotalMarks:      e.TotalMarks,
		PassMarks:       e.PassMarks,
		Instructions:    e.Instructions,
		IsPublished:     e.IsPublished,
		CreatedAt:       e.CreatedAt,
		UpdatedAt:       e.UpdatedAt,
	}
}

// NewExamListResponse maps a page of exams (used for GET /mock-exams).
func NewExamListResponse(exams []*entity.Exam) []ExamResponse {
	out := make([]ExamResponse, 0, len(exams))
	for _, e := range exams {
		out = append(out, NewExamResponse(e))
	}
	return out
}

// ExamQuestionClientResponse is one exam question shown during an attempt: prompt,
// options and per-exam marks — no answer key, no explanation.
type ExamQuestionClientResponse struct {
	QuestionID uuid.UUID              `json:"question_id"`
	Text       string                 `json:"text"`
	Type       string                 `json:"type"`
	Marks      int                    `json:"marks"`
	Subject    *string                `json:"subject"`
	Topic      *string                `json:"topic"`
	Answers    []AnswerClientResponse `json:"answers"`
}

// StartExamResponse is the live exam attempt handed to a student.
type StartExamResponse struct {
	AttemptID        uuid.UUID                    `json:"attempt_id"`
	ExamID           uuid.UUID                    `json:"exam_id"`
	Status           string                       `json:"status"`
	StartedAt        time.Time                    `json:"started_at"`
	ExpiresAt        *time.Time                   `json:"expires_at"`
	RemainingSeconds *int                         `json:"remaining_seconds"`
	TotalMarks       int                          `json:"total_marks"`
	Questions        []ExamQuestionClientResponse `json:"questions"`
}

// NewStartExamResponse maps the live exam attempt view to its response.
func NewStartExamResponse(v *exam.ExamAttemptClientView) StartExamResponse {
	questions := make([]ExamQuestionClientResponse, 0, len(v.Questions))
	for _, q := range v.Questions {
		answers := make([]AnswerClientResponse, 0, len(q.Answers))
		for _, a := range q.Answers {
			answers = append(answers, AnswerClientResponse{ID: a.ID, Text: a.Text, OrderIndex: a.OrderIndex})
		}
		questions = append(questions, ExamQuestionClientResponse{
			QuestionID: q.QuestionID,
			Text:       q.Text,
			Type:       string(q.Type),
			Marks:      q.Marks,
			Subject:    q.Subject,
			Topic:      q.Topic,
			Answers:    answers,
		})
	}
	return StartExamResponse{
		AttemptID:        v.AttemptID,
		ExamID:           v.ExamID,
		Status:           string(v.Status),
		StartedAt:        v.StartedAt,
		ExpiresAt:        v.ExpiresAt,
		RemainingSeconds: v.RemainingSeconds,
		TotalMarks:       v.TotalMarks,
		Questions:        questions,
	}
}

// ExamQuestionFeedbackResponse is the per-question outcome in a graded exam
// result. Correct answers and explanation appear here ONLY (post-submission).
type ExamQuestionFeedbackResponse struct {
	QuestionID       uuid.UUID   `json:"question_id"`
	Correct          bool        `json:"correct"`
	MarksAwarded     int         `json:"marks_awarded"`
	MarksPossible    int         `json:"marks_possible"`
	SelectedAnswers  []uuid.UUID `json:"selected_answers"`
	CorrectAnswerIDs []uuid.UUID `json:"correct_answer_ids"`
	Explanation      *string     `json:"explanation"`
}

// ExamResultResponse is the post-submission exam result. passed is decided
// against the exam's absolute pass_marks (§4.2.6).
type ExamResultResponse struct {
	AttemptID        uuid.UUID                       `json:"attempt_id"`
	ExamID           uuid.UUID                       `json:"exam_id"`
	Status           string                          `json:"status"`
	Score            float64                         `json:"score"`
	TotalMarks       int                             `json:"total_marks"`
	Percentage       float64                         `json:"percentage"`
	PassMarks        int                             `json:"pass_marks"`
	Passed           bool                            `json:"passed"`
	SubmittedAt      *time.Time                      `json:"submitted_at"`
	Feedback         []ExamQuestionFeedbackResponse  `json:"feedback"`
	SubjectBreakdown map[string]SubjectScoreResponse `json:"subject_breakdown"`
}

// NewExamResultResponse maps a graded exam result to its response.
func NewExamResultResponse(r *exam.ExamResult) ExamResultResponse {
	feedback := make([]ExamQuestionFeedbackResponse, 0, len(r.Feedback))
	for _, f := range r.Feedback {
		feedback = append(feedback, ExamQuestionFeedbackResponse{
			QuestionID:       f.QuestionID,
			Correct:          f.Correct,
			MarksAwarded:     f.MarksAwarded,
			MarksPossible:    f.MarksPossible,
			SelectedAnswers:  f.SelectedAnswers,
			CorrectAnswerIDs: f.CorrectAnswerIDs,
			Explanation:      f.Explanation,
		})
	}
	breakdown := make(map[string]SubjectScoreResponse, len(r.SubjectBreakdown))
	for k, v := range r.SubjectBreakdown {
		breakdown[k] = SubjectScoreResponse{Earned: v.Earned, Possible: v.Possible}
	}
	return ExamResultResponse{
		AttemptID:        r.AttemptID,
		ExamID:           r.ExamID,
		Status:           string(r.Status),
		Score:            r.Score,
		TotalMarks:       r.TotalMarks,
		Percentage:       r.Percentage,
		PassMarks:        r.PassMarks,
		Passed:           r.Passed,
		SubmittedAt:      r.SubmittedAt,
		Feedback:         feedback,
		SubjectBreakdown: breakdown,
	}
}

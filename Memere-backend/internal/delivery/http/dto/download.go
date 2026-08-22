package dto

import (
	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/exam"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/quiz"
)

// ─────────────────────────────────────────────────────────────────────────────
// OFFLINE DOWNLOAD PAYLOADS — the ONE sanctioned exception to Non-Negotiable #1
// ("correct answers never sent to the client"). These responses INTENTIONALLY
// serialize is_correct, the correct short-answer text, and each question's
// explanation so a downloaded quiz/exam can be graded on-device while offline.
//
// This is safe and narrowly scoped because:
//   - It is reachable ONLY via GET /quizzes/:id/download and
//     GET /mock-exams/:id/download, gated by the SAME access check as taking the
//     quiz/exam (assertCourseAccess / assertExamAccess). Re-enabling the paywall
//     (access.DisableEnrollmentCheck=false), or leaving an exam unpublished,
//     withholds these keys automatically.
//   - The normal online-play path (QuizClientResponse / ExamQuestionClientResponse,
//     built from GetQuestionsForClient) stays answer-free and server-graded.
//   - On reconnect the raw answers are re-submitted for authoritative server-side
//     grading, so on-device scores are provisional only.
//
// Do NOT reuse these types on any non-download route.
// ─────────────────────────────────────────────────────────────────────────────

// DownloadAnswerResponse is one answer option in a download payload, INCLUDING
// whether it is correct.
type DownloadAnswerResponse struct {
	ID         uuid.UUID `json:"id"`
	Text       string    `json:"text"`
	OrderIndex int       `json:"order_index"`
	IsCorrect  bool      `json:"is_correct"`
}

// DownloadQuestionResponse is one question with its full answer key and
// explanation. points carries the per-question weight (quiz points / exam marks).
type DownloadQuestionResponse struct {
	ID          uuid.UUID                `json:"id"`
	Text        string                   `json:"text"`
	Type        string                   `json:"type"`
	Points      int                      `json:"points"`
	OrderIndex  int                      `json:"order_index"`
	Subject     *string                  `json:"subject"`
	Topic       *string                  `json:"topic"`
	Explanation *string                  `json:"explanation"`
	Answers     []DownloadAnswerResponse `json:"answers"`
}

// QuizDownloadResponse is the full quiz WITH answer keys for offline grading.
type QuizDownloadResponse struct {
	ID               uuid.UUID                  `json:"id"`
	CourseID         uuid.UUID                  `json:"course_id"`
	Title            string                     `json:"title"`
	TimeLimitSeconds *int                       `json:"time_limit_seconds"`
	PassPercentage   float64                    `json:"pass_percentage"`
	MaxAttempts      *int                       `json:"max_attempts"`
	Questions        []DownloadQuestionResponse `json:"questions"`
}

// NewQuizDownloadResponse maps the quiz download usecase view to its response,
// intentionally including the answer key (see the file header).
func NewQuizDownloadResponse(d *quiz.QuizDownload) QuizDownloadResponse {
	questions := make([]DownloadQuestionResponse, 0, len(d.Questions))
	for _, q := range d.Questions {
		answers := make([]DownloadAnswerResponse, 0, len(q.Answers))
		for _, a := range q.Answers {
			answers = append(answers, DownloadAnswerResponse{
				ID:         a.ID,
				Text:       a.Text,
				OrderIndex: a.OrderIndex,
				IsCorrect:  a.IsCorrect,
			})
		}
		questions = append(questions, DownloadQuestionResponse{
			ID:          q.ID,
			Text:        q.Text,
			Type:        string(q.Type),
			Points:      q.Points,
			OrderIndex:  q.OrderIndex,
			Subject:     q.Subject,
			Topic:       q.Topic,
			Explanation: q.Explanation,
			Answers:     answers,
		})
	}
	return QuizDownloadResponse{
		ID:               d.Quiz.ID,
		CourseID:         d.Quiz.CourseID,
		Title:            d.Quiz.Title,
		TimeLimitSeconds: d.Quiz.TimeLimitSeconds,
		PassPercentage:   d.Quiz.PassPercentage,
		MaxAttempts:      d.Quiz.MaxAttempts,
		Questions:        questions,
	}
}

// ExamDownloadResponse is the full exam WITH answer keys for offline grading.
// duration_minutes lets the device run a display-only timer; pass_marks is the
// absolute pass threshold used by the on-device grader.
type ExamDownloadResponse struct {
	ID              uuid.UUID                  `json:"id"`
	CourseID        *uuid.UUID                 `json:"course_id"`
	Title           string                     `json:"title"`
	Subject         string                     `json:"subject"`
	Grade           int                        `json:"grade"`
	DurationMinutes int                        `json:"duration_minutes"`
	TotalMarks      int                        `json:"total_marks"`
	PassMarks       int                        `json:"pass_marks"`
	Instructions    *string                    `json:"instructions"`
	Questions       []DownloadQuestionResponse `json:"questions"`
}

// NewExamDownloadResponse maps the exam download usecase view to its response,
// intentionally including the answer key (see the file header). Exam marks are
// carried in the shared points field.
func NewExamDownloadResponse(d *exam.ExamDownload) ExamDownloadResponse {
	questions := make([]DownloadQuestionResponse, 0, len(d.Questions))
	for _, q := range d.Questions {
		answers := make([]DownloadAnswerResponse, 0, len(q.Answers))
		for _, a := range q.Answers {
			answers = append(answers, DownloadAnswerResponse{
				ID:         a.ID,
				Text:       a.Text,
				OrderIndex: a.OrderIndex,
				IsCorrect:  a.IsCorrect,
			})
		}
		questions = append(questions, DownloadQuestionResponse{
			ID:          q.ID,
			Text:        q.Text,
			Type:        string(q.Type),
			Points:      q.Marks,
			OrderIndex:  q.OrderIndex,
			Subject:     q.Subject,
			Topic:       q.Topic,
			Explanation: q.Explanation,
			Answers:     answers,
		})
	}
	return ExamDownloadResponse{
		ID:              d.Exam.ID,
		CourseID:        d.Exam.CourseID,
		Title:           d.Exam.Title,
		Subject:         d.Exam.Subject,
		Grade:           d.Exam.Grade,
		DurationMinutes: d.Exam.DurationMinutes,
		TotalMarks:      d.Exam.TotalMarks,
		PassMarks:       d.Exam.PassMarks,
		Instructions:    d.Exam.Instructions,
		Questions:       questions,
	}
}

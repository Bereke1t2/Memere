package dto

import (
	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/exam"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/quiz"
)

// The response types in this file are the ONE place a question's is_correct flag
// is serialized outward, because an author cannot edit an answer key they cannot
// see. This is not a hole in Non-Negotiable #1, which governs what reaches a
// STUDENT: every route returning these types is registered behind requireAuth +
// teacherOrAdmin, and the usecase additionally proves course ownership before
// building the view. The student-facing types (QuestionClientResponse,
// AnswerClientResponse) remain structurally incapable of carrying the key, and
// MUST stay that way. Do not reuse these types on any student route.

// AuthoredAnswerResponse is one option as its author sees it, answer key included.
type AuthoredAnswerResponse struct {
	ID         uuid.UUID `json:"id"`
	Text       string    `json:"text"`
	IsCorrect  bool      `json:"is_correct"`
	OrderIndex int       `json:"order_index"`
}

// AuthoredQuestionResponse is one question as its author sees it. Marks mirrors
// Points for exam questions, where the per-exam link row owns the value.
type AuthoredQuestionResponse struct {
	ID          uuid.UUID                `json:"id"`
	Text        string                   `json:"text"`
	Type        string                   `json:"type"`
	Points      int                      `json:"points"`
	Marks       int                      `json:"marks"`
	Explanation *string                  `json:"explanation"`
	OrderIndex  int                      `json:"order_index"`
	Subject     *string                  `json:"subject"`
	Topic       *string                  `json:"topic"`
	Answers     []AuthoredAnswerResponse `json:"answers"`
}

// AuthoredQuizResponse is a quiz plus its full question tree for the editor.
type AuthoredQuizResponse struct {
	Quiz      QuizResponse               `json:"quiz"`
	Questions []AuthoredQuestionResponse `json:"questions"`
}

// AuthoredExamResponse is an exam plus its full question set for the editor.
type AuthoredExamResponse struct {
	Exam      ExamResponse               `json:"exam"`
	Questions []AuthoredQuestionResponse `json:"questions"`
}

// UpdateQuestionRequest is the body of PUT /questions/:id. Nil scalar fields are
// unchanged. A nil answers field leaves the existing option set alone; a present
// one replaces it wholesale.
type UpdateQuestionRequest struct {
	Text        *string             `json:"text"`
	Type        *string             `json:"type"`
	Points      *int                `json:"points"`
	Explanation *string             `json:"explanation"`
	OrderIndex  *int                `json:"order_index"`
	Subject     *string             `json:"subject"`
	Topic       *string             `json:"topic"`
	Answers     *[]AddAnswerRequest `json:"answers"`
}

// UpdateExamQuestionRequest is the body of PUT /exams/:id/questions/:questionId.
// Marks is the per-exam mark value; changing it re-derives the exam's total.
type UpdateExamQuestionRequest struct {
	Text        *string             `json:"text"`
	Type        *string             `json:"type"`
	Marks       *int                `json:"marks"`
	Explanation *string             `json:"explanation"`
	OrderIndex  *int                `json:"order_index"`
	Subject     *string             `json:"subject"`
	Topic       *string             `json:"topic"`
	Answers     *[]AddAnswerRequest `json:"answers"`
}

// ReorderQuestionsRequest is the body of the reorder routes: the complete list
// of question ids in their intended display order.
type ReorderQuestionsRequest struct {
	QuestionIDs []uuid.UUID `json:"question_ids"`
}

// NewAuthoredQuizResponse maps the usecase authoring view to its response.
func NewAuthoredQuizResponse(v *quiz.AuthoredQuiz) AuthoredQuizResponse {
	questions := make([]AuthoredQuestionResponse, len(v.Questions))
	for i, q := range v.Questions {
		questions[i] = AuthoredQuestionResponse{
			ID:          q.ID,
			Text:        q.Text,
			Type:        string(q.Type),
			Points:      q.Points,
			Marks:       q.Points,
			Explanation: q.Explanation,
			OrderIndex:  q.OrderIndex,
			Subject:     q.Subject,
			Topic:       q.Topic,
			Answers:     newAuthoredAnswers(q.Answers),
		}
	}
	return AuthoredQuizResponse{Quiz: NewQuizResponse(v.Quiz), Questions: questions}
}

// NewAuthoredQuizQuestionResponse maps one updated quiz question to its response.
func NewAuthoredQuizQuestionResponse(q *quiz.AuthoredQuestion) AuthoredQuestionResponse {
	return AuthoredQuestionResponse{
		ID:          q.ID,
		Text:        q.Text,
		Type:        string(q.Type),
		Points:      q.Points,
		Marks:       q.Points,
		Explanation: q.Explanation,
		OrderIndex:  q.OrderIndex,
		Subject:     q.Subject,
		Topic:       q.Topic,
		Answers:     newAuthoredAnswers(q.Answers),
	}
}

// NewAuthoredExamResponse maps the usecase authoring view to its response.
func NewAuthoredExamResponse(v *exam.AuthoredExam) AuthoredExamResponse {
	questions := make([]AuthoredQuestionResponse, len(v.Questions))
	for i, q := range v.Questions {
		questions[i] = newAuthoredExamQuestion(q)
	}
	return AuthoredExamResponse{Exam: NewExamResponse(v.Exam), Questions: questions}
}

// NewAuthoredExamQuestionResponse maps one updated exam question to its response.
func NewAuthoredExamQuestionResponse(q *exam.AuthoredQuestion) AuthoredQuestionResponse {
	return newAuthoredExamQuestion(*q)
}

func newAuthoredExamQuestion(q exam.AuthoredQuestion) AuthoredQuestionResponse {
	return AuthoredQuestionResponse{
		ID:          q.ID,
		Text:        q.Text,
		Type:        string(q.Type),
		Points:      q.Marks,
		Marks:       q.Marks,
		Explanation: q.Explanation,
		OrderIndex:  q.OrderIndex,
		Subject:     q.Subject,
		Topic:       q.Topic,
		Answers: func() []AuthoredAnswerResponse {
			out := make([]AuthoredAnswerResponse, len(q.Answers))
			for i, a := range q.Answers {
				out[i] = AuthoredAnswerResponse{ID: a.ID, Text: a.Text, IsCorrect: a.IsCorrect, OrderIndex: a.OrderIndex}
			}
			return out
		}(),
	}
}

func newAuthoredAnswers(in []quiz.AuthoredAnswer) []AuthoredAnswerResponse {
	out := make([]AuthoredAnswerResponse, len(in))
	for i, a := range in {
		out[i] = AuthoredAnswerResponse{ID: a.ID, Text: a.Text, IsCorrect: a.IsCorrect, OrderIndex: a.OrderIndex}
	}
	return out
}

package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/quiz"
)

// CreateQuizRequest is the body of POST /courses/:id/quizzes. The course id comes
// from the path, not the body.
type CreateQuizRequest struct {
	LessonID           *uuid.UUID `json:"lesson_id"`
	Title              string     `json:"title"`
	TimeLimitSeconds   *int       `json:"time_limit_seconds"`
	PassPercentage     float64    `json:"pass_percentage"`
	RandomizeQuestions bool       `json:"randomize_questions"`
	MaxAttempts        *int       `json:"max_attempts"`
}

// UpdateQuizRequest is the body of PUT /quizzes/:id. Nil fields are unchanged.
type UpdateQuizRequest struct {
	Title              *string  `json:"title"`
	TimeLimitSeconds   *int     `json:"time_limit_seconds"`
	PassPercentage     *float64 `json:"pass_percentage"`
	RandomizeQuestions *bool    `json:"randomize_questions"`
	MaxAttempts        *int     `json:"max_attempts"`
}

// AddAnswerRequest is one answer option in AddQuestionRequest. is_correct is
// accepted on the AUTHORING path only (teacher input); it is never echoed back on
// any client-facing response.
type AddAnswerRequest struct {
	Text       string `json:"text"`
	IsCorrect  bool   `json:"is_correct"`
	OrderIndex int    `json:"order_index"`
}

// AddQuestionRequest is the body of POST /quizzes/:id/questions. Answers are
// created inline with the question.
type AddQuestionRequest struct {
	Text        string             `json:"text"`
	Type        string             `json:"type"`
	Points      int                `json:"points"`
	Explanation *string            `json:"explanation"`
	OrderIndex  int                `json:"order_index"`
	Subject     *string            `json:"subject"`
	Topic       *string            `json:"topic"`
	Answers     []AddAnswerRequest `json:"answers"`
}

// SaveProgressRequest is the body of PATCH /quiz-attempts/:id (auto-save). The
// answers map is keyed by question id; values are selected answer id(s) or text.
type SaveProgressRequest struct {
	Answers map[string]any `json:"answers"`
}

// SubmitAttemptRequest is the body of POST /quiz-attempts/:id/submit.
type SubmitAttemptRequest struct {
	Answers map[string]any `json:"answers"`
}

// QuizResponse is the authoring-side view of a quiz (returned on create/update).
type QuizResponse struct {
	ID                 uuid.UUID  `json:"id"`
	CourseID           uuid.UUID  `json:"course_id"`
	LessonID           *uuid.UUID `json:"lesson_id"`
	Title              string     `json:"title"`
	TimeLimitSeconds   *int       `json:"time_limit_seconds"`
	PassPercentage     float64    `json:"pass_percentage"`
	RandomizeQuestions bool       `json:"randomize_questions"`
	MaxAttempts        *int       `json:"max_attempts"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

// NewQuizResponse maps a quiz entity to its authoring response.
func NewQuizResponse(q *entity.Quiz) QuizResponse {
	return QuizResponse{
		ID:                 q.ID,
		CourseID:           q.CourseID,
		LessonID:           q.LessonID,
		Title:              q.Title,
		TimeLimitSeconds:   q.TimeLimitSeconds,
		PassPercentage:     q.PassPercentage,
		RandomizeQuestions: q.RandomizeQuestions,
		MaxAttempts:        q.MaxAttempts,
		CreatedAt:          q.CreatedAt,
		UpdatedAt:          q.UpdatedAt,
	}
}

// QuestionResponse is the authoring-side view of a question (returned on add).
// It deliberately omits answer options' is_correct flags — only the question
// shell is echoed back.
type QuestionResponse struct {
	ID          uuid.UUID `json:"id"`
	QuizID      uuid.UUID `json:"quiz_id"`
	Text        string    `json:"text"`
	Type        string    `json:"type"`
	Points      int       `json:"points"`
	Explanation *string   `json:"explanation"`
	OrderIndex  int       `json:"order_index"`
	Subject     *string   `json:"subject"`
	Topic       *string   `json:"topic"`
}

// NewQuestionResponse maps a question entity to its authoring response.
func NewQuestionResponse(q *entity.Question) QuestionResponse {
	return QuestionResponse{
		ID:          q.ID,
		QuizID:      q.QuizID,
		Text:        q.Text,
		Type:        string(q.Type),
		Points:      q.Points,
		Explanation: q.Explanation,
		OrderIndex:  q.OrderIndex,
		Subject:     q.Subject,
		Topic:       q.Topic,
	}
}

// QuizClientResponse is the student-facing quiz metadata: question COUNT only,
// never the questions. There is no field that could carry an answer key.
type QuizClientResponse struct {
	ID                 uuid.UUID `json:"id"`
	CourseID           uuid.UUID `json:"course_id"`
	Title              string    `json:"title"`
	TimeLimitSeconds   *int      `json:"time_limit_seconds"`
	PassPercentage     float64   `json:"pass_percentage"`
	RandomizeQuestions bool      `json:"randomize_questions"`
	MaxAttempts        *int      `json:"max_attempts"`
	QuestionCount      int       `json:"question_count"`
	AttemptsUsed       int       `json:"attempts_used"`
}

// NewQuizClientResponse maps the usecase client view to its response.
func NewQuizClientResponse(v *quiz.QuizClientView) QuizClientResponse {
	return QuizClientResponse{
		ID:                 v.ID,
		CourseID:           v.CourseID,
		Title:              v.Title,
		TimeLimitSeconds:   v.TimeLimitSeconds,
		PassPercentage:     v.PassPercentage,
		RandomizeQuestions: v.RandomizeQuestions,
		MaxAttempts:        v.MaxAttempts,
		QuestionCount:      v.QuestionCount,
		AttemptsUsed:       v.AttemptsUsed,
	}
}

// AnswerClientResponse is one option shown to a student. By construction it has
// no is_correct field — the answer key cannot be serialized here.
type AnswerClientResponse struct {
	ID         uuid.UUID `json:"id"`
	Text       string    `json:"text"`
	OrderIndex int       `json:"order_index"`
}

// QuestionClientResponse is one question shown during an attempt: no answer key,
// no explanation (revealed only post-submission).
type QuestionClientResponse struct {
	ID      uuid.UUID              `json:"id"`
	Text    string                 `json:"text"`
	Type    string                 `json:"type"`
	Points  int                    `json:"points"`
	Subject *string                `json:"subject"`
	Topic   *string                `json:"topic"`
	Answers []AnswerClientResponse `json:"answers"`
}

// StartAttemptResponse is the live attempt handed to a student: randomized
// questions plus the server-side timing facts. No scores, no answer keys.
type StartAttemptResponse struct {
	AttemptID        uuid.UUID                `json:"attempt_id"`
	QuizID           uuid.UUID                `json:"quiz_id"`
	AttemptNumber    int                      `json:"attempt_number"`
	Status           string                   `json:"status"`
	StartedAt        time.Time                `json:"started_at"`
	ExpiresAt        *time.Time               `json:"expires_at"`
	RemainingSeconds *int                     `json:"remaining_seconds"`
	Questions        []QuestionClientResponse `json:"questions"`
}

// NewStartAttemptResponse maps the live attempt client view to its response.
func NewStartAttemptResponse(v *quiz.AttemptClientView) StartAttemptResponse {
	return StartAttemptResponse{
		AttemptID:        v.AttemptID,
		QuizID:           v.QuizID,
		AttemptNumber:    v.AttemptNumber,
		Status:           string(v.Status),
		StartedAt:        v.StartedAt,
		ExpiresAt:        v.ExpiresAt,
		RemainingSeconds: v.RemainingSeconds,
		Questions:        newQuizQuestionClientResponses(v.Questions),
	}
}

func newQuizQuestionClientResponses(qs []quiz.QuestionClientView) []QuestionClientResponse {
	out := make([]QuestionClientResponse, 0, len(qs))
	for _, q := range qs {
		answers := make([]AnswerClientResponse, 0, len(q.Answers))
		for _, a := range q.Answers {
			answers = append(answers, AnswerClientResponse{ID: a.ID, Text: a.Text, OrderIndex: a.OrderIndex})
		}
		out = append(out, QuestionClientResponse{
			ID:      q.ID,
			Text:    q.Text,
			Type:    string(q.Type),
			Points:  q.Points,
			Subject: q.Subject,
			Topic:   q.Topic,
			Answers: answers,
		})
	}
	return out
}

// QuestionFeedbackResponse is the per-question outcome in a graded quiz result.
// CorrectAnswerIDs and explanation are present here ONLY — this is the
// post-submission reveal (§4.2.5).
type QuestionFeedbackResponse struct {
	QuestionID       uuid.UUID   `json:"question_id"`
	Correct          bool        `json:"correct"`
	PointsAwarded    int         `json:"points_awarded"`
	PointsPossible   int         `json:"points_possible"`
	SelectedAnswers  []uuid.UUID `json:"selected_answers"`
	CorrectAnswerIDs []uuid.UUID `json:"correct_answer_ids"`
	Explanation      *string     `json:"explanation"`
}

// SubjectScoreResponse is a per-subject earned/possible tally.
type SubjectScoreResponse struct {
	Earned   int `json:"earned"`
	Possible int `json:"possible"`
}

// AttemptResultResponse is the post-submission quiz result with feedback.
type AttemptResultResponse struct {
	AttemptID        uuid.UUID                       `json:"attempt_id"`
	QuizID           uuid.UUID                       `json:"quiz_id"`
	AttemptNumber    int                             `json:"attempt_number"`
	Status           string                          `json:"status"`
	Score            float64                         `json:"score"`
	TotalPoints      int                             `json:"total_points"`
	Percentage       float64                         `json:"percentage"`
	Passed           bool                            `json:"passed"`
	SubmittedAt      *time.Time                      `json:"submitted_at"`
	Feedback         []QuestionFeedbackResponse      `json:"feedback"`
	SubjectBreakdown map[string]SubjectScoreResponse `json:"subject_breakdown"`
}

// NewAttemptResultResponse maps a graded quiz result to its response.
func NewAttemptResultResponse(r *quiz.AttemptResult) AttemptResultResponse {
	feedback := make([]QuestionFeedbackResponse, 0, len(r.Feedback))
	for _, f := range r.Feedback {
		feedback = append(feedback, QuestionFeedbackResponse{
			QuestionID:       f.QuestionID,
			Correct:          f.Correct,
			PointsAwarded:    f.PointsAwarded,
			PointsPossible:   f.PointsPossible,
			SelectedAnswers:  f.SelectedAnswers,
			CorrectAnswerIDs: f.CorrectAnswerIDs,
			Explanation:      f.Explanation,
		})
	}
	breakdown := make(map[string]SubjectScoreResponse, len(r.SubjectBreakdown))
	for k, v := range r.SubjectBreakdown {
		breakdown[k] = SubjectScoreResponse{Earned: v.Earned, Possible: v.Possible}
	}
	return AttemptResultResponse{
		AttemptID:        r.AttemptID,
		QuizID:           r.QuizID,
		AttemptNumber:    r.AttemptNumber,
		Status:           string(r.Status),
		Score:            r.Score,
		TotalPoints:      r.TotalPoints,
		Percentage:       r.Percentage,
		Passed:           r.Passed,
		SubmittedAt:      r.SubmittedAt,
		Feedback:         feedback,
		SubjectBreakdown: breakdown,
	}
}

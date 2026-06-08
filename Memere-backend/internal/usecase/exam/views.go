package exam

import (
	"time"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
)

// AnswerClientView is one option as shown to a student: no is_correct field
// exists on this type, so the answer key cannot be serialized to a client.
type AnswerClientView struct {
	ID         uuid.UUID
	Text       string
	OrderIndex int
}

// QuestionClientView is one exam question as shown to a student during an
// attempt: prompt, options, and the per-exam marks — no answer key, no
// explanation (revealed only in the post-submission result).
type QuestionClientView struct {
	QuestionID uuid.UUID
	Text       string
	Type       entity.QuestionType
	Marks      int
	Subject    *string
	Topic      *string
	Answers    []AnswerClientView
}

// ExamAttemptClientView is the live attempt handed to a student: the questions in
// the attempt's order plus the server-side timing facts and the exam's total
// marks. It never contains scores or answer keys.
type ExamAttemptClientView struct {
	AttemptID        uuid.UUID
	ExamID           uuid.UUID
	Status           entity.AttemptStatus
	StartedAt        time.Time
	ExpiresAt        *time.Time
	RemainingSeconds *int
	TotalMarks       int
	Questions        []QuestionClientView
}

// QuestionFeedback is the per-question outcome in a graded result. CorrectAnswerIDs
// is populated only after grading (post-submission reveal per §4.2.5).
type QuestionFeedback struct {
	QuestionID       uuid.UUID
	Correct          bool
	MarksAwarded     int
	MarksPossible    int
	SelectedAnswers  []uuid.UUID
	CorrectAnswerIDs []uuid.UUID
	Explanation      *string
}

// SubjectScore is the per-subject (or per-topic) tally for §9.3 analytics.
type SubjectScore struct {
	Earned   int
	Possible int
}

// ExamResult is the post-submission result with feedback. Score/percentage are
// the server-computed authority; passed is decided against the exam's absolute
// pass_marks (not a percentage), per §4.2.6.
type ExamResult struct {
	AttemptID        uuid.UUID
	ExamID           uuid.UUID
	Status           entity.AttemptStatus
	Score            float64
	TotalMarks       int
	Percentage       float64
	PassMarks        int
	Passed           bool
	SubmittedAt      *time.Time
	Feedback         []QuestionFeedback
	SubjectBreakdown map[string]SubjectScore
}

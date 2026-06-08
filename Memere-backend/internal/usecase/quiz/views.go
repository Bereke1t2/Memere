package quiz

import (
	"time"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
)

// QuizClientView is the quiz metadata a student sees before starting. It carries
// the question COUNT only — never the questions themselves, which arrive (answer-
// key-free) on StartAttempt.
type QuizClientView struct {
	ID                 uuid.UUID
	CourseID           uuid.UUID
	Title              string
	TimeLimitSeconds   *int
	PassPercentage     float64
	RandomizeQuestions bool
	MaxAttempts        *int
	QuestionCount      int
	AttemptsUsed       int
}

// AnswerClientView is one answer option as shown to a student: no is_correct
// field exists on this type, so the answer key cannot be serialized to a client.
type AnswerClientView struct {
	ID         uuid.UUID
	Text       string
	OrderIndex int
}

// QuestionClientView is one question as shown to a student during an attempt:
// the prompt and its options, with no answer key and no explanation (the
// explanation is revealed only in the post-submission result).
type QuestionClientView struct {
	ID      uuid.UUID
	Text    string
	Type    entity.QuestionType
	Points  int
	Subject *string
	Topic   *string
	Answers []AnswerClientView
}

// AttemptClientView is the live attempt handed to a student: the questions in
// the attempt's randomized order plus the server-side timing facts. It never
// contains scores or answer keys.
type AttemptClientView struct {
	AttemptID        uuid.UUID
	QuizID           uuid.UUID
	AttemptNumber    int
	Status           entity.AttemptStatus
	StartedAt        time.Time
	ExpiresAt        *time.Time
	RemainingSeconds *int
	Questions        []QuestionClientView
}

// QuestionFeedback is the per-question outcome in a graded result. CorrectAnswerIDs
// is populated only after grading (post-submission reveal per §4.2.5).
type QuestionFeedback struct {
	QuestionID       uuid.UUID
	Correct          bool
	PointsAwarded    int
	PointsPossible   int
	SelectedAnswers  []uuid.UUID
	CorrectAnswerIDs []uuid.UUID
	Explanation      *string
}

// AttemptResult is the post-submission result with feedback. Score/percentage
// are the server-computed authority; the client-reported score is never used.
type AttemptResult struct {
	AttemptID        uuid.UUID
	QuizID           uuid.UUID
	AttemptNumber    int
	Status           entity.AttemptStatus
	Score            float64
	TotalPoints      int
	Percentage       float64
	Passed           bool
	SubmittedAt      *time.Time
	Feedback         []QuestionFeedback
	SubjectBreakdown map[string]SubjectScore
}

// SubjectScore is the per-subject (or per-topic) tally for §9.3 analytics.
type SubjectScore struct {
	Earned   int
	Possible int
}

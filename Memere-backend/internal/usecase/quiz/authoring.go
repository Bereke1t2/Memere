package quiz

import (
	"context"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/validator"
)

const (
	maxQuizTitleLen   = 200
	maxQuestionTexLen = 5000
	maxExplanationLen = 5000
)

// CreateQuizInput is the teacher-facing quiz create request after decoding.
type CreateQuizInput struct {
	CourseID           uuid.UUID
	LessonID           *uuid.UUID
	Title              string
	TimeLimitSeconds   *int
	PassPercentage     float64
	RandomizeQuestions bool
	MaxAttempts        *int
}

// UpdateQuizInput carries partial quiz updates; nil fields are unchanged.
type UpdateQuizInput struct {
	Title              *string
	TimeLimitSeconds   *int
	PassPercentage     *float64
	RandomizeQuestions *bool
	MaxAttempts        *int
}

// AnswerInput is one answer option for a question being authored.
type AnswerInput struct {
	Text       string
	IsCorrect  bool
	OrderIndex int
}

// QuestionInput is a question plus its options for authoring.
type QuestionInput struct {
	Text        string
	Type        entity.QuestionType
	Points      int
	Explanation *string
	OrderIndex  int
	Subject     *string
	Topic       *string
	Answers     []AnswerInput
}

// CreateQuiz creates a quiz under a course the actor owns.
func (s *Service) CreateQuiz(ctx context.Context, actor *Actor, in CreateQuizInput) (*entity.Quiz, error) {
	if _, err := s.loadOwnedCourse(ctx, actor, in.CourseID); err != nil {
		return nil, err
	}
	if details := validateQuizInput(in.Title, in.PassPercentage, in.TimeLimitSeconds, in.MaxAttempts); len(details) > 0 {
		return nil, apperror.Validation(details, nil)
	}

	quiz := &entity.Quiz{
		LessonID:           in.LessonID,
		CourseID:           in.CourseID,
		Title:              in.Title,
		TimeLimitSeconds:   in.TimeLimitSeconds,
		PassPercentage:     in.PassPercentage,
		RandomizeQuestions: in.RandomizeQuestions,
		MaxAttempts:        in.MaxAttempts,
	}
	if err := s.quizzes.Create(ctx, quiz); err != nil {
		return nil, err
	}
	return quiz, nil
}

// AddQuestion adds a question and its options to a quiz the actor owns. The
// question and its answers are written in one transaction so the answer key is
// never left half-written. Choice questions require at least one correct option.
func (s *Service) AddQuestion(ctx context.Context, actor *Actor, quizID uuid.UUID, in QuestionInput) (*entity.Question, error) {
	if _, err := s.loadOwnedQuiz(ctx, actor, quizID); err != nil {
		return nil, err
	}
	if details := validateQuestionInput(in); len(details) > 0 {
		return nil, apperror.Validation(details, nil)
	}

	question := &entity.Question{
		QuizID:      quizID,
		Text:        in.Text,
		Type:        in.Type,
		Points:      in.Points,
		Explanation: in.Explanation,
		OrderIndex:  in.OrderIndex,
		Subject:     in.Subject,
		Topic:       in.Topic,
	}
	answers := make([]*entity.Answer, len(in.Answers))
	for i, a := range in.Answers {
		answers[i] = &entity.Answer{Text: a.Text, IsCorrect: a.IsCorrect, OrderIndex: a.OrderIndex}
	}

	err := s.tx.WithinTx(ctx, func(ctx context.Context) error {
		return s.questions.Create(ctx, question, answers)
	})
	if err != nil {
		return nil, err
	}
	return question, nil
}

// UpdateQuiz applies partial updates to a quiz the actor owns.
func (s *Service) UpdateQuiz(ctx context.Context, actor *Actor, quizID uuid.UUID, in UpdateQuizInput) (*entity.Quiz, error) {
	quiz, err := s.loadOwnedQuiz(ctx, actor, quizID)
	if err != nil {
		return nil, err
	}

	if in.Title != nil {
		quiz.Title = *in.Title
	}
	if in.TimeLimitSeconds != nil {
		quiz.TimeLimitSeconds = in.TimeLimitSeconds
	}
	if in.PassPercentage != nil {
		quiz.PassPercentage = *in.PassPercentage
	}
	if in.RandomizeQuestions != nil {
		quiz.RandomizeQuestions = *in.RandomizeQuestions
	}
	if in.MaxAttempts != nil {
		quiz.MaxAttempts = in.MaxAttempts
	}

	if details := validateQuizInput(quiz.Title, quiz.PassPercentage, quiz.TimeLimitSeconds, quiz.MaxAttempts); len(details) > 0 {
		return nil, apperror.Validation(details, nil)
	}
	if err := s.quizzes.Update(ctx, quiz); err != nil {
		return nil, err
	}
	return quiz, nil
}

// DeleteQuiz soft-deletes a quiz the actor owns.
func (s *Service) DeleteQuiz(ctx context.Context, actor *Actor, quizID uuid.UUID) error {
	if _, err := s.loadOwnedQuiz(ctx, actor, quizID); err != nil {
		return err
	}
	return s.quizzes.Delete(ctx, quizID)
}

// loadOwnedCourse loads a course and asserts the actor may author under it.
func (s *Service) loadOwnedCourse(ctx context.Context, actor *Actor, courseID uuid.UUID) (*entity.Course, error) {
	if !actor.isTeacherOrAdmin() {
		return nil, apperror.Forbidden("only teachers may author quizzes", nil)
	}
	course, err := s.courses.FindByID(ctx, courseID)
	if err != nil {
		return nil, err
	}
	if !actor.ownsCourse(course) {
		return nil, apperror.Forbidden("you do not own this course", nil)
	}
	return course, nil
}

// loadOwnedQuiz loads a quiz and asserts the actor owns its parent course.
func (s *Service) loadOwnedQuiz(ctx context.Context, actor *Actor, quizID uuid.UUID) (*entity.Quiz, error) {
	if !actor.isTeacherOrAdmin() {
		return nil, apperror.Forbidden("only teachers may author quizzes", nil)
	}
	quiz, err := s.quizzes.FindByID(ctx, quizID)
	if err != nil {
		return nil, err
	}
	if _, err := s.loadOwnedCourse(ctx, actor, quiz.CourseID); err != nil {
		return nil, err
	}
	return quiz, nil
}

func validateQuizInput(title string, passPct float64, timeLimit, maxAttempts *int) map[string]any {
	v := validator.New()
	v.Required("title", title)
	v.MaxLen("title", title, maxQuizTitleLen)
	if passPct < 0 || passPct > 100 {
		v.Add("pass_percentage", "must be between 0 and 100")
	}
	if timeLimit != nil && *timeLimit <= 0 {
		v.Add("time_limit_seconds", "must be positive when set")
	}
	if maxAttempts != nil && *maxAttempts <= 0 {
		v.Add("max_attempts", "must be positive when set (omit for unlimited)")
	}
	return v.Map()
}

func validateQuestionInput(in QuestionInput) map[string]any {
	v := validator.New()
	v.Required("text", in.Text)
	v.MaxLen("text", in.Text, maxQuestionTexLen)
	if !in.Type.Valid() {
		v.Add("type", "must be multiple_choice, true_false, or short_answer")
	}
	if in.Points <= 0 {
		v.Add("points", "must be positive")
	}
	if in.Explanation != nil {
		v.MaxLen("explanation", *in.Explanation, maxExplanationLen)
	}
	if len(in.Answers) == 0 {
		v.Add("answers", "at least one answer option is required")
	}
	// Every gradable question needs at least one correct option to grade against
	// (choice types AND short_answer, whose accepted text lives in a correct
	// option).
	hasCorrect := false
	for _, a := range in.Answers {
		if a.IsCorrect {
			hasCorrect = true
			break
		}
	}
	if len(in.Answers) > 0 && !hasCorrect {
		v.Add("answers", "at least one answer must be marked correct")
	}
	return v.Map()
}

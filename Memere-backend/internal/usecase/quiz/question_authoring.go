package quiz

import (
	"context"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/pagination"
)

// The authoring views below are the ONE read path that deliberately carries
// is_correct outward, because an author cannot edit an answer key they cannot
// see. They are not an exception to Non-Negotiable #1 — that rule governs what
// reaches a STUDENT. Every method here routes through loadOwnedQuiz first, so an
// authoring view is only ever built for the teacher who owns the quiz's course
// or for an admin, and the student-facing path keeps using the separate,
// structurally key-free ClientQuestion queries. Any new route serving these
// types MUST be registered behind requireAuth + teacherOrAdmin.

// AuthoredAnswer is one answer option as its author sees it, answer key included.
type AuthoredAnswer struct {
	ID         uuid.UUID
	Text       string
	IsCorrect  bool
	OrderIndex int
}

// AuthoredQuestion is a question as its author sees it: the prompt, its tags,
// the explanation, and every option with its is_correct flag.
type AuthoredQuestion struct {
	ID          uuid.UUID
	Text        string
	Type        entity.QuestionType
	Points      int
	Explanation *string
	OrderIndex  int
	Subject     *string
	Topic       *string
	Answers     []AuthoredAnswer
}

// AuthoredQuiz is a quiz plus its full question tree for the authoring UI.
type AuthoredQuiz struct {
	Quiz      *entity.Quiz
	Questions []AuthoredQuestion
}

// UpdateQuestionInput carries a full question replacement. Scalar fields are
// applied when non-nil; Answers replaces the whole option set when non-nil
// (see repository.QuestionRepository.UpdateWithAnswers for why it replaces
// rather than diffs). A nil Answers leaves the existing key untouched.
type UpdateQuestionInput struct {
	Text        *string
	Type        *entity.QuestionType
	Points      *int
	Explanation *string
	OrderIndex  *int
	Subject     *string
	Topic       *string
	Answers     []AnswerInput
}

// ListQuizzes returns quizzes matching filter, newest first, alongside each
// quiz's question tally. Teachers see only quizzes they may author: the filter
// is narrowed to their own courses, so the cross-course listing cannot leak
// another teacher's quizzes.
func (s *Service) ListQuizzes(ctx context.Context, actor *Actor, filter repository.QuizFilter, cursor *pagination.Cursor, limit int) ([]*entity.Quiz, []int, *pagination.Cursor, error) {
	if !actor.isTeacherOrAdmin() {
		return nil, nil, nil, apperror.Forbidden("only teachers may list quizzes for authoring", nil)
	}

	quizzes, next, err := s.quizzes.List(ctx, filter, cursor, limit)
	if err != nil {
		return nil, nil, nil, err
	}

	// A teacher may only see quizzes under a course they own. Filtering here
	// rather than in SQL keeps the ownership rule in exactly one place
	// (ownsCourse) at the cost of a lookup per distinct course.
	if !actor.isAdmin() {
		owned := make([]*entity.Quiz, 0, len(quizzes))
		seen := make(map[uuid.UUID]bool, len(quizzes))
		for _, q := range quizzes {
			if q.CourseID == nil {
				continue // standalone quizzes are admin-only
			}
			allowed, known := seen[*q.CourseID]
			if !known {
				course, err := s.courses.FindByID(ctx, *q.CourseID)
				allowed = err == nil && actor.ownsCourse(course)
				seen[*q.CourseID] = allowed
			}
			if allowed {
				owned = append(owned, q)
			}
		}
		quizzes = owned
	}

	counts := make([]int, len(quizzes))
	for i, q := range quizzes {
		n, err := s.quizzes.CountQuestions(ctx, q.ID)
		if err != nil {
			return nil, nil, nil, err
		}
		counts[i] = n
	}
	return quizzes, counts, next, nil
}

// GetQuizForAuthoring returns the quiz and its full question tree WITH answer
// keys, for the author's editor. Ownership-checked.
func (s *Service) GetQuizForAuthoring(ctx context.Context, actor *Actor, quizID uuid.UUID) (*AuthoredQuiz, error) {
	quiz, err := s.loadOwnedQuiz(ctx, actor, quizID)
	if err != nil {
		return nil, err
	}

	questions, err := s.questions.ListByQuiz(ctx, quizID)
	if err != nil {
		return nil, err
	}

	out := make([]AuthoredQuestion, len(questions))
	for i, q := range questions {
		answers, err := s.questions.ListAnswers(ctx, q.ID)
		if err != nil {
			return nil, err
		}
		opts := make([]AuthoredAnswer, len(answers))
		for j, a := range answers {
			opts[j] = AuthoredAnswer{ID: a.ID, Text: a.Text, IsCorrect: a.IsCorrect, OrderIndex: a.OrderIndex}
		}
		out[i] = AuthoredQuestion{
			ID:          q.ID,
			Text:        q.Text,
			Type:        q.Type,
			Points:      q.Points,
			Explanation: q.Explanation,
			OrderIndex:  q.OrderIndex,
			Subject:     q.Subject,
			Topic:       q.Topic,
			Answers:     opts,
		}
	}
	return &AuthoredQuiz{Quiz: quiz, Questions: out}, nil
}

// UpdateQuestion applies a partial update to one quiz question, replacing its
// answer set when Answers is non-nil. The question row and its options are
// written in one transaction so the answer key is never left half-written.
func (s *Service) UpdateQuestion(ctx context.Context, actor *Actor, questionID uuid.UUID, in UpdateQuestionInput) (*AuthoredQuestion, error) {
	question, err := s.loadOwnedQuestion(ctx, actor, questionID)
	if err != nil {
		return nil, err
	}

	if in.Text != nil {
		question.Text = *in.Text
	}
	if in.Type != nil {
		question.Type = *in.Type
	}
	if in.Points != nil {
		question.Points = *in.Points
	}
	if in.Explanation != nil {
		question.Explanation = in.Explanation
	}
	if in.OrderIndex != nil {
		question.OrderIndex = *in.OrderIndex
	}
	if in.Subject != nil {
		question.Subject = in.Subject
	}
	if in.Topic != nil {
		question.Topic = in.Topic
	}

	// Validate against the answer set that will actually be stored: the incoming
	// one when replacing, otherwise the one already on disk. Skipping this when
	// Answers is nil would let a type change (say to multiple_choice) land on an
	// option set that can no longer grade it.
	answers := in.Answers
	if answers == nil {
		existing, err := s.questions.ListAnswers(ctx, questionID)
		if err != nil {
			return nil, err
		}
		answers = make([]AnswerInput, len(existing))
		for i, a := range existing {
			answers[i] = AnswerInput{Text: a.Text, IsCorrect: a.IsCorrect, OrderIndex: a.OrderIndex}
		}
	}
	candidate := QuestionInput{
		Text:        question.Text,
		Type:        question.Type,
		Points:      question.Points,
		Explanation: question.Explanation,
		OrderIndex:  question.OrderIndex,
		Subject:     question.Subject,
		Topic:       question.Topic,
		Answers:     answers,
	}
	if details := validateQuestionInput(candidate); len(details) > 0 {
		return nil, apperror.Validation(details, nil)
	}

	stored := make([]*entity.Answer, len(answers))
	for i, a := range answers {
		stored[i] = &entity.Answer{Text: a.Text, IsCorrect: a.IsCorrect, OrderIndex: a.OrderIndex}
	}

	err = s.tx.WithinTx(ctx, func(ctx context.Context) error {
		return s.questions.UpdateWithAnswers(ctx, question, stored)
	})
	if err != nil {
		return nil, err
	}

	opts := make([]AuthoredAnswer, len(stored))
	for i, a := range stored {
		opts[i] = AuthoredAnswer{ID: a.ID, Text: a.Text, IsCorrect: a.IsCorrect, OrderIndex: a.OrderIndex}
	}
	return &AuthoredQuestion{
		ID:          question.ID,
		Text:        question.Text,
		Type:        question.Type,
		Points:      question.Points,
		Explanation: question.Explanation,
		OrderIndex:  question.OrderIndex,
		Subject:     question.Subject,
		Topic:       question.Topic,
		Answers:     opts,
	}, nil
}

// DeleteQuestion removes one question from its quiz. Its answer options cascade
// via the FK. Non-Negotiable #5 governs user-facing content with a tombstone
// column; a question has none and is removed with its tree, matching the
// existing QuestionRepository.Delete contract.
func (s *Service) DeleteQuestion(ctx context.Context, actor *Actor, questionID uuid.UUID) error {
	if _, err := s.loadOwnedQuestion(ctx, actor, questionID); err != nil {
		return err
	}
	return s.questions.Delete(ctx, questionID)
}

// ReorderQuestions rewrites the display order of a quiz's questions to match the
// given id sequence. Every one of the quiz's questions must appear exactly once,
// so a stale client list cannot silently drop a question from the order. The
// whole reorder runs in one transaction.
func (s *Service) ReorderQuestions(ctx context.Context, actor *Actor, quizID uuid.UUID, ordered []uuid.UUID) error {
	if _, err := s.loadOwnedQuiz(ctx, actor, quizID); err != nil {
		return err
	}

	existing, err := s.questions.ListByQuiz(ctx, quizID)
	if err != nil {
		return err
	}
	byID := make(map[uuid.UUID]*entity.Question, len(existing))
	for _, q := range existing {
		byID[q.ID] = q
	}
	if len(ordered) != len(existing) {
		return apperror.Validation(map[string]any{
			"question_ids": "must list every question in the quiz exactly once",
		}, nil)
	}
	seen := make(map[uuid.UUID]bool, len(ordered))
	for _, id := range ordered {
		if seen[id] || byID[id] == nil {
			return apperror.Validation(map[string]any{
				"question_ids": "must list every question in the quiz exactly once",
			}, nil)
		}
		seen[id] = true
	}

	return s.tx.WithinTx(ctx, func(ctx context.Context) error {
		for i, id := range ordered {
			q := byID[id]
			q.OrderIndex = i
			if err := s.questions.Update(ctx, q); err != nil {
				return err
			}
		}
		return nil
	})
}

// loadOwnedQuestion loads a question and asserts the actor may author it, by
// delegating to the ownership rule of its parent quiz.
//
// A question whose quiz_id is NULL is a question-bank row owned by an exam
// (migration 0022), not a quiz. It is rejected here rather than silently
// permitted: without a parent quiz there is no course to derive ownership from,
// so it must be edited through its exam, where the exam's own ownership rule
// applies.
func (s *Service) loadOwnedQuestion(ctx context.Context, actor *Actor, questionID uuid.UUID) (*entity.Question, error) {
	question, err := s.questions.FindByID(ctx, questionID)
	if err != nil {
		return nil, err
	}
	if question.QuizID == uuid.Nil {
		return nil, apperror.BadRequest("this question belongs to an exam, not a quiz; edit it through its exam", nil)
	}
	if _, err := s.loadOwnedQuiz(ctx, actor, question.QuizID); err != nil {
		return nil, err
	}
	return question, nil
}

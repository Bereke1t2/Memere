package exam

import (
	"context"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// The authoring views below are the ONE exam read path that deliberately carries
// is_correct outward, because an author cannot edit an answer key they cannot
// see. They are not an exception to Non-Negotiable #1 — that rule governs what
// reaches a STUDENT. Every method here routes through loadOwnedExam first, so an
// authoring view is only ever built for the teacher who owns the exam's course or
// for an admin, and the student-facing path keeps using the separate,
// structurally key-free ClientExamQuestion queries. Any new route serving these
// types MUST be registered behind requireAuth + teacherOrAdmin.

// AuthoredAnswer is one answer option as its author sees it, answer key included.
type AuthoredAnswer struct {
	ID         uuid.UUID
	Text       string
	IsCorrect  bool
	OrderIndex int
}

// AuthoredQuestion is one exam question as its author sees it: the per-exam
// order and marks, the prompt and tags, the explanation, and every option with
// its is_correct flag.
type AuthoredQuestion struct {
	ID          uuid.UUID
	Text        string
	Type        entity.QuestionType
	Marks       int
	Explanation *string
	OrderIndex  int
	Subject     *string
	Topic       *string
	Answers     []AuthoredAnswer
}

// AuthoredExam is an exam plus its full question set for the authoring UI.
type AuthoredExam struct {
	Exam      *entity.Exam
	Questions []AuthoredQuestion
}

// UpdateExamQuestionInput carries a partial update to one exam question. Scalar
// fields apply when non-nil; Answers replaces the whole option set when non-nil.
// Marks is the per-exam mark value stored on the link row, and changing it
// re-derives the exam's total_marks.
type UpdateExamQuestionInput struct {
	Text        *string
	Type        *entity.QuestionType
	Marks       *int
	Explanation *string
	OrderIndex  *int
	Subject     *string
	Topic       *string
	Answers     []InlineAnswerInput
}

// GetExamForAuthoring returns the exam and its full question set WITH answer
// keys, for the author's editor. Ownership-checked.
func (s *Service) GetExamForAuthoring(ctx context.Context, actor *Actor, examID uuid.UUID) (*AuthoredExam, error) {
	exam, err := s.loadOwnedExam(ctx, actor, examID)
	if err != nil {
		return nil, err
	}
	if s.questions == nil {
		return nil, apperror.Internal(nil)
	}

	links, err := s.exams.ListQuestions(ctx, examID)
	if err != nil {
		return nil, err
	}

	out := make([]AuthoredQuestion, 0, len(links))
	for _, link := range links {
		q, err := s.questions.FindByID(ctx, link.QuestionID)
		if err != nil {
			return nil, err
		}
		answers, err := s.questions.ListAnswers(ctx, q.ID)
		if err != nil {
			return nil, err
		}
		opts := make([]AuthoredAnswer, len(answers))
		for j, a := range answers {
			opts[j] = AuthoredAnswer{ID: a.ID, Text: a.Text, IsCorrect: a.IsCorrect, OrderIndex: a.OrderIndex}
		}
		out = append(out, AuthoredQuestion{
			ID:          q.ID,
			Text:        q.Text,
			Type:        q.Type,
			Marks:       link.Marks,
			Explanation: q.Explanation,
			OrderIndex:  link.OrderIndex,
			Subject:     q.Subject,
			Topic:       q.Topic,
			Answers:     opts,
		})
	}
	return &AuthoredExam{Exam: exam, Questions: out}, nil
}

// UpdateExamQuestion edits one question of an exam: its body and tags, its
// answer set, and its per-exam marks and order. Everything runs in one
// transaction, and the exam's total_marks is re-derived from the link rows
// afterwards so the header total never drifts from the parts.
func (s *Service) UpdateExamQuestion(ctx context.Context, actor *Actor, examID, questionID uuid.UUID, in UpdateExamQuestionInput) (*AuthoredQuestion, error) {
	exam, err := s.loadOwnedExam(ctx, actor, examID)
	if err != nil {
		return nil, err
	}
	if s.questions == nil {
		return nil, apperror.Internal(nil)
	}

	// The link row is the proof that this question belongs to THIS exam. Without
	// it, an id from another exam could be edited through this path.
	link, err := s.findLink(ctx, examID, questionID)
	if err != nil {
		return nil, err
	}

	question, err := s.questions.FindByID(ctx, questionID)
	if err != nil {
		return nil, err
	}

	if in.Text != nil {
		question.Text = *in.Text
	}
	if in.Type != nil {
		question.Type = *in.Type
	}
	if in.Explanation != nil {
		question.Explanation = in.Explanation
	}
	if in.Subject != nil {
		question.Subject = in.Subject
	}
	if in.Topic != nil {
		question.Topic = in.Topic
	}
	if in.OrderIndex != nil {
		link.OrderIndex = *in.OrderIndex
	}
	if in.Marks != nil {
		link.Marks = *in.Marks
	}
	// points mirrors the per-exam marks so the question grades the same whether
	// it is reached through the exam or straight from the bank.
	question.Points = link.Marks

	answers := in.Answers
	if answers == nil {
		existing, err := s.questions.ListAnswers(ctx, questionID)
		if err != nil {
			return nil, err
		}
		answers = make([]InlineAnswerInput, len(existing))
		for i, a := range existing {
			answers[i] = InlineAnswerInput{Text: a.Text, IsCorrect: a.IsCorrect, OrderIndex: a.OrderIndex}
		}
	}
	if details := validateExamQuestion(question.Text, question.Type, link.Marks, answers); len(details) > 0 {
		return nil, apperror.Validation(details, nil)
	}

	stored := make([]*entity.Answer, len(answers))
	for i, a := range answers {
		stored[i] = &entity.Answer{Text: a.Text, IsCorrect: a.IsCorrect, OrderIndex: a.OrderIndex}
	}

	err = s.tx.WithinTx(ctx, func(ctx context.Context) error {
		if err := s.questions.UpdateWithAnswers(ctx, question, stored); err != nil {
			return err
		}
		if err := s.exams.UpdateQuestionLink(ctx, link); err != nil {
			return err
		}
		return s.retotal(ctx, exam)
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
		Marks:       link.Marks,
		Explanation: question.Explanation,
		OrderIndex:  link.OrderIndex,
		Subject:     question.Subject,
		Topic:       question.Topic,
		Answers:     opts,
	}, nil
}

// RemoveExamQuestion detaches a question from an exam and re-derives
// total_marks, transactionally. The question itself stays in the bank, so any
// other exam reusing it is unaffected.
func (s *Service) RemoveExamQuestion(ctx context.Context, actor *Actor, examID, questionID uuid.UUID) (*entity.Exam, error) {
	exam, err := s.loadOwnedExam(ctx, actor, examID)
	if err != nil {
		return nil, err
	}
	if _, err := s.findLink(ctx, examID, questionID); err != nil {
		return nil, err
	}

	err = s.tx.WithinTx(ctx, func(ctx context.Context) error {
		if err := s.exams.RemoveQuestion(ctx, examID, questionID); err != nil {
			return err
		}
		return s.retotal(ctx, exam)
	})
	if err != nil {
		return nil, err
	}
	return exam, nil
}

// ReorderExamQuestions rewrites the display order of an exam's questions to
// match the given id sequence. Every linked question must appear exactly once,
// so a stale client list cannot silently drop one from the order.
func (s *Service) ReorderExamQuestions(ctx context.Context, actor *Actor, examID uuid.UUID, ordered []uuid.UUID) error {
	if _, err := s.loadOwnedExam(ctx, actor, examID); err != nil {
		return err
	}

	links, err := s.exams.ListQuestions(ctx, examID)
	if err != nil {
		return err
	}
	byID := make(map[uuid.UUID]*entity.ExamQuestion, len(links))
	for _, l := range links {
		byID[l.QuestionID] = l
	}
	if len(ordered) != len(links) {
		return apperror.Validation(map[string]any{
			"question_ids": "must list every question in the exam exactly once",
		}, nil)
	}
	seen := make(map[uuid.UUID]bool, len(ordered))
	for _, id := range ordered {
		if seen[id] || byID[id] == nil {
			return apperror.Validation(map[string]any{
				"question_ids": "must list every question in the exam exactly once",
			}, nil)
		}
		seen[id] = true
	}

	return s.tx.WithinTx(ctx, func(ctx context.Context) error {
		for i, id := range ordered {
			link := byID[id]
			link.OrderIndex = i
			if err := s.exams.UpdateQuestionLink(ctx, link); err != nil {
				return err
			}
		}
		return nil
	})
}

// findLink returns the exam_questions row joining examID to questionID, or a
// 404 when the question is not part of that exam.
func (s *Service) findLink(ctx context.Context, examID, questionID uuid.UUID) (*entity.ExamQuestion, error) {
	links, err := s.exams.ListQuestions(ctx, examID)
	if err != nil {
		return nil, err
	}
	for _, l := range links {
		if l.QuestionID == questionID {
			return l, nil
		}
	}
	return nil, apperror.NotFound("question is not part of this exam", nil)
}

// retotal re-derives exam.TotalMarks from the current link rows and persists it.
// Call inside a transaction alongside whatever changed the link set.
func (s *Service) retotal(ctx context.Context, exam *entity.Exam) error {
	links, err := s.exams.ListQuestions(ctx, exam.ID)
	if err != nil {
		return err
	}
	total := 0
	for _, l := range links {
		total += l.Marks
	}
	exam.TotalMarks = total
	return s.exams.Update(ctx, exam)
}

// validateExamQuestion checks one exam question and its option set. A gradable
// question needs at least one correct option, otherwise the grader has nothing
// to score against.
func validateExamQuestion(text string, qType entity.QuestionType, marks int, answers []InlineAnswerInput) map[string]any {
	details := map[string]any{}
	if text == "" {
		details["text"] = "must not be empty"
	}
	if !qType.Valid() {
		details["type"] = "must be multiple_choice, true_false, or short_answer"
	}
	if marks <= 0 {
		details["marks"] = "must be positive"
	}
	if len(answers) < 2 {
		details["answers"] = "must have at least 2 options"
	} else {
		hasCorrect := false
		for _, a := range answers {
			if a.IsCorrect {
				hasCorrect = true
				break
			}
		}
		if !hasCorrect {
			details["answers"] = "at least one answer must be marked correct"
		}
	}
	return details
}

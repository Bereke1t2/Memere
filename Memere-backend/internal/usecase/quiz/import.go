package quiz

import (
	"context"
	"fmt"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// maxImportQuestions caps one import so a single request cannot open a
// long-running transaction that holds locks across hundreds of inserts.
const maxImportQuestions = 500

// ImportQuizInput is a whole quiz — header plus question set — created in one
// request. OrderIndex on each question is ignored: questions are numbered by
// their position in the slice, so the file's reading order is the quiz order.
type ImportQuizInput struct {
	Quiz      CreateQuizInput
	Questions []QuestionInput
}

// ImportQuiz creates a quiz and every one of its questions in a single
// transaction. Either the whole file lands or none of it does: a half-imported
// quiz with a truncated question set would look complete to a student while
// silently mis-scoring, so a failure on question 40 rolls back the 39 before it.
//
// Every question is validated BEFORE anything is written, and the returned
// error names the offending question by its index in the file, so an author
// fixes their source rather than hunting through a partially-created quiz.
func (s *Service) ImportQuiz(ctx context.Context, actor *Actor, in ImportQuizInput) (*entity.Quiz, int, error) {
	if len(in.Questions) == 0 {
		return nil, 0, apperror.Validation(map[string]any{
			"questions": "must contain at least one question",
		}, nil)
	}
	if len(in.Questions) > maxImportQuestions {
		return nil, 0, apperror.Validation(map[string]any{
			"questions": fmt.Sprintf("must contain at most %d questions per import", maxImportQuestions),
		}, nil)
	}

	details := map[string]any{}
	for i, q := range in.Questions {
		if qDetails := validateQuestionInput(q); len(qDetails) > 0 {
			details[fmt.Sprintf("questions[%d]", i)] = qDetails
		}
	}
	if len(details) > 0 {
		return nil, 0, apperror.Validation(details, nil)
	}

	var created *entity.Quiz
	err := s.tx.WithinTx(ctx, func(ctx context.Context) error {
		quiz, err := s.CreateQuiz(ctx, actor, in.Quiz)
		if err != nil {
			return err
		}
		for i, qIn := range in.Questions {
			question := &entity.Question{
				QuizID:      quiz.ID,
				Text:        qIn.Text,
				Type:        qIn.Type,
				Points:      qIn.Points,
				Explanation: qIn.Explanation,
				// Position in the file is the authoritative order.
				OrderIndex: i,
				Subject:    qIn.Subject,
				Topic:      qIn.Topic,
			}
			answers := make([]*entity.Answer, len(qIn.Answers))
			for j, a := range qIn.Answers {
				answers[j] = &entity.Answer{Text: a.Text, IsCorrect: a.IsCorrect, OrderIndex: j}
			}
			if err := s.questions.Create(ctx, question, answers); err != nil {
				return err
			}
		}
		created = quiz
		return nil
	})
	if err != nil {
		return nil, 0, err
	}
	return created, len(in.Questions), nil
}

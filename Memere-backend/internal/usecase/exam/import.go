package exam

import (
	"context"
	"fmt"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// maxImportQuestions caps one import so a single request cannot open a
// long-running transaction that holds locks across hundreds of inserts.
const maxImportQuestions = 500

// ImportExamInput is a whole exam — header plus question set — created in one
// request. OrderIndex is ignored: questions are numbered by their position in
// the slice, so the file's reading order is the exam order.
type ImportExamInput struct {
	Exam      CreateExamInput
	Questions []AddInlineQuestionInput
}

// ImportExam creates an exam and every one of its questions in a single
// transaction, deriving total_marks from the question marks. Either the whole
// file lands or none of it does: a half-imported exam would present a total_marks
// that does not match its questions and mis-score every sitting, so a failure on
// question 40 rolls back the 39 before it.
//
// Every question is validated BEFORE anything is written, and the returned error
// names the offending question by its index in the file.
func (s *Service) ImportExam(ctx context.Context, actor *Actor, in ImportExamInput) (*entity.Exam, int, error) {
	if s.questions == nil {
		return nil, 0, apperror.Internal(nil)
	}
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
	total := 0
	for i, q := range in.Questions {
		qType := q.Type
		if qType == "" {
			qType = entity.QuestionMultipleChoice
		}
		if qDetails := validateExamQuestion(q.Text, qType, q.Marks, q.Answers); len(qDetails) > 0 {
			details[fmt.Sprintf("questions[%d]", i)] = qDetails
		}
		total += q.Marks
	}
	if len(details) > 0 {
		return nil, 0, apperror.Validation(details, nil)
	}

	// pass_marks is validated against the total the question set will actually
	// produce. CreateExam cannot check this itself: it runs before the questions
	// exist, when total_marks is still 0.
	if in.Exam.PassMarks > total {
		return nil, 0, apperror.Validation(map[string]any{
			"pass_marks": fmt.Sprintf("must not exceed the imported total of %d marks", total),
		}, nil)
	}

	var created *entity.Exam
	err := s.tx.WithinTx(ctx, func(ctx context.Context) error {
		exam, err := s.CreateExam(ctx, actor, in.Exam)
		if err != nil {
			return err
		}
		for i, qIn := range in.Questions {
			qType := qIn.Type
			if qType == "" {
				qType = entity.QuestionMultipleChoice
			}
			question := &entity.Question{
				Text:        qIn.Text,
				Type:        qType,
				Points:      qIn.Marks,
				Explanation: qIn.Explanation,
				OrderIndex:  i,
				Subject:     qIn.Subject,
				Topic:       qIn.Topic,
			}
			answers := make([]*entity.Answer, len(qIn.Answers))
			for j, a := range qIn.Answers {
				answers[j] = &entity.Answer{Text: a.Text, IsCorrect: a.IsCorrect, OrderIndex: j}
			}
			if err := s.questions.Create(ctx, question, answers); err != nil {
				return err
			}
			link := &entity.ExamQuestion{
				ExamID:     exam.ID,
				QuestionID: question.ID,
				OrderIndex: i,
				Marks:      qIn.Marks,
			}
			if err := s.exams.AddQuestion(ctx, link); err != nil {
				return err
			}
		}
		exam.TotalMarks = total
		if err := s.exams.Update(ctx, exam); err != nil {
			return err
		}
		created = exam
		return nil
	})
	if err != nil {
		return nil, 0, err
	}
	return created, len(in.Questions), nil
}

package exam

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/grading"
)

// attemptView fetches the client questions for a resumed attempt and assembles
// the live view, computing the deadline from started_at + duration.
func (s *Service) attemptView(ctx context.Context, exam *entity.Exam, attempt *entity.ExamAttempt) (*ExamAttemptClientView, error) {
	clientQs, err := s.exams.GetQuestionsForClient(ctx, exam.ID)
	if err != nil {
		return nil, err
	}
	expiresAt := attempt.StartedAt.Add(time.Duration(exam.DurationMinutes) * time.Minute)
	return s.buildView(exam, attempt, &expiresAt, clientQs), nil
}

// buildView assembles the answer-key-free client view for an attempt.
func (s *Service) buildView(exam *entity.Exam, attempt *entity.ExamAttempt, expiresAt *time.Time, clientQs []repository.ClientExamQuestion) *ExamAttemptClientView {
	var remaining *int
	if expiresAt != nil {
		secs := int(expiresAt.Sub(s.now()).Seconds())
		if secs < 0 {
			secs = 0
		}
		remaining = &secs
	}

	questions := make([]QuestionClientView, len(clientQs))
	for i, q := range clientQs {
		answers := make([]AnswerClientView, len(q.Answers))
		for j, a := range q.Answers {
			answers[j] = AnswerClientView{ID: a.ID, Text: a.Text, OrderIndex: a.OrderIndex}
		}
		questions[i] = QuestionClientView{
			QuestionID: q.QuestionID,
			Text:       q.Text,
			Type:       q.Type,
			Marks:      q.Marks,
			Subject:    q.Subject,
			Topic:      q.Topic,
			Answers:    answers,
		}
	}

	return &ExamAttemptClientView{
		AttemptID:        attempt.ID,
		ExamID:           exam.ID,
		Status:           attempt.Status,
		StartedAt:        attempt.StartedAt,
		ExpiresAt:        expiresAt,
		RemainingSeconds: remaining,
		TotalMarks:       exam.TotalMarks,
		Questions:        questions,
	}
}

// gradeView re-grades a terminal attempt from its persisted answers_snapshot and
// returns the result view (used by GetExamResult for already-graded attempts).
func (s *Service) gradeView(ctx context.Context, exam *entity.Exam, attempt *entity.ExamAttempt) (*ExamResult, error) {
	gradables, err := s.loadGradables(ctx, exam.ID)
	if err != nil {
		return nil, err
	}
	answers := answersFromSnapshot(attempt.AnswersSnapshot)
	sub := parseSubmitted(gradables, answers)
	core := grading.Grade(gradables, sub)
	res := s.result(exam, attempt, core)
	res.SubmittedAt = attempt.SubmittedAt
	return res, nil
}

// loadGradables assembles the server-side grading input for an exam: each
// question's per-exam marks (from exam_questions) joined to its options WITH the
// answer key. This is the only path that reads is_correct.
func (s *Service) loadGradables(ctx context.Context, examID uuid.UUID) ([]grading.GradableQuestion, error) {
	links, err := s.exams.ListQuestions(ctx, examID)
	if err != nil {
		return nil, err
	}
	marksByQ := make(map[uuid.UUID]int, len(links))
	for _, l := range links {
		marksByQ[l.QuestionID] = l.Marks
	}

	tree, err := s.exams.GetExamWithQuestions(ctx, examID)
	if err != nil {
		return nil, err
	}

	gradables := make([]grading.GradableQuestion, 0, len(tree.Questions))
	for _, qa := range tree.Questions {
		answers := make([]grading.GradableAnswer, len(qa.Answers))
		for i, a := range qa.Answers {
			answers[i] = grading.GradableAnswer{ID: a.ID, Text: a.Text, IsCorrect: a.IsCorrect}
		}
		gradables = append(gradables, grading.GradableQuestion{
			ID:          qa.Question.ID,
			Text:        qa.Question.Text,
			Marks:       marksByQ[qa.Question.ID],
			Type:        qa.Question.Type,
			Subject:     qa.Question.Subject,
			Topic:       qa.Question.Topic,
			Explanation: qa.Question.Explanation,
			Answers:     answers,
		})
	}
	return gradables, nil
}

// result maps the shared grading output to the exam result view. Pass/fail is
// decided on the exam's absolute pass_marks (§4.2.6), not a percentage.
func (s *Service) result(exam *entity.Exam, attempt *entity.ExamAttempt, core grading.Result) *ExamResult {
	feedback := make([]QuestionFeedback, 0, len(core.Outcomes))
	for _, o := range core.Outcomes {
		answers := make([]FeedbackAnswerView, len(o.Answers))
		for i, a := range o.Answers {
			answers[i] = FeedbackAnswerView{
				ID:         a.ID,
				Text:       a.Text,
				IsCorrect:  a.IsCorrect,
				OrderIndex: i,
			}
		}
		feedback = append(feedback, QuestionFeedback{
			QuestionID:       o.QuestionID,
			QuestionText:     o.QuestionText,
			Type:             o.Type,
			Subject:          o.Subject,
			Topic:            o.Topic,
			Correct:          o.Correct,
			MarksAwarded:     o.MarksAwarded,
			MarksPossible:    o.MarksPossible,
			SelectedAnswers:  o.SelectedAnswers,
			CorrectAnswerIDs: o.CorrectAnswerIDs,
			Explanation:      o.Explanation,
			Answers:          answers,
		})
	}
	breakdown := make(map[string]SubjectScore, len(core.SubjectBreakdown))
	for k, v := range core.SubjectBreakdown {
		breakdown[k] = SubjectScore{Earned: v.Earned, Possible: v.Possible}
	}

	return &ExamResult{
		AttemptID:        attempt.ID,
		ExamID:           exam.ID,
		Status:           attempt.Status,
		Score:            core.Score,
		TotalMarks:       core.Total,
		Percentage:       core.Percentage,
		PassMarks:        exam.PassMarks,
		Passed:           int(core.Score) >= exam.PassMarks,
		Feedback:         feedback,
		SubjectBreakdown: breakdown,
	}
}

// orderSnapshot records the question presentation order for Redis resume.
func orderSnapshot(qs []repository.ClientExamQuestion) map[string]any {
	ids := make([]any, len(qs))
	for i, q := range qs {
		ids[i] = q.QuestionID.String()
	}
	return map[string]any{"questions": ids}
}

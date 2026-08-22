package quiz

import (
	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/grading"
)

// submittedAnswers maps a question ID to the answer IDs the student selected.
// For short_answer questions the free-text response is carried separately.
type submittedAnswers struct {
	selected map[uuid.UUID][]uuid.UUID
	text     map[uuid.UUID]string
}

// gradeQuiz scores a quiz attempt via the shared grading core (which is the only
// place is_correct is consulted, server-side only). A quiz question's worth is
// its points; passed is decided here on the quiz's percentage threshold (§9.1),
// which the engine-neutral core deliberately leaves to the caller.
func gradeQuiz(
	quiz *entity.Quiz,
	questions []repository.QuestionWithAnswers,
	sub submittedAnswers,
) AttemptResult {
	gradable := make([]grading.GradableQuestion, len(questions))
	for i, qa := range questions {
		gradable[i] = grading.GradableQuestion{
			ID:          qa.Question.ID,
			Text:        qa.Question.Text,
			Marks:       qa.Question.Points,
			Type:        qa.Question.Type,
			Subject:     qa.Question.Subject,
			Topic:       qa.Question.Topic,
			Explanation: qa.Question.Explanation,
			Answers:     toGradableAnswers(qa.Answers),
		}
	}

	core := grading.Grade(gradable, grading.Submitted{Selected: sub.selected, Text: sub.text})

	result := AttemptResult{
		QuizID:           quiz.ID,
		Score:            core.Score,
		TotalPoints:      core.Total,
		Percentage:       core.Percentage,
		Passed:           core.Percentage >= quiz.PassPercentage,
		Feedback:         make([]QuestionFeedback, 0, len(core.Outcomes)),
		SubjectBreakdown: toSubjectScores(core.SubjectBreakdown),
	}
	for _, o := range core.Outcomes {
		answers := make([]AnswerFeedback, 0, len(o.Answers))
		for _, a := range o.Answers {
			answers = append(answers, AnswerFeedback{ID: a.ID, Text: a.Text, IsCorrect: a.IsCorrect})
		}
		result.Feedback = append(result.Feedback, QuestionFeedback{
			QuestionID:       o.QuestionID,
			QuestionText:     o.QuestionText,
			Correct:          o.Correct,
			PointsAwarded:    o.MarksAwarded,
			PointsPossible:   o.MarksPossible,
			SelectedAnswers:  o.SelectedAnswers,
			CorrectAnswerIDs: o.CorrectAnswerIDs,
			Explanation:      o.Explanation,
			Answers:          answers,
		})
	}
	return result
}

func toGradableAnswers(answers []*entity.Answer) []grading.GradableAnswer {
	out := make([]grading.GradableAnswer, len(answers))
	for i, a := range answers {
		out[i] = grading.GradableAnswer{ID: a.ID, Text: a.Text, IsCorrect: a.IsCorrect}
	}
	return out
}

func toSubjectScores(m map[string]grading.SubjectScore) map[string]SubjectScore {
	out := make(map[string]SubjectScore, len(m))
	for k, v := range m {
		out[k] = SubjectScore{Earned: v.Earned, Possible: v.Possible}
	}
	return out
}

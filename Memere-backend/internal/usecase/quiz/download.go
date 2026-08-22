package quiz

import (
	"context"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
)

// QuizDownload is the full quiz payload for OFFLINE download, INCLUDING the
// answer key (Answer.IsCorrect) and each question's explanation.
//
// This is the single sanctioned exception to Non-Negotiable #1 ("correct answers
// never sent to the client"): keys are served ONLY for content the caller
// explicitly downloads, and the endpoint is gated by the SAME course-access check
// as taking the quiz (assertCourseAccess). Re-enabling the paywall
// (access.DisableEnrollmentCheck=false) therefore auto-restricts downloads too.
// The normal online-play path (GetQuizForStudent / GetQuestionsForClient) stays
// answer-free and server-graded, and offline scores are re-graded server-side on
// reconnect, so they are provisional only.
type QuizDownload struct {
	Quiz      *entity.Quiz
	Questions []DownloadQuestion
}

// DownloadQuestion is one question with its full answer key and explanation, for
// on-device grading while offline.
type DownloadQuestion struct {
	ID          uuid.UUID
	Text        string
	Type        entity.QuestionType
	Points      int
	OrderIndex  int
	Subject     *string
	Topic       *string
	Explanation *string
	Answers     []DownloadAnswer
}

// DownloadAnswer is one option INCLUDING whether it is correct.
type DownloadAnswer struct {
	ID         uuid.UUID
	Text       string
	OrderIndex int
	IsCorrect  bool
}

// GetQuizForDownload returns the full quiz WITH answer keys for offline use.
// Access is gated exactly like StartAttempt (assertCourseAccess); the payload is
// built from the server-only GetQuizWithQuestions tree — the one tree that
// intentionally carries Answer.IsCorrect and Question.Explanation. Everything
// else in the quiz engine keeps stripping the key.
func (s *Service) GetQuizForDownload(ctx context.Context, actor *Actor, quizID uuid.UUID) (*QuizDownload, error) {
	quiz, err := s.quizzes.FindByID(ctx, quizID)
	if err != nil {
		return nil, err
	}
	if err := s.assertCourseAccess(ctx, actor, quiz.CourseID); err != nil {
		return nil, err
	}

	tree, err := s.quizzes.GetQuizWithQuestions(ctx, quizID)
	if err != nil {
		return nil, err
	}

	out := &QuizDownload{Quiz: quiz, Questions: make([]DownloadQuestion, 0, len(tree.Questions))}
	for _, qa := range tree.Questions {
		answers := make([]DownloadAnswer, 0, len(qa.Answers))
		for _, a := range qa.Answers {
			answers = append(answers, DownloadAnswer{
				ID:         a.ID,
				Text:       a.Text,
				OrderIndex: a.OrderIndex,
				IsCorrect:  a.IsCorrect,
			})
		}
		out.Questions = append(out.Questions, DownloadQuestion{
			ID:          qa.Question.ID,
			Text:        qa.Question.Text,
			Type:        qa.Question.Type,
			Points:      qa.Question.Points,
			OrderIndex:  qa.Question.OrderIndex,
			Subject:     qa.Question.Subject,
			Topic:       qa.Question.Topic,
			Explanation: qa.Question.Explanation,
			Answers:     answers,
		})
	}
	return out, nil
}

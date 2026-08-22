package exam

import (
	"context"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
)

// ExamDownload is the full exam payload for OFFLINE download, INCLUDING the
// answer key and each question's explanation.
//
// Like QuizDownload it is the narrow, sanctioned exception to Non-Negotiable #1
// ("correct answers never sent to the client"): served ONLY for content the
// caller explicitly downloads and gated by the SAME access check as StartExam
// (assertExamAccess), so an unpublished exam — or a re-enabled paywall — still
// withholds the keys. The online-sitting path (GetQuestionsForClient) stays
// answer-free, and the offline exam's client-side timer/score is provisional:
// the raw answers are re-graded server-side on reconnect (mitigating the
// necessary Non-Negotiable #2 relaxation for offline use).
type ExamDownload struct {
	Exam      *entity.Exam
	Questions []DownloadQuestion
}

// DownloadQuestion is one exam question with its per-exam marks, full answer key,
// and explanation for on-device grading while offline.
type DownloadQuestion struct {
	ID          uuid.UUID
	Text        string
	Type        entity.QuestionType
	Marks       int
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

// GetExamForDownload returns the full exam WITH answer keys for offline use.
// Access is gated exactly like StartExam (assertExamAccess); the payload reuses
// loadGradables — the one path that joins per-exam marks to the options WITH the
// key. OrderIndex is assigned by position, matching how s.result builds feedback.
func (s *Service) GetExamForDownload(ctx context.Context, actor *Actor, examID uuid.UUID) (*ExamDownload, error) {
	exam, err := s.exams.FindByID(ctx, examID)
	if err != nil {
		return nil, err
	}
	if err := s.assertExamAccess(ctx, actor, exam); err != nil {
		return nil, err
	}

	gradables, err := s.loadGradables(ctx, examID)
	if err != nil {
		return nil, err
	}

	out := &ExamDownload{Exam: exam, Questions: make([]DownloadQuestion, 0, len(gradables))}
	for qi, g := range gradables {
		answers := make([]DownloadAnswer, len(g.Answers))
		for i, a := range g.Answers {
			answers[i] = DownloadAnswer{
				ID:         a.ID,
				Text:       a.Text,
				OrderIndex: i,
				IsCorrect:  a.IsCorrect,
			}
		}
		out.Questions = append(out.Questions, DownloadQuestion{
			ID:          g.ID,
			Text:        g.Text,
			Type:        g.Type,
			Marks:       g.Marks,
			OrderIndex:  qi,
			Subject:     g.Subject,
			Topic:       g.Topic,
			Explanation: g.Explanation,
			Answers:     answers,
		})
	}
	return out, nil
}

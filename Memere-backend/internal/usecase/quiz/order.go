package quiz

import (
	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/shuffle"
)

// buildOrder produces the durable snapshot of an attempt's question and per-
// question answer order. It is a JSONB-friendly map:
//
//	{
//	  "questions": ["<qid>", ...],            // question display order
//	  "answers":   {"<qid>": ["<aid>", ...]}  // option order within each question
//	}
//
// The order is deterministic in the attempt ID (see pkg/shuffle), so it can be
// reproduced if the Redis copy is lost. Questions are shuffled only when the
// quiz enables randomization; answer options are always shuffled (presentation
// only — the answer key is never exposed by reordering).
func buildOrder(attemptID uuid.UUID, quiz *entity.Quiz, qs []repository.ClientQuestion) map[string]any {
	qPerm := shuffle.Identity(len(qs))
	if quiz.RandomizeQuestions {
		qPerm = shuffle.Order(attemptID, questionSalt, len(qs))
	}

	questionOrder := make([]any, len(qs))
	answerOrder := make(map[string]any, len(qs))
	for newIdx, srcIdx := range qPerm {
		q := qs[srcIdx]
		questionOrder[newIdx] = q.ID.String()

		aPerm := shuffle.Order(attemptID, answerSalt+uint64(srcIdx), len(q.Answers))
		ids := make([]any, len(q.Answers))
		for ai, src := range aPerm {
			ids[ai] = q.Answers[src].ID.String()
		}
		answerOrder[q.ID.String()] = ids
	}

	return map[string]any{
		"questions": questionOrder,
		"answers":   answerOrder,
	}
}

// orderQuestions renders the client questions in the snapshot's order. Anything
// the snapshot doesn't mention falls back to the natural order so a quiz edited
// mid-attempt still renders rather than dropping questions.
func orderQuestions(qs []repository.ClientQuestion, order map[string]any) []QuestionClientView {
	byID := make(map[uuid.UUID]repository.ClientQuestion, len(qs))
	for _, q := range qs {
		byID[q.ID] = q
	}

	out := make([]QuestionClientView, 0, len(qs))
	seen := map[uuid.UUID]bool{}

	for _, raw := range toStringSlice(order["questions"]) {
		qid, err := uuid.Parse(raw)
		if err != nil {
			continue
		}
		q, ok := byID[qid]
		if !ok || seen[qid] {
			continue
		}
		seen[qid] = true
		out = append(out, questionView(q, answerOrderFor(order, raw)))
	}
	// Append any questions not covered by the snapshot, in natural order.
	for _, q := range qs {
		if !seen[q.ID] {
			out = append(out, questionView(q, nil))
		}
	}
	return out
}

// questionView maps a repo ClientQuestion to the usecase view, applying the
// per-question answer order when provided.
func questionView(q repository.ClientQuestion, answerIDOrder []string) QuestionClientView {
	answers := make([]AnswerClientView, 0, len(q.Answers))
	if len(answerIDOrder) > 0 {
		byID := make(map[uuid.UUID]repository.ClientAnswer, len(q.Answers))
		for _, a := range q.Answers {
			byID[a.ID] = a
		}
		seen := map[uuid.UUID]bool{}
		for _, raw := range answerIDOrder {
			aid, err := uuid.Parse(raw)
			if err != nil {
				continue
			}
			a, ok := byID[aid]
			if !ok || seen[aid] {
				continue
			}
			seen[aid] = true
			answers = append(answers, AnswerClientView{ID: a.ID, Text: a.Text, OrderIndex: len(answers)})
		}
		for _, a := range q.Answers {
			if !seen[a.ID] {
				answers = append(answers, AnswerClientView{ID: a.ID, Text: a.Text, OrderIndex: len(answers)})
			}
		}
	} else {
		for i, a := range q.Answers {
			answers = append(answers, AnswerClientView{ID: a.ID, Text: a.Text, OrderIndex: i})
		}
	}

	return QuestionClientView{
		ID:      q.ID,
		Text:    q.Text,
		Type:    q.Type,
		Points:  q.Points,
		Subject: q.Subject,
		Topic:   q.Topic,
		Answers: answers,
	}
}

func answerOrderFor(order map[string]any, questionID string) []string {
	answers, ok := order["answers"].(map[string]any)
	if !ok {
		return nil
	}
	return toStringSlice(answers[questionID])
}

// toStringSlice coerces a JSONB-decoded []any of strings into []string.
func toStringSlice(v any) []string {
	raw, ok := v.([]any)
	if !ok {
		return nil
	}
	out := make([]string, 0, len(raw))
	for _, e := range raw {
		if s, ok := e.(string); ok {
			out = append(out, s)
		}
	}
	return out
}

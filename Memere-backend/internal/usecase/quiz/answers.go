package quiz

import (
	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
)

// The submit/auto-save answer payload is a JSON object keyed by question ID. Each
// value is either:
//
//	{"selected": ["<answerID>", ...]}   // choice questions
//	{"text": "free response"}           // short_answer questions
//
// A bare array value (["<answerID>", ...]) or bare string is also accepted as a
// convenience shorthand for selected / text respectively.

// mergeAnswers overlays newer answers (the submit payload) on top of older ones
// (the Redis-saved auto-save). Keys present in newer win.
func mergeAnswers(older, newer map[string]any) map[string]any {
	out := make(map[string]any, len(older)+len(newer))
	for k, v := range older {
		out[k] = v
	}
	for k, v := range newer {
		out[k] = v
	}
	return out
}

// parseSubmitted converts the raw answer payload into the typed submittedAnswers
// the grader consumes, ignoring keys that don't correspond to a real question in
// the tree (so a malicious or stale client can't inject phantom questions).
func parseSubmitted(questions []repository.QuestionWithAnswers, payload map[string]any) submittedAnswers {
	valid := make(map[uuid.UUID]bool, len(questions))
	for _, qa := range questions {
		valid[qa.Question.ID] = true
	}

	sub := submittedAnswers{
		selected: map[uuid.UUID][]uuid.UUID{},
		text:     map[uuid.UUID]string{},
	}
	for rawQID, rawVal := range payload {
		qid, err := uuid.Parse(rawQID)
		if err != nil || !valid[qid] {
			continue
		}
		selected, text := extractAnswer(rawVal)
		if len(selected) > 0 {
			sub.selected[qid] = selected
		}
		if text != "" {
			sub.text[qid] = text
		}
	}
	return sub
}

// extractAnswer pulls the selected answer IDs and/or free text out of one
// question's payload value, tolerating the object and shorthand forms.
func extractAnswer(v any) (selected []uuid.UUID, text string) {
	switch val := v.(type) {
	case string:
		return nil, val
	case []any:
		return parseIDList(val), ""
	case map[string]any:
		if s, ok := val["text"].(string); ok {
			text = s
		}
		if arr, ok := val["selected"].([]any); ok {
			selected = parseIDList(arr)
		}
		return selected, text
	default:
		return nil, ""
	}
}

func parseIDList(arr []any) []uuid.UUID {
	out := make([]uuid.UUID, 0, len(arr))
	for _, e := range arr {
		s, ok := e.(string)
		if !ok {
			continue
		}
		if id, err := uuid.Parse(s); err == nil {
			out = append(out, id)
		}
	}
	return out
}

// snapshotFor builds the answers_snapshot JSONB persisted with the attempt: the
// student's raw answers plus the server-computed per-subject breakdown for §9.3
// analytics. The correct-answer key is never included.
func snapshotFor(answers map[string]any, result AttemptResult) map[string]any {
	breakdown := make(map[string]any, len(result.SubjectBreakdown))
	for k, v := range result.SubjectBreakdown {
		breakdown[k] = map[string]any{"earned": v.Earned, "possible": v.Possible}
	}
	return map[string]any{
		"answers":           answers,
		"subject_breakdown": breakdown,
	}
}

// answersFromSnapshot extracts the raw answers map persisted under "answers" in
// answers_snapshot (the inverse of snapshotFor).
func answersFromSnapshot(snapshot map[string]any) map[string]any {
	if snapshot == nil {
		return nil
	}
	if a, ok := snapshot["answers"].(map[string]any); ok {
		return a
	}
	return nil
}

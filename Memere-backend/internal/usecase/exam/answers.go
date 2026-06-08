package exam

import (
	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/grading"
)

// The submit/auto-save answer payload is a JSON object keyed by question ID, with
// each value either {"selected": ["<answerID>", ...]} for choice questions or
// {"text": "..."} for short_answer (bare array / bare string shorthands accepted).

// mergeAnswers overlays newer answers (submit payload) on older ones (Redis
// auto-save). Keys present in newer win.
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

// parseSubmitted converts the raw payload into the shared grading.Submitted,
// ignoring keys that don't match a real question in the exam (so a stale or
// malicious client can't inject phantom questions).
func parseSubmitted(questions []grading.GradableQuestion, payload map[string]any) grading.Submitted {
	valid := make(map[uuid.UUID]bool, len(questions))
	for _, q := range questions {
		valid[q.ID] = true
	}

	sub := grading.Submitted{
		Selected: map[uuid.UUID][]uuid.UUID{},
		Text:     map[uuid.UUID]string{},
	}
	for rawQID, rawVal := range payload {
		qid, err := uuid.Parse(rawQID)
		if err != nil || !valid[qid] {
			continue
		}
		selected, text := extractAnswer(rawVal)
		if len(selected) > 0 {
			sub.Selected[qid] = selected
		}
		if text != "" {
			sub.Text[qid] = text
		}
	}
	return sub
}

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
// student's raw answers plus the server-computed per-subject breakdown for §9.3.
// The correct-answer key is never included.
func snapshotFor(answers map[string]any, core grading.Result) map[string]any {
	breakdown := make(map[string]any, len(core.SubjectBreakdown))
	for k, v := range core.SubjectBreakdown {
		breakdown[k] = map[string]any{"earned": v.Earned, "possible": v.Possible}
	}
	return map[string]any{
		"answers":           answers,
		"subject_breakdown": breakdown,
	}
}

// answersFromSnapshot extracts the raw answers map persisted in answers_snapshot.
func answersFromSnapshot(snapshot map[string]any) map[string]any {
	if snapshot == nil {
		return nil
	}
	if a, ok := snapshot["answers"].(map[string]any); ok {
		return a
	}
	return nil
}

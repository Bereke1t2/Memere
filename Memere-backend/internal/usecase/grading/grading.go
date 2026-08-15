// Package grading is the shared, engine-neutral scoring core used by both the
// quiz and exam engines (spec §9.1 grading, §9.3 scoring). It is the ONLY place
// the correct-answer key (is_correct) is consulted, and it runs server-side
// only — the client-reported score is never an input here (Non-Negotiable:
// grading is server-side only).
//
// The core is unit-agnostic: a question's worth is "marks", which the quiz
// engine fills from a question's points and the exam engine fills from the
// per-exam exam_questions.marks. Pass/fail is intentionally NOT decided here —
// quizzes pass on a percentage threshold, exams on an absolute mark threshold —
// so the caller derives "passed" from Score/Total as appropriate.
package grading

import (
	"strings"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
)

// GradableAnswer is one answer option with its key. Server-internal only.
type GradableAnswer struct {
	ID        uuid.UUID
	Text      string
	IsCorrect bool
}

// GradableQuestion is a question to grade: its worth (marks), type, options with
// the key, and optional analytics tags.
type GradableQuestion struct {
	ID          uuid.UUID
	Text        string
	Marks       int
	Type        entity.QuestionType
	Answers     []GradableAnswer
	Subject     *string
	Topic       *string
	Explanation *string
}

// Submitted is the student's response set: selected option IDs per question
// (choice types) and/or free text per question (short_answer).
type Submitted struct {
	Selected map[uuid.UUID][]uuid.UUID
	Text     map[uuid.UUID]string
}

// QuestionOutcome is the per-question grading result for the feedback view.
type QuestionOutcome struct {
	QuestionID       uuid.UUID
	QuestionText     string
	Type             entity.QuestionType
	Subject          *string
	Topic            *string
	Correct          bool
	MarksAwarded     int
	MarksPossible    int
	SelectedAnswers  []uuid.UUID
	CorrectAnswerIDs []uuid.UUID
	Explanation      *string
	Answers          []GradableAnswer
}

// SubjectScore is the per-subject (or per-topic) tally for §9.3 analytics.
type SubjectScore struct {
	Earned   int
	Possible int
}

// Result is the authoritative server-computed outcome. Percentage is always
// Score/Total*100 (0 when Total is 0); the caller decides pass/fail.
type Result struct {
	Score            float64
	Total            int
	Percentage       float64
	Outcomes         []QuestionOutcome
	SubjectBreakdown map[string]SubjectScore
}

// Grade scores the submitted answers against the questions' keys. A question's
// marks are awarded only when the response is fully correct (§9.3). It also
// produces the per-subject/topic breakdown for analytics.
func Grade(questions []GradableQuestion, sub Submitted) Result {
	res := Result{
		Outcomes:         make([]QuestionOutcome, 0, len(questions)),
		SubjectBreakdown: map[string]SubjectScore{},
	}

	var earned, total int
	for _, q := range questions {
		total += q.Marks
		correctIDs := correctAnswerIDs(q.Answers)
		ok := isResponseCorrect(q, sub)

		awarded := 0
		if ok {
			awarded = q.Marks
			earned += q.Marks
		}

		res.Outcomes = append(res.Outcomes, QuestionOutcome{
			QuestionID:       q.ID,
			QuestionText:     q.Text,
			Type:             q.Type,
			Subject:          q.Subject,
			Topic:            q.Topic,
			Correct:          ok,
			MarksAwarded:     awarded,
			MarksPossible:    q.Marks,
			SelectedAnswers:  sub.Selected[q.ID],
			CorrectAnswerIDs: correctIDs,
			Explanation:      q.Explanation,
			Answers:          q.Answers,
		})

		if key := breakdownKey(q); key != "" {
			s := res.SubjectBreakdown[key]
			s.Earned += awarded
			s.Possible += q.Marks
			res.SubjectBreakdown[key] = s
		}
	}

	res.Score = float64(earned)
	res.Total = total
	if total > 0 {
		res.Percentage = float64(earned) / float64(total) * 100
	}
	return res
}

// isResponseCorrect grades a single question by type:
//
//   - multiple_choice / true_false: the set of selected option IDs must exactly
//     equal the set of correct option IDs (no missing, no extra) — handles both
//     single- and multi-correct questions.
//   - short_answer: the normalized free text must match the normalized text of
//     any option flagged is_correct. Normalization trims, lowercases, and
//     collapses internal whitespace. LIMITATION: exact normalized match only —
//     synonyms / partial credit are out of scope; a later phase may add manual
//     review.
func isResponseCorrect(q GradableQuestion, sub Submitted) bool {
	switch q.Type {
	case entity.QuestionShortAnswer:
		got := normalizeText(sub.Text[q.ID])
		if got == "" {
			return false
		}
		for _, a := range q.Answers {
			if a.IsCorrect && normalizeText(a.Text) == got {
				return true
			}
		}
		return false
	default: // multiple_choice, true_false
		want := map[uuid.UUID]bool{}
		for _, a := range q.Answers {
			if a.IsCorrect {
				want[a.ID] = true
			}
		}
		got := map[uuid.UUID]bool{}
		for _, id := range sub.Selected[q.ID] {
			got[id] = true
		}
		if len(want) == 0 || len(want) != len(got) {
			return false
		}
		for id := range want {
			if !got[id] {
				return false
			}
		}
		return true
	}
}

func correctAnswerIDs(answers []GradableAnswer) []uuid.UUID {
	var ids []uuid.UUID
	for _, a := range answers {
		if a.IsCorrect {
			ids = append(ids, a.ID)
		}
	}
	return ids
}

func breakdownKey(q GradableQuestion) string {
	if q.Subject != nil && *q.Subject != "" {
		return *q.Subject
	}
	if q.Topic != nil && *q.Topic != "" {
		return *q.Topic
	}
	return ""
}

// normalizeText trims, lowercases, and collapses internal whitespace runs to a
// single space, for short-answer comparison.
func normalizeText(s string) string {
	return strings.Join(strings.Fields(strings.ToLower(s)), " ")
}

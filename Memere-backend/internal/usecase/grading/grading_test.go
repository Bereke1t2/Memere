package grading

import (
	"testing"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
)

// mcQuestion builds a single-correct multiple-choice question worth marks, with
// answer A correct and B wrong. Returns the question and the two answer IDs.
func mcQuestion(marks int, subject string) (GradableQuestion, uuid.UUID, uuid.UUID) {
	a := uuid.New()
	b := uuid.New()
	s := subject
	return GradableQuestion{
		ID:      uuid.New(),
		Marks:   marks,
		Type:    entity.QuestionMultipleChoice,
		Subject: &s,
		Answers: []GradableAnswer{
			{ID: a, Text: "A", IsCorrect: true},
			{ID: b, Text: "B", IsCorrect: false},
		},
	}, a, b
}

func TestGrade_ScoresAndBreakdown(t *testing.T) {
	q1, a1, b1 := mcQuestion(5, "Math")
	q2, a2, _ := mcQuestion(3, "Science")

	sub := Submitted{
		Selected: map[uuid.UUID][]uuid.UUID{
			q1.ID: {a1}, // correct
			q2.ID: {a2}, // correct
		},
		Text: map[uuid.UUID]string{},
	}
	res := Grade([]GradableQuestion{q1, q2}, sub)
	if res.Score != 8 || res.Total != 8 || res.Percentage != 100 {
		t.Fatalf("want 8/8/100, got %v/%d/%v", res.Score, res.Total, res.Percentage)
	}
	if res.SubjectBreakdown["Math"].Earned != 5 || res.SubjectBreakdown["Science"].Earned != 3 {
		t.Errorf("subject breakdown wrong: %+v", res.SubjectBreakdown)
	}

	// Wrong answer on q1 → only q2's marks.
	sub.Selected[q1.ID] = []uuid.UUID{b1}
	res = Grade([]GradableQuestion{q1, q2}, sub)
	if res.Score != 3 {
		t.Errorf("want 3 after wrong q1, got %v", res.Score)
	}
}

func TestGrade_ShortAnswerNormalized(t *testing.T) {
	id := uuid.New()
	q := GradableQuestion{
		ID:      id,
		Marks:   2,
		Type:    entity.QuestionShortAnswer,
		Answers: []GradableAnswer{{ID: uuid.New(), Text: "Paris", IsCorrect: true}},
	}
	res := Grade([]GradableQuestion{q}, Submitted{
		Selected: map[uuid.UUID][]uuid.UUID{},
		Text:     map[uuid.UUID]string{id: "  pARIs  "},
	})
	if res.Score != 2 {
		t.Errorf("normalized short answer should be correct; got %v", res.Score)
	}
}

func TestGrade_MultiCorrectRequiresExactSet(t *testing.T) {
	id := uuid.New()
	a, b, c := uuid.New(), uuid.New(), uuid.New()
	q := GradableQuestion{
		ID:    id,
		Marks: 4,
		Type:  entity.QuestionMultipleChoice,
		Answers: []GradableAnswer{
			{ID: a, Text: "A", IsCorrect: true},
			{ID: b, Text: "B", IsCorrect: true},
			{ID: c, Text: "C", IsCorrect: false},
		},
	}
	// Partial selection (missing b) → wrong.
	res := Grade([]GradableQuestion{q}, Submitted{Selected: map[uuid.UUID][]uuid.UUID{id: {a}}})
	if res.Score != 0 {
		t.Errorf("partial multi-correct must score 0, got %v", res.Score)
	}
	// Exact set → correct.
	res = Grade([]GradableQuestion{q}, Submitted{Selected: map[uuid.UUID][]uuid.UUID{id: {a, b}}})
	if res.Score != 4 {
		t.Errorf("exact multi-correct must score 4, got %v", res.Score)
	}
	// Extra wrong option included → wrong.
	res = Grade([]GradableQuestion{q}, Submitted{Selected: map[uuid.UUID][]uuid.UUID{id: {a, b, c}}})
	if res.Score != 0 {
		t.Errorf("over-selection must score 0, got %v", res.Score)
	}
}

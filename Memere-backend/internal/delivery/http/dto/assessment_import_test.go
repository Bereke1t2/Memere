package dto

import (
	"encoding/json"
	"testing"

	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// decodeImport parses a document the same way the handler does, so these tests
// exercise the real JSON decoding (notably: every JSON number arrives as
// float64, which resolveCorrect must handle).
func decodeImport(t *testing.T, body string) *ImportRequest {
	t.Helper()
	var req ImportRequest
	if err := json.Unmarshal([]byte(body), &req); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return &req
}

func TestImportShorthandResolvesSelectorForms(t *testing.T) {
	cases := []struct {
		name     string
		selector string
		want     []bool
	}{
		{"zero based index", `1`, []bool{false, true, false, false}},
		{"option letter", `"B"`, []bool{false, true, false, false}},
		{"lowercase letter", `"b"`, []bool{false, true, false, false}},
		{"exact option text", `"Paris"`, []bool{false, false, true, false}},
		{"multi select array", `[0, 2]`, []bool{true, false, true, false}},
		{"mixed array forms", `["A", "Paris"]`, []bool{true, false, true, false}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := decodeImport(t, `{
				"version": 1,
				"kind": "quiz",
				"quiz": {"title": "Geography"},
				"questions": [{
					"text": "Capital of France?",
					"options": ["London", "Berlin", "Paris", "Madrid"],
					"correct": `+tc.selector+`
				}]
			}`)

			in, err := req.ToQuizImport()
			if err != nil {
				t.Fatalf("convert: %v", err)
			}
			answers := in.Questions[0].Answers
			if len(answers) != len(tc.want) {
				t.Fatalf("answers = %d, want %d", len(answers), len(tc.want))
			}
			for i, want := range tc.want {
				if answers[i].IsCorrect != want {
					t.Fatalf("option %d (%q) is_correct = %v, want %v", i, answers[i].Text, answers[i].IsCorrect, want)
				}
			}
		})
	}
}

// Option text is matched before the letter label, so an option whose text really
// is "A" selects itself rather than being read as "the first option".
func TestImportShorthandPrefersOptionTextOverLetterLabel(t *testing.T) {
	req := decodeImport(t, `{
		"version": 1,
		"kind": "quiz",
		"quiz": {"title": "Labels"},
		"questions": [{
			"text": "Which letter?",
			"options": ["Wrong", "A"],
			"correct": "A"
		}]
	}`)

	in, err := req.ToQuizImport()
	if err != nil {
		t.Fatalf("convert: %v", err)
	}
	answers := in.Questions[0].Answers
	if answers[0].IsCorrect {
		t.Fatal(`"A" was read as the first option's letter label instead of matching the option whose text is "A"`)
	}
	if !answers[1].IsCorrect {
		t.Fatal(`"A" did not match the option whose text is "A"`)
	}
}

// A selector that matches nothing must be a hard error. Importing the question
// with no correct answer would make every sitting of it unscoreable.
func TestImportShorthandRejectsUnresolvableSelector(t *testing.T) {
	cases := map[string]string{
		"index out of range": `9`,
		"unknown text":       `"Atlantis"`,
		"letter past end":    `"Z"`,
		"wrong type":         `true`,
		"fractional index":   `1.5`,
	}

	for name, selector := range cases {
		t.Run(name, func(t *testing.T) {
			req := decodeImport(t, `{
				"version": 1,
				"kind": "quiz",
				"quiz": {"title": "Geography"},
				"questions": [{
					"text": "Capital of France?",
					"options": ["London", "Paris"],
					"correct": `+selector+`
				}]
			}`)

			if _, err := req.ToQuizImport(); !apperror.IsCode(err, "VALIDATION_ERROR") {
				t.Fatalf("want VALIDATION_ERROR, got %v", err)
			}
		})
	}
}

func TestImportShorthandRequiresCorrectField(t *testing.T) {
	req := decodeImport(t, `{
		"version": 1,
		"kind": "quiz",
		"quiz": {"title": "Geography"},
		"questions": [{"text": "Capital?", "options": ["London", "Paris"]}]
	}`)

	if _, err := req.ToQuizImport(); !apperror.IsCode(err, "VALIDATION_ERROR") {
		t.Fatalf("options without correct: want VALIDATION_ERROR, got %v", err)
	}
}

// The two option forms are alternatives; mixing them makes the intended key
// ambiguous, so it is rejected rather than guessed.
func TestImportRejectsBothOptionForms(t *testing.T) {
	req := decodeImport(t, `{
		"version": 1,
		"kind": "quiz",
		"quiz": {"title": "Geography"},
		"questions": [{
			"text": "Capital?",
			"options": ["London", "Paris"],
			"correct": 1,
			"answers": [{"text": "London", "is_correct": true}]
		}]
	}`)

	if _, err := req.ToQuizImport(); !apperror.IsCode(err, "VALIDATION_ERROR") {
		t.Fatalf("want VALIDATION_ERROR, got %v", err)
	}
}

func TestImportRejectsNeitherOptionForm(t *testing.T) {
	req := decodeImport(t, `{
		"version": 1,
		"kind": "quiz",
		"quiz": {"title": "Geography"},
		"questions": [{"text": "Capital?"}]
	}`)

	if _, err := req.ToQuizImport(); !apperror.IsCode(err, "VALIDATION_ERROR") {
		t.Fatalf("want VALIDATION_ERROR, got %v", err)
	}
}

func TestImportExplicitAnswersPassThrough(t *testing.T) {
	req := decodeImport(t, `{
		"version": 1,
		"kind": "quiz",
		"quiz": {"title": "Explicit", "pass_percentage": 70, "time_limit_seconds": 600},
		"questions": [{
			"text": "Pick two",
			"type": "multiple_choice",
			"points": 4,
			"explanation": "  because  ",
			"subject": "Math",
			"topic": " Algebra ",
			"answers": [
				{"text": "one", "is_correct": true},
				{"text": "two", "is_correct": true},
				{"text": "three", "is_correct": false}
			]
		}]
	}`)

	in, err := req.ToQuizImport()
	if err != nil {
		t.Fatalf("convert: %v", err)
	}
	if in.Quiz.PassPercentage != 70 {
		t.Fatalf("pass_percentage = %v, want 70", in.Quiz.PassPercentage)
	}
	q := in.Questions[0]
	if q.Points != 4 {
		t.Fatalf("points = %d, want 4", q.Points)
	}
	// Optional strings are trimmed so a padded field does not store its padding.
	if q.Explanation == nil || *q.Explanation != "because" {
		t.Fatalf("explanation = %v, want trimmed %q", q.Explanation, "because")
	}
	if q.Topic == nil || *q.Topic != "Algebra" {
		t.Fatalf("topic = %v, want trimmed %q", q.Topic, "Algebra")
	}
	correct := 0
	for _, a := range q.Answers {
		if a.IsCorrect {
			correct++
		}
	}
	if correct != 2 {
		t.Fatalf("correct options = %d, want 2", correct)
	}
}

// A blank optional field reads as absent rather than storing an empty string.
func TestImportBlankOptionalFieldsBecomeNil(t *testing.T) {
	req := decodeImport(t, `{
		"version": 1,
		"kind": "quiz",
		"quiz": {"title": "Blanks"},
		"questions": [{
			"text": "Q",
			"explanation": "   ",
			"subject": "",
			"options": ["a", "b"],
			"correct": 0
		}]
	}`)

	in, err := req.ToQuizImport()
	if err != nil {
		t.Fatalf("convert: %v", err)
	}
	q := in.Questions[0]
	if q.Explanation != nil {
		t.Fatalf("explanation = %v, want nil", *q.Explanation)
	}
	if q.Subject != nil {
		t.Fatalf("subject = %v, want nil", *q.Subject)
	}
}

func TestImportRejectsWrongVersionOrKind(t *testing.T) {
	wrongVersion := decodeImport(t, `{
		"version": 99, "kind": "quiz", "quiz": {"title": "x"},
		"questions": [{"text": "Q", "options": ["a","b"], "correct": 0}]
	}`)
	if _, err := wrongVersion.ToQuizImport(); !apperror.IsCode(err, "VALIDATION_ERROR") {
		t.Fatalf("version 99: want VALIDATION_ERROR, got %v", err)
	}

	// Posting an exam document to the quiz endpoint must fail on kind rather than
	// being silently imported as the wrong type of assessment.
	wrongKind := decodeImport(t, `{
		"version": 1, "kind": "exam", "exam": {"title": "x", "subject": "Math", "grade": 12, "duration_minutes": 60},
		"questions": [{"text": "Q", "options": ["a","b"], "correct": 0}]
	}`)
	if _, err := wrongKind.ToQuizImport(); !apperror.IsCode(err, "VALIDATION_ERROR") {
		t.Fatalf("exam document on the quiz endpoint: want VALIDATION_ERROR, got %v", err)
	}
}

// "marks" and "points" are the same column; an exam file reads naturally with
// marks, and marks wins when both appear.
func TestImportExamUsesMarksForTotal(t *testing.T) {
	req := decodeImport(t, `{
		"version": 1,
		"kind": "exam",
		"exam": {"title": "Mock", "subject": "Physics", "grade": 12, "duration_minutes": 90, "pass_marks": 5},
		"questions": [
			{"text": "Q1", "marks": 3, "options": ["a","b"], "correct": 0},
			{"text": "Q2", "points": 1, "marks": 7, "options": ["a","b"], "correct": 1}
		]
	}`)

	in, err := req.ToExamImport()
	if err != nil {
		t.Fatalf("convert: %v", err)
	}
	if in.Questions[0].Marks != 3 {
		t.Fatalf("q1 marks = %d, want 3", in.Questions[0].Marks)
	}
	if in.Questions[1].Marks != 7 {
		t.Fatalf("q2 marks = %d, want 7 (marks must win over points)", in.Questions[1].Marks)
	}
}

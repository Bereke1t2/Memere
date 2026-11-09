package dto

import (
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/exam"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/quiz"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// This file defines the bulk-import wire format: one JSON document describing a
// whole quiz or exam, its questions, their options, which options are correct,
// and the post-submission explanations. The admin panel uploads it directly;
// authors can also hand-write or generate it.
//
// Two ways to express a question's options are accepted, and exactly one must be
// used per question:
//
//   "answers": [{"text": "...", "is_correct": true}, ...]   — explicit, canonical
//   "options": ["...", "..."], "correct": <selector>        — concise shorthand
//
// The shorthand exists because hand-authoring is_correct on every option is
// tedious and error-prone for a large file. Its selector resolves deliberately
// strictly (see resolveCorrect): a selector that matches nothing is a hard
// validation error naming the question, never a question silently imported with
// no correct answer.

// ImportFormatVersion is the only supported value of the top-level "version"
// field. It is required, so a future incompatible format can be rejected by
// version rather than by confusing field-level errors.
const ImportFormatVersion = 1

// ImportAnswerRequest is one explicitly-specified option.
type ImportAnswerRequest struct {
	Text      string `json:"text"`
	IsCorrect bool   `json:"is_correct"`
}

// ImportQuestionRequest is one question in an import document. Points and Marks
// are accepted interchangeably so the same question block works in a quiz file
// and an exam file; Marks wins when both are present.
type ImportQuestionRequest struct {
	Text        string                `json:"text"`
	Type        string                `json:"type"`
	Points      *int                  `json:"points"`
	Marks       *int                  `json:"marks"`
	Explanation *string               `json:"explanation"`
	Subject     *string               `json:"subject"`
	Topic       *string               `json:"topic"`
	Answers     []ImportAnswerRequest `json:"answers"`
	Options     []string              `json:"options"`
	// Correct is the shorthand selector. It accepts a single value or an array of
	// them, where each value is a 0-based index (number), the exact text of an
	// option (string), or a letter label such as "A"/"b" (string).
	Correct any `json:"correct"`
}

// ImportQuizHeaderRequest is the quiz-level block of a quiz import document.
type ImportQuizHeaderRequest struct {
	CourseID           *uuid.UUID `json:"course_id"`
	LessonID           *uuid.UUID `json:"lesson_id"`
	Title              string     `json:"title"`
	TimeLimitSeconds   *int       `json:"time_limit_seconds"`
	PassPercentage     float64    `json:"pass_percentage"`
	RandomizeQuestions bool       `json:"randomize_questions"`
	MaxAttempts        *int       `json:"max_attempts"`
}

// ImportExamHeaderRequest is the exam-level block of an exam import document.
type ImportExamHeaderRequest struct {
	CourseID        *uuid.UUID `json:"course_id"`
	Title           string     `json:"title"`
	Subject         string     `json:"subject"`
	Grade           int        `json:"grade"`
	DurationMinutes int        `json:"duration_minutes"`
	PassMarks       int        `json:"pass_marks"`
	Instructions    *string    `json:"instructions"`
}

// ImportRequest is the whole upload: a version, which kind of assessment it
// describes, the matching header block, and the question set.
type ImportRequest struct {
	Version   int                      `json:"version"`
	Kind      string                   `json:"kind"`
	Quiz      *ImportQuizHeaderRequest `json:"quiz"`
	Exam      *ImportExamHeaderRequest `json:"exam"`
	Questions []ImportQuestionRequest  `json:"questions"`
}

// ImportResponse reports what an import created.
type ImportResponse struct {
	Kind          string    `json:"kind"`
	ID            uuid.UUID `json:"id"`
	Title         string    `json:"title"`
	QuestionCount int       `json:"question_count"`
	TotalMarks    *int      `json:"total_marks,omitempty"`
}

// ToQuizImport validates the document as a quiz import and converts it to the
// usecase input. Errors are returned as a validation map keyed by the offending
// path (e.g. "questions[3].correct") so an author can fix their file directly.
func (r *ImportRequest) ToQuizImport() (quiz.ImportQuizInput, error) {
	var out quiz.ImportQuizInput

	details := r.validateEnvelope("quiz")
	if r.Quiz == nil {
		details["quiz"] = "is required when kind is \"quiz\""
	} else if strings.TrimSpace(r.Quiz.Title) == "" {
		details["quiz.title"] = "is required"
	}
	if len(details) > 0 {
		return out, apperror.Validation(details, nil)
	}

	questions := make([]quiz.QuestionInput, len(r.Questions))
	for i, q := range r.Questions {
		answers, qType, points, err := q.normalize(i)
		if err != nil {
			return out, err
		}
		questions[i] = quiz.QuestionInput{
			Text:        strings.TrimSpace(q.Text),
			Type:        qType,
			Points:      points,
			Explanation: trimmedPtr(q.Explanation),
			OrderIndex:  i,
			Subject:     trimmedPtr(q.Subject),
			Topic:       trimmedPtr(q.Topic),
			Answers:     answers,
		}
	}

	return quiz.ImportQuizInput{
		Quiz: quiz.CreateQuizInput{
			CourseID:           r.Quiz.CourseID,
			LessonID:           r.Quiz.LessonID,
			Title:              strings.TrimSpace(r.Quiz.Title),
			TimeLimitSeconds:   r.Quiz.TimeLimitSeconds,
			PassPercentage:     r.Quiz.PassPercentage,
			RandomizeQuestions: r.Quiz.RandomizeQuestions,
			MaxAttempts:        r.Quiz.MaxAttempts,
		},
		Questions: questions,
	}, nil
}

// ToExamImport validates the document as an exam import and converts it to the
// usecase input.
func (r *ImportRequest) ToExamImport() (exam.ImportExamInput, error) {
	var out exam.ImportExamInput

	details := r.validateEnvelope("exam")
	if r.Exam == nil {
		details["exam"] = "is required when kind is \"exam\""
	} else {
		if strings.TrimSpace(r.Exam.Title) == "" {
			details["exam.title"] = "is required"
		}
		if strings.TrimSpace(r.Exam.Subject) == "" {
			details["exam.subject"] = "is required"
		}
		if r.Exam.DurationMinutes <= 0 {
			details["exam.duration_minutes"] = "must be positive"
		}
	}
	if len(details) > 0 {
		return out, apperror.Validation(details, nil)
	}

	questions := make([]exam.AddInlineQuestionInput, len(r.Questions))
	for i, q := range r.Questions {
		answers, qType, marks, err := q.normalize(i)
		if err != nil {
			return out, err
		}
		opts := make([]exam.InlineAnswerInput, len(answers))
		for j, a := range answers {
			opts[j] = exam.InlineAnswerInput{Text: a.Text, IsCorrect: a.IsCorrect, OrderIndex: j}
		}
		questions[i] = exam.AddInlineQuestionInput{
			Text:        strings.TrimSpace(q.Text),
			Type:        qType,
			Marks:       marks,
			Explanation: trimmedPtr(q.Explanation),
			OrderIndex:  i,
			Subject:     trimmedPtr(q.Subject),
			Topic:       trimmedPtr(q.Topic),
			Answers:     opts,
		}
	}

	return exam.ImportExamInput{
		Exam: exam.CreateExamInput{
			CourseID:        r.Exam.CourseID,
			Title:           strings.TrimSpace(r.Exam.Title),
			Subject:         strings.TrimSpace(r.Exam.Subject),
			Grade:           r.Exam.Grade,
			DurationMinutes: r.Exam.DurationMinutes,
			PassMarks:       r.Exam.PassMarks,
			Instructions:    trimmedPtr(r.Exam.Instructions),
		},
		Questions: questions,
	}, nil
}

// validateEnvelope checks the version/kind/questions shell common to both kinds.
func (r *ImportRequest) validateEnvelope(want string) map[string]any {
	details := map[string]any{}
	if r.Version != ImportFormatVersion {
		details["version"] = fmt.Sprintf("must be %d", ImportFormatVersion)
	}
	if r.Kind != want {
		details["kind"] = fmt.Sprintf("must be %q for this endpoint", want)
	}
	if len(r.Questions) == 0 {
		details["questions"] = "must contain at least one question"
	}
	return details
}

// normalize turns one import question into its canonical answer set, type and
// point value. It is the single place the two option forms converge, so the
// quiz and exam paths cannot drift in how they read an answer key.
func (q *ImportQuestionRequest) normalize(index int) ([]quiz.AnswerInput, entity.QuestionType, int, error) {
	path := fmt.Sprintf("questions[%d]", index)

	qType := entity.QuestionType(strings.TrimSpace(q.Type))
	if qType == "" {
		qType = entity.QuestionMultipleChoice
	}
	if !qType.Valid() {
		return nil, "", 0, apperror.Validation(map[string]any{
			path + ".type": "must be multiple_choice, true_false, or short_answer",
		}, nil)
	}

	// Marks and points are the same column; Marks wins so an exam file reads
	// naturally. Default 1 so a minimal question block is valid.
	points := 1
	if q.Points != nil {
		points = *q.Points
	}
	if q.Marks != nil {
		points = *q.Marks
	}

	hasAnswers := len(q.Answers) > 0
	hasOptions := len(q.Options) > 0
	switch {
	case hasAnswers && hasOptions:
		return nil, "", 0, apperror.Validation(map[string]any{
			path: "use either \"answers\" or \"options\"+\"correct\", not both",
		}, nil)
	case !hasAnswers && !hasOptions:
		return nil, "", 0, apperror.Validation(map[string]any{
			path: "must provide \"answers\" or \"options\"+\"correct\"",
		}, nil)
	}

	if hasAnswers {
		answers := make([]quiz.AnswerInput, len(q.Answers))
		for i, a := range q.Answers {
			answers[i] = quiz.AnswerInput{
				Text:       strings.TrimSpace(a.Text),
				IsCorrect:  a.IsCorrect,
				OrderIndex: i,
			}
		}
		return answers, qType, points, nil
	}

	correct, err := resolveCorrect(q.Options, q.Correct, path)
	if err != nil {
		return nil, "", 0, err
	}
	answers := make([]quiz.AnswerInput, len(q.Options))
	for i, text := range q.Options {
		answers[i] = quiz.AnswerInput{
			Text:       strings.TrimSpace(text),
			IsCorrect:  correct[i],
			OrderIndex: i,
		}
	}
	return answers, qType, points, nil
}

// resolveCorrect maps a "correct" selector onto the options slice, returning a
// per-option correctness mask.
//
// Each selector value resolves in a fixed order — 0-based index for numbers,
// then exact (case-insensitive, trimmed) option text, then a single-letter
// label — and the ordering matters: text is tried before the letter label so an
// option whose text really is "A" selects itself rather than being read as
// "the first option". A value matching nothing, or a selector resolving to no
// options at all, is an error, because silently importing a question with no
// correct answer would make every sitting of it unscoreable.
func resolveCorrect(options []string, selector any, path string) (map[int]bool, error) {
	fail := func(msg string) error {
		return apperror.Validation(map[string]any{path + ".correct": msg}, nil)
	}
	if selector == nil {
		return nil, fail("is required when using \"options\"")
	}

	values, ok := selector.([]any)
	if !ok {
		values = []any{selector}
	}

	correct := map[int]bool{}
	for _, v := range values {
		switch val := v.(type) {
		case float64: // encoding/json decodes every JSON number as float64
			idx := int(val)
			if float64(idx) != val || idx < 0 || idx >= len(options) {
				return nil, fail(fmt.Sprintf("index %v is out of range for %d options", val, len(options)))
			}
			correct[idx] = true
		case string:
			needle := strings.TrimSpace(val)
			matched := -1
			for i, opt := range options {
				if strings.EqualFold(strings.TrimSpace(opt), needle) {
					matched = i
					break
				}
			}
			if matched < 0 && len(needle) == 1 {
				// Letter label: "A"/"a" -> 0, "B" -> 1, ...
				if idx := int(strings.ToUpper(needle)[0]) - 'A'; idx >= 0 && idx < len(options) {
					matched = idx
				}
			}
			if matched < 0 {
				return nil, fail(fmt.Sprintf("%q matches no option text and is not a valid option letter", val))
			}
			correct[matched] = true
		default:
			return nil, fail("must be an option index, an option's exact text, or an option letter")
		}
	}
	if len(correct) == 0 {
		return nil, fail("must select at least one option")
	}
	return correct, nil
}

// trimmedPtr normalizes an optional string: whitespace is trimmed and an empty
// result becomes nil, so a blank field in a file reads as "absent" rather than
// storing an empty explanation or tag.
func trimmedPtr(s *string) *string {
	if s == nil {
		return nil
	}
	t := strings.TrimSpace(*s)
	if t == "" {
		return nil
	}
	return &t
}

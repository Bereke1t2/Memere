package quiz

import (
	"context"
	"errors"
	"testing"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

func importQuestion(text string, correct int) QuestionInput {
	answers := []AnswerInput{
		{Text: "A", OrderIndex: 0},
		{Text: "B", OrderIndex: 1},
		{Text: "C", OrderIndex: 2},
	}
	answers[correct].IsCorrect = true
	return QuestionInput{
		Text:        text,
		Type:        entity.QuestionMultipleChoice,
		Points:      1,
		Explanation: ptrStr("because"),
		Answers:     answers,
	}
}

func TestImportQuizCreatesWholeTree(t *testing.T) {
	h := newHarness(t)
	ctx := context.Background()

	created, count, err := h.svc.ImportQuiz(ctx, admin(), ImportQuizInput{
		Quiz: CreateQuizInput{Title: "Imported", PassPercentage: 60},
		Questions: []QuestionInput{
			importQuestion("Q1", 0),
			importQuestion("Q2", 1),
			importQuestion("Q3", 2),
		},
	})
	if err != nil {
		t.Fatalf("import: %v", err)
	}
	if count != 3 {
		t.Fatalf("count = %d, want 3", count)
	}

	view, err := h.svc.GetQuizForAuthoring(ctx, admin(), created.ID)
	if err != nil {
		t.Fatalf("authoring view: %v", err)
	}
	if len(view.Questions) != 3 {
		t.Fatalf("questions = %d, want 3", len(view.Questions))
	}
	// Position in the file is the authoritative order.
	for i, q := range view.Questions {
		if q.OrderIndex != i {
			t.Fatalf("question %d has order_index %d, want %d", i, q.OrderIndex, i)
		}
	}
	if view.Questions[0].Explanation == nil || *view.Questions[0].Explanation != "because" {
		t.Fatal("explanation was not imported")
	}
}

// Validation runs over the whole file before anything is written, and the error
// names the offending question so an author can fix their source.
func TestImportQuizReportsOffendingQuestionIndex(t *testing.T) {
	h := newHarness(t)
	ctx := context.Background()

	bad := importQuestion("Q2", 0)
	for i := range bad.Answers {
		bad.Answers[i].IsCorrect = false // no key -> ungradable
	}

	_, _, err := h.svc.ImportQuiz(ctx, admin(), ImportQuizInput{
		Quiz:      CreateQuizInput{Title: "Imported", PassPercentage: 60},
		Questions: []QuestionInput{importQuestion("Q1", 0), bad},
	})
	if !apperror.IsCode(err, "VALIDATION_ERROR") {
		t.Fatalf("want VALIDATION_ERROR, got %v", err)
	}
	var ae *apperror.AppError
	if !errors.As(err, &ae) {
		t.Fatalf("error is not an AppError: %v", err)
	}
	if _, ok := ae.Details["questions[1]"]; !ok {
		t.Fatalf("details %v do not name the offending question index", ae.Details)
	}

	// Nothing was written: a half-imported quiz would look complete to a student
	// while silently mis-scoring.
	quizzes, _, _, err := h.svc.ListQuizzes(ctx, admin(), repository.QuizFilter{}, nil, 20)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(quizzes) != 0 {
		t.Fatalf("a rejected import created %d quizzes, want 0", len(quizzes))
	}
}

func TestImportQuizRejectsEmptyQuestionSet(t *testing.T) {
	h := newHarness(t)
	_, _, err := h.svc.ImportQuiz(context.Background(), admin(), ImportQuizInput{
		Quiz: CreateQuizInput{Title: "Empty", PassPercentage: 60},
	})
	if !apperror.IsCode(err, "VALIDATION_ERROR") {
		t.Fatalf("want VALIDATION_ERROR, got %v", err)
	}
}

// A teacher may import under their own course but not as a standalone quiz.
func TestImportQuizHonoursOwnership(t *testing.T) {
	h := newHarness(t)
	ctx := context.Background()

	_, _, err := h.svc.ImportQuiz(ctx, teacher(), ImportQuizInput{
		Quiz:      CreateQuizInput{Title: "Imported", PassPercentage: 60},
		Questions: []QuestionInput{importQuestion("Q1", 0)},
	})
	if !apperror.IsCode(err, "FORBIDDEN") {
		t.Fatalf("teacher importing a standalone quiz: want FORBIDDEN, got %v", err)
	}
}

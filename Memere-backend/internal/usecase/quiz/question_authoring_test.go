package quiz

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// The authoring view is the one read path that must expose is_correct, otherwise
// an author cannot see the key they are editing.
func TestGetQuizForAuthoringIncludesAnswerKey(t *testing.T) {
	h := newHarness(t)
	quiz := h.seedQuiz(nil, nil, false)

	view, err := h.svc.GetQuizForAuthoring(context.Background(), admin(), quiz.ID)
	if err != nil {
		t.Fatalf("authoring view: %v", err)
	}
	if len(view.Questions) == 0 {
		t.Fatal("authoring view returned no questions")
	}
	correct := 0
	for _, q := range view.Questions {
		for _, a := range q.Answers {
			if a.IsCorrect {
				correct++
			}
		}
	}
	if correct == 0 {
		t.Fatal("authoring view exposed no correct answer, so the key cannot be edited")
	}
}

// ...but only to someone who may author it. A student must never reach it.
func TestGetQuizForAuthoringRejectsStudent(t *testing.T) {
	h := newHarness(t)
	quiz := h.seedQuiz(nil, nil, false)

	if _, err := h.svc.GetQuizForAuthoring(context.Background(), student(), quiz.ID); !apperror.IsCode(err, "FORBIDDEN") {
		t.Fatalf("student reading the authoring view: want FORBIDDEN, got %v", err)
	}
}

func TestUpdateQuestionReplacesAnswerSet(t *testing.T) {
	h := newHarness(t)
	ctx := context.Background()
	quiz := h.seedQuiz(nil, nil, false)

	view, err := h.svc.GetQuizForAuthoring(ctx, admin(), quiz.ID)
	if err != nil {
		t.Fatalf("authoring view: %v", err)
	}
	target := view.Questions[0]

	newText := "What is 3+3?"
	updated, err := h.svc.UpdateQuestion(ctx, admin(), target.ID, UpdateQuestionInput{
		Text: &newText,
		Answers: []AnswerInput{
			{Text: "6", IsCorrect: true, OrderIndex: 0},
			{Text: "7", IsCorrect: false, OrderIndex: 1},
			{Text: "8", IsCorrect: false, OrderIndex: 2},
		},
	})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if updated.Text != newText {
		t.Fatalf("text = %q, want %q", updated.Text, newText)
	}
	if len(updated.Answers) != 3 {
		t.Fatalf("answers = %d, want 3 (the set should be replaced, not merged)", len(updated.Answers))
	}

	// Re-read to confirm the old options are gone rather than lingering with a
	// stale is_correct.
	reread, err := h.svc.GetQuizForAuthoring(ctx, admin(), quiz.ID)
	if err != nil {
		t.Fatalf("re-read: %v", err)
	}
	for _, q := range reread.Questions {
		if q.ID != target.ID {
			continue
		}
		if len(q.Answers) != 3 {
			t.Fatalf("persisted answers = %d, want 3", len(q.Answers))
		}
		for _, a := range q.Answers {
			if a.Text == "4" || a.Text == "5" {
				t.Fatalf("stale option %q survived the replacement", a.Text)
			}
		}
	}
}

// A question with no correct option cannot be graded, so the update is refused.
func TestUpdateQuestionRejectsKeylessAnswerSet(t *testing.T) {
	h := newHarness(t)
	ctx := context.Background()
	quiz := h.seedQuiz(nil, nil, false)

	view, _ := h.svc.GetQuizForAuthoring(ctx, admin(), quiz.ID)
	_, err := h.svc.UpdateQuestion(ctx, admin(), view.Questions[0].ID, UpdateQuestionInput{
		Answers: []AnswerInput{
			{Text: "6", IsCorrect: false, OrderIndex: 0},
			{Text: "7", IsCorrect: false, OrderIndex: 1},
		},
	})
	if !apperror.IsCode(err, "VALIDATION_ERROR") {
		t.Fatalf("answer set with no correct option: want VALIDATION_ERROR, got %v", err)
	}
}

// Omitting answers must still validate against the stored set, so a type change
// cannot land on an option set that can no longer grade it.
func TestUpdateQuestionValidatesStoredAnswersWhenOmitted(t *testing.T) {
	h := newHarness(t)
	ctx := context.Background()
	quiz := h.seedQuiz(nil, nil, false)

	view, _ := h.svc.GetQuizForAuthoring(ctx, admin(), quiz.ID)
	target := view.Questions[0]

	zero := 0
	if _, err := h.svc.UpdateQuestion(ctx, admin(), target.ID, UpdateQuestionInput{Points: &zero}); !apperror.IsCode(err, "VALIDATION_ERROR") {
		t.Fatalf("zero points: want VALIDATION_ERROR, got %v", err)
	}

	// A valid partial update leaves the existing key intact.
	five := 5
	updated, err := h.svc.UpdateQuestion(ctx, admin(), target.ID, UpdateQuestionInput{Points: &five})
	if err != nil {
		t.Fatalf("partial update: %v", err)
	}
	if len(updated.Answers) != len(target.Answers) {
		t.Fatalf("answers = %d, want the stored %d preserved", len(updated.Answers), len(target.Answers))
	}
}

func TestDeleteQuestion(t *testing.T) {
	h := newHarness(t)
	ctx := context.Background()
	quiz := h.seedQuiz(nil, nil, false)

	before, _ := h.svc.GetQuizForAuthoring(ctx, admin(), quiz.ID)
	if err := h.svc.DeleteQuestion(ctx, admin(), before.Questions[0].ID); err != nil {
		t.Fatalf("delete: %v", err)
	}
	after, _ := h.svc.GetQuizForAuthoring(ctx, admin(), quiz.ID)
	if len(after.Questions) != len(before.Questions)-1 {
		t.Fatalf("questions = %d, want %d", len(after.Questions), len(before.Questions)-1)
	}
}

func TestDeleteQuestionRejectsStudent(t *testing.T) {
	h := newHarness(t)
	ctx := context.Background()
	quiz := h.seedQuiz(nil, nil, false)

	view, _ := h.svc.GetQuizForAuthoring(ctx, admin(), quiz.ID)
	if err := h.svc.DeleteQuestion(ctx, student(), view.Questions[0].ID); !apperror.IsCode(err, "FORBIDDEN") {
		t.Fatalf("student deleting a question: want FORBIDDEN, got %v", err)
	}
}

// An exam bank question has no quiz_id, so there is no course to derive
// ownership from. It must be rejected here rather than silently permitted.
func TestUpdateQuestionRejectsExamBankQuestion(t *testing.T) {
	h := newHarness(t)
	ctx := context.Background()

	// quiz_id uuid.Nil is how the repository represents a NULL quiz_id.
	bank := &entity.Question{Text: "Bank question", Type: entity.QuestionMultipleChoice, Points: 1}
	h.quizzes.addQuestion(uuid.Nil, bank, []*entity.Answer{
		{Text: "yes", IsCorrect: true, OrderIndex: 0},
		{Text: "no", IsCorrect: false, OrderIndex: 1},
	})

	text := "Rewritten"
	_, err := h.svc.UpdateQuestion(ctx, admin(), bank.ID, UpdateQuestionInput{Text: &text})
	if !apperror.IsCode(err, "BAD_REQUEST") {
		t.Fatalf("editing an exam bank question through the quiz path: want BAD_REQUEST, got %v", err)
	}
}

func TestReorderQuestionsRequiresCompleteList(t *testing.T) {
	h := newHarness(t)
	ctx := context.Background()
	quiz := h.seedQuiz(nil, nil, false)

	view, _ := h.svc.GetQuizForAuthoring(ctx, admin(), quiz.ID)
	ids := make([]uuid.UUID, len(view.Questions))
	for i, q := range view.Questions {
		ids[i] = q.ID
	}

	// A partial list would silently drop a question from the order.
	if err := h.svc.ReorderQuestions(ctx, admin(), quiz.ID, ids[:1]); !apperror.IsCode(err, "VALIDATION_ERROR") {
		t.Fatalf("partial reorder list: want VALIDATION_ERROR, got %v", err)
	}
	// So would a duplicate.
	dup := []uuid.UUID{ids[0], ids[0]}
	if err := h.svc.ReorderQuestions(ctx, admin(), quiz.ID, dup); !apperror.IsCode(err, "VALIDATION_ERROR") {
		t.Fatalf("duplicate reorder list: want VALIDATION_ERROR, got %v", err)
	}

	// Reversing the full list is accepted and takes effect.
	reversed := []uuid.UUID{ids[1], ids[0]}
	if err := h.svc.ReorderQuestions(ctx, admin(), quiz.ID, reversed); err != nil {
		t.Fatalf("reorder: %v", err)
	}
	after, _ := h.svc.GetQuizForAuthoring(ctx, admin(), quiz.ID)
	if after.Questions[0].ID != ids[1] {
		t.Fatal("reorder did not move the second question to the front")
	}
}

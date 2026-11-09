package quiz

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

func admin() *Actor   { return &Actor{UserID: uuid.New(), Role: entity.RoleAdmin} }
func teacher() *Actor { return &Actor{UserID: uuid.New(), Role: entity.RoleTeacher} }

// A standalone quiz has no course owner to delegate authoring to, so only an
// admin may create one.
func TestCreateStandaloneQuizRequiresAdmin(t *testing.T) {
	h := newHarness(t)
	in := CreateQuizInput{Title: "Practice", PassPercentage: 50}

	if _, err := h.svc.CreateQuiz(context.Background(), teacher(), in); !apperror.IsCode(err, "FORBIDDEN") {
		t.Fatalf("teacher creating a standalone quiz: want forbidden, got %v", err)
	}

	created, err := h.svc.CreateQuiz(context.Background(), admin(), in)
	if err != nil {
		t.Fatalf("admin creating a standalone quiz: %v", err)
	}
	if created.CourseID != nil {
		t.Fatalf("course_id = %v, want nil", created.CourseID)
	}
	if created.IsPublished {
		t.Fatal("a newly created quiz must be unpublished so it is unreachable while being authored")
	}
}

// The publish flag is what exposes a course-less quiz. Until it is set the quiz
// must stay invisible to students, including guests.
func TestStandaloneQuizAccessGate(t *testing.T) {
	h := newHarness(t)
	ctx := context.Background()

	draft, err := h.svc.CreateQuiz(ctx, admin(), CreateQuizInput{Title: "Draft", PassPercentage: 50})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	if _, err := h.svc.GetQuizForStudent(ctx, student(), draft.ID); !apperror.IsCode(err, "FORBIDDEN") {
		t.Fatalf("student reading an unpublished standalone quiz: want forbidden, got %v", err)
	}
	if _, err := h.svc.GetQuizForStudent(ctx, nil, draft.ID); !apperror.IsCode(err, "FORBIDDEN") {
		t.Fatalf("guest reading an unpublished standalone quiz: want forbidden, got %v", err)
	}
	if _, err := h.svc.GetQuizForStudent(ctx, admin(), draft.ID); err != nil {
		t.Fatalf("admin reading its own unpublished standalone quiz: %v", err)
	}

	published := true
	if _, err := h.svc.UpdateQuiz(ctx, admin(), draft.ID, UpdateQuizInput{IsPublished: &published}); err != nil {
		t.Fatalf("publish: %v", err)
	}
	if _, err := h.svc.GetQuizForStudent(ctx, student(), draft.ID); err != nil {
		t.Fatalf("student reading a published standalone quiz: %v", err)
	}
	if _, err := h.svc.GetQuizForStudent(ctx, nil, draft.ID); err != nil {
		t.Fatalf("guest reading a published standalone quiz: %v", err)
	}
}

// Publishing a standalone quiz must not become a way around a course paywall:
// a course-attached quiz stays gated on course access regardless of the flag.
func TestPublishDoesNotBypassCoursePaywall(t *testing.T) {
	h := newHarness(t)
	ctx := context.Background()

	quiz := h.seedQuizUnder(h.paidCourse(), nil, nil, false)
	quiz.IsPublished = true
	if err := h.quizzes.Update(ctx, quiz); err != nil {
		t.Fatalf("update: %v", err)
	}

	// The course paywall answers NOT_ENROLLED; the point is that is_published did
	// not short-circuit it.
	if _, err := h.svc.GetQuizForStudent(ctx, student(), quiz.ID); !apperror.IsCode(err, "NOT_ENROLLED") {
		t.Fatalf("unenrolled student on a published paid-course quiz: want NOT_ENROLLED, got %v", err)
	}
}

// A standalone quiz belongs to no teacher, so a teacher must not be able to edit
// one even though teachers may author course quizzes.
func TestStandaloneQuizEditRequiresAdmin(t *testing.T) {
	h := newHarness(t)
	ctx := context.Background()

	created, err := h.svc.CreateQuiz(ctx, admin(), CreateQuizInput{Title: "Practice", PassPercentage: 50})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	title := "Hijacked"
	if _, err := h.svc.UpdateQuiz(ctx, teacher(), created.ID, UpdateQuizInput{Title: &title}); !apperror.IsCode(err, "FORBIDDEN") {
		t.Fatalf("teacher editing a standalone quiz: want forbidden, got %v", err)
	}
}

func TestListQuizzesFiltersStandalone(t *testing.T) {
	h := newHarness(t)
	ctx := context.Background()

	if _, err := h.svc.CreateQuiz(ctx, admin(), CreateQuizInput{Title: "Alone", PassPercentage: 50}); err != nil {
		t.Fatalf("create standalone: %v", err)
	}
	h.seedQuiz(nil, nil, false)

	standalone := true
	got, counts, _, err := h.svc.ListQuizzes(ctx, admin(), repository.QuizFilter{Standalone: &standalone}, nil, 20)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(got) != 1 || got[0].Title != "Alone" {
		t.Fatalf("standalone filter returned %d quizzes, want just the standalone one", len(got))
	}
	if len(counts) != len(got) {
		t.Fatalf("counts length %d, want %d", len(counts), len(got))
	}

	// A teacher must never see another teacher's quizzes through this listing,
	// nor the admin-only standalone ones.
	got, _, _, err = h.svc.ListQuizzes(ctx, teacher(), repository.QuizFilter{}, nil, 20)
	if err != nil {
		t.Fatalf("list as teacher: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("teacher saw %d quizzes they do not own, want 0", len(got))
	}
}

func TestListQuizzesRejectsStudent(t *testing.T) {
	h := newHarness(t)
	if _, _, _, err := h.svc.ListQuizzes(context.Background(), student(), repository.QuizFilter{}, nil, 20); !apperror.IsCode(err, "FORBIDDEN") {
		t.Fatalf("student listing quizzes for authoring: want forbidden, got %v", err)
	}
}

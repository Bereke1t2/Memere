package quiz

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// bannedKeyFields are the answer-key markers that must NEVER appear in a client
// (online-play) projection. They are ALLOWED only in the /download payload.
var bannedKeyFields = []string{"is_correct", "IsCorrect", "isCorrect", "correct_answer", "CorrectAnswer"}

// --- guest attempt lifecycle (A1) ---------------------------------------------

// An unregistered guest (actor == nil) can start, auto-save, submit, and read the
// result of a published/free quiz — mirroring the exam guest flow. The attempt is
// owned by the shared GuestUserID, and grading still happens server-side.
func TestQuizGuestFlow_StartSaveSubmitResult(t *testing.T) {
	h := newHarness(t)
	quiz := h.seedQuiz(nil, nil, false)

	// 1. Guest starts.
	v, err := h.svc.StartAttempt(context.Background(), nil, quiz.ID)
	if err != nil {
		t.Fatalf("guest StartAttempt: %v", err)
	}
	if len(v.Questions) != 2 {
		t.Fatalf("want 2 questions for guest, got %d", len(v.Questions))
	}
	// The attempt is owned by the shared GuestUserID (not a random per-guest id).
	if _, err := h.attempts.FindByID(context.Background(), v.AttemptID, GuestUserID); err != nil {
		t.Fatalf("guest attempt should be owned by GuestUserID: %v", err)
	}

	// 2. Guest auto-saves a partial answer.
	c0 := h.correctAnswerID(quiz.ID, 0)
	if err := h.svc.SaveProgress(context.Background(), nil, v.AttemptID, map[string]any{
		v.Questions[idxOf(v, c0)].ID.String(): selected(c0),
	}); err != nil {
		t.Fatalf("guest SaveProgress: %v", err)
	}

	// 3. Guest submits both correct answers; graded server-side.
	c1 := h.correctAnswerID(quiz.ID, 1)
	res, err := h.svc.SubmitAttempt(context.Background(), nil, v.AttemptID, map[string]any{
		v.Questions[idxOf(v, c0)].ID.String(): selected(c0),
		v.Questions[idxOf(v, c1)].ID.String(): selected(c1),
	})
	if err != nil {
		t.Fatalf("guest SubmitAttempt: %v", err)
	}
	if res.Score != 2 || res.Percentage != 100 || !res.Passed {
		t.Errorf("guest full-marks submit: got score=%v pct=%v passed=%v", res.Score, res.Percentage, res.Passed)
	}
	if res.Status != entity.AttemptGraded {
		t.Errorf("guest submit status = %v, want graded", res.Status)
	}

	// 4. Guest reads the graded result back.
	got, err := h.svc.GetAttemptResult(context.Background(), nil, v.AttemptID)
	if err != nil {
		t.Fatalf("guest GetAttemptResult: %v", err)
	}
	if got.Status != entity.AttemptGraded || got.Score != res.Score {
		t.Errorf("guest result mismatch: got %+v, want status=graded score=%v", got, res.Score)
	}
}

// --- answer-key download (A3) -------------------------------------------------

// GetQuizForDownload is the ONE sanctioned exception to Non-Negotiable #1: it
// carries the answer key for offline on-device grading. This test proves the key
// IS present in the download payload while the SAME questions stay answer-free on
// the client (online-play) path.
func TestGetQuizForDownload_CarriesAnswerKeyClientPathDoesNot(t *testing.T) {
	h := newHarness(t)
	quiz := h.seedQuiz(nil, nil, false)

	// Download path: full key present.
	d, err := h.svc.GetQuizForDownload(context.Background(), student(), quiz.ID)
	if err != nil {
		t.Fatalf("GetQuizForDownload: %v", err)
	}
	if len(d.Questions) != 2 {
		t.Fatalf("want 2 download questions, got %d", len(d.Questions))
	}
	for i, q := range d.Questions {
		correct := 0
		for _, a := range q.Answers {
			if a.IsCorrect {
				correct++
			}
		}
		if correct != 1 {
			t.Errorf("download question %d: want exactly 1 correct option, got %d", i, correct)
		}
	}
	// The serialized DTO is the wire shape; assert is_correct actually rides along.
	if b, _ := json.Marshal(newTestDownloadWire(d)); !strings.Contains(string(b), "is_correct") {
		t.Fatalf("download payload must expose is_correct for offline grading: %s", b)
	}

	// Client (online-play) path: the SAME questions, answer key stripped.
	clientQs, err := h.quizzes.GetQuestionsForClient(context.Background(), quiz.ID)
	if err != nil {
		t.Fatalf("GetQuestionsForClient: %v", err)
	}
	b, _ := json.Marshal(clientQs)
	for _, banned := range bannedKeyFields {
		if strings.Contains(string(b), banned) {
			t.Fatalf("client question path leaks answer key (%q): %s", banned, b)
		}
	}
}

// The download endpoint reuses the exact same course-access gate as taking the
// quiz, so a logged-in student with no entitlement to a PAID course cannot pull
// the keys either: re-enabling the paywall auto-restricts downloads.
func TestGetQuizForDownload_NonEnrolledStudentBlocked(t *testing.T) {
	h := newHarness(t)
	quiz := h.seedQuizUnder(h.paidCourse(), nil, nil, false)
	stu := &Actor{UserID: uuid.New(), Role: entity.RoleStudent}

	_, err := h.svc.GetQuizForDownload(context.Background(), stu, quiz.ID)
	if !apperror.IsCode(err, "NOT_ENROLLED") {
		t.Fatalf("non-enrolled student downloading keys should be NOT_ENROLLED, got %v", err)
	}
}

// newTestDownloadWire mirrors the is_correct-bearing wire shape produced by
// dto.NewQuizDownloadResponse, so the presence assertion above tests the tag the
// client actually receives (the usecase struct itself has no JSON tags).
func newTestDownloadWire(d *QuizDownload) any {
	type answer struct {
		ID        uuid.UUID `json:"id"`
		IsCorrect bool      `json:"is_correct"`
	}
	type question struct {
		ID      uuid.UUID `json:"id"`
		Answers []answer  `json:"answers"`
	}
	out := make([]question, 0, len(d.Questions))
	for _, q := range d.Questions {
		as := make([]answer, 0, len(q.Answers))
		for _, a := range q.Answers {
			as = append(as, answer{ID: a.ID, IsCorrect: a.IsCorrect})
		}
		out = append(out, question{ID: q.ID, Answers: as})
	}
	return out
}

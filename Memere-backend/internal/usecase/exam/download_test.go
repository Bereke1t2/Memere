package exam

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// GetExamForDownload is the sanctioned Non-Negotiable #1 exception for exams: it
// carries the answer key + per-exam marks for offline on-device grading. This
// proves the key IS present in the download payload while the SAME questions stay
// answer-free on the client (online-sitting) path.
func TestGetExamForDownload_CarriesAnswerKeyClientPathDoesNot(t *testing.T) {
	h := newHarness(t)
	exam := h.seedPublishedExam(t)

	// Download path: full key + marks present.
	d, err := h.svc.GetExamForDownload(context.Background(), student(), exam.ID)
	if err != nil {
		t.Fatalf("GetExamForDownload: %v", err)
	}
	if len(d.Questions) != 2 {
		t.Fatalf("want 2 download questions, got %d", len(d.Questions))
	}
	totalMarks := 0
	for i, q := range d.Questions {
		totalMarks += q.Marks
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
	if totalMarks != 8 {
		t.Errorf("download must carry per-exam marks (5+3); summed %d", totalMarks)
	}

	// Client (online-sitting) path: the SAME questions, answer key stripped.
	clientQs, err := h.exams.GetQuestionsForClient(context.Background(), exam.ID)
	if err != nil {
		t.Fatalf("GetQuestionsForClient: %v", err)
	}
	b, _ := json.Marshal(clientQs)
	for _, banned := range []string{"is_correct", "IsCorrect", "isCorrect", "correct_answer", "CorrectAnswer"} {
		if strings.Contains(string(b), banned) {
			t.Fatalf("client exam-question path leaks answer key (%q): %s", banned, b)
		}
	}
}

// An unpublished (draft) exam withholds its keys from a non-owner: the download
// endpoint is gated by the same assertExamAccess as StartExam, so a student
// cannot download the key for content that isn't published to them.
func TestGetExamForDownload_UnpublishedBlocked(t *testing.T) {
	h := newHarness(t)
	// A standalone (no-course) draft exam authored by an admin.
	exam, err := h.svc.CreateExam(context.Background(), admin(), CreateExamInput{
		Title:           "Draft",
		Subject:         "Math",
		Grade:           12,
		DurationMinutes: 60,
		PassMarks:       1,
	})
	if err != nil {
		t.Fatalf("CreateExam: %v", err)
	}

	_, err = h.svc.GetExamForDownload(context.Background(), student(), exam.ID)
	if !apperror.IsCode(err, "FORBIDDEN") {
		t.Fatalf("student downloading a draft exam's keys should be FORBIDDEN, got %v", err)
	}
}

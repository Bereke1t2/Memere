package exam

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/access"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

func init() {
	access.DisableEnrollmentCheck = false
}

type harness struct {
	svc      *Service
	exams    *fakeExamRepo
	attempts *fakeExamAttemptRepo
	courses  *fakeCourseRepo
	state    *fakeState
	enroll   *fakeEnrollRepo
	clock    time.Time
}

func newHarness(t *testing.T) *harness {
	t.Helper()
	exams := newFakeExamRepo()
	attempts := newFakeExamAttemptRepo()
	courses := newFakeCourseRepo()
	state := newFakeState()
	enroll := newFakeEnrollRepo()
	// Course-linked exams gate through the real access.Service so the FullAccess
	// policy is exercised end to end, not mocked.
	accessSvc := access.NewService(enroll, newFakeSubRepo(), courses, nil)
	svc := NewService(exams, attempts, courses, state, nil, fakeTxManager{}, accessSvc, nil)

	h := &harness{svc: svc, exams: exams, attempts: attempts, courses: courses, state: state, enroll: enroll}
	h.clock = time.Date(2026, 6, 8, 12, 0, 0, 0, time.UTC)
	svc.now = func() time.Time { return h.clock }
	// Let the fake attempt repo derive deadlines from exam durations, mirroring the
	// real query's JOIN on exams.duration_minutes.
	expiryLookup = func(examID uuid.UUID) (time.Duration, bool) {
		e, ok := exams.exams[examID]
		if !ok {
			return 0, false
		}
		return time.Duration(e.DurationMinutes) * time.Minute, true
	}
	return h
}

func (h *harness) advance(d time.Duration) { h.clock = h.clock.Add(d) }

func student() *Actor { return &Actor{UserID: uuid.New(), Role: entity.RoleStudent} }
func admin() *Actor   { return &Actor{UserID: uuid.New(), Role: entity.RoleAdmin} }

func ptrStr(s string) *string { return &s }

// seedPublishedExam builds a published, 60-minute exam with two MC questions
// (5 marks + 3 marks) under a free course (so the FullAccess gate admits any
// student: the lifecycle tests exercise the engine, not the paywall) and returns
// it.
func (h *harness) seedPublishedExam(t *testing.T) *entity.Exam {
	t.Helper()
	course := &entity.Course{ID: uuid.New(), TeacherID: uuid.New(), IsPublished: true, IsFree: true}
	return h.seedExamForCourse(t, course)
}

// seedExamForCourse seeds the standard two-question exam under an explicit course
// (used by the access-gate tests, which need a paid course). The course's
// teacher authors and publishes the exam.
func (h *harness) seedExamForCourse(t *testing.T, course *entity.Course) *entity.Exam {
	t.Helper()
	owner := &Actor{UserID: course.TeacherID, Role: entity.RoleTeacher}
	h.courses.add(course)

	exam, err := h.svc.CreateExam(context.Background(), owner, CreateExamInput{
		CourseID:        &course.ID,
		Title:           "Midterm",
		Subject:         "Math",
		Grade:           12,
		DurationMinutes: 60,
		PassMarks:       5,
	})
	if err != nil {
		t.Fatalf("CreateExam: %v", err)
	}

	q1 := &entity.Question{Text: "2+2?", Type: entity.QuestionMultipleChoice, OrderIndex: 0, Subject: ptrStr("Arithmetic")}
	h.exams.addBankQuestion(q1, []*entity.Answer{
		{Text: "4", IsCorrect: true, OrderIndex: 0},
		{Text: "5", IsCorrect: false, OrderIndex: 1},
	})
	q2 := &entity.Question{Text: "Capital of France?", Type: entity.QuestionShortAnswer, OrderIndex: 1, Subject: ptrStr("Geography")}
	h.exams.addBankQuestion(q2, []*entity.Answer{{Text: "Paris", IsCorrect: true}})

	if _, err := h.svc.AddExamQuestion(context.Background(), owner, exam.ID, q1.ID, 5, 0); err != nil {
		t.Fatalf("AddExamQuestion q1: %v", err)
	}
	if _, err := h.svc.AddExamQuestion(context.Background(), owner, exam.ID, q2.ID, 3, 1); err != nil {
		t.Fatalf("AddExamQuestion q2: %v", err)
	}
	if _, err := h.svc.PublishExam(context.Background(), owner, exam.ID); err != nil {
		t.Fatalf("PublishExam: %v", err)
	}
	reloaded, _ := h.exams.FindByID(context.Background(), exam.ID)
	return reloaded
}

// --- total_marks recompute ----------------------------------------------------

func TestAddExamQuestion_RecomputesTotalMarks(t *testing.T) {
	h := newHarness(t)
	exam := h.seedPublishedExam(t)
	if exam.TotalMarks != 8 {
		t.Errorf("total_marks = %d, want 8 (5+3)", exam.TotalMarks)
	}
}

// --- answer-key leak ----------------------------------------------------------

func TestStartExam_ClientViewHasNoAnswerKey(t *testing.T) {
	h := newHarness(t)
	exam := h.seedPublishedExam(t)

	view, err := h.svc.StartExam(context.Background(), student(), exam.ID)
	if err != nil {
		t.Fatalf("StartExam: %v", err)
	}
	b, _ := json.Marshal(view)
	for _, banned := range []string{"is_correct", "IsCorrect", "correct_answer", "CorrectAnswer"} {
		if strings.Contains(string(b), banned) {
			t.Fatalf("exam client view leaks answer key (%q): %s", banned, b)
		}
	}
	if view.TotalMarks != 8 || len(view.Questions) != 2 {
		t.Errorf("view total=%d questions=%d, want 8/2", view.TotalMarks, len(view.Questions))
	}
	if view.RemainingSeconds == nil || *view.RemainingSeconds != 3600 {
		t.Errorf("remaining seconds = %v, want 3600", view.RemainingSeconds)
	}
}

// --- idempotent resume --------------------------------------------------------

func TestStartExam_ResumesInProgress(t *testing.T) {
	h := newHarness(t)
	exam := h.seedPublishedExam(t)
	stu := student()

	v1, err := h.svc.StartExam(context.Background(), stu, exam.ID)
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	h.advance(10 * time.Minute)
	v2, err := h.svc.StartExam(context.Background(), stu, exam.ID)
	if err != nil {
		t.Fatalf("resume: %v", err)
	}
	if v1.AttemptID != v2.AttemptID {
		t.Errorf("resume created a new attempt (%v != %v)", v1.AttemptID, v2.AttemptID)
	}
	if v2.RemainingSeconds == nil || *v2.RemainingSeconds != 3000 {
		t.Errorf("remaining after 10m = %v, want 3000", v2.RemainingSeconds)
	}
}

// --- submit before expiry → submitted→graded ---------------------------------

func TestSubmitExam_BeforeExpiryGradesServerSide(t *testing.T) {
	h := newHarness(t)
	exam := h.seedPublishedExam(t)
	stu := student()
	v, _ := h.svc.StartExam(context.Background(), stu, exam.ID)

	q1 := h.questionByText(v, "2+2?")
	c1 := h.correctOption(exam.ID, q1.QuestionID)
	q2 := h.questionByText(v, "Capital of France?")

	res, err := h.svc.SubmitExam(context.Background(), stu, v.AttemptID, map[string]any{
		q1.QuestionID.String(): map[string]any{"selected": []any{c1.String()}},
		q2.QuestionID.String(): map[string]any{"text": " paris "},
		"score":                999, // client-supplied score must be ignored
	})
	if err != nil {
		t.Fatalf("submit: %v", err)
	}
	if res.Status != entity.AttemptGraded {
		t.Errorf("status = %v, want graded", res.Status)
	}
	if res.Score != 8 || res.TotalMarks != 8 || !res.Passed {
		t.Errorf("want full marks 8/8 passed; got score=%v total=%v passed=%v", res.Score, res.TotalMarks, res.Passed)
	}
	if len(res.SubjectBreakdown) != 2 {
		t.Errorf("subject breakdown size = %d, want 2", len(res.SubjectBreakdown))
	}
}

// --- submit after expiry → expired→graded ------------------------------------

func TestSubmitExam_AfterExpiryGradesAsExpired(t *testing.T) {
	h := newHarness(t)
	exam := h.seedPublishedExam(t)
	stu := student()
	v, _ := h.svc.StartExam(context.Background(), stu, exam.ID)

	q1 := h.questionByText(v, "2+2?")
	c1 := h.correctOption(exam.ID, q1.QuestionID)
	if err := h.svc.SaveExamProgress(context.Background(), stu, v.AttemptID, map[string]any{
		q1.QuestionID.String(): map[string]any{"selected": []any{c1.String()}},
	}); err != nil {
		t.Fatalf("save: %v", err)
	}

	h.advance(2 * time.Hour) // past the 60-minute window

	res, err := h.svc.SubmitExam(context.Background(), stu, v.AttemptID, map[string]any{})
	if err != nil {
		t.Fatalf("submit after expiry: %v", err)
	}
	if res.Status != entity.AttemptGraded {
		t.Errorf("final status = %v, want graded", res.Status)
	}
	if res.Score != 5 {
		t.Errorf("auto-saved partial answer should score 5; got %v", res.Score)
	}
	// The persisted attempt passed through EXPIRED on its way to graded.
	stored, _ := h.attempts.FindByID(context.Background(), v.AttemptID, stu.UserID)
	if stored.Status != entity.AttemptGraded {
		t.Errorf("stored status = %v, want graded", stored.Status)
	}
}

// --- state machine: illegal transition ---------------------------------------

func TestSubmitExam_ReSubmitIsIdempotent(t *testing.T) {
	h := newHarness(t)
	exam := h.seedPublishedExam(t)
	stu := student()
	v, _ := h.svc.StartExam(context.Background(), stu, exam.ID)
	first, err := h.svc.SubmitExam(context.Background(), stu, v.AttemptID, map[string]any{})
	if err != nil {
		t.Fatalf("first submit: %v", err)
	}
	// A second submit on an already-graded attempt is benign: it returns the same
	// graded result rather than erroring or re-grading.
	second, err := h.svc.SubmitExam(context.Background(), stu, v.AttemptID, map[string]any{})
	if err != nil {
		t.Fatalf("re-submit should be idempotent, got %v", err)
	}
	if second.Status != entity.AttemptGraded || second.Score != first.Score {
		t.Errorf("re-submit result differs: first=%v second=%v", first, second)
	}
}

func TestGetExamResult_InProgressIsInvalidState(t *testing.T) {
	h := newHarness(t)
	exam := h.seedPublishedExam(t)
	stu := student()
	v, _ := h.svc.StartExam(context.Background(), stu, exam.ID)

	_, err := h.svc.GetExamResult(context.Background(), stu, v.AttemptID)
	if !apperror.IsCode(err, "INVALID_ATTEMPT_STATE") {
		t.Fatalf("result while in progress should be INVALID_ATTEMPT_STATE, got %v", err)
	}
}

// --- lazy expiry on GetExamResult --------------------------------------------

func TestGetExamResult_PastExpiryFinalizes(t *testing.T) {
	h := newHarness(t)
	exam := h.seedPublishedExam(t)
	stu := student()
	v, _ := h.svc.StartExam(context.Background(), stu, exam.ID)

	h.advance(2 * time.Hour)
	res, err := h.svc.GetExamResult(context.Background(), stu, v.AttemptID)
	if err != nil {
		t.Fatalf("get result past expiry: %v", err)
	}
	if res.Status != entity.AttemptGraded {
		t.Errorf("status = %v, want graded after lazy expiry", res.Status)
	}
}

// --- IDOR ---------------------------------------------------------------------

func TestSubmitExam_OtherStudentForbidden(t *testing.T) {
	h := newHarness(t)
	exam := h.seedPublishedExam(t)
	owner := student()
	v, _ := h.svc.StartExam(context.Background(), owner, exam.ID)

	_, err := h.svc.SubmitExam(context.Background(), student(), v.AttemptID, map[string]any{})
	if !apperror.IsNotFound(err) {
		t.Fatalf("another student's attempt must not be found (IDOR), got %v", err)
	}
}

// --- pass/fail vs pass_marks --------------------------------------------------

func TestSubmitExam_FailBelowPassMarks(t *testing.T) {
	h := newHarness(t)
	exam := h.seedPublishedExam(t)
	stu := student()
	v, _ := h.svc.StartExam(context.Background(), stu, exam.ID)

	// Answer only the 3-mark short-answer correctly → score 3 < pass_marks 5.
	q2 := h.questionByText(v, "Capital of France?")
	res, err := h.svc.SubmitExam(context.Background(), stu, v.AttemptID, map[string]any{
		q2.QuestionID.String(): map[string]any{"text": "Paris"},
	})
	if err != nil {
		t.Fatalf("submit: %v", err)
	}
	if res.Score != 3 || res.Passed {
		t.Errorf("score 3 < pass_marks 5 should fail; got score=%v passed=%v", res.Score, res.Passed)
	}
}

// --- authoring authz ----------------------------------------------------------

func TestCreateExam_NonTeacherForbidden(t *testing.T) {
	h := newHarness(t)
	_, err := h.svc.CreateExam(context.Background(), student(), CreateExamInput{
		Title: "x", Subject: "Math", Grade: 12, DurationMinutes: 60, PassMarks: 1,
	})
	if !apperror.IsCode(err, "FORBIDDEN") {
		t.Fatalf("student creating an exam should be FORBIDDEN, got %v", err)
	}
}

func TestListExams_StudentSeesPublishedOnly(t *testing.T) {
	h := newHarness(t)
	_ = h.seedPublishedExam(t) // published
	// An unpublished exam by an admin.
	if _, err := h.svc.CreateExam(context.Background(), admin(), CreateExamInput{
		Title: "Hidden", Subject: "Sci", Grade: 12, DurationMinutes: 30, PassMarks: 1,
	}); err != nil {
		t.Fatalf("create hidden: %v", err)
	}
	got, _, err := h.svc.ListExams(context.Background(), student(), repository.ExamFilter{}, nil, 50)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	for _, e := range got {
		if !e.IsPublished {
			t.Errorf("student listing returned an unpublished exam %q", e.Title)
		}
	}
}

// --- sweeper + race guard -----------------------------------------------------

func TestSweepExpired_GradesAbandonedExam(t *testing.T) {
	h := newHarness(t)
	exam := h.seedPublishedExam(t)
	stu := student()
	v, _ := h.svc.StartExam(context.Background(), stu, exam.ID)

	q1 := h.questionByText(v, "2+2?")
	c1 := h.correctOption(exam.ID, q1.QuestionID)
	_ = h.svc.SaveExamProgress(context.Background(), stu, v.AttemptID, map[string]any{
		q1.QuestionID.String(): map[string]any{"selected": []any{c1.String()}},
	})

	h.advance(2 * time.Hour) // abandoned past the 60-minute window

	n, err := h.svc.SweepExpired(context.Background(), h.clock, 100)
	if err != nil {
		t.Fatalf("SweepExpired: %v", err)
	}
	if n != 1 {
		t.Fatalf("want 1 swept, got %d", n)
	}
	stored, _ := h.attempts.FindByID(context.Background(), v.AttemptID, stu.UserID)
	if stored.Status != entity.AttemptGraded {
		t.Errorf("status = %v, want graded", stored.Status)
	}
	if stored.Score == nil || *stored.Score != 5 {
		t.Errorf("auto-saved answer should score 5; got %v", stored.Score)
	}
}

func TestSubmitVsSweep_ExamGradesExactlyOnce(t *testing.T) {
	h := newHarness(t)
	exam := h.seedPublishedExam(t)
	stu := student()
	v, _ := h.svc.StartExam(context.Background(), stu, exam.ID)
	h.advance(2 * time.Hour)

	n, err := h.svc.SweepExpired(context.Background(), h.clock, 100)
	if err != nil || n != 1 {
		t.Fatalf("sweep: n=%d err=%v", n, err)
	}
	// Late submit returns the graded result, does not re-grade or error.
	res, err := h.svc.SubmitExam(context.Background(), stu, v.AttemptID, map[string]any{})
	if err != nil {
		t.Fatalf("late submit should return graded result; got %v", err)
	}
	if res.Status != entity.AttemptGraded {
		t.Errorf("late submit status = %v, want graded", res.Status)
	}
}

// --- helpers ------------------------------------------------------------------

func (h *harness) questionByText(v *ExamAttemptClientView, text string) QuestionClientView {
	for _, q := range v.Questions {
		if q.Text == text {
			return q
		}
	}
	return QuestionClientView{}
}

func (h *harness) correctOption(examID, questionID uuid.UUID) uuid.UUID {
	tree, _ := h.exams.GetExamWithQuestions(context.Background(), examID)
	for _, qa := range tree.Questions {
		if qa.Question.ID == questionID {
			for _, a := range qa.Answers {
				if a.IsCorrect {
					return a.ID
				}
			}
		}
	}
	return uuid.Nil
}

// --- access gate (Phase 4) ----------------------------------------------------

// A student is blocked from an unpublished exam unless they are the teacher owner or admin.
func TestStartExam_UnpublishedExamBlocked(t *testing.T) {
	h := newHarness(t)
	course := &entity.Course{ID: uuid.New(), TeacherID: uuid.New(), IsPublished: true, IsFree: false}
	owner := &Actor{UserID: course.TeacherID, Role: entity.RoleTeacher}
	h.courses.add(course)

	exam, err := h.svc.CreateExam(context.Background(), owner, CreateExamInput{
		CourseID:        &course.ID,
		Title:           "Draft Exam",
		Subject:         "Math",
		Grade:           12,
		DurationMinutes: 60,
		PassMarks:       5,
	})
	if err != nil {
		t.Fatalf("CreateExam: %v", err)
	}
	stu := &Actor{UserID: uuid.New(), Role: entity.RoleStudent}

	_, err = h.svc.StartExam(context.Background(), stu, exam.ID)
	if err == nil {
		t.Fatalf("student should be blocked from unpublished exam")
	}
}

// Published exams are accessible to students and guests.
func TestStartExam_PublishedExamAllowed(t *testing.T) {
	h := newHarness(t)
	course := &entity.Course{ID: uuid.New(), TeacherID: uuid.New(), IsPublished: true, IsFree: false}
	exam := h.seedExamForCourse(t, course)
	stu := &Actor{UserID: uuid.New(), Role: entity.RoleStudent}

	view, err := h.svc.StartExam(context.Background(), stu, exam.ID)
	if err != nil {
		t.Fatalf("student should start published exam: %v", err)
	}
	if view == nil || len(view.Questions) == 0 {
		t.Fatalf("expected a started attempt with questions, got %+v", view)
	}
}

// Unregistered/anonymous guest (actor == nil) can start, save, submit, and view results for a published exam.
func TestStartExam_GuestUnregisteredAllowed(t *testing.T) {
	h := newHarness(t)
	exam := h.seedPublishedExam(t)

	// 1. Guest starts exam
	view, err := h.svc.StartExam(context.Background(), nil, exam.ID)
	if err != nil {
		t.Fatalf("guest should be allowed to start published exam: %v", err)
	}
	if view == nil || len(view.Questions) != 2 {
		t.Fatalf("expected 2 questions for guest, got %+v", view)
	}
	if view.AttemptID == uuid.Nil {
		t.Fatal("expected non-nil attempt ID for guest")
	}

	q1 := view.Questions[0]
	corr := h.correctOption(exam.ID, q1.QuestionID)

	// 2. Guest auto-saves progress
	answers := map[string]any{q1.QuestionID.String(): map[string]any{"selected": []any{corr.String()}}}
	if err := h.svc.SaveExamProgress(context.Background(), nil, view.AttemptID, answers); err != nil {
		t.Fatalf("guest should be allowed to save progress: %v", err)
	}

	// 3. Guest submits exam
	res, err := h.svc.SubmitExam(context.Background(), nil, view.AttemptID, answers)
	if err != nil {
		t.Fatalf("guest should be allowed to submit exam: %v", err)
	}
	if res == nil || res.Score <= 0 {
		t.Fatalf("expected graded score > 0 for guest, got %+v", res)
	}

	// 4. Guest views result
	resView, err := h.svc.GetExamResult(context.Background(), nil, view.AttemptID)
	if err != nil {
		t.Fatalf("guest should be allowed to view result: %v", err)
	}
	if resView == nil || resView.Score != res.Score {
		t.Fatalf("expected result view matching submitted result, got %+v", resView)
	}
}

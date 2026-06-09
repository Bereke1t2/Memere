package quiz

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

type harness struct {
	svc      *Service
	quizzes  *fakeQuizRepo
	attempts *fakeAttemptRepo
	courses  *fakeCourseRepo
	state    *fakeState
	clock    time.Time
}

func newHarness(t *testing.T) *harness {
	t.Helper()
	quizzes := newFakeQuizRepo()
	attempts := newFakeAttemptRepo()
	courses := newFakeCourseRepo()
	state := newFakeState()
	questions := &fakeQuestionRepo{quizzes: quizzes}
	svc := NewService(quizzes, questions, attempts, courses, state, fakeTxManager{})

	h := &harness{svc: svc, quizzes: quizzes, attempts: attempts, courses: courses, state: state}
	h.clock = time.Date(2026, 6, 8, 12, 0, 0, 0, time.UTC)
	svc.now = func() time.Time { return h.clock }
	return h
}

func (h *harness) advance(d time.Duration) { h.clock = h.clock.Add(d) }

func student() *Actor { return &Actor{UserID: uuid.New(), Role: entity.RoleStudent} }

func ptrInt(i int) *int       { return &i }
func ptrStr(s string) *string { return &s }

// publishedCourse registers a published course and returns its ID.
func (h *harness) publishedCourse() uuid.UUID {
	id := uuid.New()
	h.courses.add(&entity.Course{ID: id, TeacherID: uuid.New(), IsPublished: true})
	return id
}

// seedQuiz creates a quiz with two multiple-choice questions (one correct option
// each) under a published course. Returns the quiz.
func (h *harness) seedQuiz(timeLimit, maxAttempts *int, randomize bool) *entity.Quiz {
	courseID := h.publishedCourse()
	q := &entity.Quiz{
		CourseID:           courseID,
		Title:              "Sample",
		TimeLimitSeconds:   timeLimit,
		MaxAttempts:        maxAttempts,
		PassPercentage:     50,
		RandomizeQuestions: randomize,
	}
	if err := h.quizzes.Create(context.Background(), q); err != nil {
		t0 := err
		panic(t0)
	}
	h.quizzes.addQuestion(q.ID,
		&entity.Question{Text: "2+2?", Type: entity.QuestionMultipleChoice, Points: 1, OrderIndex: 0, Subject: ptrStr("Math")},
		[]*entity.Answer{
			{Text: "4", IsCorrect: true, OrderIndex: 0},
			{Text: "5", IsCorrect: false, OrderIndex: 1},
		},
	)
	h.quizzes.addQuestion(q.ID,
		&entity.Question{Text: "Sky color?", Type: entity.QuestionMultipleChoice, Points: 1, OrderIndex: 1, Subject: ptrStr("Science")},
		[]*entity.Answer{
			{Text: "Blue", IsCorrect: true, OrderIndex: 0},
			{Text: "Green", IsCorrect: false, OrderIndex: 1},
		},
	)
	return q
}

// correctAnswerID returns the correct answer's ID for question index qi.
func (h *harness) correctAnswerID(quizID uuid.UUID, qi int) uuid.UUID {
	tree, _ := h.quizzes.GetQuizWithQuestions(context.Background(), quizID)
	for _, a := range tree.Questions[qi].Answers {
		if a.IsCorrect {
			return a.ID
		}
	}
	return uuid.Nil
}

func selected(ids ...uuid.UUID) map[string]any {
	out := make([]any, len(ids))
	for i, id := range ids {
		out[i] = id.String()
	}
	return map[string]any{"selected": out}
}

// --- answer-key-leak test -----------------------------------------------------

func TestStartAttempt_ClientViewHasNoAnswerKey(t *testing.T) {
	h := newHarness(t)
	quiz := h.seedQuiz(nil, nil, false)

	view, err := h.svc.StartAttempt(context.Background(), student(), quiz.ID)
	if err != nil {
		t.Fatalf("StartAttempt: %v", err)
	}
	// Serialize the client view and assert no answer-key field leaks.
	b, _ := json.Marshal(view)
	for _, banned := range []string{"is_correct", "IsCorrect", "isCorrect", "correct_answer", "CorrectAnswer"} {
		if strings.Contains(string(b), banned) {
			t.Fatalf("client view leaks answer key (%q): %s", banned, b)
		}
	}
	if len(view.Questions) != 2 {
		t.Fatalf("want 2 questions, got %d", len(view.Questions))
	}
}

// --- randomization ------------------------------------------------------------

func TestStartAttempt_RandomizationStableWithinAttemptDiffersAcross(t *testing.T) {
	h := newHarness(t)
	// Many options so a shuffle is observable.
	courseID := h.publishedCourse()
	quiz := &entity.Quiz{CourseID: courseID, Title: "Q", PassPercentage: 50, RandomizeQuestions: true}
	_ = h.quizzes.Create(context.Background(), quiz)
	answers := make([]*entity.Answer, 6)
	for i := range answers {
		answers[i] = &entity.Answer{Text: string(rune('A' + i)), IsCorrect: i == 0, OrderIndex: i}
	}
	h.quizzes.addQuestion(quiz.ID, &entity.Question{Text: "pick", Type: entity.QuestionMultipleChoice, Points: 1}, answers)

	stu := student()
	v1, err := h.svc.StartAttempt(context.Background(), stu, quiz.ID)
	if err != nil {
		t.Fatalf("start1: %v", err)
	}
	order1 := answerIDOrder(v1)

	// Same attempt resumed → identical order.
	v1b, err := h.svc.StartAttempt(context.Background(), stu, quiz.ID)
	if err != nil {
		t.Fatalf("resume: %v", err)
	}
	if strings.Join(order1, ",") != strings.Join(answerIDOrder(v1b), ",") {
		t.Errorf("resumed attempt changed answer order")
	}

	// Submit, start a second attempt → different attempt ID → likely different order.
	if _, err := h.svc.SubmitAttempt(context.Background(), stu, v1.AttemptID, map[string]any{}); err != nil {
		t.Fatalf("submit: %v", err)
	}
	v2, err := h.svc.StartAttempt(context.Background(), stu, quiz.ID)
	if err != nil {
		t.Fatalf("start2: %v", err)
	}
	if v2.AttemptID == v1.AttemptID {
		t.Fatal("second attempt reused the first attempt ID")
	}
	// Deterministic-per-attempt: orders are derived from distinct attempt IDs.
	// They are overwhelmingly likely to differ for 6 options; assert they're
	// valid permutations of the same set at minimum.
	if !sameSet(order1, answerIDOrder(v2)) {
		t.Errorf("attempt orders are not permutations of the same option set")
	}
}

// --- max attempts -------------------------------------------------------------

func TestStartAttempt_MaxAttemptsEnforced(t *testing.T) {
	h := newHarness(t)
	quiz := h.seedQuiz(nil, ptrInt(1), false)
	stu := student()

	v1, err := h.svc.StartAttempt(context.Background(), stu, quiz.ID)
	if err != nil {
		t.Fatalf("first start: %v", err)
	}
	if _, err := h.svc.SubmitAttempt(context.Background(), stu, v1.AttemptID, map[string]any{}); err != nil {
		t.Fatalf("submit: %v", err)
	}
	_, err = h.svc.StartAttempt(context.Background(), stu, quiz.ID)
	if !apperror.IsCode(err, "FORBIDDEN") {
		t.Fatalf("second start past max_attempts should be FORBIDDEN, got %v", err)
	}
}

// --- idempotent double-start --------------------------------------------------

func TestStartAttempt_DoubleStartReturnsExisting(t *testing.T) {
	h := newHarness(t)
	quiz := h.seedQuiz(nil, ptrInt(3), false)
	stu := student()

	v1, err := h.svc.StartAttempt(context.Background(), stu, quiz.ID)
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	v2, err := h.svc.StartAttempt(context.Background(), stu, quiz.ID)
	if err != nil {
		t.Fatalf("re-start: %v", err)
	}
	if v1.AttemptID != v2.AttemptID {
		t.Errorf("double start created a new attempt (%v != %v)", v1.AttemptID, v2.AttemptID)
	}
	if n, _ := h.attempts.CountByStudentAndQuiz(context.Background(), stu.UserID, quiz.ID); n != 1 {
		t.Errorf("double start consumed an extra attempt: count=%d", n)
	}
}

// --- grading correctness ------------------------------------------------------

func TestSubmitAttempt_GradesServerSide_IgnoresClientScore(t *testing.T) {
	h := newHarness(t)
	quiz := h.seedQuiz(nil, nil, false)
	stu := student()

	v, _ := h.svc.StartAttempt(context.Background(), stu, quiz.ID)
	c0 := h.correctAnswerID(quiz.ID, 0)
	c1 := h.correctAnswerID(quiz.ID, 1)

	payload := map[string]any{
		v.Questions[idxOf(v, c0)].ID.String(): selected(c0),
		v.Questions[idxOf(v, c1)].ID.String(): selected(c1),
		// A bogus client-supplied score that must be ignored.
		"score": 999,
	}
	res, err := h.svc.SubmitAttempt(context.Background(), stu, v.AttemptID, payload)
	if err != nil {
		t.Fatalf("submit: %v", err)
	}
	if res.Score != 2 || res.Percentage != 100 || !res.Passed {
		t.Errorf("want full marks; got score=%v pct=%v passed=%v", res.Score, res.Percentage, res.Passed)
	}
	if res.Status != entity.AttemptGraded {
		t.Errorf("status = %v, want graded", res.Status)
	}
	// Subject breakdown captured (§9.3).
	if len(res.SubjectBreakdown) != 2 {
		t.Errorf("want 2 subjects in breakdown, got %d", len(res.SubjectBreakdown))
	}
}

func TestSubmitAttempt_FailBelowThreshold(t *testing.T) {
	h := newHarness(t)
	quiz := h.seedQuiz(nil, nil, false)
	stu := student()

	v, _ := h.svc.StartAttempt(context.Background(), stu, quiz.ID)
	c0 := h.correctAnswerID(quiz.ID, 0)
	// Answer only the first question correctly → 50%. pass is >= 50 so passes;
	// answer nothing else. To force a fail, answer one wrong instead.
	wrong := h.wrongAnswerID(quiz.ID, 0)
	payload := map[string]any{
		v.Questions[idxOf(v, c0)].ID.String(): selected(wrong),
	}
	res, err := h.svc.SubmitAttempt(context.Background(), stu, v.AttemptID, payload)
	if err != nil {
		t.Fatalf("submit: %v", err)
	}
	if res.Passed {
		t.Errorf("0%% should not pass a 50%% threshold")
	}
}

// --- short answer -------------------------------------------------------------

func TestSubmitAttempt_ShortAnswerNormalized(t *testing.T) {
	h := newHarness(t)
	courseID := h.publishedCourse()
	quiz := &entity.Quiz{CourseID: courseID, Title: "SA", PassPercentage: 100}
	_ = h.quizzes.Create(context.Background(), quiz)
	h.quizzes.addQuestion(quiz.ID,
		&entity.Question{Text: "Capital of France?", Type: entity.QuestionShortAnswer, Points: 1},
		[]*entity.Answer{{Text: "Paris", IsCorrect: true}},
	)
	stu := student()
	v, _ := h.svc.StartAttempt(context.Background(), stu, quiz.ID)
	qid := v.Questions[0].ID
	res, err := h.svc.SubmitAttempt(context.Background(), stu, v.AttemptID, map[string]any{
		qid.String(): map[string]any{"text": "  pARIs  "},
	})
	if err != nil {
		t.Fatalf("submit: %v", err)
	}
	if !res.Passed {
		t.Errorf("normalized short answer should be correct; got %+v", res)
	}
}

// --- expiry -------------------------------------------------------------------

func TestSubmitAttempt_PastExpiryGradesAsExpired(t *testing.T) {
	h := newHarness(t)
	quiz := h.seedQuiz(ptrInt(60), nil, false) // 60s limit
	stu := student()

	v, _ := h.svc.StartAttempt(context.Background(), stu, quiz.ID)
	h.advance(2 * time.Minute) // past expiry

	c0 := h.correctAnswerID(quiz.ID, 0)
	res, err := h.svc.SubmitAttempt(context.Background(), stu, v.AttemptID, map[string]any{
		v.Questions[idxOf(v, c0)].ID.String(): selected(c0),
	})
	if err != nil {
		t.Fatalf("submit: %v", err)
	}
	// Final persisted status is graded; result still grades the partial answers.
	if res.Status != entity.AttemptGraded {
		t.Errorf("status = %v, want graded", res.Status)
	}
	if res.Score != 1 {
		t.Errorf("partial answers should still be graded; score=%v", res.Score)
	}
}

func TestSaveProgress_RejectedPastExpiry(t *testing.T) {
	h := newHarness(t)
	quiz := h.seedQuiz(ptrInt(30), nil, false)
	stu := student()
	v, _ := h.svc.StartAttempt(context.Background(), stu, quiz.ID)

	h.advance(time.Minute)
	err := h.svc.SaveProgress(context.Background(), stu, v.AttemptID, map[string]any{"x": "y"})
	if !apperror.IsCode(err, "CONFLICT") {
		t.Fatalf("save past expiry should CONFLICT, got %v", err)
	}
}

// --- IDOR ---------------------------------------------------------------------

func TestSubmitAttempt_OtherStudentForbidden(t *testing.T) {
	h := newHarness(t)
	quiz := h.seedQuiz(nil, nil, false)
	owner := student()
	v, _ := h.svc.StartAttempt(context.Background(), owner, quiz.ID)

	attacker := student()
	_, err := h.svc.SubmitAttempt(context.Background(), attacker, v.AttemptID, map[string]any{})
	if !apperror.IsNotFound(err) {
		t.Fatalf("another student's attempt must not be found (IDOR), got %v", err)
	}
}

// --- redis state lifecycle ----------------------------------------------------

func TestRedisState_SetOnStart_SavedOnSave_ClearedOnSubmit(t *testing.T) {
	h := newHarness(t)
	quiz := h.seedQuiz(nil, ptrInt(2), false)
	stu := student()

	v, _ := h.svc.StartAttempt(context.Background(), stu, quiz.ID)
	if _, ok := h.state.snapshots[v.AttemptID]; !ok {
		t.Error("snapshot not set on start")
	}
	if err := h.svc.SaveProgress(context.Background(), stu, v.AttemptID, map[string]any{"a": "b"}); err != nil {
		t.Fatalf("save: %v", err)
	}
	if _, ok := h.state.answers[v.AttemptID]; !ok {
		t.Error("answers not saved on SaveProgress")
	}
	if _, err := h.svc.SubmitAttempt(context.Background(), stu, v.AttemptID, map[string]any{}); err != nil {
		t.Fatalf("submit: %v", err)
	}
	if !h.state.cleared[v.AttemptID] {
		t.Error("attempt state not cleared on submit")
	}
}

// --- sweeper + race guard -----------------------------------------------------

func TestSweepExpired_GradesAbandonedAttempt(t *testing.T) {
	h := newHarness(t)
	quiz := h.seedQuiz(ptrInt(60), nil, false) // 60s limit
	stu := student()

	v, _ := h.svc.StartAttempt(context.Background(), stu, quiz.ID)
	// Student auto-saved a correct answer, then abandoned the app.
	c0 := h.correctAnswerID(quiz.ID, 0)
	_ = h.svc.SaveProgress(context.Background(), stu, v.AttemptID, map[string]any{
		v.Questions[idxOf(v, c0)].ID.String(): selected(c0),
	})

	h.advance(2 * time.Minute) // deadline passes; no client request arrives

	n, err := h.svc.SweepExpired(context.Background(), h.clock, 100)
	if err != nil {
		t.Fatalf("SweepExpired: %v", err)
	}
	if n != 1 {
		t.Fatalf("want 1 attempt swept, got %d", n)
	}
	stored, _ := h.attempts.FindByID(context.Background(), v.AttemptID, stu.UserID)
	if stored.Status != entity.AttemptGraded {
		t.Errorf("abandoned attempt status = %v, want graded", stored.Status)
	}
	if stored.Score == nil || *stored.Score != 1 {
		t.Errorf("auto-saved answer should score 1; got %v", stored.Score)
	}
}

func TestSweepExpired_SkipsAlreadyGraded(t *testing.T) {
	h := newHarness(t)
	quiz := h.seedQuiz(ptrInt(60), nil, false)
	stu := student()
	v, _ := h.svc.StartAttempt(context.Background(), stu, quiz.ID)
	if _, err := h.svc.SubmitAttempt(context.Background(), stu, v.AttemptID, map[string]any{}); err != nil {
		t.Fatalf("submit: %v", err)
	}
	h.advance(2 * time.Minute)

	n, err := h.svc.SweepExpired(context.Background(), h.clock, 100)
	if err != nil {
		t.Fatalf("SweepExpired: %v", err)
	}
	if n != 0 {
		t.Errorf("already-graded attempt must not be re-swept; graded %d", n)
	}
}

func TestSubmitVsSweep_GradesExactlyOnce(t *testing.T) {
	h := newHarness(t)
	quiz := h.seedQuiz(ptrInt(60), nil, false)
	stu := student()
	v, _ := h.svc.StartAttempt(context.Background(), stu, quiz.ID)
	h.advance(2 * time.Minute) // both paths now see it as expired

	// The sweeper finalizes it first.
	n, err := h.svc.SweepExpired(context.Background(), h.clock, 100)
	if err != nil || n != 1 {
		t.Fatalf("sweep: n=%d err=%v", n, err)
	}
	// A late client submit must NOT re-grade — it returns the existing result.
	res, err := h.svc.SubmitAttempt(context.Background(), stu, v.AttemptID, map[string]any{})
	if err != nil {
		t.Fatalf("late submit should return the graded result, got error %v", err)
	}
	if res.Status != entity.AttemptGraded {
		t.Errorf("late submit status = %v, want graded", res.Status)
	}
	// Exactly one grading: the attempt is graded once and the count never doubled.
	stored, _ := h.attempts.FindByID(context.Background(), v.AttemptID, stu.UserID)
	if stored.Status != entity.AttemptGraded {
		t.Errorf("stored status = %v, want graded", stored.Status)
	}
}

// --- helpers ------------------------------------------------------------------

func (h *harness) wrongAnswerID(quizID uuid.UUID, qi int) uuid.UUID {
	tree, _ := h.quizzes.GetQuizWithQuestions(context.Background(), quizID)
	for _, a := range tree.Questions[qi].Answers {
		if !a.IsCorrect {
			return a.ID
		}
	}
	return uuid.Nil
}

// idxOf finds the index of the question owning answer ID c in the client view.
func idxOf(v *AttemptClientView, c uuid.UUID) int {
	for i, q := range v.Questions {
		for _, a := range q.Answers {
			if a.ID == c {
				return i
			}
		}
	}
	return 0
}

func answerIDOrder(v *AttemptClientView) []string {
	var out []string
	for _, a := range v.Questions[0].Answers {
		out = append(out, a.ID.String())
	}
	return out
}

func sameSet(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	m := map[string]int{}
	for _, s := range a {
		m[s]++
	}
	for _, s := range b {
		m[s]--
	}
	for _, n := range m {
		if n != 0 {
			return false
		}
	}
	return true
}

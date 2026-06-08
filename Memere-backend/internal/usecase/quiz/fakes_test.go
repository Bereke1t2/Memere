package quiz

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/pagination"
)

// ---- fake course repo (only FindByID is exercised by the quiz engine) --------

type fakeCourseRepo struct {
	byID map[uuid.UUID]*entity.Course
}

func newFakeCourseRepo() *fakeCourseRepo {
	return &fakeCourseRepo{byID: map[uuid.UUID]*entity.Course{}}
}

func (f *fakeCourseRepo) add(c *entity.Course) { f.byID[c.ID] = c }

func (f *fakeCourseRepo) FindByID(_ context.Context, id uuid.UUID) (*entity.Course, error) {
	if c, ok := f.byID[id]; ok {
		cp := *c
		return &cp, nil
	}
	return nil, apperror.NotFound("course not found", nil)
}

func (f *fakeCourseRepo) Create(context.Context, *entity.Course) error { return nil }
func (f *fakeCourseRepo) FindBySlug(context.Context, string) (*entity.Course, error) {
	return nil, apperror.NotFound("course not found", nil)
}
func (f *fakeCourseRepo) List(context.Context, repository.CourseFilter, *pagination.Cursor, int) ([]*entity.Course, *pagination.Cursor, error) {
	return nil, nil, nil
}
func (f *fakeCourseRepo) Update(context.Context, *entity.Course) error       { return nil }
func (f *fakeCourseRepo) SoftDelete(context.Context, uuid.UUID) error        { return nil }
func (f *fakeCourseRepo) RecomputeCounters(context.Context, uuid.UUID) error { return nil }
func (f *fakeCourseRepo) GetCourseWithSectionsAndLessons(context.Context, uuid.UUID) (*repository.CourseWithContent, error) {
	return nil, apperror.NotFound("not found", nil)
}

// ---- fake quiz repo ----------------------------------------------------------

type fakeQuizRepo struct {
	mu        sync.Mutex
	quizzes   map[uuid.UUID]*entity.Quiz
	questions map[uuid.UUID][]repository.QuestionWithAnswers // by quizID
}

func newFakeQuizRepo() *fakeQuizRepo {
	return &fakeQuizRepo{
		quizzes:   map[uuid.UUID]*entity.Quiz{},
		questions: map[uuid.UUID][]repository.QuestionWithAnswers{},
	}
}

func (f *fakeQuizRepo) Create(_ context.Context, q *entity.Quiz) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if q.ID == uuid.Nil {
		q.ID = uuid.New()
	}
	q.CreatedAt = time.Now()
	q.UpdatedAt = q.CreatedAt
	cp := *q
	f.quizzes[q.ID] = &cp
	return nil
}

func (f *fakeQuizRepo) FindByID(_ context.Context, id uuid.UUID) (*entity.Quiz, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if q, ok := f.quizzes[id]; ok && q.DeletedAt == nil {
		cp := *q
		return &cp, nil
	}
	return nil, apperror.NotFound("quiz not found", nil)
}

func (f *fakeQuizRepo) ListByCourse(_ context.Context, courseID uuid.UUID) ([]*entity.Quiz, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	var out []*entity.Quiz
	for _, q := range f.quizzes {
		if q.CourseID == courseID && q.DeletedAt == nil {
			cp := *q
			out = append(out, &cp)
		}
	}
	return out, nil
}

func (f *fakeQuizRepo) Update(_ context.Context, q *entity.Quiz) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if _, ok := f.quizzes[q.ID]; !ok {
		return apperror.NotFound("quiz not found", nil)
	}
	q.UpdatedAt = time.Now()
	cp := *q
	f.quizzes[q.ID] = &cp
	return nil
}

func (f *fakeQuizRepo) Delete(_ context.Context, id uuid.UUID) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	q, ok := f.quizzes[id]
	if !ok {
		return apperror.NotFound("quiz not found", nil)
	}
	now := time.Now()
	q.DeletedAt = &now
	return nil
}

func (f *fakeQuizRepo) GetQuizWithQuestions(_ context.Context, quizID uuid.UUID) (*repository.QuizWithQuestions, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	q, ok := f.quizzes[quizID]
	if !ok {
		return nil, apperror.NotFound("quiz not found", nil)
	}
	cp := *q
	return &repository.QuizWithQuestions{Quiz: &cp, Questions: f.questions[quizID]}, nil
}

func (f *fakeQuizRepo) GetQuestionsForClient(_ context.Context, quizID uuid.UUID) ([]repository.ClientQuestion, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	var out []repository.ClientQuestion
	for _, qa := range f.questions[quizID] {
		answers := make([]repository.ClientAnswer, len(qa.Answers))
		for i, a := range qa.Answers {
			// NOTE: this fake deliberately mirrors the real client query — it
			// copies only non-secret fields, never IsCorrect.
			answers[i] = repository.ClientAnswer{ID: a.ID, QuestionID: a.QuestionID, Text: a.Text, OrderIndex: a.OrderIndex}
		}
		out = append(out, repository.ClientQuestion{
			ID:         qa.Question.ID,
			QuizID:     qa.Question.QuizID,
			Text:       qa.Question.Text,
			Type:       qa.Question.Type,
			Points:     qa.Question.Points,
			OrderIndex: qa.Question.OrderIndex,
			Subject:    qa.Question.Subject,
			Topic:      qa.Question.Topic,
			Answers:    answers,
		})
	}
	return out, nil
}

// addQuestion is a test helper that registers a question + answers under a quiz.
func (f *fakeQuizRepo) addQuestion(quizID uuid.UUID, q *entity.Question, answers []*entity.Answer) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if q.ID == uuid.Nil {
		q.ID = uuid.New()
	}
	q.QuizID = quizID
	for _, a := range answers {
		if a.ID == uuid.Nil {
			a.ID = uuid.New()
		}
		a.QuestionID = q.ID
	}
	f.questions[quizID] = append(f.questions[quizID], repository.QuestionWithAnswers{Question: q, Answers: answers})
}

// ---- fake question repo ------------------------------------------------------

type fakeQuestionRepo struct {
	quizzes *fakeQuizRepo
}

func (f *fakeQuestionRepo) Create(_ context.Context, q *entity.Question, answers []*entity.Answer) error {
	f.quizzes.addQuestion(q.QuizID, q, answers)
	return nil
}
func (f *fakeQuestionRepo) FindByID(context.Context, uuid.UUID) (*entity.Question, error) {
	return nil, apperror.NotFound("question not found", nil)
}
func (f *fakeQuestionRepo) ListByQuiz(context.Context, uuid.UUID) ([]*entity.Question, error) {
	return nil, nil
}
func (f *fakeQuestionRepo) ListAnswers(context.Context, uuid.UUID) ([]*entity.Answer, error) {
	return nil, nil
}
func (f *fakeQuestionRepo) Update(context.Context, *entity.Question) error { return nil }
func (f *fakeQuestionRepo) Delete(context.Context, uuid.UUID) error        { return nil }

// ---- fake quiz attempt repo --------------------------------------------------

type fakeAttemptRepo struct {
	mu    sync.Mutex
	byID  map[uuid.UUID]*entity.QuizAttempt
	order []uuid.UUID
}

func newFakeAttemptRepo() *fakeAttemptRepo {
	return &fakeAttemptRepo{byID: map[uuid.UUID]*entity.QuizAttempt{}}
}

func (f *fakeAttemptRepo) Create(_ context.Context, a *entity.QuizAttempt) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	a.CreatedAt = time.Now()
	a.UpdatedAt = a.CreatedAt
	cp := *a
	f.byID[a.ID] = &cp
	f.order = append(f.order, a.ID)
	return nil
}

func (f *fakeAttemptRepo) FindByID(_ context.Context, id, studentID uuid.UUID) (*entity.QuizAttempt, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	a, ok := f.byID[id]
	if !ok || a.StudentID != studentID {
		return nil, apperror.NotFound("quiz attempt not found", nil)
	}
	cp := *a
	return &cp, nil
}

func (f *fakeAttemptRepo) GetActive(_ context.Context, studentID, quizID uuid.UUID) (*entity.QuizAttempt, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	for i := len(f.order) - 1; i >= 0; i-- {
		a := f.byID[f.order[i]]
		if a.StudentID == studentID && a.QuizID == quizID && a.Status == entity.AttemptInProgress {
			cp := *a
			return &cp, nil
		}
	}
	return nil, apperror.NotFound("no active attempt", nil)
}

func (f *fakeAttemptRepo) ListByStudentAndQuiz(_ context.Context, studentID, quizID uuid.UUID) ([]*entity.QuizAttempt, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	var out []*entity.QuizAttempt
	for _, id := range f.order {
		a := f.byID[id]
		if a.StudentID == studentID && a.QuizID == quizID {
			cp := *a
			out = append(out, &cp)
		}
	}
	return out, nil
}

func (f *fakeAttemptRepo) CountByStudentAndQuiz(_ context.Context, studentID, quizID uuid.UUID) (int, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	n := 0
	for _, a := range f.byID {
		if a.StudentID == studentID && a.QuizID == quizID {
			n++
		}
	}
	return n, nil
}

func (f *fakeAttemptRepo) ListExpired(_ context.Context, now time.Time, limit int) ([]*entity.QuizAttempt, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	var out []*entity.QuizAttempt
	for _, id := range f.order {
		a := f.byID[id]
		if a.Status == entity.AttemptInProgress && a.ExpiresAt != nil && a.ExpiresAt.Before(now) {
			cp := *a
			out = append(out, &cp)
			if len(out) >= limit {
				break
			}
		}
	}
	return out, nil
}

func (f *fakeAttemptRepo) Update(_ context.Context, a *entity.QuizAttempt) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if _, ok := f.byID[a.ID]; !ok {
		return apperror.NotFound("quiz attempt not found", nil)
	}
	a.UpdatedAt = time.Now()
	cp := *a
	f.byID[a.ID] = &cp
	return nil
}

func (f *fakeAttemptRepo) Grade(_ context.Context, a *entity.QuizAttempt) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	cur, ok := f.byID[a.ID]
	if !ok {
		return apperror.NotFound("quiz attempt not found", nil)
	}
	cur.Score = a.Score
	cur.Percentage = a.Percentage
	cur.Passed = a.Passed
	cur.Status = entity.AttemptGraded
	cur.UpdatedAt = time.Now()
	return nil
}

// ---- fake attempt-state store (in-memory Redis) ------------------------------

type fakeState struct {
	mu        sync.Mutex
	snapshots map[uuid.UUID]map[string]any
	answers   map[uuid.UUID]map[string]any
	cleared   map[uuid.UUID]bool
}

func newFakeState() *fakeState {
	return &fakeState{
		snapshots: map[uuid.UUID]map[string]any{},
		answers:   map[uuid.UUID]map[string]any{},
		cleared:   map[uuid.UUID]bool{},
	}
}

func (f *fakeState) SetSnapshot(_ context.Context, id uuid.UUID, snap map[string]any, _ time.Duration) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.snapshots[id] = snap
	return nil
}
func (f *fakeState) GetSnapshot(_ context.Context, id uuid.UUID) (map[string]any, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.snapshots[id], nil
}
func (f *fakeState) SaveAnswers(_ context.Context, id uuid.UUID, ans map[string]any, _ time.Duration) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.answers[id] = ans
	return nil
}
func (f *fakeState) GetAnswers(_ context.Context, id uuid.UUID) (map[string]any, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.answers[id], nil
}
func (f *fakeState) DeleteAttemptState(_ context.Context, id uuid.UUID) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.snapshots, id)
	delete(f.answers, id)
	f.cleared[id] = true
	return nil
}

// ---- fake tx manager ---------------------------------------------------------

type fakeTxManager struct{}

func (fakeTxManager) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

// ---- compile-time interface checks -------------------------------------------

var (
	_ repository.CourseRepository      = (*fakeCourseRepo)(nil)
	_ repository.QuizRepository        = (*fakeQuizRepo)(nil)
	_ repository.QuestionRepository    = (*fakeQuestionRepo)(nil)
	_ repository.QuizAttemptRepository = (*fakeAttemptRepo)(nil)
	_ repository.AttemptStateStore     = (*fakeState)(nil)
	_ repository.TxManager             = fakeTxManager{}
)

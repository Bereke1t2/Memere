package exam

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

// ---- fake course repo (only FindByID is exercised) ---------------------------

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

// ---- fake exam repo ----------------------------------------------------------

type fakeExamRepo struct {
	mu        sync.Mutex
	exams     map[uuid.UUID]*entity.Exam
	links     map[uuid.UUID][]*entity.ExamQuestion         // by examID
	questions map[uuid.UUID]repository.QuestionWithAnswers // by questionID (bank)
	examQs    map[uuid.UUID][]uuid.UUID                    // examID -> ordered questionIDs
}

func newFakeExamRepo() *fakeExamRepo {
	return &fakeExamRepo{
		exams:     map[uuid.UUID]*entity.Exam{},
		links:     map[uuid.UUID][]*entity.ExamQuestion{},
		questions: map[uuid.UUID]repository.QuestionWithAnswers{},
		examQs:    map[uuid.UUID][]uuid.UUID{},
	}
}

func (f *fakeExamRepo) Create(_ context.Context, e *entity.Exam) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	e.CreatedAt = time.Now()
	e.UpdatedAt = e.CreatedAt
	cp := *e
	f.exams[e.ID] = &cp
	return nil
}

func (f *fakeExamRepo) FindByID(_ context.Context, id uuid.UUID) (*entity.Exam, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if e, ok := f.exams[id]; ok && e.DeletedAt == nil {
		cp := *e
		return &cp, nil
	}
	return nil, apperror.NotFound("exam not found", nil)
}

func (f *fakeExamRepo) List(_ context.Context, filter repository.ExamFilter, _ *pagination.Cursor, limit int) ([]*entity.Exam, *pagination.Cursor, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	var out []*entity.Exam
	for _, e := range f.exams {
		if e.DeletedAt != nil {
			continue
		}
		if filter.IsPublished != nil && e.IsPublished != *filter.IsPublished {
			continue
		}
		cp := *e
		out = append(out, &cp)
	}
	return out, nil, nil
}

func (f *fakeExamRepo) ListByCourse(_ context.Context, courseID uuid.UUID) ([]*entity.Exam, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	var out []*entity.Exam
	for _, e := range f.exams {
		if e.DeletedAt != nil {
			continue
		}
		if e.CourseID != nil && *e.CourseID == courseID {
			cp := *e
			out = append(out, &cp)
		}
	}
	return out, nil
}

func (f *fakeExamRepo) Update(_ context.Context, e *entity.Exam) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if _, ok := f.exams[e.ID]; !ok {
		return apperror.NotFound("exam not found", nil)
	}
	e.UpdatedAt = time.Now()
	cp := *e
	f.exams[e.ID] = &cp
	return nil
}

func (f *fakeExamRepo) Delete(_ context.Context, id uuid.UUID) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	e, ok := f.exams[id]
	if !ok {
		return apperror.NotFound("exam not found", nil)
	}
	now := time.Now()
	e.DeletedAt = &now
	return nil
}

func (f *fakeExamRepo) AddQuestion(_ context.Context, eq *entity.ExamQuestion) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if eq.ID == uuid.Nil {
		eq.ID = uuid.New()
	}
	for _, l := range f.links[eq.ExamID] {
		if l.QuestionID == eq.QuestionID {
			return apperror.Conflict("QUESTION_ALREADY_IN_EXAM", nil)
		}
	}
	cp := *eq
	f.links[eq.ExamID] = append(f.links[eq.ExamID], &cp)
	f.examQs[eq.ExamID] = append(f.examQs[eq.ExamID], eq.QuestionID)
	return nil
}

func (f *fakeExamRepo) ListQuestions(_ context.Context, examID uuid.UUID) ([]*entity.ExamQuestion, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	var out []*entity.ExamQuestion
	for _, l := range f.links[examID] {
		cp := *l
		out = append(out, &cp)
	}
	return out, nil
}

func (f *fakeExamRepo) GetExamWithQuestions(_ context.Context, examID uuid.UUID) (*repository.ExamWithQuestions, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	e, ok := f.exams[examID]
	if !ok {
		return nil, apperror.NotFound("exam not found", nil)
	}
	var qs []repository.QuestionWithAnswers
	for _, qid := range f.examQs[examID] {
		qs = append(qs, f.questions[qid])
	}
	cp := *e
	return &repository.ExamWithQuestions{Exam: &cp, Questions: qs}, nil
}

func (f *fakeExamRepo) GetQuestionsForClient(_ context.Context, examID uuid.UUID) ([]repository.ClientExamQuestion, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	var out []repository.ClientExamQuestion
	marks := map[uuid.UUID]int{}
	for _, l := range f.links[examID] {
		marks[l.QuestionID] = l.Marks
	}
	for _, qid := range f.examQs[examID] {
		qa := f.questions[qid]
		answers := make([]repository.ClientAnswer, len(qa.Answers))
		for i, a := range qa.Answers {
			// Mirrors the real client query: never copies IsCorrect.
			answers[i] = repository.ClientAnswer{ID: a.ID, QuestionID: a.QuestionID, Text: a.Text, OrderIndex: a.OrderIndex}
		}
		out = append(out, repository.ClientExamQuestion{
			QuestionID: qa.Question.ID,
			OrderIndex: qa.Question.OrderIndex,
			Marks:      marks[qid],
			Text:       qa.Question.Text,
			Type:       qa.Question.Type,
			Subject:    qa.Question.Subject,
			Topic:      qa.Question.Topic,
			Answers:    answers,
		})
	}
	return out, nil
}

// addBankQuestion registers a bank question with answers (test helper).
func (f *fakeExamRepo) addBankQuestion(q *entity.Question, answers []*entity.Answer) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if q.ID == uuid.Nil {
		q.ID = uuid.New()
	}
	for _, a := range answers {
		if a.ID == uuid.Nil {
			a.ID = uuid.New()
		}
		a.QuestionID = q.ID
	}
	f.questions[q.ID] = repository.QuestionWithAnswers{Question: q, Answers: answers}
}

// ---- fake exam attempt repo --------------------------------------------------

type fakeExamAttemptRepo struct {
	mu    sync.Mutex
	byID  map[uuid.UUID]*entity.ExamAttempt
	order []uuid.UUID
}

func newFakeExamAttemptRepo() *fakeExamAttemptRepo {
	return &fakeExamAttemptRepo{byID: map[uuid.UUID]*entity.ExamAttempt{}}
}

func (f *fakeExamAttemptRepo) Create(_ context.Context, a *entity.ExamAttempt) error {
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

func (f *fakeExamAttemptRepo) FindByID(_ context.Context, id, studentID uuid.UUID) (*entity.ExamAttempt, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	a, ok := f.byID[id]
	if !ok || a.StudentID != studentID {
		return nil, apperror.NotFound("exam attempt not found", nil)
	}
	cp := *a
	return &cp, nil
}

func (f *fakeExamAttemptRepo) GetActive(_ context.Context, studentID, examID uuid.UUID) (*entity.ExamAttempt, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	for i := len(f.order) - 1; i >= 0; i-- {
		a := f.byID[f.order[i]]
		if a.StudentID == studentID && a.ExamID == examID && a.Status == entity.AttemptInProgress {
			cp := *a
			return &cp, nil
		}
	}
	return nil, apperror.NotFound("no active attempt", nil)
}

func (f *fakeExamAttemptRepo) ListByStudent(_ context.Context, studentID uuid.UUID) ([]*entity.ExamAttempt, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	var out []*entity.ExamAttempt
	for _, id := range f.order {
		a := f.byID[id]
		if a.StudentID == studentID {
			cp := *a
			out = append(out, &cp)
		}
	}
	return out, nil
}

// expiryLookup lets the fake compute deadlines from exam durations; set by the
// harness so ListExpired can mirror the real JOIN on exams.duration_minutes.
var expiryLookup func(examID uuid.UUID) (time.Duration, bool)

func (f *fakeExamAttemptRepo) ListExpired(_ context.Context, now time.Time, limit int) ([]*entity.ExamAttempt, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	var out []*entity.ExamAttempt
	for _, id := range f.order {
		a := f.byID[id]
		if a.Status != entity.AttemptInProgress {
			continue
		}
		if expiryLookup == nil {
			continue
		}
		dur, ok := expiryLookup(a.ExamID)
		if !ok {
			continue
		}
		if now.After(a.StartedAt.Add(dur)) {
			cp := *a
			out = append(out, &cp)
			if len(out) >= limit {
				break
			}
		}
	}
	return out, nil
}

func (f *fakeExamAttemptRepo) Update(_ context.Context, a *entity.ExamAttempt) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if _, ok := f.byID[a.ID]; !ok {
		return apperror.NotFound("exam attempt not found", nil)
	}
	a.UpdatedAt = time.Now()
	cp := *a
	f.byID[a.ID] = &cp
	return nil
}

func (f *fakeExamAttemptRepo) ClaimForGrading(_ context.Context, a *entity.ExamAttempt) (bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	cur, ok := f.byID[a.ID]
	if !ok {
		return false, apperror.NotFound("exam attempt not found", nil)
	}
	if cur.Status != entity.AttemptInProgress {
		return false, nil
	}
	cur.Status = a.Status
	cur.AnswersSnapshot = a.AnswersSnapshot
	cur.SubmittedAt = a.SubmittedAt
	cur.UpdatedAt = time.Now()
	cp := *cur
	*a = cp
	return true, nil
}

func (f *fakeExamAttemptRepo) ListGradedBySubject(context.Context, uuid.UUID, string) ([]*entity.ExamAttempt, error) {
	return nil, nil
}

func (f *fakeExamAttemptRepo) Stats(context.Context, uuid.UUID) (repository.ExamAttemptStats, error) {
	return repository.ExamAttemptStats{}, nil
}

func (f *fakeExamAttemptRepo) Grade(_ context.Context, a *entity.ExamAttempt) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	cur, ok := f.byID[a.ID]
	if !ok {
		return apperror.NotFound("exam attempt not found", nil)
	}
	cur.Score = a.Score
	cur.Percentage = a.Percentage
	cur.Status = entity.AttemptGraded
	cur.UpdatedAt = time.Now()
	return nil
}

// ---- fake attempt-state store ------------------------------------------------

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

// ---- fake enrollment / subscription repos (drive access.Service) -------------
//
// The exam engine doesn't touch these directly; they back the real
// access.Service wired into the harness so the FullAccess gate can be exercised.

type fakeEnrollRepo struct {
	active map[[2]uuid.UUID]*entity.Enrollment // (student,course) -> active enrollment
}

func newFakeEnrollRepo() *fakeEnrollRepo {
	return &fakeEnrollRepo{active: map[[2]uuid.UUID]*entity.Enrollment{}}
}

func (f *fakeEnrollRepo) enroll(studentID, courseID uuid.UUID) {
	f.active[[2]uuid.UUID{studentID, courseID}] = &entity.Enrollment{StudentID: studentID, CourseID: courseID}
}

func (f *fakeEnrollRepo) Create(context.Context, *entity.Enrollment) error { return nil }
func (f *fakeEnrollRepo) Exists(_ context.Context, s, c uuid.UUID) (bool, error) {
	_, ok := f.active[[2]uuid.UUID{s, c}]
	return ok, nil
}
func (f *fakeEnrollRepo) GetActiveForStudent(_ context.Context, s, c uuid.UUID) (*entity.Enrollment, error) {
	if e, ok := f.active[[2]uuid.UUID{s, c}]; ok {
		return e, nil
	}
	return nil, apperror.NotFound("enrollment not found", nil)
}
func (f *fakeEnrollRepo) ListByStudent(context.Context, uuid.UUID, int) ([]*entity.Enrollment, error) {
	return nil, nil
}

type fakeSubRepo struct {
	active map[uuid.UUID]*entity.Subscription
}

func newFakeSubRepo() *fakeSubRepo {
	return &fakeSubRepo{active: map[uuid.UUID]*entity.Subscription{}}
}

func (f *fakeSubRepo) Create(context.Context, *entity.Subscription) error { return nil }
func (f *fakeSubRepo) GetByID(context.Context, uuid.UUID) (*entity.Subscription, error) {
	return nil, apperror.NotFound("subscription not found", nil)
}
func (f *fakeSubRepo) GetActiveForStudent(_ context.Context, s uuid.UUID) (*entity.Subscription, error) {
	if sub, ok := f.active[s]; ok {
		return sub, nil
	}
	return nil, apperror.NotFound("subscription not found", nil)
}
func (f *fakeSubRepo) UpdateStatus(context.Context, uuid.UUID, entity.SubscriptionStatus) error {
	return nil
}
func (f *fakeSubRepo) ListExpiring(context.Context, int) ([]*entity.Subscription, error) {
	return nil, nil
}

func (f *fakeSubRepo) ExtendPeriod(context.Context, uuid.UUID, time.Time) error { return nil }
func (f *fakeSubRepo) CancelAtPeriodEnd(context.Context, uuid.UUID) (bool, error) {
	return true, nil
}
func (f *fakeSubRepo) ExpireLapsed(context.Context, uuid.UUID) (bool, error) {
	return false, nil
}

// ---- compile-time interface checks -------------------------------------------

var (
	_ repository.CourseRepository       = (*fakeCourseRepo)(nil)
	_ repository.ExamRepository         = (*fakeExamRepo)(nil)
	_ repository.ExamAttemptRepository  = (*fakeExamAttemptRepo)(nil)
	_ repository.AttemptStateStore      = (*fakeState)(nil)
	_ repository.TxManager              = fakeTxManager{}
	_ repository.EnrollmentRepository   = (*fakeEnrollRepo)(nil)
	_ repository.SubscriptionRepository = (*fakeSubRepo)(nil)
)

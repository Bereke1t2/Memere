package certificate

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/service"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/pagination"
)

// ---- fake certificate repo --------------------------------------------------

type fakeCertRepo struct {
	mu         sync.Mutex
	byID       map[uuid.UUID]*entity.Certificate
	byPair     map[[2]uuid.UUID]*entity.Certificate // (student, course)
	bySerial   map[string]*entity.Certificate
}

func newFakeCertRepo() *fakeCertRepo {
	return &fakeCertRepo{
		byID:     map[uuid.UUID]*entity.Certificate{},
		byPair:   map[[2]uuid.UUID]*entity.Certificate{},
		bySerial: map[string]*entity.Certificate{},
	}
}

func (f *fakeCertRepo) CreateIfNotExists(_ context.Context, studentID, courseID uuid.UUID, serial string) (*entity.Certificate, bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	key := [2]uuid.UUID{studentID, courseID}
	if existing, ok := f.byPair[key]; ok {
		cp := *existing
		return &cp, false, nil
	}
	cert := &entity.Certificate{
		ID:        uuid.New(),
		StudentID: studentID,
		CourseID:  courseID,
		Serial:    serial,
		IssuedAt:  time.Now().UTC(),
	}
	f.byID[cert.ID] = cert
	f.byPair[key] = cert
	f.bySerial[serial] = cert
	cp := *cert
	return &cp, true, nil
}

func (f *fakeCertRepo) GetByID(_ context.Context, id uuid.UUID) (*entity.Certificate, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if c, ok := f.byID[id]; ok {
		cp := *c
		return &cp, nil
	}
	return nil, apperror.NotFound("certificate not found", nil)
}

func (f *fakeCertRepo) GetByStudentCourse(_ context.Context, studentID, courseID uuid.UUID) (*entity.Certificate, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if c, ok := f.byPair[[2]uuid.UUID{studentID, courseID}]; ok {
		cp := *c
		return &cp, nil
	}
	return nil, apperror.NotFound("certificate not found", nil)
}

func (f *fakeCertRepo) GetBySerial(_ context.Context, serial string) (*entity.Certificate, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if c, ok := f.bySerial[serial]; ok {
		cp := *c
		return &cp, nil
	}
	return nil, apperror.NotFound("certificate not found", nil)
}

func (f *fakeCertRepo) SetFileKey(_ context.Context, id uuid.UUID, fileKey string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if c, ok := f.byID[id]; ok {
		cp := fileKey
		c.FileKey = &cp
	}
	return nil
}

func (f *fakeCertRepo) ListByStudent(_ context.Context, studentID uuid.UUID) ([]*entity.Certificate, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	var out []*entity.Certificate
	for _, c := range f.byID {
		if c.StudentID == studentID {
			cp := *c
			out = append(out, &cp)
		}
	}
	return out, nil
}

// ---- fake progress repo (certificate tests only need GetCourseProgress) -----

type fakeProgressRepo struct {
	mu       sync.Mutex
	progress map[[2]uuid.UUID]*entity.CourseProgress
}

func newFakeProgressRepo() *fakeProgressRepo {
	return &fakeProgressRepo{progress: map[[2]uuid.UUID]*entity.CourseProgress{}}
}

func (f *fakeProgressRepo) setCourseProgress(studentID, courseID uuid.UUID, pct float64) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.progress[[2]uuid.UUID{studentID, courseID}] = &entity.CourseProgress{
		StudentID:       studentID,
		CourseID:        courseID,
		PercentComplete: decimal.NewFromFloat(pct),
	}
}

func (f *fakeProgressRepo) GetCourseProgress(_ context.Context, studentID, courseID uuid.UUID) (*entity.CourseProgress, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if cp, ok := f.progress[[2]uuid.UUID{studentID, courseID}]; ok {
		return cp, nil
	}
	return nil, apperror.NotFound("progress not found", nil)
}

func (f *fakeProgressRepo) UpsertLessonProgress(context.Context, *entity.LessonProgress) (*entity.LessonProgress, error) {
	return nil, nil
}
func (f *fakeProgressRepo) GetLessonProgress(context.Context, uuid.UUID, uuid.UUID) (*entity.LessonProgress, error) {
	return nil, apperror.NotFound("not found", nil)
}
func (f *fakeProgressRepo) ListByCourse(context.Context, uuid.UUID, uuid.UUID) ([]*entity.LessonProgress, error) {
	return nil, nil
}
func (f *fakeProgressRepo) RecomputeCourseProgress(context.Context, uuid.UUID, uuid.UUID) (*entity.CourseProgress, error) {
	return nil, nil
}
func (f *fakeProgressRepo) GetStreak(context.Context, uuid.UUID) (*entity.StudyStreak, error) {
	return &entity.StudyStreak{}, nil
}
func (f *fakeProgressRepo) UpsertStreak(context.Context, *entity.StudyStreak) error { return nil }
func (f *fakeProgressRepo) ListStreakAtRisk(context.Context, time.Time, time.Time, int) ([]*entity.StudyStreak, error) {
	return nil, nil
}
func (f *fakeProgressRepo) SetLastWarnedDate(context.Context, uuid.UUID, time.Time) error {
	return nil
}

// ---- fake course repo -------------------------------------------------------

type fakeCourseRepo struct {
	mu   sync.Mutex
	byID map[uuid.UUID]*entity.Course
}

func newFakeCourseRepo() *fakeCourseRepo {
	return &fakeCourseRepo{byID: map[uuid.UUID]*entity.Course{}}
}

func (f *fakeCourseRepo) add(c *entity.Course) {
	f.mu.Lock()
	defer f.mu.Unlock()
	cp := *c
	f.byID[c.ID] = &cp
}

func (f *fakeCourseRepo) FindByID(_ context.Context, id uuid.UUID) (*entity.Course, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if c, ok := f.byID[id]; ok {
		cp := *c
		return &cp, nil
	}
	return nil, apperror.NotFound("course not found", nil)
}
func (f *fakeCourseRepo) FindBySlug(context.Context, string) (*entity.Course, error) {
	return nil, apperror.NotFound("not found", nil)
}
func (f *fakeCourseRepo) List(context.Context, repository.CourseFilter, *pagination.Cursor, int) ([]*entity.Course, *pagination.Cursor, error) {
	return nil, nil, nil
}
func (f *fakeCourseRepo) Create(context.Context, *entity.Course) error  { return nil }
func (f *fakeCourseRepo) Update(context.Context, *entity.Course) error  { return nil }
func (f *fakeCourseRepo) SoftDelete(context.Context, uuid.UUID) error   { return nil }
func (f *fakeCourseRepo) RecomputeCounters(context.Context, uuid.UUID) error { return nil }
func (f *fakeCourseRepo) GetCourseWithSectionsAndLessons(context.Context, uuid.UUID) (*repository.CourseWithContent, error) {
	return nil, apperror.NotFound("not found", nil)
}

// ---- fake user repo ---------------------------------------------------------

type fakeUserRepo struct {
	mu   sync.Mutex
	byID map[uuid.UUID]*entity.User
}

func newFakeUserRepo() *fakeUserRepo {
	return &fakeUserRepo{byID: map[uuid.UUID]*entity.User{}}
}

func (f *fakeUserRepo) add(u *entity.User) {
	f.mu.Lock()
	defer f.mu.Unlock()
	cp := *u
	f.byID[u.ID] = &cp
}

func (f *fakeUserRepo) FindByID(_ context.Context, id uuid.UUID) (*entity.User, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if u, ok := f.byID[id]; ok {
		cp := *u
		return &cp, nil
	}
	return nil, apperror.NotFound("user not found", nil)
}
func (f *fakeUserRepo) Create(context.Context, *entity.User) error { return nil }
func (f *fakeUserRepo) FindByEmail(context.Context, string) (*entity.User, error) {
	return nil, apperror.NotFound("not found", nil)
}
func (f *fakeUserRepo) Update(context.Context, *entity.User) error       { return nil }
func (f *fakeUserRepo) SoftDelete(context.Context, uuid.UUID) error      { return nil }
func (f *fakeUserRepo) SetLastLogin(context.Context, uuid.UUID, time.Time) error { return nil }
func (f *fakeUserRepo) List(context.Context, repository.AdminUserFilter, *pagination.Cursor, int) ([]*entity.User, *pagination.Cursor, error) {
	return nil, nil, nil
}
func (f *fakeUserRepo) CountByRole(context.Context, entity.Role) (int, error) { return 0, nil }

// ---- fake object store ------------------------------------------------------

type fakeObjectStore struct {
	mu   sync.Mutex
	keys map[string][]byte
}

func newFakeObjectStore() *fakeObjectStore {
	return &fakeObjectStore{keys: map[string][]byte{}}
}

func (f *fakeObjectStore) Put(_ context.Context, key, _ string, body io.Reader) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(body)
	f.keys[key] = buf.Bytes()
	return nil
}

func (f *fakeObjectStore) PresignPut(context.Context, string, string, time.Duration) (string, error) {
	return "", nil
}
func (f *fakeObjectStore) PresignGet(_ context.Context, key string, _ time.Duration) (string, error) {
	return "https://cdn.example.com/" + key, nil
}
func (f *fakeObjectStore) Get(context.Context, string) (io.ReadCloser, error) {
	return io.NopCloser(strings.NewReader("")), nil
}
func (f *fakeObjectStore) Exists(context.Context, string) (bool, error)  { return false, nil }
func (f *fakeObjectStore) Delete(context.Context, string) error          { return nil }

func (f *fakeObjectStore) has(key string) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	_, ok := f.keys[key]
	return ok
}

// ---- fake URL signer --------------------------------------------------------

type fakeURLSigner struct{}

func (fakeURLSigner) SignGet(_ context.Context, key string, _ time.Duration) (string, error) {
	return "https://signed.cdn.example.com/" + key, nil
}

// ---- fake renderer ----------------------------------------------------------

type fakeRenderer struct {
	mu    sync.Mutex
	calls []CertData
}

func (r *fakeRenderer) Certificate(_ context.Context, data CertData) ([]byte, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.calls = append(r.calls, data)
	return []byte(fmt.Sprintf("PDF:%s:%s", data.StudentName, data.CourseTitle)), nil
}

func (r *fakeRenderer) count() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.calls)
}

// ---- fake notifier ----------------------------------------------------------

type fakeNotifier struct {
	mu    sync.Mutex
	calls int
}

func (n *fakeNotifier) CertificateReady(_ context.Context, _, _ uuid.UUID) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.calls++
}

func (n *fakeNotifier) count() int {
	n.mu.Lock()
	defer n.mu.Unlock()
	return n.calls
}

// ---- compile-time interface checks ------------------------------------------

var _ repository.CertificateRepository = (*fakeCertRepo)(nil)
var _ repository.ProgressRepository = (*fakeProgressRepo)(nil)
var _ repository.CourseRepository = (*fakeCourseRepo)(nil)
var _ repository.UserRepository = (*fakeUserRepo)(nil)
var _ service.ObjectStore = (*fakeObjectStore)(nil)
var _ service.URLSigner = fakeURLSigner{}
var _ Renderer = (*fakeRenderer)(nil)
var _ Notifier = (*fakeNotifier)(nil)

// ---- helpers ----------------------------------------------------------------

// newSvc builds a Service wired to all fakes. notify accepts *fakeNotifier or
// nil (nil becomes noopNotifier in NewService via the Notifier interface check).
func newSvc(
	certRepo *fakeCertRepo,
	progressRepo *fakeProgressRepo,
	courses *fakeCourseRepo,
	users *fakeUserRepo,
	store *fakeObjectStore,
	renderer *fakeRenderer,
	notify Notifier,
) *Service {
	return NewService(certRepo, progressRepo, courses, users, store, fakeURLSigner{}, renderer, notify, Config{})
}

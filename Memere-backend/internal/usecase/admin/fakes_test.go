package admin

import (
	"context"
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

// ---- fake user repo ---------------------------------------------------------

type fakeUserRepo struct {
	mu    sync.Mutex
	byID  map[uuid.UUID]*entity.User
	byEmail map[string]*entity.User
}

func newFakeUserRepo() *fakeUserRepo {
	return &fakeUserRepo{
		byID:    map[uuid.UUID]*entity.User{},
		byEmail: map[string]*entity.User{},
	}
}

func (f *fakeUserRepo) add(u *entity.User) {
	f.mu.Lock()
	defer f.mu.Unlock()
	cp := *u
	f.byID[u.ID] = &cp
	f.byEmail[u.Email] = &cp
}

func (f *fakeUserRepo) Create(_ context.Context, u *entity.User) error {
	f.add(u)
	return nil
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

func (f *fakeUserRepo) FindByEmail(_ context.Context, email string) (*entity.User, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if u, ok := f.byEmail[email]; ok {
		cp := *u
		return &cp, nil
	}
	return nil, apperror.NotFound("user not found", nil)
}

func (f *fakeUserRepo) Update(_ context.Context, u *entity.User) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	cp := *u
	f.byID[u.ID] = &cp
	return nil
}

func (f *fakeUserRepo) SoftDelete(_ context.Context, id uuid.UUID) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if u, ok := f.byID[id]; ok {
		now := time.Now()
		u.DeletedAt = &now
	}
	return nil
}

func (f *fakeUserRepo) SetLastLogin(_ context.Context, id uuid.UUID, t time.Time) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if u, ok := f.byID[id]; ok {
		u.LastLoginAt = &t
	}
	return nil
}

func (f *fakeUserRepo) List(_ context.Context, filter repository.AdminUserFilter, _ *pagination.Cursor, limit int) ([]*entity.User, *pagination.Cursor, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	var out []*entity.User
	for _, u := range f.byID {
		if filter.Role != nil && string(u.Role) != *filter.Role {
			continue
		}
		if filter.IsActive != nil && u.IsActive != *filter.IsActive {
			continue
		}
		cp := *u
		out = append(out, &cp)
		if limit > 0 && len(out) >= limit {
			break
		}
	}
	return out, nil, nil
}

func (f *fakeUserRepo) CountByRole(_ context.Context, role entity.Role) (int, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	count := 0
	for _, u := range f.byID {
		if u.Role == role && u.IsActive && u.DeletedAt == nil {
			count++
		}
	}
	return count, nil
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
func (f *fakeCourseRepo) List(_ context.Context, _ repository.CourseFilter, _ *pagination.Cursor, _ int) ([]*entity.Course, *pagination.Cursor, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	var out []*entity.Course
	for _, c := range f.byID {
		cp := *c
		out = append(out, &cp)
	}
	return out, nil, nil
}
func (f *fakeCourseRepo) Create(context.Context, *entity.Course) error { return nil }
func (f *fakeCourseRepo) Update(_ context.Context, c *entity.Course) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	cp := *c
	f.byID[c.ID] = &cp
	return nil
}
func (f *fakeCourseRepo) SoftDelete(_ context.Context, id uuid.UUID) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if c, ok := f.byID[id]; ok {
		now := time.Now()
		c.DeletedAt = &now
	}
	return nil
}
func (f *fakeCourseRepo) RecomputeCounters(context.Context, uuid.UUID) error { return nil }
func (f *fakeCourseRepo) GetCourseWithSectionsAndLessons(context.Context, uuid.UUID) (*repository.CourseWithContent, error) {
	return nil, apperror.NotFound("not found", nil)
}

// ---- fake payment repo ------------------------------------------------------

type fakePaymentRepo struct {
	mu   sync.Mutex
	byID map[uuid.UUID]*entity.Payment
}

func newFakePaymentRepo() *fakePaymentRepo {
	return &fakePaymentRepo{byID: map[uuid.UUID]*entity.Payment{}}
}

func (f *fakePaymentRepo) add(p *entity.Payment) {
	f.mu.Lock()
	defer f.mu.Unlock()
	cp := *p
	f.byID[p.ID] = &cp
}

func (f *fakePaymentRepo) Create(_ context.Context, p *entity.Payment) error {
	f.add(p)
	return nil
}
func (f *fakePaymentRepo) GetByID(_ context.Context, id uuid.UUID) (*entity.Payment, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if p, ok := f.byID[id]; ok {
		cp := *p
		return &cp, nil
	}
	return nil, apperror.NotFound("payment not found", nil)
}
func (f *fakePaymentRepo) GetByIdempotencyKey(context.Context, string) (*entity.Payment, error) {
	return nil, apperror.NotFound("not found", nil)
}
func (f *fakePaymentRepo) GetByProviderTxn(context.Context, entity.PaymentProvider, string) (*entity.Payment, error) {
	return nil, apperror.NotFound("not found", nil)
}
func (f *fakePaymentRepo) UpdateStatusGuarded(_ context.Context, id uuid.UUID, from, to entity.PaymentStatus, _ repository.UpdatePaymentStatusFields) (bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	p, ok := f.byID[id]
	if !ok {
		return false, nil
	}
	if p.Status != from {
		return false, nil
	}
	p.Status = to
	return true, nil
}
func (f *fakePaymentRepo) ListByStudent(context.Context, uuid.UUID, int) ([]*entity.Payment, error) {
	return nil, nil
}
func (f *fakePaymentRepo) ListAll(_ context.Context, filter repository.AdminPaymentFilter, _ *pagination.Cursor, _ int) ([]*entity.Payment, *pagination.Cursor, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	var out []*entity.Payment
	for _, p := range f.byID {
		if filter.Status != nil && string(p.Status) != *filter.Status {
			continue
		}
		cp := *p
		out = append(out, &cp)
	}
	return out, nil, nil
}

// ---- fake audit repo --------------------------------------------------------

type fakeAuditRepo struct {
	mu   sync.Mutex
	rows []*entity.AdminAuditLog
}

func (f *fakeAuditRepo) Insert(_ context.Context, l *entity.AdminAuditLog) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	cp := *l
	f.rows = append(f.rows, &cp)
	return nil
}

func (f *fakeAuditRepo) List(context.Context, repository.AdminAuditFilter, *pagination.Cursor, int) ([]*entity.AdminAuditLog, *pagination.Cursor, error) {
	return nil, nil, nil
}

func (f *fakeAuditRepo) countAction(action string) int {
	f.mu.Lock()
	defer f.mu.Unlock()
	n := 0
	for _, r := range f.rows {
		if r.Action == action {
			n++
		}
	}
	return n
}

// ---- fake revenue repo ------------------------------------------------------

type fakeRevenueRepo struct{}

func (fakeRevenueRepo) CompletedTotals(context.Context, time.Time, time.Time) (repository.RevenueTotals, error) {
	return repository.RevenueTotals{Gross: decimal.NewFromInt(1000), Units: 10}, nil
}
func (fakeRevenueRepo) RefundedTotals(context.Context, time.Time, time.Time) (repository.RevenueTotals, error) {
	return repository.RevenueTotals{Gross: decimal.NewFromInt(50), Units: 1}, nil
}
func (fakeRevenueRepo) TeacherGross(context.Context, uuid.UUID, time.Time, time.Time) (repository.RevenueTotals, error) {
	return repository.RevenueTotals{}, nil
}
func (fakeRevenueRepo) CourseSales(context.Context, uuid.UUID) (repository.RevenueTotals, error) {
	return repository.RevenueTotals{}, nil
}
func (fakeRevenueRepo) CompletedByProvider(context.Context, time.Time, time.Time) ([]repository.ProviderRevenue, error) {
	return nil, nil
}

// ---- fake subscription repo -------------------------------------------------

type fakeSubRepo struct{}

func (fakeSubRepo) Create(context.Context, *entity.Subscription) error { return nil }
func (fakeSubRepo) GetByID(context.Context, uuid.UUID) (*entity.Subscription, error) {
	return nil, apperror.NotFound("not found", nil)
}
func (fakeSubRepo) GetActiveForStudent(context.Context, uuid.UUID) (*entity.Subscription, error) {
	return nil, apperror.NotFound("not found", nil)
}
func (fakeSubRepo) UpdateStatus(context.Context, uuid.UUID, entity.SubscriptionStatus) error {
	return nil
}
func (fakeSubRepo) ListExpired(context.Context, time.Time, int) ([]*entity.Subscription, error) {
	return nil, nil
}
func (fakeSubRepo) CancelAtPeriodEnd(context.Context, uuid.UUID) (bool, error) { return false, nil }
func (fakeSubRepo) ExpireLapsed(context.Context, uuid.UUID) (bool, error)      { return false, nil }
func (fakeSubRepo) ExtendPeriod(context.Context, uuid.UUID, time.Time) error      { return nil }
func (fakeSubRepo) ListExpiring(context.Context, int) ([]*entity.Subscription, error) {
	return nil, nil
}

// ---- fake enrollment repo ---------------------------------------------------

type fakeEnrollRepo struct{}

func (fakeEnrollRepo) Create(context.Context, *entity.Enrollment) error { return nil }
func (fakeEnrollRepo) Exists(context.Context, uuid.UUID, uuid.UUID) (bool, error) {
	return false, nil
}
func (fakeEnrollRepo) GetActiveForStudent(context.Context, uuid.UUID, uuid.UUID) (*entity.Enrollment, error) {
	return nil, apperror.NotFound("not found", nil)
}
func (fakeEnrollRepo) ListByStudent(context.Context, uuid.UUID, int) ([]*entity.Enrollment, error) {
	return nil, nil
}

// ---- fake notifier ----------------------------------------------------------

type fakeNotifier struct {
	mu    sync.Mutex
	calls []service.NotifyEvent
}

func (n *fakeNotifier) Notify(_ context.Context, ev service.NotifyEvent) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.calls = append(n.calls, ev)
	return nil
}

func (n *fakeNotifier) count() int {
	n.mu.Lock()
	defer n.mu.Unlock()
	return len(n.calls)
}

// ---- compile-time checks ----------------------------------------------------

var _ repository.UserRepository = (*fakeUserRepo)(nil)
var _ repository.CourseRepository = (*fakeCourseRepo)(nil)
var _ repository.PaymentRepository = (*fakePaymentRepo)(nil)
var _ repository.AdminAuditRepository = (*fakeAuditRepo)(nil)
var _ repository.RevenueRepository = fakeRevenueRepo{}
var _ repository.SubscriptionRepository = fakeSubRepo{}
var _ repository.EnrollmentRepository = fakeEnrollRepo{}
var _ service.Notifier = (*fakeNotifier)(nil)

// ---- helpers ----------------------------------------------------------------

// newSvc builds a Service wired to all fakes.
func newSvc(users *fakeUserRepo, courses *fakeCourseRepo, payments *fakePaymentRepo, audit *fakeAuditRepo, notify *fakeNotifier) *Service {
	return NewService(users, courses, payments, fakeEnrollRepo{}, fakeSubRepo{}, fakeRevenueRepo{}, audit, notify, nil)
}

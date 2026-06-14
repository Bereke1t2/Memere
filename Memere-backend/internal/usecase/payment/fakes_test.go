package payment

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/service"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/pagination"
)

// ---- fake payment repo -------------------------------------------------------

type fakePaymentRepo struct {
	mu      sync.Mutex
	byID    map[uuid.UUID]*entity.Payment
	byKey   map[string]uuid.UUID
	creates int
}

func newFakePaymentRepo() *fakePaymentRepo {
	return &fakePaymentRepo{byID: map[uuid.UUID]*entity.Payment{}, byKey: map[string]uuid.UUID{}}
}

func (f *fakePaymentRepo) Create(_ context.Context, p *entity.Payment) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if p.IdempotencyKey != nil {
		if _, ok := f.byKey[*p.IdempotencyKey]; ok {
			return apperror.Conflict("PAYMENT_DUPLICATE", nil)
		}
	}
	f.creates++
	cp := *p
	f.byID[p.ID] = &cp
	if p.IdempotencyKey != nil {
		f.byKey[*p.IdempotencyKey] = p.ID
	}
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

func (f *fakePaymentRepo) GetByIdempotencyKey(_ context.Context, key string) (*entity.Payment, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if id, ok := f.byKey[key]; ok {
		cp := *f.byID[id]
		return &cp, nil
	}
	return nil, apperror.NotFound("payment not found", nil)
}

func (f *fakePaymentRepo) GetByProviderTxn(_ context.Context, provider entity.PaymentProvider, txnID string) (*entity.Payment, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, p := range f.byID {
		if p.Provider == provider && p.ProviderTxnID != nil && *p.ProviderTxnID == txnID {
			cp := *p
			return &cp, nil
		}
	}
	return nil, apperror.NotFound("payment not found", nil)
}

func (f *fakePaymentRepo) UpdateStatusGuarded(_ context.Context, id uuid.UUID, from, to entity.PaymentStatus, fields repository.UpdatePaymentStatusFields) (bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	p, ok := f.byID[id]
	if !ok {
		return false, apperror.NotFound("payment not found", nil)
	}
	if p.Status != from {
		return false, nil
	}
	p.Status = to
	if fields.ProviderTxnID != nil {
		p.ProviderTxnID = fields.ProviderTxnID
	}
	if fields.FailureReason != nil {
		p.FailureReason = fields.FailureReason
	}
	return true, nil
}

func (f *fakePaymentRepo) ListByStudent(_ context.Context, studentID uuid.UUID, limit int) ([]*entity.Payment, error) {
	return nil, nil
}

// ---- fake enrollment repo ----------------------------------------------------

type fakeEnrollRepo struct {
	mu      sync.Mutex
	active  map[[2]uuid.UUID]*entity.Enrollment
	creates int
}

func newFakeEnrollRepo() *fakeEnrollRepo {
	return &fakeEnrollRepo{active: map[[2]uuid.UUID]*entity.Enrollment{}}
}

func (f *fakeEnrollRepo) Create(_ context.Context, e *entity.Enrollment) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	key := [2]uuid.UUID{e.StudentID, e.CourseID}
	if _, ok := f.active[key]; ok {
		return apperror.Conflict("ENROLLMENT_DUPLICATE", nil)
	}
	f.creates++
	cp := *e
	f.active[key] = &cp
	return nil
}

func (f *fakeEnrollRepo) Exists(_ context.Context, s, c uuid.UUID) (bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	_, ok := f.active[[2]uuid.UUID{s, c}]
	return ok, nil
}

func (f *fakeEnrollRepo) GetActiveForStudent(_ context.Context, s, c uuid.UUID) (*entity.Enrollment, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if e, ok := f.active[[2]uuid.UUID{s, c}]; ok {
		return e, nil
	}
	return nil, apperror.NotFound("enrollment not found", nil)
}

func (f *fakeEnrollRepo) ListByStudent(_ context.Context, s uuid.UUID, limit int) ([]*entity.Enrollment, error) {
	return nil, nil
}

// ---- fake coupon repo --------------------------------------------------------

type fakeCouponRepo struct {
	mu         sync.Mutex
	byCode     map[string]*entity.Coupon
	byID       map[uuid.UUID]*entity.Coupon
	increments int
}

func newFakeCouponRepo() *fakeCouponRepo {
	return &fakeCouponRepo{byCode: map[string]*entity.Coupon{}, byID: map[uuid.UUID]*entity.Coupon{}}
}

func (f *fakeCouponRepo) add(c *entity.Coupon) {
	f.byCode[c.Code] = c
	f.byID[c.ID] = c
}

func (f *fakeCouponRepo) GetByCode(_ context.Context, code string) (*entity.Coupon, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if c, ok := f.byCode[code]; ok {
		cp := *c
		return &cp, nil
	}
	return nil, apperror.NotFound("coupon not found", nil)
}

func (f *fakeCouponRepo) IncrementUse(_ context.Context, id uuid.UUID) (bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	c, ok := f.byID[id]
	if !ok {
		return false, apperror.NotFound("coupon not found", nil)
	}
	if c.MaxUses != nil && c.UsedCount >= *c.MaxUses {
		return false, nil
	}
	c.UsedCount++
	f.increments++
	return true, nil
}

// ---- fake webhook ledger -----------------------------------------------------

type fakeWebhookRepo struct {
	mu        sync.Mutex
	seen      map[string]bool
	processed map[string]bool
}

func newFakeWebhookRepo() *fakeWebhookRepo {
	return &fakeWebhookRepo{seen: map[string]bool{}, processed: map[string]bool{}}
}

func (f *fakeWebhookRepo) key(provider entity.PaymentProvider, eventID string) string {
	return string(provider) + "|" + eventID
}

func (f *fakeWebhookRepo) InsertIfNew(_ context.Context, provider entity.PaymentProvider, eventID, eventType string, payload []byte) (bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	k := f.key(provider, eventID)
	if f.seen[k] {
		return false, nil
	}
	f.seen[k] = true
	return true, nil
}

func (f *fakeWebhookRepo) MarkProcessed(_ context.Context, provider entity.PaymentProvider, eventID string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.processed[f.key(provider, eventID)] = true
	return nil
}

// ---- fake course repo (only FindByID is exercised) ---------------------------

type fakeCourseRepo struct{ byID map[uuid.UUID]*entity.Course }

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
	return nil, apperror.NotFound("course not found", nil)
}

// ---- fake tx manager (runs fn inline; no real isolation needed for fakes) ----

type fakeTxManager struct{}

func (fakeTxManager) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

// ---- mock provider + registry ------------------------------------------------

type mockProvider struct {
	name        string
	checkoutURL string
	verifyErr   error
	verifyEvent *service.WebhookEvent
	verifyPay   string
}

func (m *mockProvider) Name() string { return m.name }
func (m *mockProvider) CreateCheckout(_ context.Context, req service.CheckoutRequest) (*service.CheckoutResult, error) {
	url := m.checkoutURL
	if url == "" {
		url = "https://pay.example/checkout/" + req.PaymentID
	}
	return &service.CheckoutResult{CheckoutID: req.PaymentID, RedirectURL: url}, nil
}
func (m *mockProvider) VerifyWebhook(_ context.Context, _ map[string]string, _ []byte) (*service.WebhookEvent, error) {
	if m.verifyErr != nil {
		return nil, m.verifyErr
	}
	return m.verifyEvent, nil
}
func (m *mockProvider) VerifyPayment(_ context.Context, _ string) (string, error) {
	if m.verifyPay == "" {
		return string(entity.PayPending), nil
	}
	return m.verifyPay, nil
}

type mockRegistry struct{ p service.PaymentProvider }

func (m mockRegistry) Get(entity.PaymentProvider) (service.PaymentProvider, error) {
	if m.p == nil {
		return nil, apperror.New(400, "PROVIDER_NOT_CONFIGURED", "no provider", nil)
	}
	return m.p, nil
}

// ---- coupon quoter shim (reuses the real coupon usecase over the fake repo) --

type quoterFunc func(ctx context.Context, code string, courseID uuid.UUID, base decimal.Decimal) (decimal.Decimal, uuid.UUID, error)

func (q quoterFunc) Quote(ctx context.Context, code string, courseID uuid.UUID, base decimal.Decimal) (decimal.Decimal, uuid.UUID, error) {
	return q(ctx, code, courseID, base)
}

// ---- subscription pricer + activator shims -----------------------------------

type stubPricer struct {
	price    decimal.Decimal
	currency string
	err      error
}

func (s stubPricer) PriceForPlan(string) (decimal.Decimal, string, error) {
	return s.price, s.currency, s.err
}

type fakeActivator struct {
	calls int
	last  *entity.Payment
}

func (f *fakeActivator) Activate(_ context.Context, p *entity.Payment) error {
	f.calls++
	f.last = p
	return nil
}

// ---- compile-time interface checks -------------------------------------------

var (
	_ repository.PaymentRepository      = (*fakePaymentRepo)(nil)
	_ repository.EnrollmentRepository   = (*fakeEnrollRepo)(nil)
	_ repository.CouponRepository       = (*fakeCouponRepo)(nil)
	_ repository.WebhookEventRepository = (*fakeWebhookRepo)(nil)
	_ repository.CourseRepository       = (*fakeCourseRepo)(nil)
	_ repository.TxManager              = fakeTxManager{}
	_ service.PaymentProvider           = (*mockProvider)(nil)
	_ service.PaymentProviderRegistry   = mockRegistry{}
	_ CouponQuoter                      = quoterFunc(nil)
	_ PlanPricer                        = stubPricer{}
	_ SubscriptionActivator             = (*fakeActivator)(nil)
)

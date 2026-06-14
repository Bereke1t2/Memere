package subscription

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// ---- stateful fake subscription repository -----------------------------------

// fakeSubRepo is an in-memory SubscriptionRepository. GetActiveForStudent mirrors
// the real query (status=active AND current_period_end > now), so the usecase's
// create-vs-extend branch is exercised faithfully under the fake clock.
type fakeSubRepo struct {
	byID map[uuid.UUID]*entity.Subscription
	now  time.Time

	createErr error
	extendErr error
	cancelOK  bool
	cancelErr error

	extendCalls int
	createCalls int
}

var _ repository.SubscriptionRepository = (*fakeSubRepo)(nil)

func newFakeSubRepo(now time.Time) *fakeSubRepo {
	return &fakeSubRepo{byID: map[uuid.UUID]*entity.Subscription{}, now: now, cancelOK: true}
}

func (f *fakeSubRepo) Create(_ context.Context, s *entity.Subscription) error {
	if f.createErr != nil {
		return f.createErr
	}
	f.createCalls++
	cp := *s
	f.byID[s.ID] = &cp
	return nil
}

func (f *fakeSubRepo) GetByID(_ context.Context, id uuid.UUID) (*entity.Subscription, error) {
	if s, ok := f.byID[id]; ok {
		return s, nil
	}
	return nil, apperror.NotFound("subscription not found", nil)
}

func (f *fakeSubRepo) GetActiveForStudent(_ context.Context, studentID uuid.UUID) (*entity.Subscription, error) {
	for _, s := range f.byID {
		if s.StudentID == studentID && s.Status == entity.SubActive && s.CurrentPeriodEnd.After(f.now) {
			return s, nil
		}
	}
	return nil, apperror.NotFound("subscription not found", nil)
}

func (f *fakeSubRepo) UpdateStatus(_ context.Context, id uuid.UUID, status entity.SubscriptionStatus) error {
	if s, ok := f.byID[id]; ok {
		s.Status = status
	}
	return nil
}

func (f *fakeSubRepo) ExtendPeriod(_ context.Context, id uuid.UUID, newEnd time.Time) error {
	if f.extendErr != nil {
		return f.extendErr
	}
	f.extendCalls++
	if s, ok := f.byID[id]; ok {
		s.CurrentPeriodEnd = newEnd
		s.Status = entity.SubActive
	}
	return nil
}

func (f *fakeSubRepo) CancelAtPeriodEnd(_ context.Context, id uuid.UUID) (bool, error) {
	if f.cancelErr != nil {
		return false, f.cancelErr
	}
	return f.cancelOK, nil
}

func (f *fakeSubRepo) ExpireLapsed(_ context.Context, id uuid.UUID) (bool, error) {
	s, ok := f.byID[id]
	if !ok || s.Status != entity.SubActive || s.CurrentPeriodEnd.After(f.now) {
		return false, nil
	}
	s.Status = entity.SubExpired
	return true, nil
}

func (f *fakeSubRepo) ListExpiring(_ context.Context, limit int) ([]*entity.Subscription, error) {
	var out []*entity.Subscription
	for _, s := range f.byID {
		if s.Status == entity.SubActive && !s.CurrentPeriodEnd.After(f.now) {
			out = append(out, s)
		}
	}
	if limit > 0 && len(out) > limit {
		out = out[:limit]
	}
	return out, nil
}

// ---- helpers -----------------------------------------------------------------

func testConfig() Config {
	return Config{
		Monthly: PlanSpec{Price: decimal.NewFromInt(200), Currency: "ETB", Period: 30 * 24 * time.Hour},
		Annual:  PlanSpec{Price: decimal.NewFromInt(2000), Currency: "ETB", Period: 365 * 24 * time.Hour},
	}
}

func newSvc(repo *fakeSubRepo, now time.Time) *Service {
	return NewService(repo, testConfig(), func() time.Time { return now })
}

func subPayment(studentID uuid.UUID, plan entity.SubscriptionPlan) *entity.Payment {
	return &entity.Payment{
		ID:        uuid.New(),
		StudentID: studentID,
		Provider:  entity.ProviderChapa,
		Status:    entity.PayCompleted,
		Metadata:  map[string]any{metaPlan: string(plan)},
	}
}

// ---- PriceForPlan ------------------------------------------------------------

func TestPriceForPlan(t *testing.T) {
	now := time.Date(2026, 6, 14, 12, 0, 0, 0, time.UTC)
	svc := newSvc(newFakeSubRepo(now), now)

	price, currency, err := svc.PriceForPlan("monthly")
	if err != nil {
		t.Fatalf("monthly: unexpected error: %v", err)
	}
	if !price.Equal(decimal.NewFromInt(200)) || currency != "ETB" {
		t.Fatalf("monthly: got %s %s, want 200 ETB", price, currency)
	}

	if _, _, err := svc.PriceForPlan("weekly"); !apperror.IsCode(err, "UNKNOWN_PLAN") {
		t.Fatalf("unknown plan: want UNKNOWN_PLAN, got %v", err)
	}
}

// ---- Activate: create --------------------------------------------------------

func TestActivate_CreatesWhenNone(t *testing.T) {
	now := time.Date(2026, 6, 14, 12, 0, 0, 0, time.UTC)
	repo := newFakeSubRepo(now)
	svc := newSvc(repo, now)
	student := uuid.New()

	if err := svc.Activate(context.Background(), subPayment(student, entity.PlanMonthly)); err != nil {
		t.Fatalf("activate: %v", err)
	}
	if repo.createCalls != 1 {
		t.Fatalf("want 1 create, got %d", repo.createCalls)
	}
	sub, err := repo.GetActiveForStudent(context.Background(), student)
	if err != nil {
		t.Fatalf("expected an active subscription: %v", err)
	}
	wantEnd := now.Add(30 * 24 * time.Hour)
	if !sub.CurrentPeriodEnd.Equal(wantEnd) {
		t.Fatalf("period end: got %s, want %s", sub.CurrentPeriodEnd, wantEnd)
	}
	if sub.Plan != entity.PlanMonthly || sub.Status != entity.SubActive {
		t.Fatalf("got plan=%s status=%s, want monthly/active", sub.Plan, sub.Status)
	}
}

// ---- Activate: extend (idempotent stacking) ----------------------------------

func TestActivate_ExtendsExisting(t *testing.T) {
	now := time.Date(2026, 6, 14, 12, 0, 0, 0, time.UTC)
	repo := newFakeSubRepo(now)
	svc := newSvc(repo, now)
	student := uuid.New()

	// Seed an active subscription ending 10 days out.
	existingEnd := now.Add(10 * 24 * time.Hour)
	id := uuid.New()
	repo.byID[id] = &entity.Subscription{
		ID: id, StudentID: student, Plan: entity.PlanMonthly, Status: entity.SubActive,
		CurrentPeriodStart: now.Add(-20 * 24 * time.Hour), CurrentPeriodEnd: existingEnd,
	}

	if err := svc.Activate(context.Background(), subPayment(student, entity.PlanMonthly)); err != nil {
		t.Fatalf("activate: %v", err)
	}
	if repo.createCalls != 0 {
		t.Fatalf("expected no new row, got %d creates", repo.createCalls)
	}
	if repo.extendCalls != 1 {
		t.Fatalf("expected 1 extend, got %d", repo.extendCalls)
	}
	// Stacks from the current end, not from now.
	wantEnd := existingEnd.Add(30 * 24 * time.Hour)
	if got := repo.byID[id].CurrentPeriodEnd; !got.Equal(wantEnd) {
		t.Fatalf("extended end: got %s, want %s", got, wantEnd)
	}
}

// ---- Activate: unknown plan rejected -----------------------------------------

func TestActivate_UnknownPlan(t *testing.T) {
	now := time.Date(2026, 6, 14, 12, 0, 0, 0, time.UTC)
	repo := newFakeSubRepo(now)
	svc := newSvc(repo, now)

	p := subPayment(uuid.New(), entity.SubscriptionPlan("lifetime"))
	if err := svc.Activate(context.Background(), p); !apperror.IsCode(err, "UNKNOWN_PLAN") {
		t.Fatalf("want UNKNOWN_PLAN, got %v", err)
	}
}

// ---- CancelSubscription ------------------------------------------------------

func TestCancelSubscription(t *testing.T) {
	now := time.Date(2026, 6, 14, 12, 0, 0, 0, time.UTC)
	student := uuid.New()
	other := uuid.New()
	admin := uuid.New()

	seed := func(repo *fakeSubRepo) uuid.UUID {
		id := uuid.New()
		repo.byID[id] = &entity.Subscription{
			ID: id, StudentID: student, Plan: entity.PlanMonthly, Status: entity.SubActive,
			CurrentPeriodEnd: now.Add(15 * 24 * time.Hour),
		}
		return id
	}

	t.Run("owner cancels", func(t *testing.T) {
		repo := newFakeSubRepo(now)
		id := seed(repo)
		svc := newSvc(repo, now)
		err := svc.CancelSubscription(context.Background(), Actor{UserID: student, Role: entity.RoleStudent}, id)
		if err != nil {
			t.Fatalf("owner cancel: %v", err)
		}
	})

	t.Run("admin cancels any", func(t *testing.T) {
		repo := newFakeSubRepo(now)
		id := seed(repo)
		svc := newSvc(repo, now)
		err := svc.CancelSubscription(context.Background(), Actor{UserID: admin, Role: entity.RoleAdmin}, id)
		if err != nil {
			t.Fatalf("admin cancel: %v", err)
		}
	})

	t.Run("non-owner forbidden", func(t *testing.T) {
		repo := newFakeSubRepo(now)
		id := seed(repo)
		svc := newSvc(repo, now)
		err := svc.CancelSubscription(context.Background(), Actor{UserID: other, Role: entity.RoleStudent}, id)
		if !apperror.IsCode(err, "FORBIDDEN") {
			t.Fatalf("non-owner: want FORBIDDEN, got %v", err)
		}
	})

	t.Run("unauthenticated", func(t *testing.T) {
		repo := newFakeSubRepo(now)
		id := seed(repo)
		svc := newSvc(repo, now)
		err := svc.CancelSubscription(context.Background(), Actor{}, id)
		if !apperror.IsCode(err, "UNAUTHORIZED") {
			t.Fatalf("anonymous: want UNAUTHORIZED, got %v", err)
		}
	})

	t.Run("not cancelable", func(t *testing.T) {
		repo := newFakeSubRepo(now)
		repo.cancelOK = false
		id := seed(repo)
		svc := newSvc(repo, now)
		err := svc.CancelSubscription(context.Background(), Actor{UserID: student, Role: entity.RoleStudent}, id)
		if !apperror.IsCode(err, "NOT_CANCELABLE") {
			t.Fatalf("want NOT_CANCELABLE, got %v", err)
		}
	})
}

// ---- SweepExpired ------------------------------------------------------------

func TestSweepExpired(t *testing.T) {
	now := time.Date(2026, 6, 14, 12, 0, 0, 0, time.UTC)
	repo := newFakeSubRepo(now)

	// One lapsed active sub (expires), one still-current active sub (untouched).
	lapsedID, currentID := uuid.New(), uuid.New()
	repo.byID[lapsedID] = &entity.Subscription{
		ID: lapsedID, StudentID: uuid.New(), Status: entity.SubActive,
		CurrentPeriodEnd: now.Add(-24 * time.Hour),
	}
	repo.byID[currentID] = &entity.Subscription{
		ID: currentID, StudentID: uuid.New(), Status: entity.SubActive,
		CurrentPeriodEnd: now.Add(24 * time.Hour),
	}
	svc := newSvc(repo, now)

	expired, err := svc.SweepExpired(context.Background(), 100)
	if err != nil {
		t.Fatalf("sweep: %v", err)
	}
	if len(expired) != 1 || expired[0].ID != lapsedID {
		t.Fatalf("want only the lapsed sub expired, got %+v", expired)
	}
	if repo.byID[lapsedID].Status != entity.SubExpired {
		t.Fatalf("lapsed sub should be expired, got %s", repo.byID[lapsedID].Status)
	}
	if repo.byID[currentID].Status != entity.SubActive {
		t.Fatalf("current sub must stay active, got %s", repo.byID[currentID].Status)
	}

	// Idempotent: a second sweep finds nothing (the row is no longer active).
	again, err := svc.SweepExpired(context.Background(), 100)
	if err != nil {
		t.Fatalf("second sweep: %v", err)
	}
	if len(again) != 0 {
		t.Fatalf("second sweep should expire nothing, got %d", len(again))
	}
}

// ---- GetMySubscription -------------------------------------------------------

func TestGetMySubscription(t *testing.T) {
	now := time.Date(2026, 6, 14, 12, 0, 0, 0, time.UTC)
	repo := newFakeSubRepo(now)
	svc := newSvc(repo, now)
	student := uuid.New()

	if _, err := svc.GetMySubscription(context.Background(), Actor{}); !apperror.IsCode(err, "UNAUTHORIZED") {
		t.Fatalf("anonymous: want UNAUTHORIZED, got %v", err)
	}

	if _, err := svc.GetMySubscription(context.Background(), Actor{UserID: student, Role: entity.RoleStudent}); !apperror.IsNotFound(err) {
		t.Fatalf("no subscription: want NotFound, got %v", err)
	}

	if err := svc.Activate(context.Background(), subPayment(student, entity.PlanAnnual)); err != nil {
		t.Fatalf("activate: %v", err)
	}
	sub, err := svc.GetMySubscription(context.Background(), Actor{UserID: student, Role: entity.RoleStudent})
	if err != nil {
		t.Fatalf("get mine: %v", err)
	}
	if sub.Plan != entity.PlanAnnual {
		t.Fatalf("got plan %s, want annual", sub.Plan)
	}
}

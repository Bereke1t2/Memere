package coupon

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

type fakeCouponRepo struct{ byCode map[string]*entity.Coupon }

func (f *fakeCouponRepo) GetByCode(_ context.Context, code string) (*entity.Coupon, error) {
	if c, ok := f.byCode[code]; ok {
		cp := *c
		return &cp, nil
	}
	return nil, apperror.NotFound("coupon not found", nil)
}
func (f *fakeCouponRepo) IncrementUse(context.Context, uuid.UUID) (bool, error) { return true, nil }

func newService(c *entity.Coupon, now time.Time) *Service {
	repo := &fakeCouponRepo{byCode: map[string]*entity.Coupon{}}
	if c != nil {
		repo.byCode[c.Code] = c
	}
	return NewService(repo, func() time.Time { return now })
}

func TestQuote_PercentageDiscount(t *testing.T) {
	now := time.Date(2026, 6, 12, 0, 0, 0, 0, time.UTC)
	c := &entity.Coupon{
		ID: uuid.New(), Code: "HALF", DiscountType: entity.DiscountPercentage,
		DiscountValue: decimal.NewFromInt(50), ApplicableTo: entity.ScopeAll,
	}
	s := newService(c, now)

	final, id, err := s.Quote(context.Background(), "HALF", uuid.New(), decimal.NewFromInt(200))
	if err != nil {
		t.Fatalf("quote: %v", err)
	}
	if !final.Equal(decimal.NewFromInt(100)) {
		t.Fatalf("50%% off 200 should be 100, got %s", final)
	}
	if id != c.ID {
		t.Fatalf("coupon id mismatch")
	}
}

func TestQuote_Expired(t *testing.T) {
	now := time.Date(2026, 6, 12, 0, 0, 0, 0, time.UTC)
	past := now.Add(-time.Hour)
	c := &entity.Coupon{
		ID: uuid.New(), Code: "OLD", DiscountType: entity.DiscountPercentage,
		DiscountValue: decimal.NewFromInt(10), ApplicableTo: entity.ScopeAll, ExpiresAt: &past,
	}
	s := newService(c, now)

	_, _, err := s.Quote(context.Background(), "OLD", uuid.New(), decimal.NewFromInt(100))
	if !apperror.IsCode(err, "COUPON_EXPIRED") {
		t.Fatalf("expected COUPON_EXPIRED, got %v", err)
	}
}

func TestQuote_OverUse(t *testing.T) {
	now := time.Date(2026, 6, 12, 0, 0, 0, 0, time.UTC)
	max := 1
	c := &entity.Coupon{
		ID: uuid.New(), Code: "USED", DiscountType: entity.DiscountFixed,
		DiscountValue: decimal.NewFromInt(5), ApplicableTo: entity.ScopeAll,
		MaxUses: &max, UsedCount: 1,
	}
	s := newService(c, now)

	_, _, err := s.Quote(context.Background(), "USED", uuid.New(), decimal.NewFromInt(100))
	if !apperror.IsCode(err, "COUPON_EXHAUSTED") {
		t.Fatalf("expected COUPON_EXHAUSTED, got %v", err)
	}
}

func TestQuote_ScopeMismatch(t *testing.T) {
	now := time.Date(2026, 6, 12, 0, 0, 0, 0, time.UTC)
	other := uuid.New()
	c := &entity.Coupon{
		ID: uuid.New(), Code: "SCOPED", DiscountType: entity.DiscountPercentage,
		DiscountValue: decimal.NewFromInt(10), ApplicableTo: entity.ScopeSpecificCourses,
		CourseIDs: []uuid.UUID{other}, // not the course we quote against
	}
	s := newService(c, now)

	_, _, err := s.Quote(context.Background(), "SCOPED", uuid.New(), decimal.NewFromInt(100))
	if !apperror.IsCode(err, "COUPON_NOT_APPLICABLE") {
		t.Fatalf("expected COUPON_NOT_APPLICABLE, got %v", err)
	}
}

func TestQuote_UnknownCode(t *testing.T) {
	now := time.Date(2026, 6, 12, 0, 0, 0, 0, time.UTC)
	s := newService(nil, now)

	_, _, err := s.Quote(context.Background(), "NOPE", uuid.New(), decimal.NewFromInt(100))
	if !apperror.IsNotFound(err) {
		t.Fatalf("expected NotFound for unknown code, got %v", err)
	}
}

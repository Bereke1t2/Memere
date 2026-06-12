package entity

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// TestCouponApply covers the four validation gates (expiry, max-use, scope) and
// both discount math paths, including the clamp-to-zero when a fixed discount
// exceeds the price. Apply must never mutate UsedCount.
func TestCouponApply(t *testing.T) {
	now := time.Date(2026, 6, 12, 0, 0, 0, 0, time.UTC)
	past := now.Add(-time.Hour)
	future := now.Add(time.Hour)
	courseA := uuid.New()
	courseB := uuid.New()
	dec := decimal.RequireFromString
	maxUses := func(n int) *int { return &n }

	tests := []struct {
		name    string
		coupon  Coupon
		price   decimal.Decimal
		target  CouponTarget
		want    string // expected discounted price, "" when wantErr
		wantErr error
	}{
		{
			name:   "percentage off all",
			coupon: Coupon{DiscountType: DiscountPercentage, DiscountValue: dec("50"), ApplicableTo: ScopeAll},
			price:  dec("200"), target: CouponTarget{CourseID: &courseA}, want: "100",
		},
		{
			name:   "fixed amount off all",
			coupon: Coupon{DiscountType: DiscountFixed, DiscountValue: dec("30"), ApplicableTo: ScopeAll},
			price:  dec("200"), target: CouponTarget{CourseID: &courseA}, want: "170",
		},
		{
			name:   "fixed discount larger than price clamps to zero",
			coupon: Coupon{DiscountType: DiscountFixed, DiscountValue: dec("500"), ApplicableTo: ScopeAll},
			price:  dec("200"), target: CouponTarget{CourseID: &courseA}, want: "0",
		},
		{
			name:   "expired",
			coupon: Coupon{DiscountType: DiscountPercentage, DiscountValue: dec("50"), ApplicableTo: ScopeAll, ExpiresAt: &past},
			price:  dec("200"), target: CouponTarget{CourseID: &courseA}, wantErr: ErrCouponExpired,
		},
		{
			name:   "not yet expired is allowed",
			coupon: Coupon{DiscountType: DiscountPercentage, DiscountValue: dec("10"), ApplicableTo: ScopeAll, ExpiresAt: &future},
			price:  dec("200"), target: CouponTarget{CourseID: &courseA}, want: "180",
		},
		{
			name:   "over max use",
			coupon: Coupon{DiscountType: DiscountPercentage, DiscountValue: dec("50"), ApplicableTo: ScopeAll, MaxUses: maxUses(5), UsedCount: 5},
			price:  dec("200"), target: CouponTarget{CourseID: &courseA}, wantErr: ErrCouponExhausted,
		},
		{
			name:   "specific course in scope",
			coupon: Coupon{DiscountType: DiscountPercentage, DiscountValue: dec("25"), ApplicableTo: ScopeSpecificCourses, CourseIDs: []uuid.UUID{courseA}},
			price:  dec("200"), target: CouponTarget{CourseID: &courseA}, want: "150",
		},
		{
			name:   "specific course out of scope",
			coupon: Coupon{DiscountType: DiscountPercentage, DiscountValue: dec("25"), ApplicableTo: ScopeSpecificCourses, CourseIDs: []uuid.UUID{courseA}},
			price:  dec("200"), target: CouponTarget{CourseID: &courseB}, wantErr: ErrCouponNotInScope,
		},
		{
			name:   "subscription-only rejects a course purchase",
			coupon: Coupon{DiscountType: DiscountPercentage, DiscountValue: dec("25"), ApplicableTo: ScopeSubscriptionOnly},
			price:  dec("200"), target: CouponTarget{CourseID: &courseA}, wantErr: ErrCouponNotInScope,
		},
		{
			name:   "subscription-only allows a subscription purchase",
			coupon: Coupon{DiscountType: DiscountPercentage, DiscountValue: dec("25"), ApplicableTo: ScopeSubscriptionOnly},
			price:  dec("400"), target: CouponTarget{IsSubscription: true}, want: "300",
		},
		{
			name:   "unknown discount type",
			coupon: Coupon{DiscountType: DiscountType("bogus"), DiscountValue: dec("25"), ApplicableTo: ScopeAll},
			price:  dec("200"), target: CouponTarget{CourseID: &courseA}, wantErr: ErrCouponBadDiscount,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.coupon.Apply(tt.price, now, tt.target)
			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Fatalf("Apply error = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Apply unexpected error: %v", err)
			}
			if !got.Equal(dec(tt.want)) {
				t.Fatalf("Apply = %s, want %s", got, tt.want)
			}
			if tt.coupon.UsedCount != 0 && tt.coupon.UsedCount != 5 {
				t.Fatalf("Apply must not mutate UsedCount, got %d", tt.coupon.UsedCount)
			}
		})
	}
}

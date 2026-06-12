package entity

import (
	"errors"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// DiscountType selects how a coupon's value reduces a price (spec §10.3).
type DiscountType string

const (
	DiscountPercentage DiscountType = "percentage"   // value is a percent: 50 => 50% off
	DiscountFixed      DiscountType = "fixed_amount" // value is an absolute amount off
)

// CouponScope limits which purchases a coupon applies to (spec §10.3:
// applicable_to ENUM all / specific_courses / subscription_only).
type CouponScope string

const (
	ScopeAll              CouponScope = "all"
	ScopeSpecificCourses  CouponScope = "specific_courses"
	ScopeSubscriptionOnly CouponScope = "subscription_only"
)

// Coupon validation errors returned by Apply. Usecases map these to HTTP 400/409
// via apperror; keeping them as sentinels here keeps the domain free of delivery
// concerns.
var (
	ErrCouponExpired     = errors.New("coupon has expired")
	ErrCouponExhausted   = errors.New("coupon has reached its usage limit")
	ErrCouponNotInScope  = errors.New("coupon does not apply to this purchase")
	ErrCouponBadDiscount = errors.New("coupon has an invalid discount configuration")
)

// Coupon is a discount code (spec §10.3). DiscountValue is decimal money for a
// fixed_amount coupon and a percent (0–100) for a percentage coupon. CourseIDs
// scopes a specific_courses coupon.
type Coupon struct {
	ID            uuid.UUID
	Code          string
	DiscountType  DiscountType
	DiscountValue decimal.Decimal
	MaxUses       *int // nil => unlimited
	UsedCount     int
	ExpiresAt     *time.Time
	ApplicableTo  CouponScope
	CourseIDs     []uuid.UUID
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// CouponTarget describes the purchase a coupon is being applied to so Apply can
// validate scope. CourseID is nil for a subscription purchase.
type CouponTarget struct {
	CourseID       *uuid.UUID
	IsSubscription bool
}

// Apply validates the coupon against now and target, then returns the discounted
// price. It never returns a negative price (a fixed discount larger than the
// price clamps to zero) and never mutates UsedCount — incrementing the use count
// happens only on a *successful* payment, inside the fulfillment transaction
// (spec §10.3). price must be non-negative.
func (c *Coupon) Apply(price decimal.Decimal, now time.Time, target CouponTarget) (decimal.Decimal, error) {
	if c.ExpiresAt != nil && !c.ExpiresAt.After(now) {
		return decimal.Zero, ErrCouponExpired
	}
	if c.MaxUses != nil && c.UsedCount >= *c.MaxUses {
		return decimal.Zero, ErrCouponExhausted
	}
	if !c.appliesTo(target) {
		return decimal.Zero, ErrCouponNotInScope
	}

	var discounted decimal.Decimal
	switch c.DiscountType {
	case DiscountPercentage:
		// percent off: price * (1 - value/100)
		factor := decimal.NewFromInt(1).Sub(c.DiscountValue.Div(decimal.NewFromInt(100)))
		discounted = price.Mul(factor)
	case DiscountFixed:
		discounted = price.Sub(c.DiscountValue)
	default:
		return decimal.Zero, ErrCouponBadDiscount
	}

	if discounted.IsNegative() {
		return decimal.Zero, nil
	}
	// Round to 2 places — money is stored as DECIMAL(10,2).
	return discounted.Round(2), nil
}

// appliesTo reports whether the coupon's scope admits the given purchase.
func (c *Coupon) appliesTo(target CouponTarget) bool {
	switch c.ApplicableTo {
	case ScopeAll:
		return true
	case ScopeSubscriptionOnly:
		return target.IsSubscription
	case ScopeSpecificCourses:
		if target.CourseID == nil {
			return false
		}
		return slices.Contains(c.CourseIDs, *target.CourseID)
	default:
		return false
	}
}

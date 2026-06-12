// Package coupon holds the read-only coupon usecase used during checkout: given
// a code and a purchase, it validates the coupon and computes the discounted
// price. It deliberately does NOT increment the coupon's use count — that
// happens atomically inside the successful-payment fulfillment transaction
// (Skill 3) so an abandoned checkout never burns a coupon use (spec §10.3).
package coupon

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// Service quotes coupon discounts over the coupon repository. clock is
// injectable so tests can drive expiry deterministically.
type Service struct {
	coupons repository.CouponRepository
	clock   func() time.Time
}

// NewService wires the coupon service. When clock is nil it defaults to
// time.Now.
func NewService(coupons repository.CouponRepository, clock func() time.Time) *Service {
	if clock == nil {
		clock = time.Now
	}
	return &Service{coupons: coupons, clock: clock}
}

// Quote validates the coupon code against a course purchase of basePrice and
// returns the discounted final price plus the coupon id (so the payment path
// can record and later increment it). It returns an apperror for an unknown,
// expired, exhausted, or out-of-scope coupon and never mutates use counts.
func (s *Service) Quote(
	ctx context.Context,
	code string,
	courseID uuid.UUID,
	basePrice decimal.Decimal,
) (final decimal.Decimal, couponID uuid.UUID, err error) {
	c, err := s.coupons.GetByCode(ctx, code)
	if err != nil {
		return decimal.Zero, uuid.Nil, err
	}
	target := entity.CouponTarget{CourseID: &courseID, IsSubscription: false}
	final, applyErr := c.Apply(basePrice, s.clock(), target)
	if applyErr != nil {
		return decimal.Zero, uuid.Nil, mapCouponError(applyErr)
	}
	return final, c.ID, nil
}

// mapCouponError translates the domain coupon sentinels into delivery-facing
// AppErrors, keeping the domain free of HTTP concerns.
func mapCouponError(err error) error {
	switch {
	case errors.Is(err, entity.ErrCouponExpired):
		return apperror.New(409, "COUPON_EXPIRED", "this coupon has expired", err)
	case errors.Is(err, entity.ErrCouponExhausted):
		return apperror.New(409, "COUPON_EXHAUSTED", "this coupon has reached its usage limit", err)
	case errors.Is(err, entity.ErrCouponNotInScope):
		return apperror.New(409, "COUPON_NOT_APPLICABLE", "this coupon does not apply to this purchase", err)
	case errors.Is(err, entity.ErrCouponBadDiscount):
		return apperror.New(500, "COUPON_MISCONFIGURED", "this coupon is misconfigured", err)
	default:
		return apperror.New(400, "COUPON_INVALID", "this coupon cannot be applied", err)
	}
}

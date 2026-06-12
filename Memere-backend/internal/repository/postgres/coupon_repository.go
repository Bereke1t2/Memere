package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/internal/repository/postgres/sqlcgen"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// CouponRepo is the sqlc-backed implementation of repository.CouponRepository.
type CouponRepo struct {
	q *sqlcgen.Queries
}

var _ repository.CouponRepository = (*CouponRepo)(nil)

// NewCouponRepo builds a CouponRepo over a pgx pool.
func NewCouponRepo(pool *pgxpool.Pool) *CouponRepo {
	return &CouponRepo{q: sqlcgen.New(pool)}
}

func (r *CouponRepo) GetByCode(ctx context.Context, code string) (*entity.Coupon, error) {
	row, err := queriesFor(ctx, r.q).GetCouponByCode(ctx, code)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NotFound("coupon not found", err)
		}
		return nil, apperror.Internal(err)
	}
	return couponFromRow(row), nil
}

// IncrementUse bumps used_count only while it is still below max_uses. A 0-row
// result means the coupon is already exhausted (or lost a concurrent race):
// reported as redeemed=false, not an error, so the caller can fail the
// fulfillment cleanly.
func (r *CouponRepo) IncrementUse(ctx context.Context, id uuid.UUID) (bool, error) {
	n, err := queriesFor(ctx, r.q).IncrementCouponUse(ctx, toPgUUID(id))
	if err != nil {
		return false, apperror.Internal(err)
	}
	return n > 0, nil
}

func couponFromRow(row sqlcgen.PaymentsCoupon) *entity.Coupon {
	courseIDs := make([]uuid.UUID, 0, len(row.CourseIds))
	for _, id := range row.CourseIds {
		courseIDs = append(courseIDs, fromPgUUID(id))
	}
	return &entity.Coupon{
		ID:            fromPgUUID(row.ID),
		Code:          row.Code,
		DiscountType:  entity.DiscountType(row.DiscountType),
		DiscountValue: fromPgDecimal(row.DiscountValue),
		MaxUses:       fromPgInt4Ptr(row.MaxUses),
		UsedCount:     int(row.UsedCount),
		ExpiresAt:     fromPgTimestamptz(row.ExpiresAt),
		ApplicableTo:  entity.CouponScope(row.ApplicableTo),
		CourseIDs:     courseIDs,
		CreatedAt:     fromPgTimestamptzValue(row.CreatedAt),
		UpdatedAt:     fromPgTimestamptzValue(row.UpdatedAt),
	}
}

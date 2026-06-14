package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/internal/repository/postgres/sqlcgen"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// RevenueRepo is the sqlc-backed implementation of repository.RevenueRepository:
// read-only aggregations over completed payments (spec §1.6, §10).
type RevenueRepo struct {
	q *sqlcgen.Queries
}

var _ repository.RevenueRepository = (*RevenueRepo)(nil)

// NewRevenueRepo builds a RevenueRepo over a pgx pool.
func NewRevenueRepo(pool *pgxpool.Pool) *RevenueRepo {
	return &RevenueRepo{q: sqlcgen.New(pool)}
}

func (r *RevenueRepo) CompletedTotals(ctx context.Context, from, to time.Time) (repository.RevenueTotals, error) {
	row, err := queriesFor(ctx, r.q).SumCompletedPayments(ctx, sqlcgen.SumCompletedPaymentsParams{
		PaidAt:   pgTimestamptzValue(from),
		PaidAt_2: pgTimestamptzValue(to),
	})
	if err != nil {
		return repository.RevenueTotals{}, apperror.Internal(err)
	}
	return repository.RevenueTotals{Gross: fromPgDecimal(row.Gross), Units: row.Units}, nil
}

func (r *RevenueRepo) RefundedTotals(ctx context.Context, from, to time.Time) (repository.RevenueTotals, error) {
	row, err := queriesFor(ctx, r.q).SumRefundedPayments(ctx, sqlcgen.SumRefundedPaymentsParams{
		UpdatedAt:   pgTimestamptzValue(from),
		UpdatedAt_2: pgTimestamptzValue(to),
	})
	if err != nil {
		return repository.RevenueTotals{}, apperror.Internal(err)
	}
	return repository.RevenueTotals{Gross: fromPgDecimal(row.Refunded), Units: row.Units}, nil
}

func (r *RevenueRepo) TeacherGross(ctx context.Context, teacherID uuid.UUID, from, to time.Time) (repository.RevenueTotals, error) {
	row, err := queriesFor(ctx, r.q).TeacherGrossEarnings(ctx, sqlcgen.TeacherGrossEarningsParams{
		TeacherID: toPgUUID(teacherID),
		PaidAt:    pgTimestamptzValue(from),
		PaidAt_2:  pgTimestamptzValue(to),
	})
	if err != nil {
		return repository.RevenueTotals{}, apperror.Internal(err)
	}
	return repository.RevenueTotals{Gross: fromPgDecimal(row.Gross), Units: row.Units}, nil
}

func (r *RevenueRepo) CourseSales(ctx context.Context, courseID uuid.UUID) (repository.RevenueTotals, error) {
	row, err := queriesFor(ctx, r.q).CourseSales(ctx, toPgUUID(courseID))
	if err != nil {
		return repository.RevenueTotals{}, apperror.Internal(err)
	}
	return repository.RevenueTotals{Gross: fromPgDecimal(row.Gross), Units: row.Units}, nil
}

func (r *RevenueRepo) CompletedByProvider(ctx context.Context, from, to time.Time) ([]repository.ProviderRevenue, error) {
	rows, err := queriesFor(ctx, r.q).CompletedRevenueByProvider(ctx, sqlcgen.CompletedRevenueByProviderParams{
		PaidAt:   pgTimestamptzValue(from),
		PaidAt_2: pgTimestamptzValue(to),
	})
	if err != nil {
		return nil, apperror.Internal(err)
	}
	out := make([]repository.ProviderRevenue, len(rows))
	for i, row := range rows {
		out[i] = repository.ProviderRevenue{
			Provider: entity.PaymentProvider(row.Provider),
			Gross:    fromPgDecimal(row.Gross),
			Units:    row.Units,
		}
	}
	return out, nil
}

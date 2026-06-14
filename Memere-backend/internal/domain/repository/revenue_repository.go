package repository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/shopspring/decimal"
)

// RevenueTotals is a gross money sum plus the unit count behind it.
type RevenueTotals struct {
	Gross decimal.Decimal
	Units int64
}

// ProviderRevenue is a per-provider slice of completed gross.
type ProviderRevenue struct {
	Provider entity.PaymentProvider
	Gross    decimal.Decimal
	Units    int64
}

// RevenueRepository serves the read-only revenue aggregations over completed
// payments (spec §1.6, §10). It performs no authorization — the usecase layer
// applies role/ownership scoping and the configurable teacher/platform split.
type RevenueRepository interface {
	// CompletedTotals is the platform gross + count settled in [from, to].
	CompletedTotals(ctx context.Context, from, to time.Time) (RevenueTotals, error)
	// RefundedTotals is the refunded gross + count in [from, to].
	RefundedTotals(ctx context.Context, from, to time.Time) (RevenueTotals, error)
	// TeacherGross sums completed payments for the teacher's courses in [from, to]
	// (before the platform fee).
	TeacherGross(ctx context.Context, teacherID uuid.UUID, from, to time.Time) (RevenueTotals, error)
	// CourseSales is the lifetime gross + units sold for one course.
	CourseSales(ctx context.Context, courseID uuid.UUID) (RevenueTotals, error)
	// CompletedByProvider breaks the platform gross down per provider in [from, to].
	CompletedByProvider(ctx context.Context, from, to time.Time) ([]ProviderRevenue, error)
}

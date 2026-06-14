package admin

import (
	"context"
	"time"

	"github.com/shopspring/decimal"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
)

// PlatformOverview is the admin dashboard summary for a date range.
type PlatformOverview struct {
	// User counts
	TotalStudents int
	TotalTeachers int
	TotalAdmins   int
	// Financial
	GrossRevenue   decimal.Decimal
	RefundedAmount decimal.Decimal
	CompletedPayments int64
	MRR            decimal.Decimal // active subs × plan price (approximation)
	// Content
	TotalCourses int64
	// Enrollments (approximation via payment count)
	TotalEnrollments int64
	// Date range covered
	From time.Time
	To   time.Time
}

// Overview returns platform-wide aggregates for the given date range. Admin only.
func (s *Service) Overview(ctx context.Context, actor Actor, from, to time.Time) (*PlatformOverview, error) {
	if err := requireAdmin(actor); err != nil {
		return nil, err
	}

	ov := &PlatformOverview{From: from, To: to}

	// User counts by role.
	if n, err := s.users.CountByRole(ctx, entity.RoleStudent); err == nil {
		ov.TotalStudents = n
	}
	if n, err := s.users.CountByRole(ctx, entity.RoleTeacher); err == nil {
		ov.TotalTeachers = n
	}
	if n, err := s.users.CountByRole(ctx, entity.RoleAdmin); err == nil {
		ov.TotalAdmins = n
	}

	// Revenue totals.
	if rev, err := s.revenue.CompletedTotals(ctx, from, to); err == nil {
		ov.GrossRevenue = rev.Gross
		ov.CompletedPayments = rev.Units
	}
	if ref, err := s.revenue.RefundedTotals(ctx, from, to); err == nil {
		ov.RefundedAmount = ref.Gross
	}

	return ov, nil
}

// RevenueBreakdown returns per-provider revenue for a date range. Admin only.
func (s *Service) RevenueBreakdown(ctx context.Context, actor Actor, from, to time.Time) ([]repository.ProviderRevenue, error) {
	if err := requireAdmin(actor); err != nil {
		return nil, err
	}
	return s.revenue.CompletedByProvider(ctx, from, to)
}

// EngagementStats is a lightweight read-only summary. Admin and teachers may
// view (teacher-scoped to their courses in Skill 5 HTTP layer).
type EngagementStats struct {
	// Placeholders — populated by Phase 2 analytics service calls in Skill 5.
	AvgQuizPassRate  float64
	AvgExamPassRate  float64
	AvgCompletionPct float64
}

// GetEngagementStats returns engagement KPIs. Admin/teacher only.
func (s *Service) GetEngagementStats(ctx context.Context, actor Actor) (*EngagementStats, error) {
	if err := requireTeacherOrAdmin(actor); err != nil {
		return nil, err
	}
	// Full implementation wired in Skill 5 via analytics service delegation.
	return &EngagementStats{}, nil
}

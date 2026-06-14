// Package revenue holds the read-only revenue-reporting usecases (spec §1.6, §10):
// the platform revenue report, per-teacher earnings under the configurable 70/30
// split, and per-course sales stats. Like the other usecase packages it is pure
// orchestration over domain repositories — no HTTP, no SQL — and it applies the
// role/ownership scoping the repository layer deliberately omits.
//
// The teacher/platform split is config-driven (TEACHER_REVENUE_SHARE) so the
// ratio is never hardcoded; earnings are rounded to cents and the platform fee is
// taken as the remainder so gross == earnings + fee exactly.
package revenue

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// defaultTeacherShare is the spec §1.6 teacher cut (70%) used when config leaves
// the ratio unset.
var defaultTeacherShare = decimal.NewFromFloat(0.70)

// Actor is the authenticated caller's identity. A zero Actor (UserID == Nil) is
// anonymous and may read no report.
type Actor struct {
	UserID uuid.UUID
	Role   entity.Role
}

// Config carries the configurable teacher revenue share (spec §1.6). main.go maps
// the typed config onto this so the usecase stays free of the config package.
type Config struct {
	// TeacherShare is the fraction of gross a teacher keeps (e.g. 0.70). The
	// platform fee is the remainder (1 - share).
	TeacherShare decimal.Decimal
}

// Service implements the revenue-reporting usecases over the revenue and course
// repositories. now is unused today but injected for parity with the other
// usecases (and for a future "default window = last 30 days").
type Service struct {
	revenue repository.RevenueRepository
	courses repository.CourseRepository
	share   decimal.Decimal
}

// NewService wires the revenue usecase. A zero/out-of-range share falls back to
// the spec default (0.70) so a misconfigured ratio never silently zeroes payouts.
func NewService(rev repository.RevenueRepository, courses repository.CourseRepository, cfg Config) *Service {
	share := cfg.TeacherShare
	if share.LessThanOrEqual(decimal.Zero) || share.GreaterThan(decimal.NewFromInt(1)) {
		share = defaultTeacherShare
	}
	return &Service{revenue: rev, courses: courses, share: share}
}

// ProviderBreakdown is one provider's slice of completed gross.
type ProviderBreakdown struct {
	Provider string `json:"provider"`
	Gross    string `json:"gross"`
	Units    int64  `json:"units"`
}

// PlatformReport is the admin platform-revenue view over a window. Net is gross
// minus refunds.
type PlatformReport struct {
	From        time.Time           `json:"from"`
	To          time.Time           `json:"to"`
	Gross       string              `json:"gross"`
	Refunds     string              `json:"refunds"`
	Net         string              `json:"net"`
	Units       int64               `json:"units"`
	RefundUnits int64               `json:"refund_units"`
	ByProvider  []ProviderBreakdown `json:"by_provider"`
}

// TeacherEarningsReport is a teacher's payout view over a window: gross sales of
// their courses, the share applied, and the resulting earnings + platform fee.
type TeacherEarningsReport struct {
	TeacherID    string    `json:"teacher_id"`
	From         time.Time `json:"from"`
	To           time.Time `json:"to"`
	Gross        string    `json:"gross"`
	TeacherShare string    `json:"teacher_share"`
	Earnings     string    `json:"earnings"`
	PlatformFee  string    `json:"platform_fee"`
	Units        int64     `json:"units"`
}

// CourseSalesReport is the lifetime sales view for one course.
type CourseSalesReport struct {
	CourseID string `json:"course_id"`
	Gross    string `json:"gross"`
	Units    int64  `json:"units"`
}

// PlatformRevenue returns the platform-wide revenue report for [from, to].
// Admin-only.
func (s *Service) PlatformRevenue(ctx context.Context, actor Actor, from, to time.Time) (*PlatformReport, error) {
	if actor.UserID == uuid.Nil {
		return nil, apperror.Unauthorized("authentication required", nil)
	}
	if actor.Role != entity.RoleAdmin {
		return nil, apperror.Forbidden("only an administrator may view platform revenue", nil)
	}
	if err := validWindow(from, to); err != nil {
		return nil, err
	}

	completed, err := s.revenue.CompletedTotals(ctx, from, to)
	if err != nil {
		return nil, err
	}
	refunded, err := s.revenue.RefundedTotals(ctx, from, to)
	if err != nil {
		return nil, err
	}
	byProvider, err := s.revenue.CompletedByProvider(ctx, from, to)
	if err != nil {
		return nil, err
	}

	breakdown := make([]ProviderBreakdown, len(byProvider))
	for i, p := range byProvider {
		breakdown[i] = ProviderBreakdown{
			Provider: string(p.Provider),
			Gross:    p.Gross.StringFixed(2),
			Units:    p.Units,
		}
	}

	return &PlatformReport{
		From:        from,
		To:          to,
		Gross:       completed.Gross.StringFixed(2),
		Refunds:     refunded.Gross.StringFixed(2),
		Net:         completed.Gross.Sub(refunded.Gross).StringFixed(2),
		Units:       completed.Units,
		RefundUnits: refunded.Units,
		ByProvider:  breakdown,
	}, nil
}

// TeacherEarnings returns the earnings report for teacherID over [from, to],
// applying the configured teacher share. A teacher may read only their own
// earnings; an admin may read any.
func (s *Service) TeacherEarnings(ctx context.Context, actor Actor, teacherID uuid.UUID, from, to time.Time) (*TeacherEarningsReport, error) {
	if actor.UserID == uuid.Nil {
		return nil, apperror.Unauthorized("authentication required", nil)
	}
	if actor.Role != entity.RoleAdmin && actor.UserID != teacherID {
		return nil, apperror.Forbidden("you may only view your own earnings", nil)
	}
	if err := validWindow(from, to); err != nil {
		return nil, err
	}

	gross, err := s.revenue.TeacherGross(ctx, teacherID, from, to)
	if err != nil {
		return nil, err
	}
	// Round earnings to cents and take the platform fee as the remainder so the
	// two always sum back to gross exactly.
	earnings := gross.Gross.Mul(s.share).Round(2)
	platformFee := gross.Gross.Sub(earnings)

	return &TeacherEarningsReport{
		TeacherID:    teacherID.String(),
		From:         from,
		To:           to,
		Gross:        gross.Gross.StringFixed(2),
		TeacherShare: s.share.String(),
		Earnings:     earnings.StringFixed(2),
		PlatformFee:  platformFee.StringFixed(2),
		Units:        gross.Units,
	}, nil
}

// CourseSalesStats returns lifetime sales for courseID. The course owner or an
// admin may read it.
func (s *Service) CourseSalesStats(ctx context.Context, actor Actor, courseID uuid.UUID) (*CourseSalesReport, error) {
	if actor.UserID == uuid.Nil {
		return nil, apperror.Unauthorized("authentication required", nil)
	}
	course, err := s.courses.FindByID(ctx, courseID)
	if err != nil {
		return nil, err
	}
	if actor.Role != entity.RoleAdmin && course.TeacherID != actor.UserID {
		return nil, apperror.Forbidden("you do not own this course", nil)
	}

	sales, err := s.revenue.CourseSales(ctx, courseID)
	if err != nil {
		return nil, err
	}
	return &CourseSalesReport{
		CourseID: courseID.String(),
		Gross:    sales.Gross.StringFixed(2),
		Units:    sales.Units,
	}, nil
}

// validWindow rejects an inverted reporting window.
func validWindow(from, to time.Time) error {
	if from.After(to) {
		return apperror.New(400, "INVALID_WINDOW", "'from' must not be after 'to'", nil)
	}
	return nil
}

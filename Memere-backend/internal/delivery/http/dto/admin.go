package dto

import (
	"time"

	"github.com/shopspring/decimal"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/admin"
)

// SuspendUserRequest is the body of POST /admin/users/:id/suspend.
type SuspendUserRequest struct {
	Reason string `json:"reason"`
}

// ChangeRoleRequest is the body of POST /admin/users/:id/role.
type ChangeRoleRequest struct {
	Role string `json:"role"`
}

// UnpublishCourseRequest is the body of POST /admin/courses/:id/unpublish.
type UnpublishCourseRequest struct {
	Reason string `json:"reason"`
}

// BroadcastRequest is the body of POST /admin/announcements.
type BroadcastRequest struct {
	Title   string            `json:"title"`
	Body    string            `json:"body"`
	Segment string            `json:"segment"` // all | students | teachers | subscribers
	Data    map[string]string `json:"data,omitempty"`
}

// ReconcileResponse is the response for POST /admin/payments/reconcile.
type ReconcileResponse struct {
	Reconciled int `json:"reconciled"`
}

// OverviewResponse is returned by GET /admin/analytics/overview.
type OverviewResponse struct {
	TotalStudents     int             `json:"total_students"`
	TotalTeachers     int             `json:"total_teachers"`
	TotalAdmins       int             `json:"total_admins"`
	GrossRevenue      decimal.Decimal `json:"gross_revenue"`
	RefundedAmount    decimal.Decimal `json:"refunded_amount"`
	CompletedPayments int64           `json:"completed_payments"`
	MRR               decimal.Decimal `json:"mrr"`
	From              time.Time       `json:"from"`
	To                time.Time       `json:"to"`
}

// RevenueBreakdownItem is one provider's revenue row.
type RevenueBreakdownItem struct {
	Provider string          `json:"provider"`
	Gross    decimal.Decimal `json:"gross"`
	Units    int64           `json:"units"`
}

// EngagementStatsResponse is returned by GET /admin/analytics/engagement.
type EngagementStatsResponse struct {
	AvgQuizPassRate  float64 `json:"avg_quiz_pass_rate"`
	AvgExamPassRate  float64 `json:"avg_exam_pass_rate"`
	AvgCompletionPct float64 `json:"avg_completion_pct"`
}

// AdminPaymentResponse is the wire projection of entity.Payment for admin views.
type AdminPaymentResponse struct {
	ID             string    `json:"id"`
	StudentID      string    `json:"student_id"`
	Amount         string    `json:"amount"`
	Currency       string    `json:"currency"`
	Status         string    `json:"status"`
	Provider       string    `json:"provider"`
	ProviderTxnID  *string   `json:"provider_txn_id,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

// NewOverviewResponse maps PlatformOverview to the wire DTO.
func NewOverviewResponse(ov *admin.PlatformOverview) OverviewResponse {
	return OverviewResponse{
		TotalStudents:     ov.TotalStudents,
		TotalTeachers:     ov.TotalTeachers,
		TotalAdmins:       ov.TotalAdmins,
		GrossRevenue:      ov.GrossRevenue,
		RefundedAmount:    ov.RefundedAmount,
		CompletedPayments: ov.CompletedPayments,
		MRR:               ov.MRR,
		From:              ov.From,
		To:                ov.To,
	}
}

// NewRevenueBreakdown maps []ProviderRevenue to wire DTOs.
func NewRevenueBreakdown(rows []repository.ProviderRevenue) []RevenueBreakdownItem {
	out := make([]RevenueBreakdownItem, 0, len(rows))
	for _, r := range rows {
		out = append(out, RevenueBreakdownItem{
			Provider: string(r.Provider),
			Gross:    r.Gross,
			Units:    r.Units,
		})
	}
	return out
}

// NewEngagementStatsResponse maps EngagementStats to the wire DTO.
func NewEngagementStatsResponse(s *admin.EngagementStats) EngagementStatsResponse {
	return EngagementStatsResponse{
		AvgQuizPassRate:  s.AvgQuizPassRate,
		AvgExamPassRate:  s.AvgExamPassRate,
		AvgCompletionPct: s.AvgCompletionPct,
	}
}

// NewAdminPaymentResponse maps a payment to the admin wire DTO.
func NewAdminPaymentResponse(p *entity.Payment) AdminPaymentResponse {
	return AdminPaymentResponse{
		ID:            p.ID.String(),
		StudentID:     p.StudentID.String(),
		Amount:        p.Amount.String(),
		Currency:      p.Currency,
		Status:        string(p.Status),
		Provider:      string(p.Provider),
		ProviderTxnID: p.ProviderTxnID,
		CreatedAt:     p.CreatedAt,
	}
}

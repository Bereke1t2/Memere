package dto

import (
	"time"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/subscription"
)

// PlanResponse is one entry of the public subscription-plan catalogue: the price
// to charge and the billing period in days.
type PlanResponse struct {
	Plan       string `json:"plan"`
	Price      string `json:"price"`
	Currency   string `json:"currency"`
	PeriodDays int    `json:"period_days"`
}

// NewPlanList maps the usecase plan specs to the wire catalogue.
func NewPlanList(specs []subscription.PlanSpec) []PlanResponse {
	out := make([]PlanResponse, 0, len(specs))
	for _, s := range specs {
		out = append(out, PlanResponse{
			Plan:       string(s.Plan),
			Price:      s.Price.StringFixed(2),
			Currency:   s.Currency,
			PeriodDays: int(s.Period.Hours() / 24),
		})
	}
	return out
}

// SubscriptionResponse is the client-safe view of a student's subscription. The
// provider subscription id and other provider internals are never exposed.
type SubscriptionResponse struct {
	ID                 string  `json:"id"`
	Plan               string  `json:"plan"`
	Status             string  `json:"status"`
	CurrentPeriodStart string  `json:"current_period_start"`
	CurrentPeriodEnd   string  `json:"current_period_end"`
	CanceledAt         *string `json:"canceled_at,omitempty"`
}

// NewSubscriptionResponse maps an entity subscription to the wire shape.
func NewSubscriptionResponse(s *entity.Subscription) SubscriptionResponse {
	resp := SubscriptionResponse{
		ID:                 s.ID.String(),
		Plan:               string(s.Plan),
		Status:             string(s.Status),
		CurrentPeriodStart: s.CurrentPeriodStart.UTC().Format(time.RFC3339),
		CurrentPeriodEnd:   s.CurrentPeriodEnd.UTC().Format(time.RFC3339),
	}
	if s.CanceledAt != nil {
		c := s.CanceledAt.UTC().Format(time.RFC3339)
		resp.CanceledAt = &c
	}
	return resp
}

package entity

import (
	"time"

	"github.com/google/uuid"
)

// SubscriptionPlan names a billing cadence (spec §10.2 / migration CHECK).
type SubscriptionPlan string

const (
	PlanMonthly SubscriptionPlan = "monthly"
	PlanAnnual  SubscriptionPlan = "annual"
)

// Valid reports whether p is a known plan.
func (p SubscriptionPlan) Valid() bool {
	switch p {
	case PlanMonthly, PlanAnnual:
		return true
	}
	return false
}

// SubscriptionStatus is the lifecycle state of a subscription.
type SubscriptionStatus string

const (
	SubActive   SubscriptionStatus = "active"
	SubPastDue  SubscriptionStatus = "past_due"
	SubCanceled SubscriptionStatus = "canceled"
	SubExpired  SubscriptionStatus = "expired"
)

// Valid reports whether s is a known subscription status.
func (s SubscriptionStatus) Valid() bool {
	switch s {
	case SubActive, SubPastDue, SubCanceled, SubExpired:
		return true
	}
	return false
}

// Subscription is a student's all-access plan (spec §10.2). An active, unexpired
// subscription grants access to every course via the access layer (Skill 2).
type Subscription struct {
	ID                     uuid.UUID
	StudentID              uuid.UUID
	Plan                   SubscriptionPlan
	Status                 SubscriptionStatus
	CurrentPeriodStart     time.Time
	CurrentPeriodEnd       time.Time
	Provider               PaymentProvider
	ProviderSubscriptionID *string
	CreatedAt              time.Time
	UpdatedAt              time.Time
	CanceledAt             *time.Time
}

// IsActiveAt reports whether the subscription grants access at instant t: status
// is active and the current period has not ended.
func (s *Subscription) IsActiveAt(t time.Time) bool {
	return s.Status == SubActive && s.CurrentPeriodEnd.After(t)
}

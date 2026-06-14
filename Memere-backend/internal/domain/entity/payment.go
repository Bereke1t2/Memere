// Package entity holds the pure-Go domain models. These structs carry no db or
// json tags (clean-architecture dependency rule): the repository layer maps them
// to/from sqlc models, and the delivery layer maps them to/from DTOs.
//
// Money in this file uses shopspring/decimal (exact decimal arithmetic) — never
// a binary float — so amounts never accumulate rounding drift (spec §7.3 /
// Non-Negotiable: no float in money operations).
package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// PaymentProvider names a payment backend. Values mirror the payments.provider
// CHECK constraint (spec §10.1).
type PaymentProvider string

const (
	ProviderChapa    PaymentProvider = "chapa"
	ProviderTelebirr PaymentProvider = "telebirr"
	ProviderStripe   PaymentProvider = "stripe"
	// ProviderMock is a test-only backend used by the Phase 4 smoke test: it
	// settles offline (no external HTTP) and signs its own webhook. It is only
	// usable when explicitly registered (PAYMENT_MOCK_ENABLED), never in prod.
	ProviderMock PaymentProvider = "mock"
)

// Valid reports whether p is a known provider.
func (p PaymentProvider) Valid() bool {
	switch p {
	case ProviderChapa, ProviderTelebirr, ProviderStripe, ProviderMock:
		return true
	}
	return false
}

// PaymentStatus is the lifecycle state of a payment. Values mirror the
// payments.status CHECK constraint.
type PaymentStatus string

const (
	PayPending   PaymentStatus = "pending"
	PayCompleted PaymentStatus = "completed"
	PayFailed    PaymentStatus = "failed"
	PayRefunded  PaymentStatus = "refunded"
)

// Valid reports whether s is a known status.
func (s PaymentStatus) Valid() bool {
	switch s {
	case PayPending, PayCompleted, PayFailed, PayRefunded:
		return true
	}
	return false
}

// CanTransitionTo encodes the payment state machine. A pending payment settles
// to completed or fails; a completed payment can only be refunded; failed and
// refunded are terminal. The guarded UPDATE in the repo enforces this same rule
// at the database so a re-delivered webhook cannot double-settle a row.
func (s PaymentStatus) CanTransitionTo(next PaymentStatus) bool {
	switch s {
	case PayPending:
		return next == PayCompleted || next == PayFailed
	case PayCompleted:
		return next == PayRefunded
	default:
		return false // failed / refunded are terminal
	}
}

// Payment is a single purchase attempt. A nil CourseID with a non-nil
// SubscriptionID means the payment settles a subscription rather than a course
// (spec §4.2.7: payments.course_id nullable, null = subscription). IdempotencyKey
// is the client-supplied key that makes initiate idempotent (no double-charge).
type Payment struct {
	ID             uuid.UUID
	StudentID      uuid.UUID
	CourseID       *uuid.UUID // nil => subscription payment
	SubscriptionID *uuid.UUID

	Amount   decimal.Decimal
	Currency string

	Provider           PaymentProvider
	ProviderTxnID      *string
	ProviderCheckoutID *string

	Status         PaymentStatus
	IdempotencyKey *string
	CouponID       *uuid.UUID
	FailureReason  *string
	Metadata       map[string]any

	PaidAt    *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

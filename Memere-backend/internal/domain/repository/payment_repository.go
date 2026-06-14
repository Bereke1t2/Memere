package repository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/pagination"
)

// UpdatePaymentStatusFields carries the optional fields a guarded status
// transition may set alongside the new status (provider txn id on completion,
// failure reason on failure). Nil fields are left unchanged.
type UpdatePaymentStatusFields struct {
	ProviderTxnID *string
	FailureReason *string
}

// PaymentRepository persists payments. Reads that serve a student filter by the
// authenticated student_id to prevent IDOR (Non-Negotiable #7). The
// idempotency_key UNIQUE index and the guarded status update are the hard
// guarantees against double-charge / double-settle.
type PaymentRepository interface {
	Create(ctx context.Context, p *entity.Payment) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Payment, error)
	// GetByIdempotencyKey returns the payment previously created for key, or
	// apperror.NotFound — the friendly path that lets a retried initiate return
	// the same checkout instead of charging twice.
	GetByIdempotencyKey(ctx context.Context, key string) (*entity.Payment, error)
	// GetByProviderTxn resolves a payment from a provider transaction id (webhook
	// path) scoped to its provider.
	GetByProviderTxn(ctx context.Context, provider entity.PaymentProvider, txnID string) (*entity.Payment, error)
	// UpdateStatusGuarded flips status only when the row still holds `from`,
	// returning false (not an error) on a 0-row result so a re-delivered webhook
	// cannot settle an already-settled payment.
	UpdateStatusGuarded(ctx context.Context, id uuid.UUID, from, to entity.PaymentStatus, fields UpdatePaymentStatusFields) (bool, error)
	ListByStudent(ctx context.Context, studentID uuid.UUID, limit int) ([]*entity.Payment, error)
	// ListAll returns payments matching filter for admin reconciliation views.
	ListAll(ctx context.Context, filter AdminPaymentFilter, cursor *pagination.Cursor, limit int) ([]*entity.Payment, *pagination.Cursor, error)
}

// EnrollmentRepository persists course enrollments. The (student_id, course_id)
// UNIQUE constraint makes Create idempotent at the database, so a re-processed
// payment grants access exactly once.
type EnrollmentRepository interface {
	Create(ctx context.Context, e *entity.Enrollment) error
	// Exists reports whether the student already has any enrollment row for the
	// course (regardless of expiry).
	Exists(ctx context.Context, studentID, courseID uuid.UUID) (bool, error)
	// GetActiveForStudent returns the student's enrollment for the course when it
	// grants access now (unexpired), or apperror.NotFound.
	GetActiveForStudent(ctx context.Context, studentID, courseID uuid.UUID) (*entity.Enrollment, error)
	ListByStudent(ctx context.Context, studentID uuid.UUID, limit int) ([]*entity.Enrollment, error)
}

// CouponRepository persists coupons. IncrementUse is guarded on max_uses so a
// coupon can never be redeemed past its limit even under concurrent fulfillment.
type CouponRepository interface {
	GetByCode(ctx context.Context, code string) (*entity.Coupon, error)
	// IncrementUse atomically bumps used_count when it is still below max_uses (or
	// max_uses is NULL), returning false (not an error) when the coupon is already
	// exhausted. Called only inside the successful-payment fulfillment transaction.
	IncrementUse(ctx context.Context, id uuid.UUID) (bool, error)
}

// SubscriptionRepository persists all-access subscriptions (spec §10.2).
type SubscriptionRepository interface {
	Create(ctx context.Context, s *entity.Subscription) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Subscription, error)
	// GetActiveForStudent returns the student's active, unexpired subscription or
	// apperror.NotFound — the all-access access check.
	GetActiveForStudent(ctx context.Context, studentID uuid.UUID) (*entity.Subscription, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status entity.SubscriptionStatus) error
	// ExtendPeriod pushes an active subscription's period end out and keeps it
	// active — the period-extension primitive behind idempotent Activate.
	ExtendPeriod(ctx context.Context, id uuid.UUID, newEnd time.Time) error
	// CancelAtPeriodEnd flags the subscription as canceled (won't renew) while
	// keeping access until current_period_end; returns false when there was no
	// active, not-yet-canceled row to flag (idempotent).
	CancelAtPeriodEnd(ctx context.Context, id uuid.UUID) (bool, error)
	// ListExpiring returns active subscriptions whose period has ended, for the
	// renewal/expiry sweeper (Skill 4).
	ListExpiring(ctx context.Context, limit int) ([]*entity.Subscription, error)
	// ExpireLapsed transitions an active, period-ended subscription to expired,
	// guarded on status + period so a renewal racing the sweeper is never
	// clobbered; returns false when there was no still-lapsed active row (the
	// transition the sweeper applies after ListExpiring).
	ExpireLapsed(ctx context.Context, id uuid.UUID) (bool, error)
}

// WebhookEventRepository is the dedup ledger (spec §7.3).
type WebhookEventRepository interface {
	// InsertIfNew records the event and returns true when it was newly inserted,
	// false when (provider, eventID) already existed — the duplicate signal that
	// lets the caller skip re-processing.
	InsertIfNew(ctx context.Context, provider entity.PaymentProvider, eventID, eventType string, payload []byte) (bool, error)
	// MarkProcessed stamps processed_at after the event's fulfillment commits.
	MarkProcessed(ctx context.Context, provider entity.PaymentProvider, eventID string) error
}

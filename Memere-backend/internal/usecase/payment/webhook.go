package payment

import (
	"context"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/service"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// HandleWebhook processes a provider settlement callback (spec §10.2 + §7.3). It
// verifies the signature, deduplicates the event, resolves our payment, and on
// success fulfils it in one guarded transaction: complete the payment, burn the
// coupon, and grant enrollment — exactly once, even under re-delivery.
func (s *Service) HandleWebhook(ctx context.Context, provider entity.PaymentProvider, headers map[string]string, body []byte) error {
	prov, err := s.registry.Get(provider)
	if err != nil {
		return err
	}
	// Signature verification: a tampered or unsigned body is rejected (§7.3). The
	// raw body — never a parsed/trusted struct — is what the provider checks.
	ev, err := prov.VerifyWebhook(ctx, headers, body)
	if err != nil {
		return err
	}

	// Dedup ledger: the first delivery inserts; a re-delivery returns isNew=false
	// and we ack without reprocessing, so a provider retry never double-grants.
	isNew, err := s.webhooks.InsertIfNew(ctx, provider, ev.ProviderEventID, ev.Type, ev.RawPayload)
	if err != nil {
		return err
	}
	if !isNew {
		return nil
	}

	p, err := s.resolvePayment(ctx, provider, ev)
	if err != nil {
		return err
	}

	switch ev.Status {
	case string(entity.PayCompleted):
		if err := s.fulfill(ctx, p, ev.ProviderTxnID); err != nil {
			return err
		}
	case string(entity.PayFailed):
		reason := ev.Type
		if _, err := s.payments.UpdateStatusGuarded(ctx, p.ID, entity.PayPending, entity.PayFailed,
			repository.UpdatePaymentStatusFields{FailureReason: &reason}); err != nil {
			return err
		}
	default:
		// Unknown status: nothing to settle, but the event is recorded. Ack.
	}

	if err := s.webhooks.MarkProcessed(ctx, provider, ev.ProviderEventID); err != nil {
		return err
	}
	if ev.Status == string(entity.PayCompleted) {
		// Reload so the post-commit notification sees the settled row, then fire
		// the (no-op until Phase 5) enrollment-confirmed hook.
		if fresh, gerr := s.payments.GetByID(ctx, p.ID); gerr == nil {
			p = fresh
		}
		s.notify.PurchaseConfirmed(ctx, p)
	}
	return nil
}

// fulfill settles a payment in one transaction: the guarded pending→completed
// flip wins at most once, and only that winner burns the coupon and grants
// enrollment. A racing webhook re-delivery or polling fallback flips zero rows
// and no-ops, so enrollment and coupon use are each applied exactly once.
func (s *Service) fulfill(ctx context.Context, p *entity.Payment, providerTxnID string) error {
	var txn *string
	if providerTxnID != "" {
		txn = &providerTxnID
	}
	return s.tx.WithinTx(ctx, func(ctx context.Context) error {
		ok, err := s.payments.UpdateStatusGuarded(ctx, p.ID, entity.PayPending, entity.PayCompleted,
			repository.UpdatePaymentStatusFields{ProviderTxnID: txn})
		if err != nil {
			return err
		}
		if !ok {
			// Already settled by a prior delivery/poll — idempotent no-op.
			return nil
		}
		// Burn the coupon ONLY on success, inside the same tx (spec §10.3).
		if p.CouponID != nil {
			if _, err := s.coupons.IncrementUse(ctx, *p.CouponID); err != nil {
				return err
			}
		}
		if p.CourseID != nil {
			return s.grantCourse(ctx, p.StudentID, *p.CourseID)
		}
		// Subscription fulfillment arrives in Skill 4; the payment is completed.
		return nil
	})
}

// grantCourse idempotently enrols the student in the purchased course. The
// friendly Exists check avoids a UNIQUE violation aborting the enclosing
// transaction; the (student_id, course_id) UNIQUE constraint is the hard
// backstop so a slipped race still grants exactly once.
func (s *Service) grantCourse(ctx context.Context, studentID, courseID uuid.UUID) error {
	exists, err := s.enroll.Exists(ctx, studentID, courseID)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	now := s.now()
	return s.enroll.Create(ctx, &entity.Enrollment{
		ID:         uuid.New(),
		StudentID:  studentID,
		CourseID:   courseID,
		Source:     entity.SourcePurchase,
		EnrolledAt: now,
		CreatedAt:  now,
		UpdatedAt:  now,
	})
}

// resolvePayment maps a normalized webhook event back to our payment row: by the
// tx_ref we set to the PaymentID, falling back to the provider transaction id.
func (s *Service) resolvePayment(ctx context.Context, provider entity.PaymentProvider, ev *service.WebhookEvent) (*entity.Payment, error) {
	if id, err := uuid.Parse(ev.PaymentRef); err == nil {
		if p, perr := s.payments.GetByID(ctx, id); perr == nil {
			return p, nil
		}
	}
	if ev.ProviderTxnID != "" {
		if p, perr := s.payments.GetByProviderTxn(ctx, provider, ev.ProviderTxnID); perr == nil {
			return p, nil
		}
	}
	return nil, apperror.NotFound("payment not found for webhook", nil)
}

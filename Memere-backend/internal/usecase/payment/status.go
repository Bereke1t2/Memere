package payment

import (
	"context"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// GetPaymentStatus returns the current status of one of the caller's own
// payments (IDOR-safe: a student may only read their own; an admin may read any).
// If the payment is still pending and the provider has a reference to verify, it
// runs a server-to-server confirmation as a polling fallback for a lost webhook,
// then settles through the same guarded/transactional path — so polling can
// complete a payment without ever double-granting.
func (s *Service) GetPaymentStatus(ctx context.Context, actor Actor, paymentID uuid.UUID) (*PaymentView, error) {
	if actor.UserID == uuid.Nil {
		return nil, apperror.Unauthorized("authentication required", nil)
	}
	p, err := s.payments.GetByID(ctx, paymentID)
	if err != nil {
		return nil, err
	}
	if actor.Role != entity.RoleAdmin && p.StudentID != actor.UserID {
		return nil, apperror.Forbidden("you do not have access to this payment", nil)
	}

	if p.Status == entity.PayPending {
		if ref := providerRef(p); ref != "" {
			if prov, perr := s.registry.Get(p.Provider); perr == nil {
				if status, verr := prov.VerifyPayment(ctx, ref); verr == nil {
					switch status {
					case string(entity.PayCompleted):
						if err := s.fulfill(ctx, p, derefTxn(p)); err != nil {
							return nil, err
						}
					case string(entity.PayFailed):
						reason := "verify_failed"
						if _, err := s.payments.UpdateStatusGuarded(ctx, p.ID, entity.PayPending, entity.PayFailed,
							repository.UpdatePaymentStatusFields{FailureReason: &reason}); err != nil {
							return nil, err
						}
					}
					// Reload to reflect any transition.
					if fresh, gerr := s.payments.GetByID(ctx, p.ID); gerr == nil {
						p = fresh
					}
				}
			}
		}
	}
	return viewOf(p), nil
}

// RefundPayment marks a completed payment refunded (admin only). Per spec
// decision the enrollment is intentionally KEPT — refunds are handled
// out-of-band and revoking access automatically would be surprising; the
// refunded status on the payment is the audit trail. A provider-side refund call
// is stubbed until real keys land.
func (s *Service) RefundPayment(ctx context.Context, actor Actor, paymentID uuid.UUID) (*PaymentView, error) {
	if actor.Role != entity.RoleAdmin {
		return nil, apperror.Forbidden("only an administrator may refund a payment", nil)
	}
	p, err := s.payments.GetByID(ctx, paymentID)
	if err != nil {
		return nil, err
	}
	if p.Status != entity.PayCompleted {
		return nil, apperror.New(409, "NOT_REFUNDABLE",
			"only a completed payment can be refunded", nil)
	}
	ok, err := s.payments.UpdateStatusGuarded(ctx, p.ID, entity.PayCompleted, entity.PayRefunded,
		repository.UpdatePaymentStatusFields{})
	if err != nil {
		return nil, err
	}
	if !ok {
		// Lost the race to a concurrent refund; report the current state.
		if fresh, gerr := s.payments.GetByID(ctx, p.ID); gerr == nil {
			return viewOf(fresh), nil
		}
	}
	if fresh, gerr := s.payments.GetByID(ctx, p.ID); gerr == nil {
		p = fresh
	}
	return viewOf(p), nil
}

// providerRef returns the best provider reference to verify a pending payment
// against: the provider transaction id if known, else the checkout id.
func providerRef(p *entity.Payment) string {
	if p.ProviderTxnID != nil && *p.ProviderTxnID != "" {
		return *p.ProviderTxnID
	}
	if p.ProviderCheckoutID != nil {
		return *p.ProviderCheckoutID
	}
	return ""
}

// derefTxn returns the payment's provider transaction id or "" when unset.
func derefTxn(p *entity.Payment) string {
	if p.ProviderTxnID != nil {
		return *p.ProviderTxnID
	}
	return ""
}

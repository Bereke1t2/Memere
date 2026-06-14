package admin

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/pagination"
)

// ListPayments returns payments matching filter. Admin only.
func (s *Service) ListPayments(ctx context.Context, actor Actor, filter repository.AdminPaymentFilter, cursor *pagination.Cursor, limit int) ([]*entity.Payment, *pagination.Cursor, error) {
	if err := requireAdmin(actor); err != nil {
		return nil, nil, err
	}
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	return s.payments.ListAll(ctx, filter, cursor, limit)
}

// GetPaymentDetail returns a single payment. Admin only.
func (s *Service) GetPaymentDetail(ctx context.Context, actor Actor, paymentID uuid.UUID) (*entity.Payment, error) {
	if err := requireAdmin(actor); err != nil {
		return nil, err
	}
	return s.payments.GetByID(ctx, paymentID)
}

// ReconcilePending polls the provider for every payment that is still `pending`
// after staleDuration (e.g. 10 minutes) and fulfills it if the provider reports
// success. Uses Phase 4's guarded UpdateStatusGuarded so a race cannot
// double-grant. Admin only. Audited.
func (s *Service) ReconcilePending(ctx context.Context, actor Actor, staleDuration time.Duration) (reconciled int, err error) {
	if err := requireAdmin(actor); err != nil {
		return 0, err
	}
	from := time.Now().UTC().Add(-staleDuration)
	to := time.Now().UTC()
	pending := entity.PayPending
	staleStr := staleDuration.String()
	filter := repository.AdminPaymentFilter{
		Status: &staleStr,
		From:   &from,
		To:     &to,
	}
	// Override status filter to pending.
	pendingStr := string(pending)
	filter.Status = &pendingStr

	payments, _, err := s.payments.ListAll(ctx, filter, nil, 100)
	if err != nil {
		return 0, err
	}

	for _, p := range payments {
		if time.Since(p.CreatedAt) < staleDuration {
			continue
		}
		if s.providers == nil {
			continue
		}
		prov, err := s.providers.Get(p.Provider)
		if err != nil {
			continue
		}
		status, err := prov.VerifyPayment(ctx, safeDeref(p.ProviderTxnID))
		if err != nil {
			continue
		}
		if status == "success" || status == "completed" {
			id := p.ID
			_, _ = s.payments.UpdateStatusGuarded(ctx, p.ID,
				entity.PayPending, entity.PayCompleted,
				repository.UpdatePaymentStatusFields{})
			s.writeAudit(ctx, actor, "payment.reconcile", "payment", &id, map[string]any{
				"provider": string(p.Provider),
				"status":   status,
			})
			reconciled++
		}
	}
	return reconciled, nil
}

// ManualRefund marks a completed payment as refunded. Admin only. Audited.
func (s *Service) ManualRefund(ctx context.Context, actor Actor, paymentID uuid.UUID, reason string) error {
	if err := requireAdmin(actor); err != nil {
		return err
	}
	p, err := s.payments.GetByID(ctx, paymentID)
	if err != nil {
		return err
	}
	if p.Status != entity.PayCompleted {
		return apperror.BadRequest("only completed payments can be refunded", nil)
	}
	ok, err := s.payments.UpdateStatusGuarded(ctx, paymentID,
		entity.PayCompleted, entity.PayRefunded,
		repository.UpdatePaymentStatusFields{FailureReason: &reason})
	if err != nil {
		return err
	}
	if !ok {
		return apperror.Conflict("payment status changed concurrently", nil)
	}
	s.writeAudit(ctx, actor, "payment.refund", "payment", &paymentID, map[string]any{"reason": reason})
	return nil
}

func safeDeref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// ensure fmt is used
var _ = fmt.Sprintf

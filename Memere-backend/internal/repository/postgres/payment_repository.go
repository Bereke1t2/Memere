package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/internal/repository/postgres/sqlcgen"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// PaymentRepo is the sqlc-backed implementation of repository.PaymentRepository.
type PaymentRepo struct {
	q *sqlcgen.Queries
}

var _ repository.PaymentRepository = (*PaymentRepo)(nil)

// NewPaymentRepo builds a PaymentRepo over a pgx pool.
func NewPaymentRepo(pool *pgxpool.Pool) *PaymentRepo {
	return &PaymentRepo{q: sqlcgen.New(pool)}
}

func (r *PaymentRepo) Create(ctx context.Context, p *entity.Payment) error {
	id := p.ID
	if id == uuid.Nil {
		id = uuid.New()
	}
	status := p.Status
	if status == "" {
		status = entity.PayPending
	}
	row, err := queriesFor(ctx, r.q).CreatePayment(ctx, sqlcgen.CreatePaymentParams{
		ID:                 toPgUUID(id),
		StudentID:          toPgUUID(p.StudentID),
		CourseID:           toPgUUIDPtr(p.CourseID),
		SubscriptionID:     toPgUUIDPtr(p.SubscriptionID),
		Amount:             toPgDecimal(p.Amount),
		Currency:           p.Currency,
		Provider:           string(p.Provider),
		ProviderCheckoutID: p.ProviderCheckoutID,
		Status:             string(status),
		CouponID:           toPgUUIDPtr(p.CouponID),
		IdempotencyKey:     derefString(p.IdempotencyKey),
		Metadata:           toJSONB(p.Metadata),
	})
	if err != nil {
		// The idempotency_key UNIQUE index is the hard no-double-charge guarantee:
		// a retried initiate with the same key collides here.
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolation {
			return apperror.Conflict("PAYMENT_DUPLICATE", err)
		}
		return apperror.Internal(err)
	}
	*p = *paymentFromRow(row)
	return nil
}

func (r *PaymentRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.Payment, error) {
	row, err := queriesFor(ctx, r.q).GetPaymentByID(ctx, toPgUUID(id))
	if err != nil {
		return nil, mapPaymentErr(err)
	}
	return paymentFromRow(row), nil
}

func (r *PaymentRepo) GetByIdempotencyKey(ctx context.Context, key string) (*entity.Payment, error) {
	row, err := queriesFor(ctx, r.q).GetPaymentByIdempotencyKey(ctx, key)
	if err != nil {
		return nil, mapPaymentErr(err)
	}
	return paymentFromRow(row), nil
}

func (r *PaymentRepo) GetByProviderTxn(ctx context.Context, provider entity.PaymentProvider, txnID string) (*entity.Payment, error) {
	row, err := queriesFor(ctx, r.q).GetPaymentByProviderTxn(ctx, sqlcgen.GetPaymentByProviderTxnParams{
		Provider:              string(provider),
		ProviderTransactionID: &txnID,
	})
	if err != nil {
		return nil, mapPaymentErr(err)
	}
	return paymentFromRow(row), nil
}

// UpdateStatusGuarded runs the guarded transition. A 0-row result means the row
// no longer holds `from` (e.g. a re-delivered webhook): reported as changed=false,
// not an error, so the caller skips a duplicate fulfillment.
func (r *PaymentRepo) UpdateStatusGuarded(ctx context.Context, id uuid.UUID, from, to entity.PaymentStatus, fields repository.UpdatePaymentStatusFields) (bool, error) {
	n, err := queriesFor(ctx, r.q).UpdatePaymentStatusGuarded(ctx, sqlcgen.UpdatePaymentStatusGuardedParams{
		ToStatus:      string(to),
		ProviderTxn:   fields.ProviderTxnID,
		FailureReason: fields.FailureReason,
		SetPaid:       to == entity.PayCompleted,
		ID:            toPgUUID(id),
		FromStatus:    string(from),
	})
	if err != nil {
		return false, apperror.Internal(err)
	}
	return n > 0, nil
}

func (r *PaymentRepo) ListByStudent(ctx context.Context, studentID uuid.UUID, limit int) ([]*entity.Payment, error) {
	rows, err := queriesFor(ctx, r.q).ListPaymentsByStudent(ctx, sqlcgen.ListPaymentsByStudentParams{
		StudentID: toPgUUID(studentID),
		Limit:     int32(limit),
	})
	if err != nil {
		return nil, apperror.Internal(err)
	}
	out := make([]*entity.Payment, len(rows))
	for i, row := range rows {
		out[i] = paymentFromRow(row)
	}
	return out, nil
}

func mapPaymentErr(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return apperror.NotFound("payment not found", err)
	}
	return apperror.Internal(err)
}

func paymentFromRow(row sqlcgen.PaymentsPayment) *entity.Payment {
	idemKey := row.IdempotencyKey
	return &entity.Payment{
		ID:                 fromPgUUID(row.ID),
		StudentID:          fromPgUUID(row.StudentID),
		CourseID:           fromPgUUIDPtr(row.CourseID),
		SubscriptionID:     fromPgUUIDPtr(row.SubscriptionID),
		Amount:             fromPgDecimal(row.Amount),
		Currency:           row.Currency,
		Provider:           entity.PaymentProvider(row.Provider),
		ProviderTxnID:      row.ProviderTransactionID,
		ProviderCheckoutID: row.ProviderCheckoutID,
		Status:             entity.PaymentStatus(row.Status),
		IdempotencyKey:     &idemKey,
		CouponID:           fromPgUUIDPtr(row.CouponID),
		FailureReason:      row.FailureReason,
		Metadata:           fromJSONB(row.Metadata),
		PaidAt:             fromPgTimestamptz(row.PaidAt),
		CreatedAt:          fromPgTimestamptzValue(row.CreatedAt),
		UpdatedAt:          fromPgTimestamptzValue(row.UpdatedAt),
	}
}

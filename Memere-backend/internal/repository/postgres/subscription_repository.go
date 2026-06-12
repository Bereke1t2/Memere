package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/internal/repository/postgres/sqlcgen"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// SubscriptionRepo is the sqlc-backed implementation of
// repository.SubscriptionRepository.
type SubscriptionRepo struct {
	q *sqlcgen.Queries
}

var _ repository.SubscriptionRepository = (*SubscriptionRepo)(nil)

// NewSubscriptionRepo builds a SubscriptionRepo over a pgx pool.
func NewSubscriptionRepo(pool *pgxpool.Pool) *SubscriptionRepo {
	return &SubscriptionRepo{q: sqlcgen.New(pool)}
}

func (r *SubscriptionRepo) Create(ctx context.Context, s *entity.Subscription) error {
	id := s.ID
	if id == uuid.Nil {
		id = uuid.New()
	}
	row, err := queriesFor(ctx, r.q).CreateSubscription(ctx, sqlcgen.CreateSubscriptionParams{
		ID:                     toPgUUID(id),
		StudentID:              toPgUUID(s.StudentID),
		Plan:                   string(s.Plan),
		Status:                 string(s.Status),
		CurrentPeriodStart:     pgTimestamptzValue(s.CurrentPeriodStart),
		CurrentPeriodEnd:       pgTimestamptzValue(s.CurrentPeriodEnd),
		Provider:               string(s.Provider),
		ProviderSubscriptionID: s.ProviderSubscriptionID,
	})
	if err != nil {
		return apperror.Internal(err)
	}
	*s = *subscriptionFromRow(row)
	return nil
}

func (r *SubscriptionRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.Subscription, error) {
	row, err := queriesFor(ctx, r.q).GetSubscriptionByID(ctx, toPgUUID(id))
	if err != nil {
		return nil, mapSubscriptionErr(err)
	}
	return subscriptionFromRow(row), nil
}

func (r *SubscriptionRepo) GetActiveForStudent(ctx context.Context, studentID uuid.UUID) (*entity.Subscription, error) {
	row, err := queriesFor(ctx, r.q).GetActiveSubscriptionForStudent(ctx, toPgUUID(studentID))
	if err != nil {
		return nil, mapSubscriptionErr(err)
	}
	return subscriptionFromRow(row), nil
}

func (r *SubscriptionRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status entity.SubscriptionStatus) error {
	if err := queriesFor(ctx, r.q).UpdateSubscriptionStatus(ctx, sqlcgen.UpdateSubscriptionStatusParams{
		ID:     toPgUUID(id),
		Status: string(status),
	}); err != nil {
		return apperror.Internal(err)
	}
	return nil
}

func (r *SubscriptionRepo) ListExpiring(ctx context.Context, limit int) ([]*entity.Subscription, error) {
	rows, err := queriesFor(ctx, r.q).ListExpiringSubscriptions(ctx, int32(limit))
	if err != nil {
		return nil, apperror.Internal(err)
	}
	out := make([]*entity.Subscription, len(rows))
	for i, row := range rows {
		out[i] = subscriptionFromRow(row)
	}
	return out, nil
}

func mapSubscriptionErr(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return apperror.NotFound("subscription not found", err)
	}
	return apperror.Internal(err)
}

func subscriptionFromRow(row sqlcgen.PaymentsSubscription) *entity.Subscription {
	return &entity.Subscription{
		ID:                     fromPgUUID(row.ID),
		StudentID:              fromPgUUID(row.StudentID),
		Plan:                   entity.SubscriptionPlan(row.Plan),
		Status:                 entity.SubscriptionStatus(row.Status),
		CurrentPeriodStart:     fromPgTimestamptzValue(row.CurrentPeriodStart),
		CurrentPeriodEnd:       fromPgTimestamptzValue(row.CurrentPeriodEnd),
		Provider:               entity.PaymentProvider(row.Provider),
		ProviderSubscriptionID: row.ProviderSubscriptionID,
		CreatedAt:              fromPgTimestamptzValue(row.CreatedAt),
		UpdatedAt:              fromPgTimestamptzValue(row.UpdatedAt),
		CanceledAt:             fromPgTimestamptz(row.CanceledAt),
	}
}

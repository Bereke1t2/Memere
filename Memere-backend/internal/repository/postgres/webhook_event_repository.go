package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/internal/repository/postgres/sqlcgen"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// WebhookEventRepo is the sqlc-backed implementation of
// repository.WebhookEventRepository — the webhook dedup ledger (spec §7.3).
type WebhookEventRepo struct {
	q *sqlcgen.Queries
}

var _ repository.WebhookEventRepository = (*WebhookEventRepo)(nil)

// NewWebhookEventRepo builds a WebhookEventRepo over a pgx pool.
func NewWebhookEventRepo(pool *pgxpool.Pool) *WebhookEventRepo {
	return &WebhookEventRepo{q: sqlcgen.New(pool)}
}

// InsertIfNew records the event and returns true when it was newly inserted. The
// underlying query is ON CONFLICT (provider, provider_event_id) DO NOTHING
// RETURNING id, so a duplicate delivery returns no row (pgx.ErrNoRows) — reported
// as inserted=false, not an error, which is the caller's signal to skip
// re-processing.
func (r *WebhookEventRepo) InsertIfNew(ctx context.Context, provider entity.PaymentProvider, eventID, eventType string, payload []byte) (bool, error) {
	_, err := queriesFor(ctx, r.q).InsertWebhookIfNew(ctx, sqlcgen.InsertWebhookIfNewParams{
		Provider:        string(provider),
		ProviderEventID: eventID,
		EventType:       eventType,
		Payload:         payload,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil // duplicate delivery
		}
		return false, apperror.Internal(err)
	}
	return true, nil
}

func (r *WebhookEventRepo) MarkProcessed(ctx context.Context, provider entity.PaymentProvider, eventID string) error {
	if err := queriesFor(ctx, r.q).MarkWebhookProcessed(ctx, sqlcgen.MarkWebhookProcessedParams{
		Provider:        string(provider),
		ProviderEventID: eventID,
	}); err != nil {
		return apperror.Internal(err)
	}
	return nil
}

package postgres

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
)

// TestWebhookInsertIfNewDedup proves the dedup ledger (spec §7.3): the first
// delivery of a (provider, event_id) is inserted (true) and an identical
// re-delivery is a no-op (false), so a re-sent webhook is processed at most once.
// Skips when no database is configured.
func TestWebhookInsertIfNewDedup(t *testing.T) {
	dsn := dsnFromEnv()
	if dsn == "" {
		t.Skip("no database env (DB_HOST/TEST_DATABASE_URL); skipping integration check")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Skipf("cannot create pool: %v", err)
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil {
		t.Skipf("database unreachable: %v", err)
	}

	repo := NewWebhookEventRepo(pool)
	// A random event id keeps the test self-contained / re-runnable.
	eventID := "evt_" + uuid.NewString()
	payload := []byte(`{"event":"charge.success"}`)

	first, err := repo.InsertIfNew(ctx, entity.ProviderChapa, eventID, "charge.success", payload)
	if err != nil {
		t.Fatalf("first InsertIfNew error: %v", err)
	}
	if !first {
		t.Fatalf("first delivery must insert (true), got false")
	}

	dup, err := repo.InsertIfNew(ctx, entity.ProviderChapa, eventID, "charge.success", payload)
	if err != nil {
		t.Fatalf("duplicate InsertIfNew error: %v", err)
	}
	if dup {
		t.Fatalf("re-delivered event must be a no-op (false), got true")
	}

	if err := repo.MarkProcessed(ctx, entity.ProviderChapa, eventID); err != nil {
		t.Fatalf("MarkProcessed error: %v", err)
	}
}

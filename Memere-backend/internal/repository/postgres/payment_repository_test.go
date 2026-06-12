package postgres

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
)

// TestPaymentIdempotencyAndGuardedTransition exercises the two money guarantees
// at the database: (1) a created payment is retrievable by its idempotency key
// (the friendly no-double-charge path), and (2) the guarded status transition
// settles pending->completed exactly once — a second settle of the now-completed
// row reports changed=false, which is how a re-delivered webhook is made a no-op.
// Skips when no database is configured.
func TestPaymentIdempotencyAndGuardedTransition(t *testing.T) {
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

	studentID := mustCreateTestStudent(ctx, t, pool)
	repo := NewPaymentRepo(pool)

	key := "idem_" + uuid.NewString()
	p := &entity.Payment{
		StudentID:      studentID,
		Amount:         decimal.RequireFromString("199.99"),
		Currency:       "ETB",
		Provider:       entity.ProviderChapa,
		Status:         entity.PayPending,
		IdempotencyKey: &key,
	}
	if err := repo.Create(ctx, p); err != nil {
		t.Fatalf("Create payment: %v", err)
	}
	if p.ID == uuid.Nil {
		t.Fatalf("Create must populate the payment id")
	}

	// (1) idempotency-key lookup returns the same payment.
	got, err := repo.GetByIdempotencyKey(ctx, key)
	if err != nil {
		t.Fatalf("GetByIdempotencyKey: %v", err)
	}
	if got.ID != p.ID {
		t.Fatalf("idempotency lookup returned %s, want %s", got.ID, p.ID)
	}
	if !got.Amount.Equal(decimal.RequireFromString("199.99")) {
		t.Fatalf("amount round-trip lost precision: got %s", got.Amount)
	}

	// (2) guarded settle pending -> completed succeeds once...
	txn := "txn_" + uuid.NewString()
	ok, err := repo.UpdateStatusGuarded(ctx, p.ID, entity.PayPending, entity.PayCompleted,
		repository.UpdatePaymentStatusFields{ProviderTxnID: &txn})
	if err != nil {
		t.Fatalf("guarded settle: %v", err)
	}
	if !ok {
		t.Fatalf("first settle must change a pending row")
	}

	// ...and a re-delivered webhook (same pending->completed) is a no-op.
	ok, err = repo.UpdateStatusGuarded(ctx, p.ID, entity.PayPending, entity.PayCompleted,
		repository.UpdatePaymentStatusFields{ProviderTxnID: &txn})
	if err != nil {
		t.Fatalf("duplicate settle: %v", err)
	}
	if ok {
		t.Fatalf("duplicate settle of an already-completed payment must be false")
	}

	settled, err := repo.GetByID(ctx, p.ID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if settled.Status != entity.PayCompleted || settled.PaidAt == nil {
		t.Fatalf("settled payment must be completed with paid_at set, got %s paid_at=%v", settled.Status, settled.PaidAt)
	}
}

// mustCreateTestStudent inserts a throwaway active student and returns its id.
func mustCreateTestStudent(ctx context.Context, t *testing.T, pool *pgxpool.Pool) uuid.UUID {
	t.Helper()
	id := uuid.New()
	email := "pay-test-" + id.String() + "@example.com"
	_, err := pool.Exec(ctx,
		`INSERT INTO auth.users (id, email, password_hash, role, first_name, last_name)
		 VALUES ($1, $2, 'x', 'student', 'Pay', 'Test')`, id, email)
	if err != nil {
		t.Fatalf("seed student: %v", err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM payments.payments WHERE student_id = $1`, id)
		_, _ = pool.Exec(context.Background(), `DELETE FROM auth.users WHERE id = $1`, id)
	})
	return id
}

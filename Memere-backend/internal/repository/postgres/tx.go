package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/internal/repository/postgres/sqlcgen"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// txKey is the unexported context key under which an active pgx.Tx is carried so
// repositories can join an in-flight transaction. The value type is private, so
// no other package can set or read it.
type txKey struct{}

// TxManager is the pgx-backed implementation of repository.TxManager. It opens a
// transaction, stashes it in the context, and lets the repositories below pick
// it up via queriesFor.
type TxManager struct {
	pool *pgxpool.Pool
}

var _ repository.TxManager = (*TxManager)(nil)

// NewTxManager builds a TxManager over a pgx pool.
func NewTxManager(pool *pgxpool.Pool) *TxManager {
	return &TxManager{pool: pool}
}

// WithinTx runs fn inside a single transaction. fn must use the context it is
// handed for every repository call that should join the transaction. A non-nil
// return rolls back; nil commits. A nested call reuses the outer transaction.
func (m *TxManager) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
	if _, ok := txFromContext(ctx); ok {
		// Already inside a transaction — reuse it (savepoints are unnecessary
		// for the simple grouping the usecases need).
		return fn(ctx)
	}

	tx, err := m.pool.Begin(ctx)
	if err != nil {
		return apperror.Internal(err)
	}
	// Rollback is a no-op once the tx has committed, so the deferred safety net
	// only fires on an early return / panic.
	defer func() { _ = tx.Rollback(ctx) }()

	if err := fn(context.WithValue(ctx, txKey{}, tx)); err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return apperror.Internal(err)
	}
	return nil
}

// txFromContext returns the active transaction carried in ctx, if any.
func txFromContext(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(txKey{}).(pgx.Tx)
	return tx, ok
}

// queriesFor returns sqlc queries bound to the transaction in ctx when one is
// active, otherwise the pool-backed queries. Every repository routes its calls
// through this so a method automatically participates in an enclosing
// WithinTx without the caller threading anything but the context.
func queriesFor(ctx context.Context, base *sqlcgen.Queries) *sqlcgen.Queries {
	if tx, ok := txFromContext(ctx); ok {
		return base.WithTx(tx)
	}
	return base
}

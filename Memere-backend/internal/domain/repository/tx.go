package repository

import "context"

// TxManager runs a function within a single database transaction. The usecase
// layer uses it to group multiple repository writes (e.g. create a lesson and
// recompute the parent course's counters) so they commit or roll back together,
// without the usecase importing pgx or knowing how the transaction is carried.
//
// Implementations pass the active transaction to inner repository calls via the
// context: a repository checks the context for a transaction and uses it when
// present, falling back to the pool otherwise. The fn must use the ctx it is
// given for every repository call that should join the transaction.
type TxManager interface {
	// WithinTx invokes fn inside a transaction. If fn returns an error the
	// transaction is rolled back and that error is returned; otherwise it is
	// committed.
	WithinTx(ctx context.Context, fn func(ctx context.Context) error) error
}

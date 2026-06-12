// Package idempotency provides the friendly path for client-supplied
// Idempotency-Key handling on payment initiation. The hard guarantee against a
// double-charge is the UNIQUE index on payments.idempotency_key; this helper lets
// a retried request return the *same* payment instead of surfacing a unique
// violation, so a flaky network never charges a student twice (spec §7.3).
package idempotency

import (
	"context"
	"strings"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// MaxKeyLen bounds a key to the column width (payments.idempotency_key is
// VARCHAR(255)); anything longer is rejected before it reaches the database.
const MaxKeyLen = 255

// ErrMissingKey is returned when a payment-initiating request carries no
// Idempotency-Key.
var ErrMissingKey = apperror.BadRequest("an Idempotency-Key header is required", nil)

// Normalize trims and validates a client-supplied idempotency key. An empty or
// over-long key is rejected (a payment cannot be initiated without a usable key).
func Normalize(key string) (string, error) {
	k := strings.TrimSpace(key)
	if k == "" {
		return "", ErrMissingKey
	}
	if len(k) > MaxKeyLen {
		return "", apperror.BadRequest("Idempotency-Key is too long", nil)
	}
	return k, nil
}

// Lookup is the repository hook Resolve uses to find an existing payment for a
// key. PaymentRepository.GetByIdempotencyKey satisfies it.
type Lookup func(ctx context.Context, key string) (*entity.Payment, error)

// Resolve returns the payment already created for key, or (nil, nil) when none
// exists yet. A not-found from the lookup is the expected "first request" path
// and is reported as no payment, not an error; any other error propagates.
func Resolve(ctx context.Context, key string, lookup Lookup) (*entity.Payment, error) {
	existing, err := lookup(ctx, key)
	if err != nil {
		if apperror.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return existing, nil
}

-- name: CreatePayment :one
INSERT INTO payments.payments (
    id, student_id, course_id, subscription_id, amount, currency,
    provider, provider_checkout_id, status, coupon_id, idempotency_key, metadata
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
)
RETURNING *;

-- name: GetPaymentByID :one
SELECT * FROM payments.payments WHERE id = $1;

-- GetPaymentByIdempotencyKey is the friendly idempotency path: a retried initiate
-- with the same client key returns the existing payment instead of charging again.
-- name: GetPaymentByIdempotencyKey :one
SELECT * FROM payments.payments WHERE idempotency_key = $1;

-- name: GetPaymentByProviderTxn :one
SELECT * FROM payments.payments
WHERE provider = $1 AND provider_transaction_id = $2;

-- UpdatePaymentStatusGuarded flips status only when the row still holds
-- from_status; :execrows lets the repo report changed=false (0 rows) without an
-- error so a re-delivered webhook cannot double-settle an already-settled payment.
-- paid_at is stamped only on the transition to completed.
-- name: UpdatePaymentStatusGuarded :execrows
UPDATE payments.payments
SET status = @to_status,
    provider_transaction_id = COALESCE(sqlc.narg('provider_txn'), provider_transaction_id),
    failure_reason = COALESCE(sqlc.narg('failure_reason'), failure_reason),
    paid_at = CASE WHEN sqlc.arg('set_paid')::boolean THEN now() ELSE paid_at END,
    updated_at = now()
WHERE id = @id AND status = @from_status;

-- name: ListPaymentsByStudent :many
SELECT * FROM payments.payments
WHERE student_id = $1
ORDER BY created_at DESC
LIMIT $2;

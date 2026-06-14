-- name: CreateSubscription :one
INSERT INTO payments.subscriptions (
    id, student_id, plan, status, current_period_start, current_period_end,
    provider, provider_subscription_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: GetSubscriptionByID :one
SELECT * FROM payments.subscriptions WHERE id = $1;

-- GetActiveSubscriptionForStudent returns the all-access subscription only when
-- it grants access now (active and within its period).
-- name: GetActiveSubscriptionForStudent :one
SELECT * FROM payments.subscriptions
WHERE student_id = $1 AND status = 'active' AND current_period_end > now()
ORDER BY current_period_end DESC
LIMIT 1;

-- name: UpdateSubscriptionStatus :exec
UPDATE payments.subscriptions
SET status = $2,
    canceled_at = CASE WHEN $2 = 'canceled' THEN now() ELSE canceled_at END,
    updated_at = now()
WHERE id = $1;

-- ListExpiringSubscriptions returns active subscriptions whose period has ended,
-- for the renewal/expiry sweeper (Skill 4).
-- name: ListExpiringSubscriptions :many
SELECT * FROM payments.subscriptions
WHERE status = 'active' AND current_period_end <= now()
ORDER BY current_period_end ASC
LIMIT $1;

-- ExtendSubscriptionPeriod pushes an active subscription's period end out (a
-- renewal or stacked purchase) and (re)activates it; the period-extension
-- primitive used by idempotent Activate.
-- name: ExtendSubscriptionPeriod :exec
UPDATE payments.subscriptions
SET current_period_end = $2,
    status = 'active',
    updated_at = now()
WHERE id = $1;

-- CancelSubscriptionAtPeriodEnd flags an active subscription as canceled (it will
-- not renew) while keeping access until current_period_end; the sweeper expires
-- it once the period lapses. Guarded so a repeated cancel is a no-op.
-- name: CancelSubscriptionAtPeriodEnd :execrows
UPDATE payments.subscriptions
SET canceled_at = now(),
    updated_at = now()
WHERE id = $1 AND status = 'active' AND canceled_at IS NULL;

-- ExpireLapsedSubscription transitions an active subscription whose period has
-- ended to 'expired'. Guarded on status + period so a renewal that landed between
-- the sweeper's list and this update is never clobbered, and a repeated tick is a
-- no-op (the row is no longer active).
-- name: ExpireLapsedSubscription :execrows
UPDATE payments.subscriptions
SET status = 'expired',
    updated_at = now()
WHERE id = $1 AND status = 'active' AND current_period_end <= now();

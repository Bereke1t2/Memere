-- name: GetCouponByCode :one
SELECT * FROM payments.coupons WHERE code = $1;

-- IncrementCouponUse bumps used_count only while it is still below max_uses (or
-- max_uses is NULL = unlimited). :execrows reports redeemed=false (0 rows) when
-- the coupon is already exhausted, even under concurrent fulfillment. Called only
-- inside the successful-payment transaction.
-- name: IncrementCouponUse :execrows
UPDATE payments.coupons
SET used_count = used_count + 1, updated_at = now()
WHERE id = $1
  AND (max_uses IS NULL OR used_count < max_uses);

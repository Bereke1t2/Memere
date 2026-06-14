-- Revenue reporting aggregations over completed payments (spec §1.6, §10).
-- All read-only; the usecase layer applies role/ownership scoping and the
-- configurable teacher/platform split.

-- SumCompletedPayments is the platform gross + unit count of settled payments in
-- a window (keyed on paid_at, set at completion).
-- name: SumCompletedPayments :one
SELECT COALESCE(SUM(amount), 0)::numeric AS gross,
       COUNT(*)::bigint AS units
FROM payments.payments
WHERE status = 'completed'
  AND paid_at BETWEEN $1 AND $2;

-- SumRefundedPayments is the refunded total + count in a window (keyed on the
-- updated_at that stamped the refund).
-- name: SumRefundedPayments :one
SELECT COALESCE(SUM(amount), 0)::numeric AS refunded,
       COUNT(*)::bigint AS units
FROM payments.payments
WHERE status = 'refunded'
  AND updated_at BETWEEN $1 AND $2;

-- TeacherGrossEarnings sums completed payments for the teacher's courses in a
-- window; the usecase multiplies by the configured teacher share.
-- name: TeacherGrossEarnings :one
SELECT COALESCE(SUM(p.amount), 0)::numeric AS gross,
       COUNT(*)::bigint AS units
FROM payments.payments p
JOIN courses.courses c ON c.id = p.course_id
WHERE c.teacher_id = $1
  AND p.status = 'completed'
  AND p.paid_at BETWEEN $2 AND $3;

-- CourseSales is the lifetime gross + units sold for a single course.
-- name: CourseSales :one
SELECT COALESCE(SUM(amount), 0)::numeric AS gross,
       COUNT(*)::bigint AS units
FROM payments.payments
WHERE course_id = $1
  AND status = 'completed';

-- CompletedRevenueByProvider breaks the platform gross down per provider in a
-- window (for the admin revenue report).
-- name: CompletedRevenueByProvider :many
SELECT provider,
       COALESCE(SUM(amount), 0)::numeric AS gross,
       COUNT(*)::bigint AS units
FROM payments.payments
WHERE status = 'completed'
  AND paid_at BETWEEN $1 AND $2
GROUP BY provider
ORDER BY provider;

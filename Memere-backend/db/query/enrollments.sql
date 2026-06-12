-- CreateEnrollment grants access once. ON CONFLICT (student_id, course_id) DO
-- NOTHING makes a re-processed payment a no-op insert (exactly-once grant); the
-- caller treats a 0-row result as "already enrolled", not an error.
-- name: CreateEnrollment :one
INSERT INTO payments.enrollments (id, student_id, course_id, source, expires_at)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (student_id, course_id) DO NOTHING
RETURNING *;

-- name: EnrollmentExists :one
SELECT EXISTS (
    SELECT 1 FROM payments.enrollments
    WHERE student_id = $1 AND course_id = $2
);

-- GetActiveEnrollment returns the enrollment only when it grants access now
-- (no expiry, or not yet expired).
-- name: GetActiveEnrollment :one
SELECT * FROM payments.enrollments
WHERE student_id = $1 AND course_id = $2
  AND (expires_at IS NULL OR expires_at > now());

-- name: ListEnrollmentsByStudent :many
SELECT * FROM payments.enrollments
WHERE student_id = $1
ORDER BY enrolled_at DESC
LIMIT $2;

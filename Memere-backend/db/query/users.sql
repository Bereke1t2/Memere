-- name: CreateUser :one
INSERT INTO auth.users (
    email, phone, password_hash, role, first_name, last_name, avatar_url,
    is_active, is_email_verified
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM auth.users
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetUserByEmail :one
SELECT * FROM auth.users
WHERE email = $1 AND deleted_at IS NULL;

-- name: UpdateUser :one
UPDATE auth.users
SET phone = $2,
    password_hash = $3,
    first_name = $4,
    last_name = $5,
    avatar_url = $6,
    is_active = $7,
    is_email_verified = $8,
    email_verification_token = $9,
    password_reset_token = $10,
    password_reset_expires_at = $11
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteUser :exec
UPDATE auth.users
SET deleted_at = now()
WHERE id = $1 AND deleted_at IS NULL;

-- name: SetLastLogin :exec
UPDATE auth.users
SET last_login_at = $2
WHERE id = $1 AND deleted_at IS NULL;

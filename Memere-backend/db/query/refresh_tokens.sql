-- name: CreateRefreshToken :one
INSERT INTO auth.refresh_tokens (user_id, token_hash, device_info, expires_at)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetRefreshTokenByHash :one
SELECT * FROM auth.refresh_tokens
WHERE token_hash = $1;

-- name: RevokeRefreshToken :exec
UPDATE auth.refresh_tokens
SET revoked_at = now()
WHERE id = $1 AND revoked_at IS NULL;

-- name: RevokeAllRefreshTokensForUser :exec
UPDATE auth.refresh_tokens
SET revoked_at = now()
WHERE user_id = $1 AND revoked_at IS NULL;

-- name: DeleteExpiredRefreshTokens :exec
DELETE FROM auth.refresh_tokens
WHERE expires_at < now();

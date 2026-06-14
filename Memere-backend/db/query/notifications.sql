-- name: CreateNotification :one
INSERT INTO notifications.notifications (user_id, type, title, body, data)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, user_id, type, title, body, data, read_at, created_at;

-- name: ListNotificationsByUser :many
SELECT id, user_id, type, title, body, data, read_at, created_at
FROM notifications.notifications
WHERE user_id = $1
ORDER BY read_at NULLS FIRST, created_at DESC
LIMIT $2;

-- name: GetNotificationByID :one
SELECT id, user_id, type, title, body, data, read_at, created_at
FROM notifications.notifications
WHERE id = $1;

-- name: MarkNotificationRead :exec
UPDATE notifications.notifications
SET read_at = now()
WHERE id = $1 AND user_id = $2 AND read_at IS NULL;

-- name: MarkAllNotificationsRead :exec
UPDATE notifications.notifications
SET read_at = now()
WHERE user_id = $1 AND read_at IS NULL;

-- name: UnreadNotificationCount :one
SELECT COUNT(*)::int AS count
FROM notifications.notifications
WHERE user_id = $1 AND read_at IS NULL;

-- name: UpsertDeviceToken :exec
INSERT INTO notifications.device_tokens (user_id, fcm_token, platform)
VALUES ($1, $2, $3)
ON CONFLICT (fcm_token) DO UPDATE SET user_id = EXCLUDED.user_id, platform = EXCLUDED.platform;

-- name: DeleteDeviceToken :exec
DELETE FROM notifications.device_tokens WHERE fcm_token = $1;

-- name: DeleteInvalidDeviceTokens :exec
DELETE FROM notifications.device_tokens WHERE fcm_token = ANY($1::text[]);

-- name: ListDeviceTokensByUser :many
SELECT id, user_id, fcm_token, platform, created_at
FROM notifications.device_tokens
WHERE user_id = $1;

-- name: GetNotificationPreferences :one
SELECT user_id, push_enabled, email_enabled, updated_at
FROM notifications.preferences
WHERE user_id = $1;

-- name: UpsertNotificationPreferences :exec
INSERT INTO notifications.preferences (user_id, push_enabled, email_enabled, updated_at)
VALUES ($1, $2, $3, now())
ON CONFLICT (user_id) DO UPDATE SET
    push_enabled  = EXCLUDED.push_enabled,
    email_enabled = EXCLUDED.email_enabled,
    updated_at    = now();

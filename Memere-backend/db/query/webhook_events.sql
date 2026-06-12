-- InsertWebhookIfNew is the dedup primitive (spec §7.3). ON CONFLICT
-- (provider, provider_event_id) DO NOTHING means a re-delivered event returns no
-- row, signalling the caller to skip re-processing. A returned id means this is
-- the first delivery and the caller should fulfill it.
-- name: InsertWebhookIfNew :one
INSERT INTO payments.webhook_events (id, provider, provider_event_id, event_type, payload)
VALUES (gen_random_uuid(), $1, $2, $3, $4)
ON CONFLICT (provider, provider_event_id) DO NOTHING
RETURNING id;

-- name: MarkWebhookProcessed :exec
UPDATE payments.webhook_events
SET processed_at = now()
WHERE provider = $1 AND provider_event_id = $2;

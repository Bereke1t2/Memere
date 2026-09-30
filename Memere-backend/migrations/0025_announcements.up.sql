CREATE TABLE IF NOT EXISTS notifications.announcements (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sender_id       UUID NOT NULL REFERENCES auth.users(id),
    title           TEXT NOT NULL,
    body            TEXT NOT NULL,
    segment         TEXT NOT NULL,
    data            JSONB,
    recipient_count INT NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS announcements_created_at_idx
    ON notifications.announcements (created_at DESC);

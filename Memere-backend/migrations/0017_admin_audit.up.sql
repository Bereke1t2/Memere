CREATE TABLE IF NOT EXISTS auth.admin_audit_log (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    actor_id    UUID NOT NULL REFERENCES auth.users(id),
    action      TEXT NOT NULL,
    target_type TEXT NOT NULL,
    target_id   UUID,
    details     JSONB,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS admin_audit_actor_idx
    ON auth.admin_audit_log (actor_id, created_at);

CREATE INDEX IF NOT EXISTS admin_audit_target_idx
    ON auth.admin_audit_log (target_type, target_id);

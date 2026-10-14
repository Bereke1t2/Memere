-- name: InsertAdminAuditLog :exec
INSERT INTO auth.admin_audit_log (actor_id, action, target_type, target_id, details)
VALUES ($1, $2, $3, $4, $5::jsonb);

import "server-only";
import { dbPool } from "@/lib/db/neon";
import type { Announcement } from "@/lib/api/schemas";

/**
 * Ensure notifications.announcements table exists if not created yet.
 */
async function ensureAnnouncementsTable() {
  const sql = `CREATE TABLE IF NOT EXISTS notifications.announcements (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sender_id       UUID NOT NULL REFERENCES auth.users(id),
    title           TEXT NOT NULL,
    body            TEXT NOT NULL,
    segment         TEXT NOT NULL,
    data            JSONB,
    recipient_count INT NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
  );`;
  await dbPool.query(sql).catch((err) => {
    console.warn("[ensureAnnouncementsTable] Warning:", err.message);
  });
}

/**
 * List all broadcast announcements with recipient counts.
 * Aggregates across notifications.announcements, notifications.notifications, and auth.admin_audit_log.
 */
export async function listAnnouncementsFromDB(limit = 50): Promise<{ announcements: Announcement[] }> {
  await ensureAnnouncementsTable();

  const query = `
    WITH tbl_announcements AS (
      SELECT id::text AS id, title, body, recipient_count, created_at
      FROM notifications.announcements
    ),
    notifs AS (
      SELECT MIN(id::text) AS id, title, body, COUNT(*)::int AS recipient_count, MAX(created_at) AS created_at
      FROM notifications.notifications
      WHERE type = 'announcement'
      GROUP BY title, body
    ),
    audits AS (
      SELECT 
        id::text AS id,
        COALESCE(details->>'title', 'Announcement') AS title,
        COALESCE(details->>'body', '') AS body,
        COALESCE((details->>'recipients')::int, 0) AS recipient_count,
        created_at
      FROM auth.admin_audit_log
      WHERE action = 'broadcast.send'
    ),
    combined AS (
      SELECT id, title, body, recipient_count, created_at FROM tbl_announcements
      UNION ALL
      SELECT n.id, n.title, n.body, n.recipient_count, n.created_at FROM notifs n
      WHERE NOT EXISTS (SELECT 1 FROM tbl_announcements a WHERE a.title = n.title)
      UNION ALL
      SELECT au.id, au.title, au.body, au.recipient_count, au.created_at FROM audits au
      WHERE NOT EXISTS (SELECT 1 FROM tbl_announcements a WHERE a.title = au.title)
        AND NOT EXISTS (SELECT 1 FROM notifs n WHERE n.title = au.title)
    )
    SELECT id, title, body, recipient_count, created_at
    FROM combined
    ORDER BY created_at DESC
    LIMIT $1;
  `;

  try {
    const res = await dbPool.query(query, [limit]);
    const items: Announcement[] = res.rows.map((row) => ({
      id: row.id,
      title: row.title || "Announcement",
      body: row.body || "",
      recipient_count: Number(row.recipient_count) || 0,
      created_at: row.created_at ? new Date(row.created_at).toISOString() : new Date().toISOString(),
    }));
    return { announcements: items };
  } catch (err) {
    console.error("[listAnnouncementsFromDB] Query error:", err);
    return { announcements: [] };
  }
}

/**
 * Broadcast an announcement directly to Neon PostgreSQL.
 * 1. Saves record into notifications.announcements.
 * 2. Fans out notification entries into notifications.notifications for mobile student/teacher inboxes.
 * 3. Records audit log in auth.admin_audit_log.
 */
export async function createAnnouncementInDB(
  input: {
    title: string;
    body: string;
    segment: "all" | "students" | "teachers" | "subscribers";
    data?: Record<string, string>;
  },
  actorId?: string
): Promise<void> {
  await ensureAnnouncementsTable();

  // Determine fallback admin ID if actorId is not provided
  let adminId = actorId;
  if (!adminId) {
    const adminRes = await dbPool.query("SELECT id FROM auth.users WHERE role = 'admin' LIMIT 1");
    adminId = adminRes.rows[0]?.id;
  }
  if (!adminId) {
    throw new Error("No admin user found to record announcement sender");
  }

  // Find target recipient user IDs
  let userQuery = "";
  if (input.segment === "students") {
    userQuery = "SELECT id FROM auth.users WHERE role = 'student' AND status = 'active'";
  } else if (input.segment === "teachers") {
    userQuery = "SELECT id FROM auth.users WHERE role = 'teacher' AND status = 'active'";
  } else {
    userQuery = "SELECT id FROM auth.users WHERE role IN ('student', 'teacher') AND status = 'active'";
  }

  const usersRes = await dbPool.query(userQuery);
  const recipientCount = usersRes.rows.length;
  const jsonData = JSON.stringify(input.data ?? {});

  // 1. Save in notifications.announcements
  await dbPool.query(
    `INSERT INTO notifications.announcements (sender_id, title, body, segment, data, recipient_count)
     VALUES ($1, $2, $3, $4, $5::jsonb, $6)`,
    [adminId, input.title, input.body, input.segment, jsonData, recipientCount]
  );

  // 2. Fanout to notifications.notifications
  if (usersRes.rows.length > 0) {
    const values: string[] = [];
    const params: unknown[] = [input.title, input.body, jsonData];
    let paramIdx = 4;

    for (const row of usersRes.rows) {
      values.push(`($${paramIdx}, 'announcement', $1, $2, $3::jsonb, now())`);
      params.push(row.id);
      paramIdx++;
    }

    const fanoutSql = `
      INSERT INTO notifications.notifications (user_id, type, title, body, data, created_at)
      VALUES ${values.join(", ")}
    `;
    await dbPool.query(fanoutSql, params);
  }

  // 3. Write to audit log
  const auditDetails = JSON.stringify({
    title: input.title,
    body: input.body,
    segment: input.segment,
    recipients: recipientCount,
  });
  await dbPool.query(
    `INSERT INTO auth.admin_audit_log (actor_id, action, target_type, details)
     VALUES ($1, 'broadcast.send', 'system', $2::jsonb)`,
    [adminId, auditDetails]
  );
}

/**
 * Delete an announcement group across all tables.
 * Purges notifications.announcements, notifications.notifications, and auth.admin_audit_log.
 */
export async function deleteAnnouncementFromDB(id: string, actorId?: string): Promise<void> {
  await ensureAnnouncementsTable();

  const query = `
    WITH target AS (
      SELECT title, body FROM notifications.announcements WHERE id::text = $1 LIMIT 1
      UNION ALL
      SELECT title, body FROM notifications.notifications WHERE id::text = $1 LIMIT 1
      UNION ALL
      SELECT details->>'title' AS title, COALESCE(details->>'body', '') AS body FROM auth.admin_audit_log WHERE id::text = $1 LIMIT 1
    ),
    del_ann AS (
      DELETE FROM notifications.announcements
      WHERE id::text = $1 OR title IN (SELECT title FROM target WHERE title IS NOT NULL AND title != '')
      RETURNING id
    ),
    del_notif AS (
      DELETE FROM notifications.notifications
      WHERE id::text = $1 OR title IN (SELECT title FROM target WHERE title IS NOT NULL AND title != '')
      RETURNING id
    ),
    del_audit AS (
      DELETE FROM auth.admin_audit_log
      WHERE id::text = $1 OR (action = 'broadcast.send' AND details->>'title' IN (SELECT title FROM target WHERE title IS NOT NULL AND title != ''))
      RETURNING id
    )
    SELECT 
      (SELECT COUNT(*) FROM del_ann) AS deleted_announcements,
      (SELECT COUNT(*) FROM del_notif) AS deleted_notifications,
      (SELECT COUNT(*) FROM del_audit) AS deleted_audits;
  `;

  await dbPool.query(query, [id]);

  if (actorId) {
    const auditDetails = JSON.stringify({ announcement_id: id });
    await dbPool.query(
      `INSERT INTO auth.admin_audit_log (actor_id, action, target_type, details)
       VALUES ($1, 'announcement.delete', 'announcement', $2::jsonb)`,
      [actorId, auditDetails]
    ).catch(() => {});
  }
}

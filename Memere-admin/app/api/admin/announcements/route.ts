import { listAnnouncementsFromDB, createAnnouncementInDB } from "@/lib/api/announcements-db";
import { broadcast } from "@/lib/api/endpoints";
import { getRouteAdminSession } from "@/lib/auth/session";
import { z } from "zod";

export const dynamic = "force-dynamic";

export async function GET(req: Request) {
  const session = await getRouteAdminSession();
  if (!session) {
    return Response.json({ message: "Admin authentication required." }, { status: 401 });
  }

  try {
    const { searchParams } = new URL(req.url);
    const limit = Number(searchParams.get("limit")) || 50;
    const data = await listAnnouncementsFromDB(limit);
    return Response.json(data);
  } catch (err) {
    console.error("[GET /api/admin/announcements] Error:", err);
    return Response.json({ message: "Failed to fetch announcements." }, { status: 500 });
  }
}

const BodySchema = z.object({
  title: z.string().min(1),
  body: z.string().min(1),
  segment: z.enum(["all", "students", "teachers", "subscribers"]),
  data: z.record(z.string(), z.string()).optional(),
});

export async function POST(req: Request) {
  const session = await getRouteAdminSession();
  if (!session) {
    return Response.json({ message: "Admin authentication required." }, { status: 401 });
  }

  const raw = await req.json().catch(() => null);
  const parsed = BodySchema.safeParse(raw);

  if (!parsed.success) {
    return Response.json({ message: "Invalid request body." }, { status: 400 });
  }

  try {
    // 1. Direct persistent DB creation & fanout
    await createAnnouncementInDB(parsed.data, session.user.id);

    // 2. Best-effort broadcast to backend notification dispatcher (fire-and-forget)
    broadcast(parsed.data).catch((err) => {
      console.warn("[POST /api/admin/announcements] Backend notify ping (optional):", err.message);
    });

    return new Response(null, { status: 204 });
  } catch (err) {
    console.error("[POST /api/admin/announcements] Error:", err);
    return Response.json({ message: "Failed to send announcement." }, { status: 500 });
  }
}

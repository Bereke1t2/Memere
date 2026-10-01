import { deleteAnnouncementFromDB } from "@/lib/api/announcements-db";
import { deleteAnnouncement } from "@/lib/api/endpoints";
import { getRouteAdminSession } from "@/lib/auth/session";

export const dynamic = "force-dynamic";

export async function DELETE(
  _req: Request,
  { params }: { params: Promise<{ id: string }> }
) {
  const { id } = await params;
  if (!id) {
    return Response.json({ message: "Announcement ID is required." }, { status: 400 });
  }

  try {
    const session = await getRouteAdminSession().catch(() => null);
    const actorId = session?.user?.id;

    // 1. Direct persistent DB deletion
    await deleteAnnouncementFromDB(id, actorId);

    // 2. Best-effort backend deletion ping
    deleteAnnouncement(id).catch((err) => {
      console.warn("[DELETE /api/admin/announcements/[id]] Backend delete ping (optional):", err?.message);
    });

    return new Response(null, { status: 204 });
  } catch (err) {
    console.error("[DELETE /api/admin/announcements/[id]] Error:", err);
    const msg = err instanceof Error ? err.message : "Failed to delete announcement.";
    return Response.json({ message: msg }, { status: 500 });
  }
}

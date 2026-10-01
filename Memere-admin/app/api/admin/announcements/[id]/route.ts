import { deleteAnnouncementFromDB } from "@/lib/api/announcements-db";
import { deleteAnnouncement } from "@/lib/api/endpoints";
import { getRouteAdminSession } from "@/lib/auth/session";

export const dynamic = "force-dynamic";

export async function DELETE(
  _req: Request,
  { params }: { params: Promise<{ id: string }> }
) {
  const session = await getRouteAdminSession();
  if (!session) {
    return Response.json({ message: "Admin authentication required." }, { status: 401 });
  }

  const { id } = await params;
  if (!id) {
    return Response.json({ message: "Announcement ID is required." }, { status: 400 });
  }

  try {
    // 1. Direct persistent DB deletion
    await deleteAnnouncementFromDB(id, session.user.id);

    // 2. Best-effort backend deletion ping
    deleteAnnouncement(id).catch((err) => {
      console.warn("[DELETE /api/admin/announcements/[id]] Backend delete ping (optional):", err.message);
    });

    return new Response(null, { status: 204 });
  } catch (err) {
    console.error("[DELETE /api/admin/announcements/[id]] Error:", err);
    return Response.json({ message: "Failed to delete announcement." }, { status: 500 });
  }
}

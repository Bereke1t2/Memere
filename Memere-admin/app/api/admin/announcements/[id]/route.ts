import { deleteAnnouncement } from "@/lib/api/endpoints";
import { ApiError, friendlyMessage } from "@/lib/api/errors";

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
    await deleteAnnouncement(id);
    return new Response(null, { status: 204 });
  } catch (err) {
    if (err instanceof ApiError) {
      return Response.json({ message: friendlyMessage(err) }, { status: err.status });
    }
    return Response.json({ message: "Failed to delete announcement." }, { status: 500 });
  }
}

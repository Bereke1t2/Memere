import { requireStaff } from "@/lib/auth/session";
import { canManageContent } from "@/lib/auth/roles";
import { retryTranscode } from "@/lib/api/endpoints";
import { ApiError, friendlyMessage } from "@/lib/api/errors";

export async function POST(_req: Request, { params }: { params: Promise<{ id: string }> }) {
  try {
    const { user } = await requireStaff();
    if (!canManageContent(user)) return Response.json({ message: "Forbidden" }, { status: 403 });
    const { id } = await params;
    await retryTranscode(id);
    return new Response(null, { status: 204 });
  } catch (err) {
    if (err instanceof ApiError) return Response.json({ message: friendlyMessage(err) }, { status: err.status });
    return Response.json({ message: "Failed to retry transcode." }, { status: 500 });
  }
}

import { requireStaff } from "@/lib/auth/session";
import { getVideoStatus } from "@/lib/api/endpoints";
import { ApiError, friendlyMessage } from "@/lib/api/errors";

export const dynamic = "force-dynamic";

export async function GET(_req: Request, { params }: { params: Promise<{ id: string }> }) {
  try {
    await requireStaff();
    const { id } = await params;
    return Response.json(await getVideoStatus(id));
  } catch (err) {
    if (err instanceof ApiError) return Response.json({ message: friendlyMessage(err) }, { status: err.status });
    return Response.json({ message: "Failed to get video status." }, { status: 500 });
  }
}

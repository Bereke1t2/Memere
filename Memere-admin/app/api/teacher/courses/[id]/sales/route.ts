import { getRouteStaffSession } from "@/lib/auth/session";
import { canManageContent } from "@/lib/auth/roles";
import { getCourseSales } from "@/lib/api/endpoints";
import { ApiError, friendlyMessage } from "@/lib/api/errors";

export const dynamic = "force-dynamic";

export async function GET(
  _req: Request,
  { params }: { params: Promise<{ id: string }> }
) {
  try {
    const session = await getRouteStaffSession();
    if (!session) return Response.json({ message: "Unauthorized", code: "UNAUTHORIZED" }, { status: 401 });
    if (!canManageContent(session.user)) return Response.json({ message: "Forbidden" }, { status: 403 });

    const { id } = await params;
    const sales = await getCourseSales(id);
    return Response.json(sales);
  } catch (err) {
    if (err instanceof ApiError) {
      return Response.json({ message: friendlyMessage(err), details: err.details }, { status: err.status });
    }
    const message = err instanceof Error ? err.message : "Failed to load sales.";
    return Response.json({ message }, { status: 500 });
  }
}

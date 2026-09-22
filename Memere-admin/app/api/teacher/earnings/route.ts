import { getRouteStaffSession } from "@/lib/auth/session";
import { canManageContent } from "@/lib/auth/roles";
import { getMyEarnings } from "@/lib/api/endpoints";
import { ApiError, friendlyMessage } from "@/lib/api/errors";

export const dynamic = "force-dynamic";

export async function GET(req: Request) {
  try {
    const session = await getRouteStaffSession();
    if (!session) return Response.json({ message: "Unauthorized", code: "UNAUTHORIZED" }, { status: 401 });
    if (!canManageContent(session.user)) return Response.json({ message: "Forbidden" }, { status: 403 });

    const { searchParams } = new URL(req.url);
    const today = new Date().toISOString().split("T")[0];
    const monthAgo = new Date(Date.now() - 30 * 24 * 60 * 60 * 1000).toISOString().split("T")[0];
    const from = searchParams.get("from") ?? monthAgo;
    const to = searchParams.get("to") ?? today;

    const earnings = await getMyEarnings(from, to);
    return Response.json(earnings);
  } catch (err) {
    if (err instanceof ApiError) {
      return Response.json({ message: friendlyMessage(err), details: err.details }, { status: err.status });
    }
    const message = err instanceof Error ? err.message : "Failed to load earnings.";
    return Response.json({ message }, { status: 500 });
  }
}

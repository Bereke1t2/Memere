import { listAdminCourseRequests } from "@/lib/api/endpoints";
import { ApiError, friendlyMessage } from "@/lib/api/errors";

export const dynamic = "force-dynamic";

export async function GET(req: Request) {
  const { searchParams } = new URL(req.url);
  const status = searchParams.get("status") ?? undefined;
  const page = searchParams.get("page") ? parseInt(searchParams.get("page")!, 10) : undefined;
  const limit = searchParams.get("limit") ? parseInt(searchParams.get("limit")!, 10) : undefined;

  try {
    const data = await listAdminCourseRequests({ status, page, limit });
    return Response.json(data);
  } catch (err) {
    if (err instanceof ApiError) {
      return Response.json({ message: friendlyMessage(err) }, { status: err.status });
    }
    return Response.json({ message: "Failed to load course requests." }, { status: 500 });
  }
}

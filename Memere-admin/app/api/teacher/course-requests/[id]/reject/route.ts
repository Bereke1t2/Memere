import { rejectTeacherCourseRequest } from "@/lib/api/endpoints";
import { ApiError, friendlyMessage } from "@/lib/api/errors";

export async function POST(
  req: Request,
  { params }: { params: Promise<{ id: string }> }
) {
  const { id } = await params;
  const body = await req.json().catch(() => ({}));
  const reason = (body.reason ?? "").trim();

  try {
    await rejectTeacherCourseRequest(id, reason || undefined);
    return Response.json({ message: "Request rejected." });
  } catch (err) {
    if (err instanceof ApiError) {
      return Response.json({ message: friendlyMessage(err) }, { status: err.status });
    }
    return Response.json({ message: "Failed to reject course request." }, { status: 500 });
  }
}

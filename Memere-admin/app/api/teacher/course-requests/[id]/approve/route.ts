import { approveTeacherCourseRequest } from "@/lib/api/endpoints";
import { ApiError, friendlyMessage } from "@/lib/api/errors";

export async function POST(
  _req: Request,
  { params }: { params: Promise<{ id: string }> }
) {
  const { id } = await params;

  try {
    await approveTeacherCourseRequest(id);
    return Response.json({ message: "Request approved." });
  } catch (err) {
    if (err instanceof ApiError) {
      return Response.json({ message: friendlyMessage(err) }, { status: err.status });
    }
    return Response.json({ message: "Failed to approve course request." }, { status: 500 });
  }
}

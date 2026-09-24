import { revokeCourseAccess } from "@/lib/api/endpoints";
import { ApiError, friendlyMessage } from "@/lib/api/errors";

export async function DELETE(
  _req: Request,
  { params }: { params: Promise<{ id: string; courseId: string }> }
) {
  const { id, courseId } = await params;

  try {
    await revokeCourseAccess(id, courseId);
    return new Response(null, { status: 204 });
  } catch (err) {
    if (err instanceof ApiError) {
      return Response.json({ message: friendlyMessage(err) }, { status: err.status });
    }
    return Response.json({ message: "Failed to revoke course access." }, { status: 500 });
  }
}

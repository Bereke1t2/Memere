import { getUserCourseAccess } from "@/lib/api/endpoints";
import { ApiError, friendlyMessage } from "@/lib/api/errors";

export async function GET(
  _req: Request,
  { params }: { params: Promise<{ id: string }> }
) {
  const { id } = await params;

  try {
    const data = await getUserCourseAccess(id);
    return Response.json(data);
  } catch (err) {
    if (err instanceof ApiError) {
      return Response.json({ message: friendlyMessage(err) }, { status: err.status });
    }
    return Response.json({ message: "Failed to load user course access." }, { status: 500 });
  }
}

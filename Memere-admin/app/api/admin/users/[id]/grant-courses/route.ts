import { grantCourseAccess, grantAllCoursesAccess } from "@/lib/api/endpoints";
import { ApiError, friendlyMessage } from "@/lib/api/errors";

export async function POST(
  req: Request,
  { params }: { params: Promise<{ id: string }> }
) {
  const { id } = await params;
  const body = await req.json().catch(() => ({}));

  try {
    if (body.all_courses || body.grant_all) {
      await grantAllCoursesAccess(id);
    } else if (Array.isArray(body.course_ids)) {
      await grantCourseAccess(id, body.course_ids);
    } else {
      return Response.json({ message: "Invalid request payload." }, { status: 400 });
    }
    return Response.json({ message: "Course access granted successfully." });
  } catch (err) {
    if (err instanceof ApiError) {
      return Response.json({ message: friendlyMessage(err) }, { status: err.status });
    }
    return Response.json({ message: "Failed to grant course access." }, { status: 500 });
  }
}

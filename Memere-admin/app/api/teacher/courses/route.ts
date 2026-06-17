import { requireStaff } from "@/lib/auth/session";
import { listMyCourses, createCourse } from "@/lib/api/endpoints";
import { CreateCourseInputSchema } from "@/lib/api/schemas";
import { ApiError, friendlyMessage } from "@/lib/api/errors";

export const dynamic = "force-dynamic";

export async function GET(req: Request) {
  try {
    const { user } = await requireStaff();
    if (user.role !== "teacher") return Response.json({ message: "Forbidden" }, { status: 403 });

    const { searchParams } = new URL(req.url);
    const limit = Math.min(parseInt(searchParams.get("limit") ?? "20", 10), 100);
    const next_cursor = searchParams.get("next_cursor") ?? undefined;

    const result = await listMyCourses({ teacherId: user.id, limit, next_cursor });
    return Response.json({ items: result.data, next: result.next_cursor });
  } catch (err) {
    if (err instanceof ApiError) {
      return Response.json({ message: friendlyMessage(err) }, { status: err.status });
    }
    return Response.json({ message: "Failed to load courses." }, { status: 500 });
  }
}

export async function POST(req: Request) {
  try {
    const { user } = await requireStaff();
    if (user.role !== "teacher") return Response.json({ message: "Forbidden" }, { status: 403 });

    const body = await req.json().catch(() => ({}));
    const parsed = CreateCourseInputSchema.safeParse(body);
    if (!parsed.success) {
      return Response.json({ message: "Validation failed", details: parsed.error.flatten().fieldErrors }, { status: 400 });
    }

    const course = await createCourse(parsed.data);
    return Response.json(course, { status: 201 });
  } catch (err) {
    if (err instanceof ApiError) {
      return Response.json({ message: friendlyMessage(err) }, { status: err.status });
    }
    return Response.json({ message: "Failed to create course." }, { status: 500 });
  }
}

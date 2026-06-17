import { requireStaff } from "@/lib/auth/session";
import { listLessons, addLesson } from "@/lib/api/endpoints";
import { AddLessonInputSchema } from "@/lib/api/schemas";
import { ApiError, friendlyMessage } from "@/lib/api/errors";

export const dynamic = "force-dynamic";

export async function GET(_req: Request, { params }: { params: Promise<{ id: string }> }) {
  try {
    await requireStaff();
    const { id } = await params;
    const lessons = await listLessons(id);
    return Response.json(lessons);
  } catch (err) {
    if (err instanceof ApiError) return Response.json({ message: friendlyMessage(err) }, { status: err.status });
    return Response.json({ message: "Failed to load lessons." }, { status: 500 });
  }
}

export async function POST(req: Request, { params }: { params: Promise<{ id: string }> }) {
  try {
    const { user } = await requireStaff();
    if (user.role !== "teacher") return Response.json({ message: "Forbidden" }, { status: 403 });
    const { id } = await params;
    const body = await req.json().catch(() => ({}));
    const parsed = AddLessonInputSchema.safeParse(body);
    if (!parsed.success) return Response.json({ message: "Validation failed", details: parsed.error.flatten().fieldErrors }, { status: 400 });
    const lesson = await addLesson(id, parsed.data);
    return Response.json(lesson, { status: 201 });
  } catch (err) {
    if (err instanceof ApiError) return Response.json({ message: friendlyMessage(err) }, { status: err.status });
    return Response.json({ message: "Failed to add lesson." }, { status: 500 });
  }
}

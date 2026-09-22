import { getRouteStaffSession } from "@/lib/auth/session";
import { canManageContent } from "@/lib/auth/roles";
import { listLessons, addLesson } from "@/lib/api/endpoints";
import { AddLessonInputSchema } from "@/lib/api/schemas";
import { ApiError, friendlyMessage } from "@/lib/api/errors";

export const dynamic = "force-dynamic";

export async function GET(_req: Request, { params }: { params: Promise<{ id: string }> }) {
  try {
    const session = await getRouteStaffSession();
    if (!session) return Response.json({ message: "Unauthorized", code: "UNAUTHORIZED" }, { status: 401 });
    if (!canManageContent(session.user)) return Response.json({ message: "Forbidden" }, { status: 403 });
    const { id } = await params;
    const lessons = await listLessons(id);
    return Response.json(lessons);
  } catch (err) {
    if (err instanceof ApiError) return Response.json({ message: friendlyMessage(err), details: err.details }, { status: err.status });
    const message = err instanceof Error ? err.message : "Failed to load lessons.";
    return Response.json({ message }, { status: 500 });
  }
}

export async function POST(req: Request, { params }: { params: Promise<{ id: string }> }) {
  try {
    const session = await getRouteStaffSession();
    if (!session) return Response.json({ message: "Unauthorized", code: "UNAUTHORIZED" }, { status: 401 });
    if (!canManageContent(session.user)) return Response.json({ message: "Forbidden" }, { status: 403 });
    const { id } = await params;
    const body = await req.json().catch(() => ({}));
    const parsed = AddLessonInputSchema.safeParse(body);
    if (!parsed.success) return Response.json({ message: "Validation failed", details: parsed.error.flatten().fieldErrors }, { status: 400 });
    const lesson = await addLesson(id, parsed.data);
    return Response.json(lesson, { status: 201 });
  } catch (err) {
    if (err instanceof ApiError) return Response.json({ message: friendlyMessage(err), details: err.details }, { status: err.status });
    const message = err instanceof Error ? err.message : "Failed to add lesson.";
    return Response.json({ message }, { status: 500 });
  }
}

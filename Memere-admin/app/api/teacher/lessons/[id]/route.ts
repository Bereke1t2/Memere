import { z } from "zod";
import { getRouteStaffSession } from "@/lib/auth/session";
import { canManageContent } from "@/lib/auth/roles";
import { apiFetch } from "@/lib/api/server";
import { LessonSchema } from "@/lib/api/schemas";
import { ApiError, friendlyMessage } from "@/lib/api/errors";

export const dynamic = "force-dynamic";

const UpdateLessonSchema = z.object({
  title: z.string().optional(),
  type: z.string().optional(),
  is_free_preview: z.boolean().optional(),
  duration_seconds: z.number().min(0).optional(),
  is_published: z.boolean().optional(),
  order_index: z.number().int().min(0).optional(),
  quiz_id: z.string().optional().nullable(),
  content: z.string().optional().nullable(),
  pdf_url: z.string().optional().nullable(),
}).passthrough();

export async function PUT(
  req: Request,
  { params }: { params: Promise<{ id: string }> }
) {
  try {
    const session = await getRouteStaffSession();
    if (!session) return Response.json({ message: "Unauthorized", code: "UNAUTHORIZED" }, { status: 401 });
    if (!canManageContent(session.user)) return Response.json({ message: "Forbidden" }, { status: 403 });

    const { id } = await params;
    const body = await req.json().catch(() => ({}));
    const parsed = UpdateLessonSchema.safeParse(body);
    if (!parsed.success) {
      return Response.json({ message: "Invalid payload." }, { status: 400 });
    }

    const lesson = await apiFetch(`/lessons/${id}`, {
      method: "PUT",
      body: parsed.data,
      schema: LessonSchema,
    });
    return Response.json(lesson);
  } catch (err) {
    if (err instanceof ApiError) {
      return Response.json({ message: friendlyMessage(err), details: err.details }, { status: err.status });
    }
    const message = err instanceof Error ? err.message : "Failed to update lesson.";
    return Response.json({ message }, { status: 500 });
  }
}

export const POST = PUT;

export async function DELETE(
  _req: Request,
  { params }: { params: Promise<{ id: string }> }
) {
  try {
    const session = await getRouteStaffSession();
    if (!session) return Response.json({ message: "Unauthorized", code: "UNAUTHORIZED" }, { status: 401 });
    if (!canManageContent(session.user)) return Response.json({ message: "Forbidden" }, { status: 403 });

    const { id } = await params;
    await apiFetch(`/lessons/${id}`, { method: "DELETE" });
    return Response.json({ message: "Lesson deleted." });
  } catch (err) {
    if (err instanceof ApiError) {
      return Response.json({ message: friendlyMessage(err), details: err.details }, { status: err.status });
    }
    const message = err instanceof Error ? err.message : "Failed to delete lesson.";
    return Response.json({ message }, { status: 500 });
  }
}

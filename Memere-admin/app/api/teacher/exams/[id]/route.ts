import { getRouteStaffSession } from "@/lib/auth/session";
import { canManageContent } from "@/lib/auth/roles";
import { updateExam, deleteExam } from "@/lib/api/endpoints";
import { CreateExamInputSchema } from "@/lib/api/schemas";
import { ApiError, friendlyMessage } from "@/lib/api/errors";

export async function PUT(req: Request, { params }: { params: Promise<{ id: string }> }) {
  try {
    const session = await getRouteStaffSession();
    if (!session) return Response.json({ message: "Unauthorized", code: "UNAUTHORIZED" }, { status: 401 });
    if (!canManageContent(session.user)) return Response.json({ message: "Forbidden" }, { status: 403 });
    const { id } = await params;
    const body = await req.json().catch(() => ({}));
    const parsed = CreateExamInputSchema.partial().safeParse(body);
    if (!parsed.success) return Response.json({ message: "Validation failed", details: parsed.error.flatten().fieldErrors }, { status: 400 });
    const updated = await updateExam(id, parsed.data);
    return Response.json(updated);
  } catch (err) {
    if (err instanceof ApiError) return Response.json({ message: friendlyMessage(err), details: err.details }, { status: err.status });
    const message = err instanceof Error ? err.message : "Failed to update exam.";
    return Response.json({ message }, { status: 500 });
  }
}

export async function DELETE(_req: Request, { params }: { params: Promise<{ id: string }> }) {
  try {
    const session = await getRouteStaffSession();
    if (!session) return Response.json({ message: "Unauthorized", code: "UNAUTHORIZED" }, { status: 401 });
    if (!canManageContent(session.user)) return Response.json({ message: "Forbidden" }, { status: 403 });
    const { id } = await params;
    await deleteExam(id);
    return new Response(null, { status: 204 });
  } catch (err) {
    if (err instanceof ApiError) return Response.json({ message: friendlyMessage(err), details: err.details }, { status: err.status });
    const message = err instanceof Error ? err.message : "Failed to delete exam.";
    return Response.json({ message }, { status: 500 });
  }
}

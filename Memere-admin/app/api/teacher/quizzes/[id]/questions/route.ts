import { getRouteStaffSession } from "@/lib/auth/session";
import { canManageContent } from "@/lib/auth/roles";
import { addQuizQuestion } from "@/lib/api/endpoints";
import { AddQuizQuestionInputSchema } from "@/lib/api/schemas";
import { ApiError, friendlyMessage } from "@/lib/api/errors";

export async function POST(req: Request, { params }: { params: Promise<{ id: string }> }) {
  try {
    const session = await getRouteStaffSession();
    if (!session) return Response.json({ message: "Unauthorized", code: "UNAUTHORIZED" }, { status: 401 });
    if (!canManageContent(session.user)) return Response.json({ message: "Forbidden" }, { status: 403 });
    const { id } = await params;
    const body = await req.json().catch(() => ({}));
    const parsed = AddQuizQuestionInputSchema.safeParse(body);
    if (!parsed.success) return Response.json({ message: "Validation failed", details: parsed.error.flatten().fieldErrors }, { status: 400 });
    return Response.json(await addQuizQuestion(id, parsed.data), { status: 201 });
  } catch (err) {
    if (err instanceof ApiError) return Response.json({ message: friendlyMessage(err), details: err.details }, { status: err.status });
    const message = err instanceof Error ? err.message : "Failed to add question.";
    return Response.json({ message }, { status: 500 });
  }
}

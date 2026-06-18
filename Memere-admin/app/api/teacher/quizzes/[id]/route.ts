import { requireStaff } from "@/lib/auth/session";
import { updateQuiz } from "@/lib/api/endpoints";
import { CreateQuizInputSchema } from "@/lib/api/schemas";
import { ApiError, friendlyMessage } from "@/lib/api/errors";

export async function PUT(req: Request, { params }: { params: Promise<{ id: string }> }) {
  try {
    const { user } = await requireStaff();
    if (user.role !== "teacher") return Response.json({ message: "Forbidden" }, { status: 403 });
    const { id } = await params;
    const body = await req.json().catch(() => ({}));
    const parsed = CreateQuizInputSchema.partial().safeParse(body);
    if (!parsed.success) return Response.json({ message: "Validation failed", details: parsed.error.flatten().fieldErrors }, { status: 400 });
    return Response.json(await updateQuiz(id, parsed.data));
  } catch (err) {
    if (err instanceof ApiError) return Response.json({ message: friendlyMessage(err) }, { status: err.status });
    return Response.json({ message: "Failed to update quiz." }, { status: 500 });
  }
}

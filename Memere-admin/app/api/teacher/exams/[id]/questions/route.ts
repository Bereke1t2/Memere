import { z } from "zod";
import { requireStaff } from "@/lib/auth/session";
import { addExamQuestion } from "@/lib/api/endpoints";
import { ApiError, friendlyMessage } from "@/lib/api/errors";

const ExamQuestionInputSchema = z.object({
  marks: z.number().min(1),
  order_index: z.number().min(0),
  question_id: z.string().optional(),
});

export async function POST(req: Request, { params }: { params: Promise<{ id: string }> }) {
  try {
    const { user } = await requireStaff();
    if (user.role !== "teacher") return Response.json({ message: "Forbidden" }, { status: 403 });
    const { id } = await params;
    const body = await req.json().catch(() => ({}));
    const parsed = ExamQuestionInputSchema.safeParse(body);
    if (!parsed.success) return Response.json({ message: "Validation failed" }, { status: 400 });
    await addExamQuestion(id, parsed.data);
    return new Response(null, { status: 204 });
  } catch (err) {
    if (err instanceof ApiError) return Response.json({ message: friendlyMessage(err) }, { status: err.status });
    return Response.json({ message: "Failed to add question." }, { status: 500 });
  }
}

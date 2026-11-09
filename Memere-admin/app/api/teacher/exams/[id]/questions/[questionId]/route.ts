import { updateExamQuestion, removeExamQuestion } from "@/lib/api/endpoints";
import { UpdateExamQuestionInputSchema } from "@/lib/api/schemas";
import { withStaff, validationResponse } from "@/lib/api/route-guard";

type Params = { params: Promise<{ id: string; questionId: string }> };

export async function PUT(req: Request, { params }: Params) {
  const { id, questionId } = await params;
  const body = await req.json().catch(() => ({}));
  const parsed = UpdateExamQuestionInputSchema.safeParse(body);
  if (!parsed.success) return validationResponse(parsed.error.flatten().fieldErrors);
  return withStaff("Failed to update question.", () =>
    updateExamQuestion(id, questionId, parsed.data)
  );
}

// Detaching returns the exam with its re-derived total_marks, so the editor can
// show the new total without a second round trip.
export async function DELETE(_req: Request, { params }: Params) {
  const { id, questionId } = await params;
  return withStaff("Failed to remove question.", () => removeExamQuestion(id, questionId));
}

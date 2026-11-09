import { updateQuizQuestion, deleteQuizQuestion } from "@/lib/api/endpoints";
import { UpdateQuestionInputSchema } from "@/lib/api/schemas";
import { withStaff, validationResponse } from "@/lib/api/route-guard";

export async function PUT(req: Request, { params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  const body = await req.json().catch(() => ({}));
  const parsed = UpdateQuestionInputSchema.safeParse(body);
  if (!parsed.success) return validationResponse(parsed.error.flatten().fieldErrors);
  return withStaff("Failed to update question.", () => updateQuizQuestion(id, parsed.data));
}

export async function DELETE(_req: Request, { params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  return withStaff("Failed to delete question.", async () => {
    await deleteQuizQuestion(id);
    return undefined;
  });
}

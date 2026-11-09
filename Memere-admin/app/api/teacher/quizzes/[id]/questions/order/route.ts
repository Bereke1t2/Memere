import { reorderQuizQuestions } from "@/lib/api/endpoints";
import { ReorderQuestionsInputSchema } from "@/lib/api/schemas";
import { withStaff, validationResponse } from "@/lib/api/route-guard";

export async function PUT(req: Request, { params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  const body = await req.json().catch(() => ({}));
  const parsed = ReorderQuestionsInputSchema.safeParse(body);
  if (!parsed.success) return validationResponse(parsed.error.flatten().fieldErrors);
  return withStaff("Failed to reorder questions.", async () => {
    await reorderQuizQuestions(id, parsed.data.question_ids);
    return undefined;
  });
}

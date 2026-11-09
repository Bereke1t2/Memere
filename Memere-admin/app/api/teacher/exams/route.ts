import { listAllExams, createStandaloneExam } from "@/lib/api/endpoints";
import { CreateExamInputSchema } from "@/lib/api/schemas";
import { withStaff, validationResponse } from "@/lib/api/route-guard";

export async function GET(req: Request) {
  const url = new URL(req.url);
  const published = url.searchParams.get("published");
  return withStaff("Failed to list exams.", () =>
    listAllExams({
      limit: Number(url.searchParams.get("limit")) || undefined,
      after: url.searchParams.get("after") ?? undefined,
      subject: url.searchParams.get("subject") ?? undefined,
      grade: Number(url.searchParams.get("grade")) || undefined,
      published: published === null ? undefined : published === "true",
    })
  );
}

export async function POST(req: Request) {
  const body = await req.json().catch(() => ({}));
  const parsed = CreateExamInputSchema.safeParse(body);
  if (!parsed.success) return validationResponse(parsed.error.flatten().fieldErrors);
  return withStaff("Failed to create exam.", () => createStandaloneExam(parsed.data));
}

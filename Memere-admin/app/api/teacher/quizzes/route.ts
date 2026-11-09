import { listAllQuizzes, createStandaloneQuiz } from "@/lib/api/endpoints";
import { CreateQuizInputSchema } from "@/lib/api/schemas";
import { withStaff, validationResponse } from "@/lib/api/route-guard";

export async function GET(req: Request) {
  const url = new URL(req.url);
  const standalone = url.searchParams.get("standalone");
  const published = url.searchParams.get("published");
  return withStaff("Failed to list quizzes.", () =>
    listAllQuizzes({
      limit: Number(url.searchParams.get("limit")) || undefined,
      after: url.searchParams.get("after") ?? undefined,
      course_id: url.searchParams.get("course_id") ?? undefined,
      standalone: standalone === null ? undefined : standalone === "true",
      published: published === null ? undefined : published === "true",
    })
  );
}

export async function POST(req: Request) {
  const body = await req.json().catch(() => ({}));
  const parsed = CreateQuizInputSchema.safeParse(body);
  if (!parsed.success) return validationResponse(parsed.error.flatten().fieldErrors);
  return withStaff("Failed to create quiz.", () => createStandaloneQuiz(parsed.data));
}

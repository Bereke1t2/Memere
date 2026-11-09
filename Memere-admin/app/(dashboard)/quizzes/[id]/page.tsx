import { notFound } from "next/navigation";
import { getQuizForAuthoring } from "@/lib/api/endpoints";
import { ApiError } from "@/lib/api/errors";
import { Badge } from "@/components/ui/badge";
import { BreadcrumbSetter } from "@/lib/breadcrumb-context";
import { QuestionList } from "@/components/assessments/question-editor";
import { AddQuestionButton } from "@/components/assessments/add-question-button";

export const dynamic = "force-dynamic";

export default async function QuizDetailPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = await params;

  let authored;
  try {
    authored = await getQuizForAuthoring(id);
  } catch (err) {
    if (
      err instanceof ApiError &&
      (err.status === 404 || err.code === "NOT_FOUND" || err.code === "RESOURCE_NOT_FOUND")
    ) {
      notFound();
    }
    throw err;
  }

  const { quiz, questions } = authored;
  const totalPoints = questions.reduce((sum, q) => sum + (q.points ?? q.marks ?? 1), 0);

  return (
    <>
      <BreadcrumbSetter label={quiz.title} />

      <div className="flex flex-col gap-6">
        <div className="flex flex-wrap items-start justify-between gap-4">
          <div className="min-w-0">
            <h1 className="text-2xl font-semibold tracking-tight">{quiz.title}</h1>
            <div className="mt-2 flex flex-wrap items-center gap-2 text-xs text-muted-foreground">
              <Badge variant="secondary" className="text-xs">
                {questions.length} question{questions.length === 1 ? "" : "s"}
              </Badge>
              <Badge variant="secondary" className="text-xs">
                {totalPoints} points total
              </Badge>
              {!quiz.course_id && (
                <Badge variant={quiz.is_published ? "default" : "outline"} className="text-xs">
                  {quiz.is_published ? "Published" : "Draft"}
                </Badge>
              )}
              {quiz.time_limit_seconds != null && (
                <span>{Math.round(quiz.time_limit_seconds / 60)}m limit</span>
              )}
              <span>{quiz.pass_percentage}% to pass</span>
              {quiz.max_attempts != null && <span>Max {quiz.max_attempts} attempts</span>}
              {quiz.randomize_questions && <span>Randomized</span>}
            </div>
          </div>
          <AddQuestionButton scope={{ kind: "quiz", quizId: quiz.id }} nextOrderIndex={questions.length} />
        </div>

        {!quiz.course_id && !quiz.is_published && (
          <p className="rounded-md border border-amber-500/40 bg-amber-500/5 p-3 text-xs text-amber-700 dark:text-amber-400">
            This standalone quiz is a draft, so students cannot reach it yet. Publish it once the
            questions are final.
          </p>
        )}

        <QuestionList scope={{ kind: "quiz", quizId: quiz.id }} questions={questions} />
      </div>
    </>
  );
}

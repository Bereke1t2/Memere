import { notFound } from "next/navigation";
import { getExamForAuthoring } from "@/lib/api/endpoints";
import { ApiError } from "@/lib/api/errors";
import { Badge } from "@/components/ui/badge";
import { BreadcrumbSetter } from "@/lib/breadcrumb-context";
import { QuestionList } from "@/components/assessments/question-editor";
import { AddQuestionButton } from "@/components/assessments/add-question-button";
import { ExamPublishToggle } from "@/components/assessments/exam-publish-toggle";
import { formatGrade } from "@/lib/format";

export const dynamic = "force-dynamic";

export default async function ExamDetailPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = await params;

  let authored;
  try {
    authored = await getExamForAuthoring(id);
  } catch (err) {
    if (
      err instanceof ApiError &&
      (err.status === 404 || err.code === "NOT_FOUND" || err.code === "RESOURCE_NOT_FOUND")
    ) {
      notFound();
    }
    throw err;
  }

  const { exam, questions } = authored;
  // total_marks is derived server-side from the question set, so this sum is a
  // display of the same arithmetic rather than a second source of truth.
  const totalMarks = questions.reduce((sum, q) => sum + (q.marks ?? q.points ?? 1), 0);

  return (
    <>
      <BreadcrumbSetter label={exam.title} />

      <div className="flex flex-col gap-6">
        <div className="flex flex-wrap items-start justify-between gap-4">
          <div className="min-w-0">
            <h1 className="text-2xl font-semibold tracking-tight">{exam.title}</h1>
            <div className="mt-2 flex flex-wrap items-center gap-2 text-xs text-muted-foreground">
              <Badge variant={exam.is_published ? "default" : "outline"} className="text-xs">
                {exam.is_published ? "Published" : "Draft"}
              </Badge>
              <Badge variant="secondary" className="text-xs">
                {questions.length} question{questions.length === 1 ? "" : "s"}
              </Badge>
              <Badge variant="secondary" className="text-xs">
                {totalMarks} marks total
              </Badge>
              {exam.subject && <span>{exam.subject}</span>}
              {exam.grade != null && <span>{formatGrade(exam.grade)}</span>}
              {exam.duration_minutes != null && <span>{exam.duration_minutes}m</span>}
              <span>{exam.pass_marks ?? 0} marks to pass</span>
            </div>
          </div>
          <div className="flex items-center gap-2">
            <AddQuestionButton
              scope={{ kind: "exam", examId: exam.id }}
              nextOrderIndex={questions.length}
            />
            <ExamPublishToggle examId={exam.id} isPublished={exam.is_published ?? false} />
          </div>
        </div>

        {(exam.pass_marks ?? 0) > totalMarks && (
          <p className="rounded-md border border-destructive/40 bg-destructive/5 p-3 text-xs text-destructive">
            Pass marks ({exam.pass_marks}) exceed the {totalMarks} marks this question set is worth,
            so nobody can pass. Add more questions or lower the pass mark.
          </p>
        )}

        <QuestionList scope={{ kind: "exam", examId: exam.id }} questions={questions} />
      </div>
    </>
  );
}

import Link from "next/link";
import { BarChart3 } from "lucide-react";
import { requireStaff } from "@/lib/auth/session";
import { listAllExams, listMyCourses, listCourses } from "@/lib/api/endpoints";
import { Badge } from "@/components/ui/badge";
import { EmptyState } from "@/components/ui/empty-state";
import { CreateExamDialog } from "@/components/assessments/create-exam-dialog";
import { ImportAssessmentDialog } from "@/components/assessments/import-dialog";
import { formatGrade } from "@/lib/format";
import type { Exam } from "@/lib/api/schemas";

export const dynamic = "force-dynamic";

export default async function ExamsPage() {
  const { user } = await requireStaff();
  const isAdmin = user.role === "admin";

  const [examResult, courses] = await Promise.all([
    listAllExams({ limit: 100 }).catch(() => ({ data: [] as Exam[] })),
    // The two listings have different envelopes: the admin endpoint returns
    // { courses }, the teacher one { data }.
    (isAdmin
      ? listCourses({ limit: 100 }).then((r) => r.courses)
      : listMyCourses({ teacherId: user.id, limit: 100 }).then((r) => r.data)
    )
      .then((list) => list.map((c) => ({ id: c.id, title: c.title })))
      .catch(() => [] as { id: string; title: string }[]),
  ]);

  const exams = examResult.data;
  const courseTitles = new Map(courses.map((c) => [c.id, c.title]));

  return (
    <div className="flex flex-col gap-6">
      <div className="flex flex-wrap items-start justify-between gap-4">
        <div>
          <h1 className="text-2xl font-semibold tracking-tight">Exams</h1>
          <p className="mt-1 text-sm text-muted-foreground">
            Mock exams including unpublished drafts. Open one to edit its questions.
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Link
            href="/exam-analytics"
            className="inline-flex h-9 items-center gap-1.5 rounded-md border px-3 text-sm font-medium transition-colors hover:bg-muted/40"
          >
            <BarChart3 className="h-4 w-4" aria-hidden="true" />
            Analytics
          </Link>
          <ImportAssessmentDialog kind="exam" />
          <CreateExamDialog courses={courses} canCreateStandalone={isAdmin} />
        </div>
      </div>

      {exams.length === 0 ? (
        <EmptyState
          title="No exams yet"
          description="Create an exam, or import a whole set of questions from a JSON or HTML file."
        />
      ) : (
        <div className="flex flex-col gap-3">
          {exams.map((exam) => (
            <Link
              key={exam.id}
              href={`/exams/${exam.id}`}
              className="rounded-lg border p-4 transition-colors hover:bg-muted/40"
            >
              <div className="flex flex-wrap items-start justify-between gap-3">
                <div className="min-w-0">
                  <p className="text-sm font-medium">{exam.title}</p>
                  <div className="mt-1 flex flex-wrap items-center gap-x-3 gap-y-1 text-xs text-muted-foreground">
                    {exam.subject && <span>{exam.subject}</span>}
                    {exam.grade != null && <span>{formatGrade(exam.grade)}</span>}
                    {exam.duration_minutes != null && <span>{exam.duration_minutes}m</span>}
                    <span>
                      {exam.pass_marks ?? 0}/{exam.total_marks ?? 0} to pass
                    </span>
                    <span>{exam.course_id ? (courseTitles.get(exam.course_id) ?? "Course") : "Standalone"}</span>
                  </div>
                </div>
                <Badge variant={exam.is_published ? "default" : "outline"} className="shrink-0 text-xs">
                  {exam.is_published ? "Published" : "Draft"}
                </Badge>
              </div>
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}

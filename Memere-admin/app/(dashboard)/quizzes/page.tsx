import Link from "next/link";
import { requireStaff } from "@/lib/auth/session";
import { listAllQuizzes, listMyCourses, listCourses } from "@/lib/api/endpoints";
import { Badge } from "@/components/ui/badge";
import { EmptyState } from "@/components/ui/empty-state";
import { CreateQuizDialog } from "@/components/assessments/create-quiz-dialog";
import { ImportAssessmentDialog } from "@/components/assessments/import-dialog";
import type { Quiz } from "@/lib/api/schemas";

export const dynamic = "force-dynamic";

export default async function QuizzesPage() {
  const { user } = await requireStaff();
  const isAdmin = user.role === "admin";

  // The course picker only needs courses the caller may author under: an admin
  // sees all, a teacher only their own.
  const [quizResult, courses] = await Promise.all([
    listAllQuizzes({ limit: 100 }).catch(() => ({ data: [] as Quiz[] })),
    // The two listings have different envelopes: the admin endpoint returns
    // { courses }, the teacher one { data }.
    (isAdmin
      ? listCourses({ limit: 100 }).then((r) => r.courses)
      : listMyCourses({ teacherId: user.id, limit: 100 }).then((r) => r.data)
    )
      .then((list) => list.map((c) => ({ id: c.id, title: c.title })))
      .catch(() => [] as { id: string; title: string }[]),
  ]);

  const quizzes = quizResult.data;
  const courseTitles = new Map(courses.map((c) => [c.id, c.title]));

  return (
    <div className="flex flex-col gap-6">
      <div className="flex flex-wrap items-start justify-between gap-4">
        <div>
          <h1 className="text-2xl font-semibold tracking-tight">Quizzes</h1>
          <p className="mt-1 text-sm text-muted-foreground">
            Every quiz you can author, across all courses. Open one to edit its questions.
          </p>
        </div>
        <div className="flex items-center gap-2">
          <ImportAssessmentDialog kind="quiz" />
          <CreateQuizDialog courses={courses} canCreateStandalone={isAdmin} />
        </div>
      </div>

      {quizzes.length === 0 ? (
        <EmptyState
          title="No quizzes yet"
          description="Create a quiz, or import a whole set of questions from a JSON or HTML file."
        />
      ) : (
        <div className="flex flex-col gap-3">
          {quizzes.map((quiz) => (
            <Link
              key={quiz.id}
              href={`/quizzes/${quiz.id}`}
              className="rounded-lg border p-4 transition-colors hover:bg-muted/40"
            >
              <div className="flex flex-wrap items-start justify-between gap-3">
                <div className="min-w-0">
                  <p className="text-sm font-medium">{quiz.title}</p>
                  <div className="mt-1 flex flex-wrap items-center gap-x-3 gap-y-1 text-xs text-muted-foreground">
                    {quiz.course_id ? (
                      <span>{courseTitles.get(quiz.course_id) ?? "Course"}</span>
                    ) : (
                      <span>Standalone</span>
                    )}
                    {quiz.time_limit_seconds != null && (
                      <span>{Math.round(quiz.time_limit_seconds / 60)}m limit</span>
                    )}
                    <span>{quiz.pass_percentage}% to pass</span>
                    {quiz.max_attempts != null && <span>Max {quiz.max_attempts} attempts</span>}
                    {quiz.randomize_questions && <span>Randomized</span>}
                  </div>
                </div>
                <div className="flex shrink-0 items-center gap-2">
                  {quiz.question_count != null && (
                    <Badge variant="secondary" className="text-xs">
                      {quiz.question_count} question{quiz.question_count === 1 ? "" : "s"}
                    </Badge>
                  )}
                  {/* The publish flag only governs standalone quizzes; a
                      course-attached quiz is gated by access to its course. */}
                  {!quiz.course_id && (
                    <Badge variant={quiz.is_published ? "default" : "outline"} className="text-xs">
                      {quiz.is_published ? "Published" : "Draft"}
                    </Badge>
                  )}
                </div>
              </div>
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}

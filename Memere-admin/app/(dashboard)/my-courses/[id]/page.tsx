import { notFound, redirect } from "next/navigation";
import { requireStaff } from "@/lib/auth/session";
import {
  getCourse, getCourseSections, listLessons, listQuizzes, listExams,
} from "@/lib/api/endpoints";
import { BreadcrumbSetter } from "@/lib/breadcrumb-context";
import { SectionsList } from "@/components/courses/sections-list";
import { QuizzesPanel } from "@/components/courses/quizzes-panel";
import { ExamsPanel } from "@/components/courses/exams-panel";
import { TeacherCourseActions } from "@/components/courses/teacher-course-actions";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent } from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { ApiError } from "@/lib/api/errors";
import type { Lesson } from "@/lib/api/schemas";

export default async function MyCourseDetailPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { user } = await requireStaff();
  if (user.role !== "teacher") redirect(user.role === "admin" ? "/courses" : "/");

  const { id } = await params;

  let course, sections, quizzes, exams;
  try {
    [course, sections, quizzes, exams] = await Promise.all([
      getCourse(id),
      getCourseSections(id),
      listQuizzes(id),
      listExams(id),
    ]);
  } catch (err) {
    if (err instanceof ApiError && err.status === 404) notFound();
    throw err;
  }

  // Fetch lessons for each section in parallel
  const lessonArrays: Lesson[][] = await Promise.all(
    sections.map((s) => listLessons(s.id).catch(() => []))
  );
  const lessonsBySectionId: Record<string, Lesson[]> = Object.fromEntries(
    sections.map((s, i) => [s.id, lessonArrays[i]])
  );

  return (
    <>
      <BreadcrumbSetter label={course.title} />
      <div className="flex flex-col gap-6">
        {/* Header */}
        <div className="flex items-start justify-between gap-4">
          <div className="flex flex-col gap-1">
            <h1 className="text-2xl font-semibold tracking-tight">{course.title}</h1>
            <div className="flex items-center gap-2 flex-wrap">
              <Badge variant={course.is_published ? "default" : "secondary"} className="text-xs">
                {course.is_published ? "Published" : "Draft"}
              </Badge>
              <span className="text-xs text-muted-foreground">{course.subject} · Grade {course.grade}</span>
              {course.level && (
                <span className="text-xs text-muted-foreground capitalize">{course.level}</span>
              )}
            </div>
          </div>
          <TeacherCourseActions course={course} />
        </div>

        {/* Metadata card */}
        <Card className="border shadow-none">
          <CardContent className="p-6 grid sm:grid-cols-2 lg:grid-cols-4 gap-6 text-sm">
            <div>
              <p className="text-muted-foreground mb-0.5">Price</p>
              <p className="font-medium">{course.is_free || course.price === 0 ? "Free" : `${course.currency} ${course.price}`}</p>
            </div>
            <div>
              <p className="text-muted-foreground mb-0.5">Language</p>
              <p className="font-medium uppercase">{course.language ?? "en"}</p>
            </div>
            <div>
              <p className="text-muted-foreground mb-0.5">Total Lessons</p>
              <p className="font-medium tabular-nums">{course.total_lessons ?? 0}</p>
            </div>
            <div>
              <p className="text-muted-foreground mb-0.5">Enrolled</p>
              <p className="font-medium tabular-nums">{course.enrollment_count ?? 0}</p>
            </div>
            <div className="sm:col-span-2 lg:col-span-4">
              <p className="text-muted-foreground mb-0.5">Description</p>
              <p className="leading-relaxed">{course.description}</p>
            </div>
          </CardContent>
        </Card>

        {/* Tabs: Sections & Lessons | Quizzes | Exams */}
        <Tabs defaultValue="content">
          <TabsList className="mb-2">
            <TabsTrigger value="content">
              Sections & Lessons
              <span className="ml-1.5 text-xs text-muted-foreground">({sections.length})</span>
            </TabsTrigger>
            <TabsTrigger value="quizzes">
              Quizzes
              <span className="ml-1.5 text-xs text-muted-foreground">({quizzes.length})</span>
            </TabsTrigger>
            <TabsTrigger value="exams">
              Exams
              <span className="ml-1.5 text-xs text-muted-foreground">({exams.length})</span>
            </TabsTrigger>
          </TabsList>

          <TabsContent value="content">
            <SectionsList
              courseId={id}
              sections={sections}
              lessonsBySectionId={lessonsBySectionId}
            />
          </TabsContent>

          <TabsContent value="quizzes">
            <QuizzesPanel courseId={id} quizzes={quizzes} />
          </TabsContent>

          <TabsContent value="exams">
            <ExamsPanel courseId={id} exams={exams} />
          </TabsContent>
        </Tabs>
      </div>
    </>
  );
}

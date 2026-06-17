import { notFound } from "next/navigation";
import { getCourse, getCourseSections } from "@/lib/api/endpoints";
import { ApiError } from "@/lib/api/errors";
import { Badge } from "@/components/ui/badge";
import { CourseActions } from "@/components/courses/course-actions";
import { BreadcrumbSetter } from "@/lib/breadcrumb-context";
import { formatDate, formatMoney } from "@/lib/format";

function Row({ label, children }: { label: string; children: React.ReactNode }) {
  return (
    <div className="flex flex-col sm:flex-row sm:items-start gap-1 py-3 border-b last:border-0">
      <dt className="w-44 shrink-0 text-sm text-muted-foreground">{label}</dt>
      <dd className="text-sm">{children}</dd>
    </div>
  );
}

export default async function CourseDetailPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = await params;

  let course;
  let sections: Awaited<ReturnType<typeof getCourseSections>> = [];

  try {
    [course, sections] = await Promise.all([
      getCourse(id),
      getCourseSections(id).catch(() => []),
    ]);
  } catch (err) {
    if (
      err instanceof ApiError &&
      (err.status === 404 ||
        err.code === "NOT_FOUND" ||
        err.code === "RESOURCE_NOT_FOUND")
    ) {
      notFound();
    }
    throw err;
  }

  const price =
    course.is_free || course.price === 0
      ? "Free"
      : formatMoney(String(course.price), course.currency);

  return (
    <>
      <BreadcrumbSetter label={course.title} />

      <div className="flex flex-col gap-6 max-w-2xl">
        <div className="flex items-start justify-between gap-4 flex-wrap">
          <div>
            <h1 className="text-2xl font-semibold tracking-tight">{course.title}</h1>
            {course.short_description && (
              <p className="text-sm text-muted-foreground mt-0.5">
                {course.short_description}
              </p>
            )}
          </div>
          <CourseActions course={course} />
        </div>

        {/* Metadata */}
        <div className="rounded-lg border">
          <dl>
            <Row label="Status">
              {course.is_published ? (
                <Badge
                  variant="secondary"
                  className="text-green-700 bg-green-100 dark:text-green-400 dark:bg-green-950"
                >
                  Published
                </Badge>
              ) : (
                <Badge variant="outline">Unpublished</Badge>
              )}
            </Row>
            <Row label="Subject">
              <span className="capitalize">{course.subject}</span>
            </Row>
            <Row label="Grade">Grade {course.grade}</Row>
            <Row label="Price">{price}</Row>
            {course.level && (
              <Row label="Level">
                <span className="capitalize">{course.level}</span>
              </Row>
            )}
            {course.language && (
              <Row label="Language">
                <span className="uppercase">{course.language}</span>
              </Row>
            )}
            {course.enrollment_count !== undefined && (
              <Row label="Enrollments">
                {course.enrollment_count.toLocaleString()}
              </Row>
            )}
            {course.total_lessons !== undefined && (
              <Row label="Lessons">{course.total_lessons}</Row>
            )}
            <Row label="Created">{formatDate(course.created_at)}</Row>
          </dl>
        </div>

        {/* Description */}
        {course.description && (
          <div className="flex flex-col gap-1">
            <h2 className="text-sm font-medium text-muted-foreground">Description</h2>
            <p className="text-sm leading-relaxed">{course.description}</p>
          </div>
        )}

        {/* Sections */}
        <div className="flex flex-col gap-2">
          <h2 className="text-base font-semibold">
            Sections{" "}
            <span className="text-sm font-normal text-muted-foreground">
              ({sections.length})
            </span>
          </h2>
          {sections.length === 0 ? (
            <p className="text-sm text-muted-foreground">No sections yet.</p>
          ) : (
            <ol className="flex flex-col gap-1">
              {sections.map((section, i) => (
                <li
                  key={section.id}
                  className="flex items-center gap-3 rounded-md border px-4 py-2.5 text-sm"
                >
                  <span className="tabular-nums text-muted-foreground w-5">
                    {i + 1}.
                  </span>
                  <span className="flex-1">{section.title}</span>
                  {section.is_published === false && (
                    <Badge variant="outline" className="text-xs">
                      Draft
                    </Badge>
                  )}
                  {section.total_lessons !== undefined && (
                    <span className="text-xs text-muted-foreground">
                      {section.total_lessons} lesson
                      {section.total_lessons !== 1 ? "s" : ""}
                    </span>
                  )}
                </li>
              ))}
            </ol>
          )}
        </div>
      </div>
    </>
  );
}

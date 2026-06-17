import { listCourses } from "@/lib/api/endpoints";
import { CoursesClient } from "./courses-client";

export default async function CoursesPage({
  searchParams,
}: {
  searchParams: Promise<{ after?: string; limit?: string }>;
}) {
  const params = await searchParams;
  const limit = Math.min(parseInt(params.limit ?? "20", 10), 100);
  const after = params.after;

  const result = await listCourses({ limit, after });

  return (
    <div className="flex flex-col gap-6">
      <div>
        <h1 className="text-2xl font-semibold tracking-tight">Courses</h1>
        <p className="text-sm text-muted-foreground mt-1">
          Review and moderate published courses.
        </p>
      </div>
      <CoursesClient initialData={{ items: result.courses, next: result.next }} />
    </div>
  );
}

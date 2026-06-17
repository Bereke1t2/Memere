"use client";

import { useRouter } from "next/navigation";
import { DataTable } from "@/components/data-table/data-table";
import { Pagination } from "@/components/data-table/pagination";
import { useCursorList } from "@/lib/hooks/use-cursor-list";
import { CreateCourseDialog } from "@/components/courses/create-course-dialog";
import { columns } from "./columns";
import type { Course } from "@/lib/api/schemas";
import type { CursorPage } from "@/lib/hooks/use-cursor-list";

interface MyCoursesClientProps {
  initialData: CursorPage<Course>;
}

export function MyCoursesClient({ initialData }: MyCoursesClientProps) {
  const router = useRouter();
  const { data, isLoading } = useCursorList<Course>({
    route: "/api/teacher/courses",
    queryKey: ["teacher-courses"],
    initialData,
  });

  return (
    <div className="flex flex-col gap-4">
      <div className="flex justify-end">
        <CreateCourseDialog />
      </div>
      <DataTable
        columns={columns}
        data={data?.items ?? []}
        isLoading={isLoading}
        onRowClick={(row) => router.push(`/my-courses/${row.original.id}`)}
        emptyState={
          <div className="py-12 text-center text-sm text-muted-foreground">
            No courses yet. Create your first course to get started.
          </div>
        }
      />
      <Pagination nextCursor={data?.next ?? undefined} isLoading={isLoading} />
    </div>
  );
}

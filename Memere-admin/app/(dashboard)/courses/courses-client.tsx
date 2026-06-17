"use client";

import { useRouter } from "next/navigation";
import { DataTable } from "@/components/data-table/data-table";
import { Pagination } from "@/components/data-table/pagination";
import { useCursorList, type CursorPage } from "@/lib/hooks/use-cursor-list";
import { courseColumns } from "./columns";
import type { Course } from "@/lib/api/schemas";

interface CoursesClientProps {
  initialData: CursorPage<Course>;
}

export function CoursesClient({ initialData }: CoursesClientProps) {
  const router = useRouter();

  const { data, isFetching } = useCursorList<Course>({
    route: "/api/admin/courses",
    queryKey: ["courses"],
    initialData,
  });

  const items = data?.items ?? [];
  const next = data?.next ?? "";

  return (
    <div className="flex flex-col gap-4">
      <div className="flex items-center justify-end">
        {data && (
          <p className="text-xs text-muted-foreground">
            Showing {items.length} course{items.length !== 1 ? "s" : ""}
          </p>
        )}
      </div>

      <DataTable
        columns={courseColumns}
        data={items}
        isLoading={isFetching && items.length === 0}
        onRowClick={(row) => router.push(`/courses/${row.original.id}`)}
        emptyState={
          <p className="text-sm text-muted-foreground">No courses found.</p>
        }
      />

      <Pagination nextCursor={next} isLoading={isFetching} />
    </div>
  );
}

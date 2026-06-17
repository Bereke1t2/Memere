"use client";

import { useRouter } from "next/navigation";
import { DataTable } from "@/components/data-table/data-table";
import { Pagination } from "@/components/data-table/pagination";
import { SearchFilter, SelectFilter } from "@/components/data-table/filters";
import { useCursorList, type CursorPage } from "@/lib/hooks/use-cursor-list";
import { userColumns } from "./columns";
import type { User } from "@/lib/api/schemas";

const ROLE_OPTIONS = [
  { label: "Student", value: "student" },
  { label: "Teacher", value: "teacher" },
  { label: "Admin", value: "admin" },
];

interface UsersClientProps {
  initialData: CursorPage<User>;
}

export function UsersClient({ initialData }: UsersClientProps) {
  const router = useRouter();

  const { data, isFetching } = useCursorList<User>({
    route: "/api/admin/users",
    queryKey: ["users"],
    initialData,
  });

  const items = data?.items ?? [];
  const next = data?.next ?? "";

  return (
    <div className="flex flex-col gap-4">
      <div className="flex items-center gap-2 flex-wrap">
        <SearchFilter placeholder="Search by name or email…" paramKey="q" />
        <SelectFilter options={ROLE_OPTIONS} placeholder="All roles" paramKey="role" />
        {data && (
          <p className="ml-auto text-xs text-muted-foreground">
            Showing {items.length} user{items.length !== 1 ? "s" : ""}
          </p>
        )}
      </div>

      <DataTable
        columns={userColumns}
        data={items}
        isLoading={isFetching && items.length === 0}
        onRowClick={(row) => router.push(`/users/${row.original.id}`)}
        emptyState={
          <p className="text-sm text-muted-foreground">
            No users match these filters.
          </p>
        }
      />

      <Pagination nextCursor={next} isLoading={isFetching} />
    </div>
  );
}

"use client";

import { useRouter } from "next/navigation";
import { DataTable } from "@/components/data-table/data-table";
import { Pagination } from "@/components/data-table/pagination";
import { SelectFilter } from "@/components/data-table/filters";
import { useCursorList, type CursorPage } from "@/lib/hooks/use-cursor-list";
import { ReconcileButton } from "@/components/payments/reconcile-button";
import { paymentColumns } from "./columns";
import type { AdminPayment } from "@/lib/api/schemas";

const STATUS_OPTIONS = [
  { label: "Pending", value: "pending" },
  { label: "Completed", value: "completed" },
  { label: "Failed", value: "failed" },
  { label: "Refunded", value: "refunded" },
];

interface PaymentsClientProps {
  initialData: CursorPage<AdminPayment>;
}

export function PaymentsClient({ initialData }: PaymentsClientProps) {
  const router = useRouter();

  const { data, isFetching } = useCursorList<AdminPayment>({
    route: "/api/admin/payments",
    queryKey: ["payments"],
    initialData,
  });

  const items = data?.items ?? [];
  const next = data?.next ?? "";

  return (
    <div className="flex flex-col gap-4">
      <div className="flex items-center gap-2 flex-wrap">
        <SelectFilter
          options={STATUS_OPTIONS}
          placeholder="All statuses"
          paramKey="status"
        />
        <div className="ml-auto flex items-center gap-3">
          {data && (
            <p className="text-xs text-muted-foreground">
              Showing {items.length} payment{items.length !== 1 ? "s" : ""}
            </p>
          )}
          <ReconcileButton />
        </div>
      </div>

      <DataTable
        columns={paymentColumns}
        data={items}
        isLoading={isFetching && items.length === 0}
        onRowClick={(row) => router.push(`/payments/${row.original.id}`)}
        emptyState={
          <p className="text-sm text-muted-foreground">
            No payments match these filters.
          </p>
        }
      />

      <Pagination nextCursor={next} isLoading={isFetching} />
    </div>
  );
}

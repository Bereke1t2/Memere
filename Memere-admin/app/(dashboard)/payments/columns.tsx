import { type ColumnDef } from "@tanstack/react-table";
import Link from "next/link";
import { Badge } from "@/components/ui/badge";
import { formatDate, formatMoney } from "@/lib/format";
import type { AdminPayment } from "@/lib/api/schemas";

function StatusBadge({ status }: { status: string }) {
  switch (status) {
    case "completed":
      return (
        <Badge variant="secondary" className="text-green-700 bg-green-100 dark:text-green-400 dark:bg-green-950">
          Completed
        </Badge>
      );
    case "pending":
      return (
        <Badge variant="secondary" className="text-amber-700 bg-amber-100 dark:text-amber-400 dark:bg-amber-950">
          Pending
        </Badge>
      );
    case "failed":
      return <Badge variant="destructive">Failed</Badge>;
    case "refunded":
      return <Badge variant="secondary">Refunded</Badge>;
    default:
      return (
        <Badge variant="outline" className="capitalize">
          {status}
        </Badge>
      );
  }
}

function CopyableId({ id }: { id: string }) {
  return (
    <button
      className="font-mono text-xs text-muted-foreground hover:text-foreground transition-colors"
      title={id}
      onClick={(e) => {
        e.stopPropagation();
        navigator.clipboard.writeText(id).catch(() => {});
      }}
    >
      {id.slice(0, 8)}&hellip;
    </button>
  );
}

export const paymentColumns: ColumnDef<AdminPayment>[] = [
  {
    accessorKey: "id",
    header: "Payment ID",
    cell: ({ getValue }) => <CopyableId id={getValue<string>()} />,
  },
  {
    accessorKey: "student_id",
    header: "Student",
    cell: ({ getValue }) => {
      const id = getValue<string>();
      return (
        <Link
          href={`/users/${id}`}
          className="font-mono text-xs text-muted-foreground hover:text-foreground transition-colors"
          onClick={(e) => e.stopPropagation()}
          title={id}
        >
          {id.slice(0, 8)}&hellip;
        </Link>
      );
    },
  },
  {
    id: "amount",
    header: () => <span className="block text-right">Amount</span>,
    cell: ({ row }) => (
      <span className="block text-right tabular-nums font-medium">
        {formatMoney(row.original.amount, row.original.currency)}
      </span>
    ),
  },
  {
    accessorKey: "status",
    header: "Status",
    cell: ({ getValue }) => <StatusBadge status={getValue<string>()} />,
  },
  {
    accessorKey: "provider",
    header: "Provider",
    cell: ({ getValue }) => (
      <span className="text-sm capitalize">{getValue<string>()}</span>
    ),
  },
  {
    accessorKey: "created_at",
    header: "Date",
    cell: ({ getValue }) => (
      <span className="text-sm tabular-nums text-muted-foreground">
        {formatDate(getValue<string>())}
      </span>
    ),
  },
];

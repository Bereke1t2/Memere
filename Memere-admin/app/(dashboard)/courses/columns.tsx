import { type ColumnDef } from "@tanstack/react-table";
import { Badge } from "@/components/ui/badge";
import { formatDate, formatMoney } from "@/lib/format";
import type { Course } from "@/lib/api/schemas";

function priceLabel(course: Course): string {
  if (course.is_free || course.price === 0) return "Free";
  return formatMoney(String(course.price), course.currency);
}

export const courseColumns: ColumnDef<Course>[] = [
  {
    accessorKey: "title",
    header: "Title",
    cell: ({ getValue }) => (
      <span className="font-medium">{getValue<string>()}</span>
    ),
  },
  {
    accessorKey: "subject",
    header: "Subject",
    cell: ({ getValue }) => (
      <span className="text-sm capitalize">{getValue<string>()}</span>
    ),
  },
  {
    accessorKey: "grade",
    header: "Grade",
    cell: ({ getValue }) => (
      <span className="text-sm tabular-nums">Grade {getValue<number>()}</span>
    ),
  },
  {
    id: "price",
    header: "Price",
    cell: ({ row }) => (
      <span className="text-sm tabular-nums">{priceLabel(row.original)}</span>
    ),
  },
  {
    accessorKey: "is_published",
    header: "Status",
    cell: ({ getValue }) =>
      getValue<boolean>() ? (
        <Badge
          variant="secondary"
          className="text-green-700 bg-green-100 dark:text-green-400 dark:bg-green-950"
        >
          Published
        </Badge>
      ) : (
        <Badge variant="outline">Unpublished</Badge>
      ),
  },
  {
    accessorKey: "created_at",
    header: "Created",
    cell: ({ getValue }) => (
      <span className="text-sm tabular-nums text-muted-foreground">
        {formatDate(getValue<string>())}
      </span>
    ),
  },
];

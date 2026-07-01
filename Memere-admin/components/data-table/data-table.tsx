"use client";

import {
  useReactTable,
  getCoreRowModel,
  flexRender,
  type ColumnDef,
  type Row,
} from "@tanstack/react-table";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Skeleton } from "@/components/ui/skeleton";
import { cn } from "@/lib/utils";

interface DataTableProps<TData> {
  columns: ColumnDef<TData>[];
  data: TData[];
  isLoading?: boolean;
  emptyState?: React.ReactNode;
  onRowClick?: (row: Row<TData>) => void;
  skeletonRows?: number;
}

export function DataTable<TData>({
  columns,
  data,
  isLoading,
  emptyState,
  onRowClick,
  skeletonRows = 10,
}: DataTableProps<TData>) {
  // eslint-disable-next-line react-hooks/incompatible-library
  const table = useReactTable({
    data,
    columns,
    getCoreRowModel: getCoreRowModel(),
  });

  return (
    <div className="overflow-hidden rounded-lg border bg-card shadow-sm shadow-emerald-950/5">
      <div className="overflow-auto">
        <Table>
          <TableHeader className="sticky top-0 z-10">
            {table.getHeaderGroups().map((headerGroup) => (
              <TableRow key={headerGroup.id} className="border-b bg-muted/60 hover:bg-muted/60">
                {headerGroup.headers.map((header) => (
                  <TableHead key={header.id} className="py-3 text-xs font-semibold text-muted-foreground">
                    {header.isPlaceholder
                      ? null
                      : flexRender(
                          header.column.columnDef.header,
                          header.getContext()
                        )}
                  </TableHead>
                ))}
              </TableRow>
            ))}
          </TableHeader>
          <TableBody>
            {isLoading ? (
              Array.from({ length: skeletonRows }).map((_, i) => (
                <TableRow key={i}>
                  {columns.map((_, j) => (
                    <TableCell key={j}>
                      <Skeleton className="h-4 w-full" />
                    </TableCell>
                  ))}
                </TableRow>
              ))
            ) : table.getRowModel().rows.length === 0 ? (
              <TableRow>
                <TableCell
                  colSpan={columns.length}
                  className="py-16 text-center"
                >
                  {emptyState ?? (
                    <p className="text-sm text-muted-foreground">
                      No results found.
                    </p>
                  )}
                </TableCell>
              </TableRow>
            ) : (
              table.getRowModel().rows.map((row) => (
                <TableRow
                  key={row.id}
                  onClick={() => onRowClick?.(row)}
                  className={cn(
                    "hover:bg-secondary/55 transition-colors",
                    onRowClick && "cursor-pointer"
                  )}
                >
                  {row.getVisibleCells().map((cell) => (
                    <TableCell key={cell.id}>
                      {flexRender(
                        cell.column.columnDef.cell,
                        cell.getContext()
                      )}
                    </TableCell>
                  ))}
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>
    </div>
  );
}

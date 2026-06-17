"use client";

import { type ColumnDef } from "@tanstack/react-table";
import { Badge } from "@/components/ui/badge";
import { TeacherCourseActions } from "@/components/courses/teacher-course-actions";
import type { Course } from "@/lib/api/schemas";

export const columns: ColumnDef<Course>[] = [
  {
    accessorKey: "title",
    header: "Title",
    cell: ({ row }) => (
      <span className="font-medium line-clamp-1 max-w-[260px]">{row.original.title}</span>
    ),
  },
  {
    accessorKey: "subject",
    header: "Subject",
    cell: ({ row }) => <span className="text-sm">{row.original.subject}</span>,
  },
  {
    accessorKey: "grade",
    header: "Grade",
    cell: ({ row }) => <span className="text-sm">Grade {row.original.grade}</span>,
  },
  {
    accessorKey: "price",
    header: "Price",
    cell: ({ row }) => (
      <span className="text-sm tabular-nums">
        {row.original.is_free || row.original.price === 0
          ? "Free"
          : `${row.original.currency} ${row.original.price}`}
      </span>
    ),
  },
  {
    accessorKey: "is_published",
    header: "Status",
    cell: ({ row }) => (
      <Badge variant={row.original.is_published ? "default" : "secondary"} className="text-xs">
        {row.original.is_published ? "Published" : "Draft"}
      </Badge>
    ),
  },
  {
    accessorKey: "total_lessons",
    header: "Lessons",
    cell: ({ row }) => (
      <span className="text-sm text-muted-foreground tabular-nums">{row.original.total_lessons ?? 0}</span>
    ),
  },
  {
    id: "actions",
    header: "",
    cell: ({ row }) => <TeacherCourseActions course={row.original} />,
  },
];

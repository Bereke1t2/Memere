import { type ColumnDef } from "@tanstack/react-table";
import { Badge } from "@/components/ui/badge";
import { formatDate } from "@/lib/format";
import type { User } from "@/lib/api/schemas";

function RoleBadge({ role }: { role: string }) {
  const variant =
    role === "admin"
      ? "default"
      : role === "teacher"
        ? "secondary"
        : "outline";
  return (
    <Badge variant={variant} className="capitalize">
      {role}
    </Badge>
  );
}

function StatusBadge({ active }: { active: boolean }) {
  return active ? (
    <Badge variant="secondary" className="text-green-700 bg-green-100 dark:text-green-400 dark:bg-green-950">
      Active
    </Badge>
  ) : (
    <Badge variant="destructive">Suspended</Badge>
  );
}

export const userColumns: ColumnDef<User>[] = [
  {
    id: "name",
    header: "Name",
    cell: ({ row }) => (
      <span className="font-medium whitespace-nowrap">
        {row.original.first_name} {row.original.last_name}
      </span>
    ),
  },
  {
    accessorKey: "email",
    header: "Email",
    cell: ({ getValue }) => (
      <span className="text-sm text-muted-foreground">{getValue<string>()}</span>
    ),
  },
  {
    accessorKey: "role",
    header: "Role",
    cell: ({ getValue }) => <RoleBadge role={getValue<string>()} />,
  },
  {
    id: "status",
    header: "Status",
    cell: ({ row }) => <StatusBadge active={row.original.is_active} />,
  },
  {
    accessorKey: "created_at",
    header: "Joined",
    cell: ({ getValue }) => (
      <span className="text-sm tabular-nums">{formatDate(getValue<string>())}</span>
    ),
  },
  {
    accessorKey: "last_login_at",
    header: "Last login",
    cell: ({ getValue }) => (
      <span className="text-sm tabular-nums text-muted-foreground">
        {formatDate(getValue<string | null>())}
      </span>
    ),
  },
];

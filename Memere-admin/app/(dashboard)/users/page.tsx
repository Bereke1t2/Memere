import { listUsers } from "@/lib/api/endpoints";
import { UsersClient } from "./users-client";

export default async function UsersPage({
  searchParams,
}: {
  searchParams: Promise<{ after?: string; role?: string; q?: string; limit?: string }>;
}) {
  const params = await searchParams;
  const limit = Math.min(parseInt(params.limit ?? "20", 10), 100);
  const after = params.after;
  const role = params.role;

  const result = await listUsers({ limit, after, role });

  const q = (params.q ?? "").toLowerCase().trim();
  const items = q
    ? result.users.filter(
        (u) =>
          u.email.toLowerCase().includes(q) ||
          u.first_name.toLowerCase().includes(q) ||
          u.last_name.toLowerCase().includes(q)
      )
    : result.users;

  return (
    <div className="flex flex-col gap-6">
      <div>
        <h1 className="text-2xl font-semibold tracking-tight">Users</h1>
        <p className="text-sm text-muted-foreground mt-1">
          Manage platform users. Name/email search applies to the current page only.
        </p>
      </div>
      <UsersClient initialData={{ items, next: result.next }} />
    </div>
  );
}

import { listUsers } from "@/lib/api/endpoints";
import { ApiError } from "@/lib/api/errors";
import { friendlyMessage } from "@/lib/api/errors";

export const dynamic = "force-dynamic";

export async function GET(req: Request) {
  const { searchParams } = new URL(req.url);
  const limit = Math.min(parseInt(searchParams.get("limit") ?? "20", 10), 100);
  const after = searchParams.get("after") ?? undefined;
  const role = searchParams.get("role") ?? undefined;
  const q = (searchParams.get("q") ?? "").toLowerCase().trim();

  try {
    const result = await listUsers({ limit, after, role });

    // Backend has no full-text search — filter current page only.
    const items = q
      ? result.users.filter(
          (u) =>
            u.email.toLowerCase().includes(q) ||
            u.first_name.toLowerCase().includes(q) ||
            u.last_name.toLowerCase().includes(q)
        )
      : result.users;

    return Response.json({ items, next: result.next });
  } catch (err) {
    if (err instanceof ApiError) {
      return Response.json(
        { message: friendlyMessage(err) },
        { status: err.status }
      );
    }
    return Response.json({ message: "Failed to load users." }, { status: 500 });
  }
}

import { changeRole } from "@/lib/api/endpoints";
import { ApiError, friendlyMessage } from "@/lib/api/errors";

const VALID_ROLES = new Set(["student", "teacher", "admin"]);

export async function POST(
  req: Request,
  { params }: { params: Promise<{ id: string }> }
) {
  const { id } = await params;
  const body = await req.json().catch(() => ({}));
  const role = (body.role ?? "").trim();

  if (!VALID_ROLES.has(role)) {
    return Response.json({ message: "Invalid role." }, { status: 400 });
  }

  try {
    await changeRole(id, role);
    return new Response(null, { status: 204 });
  } catch (err) {
    if (err instanceof ApiError) {
      return Response.json({ message: friendlyMessage(err) }, { status: err.status });
    }
    return Response.json({ message: "Failed to change role." }, { status: 500 });
  }
}

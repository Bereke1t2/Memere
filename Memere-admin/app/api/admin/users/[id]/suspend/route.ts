import { suspendUser } from "@/lib/api/endpoints";
import { ApiError, friendlyMessage } from "@/lib/api/errors";

export async function POST(
  req: Request,
  { params }: { params: Promise<{ id: string }> }
) {
  const { id } = await params;
  const body = await req.json().catch(() => ({}));
  const reason = (body.reason ?? "").trim();

  if (!reason) {
    return Response.json({ message: "A reason is required." }, { status: 400 });
  }

  try {
    await suspendUser(id, reason);
    return new Response(null, { status: 204 });
  } catch (err) {
    if (err instanceof ApiError) {
      return Response.json({ message: friendlyMessage(err) }, { status: err.status });
    }
    return Response.json({ message: "Failed to suspend user." }, { status: 500 });
  }
}

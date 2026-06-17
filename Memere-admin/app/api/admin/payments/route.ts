import { listPayments } from "@/lib/api/endpoints";
import { ApiError, friendlyMessage } from "@/lib/api/errors";

export const dynamic = "force-dynamic";

export async function GET(req: Request) {
  const { searchParams } = new URL(req.url);
  const limit = Math.min(parseInt(searchParams.get("limit") ?? "20", 10), 100);
  const after = searchParams.get("after") ?? undefined;
  const status = searchParams.get("status") ?? undefined;

  try {
    const result = await listPayments({ limit, after, status });
    return Response.json({ items: result.payments, next: result.next });
  } catch (err) {
    if (err instanceof ApiError) {
      return Response.json({ message: friendlyMessage(err) }, { status: err.status });
    }
    return Response.json({ message: "Failed to load payments." }, { status: 500 });
  }
}

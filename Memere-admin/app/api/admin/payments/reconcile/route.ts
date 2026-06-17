import { reconcilePayments } from "@/lib/api/endpoints";
import { ApiError, friendlyMessage } from "@/lib/api/errors";

export async function POST() {
  try {
    const result = await reconcilePayments();
    return Response.json(result);
  } catch (err) {
    if (err instanceof ApiError) {
      return Response.json({ message: friendlyMessage(err) }, { status: err.status });
    }
    return Response.json({ message: "Reconciliation failed." }, { status: 500 });
  }
}

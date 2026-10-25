import { env } from "@/lib/env";
import { getRouteAdminSession } from "@/lib/auth/session";

export const dynamic = "force-dynamic";

export async function GET() {
  const session = await getRouteAdminSession();
  if (!session) {
    return Response.json({ message: "Unauthorized" }, { status: 401 });
  }

  try {
    const res = await fetch(`${env.API_BASE_URL}/health/system`, {
      cache: "no-store",
    });

    if (!res.ok) {
      return Response.json(
        { message: `Backend returned status ${res.status}` },
        { status: res.status }
      );
    }

    const data = await res.json();
    return Response.json(data);
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : "Failed to fetch system metrics";
    return Response.json(
      { message },
      { status: 502 }
    );
  }
}

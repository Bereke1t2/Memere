import { getRouteStaffSession } from "@/lib/auth/session";
import { canManageContent } from "@/lib/auth/roles";
import { env } from "@/lib/env";
import { getAccessToken } from "@/lib/auth/cookies";
import { refreshAccessToken } from "@/lib/auth/refresh";
import { ApiError, friendlyMessage } from "@/lib/api/errors";

export const dynamic = "force-dynamic";

export async function POST(req: Request, { params }: { params: Promise<{ id: string }> }) {
  try {
    const session = await getRouteStaffSession();
    if (!session) return Response.json({ message: "Unauthorized", code: "UNAUTHORIZED" }, { status: 401 });
    if (!canManageContent(session.user)) return Response.json({ message: "Forbidden" }, { status: 403 });

    const { id } = await params;
    const formData = await req.formData();
    const file = formData.get("file") as File | null;
    if (!file) {
      return Response.json({ message: "No image file provided." }, { status: 400 });
    }

    const url = `${env.API_BASE_URL}/api/v1/courses/${id}/thumbnail`;

    async function forward(bearer: string | undefined) {
      const out = new FormData();
      out.append("file", file as File, (file as File).name);
      const headers: HeadersInit = {};
      if (bearer) headers["Authorization"] = `Bearer ${bearer}`;
      return fetch(url, { method: "POST", headers, body: out, cache: "no-store" });
    }

    let res = await forward(await getAccessToken());
    if (res.status === 401) {
      const fresh = await refreshAccessToken();
      if (fresh) res = await forward(fresh);
    }

    const json = await res.json().catch(() => null);
    if (!res.ok) {
      const message =
        (json && typeof json === "object" && "message" in json && (json as { message?: string }).message) ||
        "Thumbnail upload failed.";
      return Response.json({ message }, { status: res.status });
    }
    return Response.json(json);
  } catch (err) {
    if (err instanceof ApiError) return Response.json({ message: friendlyMessage(err) }, { status: err.status });
    return Response.json({ message: err instanceof Error ? err.message : "Failed to upload thumbnail." }, { status: 500 });
  }
}

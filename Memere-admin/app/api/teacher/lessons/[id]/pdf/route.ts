import { requireStaff } from "@/lib/auth/session";
import { canManageContent } from "@/lib/auth/roles";
import { env } from "@/lib/env";
import { getAccessToken } from "@/lib/auth/cookies";
import { refreshAccessToken } from "@/lib/auth/refresh";
import { ApiError, friendlyMessage } from "@/lib/api/errors";

export const dynamic = "force-dynamic";

// Proxies a multipart PDF upload to the backend POST /lessons/:id/pdf, which
// validates the %PDF- header, stores the bytes in object storage (Backblaze B2)
// at lessons/<id>/notes.pdf, and sets the lesson's pdf_url to that storage key.
// apiFetch can't be used here — it forces Content-Type: application/json — so we
// forward the multipart body directly, carrying the bearer token with a single
// 401 -> refresh -> retry (mirroring apiFetch's refresh behavior).
export async function POST(req: Request, { params }: { params: Promise<{ id: string }> }) {
  try {
    const { user } = await requireStaff();
    if (!canManageContent(user)) return Response.json({ message: "Forbidden" }, { status: 403 });

    const { id } = await params;
    const formData = await req.formData();
    const file = formData.get("file") as File | null;
    if (!file) {
      return Response.json({ message: "No PDF file provided." }, { status: 400 });
    }

    const url = `${env.API_BASE_URL}/api/v1/lessons/${id}/pdf`;

    // Fresh FormData per attempt so the File (a Blob, re-readable) can be sent
    // again on the refresh retry. Do NOT set Content-Type — fetch sets the
    // multipart boundary itself.
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
        "PDF upload failed.";
      return Response.json({ message }, { status: res.status });
    }
    // Backend returns the updated lesson (pdf_url now = the storage key).
    return Response.json(json);
  } catch (err) {
    if (err instanceof ApiError) return Response.json({ message: friendlyMessage(err) }, { status: err.status });
    return Response.json({ message: err instanceof Error ? err.message : "Failed to upload PDF." }, { status: 500 });
  }
}

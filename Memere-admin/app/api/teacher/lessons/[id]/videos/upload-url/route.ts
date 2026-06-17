import { z } from "zod";
import { requireStaff } from "@/lib/auth/session";
import { requestVideoUpload } from "@/lib/api/endpoints";
import { ApiError, friendlyMessage } from "@/lib/api/errors";

const Schema = z.object({
  file_name: z.string().min(1),
  content_type: z.string().min(1),
  size_bytes: z.number().positive(),
});

export async function POST(req: Request, { params }: { params: Promise<{ id: string }> }) {
  try {
    const { user } = await requireStaff();
    if (user.role !== "teacher") return Response.json({ message: "Forbidden" }, { status: 403 });
    const { id } = await params;
    const body = await req.json().catch(() => ({}));
    const parsed = Schema.safeParse(body);
    if (!parsed.success) return Response.json({ message: "file_name, content_type, and size_bytes are required." }, { status: 400 });
    return Response.json(await requestVideoUpload(id, parsed.data));
  } catch (err) {
    if (err instanceof ApiError) return Response.json({ message: friendlyMessage(err) }, { status: err.status });
    return Response.json({ message: "Failed to get upload URL." }, { status: 500 });
  }
}

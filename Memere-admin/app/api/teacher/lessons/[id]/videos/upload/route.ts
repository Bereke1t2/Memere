import { requireStaff } from "@/lib/auth/session";
import { canManageContent } from "@/lib/auth/roles";
import { requestVideoUpload, confirmVideoUpload } from "@/lib/api/endpoints";
import { ApiError, friendlyMessage } from "@/lib/api/errors";

export const dynamic = "force-dynamic";

export async function POST(req: Request, { params }: { params: Promise<{ id: string }> }) {
  try {
    const { user } = await requireStaff();
    if (!canManageContent(user)) return Response.json({ message: "Forbidden" }, { status: 403 });

    const { id: lessonId } = await params;
    const formData = await req.formData();
    const file = formData.get("file") as File | null;
    if (!file) {
      return Response.json({ message: "No video file provided." }, { status: 400 });
    }

    // 1. Get presigned upload URL from backend
    const { upload_url, video_id } = await requestVideoUpload(lessonId, {
      file_name: file.name,
      content_type: file.type || "video/mp4",
      size_bytes: file.size,
    });

    // 2. Stream/Upload file bytes directly to MinIO/S3 from Node server
    const bytes = await file.arrayBuffer();
    const uploadRes = await fetch(upload_url, {
      method: "PUT",
      headers: { "Content-Type": file.type || "video/mp4" },
      body: bytes,
    });

    if (!uploadRes.ok) {
      throw new Error(`Storage upload failed with status ${uploadRes.status}`);
    }

    // 3. Confirm upload with backend
    await confirmVideoUpload(video_id);

    return Response.json({ status: "ok", video_id, file_name: file.name });
  } catch (err) {
    if (err instanceof ApiError) return Response.json({ message: friendlyMessage(err) }, { status: err.status });
    return Response.json({ message: err instanceof Error ? err.message : "Failed to upload video." }, { status: 500 });
  }
}

import { broadcast, listAnnouncements, deleteAnnouncement } from "@/lib/api/endpoints";
import { ApiError, friendlyMessage } from "@/lib/api/errors";
import { z } from "zod";

export async function GET(req: Request) {
  try {
    const { searchParams } = new URL(req.url);
    const limit = Number(searchParams.get("limit")) || 50;
    const data = await listAnnouncements(limit);
    return Response.json(data);
  } catch (err) {
    if (err instanceof ApiError) {
      return Response.json({ message: friendlyMessage(err) }, { status: err.status });
    }
    return Response.json({ message: "Failed to fetch announcements." }, { status: 500 });
  }
}

const BodySchema = z.object({
  title: z.string().min(1),
  body: z.string().min(1),
  segment: z.enum(["all", "students", "teachers", "subscribers"]),
  data: z.record(z.string(), z.string()).optional(),
});

export async function POST(req: Request) {
  const raw = await req.json().catch(() => null);
  const parsed = BodySchema.safeParse(raw);

  if (!parsed.success) {
    return Response.json({ message: "Invalid request body." }, { status: 400 });
  }

  try {
    await broadcast(parsed.data);
    return new Response(null, { status: 204 });
  } catch (err) {
    if (err instanceof ApiError) {
      return Response.json({ message: friendlyMessage(err) }, { status: err.status });
    }
    return Response.json({ message: "Failed to send announcement." }, { status: 500 });
  }
}

export async function DELETE(req: Request) {
  const { searchParams } = new URL(req.url);
  const id = searchParams.get("id");
  if (!id) {
    return Response.json({ message: "Announcement ID is required." }, { status: 400 });
  }

  try {
    await deleteAnnouncement(id);
    return new Response(null, { status: 204 });
  } catch (err) {
    if (err instanceof ApiError) {
      return Response.json({ message: friendlyMessage(err) }, { status: err.status });
    }
    return Response.json({ message: "Failed to delete announcement." }, { status: 500 });
  }
}


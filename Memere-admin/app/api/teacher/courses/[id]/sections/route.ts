import { z } from "zod";
import { requireStaff } from "@/lib/auth/session";
import { getCourseSections } from "@/lib/api/endpoints";
import { apiFetch } from "@/lib/api/server";
import { SectionSchema } from "@/lib/api/schemas";
import { ApiError, friendlyMessage } from "@/lib/api/errors";

export const dynamic = "force-dynamic";

const AddSectionSchema = z.object({
  title: z.string().min(1, "Title is required"),
  order: z.number().int().optional(),
});

export async function GET(
  _req: Request,
  { params }: { params: Promise<{ id: string }> }
) {
  try {
    await requireStaff();
    const { id } = await params;
    const sections = await getCourseSections(id);
    return Response.json(sections);
  } catch (err) {
    if (err instanceof ApiError) {
      return Response.json({ message: friendlyMessage(err) }, { status: err.status });
    }
    return Response.json({ message: "Failed to load sections." }, { status: 500 });
  }
}

export async function POST(
  req: Request,
  { params }: { params: Promise<{ id: string }> }
) {
  try {
    const { user } = await requireStaff();
    if (user.role !== "teacher") return Response.json({ message: "Forbidden" }, { status: 403 });

    const { id } = await params;
    const body = await req.json().catch(() => ({}));
    const parsed = AddSectionSchema.safeParse(body);
    if (!parsed.success) {
      return Response.json({ message: "Title is required." }, { status: 400 });
    }

    const section = await apiFetch(`/courses/${id}/sections`, {
      method: "POST",
      body: parsed.data,
      schema: SectionSchema,
    });
    return Response.json(section, { status: 201 });
  } catch (err) {
    if (err instanceof ApiError) {
      return Response.json({ message: friendlyMessage(err) }, { status: err.status });
    }
    return Response.json({ message: "Failed to add section." }, { status: 500 });
  }
}

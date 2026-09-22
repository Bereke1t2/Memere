import { z } from "zod";
import { getRouteStaffSession } from "@/lib/auth/session";
import { canManageContent } from "@/lib/auth/roles";
import { getCourseSections } from "@/lib/api/endpoints";
import { apiFetch } from "@/lib/api/server";
import { AddSectionInputSchema, SectionSchema } from "@/lib/api/schemas";
import { ApiError, friendlyMessage } from "@/lib/api/errors";

export const dynamic = "force-dynamic";

export async function GET(
  _req: Request,
  { params }: { params: Promise<{ id: string }> }
) {
  try {
    const session = await getRouteStaffSession();
    if (!session) return Response.json({ message: "Unauthorized", code: "UNAUTHORIZED" }, { status: 401 });
    if (!canManageContent(session.user)) return Response.json({ message: "Forbidden" }, { status: 403 });
    const { id } = await params;
    const sections = await getCourseSections(id);
    return Response.json(sections);
  } catch (err) {
    if (err instanceof ApiError) {
      return Response.json({ message: friendlyMessage(err), details: err.details }, { status: err.status });
    }
    const message = err instanceof Error ? err.message : "Failed to load sections.";
    return Response.json({ message }, { status: 500 });
  }
}

export async function POST(
  req: Request,
  { params }: { params: Promise<{ id: string }> }
) {
  try {
    const session = await getRouteStaffSession();
    if (!session) return Response.json({ message: "Unauthorized", code: "UNAUTHORIZED" }, { status: 401 });
    if (!canManageContent(session.user)) return Response.json({ message: "Forbidden" }, { status: 403 });

    const { id } = await params;
    const body = await req.json().catch(() => ({}));
    const parsed = AddSectionInputSchema.safeParse(body);
    if (!parsed.success) {
      return Response.json({ message: "Validation failed", details: parsed.error.flatten().fieldErrors }, { status: 400 });
    }

    const payload = {
      title: parsed.data.title,
      description: parsed.data.description ?? undefined,
      order_index: parsed.data.order_index ?? parsed.data.order,
      is_published: parsed.data.is_published ?? true,
    };

    const section = await apiFetch(`/courses/${id}/sections`, {
      method: "POST",
      body: payload,
      schema: SectionSchema,
    });
    return Response.json(section, { status: 201 });
  } catch (err) {
    if (err instanceof ApiError) {
      return Response.json({ message: friendlyMessage(err), details: err.details }, { status: err.status });
    }
    const message = err instanceof Error ? err.message : "Failed to add section.";
    return Response.json({ message }, { status: 500 });
  }
}

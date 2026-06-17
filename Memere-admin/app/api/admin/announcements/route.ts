import { broadcast } from "@/lib/api/endpoints";
import { ApiError, friendlyMessage } from "@/lib/api/errors";
import { z } from "zod";

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

import { importAssessment } from "@/lib/api/endpoints";
import { ImportDocumentSchema } from "@/lib/api/schemas";
import { withStaff, validationResponse } from "@/lib/api/route-guard";

// The browser has already parsed and previewed the upload; this route re-validates
// the document against the same schema before forwarding, so a malformed payload
// never reaches the backend on the strength of client-side checks alone.
export async function POST(req: Request) {
  const body = await req.json().catch(() => ({}));
  const parsed = ImportDocumentSchema.safeParse(body);
  if (!parsed.success) {
    return validationResponse(parsed.error.issues.map((i) => `${i.path.join(".")}: ${i.message}`));
  }
  return withStaff("Failed to import.", () => importAssessment(parsed.data));
}

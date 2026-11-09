import { getExamForAuthoring } from "@/lib/api/endpoints";
import { withStaff } from "@/lib/api/route-guard";

// Returns answer keys. Staff-only via withStaff; the backend re-checks exam
// ownership before building the view.
export async function GET(_req: Request, { params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  return withStaff("Failed to load exam questions.", () => getExamForAuthoring(id));
}

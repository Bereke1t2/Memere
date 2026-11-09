import { unpublishExamEndpoint } from "@/lib/api/endpoints";
import { withStaff } from "@/lib/api/route-guard";

export async function POST(_req: Request, { params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  return withStaff("Failed to unpublish exam.", async () => {
    await unpublishExamEndpoint(id);
    return undefined;
  });
}

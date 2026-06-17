export const dynamic = "force-dynamic";

import { requireAdmin } from "@/lib/auth/session";
import { DashboardShell } from "@/components/layout/dashboard-shell";

export default async function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const { user } = await requireAdmin();
  return <DashboardShell user={user}>{children}</DashboardShell>;
}

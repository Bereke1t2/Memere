import { requireAdmin } from "@/lib/auth/session";
import { SystemHealthDashboard } from "@/components/dashboard/system-health-dashboard";

export const dynamic = "force-dynamic";

export default async function SystemHealthPage() {
  await requireAdmin();

  return (
    <div className="flex flex-col gap-6">
      <div>
        <h1 className="text-2xl font-bold tracking-tight">System Health & Diagnostics</h1>
        <p className="text-sm text-muted-foreground mt-0.5">
          Real-time Go runtime, memory allocation, and connection pool telemetry
        </p>
      </div>

      <SystemHealthDashboard />
    </div>
  );
}

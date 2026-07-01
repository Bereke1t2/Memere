import { Suspense } from "react";
import { requireStaff } from "@/lib/auth/session";
import { getOverview, getEngagement, getRevenueBreakdown, getMyEarnings, listMyCourses } from "@/lib/api/endpoints";
import { KpiCards } from "@/components/dashboard/kpi-cards";
import { Engagement } from "@/components/dashboard/engagement";
import { RevenueBar } from "@/components/charts/revenue-bar";
import { DateRange } from "@/components/dashboard/date-range";
import { EarningsDashboard } from "@/components/earnings/earnings-dashboard";
import { Skeleton } from "@/components/ui/skeleton";
import type { User } from "@/lib/api/schemas";

function toISO(d: Date): string {
  return d.toISOString().split("T")[0];
}

function resolveRange(params: { from?: string; to?: string }) {
  const to = params.to ?? toISO(new Date());
  const from = params.from ?? toISO(new Date(Date.now() - 30 * 24 * 60 * 60 * 1000));
  return { from, to };
}

function AdminSkeleton() {
  return (
    <div className="flex flex-col gap-6">
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {[0, 1, 2, 3, 4, 5].map((i) => <Skeleton key={i} className="h-28 rounded-lg" />)}
      </div>
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-3">
        {[0, 1, 2].map((i) => <Skeleton key={i} className="h-24 rounded-lg" />)}
      </div>
      <Skeleton className="h-64 rounded-lg" />
    </div>
  );
}

function TeacherSkeleton() {
  return (
    <div className="flex flex-col gap-6">
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
        {[0, 1, 2, 3].map((i) => <Skeleton key={i} className="h-28 rounded-lg" />)}
      </div>
      <Skeleton className="h-48 rounded-lg" />
    </div>
  );
}

async function AdminDashboardContent({ from, to }: { from: string; to: string }) {
  const [overview, engagement, breakdown] = await Promise.all([
    getOverview(from, to),
    getEngagement(),
    getRevenueBreakdown(from, to),
  ]);
  return (
    <div className="flex flex-col gap-6">
      <section className="flex flex-col gap-3">
        <KpiCards data={overview} />
      </section>
      <section className="flex flex-col gap-3">
        <p className="text-sm font-semibold text-muted-foreground">Engagement</p>
        <Engagement data={engagement} />
      </section>
      <section className="flex flex-col gap-3">
        <p className="text-sm font-semibold text-muted-foreground">Revenue Breakdown</p>
        <RevenueBar data={breakdown} />
      </section>
    </div>
  );
}

async function TeacherDashboardContent({ user, from, to }: { user: User; from: string; to: string }) {
  const [earnings, courseList] = await Promise.all([
    getMyEarnings(from, to),
    listMyCourses({ teacherId: user.id, limit: 50 }),
  ]);
  const courses = courseList.data;
  const sales = await Promise.all(
    courses.map((c) =>
      import("@/lib/api/endpoints")
        .then((m) => m.getCourseSales(c.id))
        .catch(() => ({ course_id: c.id, gross: "0", units: 0 }))
    )
  );
  return <EarningsDashboard earnings={earnings} courses={courses} sales={sales} from={from} to={to} />;
}

export default async function DashboardPage({
  searchParams,
}: {
  searchParams: Promise<{ from?: string; to?: string }>;
}) {
  const { user } = await requireStaff();
  const params = await searchParams;
  const { from, to } = resolveRange(params);

  if (user.role === "teacher") {
    return (
      <Suspense fallback={<TeacherSkeleton />}>
        <TeacherDashboardContent user={user} from={from} to={to} />
      </Suspense>
    );
  }

  return (
    <div className="flex flex-col gap-6">
      <div className="flex flex-col gap-1 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 className="text-2xl font-bold">Dashboard</h1>
          <p className="text-sm text-muted-foreground mt-0.5">Platform overview and analytics</p>
        </div>
        <DateRange from={from} to={to} />
      </div>
      <Suspense fallback={<AdminSkeleton />}>
        <AdminDashboardContent from={from} to={to} />
      </Suspense>
    </div>
  );
}

import { Suspense } from "react";
import { requireStaff } from "@/lib/auth/session";
import { getExamsAnalytics } from "@/lib/api/endpoints";
import { ExamStatTiles } from "@/components/analytics/exam-stat-tiles";
import { AttemptsTrend } from "@/components/analytics/attempts-trend";
import { ExamOverviewTable } from "@/components/analytics/exam-overview-table";
import { EmptyState } from "@/components/ui/empty-state";
import { Skeleton } from "@/components/ui/skeleton";

export const dynamic = "force-dynamic";

function AnalyticsSkeleton() {
  return (
    <div className="flex flex-col gap-6">
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-5">
        {[0, 1, 2, 3, 4].map((i) => <Skeleton key={i} className="h-24 rounded-lg" />)}
      </div>
      <Skeleton className="h-64 rounded-lg" />
      <Skeleton className="h-64 rounded-lg" />
    </div>
  );
}

async function AnalyticsContent() {
  const overview = await getExamsAnalytics();

  if (overview.exams.length === 0) {
    return (
      <EmptyState
        title="No exams yet"
        description="Once you create exams and students start taking them, their performance shows up here."
      />
    );
  }

  // Busiest exams first so the table leads with what students actually take.
  const rows = [...overview.exams].sort(
    (a, b) => b.stats.total_attempts - a.stats.total_attempts
  );

  return (
    <div className="flex flex-col gap-6">
      <ExamStatTiles
        totalAttempts={overview.total_attempts}
        distinctStudents={overview.distinct_students}
        avgPercentage={overview.avg_percentage}
        passRate={overview.pass_rate}
        inProgressCount={overview.in_progress_count}
      />
      <AttemptsTrend daily={overview.daily} />
      <ExamOverviewTable rows={rows} />
    </div>
  );
}

export default async function ExamAnalyticsPage() {
  const { user } = await requireStaff();
  const isAdmin = user.role === "admin";

  return (
    <div className="flex flex-col gap-6">
      <div>
        <h1 className="text-2xl font-semibold tracking-tight">Exam Analytics</h1>
        <p className="mt-1 text-sm text-muted-foreground">
          {isAdmin
            ? "How every exam on the platform is doing: attempts, students, scores, and pass rates."
            : "How the exams in your courses are doing: attempts, students, scores, and pass rates."}
        </p>
      </div>
      <Suspense fallback={<AnalyticsSkeleton />}>
        <AnalyticsContent />
      </Suspense>
    </div>
  );
}

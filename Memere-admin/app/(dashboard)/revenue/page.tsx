import { Suspense } from "react";
import { getOverview, getRevenueBreakdown } from "@/lib/api/endpoints";
import { DateRange } from "@/components/dashboard/date-range";
import { RevenueBar } from "@/components/charts/revenue-bar";
import { UnitsPie } from "@/components/charts/units-pie";
import { RevenueTrend, type TrendPoint } from "@/components/charts/revenue-trend";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { formatMoney, PLATFORM_CURRENCY } from "@/lib/format";
import type { RevenueBreakdownItem } from "@/lib/api/schemas";

const TREND_BUCKETS = 4;

const CURRENCY = PLATFORM_CURRENCY;

function toISO(d: Date) {
  return d.toISOString().split("T")[0];
}

function resolveRange(p: { from?: string; to?: string }) {
  return {
    to: p.to ?? toISO(new Date()),
    from: p.from ?? toISO(new Date(Date.now() - 30 * 24 * 60 * 60 * 1000)),
  };
}

/** Splits [from, to] into N equal sub-ranges for the trend chart. */
function splitRange(from: string, to: string, buckets: number) {
  const start = new Date(from).getTime();
  const end = new Date(to).getTime();
  const step = Math.max(end - start, 0) / buckets;
  return Array.from({ length: buckets }, (_, i) => ({
    from: toISO(new Date(start + i * step)),
    to: toISO(new Date(start + (i + 1) * step)),
  }));
}

function KpiCard({ title, value }: { title: string; value: string }) {
  return (
    <Card>
      <CardHeader className="pb-1">
        <CardTitle className="text-sm font-medium text-muted-foreground">
          {title}
        </CardTitle>
      </CardHeader>
      <CardContent>
        <p className="text-2xl font-bold tracking-tight tabular-nums">{value}</p>
      </CardContent>
    </Card>
  );
}

function BreakdownTable({ data }: { data: RevenueBreakdownItem[] }) {
  if (!data.length) return null;
  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-base">Provider Breakdown</CardTitle>
      </CardHeader>
      <CardContent className="p-0">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b bg-muted/50">
              <th className="px-4 py-2.5 text-left font-medium text-muted-foreground">
                Provider
              </th>
              <th className="px-4 py-2.5 text-right font-medium text-muted-foreground">
                Gross
              </th>
              <th className="px-4 py-2.5 text-right font-medium text-muted-foreground">
                Payments
              </th>
            </tr>
          </thead>
          <tbody>
            {data.map((item) => (
              <tr key={item.provider} className="border-b last:border-0">
                <td className="px-4 py-2.5 capitalize">{item.provider}</td>
                <td className="px-4 py-2.5 text-right tabular-nums font-medium">
                  {formatMoney(item.gross, CURRENCY)}
                </td>
                <td className="px-4 py-2.5 text-right tabular-nums text-muted-foreground">
                  {item.units.toLocaleString()}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </CardContent>
    </Card>
  );
}

function RevenueSkeleton() {
  return (
    <div className="flex flex-col gap-6">
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
        {[0, 1, 2, 3].map((i) => (
          <Skeleton key={i} className="h-28 rounded-lg" />
        ))}
      </div>
      <div className="grid grid-cols-1 gap-4 lg:grid-cols-2">
        <Skeleton className="h-64 rounded-lg" />
        <Skeleton className="h-64 rounded-lg" />
      </div>
      <Skeleton className="h-48 rounded-lg" />
      <Skeleton className="h-56 rounded-lg" />
    </div>
  );
}

async function RevenueContent({ from, to }: { from: string; to: string }) {
  const buckets = splitRange(from, to, TREND_BUCKETS);

  const [overview, breakdown, bucketOverviews] = await Promise.all([
    getOverview(from, to),
    getRevenueBreakdown(from, to),
    Promise.all(buckets.map((b) => getOverview(b.from, b.to))),
  ]);

  const trend: TrendPoint[] = bucketOverviews.map((o, i) => ({
    label: buckets[i].to,
    gross: parseFloat(o.gross_revenue) || 0,
    grossLabel: formatMoney(o.gross_revenue, CURRENCY),
  }));

  return (
    <div className="flex flex-col gap-6">
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <KpiCard
          title="Gross Revenue"
          value={formatMoney(overview.gross_revenue, CURRENCY)}
        />
        <KpiCard
          title="Refunded"
          value={formatMoney(overview.refunded_amount, CURRENCY)}
        />
        <KpiCard title="MRR" value={formatMoney(overview.mrr, CURRENCY)} />
        <KpiCard
          title="Completed Payments"
          value={overview.completed_payments.toLocaleString()}
        />
      </div>

      <div className="grid grid-cols-1 gap-4 lg:grid-cols-2">
        <RevenueBar data={breakdown} />
        <UnitsPie data={breakdown} />
      </div>

      <BreakdownTable data={breakdown} />

      <div className="flex flex-col gap-2">
        <div>
          <h2 className="text-lg font-semibold tracking-tight">
            Financial Overview
          </h2>
          <p className="text-xs text-muted-foreground mt-0.5">
            KPI-level summary only — the backend has no admin endpoint to list
            individual subscriptions in v1. The trend below buckets the selected
            range into {TREND_BUCKETS} periods using the same overview endpoint.
          </p>
        </div>
        <RevenueTrend data={trend} />
      </div>
    </div>
  );
}

export default async function RevenuePage({
  searchParams,
}: {
  searchParams: Promise<{ from?: string; to?: string }>;
}) {
  const params = await searchParams;
  const { from, to } = resolveRange(params);

  return (
    <div className="flex flex-col gap-6">
      <div className="flex items-center justify-between flex-wrap gap-3">
        <div>
          <h1 className="text-2xl font-semibold tracking-tight">Revenue</h1>
          <p className="text-sm text-muted-foreground mt-1">
            Financial overview for the selected period.
          </p>
        </div>
        <DateRange from={from} to={to} />
      </div>

      <Suspense fallback={<RevenueSkeleton />}>
        <RevenueContent from={from} to={to} />
      </Suspense>
    </div>
  );
}

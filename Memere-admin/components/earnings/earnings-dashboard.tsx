import { Card, CardContent } from "@/components/ui/card";
import { DateRange } from "@/components/dashboard/date-range";
import { formatMoney, PLATFORM_CURRENCY } from "@/lib/format";
import type { Earnings, Course, CourseSales } from "@/lib/api/schemas";

interface EarningsDashboardProps {
  earnings: Earnings;
  courses: Course[];
  sales: CourseSales[];
  from: string;
  to: string;
}

function KpiCard({ title, value, sub }: { title: string; value: string; sub?: string }) {
  return (
    <Card>
      <CardContent className="p-6">
        <p className="text-sm text-muted-foreground">{title}</p>
        <p className="mt-2 text-2xl font-semibold tabular-nums">{value}</p>
        {sub && <p className="mt-1 text-xs text-muted-foreground">{sub}</p>}
      </CardContent>
    </Card>
  );
}

export function EarningsDashboard({ earnings, courses, sales, from, to }: EarningsDashboardProps) {
  const sharePct = Math.round(parseFloat(earnings.teacher_share) * 100);

  const salesMap = new Map(sales.map((s) => [s.course_id, s]));

  return (
    <div className="flex flex-col gap-6">
      <div className="flex flex-col gap-1 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 className="text-2xl font-semibold">Earnings</h1>
          <p className="text-sm text-muted-foreground mt-0.5">Your revenue for the selected period</p>
        </div>
        <DateRange from={from} to={to} />
      </div>

      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <KpiCard
          title="Gross Revenue"
          value={formatMoney(earnings.gross, PLATFORM_CURRENCY)}
          sub="Total sales before fee"
        />
        <KpiCard
          title={`Your Earnings (${sharePct}%)`}
          value={formatMoney(earnings.earnings, PLATFORM_CURRENCY)}
          sub={`Platform keeps ${100 - sharePct}%`}
        />
        <KpiCard
          title="Platform Fee"
          value={formatMoney(earnings.platform_fee, PLATFORM_CURRENCY)}
        />
        <KpiCard
          title="Payments"
          value={earnings.units.toLocaleString()}
          sub="Completed transactions"
        />
      </div>

      {courses.length > 0 && (
        <Card>
          <CardContent className="p-0">
            <div className="px-6 py-4 border-b">
              <p className="text-sm font-semibold">Per-Course Sales</p>
            </div>
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b bg-muted/50">
                  <th className="px-6 py-3 text-left text-xs font-semibold text-muted-foreground">Course</th>
                  <th className="px-6 py-3 text-right text-xs font-semibold text-muted-foreground">Gross</th>
                  <th className="px-6 py-3 text-right text-xs font-semibold text-muted-foreground">Your Cut</th>
                  <th className="px-6 py-3 text-right text-xs font-semibold text-muted-foreground">Sales</th>
                </tr>
              </thead>
              <tbody>
                {courses.map((course) => {
                  const s = salesMap.get(course.id);
                  const gross = parseFloat(s?.gross ?? "0");
                  const cut = gross * parseFloat(earnings.teacher_share);
                  return (
                    <tr key={course.id} className="border-b last:border-0 hover:bg-muted/30 transition-colors">
                      <td className="px-6 py-3 font-medium">{course.title}</td>
                      <td className="px-6 py-3 text-right tabular-nums">{formatMoney(String(gross), PLATFORM_CURRENCY)}</td>
                      <td className="px-6 py-3 text-right tabular-nums">{formatMoney(String(cut), PLATFORM_CURRENCY)}</td>
                      <td className="px-6 py-3 text-right tabular-nums text-muted-foreground">{s?.units ?? 0}</td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </CardContent>
        </Card>
      )}
    </div>
  );
}

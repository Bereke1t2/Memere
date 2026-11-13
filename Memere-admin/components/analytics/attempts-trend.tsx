"use client";

import {
  AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer,
} from "recharts";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import type { DailyAttempts } from "@/lib/api/schemas";

interface AttemptsTrendProps {
  daily: DailyAttempts[];
  /** Window to render; the backend series covers the last 30 days. */
  days?: number;
  title?: string;
}

interface Point {
  date: string;
  label: string;
  attempts: number;
}

/** The backend omits days with no attempts, so densify to a continuous axis. */
function fillDays(daily: DailyAttempts[], days: number): Point[] {
  const byDate = new Map(daily.map((d) => [d.date, d.attempts]));
  const out: Point[] = [];
  const today = new Date();
  for (let i = days - 1; i >= 0; i--) {
    const d = new Date(today);
    d.setDate(today.getDate() - i);
    const iso = d.toISOString().split("T")[0];
    out.push({
      date: iso,
      label: d.toLocaleDateString("en-US", { month: "short", day: "numeric" }),
      attempts: byDate.get(iso) ?? 0,
    });
  }
  return out;
}

function CustomTooltip({ active, payload }: { active?: boolean; payload?: { payload: Point }[] }) {
  if (!active || !payload?.length) return null;
  const p = payload[0].payload;
  return (
    <div className="rounded-lg border bg-card px-3 py-2.5 text-sm shadow-md">
      <p className="font-medium">{p.label}</p>
      <p className="text-muted-foreground mt-0.5">
        {p.attempts} attempt{p.attempts === 1 ? "" : "s"}
      </p>
    </div>
  );
}

export function AttemptsTrend({ daily, days = 30, title = "Attempts — last 30 days" }: AttemptsTrendProps) {
  const data = fillDays(daily, days);
  const total = data.reduce((sum, p) => sum + p.attempts, 0);

  return (
    <Card>
      <CardHeader className="pb-2">
        <CardTitle className="text-sm font-semibold">{title}</CardTitle>
      </CardHeader>
      <CardContent>
        {total === 0 ? (
          <div className="flex items-center justify-center py-12">
            <p className="text-sm text-muted-foreground">No attempts in this period.</p>
          </div>
        ) : (
          <ResponsiveContainer width="100%" height={240}>
            <AreaChart data={data} margin={{ top: 4, right: 8, left: 0, bottom: 0 }}>
              <CartesianGrid strokeDasharray="3 3" stroke="hsl(var(--border))" vertical={false} />
              <XAxis
                dataKey="label"
                tick={{ fontSize: 11, fill: "hsl(var(--muted-foreground))" }}
                axisLine={false}
                tickLine={false}
                interval="preserveStartEnd"
                minTickGap={32}
              />
              <YAxis
                allowDecimals={false}
                tick={{ fontSize: 11, fill: "hsl(var(--muted-foreground))" }}
                axisLine={false}
                tickLine={false}
                width={32}
              />
              <Tooltip content={<CustomTooltip />} cursor={{ stroke: "hsl(var(--border))" }} />
              <Area
                type="monotone"
                dataKey="attempts"
                stroke="hsl(var(--chart-1))"
                strokeWidth={2}
                fill="hsl(var(--chart-1))"
                fillOpacity={0.12}
                activeDot={{ r: 4 }}
              />
            </AreaChart>
          </ResponsiveContainer>
        )}
      </CardContent>
    </Card>
  );
}

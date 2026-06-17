"use client";

import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from "recharts";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { formatMoney, PLATFORM_CURRENCY } from "@/lib/format";
import type { RevenueBreakdownItem } from "@/lib/api/schemas";

interface RevenueBarProps {
  data: RevenueBreakdownItem[];
}

interface TooltipPayloadEntry {
  payload: RevenueBreakdownItem & { grossNum: number };
}

function CustomTooltip({
  active,
  payload,
}: {
  active?: boolean;
  payload?: TooltipPayloadEntry[];
}) {
  if (!active || !payload?.length) return null;
  const item = payload[0].payload;
  return (
    <div className="rounded-xl border-0 shadow-lg bg-white dark:bg-zinc-900 px-4 py-3 text-sm">
      <p className="font-medium capitalize">{item.provider}</p>
      <p className="text-muted-foreground">{formatMoney(item.gross, PLATFORM_CURRENCY)}</p>
      <p className="text-muted-foreground">{item.units} payments</p>
    </div>
  );
}

export function RevenueBar({ data }: RevenueBarProps) {
  if (!data.length) {
    return (
      <Card className="border-0 shadow-sm">
        <CardHeader>
          <CardTitle className="text-sm font-semibold">Revenue by Provider</CardTitle>
          <p className="text-xs text-muted-foreground">Revenue split by payment provider</p>
        </CardHeader>
        <CardContent className="flex items-center justify-center py-12">
          <p className="text-sm text-muted-foreground">No revenue data for this period.</p>
        </CardContent>
      </Card>
    );
  }

  const chartData = data.map((item) => ({
    ...item,
    grossNum: parseFloat(item.gross) || 0,
  }));

  return (
    <Card className="border-0 shadow-sm">
      <CardHeader>
        <CardTitle className="text-sm font-semibold">Revenue by Provider</CardTitle>
        <p className="text-xs text-muted-foreground">Revenue split by payment provider</p>
      </CardHeader>
      <CardContent>
        <ResponsiveContainer width="100%" height={280}>
          <BarChart data={chartData} margin={{ top: 4, right: 16, left: 0, bottom: 0 }}>
            <defs>
              <linearGradient id="revenueGrad" x1="0" y1="0" x2="0" y2="1">
                <stop offset="0%" stopColor="#6366f1" stopOpacity={1} />
                <stop offset="100%" stopColor="#818cf8" stopOpacity={0.7} />
              </linearGradient>
            </defs>
            <CartesianGrid strokeDasharray="3 3" stroke="#e4e4e7" />
            <XAxis
              dataKey="provider"
              tick={{ fontSize: 11, fill: '#71717a' }}
              axisLine={false}
              tickLine={false}
              tickFormatter={(v: string) =>
                v.charAt(0).toUpperCase() + v.slice(1)
              }
            />
            <YAxis
              tick={{ fontSize: 11, fill: '#71717a' }}
              axisLine={false}
              tickLine={false}
              tickFormatter={(v: number) =>
                v >= 1000 ? `${(v / 1000).toFixed(0)}k` : String(v)
              }
            />
            <Tooltip content={<CustomTooltip />} />
            <Bar dataKey="grossNum" name="Revenue" radius={[6, 6, 0, 0]} fill="url(#revenueGrad)" />
          </BarChart>
        </ResponsiveContainer>
      </CardContent>
    </Card>
  );
}

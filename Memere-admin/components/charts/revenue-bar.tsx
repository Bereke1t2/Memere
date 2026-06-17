"use client";

import {
  BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer,
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

function CustomTooltip({ active, payload }: { active?: boolean; payload?: TooltipPayloadEntry[] }) {
  if (!active || !payload?.length) return null;
  const item = payload[0].payload;
  return (
    <div className="rounded-lg border bg-card px-3 py-2.5 text-sm shadow-md">
      <p className="font-medium capitalize">{item.provider}</p>
      <p className="text-muted-foreground mt-0.5">{formatMoney(item.gross, PLATFORM_CURRENCY)}</p>
      <p className="text-muted-foreground">{item.units} payments</p>
    </div>
  );
}

export function RevenueBar({ data }: RevenueBarProps) {
  if (!data.length) {
    return (
      <Card className="border shadow-none">
        <CardHeader className="pb-2">
          <CardTitle className="text-sm font-semibold">Revenue by Provider</CardTitle>
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
    <Card className="border shadow-none">
      <CardHeader className="pb-2">
        <CardTitle className="text-sm font-semibold">Revenue by Provider</CardTitle>
      </CardHeader>
      <CardContent>
        <ResponsiveContainer width="100%" height={260}>
          <BarChart data={chartData} margin={{ top: 4, right: 8, left: 0, bottom: 0 }}>
            <CartesianGrid strokeDasharray="3 3" stroke="hsl(var(--border))" vertical={false} />
            <XAxis
              dataKey="provider"
              tick={{ fontSize: 11, fill: "hsl(var(--muted-foreground))" }}
              axisLine={false}
              tickLine={false}
              tickFormatter={(v: string) => v.charAt(0).toUpperCase() + v.slice(1)}
            />
            <YAxis
              tick={{ fontSize: 11, fill: "hsl(var(--muted-foreground))" }}
              axisLine={false}
              tickLine={false}
              tickFormatter={(v: number) => v >= 1000 ? `${(v / 1000).toFixed(0)}k` : String(v)}
            />
            <Tooltip content={<CustomTooltip />} cursor={{ fill: "hsl(var(--muted))" }} />
            <Bar dataKey="grossNum" name="Revenue" radius={[4, 4, 0, 0]} fill="hsl(var(--foreground))" />
          </BarChart>
        </ResponsiveContainer>
      </CardContent>
    </Card>
  );
}

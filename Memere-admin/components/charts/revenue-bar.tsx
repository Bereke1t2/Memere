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
    <div className="rounded-md border bg-popover px-3 py-2 text-sm shadow-md">
      <p className="font-medium capitalize">{item.provider}</p>
      <p className="text-muted-foreground">{formatMoney(item.gross, PLATFORM_CURRENCY)}</p>
      <p className="text-muted-foreground">{item.units} payments</p>
    </div>
  );
}

export function RevenueBar({ data }: RevenueBarProps) {
  if (!data.length) {
    return (
      <Card>
        <CardHeader>
          <CardTitle className="text-base">Revenue by Provider</CardTitle>
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
    <Card>
      <CardHeader>
        <CardTitle className="text-base">Revenue by Provider</CardTitle>
      </CardHeader>
      <CardContent>
        <ResponsiveContainer width="100%" height={240}>
          <BarChart data={chartData} margin={{ top: 4, right: 16, left: 0, bottom: 0 }}>
            <CartesianGrid strokeDasharray="3 3" className="stroke-border" />
            <XAxis
              dataKey="provider"
              tick={{ fontSize: 12 }}
              tickFormatter={(v: string) =>
                v.charAt(0).toUpperCase() + v.slice(1)
              }
            />
            <YAxis
              tick={{ fontSize: 12 }}
              tickFormatter={(v: number) =>
                v >= 1000 ? `${(v / 1000).toFixed(0)}k` : String(v)
              }
            />
            <Tooltip content={<CustomTooltip />} />
            <Bar dataKey="grossNum" name="Revenue" radius={[4, 4, 0, 0]} className="fill-primary" />
          </BarChart>
        </ResponsiveContainer>
      </CardContent>
    </Card>
  );
}

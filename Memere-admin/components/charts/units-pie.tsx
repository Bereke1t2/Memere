"use client";

import {
  PieChart,
  Pie,
  Cell,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from "recharts";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import type { RevenueBreakdownItem } from "@/lib/api/schemas";

const PALETTE = ["#2563eb", "#16a34a", "#d97706", "#9333ea", "#ef4444", "#0891b2"];

interface UnitsPieProps {
  data: RevenueBreakdownItem[];
}

interface TooltipPayloadEntry {
  name: string;
  value: number;
}

function CustomTooltip({
  active,
  payload,
}: {
  active?: boolean;
  payload?: TooltipPayloadEntry[];
}) {
  if (!active || !payload?.length) return null;
  return (
    <div className="rounded-md border bg-popover px-3 py-2 text-sm shadow-md">
      <p className="font-medium capitalize">{payload[0].name}</p>
      <p className="text-muted-foreground">
        {payload[0].value} payment{payload[0].value !== 1 ? "s" : ""}
      </p>
    </div>
  );
}

export function UnitsPie({ data }: UnitsPieProps) {
  if (!data.length) {
    return (
      <Card>
        <CardHeader>
          <CardTitle className="text-base">Payments by Provider</CardTitle>
        </CardHeader>
        <CardContent className="flex items-center justify-center py-12">
          <p className="text-sm text-muted-foreground">
            No payment data for this period.
          </p>
        </CardContent>
      </Card>
    );
  }

  const chartData = data.map((item) => ({
    name: item.provider.charAt(0).toUpperCase() + item.provider.slice(1),
    value: item.units,
  }));

  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-base">Payments by Provider</CardTitle>
      </CardHeader>
      <CardContent>
        <ResponsiveContainer width="100%" height={240}>
          <PieChart>
            <Pie
              data={chartData}
              dataKey="value"
              nameKey="name"
              cx="50%"
              cy="50%"
              outerRadius={80}
              label={({ name, percent }: { name?: string; percent?: number }) =>
                `${name ?? ""} ${((percent ?? 0) * 100).toFixed(0)}%`
              }
              labelLine={false}
            >
              {chartData.map((_, i) => (
                <Cell
                  key={i}
                  fill={PALETTE[i % PALETTE.length]}
                />
              ))}
            </Pie>
            <Tooltip content={<CustomTooltip />} />
            <Legend />
          </PieChart>
        </ResponsiveContainer>
      </CardContent>
    </Card>
  );
}

import { Users, GraduationCap, DollarSign, TrendingUp, CheckCircle2, RotateCcw } from "lucide-react";
import { Card, CardContent } from "@/components/ui/card";
import { formatMoney, formatCount, PLATFORM_CURRENCY } from "@/lib/format";
import type { Overview } from "@/lib/api/schemas";
import type { LucideIcon } from "lucide-react";

interface KpiCardsProps {
  data: Overview;
}

interface CardConfig {
  title: string;
  render: (d: Overview) => string;
  icon: LucideIcon;
}

const CARDS: CardConfig[] = [
  { title: "Total Students",      render: (d) => formatCount(d.total_students),              icon: Users },
  { title: "Total Teachers",      render: (d) => formatCount(d.total_teachers),              icon: GraduationCap },
  { title: "Gross Revenue",       render: (d) => formatMoney(d.gross_revenue, PLATFORM_CURRENCY), icon: DollarSign },
  { title: "MRR",                 render: (d) => formatMoney(d.mrr, PLATFORM_CURRENCY),      icon: TrendingUp },
  { title: "Completed Payments",  render: (d) => formatCount(d.completed_payments),          icon: CheckCircle2 },
  { title: "Refunded Amount",     render: (d) => formatMoney(d.refunded_amount, PLATFORM_CURRENCY), icon: RotateCcw },
];

export function KpiCards({ data }: KpiCardsProps) {
  return (
    <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
      {CARDS.map(({ title, render, icon: Icon }) => (
        <Card key={title} className="border shadow-none">
          <CardContent className="p-6">
            <div className="flex items-start justify-between">
              <div>
                <p className="text-sm text-muted-foreground">{title}</p>
                <p className="mt-2 text-2xl font-semibold tracking-tight tabular-nums">{render(data)}</p>
              </div>
              <div className="rounded-lg bg-muted p-2">
                <Icon className="h-4 w-4 text-muted-foreground" />
              </div>
            </div>
          </CardContent>
        </Card>
      ))}
    </div>
  );
}

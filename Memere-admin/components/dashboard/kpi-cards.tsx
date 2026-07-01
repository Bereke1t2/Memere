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
  tone: string;
}

const CARDS: CardConfig[] = [
  { title: "Total Students", render: (d) => formatCount(d.total_students), icon: Users, tone: "bg-emerald-50 text-emerald-700 border-emerald-100" },
  { title: "Total Teachers", render: (d) => formatCount(d.total_teachers), icon: GraduationCap, tone: "bg-sky-50 text-sky-700 border-sky-100" },
  { title: "Gross Revenue", render: (d) => formatMoney(d.gross_revenue, PLATFORM_CURRENCY), icon: DollarSign, tone: "bg-amber-50 text-amber-700 border-amber-100" },
  { title: "MRR", render: (d) => formatMoney(d.mrr, PLATFORM_CURRENCY), icon: TrendingUp, tone: "bg-teal-50 text-teal-700 border-teal-100" },
  { title: "Completed Payments", render: (d) => formatCount(d.completed_payments), icon: CheckCircle2, tone: "bg-lime-50 text-lime-700 border-lime-100" },
  { title: "Refunded Amount", render: (d) => formatMoney(d.refunded_amount, PLATFORM_CURRENCY), icon: RotateCcw, tone: "bg-rose-50 text-rose-700 border-rose-100" },
];

export function KpiCards({ data }: KpiCardsProps) {
  return (
    <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
      {CARDS.map(({ title, render, icon: Icon, tone }) => (
        <Card key={title}>
          <CardContent className="p-6">
            <div className="flex items-start justify-between">
              <div>
                <p className="text-sm text-muted-foreground">{title}</p>
                <p className="mt-2 text-2xl font-semibold tabular-nums">{render(data)}</p>
              </div>
              <div className={`rounded-lg border p-2 ${tone}`}>
                <Icon className="h-4 w-4" />
              </div>
            </div>
          </CardContent>
        </Card>
      ))}
    </div>
  );
}

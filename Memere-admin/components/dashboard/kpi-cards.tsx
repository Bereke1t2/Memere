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
  colorClass: string;
}

const CARDS: CardConfig[] = [
  {
    title: "Total Students",
    render: (d: Overview) => formatCount(d.total_students),
    icon: Users,
    colorClass: "bg-blue-50 text-blue-600 dark:bg-blue-950 dark:text-blue-400",
  },
  {
    title: "Total Teachers",
    render: (d: Overview) => formatCount(d.total_teachers),
    icon: GraduationCap,
    colorClass: "bg-emerald-50 text-emerald-600 dark:bg-emerald-950 dark:text-emerald-400",
  },
  {
    title: "Gross Revenue",
    render: (d: Overview) => formatMoney(d.gross_revenue, PLATFORM_CURRENCY),
    icon: DollarSign,
    colorClass: "bg-violet-50 text-violet-600 dark:bg-violet-950 dark:text-violet-400",
  },
  {
    title: "MRR",
    render: (d: Overview) => formatMoney(d.mrr, PLATFORM_CURRENCY),
    icon: TrendingUp,
    colorClass: "bg-indigo-50 text-indigo-600 dark:bg-indigo-950 dark:text-indigo-400",
  },
  {
    title: "Completed Payments",
    render: (d: Overview) => formatCount(d.completed_payments),
    icon: CheckCircle2,
    colorClass: "bg-green-50 text-green-600 dark:bg-green-950 dark:text-green-400",
  },
  {
    title: "Refunded Amount",
    render: (d: Overview) => formatMoney(d.refunded_amount, PLATFORM_CURRENCY),
    icon: RotateCcw,
    colorClass: "bg-orange-50 text-orange-600 dark:bg-orange-950 dark:text-orange-400",
  },
];

export function KpiCards({ data }: KpiCardsProps) {
  return (
    <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
      {CARDS.map(({ title, render, icon: Icon, colorClass }) => (
        <Card key={title} className="relative overflow-hidden border-0 shadow-sm">
          <CardContent className="p-6">
            <div className="flex items-start justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">{title}</p>
                <p className="mt-2 text-3xl font-bold tracking-tight">{render(data)}</p>
              </div>
              <div className={`rounded-xl p-2.5 ${colorClass}`}>
                <Icon className="h-5 w-5" />
              </div>
            </div>
          </CardContent>
        </Card>
      ))}
    </div>
  );
}

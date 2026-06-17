import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { formatMoney, formatCount, PLATFORM_CURRENCY } from "@/lib/format";
import type { Overview } from "@/lib/api/schemas";

interface KpiCardsProps {
  data: Overview;
}

const CARDS = [
  {
    title: "Total Students",
    render: (d: Overview) => formatCount(d.total_students),
  },
  {
    title: "Total Teachers",
    render: (d: Overview) => formatCount(d.total_teachers),
  },
  {
    title: "Gross Revenue",
    render: (d: Overview) => formatMoney(d.gross_revenue, PLATFORM_CURRENCY),
  },
  {
    title: "MRR",
    render: (d: Overview) => formatMoney(d.mrr, PLATFORM_CURRENCY),
  },
  {
    title: "Completed Payments",
    render: (d: Overview) => formatCount(d.completed_payments),
  },
  {
    title: "Refunded Amount",
    render: (d: Overview) => formatMoney(d.refunded_amount, PLATFORM_CURRENCY),
  },
] as const;

export function KpiCards({ data }: KpiCardsProps) {
  return (
    <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
      {CARDS.map(({ title, render }) => (
        <Card key={title}>
          <CardHeader className="pb-1">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              {title}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-2xl font-bold tracking-tight">{render(data)}</p>
          </CardContent>
        </Card>
      ))}
    </div>
  );
}

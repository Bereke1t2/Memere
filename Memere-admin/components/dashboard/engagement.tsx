import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { formatPercent } from "@/lib/format";
import type { EngagementStats } from "@/lib/api/schemas";

interface EngagementProps {
  data: EngagementStats;
}

const TILES = [
  {
    title: "Quiz Pass Rate",
    value: (d: EngagementStats) => d.avg_quiz_pass_rate,
  },
  {
    title: "Exam Pass Rate",
    value: (d: EngagementStats) => d.avg_exam_pass_rate,
  },
  {
    title: "Course Completion",
    value: (d: EngagementStats) => d.avg_completion_pct,
  },
] as const;

export function Engagement({ data }: EngagementProps) {
  return (
    <div className="grid grid-cols-1 gap-4 sm:grid-cols-3">
      {TILES.map(({ title, value }) => {
        const pct = value(data);
        const width = Math.min(100, Math.max(0, pct * 100));
        return (
          <Card key={title}>
            <CardHeader className="pb-1">
              <CardTitle className="text-sm font-medium text-muted-foreground">
                {title}
              </CardTitle>
            </CardHeader>
            <CardContent className="flex flex-col gap-2">
              <p className="text-2xl font-bold">{formatPercent(pct)}</p>
              <div
                className="h-1.5 rounded-full bg-muted overflow-hidden"
                role="progressbar"
                aria-valuenow={Math.round(width)}
                aria-valuemin={0}
                aria-valuemax={100}
                aria-label={title}
              >
                <div
                  className="h-full rounded-full bg-primary transition-all"
                  style={{ width: `${width}%` }}
                />
              </div>
            </CardContent>
          </Card>
        );
      })}
    </div>
  );
}

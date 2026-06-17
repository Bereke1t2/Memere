import { Card, CardContent } from "@/components/ui/card";
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
          <Card key={title} className="border-0 shadow-sm">
            <CardContent className="p-6">
              <div className="flex items-center justify-between mb-3">
                <p className="text-sm font-medium text-muted-foreground">{title}</p>
                <span className="text-2xl font-bold">{formatPercent(pct)}</span>
              </div>
              <div
                className="h-2 rounded-full bg-zinc-100 dark:bg-zinc-800 overflow-hidden"
                role="progressbar"
                aria-valuenow={Math.round(width)}
                aria-valuemin={0}
                aria-valuemax={100}
                aria-label={title}
              >
                <div
                  className="h-full rounded-full bg-gradient-to-r from-indigo-500 to-violet-500 transition-all duration-700"
                  style={{ width: `${width}%` }}
                />
              </div>
              <p className="mt-2 text-xs text-muted-foreground">{Math.round(width)}% of target</p>
            </CardContent>
          </Card>
        );
      })}
    </div>
  );
}

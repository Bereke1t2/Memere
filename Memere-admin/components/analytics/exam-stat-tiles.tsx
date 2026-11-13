import { ClipboardList, Users, Percent, CheckCircle2, Timer } from "lucide-react";
import { Card, CardContent } from "@/components/ui/card";
import { formatCount } from "@/lib/format";
import type { LucideIcon } from "lucide-react";

interface ExamStatTilesProps {
  totalAttempts: number;
  distinctStudents: number;
  avgPercentage: number;
  passRate: number;
  inProgressCount: number;
}

interface Tile {
  title: string;
  value: string;
  icon: LucideIcon;
  tone: string;
}

export function ExamStatTiles(props: ExamStatTilesProps) {
  const noAttempts = props.totalAttempts === 0;
  const tiles: Tile[] = [
    { title: "Graded Attempts", value: formatCount(props.totalAttempts), icon: ClipboardList, tone: "bg-emerald-50 text-emerald-700 border-emerald-100" },
    { title: "Students", value: formatCount(props.distinctStudents), icon: Users, tone: "bg-sky-50 text-sky-700 border-sky-100" },
    { title: "Average Score", value: noAttempts ? "—" : `${props.avgPercentage.toFixed(1)}%`, icon: Percent, tone: "bg-amber-50 text-amber-700 border-amber-100" },
    { title: "Pass Rate", value: noAttempts ? "—" : `${props.passRate.toFixed(0)}%`, icon: CheckCircle2, tone: "bg-teal-50 text-teal-700 border-teal-100" },
    { title: "Taking Now", value: formatCount(props.inProgressCount), icon: Timer, tone: "bg-rose-50 text-rose-700 border-rose-100" },
  ];

  return (
    <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-5">
      {tiles.map(({ title, value, icon: Icon, tone }) => (
        <Card key={title}>
          <CardContent className="p-5">
            <div className="flex items-start justify-between">
              <div>
                <p className="text-sm text-muted-foreground">{title}</p>
                <p className="mt-2 text-2xl font-semibold tabular-nums">{value}</p>
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

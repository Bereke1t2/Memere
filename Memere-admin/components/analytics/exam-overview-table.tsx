import Link from "next/link";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { formatGrade } from "@/lib/format";
import type { Exam, ExamOverviewStat } from "@/lib/api/schemas";

interface ExamOverviewTableProps {
  rows: { exam: Exam; stats: ExamOverviewStat }[];
}

function formatLastActivity(iso?: string | null): string {
  if (!iso) return "—";
  const d = new Date(iso);
  if (Number.isNaN(d.getTime())) return "—";
  return d.toLocaleDateString("en-US", { month: "short", day: "numeric", year: "numeric" });
}

/** Inline meter so the score column can be scanned without reading each number. */
function ScoreMeter({ pct }: { pct: number }) {
  return (
    <div className="flex items-center gap-2">
      <span className="w-12 text-right tabular-nums">{pct.toFixed(1)}%</span>
      <div className="h-1.5 w-16 overflow-hidden rounded-full bg-muted" aria-hidden="true">
        <div
          className="h-full rounded-full bg-[hsl(var(--chart-1))]"
          style={{ width: `${Math.min(100, Math.max(0, pct))}%` }}
        />
      </div>
    </div>
  );
}

export function ExamOverviewTable({ rows }: ExamOverviewTableProps) {
  return (
    <Card>
      <CardHeader className="pb-2">
        <CardTitle className="text-sm font-semibold">Per-Exam Performance</CardTitle>
      </CardHeader>
      <CardContent className="overflow-x-auto">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Exam</TableHead>
              <TableHead className="text-right">Attempts</TableHead>
              <TableHead className="text-right">Students</TableHead>
              <TableHead>Avg score</TableHead>
              <TableHead className="text-right">Pass rate</TableHead>
              <TableHead className="text-right">Taking now</TableHead>
              <TableHead className="text-right">Last activity</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {rows.map(({ exam, stats }) => {
              const hasAttempts = stats.total_attempts > 0;
              return (
                <TableRow key={exam.id}>
                  <TableCell>
                    <Link href={`/exams/${exam.id}`} className="group flex flex-col gap-0.5">
                      <span className="text-sm font-medium group-hover:underline">{exam.title}</span>
                      <span className="flex items-center gap-2 text-xs text-muted-foreground">
                        {exam.subject && <span>{exam.subject}</span>}
                        {exam.grade != null && <span>{formatGrade(exam.grade)}</span>}
                        <Badge
                          variant={exam.is_published ? "default" : "outline"}
                          className="px-1.5 py-0 text-[10px]"
                        >
                          {exam.is_published ? "Published" : "Draft"}
                        </Badge>
                      </span>
                    </Link>
                  </TableCell>
                  <TableCell className="text-right tabular-nums">{stats.total_attempts}</TableCell>
                  <TableCell className="text-right tabular-nums">{stats.distinct_students}</TableCell>
                  <TableCell>
                    {hasAttempts ? (
                      <ScoreMeter pct={stats.avg_percentage} />
                    ) : (
                      <span className="text-muted-foreground">—</span>
                    )}
                  </TableCell>
                  <TableCell className="text-right tabular-nums">
                    {hasAttempts ? `${stats.pass_rate.toFixed(0)}%` : "—"}
                  </TableCell>
                  <TableCell className="text-right tabular-nums">{stats.in_progress_count}</TableCell>
                  <TableCell className="text-right text-xs text-muted-foreground">
                    {formatLastActivity(stats.last_submitted_at)}
                  </TableCell>
                </TableRow>
              );
            })}
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  );
}

"use client";

import { useRouter, useSearchParams } from "next/navigation";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

const PRESETS = [
  { label: "7d", days: 7 },
  { label: "30d", days: 30 },
  { label: "90d", days: 90 },
  { label: "This year", days: null },
] as const;

function toISO(d: Date): string {
  return d.toISOString().split("T")[0];
}

function presetRange(days: number | null): { from: string; to: string } {
  const to = new Date();
  if (days === null) {
    return { from: `${to.getFullYear()}-01-01`, to: toISO(to) };
  }
  const from = new Date(to);
  from.setDate(from.getDate() - days);
  return { from: toISO(from), to: toISO(to) };
}

function isActivePreset(days: number | null, from: string, to: string): boolean {
  const r = presetRange(days);
  return r.from === from && r.to === to;
}

interface DateRangeProps {
  from: string;
  to: string;
}

export function DateRange({ from, to }: DateRangeProps) {
  const router = useRouter();
  const searchParams = useSearchParams();

  function apply(days: number | null) {
    const { from: f, to: t } = presetRange(days);
    const params = new URLSearchParams(searchParams.toString());
    params.set("from", f);
    params.set("to", t);
    router.replace(`?${params.toString()}`);
  }

  return (
    <div className="flex items-center gap-1" role="group" aria-label="Date range presets">
      {PRESETS.map((p) => (
        <Button
          key={p.label}
          size="sm"
          variant={isActivePreset(p.days, from, to) ? "default" : "outline"}
          className={cn("text-xs h-7 px-3")}
          onClick={() => apply(p.days)}
        >
          {p.label}
        </Button>
      ))}
    </div>
  );
}

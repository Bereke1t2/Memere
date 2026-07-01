"use client";

import { useRouter, useSearchParams } from "next/navigation";

const PRESETS = [
  { label: "7d",       days: 7 },
  { label: "30d",      days: 30 },
  { label: "90d",      days: 90 },
  { label: "This year", days: null },
] as const;

function toISO(d: Date): string {
  return d.toISOString().split("T")[0];
}

function presetRange(days: number | null): { from: string; to: string } {
  const to = new Date();
  if (days === null) return { from: `${to.getFullYear()}-01-01`, to: toISO(to) };
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
    <div
      className="inline-flex items-center rounded-lg border bg-card p-0.5 gap-0.5"
      role="group"
      aria-label="Date range presets"
    >
      {PRESETS.map((p) => {
        const active = isActivePreset(p.days, from, to);
        return (
          <button
            key={p.label}
            onClick={() => apply(p.days)}
            className={
              active
                ? "rounded-md bg-primary text-primary-foreground h-7 px-3 text-xs font-medium transition-colors"
                : "rounded-md text-muted-foreground hover:text-foreground hover:bg-muted h-7 px-3 text-xs font-medium transition-colors"
            }
          >
            {p.label}
          </button>
        );
      })}
    </div>
  );
}

# Phase 5 · Skill 4 — Earnings Dashboard (Teacher)

> **Prerequisite:** Skill 3 done.

---

## Goal

Teachers see their earnings at `/earnings`: total gross, their 70% share,
platform fee, and a per-course sales breakdown.

---

## Backend endpoints

| Method | Path | Response |
|--------|------|----------|
| GET | `/me/earnings?from=&to=` | `{ teacher_id, from, to, gross, teacher_share, earnings, platform_fee, units }` |
| GET | `/courses/:id/sales` | `{ course_id, gross, units }` |

`teacher_share` is a decimal string (e.g. `"0.7"` = 70%).

---

## Tasks

### 4.1 — New Zod schemas

`lib/api/schemas.ts`

```ts
export const EarningsSchema = z.object({
  teacher_id: z.string(),
  from: z.string(),
  to: z.string(),
  gross: z.string(),
  teacher_share: z.string(),
  earnings: z.string(),
  platform_fee: z.string(),
  units: z.number(),
});
export type Earnings = z.infer<typeof EarningsSchema>;

export const CourseSalesSchema = z.object({
  course_id: z.string(),
  gross: z.string(),
  units: z.number(),
});
export type CourseSales = z.infer<typeof CourseSalesSchema>;
```

### 4.2 — API endpoint helpers

`lib/api/endpoints.ts`

```ts
export async function getMyEarnings(from: string, to: string): Promise<Earnings>
export async function getCourseSales(courseId: string): Promise<CourseSales>
```

### 4.3 — Route Handler

`app/api/teacher/earnings/route.ts` — GET, proxies `/me/earnings?from=&to=`
`app/api/teacher/courses/[id]/sales/route.ts` — GET, proxies `/courses/:id/sales`

Both verify `user.role === "teacher"` (403 otherwise).

### 4.4 — Earnings page

`app/(dashboard)/earnings/page.tsx`

Server component:
1. `requireStaff()` + teacher-only guard.
2. Resolve date range from searchParams (same 30d default as dashboard).
3. Fetch `getMyEarnings(from, to)` and teacher's courses for the breakdown.
4. Render `<EarningsDashboard>` client component with initial data.

`app/(dashboard)/earnings/loading.tsx` — skeleton

### 4.5 — EarningsDashboard component

`components/earnings/earnings-dashboard.tsx`

Layout:
```
┌─────────────────────────────────────────────────────┐
│ Earnings            [DateRange picker]               │
├──────────┬──────────┬───────────────┬───────────────┤
│ Gross    │ Earnings │ Platform Fee  │ Payments      │
│ Revenue  │ (70%)    │ (30%)         │               │
├─────────────────────────────────────────────────────┤
│ Per-Course Sales                                    │
│ Course title | Gross | Your cut | Payments         │
└─────────────────────────────────────────────────────┘
```

- KPI cards: same style as main dashboard (border, muted icon, no color).
- Per-course table: simple, no pagination needed (teachers typically have
  few courses). Columns: Title, Gross, Earnings (gross × teacher_share),
  Payments.
- Teacher share displayed as "Your cut: 70%".
- All money formatted via `formatMoney(amount, PLATFORM_CURRENCY)`.

---

## Definition of Done

- [ ] Teacher sees earnings KPIs at `/earnings`.
- [ ] Date range changes update all numbers.
- [ ] Per-course breakdown renders with correct math.
- [ ] Admin visiting `/earnings` is redirected to `/revenue`.
- [ ] `pnpm build` + `pnpm typecheck` pass.

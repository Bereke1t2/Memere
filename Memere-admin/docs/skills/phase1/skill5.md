# Phase 1 · Skill 5 — Dashboard Home (KPIs + charts) & Phase 1 smoke test

> **Prerequisite:** Skill 4 done (shell + nav). Read [`docs/skill.md`](../../skill.md)
> §2 (#4 no client business logic, #6 error envelope), [`design spec`](../../Memere_Admin_Design_Specification.md)
> §4.2 (dashboard), §2.2/§2.3 (analytics DTOs), §7 (money/date conventions).

---

## Goal

Replace the dashboard placeholder with the real home screen: KPI cards,
engagement tiles, and a revenue-by-provider chart, all wired to the live
`/admin/analytics/*` endpoints via the server API client — numbers rendered
exactly as the backend returns them. Close Phase 1 with a smoke test.

---

## Tasks

### 5.1 — Date-range control (`components/dashboard/date-range.tsx`)

- Client component: presets (Last 7 / 30 / 90 days, This year) → `from`/`to` ISO.
- Default: last 30 days. Drives the analytics queries via URL search params so the
  Server Component can read them.

### 5.2 — Data loading (Server Component first paint)

- `app/(dashboard)/page.tsx` reads `from`/`to` from `searchParams`, calls
  `getOverview(from, to)`, `getEngagement()`, `getRevenueBreakdown(from, to)`
  server-side, and passes typed data to presentational components.
- Wrap with `<Suspense>` + the skeletons from Skill 4.

### 5.3 — KPI cards (`components/dashboard/kpi-cards.tsx`)

- Cards: Total Students, Total Teachers, Gross Revenue, MRR, Completed Payments,
  Refunded Amount.
- Money via `formatMoney(value, currency)` — **value stays a string from the
  backend**; currency labelled from the data, not hardcoded. Counts via locale
  number formatting.

### 5.4 — Engagement tiles (`components/dashboard/engagement.tsx`)

- avg quiz pass rate, avg exam pass rate, avg completion % → `formatPercent`.
- Small progress bars or stat tiles.

### 5.5 — Revenue chart (`components/charts/revenue-bar.tsx`)

- Recharts bar chart of provider breakdown (`gross` per provider). Parse the money
  string to number **only for plotting**; tooltip shows the formatted string.
- Empty state when no revenue in range.

### 5.6 — Client refresh

- TanStack Query hooks so changing the date range refetches without full reload
  (server action or route revalidation acceptable; keep tokens server-side).

### 5.7 — Phase 1 smoke test (`docs/SMOKE.md` + manual script)

Document and run the end-to-end check:
1. `pnpm build && pnpm start` against a running backend with seeded data.
2. Login as admin → dashboard shows real KPIs.
3. Change date range → numbers update.
4. Engagement + revenue chart render; empty range handled.
5. Logout → `/login`.

---

## Definition of Done

- [ ] Dashboard shows **live** KPIs, engagement, and a revenue chart from the
      backend (no mocked numbers).
- [ ] Every money/percent value is backend-sourced and formatted via the shared
      helpers; currency is read from data, not hardcoded.
- [ ] Changing the date range refetches and updates the figures.
- [ ] Empty/zero ranges render gracefully (no NaN, no crash).
- [ ] `pnpm typecheck` + `pnpm build` + `pnpm lint` pass; smoke test passes.

## Verification commands

```bash
pnpm build && pnpm start
# Manual smoke (docs/SMOKE.md): login → KPIs → change range → logout.
```

---

## 🎉 Phase 1 complete

The admin panel now: runs, authenticates admins with tokens kept entirely
server-side, presents a guarded responsive shell, and shows a live KPI dashboard.
Proceed to **Phase 2** (`docs/prompts/phase2-kickoff.md`) — Users + Courses.

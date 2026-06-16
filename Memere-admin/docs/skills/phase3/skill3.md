# Phase 3 · Skill 3 — Revenue Dashboard

> **Prerequisite:** Skill 2 done. Read [`docs/skill.md`](../../skill.md) §2 (#4),
> §7 (money), [`design spec`](../../Memere_Admin_Design_Specification.md) §4.6,
> §2.2/§2.3 (overview + revenue breakdown DTOs).

---

## Goal

The Revenue screen: header KPIs (gross, refunded, MRR), provider-breakdown charts,
and a breakdown table — every figure straight from the backend, over a selectable
date range.

---

## Tasks

### 3.1 — Endpoints

- Reuse `getOverview(from, to)` for gross/refunded/MRR.
- `getRevenueBreakdown(from, to)` → `RevenueBreakdownItem[]` (`provider, gross,
  units`). Also wire `/admin/revenue` totals if distinct from overview.

### 3.2 — Page (`app/(dashboard)/revenue/page.tsx`)

- Date-range control (reuse the dashboard one) writing to search params.
- KPI tiles: Gross Revenue, Refunded, MRR (string money, currency-labelled).
- Charts (`components/charts/`): revenue-by-provider (bar) and units-by-provider
  (bar or pie). Parse money to number **only for plotting**; tooltips show the
  formatted string.
- Breakdown table: provider, gross, units.

### 3.3 — Empty/zero handling

- A range with no revenue renders a clean empty state (no NaN, no broken chart).

---

## Definition of Done

- [ ] Revenue KPIs + charts + table reflect backend data for the chosen range.
- [ ] All money is backend-sourced strings, formatted via helpers, currency from data.
- [ ] Charts handle empty ranges gracefully.
- [ ] `pnpm typecheck` + `pnpm build` pass.

## Verification commands

```bash
pnpm build && pnpm start
# Manual: change ranges, confirm KPIs + charts update and zero ranges are clean.
```

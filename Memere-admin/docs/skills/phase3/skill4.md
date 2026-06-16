# Phase 3 · Skill 4 — Subscriptions / Financial KPI Overview

> **Prerequisite:** Skill 3 done. Read [`design spec`](../../Memere_Admin_Design_Specification.md)
> §4.6, §2 (analytics endpoints).

---

## Goal

A financial overview surfacing subscription-driven metrics (MRR, completed
payments, refunded amount) and any subscription KPIs the backend exposes via the
analytics/overview endpoint — presented as a focused tiles + trend view. (The
backend has no admin subscription-list endpoint in v1, so this is KPI-level, not a
per-subscription table — note that limitation rather than inventing an endpoint.)

---

## Tasks

### 4.1 — KPI tiles

- From `getOverview`: MRR, completed payments, gross vs refunded ratio.
- Present as a compact financial summary, optionally embedded on the Revenue page
  as a second section or its own `/revenue` tab.

### 4.2 — Trend (optional, if data allows)

- If overview supports ranged queries, show MRR/gross across a few sub-ranges as a
  small line chart. If not feasible from the current API, render the point-in-time
  KPIs and document the gap as a future backend enhancement.

### 4.3 — Honesty about gaps

- Where the backend lacks an endpoint (e.g. list individual subscriptions),
  clearly label the section as KPI-only and do **not** fabricate data.

---

## Definition of Done

- [ ] Subscription/financial KPIs render from backend analytics data.
- [ ] No invented endpoints or client-computed financials; gaps are labelled.
- [ ] `pnpm typecheck` + `pnpm build` pass.

## Verification commands

```bash
pnpm build && pnpm start
# Manual: confirm MRR + payment KPIs match the backend overview values.
```

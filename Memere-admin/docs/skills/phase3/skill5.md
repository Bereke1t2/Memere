# Phase 3 · Skill 5 — Money/Date Formatting Audit + Phase 3 Smoke Test

> **Prerequisite:** Skill 4 done. Read [`docs/skill.md`](../../skill.md) §2 (#4),
> §7.

---

## Goal

Lock down financial correctness: audit every money/date render, ensure all values
are backend-sourced and consistently formatted, then run the Phase-3 smoke test.

---

## Tasks

### 5.1 — Formatting audit

- Grep for any `Number(...)` / arithmetic on money fields used for **display**;
  replace with the string-preserving `formatMoney`. Arithmetic is allowed **only**
  for chart plotting, never for shown totals.
- Confirm currency is always read from the data field, never hardcoded "ETB".
- Confirm all dates use `formatDate` with an ISO tooltip.

### 5.2 — Cross-check against backend

- Spot-check a few payments + the revenue totals in the UI against the backend
  Postman responses to confirm the panel displays exactly what the API returns.

### 5.3 — Phase 3 smoke test (append to `docs/SMOKE.md`)

1. Payments: list → filter by status → detail → refund (confirm) → reconcile.
2. Revenue: KPIs + charts + breakdown table for a range; zero range clean.
3. All money/date renders consistent; numbers match backend.

---

## Definition of Done

- [ ] No client-side money arithmetic for displayed figures; currency from data.
- [ ] UI figures match backend responses on spot-check.
- [ ] Phase-3 smoke test passes.
- [ ] `pnpm typecheck` + `pnpm build` + `pnpm lint` pass.

## Verification commands

```bash
pnpm build && pnpm start
# Run the docs/SMOKE.md Phase-3 checklist; spot-check values vs Postman.
```

---

## 🎉 Phase 3 complete

The financial surface is done. Proceed to **Phase 4**
(`docs/prompts/phase4-kickoff.md`) — Announcements + Hardening + Deployment.

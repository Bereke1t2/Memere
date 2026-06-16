# Phase 3 · Skill 1 — Payments List

> **Prerequisite:** Phase 2 green. Read [`docs/skill.md`](../../skill.md) §2 (#4
> no client business logic, #8 cursor), [`design spec`](../../Memere_Admin_Design_Specification.md)
> §4.5 (Payments), §2.2/§2.3 (AdminPayment DTO; amount is a string).

---

## Goal

The Payments list: a paginated, status-filterable table on `/admin/payments`,
amounts rendered exactly as the backend returns them (string money, labelled
currency).

---

## Tasks

### 1.1 — Endpoint

- `listPayments({ limit, after, status })` → `{ payments: AdminPayment[], next }`.
  Validate with the paginated schema; `amount` stays a **string**.

### 1.2 — Columns (`app/(dashboard)/payments/columns.tsx`)

- Short payment id (monospace, copyable), short student id (link to user detail
  if resolvable), amount + currency (right-aligned, via `formatMoney(amount,
  currency)`), status badge (pending=amber, completed=green, failed=red,
  refunded=slate), provider, created date.

### 1.3 — Page (`app/(dashboard)/payments/page.tsx`)

- Server Component reads `after` + `status` from params; status filter select
  (all/pending/completed/failed/refunded); cursor pagination; row → `/payments/[id]`.

---

## Definition of Done

- [ ] `/payments` lists real payments with status filter + cursor pagination.
- [ ] Amount shown as backend string with correct currency; status badges consistent.
- [ ] `pnpm typecheck` + `pnpm build` pass.

## Verification commands

```bash
pnpm build && pnpm start
# Manual: filter by status, paginate, click a payment row.
```

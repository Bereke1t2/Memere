# Phase 3 · Skill 2 — Payment Detail + Refund + Reconcile

> **Prerequisite:** Skill 1 done. Read [`docs/skill.md`](../../skill.md) §2 (#7
> confirm destructive), [`design spec`](../../Memere_Admin_Design_Specification.md)
> §4.5.

---

## Goal

The payment detail page with the **refund** action (admin-only, irreversible) and
a toolbar **reconcile-pending** action — both awaiting the backend and reporting
the real result.

---

## Tasks

### 2.1 — Endpoints

- `getPayment(id)` → `AdminPayment` (full record).
- `refundPayment(id)` (POST `/payments/:id/refund`) → refunded view.
- `reconcilePayments()` (POST `/admin/payments/reconcile`) → `{ reconciled: number }`.
  All via Route Handler / Server Action (token server-side).

### 2.2 — Detail page (`app/(dashboard)/payments/[id]/page.tsx`)

- Server Component: full payment record (id, student, amount+currency, status,
  provider, provider txn id, created). 404 → friendly not-found.
- **Refund** button shown only when `status === "completed"`. Strong confirm
  dialog ("This refunds <amount> <currency> and cannot be undone"). On success →
  toast + invalidate; surface backend errors faithfully (e.g. already refunded).

### 2.3 — Reconcile action (payments list toolbar)

- A "Reconcile pending" button → `reconcilePayments()`; toast the returned
  `{reconciled}` count; refetch the list.

---

## Definition of Done

- [ ] Payment detail shows the full record; refund only for completed payments.
- [ ] Refund requires a strong confirm, awaits backend, toasts result, refreshes.
- [ ] Reconcile reports the real `{reconciled}` count and refetches.
- [ ] No financial number is computed client-side.
- [ ] `pnpm typecheck` + `pnpm build` pass.

## Verification commands

```bash
pnpm build && pnpm start
# Manual: open a completed payment → refund (confirm) → status flips to refunded;
# run reconcile → count toast; force an error and confirm faithful reporting.
```

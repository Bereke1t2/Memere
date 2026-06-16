# Phase 4 · Skill 2 — Global Hardening (states, toasts, 401→refresh→retry)

> **Prerequisite:** Skill 1 done. Read [`docs/skill.md`](../../skill.md) §2 (#1,
> #6), [`design spec`](../../Memere_Admin_Design_Specification.md) §5 (refresh), §7.

---

## Goal

A consistency pass so every screen handles loading, empty, and error states
uniformly, mutations toast, and an expired access token transparently refreshes
once before forcing re-login.

---

## Tasks

### 2.1 — State coverage audit

- Every list/detail page has: skeleton loading, friendly empty state, and an error
  boundary that maps `ApiError.code` → message (never a raw stack).
- Standardize on shared `<EmptyState>`, `<LoadingTable>`, `<ErrorState>` components.

### 2.2 — 401 → refresh → retry (finalize)

- Confirm the Skill 2/3 refresh hook is wired everywhere `apiFetch` is used: one
  silent refresh on 401, retry once, and on refresh failure clear cookies +
  redirect to `/login` with a "session expired" toast.

### 2.3 — Toasts + mutation UX

- All mutations (suspend, reactivate, role, unpublish, refund, reconcile,
  broadcast) toast success/failure consistently and disable their trigger while
  pending.

### 2.4 — Network resilience

- TanStack Query: sensible `retry` (1) with backoff for idempotent GETs; **no**
  auto-retry on mutations. Stale-while-revalidate for lists.

---

## Definition of Done

- [ ] Every page has loading/empty/error states from shared components.
- [ ] Expired token refreshes once transparently; failed refresh → clean re-login.
- [ ] All mutations toast and disable-while-pending.
- [ ] `pnpm typecheck` + `pnpm build` + `pnpm lint` pass.

## Verification commands

```bash
pnpm build && pnpm start
# Manual: let an access token expire (or shorten TTL) → action still works via
# silent refresh; revoke refresh → forced clean re-login.
```

# Phase 2 · Skill 3 — User Detail + Actions (suspend / reactivate / change role)

> **Prerequisite:** Skill 2 (Users list) done. Read [`docs/skill.md`](../../skill.md)
> §2 (#7 destructive actions confirm), [`design spec`](../../Memere_Admin_Design_Specification.md)
> §4.3.

---

## Goal

The user detail page with the three moderation actions, each behind a confirm
dialog, each awaiting the backend, toasting the real result, and invalidating the
list/detail queries.

---

## Tasks

### 3.1 — Endpoints

- `getUser(id)` → `User`.
- `suspendUser(id, reason)` (POST, 204), `reactivateUser(id)` (POST, 204),
  `changeRole(id, role)` (POST, 204). These go through Route Handlers or Server
  Actions so the token stays server-side.

### 3.2 — Detail page (`app/(dashboard)/users/[id]/page.tsx`)

- Server Component: `getUser(id)`, render full profile (name, email, phone, role,
  status, verified, last login, joined). 404 → friendly not-found.
- Action bar (client component) conditioned on state:
  - **Suspend** (when `is_active`): dialog with a **required reason** textarea.
  - **Reactivate** (when suspended): confirm dialog.
  - **Change role**: dialog with a role select (student/teacher/admin) + confirm.

### 3.3 — Mutations (`components/users/user-actions.tsx`, `"use client"`)

- TanStack Query mutations calling the Route Handlers. On success: toast, close
  dialog, `invalidate` `qk.user(id)` + `qk.users(...)`. On error: toast the
  mapped envelope message, keep the dialog open.
- Guard against acting on yourself for role changes if the backend forbids it
  (surface the backend error faithfully).

---

## Definition of Done

- [ ] Detail page shows the full user; missing id → friendly not-found.
- [ ] Suspend requires a reason; reactivate + change-role confirm before calling.
- [ ] Each action awaits the backend, toasts the real outcome, and refreshes the
      view (no optimistic lie on failure).
- [ ] `pnpm typecheck` + `pnpm build` pass.

## Verification commands

```bash
pnpm build && pnpm start
# Manual: suspend a test user (reason required) → status flips; reactivate;
# change role; force a backend error and confirm the toast + no false success.
```

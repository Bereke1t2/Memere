# Phase 1 · Skill 4 — App Shell (protected layout, sidebar, header, nav)

> **Prerequisite:** Skill 3 done (auth loop works). Read [`docs/skill.md`](../../skill.md)
> §3–4, [`design spec`](../../Memere_Admin_Design_Specification.md) §3 (IA), §7 (UX).

---

## Goal

Build the protected `(dashboard)` shell every admin screen lives inside: a
server-side role-guarded layout, a collapsible sidebar with the nav, a header with
breadcrumb + user menu + theme toggle, and the loading/error boundaries. After
this skill the app has a real "logged-in" frame with placeholder pages for each
section.

---

## Tasks

### 4.1 — Protected layout (`app/(dashboard)/layout.tsx`)

- **Server Component.** Call `getSession()`; if no session or `role !== "admin"`
  → `redirect("/login")`. This is the authoritative role gate (middleware only
  checks cookie presence).
- Render the shell: `<Sidebar />` + `<Header user={...} />` + `<main>{children}</main>`.
- Pass the authenticated user down for the header user-menu.

### 4.2 — Sidebar (`components/layout/sidebar.tsx`)

- Nav items (icon + label) for the full IA, even though later pages are
  placeholders now: Dashboard `/`, Users `/users`, Courses `/courses`, Payments
  `/payments`, Revenue `/revenue`, Announcements `/announcements`.
- Active-route highlight (`usePathname`). Collapsible on tablet widths; icon-only
  collapsed state with tooltips.
- lucide-react icons.

### 4.3 — Header (`components/layout/header.tsx`)

- Breadcrumb derived from the route.
- Theme toggle (light/dark).
- User menu (avatar + name): shows email/role, **Logout** (posts to
  `/api/auth/logout`, then `router.push("/login")`).

### 4.4 — Placeholder section pages

Create minimal Server Component pages so the nav works end-to-end:
`app/(dashboard)/page.tsx` (dashboard — replaced in Skill 5),
`users/page.tsx`, `courses/page.tsx`, `payments/page.tsx`, `revenue/page.tsx`,
`announcements/page.tsx` — each a titled empty state ("Coming in Phase N").

### 4.5 — Boundaries

- `app/(dashboard)/loading.tsx` — skeleton shell.
- `app/(dashboard)/error.tsx` (`"use client"`) — friendly error with a retry
  button; maps `ApiError.code` to a message; never dumps raw stack to the user.

### 4.6 — Responsive + a11y

- Works down to tablet (sidebar collapses to a sheet/drawer on small screens).
- Keyboard-navigable nav; aria-labels on icon-only controls; visible focus rings.

---

## Definition of Done

- [ ] Logged-in admin sees the full shell; every sidebar link routes to its
      placeholder page with active-state highlighting.
- [ ] Direct navigation to a dashboard route while unauthenticated → `/login`;
      while authenticated-but-non-admin (forged cookie) → rejected by `getSession()`.
- [ ] Header user menu shows the admin's name/email/role; Logout works.
- [ ] Theme toggle persists; layout is usable at tablet width.
- [ ] `loading.tsx` and `error.tsx` render; error boundary shows friendly copy.
- [ ] `pnpm typecheck` + `pnpm build` + `pnpm lint` pass.

## Verification commands

```bash
pnpm build && pnpm start
# Manual: navigate every sidebar link, collapse sidebar, toggle theme, logout,
# resize to tablet width, trigger the error boundary (temporary throw) and recover.
```

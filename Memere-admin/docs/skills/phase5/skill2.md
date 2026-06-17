# Phase 5 · Skill 2 — Role-Aware Navigation + Page Guards

> **Prerequisite:** Skill 1 done. Teacher can log in.

---

## Goal

Admins and teachers see different nav items. Admin-only pages redirect
teachers with a friendly "access denied" message instead of a raw error.

---

## Navigation matrix

| Route               | Admin | Teacher |
|---------------------|-------|---------|
| `/`                 | ✓     | ✓ (their own stats) |
| `/users`            | ✓     | ✗ redirect |
| `/courses`          | ✓     | ✗ redirect |
| `/payments`         | ✓     | ✗ redirect |
| `/revenue`          | ✓     | ✗ redirect |
| `/announcements`    | ✓     | ✗ redirect |
| `/my-courses`       | ✗ (redirects to /courses) | ✓ |
| `/earnings`         | ✗ (redirects to /revenue) | ✓ |

---

## Tasks

### 2.1 — Role-aware SidebarNav

`components/layout/sidebar-nav.tsx`

Define two nav configs and select by role:

```ts
const ADMIN_NAV = [
  { href: "/",              label: "Dashboard",     icon: LayoutDashboard },
  { href: "/users",         label: "Users",         icon: Users },
  { href: "/courses",       label: "Courses",       icon: BookOpen },
  { href: "/payments",      label: "Payments",      icon: CreditCard },
  { href: "/revenue",       label: "Revenue",       icon: TrendingUp },
  { href: "/announcements", label: "Announcements", icon: Megaphone },
];

const TEACHER_NAV = [
  { href: "/",            label: "Dashboard",  icon: LayoutDashboard },
  { href: "/my-courses",  label: "My Courses", icon: BookOpen },
  { href: "/earnings",    label: "Earnings",   icon: DollarSign },
];
```

`SidebarNav` receives `role: string` prop. `DashboardShell` passes
`user.role` down through `Sidebar` → `SidebarNav`.

### 2.2 — Sidebar + DashboardShell threading

- `DashboardShell` passes `role={user.role}` to `<Sidebar role={...}>`.
- `Sidebar` passes `role` to `<SidebarNav role={...}>`.
- `MobileSidebar` same threading.

### 2.3 — Admin-only page guards

Each admin-only page (`/users`, `/courses`, `/payments`, `/revenue`,
`/announcements`) adds at the top of the server component:

```ts
const { user } = await requireStaff();
if (user.role !== "admin") redirect("/");
```

Use a shared helper `requireAdmin` (already exists in `lib/auth/session.ts`) —
reuse it directly on these pages. No behaviour change for admins.

### 2.4 — Teacher-only redirects for admins

`app/(dashboard)/my-courses/page.tsx` and `app/(dashboard)/earnings/page.tsx`:

```ts
const { user } = await requireStaff();
if (user.role !== "teacher") redirect(user.role === "admin" ? "/courses" : "/");
```

### 2.5 — Dashboard page: role-aware content

`app/(dashboard)/page.tsx`

- Admin: unchanged (overview KPIs + engagement + revenue chart).
- Teacher: show only engagement stats + a "My Earnings" summary card
  (fetched from `GET /me/earnings`). KPIs that require admin-only data
  (total_students, etc.) are hidden.
- Use `const { user } = await requireStaff()` at the top and branch on `user.role`.

---

## Definition of Done

- [ ] Teacher sees only Dashboard / My Courses / Earnings in sidebar.
- [ ] Admin sees original 6-item nav unchanged.
- [ ] Teacher visiting `/users` or `/courses` gets redirected to `/`.
- [ ] Admin visiting `/my-courses` is redirected to `/courses`.
- [ ] `pnpm build` + `pnpm typecheck` pass.

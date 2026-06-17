# Phase 5 · Skill 1 — Teacher Auth + Role-Aware Session

> **Prerequisite:** Phase 4 complete. Read `docs/skill.md` §2 Non-Negotiables.

---

## Goal

Allow teachers to log in to the admin panel with their own role-scoped view.
Currently the login Route Handler hard-gates on `role === "admin"` — teachers
get a 403. Session helpers must be updated to support both roles cleanly.

---

## Context

Backend already distinguishes roles. Endpoints relevant to this phase:
- `POST /auth/login` → returns `{access_token, refresh_token, expires_in, user}`
- `GET /auth/me` → returns the caller's user object (role embedded in JWT)
- Teacher role string: `"teacher"` (entity.RoleTeacher in Go)

---

## Tasks

### 1.1 — Update login Route Handler

`app/api/auth/login/route.ts`

- Remove the `role !== "admin"` gate that returns 403 NOT_ADMIN.
- Allow `role === "admin"` OR `role === "teacher"`.
- Reject any other role with `{ code: "NOT_STAFF", message: "Only admins and teachers can access this panel." }` and status 403.

### 1.2 — Update login form error copy

`components/auth/login-form.tsx`

- Add `NOT_STAFF` → `"Only admins and teachers can access this panel."` to `CODE_COPY`.
- Change existing `NOT_ADMIN` copy to match new message.

### 1.3 — Session helpers

`lib/auth/session.ts`

Add alongside the existing `requireAdmin()`:

```ts
/** Allows both admin and teacher. Returns the session or redirects to /login. */
export async function requireStaff(): Promise<Session> {
  const session = await getSession();
  if (!session) redirect("/login");
  if (session.user.role !== "admin" && session.user.role !== "teacher")
    redirect("/login");
  return session;
}

/** Returns true when the session user is an admin. */
export function isAdmin(session: Session): boolean {
  return session.user.role === "admin";
}
```

### 1.4 — Dashboard layout uses requireStaff

`app/(dashboard)/layout.tsx`

- Replace `requireAdmin()` → `requireStaff()`.
- Pass the session (including role) down to `DashboardShell`.

### 1.5 — DashboardShell + Sidebar receive role

`components/layout/dashboard-shell.tsx` and `components/layout/sidebar.tsx`

- `DashboardShell` already receives `user: User`. No structural change needed.
- The `user.role` is already available in all child components via the `user` prop passed to `Header`.

### 1.6 — seed-teacher Makefile target

`Memere-backend/cmd/seed-teacher/main.go` + `Makefile`

Mirror `cmd/seed-admin` but with `role = 'teacher'` and defaults:
- `TEACHER_EMAIL=teacher@memere.app`
- `TEACHER_PASSWORD=Teacher1234!`
- `TEACHER_FIRST_NAME=Sample`, `TEACHER_LAST_NAME=Teacher`

---

## Definition of Done

- [ ] `POST /api/auth/login` with teacher credentials returns `{ok: true}` and
      sets `mm_access` + `mm_refresh` cookies.
- [ ] Admin login still works unchanged.
- [ ] Student login returns 403 with `NOT_STAFF` code.
- [ ] `make seed-teacher` creates the teacher account idempotently.
- [ ] `pnpm build` + `pnpm typecheck` pass.

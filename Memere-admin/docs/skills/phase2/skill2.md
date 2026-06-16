# Phase 2 · Skill 2 — Users List

> **Prerequisite:** Skill 1 (DataTable) done. Read [`design spec`](../../Memere_Admin_Design_Specification.md)
> §4.3 (Users), §2.2 (`/admin/users`), §2.3 (User DTO). Skim the Postman "Admin"
> folder.

---

## Goal

The Users list screen: a paginated, searchable, role-filterable table on
`/admin/users`, wired through the server API client and the reusable DataTable.

---

## Tasks

### 2.1 — Endpoint (`lib/api/endpoints.ts`)

- Implement `listUsers({ limit, after, role })` → `{ users: User[], next: string }`
  (match the backend envelope: `{ users, next }`). Validate with the paginated
  schema.

### 2.2 — Columns (`app/(dashboard)/users/columns.tsx`)

- Name (`first_name last_name`), email, role (badge), status (active vs suspended
  badge from `is_active`), joined (`created_at` via `formatDate`), last login
  (`last_login_at`, "—" if null).

### 2.3 — Page (`app/(dashboard)/users/page.tsx`)

- Server Component: read `searchParams` (`after`, `role`, `q`), call `listUsers`,
  render `<DataTable>`.
- Search box filters by email (backend supports `role` filter; if email search
  isn't a backend param, filter client-side within the page and note the
  limitation — do **not** fetch-all to filter).
- Role filter select (all/student/teacher/admin).
- Row click → `/users/[id]`.

### 2.4 — Loading & empty

- `users/loading.tsx` skeleton; empty state ("No users match these filters").

---

## Definition of Done

- [ ] `/users` lists real users with working role filter and cursor pagination.
- [ ] Status + role render as consistent badges; dates via `formatDate`.
- [ ] Row click navigates to the detail route.
- [ ] `pnpm typecheck` + `pnpm build` pass.

## Verification commands

```bash
pnpm build && pnpm start
# Manual: filter by role, paginate with the next cursor, click a row.
```

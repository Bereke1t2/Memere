# Phase 2 · Skill 4 — Courses List + Moderation (unpublish)

> **Prerequisite:** Skill 3 done. Read [`design spec`](../../Memere_Admin_Design_Specification.md)
> §4.4 (Courses), §2.2 (`/admin/courses`, unpublish), §2.3 (Course DTO).

---

## Goal

The Courses moderation list: a paginated table on `/admin/courses` with an
**unpublish-with-reason** action behind a confirm dialog.

---

## Tasks

### 4.1 — Endpoints

- `listCourses({ limit, after })` → `{ courses: Course[], next: string }`.
- `unpublishCourse(id, reason)` (POST, 204) via Route Handler / Server Action.

### 4.2 — Columns + page (`app/(dashboard)/courses/page.tsx`)

- Columns: title, subject, grade, price (or "Free" when `is_free`/price 0,
  currency-labelled), published badge (`is_published`), created date.
- Server Component reads cursor params, renders `<DataTable>`, row click →
  `/courses/[id]`.

### 4.3 — Unpublish action (`components/courses/course-actions.tsx`)

- Shown only for published courses. Confirm dialog with a **required reason**.
- Mutation → toast → invalidate `qk.courses(...)`.

---

## Definition of Done

- [ ] `/courses` lists real courses with cursor pagination and consistent badges.
- [ ] Price renders correctly (Free vs amount + currency from data).
- [ ] Unpublish requires a reason, confirms, awaits backend, toasts, refreshes.
- [ ] `pnpm typecheck` + `pnpm build` pass.

## Verification commands

```bash
pnpm build && pnpm start
# Manual: paginate courses, unpublish a published test course (reason required).
```

# Phase 2 · Skill 1 — Reusable DataTable (cursor pagination)

> **Prerequisite:** Phase 1 green. Read [`docs/skill.md`](../../skill.md) §2 (#8
> cursor pagination), §7; [`design spec`](../../Memere_Admin_Design_Specification.md)
> §7 (tables UX).

---

## Goal

Build one reusable `DataTable` (TanStack Table v8) that every list screen reuses:
typed columns, cursor pagination (`limit` + `after` → `next`), column filters,
sticky header, loading skeleton, and a friendly empty state. Get this right once;
Users/Courses/Payments all consume it.

---

## Tasks

### 1.1 — DataTable component (`components/data-table/data-table.tsx`)

- Generic `DataTable<TData>({ columns, data, isLoading, emptyState })`.
- Sticky header, row hover, optional row `onClick` (navigate to detail).
- Skeleton rows while `isLoading`; empty-state slot when `data.length === 0`.

### 1.2 — Cursor pagination (`components/data-table/pagination.tsx`)

- The backend returns a `next` cursor (opaque string) — **not** page numbers.
- Implement "Next page" / "Previous" by keeping a cursor stack in URL search
  params (`?after=<cursor>`), plus a `limit` selector (10/20/50).
- "Load more" style is also acceptable; pick one and document it. Never fetch-all.

### 1.3 — Filter primitives (`components/data-table/filters.tsx`)

- A search input (debounced) and a select-filter, both writing to URL search
  params so server components can read them and refetch.

### 1.4 — Query integration

- A `useCursorList` hook (TanStack Query) keyed by `(endpoint, params)` that
  drives client-side refetch on filter/cursor change while first paint stays
  server-rendered.

---

## Definition of Done

- [ ] `DataTable` renders typed columns with loading + empty states.
- [ ] Cursor pagination advances via the backend `next` cursor (no offset math).
- [ ] Search + select filters update URL params and refetch.
- [ ] `pnpm typecheck` + `pnpm build` + `pnpm lint` pass.

## Verification commands

```bash
pnpm build
# Exercised fully in skill2 against /admin/users.
```

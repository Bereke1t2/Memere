# Phase 2 · Skill 5 — Course Detail + Phase 2 Polish & Smoke Test

> **Prerequisite:** Skill 4 done. Read [`design spec`](../../Memere_Admin_Design_Specification.md)
> §4.4.

---

## Goal

A read-only course detail view (metadata + sections) for moderation context, plus
a Phase-2 polish pass and smoke test.

---

## Tasks

### 5.1 — Course detail (`app/(dashboard)/courses/[id]/page.tsx`)

- Server Component: fetch course via public `GET /courses/:id` and sections via
  `GET /courses/:id/sections` (read-only inspection). 404 → friendly not-found.
- Show metadata (title, description, subject, grade, price, published, teacher if
  available) + sections list. Include the **Unpublish** action here too.

### 5.2 — Polish

- Consistent badges/spacing across Users + Courses.
- Empty/loading/error states verified on every Phase-2 page.
- Breadcrumbs reflect detail routes (Users / [name], Courses / [title]).

### 5.3 — Phase 2 smoke test (append to `docs/SMOKE.md`)

1. Users: list → filter → paginate → open detail → suspend → reactivate → change
   role (each confirmed, toasted, reflected).
2. Courses: list → paginate → open detail → unpublish (reason required).
3. All destructive actions show the real backend outcome.

---

## Definition of Done

- [ ] Course detail shows metadata + sections; unpublish available + working.
- [ ] Every Phase-2 page has loading/empty/error states.
- [ ] Phase-2 smoke test passes end-to-end against the live backend.
- [ ] `pnpm typecheck` + `pnpm build` + `pnpm lint` pass.

## Verification commands

```bash
pnpm build && pnpm start
# Run the docs/SMOKE.md Phase-2 checklist.
```

---

## 🎉 Phase 2 complete

Users and Courses are fully operable. Proceed to **Phase 3**
(`docs/prompts/phase3-kickoff.md`) — Payments + Revenue + Subscriptions.

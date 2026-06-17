# Phase 5 · Skill 5 — Course Detail + Section Management (Teacher)

> **Prerequisite:** Skill 4 done.

---

## Goal

Teachers can drill into a course from `/my-courses` to see its sections and
add new ones.

---

## Backend endpoints

| Method | Path | Purpose |
|--------|------|---------|
| GET | `/courses/:id` | Course detail |
| GET | `/courses/:id/sections` | List sections |
| POST | `/courses/:id/sections` | Add section |

Add section body: `{ title: string, order?: number }`.

---

## Tasks

### 5.1 — Route Handlers

`app/api/teacher/courses/[id]/sections/route.ts`
- GET: proxy `GET /courses/:id/sections` (no role gate — public endpoint, but
  keep consistent routing through the panel).
- POST: verify teacher role, proxy `POST /courses/:id/sections`.

### 5.2 — Course detail page

`app/(dashboard)/my-courses/[id]/page.tsx`

Server component:
1. `requireStaff()` + teacher-only guard.
2. Fetch `getCourse(id)` + `getCourseSections(id)` in parallel.
3. Render course metadata header + sections list.
4. Include `<TeacherCourseActions>` (publish/edit/delete buttons from skill 3).

`app/(dashboard)/my-courses/[id]/loading.tsx` — skeleton
`app/(dashboard)/my-courses/[id]/not-found.tsx` — not-found page

### 5.3 — Sections list + add section

`components/courses/sections-list.tsx`

- Renders each section as a card row: order number, title, lesson count.
- "Add section" button at bottom opens `<AddSectionDialog>`.

`components/courses/add-section-dialog.tsx`

- Form: Title (required), Order (optional number).
- On submit → POST `/api/teacher/courses/:id/sections` → toast → router.refresh().

### 5.4 — Breadcrumb

`app/(dashboard)/my-courses/[id]/page.tsx`

Use `BreadcrumbSetter` (same pattern as admin course detail) to set the
breadcrumb label to the course title.

---

## Definition of Done

- [ ] Clicking a course row in `/my-courses` navigates to `/my-courses/:id`.
- [ ] Course metadata and sections render.
- [ ] "Add section" dialog creates a section visible after refresh.
- [ ] Breadcrumb shows course title.
- [ ] `pnpm build` + `pnpm typecheck` pass.

---

## Phase 5 complete — Teacher portal is done

All five skills delivered:
1. Auth: teachers can log in, role-gated session helpers.
2. Navigation: role-aware sidebar, admin-only page guards.
3. My Courses: list, create, edit, publish, delete.
4. Earnings: KPIs + per-course sales breakdown.
5. Course Detail: sections list + add section.

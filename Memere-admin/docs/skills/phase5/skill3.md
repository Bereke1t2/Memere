# Phase 5 · Skill 3 — My Courses Page (Teacher)

> **Prerequisite:** Skill 2 done. Teacher sees correct nav.

---

## Goal

Teachers can view, create, edit, publish/unpublish, and delete their own
courses from `/my-courses`. All mutations go via Next.js Route Handlers
(browser never calls the Go API directly).

---

## Backend endpoints used

| Method | Path | Purpose |
|--------|------|---------|
| GET | `/courses?teacher_id=<id>&limit=N&next_cursor=X` | List teacher's courses |
| POST | `/courses` | Create course |
| PUT | `/courses/:id` | Update title/description/price/etc. |
| DELETE | `/courses/:id` | Delete draft course |
| POST | `/courses/:id/publish` | Publish course |

Public `/courses` response shape (differs from admin endpoint):
```json
{ "data": [...], "next_cursor": "uuid-or-empty", "limit": 20 }
```

Course object fields (from GET /courses): `id, teacher_id, title, slug,
description, subject, grade, price, currency, is_free, is_published,
language, level, total_lessons, rating_avg, enrollment_count, created_at,
updated_at`.

Create course required fields: `title`, `description`, `subject`, `grade`
(1–12), `level` (`beginner|intermediate|advanced`).
Optional: `price` (number, default 0 = free), `language` (default "en").

---

## Tasks

### 3.1 — New Zod schemas

`lib/api/schemas.ts`

```ts
export const TeacherCourseListSchema = z.object({
  data: CourseSchema.array(),
  next_cursor: z.string().nullable(),
  limit: z.number(),
});
export type TeacherCourseList = z.infer<typeof TeacherCourseListSchema>;

export const CreateCourseInputSchema = z.object({
  title: z.string().min(3).max(200),
  description: z.string().min(10),
  subject: z.string().min(1),
  grade: z.coerce.number().int().min(1).max(12),
  level: z.enum(["beginner", "intermediate", "advanced"]),
  price: z.coerce.number().min(0).default(0),
  language: z.string().default("en"),
});
export type CreateCourseInput = z.infer<typeof CreateCourseInputSchema>;
```

### 3.2 — API endpoint helpers

`lib/api/endpoints.ts`

```ts
// Teacher: list their own courses via public endpoint
export async function listMyCourses(params: {
  limit?: number;
  next_cursor?: string;
  teacher_id: string;
}): Promise<TeacherCourseList>

// Teacher: create course
export async function createCourse(input: CreateCourseInput): Promise<Course>

// Teacher: update course  
export async function updateCourse(id: string, input: Partial<CreateCourseInput>): Promise<Course>

// Teacher: delete course
export async function deleteCourse(id: string): Promise<void>

// Teacher: publish course
export async function publishCourse(id: string): Promise<void>
```

### 3.3 — Route Handlers

`app/api/teacher/courses/route.ts` — GET (list) + POST (create)
`app/api/teacher/courses/[id]/route.ts` — PUT (update) + DELETE
`app/api/teacher/courses/[id]/publish/route.ts` — POST

All Route Handlers:
1. Read session from httpOnly cookie (same pattern as admin handlers).
2. Verify `user.role === "teacher"` (403 otherwise).
3. Proxy to backend with the access token.

### 3.4 — Page + client component

`app/(dashboard)/my-courses/page.tsx` — server component, fetches initial page
`app/(dashboard)/my-courses/columns.tsx` — table columns: title, subject, grade, status badge, lessons, enrolled, date
`app/(dashboard)/my-courses/my-courses-client.tsx` — DataTable with pagination, "New Course" button
`app/(dashboard)/my-courses/loading.tsx` — skeleton

### 3.5 — Create Course dialog

`components/courses/create-course-dialog.tsx`

- Triggered by "New Course" button on the list page.
- Form fields: Title, Description (textarea), Subject, Grade (1–12 select),
  Level (select), Price (number, 0 = free), Language.
- RHF + Zod (`CreateCourseInputSchema`).
- On submit → POST `/api/teacher/courses` → toast success → router.refresh().

### 3.6 — Course actions (teacher)

`components/courses/teacher-course-actions.tsx`

Actions available per course:
- **Publish** — only on unpublished courses; confirm dialog → POST publish.
- **Edit** — opens edit sheet pre-filled with current values → PUT update.
- **Delete** — only on unpublished/draft courses; confirm dialog → DELETE.

---

## Definition of Done

- [ ] Teacher can see their own courses at `/my-courses`.
- [ ] "New Course" dialog creates a course → appears in list.
- [ ] Publish action makes `is_published` flip to true.
- [ ] Edit updates title/description.
- [ ] Delete removes the course from the list.
- [ ] Non-teacher (admin/student token) → Route Handlers return 403.
- [ ] `pnpm build` + `pnpm typecheck` pass.

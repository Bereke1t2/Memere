# Phase 1 · Skill 4 — Course Service (CRUD + Sections + Lessons)

> **Prerequisite:** Skills 1–3 done. Read [`docs/skill.md`](../../skill.md).
> The course service is where the **IDOR** and **soft-delete** non-negotiables
> get their first real workout.
>
> **Spec references:** `memere_Design_Specification.md` §4.2.2 (courses),
> §4.2.3 (sections & lessons), §4.2.4 (videos metadata), §5.4 (API conventions),
> README "Course Endpoints", §7.2 (teacher owns only their courses).

---

## Goal

Implement the **course** vertical as usecases + repository implementations
(no HTTP yet — Skill 5). Covers: create/list/get/update/soft-delete courses,
manage sections within a course, manage lessons within a section, and assemble
the full nested course view (course → sections → lessons). Video *metadata* rows
may be created here, but upload/transcode/streaming is Phase 3 — only stub the
`videos` row linkage if a lesson is `type=video`.

---

## Authorization model for this skill (from §7.2)

- **Anyone (even unauthenticated)** can list/get **published** courses.
- **Teachers** can create courses and edit/delete **only courses where
  `teacher_id == authenticated user`**.
- **Admins** can edit/delete any course.
- Unpublished courses are visible only to their teacher and admins.

The usecase enforces ownership; never rely on the client-sent IDs for authz.
The actual `userID`/`role` arrives from the JWT middleware (Skill 5) and is
passed into each usecase method as explicit arguments — **the usecase signature
must include the caller's identity** (e.g. `actor Actor{UserID, Role}`).

---

## Tasks

### 4.1 — Repository implementations (`internal/repository/postgres/`)

Implement the course/section/lesson interfaces from Skill 2 over the sqlc queries:

- `course_repository.go` — Create, FindByID, FindBySlug, List (with
  `CourseFilter` + cursor pagination), Update, SoftDelete. Map `23505` on `slug`
  → `apperror.Conflict("SLUG_TAKEN")`; `pgx.ErrNoRows` → `apperror.NotFound`.
- `section_repository.go` — CreateSection, ListSectionsByCourse (ordered by
  `order_index`), UpdateSection, SoftDelete/Delete.
- `lesson_repository.go` — CreateLesson, ListLessonsBySection, ListLessonsByCourse,
  UpdateLesson, SoftDelete.
- `GetCourseWithSectionsAndLessons` — either one query with joins assembled in Go,
  or three queries (course, sections, lessons-by-course) stitched in the repo.
  Prefer the 3-query approach for clarity; assemble into the nested entity.

All list/read queries **must** include `deleted_at IS NULL` and respect
`is_published` visibility based on the caller.

### 4.2 — Cursor pagination (`pkg/pagination`)

If not finished in Skill 2, implement here:

- `EncodeCursor(createdAt time.Time, id uuid.UUID) string` (base64 of a stable
  composite), `DecodeCursor(string) (time.Time, uuid.UUID, error)`.
- Keyset pagination: `WHERE (created_at, id) < ($after_ts, $after_id) ORDER BY
  created_at DESC, id DESC LIMIT $limit+1`. The extra row tells you if there's a
  next page; return `nextCursor` only when it exists.
- Default `limit=20`, max `limit=100` (clamp).

### 4.3 — Course usecases (`internal/usecase/course/`)

A `Service` holding the course/section/lesson repos. Methods (each takes
`ctx` + an `actor` identity):

- `CreateCourse(ctx, actor, CreateCourseInput) (*entity.Course, error)`
  - require `actor.Role` ∈ {teacher, admin}; set `teacher_id = actor.UserID`
    (teachers cannot create on behalf of others).
  - generate `slug` from title (slugify + ensure uniqueness, append short suffix
    on collision); default `currency=ETB`, `language=en`, `is_published=false`,
    counters = 0.
  - validate: title length, price ≥ 0, grade in sane range, subject non-empty,
    level valid.
- `ListCourses(ctx, actor, CourseFilter, cursor, limit) ([]entity.Course, nextCursor, error)`
  - if `actor` is nil/anonymous or a student → force `is_published=true`.
  - if teacher → may also see their own unpublished (filter
    `is_published=true OR teacher_id=actor.UserID`).
  - if admin → no published restriction.
- `GetCourse(ctx, actor, id) (*entity.Course, error)` — nested view via
  `GetCourseWithSectionsAndLessons`; apply the same visibility rule; hide
  unpublished sections/lessons from non-owners.
- `GetCourseBySlug` — same rules.
- `UpdateCourse(ctx, actor, id, UpdateCourseInput)` — load, **assert ownership**
  (teacher owns it, or admin), apply partial updates, save.
- `PublishCourse` / `UnpublishCourse(ctx, actor, id)` — ownership-checked toggle.
- `DeleteCourse(ctx, actor, id)` — ownership-checked **soft delete**.

Sections:
- `AddSection(ctx, actor, courseID, SectionInput)` — assert course ownership;
  `order_index` auto-next if not supplied.
- `ListSections(ctx, actor, courseID)`.
- `UpdateSection`, `DeleteSection` (soft) — ownership via parent course.
- `ReorderSections(ctx, actor, courseID, []sectionID)` — optional but nice.

Lessons:
- `AddLesson(ctx, actor, sectionID, LessonInput)` — assert ownership via section's
  course; validate `type` ∈ {video,note,quiz,mixed}; `is_free_preview` default
  false.
- `ListLessons`, `UpdateLesson`, `DeleteLesson` (soft) — ownership-checked.

### 4.4 — Denormalized counters

The spec keeps `total_lessons`, `total_duration_seconds`, `enrollment_count`,
`rating_avg` on `courses`. For Phase 1:

- Maintain `total_lessons` and `total_duration_seconds` when lessons are
  added/removed/updated (recompute in the same transaction, or increment).
- `enrollment_count` / `rating_avg` stay 0 until Phase 4/5 — leave them.
- Wrap multi-write operations (add lesson + bump course counters) in a **pgx
  transaction** so they stay consistent. Add a `WithTx` helper to the repo layer
  if not present.

### 4.5 — Input validation (`pkg/validator`)

Centralize field validation (lengths, enums, ranges, price ≥ 0, slug charset).
Return `apperror.Validation(details)` with a `details` map of field→message so the
HTTP layer (Skill 5) can surface per-field errors in the standard envelope.

### 4.6 — Unit tests

- `usecase/course` with mocked repos: create-as-teacher sets teacher_id;
  create-as-student rejected; update-not-owner rejected (`FORBIDDEN`);
  list-as-anonymous hides unpublished; delete is soft (repo SoftDelete called,
  not hard delete); add-lesson bumps course counters; pagination returns
  nextCursor only when more rows exist.

---

## Definition of Done

- [ ] `go build ./...` clean; `golangci-lint run` clean.
- [ ] `go test ./internal/usecase/course/... ./pkg/pagination/...` passes.
- [ ] Repos satisfy interfaces (compile-time `var _ repository.CourseRepository =
      …`).
- [ ] Ownership enforced in usecases (teacher cannot edit another's course;
      student cannot create) — proven by tests.
- [ ] Deletes are soft (`deleted_at` set); reads exclude soft-deleted rows.
- [ ] Anonymous/student list+get never returns unpublished content.
- [ ] Lesson add/remove keeps `total_lessons`/`total_duration_seconds` correct,
      inside a transaction.
- [ ] Cursor pagination: stable order, correct `nextCursor`, clamped limit.

## Verification commands

```bash
go build ./...
go test ./internal/usecase/course/... ./pkg/pagination/... -v
golangci-lint run
# Optional integration smoke against the running DB (after Skill 5 wires HTTP):
#   create course as teacher, list as anonymous, confirm unpublished hidden.
```

## Hand-off to Skill 5

Both verticals (auth + course) now exist as tested, HTTP-free business logic.
Skill 5 is the **delivery layer**: Gin handlers, DTOs, the full middleware stack
(auth/JWT, RBAC, CORS, logging, request-id, rate-limit, recovery), the router,
and final dependency wiring in `cmd/api/main.go` — turning everything into a live
API.

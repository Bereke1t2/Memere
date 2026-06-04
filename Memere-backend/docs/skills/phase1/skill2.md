# Phase 1 · Skill 2 — Domain Model, Migrations & sqlc

> **Prerequisite:** Skill 1's Definition of Done passes (service runs, DB+Redis
> connect, `apperror` exists). Read [`docs/skill.md`](../../skill.md) first.
>
> **Spec references:** `memere_Design_Specification.md` §4 (entire Database
> Design section — §4.1 principles, §4.2.1–4.2.8 tables, the ER diagram),
> §5.3 (layer responsibilities — domain has no DB tags).

---

## Goal

Define the **data foundation** for Phase 1: the domain entities (pure Go),
the repository interfaces (contracts), the SQL schema as golang-migrate
migrations, and sqlc wired up to generate type-safe query code. After this skill
the database can be migrated up/down and sqlc generates compiling Go — but no
usecases consume it yet.

**Phase 1 scope of tables:** we create the schema for the *whole* product (so
migrations are coherent), but we only write **entities + repository interfaces +
sqlc queries** for the tables Phase 1 actually uses: `users`, `refresh_tokens`,
`courses`, `course_sections`, `lessons`. The quiz/exam/payment/progress tables
are created in migrations but their Go layer is deferred to later phases.

---

## Design rules carried from the spec (§4.1)

- **UUID PKs** via `gen_random_uuid()` (enable `pgcrypto` or use `gen_random_uuid`
  built into PG13+ — prefer the built-in).
- **Soft deletes:** `deleted_at TIMESTAMPTZ NULL` on user-facing tables.
- **Audit columns:** `created_at`, `updated_at` (`TIMESTAMPTZ NOT NULL DEFAULT now()`).
- **Per-domain schemas:** `auth`, `courses`, `payments`, `progress`. (For the
  monolith we *may* keep everything in `public` to simplify sqlc — **decision:**
  use Postgres schemas as the spec says; prefix tables in sqlc config. If schema
  juggling slows you down, fall back to `public` and note it. Ask if unsure.)
- **JSONB** for flexible metadata columns.
- An `updated_at` auto-touch trigger on every table.

---

## Tasks

### 2.1 — Domain entities (`internal/domain/entity/`)

Pure Go structs. **No `db:` tags, no `json:` tags** that bind to frameworks — the
spec is explicit that domain entities are framework-free (§5.3). (DTOs with json
tags live in the delivery layer in Skill 5; sqlc models live in the repository
layer.) Use `uuid.UUID`, `time.Time`, and pointers / `*string` for nullable
fields.

Create one file per aggregate:

- `user.go` — `User` (id, email, phone, passwordHash, role, firstName, lastName,
  avatarURL, isActive, isEmailVerified, email/password reset token fields,
  lastLoginAt, created/updated/deletedAt) + a `Role` typed string with constants
  `RoleStudent`, `RoleTeacher`, `RoleAdmin`.
- `refresh_token.go` — `RefreshToken` (id, userID, tokenHash, deviceInfo,
  expiresAt, revokedAt, createdAt).
- `course.go` — `Course` (all columns from spec §4.2.2: teacherID, title, slug,
  description, shortDescription, subject, grade, thumbnailURL, price, currency,
  isFree, isPublished, language, level, totalDurationSeconds, totalLessons,
  ratingAvg, enrollmentCount, metadata, audit) + `Level` typed string
  (`beginner`/`intermediate`/`advanced`).
- `section.go` — `CourseSection` (§4.2.3).
- `lesson.go` — `Lesson` (§4.2.3) + `LessonType` typed string
  (`video`/`note`/`quiz`/`mixed`).

Keep enums as typed strings with a `Valid()` method.

### 2.2 — Repository interfaces (`internal/domain/repository/`)

Interfaces only — the contracts the usecases depend on. Implementations come in
Skills 3–4. Each method takes `context.Context` first and returns domain entities
+ error. Define exactly the methods Phase 1 needs:

`user_repository.go`:
```go
type UserRepository interface {
    Create(ctx, *entity.User) error
    FindByID(ctx, id uuid.UUID) (*entity.User, error)
    FindByEmail(ctx, email string) (*entity.User, error)
    Update(ctx, *entity.User) error
    SoftDelete(ctx, id uuid.UUID) error
    SetLastLogin(ctx, id uuid.UUID, t time.Time) error
}
```

`refresh_token_repository.go`: `Create`, `FindByHash`, `Revoke(id)`,
`RevokeAllForUser(userID)`, `DeleteExpired`.

`course_repository.go`: `Create`, `FindByID`, `FindBySlug`,
`List(ctx, filter, cursor, limit)` (returns courses + next cursor),
`Update`, `SoftDelete`, plus section & lesson methods (or split into
`section_repository.go` / `lesson_repository.go` — prefer splitting for clarity):
`CreateSection`, `ListSectionsByCourse`, `CreateLesson`, `ListLessonsBySection`,
`GetCourseWithSectionsAndLessons(ctx, courseID)`.

Define a small `CourseFilter` struct (subject, grade, isPublished, teacherID) and
a `Cursor` type in `pkg/pagination` (Skill 1 created the folder; implement the
encode/decode here or in Skill 4 — do it here so repos can use it).

### 2.3 — Migrations (`migrations/`, golang-migrate format)

Files are `NNNN_description.up.sql` + matching `.down.sql`. Suggested set:

1. `0001_init_extensions_and_schemas` — create schemas `auth, courses, payments,
   progress`; the shared `set_updated_at()` trigger function.
2. `0002_users` — `auth.users` (+ `auth.refresh_tokens`) with all columns,
   indexes (email UNIQUE, phone UNIQUE, role, deleted_at, created_at), and the
   updated_at trigger.
3. `0003_courses` — `courses.courses`, `courses.course_sections`,
   `courses.lessons`, `courses.videos` (metadata table per §4.2.4) with indexes
   from spec §4.2.2.
4. `0004_quizzes` — `courses.quizzes/questions/answers` (§4.2.5) — schema only,
   no Go yet.
5. `0005_exams` — `courses.exams/exam_attempts` (§4.2.6).
6. `0006_enrollments_payments` — `payments.enrollments`, `payments.payments`,
   `payments.coupons` (§4.2.7, §10.3).
7. `0007_progress` — `progress.progress` (§4.2.8).

Every table: UUID PK default `gen_random_uuid()`, audit columns, `deleted_at`
where the spec marks soft-delete, FKs with sensible `ON DELETE` (prefer
`RESTRICT` + soft delete over cascade), and the indexes the spec lists.

**Down migrations must cleanly reverse** (drop in reverse dependency order).

### 2.4 — Wire golang-migrate into the Makefile + a runner

- Implement `make migrate-up` / `make migrate-down` using the `migrate` CLI, or a
  `cmd/migrate/main.go` that embeds migrations via `//go:embed` and runs them with
  the `golang-migrate` library (preferred — no extra binary needed in prod).
- Add `make migrate-create name=...` to scaffold a new pair.
- On `make migrate-up`, migrations run against the `.env` DSN.

### 2.5 — sqlc setup

- `sqlc.yaml` (v2): engine `postgresql`, driver `pgx/v5`, point `schema:` at
  `migrations/`, `queries:` at `db/query/`, output package
  `internal/repository/postgres/sqlcgen` (or `gen`). `emit_pointers_for_null_types:
  true`, `emit_json_tags: false`.
- Write Phase-1 queries in `db/query/`:
  - `users.sql` — CreateUser, GetUserByID, GetUserByEmail, UpdateUser,
    SoftDeleteUser, SetLastLogin (all `:one`/`:exec`).
  - `refresh_tokens.sql` — Create, GetByHash, Revoke, RevokeAllForUser,
    DeleteExpired.
  - `courses.sql` — CreateCourse, GetCourseByID, GetCourseBySlug, ListCourses
    (with keyset/cursor `WHERE (created_at, id) < (...)` pattern + `LIMIT`),
    UpdateCourse, SoftDeleteCourse.
  - `sections.sql`, `lessons.sql` — create + list-by-parent.
  - **Every** read query includes `AND deleted_at IS NULL` where applicable.
- `make sqlc` runs `sqlc generate`. Generated code must compile.

> Note: sqlc + Postgres schemas (`courses.courses`) work but table references
> must be schema-qualified consistently in queries. If it fights you, this is the
> point to fall back to `public` schema — document the choice in a code comment.

---

## Definition of Done

- [ ] `make migrate-up` creates all tables; `make migrate-down` removes them
      cleanly; re-running up is idempotent from a clean DB.
- [ ] `\dt *.*` in psql shows users, refresh_tokens, courses, course_sections,
      lessons, videos, quizzes, questions, answers, exams, exam_attempts,
      enrollments, payments, coupons, progress.
- [ ] All spec-listed indexes exist (verify with `\d auth.users` etc.).
- [ ] `make sqlc` generates code that `go build ./...` accepts.
- [ ] Domain entities have **no** db/json struct tags (grep to confirm).
- [ ] Repository interfaces exist and compile; no implementations yet (that's OK).
- [ ] `updated_at` trigger works: an `UPDATE` bumps `updated_at` automatically.

## Verification commands

```bash
make migrate-up
psql "$DB_DSN" -c '\dt auth.*; \dt courses.*; \dt payments.*; \dt progress.*'
make sqlc
go build ./...
grep -rn 'db:"\|json:"' internal/domain/entity/ && echo "FAIL: tags in domain" || echo "OK: clean domain"
make migrate-down && make migrate-up
```

## Hand-off to Skill 3

Schema + entities + interfaces + sqlc are ready. Skill 3 implements the **auth**
vertical: the postgres+redis repository implementations for users/refresh-tokens,
plus `pkg/password`, `pkg/jwt`, and the auth usecases.

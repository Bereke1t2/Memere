# Phase 1 · Skill 5 — HTTP Delivery, Middleware & Wiring

> **Prerequisite:** Skills 1–4 done (auth + course usecases tested, repos
> implemented, DB migrated). Read [`docs/skill.md`](../../skill.md). This skill
> turns the business logic into a **live, secured API**.
>
> **Spec references:** `memere_Design_Specification.md` §5.4 (API conventions),
> §3.2 (gateway responsibilities — for the monolith, this middleware *is* the
> gateway), §7 (auth/RBAC/rate-limit), README "API Overview" (endpoint tables).

---

## Goal

Expose the auth and course verticals over HTTP with Gin: request/response DTOs,
handlers, the full production middleware stack, the versioned router, and the
final dependency wiring in `cmd/api/main.go`. After this skill, the Phase 1 MVP
backend is **fully runnable end-to-end** — you can register, log in, refresh,
create a course as a teacher, and browse published courses.

---

## API surface for Phase 1 (subset of README tables)

Base path `/api/v1`. Implement exactly these in Phase 1:

**Auth** (`auth_handler.go`)
| Method | Path | Auth | Body |
|---|---|---|---|
| POST | `/auth/register` | public | email, password, first_name, last_name, phone? |
| POST | `/auth/login` | public | email, password |
| POST | `/auth/refresh` | public (refresh token in body) | refresh_token |
| POST | `/auth/logout` | bearer | refresh_token |
| GET | `/auth/me` | bearer | — (returns current user) |

**Courses** (`course_handler.go`)
| Method | Path | Auth |
|---|---|---|
| GET | `/courses` | optional bearer (affects visibility) |
| GET | `/courses/:id` | optional bearer |
| POST | `/courses` | bearer + teacher/admin |
| PUT | `/courses/:id` | bearer + owner/admin |
| DELETE | `/courses/:id` | bearer + owner/admin (soft) |
| POST | `/courses/:id/publish` | bearer + owner/admin |
| POST | `/courses/:id/sections` | bearer + owner/admin |
| GET | `/courses/:id/sections` | optional bearer |
| POST | `/sections/:id/lessons` | bearer + owner/admin |
| GET | `/sections/:id/lessons` | optional bearer |

(Quiz/exam/payment/video/progress endpoints from the README are **later phases** —
do not implement them now.)

---

## Tasks

### 5.1 — DTOs (`internal/delivery/http/dto/` or per-handler)

Request and response structs **with `json` tags** (this is the layer where json
tags belong — not the domain). 

- Request DTOs: `RegisterRequest`, `LoginRequest`, `RefreshRequest`,
  `CreateCourseRequest`, `UpdateCourseRequest`, `CreateSectionRequest`,
  `CreateLessonRequest`, plus list query params (`limit`, `after`, `subject`,
  `grade`).
- Response DTOs: `UserResponse` (**no** password hash / tokens / reset fields),
  `AuthResponse` (`access_token`, `refresh_token`, `expires_in`, `user`),
  `CourseResponse`, `SectionResponse`, `LessonResponse`,
  `CourseDetailResponse` (nested), and a generic paginated envelope
  `{ "data": [...], "next_cursor": "...", "limit": N }`.
- Mapper functions `entity → response DTO` (the only place the hash could leak —
  make it structurally impossible: `UserResponse` simply has no such field).

### 5.2 — Standard response helpers

- `respondJSON(c, status, payload)` and `respondError(c, err)`.
- `respondError` type-switches on `*apperror.AppError` → uses its `HTTPStatus`,
  `Code`, `Message`, `Details`. Any non-AppError → log it + return generic
  `500 INTERNAL` (never leak internals to the client — §7.3).
- Always emit the spec envelope: `{ "code", "message", "details" }` on errors.

### 5.3 — Middleware stack (`internal/delivery/middleware/`)

Order matters. Apply globally in this order:

1. **Recovery** — recover from panics → `500 INTERNAL` (replace the Skill 1 stub
   with one that uses the error envelope and logs the stack).
2. **RequestID** — generate/propagate `X-Request-ID`; put it in the context and
   response header.
3. **Logger** — structured access log (method, path, status, latency,
   request-id). **Never log Authorization headers, bodies with passwords, or
   tokens** (§ non-negotiable #8).
4. **CORS** — configurable allowed origins (from config; default permissive in
   dev, locked in prod).
5. **RateLimit** — Redis-backed fixed/sliding window. Global sane default; a
   **stricter limiter on `/auth/login`: 5 attempts / 15 min per IP** (spec §7.3).
   Return `429 TOO_MANY_REQUESTS` with `Retry-After`.

Per-route middleware:

6. **Auth (JWT)** — `RequireAuth`: parse `Authorization: Bearer`, `jwt.Verify`
   (type must be "access"), load claims into context (`actor` = userID+role).
   On failure → `401 UNAUTHORIZED`.
7. **OptionalAuth** — same but doesn't fail when no token; sets actor if present.
   Used on `GET /courses*` so visibility adapts.
8. **RequireRole(roles...)** — checks the context actor's role; `403 FORBIDDEN`
   otherwise. (Ownership checks stay in the usecase — middleware only gates by
   role.)

Provide a context helper `actorFromContext(c) (course.Actor, bool)` reused by all
handlers.

### 5.4 — Handlers

- `auth_handler.go` — bind+validate request DTO, call the auth usecase, map to
  response DTO, set status (`201` for register, `200` for login/refresh).
  `/auth/me` reads the actor, loads the user, returns `UserResponse`.
- `course_handler.go` — bind, pull `actor` (or anonymous) from context, call the
  course usecase, map results. Parse pagination params; clamp limit. Return the
  paginated envelope for list endpoints. `:id` params parsed as UUID with a clean
  `400 INVALID_ID` on parse failure.
- Validation errors from binding → `apperror.Validation(details)` → `422`/`400`
  per your convention (pick one; spec uses generic error envelope — use `400`
  with `code=VALIDATION_ERROR` and per-field `details`).

### 5.5 — Router (`internal/delivery/http/router.go`)

- A `NewRouter(deps) *gin.Engine` that: applies global middleware, mounts
  `/health` (from Skill 1), groups `/api/v1`, and registers auth + course routes
  with the correct per-route middleware.
- Keep route registration declarative and grouped by resource.

### 5.6 — Dependency wiring (`cmd/api/main.go`)

Replace Skill 1's minimal bootstrap with full wiring:

1. `config.Load()`
2. connect Postgres (pgxpool) + Redis
3. build sqlc `Queries` from the pool
4. construct repositories (postgres + redis impls)
5. construct `jwt.Manager`
6. construct usecases (auth, course) injecting repos + jwt + config
7. construct handlers injecting usecases
8. build the router with all deps + middleware
9. `http.Server` with timeouts + graceful shutdown (keep from Skill 1)

Keep `main` thin — consider a `wire.go`/`app.go` builder function if it grows, but
**no DI framework**; plain constructor wiring (matches spec's explicit-wiring
intent).

### 5.7 — End-to-end smoke test

Add a `scripts/smoke.sh` (curl-based) or a Go integration test that, against a
running stack:

1. register a teacher, login → capture access+refresh.
2. `GET /auth/me` with the token → 200.
3. create a course → 201; it's unpublished.
4. `GET /courses` **without** token → course NOT listed (unpublished hidden).
5. publish it → `GET /courses` anonymous now lists it.
6. add a section, add a lesson → counters update.
7. refresh the access token → 200, new token works, old refresh revoked.
8. login 6× with wrong password → 6th returns `429`.

### 5.8 — Update docs

- Update root `README.md` "Getting Started" if any command changed.
- Optionally generate/commit an OpenAPI stub under `api/` documenting the Phase 1
  endpoints (nice-to-have; not blocking).

---

## Definition of Done

- [ ] `go build ./...` clean; `golangci-lint run` clean; `go test ./...` passes.
- [ ] `make up && make migrate-up && make run` brings up a working API.
- [ ] The §5.7 smoke flow passes end to end (script committed).
- [ ] Error responses always use the `{code,message,details}` envelope; 500s
      never leak internal errors.
- [ ] Login rate limit (5 / 15 min / IP) returns `429` with `Retry-After`.
- [ ] No log line ever contains a password, token, or Authorization header
      (grep the logger; review).
- [ ] Unpublished courses are invisible to anonymous/student callers via the API.
- [ ] `UserResponse` cannot serialize a password hash (structurally absent).

## Verification commands

```bash
make up && make migrate-up
make run &
bash scripts/smoke.sh          # the §5.7 flow; all assertions green
go test ./... && golangci-lint run
```

---

## 🎉 Phase 1 complete — what now

When this skill's Definition of Done passes, the Phase 1 MVP backend is live:
config, DB+Redis, migrations, auth (JWT+refresh+RBAC+rate-limit), and the course
service (CRUD + sections + lessons) — all behind a secured Gin API following the
clean-architecture and non-negotiable rules.

**To proceed to Phase 2:**
1. Report "Phase 1 complete" to the user with a one-paragraph summary of what
   runs and the smoke-test result.
2. Ask the user (Claude) to author the **Phase 2** skill files into
   `docs/skills/phase2/` — these will cover the **Quiz & Exam engines**
   (server-side grading, server-enforced timers, Redis attempt state,
   randomization, scoring & analytics) per spec §9.
3. Do **not** scaffold Phase 2 code before its skills are written and reviewed.

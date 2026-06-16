# Phase 1 · Skill 2 — API Layer (server client, Zod DTOs, typed endpoints)

> **Prerequisite:** Skill 1 done (app scaffolds, builds, themed). Read
> [`docs/skill.md`](../../skill.md) §2 (#1, #2, #4, #6, #8), §3, §7.
>
> **Spec references:** design spec §2 (backend contract), §2.2 (endpoints), §2.3
> (DTO shapes), §5 (auth/session — for the refresh hook this layer calls).

---

## Goal

Build the **server-only** API layer: a fetch client that attaches the Bearer
token from the httpOnly cookie, maps the backend's error envelope to typed errors,
mirrors backend DTOs as Zod schemas, and exposes typed endpoint functions. This is
the single seam between the panel and the Go API. **No React here** — pure TS.

This skill writes the client; Skill 3 sets the cookies it reads. Until then, the
client can be unit-exercised against the backend with a token passed in tests.

---

## Tasks

### 2.1 — Zod DTO schemas (`lib/api/schemas.ts`)

Mirror the backend DTOs from design spec §2.3 exactly — no extra fields:

- `UserSchema`, `AuthResponseSchema`, `AdminPaymentSchema`, `OverviewSchema`,
  `RevenueBreakdownItemSchema`, `EngagementStatsSchema`, `CourseSchema`.
- A generic `PaginatedSchema<T>` helper for `{ data/items, next }` list shapes
  (match the backend's actual envelope: e.g. `{ users, next }`, `{ courses, next }`,
  `{ payments, next }` — verify against the Postman collection and model each).
- An `ErrorEnvelopeSchema` = `{ code: string, message: string, details?: unknown }`.
- Export inferred TS types (`export type User = z.infer<typeof UserSchema>`).

> Bind defensively: mark fields the backend may omit as `.optional()`. Money
> fields (`amount`, `gross_revenue`, `gross`, `mrr`, `refunded_amount`) are
> **strings** — schema them as `z.string()` (or `z.coerce.string()`), never
> `z.number()`.

### 2.2 — Error type + envelope mapping (`lib/api/errors.ts`)

- `class ApiError extends Error { code: string; status: number; details?: unknown }`.
- A `CODE_MESSAGES` dictionary mapping known backend codes to friendly copy
  (`INVALID_CREDENTIALS`, `RESOURCE_NOT_FOUND`, `FORBIDDEN`,
  `IDEMPOTENCY_KEY_REQUIRED`, `VALIDATION_ERROR`, …). `friendlyMessage(err)` falls
  back to `err.message`.

### 2.3 — Server fetch client (`lib/api/server.ts`)

`import "server-only";` at the top (build fails if a client component imports it).

- `apiFetch<T>(path, { method, body, schema, token? })`:
  - Builds `${API_BASE_URL}/api/v1${path}`.
  - Attaches `Authorization: Bearer <token>` where token comes from the httpOnly
    cookie (via a `getAccessToken()` helper; Skill 3 implements the cookie read —
    stub it to read `cookies()` now).
  - Sets `Content-Type: application/json` for bodies; `cache: "no-store"` for
    dynamic admin data.
  - On non-2xx: parse the error envelope, throw `ApiError` (preserve `status` +
    `code`). On 204: return `undefined`.
  - On 2xx with a `schema`: `schema.parse(json)` and return typed data.
- **401 refresh hook (single retry):** on 401, call the refresh routine
  (`refreshAccessToken()` — stubbed in this skill, implemented in Skill 3), then
  retry once. If still 401, throw a `SESSION_EXPIRED` `ApiError`.

### 2.4 — Typed endpoint functions (`lib/api/endpoints.ts`)

Thin, typed wrappers used by pages/handlers. Each calls `apiFetch` with the right
schema. Author the ones Phase 1 needs now; stub signatures for later phases:

```ts
// Phase 1
login(email, password): Promise<AuthResponse>
me(): Promise<User>
getOverview(from, to): Promise<Overview>
getEngagement(): Promise<EngagementStats>
getRevenueBreakdown(from, to): Promise<RevenueBreakdownItem[]>
// declared now, bodies filled in later phases:
listUsers(params), getUser(id), suspendUser(id, reason), reactivateUser(id),
changeRole(id, role), listCourses(params), unpublishCourse(id, reason),
listPayments(params), getPayment(id), refundPayment(id), reconcilePayments(),
broadcast(input)
```

### 2.5 — Query keys + provider helpers (`lib/query/keys.ts`)

- Centralized query-key factory (`qk.overview(range)`, `qk.users(params)`, …) so
  invalidation is consistent across mutations.

---

## Definition of Done

- [ ] `lib/api/schemas.ts` mirrors every design-spec §2.3 DTO; money fields are strings.
- [ ] `lib/api/server.ts` is `server-only`, attaches Bearer from cookie, maps the
      error envelope to `ApiError`, handles 204, and retries once on 401.
- [ ] `lib/api/endpoints.ts` exposes typed Phase-1 functions; later ones declared.
- [ ] A `"use client"` component importing `lib/api/server` fails the build
      (proves the server-only boundary).
- [ ] `pnpm typecheck` + `pnpm build` pass.

## Verification commands

```bash
pnpm typecheck
pnpm build
# Optional: a temporary script calling getOverview() with a real admin token
# against a running backend prints typed data, then is removed.
```

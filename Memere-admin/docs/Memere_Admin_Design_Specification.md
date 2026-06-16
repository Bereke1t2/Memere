# Memere Admin Web Panel — Design Specification

> **Source of truth for the Memere Admin project.** The orchestration map is
> [`docs/skill.md`](./skill.md). The backend contract is the Go API in
> `Memere-backend` and its
> [`Memere.postman_collection.json`](../../Memere-backend/docs/Memere.postman_collection.json).
> Where this document and a skill file disagree, this document wins.

---

## 1. Purpose & scope

### 1.1 What this is

A **web-based administration console** for Memere (ExamPrep). Platform operators
(admins) use it from a desktop browser to:

- Monitor platform health and KPIs (users, revenue, engagement).
- Manage users (search, inspect, suspend/reactivate, change roles).
- Moderate content (review courses, unpublish violating content).
- Operate payments (inspect, refund, reconcile pending).
- View revenue and subscription financials.
- Broadcast announcements to user segments.

### 1.2 What this is NOT

- **Not a backend.** It owns no database and no business logic. It is a client of
  the existing Go API.
- **Not for students or teachers.** They use the Flutter app (`memere_mobile`).
  Teacher-facing course-authoring is out of scope (teachers author on mobile/their
  own surfaces; admins only moderate).
- **Not a place where money or grades are computed.** Every financial and
  analytical number is rendered exactly as the backend returns it.

### 1.3 Primary persona

**Platform Admin** — operations/support staff. Trusts the data, needs speed,
clarity, and safe destructive actions. Works on a laptop (1280px+), occasionally a
tablet. Not necessarily technical; the UI must be self-explanatory.

---

## 2. The backend contract (what we consume)

Base URL: `<API_BASE_URL>/api/v1` (server-side env, e.g.
`https://api.memere.app/api/v1`). Auth: `Authorization: Bearer <jwt>`. Errors:
`{ "code": "SNAKE_CASE", "message": "...", "details": {} }`. Lists: cursor
pagination via `limit` + `after`, responses carry a `next` cursor.

### 2.1 Auth endpoints

| Method | Path | Use |
|---|---|---|
| POST | `/auth/login` | Exchange email+password → `{ access_token, refresh_token, expires_in, user }`. Admin panel rejects non-admin `user.role`. |
| POST | `/auth/refresh` | Rotate access token using refresh token. |
| POST | `/auth/logout` | Invalidate the refresh token / session. |
| GET | `/auth/me` | Current user; used to re-verify role on the server. |

Access token TTL ≈ 15 min; refresh ≈ 30 days. The panel transparently refreshes
on a 401 (once) before forcing re-login.

### 2.2 Admin endpoints (the panel's working set)

**Analytics**
| Method | Path | Returns |
|---|---|---|
| GET | `/admin/analytics/overview?from&to` | `total_students, total_teachers, total_admins, gross_revenue, refunded_amount, completed_payments, mrr, from, to` |
| GET | `/admin/analytics/revenue?from&to` | `breakdown[]`: `{ provider, gross, units }` |
| GET | `/admin/analytics/engagement` | `avg_quiz_pass_rate, avg_exam_pass_rate, avg_completion_pct` |
| GET | `/admin/revenue?from&to` | platform revenue totals (admin-only) |

**Users**
| Method | Path | Action |
|---|---|---|
| GET | `/admin/users?limit&after&role` | paginated users + `next` |
| GET | `/admin/users/:id` | one user |
| POST | `/admin/users/:id/suspend` `{reason}` | suspend (204) |
| POST | `/admin/users/:id/reactivate` | reactivate (204) |
| POST | `/admin/users/:id/role` `{role}` | change role (204) |

**Courses**
| Method | Path | Action |
|---|---|---|
| GET | `/admin/courses?limit&after` | paginated courses + `next` |
| POST | `/admin/courses/:id/unpublish` `{reason}` | unpublish (204) |
| GET | `/courses/:id` (public) | course detail for inspection |
| GET | `/courses/:id/sections` (public) | sections |

**Payments**
| Method | Path | Action |
|---|---|---|
| GET | `/admin/payments?limit&after&status` | paginated payments + `next` |
| GET | `/admin/payments/:id` | one payment |
| POST | `/admin/payments/reconcile` | reconcile pending → `{reconciled}` |
| POST | `/payments/:id/refund` (admin) | refund a completed payment |

**Announcements**
| Method | Path | Action |
|---|---|---|
| POST | `/admin/announcements` `{title, body, segment, data}` | broadcast (204). `segment ∈ {all, students, teachers, subscribers}` |

### 2.3 DTO shapes (mirror these in Zod/TS)

These are the exact response fields the panel binds to (from the backend `dto`
package). Mirror them in `lib/api/schemas.ts`; do not add fields the backend
doesn't send.

- **User**: `id, email, phone?, role, first_name, last_name, avatar_url?,
  is_active, is_email_verified, last_login_at?, created_at, updated_at`
- **AuthResponse**: `access_token, refresh_token, expires_in, user`
- **AdminPayment**: `id, student_id, amount (string), currency, status, provider,
  provider_txn_id?, created_at`
- **Overview**: `total_students, total_teachers, total_admins, gross_revenue
  (string/decimal), refunded_amount, completed_payments, mrr, from, to`
- **RevenueBreakdownItem**: `provider, gross (decimal), units`
- **EngagementStats**: `avg_quiz_pass_rate, avg_exam_pass_rate,
  avg_completion_pct`
- **Course**: `id, title, description, subject, grade, price, currency, is_free,
  is_published, …` (bind defensively — treat unknown fields as optional)

> **Money fields arrive as strings/decimals.** Keep them as strings for display;
> never coerce with `Number()` for anything other than chart plotting, and label
> the currency from the field, not a hardcoded "ETB".

---

## 3. Information architecture (screens)

```
/login                         Public. Email + password. Admin-only gate.
/(dashboard)                   Protected shell (sidebar + header).
  ├── /                        Dashboard home — KPI cards + charts
  ├── /users                   Users table (search, role filter, paginate)
  │     └── /users/[id]        User detail + suspend/reactivate/role actions
  ├── /courses                 Courses table (filter, paginate)
  │     └── /courses/[id]      Course detail + unpublish
  ├── /payments                Payments table (status filter, paginate)
  │     └── /payments/[id]     Payment detail + refund
  ├── /revenue                 Revenue & financials (charts, breakdown, date range)
  └── /announcements           Broadcast composer + (optional) history
```

Navigation: persistent left sidebar (collapsible on tablet), top header with
breadcrumb, global search (users) optional, user menu (logout, theme).

---

## 4. Screen specifications

### 4.1 Login (`/login`)
- Email + password form (RHF + Zod). Submit → `POST /api/auth/login` (Route
  Handler) which calls the backend, verifies `user.role === "admin"`, sets
  httpOnly cookies, returns `{ ok }`. Non-admin → 403 with a clear message
  ("This account is not an administrator").
- Loading state on submit; field + form-level error display from the error
  envelope (`INVALID_CREDENTIALS` → "Email or password is incorrect").
- Already-authenticated admins visiting `/login` are redirected to `/`.

### 4.2 Dashboard home (`/`)
- **KPI cards**: Total Students, Total Teachers, Gross Revenue, MRR, Completed
  Payments, Refunded Amount — from `/admin/analytics/overview`.
- **Engagement tiles**: avg quiz pass rate, avg exam pass rate, avg completion %
  — from `/admin/analytics/engagement` (render as percentages).
- **Revenue chart**: provider breakdown (bar) from `/admin/analytics/revenue`.
- **Date range** control (default: last 30 days) feeding `from`/`to`.
- Server-rendered first paint; client refresh on range change.

### 4.3 Users (`/users`, `/users/[id]`)
- Table columns: name (first+last), email, role (badge), status (active/suspended
  badge), joined date, last login. Search box (by email), role filter
  (all/student/teacher/admin), cursor pagination.
- Row → detail page. Detail shows full profile + action buttons:
  - **Suspend** (active users) — confirm dialog with required reason.
  - **Reactivate** (suspended users) — confirm dialog.
  - **Change role** — select dialog (student/teacher/admin) with confirm.
  - Each action: optimistic-safe (await backend), toast result, invalidate query.

### 4.4 Courses (`/courses`, `/courses/[id]`)
- Table: title, subject, grade, price (or "Free"), published badge, created date.
  Cursor pagination.
- Detail: course metadata + sections list (read-only inspection). **Unpublish**
  action (published courses) — confirm dialog with required reason.

### 4.5 Payments (`/payments`, `/payments/[id]`)
- Table: payment id (short), student id (short/link), amount + currency, status
  badge, provider, created date. Status filter (pending/completed/failed/
  refunded). Cursor pagination.
- Toolbar action: **Reconcile pending** → `POST /admin/payments/reconcile`, toast
  the `{reconciled}` count, refetch.
- Detail: full payment record. **Refund** action (completed payments only) —
  confirm dialog, toast result.

### 4.6 Revenue (`/revenue`)
- Header KPIs: gross revenue, refunded, MRR (from overview).
- Charts: revenue-by-provider (bar), units-by-provider (bar/pie). Date range.
- Table: provider breakdown (`provider, gross, units`).

### 4.7 Announcements (`/announcements`)
- Composer form: title, body (textarea), segment (all/students/teachers/
  subscribers), optional key/value `data`. Submit → `POST /admin/announcements`,
  toast success. Confirmation before sending to "all".
- Optional: local history of sends this session (backend has no list endpoint).

---

## 5. Auth & session model

1. **Login**: browser posts credentials to `/api/auth/login` (Route Handler). The
   handler calls the Go `/auth/login`, checks `user.role === "admin"`, and on
   success sets two httpOnly cookies: `mm_access` (access token, ~15m) and
   `mm_refresh` (refresh token, ~30d), both `Secure` + `SameSite=Lax`. Returns no
   token to the browser.
2. **Authenticated requests**: Server Components and Route Handlers read `mm_access`
   from cookies and attach it as a Bearer to the Go API call.
3. **Refresh**: on a 401 from the Go API, the server API client calls
   `/auth/refresh` with `mm_refresh`, updates `mm_access`, and retries once. If
   refresh fails, cookies are cleared and the user is sent to `/login`.
4. **`middleware.ts`**: edge gate — requests to `(dashboard)` without an `mm_access`
   (or `mm_refresh`) cookie redirect to `/login`. Role is re-verified server-side
   in the dashboard layout via `/auth/me` (defense in depth; a forged cookie can't
   pass the backend).
5. **Logout**: `/api/auth/logout` calls the Go `/auth/logout` and clears both
   cookies.

---

## 6. Security requirements

- httpOnly + Secure + SameSite cookies; no token in JS-readable storage.
- `API_BASE_URL` and any secret are server-only env (never `NEXT_PUBLIC_`).
- CSRF: mutations go through same-site Route Handlers; use SameSite=Lax and a
  per-session CSRF token for state-changing POSTs (Phase 4 hardening).
- CSP + security headers (`X-Frame-Options`, `X-Content-Type-Options`,
  `Referrer-Policy`) via `next.config.ts` headers (Phase 4).
- Confirm dialogs on every destructive action.
- Error envelope mapping; never leak raw backend internals to the UI.
- Dependency scanning in CI before deploy.

---

## 7. UX & visual conventions

- **Design system**: shadcn/ui + Tailwind tokens. Neutral, data-dense, clean.
  Light + dark mode from Phase 1.
- **Layout**: fixed left sidebar (icons + labels), top header (breadcrumb + user
  menu). Content max-width comfortable for tables.
- **Status badges**: consistent color language — active/completed = green,
  pending = amber, suspended/failed = red, refunded = slate, draft/unpublished =
  muted.
- **Tables**: sticky header, zebra optional, row hover, click-through, skeleton on
  load, friendly empty state, cursor "Load more" / next-page control.
- **Forms**: inline validation, disabled submit while pending, success/error
  toasts.
- **Money**: right-aligned, currency-labelled, grouped thousands.
- **Dates**: human format (e.g. "15 Jun 2026, 14:32") with raw ISO in a tooltip.
- **Accessibility**: focus rings, keyboard nav, aria labels on icon buttons.

---

## 8. Non-goals (explicitly out of scope for v1)

- Course/quiz/exam **authoring** (teachers do this; admins only moderate).
- Video upload/transcode controls (backend worker concern).
- Real-time updates / websockets (poll/refetch is enough at admin scale).
- Multi-language UI (English only for v1; Amharic is a later roadmap item).
- Fine-grained RBAC within admin (single "admin" role for v1).

---

## 9. Environment contract

| Var | Scope | Example | Notes |
|---|---|---|---|
| `API_BASE_URL` | server | `https://api.memere.app` | Go backend origin; `/api/v1` appended by the client. |
| `COOKIE_SECRET` | server | 32+ random bytes | Signs/verifies session cookie integrity if used. |
| `NODE_ENV` | server | `production` | Standard. |
| `NEXT_PUBLIC_APP_NAME` | public | `Memere Admin` | Cosmetic only — never a secret. |

The backend must include the admin panel's origin in its `CORS_ALLOWED_ORIGINS`
**only** if any direct browser→API calls are made; in this design all API calls
are server-side, so CORS is not strictly required — but coordinate at deploy time
(Phase 4 skill 5).

---

## 10. Definition of "done" for the product (v1)

- An admin can log in, see real platform KPIs, manage users, moderate courses,
  operate payments/refunds, view revenue, and send announcements — all against the
  live backend, with tokens never exposed to the browser, destructive actions
  confirmed, and every financial/grading number sourced from the backend.

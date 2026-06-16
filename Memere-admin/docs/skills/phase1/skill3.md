# Phase 1 · Skill 3 — Auth (login, cookies, session, middleware, admin gate)

> **Prerequisite:** Skill 2 done (API layer with stubbed cookie/refresh hooks).
> Read [`docs/skill.md`](../../skill.md) §2 (#1–#3), [`design spec`](../../Memere_Admin_Design_Specification.md)
> §4.1 (login), §5 (auth/session model), §6 (security).

---

## Goal

Implement the full auth loop with **tokens never reaching the browser**: a login
page, Route Handlers that set/clear/rotate httpOnly cookies, server-side session
helpers, the admin-only role gate, and `middleware.ts`. Replace the cookie/refresh
stubs from Skill 2 with real implementations.

---

## Tasks

### 3.1 — Cookie helpers (`lib/auth/cookies.ts`)

- Constants: `ACCESS_COOKIE = "mm_access"`, `REFRESH_COOKIE = "mm_refresh"`.
- `setAuthCookies(res, { access, refresh, expiresIn })` — both **httpOnly**,
  **Secure** (in prod), **SameSite=Lax**, `path=/`. Access cookie maxAge ≈
  `expiresIn`; refresh ≈ 30 days.
- `clearAuthCookies(res)`.
- `getAccessToken()` / `getRefreshToken()` — server-side reads via `cookies()`.

### 3.2 — Login Route Handler (`app/api/auth/login/route.ts`)

- `POST { email, password }`:
  1. Call backend `login(email, password)` (Skill 2 endpoint).
  2. **Verify `user.role === "admin"`.** If not → 403 JSON
     `{ code: "NOT_ADMIN", message: "This account is not an administrator." }` and
     set **no** cookies.
  3. On success → set httpOnly cookies, return `{ ok: true }` (no tokens in body).
  4. Map backend `INVALID_CREDENTIALS` → 401 with friendly message.

### 3.3 — Logout + Refresh handlers

- `app/api/auth/logout/route.ts`: call backend `/auth/logout` with refresh token,
  clear cookies, return `{ ok: true }`.
- `app/api/auth/refresh/route.ts` and a server util `refreshAccessToken()`:
  exchange `mm_refresh` → new access token via backend `/auth/refresh`, update the
  access cookie. Wire this into the Skill 2 `apiFetch` 401-retry hook. On failure,
  clear cookies and signal session expiry.

### 3.4 — Session helpers (`lib/auth/session.ts`)

- `getSession()` — server-only: returns `{ user }` by calling `me()` with the
  access token, or `null`. Used by the dashboard layout to re-verify role
  (defense in depth — a forged cookie can't pass the backend).
- `requireAdmin()` — returns the user or `redirect("/login")` / `redirect("/")`
  with a forbidden flag if not admin.

### 3.5 — Login page (`app/login/page.tsx`)

- RHF + Zod form (email, password). Submits to `/api/auth/login`.
- Loading state, field + form errors from the envelope, success → `router.push("/")`.
- If already authenticated (session present), redirect to `/` (check server-side).
- Branded, centered card; light/dark aware.

### 3.6 — Edge middleware (`middleware.ts`)

- `matcher` covers the `(dashboard)` routes (everything except `/login`,
  `/api/auth/*`, static assets).
- If no `mm_access` **and** no `mm_refresh` cookie → redirect to `/login`.
- Do NOT verify role in middleware (no backend call at the edge) — role is
  enforced in the dashboard layout via `getSession()`.

---

## Definition of Done

- [ ] Login with valid **admin** creds sets httpOnly cookies and lands on `/`;
      the browser has **no** token in JS/localStorage (verify in DevTools).
- [ ] Login with a **non-admin** account is rejected with "not an administrator"
      and sets no cookies.
- [ ] Wrong password shows a friendly `INVALID_CREDENTIALS` message.
- [ ] Visiting a dashboard route with no cookies redirects to `/login`.
- [ ] Logout clears cookies and returns to `/login`.
- [ ] A 401 from the backend triggers one transparent refresh+retry; a failed
      refresh forces re-login.
- [ ] `pnpm typecheck` + `pnpm build` pass.

## Verification commands

```bash
pnpm build && pnpm start
# Manual: login as admin (cookies set, no token in JS), login as student (403),
# wrong password (401 friendly), hit /users with no cookies (→ /login), logout.
```

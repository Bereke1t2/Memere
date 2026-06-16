# Memere Admin — Skill Index (Master)

> **Read this file first.** It is the orchestration map for building the **Memere
> Admin Web Panel** — the browser-based dashboard that platform admins use to run
> Memere (ExamPrep). It tells you *what order* to execute skills in, *what rules
> never change*, and *how phases relate*. Every individual skill file assumes you
> have read and internalized this master file.

---

## 1. What we are building

A **web admin panel** (desktop-first, responsive) for **Memere (ExamPrep)**, the
Grade-12 Ethiopian university-entrance-exam prep platform. The students and
teachers use the **Flutter mobile app** (`memere_mobile`); **admins use this web
panel**. Both clients talk to the **same Go backend** (`Memere-backend`) over its
REST API — this project ships **no backend of its own**.

The full product vision for the backend lives in
[`../../Memere-backend/docs/memere_Design_Specification.md`]. The admin-specific
design is [`docs/Memere_Admin_Design_Specification.md`](./Memere_Admin_Design_Specification.md) —
the **single source of truth** for this project. If a skill file and the design
spec ever disagree, the design spec wins — stop and flag the conflict.

The backend API contract is fixed and documented in
[`../../Memere-backend/docs/Memere.postman_collection.json`]. **We consume it; we
do not change it.** If the admin needs an endpoint the backend does not expose,
stop and flag it as a backend change request — do not invent client-side
work-arounds that violate a Non-Negotiable (e.g. never grade or compute revenue
in the browser).

### Locked technical decisions (do not deviate without asking)

| Decision | Choice | Reason |
|---|---|---|
| Framework | **Next.js 15 (App Router, React 19)** | SSR + Route Handlers let us keep JWTs in httpOnly cookies, out of browser JS (XSS-safe). |
| Language | **TypeScript (strict)** | Type-safe API DTOs mirrored from the Go backend. |
| Styling | **Tailwind CSS v4** | Utility-first; pairs with shadcn/ui. |
| Components | **shadcn/ui (Radix primitives)** | Accessible, unstyled-then-themed; copy-in, not a dep lock-in. |
| Server state | **TanStack Query v5** | Caching, pagination, mutation/invalidation for the API. |
| Forms | **React Hook Form + Zod** | Typed, validated forms; Zod schemas double as runtime guards. |
| Charts | **Recharts** | Revenue/engagement dashboards. |
| Tables | **TanStack Table v8** | Sorting, cursor pagination, column filters. |
| Auth transport | **httpOnly cookies + Next.js Route Handler proxy** | Browser never holds the access/refresh token. |
| Icons | **lucide-react** | Ships with shadcn/ui. |
| Package manager | **pnpm** | Fast, disk-efficient. |
| Node version | **20 LTS+** | Required by Next.js 15. |

---

## 2. The Non-Negotiables (apply to EVERY skill, every phase)

Violating any of these is a build-breaking bug:

1. **The browser never holds a raw JWT.** Access and refresh tokens live in
   **httpOnly, Secure, SameSite cookies** set by a Next.js Route Handler. Client
   components never read `document.cookie` for a token, never put a token in
   `localStorage`/`sessionStorage`, and never receive it in a JSON body.
2. **All backend calls are server-side.** Browser → Next.js Route Handler /
   Server Component / Server Action → Go API. The Go API base URL and any service
   credentials stay server-only (no `NEXT_PUBLIC_` token vars).
3. **Admin-only.** Every route under the dashboard is gated: a valid session AND
   `role === "admin"`. A teacher/student token reaching the panel is rejected
   (the backend already returns 403; the panel must also refuse to render).
4. **No business logic in the client.** Grading, revenue math, MRR, pass-rates,
   refunds, reconciliation — all computed by the backend. The panel **displays**
   backend results; it never recomputes a financial or grading figure locally.
5. **Never log or surface secrets.** No tokens, no raw payment data, no password
   fields in logs, error toasts, or the DOM. Mirror the backend's redaction rule.
6. **Respect the error envelope.** The backend returns
   `{ "code", "message", "details" }`. The panel maps `code` → friendly UI; it
   never blindly renders raw error strings to users.
7. **Destructive actions confirm.** Suspend user, unpublish course, refund
   payment, change role — each requires an explicit confirmation dialog and shows
   the backend's success/failure faithfully (no optimistic lie on failure).
8. **Cursor pagination only.** Lists use the backend's `limit` + `after` cursor
   contract. Never assume offset paging or fetch-all.
9. **Accessibility + responsiveness.** Keyboard-navigable, semantic HTML (shadcn/
   Radix gives this), works down to tablet width. Admins may use it on the go.

---

## 3. Architecture layering (memorize this)

```
  app/(dashboard)/**           ← Server Components: fetch via server API client
        │   page.tsx, layout.tsx, loading.tsx, error.tsx
        ▼
  components/**                ← Client Components: tables, forms, charts, dialogs
        │   ("use client"); receive data as props or via TanStack Query
        ▼
  lib/api/**                   ← typed API client (server-only) + Zod DTO schemas
        │   attaches Bearer from httpOnly cookie; maps error envelope
        ▼
  app/api/**  (Route Handlers) ← /api/auth/login, /api/auth/logout, /api/proxy/*
        │   the ONLY place that sets/reads token cookies & calls the Go API
        ▼
  Go backend  (Memere-backend) ← REST /api/v1/** — the source of all truth
```

**Rule of thumb:** a Client Component must never `import` the server API client
or read a token. If a `"use client"` file imports `lib/api/server`, you broke it.
Data crosses the boundary as plain serializable props or through a Route Handler.

---

## 4. Directory layout (target for end of Phase 1)

```
Memere-admin/
├── app/
│   ├── layout.tsx                  # root layout, providers (Query, theme, toast)
│   ├── globals.css                 # Tailwind v4 + theme tokens
│   ├── login/page.tsx              # public login screen
│   ├── (dashboard)/                # protected group — all admin pages
│   │   ├── layout.tsx              # sidebar + header shell; server-side role guard
│   │   ├── page.tsx                # dashboard home (analytics overview)
│   │   ├── users/                  # Phase 2
│   │   ├── courses/                # Phase 2
│   │   ├── payments/               # Phase 3
│   │   ├── revenue/                # Phase 3
│   │   └── announcements/          # Phase 4
│   └── api/
│       ├── auth/login/route.ts     # POST → backend login, set httpOnly cookies
│       ├── auth/logout/route.ts    # clear cookies + backend logout
│       └── auth/refresh/route.ts   # rotate access token via refresh cookie
├── components/
│   ├── ui/                         # shadcn/ui primitives (button, table, dialog…)
│   ├── layout/                     # sidebar, header, nav, user-menu
│   ├── data-table/                 # reusable TanStack Table wrapper
│   └── charts/                     # Recharts wrappers
├── lib/
│   ├── api/
│   │   ├── server.ts               # server-only fetch client (Bearer from cookie)
│   │   ├── client.ts               # browser → Route Handler helpers (no token)
│   │   ├── endpoints.ts            # typed endpoint functions (listUsers, …)
│   │   └── schemas.ts              # Zod schemas mirroring backend DTOs
│   ├── auth/
│   │   ├── session.ts              # read/verify session cookie server-side
│   │   └── cookies.ts              # cookie name + set/clear helpers
│   ├── query/                      # TanStack Query keys + provider
│   └── utils.ts                    # cn(), formatters (money, dates, percentages)
├── middleware.ts                   # edge gate: no session cookie → /login
├── types/                          # shared TS types (mirror Go DTOs)
├── .env.local.example              # API_BASE_URL, COOKIE_SECRET, …
├── components.json                 # shadcn/ui config
├── next.config.ts
├── tailwind.config.ts (or v4 css)
├── tsconfig.json
├── package.json
└── README.md
```

---

## 5. Phase map (admin panel)

We build **phase by phase**, strictly in order. Each phase is gated on the
previous phase's final Definition of Done. Each phase has a kickoff prompt in
`docs/prompts/`.

| Phase | Theme | Status |
|---|---|---|
| **Phase 1** | **Foundation + Auth + App Shell + Dashboard** | ✅ Authored |
| **Phase 2** | **User Management + Course Moderation** | ✅ Authored |
| **Phase 3** | **Payments + Revenue + Subscriptions** | ✅ Authored |
| **Phase 4** | **Announcements + Hardening + Deployment (final)** | ✅ Authored |

### Phase 1 skill order (run strictly in sequence)

| # | File | Builds |
|---|---|---|
| 1 | [`skills/phase1/skill1.md`](./skills/phase1/skill1.md) | Project scaffold: Next.js 15 + TS + Tailwind v4 + shadcn/ui, env config, base layout, providers, README |
| 2 | [`skills/phase1/skill2.md`](./skills/phase1/skill2.md) | API layer: server-only fetch client, Zod DTO schemas mirroring the backend, typed endpoint functions, error-envelope mapping |
| 3 | [`skills/phase1/skill3.md`](./skills/phase1/skill3.md) | Auth: login page, Route Handlers (login/logout/refresh), httpOnly token cookies, session helpers, `middleware.ts` gate, admin-role enforcement |
| 4 | [`skills/phase1/skill4.md`](./skills/phase1/skill4.md) | App shell: protected `(dashboard)` layout, sidebar + header nav, user menu, theme toggle, responsive behavior, loading/error boundaries |
| 5 | [`skills/phase1/skill5.md`](./skills/phase1/skill5.md) | Dashboard home: analytics overview cards + engagement stats wired to `/admin/analytics/*`, first charts, Phase 1 smoke test |

### Phase 2 skill order (only after Phase 1 is green)

| # | File | Builds |
|---|---|---|
| 1 | [`skills/phase2/skill1.md`](./skills/phase2/skill1.md) | Reusable `DataTable` (TanStack Table) with cursor pagination, column filters, empty/loading states |
| 2 | [`skills/phase2/skill2.md`](./skills/phase2/skill2.md) | Users list page: search, role filter, paginated table wired to `/admin/users` |
| 3 | [`skills/phase2/skill3.md`](./skills/phase2/skill3.md) | User detail + actions: suspend, reactivate, change-role with confirm dialogs and mutation/invalidation |
| 4 | [`skills/phase2/skill4.md`](./skills/phase2/skill4.md) | Courses list + moderation: list, filter, unpublish with reason |
| 5 | [`skills/phase2/skill5.md`](./skills/phase2/skill5.md) | Course detail view + Phase 2 polish & smoke test |

### Phase 3 skill order (only after Phase 2 is green)

| # | File | Builds |
|---|---|---|
| 1 | [`skills/phase3/skill1.md`](./skills/phase3/skill1.md) | Payments list: status filter, paginated table wired to `/admin/payments` |
| 2 | [`skills/phase3/skill2.md`](./skills/phase3/skill2.md) | Payment detail + refund + bulk reconcile-pending action |
| 3 | [`skills/phase3/skill3.md`](./skills/phase3/skill3.md) | Revenue dashboard: platform totals, MRR, provider breakdown charts, date-range picker |
| 4 | [`skills/phase3/skill4.md`](./skills/phase3/skill4.md) | Subscriptions overview + financial KPI tiles |
| 5 | [`skills/phase3/skill5.md`](./skills/phase3/skill5.md) | Phase 3 polish, money/date formatting audit, smoke test |

### Phase 4 skill order (only after Phase 3 is green) — FINAL

| # | File | Builds |
|---|---|---|
| 1 | [`skills/phase4/skill1.md`](./skills/phase4/skill1.md) | Announcements: broadcast composer (segment targeting) wired to `/admin/announcements`, history view |
| 2 | [`skills/phase4/skill2.md`](./skills/phase4/skill2.md) | Global hardening: error boundaries, toasts, loading skeletons, empty states, 401/refresh handling, retry/backoff |
| 3 | [`skills/phase4/skill3.md`](./skills/phase4/skill3.md) | Security review: cookie flags, CSRF on mutations, CSP headers, no-token-in-client audit, dependency scan |
| 4 | [`skills/phase4/skill4.md`](./skills/phase4/skill4.md) | Build & containerize: production build, Dockerfile (standalone output), env contract, image hardening |
| 5 | [`skills/phase4/skill5.md`](./skills/phase4/skill5.md) | Deployment: Vercel and self-hosted (VPS/Docker) runbooks, CORS coordination with backend, final smoke test, docs |

Each skill ends with a **Definition of Done** checklist and **verification
commands**. A skill is not complete until every box is checked and the
verification commands pass.

---

## 6. How to work a skill (the loop to follow)

For each skill file, in order:

1. **Read** the whole skill file before writing any code.
2. **Re-read** the relevant section of `Memere_Admin_Design_Specification.md` and
   the backend Postman collection it points to.
3. **Plan** the files you'll create/modify (list them).
4. **Implement** following the file's "Tasks" section exactly.
5. **Self-check** against the "Non-Negotiables" (Section 2 above).
6. **Verify** by running the skill's verification commands (`pnpm build`,
   `pnpm lint`, `pnpm typecheck`, manual click-through).
7. **Tick** the Definition of Done. Only then move to the next skill.

Do not jump ahead. Skill N assumes Skill N-1's Definition of Done is satisfied.

---

## 7. Conventions used in every skill file

- **API access:** server-side only, through `lib/api`. Bearer token pulled from
  the httpOnly cookie inside Route Handlers / Server Components.
- **Data fetching:** Server Components for first paint; TanStack Query for
  client-side refetch, mutations, and cache invalidation.
- **Error UI:** map backend `code` → friendly message via a small dictionary;
  fall back to `message`; toast on mutation failure, inline on load failure.
- **Money & dates:** format with shared helpers (`formatETB`, `formatDate`);
  never hand-roll. Currency/amount values come from the backend as strings —
  parse with care, never with `Number()` on a money field for display math.
- **Commits:** Conventional Commits (`feat:`, `fix:`, `docs:`, `refactor:`,
  `test:`, `chore:`). One commit per skill task where practical.
- **Styling:** Tailwind utilities + shadcn/ui tokens; no inline hex colors, use
  theme variables. Dark mode supported from Phase 1.

---

## 8. How to start (the very first prompt)

The exact wording is in [`docs/prompts/phase1-kickoff.md`](./prompts/phase1-kickoff.md).
In short: read this `skill.md`, read the admin design spec, skim the backend
Postman collection, then execute `phase1/skill1.md` through `skill5.md` in order,
stopping after each skill's Definition of Done for review.

---

## 9. Moving between phases

1. Start a phase with its kickoff prompt (e.g. `docs/prompts/phase1-kickoff.md`).
2. Run its skills `skill1 → skill5` strictly in sequence, stopping after each
   skill's Definition of Done for review.
3. When a phase's final skill passes, report "Phase N complete" and move to the
   next phase's kickoff prompt.

**Never** scaffold a later phase's pages while building an earlier one. Phases are
gated on review. The backend is fixed — if a phase needs a new endpoint, raise it
as a backend change, do not fake it client-side.

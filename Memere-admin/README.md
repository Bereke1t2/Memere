# Memere Admin

Web admin panel for the **Memere (ExamPrep)** platform — a Next.js 15 client of
the Go backend (`../Memere-backend`). It owns no database and no business logic;
all data is sourced from the Go REST API via server-side calls.

---

## Architecture

```
Browser
  │  (httpOnly cookies only — no token in JS)
  ▼
Next.js 15 (App Router)
  ├── Server Components  → fetch data server-side → render HTML
  ├── Route Handlers     → proxy mutations to Go API (tokens from cookie)
  └── Client Components → UI, TanStack Query (calls Route Handlers, never Go API)
            │
            ▼ (server-to-server, Bearer token from httpOnly cookie)
      Go Backend (Memere-backend)
            │
            ▼
      PostgreSQL · Redis
```

**Key invariants:**
- The browser **never** holds a raw JWT — access/refresh tokens live in `httpOnly, Secure, SameSite=Lax` cookies set by Route Handlers.
- All Go API calls originate from the Next.js server process — browser CORS is not required.
- No financial or grading logic is computed client-side; the panel displays backend-returned values only.

---

## Screens

| Route | Description |
|---|---|
| `/` | Dashboard — live KPI tiles, engagement stats, revenue chart |
| `/users` | User management — list, filter, paginate, detail, suspend/reactivate/change-role |
| `/courses` | Course moderation — list, paginate, detail with sections, unpublish |
| `/payments` | Payment records — list, filter by status, detail, refund, reconcile |
| `/revenue` | Revenue dashboard — KPIs, provider charts, trend, breakdown table |
| `/announcements` | Broadcast composer — segment targeting, double-confirm for "all" |

---

## Prerequisites

- Node.js ≥ 20 LTS
- pnpm ≥ 11

---

## Local development

```bash
pnpm install
cp .env.local.example .env.local   # fill in API_BASE_URL and COOKIE_SECRET
pnpm dev                            # http://localhost:3000
```

### Environment variables

| Variable | Required | Notes |
|---|---|---|
| `API_BASE_URL` | **Yes** | Go backend base URL — server-only, never exposed to browser |
| `COOKIE_SECRET` | **Yes** | ≥ 32 random bytes: `openssl rand -hex 32` |
| `NODE_ENV` | **Yes** | Set to `production` in prod |
| `NEXT_PUBLIC_APP_NAME` | No | Display name (safe to expose, default: "Memere Admin") |

### Scripts

```bash
pnpm dev        # start dev server (Turbopack)
pnpm build      # production build (standalone output)
pnpm start      # serve production build
pnpm typecheck  # TypeScript strict check (no emit)
pnpm lint       # ESLint
make docker-build  # build container image
make docker-run    # run image locally with .env.local
```

---

## Deployment

See [`docs/DEPLOYMENT.md`](docs/DEPLOYMENT.md) for the full runbook. Quick summary:

### Vercel

1. Import repo → set `API_BASE_URL`, `COOKIE_SECRET`, `NODE_ENV=production` as server env vars.
2. Deploy. Add custom domain `admin.memere.app`.

### Self-hosted (Docker + nginx)

```bash
make docker-build
docker tag memere-admin:local registry.example.com/memere-admin:latest
docker push registry.example.com/memere-admin:latest
# On server: docker compose up -d (see docs/DEPLOYMENT.md for compose + nginx config)
```

---

## Tech stack

| Layer | Choice |
|---|---|
| Framework | Next.js 15, React 19, TypeScript strict |
| Styling | Tailwind CSS v4 + shadcn/ui (Radix primitives) |
| Server state | TanStack Query v5 |
| Tables | TanStack Table v8 |
| Forms | React Hook Form + Zod |
| Charts | Recharts |
| Icons | lucide-react |
| Package manager | pnpm |

---

## Related

- **Backend**: `../Memere-backend` — Go REST API
- **API contract**: `../Memere-backend/docs/Memere.postman_collection.json`
- **Design spec**: `docs/Memere_Admin_Design_Specification.md`
- **Smoke test**: `docs/SMOKE.md`
- **Deployment**: `docs/DEPLOYMENT.md`

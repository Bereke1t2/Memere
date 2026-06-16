# Memere Admin

Web admin panel for the **Memere (ExamPrep)** platform. It is a Next.js client
of the existing Go backend (`../Memere-backend`). It owns no database and no
business logic — all data is sourced from the Go REST API.

## Prerequisites

- Node.js ≥ 20 LTS
- pnpm ≥ 8

## Setup

```bash
pnpm install
cp .env.local.example .env.local   # fill in API_BASE_URL and COOKIE_SECRET
pnpm dev
```

Visit [http://localhost:3000](http://localhost:3000).

## Environment variables

See `.env.local.example` for the full contract. The critical ones:

| Variable | Required | Notes |
|---|---|---|
| `API_BASE_URL` | Yes | Go backend base URL (server-side only) |
| `COOKIE_SECRET` | Yes | Signs/verifies session cookie integrity |
| `NEXT_PUBLIC_APP_NAME` | No | Display name (safe to expose) |

## Scripts

```bash
pnpm dev        # start dev server
pnpm build      # production build
pnpm start      # serve production build
pnpm typecheck  # TypeScript type check (no emit)
pnpm lint       # ESLint
pnpm format     # Prettier
```

## Architecture

```
Browser → Next.js Route Handler → Go API (Memere-backend)
```

All backend calls are server-side. JWTs live in httpOnly cookies; the browser
never holds a raw token. See `docs/skill.md` for the full architecture layering.

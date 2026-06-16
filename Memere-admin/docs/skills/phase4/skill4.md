# Phase 4 · Skill 4 — Build & Containerize

> **Prerequisite:** Skill 3 done. Read [`docs/skill.md`](../../skill.md) §1, §4,
> [`design spec`](../../Memere_Admin_Design_Specification.md) §9 (env contract).

---

## Goal

A production build and a hardened container image, with a clear server-only env
contract — ready for both Vercel and self-hosted (Docker/VPS) deploys.

---

## Tasks

### 4.1 — Standalone output

- `next.config.ts`: `output: "standalone"` for a minimal self-contained server
  bundle. Confirm `pnpm build` emits `.next/standalone`.

### 4.2 — Dockerfile (multi-stage)

```
# deps → build → runtime (node:20-alpine, non-root)
```

- Stage 1 `deps`: `pnpm install --frozen-lockfile`.
- Stage 2 `build`: `pnpm build` (standalone).
- Stage 3 `runtime`: copy `.next/standalone` + `.next/static` + `public`; run as a
  non-root user; `EXPOSE 3000`; `CMD ["node","server.js"]`.
- Pin the base image; no secrets baked in. `API_BASE_URL` injected at runtime.

### 4.3 — Env contract + `.dockerignore`

- Document required runtime env: `API_BASE_URL`, `COOKIE_SECRET`, `NODE_ENV`,
  `NEXT_PUBLIC_APP_NAME`. `.dockerignore` excludes `node_modules`, `.next`,
  `.env*`, `docs`.

### 4.4 — `Makefile` / scripts

- `docker-build`, `docker-run` (with env file), `build`, `start` targets.

---

## Definition of Done

- [ ] `pnpm build` produces a standalone server; app runs via `node server.js`.
- [ ] `docker build` succeeds; the image runs non-root and serves on :3000 with
      `API_BASE_URL` provided at runtime.
- [ ] No secret/token baked into the image; `.dockerignore` excludes env + docs.
- [ ] Image starts and reaches the login page against a reachable backend.

## Verification commands

```bash
pnpm build
docker build -t memere-admin:local .
docker run --rm -p 3000:3000 \
  -e API_BASE_URL=https://api.memere.app \
  -e COOKIE_SECRET=$(openssl rand -hex 32) \
  -e NODE_ENV=production \
  memere-admin:local
# visit http://localhost:3000/login
```

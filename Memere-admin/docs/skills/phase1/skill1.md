# Phase 1 · Skill 1 — Project Scaffold

> **Prerequisite:** Node ≥ 20, pnpm installed. Read [`docs/skill.md`](../../skill.md)
> §1–4 and [`Memere_Admin_Design_Specification.md`](../../Memere_Admin_Design_Specification.md)
> §1, §7, §9.
>
> **Spec references:** skill.md §1 (stack), §3 (layering), §4 (directory layout);
> design spec §7 (UX conventions), §9 (environment contract).

---

## Goal

Stand up a runnable Next.js 15 + TypeScript app with Tailwind v4, shadcn/ui, the
global providers (TanStack Query, theme, toaster), the env contract, and the base
layout. No business screens yet — just a clean, themed, type-checked shell that
builds.

---

## Tasks

### 1.1 — Initialize the app

```bash
pnpm create next-app@latest . \
  --typescript --tailwind --eslint --app --src-dir=false \
  --import-alias "@/*" --use-pnpm
```

- App Router, TypeScript strict (`tsconfig.json`: `"strict": true`).
- Confirm `pnpm dev` serves a page, `pnpm build` succeeds.

### 1.2 — Tailwind v4 + theme tokens

- Configure `app/globals.css` with Tailwind v4 and CSS theme variables (light +
  dark) using the shadcn token set (`--background`, `--foreground`, `--primary`,
  `--muted`, `--border`, etc.).
- Dark mode via `class` strategy.

### 1.3 — shadcn/ui

```bash
pnpm dlx shadcn@latest init      # neutral base color, CSS variables
pnpm dlx shadcn@latest add button card input label badge dropdown-menu \
  dialog sonner skeleton table tabs select avatar tooltip
```

- `components.json` committed. Primitives land in `components/ui/`.

### 1.4 — Providers

Create `app/providers.tsx` (`"use client"`) wrapping children with:
- TanStack Query `QueryClientProvider` (a single `QueryClient`, sensible defaults:
  `staleTime` 30s, `retry` 1).
- Theme provider (`next-themes`) for light/dark.
- `<Toaster />` from sonner.

Wire it into `app/layout.tsx`. Set `<html suppressHydrationWarning>` and base
metadata (`title: "Memere Admin"`).

### 1.5 — Utilities

- `lib/utils.ts`: `cn()` (clsx + tailwind-merge), plus formatter stubs
  `formatDate(iso)`, `formatMoney(amount, currency)`, `formatPercent(n)` (filled
  out later but defined now).

### 1.6 — Env contract

- `.env.local.example` with `API_BASE_URL`, `COOKIE_SECRET`,
  `NEXT_PUBLIC_APP_NAME` (per design spec §9). Add `.env.local` to `.gitignore`
  (Next's template already does).
- A tiny `lib/env.ts` that reads + validates server env at startup (throw if
  `API_BASE_URL` missing in production). **Never** expose `API_BASE_URL` to the
  client.

### 1.7 — Scripts & README

- `package.json` scripts: `dev`, `build`, `start`, `lint`, `typecheck`
  (`tsc --noEmit`), `format`.
- `README.md`: what this is (admin client of the Go backend), prerequisites,
  `pnpm install`, copy `.env.local.example` → `.env.local`, `pnpm dev`.

---

## Definition of Done

- [ ] `pnpm dev` renders a themed placeholder page with working light/dark toggle.
- [ ] `pnpm build` and `pnpm typecheck` pass with zero errors.
- [ ] `pnpm lint` is clean.
- [ ] shadcn/ui primitives present in `components/ui/`; a sample `<Button>` renders.
- [ ] `.env.local.example` documents the env contract; `API_BASE_URL` is server-only.
- [ ] README explains setup and the "client of the Go backend" framing.

## Verification commands

```bash
pnpm install
pnpm typecheck
pnpm lint
pnpm build
pnpm dev   # visit http://localhost:3000, toggle theme
```

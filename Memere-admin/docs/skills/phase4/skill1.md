# Phase 4 · Skill 1 — Announcements Broadcast Composer

> **Prerequisite:** Phase 3 green. Read [`docs/skill.md`](../../skill.md) §2 (#7),
> [`design spec`](../../Memere_Admin_Design_Specification.md) §4.7, §2.2
> (`/admin/announcements`).

---

## Goal

The Announcements screen: a composer that broadcasts to a user segment via
`POST /admin/announcements`, with an extra confirmation for "all" and a
session-local history of sends (the backend has no list endpoint in v1).

---

## Tasks

### 1.1 — Endpoint

- `broadcast({ title, body, segment, data })` (POST, 204) via Route Handler /
  Server Action. `segment ∈ {all, students, teachers, subscribers}`.

### 1.2 — Composer (`app/(dashboard)/announcements/page.tsx`)

- RHF + Zod form: title (required), body (required textarea), segment (select),
  optional key/value `data` pairs.
- Submit → confirm dialog; if `segment === "all"` require an **extra** explicit
  confirmation ("This sends to ALL users"). On success → toast + reset form.

### 1.3 — Session history (optional)

- Keep an in-memory/session list of sends (title, segment, time) shown below the
  composer. Clearly labelled as "this session only" since the backend stores no
  retrievable list.

---

## Definition of Done

- [ ] An admin can compose + send an announcement to a segment.
- [ ] "All" requires a second confirmation; success toasts and resets.
- [ ] Validation prevents empty title/body; backend errors surface faithfully.
- [ ] Session history (if built) is labelled session-only.
- [ ] `pnpm typecheck` + `pnpm build` pass.

## Verification commands

```bash
pnpm build && pnpm start
# Manual: send to "students", then attempt "all" (double-confirm), verify toasts.
```

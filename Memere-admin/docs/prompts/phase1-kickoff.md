# Phase 1 — Kickoff Prompt (Memere Admin)

This file gives (1) how the skills are structured, (2) the **exact first prompt**
to start the build, and (3) how to drive it skill-by-skill.

---

## 1. How the skill files are structured

```
Memere-admin/
└── docs/
    ├── Memere_Admin_Design_Specification.md  ← source of truth
    ├── skill.md                              ← MASTER index: rules, layering, phase map
    ├── prompts/
    │   ├── phase1-kickoff.md   ← this file
    │   ├── phase2-kickoff.md
    │   ├── phase3-kickoff.md
    │   └── phase4-kickoff.md
    └── skills/
        ├── phase1/  skill1..5  ← Foundation, API layer, auth, shell, dashboard
        ├── phase2/  skill1..5  ← Users + Courses
        ├── phase3/  skill1..5  ← Payments + Revenue + Subscriptions
        └── phase4/  skill1..5  ← Announcements + Hardening + Deployment
```

> **All four phases are authored.** Build them in order, gating each on the
> previous phase's final Definition of Done.

**Reading order:** `skill.md` → admin design spec → backend Postman collection →
`skill1` → `skill2` → … strictly in sequence.

---

## 2. The first prompt

> Copy the box below to start the build.

```
You are building the WEB ADMIN PANEL for "Memere (ExamPrep)". It is a Next.js
client of the existing Go backend (../Memere-backend). It owns no backend.

STEP 0 — Read, in this order, before writing any code:
  1. docs/skill.md                                  (master index — rules & phase map)
  2. docs/Memere_Admin_Design_Specification.md       (source of truth)
  3. ../Memere-backend/docs/Memere.postman_collection.json  (the API contract)
  4. docs/skills/phase1/skill1.md                    (the first build skill)

GROUND RULES (docs/skill.md §2 — never violate):
  - The browser NEVER holds a raw JWT. Tokens live in httpOnly cookies set by a
    Next.js Route Handler. No tokens in localStorage / JS / JSON to the client.
  - All backend calls are server-side (Server Components / Route Handlers).
    API_BASE_URL and secrets are server-only (no NEXT_PUBLIC_ token vars).
  - Admin-only: every dashboard route requires a valid session AND role==="admin".
  - No business logic in the client: grading, revenue, MRR, refunds are backend
    results we DISPLAY, never recompute.
  - Respect the error envelope { code, message, details }. Confirm destructive
    actions. Cursor pagination only.
  - Stack is locked: Next.js 15 (App Router) + TS strict + Tailwind v4 +
    shadcn/ui + TanStack Query/Table + RHF + Zod + Recharts + pnpm. Do not
    substitute without asking.

HOW TO WORK:
  - Execute ONLY skill1.md now. Follow its "Tasks" exactly.
  - Run its "Verification commands" and check every Definition-of-Done box.
  - Then STOP and report what you built, the verification output, and any
    deviations. Wait for "continue" before skill2.md.

Begin with STEP 0, then implement docs/skills/phase1/skill1.md.
```

---

## 3. Driving the build skill-by-skill

After each skill, stop and report. Reusable continue prompt:

```
Verification looks good. Now execute docs/skills/phase1/skill<N>.md:
read it fully, re-read the design-spec sections it references, follow its Tasks,
run its Verification commands, check its Definition of Done, then stop and report.
Honor all Non-Negotiables in docs/skill.md §2.
```

Why stop between skills? Each skill is a clean, reviewable unit (ideally one
commit). It catches a violation at skill 2 instead of skill 5.

---

## 4. When Phase 1 is finished

Skill 5's Definition of Done is the Phase-1 gate. When it passes:

1. Report "Phase 1 complete" with a summary + smoke-test result.
2. Proceed to **Phase 2** using `docs/prompts/phase2-kickoff.md` (User & Course
   management), which assumes the Phase 1 shell + auth + API layer are green.

**Phase roadmap:**

| Phase | Theme | Kickoff |
|---|---|---|
| 1 | Foundation + Auth + Shell + Dashboard | `prompts/phase1-kickoff.md` |
| 2 | User Management + Course Moderation | `prompts/phase2-kickoff.md` |
| 3 | Payments + Revenue + Subscriptions | `prompts/phase3-kickoff.md` |
| 4 | Announcements + Hardening + Deployment | `prompts/phase4-kickoff.md` |

---

## 5. Tips

- **Keep the backend Postman collection open.** Bind to the exact DTO fields it
  returns; don't invent fields.
- **Don't let it call the API from the browser.** If a Client Component fetches
  the Go API directly with a token, stop it — that's a Non-Negotiable violation.
- **Confirm Node ≥ 20 and pnpm installed** before skill 1.
- **One commit per skill** (Conventional Commits, e.g.
  `feat: phase1 skill1 admin scaffold`).

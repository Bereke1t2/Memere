# Phase 1 — Antigravity Kickoff Prompt

This file gives you (1) how to structure the skills for Antigravity, (2) the
**exact first prompt** to paste into Antigravity to start the build, and (3) how
to drive it skill-by-skill and phase-by-phase.

---

## 1. How the skill files are structured

```
memere-backend/
└── docs/
    ├── memere_Design_Specification.md   ← source of truth (the full design doc)
    ├── skill.md                         ← MASTER index: rules, layering, phase map
    ├── prompts/
    │   └── phase1-kickoff.md            ← this file
    └── skills/
        ├── phase1/   skill1..5  ← Foundation, domain, auth, course, HTTP (build first)
        ├── phase2/   skill1..5  ← Quiz & Exam engines
        ├── phase3/   skill1..5  ← Video pipeline
        ├── phase4/   skill1..5  ← Payments & Enrollments
        ├── phase5/   skill1..5  ← Progress, Notifications, Admin, Certificates
        └── phase6/   skill1..5  ← Hardening, observability, CI/CD, k8s (final)
```

> **All six phases are already authored.** Each phase has its own kickoff prompt
> (`docs/prompts/phaseN-kickoff.md`). This file kicks off **Phase 1**; build the
> phases in order, gating each on the previous phase's final Definition of Done.

**Reading order Antigravity must follow:** `skill.md` → design spec → `skill1`
→ `skill2` → `skill3` → `skill4` → `skill5`, strictly in sequence. Each skill has
a **Definition of Done** that gates the next.

---

## 2. The first prompt to paste into Antigravity

> Copy everything in the box below into Antigravity as your opening message.

```
You are building the Go backend for "Memere (ExamPrep)". All instructions live in
this repo's docs folder. Work strictly from them.

STEP 0 — Read, in this order, and do not write any code until you have:
  1. docs/skill.md                          (master index — rules & phase map)
  2. docs/memere_Design_Specification.md     (the full design; the source of truth)
  3. docs/skills/phase1/skill1.md            (the first build skill)

GROUND RULES (from docs/skill.md §2 "Non-Negotiables" — never violate):
  - Clean Architecture: dependencies point inward only; domain has zero external
    deps and no db/json tags.
  - Server-side grading & timers; pre-signed video URLs; payment idempotency;
    soft deletes only; filter every query by the authenticated user_id (no IDOR);
    never log passwords/tokens; UUID PKs; HTTPS at the edge.
  - Stack is locked: Go 1.22+, Gin, pgx v5, sqlc, golang-migrate, go-redis,
    golang-jwt, bcrypt. Do not substitute libraries without asking.

HOW TO WORK:
  - Execute ONLY skill1.md now. Follow its "Tasks" exactly.
  - When done, run its "Verification commands" and check every box in its
    "Definition of Done".
  - Then STOP and report: what you built, the verification output, and any
    deviations. Wait for my "continue" before starting skill2.md.

Begin with STEP 0, then implement docs/skills/phase1/skill1.md.
```

---

## 3. Driving the build skill-by-skill

After each skill, Antigravity stops and reports. You review, then continue:

- To proceed: reply **`continue with docs/skills/phase1/skill2.md`** (then skill3,
  4, 5 the same way).
- If something's off: point at the specific Definition-of-Done item that failed
  and ask it to fix before moving on.

A short per-skill continue prompt you can reuse:

```
Verification looks good. Now execute docs/skills/phase1/skill<N>.md:
read it fully, re-read the design-spec sections it references, follow its Tasks,
run its Verification commands, check its Definition of Done, then stop and report.
Honor all Non-Negotiables in docs/skill.md §2.
```

Why stop between skills? Each skill is a clean, reviewable unit (ideally one
commit). It keeps the model from drifting and lets you catch an architecture
violation at skill 2 instead of discovering it at skill 5.

---

## 4. When Phase 1 is finished

Skill 5's Definition of Done is the Phase-1 gate. When it passes:

1. Antigravity reports "Phase 1 complete" with a summary + smoke-test result.
2. Proceed to **Phase 2** using its kickoff prompt:
   [`docs/prompts/phase2-kickoff.md`](./phase2-kickoff.md). Phase 2 covers the
   **Quiz & Exam engines** (spec §9) and assumes the Phase 1 codebase is green.
3. Run each subsequent phase the same way: open `phaseN-kickoff.md` → skill-by-skill
   → gate on the final Definition of Done.

**Phase roadmap (backend only) — all phases authored:**

| Phase | Theme | Kickoff prompt |
|---|---|---|
| 1 | Foundation + Auth + Course CRUD + HTTP wiring | `prompts/phase1-kickoff.md` |
| 2 | Quiz & Exam engines (grading, timers, attempt state) | `prompts/phase2-kickoff.md` |
| 3 | Video pipeline (pre-signed upload, FFmpeg→HLS, signed delivery) | `prompts/phase3-kickoff.md` |
| 4 | Payments (Chapa/Telebirr/Stripe) + enrollments + access | `prompts/phase4-kickoff.md` |
| 5 | Progress + Notifications + Admin analytics + Certificates | `prompts/phase5-kickoff.md` |
| 6 | Hardening, observability, performance, CI/CD, k8s (final) | `prompts/phase6-kickoff.md` |

Dependency notes carried across phases (enforced by grep gates in the skills):
Phases 2–3 leave `TODO(phase4)` access hooks that **Phase 4 removes**; Phases 3–4
leave `notifyNoop` hooks that **Phase 5 wires up**; Phase 6 changes no business
behavior.

---

## 5. Tips for a smooth run

- **Keep the design spec open.** If Antigravity proposes something that conflicts
  with `memere_Design_Specification.md`, the spec wins — tell it to re-read.
- **Don't let it skip ahead.** If it starts writing quiz/payment code in Phase 1,
  stop it; that's a later phase.
- **One commit per skill** (Conventional Commits, e.g.
  `feat: phase1 skill1 project foundation`). Makes review and rollback easy.
- **Go version:** confirm `go version` ≥ 1.22 before skill 1 (the dev box had
  1.20.5).
- If a skill's scope is too big for one Antigravity run, it's fine to split a
  single skill across multiple turns — just don't cross into the next skill until
  the current Definition of Done is green.
```

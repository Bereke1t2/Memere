# Phase 2 — Antigravity Kickoff Prompt

Use this **only after Phase 1 is complete and green** (auth + course API live,
schema migrated, all Phase 1 Definitions of Done checked). Phase 2 builds the
**Quiz & Exam engines** on top of the Phase 1 codebase.

---

## 1. The first prompt to paste into Antigravity

> Copy everything in the box below.

```
We are starting Phase 2 of the Memere backend (Quiz & Exam engines). Phase 1 is
already complete and green. Work strictly from the repo docs.

STEP 0 — Read, in this order, before writing any code:
  1. docs/skill.md                          (master index — rules & phase map)
  2. docs/memere_Design_Specification.md §9  (Quiz & Mock Exam Engine) and §4.2.5–4.2.6 (tables)
  3. docs/skills/phase2/skill1.md            (the first Phase 2 build skill)

CARRY-OVER GROUND RULES (docs/skill.md §2 — never violate):
  - Correct answers NEVER leave the server; grading is server-side only.
  - Exam/quiz timers are enforced server-side (started_at/expires_at in DB);
    the client timer is display-only; abandoned attempts auto-grade.
  - Clean Architecture (deps inward only; domain has no db/json tags).
  - Soft deletes, UUID PKs, filter every query by the authenticated user_id
    (no IDOR), never log answers/tokens.
  - Reuse Phase 1's delivery layer (apperror envelope, middleware, router,
    constructor wiring) — extend it, do not rebuild it.
  - Stack stays locked: Gin, pgx v5, sqlc, golang-migrate, go-redis, golang-jwt.

HOW TO WORK:
  - Execute ONLY skill1.md now. Migrations in Phase 2 are ADDITIVE (continue
    numbering from Phase 1's 0007). Do not edit or re-run Phase 1 migrations.
  - When done, run its "Verification commands" and check every box in its
    "Definition of Done". Then STOP and report. Wait for my "continue" before
    skill2.md.

Begin with STEP 0, then implement docs/skills/phase2/skill1.md.
```

---

## 2. Driving it skill-by-skill

Same loop as Phase 1. After each skill Antigravity stops and reports; you review,
then send:

```
Verification looks good. Now execute docs/skills/phase2/skill<N>.md:
read it fully, re-read the design-spec sections it references, follow its Tasks,
run its Verification commands, check its Definition of Done, then stop and report.
Honor all Non-Negotiables in docs/skill.md §2 and reuse the Phase 1 delivery layer.
```

Phase 2 order: skill1 (data) → skill2 (quiz engine) → skill3 (exam engine) →
skill4 (sweeper + analytics) → skill5 (HTTP + wiring + smoke).

---

## 3. Phase-2-specific watch-outs

- **Additive migrations only.** Phase 2 adds `quiz_attempts`, `exam_questions`,
  and a few `ALTER`s. It must not modify Phase 1 migration files.
- **The answer-key leak is the #1 risk.** At skill2/skill3/skill5, verify the
  client-facing question DTO has no `is_correct`/correct-answer field. The skills
  bake in grep checks and smoke assertions — make sure they pass.
- **Shared grading core.** Don't let the quiz and exam engines duplicate grading;
  skill3 refactors it into one package.
- **The sweeper** (skill4) is what makes "server timer fires" true for abandoned
  attempts — its race-safety test (sweeper vs late submit → exactly one grade) is
  not optional.
- **Don't skip ahead** into video/payment endpoints — those are Phase 3+.

---

## 4. When Phase 2 is finished

Skill 5's Definition of Done is the Phase-2 gate. When it passes:

1. Antigravity reports "Phase 2 complete" with a summary + smoke-test result.
2. Come back to Claude and say **"author Phase 3 skills"** — I'll write
   `docs/skills/phase3/` for the **Video pipeline** (pre-signed S3 upload, HLS
   transcode trigger/metadata, streaming/download URLs; spec §8), including the
   `videos` table Go layer stubbed back in Phase 1.
3. Do not scaffold Phase 3 code before its skills are written and reviewed.

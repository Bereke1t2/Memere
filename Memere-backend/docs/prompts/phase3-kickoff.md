# Phase 3 — Antigravity Kickoff Prompt

Use this **only after Phase 2 is complete and green**. Phase 3 builds the **Video
pipeline** (spec §8): pre-signed upload, FFmpeg→HLS transcoding, and secure
signed delivery, on top of the Phase 1/2 codebase.

> **New local prerequisites:** **FFmpeg + ffprobe** installed on the machine that
> runs the transcode worker, and **MinIO** (added to docker-compose in Skill 1)
> for S3-compatible local storage. Confirm `ffmpeg -version` works before Skill 3.

---

## 1. The first prompt to paste into Antigravity

> Copy everything in the box below.

```
We are starting Phase 3 of the Memere backend (Video pipeline). Phases 1 and 2
are already complete and green. Work strictly from the repo docs.

STEP 0 — Read, in this order, before writing any code:
  1. docs/skill.md                          (master index — rules & phase map)
  2. docs/memere_Design_Specification.md §8  (Video Learning System) and §4.2.4 (videos table)
  3. docs/skills/phase3/skill1.md            (the first Phase 3 build skill)

CARRY-OVER GROUND RULES (docs/skill.md §2 — never violate):
  - Pre-signed CDN URLs ONLY for video; no public S3; URLs are short-lived (≈2h).
  - Resolve video -> lesson -> course server-side and authorize the AUTHENTICATED
    user (no IDOR). Never return raw storage keys to clients.
  - Clean Architecture: storage/queue/transcoder are PORTS (interfaces) in the
    domain; S3/MinIO/FFmpeg/in-proc-queue are infrastructure impls. Usecases
    depend on the interfaces only.
  - Migrations are ADDITIVE (continue numbering after Phase 2's 0010).
  - Reuse Phase 1/2 delivery layer (apperror envelope, middleware, router,
    constructor wiring) — extend it, do not rebuild it.
  - Each skill file contains TEMPLATE CODE. Treat it as a starting scaffold to
    adapt to the real package paths/signatures — not as literal copy-paste.

HOW TO WORK:
  - Execute ONLY skill1.md now. Follow its Tasks; run its Verification commands;
    check every box in its Definition of Done. Then STOP and report. Wait for my
    "continue" before skill2.md.

Begin with STEP 0, then implement docs/skills/phase3/skill1.md.
```

---

## 2. Driving it skill-by-skill

Same loop as before. After each skill, send:

```
Verification looks good. Now execute docs/skills/phase3/skill<N>.md:
read it fully, re-read the design-spec sections it references, adapt its template
code to our actual package paths, follow its Tasks, run its Verification commands,
check its Definition of Done, then stop and report. Honor docs/skill.md §2 and
reuse the Phase 1/2 delivery layer.
```

Phase 3 order: skill1 (data + storage port) → skill2 (upload + queue) →
skill3 (FFmpeg→HLS worker) → skill4 (signed delivery) → skill5 (HTTP + wiring +
smoke).

---

## 3. Phase-3-specific watch-outs

- **The template code is a scaffold, not gospel.** It uses `...` placeholders and
  assumes module/package paths — Antigravity must wire it to the real ones and
  fill the elided bits. Don't let it paste blindly and leave `...` in the tree.
- **Ports vs impls.** Storage (`ObjectStore`), queue (`JobQueue`), transcoder
  (`Transcoder`), and URL signing (`URLSigner`) are interfaces. If business logic
  ever imports `aws-sdk` or `os/exec` directly, that's an architecture break.
- **No public buckets.** Dev uses MinIO; production must use CloudFront signed
  URLs/cookies so HLS **segments** (not just the manifest) are protected. Skill 4
  forces a documented decision here — make sure it's honored.
- **Single-use downloads** use Redis `GETDEL`. Streaming is *not* single-use.
- **FFmpeg flags may need tuning** for real media; keep them isolated in the
  runner. The `//go:build ffmpeg` integration test is the proof they work.
- **Enrollment is still Phase 4.** Paid-content access is gated to owner/admin
  with `TODO(phase4)` hooks; free/preview content is open to authenticated
  students. Don't invent an enrollment system here.
- **Don't skip ahead** into payments/notifications — Phase 4+.

---

## 4. When Phase 3 is finished

Skill 5's Definition of Done is the Phase-3 gate. When it passes:

1. Antigravity reports "Phase 3 complete" with a summary + smoke-test result.
2. Come back to Claude and say **"author Phase 4 skills"** — I'll write
   `docs/skills/phase4/` for **Payments & Enrollments** (Chapa/Telebirr/Stripe
   behind a provider port, idempotency keys, webhook verify/dedup, the enrollment
   flow) and the replacement of every `TODO(phase4)` hook left in Phases 2–3
   (quiz/exam taking + paid video access) with real enrollment checks (spec §10).
3. Do not scaffold Phase 4 code before its skills are written and reviewed.

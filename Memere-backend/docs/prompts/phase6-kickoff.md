# Phase 6 — Antigravity Kickoff Prompt (Final Phase)

Use this **only after Phase 5 is complete and green**. Phase 6 is the **final
phase**: it adds **no product features** — it makes the feature-complete monolith
**production-ready and scalable** (spec §12, §13, §3.1–3.2).

---

## 1. The first prompt to paste into Antigravity

```
We are starting Phase 6 (final) of the Memere backend: production hardening,
observability, performance, CI/CD, and Kubernetes. Phases 1–5 are complete and
green. Work strictly from the repo docs.

STEP 0 — Read, in this order, before writing any code:
  1. docs/skill.md                          (master index — rules & phase map)
  2. docs/memere_Design_Specification.md §12 (Scalability) + §13 (DevOps/CI/CD) + §3.1–3.2
  3. docs/skills/phase6/skill1.md            (first Phase 6 build skill)

GROUND RULES (docs/skill.md §2 — still apply):
  - No new product features. Do not change business behavior; only add
    observability, security, performance, and deployment infrastructure.
  - Measure before optimizing: Skill 1 (observability) comes first so Skills 2–3
    are driven by real metrics/traces, not guesses.
  - Keep the clean-architecture boundaries intact — they are what make the
    microservices extraction (Skill 5) mechanical. Do not introduce cross-domain
    coupling.
  - Secrets never in code/git; redaction of secrets/PII in logs is mandatory.
  - Template code in skills is a scaffold — adapt to real package paths; no `...`
    left behind.

HOW TO WORK:
  - Execute ONLY skill1.md now. Run Verification; check Definition of Done; STOP
    and report. Wait for "continue" before skill2.md.

Begin with STEP 0, then implement docs/skills/phase6/skill1.md.
```

---

## 2. Per-skill continue prompt

```
Verification looks good. Now execute docs/skills/phase6/skill<N>.md:
read it fully, re-read the referenced spec sections, adapt template code to our
package paths, follow its Tasks, run Verification, check Definition of Done, then
stop and report. Honor docs/skill.md §2 and change NO business behavior.
```

Order: skill1 (observability) → skill2 (security hardening) → skill3 (performance
+ caching) → skill4 (CI/CD + containers) → skill5 (k8s + microservices plan).

---

## 3. Phase-6 watch-outs

- **Observability first.** Don't let Antigravity jump to caching/perf before
  metrics+traces exist — Skills 2–3 depend on being able to measure.
- **Hardening is an audit, not just new code.** Skill 2 walks the §7.3 table and
  verifies each mitigation is actually present (several were built earlier).
- **Two images.** The transcode worker needs ffmpeg; the api image is distroless.
  Skill 4 builds separate `api` and `worker` images — don't bloat the api image.
- **Migrations run as a pre-rollout Job/init container**, never in app start
  (Skill 4/5), and must be backward-compatible (expand/contract) for zero-downtime.
- **Don't split into microservices.** Skill 5 produces a *plan* + the worker as
  the first separate deployment. Actual extraction happens later, when §12.2
  triggers fire. The monolith is correct until metrics say otherwise.
- **p95 < 200ms** is the load-test gate in Skill 3 — record results in
  `docs/perf.md`.

---

## 4. When Phase 6 is finished — the build is done

Skill 5's Definition of Done is the final gate. When it passes:

1. Report "Phase 6 complete — build done" with a summary + the load-test result and
   k8s deploy status.
2. The backend is production-ready: observable, hardened, performant,
   containerized, CI/CD-automated, running on Kubernetes, with `docs/scaling.md`
   and `docs/microservices-plan.md` guiding future growth.
3. Future work is product iteration (spec roadmap: AI tutor, adaptive learning,
   live sessions, Amharic localization, parent portal, B2B). Each new initiative
   can be authored as a fresh `docs/skills/phaseN/` set in this same format — ask
   Claude to author it when you're ready.

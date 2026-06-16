# Phase 4 — Kickoff Prompt (Memere Admin): Announcements + Hardening + Deployment

> **Prerequisite:** Phase 3 is green — payments, refunds, reconcile, revenue
> dashboards all work against the live backend with backend-sourced numbers.

Phase 4 finishes the product: the **Announcements** broadcast composer, a
**global hardening** pass (error/loading/empty states, 401-refresh, toasts), a
**security review**, and **build + deployment** runbooks. This is the FINAL phase.

---

## The prompt

```
Continue building the Memere Admin panel. Phase 3 is complete and green.

STEP 0 — Re-read:
  1. docs/skill.md §2 (Non-Negotiables) and §9 (moving between phases)
  2. docs/Memere_Admin_Design_Specification.md §4.7 (Announcements), §5 (auth/
     session), §6 (security), §9 (environment contract)
  3. ../Memere-backend/docs/Memere.postman_collection.json — "Admin" announcements
     endpoint; and ../Memere-backend/k8s + docs/scaling.md for deploy coordination

GROUND RULES (in addition to §2):
  - Broadcasting to "all" requires an explicit extra confirmation.
  - Security pass must prove: no token in client JS, httpOnly+Secure+SameSite
    cookies, CSP/security headers set, CSRF protection on mutations, deps scanned.
  - Deployment keeps API_BASE_URL + secrets server-only and coordinates CORS with
    the backend if any direct browser→API path exists (there should be none).

HOW TO WORK: execute docs/skills/phase4/skill1.md only, verify, check Definition
of Done, STOP and report. Continue skill-by-skill. When skill5 passes, the admin
panel build is DONE.
```

---

## Skill order (strictly in sequence) — FINAL

| # | File | Builds |
|---|---|---|
| 1 | `skills/phase4/skill1.md` | Announcements broadcast composer (segment targeting) + session history |
| 2 | `skills/phase4/skill2.md` | Global hardening: error boundaries, toasts, skeletons, empty states, 401→refresh→retry |
| 3 | `skills/phase4/skill3.md` | Security review: cookie flags, CSRF, CSP/headers, no-token audit, dependency scan |
| 4 | `skills/phase4/skill4.md` | Production build + Dockerfile (standalone), env contract, image hardening |
| 5 | `skills/phase4/skill5.md` | Deployment runbooks (Vercel + self-hosted), CORS coordination, final smoke test, README |

When skill 5's Definition of Done passes: **the Memere Admin panel is complete.**
Future work is product iteration (audit-log viewer, fine-grained admin RBAC,
Amharic localization, saved report exports) — add as a new `docs/skills/phaseN/`
set using this same format.

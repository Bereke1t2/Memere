# Phase 5 — Antigravity Kickoff Prompt

Use this **only after Phase 4 is complete and green**. Phase 5 builds **Progress
tracking, Notifications, Admin operations, and Certificates** (spec §11, §4.2.8)
and wires every `notify` no-op hook left in Phases 3–4 into a real system.

---

## 1. The first prompt to paste into Antigravity

```
We are starting Phase 5 of the Memere backend (Progress, Notifications, Admin,
Certificates). Phases 1–4 are complete and green. Work strictly from the repo docs.

STEP 0 — Read, in this order, before writing any code:
  1. docs/skill.md                          (master index — rules & phase map)
  2. docs/memere_Design_Specification.md §11 (Notifications) + §4.2.8 (progress) + §9.3
  3. docs/skills/phase5/skill1.md            (first Phase 5 build skill)

CARRY-OVER GROUND RULES (docs/skill.md §2 — never violate):
  - Notifications send ASYNC via the job queue — a slow/failing provider must never
    block or fail the originating API request.
  - Replace every notifyNoop/NoopNotifier across Phases 2–5 with the real Notifier;
    grep must show zero noop at the end.
  - All progress/notification/cert reads & writes are scoped to the authenticated
    user (no IDOR). Admin endpoints are admin-gated and audited.
  - Certificate downloads use short-lived signed URLs (reuse Phase 3 signer); no
    raw keys to clients. Public cert verification exposes minimal data only.
  - Clean Architecture: Notifier/PushSender/EmailSender/Renderer are PORTS;
    FCM/SendGrid/PDF-lib are infrastructure. Migrations ADDITIVE (after 0012).
  - Template code in skills is a scaffold — adapt to real package paths; no `...`
    left behind.

HOW TO WORK:
  - Execute ONLY skill1.md now. Run Verification; check Definition of Done; STOP
    and report. Wait for "continue" before skill2.md.

Begin with STEP 0, then implement docs/skills/phase5/skill1.md.
```

---

## 2. Per-skill continue prompt

```
Verification looks good. Now execute docs/skills/phase5/skill<N>.md:
read it fully, re-read the referenced spec sections, adapt template code to our
package paths, follow its Tasks, run Verification, check Definition of Done, then
stop and report. Honor docs/skill.md §2.
```

Order: skill1 (progress + streaks) → skill2 (notifications + wire all hooks) →
skill3 (admin ops + analytics) → skill4 (certificates + engagement sweeper) →
skill5 (HTTP + worker wiring + smoke).

---

## 3. Phase-5 watch-outs

- **Async or nothing:** the notification path writes the in-app row immediately
  and enqueues push/email — it never calls FCM/SendGrid inline on the request
  path. The smoke test proves a failing provider doesn't fail the API.
- **Hook completeness:** by end of skill2, video_ready (P3), purchase_confirmed +
  subscription_expired (P4), exam_graded (P2), certificate_ready + streak_warning
  (P5) all fire the real notifier.
- **Dev fallback:** FCM/SendGrid use a `LogSender` when keys are absent — Phase 5
  must run end-to-end locally without real provider keys.
- **Admin = audited:** every state-changing admin action writes an
  `admin_audit_log` row; broadcasts page/batch recipients (no unbounded list).
- **Streaks** use a configurable timezone (default Africa/Addis_Ababa); test day
  boundaries with the fake clock.

---

## 4. When Phase 5 is finished

1. Report "Phase 5 complete" + smoke result. The platform is now feature-complete
   for the MVP + growth scope.
2. Proceed to the **final phase, Phase 6** (already authored): use
   `docs/prompts/phase6-kickoff.md`. Phase 6 adds no features — it makes the system
   production-ready (observability, hardening, performance, CI/CD, k8s, and the
   microservices roadmap).

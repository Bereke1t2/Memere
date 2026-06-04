# Phase 4 — Antigravity Kickoff Prompt

Use this **only after Phase 3 is complete and green**. Phase 4 builds **Payments &
Enrollments** (spec §10) and replaces every `TODO(phase4)` access hook left in
Phases 2–3.

---

## 1. The first prompt to paste into Antigravity

```
We are starting Phase 4 of the Memere backend (Payments & Enrollments). Phases 1–3
are complete and green. Work strictly from the repo docs.

STEP 0 — Read, in this order, before writing any code:
  1. docs/skill.md                          (master index — rules & phase map)
  2. docs/memere_Design_Specification.md §10 (Payment System) + §4.2.7 + §7.3
  3. docs/skills/phase4/skill1.md            (first Phase 4 build skill)

CARRY-OVER GROUND RULES (docs/skill.md §2 — never violate):
  - ALL payments use a client Idempotency-Key; webhooks are signature-verified and
    deduplicated; fulfillment is ONE transaction that grants enrollment exactly
    once. No double-charge, no double-grant.
  - Money is decimal, never float. Never log provider secrets or raw payloads.
  - Replace every TODO(phase4) hook (Phase 2 quiz/exam taking, Phase 3 paid video)
    with the real access.Service — grep must show zero TODO(phase4) at the end.
  - Clean Architecture: PaymentProvider is a domain PORT; Chapa/Telebirr/Stripe
    are infrastructure impls. Migrations are ADDITIVE (after 0011).
  - Each skill contains TEMPLATE CODE — adapt it to real package paths; don't
    leave `...` placeholders in the tree.

HOW TO WORK:
  - Execute ONLY skill1.md now. Run its Verification commands; check its Definition
    of Done; STOP and report. Wait for "continue" before skill2.md.

Begin with STEP 0, then implement docs/skills/phase4/skill1.md.
```

---

## 2. Per-skill continue prompt

```
Verification looks good. Now execute docs/skills/phase4/skill<N>.md:
read it fully, re-read the referenced spec sections, adapt its template code to our
package paths, follow its Tasks, run its Verification commands, check its Definition
of Done, then stop and report. Honor docs/skill.md §2.
```

Order: skill1 (data + provider port) → skill2 (enrollment/access + remove hooks) →
skill3 (payment flow + webhooks) → skill4 (subscriptions + revenue) → skill5
(HTTP + webhook route + wiring + smoke).

---

## 3. Phase-4 watch-outs

- **The webhook route is special:** raw body, signature-verified, unauthenticated,
  outside JSON/auth middleware. Getting it wrong breaks signature verification.
- **Idempotency + dedup are the whole point** — the smoke test re-sends the same
  initiate (same key) and re-delivers the same webhook to prove no double
  charge/grant. Don't skip those assertions.
- **Coupons increment only on success**, inside the fulfillment transaction.
- **Active subscription = all-access** via `access.Service` — confirm quiz/exam/
  video access all flow through it now.
- Telebirr/Stripe may be sandbox/stub; Chapa is primary. All three satisfy the
  port and verify signatures.

---

## 4. When Phase 4 is finished

1. Report "Phase 4 complete" + smoke result.
2. Proceed to **Phase 5** (already authored): use
   `docs/prompts/phase5-kickoff.md`. Phase 5 wires the `notify` no-op hooks left
   here and in Phase 3 into a real notification system.

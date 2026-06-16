# Phase 3 — Kickoff Prompt (Memere Admin): Payments + Revenue + Subscriptions

> **Prerequisite:** Phase 2 is green — the `DataTable`, Users management, and
> Course moderation flows work against the live backend with confirmed destructive
> actions and cursor pagination.

Phase 3 builds the financial surface: **Payments** (list → detail → refund +
reconcile), the **Revenue** dashboard (totals, MRR, provider breakdown charts),
and a **Subscriptions/financial KPI** overview. This is the most sensitive area —
every number comes from the backend, never recomputed in the browser.

---

## The prompt

```
Continue building the Memere Admin panel. Phase 2 is complete and green.

STEP 0 — Re-read:
  1. docs/skill.md §2 (Non-Negotiables — especially #4 "no business logic in the
     client" and the money-formatting convention in §7)
  2. docs/Memere_Admin_Design_Specification.md §4.5 (Payments), §4.6 (Revenue),
     §2.2/§2.3 (payment + revenue DTOs; money fields arrive as strings/decimals)
  3. ../Memere-backend/docs/Memere.postman_collection.json — "Payments",
     "Revenue", and "Admin" (analytics) folders

GROUND RULES (in addition to §2):
  - NEVER compute or re-derive a financial figure in the browser. Display
    gross_revenue, mrr, refunded_amount, breakdown.gross/units exactly as returned.
  - Money values are strings/decimals from the backend — keep as strings for
    display; only parse to number for chart plotting, and label currency from the
    field, not a hardcoded "ETB".
  - Refund is admin-only and irreversible — strong confirm dialog, await backend,
    toast real result, refetch. Reconcile-pending shows the returned {reconciled}.

HOW TO WORK: execute docs/skills/phase3/skill1.md only, verify, check Definition
of Done, STOP and report. Continue skill-by-skill.
```

---

## Skill order (strictly in sequence)

| # | File | Builds |
|---|---|---|
| 1 | `skills/phase3/skill1.md` | Payments list: status filter, paginated table on `/admin/payments` |
| 2 | `skills/phase3/skill2.md` | Payment detail + refund + reconcile-pending toolbar action |
| 3 | `skills/phase3/skill3.md` | Revenue dashboard: overview KPIs, provider breakdown charts, date range |
| 4 | `skills/phase3/skill4.md` | Subscriptions / financial KPI overview tiles |
| 5 | `skills/phase3/skill5.md` | Money/date formatting audit + Phase 3 smoke test |

When skill 5 passes, report "Phase 3 complete" and move to
`docs/prompts/phase4-kickoff.md`.

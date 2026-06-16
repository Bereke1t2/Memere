# Phase 2 — Kickoff Prompt (Memere Admin): User Management + Course Moderation

> **Prerequisite:** Phase 1 is green — the Next.js app runs, an admin can log in
> (httpOnly cookies), the protected `(dashboard)` shell renders with sidebar/header,
> the server API client + Zod schemas exist, and the dashboard home shows live KPIs.

Phase 2 builds the first real operational screens: a reusable data table, the
**Users** management flow (list → detail → suspend/reactivate/change-role) and
**Course** moderation (list → detail → unpublish).

---

## The prompt

```
Continue building the Memere Admin panel. Phase 1 is complete and green.

STEP 0 — Re-read:
  1. docs/skill.md §2 (Non-Negotiables) and §5 (Phase 2 skill table)
  2. docs/Memere_Admin_Design_Specification.md §4.3 (Users) and §4.4 (Courses)
  3. ../Memere-backend/docs/Memere.postman_collection.json — the "Admin" folder
     (users + courses endpoints) and the User/Course DTO shapes

GROUND RULES: unchanged from Phase 1 (docs/skill.md §2). In particular:
  - All /admin/* calls are server-side through lib/api; Bearer from httpOnly cookie.
  - Cursor pagination (limit + after, read `next`). No offset, no fetch-all.
  - Every destructive action (suspend, reactivate, change-role, unpublish) shows
    a confirm dialog, awaits the backend, toasts the real result, invalidates the
    query. No optimistic lie on failure.
  - Map the error envelope `code` → friendly message.

HOW TO WORK: execute docs/skills/phase2/skill1.md only, run its verification,
check its Definition of Done, then STOP and report. Continue skill-by-skill on my
"continue".
```

---

## Skill order (strictly in sequence)

| # | File | Builds |
|---|---|---|
| 1 | `skills/phase2/skill1.md` | Reusable `DataTable` (TanStack Table): cursor pagination, filters, loading skeleton, empty state |
| 2 | `skills/phase2/skill2.md` | Users list: search, role filter, paginated table on `/admin/users` |
| 3 | `skills/phase2/skill3.md` | User detail + suspend/reactivate/change-role with confirm dialogs + mutations |
| 4 | `skills/phase2/skill4.md` | Courses list + unpublish-with-reason moderation |
| 5 | `skills/phase2/skill5.md` | Course detail (sections inspect) + Phase 2 polish & smoke test |

When skill 5's Definition of Done passes, report "Phase 2 complete" and move to
`docs/prompts/phase3-kickoff.md`.

# Phase 4 · Skill 5 — Deployment Runbooks + Final Smoke Test

> **Prerequisite:** Skill 4 done. Read [`docs/skill.md`](../../skill.md) §9,
> [`design spec`](../../Memere_Admin_Design_Specification.md) §9, and the backend
> `../Memere-backend/docs/scaling.md` for deploy coordination. This is the **final
> skill**.

---

## Goal

Document and validate two deployment paths (Vercel and self-hosted Docker/VPS),
coordinate CORS/origins with the backend, run the full end-to-end smoke test, and
finalize the README.

---

## Tasks

### 5.1 — `docs/DEPLOYMENT.md`

**Path A — Vercel (fastest):**
- Import the repo; set env (`API_BASE_URL`, `COOKIE_SECRET`, `NEXT_PUBLIC_APP_NAME`)
  as Vercel project env (server-side). Deploy. Set the custom domain
  (`admin.memere.app`).

**Path B — Self-hosted (Docker on a VPS):**
- Run the Phase-4 image behind nginx + TLS (Certbot), env injected at runtime.
- Example `docker run` / compose snippet + nginx reverse-proxy block.

### 5.2 — Backend coordination

- Since all API calls are **server-side**, browser CORS is not required. If any
  direct browser→API call exists, add `https://admin.memere.app` to the backend's
  `CORS_ALLOWED_ORIGINS` (see backend `.env` / `k8s/configmaps/app-config.yaml`).
- Confirm the backend is reachable from the admin host over HTTPS.

### 5.3 — Final smoke test (`docs/SMOKE.md` — full run)

Against staging/production with a real admin account:
1. Login (admin only) → cookies set, no token in JS.
2. Dashboard KPIs + charts live.
3. Users: list/filter/paginate/detail/suspend/reactivate/role.
4. Courses: list/detail/unpublish.
5. Payments: list/filter/detail/refund/reconcile.
6. Revenue: KPIs + charts for a range.
7. Announcements: send to a segment (double-confirm "all").
8. Token refresh + logout.

### 5.4 — README finalize

- Architecture summary, env contract, local dev, both deploy paths, link to the
  backend repo + Postman collection.

---

## Definition of Done

- [ ] `docs/DEPLOYMENT.md` covers Vercel + self-hosted with env + TLS + CORS notes.
- [ ] Backend origin reachable; CORS coordinated if needed.
- [ ] Full smoke test passes end-to-end on a deployed environment.
- [ ] README is complete and accurate.
- [ ] `pnpm build` + `pnpm lint` + `pnpm typecheck` pass.

## Verification commands

```bash
# After deploy:
curl -sI https://admin.memere.app/login        # 200, security headers present
# Run the full docs/SMOKE.md checklist against the deployed panel.
```

---

## 🎉 Phase 4 complete — the admin panel is DONE

All four phases delivered:

1. **Foundation** — Next.js + auth (server-only tokens) + guarded shell + live dashboard.
2. **User & Course management** — paginated tables, confirmed moderation actions.
3. **Payments & Revenue** — refunds, reconcile, backend-sourced financials.
4. **Announcements + Hardening + Deployment** — broadcast, security, container, runbooks.

Future work (new `docs/skills/phaseN/` sets): audit-log viewer, fine-grained admin
RBAC, Amharic localization, saved/exported reports, real-time updates.

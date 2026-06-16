# Phase 4 · Skill 3 — Security Review

> **Prerequisite:** Skill 2 done. Read [`docs/skill.md`](../../skill.md) §2 (all),
> [`design spec`](../../Memere_Admin_Design_Specification.md) §6 (security).

---

## Goal

Prove the panel meets its security bar: tokens never in client JS, hardened
cookies, CSRF protection on mutations, security headers/CSP, and a clean
dependency scan.

---

## Tasks

### 3.1 — No-token-in-client audit

- Grep the client bundle / source for token usage: no `localStorage`/
  `sessionStorage` tokens, no `document.cookie` token reads, no `NEXT_PUBLIC_`
  containing a secret or `API_BASE_URL`. Confirm in built output, not just source.

### 3.2 — Cookie flags

- Verify `mm_access` / `mm_refresh` are `httpOnly`, `Secure` (prod),
  `SameSite=Lax`, scoped `path=/`, with sane max-age. Confirm via response headers.

### 3.3 — CSRF on mutations

- All state-changing requests go through same-site Route Handlers. Add a
  double-submit CSRF token (or `Origin`/`Referer` check) for POST Route Handlers
  to defend against cross-site form posts.

### 3.4 — Security headers / CSP (`next.config.ts`)

- Add headers: `Content-Security-Policy` (restrict script/connect/img/style to
  self + needed origins), `X-Frame-Options: DENY`, `X-Content-Type-Options:
  nosniff`, `Referrer-Policy: strict-origin-when-cross-origin`,
  `Permissions-Policy`.

### 3.5 — Dependency scan

- Run `pnpm audit` (and/or `osv-scanner`); resolve high/critical. Document any
  accepted advisories.

---

## Definition of Done

- [ ] Built client contains no token and no `API_BASE_URL`/secret.
- [ ] Cookies carry httpOnly + Secure + SameSite + path; verified in headers.
- [ ] Mutation Route Handlers are CSRF-protected.
- [ ] CSP + security headers present on responses.
- [ ] `pnpm audit` clean of high/critical (or documented).
- [ ] `pnpm build` passes with headers applied.

## Verification commands

```bash
pnpm build && pnpm start
curl -sI http://localhost:3000/ | grep -iE 'content-security-policy|x-frame|x-content|referrer'
pnpm audit
# DevTools: confirm cookies' flags and that no token is reachable from JS.
```

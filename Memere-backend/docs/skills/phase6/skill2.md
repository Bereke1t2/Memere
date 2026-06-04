# Phase 6 · Skill 2 — Security Hardening & API Gateway Concerns

> **Prerequisite:** Phase 6 Skill 1 done (observability — so you can measure
> hardening impact). Read [`docs/skill.md`](../../skill.md) §2.
>
> **Spec references:** `memere_Design_Specification.md` §7.3 (attack vectors &
> mitigations — the full table), §3.2 (API gateway responsibilities), §1.5
> (security NFRs), README "Security Measures".

---

## Goal

Close the gap between "works" and "safe in production". Implement the gateway-level
hardening from §7.3 that the monolith must own until a real gateway fronts it:
security headers, tightened/segmented rate limiting, request size & timeout limits,
secret management, CORS lockdown, and dependency/secret scanning in CI.

Walk the **§7.3 attack-vector table** and verify each mitigation is actually
present (some were built earlier — this skill audits + completes them).

---

## §7.3 audit checklist (verify each is real, not assumed)

| Vector | Mitigation | Built in | This skill |
|---|---|---|---|
| Brute-force login | 5/15min/IP + lockout | P1 mw | verify + add account lockout after 10 fails |
| JWT theft | short access TTL, rotation | P1/P3 | verify; add token-revocation list option |
| SQL injection | parameterized (sqlc) | all | verify no raw string SQL anywhere |
| XSS (admin/web) | output encoding, CSP | — | add CSP + security headers |
| IDOR | filter by user_id | all | audit every handler/usecase |
| Payment replay | idempotency + webhook dedup | P4 | verify |
| Video hotlinking | signed 2h URLs | P3 | verify; add segment protection note |
| Exam answer leak | server-side grade only | P2 | verify (grep) |
| DDoS | rate limit + edge | P1 | tighten + per-route limits |
| Credential stuffing | breach detection, alerting | — | add suspicious-login metric/alert |

---

## Tasks

### 2.1 — Security headers middleware (`internal/delivery/middleware/security.go`)

```go
func SecurityHeaders() gin.HandlerFunc {
    return func(c *gin.Context) {
        h := c.Writer.Header()
        h.Set("X-Content-Type-Options", "nosniff")
        h.Set("X-Frame-Options", "DENY")
        h.Set("Referrer-Policy", "no-referrer")
        h.Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'")
        h.Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
        c.Next()
    }
}
```

(API is JSON-only, so CSP can be very strict. HSTS assumes TLS-terminated edge.)

### 2.2 — Request limits

- **Body size cap** middleware: `http.MaxBytesReader` (e.g. 1 MB for JSON routes;
  the upload route uses pre-signed S3 so it never proxies large bodies — keep its
  JSON small too). Webhook already capped (Phase 4).
- **Server timeouts** (verify Phase 1 set them): `ReadHeaderTimeout`,
  `ReadTimeout`, `WriteTimeout`, `IdleTimeout`. Add per-handler `context` deadlines
  for slow downstreams (provider calls) so a hung provider can't pin a goroutine.
- **Max concurrent in-flight** option (graceful 503 under overload).

### 2.3 — Rate limiting v2 (`internal/delivery/middleware/ratelimit.go`)

Replace the Phase 1 limiter with a segmented Redis sliding-window:
- global per-IP default,
- stricter on `/auth/login`, `/auth/register`, `/auth/refresh`,
  `/payments/initiate`,
- per-**user** limits (not just per-IP) on authenticated expensive routes,
- webhook route: per-provider limit.
Return `429` + `Retry-After`; emit a `rate_limit_block_total` metric (Skill 1).

### 2.4 — Account lockout + suspicious login (extend auth)

- After N (e.g. 10) consecutive failed logins for an account, lock for a cooldown
  (store counter in Redis keyed by user, reset on success). Distinct from the
  IP-based limiter. Return a generic error (don't reveal lockout to avoid user
  enumeration; §7.3).
- Emit a `suspicious_login_total` metric and a Sentry breadcrumb on lockout —
  basis for the §13.3 alert.
- Optional: an access-token denylist in Redis (`revoked:jti`) checked by
  `RequireAuth`, so logout/suspend can revoke access tokens before their 15-min
  expiry. Wire admin `SuspendUser` (Phase 5) to revoke active tokens.

### 2.5 — CORS lockdown

Make CORS strict in production: explicit allowed origins from config (the Flutter
app's web build / admin domain), specific methods/headers, credentials handling
correct. Dev stays permissive. Fail closed if `APP_ENV=production` and no origins
configured.

### 2.6 — Secret management

- No secrets in code/repo (verify). All from env; in k8s they come from Secrets
  (Skill 5) / a secrets manager.
- Add startup validation: in `production`, refuse to boot if `JWT_SECRET` is weak
  (length check), or if any required provider secret is missing while its provider
  is enabled.
- Document a secret-rotation note for `JWT_SECRET` (support two valid secrets
  during rotation — optional but note it).

### 2.7 — CI security scanning

Add to the pipeline (Skill 4 builds the full CI): `govulncheck`, `gosec`,
`gitleaks` (secret scan), and `go mod verify`. Fail the build on high-severity
findings.

### 2.8 — IDOR + injection audit

- Grep/review every handler: confirm object access resolves ownership through the
  usecase against the authenticated user (no trust of client IDs). Add a test per
  resource that user B cannot touch user A's object.
- Confirm **zero** raw/concatenated SQL (everything via sqlc); `grep` for
  `fmt.Sprintf` near SQL.
- Confirm **zero** answer-key exposure (re-run the Phase 2 grep).

### 2.9 — Tests

Security headers present on every response; oversize body → 413; login lockout
after N fails (generic error); revoked token rejected; CORS rejects unknown origin
in prod mode; production boot fails on weak/missing secrets; IDOR tests per
resource.

---

## Definition of Done

- [ ] `go build ./...` clean; `golangci-lint run` clean; `go test ./...` passes.
- [ ] Every §7.3 row is verified present (the audit table is fully checked).
- [ ] Security headers + strict CSP/HSTS on all responses; body-size + timeouts
      enforced.
- [ ] Segmented rate limits (IP + user + route) with `429`+`Retry-After` and a
      metric; login lockout after N failures (generic error).
- [ ] Token revocation works (logout/suspend invalidates access before expiry).
- [ ] CORS fails closed in production without configured origins; weak/missing
      secrets block production boot.
- [ ] CI runs `govulncheck`/`gosec`/`gitleaks` and fails on high severity.
- [ ] IDOR tests per resource pass; no raw SQL; no answer-key exposure.

## Verification commands

```bash
go build ./... && golangci-lint run && go test ./...
govulncheck ./... ; gosec ./... ; gitleaks detect --no-banner
grep -rn 'fmt.Sprintf' internal/repository | grep -i 'select\|insert\|update\|delete' && echo "REVIEW: raw SQL?" || echo OK
curl -sI localhost:8080/healthz | grep -E 'X-Content-Type-Options|Content-Security-Policy'
```

## Hand-off to Skill 3

The app is hardened. Skill 3 tackles **performance & scalability**: caching
(cache-aside with Redis), query/index tuning, the leaderboard, connection pool
sizing, and load testing against the §1.5 SLOs (p95 < 200ms).

# Phase 1 · Skill 3 — Authentication (JWT + Refresh Tokens + RBAC)

> **Prerequisite:** Skills 1–2 done (schema migrated, sqlc generates, entities &
> repository interfaces exist). Read [`docs/skill.md`](../../skill.md) first, and
> re-read the **Non-Negotiables** — auth is where IDOR, token leakage, and
> password handling rules bite hardest.
>
> **Spec references:** `memere_Design_Specification.md` §7 (entire Auth &
> Security section — §7.1 JWT flow sequence, §7.2 RBAC, §7.3 attack vectors),
> README "Authentication & Security", §4.2.1 users table.

---

## Goal

Implement the complete authentication vertical as **usecases + supporting pkg +
repository implementations** — everything except the HTTP handlers (those land in
Skill 5). After this skill, the auth usecases can be unit-tested in isolation:
register, login, refresh, logout, plus password hashing, JWT issue/verify, and
the refresh-token store in Redis + Postgres.

---

## The JWT flow we are implementing (spec §7.1)

1. **Login** → verify bcrypt password → issue **access token** (15 min) +
   **refresh token** (30 days). Store a **hash** of the refresh token in
   Postgres (`auth.refresh_tokens`) and in Redis (`session:{user_id}` →
   token-hash, TTL 30d). Return both tokens + sanitized user.
2. **Authenticated request** → access token in `Authorization: Bearer`.
3. **Access expired** → client calls `/auth/refresh` with the refresh token →
   server checks the hash is present & not revoked in Redis/PG → issues a new
   access token (and rotates the refresh token — see 3.5).
4. **Logout** → revoke the refresh token (delete from Redis, set `revoked_at` in
   PG).

---

## Tasks

### 3.1 — `pkg/password`

- `Hash(plain string) (string, error)` using `bcrypt.GenerateFromPassword` with
  cost ≥ 12.
- `Compare(hash, plain string) error` (constant-time via bcrypt).
- Never log or return the plaintext anywhere.

### 3.2 — `pkg/jwt`

- A `Manager` constructed from `secret`, `accessTTL`, `refreshTTL`, `issuer`.
- `Claims` embedding `jwt.RegisteredClaims` plus custom: `UserID uuid.UUID`,
  `Role string`, `TokenType string` ("access" | "refresh").
- `GenerateAccessToken(user) (string, error)` — short TTL, signed HS256.
- `GenerateRefreshToken(user) (string, error)` — long TTL. (Refresh token is a
  signed JWT *and* we store a hash of it; both checks must pass.)
- `Verify(token string) (*Claims, error)` — validates signature, expiry, issuer,
  and `TokenType`. Returns `apperror.Unauthorized` on any failure (never leak why).
- A `HashToken(token string) string` helper (SHA-256 hex) used to store refresh
  tokens without storing the raw token. **Raw refresh tokens are never persisted.**

### 3.3 — Repository implementations

Implement the interfaces from Skill 2 in `internal/repository/`:

- `postgres/user_repository.go` — wraps the sqlc-generated queries, maps sqlc
  models ↔ domain entities. Handle `pgx.ErrNoRows` → `apperror.NotFound`. Map
  unique-violation on email (`23505`) → `apperror.Conflict("EMAIL_TAKEN")`.
- `postgres/refresh_token_repository.go` — Create / FindByHash / Revoke /
  RevokeAllForUser / DeleteExpired.
- `redis/session_repository.go` — a thin store:
  `SetSession(ctx, userID, tokenHash, ttl)`, `GetSession(ctx, userID)`,
  `DeleteSession(ctx, userID)`. Key format `session:{user_id}`.

Mapping helpers (sqlc model → entity) live next to each repo. **The usecase layer
must only ever see domain entities**, never sqlc types (dependency rule).

### 3.4 — Auth usecases (`internal/usecase/auth/`)

Construct a `Service` (or `UseCase`) struct holding the repository interfaces +
`jwt.Manager` + config. Methods:

- `Register(ctx, RegisterInput) (*entity.User, error)`
  - validate input (email format, password strength, names present, role —
    default `student`; reject self-assigning `admin`).
  - check email not taken; hash password; create user with
    `is_email_verified=false`; generate an email-verification token (store on
    user). Return sanitized user. (Sending the email is a later phase — just
    persist the token.)
- `Login(ctx, LoginInput) (*AuthTokens, *entity.User, error)`
  - find by email (generic `INVALID_CREDENTIALS` error if not found — **don't**
    reveal whether the email exists, per §7.3).
  - `password.Compare`; on success issue access+refresh, store refresh hash in
    PG + Redis, update `last_login_at`.
- `Refresh(ctx, refreshToken string) (*AuthTokens, error)`
  - `jwt.Verify` (must be type "refresh") → hash → confirm present & not revoked
    in Redis (fast path) and PG (source of truth) → issue new access token →
    **rotate** refresh token (3.5).
- `Logout(ctx, userID, refreshToken string) error`
  - revoke in PG + delete Redis session.
- `AuthTokens` = `{ AccessToken, RefreshToken string, ExpiresIn int }`.

**Sanitize user** before returning anywhere: strip `passwordHash`, reset tokens,
verification tokens. Provide an `entity.User.Sanitized()` or a delivery-layer DTO
(prefer the DTO in Skill 5; for now ensure usecases never return the hash to
callers that serialize it).

### 3.5 — Refresh-token rotation

On every successful `Refresh`: revoke the old refresh token (PG `revoked_at` +
Redis overwrite) and issue a fresh one. This limits replay if a refresh token
leaks (§7.3 "Token rotation").

### 3.6 — RBAC helper

Define role constants (already on `entity.Role` from Skill 2) and a permission
helper the middleware will use in Skill 5:
`func (r Role) Can(action) bool` or a simpler `RequireRole(claims, allowed...)`.
Map the spec §7.2 matrix:

| Role | Can |
|---|---|
| student | read enrolled courses, submit attempts, view own progress, purchase |
| teacher | student perms + create/edit *own* courses, view own-course analytics |
| admin | everything + manage all users/courses, payments |

Phase 1 only needs the *mechanism*; course-ownership enforcement is applied in
Skill 4/5.

### 3.7 — Unit tests

- `usecase/auth` tests with **mocked repository interfaces** (hand-written fakes
  or `gomock`/`mockery` — prefer simple hand-written fakes to avoid tooling).
  Cover: register happy path, duplicate email, login success, wrong password,
  refresh success, refresh with revoked token, logout.
- `pkg/jwt` tests: round-trip generate→verify, expired token rejected, wrong
  type rejected, tampered signature rejected.
- `pkg/password`: hash≠plain, compare success/failure.

---

## Security checklist (must all hold — from §7.3)

- [ ] Login returns identical error/timing for "no such email" and "wrong
      password" (`INVALID_CREDENTIALS`).
- [ ] Raw refresh tokens are never stored — only SHA-256 hashes.
- [ ] Access token TTL ≈ 15 min; refresh ≈ 30 days (from config, not hardcoded).
- [ ] Refresh tokens rotate on use; old ones are revoked.
- [ ] Passwords hashed with bcrypt cost ≥ 12; never logged.
- [ ] `passwordHash` and any token field never appear in a returned/serialized
      user.
- [ ] JWT verification failures are opaque (`UNAUTHORIZED`, no internal detail).

## Definition of Done

- [ ] `go build ./...` clean; `golangci-lint run` clean.
- [ ] `go test ./internal/usecase/auth/... ./pkg/jwt/... ./pkg/password/...`
      passes.
- [ ] Repository implementations satisfy the Skill 2 interfaces (compile-time
      assert: `var _ repository.UserRepository = (*postgres.UserRepo)(nil)`).
- [ ] Every item in the Security checklist above is verified by a test or code
      review.

## Verification commands

```bash
go build ./...
go test ./internal/usecase/auth/... ./pkg/... -v
golangci-lint run
grep -rn 'PasswordHash' internal/usecase internal/delivery && echo "review: ensure never serialized"
```

## Hand-off to Skill 4

Auth logic is ready and tested in isolation. Skill 4 builds the **course**
vertical (CRUD + sections + lessons usecases and repository implementations),
reusing the same patterns: interfaces from Skill 2, sqlc queries, domain-only
usecases. Skill 5 then exposes both auth and course over HTTP.

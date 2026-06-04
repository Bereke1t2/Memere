# Phase 1 · Skill 1 — Project Foundation & Infrastructure

> **Prerequisite:** You have read [`docs/skill.md`](../../skill.md) (the master
> index) and the **Non-Negotiables** in Section 2 there. This is the first build
> skill — nothing exists yet except docs.
>
> **Spec references:** `memere_Design_Specification.md` §5.2 (folder structure),
> §3.3 (infrastructure components), §12 (scaling — we target the monolith tier),
> README "Environment Variables" table.

---

## Goal

Stand up an empty-but-runnable Go service: it compiles, connects to Postgres and
Redis, exposes a `/health` endpoint, and has the full clean-architecture folder
skeleton plus Docker Compose for local infra. **No business logic yet.**

By the end, `make run` starts a server on `:8080` and `GET /health` returns
`200 {"status":"ok"}` with live DB + Redis pings.

---

## Pre-flight

⚠️ **Go version:** The dev machine reported Go 1.20.5, but we target **1.22+**.
Run `go version`; if below 1.22, install/upgrade Go before continuing. `go.mod`
will declare `go 1.22`.

---

## Tasks

### 1.1 — Initialize the module

- `go mod init github.com/Bereke1t2/Memere/memere-backend` (match the repo path;
  if the user has a different module path preference, ask).
- Set `go 1.22` in `go.mod`.

### 1.2 — Create the full directory skeleton

Create every directory from the master index Section 4 layout, each with a
`.gitkeep` (or a real file where this skill produces one). At minimum:

```
cmd/api/
internal/domain/entity/      internal/domain/repository/
internal/usecase/auth/       internal/usecase/course/
internal/repository/postgres/ internal/repository/redis/
internal/delivery/http/      internal/delivery/middleware/
internal/infrastructure/database/ internal/infrastructure/cache/
pkg/jwt/ pkg/password/ pkg/validator/ pkg/pagination/ pkg/apperror/
config/ migrations/ db/query/
```

### 1.3 — Typed configuration (`config/config.go`)

Load **all** env vars from the README table into a typed `Config` struct. Group
them: `App`, `DB`, `Redis`, `JWT`, plus placeholders (commented or zero-value) for
`AWS`, `Chapa`, `Stripe`, `FCM`, `SendGrid` that later phases will use.

- Provide `Load() (*Config, error)` that reads env (use
  `github.com/kelseyhightower/envconfig` or hand-rolled `os.Getenv` with
  defaults — pick envconfig for less boilerplate).
- Parse durations (`JWT_ACCESS_TTL=15m`, `JWT_REFRESH_TTL=720h`) into
  `time.Duration`.
- **Fail fast:** return an error if any required var (DB_*, JWT_SECRET) is missing.
- Build a Postgres DSN helper: `cfg.DB.DSN()` →
  `postgres://user:pass@host:port/name?sslmode=...`.

Required env keys (from README — reproduce in `.env.example`):

```
APP_ENV, APP_PORT,
DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, DB_SSL_MODE,
REDIS_HOST, REDIS_PORT, REDIS_PASSWORD,
JWT_SECRET, JWT_ACCESS_TTL, JWT_REFRESH_TTL,
AWS_S3_BUCKET, AWS_REGION,
CHAPA_SECRET_KEY, STRIPE_SECRET_KEY, FCM_SERVER_KEY, SENDGRID_API_KEY
```

### 1.4 — Postgres connection (`internal/infrastructure/database/postgres.go`)

- Use `pgxpool` (pgx v5). Function `Connect(ctx, cfg) (*pgxpool.Pool, error)`.
- Sensible pool config: `MaxConns` ~ 10–25, `MaxConnIdleTime`, `HealthCheckPeriod`.
- Run `pool.Ping(ctx)` on startup; fail fast if unreachable.
- Provide a `Close()` path.

### 1.5 — Redis connection (`internal/infrastructure/cache/redis.go`)

- Use `go-redis/v9`. Function `Connect(ctx, cfg) (*redis.Client, error)`.
- `Ping` on startup; fail fast.

### 1.6 — `pkg/apperror` (typed errors + HTTP mapping)

This is foundational — every later skill uses it. Implement:

- An `AppError` type: `Code string` (e.g. `RESOURCE_NOT_FOUND`), `Message string`,
  `HTTPStatus int`, `Details map[string]any`, and an unexported wrapped `err`.
- Constructors for the common cases: `NotFound`, `Unauthorized`, `Forbidden`,
  `BadRequest`, `Conflict`, `Internal`, `TooManyRequests`, `Validation(details)`.
- `Error() string` + `Unwrap() error`.
- The JSON shape on the wire must be exactly:
  `{ "code": "...", "message": "...", "details": {} }` (per spec §5.4).

### 1.7 — Minimal HTTP bootstrap (`internal/delivery/http/server.go` + health)

- A tiny Gin engine with **only**: a recovery middleware, a request-logger stub,
  and `GET /health`.
- `/health` pings DB and Redis; returns `200 {"status":"ok"}` if both healthy,
  `503` otherwise. (Full middleware stack & routes come in Skill 5 — keep this
  minimal but real.)

### 1.8 — Entry point (`cmd/api/main.go`)

Wire it together with graceful shutdown:

1. `config.Load()`
2. connect Postgres, connect Redis (defer closes)
3. build the minimal Gin engine + `/health`
4. `http.Server` with timeouts; **keep-alive / read / write timeouts set**
5. Run in a goroutine; listen for `SIGINT`/`SIGTERM`; `Shutdown(ctx)` with a
   timeout.

### 1.9 — `docker-compose.yml` (local infra)

Services: `postgres:15` and `redis:7`. Postgres env from `.env`. Expose
`5432`/`6379`. Named volumes for persistence. A healthcheck on each. (App itself
can run on host via `make run` during dev — containerizing the app is Skill 5 /
later.)

### 1.10 — `Dockerfile` (multi-stage, for the app)

- Stage 1: `golang:1.22` build → static binary (`CGO_ENABLED=0`).
- Stage 2: `gcr.io/distroless/static` or `alpine`; copy binary; non-root user;
  `EXPOSE 8080`; `ENTRYPOINT ["/app"]`.

### 1.11 — `Makefile`

Targets (the master index references these):

```
make run            # go run ./cmd/api
make build          # go build -o bin/api ./cmd/api
make tidy           # go mod tidy
make lint           # golangci-lint run
make test           # go test ./...
make up             # docker-compose up -d
make down           # docker-compose down
make migrate-up     # (stub now; real in Skill 2)
make migrate-down   # (stub now; real in Skill 2)
make sqlc           # (stub now; real in Skill 2)
```

### 1.12 — `.env.example`, `.gitignore`, `golangci-lint`

- `.env.example` with every key from 1.3 (placeholder values, never real secrets).
- `.gitignore`: `bin/`, `.env`, `*.out`, coverage files, OS junk.
- `.golangci.yml`: enable `govet, staticcheck, errcheck, gosimple, ineffassign,
  unused, gofmt, goimports`.

---

## Definition of Done

- [ ] `go build ./...` succeeds with zero errors.
- [ ] `go vet ./...` and `golangci-lint run` are clean.
- [ ] `make up` brings Postgres + Redis healthy.
- [ ] `make run` starts the server; logs show successful DB and Redis pings.
- [ ] `curl localhost:8080/health` → `200 {"status":"ok"}`.
- [ ] Stopping Postgres makes `/health` return `503` (proves real pings).
- [ ] `Ctrl-C` triggers graceful shutdown (no panic, server closes cleanly).
- [ ] `.env.example` lists every config key; no secrets committed.
- [ ] All folders from master index Section 4 exist.

## Verification commands

```bash
go build ./...
golangci-lint run
make up
make run &           # then, in another shell:
curl -i localhost:8080/health
docker compose stop postgres && curl -i localhost:8080/health   # expect 503
docker compose start postgres
```

## Hand-off to Skill 2

You now have a running shell with DB + Redis + config + error type. Skill 2 adds
the **data model**: domain entities, repository interfaces, SQL migrations, and
sqlc codegen. Do not write any auth or course logic yet.

# Memere Backend — Antigravity Skill Index (Master)

> **Read this file first.** It is the orchestration map for building the Memere
> (ExamPrep) Go backend with Antigravity. It tells you *what order* to execute
> skills in, *what rules never change*, and *how phases relate*. Every individual
> skill file assumes you have read and internalized this master file.

---

## 1. What we are building

A Go backend for **Memere (ExamPrep)** — a Grade 12 university-entrance-exam prep
platform for Ethiopian students. The full product vision lives in
[`docs/memere_Design_Specification.md`](./memere_Design_Specification.md). That
document is the **single source of truth**. If a skill file and the design spec
ever disagree, the design spec wins — stop and flag the conflict.

### Locked technical decisions (do not deviate without asking)

| Decision | Choice | Reason |
|---|---|---|
| Build shape (Phase 1) | **Modular monolith** — one Go binary, clean-architecture modules | Matches the spec's "1K users = monolith" tier; fastest path to a working MVP. Microservice split is a later phase. |
| HTTP framework | **Gin** | Spec allows Gin or Echo; we standardize on Gin. |
| DB driver | **pgx (v5)** | Highest-performance Postgres driver in Go. |
| Query layer | **sqlc** | Generates type-safe Go from raw SQL; satisfies the spec's "parameterized queries only" security rule with compile-time checks. |
| Migrations | **golang-migrate** | Spec-mandated. |
| Cache/session store | **Redis 7 (go-redis v9)** | Spec-mandated. |
| Auth | **JWT (golang-jwt v5) + bcrypt** | Access 15 min, refresh 30 days. |
| Config | **env vars via a typed Config struct** (envconfig or manual) | 12-factor. |
| Go version | **1.22+** | ⚠️ The dev machine currently has Go 1.20.5 — upgrade before starting Skill 1. |

---

## 2. The Non-Negotiables (apply to EVERY skill, every phase)

These come straight from the design spec's "Architecture Non-Negotiables" and
security sections. Violating any of these is a build-breaking bug:

1. **Correct answers NEVER leave the server.** Quiz/exam answer keys are never in
   any API response. Grading is server-side only.
2. **Exam/quiz timers are enforced server-side.** `started_at` lives in the DB;
   the client timer is display-only. The server auto-submits at expiry.
3. **Video access is via pre-signed CDN URLs only** (2-hour expiry). No public S3.
4. **All payments use idempotency keys.** Webhooks are deduplicated.
5. **Soft deletes only.** Every user-facing table has `deleted_at`; never hard
   `DELETE`. Every read query filters `deleted_at IS NULL`.
6. **HTTPS only** at the gateway (HTTP → HTTPS redirect).
7. **Every data query filters by the authenticated `user_id`** to prevent IDOR.
   Never trust an ID from the client as authorization.
8. **Never log** raw passwords, tokens, or payment card data.
9. **UUID primary keys everywhere** (`gen_random_uuid()`), never sequential ints.
10. **Clean Architecture dependency rule:** dependencies point *inward only*.
    `domain` imports nothing external. `usecase` imports only `domain`. Handlers
    and repositories import `usecase`/`domain`. Infrastructure is the outermost
    ring. A `domain` entity must have **no** db/json struct tags.

---

## 3. Clean Architecture layering (memorize this)

```
        delivery/http  delivery/middleware      ← parse HTTP, call usecases
                     │
                     ▼
              usecase/<domain>                   ← business logic, stateless
                     │
                     ▼
   domain/entity   domain/repository (interfaces) ← pure Go, zero deps
                     ▲
                     │ (implements the interfaces)
   repository/postgres   repository/redis        ← concrete data access
                     │
                     ▼
            infrastructure/*                      ← db, cache, s3, fcm clients
```

**Rule of thumb:** an inner layer must never `import` an outer layer. The
compiler is your friend — if `domain` ever imports `gin` or `pgx`, you broke it.

---

## 4. Directory layout (target for end of Phase 1)

```
memere-backend/
├── cmd/api/main.go            # entry point, dependency wiring
├── cmd/migrate/main.go        # migration runner (optional wrapper)
├── internal/
│   ├── domain/
│   │   ├── entity/            # User, Course, Section, Lesson … (pure structs)
│   │   └── repository/        # UserRepository, CourseRepository … (interfaces)
│   ├── usecase/
│   │   ├── auth/              # Register, Login, RefreshToken, Logout
│   │   └── course/           # CreateCourse, ListCourses, GetCourse, …
│   ├── repository/
│   │   ├── postgres/          # sqlc-generated + thin adapters to interfaces
│   │   └── redis/             # session / refresh-token store
│   ├── delivery/
│   │   ├── http/              # gin handlers: auth_handler.go, course_handler.go
│   │   └── middleware/        # auth, cors, logging, request-id, rate-limit, recover
│   └── infrastructure/
│       ├── database/          # pgxpool connect, migrate runner
│       └── cache/             # redis client wrapper
├── pkg/
│   ├── jwt/                   # token create/validate
│   ├── password/             # bcrypt hash/compare
│   ├── validator/            # request validation helpers
│   ├── pagination/           # cursor encode/decode
│   └── apperror/             # typed errors + HTTP mapping  (spec calls it errors)
├── config/                   # Config struct + env loading
├── migrations/               # NNNN_name.up.sql / .down.sql  (golang-migrate)
├── db/query/                 # *.sql files consumed by sqlc
├── sqlc.yaml                 # sqlc config
├── .env.example
├── Dockerfile
├── docker-compose.yml
├── Makefile
└── go.mod
```

---

## 5. Phase map (backend only)

We build **phase by phase**. **All six phases are now authored.** Build them
strictly in order — each phase is gated on the previous phase's final
Definition of Done. Each phase has a kickoff prompt in `docs/prompts/`.

| Phase | Theme | Status |
|---|---|---|
| **Phase 1** | **Foundation + Auth + Course CRUD + HTTP wiring** | ✅ Authored |
| **Phase 2** | **Quiz & Exam engines (server-side grading, timers, Redis attempt state, analytics)** | ✅ Authored |
| **Phase 3** | **Video pipeline (pre-signed S3 upload, FFmpeg→HLS transcode, signed streaming/download)** | ✅ Authored |
| **Phase 4** | **Payments (Chapa/Telebirr/Stripe, idempotency, webhooks) + enrollments + access control** | ✅ Authored |
| **Phase 5** | **Progress tracking + Notifications (FCM/SendGrid/in-app) + Admin analytics + Certificates** | ✅ Authored |
| **Phase 6** | **Hardening + observability + performance + CI/CD + k8s + microservice plan (final)** | ✅ Authored |

### Phase 1 skill order (run strictly in sequence)

| # | File | Builds |
|---|---|---|
| 1 | [`skills/phase1/skill1.md`](./skills/phase1/skill1.md) | Project scaffold, `go.mod`, config, Docker Compose, pgx pool + Redis connect, Makefile, healthcheck |
| 2 | [`skills/phase1/skill2.md`](./skills/phase1/skill2.md) | Domain entities, repository interfaces, SQL migrations, sqlc setup |
| 3 | [`skills/phase1/skill3.md`](./skills/phase1/skill3.md) | Auth: bcrypt, JWT, refresh tokens in Redis+PG, register/login/refresh/logout usecases |
| 4 | [`skills/phase1/skill4.md`](./skills/phase1/skill4.md) | Course service: CRUD + sections + lessons usecases & repositories |
| 5 | [`skills/phase1/skill5.md`](./skills/phase1/skill5.md) | HTTP layer: Gin handlers, middleware stack, router, dependency wiring in `main.go` |

Each skill ends with a **Definition of Done** checklist and **verification
commands**. A skill is not complete until every box is checked and the
verification commands pass.

### Phase 2 skill order (run strictly in sequence — only after Phase 1 is green)

| # | File | Builds |
|---|---|---|
| 1 | [`skills/phase2/skill1.md`](./skills/phase2/skill1.md) | Quiz/exam data layer: additive migrations (quiz_attempts, exam_questions, enhancements), entities, repo interfaces, sqlc (client-vs-grading queries) |
| 2 | [`skills/phase2/skill2.md`](./skills/phase2/skill2.md) | Quiz engine: attempts, Redis snapshot/randomization, auto-save, server-side grading, attempt limits, feedback |
| 3 | [`skills/phase2/skill3.md`](./skills/phase2/skill3.md) | Exam engine: §9.2 state machine, server-enforced timer, exam assembly, shared grading core |
| 4 | [`skills/phase2/skill4.md`](./skills/phase2/skill4.md) | Background expiry sweeper (auto-grade abandoned attempts) + §9.3 scoring/analytics (percentile, weak areas, trend) |
| 5 | [`skills/phase2/skill5.md`](./skills/phase2/skill5.md) | HTTP delivery for quiz/exam/analytics (answer-key-free DTOs), routes, RBAC, sweeper wiring, Phase 2 smoke test |

### Phase 3 skill order (run strictly in sequence — only after Phase 2 is green)

| # | File | Builds |
|---|---|---|
| 1 | [`skills/phase3/skill1.md`](./skills/phase3/skill1.md) | Video data layer: `videos` Go layer, additive migration, `ObjectStore` port + S3/MinIO impl, key scheme, local MinIO |
| 2 | [`skills/phase3/skill2.md`](./skills/phase3/skill2.md) | Upload flow: pre-signed PUT usecase, `JobQueue` port + in-proc impl, confirm-upload enqueues transcode, boot reconciler |
| 3 | [`skills/phase3/skill3.md`](./skills/phase3/skill3.md) | Transcode worker: FFmpeg→HLS ladder (480/720/1080) + thumbnail, guarded state machine, bounded retry |
| 4 | [`skills/phase3/skill4.md`](./skills/phase3/skill4.md) | Secure delivery: CloudFront/presign signer, 2h stream URLs, single-use download tokens, access control |
| 5 | [`skills/phase3/skill5.md`](./skills/phase3/skill5.md) | HTTP delivery for video API, worker + queue + signer wiring, Phase 3 smoke test (MinIO + FFmpeg) |

### Phase 4 skill order (run strictly in sequence — only after Phase 3 is green)

| # | File | Builds |
|---|---|---|
| 1 | [`skills/phase4/skill1.md`](./skills/phase4/skill1.md) | Payment/enrollment/coupon data layer, `PaymentProvider` port, idempotency + webhook-dedup tables |
| 2 | [`skills/phase4/skill2.md`](./skills/phase4/skill2.md) | Enrollment engine + unified `access.Service`; replace every `TODO(phase4)` hook from Phases 2–3 |
| 3 | [`skills/phase4/skill3.md`](./skills/phase4/skill3.md) | Payment flow: idempotent initiate, provider checkout, signature-verified webhooks, transactional fulfillment |
| 4 | [`skills/phase4/skill4.md`](./skills/phase4/skill4.md) | Subscriptions (plans, renewal/expiry sweeper) + revenue reporting (70/30 split) |
| 5 | [`skills/phase4/skill5.md`](./skills/phase4/skill5.md) | HTTP delivery (raw-body webhook route), provider + sweeper wiring, Phase 4 smoke test |

### Phase 5 skill order (run strictly in sequence — only after Phase 4 is green)

| # | File | Builds |
|---|---|---|
| 1 | [`skills/phase5/skill1.md`](./skills/phase5/skill1.md) | Progress tracking: completion %, video position, study streaks, dashboard |
| 2 | [`skills/phase5/skill2.md`](./skills/phase5/skill2.md) | Notification system: `Notifier` port, FCM/SendGrid/in-app, async dispatch; wire all no-op hooks |
| 3 | [`skills/phase5/skill3.md`](./skills/phase5/skill3.md) | Admin operations + platform analytics + audit log + broadcast announcements |
| 4 | [`skills/phase5/skill4.md`](./skills/phase5/skill4.md) | Certificates (PDF + signed download + public verify) + engagement/streak-warning sweeper |
| 5 | [`skills/phase5/skill5.md`](./skills/phase5/skill5.md) | HTTP delivery (progress/notifications/admin/certs), worker wiring, Phase 5 smoke test |

### Phase 6 skill order (run strictly in sequence — only after Phase 5 is green) — FINAL

| # | File | Builds |
|---|---|---|
| 1 | [`skills/phase6/skill1.md`](./skills/phase6/skill1.md) | Observability: slog + redaction, Prometheus RED+business metrics, OTel tracing, Sentry, healthz/readyz |
| 2 | [`skills/phase6/skill2.md`](./skills/phase6/skill2.md) | Security hardening: §7.3 audit, security headers, segmented rate limits, lockout, token revocation, CI scans |
| 3 | [`skills/phase6/skill3.md`](./skills/phase6/skill3.md) | Performance: cache-aside, query/index tuning, leaderboard, pool sizing, load test (p95<200ms) |
| 4 | [`skills/phase6/skill4.md`](./skills/phase6/skill4.md) | CI/CD: hardened multi-stage images (api + worker), GitHub Actions per §13.1, migrations in pipeline |
| 5 | [`skills/phase6/skill5.md`](./skills/phase6/skill5.md) | Kubernetes manifests + HPA + ingress/TLS + alerts; scaling playbook + microservices extraction plan |

Each skill ends with a **Definition of Done** checklist and **verification
commands**. A skill is not complete until every box is checked and the
verification commands pass.

---

## 6. How to work a skill (the loop Antigravity should follow)

For each skill file, in order:

1. **Read** the whole skill file before writing any code.
2. **Re-read** the relevant section of `memere_Design_Specification.md` it points to.
3. **Plan** the files you'll create/modify (list them).
4. **Implement** following the file's "Tasks" section exactly.
5. **Self-check** against the "Non-Negotiables" (Section 2 above).
6. **Verify** by running the skill's verification commands.
7. **Tick** the Definition of Done. Only then move to the next skill.

Do not jump ahead. Skill N assumes Skill N-1's Definition of Done is satisfied.

---

## 7. Conventions used in every skill file

- **API:** base path `/api/v1`, plural kebab-case resources, cursor pagination
  (`limit` + `after`), `Authorization: Bearer <jwt>`.
- **Error envelope:** `{ "code": "SNAKE_CASE_CODE", "message": "...", "details": {} }`.
- **Timestamps:** `created_at` / `updated_at` on every table; `deleted_at` where
  soft-deletable. All `TIMESTAMPTZ`.
- **Commits:** Conventional Commits (`feat:`, `fix:`, `docs:`, `refactor:`,
  `test:`, `chore:`). One commit per logical unit, ideally one per skill task.
- **Tests:** table-driven, `testify` for assertions. Usecases get unit tests with
  mocked repositories; repositories get integration tests against a real
  Postgres (Docker).

---

## 8. How to start (the very first prompt to Antigravity)

The exact wording is in [`docs/prompts/phase1-kickoff.md`](./prompts/phase1-kickoff.md).
In short: point Antigravity at this `skill.md`, tell it to read the design spec,
then execute `phase1/skill1.md` through `skill5.md` in order, stopping after each
skill's Definition of Done for your review.

---

## 9. Moving between phases

**All six phases are authored.** Each phase has a kickoff prompt in
`docs/prompts/phaseN-kickoff.md`. Build them in order:

1. Start a phase with its kickoff prompt (e.g. `docs/prompts/phase1-kickoff.md`).
2. Run its skills `skill1 → skill5` strictly in sequence, stopping after each
   skill's Definition of Done for review.
3. When a phase's final skill (skill5) passes its Definition of Done, report
   "Phase N complete" and move to the next phase's kickoff prompt.

Phase dependency notes baked into the skills:
- Phases 2–3 leave `TODO(phase4)` access hooks; **Phase 4 must remove them all**
  (grep gate).
- Phases 3–4 leave `notifyNoop` hooks; **Phase 5 must wire them all** (grep gate).
- Phase 6 changes **no business behavior** — only observability, security,
  performance, and deployment.

**Never** scaffold a later phase's code while building an earlier one. Phases are
gated on review. The microservice split (Phase 6 Skill 5) is a **plan** — do not
execute it until the §12.2 scaling triggers fire.

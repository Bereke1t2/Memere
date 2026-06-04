# Phase 6 · Skill 4 — CI/CD Pipeline & Containerization

> **Prerequisite:** Phase 6 Skills 1–3 done (observable, hardened, performant).
> Read [`docs/skill.md`](../../skill.md) §2.
>
> **Spec references:** `memere_Design_Specification.md` §13.1 (CI/CD pipeline
> diagram + steps), §13 (DevOps), README "CI/CD Pipeline" & "Deployment".

---

## Goal

Automate build → test → scan → image → deploy. Produce a small, secure container
image; a GitHub Actions pipeline implementing the §13.1 branch strategy (feature →
CI; main → CI + staging; tag → CI + production); and a repeatable migration step.
After this skill, every merge is tested and shippable.

---

## §13.1 branch → action map (implement exactly)

| Trigger | Actions |
|---|---|
| feature branch / PR | lint, `go test ./...`, `govulncheck`/`gosec`/`gitleaks`, `docker build` (no push) |
| push to `main` | all CI + build & push image + deploy to **staging** + integration smoke + notify |
| release tag `v*` | all CI + build & push prod-tagged image + deploy to **production** + smoke + notify |

---

## Tasks

### 4.1 — Harden the Dockerfile (multi-stage, minimal, non-root)

```dockerfile
# build stage
FROM golang:1.22 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w -X main.version=${VERSION}" \
    -o /out/api ./cmd/api
# (separate binary for the migration runner if you have cmd/migrate)
RUN CGO_ENABLED=0 go build -trimpath -o /out/migrate ./cmd/migrate

# runtime stage
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/api /api
COPY --from=build /out/migrate /migrate
COPY --from=build /src/migrations /migrations
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/api"]
```

- FFmpeg note: the **transcode worker** needs ffmpeg, which distroless lacks. Two
  options: (a) a **separate worker image** based on a small image with ffmpeg
  installed (recommended — foreshadows Skill 5's worker deployment), or (b) keep
  transcoding out of the API image entirely. **Decision:** build two images —
  `api` (distroless) and `worker` (alpine + ffmpeg). Document it.
- Pin base image digests; scan the image (Trivy) in CI.

### 4.2 — Compose for full local stack

Extend `docker-compose.yml` so `docker compose up` runs the **whole** app
(api + worker + postgres + redis + minio) locally — useful for integration tests.
Keep the dev `make run` path too.

### 4.3 — Migrations in the pipeline

- A CI/CD step runs `/migrate up` against the target DB **before** the new app
  rolls out (init container in k8s — Skill 5; a pipeline step for now).
- Migrations must be **backward-compatible** for zero-downtime rollouts
  (expand/contract pattern: add columns nullable, deploy, backfill, then tighten
  in a later migration). Document the rule; the additive migrations from Phases
  2–5 already follow it.

### 4.4 — GitHub Actions workflows (`.github/workflows/`)

`ci.yml` (PR + push):

```yaml
name: ci
on: { pull_request: {}, push: { branches: [main] } }
jobs:
  test:
    runs-on: ubuntu-latest
    services:
      postgres: { image: postgres:15, env: { POSTGRES_PASSWORD: test }, ports: ["5432:5432"],
        options: >-
          --health-cmd "pg_isready" --health-interval 10s --health-timeout 5s --health-retries 5 }
      redis: { image: redis:7, ports: ["6379:6379"] }
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with: { go-version: '1.22' }
      - run: go mod download
      - run: make migrate-up           # against the service DB
      - run: make sqlc && git diff --exit-code   # generated code is committed & current
      - run: golangci-lint run
      - run: go test ./... -race -coverprofile=cover.out
      - run: govulncheck ./...
      - uses: gitleaks/gitleaks-action@v2
  build:
    needs: test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - run: docker build -t memere-api:${{ github.sha }} .
      # on main/tags: login to ECR/GHCR + push + Trivy scan
```

`deploy-staging.yml` (on `main`) and `deploy-prod.yml` (on `v*` tags): build+push
the `api` and `worker` images, run migrations, `kubectl apply`/`rollout`
(Skill 5 provides manifests), smoke test, notify Slack. Use environment
protection rules + secrets for registry/kubeconfig.

### 4.5 — Makefile targets for parity

`make ci` runs the same lint+test+scan locally so devs reproduce CI. `make
docker-build` / `make docker-build-worker`. `make compose-up` for the full stack.

### 4.6 — Versioning & release

- Inject version/commit via ldflags into a `/version` endpoint (or readyz body).
- Tag releases `vX.Y.Z` (semver); conventional-commit history drives the changelog
  (optional: `git-cliff`).

### 4.7 — Tests / validation

CI runs green on a sample PR; `sqlc`-generated code is verified current (diff
gate); image builds and runs (`docker run` → `/healthz` 200); Trivy finds no
high/critical; migrations apply cleanly in the CI service DB.

---

## Definition of Done

- [ ] `docker build` produces a small (<30 MB ideal) non-root distroless **api**
      image; a separate **worker** image includes ffmpeg.
- [ ] `docker run` of the api image serves `/healthz` 200.
- [ ] `ci.yml` runs lint + `-race` tests (with PG/Redis services) + `govulncheck` +
      `gitleaks` + sqlc-drift gate, and is green on a sample PR.
- [ ] `main` → staging deploy and `v*` tag → prod deploy workflows exist with
      migration step, smoke test, and notification (per §13.1).
- [ ] Migrations run as a pre-rollout step and follow the expand/contract
      (zero-downtime) rule.
- [ ] Image scanned (Trivy); no high/critical vulns; base images digest-pinned.
- [ ] `/version` (or readyz) reports build version/commit.

## Verification commands

```bash
docker build -t memere-api:test .
docker run --rm -p 8080:8080 --env-file .env memere-api:test &
curl -i localhost:8080/healthz
trivy image memere-api:test --severity HIGH,CRITICAL
make ci   # lint + test + scan locally
```

## Hand-off to Skill 5

The app ships as scanned images via automated pipelines. Skill 5 — the **final
skill** — delivers the Kubernetes manifests, autoscaling, and the
**modular-monolith → microservices** extraction plan (the §3.1/§12 endgame).

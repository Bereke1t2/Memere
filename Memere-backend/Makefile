.PHONY: run build tidy lint test up down migrate-up migrate-down migrate-create sqlc minio-bucket \
        ci docker-build docker-build-worker compose-up compose-down seed seed-admin seed-teacher \
        k8s-dry-run k8s-apply k8s-migrate k8s-status k8s-rollout kind-up kind-down

# Auto-load .env so config.Load() (which reads the process environment) sees
# DB_HOST, JWT_SECRET, etc. without the caller exporting them by hand. The
# variables are exported into the environment of every recipe below.
ifneq (,$(wildcard ./.env))
include .env
export
endif

# Build metadata injected into the binary via ldflags (matches Dockerfile).
VERSION  ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT   ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
BUILD_TIME ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

LDFLAGS := -s -w \
  -X 'github.com/Bereke1t2/Memere/memere-backend/internal/version.Version=$(VERSION)' \
  -X 'github.com/Bereke1t2/Memere/memere-backend/internal/version.Commit=$(COMMIT)' \
  -X 'github.com/Bereke1t2/Memere/memere-backend/internal/version.BuildTime=$(BUILD_TIME)'

run:
	go run ./cmd/api

build:
	CGO_ENABLED=0 go build -trimpath -ldflags="$(LDFLAGS)" -o bin/api ./cmd/api

tidy:
	go mod tidy

lint:
	golangci-lint run

test:
	go test ./...

# ci — reproduces the full CI pipeline locally (lint + race tests + security scans).
# Requires: golangci-lint, govulncheck, gosec in PATH.
ci: lint
	go test ./... -race -timeout 120s
	govulncheck ./... || true
	gosec -severity high -confidence high ./... || true
	@echo "CI checks complete."

up:
	docker-compose up -d postgres redis minio

down:
	docker-compose down

# compose-up — start the entire application stack (infra + api + worker).
compose-up:
	docker-compose up -d

# compose-down — stop the full stack and remove containers (data volumes preserved).
compose-down:
	docker-compose down

# docker-build — build the API image with version metadata.
docker-build:
	docker build \
	  --build-arg VERSION=$(VERSION) \
	  --build-arg COMMIT=$(COMMIT) \
	  --build-arg BUILD_TIME=$(BUILD_TIME) \
	  -t memere-api:$(VERSION) \
	  -f Dockerfile .

# docker-build-worker — build the transcode worker image (includes FFmpeg).
docker-build-worker:
	docker build \
	  --build-arg VERSION=$(VERSION) \
	  --build-arg COMMIT=$(COMMIT) \
	  --build-arg BUILD_TIME=$(BUILD_TIME) \
	  -t memere-worker:$(VERSION) \
	  -f Dockerfile.worker .

# migrate-up / migrate-down apply the embedded migrations via the library-based
# runner in cmd/migrate (no external `migrate` binary required). DSN is read
# from the environment / .env by config.Load().
migrate-up:
	go run ./cmd/migrate -direction up

migrate-down:
	go run ./cmd/migrate -direction down

# migrate-create scaffolds a new up/down pair:  make migrate-create name=add_widgets
migrate-create:
	@if [ -z "$(name)" ]; then echo "usage: make migrate-create name=<description>"; exit 1; fi
	@last=$$(ls migrations/[0-9]*.up.sql 2>/dev/null | sed -E 's#migrations/([0-9]+)_.*#\1#' | sort -n | tail -1); \
	next=$$(printf '%04d' $$((10#$${last:-0} + 1))); \
	up="migrations/$${next}_$(name).up.sql"; \
	down="migrations/$${next}_$(name).down.sql"; \
	echo "-- $(name) (up)" > $$up; \
	echo "-- $(name) (down)" > $$down; \
	echo "created $$up"; \
	echo "created $$down"

# sqlc generate runs the pinned sqlc via the go tool directive (go.mod), so no
# global install is required.
sqlc:
	go tool sqlc generate

# seed — populate the database with realistic development data.
seed:
	go run ./scripts/seed

# seed-admin — create (or reset) the admin account used for local dev and smoke
# tests. Idempotent: safe to run multiple times. Override defaults with env vars:
#   ADMIN_EMAIL=me@example.com ADMIN_PASSWORD=secret make seed-admin
seed-admin:
	go run ./cmd/seed-admin

# seed-teacher — create (or reset) a teacher account for local dev and smoke tests.
# Idempotent. Override defaults:
#   TEACHER_EMAIL=me@example.com TEACHER_PASSWORD=secret make seed-teacher
seed-teacher:
	go run ./cmd/seed-teacher

# ---------------------------------------------------------------------------
# Kubernetes targets (Phase 6 Skill 5)
# ---------------------------------------------------------------------------
NAMESPACE ?= production

# k8s-dry-run — validate manifests against a running cluster without applying.
# Requires kubectl and a configured kubeconfig (or KUBECONFIG env var).
k8s-dry-run:
	kubectl apply -k k8s/ --dry-run=server --namespace=$(NAMESPACE)

# k8s-apply — apply all manifests. Run k8s-migrate first on a fresh deploy.
k8s-apply:
	kubectl apply -k k8s/ --namespace=$(NAMESPACE)

# k8s-migrate — run the migration Job and wait for it to complete.
# Always run before deploying new image tags that include schema changes.
k8s-migrate:
	kubectl apply -f k8s/jobs/migrate-job.yaml --namespace=$(NAMESPACE)
	kubectl wait --for=condition=complete job/migrate --namespace=$(NAMESPACE) --timeout=120s

# k8s-status — print pod and HPA status for the namespace.
k8s-status:
	@echo "=== Pods ==="
	kubectl get pods --namespace=$(NAMESPACE)
	@echo ""
	@echo "=== HPA ==="
	kubectl get hpa --namespace=$(NAMESPACE)
	@echo ""
	@echo "=== Ingress ==="
	kubectl get ingress --namespace=$(NAMESPACE)

# k8s-rollout — watch rollout status for api and worker deployments.
k8s-rollout:
	kubectl rollout status deployment/api --namespace=$(NAMESPACE)
	kubectl rollout status deployment/worker --namespace=$(NAMESPACE)

# kind-up — create a local kind cluster and apply all manifests for testing.
# Requires: kind (https://kind.sigs.k8s.io) and kubectl.
kind-up:
	kind create cluster --name memere-local
	kubectl apply -k k8s/ --dry-run=client   # validate before applying
	@echo "Cluster ready. Apply manifests with: kubectl apply -k k8s/"

# kind-down — delete the local kind cluster.
kind-down:
	kind delete cluster --name memere-local

# ---------------------------------------------------------------------------
# minio-bucket creates the dev media bucket in the running MinIO (make up first)
# and ensures it is private (no anonymous access — Non-Negotiable #3).
minio-bucket:
	docker run --rm --network host --entrypoint sh minio/mc -c \
	 'mc alias set local http://localhost:9000 $(AWS_ACCESS_KEY_ID) $(AWS_SECRET_ACCESS_KEY) && \
	  mc mb -p local/$(AWS_S3_BUCKET) && mc anonymous set none local/$(AWS_S3_BUCKET)'

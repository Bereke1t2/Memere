.PHONY: run build tidy lint test up down migrate-up migrate-down migrate-create sqlc minio-bucket \
        ci docker-build docker-build-worker compose-up compose-down seed

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

# minio-bucket creates the dev media bucket in the running MinIO (make up first)
# and ensures it is private (no anonymous access — Non-Negotiable #3).
minio-bucket:
	docker run --rm --network host --entrypoint sh minio/mc -c \
	 'mc alias set local http://localhost:9000 $(AWS_ACCESS_KEY_ID) $(AWS_SECRET_ACCESS_KEY) && \
	  mc mb -p local/$(AWS_S3_BUCKET) && mc anonymous set none local/$(AWS_S3_BUCKET)'

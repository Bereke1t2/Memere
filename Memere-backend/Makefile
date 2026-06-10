.PHONY: run build tidy lint test up down migrate-up migrate-down migrate-create sqlc minio-bucket

# Auto-load .env so config.Load() (which reads the process environment) sees
# DB_HOST, JWT_SECRET, etc. without the caller exporting them by hand. The
# variables are exported into the environment of every recipe below.
ifneq (,$(wildcard ./.env))
include .env
export
endif

run:
	go run ./cmd/api

build:
	go build -o bin/api ./cmd/api

tidy:
	go mod tidy

lint:
	golangci-lint run

test:
	go test ./...

up:
	docker-compose up -d

down:
	docker-compose down

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

# minio-bucket creates the dev media bucket in the running MinIO (make up first)
# and ensures it is private (no anonymous access — Non-Negotiable #3).
minio-bucket:
	docker run --rm --network host --entrypoint sh minio/mc -c \
	 'mc alias set local http://localhost:9000 $(AWS_ACCESS_KEY_ID) $(AWS_SECRET_ACCESS_KEY) && \
	  mc mb -p local/$(AWS_S3_BUCKET) && mc anonymous set none local/$(AWS_S3_BUCKET)'

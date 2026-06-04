.PHONY: run build tidy lint test up down migrate-up migrate-down sqlc

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

migrate-up:
	@echo "Stub for Skill 2"

migrate-down:
	@echo "Stub for Skill 2"

sqlc:
	@echo "Stub for Skill 2"

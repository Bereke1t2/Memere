# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Memere (ExamPrep) is a mobile-first learning platform for Grade 12 students in Ethiopia preparing for university entrance exams. This repository is the **Go backend** — still in early scaffolding phase (all `.go` files are empty stubs). The Flutter mobile app lives in a sibling `memere_mobile` directory.

**Tech stack:** Go 1.22+, PostgreSQL 15, Redis 7, Docker/Kubernetes, JWT auth, HLS video streaming, multi-provider payments (Chapa/Telebirr/Stripe).

## Architecture

The backend follows **Uncle Bob's Clean Architecture** with strict dependency inversion. Outer layers depend on inner layers — never the reverse:

```
HTTP Handlers (delivery/http, delivery/middleware)
  → Use Cases (usecase/) — stateless business logic orchestrators
    → Domain (domain/entity, domain/repository) — pure Go, zero external deps
      ← Repository Implementations (repository/postgres, repository/redis)
        ← Infrastructure (infrastructure/) — DB, S3, FCM, payment clients
```

**Microservices** (one per bounded domain, each on its own port):
| Service | Port | Purpose |
|---|---|---|
| API Gateway | 8080 | Rate limiting, routing, auth validation, SSL termination |
| Auth Service | 8081 | Registration, login, JWT, refresh tokens, password reset |
| Course Service | 8082 | Course CRUD, section/lesson management, video/note metadata |
| Quiz Service | 8083 | Quiz creation, question bank, attempt recording, auto-grading |
| Exam Service | 8084 | Mock exam engine, timer management, exam sessions, scoring |
| Payment Service | 8085 | Purchase flow, subscription management, webhook handling |
| Notification Service | 8086 | Push (FCM), email (SendGrid), in-app notifications |
| Progress Service | 8087 | Completion tracking, streak calculation, analytics |

## Commands

```bash
# Infrastructure (PostgreSQL + Redis)
docker-compose up -d

# Database migrations
make migrate-up

# Seed development data
make seed

# Run the server
make run

# Run tests
go test ./...

# Lint
golangci-lint run

# Build Docker image
docker build -t memere-backend .
```

## Key Design Rules (Non-Negotiable)

1. **Correct answers NEVER sent to client** — always grade server-side
2. **Exam timer enforced server-side** — client timer is display-only
3. **Pre-signed CDN URLs for video** — no public S3 access; 2-hour expiry
4. **Idempotency keys on all payment requests** — prevent double-charge
5. **Soft deletes only** — `deleted_at` timestamp, never hard DELETE user data
6. **HTTPS only** — HTTP redirects to HTTPS at gateway
7. **All DB queries must filter by authenticated `user_id`** — prevent IDOR
8. **Never log raw passwords, tokens, or payment card data**

## API Conventions

- Base path: `/api/v1`
- Resource naming: plural nouns, kebab-case
- Pagination: cursor-based with `limit`/`after` params
- Auth: `Authorization: Bearer <jwt>` (access token 15min TTL, refresh token 30 days)
- Error format: `{ "code": "RESOURCE_NOT_FOUND", "message": "...", "details": {} }`
- UUID primary keys on all entities
- `created_at`/`updated_at` audit columns on every table

## Database Design

- One PostgreSQL cluster, schemas per domain (`auth`, `courses`, `payments`, `progress`)
- JSONB for flexible metadata; typed columns for queryable fields
- RBAC roles: `student`, `teacher`, `admin`

## Commit Convention

[Conventional Commits](https://www.conventionalcommits.org/): `feat:`, `fix:`, `docs:`, `refactor:`, `test:`, `chore:`

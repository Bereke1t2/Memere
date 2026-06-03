<p align="center">
  <h1 align="center">📚 Memere — ExamPrep Backend</h1>
  <p align="center">
    A high-performance Go backend for Ethiopia's Grade 12 university entrance exam preparation platform.
  </p>
  <p align="center">
    <img src="https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go">
    <img src="https://img.shields.io/badge/PostgreSQL-15-4169E1?style=for-the-badge&logo=postgresql&logoColor=white" alt="PostgreSQL">
    <img src="https://img.shields.io/badge/Redis-7-DC382D?style=for-the-badge&logo=redis&logoColor=white" alt="Redis">
    <img src="https://img.shields.io/badge/Docker-Ready-2496ED?style=for-the-badge&logo=docker&logoColor=white" alt="Docker">
    <img src="https://img.shields.io/badge/License-MIT-green?style=for-the-badge" alt="License">
  </p>
</p>

---

## 📖 Overview

**Memere (ExamPrep)** is a mobile-first learning platform built specifically for **Grade 12 students in Ethiopia** preparing for university entrance examinations. The platform delivers structured video lessons, study notes, adaptive quizzes, and timed mock examinations — specialized for high-stakes national exam preparation.

This repository contains the **Go backend** powering the platform, built with **Clean Architecture** principles (Uncle Bob) and designed to scale from 1,000 to 1,000,000+ users.

### 🎯 Business Goals

| Goal | Target |
|------|--------|
| Market Entry | 500 registered users in Month 1 |
| Revenue | ETB 50,000 MRR by Month 6 |
| Student Outcomes | 15% average score improvement |
| Content Scale | 50+ courses covering all subjects by Month 6 |
| Retention | DAU/MAU ratio > 40% |

---

## ✨ Key Features

### Student-Facing

- 🎥 **Video Lessons** — Stream and download HD video lessons per subject (HLS adaptive bitrate)
- 📝 **Study Notes** — Access PDF/document notes per lesson
- ❓ **Quizzes** — Per-lesson quizzes with immediate feedback
- 📋 **Mock Exams** — Full timed mock exams simulating national exam format
- 📊 **Progress Tracking** — Dashboard showing completion %, scores, weak areas
- 💳 **Course Purchase** — Buy individual courses or bundles via Chapa/Telebirr/Stripe
- 📱 **Offline Mode** — Download videos and notes for offline use
- 🔔 **Push Notifications** — Reminders, new content alerts, score notifications
- 🏆 **Leaderboard** — Class ranking based on quiz/exam scores
- 📜 **Certificates** — Downloadable completion certificates

### Teacher & Admin Features

- 📤 **Video Upload** — Upload and manage video lessons with automatic HLS transcoding
- 📄 **Notes Upload** — Upload PDFs and documents
- 🧩 **Quiz Builder** — Create multiple-choice and short-answer quizzes
- 🧪 **Exam Builder** — Build timed mock exams with question randomization
- 📚 **Course Manager** — Create, edit, and publish courses
- 📈 **Student Analytics** — View per-student performance data
- 💰 **Payment Dashboard** — Track revenue, payouts, and subscriptions
- 🎟️ **Coupon Manager** — Create discount coupons and promotional codes

---

## 🛠️ Tech Stack

| Component | Technology | Purpose |
|-----------|-----------|---------|
| **Language** | Go 1.22+ | High performance, low memory, fast compile |
| **HTTP Framework** | Gin / Echo | API routing and middleware |
| **Database** | PostgreSQL 15 | Primary relational data store (ACID, JSONB) |
| **Cache** | Redis 7 | Sessions, exam state, leaderboard sorted sets |
| **Object Storage** | AWS S3 / GCS | Video files, PDFs, images, certificates |
| **CDN** | CloudFront / GCP CDN | Low-latency video and static content delivery |
| **Video Processing** | FFmpeg / AWS MediaConvert | HLS adaptive bitrate transcoding |
| **Message Queue** | SQS / RabbitMQ | Async video processing, emails, notifications |
| **Auth** | JWT + bcrypt | Access tokens (15min) + refresh tokens (30 days) |
| **Payments** | Chapa + Telebirr + Stripe | ETB mobile money + international cards |
| **Push Notifications** | Firebase Cloud Messaging | Mobile push notifications |
| **Email** | SendGrid | Transactional emails |
| **Containers** | Docker + Kubernetes | Deployment and orchestration |
| **CI/CD** | GitHub Actions | Automated build, test, and deploy |
| **Monitoring** | Prometheus + Grafana + Sentry | Metrics, dashboards, and error tracking |

---

## 🏗️ Architecture

The backend follows **Uncle Bob's Clean Architecture** with strict dependency inversion. Outer layers depend on inner layers — never the reverse. The Domain layer has zero external dependencies.

```
┌─────────────────────────────────────────────────┐
│                 HTTP Handlers                    │  ← Outermost (delivery)
│              (Gin/Echo routes)                   │
├─────────────────────────────────────────────────┤
│               Use Cases                          │  ← Business Logic
│         (Stateless orchestrators)                │
├─────────────────────────────────────────────────┤
│            Domain Entities                       │  ← Innermost (pure Go)
│     (User, Course, Quiz, Exam structs)           │
├─────────────────────────────────────────────────┤
│          Repository Interfaces                   │  ← Contracts
│       (Go interfaces for DB ops)                 │
├─────────────────────────────────────────────────┤
│        Infrastructure Layer                      │  ← External concerns
│   (PostgreSQL, Redis, S3, FCM clients)           │
└─────────────────────────────────────────────────┘
```

### Microservices

| Service | Port | Responsibilities |
|---------|------|-----------------|
| API Gateway | `8080` | Rate limiting, routing, auth validation, SSL termination |
| Auth Service | `8081` | Registration, login, JWT issuance, refresh tokens, password reset |
| Course Service | `8082` | Course CRUD, section/lesson management, video/note metadata |
| Quiz Service | `8083` | Quiz creation, question bank, attempt recording, auto-grading |
| Exam Service | `8084` | Mock exam engine, timer management, exam sessions, scoring |
| Payment Service | `8085` | Purchase flow, subscription management, webhook handling |
| Notification Service | `8086` | Push (FCM), email (SendGrid), in-app notifications |
| Progress Service | `8087` | Completion tracking, streak calculation, analytics |

---

## 📁 Project Structure

```
memere-backend/
├── cmd/
│   ├── api/                # main.go — entry point, dependency wiring
│   └── migrate/            # Database migration runner
├── internal/
│   ├── domain/             # ← INNERMOST LAYER (no external imports)
│   │   ├── entity/         # User, Course, Quiz, Exam structs
│   │   ├── repository/     # Repository INTERFACES (contracts)
│   │   └── service/        # Domain service interfaces
│   ├── usecase/            # ← BUSINESS LOGIC
│   │   ├── auth/           # RegisterUser, LoginUser, RefreshToken
│   │   ├── course/         # CreateCourse, EnrollStudent, GetCourseById
│   │   ├── quiz/           # SubmitQuizAttempt, GradeQuiz
│   │   ├── exam/           # StartExam, SubmitExam, GetExamResults
│   │   ├── payment/        # InitiatePayment, HandleWebhook
│   │   └── progress/       # MarkLessonComplete, GetCourseProgress
│   ├── repository/         # ← REPOSITORY IMPLEMENTATIONS
│   │   ├── postgres/       # SQL queries per entity
│   │   └── redis/          # Cache operations
│   ├── delivery/           # ← HTTP HANDLERS
│   │   ├── http/           # Gin/Echo route handlers
│   │   └── middleware/     # Auth, CORS, logging, rate-limit
│   └── infrastructure/     # ← EXTERNAL CONCERNS
│       ├── database/       # PostgreSQL connection, migrations
│       ├── cache/          # Redis client wrapper
│       ├── storage/        # S3/GCS client wrapper
│       ├── messaging/      # SQS/RabbitMQ client
│       ├── payment/        # Stripe, Chapa, Telebirr clients
│       └── notification/   # FCM, SendGrid wrappers
├── pkg/
│   ├── jwt/                # JWT creation and validation helpers
│   ├── password/           # bcrypt helpers
│   ├── validator/          # Input validation utilities
│   ├── pagination/         # Cursor-based pagination helpers
│   └── errors/             # Custom error types and codes
├── config/                 # Config structs, env loading
├── migrations/             # SQL migration files (golang-migrate)
├── api/                    # OpenAPI/Swagger spec
├── scripts/
│   ├── seed.go             # Test data seeder
│   └── generate.go         # Code generation scripts
├── docs/                   # Design specification & documentation
├── Dockerfile
├── docker-compose.yml
├── Makefile
└── README.md
```

---

## 📋 Prerequisites

Before you begin, ensure you have the following installed:

- **Go** 1.22 or higher — [Install Go](https://go.dev/doc/install)
- **PostgreSQL** 15+ — [Install PostgreSQL](https://www.postgresql.org/download/)
- **Redis** 7+ — [Install Redis](https://redis.io/docs/getting-started/)
- **Docker** & **Docker Compose** — [Install Docker](https://docs.docker.com/get-docker/)
- **golang-migrate** — [Install](https://github.com/golang-migrate/migrate)

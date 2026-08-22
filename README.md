# 📚 Memere (ExamPrep) — Platform Ecosystem

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go">
  <img src="https://img.shields.io/badge/Next.js-15+-000000?style=for-the-badge&logo=nextdotjs&logoColor=white" alt="Next.js">
  <img src="https://img.shields.io/badge/Flutter-3.13+-02569B?style=for-the-badge&logo=flutter&logoColor=white" alt="Flutter">
  <img src="https://img.shields.io/badge/PostgreSQL-15-4169E1?style=for-the-badge&logo=postgresql&logoColor=white" alt="PostgreSQL">
  <img src="https://img.shields.io/badge/Redis-7-DC382D?style=for-the-badge&logo=redis&logoColor=white" alt="Redis">
  <img src="https://img.shields.io/badge/Docker-Ready-2496ED?style=for-the-badge&logo=docker&logoColor=white" alt="Docker">
  <img src="https://img.shields.io/badge/License-MIT-green?style=for-the-badge" alt="License">
</p>

---


## 📖 Overview

**Memere (ExamPrep)** is an end-to-end digital learning and examination preparation ecosystem designed specifically for **Grade 12 students in Ethiopia** preparing for national university entrance examinations.

The platform provides structured subject courses, low-latency adaptive HLS video streaming, downloadable PDF study notes, per-lesson quizzes, full timed mock entrance exams, real-time analytics, and mobile money payment integrations (Chapa and Telebirr).

---

## 📁 Repository Structure

This monorepo repository unites all three core pillars of the Memere platform:

| Component | Path | Tech Stack | Role & Purpose |
|---|---|---|---|
| ⚙️ **Backend** | [`/Memere-backend`](Memere-backend/README.md) | Go 1.22+, Postgres 15, Redis 7, MinIO/S3 | Central REST API engine, async media transcoding worker, auth, payments & core database |
| 💻 **Admin Web** | [`/Memere-admin`](Memere-admin/README.md) | Next.js 15, React 19, Tailwind v4, TypeScript | Secure management portal for course moderation, user management, financial revenue & broadcasts |
| 📱 **Mobile App** | [`/memere_mobile`](memere_mobile/README.md) | Flutter 3.13+, Riverpod 2.5, GoRouter, Hive | Student-facing app for HLS video learning, offline PDF notes, quizzes, mock exams & mobile payments |

---

## 🏗️ System Architecture

```
                                 ┌─────────────────────────┐
                                 │   Grade 12 Students     │
                                 └────────────┬────────────┘
                                              │
                                              ▼ (HTTP/REST - Bearer JWT Token)
                                   ┌─────────────────────┐
                                   │ Flutter Mobile App  │
                                   └─────────────────────┘
                                              │
                                              │
┌────────────────────────┐                    │
│ Admin / Content Creator│                    │
└───────────┬────────────┘                    │
            │                                 │
            ▼ (HTTP - Session Cookie)         │
 ┌─────────────────────┐                      │
 │ Next.js Admin Panel │                      │
 └──────────┬──────────┘                      │
            │ (Server-to-Server Proxy)        │
            └─────────────────┬───────────────┘
                              │
                              ▼
                ┌───────────────────────────┐
                │   Go REST API (Backend)   │
                └─────────────┬─────────────┘
                              │
       ┌──────────────────────┼──────────────────────┐
       ▼                      ▼                      ▼
┌──────────────┐      ┌──────────────┐      ┌──────────────┐
│ PostgreSQL 15│      │   Redis 7    │      │  MinIO / S3  │
│ (Primary DB) │      │(Cache/Queue) │      │(Video & PDF) │
└──────────────┘      └──────────────┘      └──────────────┘
                              ▲
                              │ (Async Video Processing / Transcoding)
                      ┌───────┴──────┐
                      │  Go Worker   │
                      └──────────────┘
```

---

## ✨ Component Breakdown & Features

### 1. ⚙️ Go Backend (`Memere-backend`)
- **Clean Architecture**: Domain entities, use cases, HTTP handlers, and repository interfaces.
- **High Performance**: Built with Go 1.22+ for microsecond latency and concurrent student connections.
- **Media Transcoding**: Asynchronous FFmpeg worker service to process raw video uploads into HLS adaptive bitrate streams.
- **Payment Gateway**: Integration with Chapa, Telebirr, and Stripe webhooks.
- **Database & Cache**: PostgreSQL 15 with SQLC generated type-safe queries, Redis 7 for user sessions, quiz states, and leaderboards.
- **Documentation**: Includes Postman collections and OpenAPI specifications.

### 2. 💻 Admin Web Portal (`Memere-admin`)
- **Security-First SSR Architecture**: Next.js 15 App Router acting as a secure server-side proxy. Raw JWTs are stored strictly in `httpOnly` cookies; browser JS never handles tokens directly.
- **Live Dashboards**: Real-time revenue charts (Recharts), engagement KPIs, platform statistics.
- **Course & User Moderation**: Content publishing workflows, lesson structure management, user role assignments, suspensions.
- **Financial Management**: Revenue reconciliation, transaction logs, payment status filtering, refund processing.
- **Broadcast System**: Audience segment composer for targeted push notifications and announcements.

### 3. 📱 Mobile Application (`memere_mobile`)
- **Student-First UX**: Clean, accessible Material 3 interface optimized for phone screens and low-bandwidth connections.
- **Adaptive Video Player**: HLS streaming engine with resolution switching, playback speed control, and video caching.
- **PDF Reader**: Built-in document viewer for downloadable lesson notes and study guides.
- **Quiz & Timed Exam Engine**: Practice per-lesson quizzes or simulate national university entrance exams with real-time timers and performance scoring.
- **Offline Mode**: Save videos and PDF notes to local storage (Hive) for offline study sessions.
- **Mobile Payments**: Seamless course purchasing via Chapa and Telebirr WebViews.

---

## 🛠️ Global Prerequisites

To run all components of the Memere platform locally, ensure you have installed:

- **Go**: `1.22+`
- **Node.js**: `20.x LTS` & **pnpm**: `11+`
- **Flutter SDK**: `3.13+` & **Dart**: `3.0+`
- **PostgreSQL**: `15`
- **Redis**: `7`
- **Docker** & **Docker Compose** (recommended for full infrastructure)

---

## 🚀 Quick Start Guide

### Step 1: Start Infrastructure & Backend Services

Navigate to `Memere-backend` and launch PostgreSQL, Redis, MinIO, and the Go API:

```bash
cd Memere-backend

# Start infrastructure (Postgres, Redis, MinIO)
make up

# Run database migrations
make migrate-up

# Start the Go REST API
make run
```
> The Go API will be running at `http://localhost:8080`.

---

### Step 2: Launch the Admin Web Panel

In a new terminal window:

```bash
cd Memere-admin

# Install dependencies
pnpm install

# Copy environment variables
cp .env.local.example .env.local

# Run Next.js dev server
pnpm dev
```
> The Admin panel will be available at `http://localhost:3000`.

---

### Step 3: Run the Flutter Mobile App

In a third terminal window:

```bash
cd memere_mobile

# Install Flutter packages
flutter pub get

# Generate compile-time env and riverpod code
dart run build_runner build --delete-conflicting-outputs

# Start app on emulator or target device
flutter run
```

---

## 🔐 Environment Configuration Overview

Each subsystem relies on its own environment file:

| Service | Environment File | Key Variables |
|---|---|---|
| **Backend** | `Memere-backend/.env` | `PORT`, `DB_HOST`, `DB_PORT`, `DB_NAME`, `JWT_SECRET`, `CHAPA_SECRET_KEY`, `S3_BUCKET` |
| **Admin** | `Memere-admin/.env.local` | `API_BASE_URL` (Go API), `COOKIE_SECRET`, `NODE_ENV` |
| **Mobile** | `memere_mobile/.env` | `BASE_URL` (Loaded via `envied`), `ENABLE_ANALYTICS` |

---

## 🧪 Verification & Testing

Run tests across all modules:

```bash
# Test Backend Go Code
cd Memere-backend && make test

# Typecheck & Lint Admin Web Panel
cd Memere-admin && pnpm typecheck && pnpm lint

# Test Mobile Flutter Code
cd memere_mobile && flutter test
```

---

## 📄 License & Documentation

- **License**: MIT License
- **Backend Documentation**: [`Memere-backend/README.md`](Memere-backend/README.md)
- **Admin Documentation**: [`Memere-admin/README.md`](Memere-admin/README.md)
- **Mobile Documentation**: [`memere_mobile/README.md`](memere_mobile/README.md)
- **API Spec / Postman**: [`Memere-backend/docs/Memere.postman_collection.json`](Memere-backend/docs/Memere.postman_collection.json)

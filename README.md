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

## ✨ Component Breakdown & Features

### 1. ⚙️ Go Backend (`Memere-backend`)
- **Clean Architecture**: Domain entities, use cases, HTTP handlers, and repository interfaces.
- **High Performance**: Built with Go 1.22+ for microsecond latency and concurrent student connections.
- **Media Transcoding**: Asynchronous FFmpeg worker service to process raw video uploads into HLS adaptive bitrate streams.
- **Payment Gateway**: Integration with Chapa, Telebirr, and Stripe webhooks.
- **Database & Cache**: PostgreSQL 15 with SQLC generated type-safe queries, Redis 7 for user sessions, quiz states, and leaderboards.

### 2. 💻 Admin Web Portal (`Memere-admin`)
- **Security-First SSR Architecture**: Next.js 15 App Router acting as a secure server-side proxy. Raw JWTs are stored strictly in `httpOnly` cookies; browser JS never handles tokens directly.
- **Live Dashboards**: Real-time revenue charts (Recharts), engagement KPIs, platform statistics.
- **Course & User Moderation**: Content publishing workflows, lesson structure management, user role assignments, suspensions.

### 3. 📱 Mobile Application (`memere_mobile`)
- **Student-First UX**: Clean, accessible Material 3 interface optimized for phone screens and low-bandwidth connections.
- **Adaptive Video Player**: HLS streaming engine with resolution switching, playback speed control, and video caching.
- **PDF Reader**: Built-in document viewer for downloadable lesson notes and study guides.
- **Quiz & Timed Exam Engine**: Practice per-lesson quizzes or simulate national university entrance exams with real-time timers and performance scoring.

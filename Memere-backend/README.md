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

# 📱 Memere Mobile — Grade 12 Exam Prep Flutter App

<p align="center">
  <img src="https://img.shields.io/badge/Flutter-3.13+-02569B?style=for-the-badge&logo=flutter&logoColor=white" alt="Flutter">
  <img src="https://img.shields.io/badge/Dart-3.0+-0175C2?style=for-the-badge&logo=dart&logoColor=white" alt="Dart">
  <img src="https://img.shields.io/badge/Riverpod-2.5-blue?style=for-the-badge&logo=flutter&logoColor=white" alt="Riverpod">
  <img src="https://img.shields.io/badge/Platform-Android%20%7C%20iOS-brightgreen?style=for-the-badge" alt="Platform">
  <img src="https://img.shields.io/badge/License-MIT-green?style=for-the-badge" alt="License">
</p>

---

## 📖 Overview

**Memere Mobile** is a mobile-first learning application built with Flutter specifically for Grade 12 students in Ethiopia preparing for university entrance examinations. The app provides low-latency adaptive video streaming, downloadable offline study materials, interactive quizzes, timed national mock entrance exams, and seamless mobile money payment options (Chapa and Telebirr).

---

## ✨ Key Features

### 🎓 Learning & Examination
- 🎥 **Adaptive HLS Video Streaming** — Smooth video streaming with playback speed control, resolution switching, and Chewie/VideoPlayer integration.
- 📄 **PDF Study Notes Viewer** — In-app PDF reader powered by Syncfusion PDF Engine for chapter summaries and study sheets.
- ❓ **Per-Lesson Quizzes** — Interactive quizzes with instant feedback, explanations, and score tracking.
- 📋 **Timed Mock Entrance Exams** — Simulated national university entrance exams with countdown timers, question navigation, and detailed score breakdown.
- 📊 **Progress & Analytics** — Student progress dashboard tracking completed lessons, subject mastery, quiz history, and weak areas.

### 💳 Payment & Subscriptions
- 📱 **Mobile Money WebView Flow** — Integrated payment flow supporting Ethiopian payment gateways (Chapa & Telebirr) via `webview_flutter`.
- 🎟️ **Coupons & Bundles** — Promo code redemption and course bundle purchases.

### 📴 Offline Support & Media Management
- 💾 **Offline Saved Lessons** — Download video lessons and PDF documents for offline access in areas with limited internet.
- 🔄 **Background Sync** — Periodic offline progress sync using `workmanager` and `connectivity_plus`.

### 🔔 Notifications & Security
- 🔔 **Push Notifications** — Reminders and content updates powered by Firebase Cloud Messaging (FCM) and `flutter_local_notifications`.
- 🔐 **Secure Token Storage** — Authentication tokens securely saved via `flutter_secure_storage`.

---

## 🏗️ Architecture & Tech Stack

The application follows **Clean Architecture** combined with a **Feature-First** folder structure to maximize testability and maintainability.

```
lib/
├── app.dart                   # Main App widget & router setup
├── main.dart                  # Entry point (Firebase & DI init)
├── core/                      # Global core utilities & services
│   ├── constants/             # API routes, app colors, strings
│   ├── di/                    # Dependency injection providers
│   ├── errors/                # Failure & Exception classes
│   ├── network/               # Dio HTTP client, interceptors, auth token refresh
│   ├── router/                # GoRouter route definitions & guards
│   ├── storage/               # Hive boxes & SecureStorage wrappers
│   ├── theme/                 # Material 3 dark/light themes & Google Fonts
│   └── utils/                 # Helpers, formatters, validators
├── features/                  # Feature-based modular domains
│   ├── auth/                  # Login, Register, OTP Verification, Password Reset
│   ├── courses/               # Catalog, Subject Filtering, Course Detail
│   ├── exam/                  # Timed Mock Exams, Question Runner, Scoring
│   ├── learning/              # Lesson Viewer (Video + PDF Notes)
│   ├── notifications/         # In-App & FCM Push Notification Handlers
│   ├── payment/               # Chapa / Telebirr WebView Flow & Receipt
│   ├── profile/               # User Settings, Dark Mode, Account Management
│   ├── progress/              # Student Performance Dashboard & Analytics
│   ├── quiz/                  # Per-Lesson Quiz Runner & Feedback
│   ├── saved/                 # Bookmarked & Offline Content Management
│   └── video_player/          # HLS Video Controller & Player UI
└── shared/                    # Reusable widgets, state indicators, buttons
```

### Tech Stack Details

| Layer / Subsystem | Library / Technology |
|---|---|
| **Framework & Language** | Flutter `^3.13.0`, Dart `^3.0.0` |
| **State Management** | `flutter_riverpod` (v2.5.1) + `riverpod_annotation` |
| **Routing & Navigation** | `go_router` (v13.2.0) |
| **Networking & HTTP** | `dio` (v5.4.3) with token refresh interceptors |
| **Local Storage** | `hive_flutter`, `flutter_secure_storage`, `shared_preferences` |
| **Immutability & Code Gen** | `freezed`, `json_annotation`, `build_runner` |
| **Functional Error Handling**| `fpdart` |
| **Video Player** | `chewie`, `video_player`, `flutter_hls_parser` |
| **PDF Viewer** | `syncfusion_flutter_pdfviewer`, `flutter_pdfview` |
| **Push Notifications** | `firebase_messaging`, `flutter_local_notifications` |
| **Offline Background Sync** | `workmanager`, `connectivity_plus` |
| **UI Components** | Google Fonts, Lucide Icons / Flutter SVG, Shimmer, Lottie |

---

## 🛠️ Prerequisites & Setup

### Prerequisites

- **Flutter SDK**: `>= 3.13.0`
- **Dart SDK**: `>= 3.0.0`
- **Android Studio** (for Android build) / **Xcode 15+** (for iOS build)
- **CocoaPods** (for iOS dependencies)

### Installation

1. **Clone the repository & navigate to the mobile folder:**
   ```bash
   cd memere_mobile
   ```

2. **Install Flutter dependencies:**
   ```bash
   flutter pub get
   ```

3. **Configure Environment Variables:**
   Create a `.env` file in the `memere_mobile/` root directory (loaded at compile-time via `envied`):
   ```ini
   BASE_URL=http://10.0.2.2:8080/api/v1
   ENABLE_ANALYTICS=false
   ```

4. **Run Code Generation:**
   Generate Riverpod providers, Freezed models, JSON serializers, and Envied variables:
   ```bash
   dart run build_runner build --delete-conflicting-outputs
   ```

---

## 🚀 Running the App

### Development Mode

Run on an connected Android emulator or physical device:
```bash
flutter run
```

Run on a specific device:
```bash
flutter devices
flutter run -d <device-id>
```

### Code Generation & Watch Mode

If you are modifying models or state providers, run build runner in watch mode:
```bash
dart run build_runner watch --delete-conflicting-outputs
```

---

## 🧪 Testing & Quality Assurance

### Run Unit & Widget Tests

```bash
flutter test
```

### Static Code Analysis

```bash
flutter analyze
```

---

## 📦 Building for Production

### Android APK
```bash
flutter build apk --release
```

### Android App Bundle (AAB for Play Store)
```bash
flutter build appbundle --release
```

### iOS IPA (App Store)
```bash
flutter build ipa --release
```

---

## 🔗 Related Components

- ⚙️ **Backend API**: [`../Memere-backend`](../Memere-backend/README.md) — High-performance Go REST API & worker service
- 💻 **Admin Web Panel**: [`../Memere-admin`](../Memere-admin/README.md) — Next.js 15 management portal for course creators and administrators
- 🌐 **Project Root**: [`../README.md`](../README.md) — Full ecosystem overview and architecture diagram

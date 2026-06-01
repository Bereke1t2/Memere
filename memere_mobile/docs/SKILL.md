# SKILL.md — Antigravity Global Rules
# ExamPrep Mobile (memere_mobile)
# READ THIS BEFORE EVERY PROMPT IN EVERY PHASE

---

## 1. PROJECT IDENTITY

| Key | Value |
|-----|-------|
| App Name | ExamPrep (internal codename: memere) |
| Platform | Flutter (Android-first, then iOS Month 6) |
| State Management | Riverpod (AsyncNotifierProvider, FutureProvider, StateNotifierProvider) |
| Navigation | GoRouter with redirect guards |
| Architecture | Clean Architecture — strict layer separation |
| Local Storage | Hive (course cache) + Flutter Secure Storage (tokens) + SharedPreferences (prefs) |
| HTTP Client | Dio with interceptors (JWT refresh on 401) |
| Code Generation | Freezed + json_serializable + build_runner |
| Design Language | Dark-first, minimal, high-contrast — inspired by ChatGPT/Claude mobile UI |
| Language | Dart (null-safe, no dynamic types) |
| Min SDK | Android 21 (5.0), iOS 13 |

---

## 2. MASTER FOLDER STRUCTURE

```
memere_mobile/
├── lib/
│   ├── main.dart
│   ├── app.dart
│   ├── core/
│   │   ├── constants/          # AppColors, AppTextStyles, AppSizes, Endpoints
│   │   ├── errors/             # Failure classes (ServerFailure, CacheFailure, NetworkFailure)
│   │   ├── network/            # DioClient, interceptors, connectivity
│   │   ├── storage/            # HiveService, SecureStorageService, PrefsService
│   │   ├── theme/              # AppTheme (dark + light)
│   │   ├── router/             # AppRouter (GoRouter), route names, guards
│   │   └── di/                 # All Riverpod providers wired here
│   ├── features/
│   │   ├── auth/
│   │   │   ├── data/
│   │   │   │   ├── datasources/   # AuthRemoteDataSource, AuthLocalDataSource
│   │   │   │   ├── models/        # UserModel (extends UserEntity, has fromJson/toJson)
│   │   │   │   └── repositories/  # AuthRepositoryImpl
│   │   │   ├── domain/
│   │   │   │   ├── entities/      # UserEntity (pure Dart, no JSON deps)
│   │   │   │   ├── repositories/  # AuthRepository (abstract interface)
│   │   │   │   └── usecases/      # LoginUseCase, RegisterUseCase, LogoutUseCase
│   │   │   └── presentation/
│   │   │       ├── providers/     # authStateProvider, loginNotifierProvider
│   │   │       ├── screens/       # SplashScreen, LoginScreen, RegisterScreen
│   │   │       └── widgets/       # AuthTextField, PrimaryButton, SocialLoginRow
│   │   ├── courses/               # same 3-layer structure
│   │   ├── video_player/
│   │   ├── quiz/
│   │   ├── exam/
│   │   ├── payment/
│   │   ├── progress/
│   │   └── notifications/
│   └── shared/
│       ├── widgets/               # AppButton, AppTextField, LoadingOverlay, ErrorWidget
│       ├── extensions/            # String, DateTime, int, BuildContext extensions
│       └── utils/                 # Formatters, validators, debouncer
├── docs/
│   └── memere_Design_Specification.md   # ← ALREADY EXISTS — always reference this
├── test/
│   ├── unit/
│   ├── widget/
│   └── integration/
├── assets/
│   ├── images/
│   ├── icons/
│   ├── fonts/
│   └── lottie/
├── pubspec.yaml
├── analysis_options.yaml
├── .env.example
└── README.md
```

---

## 3. ARCHITECTURE NON-NEGOTIABLES (FROM SPEC DOC)

These rules are ABSOLUTE. Never violate them:

1. **Correct answers NEVER sent to client** — always grade server-side
2. **Exam timer MUST be enforced server-side** — client timer is display-only
3. **Pre-signed CDN URLs ONLY for video** — no public S3 access
4. **All payments must use idempotency keys** — no double-charge risk
5. **Soft deletes only** — never hard DELETE user-facing data
6. **HTTPS-only** — HTTP redirects to HTTPS at gateway
7. **Never log raw passwords, tokens, or payment card data**
8. **All DB queries filter by authenticated user_id** (prevent IDOR)

---

## 4. FLUTTER LAYER RULES

### Domain Layer (`domain/`)
- Pure Dart only — ZERO external package imports
- Entities are plain Dart classes (no JSON, no DB annotations)
- Repository interfaces define contracts only — no implementation
- Use cases: single responsibility, one public `call()` method

### Data Layer (`data/`)
- Models extend entities and add `fromJson` / `toJson`
- Remote datasources use `DioClient` only — never call Dio directly in repos
- Local datasources use `HiveService` or `SecureStorageService`
- Repository impls handle error mapping → Failure types

### Presentation Layer (`presentation/`)
- Screens contain ZERO business logic — delegate to providers
- Providers call use cases only — never call repos directly
- Every async provider returns `AsyncValue<T>` — handle loading/error/data
- No `setState` in any feature screen — Riverpod only

---

## 5. NAMING CONVENTIONS

| Type | Convention | Example |
|------|-----------|---------|
| Files | snake_case | `login_screen.dart` |
| Classes | PascalCase | `LoginScreen`, `AuthRepository` |
| Providers | camelCase + suffix | `authStateProvider`, `loginNotifierProvider` |
| Use Cases | PascalCase + UseCase | `LoginUseCase` |
| Models | PascalCase + Model | `UserModel` |
| Entities | PascalCase + Entity | `UserEntity` |
| Routes | `/kebab-case` | `/course-detail` |
| Constants | SCREAMING_SNAKE | `BASE_URL`, `ACCESS_TOKEN_KEY` |
| Private vars | `_camelCase` | `_isLoading` |

---

## 6. DESIGN SYSTEM (ALWAYS REFERENCE memere_Design_Specification.md)

### Color Palette (Dark-First)
```dart
// Primary Background
static const Color bgPrimary    = Color(0xFF0D0D0D);  // near-black
static const Color bgSecondary  = Color(0xFF1A1A1A);  // card bg
static const Color bgTertiary   = Color(0xFF252525);  // elevated surface

// Accent
static const Color accentPrimary   = Color(0xFF6C63FF); // purple-blue
static const Color accentSecondary = Color(0xFF03DAC6); // teal
static const Color accentWarning   = Color(0xFFFFB347); // amber

// Text
static const Color textPrimary   = Color(0xFFF5F5F5);
static const Color textSecondary = Color(0xFF9E9E9E);
static const Color textDisabled  = Color(0xFF4A4A4A);

// Status
static const Color success = Color(0xFF4CAF50);
static const Color error   = Color(0xFFCF6679);
static const Color warning = Color(0xFFFFB347);
```

### Typography (Sora display + DM Sans body)
```dart
// Display / Headlines — Sora
headlineLarge:  Sora, 28sp, w700
headlineMedium: Sora, 22sp, w600
headlineSmall:  Sora, 18sp, w600

// Body — DM Sans
bodyLarge:   DM Sans, 16sp, w400
bodyMedium:  DM Sans, 14sp, w400
bodySmall:   DM Sans, 12sp, w400

// Labels / Buttons
labelLarge:  DM Sans, 14sp, w600
labelMedium: DM Sans, 12sp, w500
```

### Spacing System (8px grid)
```dart
xs  = 4.0
sm  = 8.0
md  = 16.0
lg  = 24.0
xl  = 32.0
xxl = 48.0
```

### Border Radius
```dart
radiusSm = 8.0
radiusMd = 12.0
radiusLg = 16.0
radiusXl = 24.0
radiusFull = 999.0
```

---

## 7. ERROR HANDLING PATTERN

Always use `Either<Failure, T>` from `fpdart` or a custom Result type:

```dart
// In repository impl:
try {
  final result = await remoteDataSource.login(params);
  return Right(result.toEntity());
} on DioException catch (e) {
  return Left(ServerFailure.fromDioError(e));
} on HiveError catch (e) {
  return Left(CacheFailure(e.message));
} catch (e) {
  return Left(UnknownFailure(e.toString()));
}

// In use case:
final result = await repository.login(params);
return result.fold(
  (failure) => Left(failure),
  (user) => Right(user),
);
```

---

## 8. PROVIDER PATTERN

```dart
// AsyncNotifierProvider pattern (preferred for complex state)
@riverpod
class LoginNotifier extends _$LoginNotifier {
  @override
  FutureOr<void> build() {}

  Future<void> login(String email, String password) async {
    state = const AsyncLoading();
    final useCase = ref.read(loginUseCaseProvider);
    final result = await useCase(LoginParams(email: email, password: password));
    result.fold(
      (failure) => state = AsyncError(failure, StackTrace.current),
      (_) => state = const AsyncData(null),
    );
  }
}
```

---

## 9. SECURITY RULES (FLUTTER SIDE)

- JWT access token stored in **Flutter Secure Storage** (never SharedPreferences)
- Refresh token stored in **Flutter Secure Storage**
- On 401 response: Dio interceptor silently refreshes token, retries request
- On refresh failure: clear all tokens → redirect to login
- Never log token values, even in debug mode
- Video download: encrypt with AES-256 using device-bound key from Secure Storage
- Downloaded content: expires 30 days — check timestamp on every play

---

## 10. OFFLINE STRATEGY

```
On app launch:
  1. Check connectivity (connectivity_plus)
  2. If offline → serve from Hive cache, queue writes to SyncQueue
  3. If online → fetch fresh data, update Hive cache (TTL = 1 hour)

On reconnect:
  1. Flush SyncQueue to backend (idempotent endpoints)
  2. Conflict resolution: server-wins except quiz answers (client-wins)
  3. Background sync via WorkManager (Android)
```

---

## 11. API CONVENTIONS

- Base URL: `https://api.examprep.et/api/v1`
- All requests: `Authorization: Bearer <access_token>`
- Pagination: `?limit=20&after=<cursor>`
- Error response: `{ "code": "ERROR_CODE", "message": "...", "details": {} }`
- File uploads: pre-signed S3 URLs — never upload directly through backend

---

## 12. PHASE OVERVIEW

| Phase | Focus | Key Deliverables |
|-------|-------|-----------------|
| **Phase 1** (current) | Foundation | Scaffold, design system, auth feature, splash/onboarding/login/register screens |
| **Phase 2** | Course browsing | Course list, course detail, section/lesson list, search/filter |
| **Phase 3** | Video player | HLS streaming, offline download, progress saving |
| **Phase 4** | Quiz engine | Quiz screen, question flow, submission, results |
| **Phase 5** | Mock exam | Timed exam, server-synced timer, results + analytics |
| **Phase 6** | Payment | Chapa/Telebirr integration, enrollment, purchase flow |
| **Phase 7** | Progress + Dashboard | Student dashboard, streaks, completion tracking |
| **Phase 8** | Notifications | FCM push, in-app notification center |
| **Phase 9** | Teacher flow | Content management (teacher-only screens) |
| **Phase 10** | Polish + Release | Performance, accessibility, Play Store submission |

---

## 13. HOW TO USE THESE SKILLS IN ANTIGRAVITY

1. **Start every session** by telling Antigravity: `"Read SKILL.md and all phase1/ skill files before proceeding"`
2. **Reference the spec**: `"The full architecture is documented in memere_mobile/docs/memere_Design_Specification.md"`
3. **Be explicit about phase**: `"We are on Phase 1. Build only what Phase 1 specifies."`
4. **Divide large tasks**: Ask Antigravity to build one feature at a time within a phase
5. **Moving to next phase**: Say `"Phase 1 is complete. Read phase2/SKILL_1.md and begin Phase 2"`

---

## 14. WHAT ANTIGRAVITY MUST NEVER DO

- Import packages not in `pubspec.yaml` without asking first
- Use `setState` in feature screens (Riverpod only)
- Skip the domain layer and call repositories from presentation
- Put business logic in widgets or screens
- Use `dynamic` type
- Use `print()` — use `debugPrint()` or a logger package
- Hard-delete data
- Send correct quiz/exam answers in API responses
- Hardcode API URLs or secrets — use `AppConstants` and `.env`

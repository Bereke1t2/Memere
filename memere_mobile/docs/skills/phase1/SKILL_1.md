# phase1/SKILL_1.md — Project Scaffold, pubspec.yaml & Folder Structure
# ExamPrep Mobile (memere_mobile) — Phase 1, Part 1
# READ SKILL.md FIRST, then this file.

---

## OBJECTIVE

Build the complete project skeleton with zero business logic. By the end of this skill,
the project must compile cleanly, run on an Android emulator, and show a blank dark screen
with the correct background color. No features yet — just the foundation everything else
will be built on.

---

## STEP 1 — pubspec.yaml (Complete)

Replace the entire `pubspec.yaml` with this exact content:

```yaml
name: memere_mobile
description: ExamPrep — Grade 12 University Entrance Exam Prep Platform
publish_to: 'none'
version: 1.0.0+1

environment:
  sdk: '>=3.0.0 <4.0.0'
  flutter: '>=3.13.0'

dependencies:
  flutter:
    sdk: flutter

  # ── State Management ─────────────────────────────────────
  flutter_riverpod: ^2.5.1
  riverpod_annotation: ^2.3.5

  # ── Navigation ───────────────────────────────────────────
  go_router: ^13.2.0

  # ── Networking ───────────────────────────────────────────
  dio: ^5.4.3
  connectivity_plus: ^6.0.3

  # ── Local Storage ────────────────────────────────────────
  hive_flutter: ^1.1.0
  flutter_secure_storage: ^9.0.0
  shared_preferences: ^2.2.3
  path_provider: ^2.1.3

  # ── Code Generation & Utilities ──────────────────────────
  freezed_annotation: ^2.4.1
  json_annotation: ^4.9.0

  # ── Functional Error Handling ────────────────────────────
  fpdart: ^1.1.0

  # ── UI & Design ──────────────────────────────────────────
  google_fonts: ^6.2.1
  flutter_svg: ^2.0.10+1
  cached_network_image: ^3.3.1
  shimmer: ^3.0.0
  lottie: ^3.1.2

  # ── Video ────────────────────────────────────────────────
  video_player: ^2.8.6
  chewie: ^1.7.5
  flutter_hls_parser: ^2.2.1

  # ── Utilities ────────────────────────────────────────────
  intl: ^0.19.0
  uuid: ^4.4.0
  equatable: ^2.0.5
  logger: ^2.3.0
  envied: ^0.5.4+1
  package_info_plus: ^8.0.0

  # ── Offline / Sync ───────────────────────────────────────
  workmanager: ^0.5.2

  # ── Notifications ────────────────────────────────────────
  firebase_core: ^3.3.0
  firebase_messaging: ^15.1.0
  flutter_local_notifications: ^17.2.2

  # ── Payments ─────────────────────────────────────────────
  webview_flutter: ^4.8.0        # for Chapa/Telebirr WebView payment flow
  url_launcher: ^6.3.0

dev_dependencies:
  flutter_test:
    sdk: flutter
  flutter_lints: ^4.0.0
  build_runner: ^2.4.9
  riverpod_generator: ^2.4.3
  freezed: ^2.5.2
  json_serializable: ^6.8.0
  hive_generator: ^2.0.1
  envied_generator: ^0.5.4+1
  mockito: ^5.4.4
  flutter_launcher_icons: ^0.13.1

flutter:
  uses-material-design: true

  assets:
    - assets/images/
    - assets/icons/
    - assets/lottie/
    - .env                        # loaded via envied — not committed to git

  fonts:
    - family: Sora
      fonts:
        - asset: assets/fonts/Sora-Regular.ttf
          weight: 400
        - asset: assets/fonts/Sora-SemiBold.ttf
          weight: 600
        - asset: assets/fonts/Sora-Bold.ttf
          weight: 700
    - family: DM_Sans
      fonts:
        - asset: assets/fonts/DMSans-Regular.ttf
          weight: 400
        - asset: assets/fonts/DMSans-Medium.ttf
          weight: 500
        - asset: assets/fonts/DMSans-SemiBold.ttf
          weight: 600
```

---

## STEP 2 — Create ALL Folders

Run this exact command from the project root (`memere_mobile/`):

```bash
# Core folders
mkdir -p lib/core/constants
mkdir -p lib/core/errors
mkdir -p lib/core/network
mkdir -p lib/core/storage
mkdir -p lib/core/theme
mkdir -p lib/core/router
mkdir -p lib/core/di

# Auth feature (full 3-layer clean architecture)
mkdir -p lib/features/auth/data/datasources
mkdir -p lib/features/auth/data/models
mkdir -p lib/features/auth/data/repositories
mkdir -p lib/features/auth/domain/entities
mkdir -p lib/features/auth/domain/repositories
mkdir -p lib/features/auth/domain/usecases
mkdir -p lib/features/auth/presentation/providers
mkdir -p lib/features/auth/presentation/screens
mkdir -p lib/features/auth/presentation/widgets

# Stub folders for future phases (empty — just structure)
mkdir -p lib/features/courses/data/datasources
mkdir -p lib/features/courses/data/models
mkdir -p lib/features/courses/data/repositories
mkdir -p lib/features/courses/domain/entities
mkdir -p lib/features/courses/domain/repositories
mkdir -p lib/features/courses/domain/usecases
mkdir -p lib/features/courses/presentation/providers
mkdir -p lib/features/courses/presentation/screens
mkdir -p lib/features/courses/presentation/widgets

mkdir -p lib/features/video_player/presentation/screens
mkdir -p lib/features/video_player/presentation/widgets
mkdir -p lib/features/quiz/presentation/screens
mkdir -p lib/features/exam/presentation/screens
mkdir -p lib/features/payment/presentation/screens
mkdir -p lib/features/progress/presentation/screens
mkdir -p lib/features/notifications/presentation/screens

# Shared
mkdir -p lib/shared/widgets
mkdir -p lib/shared/extensions
mkdir -p lib/shared/utils

# Assets
mkdir -p assets/images
mkdir -p assets/icons
mkdir -p assets/lottie
mkdir -p assets/fonts

# Tests (mirrors lib structure)
mkdir -p test/unit/features/auth
mkdir -p test/widget/features/auth
mkdir -p test/integration

# Docs (already exists)
mkdir -p docs
```

---

## STEP 3 — Environment Config

### `.env` (create in project root — NEVER commit to git)
```
BASE_URL=https://api.examprep.et/api/v1
CHAPA_PUBLIC_KEY=your_chapa_public_key_here
FIREBASE_PROJECT_ID=your_firebase_project_id
```

### `.env.example` (safe to commit)
```
BASE_URL=https://api.examprep.et/api/v1
CHAPA_PUBLIC_KEY=
FIREBASE_PROJECT_ID=
```

### `lib/core/constants/env.dart`
```dart
import 'package:envied/envied.dart';

part 'env.g.dart';

@Envied(path: '.env')
abstract class Env {
  @EnviedField(varName: 'BASE_URL')
  static const String baseUrl = _Env.baseUrl;

  @EnviedField(varName: 'CHAPA_PUBLIC_KEY', obfuscate: true)
  static const String chapaPublicKey = _Env.chapaPublicKey;

  @EnviedField(varName: 'FIREBASE_PROJECT_ID')
  static const String firebaseProjectId = _Env.firebaseProjectId;
}
```

### `.gitignore` (add these lines)
```
.env
*.g.dart           # generated — can be toggled based on team preference
lib/core/constants/env.g.dart
```

---

## STEP 4 — `analysis_options.yaml`

```yaml
include: package:flutter_lints/flutter.yaml

analyzer:
  errors:
    missing_required_param: error
    missing_return: error
    dead_code: warning
  exclude:
    - '**/*.g.dart'
    - '**/*.freezed.dart'

linter:
  rules:
    - avoid_dynamic_calls
    - avoid_print              # use logger or debugPrint
    - prefer_const_constructors
    - prefer_final_fields
    - prefer_single_quotes
    - sort_child_properties_last
    - use_key_in_widget_constructors
```

---

## STEP 5 — `main.dart`

```dart
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:hive_flutter/hive_flutter.dart';
import 'package:firebase_core/firebase_core.dart';
import 'app.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();

  // Lock orientation to portrait
  await SystemChrome.setPreferredOrientations([
    DeviceOrientation.portraitUp,
    DeviceOrientation.portraitDown,
  ]);

  // Dark status bar
  SystemChrome.setSystemUIOverlayStyle(
    const SystemUiOverlayStyle(
      statusBarColor: Colors.transparent,
      statusBarIconBrightness: Brightness.light,
      systemNavigationBarColor: Color(0xFF0D0D0D),
      systemNavigationBarIconBrightness: Brightness.light,
    ),
  );

  // Initialize Hive
  await Hive.initFlutter();

  // Initialize Firebase
  await Firebase.initializeApp();

  runApp(
    const ProviderScope(
      child: ExamPrepApp(),
    ),
  );
}
```

---

## STEP 6 — `app.dart`

```dart
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'core/router/app_router.dart';
import 'core/theme/app_theme.dart';

class ExamPrepApp extends ConsumerWidget {
  const ExamPrepApp({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final router = ref.watch(appRouterProvider);

    return MaterialApp.router(
      title: 'ExamPrep',
      debugShowCheckedModeBanner: false,
      theme: AppTheme.dark(),
      darkTheme: AppTheme.dark(),
      themeMode: ThemeMode.dark,
      routerConfig: router,
    );
  }
}
```

---

## STEP 7 — Verify Build

After creating all files and running `flutter pub get`, confirm:

```bash
flutter pub get
dart run build_runner build --delete-conflicting-outputs
flutter analyze
flutter run
```

**Expected result**: App launches, shows black screen (0xFF0D0D0D), no errors, no warnings (except "missing fonts" until font files are added).

---

## PHASE 1 CHECKLIST — SKILL_1 COMPLETE WHEN:

- [ ] `pubspec.yaml` matches exactly as above
- [ ] All folders created
- [ ] `.env` and `.env.example` created
- [ ] `main.dart` compiles with ProviderScope + Hive + Firebase init
- [ ] `app.dart` references router and theme (stubs OK at this point)
- [ ] `flutter analyze` passes with 0 errors
- [ ] App runs on emulator showing dark background

---

## NEXT: Go to `phase1/SKILL_2.md` — Design System

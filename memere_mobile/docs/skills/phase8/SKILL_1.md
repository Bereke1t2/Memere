# phase8/SKILL_1.md - Notifications Domain, Data Layer, Backend Contract & FCM Platform Setup
# Memere Mobile (memere_mobile) - Phase 8, Part 1
# READ SKILL.md -> all phase1 skill files -> all phase2 skill files -> all phase3 skill files -> all phase4 skill files -> all phase5 skill files -> all phase6 skill files -> all phase7 skill files -> then this file.

---

## OBJECTIVE

Build the Notifications foundation:
notification entity -> preferences entity -> repository contract -> use cases ->
backend DTO models -> remote datasource -> repository implementation ->
Firebase/FCM platform setup -> device token registration service ->
Riverpod dependency providers.

By the end of this skill, the app can register its FCM device token with the
local Memere backend, load the in-app notification list, unread count, and
notification preferences. Screens, the notification center, push handling, and
settings UI are built in `SKILL_2.md` and `SKILL_3.md`.

---

## PHASE 8 BOUNDARY

Phase 8 builds **student-facing notifications only**:

- FCM push setup (Firebase init, permission request)
- device token registration / unregistration
- in-app notification center
- notification read / unread state and unread badge
- notification preferences (push / email toggles)
- foreground / background / terminated push handling
- deep-link routing from a notification tap

Do **not** build these yet:

- composing or sending notifications from the app
- teacher/admin broadcast UI
- email or SMS channel UI
- scheduled local reminders the backend does not drive
- notification grouping/threading beyond a flat list
- websocket realtime stream (Phase 8 uses fetch + refresh + push wake-up only)

The backend owns notification creation. The app only **reads** notifications,
**marks** them read, **registers** device tokens, and **reads/writes**
preferences.

---

## BACKEND CONTRACT

These endpoints already exist in the Memere backend (verified, authenticated):

```text
GET    /me/notifications?limit=50
GET    /me/notifications/unread-count
POST   /me/notifications/:id/read
POST   /me/notifications/read-all
POST   /me/devices
DELETE /me/devices/:token
GET    /me/notification-preferences
PUT    /me/notification-preferences
```

Rules:

- All are behind auth — the auth interceptor attaches the bearer token.
- Do not invent new notification endpoints.
- Do not compute unread count client-side when the backend exposes it; use
  `/me/notifications/unread-count` as the source of truth, and only derive
  locally for optimistic UI between refreshes.

---

## RESPONSE SHAPES

### `GET /me/notifications`

```json
{
  "notifications": [
    {
      "id": "uuid",
      "type": "exam_graded",
      "title": "Exam graded",
      "body": "Your exam has been graded. Check your results.",
      "data": { "attempt_id": "uuid" },
      "read_at": "2026-01-01T00:00:00Z",
      "created_at": "2026-01-01T00:00:00Z"
    }
  ]
}
```

Important:

- `data` is a `Map<String, String>` and may be omitted/empty.
- `read_at` is omitted or null when unread. Treat **absent OR null** as unread.
- `created_at` is always present.
- The list is wrapped in a `notifications` key — do not assume a bare array.

### `GET /me/notifications/unread-count`

```json
{ "count": 3 }
```

### `POST /me/notifications/:id/read` — `204 No Content`
### `POST /me/notifications/read-all` — `204 No Content`

### `POST /me/devices` — `201 Created`

Request body:

```json
{ "fcm_token": "<token>", "platform": "android" }
```

`platform` is `android` or `ios`.

### `DELETE /me/devices/:token` — `204 No Content`

### `GET /me/notification-preferences`

```json
{ "push_enabled": true, "email_enabled": true }
```

### `PUT /me/notification-preferences`

Request body:

```json
{ "push_enabled": true, "email_enabled": false }
```

Returns the updated preferences in the same shape as the GET.

---

## NOTIFICATION TYPE CATALOG

The backend emits these `type` values with these `data` keys. Use them for
icon selection and deep-link routing (`SKILL_2.md`/`SKILL_3.md`). Treat unknown
types gracefully — never crash on a new type.

| type                  | data keys                     | tap destination (if route exists) |
|-----------------------|-------------------------------|-----------------------------------|
| `exam_graded`         | `attempt_id`                  | exam result/analytics screen      |
| `lesson_published`    | `lesson_id`, `course_id`      | course detail / lesson            |
| `purchase_confirmed`  | `payment_id`                  | purchases / course detail         |
| `certificate_ready`   | `course_id`                   | course progress / profile         |
| `streak_warning`      | `student_id`                  | dashboard                         |
| `subscription_expired`| `subscription_id`             | subscription/purchases            |
| `video_ready`         | `video_id`, `course_id`       | teacher-facing — ignore for student build |

Rules:

- Do not hardcode title/body locally — render the backend `title`/`body`.
- Use `type` only to pick the icon and the deep-link target.
- If a `data` key is missing, fall back to opening the notification center (no
  navigation), never crash.

---

## REQUIRED FOLDER STRUCTURE

Create or verify:

```bash
mkdir -p lib/features/notifications/data/datasources
mkdir -p lib/features/notifications/data/models
mkdir -p lib/features/notifications/data/repositories
mkdir -p lib/features/notifications/domain/entities
mkdir -p lib/features/notifications/domain/repositories
mkdir -p lib/features/notifications/domain/usecases
mkdir -p lib/features/notifications/presentation/providers
mkdir -p lib/features/notifications/presentation/screens
mkdir -p lib/features/notifications/presentation/widgets
mkdir -p lib/core/services
```

The `lib/features/notifications/{data,domain,presentation}` directories already
exist but are empty — fill them.

---

## PART A - DOMAIN ENTITIES

### FILE A1 - `lib/features/notifications/domain/entities/notification_type.dart`

Create an enum mapping the backend strings to a closed Dart type plus an
`unknown` fallback:

```dart
enum NotificationType {
  examGraded,
  lessonPublished,
  purchaseConfirmed,
  certificateReady,
  streakWarning,
  subscriptionExpired,
  videoReady,
  unknown;

  static NotificationType fromRaw(String raw) { /* map snake_case -> enum */ }
}
```

Rules:

- `fromRaw` returns `unknown` for any unrecognized value.
- Keep this enum pure Dart, no Flutter imports.

---

### FILE A2 - `lib/features/notifications/domain/entities/notification_entity.dart`

Fields:

- `id`
- `type` (`NotificationType`)
- `rawType` (`String` — keep the original backend string)
- `title`
- `body`
- `data` (`Map<String, String>`)
- `readAt` (`DateTime?`)
- `createdAt` (`DateTime`)

Helpers:

- `bool get isRead` -> `readAt != null`
- `bool get isUnread` -> `!isRead`
- `String? dataValue(String key)` -> safe lookup into `data`
- `String get timeAgoLabel` -> relative time using the existing DateTime
  extension from `shared/extensions`

Keep this entity independent from JSON and Flutter imports.

---

### FILE A3 - `lib/features/notifications/domain/entities/notification_preferences_entity.dart`

Fields:

- `pushEnabled`
- `emailEnabled`

Helpers:

- `NotificationPreferencesEntity copyWith({bool? pushEnabled, bool? emailEnabled})`

Keep pure Dart.

---

### FILE A4 - `lib/features/notifications/domain/entities/device_registration_entity.dart`

Fields:

- `fcmToken`
- `platform` (`String` — `android` or `ios`)

This is the value object passed to the register-device use case. Keep pure Dart.

---

## PART B - REPOSITORY CONTRACT

### FILE B1 - `lib/features/notifications/domain/repositories/notification_repository.dart`

Create:

```dart
abstract class NotificationRepository {
  Future<Either<Failure, List<NotificationEntity>>> getNotifications({int limit});
  Future<Either<Failure, int>> getUnreadCount();
  Future<Either<Failure, Unit>> markRead(String notificationId);
  Future<Either<Failure, Unit>> markAllRead();
  Future<Either<Failure, Unit>> registerDevice(DeviceRegistrationEntity device);
  Future<Either<Failure, Unit>> unregisterDevice(String fcmToken);
  Future<Either<Failure, NotificationPreferencesEntity>> getPreferences();
  Future<Either<Failure, NotificationPreferencesEntity>> updatePreferences(
    NotificationPreferencesEntity preferences,
  );
}
```

Rules:

- Return `Either<Failure, T>` like earlier phases (`Unit` for 204 endpoints).
- Do not expose DTO models outside the data layer.
- Do not make screens call this repository directly.

---

## PART C - USE CASES

Create one use case per repository method.

### FILE C1 - `get_notifications_usecase.dart`
Calls `repository.getNotifications(limit: ...)`. Default `limit = 50`.

### FILE C2 - `get_unread_count_usecase.dart`
Calls `repository.getUnreadCount()`.

### FILE C3 - `mark_notification_read_usecase.dart`
Accepts `notificationId`; validates non-empty before calling repository.

### FILE C4 - `mark_all_read_usecase.dart`
Calls `repository.markAllRead()`.

### FILE C5 - `register_device_usecase.dart`
Accepts `DeviceRegistrationEntity`; validates non-empty `fcmToken` and that
`platform` is `android` or `ios`.

### FILE C6 - `unregister_device_usecase.dart`
Accepts `fcmToken`; validates non-empty.

### FILE C7 - `get_notification_preferences_usecase.dart`
Calls `repository.getPreferences()`.

### FILE C8 - `update_notification_preferences_usecase.dart`
Accepts `NotificationPreferencesEntity`; calls `repository.updatePreferences(...)`.

Failure behavior:

- empty `notificationId` / `fcmToken` -> validation failure
- invalid `platform` -> validation failure
- network/server errors -> map through repository implementation

---

## PART D - DATA MODELS

### FILE D1 - `lib/features/notifications/data/models/notification_model.dart`

Parse:

- `id`
- `type`
- `title`
- `body`
- `data`
- `read_at`
- `created_at`

Parsing rules:

- `data` may be absent -> default to `<String, String>{}`. Coerce each value to
  `String` (`value?.toString() ?? ''`); backend sends strings but parse
  defensively.
- `read_at` absent OR null -> `null`.
- Map `type` through `NotificationType.fromRaw`, and keep `rawType`.
- Map to `NotificationEntity`.

---

### FILE D2 - `lib/features/notifications/data/models/notification_preferences_model.dart`

Parse:

- `push_enabled`
- `email_enabled`

Add `toJson()` producing exactly:

```json
{ "push_enabled": <bool>, "email_enabled": <bool> }
```

Map to `NotificationPreferencesEntity`.

---

### FILE D3 - `lib/features/notifications/data/models/register_device_request.dart`

Build the request body:

```json
{ "fcm_token": "<token>", "platform": "<android|ios>" }
```

Provide a `toJson()` and a constructor from `DeviceRegistrationEntity`.

---

## PART E - REMOTE DATASOURCE

### FILE E1 - `lib/features/notifications/data/datasources/notification_remote_datasource.dart`

Create:

```dart
abstract class NotificationRemoteDataSource {
  Future<List<NotificationModel>> getNotifications({int limit});
  Future<int> getUnreadCount();
  Future<void> markRead(String notificationId);
  Future<void> markAllRead();
  Future<void> registerDevice(RegisterDeviceRequest request);
  Future<void> unregisterDevice(String fcmToken);
  Future<NotificationPreferencesModel> getPreferences();
  Future<NotificationPreferencesModel> updatePreferences(
    NotificationPreferencesModel preferences,
  );
}
```

Implementation paths:

```dart
dio.get('/me/notifications', queryParameters: {'limit': limit});
dio.get('/me/notifications/unread-count');
dio.post('/me/notifications/$notificationId/read');
dio.post('/me/notifications/read-all');
dio.post('/me/devices', data: request.toJson());
dio.delete('/me/devices/$fcmToken');
dio.get('/me/notification-preferences');
dio.put('/me/notification-preferences', data: preferences.toJson());
```

Rules:

- Read the list from `response.data['notifications']` (wrapped), not the root.
- Read the count from `response.data['count']`.
- Use the shared Dio client from core providers.
- Let the auth interceptor attach tokens.
- Do not hardcode the full base URL.
- 204 endpoints return no body — do not parse a body for them.

---

## PART F - REPOSITORY IMPLEMENTATION

### FILE F1 - `lib/features/notifications/data/repositories/notification_repository_impl.dart`

Implement `NotificationRepository`.

Pattern:

- call remote datasource
- convert model to entity (or return `unit` for void endpoints)
- catch `ServerException` -> `ServerFailure`
- catch `NetworkException` -> `NetworkFailure`
- catch unknown exceptions -> `UnexpectedFailure`

Keep failure mapping consistent with previous feature repositories.

---

## PART G - FCM / FIREBASE PLATFORM SETUP

The packages are already in `pubspec.yaml`:

```yaml
firebase_core: ^3.3.0
firebase_messaging: ^15.1.0
flutter_local_notifications: ^17.2.2
```

Do not bump versions without asking.

### G1 - Android native config

- Add `android/app/google-services.json` (provided by the project owner / Firebase
  console). Do **not** commit a real one with secrets if the repo policy forbids
  it — confirm with the owner. If absent, document that push will be inert until
  it is added, but the in-app notification center must still work.
- In `android/build.gradle` (project), add the Google services classpath.
- In `android/app/build.gradle`, apply `com.google.gms.google-services` and
  confirm `minSdkVersion` is at least 21 (already required by `SKILL.md`).
- Add the `POST_NOTIFICATIONS` permission to `AndroidManifest.xml` (Android 13+).

### G2 - iOS native config (best-effort, Android-first phase)

- Note in code comments that iOS APNs setup is deferred to the iOS milestone,
  but the Dart code path must not break on iOS. Guard platform-specific calls.

### G3 - Firebase init

- Initialize `Firebase.initializeApp()` in `main.dart` before `runApp`, guarded
  so a missing/invalid config logs a warning (debugPrint) and the app still
  launches. The in-app notification center must not depend on Firebase init
  succeeding.

### FILE G4 - `lib/core/services/push_messaging_service.dart`

Create a thin wrapper around `FirebaseMessaging` and
`flutter_local_notifications`. Responsibilities only (no UI, no business logic):

- `Future<bool> requestPermission()` -> request notification permission
- `Future<String?> getToken()` -> current FCM token (null if unavailable)
- `Stream<String> get onTokenRefresh`
- `Future<void> deleteToken()`
- `Stream<RemoteMessage> get onForegroundMessage`
- `Stream<RemoteMessage> get onMessageOpenedApp`
- `Future<RemoteMessage?> getInitialMessage()` (terminated-state tap)
- local-notification display for foreground messages

Rules:

- This service touches Firebase plugins **only**. It must not import the domain
  or data layers, and must not call the repository.
- Every method must no-op safely if Firebase is not initialized (Android without
  `google-services.json`, iOS without APNs). Return null/empty, log with
  `debugPrint`, never throw out of the service.
- Do not register a background handler that imports heavy app code; keep the
  top-level background handler minimal.

---

## PART H - RIVERPOD DEPENDENCY PROVIDERS

### FILE H1 - `lib/features/notifications/presentation/providers/notification_providers.dart`

Create providers:

- `notificationRemoteDataSourceProvider`
- `notificationRepositoryProvider`
- `getNotificationsUseCaseProvider`
- `getUnreadCountUseCaseProvider`
- `markNotificationReadUseCaseProvider`
- `markAllReadUseCaseProvider`
- `registerDeviceUseCaseProvider`
- `unregisterDeviceUseCaseProvider`
- `getNotificationPreferencesUseCaseProvider`
- `updateNotificationPreferencesUseCaseProvider`
- `pushMessagingServiceProvider` (exposes the `PushMessagingService`)

Rules:

- providers expose use cases / service to presentation
- screens consume feature providers from `SKILL_2.md` and `SKILL_3.md`
- avoid one giant provider that loads everything unconditionally

---

## PART I - DATA-LAYER TESTING NOTES

Unit tests:

- notification model parses with and without `data`
- notification model treats absent/null `read_at` as unread
- unknown `type` maps to `NotificationType.unknown` without throwing
- preferences model round-trips `toJson`/`fromJson`
- register-device request serializes `fcm_token` and `platform`
- repository maps network/server failures
- repository returns `unit` for 204 endpoints

---

## VALIDATION

Run:

```bash
dart format lib/features/notifications lib/core/services
flutter analyze
```

Manual backend run:

```bash
cd ../Memere-backend
make up
make migrate-up
make seed
make run
```

Physical Android:

```bash
adb reverse tcp:8080 tcp:8080
flutter run
```

Manual checks after UI files are built:

- `/me/notifications` loads after login
- `/me/notifications/unread-count` returns a count
- `POST /me/devices` succeeds with a real or stub FCM token
- `/me/notification-preferences` loads and updates

---

## SKILL_1 COMPLETE WHEN

- [ ] Notification domain entities + type enum compile
- [ ] Repository contract exists
- [ ] Use cases exist and validate required IDs/tokens/platform
- [ ] DTO models parse the wrapped list, count, and preferences shapes
- [ ] `data` map and null `read_at` parse safely
- [ ] Unknown notification types fall back without crashing
- [ ] Remote datasource calls the correct endpoints
- [ ] Repository implementation maps failures consistently
- [ ] Firebase init is guarded and non-fatal
- [ ] `PushMessagingService` wraps FCM and no-ops safely when uninitialized
- [ ] Riverpod dependency providers are wired

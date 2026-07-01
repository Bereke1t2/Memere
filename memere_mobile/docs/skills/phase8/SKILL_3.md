# phase8/SKILL_3.md - Push Handling, Deep-Link Routing, Notification Settings & Phase 8 Checklist
# Memere Mobile (memere_mobile) - Phase 8, Part 3
# READ SKILL.md -> phase8/SKILL_1.md -> phase8/SKILL_2.md -> then this file.

---

## OBJECTIVE

Finish Phase 8 by adding:
foreground/background/terminated push handling -> local notification display ->
deep-link routing from a notification tap -> notification preferences (settings)
screen -> refresh integration -> final validation.

By the end of this skill, push notifications wake the app to the right screen,
unread badges update live, and students can control push/email channels.

---

## PHASE 8 FINAL BOUNDARY

Included:

- foreground push -> local notification + badge/list refresh
- background/terminated push tap -> deep-link routing
- deep-link from notification `type` + `data`
- notification preferences screen (push/email toggles)
- refresh of badge/list after relevant app events

Deferred:

- composing/sending notifications from the app
- teacher/admin broadcast UI
- websocket realtime stream
- scheduled local reminders not driven by the backend
- per-event-type granular preferences (backend exposes only push/email)
- rich media / image notifications

---

## PART A - PUSH MESSAGE HANDLING

Use `PushMessagingService` from `SKILL_1.md` Part G. All handling lives in one
place so routing and refresh stay consistent.

### FILE A1 - `lib/features/notifications/presentation/providers/push_handler_provider.dart`

Create a provider/controller initialized once at app start (after Firebase init
and after the router is available) that wires the three delivery states:

1. **Foreground** (`onForegroundMessage`):
   - show a local notification via `flutter_local_notifications`
   - invalidate `notificationListProvider` and `unreadCountProvider`
   - do **not** auto-navigate; the user is already using the app

2. **Background tap** (`onMessageOpenedApp`):
   - parse the `RemoteMessage` data
   - deep-link via the router (Part B)

3. **Terminated tap** (`getInitialMessage` at startup):
   - if a message launched the app, deep-link once after the first frame
   - guard so it only fires once

Rules:

- the background message handler must be a top-level function (FCM requirement);
  keep it minimal — it must not touch Riverpod or navigation.
- never crash on missing/garbage `data`; fall back to opening `/notifications`.
- if Firebase is not initialized, this provider is inert (no streams) and the
  in-app center still works.

---

## PART B - DEEP-LINK ROUTING

### FILE B1 - `lib/features/notifications/presentation/services/notification_router.dart`

Create a pure mapping from a notification to a route, driven by
`NotificationType` + `data`. It must not import Firebase types directly — accept
a normalized `{ type, data }` so it is unit-testable.

Routing table (route only if the target route already exists from earlier
phases; otherwise fall back to `/notifications`):

| type                  | data used        | destination                                  |
|-----------------------|------------------|----------------------------------------------|
| `examGraded`          | `attempt_id`     | exam result/analytics route (Phase 5)        |
| `lessonPublished`     | `course_id`      | course detail (Phase 2) / lesson             |
| `purchaseConfirmed`   | `payment_id`     | purchases/orders (Phase 6) or course detail  |
| `certificateReady`    | `course_id`      | course progress (Phase 7) / profile          |
| `streakWarning`       | —                | `/dashboard` (Phase 7)                        |
| `subscriptionExpired` | `subscription_id`| subscription/purchases (Phase 6)             |
| `videoReady`/`unknown`| —                | `/notifications`                             |

Rules:

- if a required `data` key is missing, fall back to `/notifications` — never
  navigate to a malformed route.
- do not deep-link to teacher-only screens in the student build.
- use the `AppRoutes` path builders from earlier phases (e.g.
  `AppRoutes.courseProgressPath(courseId)`), never raw string concatenation when
  a builder exists.

---

## PART C - NOTIFICATION PREFERENCES PROVIDER

### FILE C1 - `lib/features/notifications/presentation/providers/notification_preferences_provider.dart`

Create:

```dart
final notificationPreferencesProvider =
    FutureProvider<NotificationPreferencesEntity>((ref) async {
  final useCase = ref.watch(getNotificationPreferencesUseCaseProvider);
  final result = await useCase();
  return result.fold((failure) => throw failure, (prefs) => prefs);
});
```

Add an update controller exposing:

```dart
Future<void> setPush(WidgetRef ref, bool enabled);
Future<void> setEmail(WidgetRef ref, bool enabled);
```

Behavior:

- optimistic toggle, then call `updateNotificationPreferencesUseCase`
- on success: keep the new value (the PUT returns the updated prefs)
- on failure: revert the toggle and show a snackbar

---

## PART D - NOTIFICATION SETTINGS SCREEN

### FILE D1 - `lib/features/notifications/presentation/screens/notification_settings_screen.dart`

Create:

```dart
class NotificationSettingsScreen extends ConsumerWidget {
  const NotificationSettingsScreen({super.key});
}
```

Route:

```dart
static const notificationSettings = '/settings/notifications';
```

Required UI states:

- loading
- error with retry
- loaded toggles

Loaded layout:

1. app bar titled `Notification settings`
2. `Push notifications` toggle -> `setPush`
3. `Email notifications` toggle -> `setEmail`
4. helper text explaining what each channel covers (short, factual)
5. if OS-level notification permission is denied, show a banner with an
   `Open settings` action (use the platform settings deep-link if available;
   otherwise explain how to enable)

Rules:

- toggles reflect backend state; do not keep a separate local source of truth
- turning push off in-app updates the backend pref — it does **not** unregister
  the device token (the backend decides delivery from the pref)
- keep copy short and non-marketing

Entry point: an `Account`/profile list row -> `/settings/notifications`.

---

## PART E - SETTINGS / PERMISSION INTERACTION

- App-level `push_enabled` (backend pref) and OS-level permission are different
  things. Show both honestly:
  - OS permission denied + app pref on -> "Enable notifications in system
    settings to receive push."
  - OS permission granted + app pref off -> push suppressed by the user's choice.
- Do not silently flip the backend pref based on OS permission, or vice versa.

---

## PART F - REFRESH INTEGRATION

Invalidate `unreadCountProvider` (and `notificationListProvider` when the center
is open) after:

- a foreground push arrives
- marking read / read-all
- returning from a deep-linked screen
- login (initial load)

Do not poll on a tight timer. Phase 8 relies on:

- fetch on screen open
- pull-to-refresh
- push wake-up

A light refresh on app-resume (`AppLifecycleState.resumed`) for the unread count
is acceptable; do not refresh on every lifecycle tick.

---

## PART G - ERROR HANDLING

Handle:

- unauthenticated -> route guard/login flow (do not call notif endpoints)
- Firebase uninitialized -> push inert, in-app center still works
- permission denied -> settings banner, no crash
- token registration failure -> non-fatal, logged
- malformed push `data` -> fall back to `/notifications`
- backend 5xx on list/count -> error state with retry / hidden badge

Do not expose raw backend errors directly in user-facing text.

---

## PART H - TESTING NOTES

Widget tests:

- settings screen loads toggles from backend
- toggling push optimistically flips and reverts on failure
- permission-denied banner renders when OS permission is denied

Unit tests (notification_router):

- `examGraded` with `attempt_id` -> exam result route
- `lessonPublished` with `course_id` -> course detail route
- missing required `data` key -> `/notifications`
- `unknown` type -> `/notifications`
- teacher-only `videoReady` -> `/notifications` (no teacher route in student build)

Integration/manual tests:

- foreground push shows a local notification and bumps the badge
- background tap opens the correct screen
- terminated-state tap (cold start from notification) opens the correct screen
- preference toggle persists across app restart

---

## VALIDATION

Run:

```bash
dart format lib/core/router lib/features/notifications lib/core/services
flutter analyze
flutter build apk --debug
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

Send a test push (Firebase console or backend trigger) and verify:

- foreground: local notification + badge update
- background: tap routes to the right screen
- terminated: cold-start tap routes to the right screen
- settings toggles read and update `/me/notification-preferences`
- no screen shows an old project name

---

## PHASE 8 FINAL CHECKLIST

### SKILL_1 complete when:

- [ ] Notification entities/models + type enum compile
- [ ] Repository/use cases cover list, count, read, devices, preferences
- [ ] Remote datasource calls all notification endpoints (wrapped list, count)
- [ ] `data` map and null `read_at` parse safely; unknown types fall back
- [ ] Firebase init is guarded and non-fatal
- [ ] `PushMessagingService` wraps FCM and no-ops when uninitialized
- [ ] Failure mapping is consistent

### SKILL_2 complete when:

- [ ] Notifications route exists and is auth-guarded
- [ ] Notification center loads backend data with all UI states
- [ ] Read vs unread tiles render correctly
- [ ] Mark-read / mark-all-read update list + badge
- [ ] Unread bell badge renders in the dashboard app bar
- [ ] Device token registers on login and unregisters on logout
- [ ] Registration failures are non-fatal

### SKILL_3 complete when:

- [ ] Foreground / background / terminated push are all handled
- [ ] Deep-link routing maps type + data to the correct route
- [ ] Malformed/unknown notifications fall back to `/notifications`
- [ ] Notification settings screen reads and updates preferences
- [ ] OS-permission vs app-pref are shown honestly
- [ ] Badge/list refresh on push, read, resume, and login
- [ ] `flutter analyze` has 0 errors
- [ ] Android debug APK builds

---

## PHASE 8 -> PHASE 9 HANDOFF

**Phase 8 is complete when all 3 SKILL files are done and the checklist above passes.**

To start Phase 9, tell Antigravity:

```text
Phase 8 is complete. Read SKILL.md and phase9/SKILL_1.md.
We are starting Phase 9: Teacher flow (content management).
Reference: memere_mobile/docs/memere_Design_Specification.md
```

**What Phase 9 will build:**

- teacher-only content management screens
- course / section / lesson create + edit
- video upload via pre-signed URLs
- quiz / exam authoring
- teacher role gating and route guards

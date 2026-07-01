# phase8/SKILL_2.md - Notification Center Screen, Unread Badge, Routing & Token Lifecycle
# Memere Mobile (memere_mobile) - Phase 8, Part 2
# READ SKILL.md -> phase8/SKILL_1.md -> then this file.

---

## OBJECTIVE

Build the in-app notification experience:
notifications route -> notification list providers -> notification center
screen -> unread badge -> notification tiles -> mark-read interactions ->
device-token lifecycle (register on login, unregister on logout) ->
loading, empty, and error states.

By the end of this skill, a logged-in student can open a notification center,
see read/unread notifications from the backend, mark them read, and the app
registers/refreshes its FCM device token at the right moments.

---

## NOTIFICATION CENTER RULES

- List data comes from `GET /me/notifications`.
- Unread count comes from `GET /me/notifications/unread-count` (source of truth).
- Marking read uses `POST /me/notifications/:id/read` and `read-all`.
- Render the backend `title`/`body` — never fabricate copy.
- Use `type` only to pick an icon and a deep-link target.
- Do not show any old project name. Use `Memere` only for branding.
- Keep the list dense and scannable; no large hero sections.

---

## PART A - ROUTING

### FILE A1 - `lib/core/router/app_router.dart`

Add route:

```dart
static const notifications = '/notifications';
```

Route builder:

- `/notifications` -> `NotificationCenterScreen`

Keep the route behind auth.

Entry points to the notification center:

- a bell icon with an unread badge in the dashboard/home app bar (Phase 7)
- optionally an `Account`/profile list row

If the home shell has a bottom nav (Phase 7 order Dashboard / Courses / Exams /
Account), do **not** add a 5th tab for notifications — surface it as the app-bar
bell so the tab bar stays at four.

---

## PART B - NOTIFICATION LIST PROVIDERS

### FILE B1 - `lib/features/notifications/presentation/providers/notification_list_provider.dart`

Create:

```dart
final notificationListProvider =
    FutureProvider<List<NotificationEntity>>((ref) async {
  final useCase = ref.watch(getNotificationsUseCaseProvider);
  final result = await useCase();
  return result.fold((failure) => throw failure, (items) => items);
});
```

Refresh by `ref.invalidate(notificationListProvider)` after:

- marking one/all read
- a foreground push arrives (`SKILL_3.md`)
- returning to the screen after a deep-link tap

---

### FILE B2 - `lib/features/notifications/presentation/providers/unread_count_provider.dart`

Create:

```dart
final unreadCountProvider = FutureProvider<int>((ref) async {
  final useCase = ref.watch(getUnreadCountUseCaseProvider);
  final result = await useCase();
  return result.fold((failure) => throw failure, (count) => count);
});
```

Rules:

- The app-bar badge watches this provider.
- Invalidate it whenever read-state changes or a push arrives.
- On failure, the badge hides (treat error as "no badge"), never blocks the UI.

---

### FILE B3 - `lib/features/notifications/presentation/providers/notification_actions_provider.dart`

Optional controller (AsyncNotifier or plain provider) exposing:

```dart
Future<void> markRead(WidgetRef ref, String id);
Future<void> markAllRead(WidgetRef ref);
```

Behavior:

- optimistic: flip the tile to read immediately, then call the use case
- on success: invalidate `notificationListProvider` and `unreadCountProvider`
- on failure: revert optimistic state and surface a snackbar

Use whatever controller pattern the project already uses; do not invent a new
state-management style.

---

## PART C - NOTIFICATION CENTER SCREEN

### FILE C1 - `lib/features/notifications/presentation/screens/notification_center_screen.dart`

Create:

```dart
class NotificationCenterScreen extends ConsumerWidget {
  const NotificationCenterScreen({super.key});
}
```

Required UI states:

- loading skeleton
- empty state ("You're all caught up")
- error state with retry
- loaded list

Loaded layout:

1. app bar titled `Notifications`
2. `Mark all read` action in the app bar (disabled when unread count is 0)
3. pull-to-refresh list of notification tiles
4. unread tiles visually distinct from read tiles

Screen behavior:

- pull to refresh invalidates list + unread count providers
- tapping a tile marks it read, then deep-links per `type`/`data`
- `Mark all read` calls the read-all action

---

## PART D - NOTIFICATION TILE

### FILE D1 - `lib/features/notifications/presentation/widgets/notification_tile.dart`

Props:

- `NotificationEntity notification`
- `VoidCallback onTap`

Show:

- leading icon chosen by `NotificationType` (default icon for `unknown`)
- `title` (semibold)
- `body` (secondary, max 2 lines, ellipsis)
- `timeAgoLabel` from `createdAt`
- unread indicator dot when `isUnread`

Rules:

- read tiles use muted/secondary text; unread tiles use full-contrast text and
  the unread dot
- stable row height so refreshes don't make rows jump
- never render raw `data` map to the user

---

### FILE D2 - `lib/features/notifications/presentation/widgets/notification_icon.dart`

Map `NotificationType` -> icon + accent color:

- `examGraded` -> results icon
- `lessonPublished` -> play/lesson icon
- `purchaseConfirmed` -> receipt icon
- `certificateReady` -> award icon
- `streakWarning` -> flame icon (use warning accent)
- `subscriptionExpired` -> alert icon
- `videoReady` / `unknown` -> generic bell icon

Use the Phase design tokens (`AppColors`) — no hardcoded hex.

---

## PART E - UNREAD BADGE

### FILE E1 - `lib/features/notifications/presentation/widgets/notification_bell.dart`

Props:

- `VoidCallback onTap`

Behavior:

- watches `unreadCountProvider`
- shows a bell icon with a small count badge when count > 0
- shows `99+` when count > 99
- hides the badge on loading/error/zero
- tap routes to `/notifications`

Use stable layout so the badge appearing/disappearing does not shift the app bar.

Integrate the bell into:

- Phase 7 dashboard/home app bar (primary)
- profile/account header if one exists

---

## PART F - SKELETON, EMPTY & ERROR STATES

### FILE F1 - `notification_list_skeleton.dart`
Shimmer/skeleton of ~6 tile rows using the existing skeleton pattern.

### FILE F2 - `notification_empty_state.dart`
Shown when the list loads with zero notifications:

```text
You're all caught up.
```

Optional secondary line: `New notifications will appear here.`

### FILE F3 - `notification_error_state.dart`
Props: `Object error`, `VoidCallback onRetry`. Short message + retry button.
Never expose raw stack traces.

---

## PART G - DEVICE TOKEN LIFECYCLE

This is the bridge between `PushMessagingService` (`SKILL_1.md` Part G) and the
backend `/me/devices` endpoints. Token registration must be tied to auth state.

### FILE G1 - `lib/features/notifications/presentation/providers/device_token_provider.dart`

Create a controller that:

1. on successful login / app-start-while-authenticated:
   - request notification permission (once)
   - read the FCM token via `PushMessagingService.getToken()`
   - call `registerDeviceUseCase` with `{ fcmToken, platform }`
2. listens to `PushMessagingService.onTokenRefresh` and re-registers
3. on logout:
   - call `unregisterDeviceUseCase(currentToken)` **before** tokens are cleared
   - optionally `deleteToken()` on the FCM side

Rules:

- platform string is `android` or `ios` — derive from `Platform.isAndroid` /
  `Platform.isIOS`, guarded for web/tests.
- registration failure is non-fatal: log and continue. The user can still use
  the app and the in-app center without a registered device.
- never block login/logout on device registration.
- do not register a token while unauthenticated.

### G2 - Wire into auth flow

- after the Phase 1 login success path, trigger device registration
- in the logout path, trigger unregister **before** clearing secure storage
  tokens, so the authenticated `DELETE /me/devices/:token` call still succeeds
- on app launch while already authenticated (Phase 1 splash/route guard),
  trigger registration once

Do not duplicate registration on every rebuild — guard with a one-shot flag or a
provider that runs once per session.

---

## PART H - HOME / DASHBOARD INTEGRATION

- Add the `NotificationBell` to the dashboard app bar (Phase 7 `DashboardHeader`
  or shell app bar).
- Invalidate `unreadCountProvider` when:
  - the dashboard refreshes
  - returning from the notification center
  - a foreground push arrives (`SKILL_3.md`)

Keep catalog / exams / account reachable exactly as Phase 7 left them.

---

## PART I - TESTING NOTES

Widget tests:

- notification list loading skeleton renders
- empty list shows the caught-up state
- unread tile shows the unread dot; read tile does not
- tapping a tile triggers mark-read
- bell shows a badge for count > 0 and `99+` above 99
- bell hides the badge on zero/error

Provider tests:

- list provider returns entities on success, throws failure on error
- unread count provider returns count, hides badge on failure
- mark-read invalidates list + unread providers
- device-token controller registers on login, unregisters on logout
- device registration failure does not throw out of the auth flow

---

## VALIDATION

Run:

```bash
dart format lib/core/router lib/features/notifications
flutter analyze
flutter build apk --debug
```

Manual checks:

- bell with unread badge appears on the dashboard
- tapping the bell opens the notification center
- list loads backend notifications after login
- pull-to-refresh reloads list + badge
- tapping a notification marks it read and updates the badge
- `Mark all read` clears the badge
- logging in registers a device token (check backend `/me/devices`)
- logging out unregisters the device token
- no screen shows an old project name

---

## SKILL_2 COMPLETE WHEN

- [ ] Notifications route exists and is auth-guarded
- [ ] List + unread-count providers load backend data
- [ ] Notification center screen has loading/empty/error/loaded states
- [ ] Notification tiles distinguish read vs unread
- [ ] Mark-read and mark-all-read update list + badge
- [ ] Unread bell badge renders in the dashboard app bar
- [ ] Device token registers on login and app-start-while-authenticated
- [ ] Device token re-registers on FCM token refresh
- [ ] Device token unregisters on logout before tokens are cleared
- [ ] Registration failures are non-fatal

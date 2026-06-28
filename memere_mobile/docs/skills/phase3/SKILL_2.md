# phase3/SKILL_2.md — HLS Player Screen, Controls & Progress Saving
# Memere Mobile (memere_mobile) — Phase 3, Part 2
# READ SKILL.md → phase3/SKILL_1.md → then this file.

---

## OBJECTIVE

Build the student video playback experience:
course lesson tap → video player route → signed HLS stream fetch →
`video_player` + `chewie` playback → progress autosave → completion handoff.

By the end of this skill, a user can play a ready lesson video when a valid
`videoId` is available.

---

## CRITICAL PRECONDITION

Phase 3 playback needs a `video_id`.

If Phase 2 lesson DTOs do not include `video_id`, do **not** fake one. Add a
clear placeholder state:

```text
This lesson does not have a video attached yet.
```

Then either:

- update backend course/lesson DTOs to include `video_id`, or
- add a backend endpoint to resolve video by lesson ID.

The mobile app must not guess or hardcode video IDs.

---

## DESIGN REQUIREMENTS

Use a functional app-screen layout:

- black video surface at top
- compact lesson/course metadata below
- lesson list context or next lesson control
- download action as an icon button
- no marketing hero
- no decorative cards around the video player

The video player itself may be full width and fixed aspect ratio.

---

## PART A — ROUTING UPDATE

### FILE A1 — `lib/core/router/app_router.dart`

Add route:

```dart
static const videoPlayer = '/videos/:videoId';

static String videoPlayerPath({
  required String videoId,
  required String lessonId,
  required String courseId,
}) {
  return '/videos/$videoId?lessonId=$lessonId&courseId=$courseId';
}
```

Route builder:

```dart
GoRoute(
  path: AppRoutes.videoPlayer,
  builder: (_, state) {
    final videoId = state.pathParameters['videoId']!;
    final lessonId = state.uri.queryParameters['lessonId'] ?? '';
    final courseId = state.uri.queryParameters['courseId'] ?? '';
    return VideoPlayerScreen(
      videoId: videoId,
      lessonId: lessonId,
      courseId: courseId,
    );
  },
),
```

Keep route behind auth.

---

## PART B — PLAYER STATE PROVIDER

### FILE B1 — `lib/features/video_player/presentation/providers/video_player_controller_provider.dart`

Create an `AsyncNotifierProvider.family` or `StateNotifierProvider.family` that
owns the player lifecycle.

State shape:

```dart
class VideoPlaybackState {
  const VideoPlaybackState({
    this.stream,
    this.isInitializing = false,
    this.isSavingProgress = false,
    this.lastSavedPositionSeconds = 0,
    this.errorMessage,
  });
}
```

Provider responsibility:

- call `GetVideoStreamUseCase`
- initialize `VideoPlayerController.networkUrl(Uri.parse(masterUrl))`
- create `ChewieController`
- expose both controllers to screen/widgets
- dispose both controllers
- save progress every 30 seconds of watch movement
- save progress on pause
- save progress on dispose
- mark lesson complete only through backend progress response or explicit use case

Do not store controllers in domain/data layers.

### Save cadence

Use:

```dart
static const progressSaveIntervalSeconds = 30;
```

Only call backend when:

- `lessonId` is not empty
- current position is at least 30 seconds beyond last saved position, or
- playback is pausing/disposed and position changed

---

## PART C — VIDEO PLAYER SCREEN

### FILE C1 — `lib/features/video_player/presentation/screens/video_player_screen.dart`

Constructor:

```dart
class VideoPlayerScreen extends ConsumerWidget {
  const VideoPlayerScreen({
    super.key,
    required this.videoId,
    required this.lessonId,
    required this.courseId,
  });

  final String videoId;
  final String lessonId;
  final String courseId;
}
```

Required UI states:

- loading stream/player
- video processing/not ready
- playback ready
- stream error
- missing video ID/lesson ID

Layout:

1. black scaffold background
2. top app bar overlay/back button
3. player area with 16:9 aspect ratio
4. lesson metadata below
5. download icon button
6. progress save indicator only when useful

If stream URL expires during playback, show retry action that invalidates the
provider and requests a fresh stream URL.

---

## PART D — PLAYER WIDGETS

### FILE D1 — `lib/features/video_player/presentation/widgets/hls_video_player.dart`

Props:

- `VideoPlayerController controller`
- `ChewieController chewieController`

Render:

```dart
AspectRatio(
  aspectRatio: controller.value.aspectRatio == 0
      ? 16 / 9
      : controller.value.aspectRatio,
  child: Chewie(controller: chewieController),
)
```

Handle uninitialized state gracefully.

---

### FILE D2 — `lib/features/video_player/presentation/widgets/video_player_error_state.dart`

Props:

- `String title`
- `String message`
- optional retry callback
- optional back callback

Use for:

- missing video
- failed stream fetch
- not ready/processing
- access denied

---

### FILE D3 — `lib/features/video_player/presentation/widgets/video_action_bar.dart`

Actions:

- download
- refresh stream
- mark complete if appropriate

Use icons, not text-only command buttons.

Download action calls provider from `SKILL_3.md`.

---

## PART E — COURSE DETAIL LESSON TAP INTEGRATION

### FILE E1 — `lib/features/courses/presentation/widgets/lesson_tile.dart`

Update lesson tap behavior from Phase 2 placeholder.

If lesson has `videoId`:

```dart
context.go(
  AppRoutes.videoPlayerPath(
    videoId: lesson.videoId!,
    lessonId: lesson.id,
    courseId: lesson.courseId,
  ),
);
```

If lesson has no `videoId`:

```dart
ScaffoldMessenger.of(context).showSnackBar(
  const SnackBar(content: Text('This lesson does not have a video attached yet.')),
);
```

If lesson is locked and the user is not enrolled, keep Phase 2 enrollment
placeholder until Phase 6.

---

## PART F — STATUS BEFORE STREAM

Before requesting stream, optionally call `GET /videos/:id/status`.

Rules:

- `ready` → request stream
- `pending` or `processing` → show processing state
- `failed` → show failed state

If you skip status and call stream directly, handle backend errors cleanly.

Preferred implementation:

1. status
2. stream only when ready
3. initialize player

---

## PART G — PLATFORM REQUIREMENTS

Android must have internet permission from Phase 1.

For cleartext local backend:

```xml
android:usesCleartextTraffic="true"
```

Do not add platform-specific video hacks unless the player fails on a real
device and the failure is confirmed.

---

## VALIDATION

Run:

```bash
dart format lib/core/router lib/features/video_player lib/features/courses
flutter analyze
flutter build apk --debug
```

Manual test on physical Android:

```bash
adb reverse tcp:8080 tcp:8080
flutter run
```

Use a backend course with a ready video record. If no ready video exists, verify
the missing/not-ready states instead of hardcoding IDs.

---

## SKILL_2 CHECKLIST

- [ ] `/videos/:videoId` route exists
- [ ] Player screen receives `videoId`, `lessonId`, and `courseId`
- [ ] Player fetches signed stream URL from backend
- [ ] `video_player` + `chewie` initialize and dispose correctly
- [ ] Loading/error/not-ready states render
- [ ] Progress saves every 30 seconds
- [ ] Progress saves on pause/dispose
- [ ] Course lesson tap navigates only when `videoId` exists
- [ ] Missing `videoId` is handled honestly
- [ ] `flutter analyze` has 0 errors
- [ ] Android debug APK builds

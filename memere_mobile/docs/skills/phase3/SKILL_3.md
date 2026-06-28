# phase3/SKILL_3.md — Offline Download Foundation, Progress Hooks & Phase 3 Checklist
# Memere Mobile (memere_mobile) — Phase 3, Part 3
# READ SKILL.md → phase3/SKILL_1.md → phase3/SKILL_2.md → then this file.

---

## OBJECTIVE

Build the offline video foundation and final Phase 3 polish:
signed download URL request → download queue state → local metadata storage →
expiry handling → progress hooks → final validation.

By the end of this skill, the app should support online HLS playback and have a
safe offline-download foundation that can be expanded without breaking security
rules.

---

## PHASE 3 OFFLINE BOUNDARY

Allowed in Phase 3:

- request signed download URL
- represent download queue state
- persist downloaded-video metadata
- show download progress UI
- expire local downloads after 30 days
- optionally download a single manifest/file as a proof of plumbing

Not required in Phase 3:

- full HLS segment parser/downloader
- AES-256 encryption implementation
- custom local HLS playback from rewritten manifests
- background worker retries across app restarts

Reason: secure offline HLS requires careful manifest parsing, segment download,
path rewriting, encryption, and expiry. Phase 3 should lay the architecture and
avoid pretending a partial file is a complete offline video.

---

## SECURITY RULES

Never store or log:

- access token
- refresh token
- raw S3 object keys
- signed URLs beyond their useful expiry

Downloaded content:

- must be tied to authenticated user/device
- must expire after 30 days
- should be encrypted before production release

Do not bypass backend access checks. Always request stream/download URLs from
the authenticated API.

---

## PART A — OFFLINE DOMAIN ENTITIES

### FILE A1 — `lib/features/video_player/domain/entities/offline_video_entity.dart`

Fields:

- `videoId`
- `lessonId`
- `courseId`
- `title`
- `localPath`
- `downloadedAt`
- `expiresAt`
- `fileSizeBytes`
- `status`

Use enum:

```dart
enum OfflineVideoStatus { queued, downloading, downloaded, failed, expired }
```

Helpers:

- `bool get isExpired`
- `bool get isPlayable`

---

### FILE A2 — `lib/features/video_player/domain/entities/download_task_entity.dart`

Fields:

- `videoId`
- `lessonId`
- `courseId`
- `title`
- `progress`
- `status`
- `errorMessage`

`progress` is `double` from `0.0` to `1.0`.

---

## PART B — LOCAL STORAGE CONTRACT

### FILE B1 — `lib/features/video_player/data/datasources/offline_video_local_datasource.dart`

Use Hive or SharedPreferences metadata. Prefer Hive because the project already
uses Hive for cache/offline features.

Required methods:

```dart
Future<List<OfflineVideoModel>> getDownloads();
Future<OfflineVideoModel?> getDownload(String videoId);
Future<void> saveDownload(OfflineVideoModel video);
Future<void> removeDownload(String videoId);
Future<void> clearExpiredDownloads();
```

Do not store signed URLs permanently.

Store only:

- IDs
- local file/manifest path
- title
- downloaded timestamp
- expiry timestamp
- status
- file size

---

## PART C — OFFLINE MODELS

### FILE C1 — `lib/features/video_player/data/models/offline_video_model.dart`

`OfflineVideoModel extends OfflineVideoEntity`.

Implement:

- `factory OfflineVideoModel.fromJson(Map<String, dynamic> json)`
- `Map<String, dynamic> toJson()`
- status parsing fallback to `OfflineVideoStatus.failed`

---

### FILE C2 — `lib/features/video_player/data/models/download_task_model.dart`

`DownloadTaskModel extends DownloadTaskEntity`.

This can be in-memory only. It does not need JSON unless persisted queue retries
are implemented.

---

## PART D — DOWNLOAD REPOSITORY

### FILE D1 — `lib/features/video_player/domain/repositories/offline_video_repository.dart`

Required methods:

```dart
Future<Either<Failure, List<OfflineVideoEntity>>> getDownloads();

Future<Either<Failure, OfflineVideoEntity?>> getDownload(String videoId);

Future<Either<Failure, OfflineVideoEntity>> startDownload({
  required String videoId,
  required String lessonId,
  required String courseId,
  required String title,
});

Future<Either<Failure, void>> removeDownload(String videoId);

Future<Either<Failure, void>> clearExpiredDownloads();
```

---

### FILE D2 — `lib/features/video_player/data/repositories/offline_video_repository_impl.dart`

Responsibilities:

- call `VideoRepository.getDownloadUrl(videoId)`
- create queue/downloading state
- write downloaded metadata to local datasource
- clear expired local metadata
- map errors to `Failure`

If implementing actual file download in Phase 3, use Dio download APIs with a
short-lived signed URL. Keep the implementation cancellable and scoped.

Do not download with raw backend auth headers against CDN signed URLs unless the
signed URL explicitly requires no bearer token.

---

## PART E — DOWNLOAD PROVIDERS

### FILE E1 — `lib/features/video_player/presentation/providers/offline_video_provider.dart`

Create:

- `offlineVideoRepositoryProvider`
- `offlineDownloadsProvider`
- `downloadTaskProvider`

Required notifier methods:

- `loadDownloads()`
- `startDownload(videoId, lessonId, courseId, title)`
- `removeDownload(videoId)`
- `clearExpiredDownloads()`

UI should reflect:

- queued
- downloading with progress
- downloaded
- failed
- expired

---

## PART F — DOWNLOAD UI

### FILE F1 — `lib/features/video_player/presentation/widgets/download_button.dart`

Props:

- `videoId`
- `lessonId`
- `courseId`
- `title`

States:

- not downloaded: download icon
- queued/downloading: circular progress
- downloaded: check/downloaded icon
- failed: retry icon
- expired: refresh icon

Use tooltip/semantic label.

---

### FILE F2 — integrate with `video_action_bar.dart`

Replace placeholder download action with `DownloadButton`.

If `videoId` is missing, disable the button.

---

## PART G — PROGRESS COMPLETION HOOKS

Backend behavior:

- `PUT /lessons/:id/video-progress` stores monotonic video progress.
- Backend auto-completes when watched progress reaches its threshold.
- `POST /lessons/:id/complete` is available for explicit completion.

Mobile rules:

- Save progress during playback.
- Trust backend `is_completed` response.
- Do not mark complete locally without backend confirmation.
- If save progress fails, keep playback usable and retry on next interval.

Player provider should expose a small status:

- last saved position
- save pending/error
- completed flag

Do not block playback on progress save failure.

---

## PART H — OPTIONAL DOWNLOAD IMPLEMENTATION DETAIL

If implementing a first real download pass:

1. Request `/videos/:id/download-url`.
2. Download manifest file to app documents directory.
3. Persist metadata with status `downloaded`.
4. Mark clearly in code/comments that full HLS segment download is not complete.

Do not present this as fully playable offline HLS unless segment download and
manifest rewriting are implemented.

Production-ready offline HLS requires:

- parse master manifest
- choose variant or download all variants
- download media playlist
- download each `.ts`/segment file
- rewrite segment URLs to local paths
- encrypt files
- play local manifest
- clear expired files

---

## PART I — ROUTE AND SCREEN POLISH

Before final validation:

- Course detail lesson rows should show video availability.
- Player screen should show lesson title when passed or resolvable.
- Back navigation should return to course detail.
- Expired stream URL should have retry.
- App should recover gracefully when backend has no video service configured.

---

## VALIDATION

Run:

```bash
dart format lib/features/video_player lib/features/courses lib/core/router
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

Manual checks:

- Open course detail.
- Tap lesson with a valid video ID.
- Player requests `/videos/:id/status`.
- Ready video requests `/videos/:id/stream`.
- HLS plays.
- Progress saves to `/lessons/:id/video-progress`.
- Download button requests `/videos/:id/download-url`.
- Missing/no-ready-video states are clear.

---

## PHASE 3 FINAL CHECKLIST

### SKILL_1 complete when:

- [ ] Video/progress entities compile
- [ ] Video models parse backend JSON
- [ ] Video repository uses backend endpoints
- [ ] Stream/download/progress use cases validate inputs
- [ ] Dependency providers are wired

### SKILL_2 complete when:

- [ ] `/videos/:videoId` route exists
- [ ] Video player screen initializes from signed stream URL
- [ ] Chewie/video_player controllers dispose correctly
- [ ] Loading, error, processing, and failed states render
- [ ] Progress saves every 30 seconds and on pause/dispose
- [ ] Course lesson tap handles missing `videoId` honestly

### SKILL_3 complete when:

- [ ] Download URL flow is wired
- [ ] Offline metadata is stored locally
- [ ] Download UI reflects queued/downloading/downloaded/failed/expired
- [ ] Expired downloads are cleared
- [ ] No signed URLs are persisted beyond active download
- [ ] `flutter analyze` has 0 errors
- [ ] Android debug APK builds

---

## PHASE 3 → PHASE 4 HANDOFF

**Phase 3 is complete when all 3 SKILL files are done and the checklist above passes.**

To start Phase 4, tell Antigravity:

```text
Phase 3 is complete. Read SKILL.md and all phase4 skill files.
We are starting Phase 4: Quiz Engine.
Reference: memere_mobile/docs/memere_Design_Specification.md
```

**What Phase 4 will build:**

- quiz list/detail access from lessons
- quiz attempt start
- question flow UI
- answer selection
- autosave/submission
- server-graded result screen

# phase3/SKILL_1.md — Video Domain, Data Layer & Backend Contract
# Memere Mobile (memere_mobile) — Phase 3, Part 1
# READ SKILL.md → all phase1 skill files → all phase2 skill files → then this file.

---

## OBJECTIVE

Build the Video Learning foundation:
video stream/status/download entities → progress entities → repository contracts →
use cases → backend DTO models → remote datasources → repository implementations →
Riverpod dependency providers.

By the end of this skill, the app can request signed HLS stream URLs, check video
status, request download URLs, and save lesson video progress through the local
backend. UI/player work comes in `SKILL_2.md`.

---

## PHASE 3 BOUNDARY

Phase 3 builds **student video playback + offline download foundation only**:

- HLS stream URL fetch
- video player screen
- playback controls
- periodic progress save
- lesson completion at backend threshold
- signed download URL flow
- offline download queue foundation

Do **not** build these yet:

- quiz/exam flows
- payment/enrollment WebView
- teacher video upload screens
- progress dashboard
- push notifications
- AI tutor

Teacher upload routes exist in the backend but are Phase 9 for mobile.

---

## BACKEND CONTRACT

Use the existing local backend API:

```text
GET /api/v1/videos/:id/status
GET /api/v1/videos/:id/stream
GET /api/v1/videos/:id/download-url
GET /api/v1/videos/download/:token
PUT /api/v1/lessons/:id/video-progress
POST /api/v1/lessons/:id/complete
```

Important:

- `:id` for video routes is `video_id`, not `lesson_id`.
- `PUT /lessons/:id/video-progress` uses `lesson_id`.
- Phase 2 `LessonEntity` does not currently expose `video_id`.
- Add a nullable `videoId` to lesson models only if the backend course/detail DTO
  includes it. If not, Phase 3 needs a backend change or a new lesson-video lookup
  endpoint before real playback can start from a lesson tap.

---

## RESPONSE SHAPES

### `GET /videos/:id/status`

```json
{
  "video_id": "uuid",
  "status": "ready",
  "duration_seconds": 900,
  "has_thumbnail": true,
  "error": "optional teacher/admin-only processing error"
}
```

### `GET /videos/:id/stream`

```json
{
  "master_url": "https://signed-cdn-url/master.m3u8",
  "expires_in": 7200,
  "thumbnail_url": "https://...",
  "duration_seconds": 900
}
```

### `GET /videos/:id/download-url`

```json
{
  "download_url": "https://signed-cdn-url/master.m3u8",
  "token": "single-use-token",
  "expires_in": 7200
}
```

### `PUT /lessons/:id/video-progress`

Request:

```json
{
  "position_seconds": 420
}
```

Response:

```json
{
  "lesson_id": "uuid",
  "is_completed": false,
  "completed_at": "2026-01-01T00:00:00Z",
  "video_progress_seconds": 420
}
```

---

## REQUIRED FOLDER STRUCTURE

Create or verify:

```bash
mkdir -p lib/features/video_player/data/datasources
mkdir -p lib/features/video_player/data/models
mkdir -p lib/features/video_player/data/repositories
mkdir -p lib/features/video_player/domain/entities
mkdir -p lib/features/video_player/domain/repositories
mkdir -p lib/features/video_player/domain/usecases
mkdir -p lib/features/video_player/presentation/providers
mkdir -p lib/features/video_player/presentation/screens
mkdir -p lib/features/video_player/presentation/widgets
```

---

## PART A — DOMAIN ENTITIES

### FILE A1 — `lib/features/video_player/domain/entities/video_status_entity.dart`

Fields:

- `videoId`
- `status`
- `durationSeconds`
- `hasThumbnail`
- `error`

Use enum:

```dart
enum VideoProcessingStatus { pending, processing, ready, failed }
```

Add helper:

- `bool get isReady`
- `bool get isProcessing`
- `bool get isFailed`

---

### FILE A2 — `lib/features/video_player/domain/entities/video_stream_entity.dart`

Fields:

- `masterUrl`
- `expiresIn`
- `thumbnailUrl`
- `durationSeconds`

Add helper:

- `DateTime get expiresAt`

---

### FILE A3 — `lib/features/video_player/domain/entities/video_download_entity.dart`

Fields:

- `downloadUrl`
- `token`
- `expiresIn`

Add helper:

- `DateTime get expiresAt`

---

### FILE A4 — `lib/features/video_player/domain/entities/lesson_progress_entity.dart`

Fields:

- `lessonId`
- `isCompleted`
- `completedAt`
- `videoProgressSeconds`

---

## PART B — REPOSITORY CONTRACT

### FILE B1 — `lib/features/video_player/domain/repositories/video_repository.dart`

Use `Either<Failure, T>`.

Required methods:

```dart
Future<Either<Failure, VideoStatusEntity>> getVideoStatus(String videoId);

Future<Either<Failure, VideoStreamEntity>> getStream(String videoId);

Future<Either<Failure, VideoDownloadEntity>> getDownloadUrl(String videoId);

Future<Either<Failure, LessonProgressEntity>> saveVideoProgress({
  required String lessonId,
  required int positionSeconds,
});

Future<Either<Failure, void>> markLessonComplete(String lessonId);
```

Do not expose Dio, JSON, Chewie, or `VideoPlayerController` from the domain layer.

---

## PART C — USE CASES

Create one file per use case:

```text
lib/features/video_player/domain/usecases/get_video_status_usecase.dart
lib/features/video_player/domain/usecases/get_video_stream_usecase.dart
lib/features/video_player/domain/usecases/get_video_download_url_usecase.dart
lib/features/video_player/domain/usecases/save_video_progress_usecase.dart
lib/features/video_player/domain/usecases/mark_lesson_complete_usecase.dart
```

Rules:

- validate IDs are not empty
- validate `positionSeconds >= 0`
- return `ValidationFailure` for invalid input
- do not handle player state inside use cases

---

## PART D — DATA MODELS

### FILE D1 — `lib/features/video_player/data/models/video_status_model.dart`

`VideoStatusModel extends VideoStatusEntity`.

Map:

| Dart | JSON |
|------|------|
| `videoId` | `video_id` |
| `status` | `status` |
| `durationSeconds` | `duration_seconds` |
| `hasThumbnail` | `has_thumbnail` |
| `error` | `error` |

Fallback unknown status to `VideoProcessingStatus.pending`.

---

### FILE D2 — `lib/features/video_player/data/models/video_stream_model.dart`

`VideoStreamModel extends VideoStreamEntity`.

Map:

| Dart | JSON |
|------|------|
| `masterUrl` | `master_url` |
| `expiresIn` | `expires_in` |
| `thumbnailUrl` | `thumbnail_url` |
| `durationSeconds` | `duration_seconds` |

---

### FILE D3 — `lib/features/video_player/data/models/video_download_model.dart`

`VideoDownloadModel extends VideoDownloadEntity`.

Map:

| Dart | JSON |
|------|------|
| `downloadUrl` | `download_url` |
| `token` | `token` |
| `expiresIn` | `expires_in` |

---

### FILE D4 — `lib/features/video_player/data/models/lesson_progress_model.dart`

`LessonProgressModel extends LessonProgressEntity`.

Map:

| Dart | JSON |
|------|------|
| `lessonId` | `lesson_id` |
| `isCompleted` | `is_completed` |
| `completedAt` | `completed_at` |
| `videoProgressSeconds` | `video_progress_seconds` |

---

## PART E — REMOTE DATASOURCE

### FILE E1 — `lib/features/video_player/data/datasources/video_remote_datasource.dart`

Use `DioClient` only.

Required methods:

```dart
Future<VideoStatusModel> getVideoStatus(String videoId);
Future<VideoStreamModel> getStream(String videoId);
Future<VideoDownloadModel> getDownloadUrl(String videoId);
Future<LessonProgressModel> saveVideoProgress({
  required String lessonId,
  required int positionSeconds,
});
Future<void> markLessonComplete(String lessonId);
```

Endpoint examples:

```dart
await _client.get<Map<String, dynamic>>('/videos/$videoId/status');
await _client.get<Map<String, dynamic>>('/videos/$videoId/stream');
await _client.get<Map<String, dynamic>>('/videos/$videoId/download-url');
await _client.put<Map<String, dynamic>>(
  '/lessons/$lessonId/video-progress',
  data: {'position_seconds': positionSeconds},
);
await _client.post('/lessons/$lessonId/complete');
```

Throw `FormatException` for missing response bodies where a body is expected.

---

## PART F — REPOSITORY IMPLEMENTATION

### FILE F1 — `lib/features/video_player/data/repositories/video_repository_impl.dart`

Map errors consistently:

- `DioException` → `ServerFailure.fromDioError`
- any other exception → `UnknownFailure`

Do not instantiate player controllers in repositories.

---

## PART G — PROVIDERS

### FILE G1 — `lib/features/video_player/presentation/providers/video_providers.dart`

Create dependency providers:

- `videoRemoteDataSourceProvider`
- `videoRepositoryProvider`
- `getVideoStatusUseCaseProvider`
- `getVideoStreamUseCaseProvider`
- `getVideoDownloadUrlUseCaseProvider`
- `saveVideoProgressUseCaseProvider`
- `markLessonCompleteUseCaseProvider`

Mirror the auth/courses provider style.

---

## VALIDATION

After this skill:

```bash
dart format lib/features/video_player
flutter analyze
```

Do not continue to player UI until this passes.

---

## SKILL_1 CHECKLIST

- [ ] Video entities compile
- [ ] Lesson progress entity compiles
- [ ] Repository contract has status/stream/download/progress methods
- [ ] Use cases validate IDs and progress values
- [ ] Models parse backend JSON exactly
- [ ] Remote datasource calls video/progress endpoints
- [ ] Repository maps failures consistently
- [ ] Dependency providers are wired
- [ ] `flutter analyze` has 0 errors

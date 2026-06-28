# phase7/SKILL_1.md - Progress Dashboard Domain, Data Layer & Backend Contract
# Memere Mobile (memere_mobile) - Phase 7, Part 1
# READ SKILL.md -> all phase1 skill files -> all phase2 skill files -> all phase3 skill files -> all phase4 skill files -> all phase5 skill files -> all phase6 skill files -> then this file.

---

## OBJECTIVE

Build the Progress Dashboard foundation:
lesson progress entity -> course progress entity -> dashboard entities ->
streak entity -> trend entity -> repository contract -> use cases ->
backend DTO models -> remote datasource -> repository implementation ->
Riverpod dependency providers.

By the end of this skill, the app can load the student's dashboard, enrolled
course progress, study streak, and subject trend data from the local Memere
backend. Screens and widgets are built in `SKILL_2.md` and `SKILL_3.md`.

---

## PHASE 7 BOUNDARY

Phase 7 builds **student progress visibility only**:

- student dashboard data
- enrolled course progress summaries
- per-course lesson progress
- current and longest study streak
- completion percentages
- recent/continue learning data derived from enrolled course progress
- subject trend summary from mock-exam analytics
- weak-area summary only from already available exam analytics data

Do **not** build these yet:

- push notifications
- in-app notification center
- certificates
- adaptive learning engine
- AI recommendations
- teacher/admin analytics dashboards
- new backend analytics endpoints

---

## BACKEND CONTRACT

Use the existing authenticated backend API:

```text
GET /api/v1/me/dashboard
GET /api/v1/me/streak
GET /api/v1/courses/:id/progress
GET /api/v1/me/trend?subject=<subject>
```

Related endpoints already used by earlier phases:

```text
POST /api/v1/lessons/:id/complete
PUT /api/v1/lessons/:id/video-progress
GET /api/v1/exam-attempts/:id/analytics
```

Do not create client-only progress totals. Dashboard totals must come from the
backend response. Lesson completion and video progress are written by Phase 3
learning flow, then read here.

---

## RESPONSE SHAPES

### `GET /me/dashboard`

```json
{
  "courses": [
    {
      "course_id": "uuid",
      "course_title": "Mathematics Grade 12",
      "percent_complete": 42.5,
      "completed_at": "2026-01-01T00:00:00Z",
      "completed_lessons": 17,
      "total_lessons": 40
    }
  ],
  "current_streak": 4,
  "longest_streak": 12
}
```

### `GET /me/streak`

```json
{
  "current_streak": 4,
  "longest_streak": 12,
  "last_study_date": "2026-01-01T00:00:00Z"
}
```

### `GET /courses/:id/progress`

```json
{
  "course_id": "uuid",
  "completed_lessons": 17,
  "total_lessons": 40,
  "percent_complete": "42.50",
  "completed_at": "2026-01-01T00:00:00Z",
  "lessons": [
    {
      "lesson_id": "uuid",
      "is_completed": true,
      "completed_at": "2026-01-01T00:00:00Z",
      "video_progress_seconds": 480
    }
  ]
}
```

Important:

- `CourseProgressResponse.percent_complete` is a string.
- `DashboardItemResponse.percent_complete` is a number.
- Parse both safely into `double` in the domain layer.

### `GET /me/trend?subject=<subject>`

```json
{
  "subject": "Mathematics",
  "points": [
    {
      "attempt_id": "uuid",
      "exam_id": "uuid",
      "percentage": 72.5
    }
  ]
}
```

---

## REQUIRED FOLDER STRUCTURE

Create or verify:

```bash
mkdir -p lib/features/progress/data/datasources
mkdir -p lib/features/progress/data/models
mkdir -p lib/features/progress/data/repositories
mkdir -p lib/features/progress/domain/entities
mkdir -p lib/features/progress/domain/repositories
mkdir -p lib/features/progress/domain/usecases
mkdir -p lib/features/progress/presentation/providers
mkdir -p lib/features/progress/presentation/screens
mkdir -p lib/features/progress/presentation/widgets
```

---

## PART A - DOMAIN ENTITIES

### FILE A1 - `lib/features/progress/domain/entities/lesson_progress_entity.dart`

Fields:

- `lessonId`
- `isCompleted`
- `completedAt`
- `videoProgressSeconds`

Helpers:

- `bool get hasVideoProgress`
- `String get progressLabel`

Keep this entity independent from JSON and Flutter imports.

---

### FILE A2 - `lib/features/progress/domain/entities/course_progress_entity.dart`

Fields:

- `courseId`
- `completedLessons`
- `totalLessons`
- `percentComplete`
- `completedAt`
- `lessons`

`lessons` is `List<LessonProgressEntity>`.

Helpers:

- `bool get isStarted`
- `bool get isComplete`
- `double get normalizedPercent`
- `String get percentLabel`
- `String get lessonCountLabel`

Rules:

- guard against division by zero
- clamp visual percent between `0` and `100`
- do not mark complete locally unless backend `percentComplete >= 100` or `completedAt != null`

---

### FILE A3 - `lib/features/progress/domain/entities/dashboard_course_entity.dart`

Fields:

- `courseId`
- `courseTitle`
- `percentComplete`
- `completedAt`
- `completedLessons`
- `totalLessons`

Helpers:

- `bool get isComplete`
- `String get percentLabel`
- `String get lessonCountLabel`

---

### FILE A4 - `lib/features/progress/domain/entities/student_dashboard_entity.dart`

Fields:

- `courses`
- `currentStreak`
- `longestStreak`

`courses` is `List<DashboardCourseEntity>`.

Helpers:

- `int get enrolledCourseCount`
- `int get completedCourseCount`
- `double get averageCompletionPercent`
- `DashboardCourseEntity? get nextCourseToContinue`

`nextCourseToContinue` should prefer:

1. in-progress courses
2. most recently completed courses only if no in-progress courses exist
3. `null` when there are no courses

Do not invent timestamps that are not in the response.

---

### FILE A5 - `lib/features/progress/domain/entities/study_streak_entity.dart`

Fields:

- `currentStreak`
- `longestStreak`
- `lastStudyDate`

Helpers:

- `bool get hasStudiedBefore`
- `String get currentStreakLabel`
- `String get longestStreakLabel`

---

### FILE A6 - `lib/features/progress/domain/entities/subject_trend_entity.dart`

Create:

```dart
class TrendPointEntity {
  const TrendPointEntity({
    required this.attemptId,
    required this.examId,
    required this.percentage,
  });

  final String attemptId;
  final String examId;
  final double percentage;
}
```

And:

```dart
class SubjectTrendEntity {
  const SubjectTrendEntity({
    required this.subject,
    required this.points,
  });

  final String subject;
  final List<TrendPointEntity> points;
}
```

Helpers:

- `bool get hasTrend`
- `double? get latestPercentage`
- `double? get bestPercentage`
- `double? get averagePercentage`
- `double? get deltaFromFirst`

---

## PART B - REPOSITORY CONTRACT

### FILE B1 - `lib/features/progress/domain/repositories/progress_repository.dart`

Create:

```dart
abstract class ProgressRepository {
  Future<Either<Failure, StudentDashboardEntity>> getDashboard();
  Future<Either<Failure, StudyStreakEntity>> getStreak();
  Future<Either<Failure, CourseProgressEntity>> getCourseProgress(String courseId);
  Future<Either<Failure, SubjectTrendEntity>> getSubjectTrend(String subject);
}
```

Rules:

- Return `Either<Failure, T>` like earlier phases.
- Do not expose DTO models outside the data layer.
- Do not make screens call this repository directly.

---

## PART C - USE CASES

Create one use case per repository method.

### FILE C1 - `lib/features/progress/domain/usecases/get_student_dashboard_usecase.dart`

Calls `repository.getDashboard()`.

### FILE C2 - `lib/features/progress/domain/usecases/get_study_streak_usecase.dart`

Calls `repository.getStreak()`.

### FILE C3 - `lib/features/progress/domain/usecases/get_course_progress_usecase.dart`

Accepts `courseId`; validates non-empty before calling repository.

### FILE C4 - `lib/features/progress/domain/usecases/get_subject_trend_usecase.dart`

Accepts `subject`; validates non-empty before calling repository.

Failure behavior:

- empty `courseId` -> validation failure
- empty `subject` -> validation failure
- network/server errors -> map through repository implementation

---

## PART D - DATA MODELS

### FILE D1 - `lib/features/progress/data/models/lesson_progress_model.dart`

Parse:

- `lesson_id`
- `is_completed`
- `completed_at`
- `video_progress_seconds`

Map to `LessonProgressEntity`.

---

### FILE D2 - `lib/features/progress/data/models/course_progress_model.dart`

Parse:

- `course_id`
- `completed_lessons`
- `total_lessons`
- `percent_complete`
- `completed_at`
- `lessons`

`percent_complete` is a backend string for this endpoint. Parse with:

```dart
double.tryParse(value.toString()) ?? 0
```

Map nested lessons to `LessonProgressEntity`.

---

### FILE D3 - `lib/features/progress/data/models/dashboard_model.dart`

Create both:

- `DashboardCourseModel`
- `StudentDashboardModel`

Parse:

- `course_id`
- `course_title`
- `percent_complete`
- `completed_at`
- `completed_lessons`
- `total_lessons`
- `courses`
- `current_streak`
- `longest_streak`

`percent_complete` is numeric here, but still parse defensively with
`num`/`toDouble`.

---

### FILE D4 - `lib/features/progress/data/models/study_streak_model.dart`

Parse:

- `current_streak`
- `longest_streak`
- `last_study_date`

---

### FILE D5 - `lib/features/progress/data/models/subject_trend_model.dart`

Create both:

- `TrendPointModel`
- `SubjectTrendModel`

Parse:

- `attempt_id`
- `exam_id`
- `percentage`
- `subject`
- `points`

---

## PART E - REMOTE DATASOURCE

### FILE E1 - `lib/features/progress/data/datasources/progress_remote_datasource.dart`

Create:

```dart
abstract class ProgressRemoteDataSource {
  Future<StudentDashboardModel> getDashboard();
  Future<StudyStreakModel> getStreak();
  Future<CourseProgressModel> getCourseProgress(String courseId);
  Future<SubjectTrendModel> getSubjectTrend(String subject);
}
```

Implementation paths:

```dart
dio.get('/me/dashboard')
dio.get('/me/streak')
dio.get('/courses/$courseId/progress')
dio.get('/me/trend', queryParameters: {'subject': subject})
```

Rules:

- Use the shared Dio client from core providers.
- Let the auth interceptor attach tokens.
- Do not hardcode the full base URL.
- Do not swallow server error bodies before failure mapping.

---

## PART F - REPOSITORY IMPLEMENTATION

### FILE F1 - `lib/features/progress/data/repositories/progress_repository_impl.dart`

Implement `ProgressRepository`.

Pattern:

- call remote datasource
- convert model to entity
- catch `ServerException`
- catch `NetworkException`
- catch unknown exceptions as `UnexpectedFailure`

Keep failure mapping consistent with previous feature repositories.

---

## PART G - RIVERPOD DEPENDENCY PROVIDERS

### FILE G1 - `lib/features/progress/presentation/providers/progress_providers.dart`

Create providers:

- `progressRemoteDataSourceProvider`
- `progressRepositoryProvider`
- `getStudentDashboardUseCaseProvider`
- `getStudyStreakUseCaseProvider`
- `getCourseProgressUseCaseProvider`
- `getSubjectTrendUseCaseProvider`

Rules:

- providers expose use cases to presentation
- screens consume feature providers from `SKILL_2.md` and `SKILL_3.md`
- avoid one giant provider that loads all progress data unconditionally

---

## PART H - DATA-LAYER TESTING NOTES

Unit tests:

- dashboard model parses numeric `percent_complete`
- course progress model parses string `percent_complete`
- missing optional date fields parse as null
- empty lessons list is valid
- trend response with empty points is valid
- repository maps network/server failures

---

## VALIDATION

Run:

```bash
dart format lib/features/progress
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

- `/me/dashboard` loads after login
- `/me/streak` loads after login
- `/courses/:id/progress` loads for enrolled course
- `/me/trend?subject=Mathematics` handles both data and empty points

---

## SKILL_1 COMPLETE WHEN

- [ ] Progress domain entities compile
- [ ] Repository contract exists
- [ ] Use cases exist and validate required IDs/subjects
- [ ] DTO models parse all backend response shapes
- [ ] Course progress string percent parses safely
- [ ] Dashboard numeric percent parses safely
- [ ] Remote datasource calls the correct endpoints
- [ ] Repository implementation maps failures consistently
- [ ] Riverpod dependency providers are wired

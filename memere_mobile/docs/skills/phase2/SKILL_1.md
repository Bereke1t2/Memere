# phase2/SKILL_1.md — Course Domain, Data Layer & Backend Contract
# Memere Mobile (memere_mobile) — Phase 2, Part 1
# READ SKILL.md → phase1/SKILL_1.md → phase1/SKILL_2.md → phase1/SKILL_3.md → then this file.

---

## OBJECTIVE

Build the Course Browsing foundation:
Course entities → repository interface → list/detail use cases →
backend DTO models → remote datasource → repository implementation →
Riverpod dependency providers.

By the end of this skill, the app can load course lists and course details from
the local Memere backend, but the UI can still be a placeholder.

---

## PHASE 2 BOUNDARY

Phase 2 builds **student course browsing only**:

- Course catalog/list
- Search and subject/grade filters
- Course detail
- Sections and lesson list
- Loading, error, empty states

Do **not** build these yet:

- HLS video playback
- Offline downloads
- Quiz/exam flows
- Payment/enrollment WebView
- Teacher/admin course authoring
- Progress dashboard

---

## BACKEND CONTRACT

Use the existing local backend API:

```text
GET /api/v1/courses?limit=20&after=<cursor>&subject=<subject>&grade=12
GET /api/v1/courses/:id
GET /api/v1/courses/:id/sections
GET /api/v1/sections/:id/lessons
```

Preferred Phase 2 detail call:

```text
GET /api/v1/courses/:id
```

It already returns nested sections and lessons.

### Course list response

```json
{
  "data": [
    {
      "id": "uuid",
      "teacher_id": "uuid",
      "title": "Mathematics Grade 12",
      "slug": "mathematics-grade-12",
      "description": "Full description",
      "short_description": "Card description",
      "subject": "Mathematics",
      "grade": 12,
      "thumbnail_url": "https://...",
      "price": 250,
      "currency": "ETB",
      "is_free": false,
      "is_published": true,
      "language": "en",
      "level": "beginner",
      "total_duration_seconds": 36000,
      "total_lessons": 40,
      "rating_avg": 4.7,
      "enrollment_count": 1200,
      "metadata": {},
      "created_at": "2026-01-01T00:00:00Z",
      "updated_at": "2026-01-01T00:00:00Z"
    }
  ],
  "next_cursor": "",
  "limit": 20
}
```

### Course detail response

```json
{
  "id": "uuid",
  "title": "Mathematics Grade 12",
  "...": "same course fields as list",
  "sections": [
    {
      "id": "uuid",
      "course_id": "uuid",
      "title": "Algebra Foundations",
      "description": "Section overview",
      "order_index": 1,
      "is_published": true,
      "created_at": "2026-01-01T00:00:00Z",
      "updated_at": "2026-01-01T00:00:00Z",
      "lessons": [
        {
          "id": "uuid",
          "section_id": "uuid",
          "course_id": "uuid",
          "title": "Linear Equations",
          "type": "video",
          "order_index": 1,
          "is_free_preview": true,
          "duration_seconds": 900,
          "is_published": true,
          "created_at": "2026-01-01T00:00:00Z",
          "updated_at": "2026-01-01T00:00:00Z"
        }
      ]
    }
  ]
}
```

### Search behavior

The backend course handler supports `subject` and `grade`. It does **not**
currently expose a `q` search parameter. Implement text search locally over the
currently loaded course list unless the backend later adds `q`.

---

## REQUIRED FOLDER STRUCTURE

Create or verify these folders:

```bash
mkdir -p lib/features/courses/data/datasources
mkdir -p lib/features/courses/data/models
mkdir -p lib/features/courses/data/repositories
mkdir -p lib/features/courses/domain/entities
mkdir -p lib/features/courses/domain/repositories
mkdir -p lib/features/courses/domain/usecases
mkdir -p lib/features/courses/presentation/providers
mkdir -p lib/features/courses/presentation/screens
mkdir -p lib/features/courses/presentation/widgets
```

---

## PART A — DOMAIN ENTITIES

### FILE A1 — `lib/features/courses/domain/entities/course_entity.dart`

Create a pure Dart entity. No JSON imports.

Fields:

- `id`
- `teacherId`
- `title`
- `slug`
- `description`
- `shortDescription`
- `subject`
- `grade`
- `thumbnailUrl`
- `price`
- `currency`
- `isFree`
- `isPublished`
- `language`
- `level`
- `totalDurationSeconds`
- `totalLessons`
- `ratingAvg`
- `enrollmentCount`
- `createdAt`
- `updatedAt`

Add helpers:

- `bool get isPaid`
- `String get priceLabel`
- `String get durationLabel`
- `String get lessonCountLabel`

Use enums:

```dart
enum CourseLevel { beginner, intermediate, advanced }
```

---

### FILE A2 — `lib/features/courses/domain/entities/course_section_entity.dart`

Fields:

- `id`
- `courseId`
- `title`
- `description`
- `orderIndex`
- `isPublished`
- `createdAt`
- `updatedAt`
- `lessons`

`lessons` is a `List<LessonEntity>`.

---

### FILE A3 — `lib/features/courses/domain/entities/lesson_entity.dart`

Fields:

- `id`
- `sectionId`
- `courseId`
- `title`
- `type`
- `orderIndex`
- `isFreePreview`
- `durationSeconds`
- `isPublished`
- `createdAt`
- `updatedAt`

Use enum:

```dart
enum LessonType { video, note, quiz, mixed }
```

Add helper:

- `String get durationLabel`

---

### FILE A4 — `lib/features/courses/domain/entities/course_detail_entity.dart`

Wrap the selected course plus nested content:

```dart
class CourseDetailEntity {
  const CourseDetailEntity({
    required this.course,
    required this.sections,
  });

  final CourseEntity course;
  final List<CourseSectionEntity> sections;
}
```

---

### FILE A5 — `lib/features/courses/domain/entities/paginated_courses_entity.dart`

```dart
class PaginatedCoursesEntity {
  const PaginatedCoursesEntity({
    required this.courses,
    required this.nextCursor,
    required this.limit,
  });

  final List<CourseEntity> courses;
  final String? nextCursor;
  final int limit;

  bool get hasMore => nextCursor != null && nextCursor!.isNotEmpty;
}
```

---

## PART B — REPOSITORY CONTRACT

### FILE B1 — `lib/features/courses/domain/repositories/courses_repository.dart`

Use `Either<Failure, T>` just like auth.

Required methods:

```dart
Future<Either<Failure, PaginatedCoursesEntity>> listCourses({
  int limit = 20,
  String? after,
  String? subject,
  int? grade,
});

Future<Either<Failure, CourseDetailEntity>> getCourseDetail(String courseId);
```

Do not expose Dio, JSON, or models from the domain layer.

---

## PART C — USE CASES

### FILE C1 — `lib/features/courses/domain/usecases/list_courses_usecase.dart`

Create `ListCoursesParams`:

- `limit`
- `after`
- `subject`
- `grade`
- `searchQuery`

`searchQuery` is presentation-only for local filtering. It should not be sent to
the backend until the backend supports it.

Use case calls repository `listCourses`.

---

### FILE C2 — `lib/features/courses/domain/usecases/get_course_detail_usecase.dart`

Validate `courseId` is not empty.

Return `ValidationFailure` when empty.

---

## PART D — DATA MODELS

### FILE D1 — `lib/features/courses/data/models/course_model.dart`

`CourseModel extends CourseEntity`.

Implement:

- `factory CourseModel.fromJson(Map<String, dynamic> json)`
- `Map<String, dynamic> toJson()`
- role-safe parsing for nullable backend fields
- level parsing fallback to `CourseLevel.beginner`

Backend field mapping:

| Dart | JSON |
|------|------|
| `teacherId` | `teacher_id` |
| `shortDescription` | `short_description` |
| `thumbnailUrl` | `thumbnail_url` |
| `isFree` | `is_free` |
| `isPublished` | `is_published` |
| `totalDurationSeconds` | `total_duration_seconds` |
| `totalLessons` | `total_lessons` |
| `ratingAvg` | `rating_avg` |
| `enrollmentCount` | `enrollment_count` |
| `createdAt` | `created_at` |
| `updatedAt` | `updated_at` |

---

### FILE D2 — `lib/features/courses/data/models/lesson_model.dart`

`LessonModel extends LessonEntity`.

Parse:

- `section_id`
- `course_id`
- `order_index`
- `is_free_preview`
- `duration_seconds`
- `is_published`

Fallback lesson type to `LessonType.video`.

---

### FILE D3 — `lib/features/courses/data/models/course_section_model.dart`

`CourseSectionModel extends CourseSectionEntity`.

Parse nested lessons:

```dart
final lessonsJson = json['lessons'];
final lessons = lessonsJson is List
    ? lessonsJson
        .whereType<Map<String, dynamic>>()
        .map(LessonModel.fromJson)
        .toList()
    : <LessonModel>[];
```

---

### FILE D4 — `lib/features/courses/data/models/course_detail_model.dart`

`CourseDetailModel extends CourseDetailEntity`.

Use:

- `CourseModel.fromJson(json)` for the course
- `CourseSectionModel.fromJson` for `sections`

---

### FILE D5 — `lib/features/courses/data/models/paginated_courses_model.dart`

Parse:

- `data`
- `next_cursor`
- `limit`

Return an empty course list if `data` is missing or not a list.

---

## PART E — REMOTE DATASOURCE

### FILE E1 — `lib/features/courses/data/datasources/courses_remote_datasource.dart`

Use `DioClient` only. Do not create a raw Dio instance.

Required methods:

```dart
Future<PaginatedCoursesModel> listCourses({
  int limit = 20,
  String? after,
  String? subject,
  int? grade,
});

Future<CourseDetailModel> getCourseDetail(String courseId);
```

Endpoint behavior:

```dart
await _client.get<Map<String, dynamic>>(
  '/courses',
  queryParameters: {
    'limit': limit,
    if (after != null && after.isNotEmpty) 'after': after,
    if (subject != null && subject.isNotEmpty) 'subject': subject,
    if (grade != null) 'grade': grade,
  },
);
```

Detail:

```dart
await _client.get<Map<String, dynamic>>('/courses/$courseId');
```

Throw `FormatException` for missing response bodies.

---

## PART F — REPOSITORY IMPLEMENTATION

### FILE F1 — `lib/features/courses/data/repositories/courses_repository_impl.dart`

Map errors consistently:

- `DioException` → `ServerFailure.fromDioError`
- any other exception → `UnknownFailure`

Do not do UI filtering in repository. Keep local search in provider/presentation
state.

---

## PART G — PROVIDER REGISTRATION

### FILE G1 — `lib/features/courses/presentation/providers/courses_providers.dart`

Create dependency providers:

- `coursesRemoteDataSourceProvider`
- `coursesRepositoryProvider`
- `listCoursesUseCaseProvider`
- `getCourseDetailUseCaseProvider`

These mirror the auth provider style.

---

## VALIDATION

After this skill:

```bash
dart format lib/features/courses
flutter analyze
```

Do not continue to Phase 2 UI until this passes.

---

## SKILL_1 CHECKLIST

- [ ] Course domain entities compile
- [ ] Repository interface has list/detail methods
- [ ] Use cases return `Either<Failure, T>`
- [ ] Models parse backend JSON shape exactly
- [ ] Remote datasource calls `/courses` and `/courses/:id`
- [ ] Repository impl maps failures consistently
- [ ] Dependency providers are wired
- [ ] `flutter analyze` has 0 errors

# phase2/SKILL_3.md — Course Detail, Sections, Lesson List & Phase 2 Checklist
# Memere Mobile (memere_mobile) — Phase 2, Part 3
# READ SKILL.md → phase2/SKILL_1.md → phase2/SKILL_2.md → then this file.

---

## OBJECTIVE

Build the course detail experience:
course detail route → backend detail fetch → course header → stats → description →
sections accordion → lesson list → Phase 2-safe CTA placeholders.

By the end of this skill, users can browse from catalog to a course detail page
and inspect the course curriculum from backend data.

---

## PHASE 2 DETAIL BOUNDARY

The detail screen may show lessons, preview labels, and CTA buttons, but it must
not implement Phase 3+ behavior.

Allowed:

- open course detail
- inspect sections
- inspect lesson titles
- show locked/free-preview visual states
- show `Start learning` / `Enroll` placeholder CTA
- show snackbar saying the next phase handles playback/enrollment

Not allowed:

- launching video player
- starting quiz
- initiating payment
- enrolling a user
- marking progress complete

---

## PART A — DETAIL PROVIDER

### FILE A1 — `lib/features/courses/presentation/providers/course_detail_provider.dart`

Create:

```dart
final courseDetailProvider = FutureProvider.family<CourseDetailEntity, String>(
  (ref, courseId) async {
    final useCase = ref.watch(getCourseDetailUseCaseProvider);
    final result = await useCase(courseId);
    return result.fold(
      (failure) => throw failure,
      (detail) => detail,
    );
  },
);
```

If you prefer `AsyncNotifierProvider.family`, keep it simple. Detail only needs
initial load and refresh in Phase 2.

---

## PART B — COURSE DETAIL SCREEN

### FILE B1 — `lib/features/courses/presentation/screens/course_detail_screen.dart`

Constructor:

```dart
class CourseDetailScreen extends ConsumerWidget {
  const CourseDetailScreen({
    super.key,
    required this.courseId,
  });

  final String courseId;
}
```

Required states:

- loading skeleton
- error state with retry/back
- loaded detail view

Use:

```dart
ref.watch(courseDetailProvider(courseId))
```

Retry:

```dart
ref.invalidate(courseDetailProvider(courseId));
```

---

## PART C — DETAIL LAYOUT

### Loaded layout order

1. Collapsing or fixed top bar with back button
2. Course visual/header
3. Subject, grade, level chips
4. Title
5. Short description or description preview
6. Stats row
7. Price/free CTA band
8. About course
9. Curriculum section list
10. Bottom CTA button

### Header requirements

Use `CachedNetworkImage` for `thumbnailUrl` if present.

Fallback when no thumbnail:

- dark subject-color surface
- subject icon
- subject name

Avoid giant marketing hero treatment. Detail should feel like an app screen.

### CTA behavior

Use a placeholder action:

```dart
ScaffoldMessenger.of(context).showSnackBar(
  const SnackBar(content: Text('Enrollment starts in Phase 6.')),
);
```

If course is free:

```text
Start learning
```

If course is paid:

```text
Enroll for ETB <price>
```

Do not connect payment in Phase 2.

---

## PART D — DETAIL WIDGETS

### FILE D1 — `lib/features/courses/presentation/widgets/course_detail_header.dart`

Props:

- `CourseEntity course`

Shows:

- thumbnail/fallback visual
- subject chip
- grade chip
- level chip
- title
- short description

---

### FILE D2 — `lib/features/courses/presentation/widgets/course_stats_row.dart`

Props:

- `CourseEntity course`

Stats:

- lessons
- duration
- rating
- students

Use compact icon + text items. Keep dimensions stable.

---

### FILE D3 — `lib/features/courses/presentation/widgets/course_section_tile.dart`

Props:

- `CourseSectionEntity section`
- `int sectionNumber`

Use `ExpansionTile` or custom expandable tile.

Show:

- section title
- lesson count
- optional description
- nested lesson rows

Default behavior:

- first section expanded
- other sections collapsed

If using expansion state, keep it local inside the widget only when it is purely
UI state. Do not put business logic in widgets.

---

### FILE D4 — `lib/features/courses/presentation/widgets/lesson_tile.dart`

Props:

- `LessonEntity lesson`
- `int lessonNumber`

Show:

- lesson number
- lesson type icon
- title
- duration
- free preview badge when `isFreePreview == true`
- lock icon for non-preview lesson placeholder

Tap behavior:

- free preview: snackbar `Video playback starts in Phase 3.`
- locked/paid: snackbar `Enrollment starts in Phase 6.`

Do not navigate to video/quiz screens in Phase 2.

---

### FILE D5 — `lib/features/courses/presentation/widgets/course_detail_skeleton.dart`

Use `shimmer`.

Skeleton sections:

- header block
- title rows
- stats row
- 3 section placeholders

---

## PART E — FORMATTERS

If helpers do not already exist, create or extend:

### FILE E1 — `lib/shared/utils/formatters.dart`

Add pure functions:

```dart
String formatDurationSeconds(int seconds)
String formatCompactCount(int count)
String formatPrice(double price, String currency)
```

Expected output:

- `formatDurationSeconds(65)` → `1m`
- `formatDurationSeconds(3600)` → `1h`
- `formatDurationSeconds(5400)` → `1h 30m`
- `formatCompactCount(1200)` → `1.2k`
- `formatPrice(250, 'ETB')` → `ETB 250`

Use these helpers in entities or widgets. Avoid duplicating formatter logic.

---

## PART F — ROUTER DETAIL WIRING

Update `app_router.dart` route:

```dart
GoRoute(
  path: AppRoutes.courseDetail,
  builder: (_, state) {
    final courseId = state.pathParameters['courseId']!;
    return CourseDetailScreen(courseId: courseId);
  },
),
```

Course cards should navigate with:

```dart
context.go(AppRoutes.courseDetailPath(course.id));
```

---

## PART G — ERROR HANDLING

Failure display:

```dart
final message = error is Failure
    ? error.message
    : 'Could not load course. Please try again.';
```

For 404:

- title: `Course not found`
- body: `This course may have been removed or unpublished.`

For other errors:

- title: `Could not load course`
- retry button

---

## PART H — ACCESSIBILITY AND MOBILE FIT

Before marking Phase 2 complete:

- Test narrow Android width.
- Ensure course titles wrap cleanly.
- Ensure filter chips do not overflow.
- Ensure bottom CTA is not hidden behind system navigation.
- Ensure skeletons have stable sizes.
- Ensure empty/error states fit without clipping.

---

## VALIDATION

Run:

```bash
dart format lib/shared lib/core/router lib/features/courses
flutter analyze
flutter build apk --debug
```

Manual local backend run:

```bash
cd ../Memere-backend
make up
make migrate-up
make seed
make run
```

Physical Android bridge:

```bash
adb reverse tcp:8080 tcp:8080
flutter run
```

---

## PHASE 2 FINAL CHECKLIST

### SKILL_1 complete when:

- [ ] Domain entities compile
- [ ] Course models parse backend JSON
- [ ] List/detail use cases work
- [ ] Repository maps backend errors to `Failure`
- [ ] Course dependency providers are wired

### SKILL_2 complete when:

- [ ] `/home` opens the course catalog
- [ ] Course list loads from `/courses`
- [ ] Subject and grade filters work
- [ ] Search works locally
- [ ] Pull-to-refresh works
- [ ] Pagination uses `next_cursor`
- [ ] Loading, empty, and error states render

### SKILL_3 complete when:

- [ ] Course cards navigate to `/courses/:courseId`
- [ ] Detail loads from `/courses/:id`
- [ ] Header, stats, about, sections, and lessons render
- [ ] Free preview/locked lesson states are visible
- [ ] Phase 3/6 CTAs show placeholders only
- [ ] `flutter analyze` has 0 errors
- [ ] Android debug APK builds

---

## PHASE 2 → PHASE 3 HANDOFF

**Phase 2 is complete when all 3 SKILL files are done and the checklist above passes.**

To start Phase 3, tell Antigravity:

```text
Phase 2 is complete. Read SKILL.md and all phase3 skill files.
We are starting Phase 3: HLS Video Player + Offline Download.
Reference: memere_mobile/docs/memere_Design_Specification.md
```

**What Phase 3 will build:**

- Video player screen
- HLS playback with `video_player`/`chewie`
- stream URL fetch from backend
- video progress save hooks
- download URL flow
- offline download queue foundation

# phase7/SKILL_3.md - Course Progress Detail, Trends, Weak Areas & Phase 7 Checklist
# Memere Mobile (memere_mobile) - Phase 7, Part 3
# READ SKILL.md -> phase7/SKILL_1.md -> phase7/SKILL_2.md -> then this file.

---

## OBJECTIVE

Finish Phase 7 by adding:
course progress detail -> lesson progress list -> subject trend summary ->
weak areas preview -> dashboard refresh integration -> final validation.

By the end of this skill, students can inspect progress for an enrolled course
and see lightweight study analytics without adding new backend endpoints.

---

## PHASE 7 FINAL BOUNDARY

Included:

- course progress detail screen
- lesson completion list
- course completion percent
- video progress seconds display
- subject trend preview
- weak areas preview from existing exam analytics when available
- dashboard refresh after learning activity

Deferred:

- full adaptive learning
- AI-generated study plans
- certificates
- notification center
- teacher/admin analytics
- global weak-areas backend endpoint

---

## PART A - COURSE PROGRESS PROVIDER

### FILE A1 - `lib/features/progress/presentation/providers/course_progress_provider.dart`

Create:

```dart
final courseProgressProvider =
    FutureProvider.family<CourseProgressEntity, String>((ref, courseId) async {
  final useCase = ref.watch(getCourseProgressUseCaseProvider);
  final result = await useCase(courseId);
  return result.fold((failure) => throw failure, (progress) => progress);
});
```

Refresh this provider after:

- marking lesson complete
- updating video progress
- returning from video player
- successful enrollment/payment when opening the detail page

---

## PART B - SUBJECT TREND PROVIDER

### FILE B1 - `lib/features/progress/presentation/providers/subject_trend_provider.dart`

Create:

```dart
final subjectTrendProvider =
    FutureProvider.family<SubjectTrendEntity, String>((ref, subject) async {
  final useCase = ref.watch(getSubjectTrendUseCaseProvider);
  final result = await useCase(subject);
  return result.fold((failure) => throw failure, (trend) => trend);
});
```

Rules:

- subject must come from selected course/exam context, not a hardcoded value
- empty trend points are valid
- trend preview must not block the course progress screen

---

## PART C - COURSE PROGRESS SCREEN

### FILE C1 - `lib/features/progress/presentation/screens/course_progress_screen.dart`

Create:

```dart
class CourseProgressScreen extends ConsumerWidget {
  const CourseProgressScreen({
    super.key,
    required this.courseId,
    this.courseTitle,
    this.subject,
  });

  final String courseId;
  final String? courseTitle;
  final String? subject;
}
```

Route:

```dart
/courses/:courseId/progress
```

Required UI states:

- loading skeleton
- error with retry
- loaded progress
- empty lesson progress state

Loaded layout:

1. app bar
2. course progress summary
3. lesson progress list
4. subject trend card when subject is known
5. weak areas preview when source data is available
6. actions:
   - continue course
   - back to dashboard

Navigation:

- continue course routes to course detail or current lesson if known
- lesson rows route to their lesson screen only if route exists
- otherwise rows are read-only progress rows

---

## PART D - COURSE PROGRESS SUMMARY WIDGET

### FILE D1 - `lib/features/progress/presentation/widgets/course_progress_summary_card.dart`

Props:

- `CourseProgressEntity progress`
- optional `courseTitle`

Show:

- course title or `Course progress`
- percent complete
- completed lessons / total lessons
- completed date if present
- progress bar or compact ring

Rules:

- if `totalLessons == 0`, show `No lessons yet`
- if complete, show `Completed`
- do not use oversized hero styling

---

## PART E - LESSON PROGRESS LIST

### FILE E1 - `lib/features/progress/presentation/widgets/lesson_progress_list.dart`

Props:

- `List<LessonProgressEntity> lessons`
- optional lesson metadata map from course detail if available

Show:

- completed state
- lesson title if known
- fallback `Lesson <number>`
- video progress label if `videoProgressSeconds > 0`
- completed date if present

Do not fabricate lesson titles. If the course detail entity is available, merge
title/duration by `lessonId`; otherwise use the fallback label.

---

### FILE E2 - `lib/features/progress/presentation/widgets/lesson_progress_tile.dart`

Props:

- `LessonProgressEntity progress`
- `String title`
- optional `VoidCallback onTap`

Visual states:

- completed: success icon/check
- started video: progress indicator or secondary text
- not started: neutral state

Use stable row height so list items do not jump when data refreshes.

---

## PART F - SUBJECT TREND CARD

### FILE F1 - `lib/features/progress/presentation/widgets/subject_trend_card.dart`

Props:

- `SubjectTrendEntity trend`

Show:

- subject
- latest percentage
- best percentage
- delta from first attempt when available
- compact trend visual

Implementation options:

- use a simple custom painter line if no chart package exists
- use a vertical list of recent attempts if charting would add unnecessary dependency

Do not add a new chart dependency just for Phase 7 unless the project already
uses one.

Empty state:

```text
No exam trend yet.
```

---

## PART G - WEAK AREAS PREVIEW

### FILE G1 - `lib/features/progress/presentation/widgets/weak_areas_preview_card.dart`

Phase 7 does not have a global weak-areas endpoint. Build this card only from
data already available through Phase 5 exam analytics, such as
`GET /exam-attempts/:id/analytics`, when the current screen has a recent
attempt ID.

If no recent attempt analytics are available, show a compact empty state:

```text
Complete a mock exam to see weak areas.
```

Do not:

- invent weak areas from course progress
- infer correctness from quiz options
- show answer keys before result/analytics screens
- create fake recommendation data

---

## PART H - RECENT ACTIVITY SUMMARY

### FILE H1 - `lib/features/progress/presentation/widgets/recent_activity_summary.dart`

Use existing backend data only.

Allowed sources:

- completed course dates from dashboard/course progress
- lesson completed dates from course progress
- last study date from streak
- recent exam result/analytics if Phase 5 providers expose it

If there is not enough data:

```text
Recent activity will appear as you study.
```

Do not create a local fake activity feed.

---

## PART I - VIDEO/LESSON PROGRESS REFRESH INTEGRATION

Update Phase 3 learning flow where appropriate:

- after `POST /lessons/:id/complete`, invalidate:
  - `studentDashboardProvider`
  - `courseProgressProvider(courseId)`
  - `studyStreakProvider`
- after `PUT /lessons/:id/video-progress`, invalidate:
  - `courseProgressProvider(courseId)` when returning from video player

Avoid refreshing dashboard on every video-progress tick. Save video progress as
Phase 3 requires, then refresh dashboard only when the user leaves the lesson or
when a lesson is completed.

---

## PART J - PAYMENT/ENROLLMENT REFRESH INTEGRATION

Update Phase 6 success paths:

- after free enrollment:
  - invalidate `studentDashboardProvider`
  - invalidate course access provider
- after paid course payment completed:
  - invalidate `studentDashboardProvider`
  - invalidate course access provider
- after subscription payment completed:
  - invalidate course access provider
  - refresh dashboard only if backend dashboard includes newly available courses

Do not assume subscription means every course should appear in dashboard
progress. Let backend dashboard decide.

---

## PART K - ERROR HANDLING

Handle:

- unauthenticated -> route guard/login flow
- not enrolled -> show access required state and link to course detail
- course not found -> show friendly error and back action
- empty trend -> normal empty state
- analytics unavailable -> hide weak areas or show empty state
- backend percent parse issue -> default to 0 and log in debug only

Do not expose raw backend errors directly in user-facing text.

---

## PART L - TESTING NOTES

Widget tests:

- course progress loading skeleton renders
- course progress summary handles 0, partial, and 100 percent
- lesson list uses fallback labels without metadata
- trend card handles empty, one-point, and multi-point trend data
- weak areas card shows empty state without recent attempt analytics

Provider tests:

- course progress provider loads by course ID
- subject trend provider loads by subject
- invalid course ID maps to validation failure
- invalid subject maps to validation failure

Integration/manual tests:

- complete lesson updates course detail progress after refresh
- video progress appears after returning from player
- dashboard average changes after course progress changes
- completed course shows completed state

---

## VALIDATION

Run:

```bash
dart format lib/core/router lib/features/progress lib/features/video_player lib/features/payment
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

- dashboard loads after login
- course progress detail opens from dashboard
- lesson progress list displays backend data
- completing a lesson refreshes dashboard and course progress
- returning from video player refreshes course progress
- subject trend empty state works
- weak areas preview does not invent data
- no screen shows an old project name

---

## PHASE 7 FINAL CHECKLIST

### SKILL_1 complete when:

- [ ] Progress entities/models compile
- [ ] Repository/use cases cover dashboard, streak, course progress, trend
- [ ] Remote datasource calls all progress endpoints
- [ ] Course progress parses string percent safely
- [ ] Dashboard parses numeric percent safely
- [ ] Failure mapping is consistent

### SKILL_2 complete when:

- [ ] Dashboard route exists
- [ ] Student dashboard screen loads backend data
- [ ] Streak/progress summary/continue learning widgets render
- [ ] Enrolled course progress cards render
- [ ] Empty and error states are implemented
- [ ] Home/bottom nav exposes dashboard

### SKILL_3 complete when:

- [ ] Course progress route exists
- [ ] Course progress screen loads backend data
- [ ] Lesson progress list renders
- [ ] Subject trend card handles data and empty states
- [ ] Weak areas preview uses only existing analytics data
- [ ] Progress refreshes after lesson/video/enrollment changes
- [ ] `flutter analyze` has 0 errors
- [ ] Android debug APK builds

---

## PHASE 7 -> PHASE 8 HANDOFF

**Phase 7 is complete when all 3 SKILL files are done and the checklist above passes.**

To start Phase 8, tell Antigravity:

```text
Phase 7 is complete. Read SKILL.md and phase8/SKILL_1.md.
We are starting Phase 8: Notifications.
Reference: memere_mobile/docs/memere_Design_Specification.md
```

**What Phase 8 will build:**

- FCM push setup
- notification device token registration
- in-app notification center
- notification read/unread state
- notification settings

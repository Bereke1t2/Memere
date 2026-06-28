# phase7/SKILL_2.md - Student Dashboard Screen, Home Integration & Progress Widgets
# Memere Mobile (memere_mobile) - Phase 7, Part 2
# READ SKILL.md -> phase7/SKILL_1.md -> then this file.

---

## OBJECTIVE

Build the student dashboard experience:
dashboard route -> dashboard providers -> home/bottom-nav integration ->
streak summary -> enrolled course progress cards -> continue learning ->
recent activity style summaries -> loading, empty, and error states.

By the end of this skill, a logged-in student can open a useful progress
dashboard and continue learning from enrolled course progress.

---

## DASHBOARD RULES

- Dashboard data comes from `GET /me/dashboard`.
- Streak data comes from `GET /me/streak` or the dashboard response.
- Course cards must use backend completion counts and percentages.
- Do not calculate hidden lesson completion from cached course detail alone.
- Do not show paid locked courses as enrolled progress.
- Do not add marketing copy or large hero sections.
- Keep the dashboard dense, clear, and useful.

---

## PART A - ROUTING

### FILE A1 - `lib/core/router/app_router.dart`

Add routes:

```dart
static const dashboard = '/dashboard';
static const courseProgress = '/courses/:courseId/progress';

static String courseProgressPath(String courseId) {
  return '/courses/$courseId/progress';
}
```

Route builders:

- `/dashboard` -> `StudentDashboardScreen`
- `/courses/:courseId/progress` -> `CourseProgressScreen` from `SKILL_3.md`

Keep both routes behind auth.

If Phase 2 home already has bottom navigation:

- add dashboard as a tab or replace placeholder home content with dashboard
- keep catalog/search/mock exams reachable

If bottom navigation does not exist yet:

- make `/home` show the dashboard-first student shell
- include navigation entries for courses, mock exams, dashboard, profile/settings when available

---

## PART B - DASHBOARD FEATURE PROVIDERS

### FILE B1 - `lib/features/progress/presentation/providers/student_dashboard_provider.dart`

Create:

```dart
final studentDashboardProvider = FutureProvider<StudentDashboardEntity>((ref) async {
  final useCase = ref.watch(getStudentDashboardUseCaseProvider);
  final result = await useCase();
  return result.fold((failure) => throw failure, (dashboard) => dashboard);
});
```

Add a refresh helper if the codebase uses provider extensions or controller
classes. Minimum accepted behavior:

```dart
ref.invalidate(studentDashboardProvider);
```

---

### FILE B2 - `lib/features/progress/presentation/providers/study_streak_provider.dart`

Create:

```dart
final studyStreakProvider = FutureProvider<StudyStreakEntity>((ref) async {
  final useCase = ref.watch(getStudyStreakUseCaseProvider);
  final result = await useCase();
  return result.fold((failure) => throw failure, (streak) => streak);
});
```

Use this provider when the screen needs `lastStudyDate`. If the dashboard
response is enough for the current widget, do not force an extra network call.

---

### FILE B3 - `lib/features/progress/presentation/providers/dashboard_refresh_provider.dart`

Optional but recommended for pull-to-refresh:

```dart
final dashboardRefreshProvider = Provider<void Function(WidgetRef)>((ref) {
  return (widgetRef) {
    widgetRef.invalidate(studentDashboardProvider);
    widgetRef.invalidate(studyStreakProvider);
  };
});
```

Use whatever refresh pattern already exists in the project if different.

---

## PART C - STUDENT DASHBOARD SCREEN

### FILE C1 - `lib/features/progress/presentation/screens/student_dashboard_screen.dart`

Create:

```dart
class StudentDashboardScreen extends ConsumerWidget {
  const StudentDashboardScreen({super.key});
}
```

Required UI states:

- loading skeleton
- empty enrolled state
- error state with retry
- loaded dashboard

Loaded layout:

1. top app bar or shell header
2. streak card
3. progress summary row
4. continue learning card
5. enrolled course progress list
6. weak areas/trend preview from `SKILL_3.md` if available

Screen behavior:

- pull to refresh invalidates dashboard/streak providers
- tapping a course opens `AppRoutes.courseProgressPath(courseId)`
- tapping continue learning opens course detail or first available lesson route
- do not block dashboard rendering if trend preview is empty

---

## PART D - DASHBOARD HEADER

### FILE D1 - `lib/features/progress/presentation/widgets/dashboard_header.dart`

Show:

- greeting or title: `Dashboard`
- optional profile/avatar action if already available
- compact subtitle derived from progress, for example:

```text
Keep going with your enrolled courses.
```

Rules:

- do not show any old project name anywhere
- use `Memere` only where app branding is needed
- keep header compact on mobile

---

## PART E - STREAK CARD

### FILE E1 - `lib/features/progress/presentation/widgets/streak_card.dart`

Props:

- `currentStreak`
- `longestStreak`
- optional `lastStudyDate`

Show:

- current streak days
- longest streak
- last study date when available

Copy:

- `0 days` when current streak is zero
- `Study today to start a streak.` only in the empty state

Do not shame the user for a broken streak.

---

## PART F - PROGRESS SUMMARY

### FILE F1 - `lib/features/progress/presentation/widgets/progress_summary_row.dart`

Props:

- `StudentDashboardEntity dashboard`

Show three compact metrics:

- enrolled courses
- completed courses
- average completion

Use stable dimensions so numbers changing from one to three digits do not shift
the layout.

---

## PART G - CONTINUE LEARNING CARD

### FILE G1 - `lib/features/progress/presentation/widgets/continue_learning_card.dart`

Props:

- `DashboardCourseEntity? course`
- `VoidCallback? onTap`

When course exists, show:

- course title
- completion percent
- completed lessons / total lessons
- primary action: `Continue`

When no course exists:

- show a compact empty state
- action: browse courses

Do not invent lesson titles. If the backend does not provide a current lesson,
route to course detail and let the course detail screen decide the next lesson.

---

## PART H - COURSE PROGRESS CARDS

### FILE H1 - `lib/features/progress/presentation/widgets/dashboard_course_progress_card.dart`

Props:

- `DashboardCourseEntity course`
- `VoidCallback onTap`

Show:

- course title
- percent complete
- progress bar or ring
- completed lessons / total lessons
- completed badge when complete

Rules:

- progress visuals clamp between 0 and 100
- never show NaN or Infinity
- if `totalLessons == 0`, show `No lessons yet`

---

## PART I - ENROLLED COURSES SECTION

### FILE I1 - `lib/features/progress/presentation/widgets/enrolled_progress_section.dart`

Props:

- `List<DashboardCourseEntity> courses`

Behavior:

- sort in-progress courses before completed courses
- show completed courses lower in the list
- keep original backend order inside each group unless a local sort is clearly better

Empty state:

```text
No enrolled courses yet.
```

Action:

- browse courses

---

## PART J - DASHBOARD SKELETON AND EMPTY STATE

### FILE J1 - `lib/features/progress/presentation/widgets/dashboard_skeleton.dart`

Skeleton blocks:

- header
- streak card
- summary row
- continue learning card
- 3 course cards

Use the same shimmer/skeleton pattern from earlier phases.

---

### FILE J2 - `lib/features/progress/presentation/widgets/dashboard_empty_state.dart`

Show when dashboard loads with no courses.

Required actions:

- browse course catalog
- browse mock exams if Phase 5 routes exist

Keep copy short.

---

### FILE J3 - `lib/features/progress/presentation/widgets/dashboard_error_state.dart`

Props:

- `Object error`
- `VoidCallback onRetry`

Show:

- short message
- retry button

Avoid exposing raw stack traces in UI.

---

## PART K - HOME/BOTTOM NAV INTEGRATION

If `HomeScreen` exists from Phase 2:

- add dashboard as the default first tab
- keep course catalog as a separate tab
- keep mock exam tab from Phase 5
- keep purchases/subscription reachable from profile or account section

Recommended tab order:

1. Dashboard
2. Courses
3. Exams
4. Account

If a lesson/video screen returns after progress updates:

- invalidate `studentDashboardProvider`
- invalidate `courseProgressProvider(courseId)` from `SKILL_3.md`

---

## PART L - ACCESS AND ENROLLMENT INTERACTION

Course progress should only appear for enrolled/access-granted courses returned
by `/me/dashboard`.

When a user completes payment or free enrollment in Phase 6:

- refresh enrollments
- refresh course access
- refresh dashboard

When subscription grants access:

- do not show every subscribed course as progress unless backend includes it in dashboard
- course catalog still shows accessible courses through access logic

---

## PART M - TESTING NOTES

Widget tests:

- dashboard loading skeleton renders
- empty dashboard shows browse courses action
- streak card handles zero and non-zero streaks
- progress summary calculates average safely
- course progress card clamps percent visuals
- tapping course card navigates to course progress route

Provider tests:

- dashboard provider returns entity on success
- dashboard provider throws failure on repository failure
- refresh invalidates dashboard and streak providers

---

## VALIDATION

Run:

```bash
dart format lib/core/router lib/features/progress
flutter analyze
flutter build apk --debug
```

Manual checks:

- login lands on home/dashboard
- dashboard loads from local backend
- pull to refresh reloads dashboard
- enrolled course card opens course progress detail
- empty account shows browse course action
- no screen shows an old project name

---

## SKILL_2 COMPLETE WHEN

- [ ] Dashboard route exists
- [ ] Dashboard provider loads `/me/dashboard`
- [ ] Streak provider loads `/me/streak` only when needed
- [ ] Student dashboard screen has loading/empty/error/loaded states
- [ ] Streak card renders correctly
- [ ] Progress summary renders safely
- [ ] Continue learning card routes correctly
- [ ] Enrolled course progress cards route correctly
- [ ] Home/bottom navigation exposes the dashboard
- [ ] Dashboard refresh works after enrollment/payment/progress changes

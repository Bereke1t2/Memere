# phase2/SKILL_2.md — Home Shell, Course Catalog, Search & Filters
# Memere Mobile (memere_mobile) — Phase 2, Part 2
# READ SKILL.md → phase2/SKILL_1.md → then this file.

---

## OBJECTIVE

Build the student-facing course catalog UI:
authenticated home route → course list screen → search field → subject chips →
grade filter → course cards → pull-to-refresh → load-more pagination.

By the end of this skill, a logged-in student should land on the Memere course
catalog and browse courses from the local backend.

---

## DESIGN REQUIREMENTS

Use the Phase 1 design system:

- `AppColors`
- `AppTextStyles`
- `AppSizes`
- `AppButton`
- `AppTextField`

The UI must stay dark-first, compact, and study-focused. This is an operational
learning app, not a marketing landing page.

Avoid:

- oversized hero sections
- nested cards
- decorative gradients as the main content
- placeholder text explaining app features

Use:

- course cards for repeated course items
- chips/segmented controls for filters
- icon buttons for toolbar actions
- skeletons/shimmer for loading lists
- explicit empty and error states

---

## PART A — ROUTING UPDATE

### FILE A1 — `lib/core/router/app_router.dart`

Replace the temporary `_PhaseOneHomeScreen` with the real course catalog screen.

Required routes:

```dart
abstract class AppRoutes {
  static const splash = '/';
  static const onboarding = '/onboarding';
  static const login = '/login';
  static const register = '/register';
  static const home = '/home';
  static const courseDetail = '/courses/:courseId';

  static String courseDetailPath(String courseId) => '/courses/$courseId';
}
```

Route builders:

- `/home` → `CourseListScreen`
- `/courses/:courseId` → `CourseDetailScreen` (implemented in `SKILL_3`)

Keep auth redirects from Phase 1:

- logged out user cannot access `/home` or course detail
- logged in user should not remain on login/register/onboarding
- splash remains the startup decision screen

---

## PART B — COURSE LIST STATE

### FILE B1 — `lib/features/courses/presentation/providers/course_list_provider.dart`

Create an `AsyncNotifierProvider` for catalog state.

State shape:

```dart
class CourseListState {
  const CourseListState({
    this.courses = const [],
    this.filteredCourses = const [],
    this.nextCursor,
    this.selectedSubject,
    this.selectedGrade = 12,
    this.searchQuery = '',
    this.isLoadingMore = false,
  });
}
```

Required notifier methods:

- `Future<CourseListState> build()`
- `Future<void> refresh()`
- `Future<void> loadMore()`
- `void setSearchQuery(String value)`
- `Future<void> setSubject(String? subject)`
- `Future<void> setGrade(int? grade)`
- `void clearFilters()`

Rules:

- Fetch backend with `subject` and `grade`.
- Apply `searchQuery` locally against title, subject, and short description.
- Preserve old courses while loading more.
- Prevent duplicate `loadMore()` calls while `isLoadingMore == true`.
- If `nextCursor` is null/empty, do not call backend again.

Recommended subject list:

```dart
const phase2Subjects = [
  'Mathematics',
  'Physics',
  'Chemistry',
  'Biology',
  'English',
  'History',
  'Geography',
  'Economics',
];
```

---

## PART C — COURSE LIST SCREEN

### FILE C1 — `lib/features/courses/presentation/screens/course_list_screen.dart`

Build a `ConsumerWidget` or `ConsumerStatefulWidget`.

Required layout:

1. Safe area scaffold with `AppColors.bgPrimary`
2. Top header:
   - Memere app mark/icon
   - greeting line
   - logout icon button or profile placeholder
3. Search field
4. Horizontally scrolling subject filter chips
5. Grade selector (default Grade 12)
6. Course list
7. Load-more trigger near bottom

Do not use `setState`. Use Riverpod provider methods.

### Header copy

Use Memere branding:

- Title: `Explore courses`
- Subtitle: `Build your Grade 12 exam plan`

### Refresh behavior

Wrap list content in `RefreshIndicator`.

Call:

```dart
ref.read(courseListProvider.notifier).refresh()
```

### Error state

Show:

- concise message from `Failure.message`
- retry button

### Empty state

Show when `filteredCourses` is empty:

- icon
- title: `No courses found`
- body: `Try another subject or search term.`
- clear filters button when any filter/search is active

---

## PART D — WIDGETS

### FILE D1 — `lib/features/courses/presentation/widgets/course_card.dart`

Course card requirements:

- thumbnail area with `CachedNetworkImage` when `thumbnailUrl` exists
- fallback subject color block when no thumbnail
- subject chip
- course title
- short description
- lesson count
- duration
- rating
- enrollment count
- price/free badge

Tap behavior:

```dart
context.go(AppRoutes.courseDetailPath(course.id));
```

Do not start videos from the card in Phase 2.

---

### FILE D2 — `lib/features/courses/presentation/widgets/subject_filter_chips.dart`

Use horizontally scrolling chips.

Required behavior:

- `All` chip clears subject
- selected chip uses `AppColors.accentPrimary`
- unselected chips use `AppColors.bgTertiary` and border

---

### FILE D3 — `lib/features/courses/presentation/widgets/course_list_skeleton.dart`

Use `shimmer`.

Show 4-6 skeleton course cards while initial list loads.

Skeleton must use stable dimensions so the layout does not jump.

---

### FILE D4 — `lib/features/courses/presentation/widgets/course_empty_state.dart`

Reusable empty/error-style state widget with:

- icon
- title
- body
- optional button label
- optional callback

Use for no data and recoverable errors.

---

## PART E — LOGOUT ACCESS

The home/catalog screen should include a logout affordance for Phase 2 testing.

Use:

```dart
await ref.read(authStateProvider.notifier).logout();
if (context.mounted) context.go(AppRoutes.login);
```

This is temporary but useful until profile/settings screens exist.

---

## PART F — LOCAL BACKEND TEST DATA

Course list will be empty unless the backend has published courses.

Before testing, run backend seed commands if available:

```bash
cd ../Memere-backend
make up
make migrate-up
make seed
make run
```

For a physical Android phone, also run:

```bash
adb reverse tcp:8080 tcp:8080
```

---

## VALIDATION

After this skill:

```bash
dart format lib/core/router lib/features/courses
flutter analyze
flutter build apk --debug
```

---

## SKILL_2 CHECKLIST

- [ ] `/home` renders real `CourseListScreen`
- [ ] Course list loads from backend
- [ ] Search filters loaded courses locally
- [ ] Subject filter calls backend with `subject`
- [ ] Grade filter calls backend with `grade`
- [ ] Pull-to-refresh works
- [ ] Load-more uses backend `next_cursor`
- [ ] Loading skeleton appears
- [ ] Empty state appears for no results
- [ ] Error state has retry
- [ ] Logout returns to login
- [ ] `flutter analyze` has 0 errors
- [ ] Android debug APK builds

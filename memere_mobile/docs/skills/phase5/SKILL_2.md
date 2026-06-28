# phase5/SKILL_2.md — Mock Exam Catalog, Timed Attempt Screen & Autosave
# Memere Mobile (memere_mobile) — Phase 5, Part 2
# READ SKILL.md → phase5/SKILL_1.md → then this file.

---

## OBJECTIVE

Build the mock exam taking experience:
exam catalog route → filters/search → exam start screen → live timed exam attempt →
question navigation → answer state → autosave → submit confirmation.

By the end of this skill, a student can start a published mock exam, answer
questions under a server-synced timer display, autosave progress, and submit.

---

## TIMER RULES

- Backend `expires_at` and `remaining_seconds` are authoritative.
- UI countdown is display-only.
- When local timer reaches zero, submit or show expired state, but backend still decides.
- Never extend or reset the timer locally.
- Save progress before submit when possible.

---

## PART A — ROUTING UPDATE

### FILE A1 — `lib/core/router/app_router.dart`

Add routes:

```dart
static const mockExams = '/mock-exams';
static const examAttempt = '/exam-attempts/:attemptId';

static String examAttemptPath({
  required String attemptId,
  required String examId,
}) {
  return '/exam-attempts/$attemptId?examId=$examId';
}
```

Route builders:

- `/mock-exams` → `MockExamCatalogScreen`
- `/exam-attempts/:attemptId?examId=...` → `ExamAttemptScreen`

Keep both routes behind auth.

Add a way to reach `/mock-exams` from the home/catalog area:

- bottom navigation item, or
- toolbar action, or
- dedicated section on home

Do not bury mock exams behind course detail only.

---

## PART B — CATALOG STATE PROVIDER

### FILE B1 — `lib/features/exam/presentation/providers/mock_exam_catalog_provider.dart`

State shape:

```dart
class MockExamCatalogState {
  const MockExamCatalogState({
    this.exams = const [],
    this.filteredExams = const [],
    this.nextCursor,
    this.selectedSubject,
    this.selectedGrade = 12,
    this.searchQuery = '',
    this.isLoadingMore = false,
  });
}
```

Methods:

- `Future<MockExamCatalogState> build()`
- `refresh()`
- `loadMore()`
- `setSearchQuery(String value)`
- `setSubject(String? subject)`
- `setGrade(int? grade)`
- `clearFilters()`

Backend filters:

- `subject`
- `grade`
- `limit`
- `after`

Local search:

- title
- subject
- instructions

---

## PART C — MOCK EXAM CATALOG SCREEN

### FILE C1 — `lib/features/exam/presentation/screens/mock_exam_catalog_screen.dart`

Required UI:

1. safe-area scaffold
2. top app bar/title: `Mock exams`
3. subtitle: `Practice under real exam timing`
4. search field
5. subject chips
6. grade filter
7. exam cards
8. loading skeleton
9. empty state
10. error state with retry
11. pagination/load more

Use Phase 1 design system. Keep screen dense and study-focused.

---

## PART D — CATALOG WIDGETS

### FILE D1 — `lib/features/exam/presentation/widgets/mock_exam_card.dart`

Show:

- title
- subject
- grade
- duration
- total marks
- pass marks
- instructions preview
- start button

Tap/start behavior:

- show start confirmation
- call `StartExamUseCase`
- navigate to attempt route on success

---

### FILE D2 — `lib/features/exam/presentation/widgets/mock_exam_catalog_skeleton.dart`

Use shimmer with stable card heights.

---

### FILE D3 — `lib/features/exam/presentation/widgets/exam_empty_state.dart`

Reusable empty/error state.

---

## PART E — START CONFIRMATION

Before starting an exam, show a bottom sheet:

- duration
- total marks
- pass marks
- instructions
- warning: timer starts immediately
- cancel
- start exam

Copy:

```text
The timer starts when you begin. Keep the app open and submit before time runs out.
```

---

## PART F — EXAM ATTEMPT PROVIDER

### FILE F1 — `lib/features/exam/presentation/providers/exam_attempt_provider.dart`

State shape:

```dart
class ExamAttemptState {
  const ExamAttemptState({
    required this.attempt,
    this.currentIndex = 0,
    this.answers = const {},
    this.isSaving = false,
    this.lastSavedAt,
    this.saveError,
    this.isSubmitting = false,
    this.hasUnsavedChanges = false,
  });
}
```

Methods:

- `selectSingleAnswer(questionId, answerId)`
- `toggleMultiAnswer(questionId, answerId)`
- `setShortAnswer(questionId, value)`
- `goToQuestion(index)`
- `nextQuestion()`
- `previousQuestion()`
- `saveProgress()`
- `submit()`

Autosave:

- save every 15-30 seconds while dirty
- save before leaving if possible
- save before submit
- do not block navigation on autosave failure

---

## PART G — EXAM TIMER PROVIDER

### FILE G1 — `lib/features/exam/presentation/providers/exam_timer_provider.dart`

Use `expiresAt` when present.

Rules:

- calculate remaining display from server expiry
- fallback to initial `remainingSeconds`
- warning state below 5 minutes
- critical state below 60 seconds
- at zero, trigger submit flow or show expired state

Do not extend time locally.

---

## PART H — EXAM ATTEMPT SCREEN

### FILE H1 — `lib/features/exam/presentation/screens/exam_attempt_screen.dart`

Constructor:

```dart
class ExamAttemptScreen extends ConsumerWidget {
  const ExamAttemptScreen({
    super.key,
    required this.attemptId,
    required this.examId,
  });

  final String attemptId;
  final String examId;
}
```

Required UI:

1. locked top bar with timer
2. question progress: `Question 12 of 100`
3. marks badge
4. question text
5. answer controls
6. question palette/navigator
7. previous/next
8. submit button
9. autosave status

Back behavior:

- warn user that exam timer continues server-side
- save progress before leaving when possible

---

## PART I — ATTEMPT WIDGETS

### FILE I1 — `lib/features/exam/presentation/widgets/exam_question_card.dart`

Show:

- question text
- marks
- subject/topic chips
- answer control based on question type

No correct/wrong styling during attempt.

---

### FILE I2 — `lib/features/exam/presentation/widgets/exam_answer_option_tile.dart`

Support:

- single choice
- true/false
- multi select

---

### FILE I3 — `lib/features/exam/presentation/widgets/exam_question_palette.dart`

Show grid of question numbers.

States:

- current
- answered
- unanswered
- marked for review if implemented

Do not show correct/wrong state.

---

### FILE I4 — `lib/features/exam/presentation/widgets/exam_timer_badge.dart`

Display:

- `h:mm:ss`

Warning colors:

- normal: text secondary/accent
- under 5 minutes: warning
- under 60 seconds: error

---

### FILE I5 — `lib/features/exam/presentation/widgets/exam_save_status.dart`

Show:

- `Saved`
- `Saving...`
- `Offline - will retry` or `Save failed`

Keep compact.

---

## PART J — SUBMIT CONFIRMATION

Before submit, show:

- answered count
- unanswered count
- remaining time
- final submission warning
- cancel
- submit

On submit:

1. save progress
2. call submit use case
3. navigate to result route

If submit fails:

- show snackbar
- keep answers intact
- let user retry

---

## PART K — EXPIRED STATE

When timer reaches zero:

- call submit automatically if answers exist and backend accepts it
- otherwise show expired state and submit button disabled
- backend remains final authority

Copy:

```text
Time is up. Submitting your saved answers.
```

---

## VALIDATION

Run:

```bash
dart format lib/core/router lib/features/exam
flutter analyze
flutter build apk --debug
```

Manual test:

- open mock exam catalog
- filter by subject
- start exam
- answer questions
- navigate with palette
- autosave fires
- timer displays countdown
- submit confirmation appears
- submit navigates to result

---

## SKILL_2 CHECKLIST

- [ ] `/mock-exams` route exists
- [ ] Mock exam catalog loads backend data
- [ ] Subject/grade filters work
- [ ] Local search works
- [ ] Start confirmation warns timer starts immediately
- [ ] Start exam calls `/mock-exams/:id/start`
- [ ] Attempt screen renders live questions
- [ ] Answer selection/input works
- [ ] Timer display uses backend timing
- [ ] Autosave calls `PATCH /exam-attempts/:id`
- [ ] Submit confirmation works
- [ ] No correct/wrong state appears before submit
- [ ] `flutter analyze` has 0 errors
- [ ] Android debug APK builds

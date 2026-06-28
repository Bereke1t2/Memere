# phase4/SKILL_2.md — Quiz Attempt Screen, Timer, Answer Flow & Autosave
# Memere Mobile (memere_mobile) — Phase 4, Part 2
# READ SKILL.md → phase4/SKILL_1.md → then this file.

---

## OBJECTIVE

Build the live quiz-taking experience:
quiz entry route → metadata screen → start attempt → question flow →
answer selection/input → local answer state → autosave → submit confirmation.

By the end of this skill, a student can start a quiz attempt, answer questions,
autosave progress, and submit to the backend.

---

## SECURITY AND TIMER RULES

- The UI timer is display-only.
- Backend `expires_at` and `remaining_seconds` are authoritative.
- Do not reveal correct answers during the attempt.
- Do not show explanations during the attempt.
- Do not score locally.
- Disable local submit only for obvious invalid states, but let backend enforce final rules.

---

## PART A — ROUTING UPDATE

### FILE A1 — `lib/core/router/app_router.dart`

Add routes:

```dart
static const quizDetail = '/quizzes/:quizId';
static const quizAttempt = '/quiz-attempts/:attemptId';

static String quizDetailPath(String quizId) => '/quizzes/$quizId';

static String quizAttemptPath({
  required String attemptId,
  required String quizId,
}) {
  return '/quiz-attempts/$attemptId?quizId=$quizId';
}
```

Route builders:

- `/quizzes/:quizId` → `QuizDetailScreen`
- `/quiz-attempts/:attemptId?quizId=...` → `QuizAttemptScreen`

Keep both routes behind auth.

---

## PART B — LESSON/COURSE ENTRY POINT

Phase 4 needs a way to discover quiz IDs.

Preferred path:

- lessons with `type == quiz` expose quiz metadata or `quizId`
- tapping quiz lesson navigates to `QuizDetailScreen`

If backend/course DTOs do not expose quiz IDs yet:

- show a clear placeholder on quiz lesson tap
- do not hardcode quiz IDs
- document the backend DTO gap

Placeholder:

```text
This lesson does not have a quiz attached yet.
```

---

## PART C — QUIZ DETAIL PROVIDER

### FILE C1 — `lib/features/quiz/presentation/providers/quiz_detail_provider.dart`

Create:

```dart
final quizDetailProvider = FutureProvider.family<QuizEntity, String>(
  (ref, quizId) async {
    final useCase = ref.watch(getQuizUseCaseProvider);
    final result = await useCase(quizId);
    return result.fold((failure) => throw failure, (quiz) => quiz);
  },
);
```

---

## PART D — QUIZ DETAIL SCREEN

### FILE D1 — `lib/features/quiz/presentation/screens/quiz_detail_screen.dart`

Constructor:

```dart
class QuizDetailScreen extends ConsumerWidget {
  const QuizDetailScreen({
    super.key,
    required this.quizId,
  });

  final String quizId;
}
```

Required UI:

- quiz title
- question count
- pass percentage
- time limit
- attempts used/max attempts
- start button
- loading/error states

Start action:

1. call `StartQuizAttemptUseCase`
2. on success navigate to attempt screen
3. on failure show snackbar

Do not load questions with `GET /quizzes/:id`; only `POST /quizzes/:id/attempts`
returns live questions.

---

## PART E — ATTEMPT STATE PROVIDER

### FILE E1 — `lib/features/quiz/presentation/providers/quiz_attempt_provider.dart`

Create an `AsyncNotifierProvider.family`.

State shape:

```dart
class QuizAttemptState {
  const QuizAttemptState({
    required this.attempt,
    this.currentIndex = 0,
    this.answers = const {},
    this.isSaving = false,
    this.lastSavedAt,
    this.saveError,
    this.isSubmitting = false,
  });
}
```

Required methods:

- `selectSingleAnswer(questionId, answerId)`
- `toggleMultiAnswer(questionId, answerId)`
- `setShortAnswer(questionId, value)`
- `goToQuestion(index)`
- `nextQuestion()`
- `previousQuestion()`
- `saveProgress()`
- `submit()`

Autosave:

- debounce answer changes
- save every 15-30 seconds while there are changes
- save before submit
- do not block question navigation on save failure

---

## PART F — QUIZ ATTEMPT SCREEN

### FILE F1 — `lib/features/quiz/presentation/screens/quiz_attempt_screen.dart`

Constructor:

```dart
class QuizAttemptScreen extends ConsumerWidget {
  const QuizAttemptScreen({
    super.key,
    required this.attemptId,
    required this.quizId,
  });

  final String attemptId;
  final String quizId;
}
```

Required UI:

1. top bar with close/back confirmation
2. timer display when `remainingSeconds`/`expiresAt` exists
3. progress indicator: `Question 3 of 10`
4. question card
5. answer options/input
6. previous/next buttons
7. question navigator
8. submit button on final question
9. save status indicator

Back/close behavior:

- if answers changed, show confirmation sheet
- save progress before leaving when possible

---

## PART G — TIMER PROVIDER/WIDGET

### FILE G1 — `lib/features/quiz/presentation/providers/quiz_timer_provider.dart`

Use backend timing values to create a display countdown.

Rules:

- calculate countdown from `expiresAt` when present
- fallback to `remainingSeconds`
- when display reaches zero, call submit or show expired state
- backend still validates expiry

Do not extend time locally.

---

## PART H — ATTEMPT WIDGETS

### FILE H1 — `lib/features/quiz/presentation/widgets/quiz_question_card.dart`

Props:

- `QuizQuestionEntity question`
- `Object? selectedAnswer`
- callbacks for answer changes

Show:

- question text
- points
- subject/topic chips when available
- answer control based on question type

---

### FILE H2 — `lib/features/quiz/presentation/widgets/answer_option_tile.dart`

Support:

- single choice
- true/false
- multi select

No correct/incorrect styling during attempt.

Selected state should be visually clear.

---

### FILE H3 — `lib/features/quiz/presentation/widgets/short_answer_field.dart`

For short-answer questions.

Debounce text updates through provider.

---

### FILE H4 — `lib/features/quiz/presentation/widgets/question_navigator.dart`

Show question dots/numbers.

States:

- current
- answered
- unanswered

Do not show correct/wrong state.

---

### FILE H5 — `lib/features/quiz/presentation/widgets/quiz_timer_badge.dart`

Display:

- mm:ss for under 1 hour
- h:mm:ss for longer

Use warning color when under 60 seconds.

---

## PART I — SUBMIT CONFIRMATION

Before submit, show a bottom sheet/dialog:

- answered count
- unanswered count
- warning that submission is final
- cancel
- submit

On submit:

1. save progress
2. call submit use case
3. navigate to result screen with attempt ID

If submit returns result directly, pass/store it via provider or navigate and let
result screen refetch using `GET /quiz-attempts/:id/result`.

Preferred: navigate by attempt ID and let result screen fetch/refresh.

---

## PART J — ERROR HANDLING

Failure messages:

```dart
final message = error is Failure
    ? error.message
    : 'Could not load quiz. Please try again.';
```

Common states:

- attempts exhausted
- quiz unavailable
- attempt expired
- network error during autosave
- submit failed

Autosave failure should show a subtle status, not a blocking dialog.

Submit failure should show snackbar and keep answers intact.

---

## VALIDATION

Run:

```bash
dart format lib/core/router lib/features/quiz lib/features/courses
flutter analyze
flutter build apk --debug
```

Manual test:

- open quiz detail
- start attempt
- answer single-choice question
- navigate between questions
- autosave fires
- submit confirmation appears
- submit returns result or navigates to result route

---

## SKILL_2 CHECKLIST

- [ ] Quiz routes exist
- [ ] Quiz detail screen loads metadata only
- [ ] Start attempt returns live questions
- [ ] Attempt state stores answers locally
- [ ] Single-choice answers work
- [ ] Multi-select answers work when present
- [ ] Short-answer input works when present
- [ ] Timer display is based on backend timing
- [ ] Autosave calls `PATCH /quiz-attempts/:id`
- [ ] Submit calls backend and keeps answers on failure
- [ ] No correct/wrong state appears before submit
- [ ] `flutter analyze` has 0 errors
- [ ] Android debug APK builds

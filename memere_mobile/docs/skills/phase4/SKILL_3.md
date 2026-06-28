# phase4/SKILL_3.md — Quiz Results, Feedback Review & Phase 4 Checklist
# Memere Mobile (memere_mobile) — Phase 4, Part 3
# READ SKILL.md → phase4/SKILL_1.md → phase4/SKILL_2.md → then this file.

---

## OBJECTIVE

Build the post-submission quiz result experience:
result route → fetch graded result → score summary → pass/fail state →
subject breakdown → per-question feedback → review answers.

By the end of this skill, a submitted quiz attempt shows server-graded results
and feedback without any client-side grading.

---

## RESULT SECURITY RULES

- Correct answers and explanations are allowed only on the result screen.
- Result screen must use backend response only.
- Do not infer correctness locally by comparing selected and correct IDs.
- Do not mutate result score locally.
- Treat backend result as immutable display data.

---

## PART A — RESULT ROUTING

### FILE A1 — `lib/core/router/app_router.dart`

Add:

```dart
static const quizResult = '/quiz-attempts/:attemptId/result';

static String quizResultPath(String attemptId) {
  return '/quiz-attempts/$attemptId/result';
}
```

Route builder:

```dart
GoRoute(
  path: AppRoutes.quizResult,
  builder: (_, state) {
    final attemptId = state.pathParameters['attemptId']!;
    return QuizResultScreen(attemptId: attemptId);
  },
),
```

Keep behind auth.

---

## PART B — RESULT PROVIDER

### FILE B1 — `lib/features/quiz/presentation/providers/quiz_result_provider.dart`

Create:

```dart
final quizResultProvider = FutureProvider.family<QuizResultEntity, String>(
  (ref, attemptId) async {
    final useCase = ref.watch(getQuizResultUseCaseProvider);
    final result = await useCase(attemptId);
    return result.fold((failure) => throw failure, (quizResult) => quizResult);
  },
);
```

Retry:

```dart
ref.invalidate(quizResultProvider(attemptId));
```

---

## PART C — RESULT SCREEN

### FILE C1 — `lib/features/quiz/presentation/screens/quiz_result_screen.dart`

Constructor:

```dart
class QuizResultScreen extends ConsumerWidget {
  const QuizResultScreen({
    super.key,
    required this.attemptId,
  });

  final String attemptId;
}
```

Required UI states:

- loading skeleton
- error with retry
- result loaded

Loaded layout:

1. top bar with close/back
2. pass/fail summary
3. percentage score
4. points earned / total points
5. attempt number
6. submitted time
7. subject breakdown
8. question feedback list
9. actions:
   - back to course
   - retry quiz when allowed later

Retry quiz behavior can be a placeholder if attempt limit metadata is not
available on result.

---

## PART D — RESULT WIDGETS

### FILE D1 — `lib/features/quiz/presentation/widgets/quiz_score_summary.dart`

Props:

- `QuizResultEntity result`

Show:

- pass/fail icon
- percentage
- score/total points
- concise status text

Use:

- success color when passed
- error/warning color when failed

---

### FILE D2 — `lib/features/quiz/presentation/widgets/subject_breakdown_card.dart`

Props:

- `Map<String, SubjectScoreEntity> subjectBreakdown`

Show each subject:

- subject name
- earned/possible
- percentage bar

If empty, hide the section.

---

### FILE D3 — `lib/features/quiz/presentation/widgets/question_feedback_tile.dart`

Props:

- `QuestionFeedbackEntity feedback`
- optional matching `QuizQuestionEntity` if available

Show:

- correct/incorrect status
- points awarded/possible
- selected answer IDs or selected answer text when available
- correct answer IDs
- explanation when present

If question text is unavailable from result, show:

```text
Question <number>
```

Do not fabricate question text.

---

### FILE D4 — `lib/features/quiz/presentation/widgets/quiz_result_skeleton.dart`

Use `shimmer`.

Skeleton areas:

- score summary block
- subject breakdown rows
- 4 feedback rows

---

## PART E — PASS/FAIL COPY

Use concise copy:

Passed:

```text
Passed
Good work. Review the explanations before moving on.
```

Failed:

```text
Keep practicing
Review the missed questions and try again when ready.
```

Avoid shame-heavy language.

---

## PART F — ATTEMPT SCREEN INTEGRATION

After successful submit:

```dart
context.go(AppRoutes.quizResultPath(result.attemptId));
```

If using result refetch:

- ignore direct result payload after navigation, or
- seed provider cache only if the local Riverpod setup makes that clean

Prefer simple refetch by attempt ID for reliability.

---

## PART G — COURSE/LESSON INTEGRATION

Phase 4 should update quiz lesson tap behavior:

- if lesson has `quizId`, navigate to `QuizDetailScreen`
- if lesson has no `quizId`, show placeholder
- if lesson type is video, keep Phase 3 behavior
- if lesson type is note/mixed, keep placeholder until later phase

Do not hardcode quiz IDs.

---

## PART H — OFFLINE/RESUME BEHAVIOR

Minimum Phase 4:

- preserve in-memory answers during active attempt
- autosave to backend regularly
- show autosave failure state

Optional:

- store draft answers locally with SharedPreferences/Hive
- clear draft after successful submit

If local draft storage is added:

- key by `attemptId`
- do not store correct answers
- clear when attempt submits/expires

---

## PART I — TESTING NOTES

Widget tests should cover:

- live answer option does not expose correct state
- submit confirmation counts unanswered questions
- result screen displays backend percentage/pass state
- feedback appears only on result screen

Unit tests should cover:

- answer payload building
- timer display formatting
- model JSON parsing

---

## VALIDATION

Run:

```bash
dart format lib/core/router lib/features/quiz lib/features/courses
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

- open quiz detail
- start attempt
- answer some questions
- autosave succeeds
- submit attempt
- result loads
- feedback is visible only after submit
- no answer key appears during attempt

---

## PHASE 4 FINAL CHECKLIST

### SKILL_1 complete when:

- [ ] Quiz entities and models compile
- [ ] Live answer models have no answer-key field
- [ ] Repository/use cases cover quiz lifecycle
- [ ] Remote datasource calls all quiz endpoints
- [ ] Failure mapping is consistent

### SKILL_2 complete when:

- [ ] Quiz detail screen loads metadata
- [ ] Start attempt loads live questions
- [ ] Answer selection/input works
- [ ] Timer display uses backend timing
- [ ] Autosave works
- [ ] Submit confirmation works
- [ ] Submit sends answers to backend
- [ ] No live attempt UI reveals correct answers

### SKILL_3 complete when:

- [ ] Result route exists
- [ ] Result screen fetches backend result
- [ ] Score summary displays pass/fail
- [ ] Subject breakdown displays
- [ ] Feedback list displays selected/correct answers and explanation
- [ ] Course/lesson quiz tap integration handles missing quiz IDs honestly
- [ ] `flutter analyze` has 0 errors
- [ ] Android debug APK builds

---

## PHASE 4 → PHASE 5 HANDOFF

**Phase 4 is complete when all 3 SKILL files are done and the checklist above passes.**

To start Phase 5, tell Antigravity:

```text
Phase 4 is complete. Read SKILL.md and all phase5 skill files.
We are starting Phase 5: Mock Exam Engine.
Reference: memere_mobile/docs/memere_Design_Specification.md
```

**What Phase 5 will build:**

- mock exam catalog
- timed exam attempt start
- server-synced exam timer
- exam answer autosave
- exam submission
- exam result and analytics screens

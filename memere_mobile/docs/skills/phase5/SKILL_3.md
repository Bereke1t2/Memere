# phase5/SKILL_3.md — Mock Exam Results, Analytics & Phase 5 Checklist
# Memere Mobile (memere_mobile) — Phase 5, Part 3
# READ SKILL.md → phase5/SKILL_1.md → phase5/SKILL_2.md → then this file.

---

## OBJECTIVE

Build the post-exam experience:
result route → fetch server-graded result → score/pass summary →
subject breakdown → weak areas analytics → percentile → question feedback review.

By the end of this skill, a submitted mock exam attempt shows authoritative
results and useful analytics without any client-side grading.

---

## RESULT SECURITY RULES

- Correct answers/explanations are result-only.
- Result screen must use backend response only.
- Analytics screen must use backend response only.
- Do not infer correctness locally.
- Do not mutate score, pass/fail, percentile, or weak areas locally.

---

## PART A — RESULT ROUTING

### FILE A1 — `lib/core/router/app_router.dart`

Add routes:

```dart
static const examResult = '/exam-attempts/:attemptId/results';
static const examAnalytics = '/exam-attempts/:attemptId/analytics';

static String examResultPath(String attemptId) {
  return '/exam-attempts/$attemptId/results';
}

static String examAnalyticsPath(String attemptId) {
  return '/exam-attempts/$attemptId/analytics';
}
```

Route builders:

- `/exam-attempts/:attemptId/results` → `ExamResultScreen`
- `/exam-attempts/:attemptId/analytics` → `ExamAnalyticsScreen`

Keep behind auth.

---

## PART B — RESULT PROVIDERS

### FILE B1 — `lib/features/exam/presentation/providers/exam_result_provider.dart`

Create:

```dart
final examResultProvider = FutureProvider.family<ExamResultEntity, String>(
  (ref, attemptId) async {
    final useCase = ref.watch(getExamResultUseCaseProvider);
    final result = await useCase(attemptId);
    return result.fold((failure) => throw failure, (examResult) => examResult);
  },
);
```

### FILE B2 — `lib/features/exam/presentation/providers/exam_analytics_provider.dart`

Create:

```dart
final examAnalyticsProvider =
    FutureProvider.family<ExamAttemptAnalyticsEntity, String>(
  (ref, attemptId) async {
    final useCase = ref.watch(getExamAnalyticsUseCaseProvider);
    final result = await useCase(attemptId);
    return result.fold((failure) => throw failure, (analytics) => analytics);
  },
);
```

---

## PART C — RESULT SCREEN

### FILE C1 — `lib/features/exam/presentation/screens/exam_result_screen.dart`

Constructor:

```dart
class ExamResultScreen extends ConsumerWidget {
  const ExamResultScreen({
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

1. top bar
2. pass/fail summary
3. percentage score
4. marks earned / total marks
5. pass marks
6. submitted time
7. subject breakdown
8. analytics preview button
9. question feedback list
10. actions:
   - back to mock exams
   - review analytics

---

## PART D — ANALYTICS SCREEN

### FILE D1 — `lib/features/exam/presentation/screens/exam_analytics_screen.dart`

Constructor:

```dart
class ExamAnalyticsScreen extends ConsumerWidget {
  const ExamAnalyticsScreen({
    super.key,
    required this.attemptId,
  });

  final String attemptId;
}
```

Show:

- score
- percentage
- percentile if present
- subject breakdown
- weak areas
- recommendation placeholder

Recommendation placeholder copy:

```text
Focus revision on the weak areas above before your next attempt.
```

Do not implement AI/adaptive learning in Phase 5.

---

## PART E — RESULT WIDGETS

### FILE E1 — `lib/features/exam/presentation/widgets/exam_score_summary.dart`

Props:

- `ExamResultEntity result`

Show:

- pass/fail icon
- percentage
- score/total marks
- pass marks

---

### FILE E2 — `lib/features/exam/presentation/widgets/exam_subject_breakdown.dart`

Props:

- subject breakdown map/list

Show:

- subject
- earned/possible
- progress bar

---

### FILE E3 — `lib/features/exam/presentation/widgets/exam_question_feedback_tile.dart`

Props:

- `ExamQuestionFeedbackEntity feedback`
- optional question text if available

Show:

- correct/incorrect
- marks awarded/possible
- selected answer IDs or text if available
- correct answer IDs
- explanation

If question text is unavailable:

```text
Question <number>
```

Do not fabricate text.

---

### FILE E4 — `lib/features/exam/presentation/widgets/weak_area_card.dart`

Props:

- `ExamSubjectScoreEntity weakArea`

Show:

- topic/subject key
- earned/possible
- percentage
- visual severity:
  - under 40%: error
  - under 70%: warning
  - otherwise: normal

---

### FILE E5 — `lib/features/exam/presentation/widgets/percentile_card.dart`

Props:

- `double? percentile`

If null, hide or show:

```text
Percentile will appear after more students attempt this exam.
```

---

### FILE E6 — `lib/features/exam/presentation/widgets/exam_result_skeleton.dart`

Use shimmer.

Skeleton:

- score summary block
- subject rows
- feedback rows

---

## PART F — RESULT COPY

Passed:

```text
Passed
You reached the pass mark. Review the feedback to strengthen weak areas.
```

Failed:

```text
Keep practicing
Review weak areas and try another mock exam when ready.
```

Use direct, respectful language.

---

## PART G — ATTEMPT SCREEN INTEGRATION

After successful submit:

```dart
context.go(AppRoutes.examResultPath(result.attemptId));
```

Preferred:

- navigate by attempt ID
- result screen fetches result from backend

This avoids keeping large result payloads in navigation state.

---

## PART H — OPTIONAL LOCAL DRAFT RECOVERY

Minimum Phase 5:

- autosave to backend
- keep answers in memory during active attempt

Optional:

- store local draft answers with Hive/SharedPreferences keyed by `attemptId`
- clear local draft after successful submit

Do not store correct answers locally.

---

## PART I — TESTING NOTES

Widget tests:

- live exam answer option has no correct state
- timer warning color changes under thresholds
- submit confirmation counts unanswered questions
- result screen displays backend score/pass state
- analytics screen displays weak areas

Unit tests:

- model parsing
- answer payload building
- timer formatting
- subject score percentage helper

---

## VALIDATION

Run:

```bash
dart format lib/core/router lib/features/exam
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

- open mock exam catalog
- start mock exam
- answer questions
- autosave succeeds
- submit
- result loads
- analytics loads
- feedback appears only after submit
- no answer key appears during attempt

---

## PHASE 5 FINAL CHECKLIST

### SKILL_1 complete when:

- [ ] Mock exam entities/models compile
- [ ] Live exam models have no answer-key field
- [ ] Repository/use cases cover catalog/start/save/submit/result/analytics
- [ ] Remote datasource calls exam endpoints
- [ ] Failure mapping is consistent

### SKILL_2 complete when:

- [ ] Mock exam catalog loads and filters
- [ ] Start confirmation works
- [ ] Start exam returns live questions
- [ ] Attempt screen renders timer, questions, answers, palette
- [ ] Autosave works
- [ ] Submit confirmation works
- [ ] No live attempt UI reveals correct answers

### SKILL_3 complete when:

- [ ] Result route exists
- [ ] Result screen fetches backend result
- [ ] Analytics route exists
- [ ] Analytics screen fetches backend analytics
- [ ] Score summary displays pass/fail
- [ ] Subject breakdown and weak areas display
- [ ] Feedback list displays only post-submit
- [ ] `flutter analyze` has 0 errors
- [ ] Android debug APK builds

---

## PHASE 5 → PHASE 6 HANDOFF

**Phase 5 is complete when all 3 SKILL files are done and the checklist above passes.**

To start Phase 6, tell Antigravity:

```text
Phase 5 is complete. Read SKILL.md and all phase6 skill files.
We are starting Phase 6: Payment & Enrollment.
Reference: memere_mobile/docs/memere_Design_Specification.md
```

**What Phase 6 will build:**

- free course enrollment
- payment initiation
- Chapa/Telebirr/Stripe WebView flow
- payment status polling
- enrollment confirmation
- purchase history

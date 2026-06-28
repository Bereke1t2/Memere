# phase4/SKILL_1.md — Quiz Domain, Data Layer & Backend Contract
# Memere Mobile (memere_mobile) — Phase 4, Part 1
# READ SKILL.md → all phase1 skill files → all phase2 skill files → all phase3 skill files → then this file.

---

## OBJECTIVE

Build the Quiz Engine foundation:
quiz metadata entities → live attempt entities → answer payloads → result entities →
repository contract → use cases → backend DTO models → remote datasource →
repository implementation → Riverpod dependency providers.

By the end of this skill, the app can list lesson/course quizzes, start attempts,
autosave answers, submit attempts, and fetch graded results from the backend.
Screens are built in `SKILL_2.md` and `SKILL_3.md`.

---

## PHASE 4 BOUNDARY

Phase 4 builds **lesson/course quizzes only**:

- quiz metadata
- start quiz attempt
- live question flow
- answer selection/input
- autosave progress
- submit attempt
- server-graded result screen

Do **not** build these yet:

- full mock exam catalog/timed exam engine
- payment/enrollment
- leaderboard
- advanced dashboard analytics
- teacher quiz builder UI

Teacher authoring endpoints exist but are Phase 9 for mobile.

---

## SECURITY RULES

These are non-negotiable:

- Correct answers are never available during live attempts.
- Do not add `isCorrect` to live answer models.
- Do not grade on the client.
- Do not calculate pass/fail locally.
- Only show correct answers/explanations from `/quiz-attempts/:id/result` or submit response.
- Do not log answer payloads in production.

The client may track selected answer IDs for UI state, but the backend is the
source of truth for scores.

---

## BACKEND CONTRACT

Student quiz routes:

```text
GET /api/v1/courses/:id/quizzes
GET /api/v1/quizzes/:id
POST /api/v1/quizzes/:id/attempts
PATCH /api/v1/quiz-attempts/:id
POST /api/v1/quiz-attempts/:id/submit
GET /api/v1/quiz-attempts/:id/result
```

Important:

- `GET /courses/:id/quizzes` is currently teacher/admin-gated in router. For student Phase 4, either expose a student-safe course quiz list route or enter quizzes from lesson metadata once available.
- `GET /quizzes/:id` returns metadata and question count only.
- `POST /quizzes/:id/attempts` returns live questions and answer options with no answer keys.
- `PATCH /quiz-attempts/:id` autosaves answers and returns `204`.
- `POST /quiz-attempts/:id/submit` returns graded result.

---

## RESPONSE SHAPES

### `GET /quizzes/:id`

```json
{
  "id": "uuid",
  "course_id": "uuid",
  "title": "Algebra Basics Quiz",
  "time_limit_seconds": 600,
  "pass_percentage": 70,
  "randomize_questions": true,
  "max_attempts": 3,
  "question_count": 10,
  "attempts_used": 1
}
```

### `POST /quizzes/:id/attempts`

```json
{
  "attempt_id": "uuid",
  "quiz_id": "uuid",
  "attempt_number": 2,
  "status": "in_progress",
  "started_at": "2026-01-01T00:00:00Z",
  "expires_at": "2026-01-01T00:10:00Z",
  "remaining_seconds": 600,
  "questions": [
    {
      "id": "uuid",
      "text": "Solve x + 2 = 5",
      "type": "multiple_choice",
      "points": 1,
      "subject": "Mathematics",
      "topic": "Algebra",
      "answers": [
        {
          "id": "uuid",
          "text": "3",
          "order_index": 1
        }
      ]
    }
  ]
}
```

### `PATCH /quiz-attempts/:id`

Request:

```json
{
  "answers": {
    "question_uuid": "answer_uuid"
  }
}
```

Response: `204 No Content`

### `POST /quiz-attempts/:id/submit`

Request:

```json
{
  "answers": {
    "question_uuid": "answer_uuid"
  }
}
```

Response: same shape as result below.

### `GET /quiz-attempts/:id/result`

```json
{
  "attempt_id": "uuid",
  "quiz_id": "uuid",
  "attempt_number": 2,
  "status": "submitted",
  "score": 8,
  "total_points": 10,
  "percentage": 80,
  "passed": true,
  "submitted_at": "2026-01-01T00:09:00Z",
  "feedback": [
    {
      "question_id": "uuid",
      "correct": true,
      "points_awarded": 1,
      "points_possible": 1,
      "selected_answers": ["uuid"],
      "correct_answer_ids": ["uuid"],
      "explanation": "x = 3"
    }
  ],
  "subject_breakdown": {
    "Mathematics": {
      "earned": 8,
      "possible": 10
    }
  }
}
```

---

## REQUIRED FOLDER STRUCTURE

Create or verify:

```bash
mkdir -p lib/features/quiz/data/datasources
mkdir -p lib/features/quiz/data/models
mkdir -p lib/features/quiz/data/repositories
mkdir -p lib/features/quiz/domain/entities
mkdir -p lib/features/quiz/domain/repositories
mkdir -p lib/features/quiz/domain/usecases
mkdir -p lib/features/quiz/presentation/providers
mkdir -p lib/features/quiz/presentation/screens
mkdir -p lib/features/quiz/presentation/widgets
```

---

## PART A — DOMAIN ENTITIES

### FILE A1 — `lib/features/quiz/domain/entities/quiz_entity.dart`

Fields:

- `id`
- `courseId`
- `title`
- `timeLimitSeconds`
- `passPercentage`
- `randomizeQuestions`
- `maxAttempts`
- `questionCount`
- `attemptsUsed`

Helpers:

- `bool get hasTimeLimit`
- `bool get hasAttemptsRemaining`
- `int? get attemptsRemaining`

---

### FILE A2 — `lib/features/quiz/domain/entities/quiz_attempt_entity.dart`

Fields:

- `attemptId`
- `quizId`
- `attemptNumber`
- `status`
- `startedAt`
- `expiresAt`
- `remainingSeconds`
- `questions`

Use enum:

```dart
enum QuizAttemptStatus { inProgress, submitted, expired, graded }
```

---

### FILE A3 — `lib/features/quiz/domain/entities/quiz_question_entity.dart`

Fields:

- `id`
- `text`
- `type`
- `points`
- `subject`
- `topic`
- `answers`

Use enum:

```dart
enum QuizQuestionType { multipleChoice, multipleSelect, trueFalse, shortAnswer }
```

No answer key fields.

---

### FILE A4 — `lib/features/quiz/domain/entities/quiz_answer_entity.dart`

Fields:

- `id`
- `text`
- `orderIndex`

Do not add `isCorrect`.

---

### FILE A5 — `lib/features/quiz/domain/entities/quiz_result_entity.dart`

Fields:

- `attemptId`
- `quizId`
- `attemptNumber`
- `status`
- `score`
- `totalPoints`
- `percentage`
- `passed`
- `submittedAt`
- `feedback`
- `subjectBreakdown`

---

### FILE A6 — `lib/features/quiz/domain/entities/question_feedback_entity.dart`

Fields:

- `questionId`
- `correct`
- `pointsAwarded`
- `pointsPossible`
- `selectedAnswers`
- `correctAnswerIds`
- `explanation`

This entity is result-only.

---

### FILE A7 — `lib/features/quiz/domain/entities/subject_score_entity.dart`

Fields:

- `earned`
- `possible`

Helper:

- `double get percentage`

---

## PART B — ANSWER PAYLOAD

### FILE B1 — `lib/features/quiz/domain/entities/quiz_answer_payload.dart`

Represent answers as:

```dart
typedef QuizAnswerPayload = Map<String, Object>;
```

Rules:

- single choice: `questionId -> answerId`
- multi select: `questionId -> List<String>`
- short answer: `questionId -> String`

Keep this flexible because backend accepts `map[string]any`.

---

## PART C — REPOSITORY CONTRACT

### FILE C1 — `lib/features/quiz/domain/repositories/quiz_repository.dart`

Use `Either<Failure, T>`.

Methods:

```dart
Future<Either<Failure, List<QuizEntity>>> listQuizzesByCourse(String courseId);

Future<Either<Failure, QuizEntity>> getQuiz(String quizId);

Future<Either<Failure, QuizAttemptEntity>> startAttempt(String quizId);

Future<Either<Failure, void>> saveProgress({
  required String attemptId,
  required QuizAnswerPayload answers,
});

Future<Either<Failure, QuizResultEntity>> submitAttempt({
  required String attemptId,
  required QuizAnswerPayload answers,
});

Future<Either<Failure, QuizResultEntity>> getResult(String attemptId);
```

---

## PART D — USE CASES

Create:

```text
list_quizzes_by_course_usecase.dart
get_quiz_usecase.dart
start_quiz_attempt_usecase.dart
save_quiz_progress_usecase.dart
submit_quiz_attempt_usecase.dart
get_quiz_result_usecase.dart
```

Rules:

- validate IDs are not empty
- save/submit answers may be empty only when backend allows it
- return `ValidationFailure` for invalid input
- do not grade locally

---

## PART E — DATA MODELS

Create models matching each entity:

```text
quiz_model.dart
quiz_attempt_model.dart
quiz_question_model.dart
quiz_answer_model.dart
quiz_result_model.dart
question_feedback_model.dart
subject_score_model.dart
```

Parsing rules:

- UUIDs are strings in Dart.
- Parse nullable `time_limit_seconds`, `max_attempts`, `expires_at`, `remaining_seconds`.
- Parse result feedback only from result responses.
- Fallback unknown question type to `multipleChoice`.
- Fallback unknown attempt status to `inProgress`.

Backend field mapping examples:

| Dart | JSON |
|------|------|
| `courseId` | `course_id` |
| `timeLimitSeconds` | `time_limit_seconds` |
| `passPercentage` | `pass_percentage` |
| `randomizeQuestions` | `randomize_questions` |
| `maxAttempts` | `max_attempts` |
| `questionCount` | `question_count` |
| `attemptsUsed` | `attempts_used` |
| `attemptId` | `attempt_id` |
| `attemptNumber` | `attempt_number` |
| `startedAt` | `started_at` |
| `expiresAt` | `expires_at` |
| `remainingSeconds` | `remaining_seconds` |
| `totalPoints` | `total_points` |
| `submittedAt` | `submitted_at` |
| `subjectBreakdown` | `subject_breakdown` |

---

## PART F — REMOTE DATASOURCE

### FILE F1 — `lib/features/quiz/data/datasources/quiz_remote_datasource.dart`

Use `DioClient` only.

Methods:

```dart
Future<List<QuizModel>> listQuizzesByCourse(String courseId);
Future<QuizModel> getQuiz(String quizId);
Future<QuizAttemptModel> startAttempt(String quizId);
Future<void> saveProgress({
  required String attemptId,
  required QuizAnswerPayload answers,
});
Future<QuizResultModel> submitAttempt({
  required String attemptId,
  required QuizAnswerPayload answers,
});
Future<QuizResultModel> getResult(String attemptId);
```

Endpoint examples:

```dart
await _client.get<Map<String, dynamic>>('/courses/$courseId/quizzes');
await _client.get<Map<String, dynamic>>('/quizzes/$quizId');
await _client.post<Map<String, dynamic>>('/quizzes/$quizId/attempts');
await _client.patch('/quiz-attempts/$attemptId', data: {'answers': answers});
await _client.post<Map<String, dynamic>>(
  '/quiz-attempts/$attemptId/submit',
  data: {'answers': answers},
);
await _client.get<Map<String, dynamic>>('/quiz-attempts/$attemptId/result');
```

Throw `FormatException` for missing bodies where expected.

---

## PART G — REPOSITORY IMPLEMENTATION

### FILE G1 — `lib/features/quiz/data/repositories/quiz_repository_impl.dart`

Map errors:

- `DioException` → `ServerFailure.fromDioError`
- any other exception → `UnknownFailure`

Do not do scoring logic in repository.

---

## PART H — PROVIDERS

### FILE H1 — `lib/features/quiz/presentation/providers/quiz_providers.dart`

Create:

- `quizRemoteDataSourceProvider`
- `quizRepositoryProvider`
- all use case providers

Mirror auth/courses/video provider style.

---

## VALIDATION

After this skill:

```bash
dart format lib/features/quiz
flutter analyze
```

Do not continue to quiz UI until this passes.

---

## SKILL_1 CHECKLIST

- [ ] Quiz entities compile
- [ ] Live question/answer entities contain no answer keys
- [ ] Result feedback entity is result-only
- [ ] Repository contract covers list/get/start/save/submit/result
- [ ] Use cases validate inputs
- [ ] Models parse backend JSON exactly
- [ ] Remote datasource calls quiz endpoints
- [ ] Repository maps failures consistently
- [ ] Dependency providers are wired
- [ ] `flutter analyze` has 0 errors

# phase5/SKILL_1.md — Mock Exam Domain, Data Layer & Backend Contract
# Memere Mobile (memere_mobile) — Phase 5, Part 1
# READ SKILL.md → all phase1 skill files → all phase2 skill files → all phase3 skill files → all phase4 skill files → then this file.

---

## OBJECTIVE

Build the Mock Exam foundation:
exam catalog entities → live exam attempt entities → answer payloads →
result/analytics entities → repository contract → use cases → backend DTO models →
remote datasource → repository implementation → Riverpod dependency providers.

By the end of this skill, the app can list published mock exams, start an exam
attempt, autosave answers, submit, fetch results, and fetch attempt analytics
from the local backend. Screens are built in `SKILL_2.md` and `SKILL_3.md`.

---

## PHASE 5 BOUNDARY

Phase 5 builds **mock exam taking only**:

- mock exam catalog
- subject/grade filters
- start exam attempt
- server-synced timer display
- live exam question flow
- autosave answers
- submit exam
- result screen
- attempt analytics screen

Do **not** build these yet:

- payment/enrollment
- leaderboard UI
- teacher exam builder
- advanced dashboard/trends
- certificates
- notifications

Backend leaderboard/trend endpoints can remain for later phases.

---

## SECURITY AND TIMER RULES

These are non-negotiable:

- Correct answers are never available during live exam attempts.
- Do not add `isCorrect` to live answer models.
- Do not grade on the client.
- Do not calculate pass/fail locally.
- Server timer is authoritative.
- UI timer is display-only.
- Backend decides expiration and final score.
- Only show correct answers/explanations after result fetch/submission.

---

## BACKEND CONTRACT

Student mock exam routes:

```text
GET /api/v1/mock-exams?limit=20&after=<cursor>&subject=<subject>&grade=12
POST /api/v1/mock-exams/:id/start
PATCH /api/v1/exam-attempts/:id
POST /api/v1/exam-attempts/:id/submit
GET /api/v1/exam-attempts/:id/results
GET /api/v1/exam-attempts/:id/analytics
```

Important:

- `GET /mock-exams` returns published exams only.
- `POST /mock-exams/:id/start` returns questions and server timing.
- `PATCH /exam-attempts/:id` autosaves and returns `204`.
- `POST /exam-attempts/:id/submit` returns graded result.
- `GET /exam-attempts/:id/analytics` returns weak areas/percentile when available.

---

## RESPONSE SHAPES

### `GET /mock-exams`

```json
{
  "data": [
    {
      "id": "uuid",
      "course_id": "uuid",
      "title": "Grade 12 Mathematics Mock Exam",
      "subject": "Mathematics",
      "grade": 12,
      "duration_minutes": 120,
      "total_marks": 100,
      "pass_marks": 50,
      "instructions": "Read each question carefully.",
      "is_published": true,
      "created_at": "2026-01-01T00:00:00Z",
      "updated_at": "2026-01-01T00:00:00Z"
    }
  ],
  "next_cursor": "",
  "limit": 20
}
```

### `POST /mock-exams/:id/start`

```json
{
  "attempt_id": "uuid",
  "exam_id": "uuid",
  "status": "in_progress",
  "started_at": "2026-01-01T00:00:00Z",
  "expires_at": "2026-01-01T02:00:00Z",
  "remaining_seconds": 7200,
  "total_marks": 100,
  "questions": [
    {
      "question_id": "uuid",
      "text": "Solve the problem.",
      "type": "multiple_choice",
      "marks": 2,
      "subject": "Mathematics",
      "topic": "Algebra",
      "answers": [
        {
          "id": "uuid",
          "text": "Answer option",
          "order_index": 1
        }
      ]
    }
  ]
}
```

### `PATCH /exam-attempts/:id`

Request:

```json
{
  "answers": {
    "question_uuid": "answer_uuid"
  }
}
```

Response: `204 No Content`

### `POST /exam-attempts/:id/submit`

Request:

```json
{
  "answers": {
    "question_uuid": "answer_uuid"
  }
}
```

Response: same shape as result below.

### `GET /exam-attempts/:id/results`

```json
{
  "attempt_id": "uuid",
  "exam_id": "uuid",
  "status": "submitted",
  "score": 76,
  "total_marks": 100,
  "percentage": 76,
  "pass_marks": 50,
  "passed": true,
  "submitted_at": "2026-01-01T01:50:00Z",
  "feedback": [
    {
      "question_id": "uuid",
      "correct": true,
      "marks_awarded": 2,
      "marks_possible": 2,
      "selected_answers": ["uuid"],
      "correct_answer_ids": ["uuid"],
      "explanation": "Explanation appears only after submission."
    }
  ],
  "subject_breakdown": {
    "Mathematics": {
      "earned": 76,
      "possible": 100
    }
  }
}
```

### `GET /exam-attempts/:id/analytics`

```json
{
  "attempt_id": "uuid",
  "exam_id": "uuid",
  "score": 76,
  "percentage": 76,
  "subject_breakdown": [
    {
      "key": "Mathematics",
      "earned": 76,
      "possible": 100
    }
  ],
  "weak_areas": [
    {
      "key": "Algebra",
      "earned": 4,
      "possible": 10
    }
  ],
  "percentile": 82.5
}
```

---

## REQUIRED FOLDER STRUCTURE

Create or verify:

```bash
mkdir -p lib/features/exam/data/datasources
mkdir -p lib/features/exam/data/models
mkdir -p lib/features/exam/data/repositories
mkdir -p lib/features/exam/domain/entities
mkdir -p lib/features/exam/domain/repositories
mkdir -p lib/features/exam/domain/usecases
mkdir -p lib/features/exam/presentation/providers
mkdir -p lib/features/exam/presentation/screens
mkdir -p lib/features/exam/presentation/widgets
```

---

## PART A — DOMAIN ENTITIES

### FILE A1 — `lib/features/exam/domain/entities/mock_exam_entity.dart`

Fields:

- `id`
- `courseId`
- `title`
- `subject`
- `grade`
- `durationMinutes`
- `totalMarks`
- `passMarks`
- `instructions`
- `isPublished`
- `createdAt`
- `updatedAt`

Helpers:

- `String get durationLabel`
- `String get marksLabel`

---

### FILE A2 — `lib/features/exam/domain/entities/paginated_mock_exams_entity.dart`

Fields:

- `exams`
- `nextCursor`
- `limit`

Helper:

- `bool get hasMore`

---

### FILE A3 — `lib/features/exam/domain/entities/exam_attempt_entity.dart`

Fields:

- `attemptId`
- `examId`
- `status`
- `startedAt`
- `expiresAt`
- `remainingSeconds`
- `totalMarks`
- `questions`

Use enum:

```dart
enum ExamAttemptStatus { inProgress, submitted, expired, graded }
```

---

### FILE A4 — `lib/features/exam/domain/entities/exam_question_entity.dart`

Fields:

- `questionId`
- `text`
- `type`
- `marks`
- `subject`
- `topic`
- `answers`

Use the same question type concepts from Phase 4. No answer key fields.

---

### FILE A5 — `lib/features/exam/domain/entities/exam_answer_entity.dart`

Fields:

- `id`
- `text`
- `orderIndex`

Do not add `isCorrect`.

---

### FILE A6 — `lib/features/exam/domain/entities/exam_result_entity.dart`

Fields:

- `attemptId`
- `examId`
- `status`
- `score`
- `totalMarks`
- `percentage`
- `passMarks`
- `passed`
- `submittedAt`
- `feedback`
- `subjectBreakdown`

---

### FILE A7 — `lib/features/exam/domain/entities/exam_question_feedback_entity.dart`

Fields:

- `questionId`
- `correct`
- `marksAwarded`
- `marksPossible`
- `selectedAnswers`
- `correctAnswerIds`
- `explanation`

Result-only.

---

### FILE A8 — `lib/features/exam/domain/entities/exam_attempt_analytics_entity.dart`

Fields:

- `attemptId`
- `examId`
- `score`
- `percentage`
- `subjectBreakdown`
- `weakAreas`
- `percentile`

---

### FILE A9 — `lib/features/exam/domain/entities/exam_subject_score_entity.dart`

Fields:

- `key`
- `earned`
- `possible`

Helper:

- `double get percentage`

---

## PART B — ANSWER PAYLOAD

### FILE B1 — `lib/features/exam/domain/entities/exam_answer_payload.dart`

Use:

```dart
typedef ExamAnswerPayload = Map<String, Object>;
```

Rules:

- single choice: `questionId -> answerId`
- multi select: `questionId -> List<String>`
- short answer: `questionId -> String`

Backend accepts `map[string]any`.

---

## PART C — REPOSITORY CONTRACT

### FILE C1 — `lib/features/exam/domain/repositories/exam_repository.dart`

Use `Either<Failure, T>`.

Methods:

```dart
Future<Either<Failure, PaginatedMockExamsEntity>> listMockExams({
  int limit = 20,
  String? after,
  String? subject,
  int? grade,
});

Future<Either<Failure, ExamAttemptEntity>> startExam(String examId);

Future<Either<Failure, void>> saveProgress({
  required String attemptId,
  required ExamAnswerPayload answers,
});

Future<Either<Failure, ExamResultEntity>> submitExam({
  required String attemptId,
  required ExamAnswerPayload answers,
});

Future<Either<Failure, ExamResultEntity>> getResult(String attemptId);

Future<Either<Failure, ExamAttemptAnalyticsEntity>> getAnalytics(String attemptId);
```

---

## PART D — USE CASES

Create:

```text
list_mock_exams_usecase.dart
start_exam_usecase.dart
save_exam_progress_usecase.dart
submit_exam_usecase.dart
get_exam_result_usecase.dart
get_exam_analytics_usecase.dart
```

Rules:

- validate IDs are not empty
- validate pagination limit is reasonable
- save/submit answers may be empty only when backend allows it
- return `ValidationFailure` for invalid input
- do not grade locally

---

## PART E — DATA MODELS

Create models matching each entity:

```text
mock_exam_model.dart
paginated_mock_exams_model.dart
exam_attempt_model.dart
exam_question_model.dart
exam_answer_model.dart
exam_result_model.dart
exam_question_feedback_model.dart
exam_attempt_analytics_model.dart
exam_subject_score_model.dart
```

Parsing rules:

- UUIDs are strings in Dart.
- `course_id` is nullable.
- Parse nullable `instructions`, `expires_at`, `remaining_seconds`, `submitted_at`, and `percentile`.
- Parse live questions without correct answer fields.
- Parse result feedback only from result responses.
- Fallback unknown attempt status to `inProgress`.

---

## PART F — REMOTE DATASOURCE

### FILE F1 — `lib/features/exam/data/datasources/exam_remote_datasource.dart`

Use `DioClient` only.

Methods:

```dart
Future<PaginatedMockExamsModel> listMockExams({
  int limit = 20,
  String? after,
  String? subject,
  int? grade,
});

Future<ExamAttemptModel> startExam(String examId);

Future<void> saveProgress({
  required String attemptId,
  required ExamAnswerPayload answers,
});

Future<ExamResultModel> submitExam({
  required String attemptId,
  required ExamAnswerPayload answers,
});

Future<ExamResultModel> getResult(String attemptId);

Future<ExamAttemptAnalyticsModel> getAnalytics(String attemptId);
```

Endpoint examples:

```dart
await _client.get<Map<String, dynamic>>(
  '/mock-exams',
  queryParameters: {
    'limit': limit,
    if (after != null && after.isNotEmpty) 'after': after,
    if (subject != null && subject.isNotEmpty) 'subject': subject,
    if (grade != null) 'grade': grade,
  },
);
await _client.post<Map<String, dynamic>>('/mock-exams/$examId/start');
await _client.patch('/exam-attempts/$attemptId', data: {'answers': answers});
await _client.post<Map<String, dynamic>>(
  '/exam-attempts/$attemptId/submit',
  data: {'answers': answers},
);
await _client.get<Map<String, dynamic>>('/exam-attempts/$attemptId/results');
await _client.get<Map<String, dynamic>>('/exam-attempts/$attemptId/analytics');
```

Throw `FormatException` for missing bodies where expected.

---

## PART G — REPOSITORY IMPLEMENTATION

### FILE G1 — `lib/features/exam/data/repositories/exam_repository_impl.dart`

Map errors:

- `DioException` → `ServerFailure.fromDioError`
- any other exception → `UnknownFailure`

Do not do scoring/timer enforcement in repository.

---

## PART H — PROVIDERS

### FILE H1 — `lib/features/exam/presentation/providers/exam_providers.dart`

Create:

- `examRemoteDataSourceProvider`
- `examRepositoryProvider`
- all use case providers

Mirror auth/courses/video/quiz provider style.

---

## VALIDATION

After this skill:

```bash
dart format lib/features/exam
flutter analyze
```

Do not continue to mock exam UI until this passes.

---

## SKILL_1 CHECKLIST

- [ ] Mock exam entities compile
- [ ] Live exam question/answer entities contain no answer keys
- [ ] Result feedback is result-only
- [ ] Analytics entities compile
- [ ] Repository contract covers catalog/start/save/submit/result/analytics
- [ ] Use cases validate inputs
- [ ] Models parse backend JSON exactly
- [ ] Remote datasource calls exam endpoints
- [ ] Repository maps failures consistently
- [ ] Dependency providers are wired
- [ ] `flutter analyze` has 0 errors

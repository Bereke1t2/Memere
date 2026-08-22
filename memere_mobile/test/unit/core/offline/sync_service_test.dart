import 'dart:io';

import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:fpdart/fpdart.dart';
import 'package:hive_flutter/hive_flutter.dart';
import 'package:memere_mobile/core/constants/app_constants.dart';
import 'package:memere_mobile/core/errors/failures.dart';
import 'package:memere_mobile/core/network/connectivity_service.dart';
import 'package:memere_mobile/core/offline/offline_attempt_store.dart';
import 'package:memere_mobile/core/storage/hive/models/pending_submission.dart';
import 'package:memere_mobile/features/auth/presentation/providers/auth_state_provider.dart';
import 'package:memere_mobile/core/offline/sync_service.dart';
import 'package:memere_mobile/features/exam/data/models/exam_result_model.dart';
import 'package:memere_mobile/features/exam/domain/entities/exam_answer_payload.dart';
import 'package:memere_mobile/features/exam/domain/entities/exam_attempt_entity.dart';
import 'package:memere_mobile/features/exam/domain/entities/exam_result_entity.dart';
import 'package:memere_mobile/features/exam/domain/repositories/exam_repository.dart';
import 'package:memere_mobile/features/exam/presentation/providers/exam_providers.dart';
import 'package:memere_mobile/features/quiz/data/models/quiz_result_model.dart';
import 'package:memere_mobile/features/quiz/domain/entities/quiz_answer_payload.dart';
import 'package:memere_mobile/features/quiz/domain/entities/quiz_attempt_entity.dart';
import 'package:memere_mobile/features/quiz/domain/entities/quiz_result_entity.dart';
import 'package:memere_mobile/features/quiz/domain/repositories/quiz_repository.dart';
import 'package:memere_mobile/features/quiz/presentation/providers/quiz_providers.dart';

/// Drives [SyncService.drain] against fake repositories and a real (plaintext)
/// Hive-backed [OfflineAttemptStore], exercising the replay state machine end to
/// end: start→submit→reconcile, the serverAttemptId skip-start path, crash-safe
/// persistence of the server attempt id before submit, bounded retry/give-up,
/// and the authed-only + online-only guards.
///
/// The use-case providers wrap the (overridden) repository providers, so faking
/// the repos is the narrowest seam that still runs the real use cases.
void main() {
  late Directory tempDir;
  const store = OfflineAttemptStore();

  setUpAll(() async {
    tempDir = await Directory.systemTemp.createTemp('sync_service');
    Hive.init(tempDir.path);
    await Hive.openBox<String>(AppConstants.offlineAttemptResultsBoxKey);
    await Hive.openBox<String>(AppConstants.syncQueueKey);
  });

  tearDown(() async {
    await Hive.box<String>(AppConstants.offlineAttemptResultsBoxKey).clear();
    await Hive.box<String>(AppConstants.syncQueueKey).clear();
  });

  tearDownAll(() async {
    await Hive.close();
    await tempDir.delete(recursive: true);
  });

  group('quiz replay', () {
    test('start → submit → reconcile → dequeue (no serverAttemptId yet)',
        () async {
      await store.saveResult('local-1', {'attempt_id': 'local-1', 'score': 1});
      await store.enqueue(_pending(localId: 'local-1'));

      final repo = _FakeQuizRepository(
        startResult: Right(_quizAttempt('srv-1')),
        submitResult: Right(_quizResult(attemptId: 'srv-1', score: 9)),
      );
      final container = await _containerFor(quizRepo: repo);
      addTearDown(container.dispose);

      await container.read(syncServiceProvider).drain();

      // Attempt was created once, submitted to the returned server id, with the
      // persisted raw answers forwarded verbatim.
      expect(repo.startCalls, 1);
      expect(repo.submitCalls, 1);
      expect(repo.lastSubmitAttemptId, 'srv-1');
      expect(repo.lastSubmitAnswers?['q1'], ['a1']);

      // Provisional result replaced by the server's; queue drained.
      final reconciled = store.getResult('local-1');
      expect(reconciled?['attempt_id'], 'srv-1');
      expect(reconciled?['score'], 9);
      expect(store.listQueue(), isEmpty);
    });

    test('reuses an existing serverAttemptId and skips start', () async {
      await store.enqueue(_pending(localId: 'local-2', serverAttemptId: 'srv-9'));

      final repo = _FakeQuizRepository(
        submitResult: Right(_quizResult(attemptId: 'srv-9', score: 7)),
      );
      final container = await _containerFor(quizRepo: repo);
      addTearDown(container.dispose);

      await container.read(syncServiceProvider).drain();

      expect(repo.startCalls, 0);
      expect(repo.submitCalls, 1);
      expect(repo.lastSubmitAttemptId, 'srv-9');
      expect(store.getResult('local-2')?['score'], 7);
      expect(store.listQueue(), isEmpty);
    });

    test('start succeeds but submit fails: persists serverAttemptId, bumps, keeps queued',
        () async {
      await store.enqueue(_pending(localId: 'local-3'));

      final repo = _FakeQuizRepository(
        startResult: Right(_quizAttempt('srv-1')),
        submitResult: const Left<Failure, QuizResultEntity>(
          ServerFailure('offline', code: 'NO_INTERNET'),
        ),
      );
      final container = await _containerFor(quizRepo: repo);
      addTearDown(container.dispose);

      await container.read(syncServiceProvider).drain();

      expect(repo.startCalls, 1);
      expect(repo.submitCalls, 1);
      // The server attempt id is captured so the next drain re-submits to the
      // SAME attempt instead of starting a duplicate; the failure just bumps.
      final queued = store.listQueue();
      expect(queued, hasLength(1));
      expect(queued.single.serverAttemptId, 'srv-1');
      expect(queued.single.attempts, 1);
      expect(store.getResult('local-3'), isNull); // not reconciled
    });

    test('start failure bumps without submitting', () async {
      await store.enqueue(_pending(localId: 'local-4'));

      final repo = _FakeQuizRepository(
        startResult: const Left<Failure, QuizAttemptEntity>(
          ServerFailure('offline', code: 'NO_INTERNET'),
        ),
      );
      final container = await _containerFor(quizRepo: repo);
      addTearDown(container.dispose);

      await container.read(syncServiceProvider).drain();

      expect(repo.startCalls, 1);
      expect(repo.submitCalls, 0);
      final queued = store.listQueue();
      expect(queued, hasLength(1));
      expect(queued.single.serverAttemptId, isNull);
      expect(queued.single.attempts, 1);
    });

    test('drops the item after the max-attempts threshold, keeping its result',
        () async {
      await store.saveResult('local-5', {'attempt_id': 'local-5', 'score': 3});
      await store.enqueue(_pending(localId: 'local-5', attempts: 5));

      final repo = _FakeQuizRepository(
        startResult: Right(_quizAttempt('srv-1')),
        submitResult: Right(_quizResult(attemptId: 'srv-1', score: 9)),
      );
      final container = await _containerFor(quizRepo: repo);
      addTearDown(container.dispose);

      await container.read(syncServiceProvider).drain();

      // Given up on: no network calls, removed from the queue, but the
      // provisional on-device result is retained.
      expect(repo.startCalls, 0);
      expect(repo.submitCalls, 0);
      expect(store.listQueue(), isEmpty);
      expect(store.getResult('local-5')?['score'], 3);
    });
  });

  group('exam replay', () {
    test('start → submit → reconcile → dequeue', () async {
      await store.saveResult('local-e', {'attempt_id': 'local-e', 'score': 1});
      await store.enqueue(
        _pending(
          localId: 'local-e',
          kind: SubmissionKind.exam,
          contentId: 'exam-1',
        ),
      );

      final repo = _FakeExamRepository(
        startResult: Right(_examAttempt('srv-e')),
        submitResult: Right(_examResult(attemptId: 'srv-e', score: 20)),
      );
      final container = await _containerFor(examRepo: repo);
      addTearDown(container.dispose);

      await container.read(syncServiceProvider).drain();

      expect(repo.startCalls, 1);
      expect(repo.submitCalls, 1);
      expect(repo.lastSubmitAttemptId, 'srv-e');
      final reconciled = store.getResult('local-e');
      expect(reconciled?['attempt_id'], 'srv-e');
      expect(reconciled?['score'], 20);
      expect(store.listQueue(), isEmpty);
    });
  });

  group('guards', () {
    test('guests never drain (no token to replay with)', () async {
      await store.enqueue(_pending(localId: 'local-g'));

      final repo = _FakeQuizRepository(
        startResult: Right(_quizAttempt('srv-1')),
        submitResult: Right(_quizResult(attemptId: 'srv-1', score: 9)),
      );
      final container = await _containerFor(quizRepo: repo, authed: false);
      addTearDown(container.dispose);

      await container.read(syncServiceProvider).drain();

      expect(repo.startCalls, 0);
      expect(store.listQueue(), hasLength(1));
    });

    test('offline is a clean no-op (does not burn a retry)', () async {
      await store.enqueue(_pending(localId: 'local-o'));

      final repo = _FakeQuizRepository(
        startResult: Right(_quizAttempt('srv-1')),
        submitResult: Right(_quizResult(attemptId: 'srv-1', score: 9)),
      );
      final container = await _containerFor(quizRepo: repo, online: false);
      addTearDown(container.dispose);

      await container.read(syncServiceProvider).drain();

      expect(repo.startCalls, 0);
      final queued = store.listQueue();
      expect(queued, hasLength(1));
      expect(queued.single.attempts, 0); // untouched, not bumped
    });
  });
}

// ── Container / provider overrides ────────────────────────────────────────────

Future<ProviderContainer> _containerFor({
  _FakeQuizRepository? quizRepo,
  _FakeExamRepository? examRepo,
  bool authed = true,
  bool online = true,
}) async {
  final container = ProviderContainer(
    overrides: [
      quizRepositoryProvider.overrideWithValue(quizRepo ?? _FakeQuizRepository()),
      examRepositoryProvider.overrideWithValue(examRepo ?? _FakeExamRepository()),
      connectivityServiceProvider
          .overrideWithValue(_FakeConnectivityService(online)),
      authStateProvider.overrideWith(() => _FakeAuthNotifier(authed)),
    ],
  );
  // Resolve the async build so `drain`'s synchronous `.read` sees AsyncData.
  await container.read(authStateProvider.future);
  return container;
}

// ── Fakes ─────────────────────────────────────────────────────────────────────

class _FakeAuthNotifier extends AuthStateNotifier {
  _FakeAuthNotifier(this.authed);

  final bool authed;

  @override
  Future<AuthState> build() async => AuthState(isAuthenticated: authed);
}

class _FakeConnectivityService extends ConnectivityService {
  _FakeConnectivityService(this.online);

  final bool online;

  @override
  Future<bool> isOnline() async => online;
}

class _FakeQuizRepository implements QuizRepository {
  _FakeQuizRepository({this.startResult, this.submitResult});

  Either<Failure, QuizAttemptEntity>? startResult;
  Either<Failure, QuizResultEntity>? submitResult;

  int startCalls = 0;
  int submitCalls = 0;
  String? lastSubmitAttemptId;
  QuizAnswerPayload? lastSubmitAnswers;

  @override
  Future<Either<Failure, QuizAttemptEntity>> startAttempt(String quizId) async {
    startCalls++;
    return startResult ??
        const Left<Failure, QuizAttemptEntity>(
          ServerFailure('no start configured'),
        );
  }

  @override
  Future<Either<Failure, QuizResultEntity>> submitAttempt({
    required String attemptId,
    required QuizAnswerPayload answers,
  }) async {
    submitCalls++;
    lastSubmitAttemptId = attemptId;
    lastSubmitAnswers = answers;
    return submitResult ??
        const Left<Failure, QuizResultEntity>(
          ServerFailure('no submit configured'),
        );
  }

  @override
  dynamic noSuchMethod(Invocation invocation) => throw UnimplementedError();
}

class _FakeExamRepository implements ExamRepository {
  _FakeExamRepository({this.startResult, this.submitResult});

  Either<Failure, ExamAttemptEntity>? startResult;
  Either<Failure, ExamResultEntity>? submitResult;

  int startCalls = 0;
  int submitCalls = 0;
  String? lastSubmitAttemptId;
  ExamAnswerPayload? lastSubmitAnswers;

  @override
  Future<Either<Failure, ExamAttemptEntity>> startExam(String examId) async {
    startCalls++;
    return startResult ??
        const Left<Failure, ExamAttemptEntity>(
          ServerFailure('no start configured'),
        );
  }

  @override
  Future<Either<Failure, ExamResultEntity>> submitExam({
    required String attemptId,
    required ExamAnswerPayload answers,
  }) async {
    submitCalls++;
    lastSubmitAttemptId = attemptId;
    lastSubmitAnswers = answers;
    return submitResult ??
        const Left<Failure, ExamResultEntity>(
          ServerFailure('no submit configured'),
        );
  }

  @override
  dynamic noSuchMethod(Invocation invocation) => throw UnimplementedError();
}

// ── Builders ──────────────────────────────────────────────────────────────────

PendingSubmission _pending({
  required String localId,
  SubmissionKind kind = SubmissionKind.quiz,
  String contentId = 'quiz-1',
  String? serverAttemptId,
  int attempts = 0,
}) =>
    PendingSubmission(
      localId: localId,
      kind: kind,
      contentId: contentId,
      rawAnswers: {
        'q1': ['a1'],
      },
      provisionalScore: 1,
      provisionalPct: 10.0,
      takenAt: DateTime.parse('2026-08-21T10:00:00.000Z'),
      serverAttemptId: serverAttemptId,
      attempts: attempts,
    );

QuizAttemptEntity _quizAttempt(String attemptId) => QuizAttemptEntity(
      attemptId: attemptId,
      quizId: 'quiz-1',
      attemptNumber: 1,
      status: QuizAttemptStatus.inProgress,
      startedAt: DateTime.parse('2026-08-21T11:00:00.000Z'),
      questions: const [],
    );

QuizResultModel _quizResult({
  required String attemptId,
  required double score,
}) =>
    QuizResultModel(
      attemptId: attemptId,
      quizId: 'quiz-1',
      attemptNumber: 1,
      status: QuizAttemptStatus.graded,
      score: score,
      totalPoints: 10,
      percentage: score * 10,
      passed: score >= 5,
      submittedAt: DateTime.parse('2026-08-21T12:00:00.000Z'),
      feedback: const [],
      subjectBreakdown: const {},
    );

ExamAttemptEntity _examAttempt(String attemptId) => ExamAttemptEntity(
      attemptId: attemptId,
      examId: 'exam-1',
      status: ExamAttemptStatus.inProgress,
      totalMarks: 30,
      questions: const [],
    );

ExamResultModel _examResult({
  required String attemptId,
  required double score,
}) =>
    ExamResultModel(
      attemptId: attemptId,
      examId: 'exam-1',
      status: ExamAttemptStatus.graded,
      score: score,
      totalMarks: 30,
      percentage: score,
      passMarks: 15,
      passed: score >= 15,
      submittedAt: DateTime.parse('2026-08-21T12:00:00.000Z'),
      feedback: const [],
      subjectBreakdown: const {},
    );

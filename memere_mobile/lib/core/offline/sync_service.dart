import 'dart:developer' as developer;

import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../features/auth/presentation/providers/auth_state_provider.dart';
import '../../features/exam/data/models/exam_result_model.dart';
import '../../features/exam/domain/usecases/submit_exam_usecase.dart';
import '../../features/exam/presentation/providers/exam_providers.dart';
import '../../features/quiz/data/models/quiz_result_model.dart';
import '../../features/quiz/domain/entities/quiz_answer_payload.dart';
import '../../features/quiz/domain/usecases/submit_quiz_attempt_usecase.dart';
import '../../features/quiz/presentation/providers/quiz_providers.dart';
import '../network/connectivity_service.dart';
import '../storage/hive/models/pending_submission.dart';
import 'offline_attempt_store.dart';
import 'offline_providers.dart';

/// Foreground drainer for the offline-submission queue (authenticated users
/// only — guests never enqueue, having no token to replay with).
///
/// On a reconnect (connectivity `off→on`) or app-resume, [drain] replays each
/// queued [PendingSubmission] through the normal `start`→`submit` endpoints so
/// the server grades it authoritatively, then overwrites the provisional
/// on-device result with the server's. The `local-…` id keeps the client-side
/// record stable across the replay; the server attempt id is captured once so a
/// crash between `start` and `submit` never starts a duplicate attempt.
///
/// This is v1's deliberately simple, timer-free strategy: it runs only when the
/// app is foregrounded and something triggers it (no `workmanager`). Retries are
/// bounded — after [_maxAttempts] failures an item is dropped (and logged, never
/// silently), leaving its provisional result in place.
class SyncService {
  SyncService(this._ref);

  final Ref _ref;

  /// Guards against overlapping drains (a reconnect and a resume can fire almost
  /// together); a second call while one is in flight is a no-op.
  bool _draining = false;

  /// Give-up threshold. Each failed replay bumps `PendingSubmission.attempts`;
  /// past this the item is dropped so a permanently-rejected submission can't
  /// wedge the queue forever.
  static const int _maxAttempts = 5;

  Future<void> drain() async {
    if (_draining) return;

    // Authed-only: without a token the replay would 401. Guests keep their
    // on-device score as final and never reach this path.
    final authed =
        _ref.read(authStateProvider).valueOrNull?.isAuthenticated ?? false;
    if (!authed) return;

    // Reachability gate: skip cleanly when offline rather than burning a retry
    // attempt on a submit that would just fail with NO_INTERNET.
    final online = await _ref.read(connectivityServiceProvider).isOnline();
    if (!online) return;

    _draining = true;
    try {
      final store = _ref.read(offlineAttemptStoreProvider);
      for (final item in store.listQueue()) {
        await _process(item, store);
      }
    } finally {
      _draining = false;
    }
  }

  Future<void> _process(PendingSubmission item, OfflineAttemptStore store) async {
    if (item.attempts >= _maxAttempts) {
      developer.log(
        'dropping ${item.kind.name} ${item.localId} after ${item.attempts} '
        'failed replays; provisional result kept',
        name: 'SyncService',
      );
      await store.dequeue(item.localId);
      return;
    }

    switch (item.kind) {
      case SubmissionKind.quiz:
        await _syncQuiz(item, store);
      case SubmissionKind.exam:
        await _syncExam(item, store);
    }
  }

  Future<void> _syncQuiz(PendingSubmission item, OfflineAttemptStore store) async {
    var current = item;

    // 1. Ensure a server-side attempt exists, persisting its id before we submit
    //    so a crash in between resumes from here instead of starting a new one.
    if (current.serverAttemptId == null) {
      final started =
          await _ref.read(startQuizAttemptUseCaseProvider)(current.contentId);
      final withId = started.fold<PendingSubmission?>(
        (_) => null,
        (attempt) => current.copyWith(serverAttemptId: attempt.attemptId),
      );
      if (withId == null) {
        await _bump(current, store);
        return;
      }
      current = withId;
      await store.enqueue(current);
    }

    // 2. Submit for authoritative grading; reconcile the provisional result.
    final result = await _ref.read(submitQuizAttemptUseCaseProvider)(
      SubmitQuizAttemptParams(
        attemptId: current.serverAttemptId!,
        answers: _payload(current.rawAnswers),
      ),
    );
    await result.fold(
      (_) => _bump(current, store),
      (serverResult) async {
        if (serverResult is QuizResultModel) {
          await store.saveResult(current.localId, serverResult.toJson());
        }
        await store.dequeue(current.localId);
      },
    );
  }

  Future<void> _syncExam(PendingSubmission item, OfflineAttemptStore store) async {
    var current = item;

    if (current.serverAttemptId == null) {
      final started =
          await _ref.read(startExamUseCaseProvider)(current.contentId);
      final withId = started.fold<PendingSubmission?>(
        (_) => null,
        (attempt) => current.copyWith(serverAttemptId: attempt.attemptId),
      );
      if (withId == null) {
        await _bump(current, store);
        return;
      }
      current = withId;
      await store.enqueue(current);
    }

    final result = await _ref.read(submitExamUseCaseProvider)(
      SubmitExamParams(
        attemptId: current.serverAttemptId!,
        answers: _payload(current.rawAnswers),
      ),
    );
    await result.fold(
      (_) => _bump(current, store),
      (serverResult) async {
        if (serverResult is ExamResultModel) {
          await store.saveResult(current.localId, serverResult.toJson());
        }
        await store.dequeue(current.localId);
      },
    );
  }

  /// Records one failed replay; the item stays queued and is retried on the next
  /// drain until it succeeds or hits [_maxAttempts].
  Future<void> _bump(PendingSubmission current, OfflineAttemptStore store) =>
      store.enqueue(current.copyWith(attempts: current.attempts + 1));

  /// The persisted raw answers are `Map<String, dynamic>` (JSON-decoded); the
  /// submit endpoints want a non-null `Map<String, Object>`. Values are option
  /// id lists or short-answer strings — never null once stored — but we drop any
  /// null defensively so the cast is total.
  QuizAnswerPayload _payload(Map<String, dynamic> raw) => {
        for (final entry in raw.entries)
          if (entry.value != null) entry.key: entry.value as Object,
      };
}

final syncServiceProvider = Provider<SyncService>((ref) => SyncService(ref));

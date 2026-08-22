import 'dart:async';

import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/errors/failures.dart';
import '../../../../core/grading/offline_grader.dart';
import '../../../../core/offline/offline_attempt_factory.dart';
import '../../../../core/offline/offline_providers.dart';
import '../../../../core/storage/hive/models/downloaded_item.dart';
import '../../../../core/storage/hive/models/pending_submission.dart';
import '../../../auth/presentation/providers/auth_state_provider.dart';
import '../../domain/entities/quiz_answer_payload.dart';
import '../../domain/entities/quiz_attempt_entity.dart';
import '../../domain/usecases/save_quiz_progress_usecase.dart';
import '../../domain/usecases/submit_quiz_attempt_usecase.dart';
import 'quiz_providers.dart';

class QuizAttemptParams {
  const QuizAttemptParams({
    required this.attemptId,
    required this.quizId,
  });

  final String attemptId;
  final String quizId;

  @override
  bool operator ==(Object other) {
    return other is QuizAttemptParams &&
        other.attemptId == attemptId &&
        other.quizId == quizId;
  }

  @override
  int get hashCode => Object.hash(attemptId, quizId);
}

class QuizAttemptState {
  const QuizAttemptState({
    required this.attempt,
    this.currentIndex = 0,
    this.answers = const {},
    this.isSaving = false,
    this.lastSavedAt,
    this.saveError,
    this.isSubmitting = false,
    this.hasUnsavedChanges = false,
  });

  final QuizAttemptEntity attempt;
  final int currentIndex;
  final QuizAnswerPayload answers;
  final bool isSaving;
  final DateTime? lastSavedAt;
  final String? saveError;
  final bool isSubmitting;
  final bool hasUnsavedChanges;

  int get questionCount => attempt.questions.length;
  int get answeredCount => answers.values.where(_answerValueHasContent).length;
  int get unansweredCount => questionCount - answeredCount;

  QuizAttemptState copyWith({
    QuizAttemptEntity? attempt,
    int? currentIndex,
    QuizAnswerPayload? answers,
    bool? isSaving,
    DateTime? lastSavedAt,
    String? saveError,
    bool? isSubmitting,
    bool? hasUnsavedChanges,
    bool clearSaveError = false,
  }) {
    return QuizAttemptState(
      attempt: attempt ?? this.attempt,
      currentIndex: currentIndex ?? this.currentIndex,
      answers: answers ?? this.answers,
      isSaving: isSaving ?? this.isSaving,
      lastSavedAt: lastSavedAt ?? this.lastSavedAt,
      saveError: clearSaveError ? null : saveError ?? this.saveError,
      isSubmitting: isSubmitting ?? this.isSubmitting,
      hasUnsavedChanges: hasUnsavedChanges ?? this.hasUnsavedChanges,
    );
  }
}

final quizAttemptProvider = AsyncNotifierProviderFamily<QuizAttemptNotifier,
    QuizAttemptState, QuizAttemptParams>(QuizAttemptNotifier.new);

class QuizAttemptNotifier
    extends FamilyAsyncNotifier<QuizAttemptState, QuizAttemptParams> {
  Timer? _autosaveTimer;

  @override
  Future<QuizAttemptState> build(QuizAttemptParams arg) async {
    ref.onDispose(() {
      _autosaveTimer?.cancel();
      unawaited(saveProgress(force: true));
    });

    // Offline / guest path: a locally-minted attempt id means grade this attempt
    // on-device from the downloaded copy. No network, and no server-side attempt
    // is ever created (this also sidesteps the online double-start).
    if (isLocalAttemptId(arg.attemptId)) {
      final quiz = ref.read(downloadStoreProvider).getOfflineQuiz(arg.quizId);
      if (quiz == null) {
        throw const CacheFailure(
          'This quiz is not downloaded. Download it first to take it offline.',
        );
      }
      return QuizAttemptState(
        attempt: buildLocalQuizAttempt(quiz, attemptId: arg.attemptId),
      );
    }

    final result = await ref.read(startQuizAttemptUseCaseProvider)(arg.quizId);
    final attempt = result.fold((failure) => throw failure, (value) => value);
    return QuizAttemptState(attempt: attempt);
  }

  void selectSingleAnswer(String questionId, String answerId) {
    // Store choice answers as a one-element list so the backend files them as
    // "selected" option IDs (its bare-array shorthand). A bare String is reserved
    // for short_answer text — sending it here made the grader treat the answer as
    // free text, leaving Selected empty and scoring every choice question 0.
    _updateAnswer(questionId, <String>[answerId]);
  }

  void toggleMultiAnswer(String questionId, String answerId) {
    final current = state.valueOrNull;
    if (current == null) return;
    final existing = current.answers[questionId];
    final values = existing is List
        ? existing.map((item) => item.toString()).toList()
        : <String>[];
    if (values.contains(answerId)) {
      values.remove(answerId);
    } else {
      values.add(answerId);
    }
    _updateAnswer(questionId, List<String>.unmodifiable(values));
  }

  void setShortAnswer(String questionId, String value) {
    final trimmed = value.trim();
    if (trimmed.isEmpty) {
      _removeAnswer(questionId);
      return;
    }
    _updateAnswer(questionId, trimmed);
  }

  void goToQuestion(int index) {
    final current = state.valueOrNull;
    if (current == null || current.questionCount == 0) return;
    final nextIndex = index.clamp(0, current.questionCount - 1).toInt();
    state = AsyncData(current.copyWith(currentIndex: nextIndex));
  }

  void nextQuestion() {
    final current = state.valueOrNull;
    if (current == null) return;
    goToQuestion(current.currentIndex + 1);
  }

  void previousQuestion() {
    final current = state.valueOrNull;
    if (current == null) return;
    goToQuestion(current.currentIndex - 1);
  }

  Future<void> saveProgress({bool force = false}) async {
    final current = state.valueOrNull;
    if (current == null || current.isSaving) return;
    if (!force && !current.hasUnsavedChanges) return;

    // Local attempts keep their answers in memory until submit — there is no
    // server progress endpoint for a `local-…` id, so just mark them saved.
    if (isLocalAttemptId(current.attempt.attemptId)) {
      state = AsyncData(
        current.copyWith(
          lastSavedAt: DateTime.now(),
          hasUnsavedChanges: false,
          clearSaveError: true,
        ),
      );
      return;
    }

    state = AsyncData(current.copyWith(isSaving: true, clearSaveError: true));
    final result = await ref.read(saveQuizProgressUseCaseProvider)(
      SaveQuizProgressParams(
        attemptId: current.attempt.attemptId,
        answers: current.answers,
      ),
    );
    final latest = state.valueOrNull ?? current;
    result.fold(
      (failure) {
        state = AsyncData(
          latest.copyWith(
            isSaving: false,
            saveError: failure.message,
          ),
        );
      },
      (_) {
        state = AsyncData(
          latest.copyWith(
            isSaving: false,
            lastSavedAt: DateTime.now(),
            hasUnsavedChanges: false,
            clearSaveError: true,
          ),
        );
      },
    );
  }

  Future<String?> submit() async {
    final current = state.valueOrNull;
    if (current == null || current.isSubmitting) return null;

    await saveProgress(force: true);

    final latestBeforeSubmit = state.valueOrNull ?? current;
    state = AsyncData(
      latestBeforeSubmit.copyWith(isSubmitting: true, clearSaveError: true),
    );

    if (isLocalAttemptId(latestBeforeSubmit.attempt.attemptId)) {
      return _submitLocally(latestBeforeSubmit);
    }

    final result = await ref.read(submitQuizAttemptUseCaseProvider)(
      SubmitQuizAttemptParams(
        attemptId: latestBeforeSubmit.attempt.attemptId,
        answers: latestBeforeSubmit.answers,
      ),
    );
    final latest = state.valueOrNull ?? latestBeforeSubmit;
    return result.fold(
      (failure) async {
        // Started on the server but the network dropped at submit time: grade
        // the downloaded copy on-device now, and queue the answers to re-submit
        // to the SAME server attempt on reconnect (no duplicate attempt).
        if (isOfflineFailure(failure) &&
            ref.read(downloadStoreProvider).isDownloaded(
                  DownloadType.quiz,
                  latestBeforeSubmit.attempt.quizId,
                )) {
          return _submitLocally(
            latest,
            serverAttemptId: latestBeforeSubmit.attempt.attemptId,
          );
        }
        state = AsyncData(
          latest.copyWith(
            isSubmitting: false,
            saveError: failure.message,
          ),
        );
        return null;
      },
      (quizResult) async {
        state = AsyncData(
          latest.copyWith(
            isSubmitting: false,
            hasUnsavedChanges: false,
            clearSaveError: true,
          ),
        );
        return quizResult.attemptId;
      },
    );
  }

  /// Grades a downloaded quiz on-device, stores the provisional result under the
  /// `local-…` id, and — for signed-in users only — enqueues a [PendingSubmission]
  /// for authoritative server replay on reconnect. Guests never enqueue (no token
  /// to submit with); their on-device score stands as final. Returns the local id
  /// so the result screen resolves it from [OfflineAttemptStore].
  ///
  /// [serverAttemptId] distinguishes the two callers: `null` for the pure
  /// offline/guest path (the attempt id is already a `local-…` id); non-null for
  /// the submit-time fallback (an online attempt whose submit hit the network) —
  /// then a fresh local id anchors the provisional record and the existing server
  /// attempt id rides along so sync skips `start` and re-submits straight to it.
  Future<String?> _submitLocally(
    QuizAttemptState current, {
    String? serverAttemptId,
  }) async {
    final localId = serverAttemptId == null
        ? current.attempt.attemptId
        : newLocalAttemptId();
    final quiz =
        ref.read(downloadStoreProvider).getOfflineQuiz(current.attempt.quizId);
    if (quiz == null) {
      final latest = state.valueOrNull ?? current;
      state = AsyncData(
        latest.copyWith(
          isSubmitting: false,
          saveError:
              'This quiz is no longer downloaded. Reconnect or re-download to submit.',
        ),
      );
      return null;
    }

    final rawAnswers = Map<String, dynamic>.from(current.answers);
    final result = OfflineGrader.gradeQuiz(
      quiz: quiz,
      submission: offlineSubmissionFromRaw(quiz.questions, rawAnswers),
      attemptId: localId,
    );

    final store = ref.read(offlineAttemptStoreProvider);
    await store.saveResult(localId, result.toJson());

    final isAuthed =
        ref.read(authStateProvider).valueOrNull?.isAuthenticated ?? false;
    if (isAuthed) {
      await store.enqueue(
        PendingSubmission(
          localId: localId,
          kind: SubmissionKind.quiz,
          contentId: quiz.id,
          rawAnswers: rawAnswers,
          provisionalScore: result.score.round(),
          provisionalPct: result.percentage,
          takenAt: result.submittedAt ?? DateTime.now(),
          serverAttemptId: serverAttemptId,
        ),
      );
    }

    final latest = state.valueOrNull ?? current;
    state = AsyncData(
      latest.copyWith(
        isSubmitting: false,
        hasUnsavedChanges: false,
        clearSaveError: true,
      ),
    );
    return localId;
  }

  void _updateAnswer(String questionId, Object value) {
    final current = state.valueOrNull;
    if (current == null) return;
    final updated = Map<String, Object>.from(current.answers)
      ..[questionId] = value;
    state = AsyncData(
      current.copyWith(
        answers: Map<String, Object>.unmodifiable(updated),
        hasUnsavedChanges: true,
        clearSaveError: true,
      ),
    );
    _scheduleAutosave();
  }

  void _removeAnswer(String questionId) {
    final current = state.valueOrNull;
    if (current == null) return;
    final updated = Map<String, Object>.from(current.answers)
      ..remove(questionId);
    state = AsyncData(
      current.copyWith(
        answers: Map<String, Object>.unmodifiable(updated),
        hasUnsavedChanges: true,
        clearSaveError: true,
      ),
    );
    _scheduleAutosave();
  }

  void _scheduleAutosave() {
    _autosaveTimer?.cancel();
    _autosaveTimer = Timer(const Duration(seconds: 20), () {
      unawaited(saveProgress());
    });
  }
}

bool _answerValueHasContent(Object value) {
  if (value is String) return value.trim().isNotEmpty;
  if (value is List) return value.isNotEmpty;
  return true;
}

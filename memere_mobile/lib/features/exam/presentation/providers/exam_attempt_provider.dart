import 'dart:async';

import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/constants/app_constants.dart';
import '../../../../core/errors/failures.dart';
import '../../../../core/grading/offline_grader.dart';
import '../../../../core/offline/offline_attempt_factory.dart';
import '../../../../core/offline/offline_providers.dart';
import '../../../../core/storage/hive/models/downloaded_item.dart';
import '../../../../core/storage/hive/models/pending_submission.dart';
import '../../../auth/presentation/providers/auth_state_provider.dart';
import '../../domain/entities/exam_answer_payload.dart';
import '../../domain/entities/exam_attempt_entity.dart';
import '../../domain/usecases/save_exam_progress_usecase.dart';
import '../../domain/usecases/submit_exam_usecase.dart';
import 'exam_providers.dart';

/// Holds a freshly started attempt so the attempt screen can adopt it without
/// starting a second attempt. Set by the catalog card right before navigating.
final pendingExamAttemptProvider =
    StateProvider<ExamAttemptEntity?>((_) => null);

class ExamAttemptParams {
  const ExamAttemptParams({
    required this.attemptId,
    required this.examId,
  });

  final String attemptId;
  final String examId;

  @override
  bool operator ==(Object other) {
    return other is ExamAttemptParams &&
        other.attemptId == attemptId &&
        other.examId == examId;
  }

  @override
  int get hashCode => Object.hash(attemptId, examId);
}

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

  final ExamAttemptEntity attempt;
  final int currentIndex;
  final ExamAnswerPayload answers;
  final bool isSaving;
  final DateTime? lastSavedAt;
  final String? saveError;
  final bool isSubmitting;
  final bool hasUnsavedChanges;

  int get questionCount => attempt.questions.length;
  int get answeredCount => answers.values.where(_answerValueHasContent).length;
  int get unansweredCount => questionCount - answeredCount;

  ExamAttemptState copyWith({
    ExamAttemptEntity? attempt,
    int? currentIndex,
    ExamAnswerPayload? answers,
    bool? isSaving,
    DateTime? lastSavedAt,
    String? saveError,
    bool? isSubmitting,
    bool? hasUnsavedChanges,
    bool clearSaveError = false,
  }) {
    return ExamAttemptState(
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

final examAttemptProvider = AsyncNotifierProviderFamily<ExamAttemptNotifier,
    ExamAttemptState, ExamAttemptParams>(ExamAttemptNotifier.new);

class ExamAttemptNotifier
    extends FamilyAsyncNotifier<ExamAttemptState, ExamAttemptParams> {
  Timer? _autosaveTimer;

  @override
  Future<ExamAttemptState> build(ExamAttemptParams arg) async {
    ref.onDispose(() {
      _autosaveTimer?.cancel();
      unawaited(saveProgress(force: true));
    });

    final pending = ref.read(pendingExamAttemptProvider);
    if (pending != null && pending.attemptId == arg.attemptId) {
      // Adopt the attempt that the catalog already started; clear the handoff.
      // This is also the offline/guest happy path — the start site stashes the
      // locally-synthesized attempt here before navigating.
      ref.read(pendingExamAttemptProvider.notifier).state = null;
      return ExamAttemptState(attempt: pending);
    }

    // Local attempt whose handoff was lost (rebuild / deep link): re-synthesize
    // from the downloaded copy instead of hitting the network.
    if (isLocalAttemptId(arg.attemptId)) {
      final exam = ref.read(downloadStoreProvider).getOfflineExam(arg.examId);
      if (exam == null) {
        throw const CacheFailure(
          'This exam is not downloaded. Download it first to take it offline.',
        );
      }
      return ExamAttemptState(
        attempt: buildLocalExamAttempt(exam, attemptId: arg.attemptId),
      );
    }

    // Fallback (deep link / retry): start a fresh attempt from the exam id.
    final result = await ref.read(startExamUseCaseProvider)(arg.examId);
    final attempt = result.fold((failure) => throw failure, (value) => value);
    return ExamAttemptState(attempt: attempt);
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
    if (values.isEmpty) {
      _removeAnswer(questionId);
      return;
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
    final result = await ref.read(saveExamProgressUseCaseProvider)(
      SaveExamProgressParams(
        attemptId: current.attempt.attemptId,
        answers: current.answers,
      ),
    );
    final latest = state.valueOrNull ?? current;
    result.fold(
      (failure) {
        state = AsyncData(
          latest.copyWith(isSaving: false, saveError: failure.message),
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

    final result = await ref.read(submitExamUseCaseProvider)(
      SubmitExamParams(
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
                  DownloadType.exam,
                  latestBeforeSubmit.attempt.examId,
                )) {
          return _submitLocally(
            latest,
            serverAttemptId: latestBeforeSubmit.attempt.attemptId,
          );
        }
        state = AsyncData(
          latest.copyWith(isSubmitting: false, saveError: failure.message),
        );
        return null;
      },
      (examResult) async {
        state = AsyncData(
          latest.copyWith(
            isSubmitting: false,
            hasUnsavedChanges: false,
            clearSaveError: true,
          ),
        );
        return examResult.attemptId;
      },
    );
  }

  /// Grades a downloaded exam on-device, stores the provisional result under the
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
    ExamAttemptState current, {
    String? serverAttemptId,
  }) async {
    final localId = serverAttemptId == null
        ? current.attempt.attemptId
        : newLocalAttemptId();
    final exam =
        ref.read(downloadStoreProvider).getOfflineExam(current.attempt.examId);
    if (exam == null) {
      final latest = state.valueOrNull ?? current;
      state = AsyncData(
        latest.copyWith(
          isSubmitting: false,
          saveError:
              'This exam is no longer downloaded. Reconnect or re-download to submit.',
        ),
      );
      return null;
    }

    final rawAnswers = Map<String, dynamic>.from(current.answers);
    final result = OfflineGrader.gradeExam(
      exam: exam,
      submission: offlineSubmissionFromRaw(exam.questions, rawAnswers),
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
          kind: SubmissionKind.exam,
          contentId: exam.id,
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
    _autosaveTimer = Timer(
      const Duration(seconds: AppConstants.examAutoSaveIntervalSeconds),
      () => unawaited(saveProgress()),
    );
  }
}

bool _answerValueHasContent(Object value) {
  if (value is String) return value.trim().isNotEmpty;
  if (value is List) return value.isNotEmpty;
  return true;
}

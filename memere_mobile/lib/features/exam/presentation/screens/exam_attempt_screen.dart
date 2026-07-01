import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../../../core/errors/failures.dart';
import '../../../../core/router/app_router.dart';
import '../../../../shared/widgets/app_button.dart';
import '../../domain/entities/exam_question_entity.dart';
import '../providers/exam_attempt_provider.dart';
import '../providers/exam_timer_provider.dart';
import '../widgets/exam_question_card.dart';
import '../widgets/exam_question_palette.dart';
import '../widgets/exam_save_status.dart';
import '../widgets/exam_timer_badge.dart';

class ExamAttemptScreen extends ConsumerStatefulWidget {
  const ExamAttemptScreen({
    super.key,
    required this.attemptId,
    required this.examId,
  });

  final String attemptId;
  final String examId;

  @override
  ConsumerState<ExamAttemptScreen> createState() => _ExamAttemptScreenState();
}

class _ExamAttemptScreenState extends ConsumerState<ExamAttemptScreen> {
  bool _autoSubmitTriggered = false;

  ExamAttemptParams get _params =>
      ExamAttemptParams(attemptId: widget.attemptId, examId: widget.examId);

  @override
  Widget build(BuildContext context) {
    final attemptAsync = ref.watch(examAttemptProvider(_params));

    return PopScope(
      canPop: false,
      onPopInvokedWithResult: (didPop, _) async {
        if (didPop) return;
        final router = GoRouter.of(context);
        if (await _confirmLeave() && mounted) router.pop();
      },
      child: Scaffold(
        backgroundColor: AppColors.bgPrimary,
        appBar: AppBar(
          automaticallyImplyLeading: false,
          leading: IconButton(
            tooltip: 'Close',
            onPressed: () async {
              final router = GoRouter.of(context);
              if (await _confirmLeave() && mounted) router.pop();
            },
            icon: const Icon(Icons.close_rounded),
          ),
          title: const Text('Mock exam'),
        ),
        body: SafeArea(
          top: false,
          child: attemptAsync.when(
            loading: () => const Center(
              child: CircularProgressIndicator(
                color: AppColors.accentPrimary,
                strokeWidth: 2,
              ),
            ),
            error: (error, _) => _AttemptErrorState(
              message: error is Failure
                  ? error.message
                  : 'Could not start exam. Please try again.',
              onRetry: () => ref.invalidate(examAttemptProvider(_params)),
            ),
            data: (state) {
              if (state.attempt.questions.isEmpty) {
                return const _AttemptErrorState(
                  message: 'This exam has no questions yet.',
                );
              }
              return _AttemptBody(
                params: _params,
                state: state,
                onTimerZero: () => _handleExpiry(state),
                onSubmit: () => _confirmSubmit(state),
              );
            },
          ),
        ),
      ),
    );
  }

  Future<bool> _confirmLeave() async {
    final notifier = ref.read(examAttemptProvider(_params).notifier);
    final leave = await showDialog<bool>(
      context: context,
      builder: (dialogContext) => AlertDialog(
        backgroundColor: AppColors.bgSecondary,
        title: const Text('Leave exam?'),
        content: const Text(
          'The exam timer keeps running on the server. Your latest answers '
          'will be saved before you leave.',
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(dialogContext).pop(false),
            child: const Text('Stay'),
          ),
          TextButton(
            onPressed: () => Navigator.of(dialogContext).pop(true),
            child: const Text('Leave'),
          ),
        ],
      ),
    );
    if (leave == true) {
      await notifier.saveProgress(force: true);
      return true;
    }
    return false;
  }

  Future<void> _handleExpiry(ExamAttemptState state) async {
    if (_autoSubmitTriggered) return;
    _autoSubmitTriggered = true;

    final messenger = ScaffoldMessenger.of(context);
    messenger.showSnackBar(
      const SnackBar(
          content: Text('Time is up. Submitting your saved answers.')),
    );

    final resultAttemptId =
        await ref.read(examAttemptProvider(_params).notifier).submit();
    if (!mounted) return;
    if (resultAttemptId != null) {
      context.go(AppRoutes.examResultPath(resultAttemptId));
    }
  }

  Future<void> _confirmSubmit(ExamAttemptState state) async {
    final seconds = ref.read(examTimerProvider(_seedFor(state))).valueOrNull;

    final submit = await showModalBottomSheet<bool>(
      context: context,
      backgroundColor: AppColors.bgSecondary,
      builder: (sheetContext) => SafeArea(
        child: Padding(
          padding: const EdgeInsets.all(AppSizes.screenPaddingH),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              const Text('Submit exam?', style: AppTextStyles.headlineSmall),
              const SizedBox(height: AppSizes.sm),
              Text(
                '${state.answeredCount} answered · '
                '${state.unansweredCount} unanswered',
                style: AppTextStyles.bodyMedium,
              ),
              if (seconds != null) ...[
                const SizedBox(height: AppSizes.xs),
                Text(
                  'Time remaining ${formatExamTimer(seconds)}',
                  style: AppTextStyles.bodySmall.copyWith(
                    color: AppColors.textSecondary,
                  ),
                ),
              ],
              const SizedBox(height: AppSizes.sm),
              Text(
                'Submission is final and cannot be undone.',
                style: AppTextStyles.bodySmall.copyWith(
                  color: AppColors.textSecondary,
                ),
              ),
              const SizedBox(height: AppSizes.lg),
              Row(
                children: [
                  Expanded(
                    child: AppButton(
                      label: 'Cancel',
                      onPressed: () => Navigator.of(sheetContext).pop(false),
                      variant: AppButtonVariant.secondary,
                    ),
                  ),
                  const SizedBox(width: AppSizes.md),
                  Expanded(
                    child: AppButton(
                      label: 'Submit',
                      onPressed: () => Navigator.of(sheetContext).pop(true),
                    ),
                  ),
                ],
              ),
            ],
          ),
        ),
      ),
    );

    if (submit != true) return;

    final resultAttemptId =
        await ref.read(examAttemptProvider(_params).notifier).submit();
    if (!mounted) return;
    if (resultAttemptId == null) {
      final message =
          ref.read(examAttemptProvider(_params)).valueOrNull?.saveError ??
              'Could not submit exam. Please try again.';
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text(message)),
      );
      return;
    }
    context.go(AppRoutes.examResultPath(resultAttemptId));
  }

  ExamTimerSeed _seedFor(ExamAttemptState state) => ExamTimerSeed(
        attemptId: state.attempt.attemptId,
        expiresAt: state.attempt.expiresAt,
        remainingSeconds: state.attempt.remainingSeconds,
      );
}

class _AttemptBody extends ConsumerWidget {
  const _AttemptBody({
    required this.params,
    required this.state,
    required this.onTimerZero,
    required this.onSubmit,
  });

  final ExamAttemptParams params;
  final ExamAttemptState state;
  final VoidCallback onTimerZero;
  final Future<void> Function() onSubmit;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final notifier = ref.read(examAttemptProvider(params).notifier);
    final seed = ExamTimerSeed(
      attemptId: state.attempt.attemptId,
      expiresAt: state.attempt.expiresAt,
      remainingSeconds: state.attempt.remainingSeconds,
    );
    final seconds = ref.watch(examTimerProvider(seed)).valueOrNull;

    ref.listen(examTimerProvider(seed), (_, next) {
      if (next.valueOrNull == 0) onTimerZero();
    });

    final currentQuestion = state.attempt.questions[state.currentIndex];
    final answeredIndexes = _answeredIndexes(state);

    return ListView(
      padding: const EdgeInsets.all(AppSizes.screenPaddingH),
      children: [
        Row(
          children: [
            Expanded(
              child: Text(
                'Question ${state.currentIndex + 1} of ${state.questionCount}',
                style: AppTextStyles.titleMedium,
              ),
            ),
            ExamTimerBadge(seconds: seconds),
          ],
        ),
        const SizedBox(height: AppSizes.sm),
        LinearProgressIndicator(
          minHeight: 6,
          value: (state.currentIndex + 1) / state.questionCount,
          backgroundColor: AppColors.bgTertiary,
          valueColor: const AlwaysStoppedAnimation(AppColors.accentPrimary),
        ),
        const SizedBox(height: AppSizes.lg),
        ExamQuestionCard(
          question: currentQuestion,
          selectedAnswer: state.answers[currentQuestion.questionId],
          onSingleAnswerSelected: notifier.selectSingleAnswer,
          onMultiAnswerToggled: notifier.toggleMultiAnswer,
          onShortAnswerChanged: notifier.setShortAnswer,
        ),
        const SizedBox(height: AppSizes.md),
        ExamSaveStatus(
          isSaving: state.isSaving,
          saveError: state.saveError,
          lastSavedAt: state.lastSavedAt,
        ),
        const SizedBox(height: AppSizes.lg),
        ExamQuestionPalette(
          count: state.questionCount,
          currentIndex: state.currentIndex,
          answeredIndexes: answeredIndexes,
          onSelected: notifier.goToQuestion,
        ),
        const SizedBox(height: AppSizes.lg),
        Row(
          children: [
            Expanded(
              child: AppButton(
                label: 'Previous',
                onPressed:
                    state.currentIndex == 0 ? null : notifier.previousQuestion,
                variant: AppButtonVariant.secondary,
              ),
            ),
            const SizedBox(width: AppSizes.md),
            Expanded(
              child: state.currentIndex == state.questionCount - 1
                  ? AppButton(
                      label: 'Submit',
                      isLoading: state.isSubmitting,
                      onPressed: () => onSubmit(),
                    )
                  : AppButton(
                      label: 'Next',
                      onPressed: notifier.nextQuestion,
                      suffixIcon: Icons.arrow_forward_rounded,
                    ),
            ),
          ],
        ),
        const SizedBox(height: AppSizes.md),
        AppButton(
          label: 'Submit exam',
          isLoading: state.isSubmitting,
          onPressed: () => onSubmit(),
          variant: AppButtonVariant.outline,
        ),
        if (seconds == 0) ...[
          const SizedBox(height: AppSizes.md),
          Text(
            'Time is up. Submitting your saved answers.',
            style: AppTextStyles.bodySmall.copyWith(color: AppColors.error),
          ),
        ],
        const SizedBox(height: AppSizes.lg),
      ],
    );
  }
}

class _AttemptErrorState extends StatelessWidget {
  const _AttemptErrorState({
    required this.message,
    this.onRetry,
  });

  final String message;
  final VoidCallback? onRetry;

  @override
  Widget build(BuildContext context) {
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(AppSizes.screenPaddingH),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            const Icon(
              Icons.assignment_late_outlined,
              color: AppColors.textSecondary,
              size: AppSizes.iconXl,
            ),
            const SizedBox(height: AppSizes.md),
            Text(
              message,
              textAlign: TextAlign.center,
              style: AppTextStyles.bodyMedium,
            ),
            if (onRetry != null) ...[
              const SizedBox(height: AppSizes.lg),
              AppButton(
                label: 'Retry',
                onPressed: onRetry,
                variant: AppButtonVariant.secondary,
              ),
            ],
          ],
        ),
      ),
    );
  }
}

Set<int> _answeredIndexes(ExamAttemptState state) {
  final indexes = <int>{};
  for (var i = 0; i < state.attempt.questions.length; i++) {
    final question = state.attempt.questions[i];
    final answer = state.answers[question.questionId];
    if (_hasAnswer(question, answer)) indexes.add(i);
  }
  return indexes;
}

bool _hasAnswer(ExamQuestionEntity question, Object? answer) {
  if (answer == null) return false;
  if (question.type == ExamQuestionType.shortAnswer) {
    return answer is String && answer.trim().isNotEmpty;
  }
  if (answer is List) return answer.isNotEmpty;
  return answer is String && answer.trim().isNotEmpty;
}

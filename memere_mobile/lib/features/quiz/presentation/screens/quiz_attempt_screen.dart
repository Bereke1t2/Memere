import 'dart:async';

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../../../core/errors/failures.dart';
import '../../../../core/router/app_router.dart';
import '../../../../shared/widgets/app_button.dart';
import '../../domain/entities/quiz_question_entity.dart';
import '../providers/quiz_attempt_provider.dart';
import '../providers/quiz_timer_provider.dart';
import '../widgets/question_navigator.dart';
import '../widgets/quiz_question_card.dart';
import '../widgets/quiz_timer_badge.dart';

class QuizAttemptScreen extends ConsumerWidget {
  const QuizAttemptScreen({
    super.key,
    required this.attemptId,
    required this.quizId,
  });

  final String attemptId;
  final String quizId;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final params = QuizAttemptParams(attemptId: attemptId, quizId: quizId);
    final attemptAsync = ref.watch(quizAttemptProvider(params));

    return PopScope(
      canPop: false,
      onPopInvokedWithResult: (didPop, _) async {
        if (didPop) return;
        if (await _confirmLeave(context, ref, params) && context.mounted) {
          context.pop();
        }
      },
      child: Scaffold(
        backgroundColor: AppColors.bgPrimary,
        appBar: AppBar(
          leading: IconButton(
            tooltip: 'Close',
            onPressed: () async {
              if (await _confirmLeave(context, ref, params) &&
                  context.mounted) {
                context.pop();
              }
            },
            icon: const Icon(Icons.close_rounded),
          ),
          title: const Text('Quiz attempt'),
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
                  : 'Could not start quiz. Please try again.',
              onRetry: () => ref.invalidate(quizAttemptProvider(params)),
            ),
            data: (state) {
              if (state.attempt.questions.isEmpty) {
                return const _AttemptErrorState(
                  message: 'This quiz has no questions yet.',
                );
              }

              final seed = QuizTimerSeed(
                attemptId: state.attempt.attemptId,
                expiresAt: state.attempt.expiresAt,
                remainingSeconds: state.attempt.remainingSeconds,
              );
              final seconds = ref.watch(quizTimerProvider(seed)).valueOrNull;
              final currentQuestion =
                  state.attempt.questions[state.currentIndex];
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
                      QuizTimerBadge(seconds: seconds),
                    ],
                  ),
                  const SizedBox(height: AppSizes.sm),
                  LinearProgressIndicator(
                    minHeight: 6,
                    value: (state.currentIndex + 1) / state.questionCount,
                    backgroundColor: AppColors.bgTertiary,
                    valueColor: const AlwaysStoppedAnimation(
                      AppColors.accentPrimary,
                    ),
                  ),
                  const SizedBox(height: AppSizes.lg),
                  QuizQuestionCard(
                    question: currentQuestion,
                    selectedAnswer: state.answers[currentQuestion.id],
                    onSingleAnswerSelected: ref
                        .read(quizAttemptProvider(params).notifier)
                        .selectSingleAnswer,
                    onMultiAnswerToggled: ref
                        .read(quizAttemptProvider(params).notifier)
                        .toggleMultiAnswer,
                    onShortAnswerChanged: ref
                        .read(quizAttemptProvider(params).notifier)
                        .setShortAnswer,
                  ),
                  const SizedBox(height: AppSizes.md),
                  _SaveStatus(state: state),
                  const SizedBox(height: AppSizes.lg),
                  QuestionNavigator(
                    count: state.questionCount,
                    currentIndex: state.currentIndex,
                    answeredIndexes: answeredIndexes,
                    onSelected: ref
                        .read(quizAttemptProvider(params).notifier)
                        .goToQuestion,
                  ),
                  const SizedBox(height: AppSizes.lg),
                  Row(
                    children: [
                      Expanded(
                        child: AppButton(
                          label: 'Previous',
                          onPressed: state.currentIndex == 0
                              ? null
                              : ref
                                  .read(quizAttemptProvider(params).notifier)
                                  .previousQuestion,
                          variant: AppButtonVariant.secondary,
                        ),
                      ),
                      const SizedBox(width: AppSizes.md),
                      Expanded(
                        child: state.currentIndex == state.questionCount - 1
                            ? AppButton(
                                label: 'Submit',
                                isLoading: state.isSubmitting,
                                onPressed: () => _confirmSubmit(
                                  context,
                                  ref,
                                  params,
                                  state,
                                ),
                              )
                            : AppButton(
                                label: 'Next',
                                onPressed: ref
                                    .read(
                                      quizAttemptProvider(params).notifier,
                                    )
                                    .nextQuestion,
                                suffixIcon: Icons.arrow_forward_rounded,
                              ),
                      ),
                    ],
                  ),
                  if (seconds == 0) ...[
                    const SizedBox(height: AppSizes.md),
                    Text(
                      'Time is up. Submit to finalize this attempt.',
                      style: AppTextStyles.bodySmall.copyWith(
                        color: AppColors.warning,
                      ),
                    ),
                  ],
                ],
              );
            },
          ),
        ),
      ),
    );
  }

  Future<bool> _confirmLeave(
    BuildContext context,
    WidgetRef ref,
    QuizAttemptParams params,
  ) async {
    final state = ref.read(quizAttemptProvider(params)).valueOrNull;
    if (state == null || !state.hasUnsavedChanges) return true;

    final leave = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Leave quiz?'),
        content:
            const Text('Your latest answers will be saved before leaving.'),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(false),
            child: const Text('Cancel'),
          ),
          TextButton(
            onPressed: () => Navigator.of(context).pop(true),
            child: const Text('Leave'),
          ),
        ],
      ),
    );
    if (leave == true) {
      await ref.read(quizAttemptProvider(params).notifier).saveProgress(
            force: true,
          );
      return true;
    }
    return false;
  }

  Future<void> _confirmSubmit(
    BuildContext context,
    WidgetRef ref,
    QuizAttemptParams params,
    QuizAttemptState state,
  ) async {
    final submit = await showModalBottomSheet<bool>(
      context: context,
      backgroundColor: AppColors.bgSecondary,
      builder: (context) => SafeArea(
        child: Padding(
          padding: const EdgeInsets.all(AppSizes.screenPaddingH),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              const Text('Submit quiz?', style: AppTextStyles.headlineSmall),
              const SizedBox(height: AppSizes.sm),
              Text(
                '${state.answeredCount} answered, ${state.unansweredCount} unanswered. Submission is final.',
                style: AppTextStyles.bodySmall,
              ),
              const SizedBox(height: AppSizes.lg),
              Row(
                children: [
                  Expanded(
                    child: AppButton(
                      label: 'Cancel',
                      onPressed: () => Navigator.of(context).pop(false),
                      variant: AppButtonVariant.secondary,
                    ),
                  ),
                  const SizedBox(width: AppSizes.md),
                  Expanded(
                    child: AppButton(
                      label: 'Submit',
                      onPressed: () => Navigator.of(context).pop(true),
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
        await ref.read(quizAttemptProvider(params).notifier).submit();
    if (!context.mounted) return;
    if (resultAttemptId == null) {
      final message =
          ref.read(quizAttemptProvider(params)).valueOrNull?.saveError ??
              'Could not submit quiz. Please try again.';
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text(message)),
      );
      return;
    }
    context.go(AppRoutes.quizResultPath(resultAttemptId));
  }
}

class _SaveStatus extends StatelessWidget {
  const _SaveStatus({required this.state});

  final QuizAttemptState state;

  @override
  Widget build(BuildContext context) {
    final text = state.isSaving
        ? 'Saving...'
        : state.saveError != null
            ? 'Autosave failed: ${state.saveError}'
            : state.lastSavedAt != null
                ? 'Saved'
                : 'Autosave ready';
    final color = state.saveError != null
        ? AppColors.warning
        : state.isSaving
            ? AppColors.accentPrimary
            : AppColors.textSecondary;
    return Text(
      text,
      style: AppTextStyles.caption.copyWith(color: color),
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
              Icons.quiz_outlined,
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

Set<int> _answeredIndexes(QuizAttemptState state) {
  final indexes = <int>{};
  for (var i = 0; i < state.attempt.questions.length; i++) {
    final question = state.attempt.questions[i];
    final answer = state.answers[question.id];
    if (_hasAnswer(question, answer)) indexes.add(i);
  }
  return indexes;
}

bool _hasAnswer(QuizQuestionEntity question, Object? answer) {
  if (answer == null) return false;
  if (question.type == QuizQuestionType.shortAnswer) {
    return answer is String && answer.trim().isNotEmpty;
  }
  if (answer is List) return answer.isNotEmpty;
  return answer is String && answer.trim().isNotEmpty;
}

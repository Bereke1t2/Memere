import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../../../core/errors/failures.dart';
import '../../../../core/router/app_router.dart';
import '../../../../shared/utils/formatters.dart';
import '../../../../shared/widgets/app_button.dart';
import '../providers/quiz_detail_provider.dart';
import '../providers/quiz_providers.dart';

class QuizDetailScreen extends ConsumerWidget {
  const QuizDetailScreen({
    super.key,
    required this.quizId,
  });

  final String quizId;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final quizAsync = ref.watch(quizDetailProvider(quizId));
    final isStarting = ref.watch(quizStartLoadingProvider(quizId));

    return Scaffold(
      backgroundColor: AppColors.bgPrimary,
      appBar: AppBar(title: const Text('Quiz')),
      body: SafeArea(
        top: false,
        child: quizAsync.when(
          loading: () => const Center(
            child: CircularProgressIndicator(
              color: AppColors.accentPrimary,
              strokeWidth: 2,
            ),
          ),
          error: (error, _) => _QuizErrorState(
            title: 'Could not load quiz',
            message: error is Failure ? error.message : 'Please try again.',
            onRetry: () => ref.invalidate(quizDetailProvider(quizId)),
          ),
          data: (quiz) => ListView(
            padding: const EdgeInsets.all(AppSizes.screenPaddingH),
            children: [
              Container(
                padding: const EdgeInsets.all(AppSizes.lg),
                decoration: BoxDecoration(
                  color: AppColors.bgSecondary,
                  borderRadius: BorderRadius.circular(AppSizes.radiusMd),
                  border: Border.all(color: AppColors.border),
                ),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Container(
                      width: 48,
                      height: 48,
                      decoration: BoxDecoration(
                        color: AppColors.accentGlow,
                        borderRadius: BorderRadius.circular(AppSizes.radiusSm),
                      ),
                      child: const Icon(
                        Icons.quiz_outlined,
                        color: AppColors.accentPrimary,
                      ),
                    ),
                    const SizedBox(height: AppSizes.md),
                    Text(quiz.title, style: AppTextStyles.headlineMedium),
                    const SizedBox(height: AppSizes.sm),
                    Text(
                      quiz.hasAttemptsRemaining
                          ? 'Answer each question, then submit for server grading.'
                          : 'No attempts remaining.',
                      style: AppTextStyles.bodySmall,
                    ),
                    const SizedBox(height: AppSizes.lg),
                    _MetaRow(
                      icon: Icons.help_outline_rounded,
                      label: 'Questions',
                      value: quiz.questionCount.toString(),
                    ),
                    _MetaRow(
                      icon: Icons.percent_rounded,
                      label: 'Pass mark',
                      value: '${quiz.passPercentage.toStringAsFixed(0)}%',
                    ),
                    _MetaRow(
                      icon: Icons.timer_outlined,
                      label: 'Time limit',
                      value: quiz.timeLimitSeconds == null
                          ? 'Untimed'
                          : formatDurationSeconds(quiz.timeLimitSeconds!),
                    ),
                    _MetaRow(
                      icon: Icons.replay_rounded,
                      label: 'Attempts',
                      value: quiz.maxAttempts == null
                          ? '${quiz.attemptsUsed} used'
                          : '${quiz.attemptsUsed}/${quiz.maxAttempts}',
                    ),
                  ],
                ),
              ),
              const SizedBox(height: AppSizes.lg),
              AppButton(
                label: quiz.hasAttemptsRemaining
                    ? 'Start quiz'
                    : 'Attempts exhausted',
                isLoading: isStarting,
                isDisabled: !quiz.hasAttemptsRemaining,
                onPressed: () => _startAttempt(context, ref),
                suffixIcon: Icons.arrow_forward_rounded,
              ),
            ],
          ),
        ),
      ),
    );
  }

  Future<void> _startAttempt(BuildContext context, WidgetRef ref) async {
    ref.read(quizStartLoadingProvider(quizId).notifier).state = true;
    final result = await ref.read(startQuizAttemptUseCaseProvider)(quizId);
    ref.read(quizStartLoadingProvider(quizId).notifier).state = false;

    result.fold(
      (failure) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text(failure.message)),
        );
      },
      (attempt) {
        context.go(
          AppRoutes.quizAttemptPath(
            attemptId: attempt.attemptId,
            quizId: attempt.quizId,
          ),
        );
      },
    );
  }
}

class _MetaRow extends StatelessWidget {
  const _MetaRow({
    required this.icon,
    required this.label,
    required this.value,
  });

  final IconData icon;
  final String label;
  final String value;

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.only(bottom: AppSizes.sm),
      child: Row(
        children: [
          Icon(icon, size: AppSizes.iconSm, color: AppColors.textSecondary),
          const SizedBox(width: AppSizes.sm),
          Expanded(
            child: Text(
              label,
              style: AppTextStyles.bodySmall,
            ),
          ),
          Text(value, style: AppTextStyles.labelMedium),
        ],
      ),
    );
  }
}

class _QuizErrorState extends StatelessWidget {
  const _QuizErrorState({
    required this.title,
    required this.message,
    required this.onRetry,
  });

  final String title;
  final String message;
  final VoidCallback onRetry;

  @override
  Widget build(BuildContext context) {
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(AppSizes.screenPaddingH),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            const Icon(
              Icons.error_outline_rounded,
              color: AppColors.error,
              size: AppSizes.iconXl,
            ),
            const SizedBox(height: AppSizes.md),
            Text(title, style: AppTextStyles.headlineSmall),
            const SizedBox(height: AppSizes.sm),
            Text(
              message,
              textAlign: TextAlign.center,
              style: AppTextStyles.bodySmall,
            ),
            const SizedBox(height: AppSizes.lg),
            AppButton(
              label: 'Retry',
              onPressed: onRetry,
              variant: AppButtonVariant.secondary,
            ),
          ],
        ),
      ),
    );
  }
}

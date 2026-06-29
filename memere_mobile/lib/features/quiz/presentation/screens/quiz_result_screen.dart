import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../../../core/errors/failures.dart';
import '../../../../core/router/app_router.dart';
import '../../../../shared/widgets/app_button.dart';
import '../providers/quiz_result_provider.dart';
import '../widgets/question_feedback_tile.dart';
import '../widgets/quiz_result_skeleton.dart';
import '../widgets/quiz_score_summary.dart';
import '../widgets/subject_breakdown_card.dart';

class QuizResultScreen extends ConsumerWidget {
  const QuizResultScreen({
    super.key,
    required this.attemptId,
  });

  final String attemptId;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final resultAsync = ref.watch(quizResultProvider(attemptId));

    return Scaffold(
      backgroundColor: AppColors.bgPrimary,
      appBar: AppBar(
        title: const Text('Quiz result'),
        leading: IconButton(
          tooltip: 'Close',
          onPressed: () => context.go(AppRoutes.home),
          icon: const Icon(Icons.close_rounded),
        ),
      ),
      body: SafeArea(
        top: false,
        child: resultAsync.when(
          loading: () => const QuizResultSkeleton(),
          error: (error, _) => _ResultErrorState(
            message: error is Failure
                ? error.message
                : 'Could not load quiz result. Please try again.',
            onRetry: () => ref.invalidate(quizResultProvider(attemptId)),
          ),
          data: (result) => ListView(
            padding: const EdgeInsets.all(AppSizes.screenPaddingH),
            children: [
              QuizScoreSummary(result: result),
              const SizedBox(height: AppSizes.md),
              Container(
                padding: const EdgeInsets.all(AppSizes.md),
                decoration: BoxDecoration(
                  color: AppColors.bgSecondary,
                  borderRadius: BorderRadius.circular(AppSizes.radiusMd),
                  border: Border.all(color: AppColors.border),
                ),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    _ResultMeta(
                      label: 'Attempt',
                      value: result.attemptNumber.toString(),
                    ),
                    _ResultMeta(
                      label: 'Submitted',
                      value: _formatSubmittedAt(result.submittedAt),
                    ),
                  ],
                ),
              ),
              const SizedBox(height: AppSizes.md),
              SubjectBreakdownCard(
                subjectBreakdown: result.subjectBreakdown,
              ),
              const SizedBox(height: AppSizes.lg),
              const Text('Review', style: AppTextStyles.headlineSmall),
              const SizedBox(height: AppSizes.md),
              ...result.feedback.asMap().entries.map(
                    (entry) => QuestionFeedbackTile(
                      feedback: entry.value,
                      index: entry.key,
                    ),
                  ),
              const SizedBox(height: AppSizes.lg),
              AppButton(
                label: 'Back to courses',
                onPressed: () => context.go(AppRoutes.home),
                variant: AppButtonVariant.secondary,
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class _ResultMeta extends StatelessWidget {
  const _ResultMeta({
    required this.label,
    required this.value,
  });

  final String label;
  final String value;

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.only(bottom: AppSizes.sm),
      child: Row(
        children: [
          Expanded(
            child: Text(label, style: AppTextStyles.bodySmall),
          ),
          Text(value, style: AppTextStyles.labelMedium),
        ],
      ),
    );
  }
}

class _ResultErrorState extends StatelessWidget {
  const _ResultErrorState({
    required this.message,
    required this.onRetry,
  });

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
            Text(
              message,
              textAlign: TextAlign.center,
              style: AppTextStyles.bodyMedium,
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

String _formatSubmittedAt(DateTime? value) {
  if (value == null) return 'Not available';
  final local = value.toLocal();
  final date =
      '${local.year}-${local.month.toString().padLeft(2, '0')}-${local.day.toString().padLeft(2, '0')}';
  final time =
      '${local.hour.toString().padLeft(2, '0')}:${local.minute.toString().padLeft(2, '0')}';
  return '$date $time';
}

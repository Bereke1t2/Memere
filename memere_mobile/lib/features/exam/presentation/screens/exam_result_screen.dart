import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../../../core/errors/failures.dart';
import '../../../../core/router/app_router.dart';
import '../../../../shared/widgets/app_button.dart';
import '../../domain/entities/exam_subject_score_entity.dart';
import '../providers/exam_result_provider.dart';
import '../widgets/exam_question_feedback_tile.dart';
import '../widgets/exam_result_skeleton.dart';
import '../widgets/exam_score_summary.dart';
import '../widgets/exam_subject_breakdown.dart';

class ExamResultScreen extends ConsumerWidget {
  const ExamResultScreen({
    super.key,
    required this.attemptId,
  });

  final String attemptId;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final resultAsync = ref.watch(examResultProvider(attemptId));

    return Scaffold(
      backgroundColor: AppColors.bgPrimary,
      appBar: AppBar(
        title: const Text('Exam result'),
        leading: IconButton(
          tooltip: 'Close',
          onPressed: () => context.go(AppRoutes.mockExams),
          icon: const Icon(Icons.close_rounded),
        ),
      ),
      body: SafeArea(
        top: false,
        child: resultAsync.when(
          loading: () => const ExamResultSkeleton(),
          error: (error, _) => _ResultErrorState(
            message: error is Failure
                ? error.message
                : 'Could not load exam result. Please try again.',
            onRetry: () => ref.invalidate(examResultProvider(attemptId)),
          ),
          data: (result) {
            final breakdown = result.subjectBreakdown.entries
                .map(
                  (entry) => ExamSubjectScoreEntity(
                    key: entry.key,
                    earned: entry.value.earned,
                    possible: entry.value.possible,
                  ),
                )
                .toList();

            return ListView(
              padding: const EdgeInsets.all(AppSizes.screenPaddingH),
              children: [
                ExamScoreSummary(result: result),
                const SizedBox(height: AppSizes.md),
                Container(
                  padding: const EdgeInsets.all(AppSizes.md),
                  decoration: BoxDecoration(
                    color: AppColors.bgSecondary,
                    borderRadius: BorderRadius.circular(AppSizes.radiusMd),
                    border: Border.all(color: AppColors.border),
                  ),
                  child: _ResultMeta(
                    label: 'Submitted',
                    value: _formatSubmittedAt(result.submittedAt),
                  ),
                ),
                const SizedBox(height: AppSizes.md),
                ExamSubjectBreakdown(scores: breakdown),
                const SizedBox(height: AppSizes.lg),
                Align(
                  alignment: Alignment.centerRight,
                  child: AppButton(
                    label: 'Review analytics',
                    onPressed: () =>
                        context.push(AppRoutes.examAnalyticsPath(attemptId)),
                    suffixIcon: Icons.insights_outlined,
                    width: 184,
                    height: AppSizes.buttonHeightSm,
                  ),
                ),
                const SizedBox(height: AppSizes.lg),
                if (result.feedback.isNotEmpty) ...[
                  const Text('Review', style: AppTextStyles.headlineSmall),
                  const SizedBox(height: AppSizes.md),
                  ...result.feedback.asMap().entries.map(
                        (entry) => ExamQuestionFeedbackTile(
                          feedback: entry.value,
                          index: entry.key,
                        ),
                      ),
                ],
                const SizedBox(height: AppSizes.sm),
                Align(
                  alignment: Alignment.centerRight,
                  child: AppButton(
                    label: 'Back to exams',
                    onPressed: () => context.go(AppRoutes.mockExams),
                    variant: AppButtonVariant.secondary,
                    width: 148,
                    height: AppSizes.buttonHeightSm,
                  ),
                ),
                const SizedBox(height: AppSizes.lg),
              ],
            );
          },
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
    return Row(
      children: [
        Expanded(
          child: Text(label, style: AppTextStyles.bodySmall),
        ),
        Text(value, style: AppTextStyles.labelMedium),
      ],
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
              width: 120,
              height: AppSizes.buttonHeightSm,
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
  final date = '${local.year}-${local.month.toString().padLeft(2, '0')}-'
      '${local.day.toString().padLeft(2, '0')}';
  final time = '${local.hour.toString().padLeft(2, '0')}:'
      '${local.minute.toString().padLeft(2, '0')}';
  return '$date $time';
}

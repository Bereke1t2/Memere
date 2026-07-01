import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../../../core/errors/failures.dart';
import '../../../../shared/widgets/app_button.dart';
import '../providers/exam_analytics_provider.dart';
import '../widgets/exam_result_skeleton.dart';
import '../widgets/exam_subject_breakdown.dart';
import '../widgets/percentile_card.dart';
import '../widgets/weak_area_card.dart';

class ExamAnalyticsScreen extends ConsumerWidget {
  const ExamAnalyticsScreen({
    super.key,
    required this.attemptId,
  });

  final String attemptId;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final analyticsAsync = ref.watch(examAnalyticsProvider(attemptId));

    return Scaffold(
      backgroundColor: AppColors.bgPrimary,
      appBar: AppBar(title: const Text('Exam analytics')),
      body: SafeArea(
        top: false,
        child: analyticsAsync.when(
          loading: () => const ExamResultSkeleton(),
          error: (error, _) => _AnalyticsErrorState(
            message: error is Failure
                ? error.message
                : 'Could not load analytics. Please try again.',
            onRetry: () => ref.invalidate(examAnalyticsProvider(attemptId)),
          ),
          data: (analytics) => ListView(
            padding: const EdgeInsets.all(AppSizes.screenPaddingH),
            children: [
              Container(
                width: double.infinity,
                padding: const EdgeInsets.all(AppSizes.lg),
                decoration: BoxDecoration(
                  color: AppColors.bgSecondary,
                  borderRadius: BorderRadius.circular(AppSizes.radiusMd),
                  border: Border.all(color: AppColors.border),
                ),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const Text('Overall score',
                        style: AppTextStyles.bodySmall),
                    const SizedBox(height: AppSizes.xs),
                    Text(
                      '${analytics.percentage.toStringAsFixed(0)}%',
                      style: AppTextStyles.displayMedium.copyWith(
                        color: AppColors.accentPrimary,
                      ),
                    ),
                    const SizedBox(height: AppSizes.xs),
                    Text(
                      '${_trimNumber(analytics.score)} marks earned',
                      style: AppTextStyles.bodyMedium.copyWith(
                        color: AppColors.textSecondary,
                      ),
                    ),
                  ],
                ),
              ),
              const SizedBox(height: AppSizes.md),
              PercentileCard(percentile: analytics.percentile),
              const SizedBox(height: AppSizes.md),
              ExamSubjectBreakdown(scores: analytics.subjectBreakdown),
              const SizedBox(height: AppSizes.lg),
              if (analytics.weakAreas.isNotEmpty) ...[
                const Text('Weak areas', style: AppTextStyles.headlineSmall),
                const SizedBox(height: AppSizes.md),
                ...analytics.weakAreas.map(
                  (area) => WeakAreaCard(weakArea: area),
                ),
                const SizedBox(height: AppSizes.md),
                Container(
                  padding: const EdgeInsets.all(AppSizes.md),
                  decoration: BoxDecoration(
                    color: AppColors.infoSurface,
                    borderRadius: BorderRadius.circular(AppSizes.radiusMd),
                  ),
                  child: Row(
                    children: [
                      const Icon(
                        Icons.lightbulb_outline_rounded,
                        color: AppColors.info,
                        size: AppSizes.iconSm,
                      ),
                      const SizedBox(width: AppSizes.sm),
                      Expanded(
                        child: Text(
                          'Focus revision on the weak areas above before '
                          'your next attempt.',
                          style: AppTextStyles.bodySmall.copyWith(
                            color: AppColors.info,
                          ),
                        ),
                      ),
                    ],
                  ),
                ),
              ],
              const SizedBox(height: AppSizes.lg),
              AppButton(
                label: 'Back to result',
                onPressed: () => context.pop(),
                variant: AppButtonVariant.secondary,
              ),
              const SizedBox(height: AppSizes.lg),
            ],
          ),
        ),
      ),
    );
  }
}

class _AnalyticsErrorState extends StatelessWidget {
  const _AnalyticsErrorState({
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
              Icons.insights_outlined,
              color: AppColors.textSecondary,
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

String _trimNumber(double value) {
  return value % 1 == 0 ? value.toInt().toString() : value.toStringAsFixed(1);
}

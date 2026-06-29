import 'package:flutter/material.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../domain/entities/quiz_result_entity.dart';

class QuizScoreSummary extends StatelessWidget {
  const QuizScoreSummary({
    super.key,
    required this.result,
  });

  final QuizResultEntity result;

  @override
  Widget build(BuildContext context) {
    final color = result.passed ? AppColors.success : AppColors.warning;
    return Container(
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
          Icon(
            result.passed
                ? Icons.check_circle_rounded
                : Icons.flag_circle_rounded,
            color: color,
            size: AppSizes.iconLg,
          ),
          const SizedBox(height: AppSizes.md),
          Text(
            result.passed ? 'Passed' : 'Keep practicing',
            style: AppTextStyles.headlineMedium.copyWith(color: color),
          ),
          const SizedBox(height: AppSizes.xs),
          Text(
            result.passed
                ? 'Good work. Review the explanations before moving on.'
                : 'Review the missed questions and try again when ready.',
            style: AppTextStyles.bodySmall,
          ),
          const SizedBox(height: AppSizes.lg),
          Text(
            '${result.percentage.toStringAsFixed(0)}%',
            style: AppTextStyles.displayMedium.copyWith(color: color),
          ),
          const SizedBox(height: AppSizes.xs),
          Text(
            '${_trimNumber(result.score)} / ${result.totalPoints} points',
            style: AppTextStyles.bodyMedium.copyWith(
              color: AppColors.textSecondary,
            ),
          ),
        ],
      ),
    );
  }
}

String _trimNumber(double value) {
  return value % 1 == 0 ? value.toInt().toString() : value.toStringAsFixed(1);
}

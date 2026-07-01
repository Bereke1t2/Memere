import 'package:flutter/material.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../domain/entities/exam_subject_score_entity.dart';

class WeakAreaCard extends StatelessWidget {
  const WeakAreaCard({
    super.key,
    required this.weakArea,
  });

  final ExamSubjectScoreEntity weakArea;

  @override
  Widget build(BuildContext context) {
    final percent = weakArea.percentage;
    final (color, surface) = _severity(percent);

    return Container(
      margin: const EdgeInsets.only(bottom: AppSizes.sm),
      padding: const EdgeInsets.all(AppSizes.md),
      decoration: BoxDecoration(
        color: surface,
        borderRadius: BorderRadius.circular(AppSizes.radiusMd),
        border: Border.all(color: color),
      ),
      child: Row(
        children: [
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(weakArea.key, style: AppTextStyles.titleMedium),
                const SizedBox(height: AppSizes.xs),
                Text(
                  '${weakArea.earned}/${weakArea.possible} marks',
                  style: AppTextStyles.bodySmall.copyWith(
                    color: AppColors.textSecondary,
                  ),
                ),
              ],
            ),
          ),
          Text(
            '${percent.toStringAsFixed(0)}%',
            style: AppTextStyles.titleLarge.copyWith(color: color),
          ),
        ],
      ),
    );
  }

  (Color, Color) _severity(double percent) {
    if (percent < 40) return (AppColors.error, AppColors.errorSurface);
    if (percent < 70) return (AppColors.warning, AppColors.warningSurface);
    return (AppColors.success, AppColors.successSurface);
  }
}

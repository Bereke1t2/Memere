import 'package:flutter/material.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../domain/entities/subject_score_entity.dart';

class SubjectBreakdownCard extends StatelessWidget {
  const SubjectBreakdownCard({
    super.key,
    required this.subjectBreakdown,
  });

  final Map<String, SubjectScoreEntity> subjectBreakdown;

  @override
  Widget build(BuildContext context) {
    if (subjectBreakdown.isEmpty) return const SizedBox.shrink();

    return Container(
      padding: const EdgeInsets.all(AppSizes.md),
      decoration: BoxDecoration(
        color: AppColors.bgSecondary,
        borderRadius: BorderRadius.circular(AppSizes.radiusMd),
        border: Border.all(color: AppColors.border),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text('Subject breakdown', style: AppTextStyles.titleMedium),
          const SizedBox(height: AppSizes.md),
          ...subjectBreakdown.entries.map((entry) {
            final score = entry.value;
            final percent = (score.percentage / 100).clamp(0.0, 1.0);
            return Padding(
              padding: const EdgeInsets.only(bottom: AppSizes.md),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    children: [
                      Expanded(
                        child: Text(entry.key, style: AppTextStyles.bodyMedium),
                      ),
                      Text(
                        '${score.earned}/${score.possible}',
                        style: AppTextStyles.labelMedium.copyWith(
                          color: AppColors.textSecondary,
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: AppSizes.xs),
                  ClipRRect(
                    borderRadius: BorderRadius.circular(AppSizes.radiusFull),
                    child: LinearProgressIndicator(
                      minHeight: 6,
                      value: percent,
                      backgroundColor: AppColors.bgTertiary,
                      valueColor: const AlwaysStoppedAnimation(
                        AppColors.accentPrimary,
                      ),
                    ),
                  ),
                ],
              ),
            );
          }),
        ],
      ),
    );
  }
}

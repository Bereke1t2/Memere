import 'package:flutter/material.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';

class PercentileCard extends StatelessWidget {
  const PercentileCard({
    super.key,
    required this.percentile,
  });

  final double? percentile;

  @override
  Widget build(BuildContext context) {
    return Container(
      width: double.infinity,
      padding: const EdgeInsets.all(AppSizes.md),
      decoration: BoxDecoration(
        color: AppColors.bgSecondary,
        borderRadius: BorderRadius.circular(AppSizes.radiusMd),
        border: Border.all(color: AppColors.border),
      ),
      child: Row(
        children: [
          Container(
            width: 44,
            height: 44,
            decoration: BoxDecoration(
              color: AppColors.infoSurface,
              borderRadius: BorderRadius.circular(AppSizes.radiusSm),
            ),
            child: const Icon(
              Icons.leaderboard_outlined,
              color: AppColors.info,
            ),
          ),
          const SizedBox(width: AppSizes.md),
          Expanded(
            child: percentile == null
                ? Text(
                    'Percentile will appear after more students attempt '
                    'this exam.',
                    style: AppTextStyles.bodySmall.copyWith(
                      color: AppColors.textSecondary,
                    ),
                  )
                : Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      const Text('Percentile', style: AppTextStyles.bodySmall),
                      const SizedBox(height: AppSizes.xs),
                      Text(
                        'Top ${(100 - percentile!).clamp(0, 100).toStringAsFixed(0)}% '
                        '· ${percentile!.toStringAsFixed(1)}th percentile',
                        style: AppTextStyles.titleMedium.copyWith(
                          color: AppColors.info,
                        ),
                      ),
                    ],
                  ),
          ),
        ],
      ),
    );
  }
}

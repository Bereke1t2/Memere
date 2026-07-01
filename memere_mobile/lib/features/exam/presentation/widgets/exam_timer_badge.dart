import 'package:flutter/material.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';

class ExamTimerBadge extends StatelessWidget {
  const ExamTimerBadge({
    super.key,
    required this.seconds,
  });

  final int? seconds;

  @override
  Widget build(BuildContext context) {
    if (seconds == null) {
      return const SizedBox.shrink();
    }
    final value = seconds!;
    final critical = value <= 60;
    final warning = value <= 300;

    final color = critical
        ? AppColors.error
        : warning
            ? AppColors.warning
            : AppColors.textPrimary;
    final surface = critical
        ? AppColors.errorSurface
        : warning
            ? AppColors.warningSurface
            : AppColors.bgTertiary;
    final borderColor = critical
        ? AppColors.error
        : warning
            ? AppColors.warning
            : AppColors.border;

    return Container(
      padding: const EdgeInsets.symmetric(
        horizontal: AppSizes.sm,
        vertical: AppSizes.xs,
      ),
      decoration: BoxDecoration(
        color: surface,
        borderRadius: BorderRadius.circular(AppSizes.radiusFull),
        border: Border.all(color: borderColor),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(
            Icons.timer_outlined,
            size: AppSizes.iconXs,
            color: critical || warning ? color : AppColors.textSecondary,
          ),
          const SizedBox(width: AppSizes.xs),
          Text(
            formatExamTimer(value),
            style: AppTextStyles.labelMedium.copyWith(color: color),
          ),
        ],
      ),
    );
  }
}

/// Formats seconds as `h:mm:ss` (or `m:ss` under an hour).
String formatExamTimer(int seconds) {
  final safe = seconds < 0 ? 0 : seconds;
  final hours = safe ~/ 3600;
  final minutes = (safe % 3600) ~/ 60;
  final secs = safe % 60;
  if (hours > 0) {
    return '$hours:${minutes.toString().padLeft(2, '0')}:'
        '${secs.toString().padLeft(2, '0')}';
  }
  return '$minutes:${secs.toString().padLeft(2, '0')}';
}

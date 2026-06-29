import 'package:flutter/material.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';

class QuizTimerBadge extends StatelessWidget {
  const QuizTimerBadge({
    super.key,
    required this.seconds,
  });

  final int? seconds;

  @override
  Widget build(BuildContext context) {
    if (seconds == null) {
      return const SizedBox.shrink();
    }
    final urgent = seconds! <= 60;
    return Container(
      padding: const EdgeInsets.symmetric(
        horizontal: AppSizes.sm,
        vertical: AppSizes.xs,
      ),
      decoration: BoxDecoration(
        color: urgent ? AppColors.warningSurface : AppColors.bgTertiary,
        borderRadius: BorderRadius.circular(AppSizes.radiusFull),
        border: Border.all(
          color: urgent ? AppColors.warning : AppColors.border,
        ),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(
            Icons.timer_outlined,
            size: AppSizes.iconXs,
            color: urgent ? AppColors.warning : AppColors.textSecondary,
          ),
          const SizedBox(width: AppSizes.xs),
          Text(
            formatQuizTimer(seconds!),
            style: AppTextStyles.labelMedium.copyWith(
              color: urgent ? AppColors.warning : AppColors.textPrimary,
            ),
          ),
        ],
      ),
    );
  }
}

String formatQuizTimer(int seconds) {
  final safe = seconds < 0 ? 0 : seconds;
  final hours = safe ~/ 3600;
  final minutes = (safe % 3600) ~/ 60;
  final secs = safe % 60;
  if (hours > 0) {
    return '$hours:${minutes.toString().padLeft(2, '0')}:${secs.toString().padLeft(2, '0')}';
  }
  return '$minutes:${secs.toString().padLeft(2, '0')}';
}

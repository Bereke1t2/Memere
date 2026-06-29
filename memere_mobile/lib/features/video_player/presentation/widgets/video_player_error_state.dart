import 'package:flutter/material.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../../../shared/widgets/app_button.dart';

class VideoPlayerErrorState extends StatelessWidget {
  const VideoPlayerErrorState({
    super.key,
    required this.title,
    required this.message,
    this.onRetry,
    this.onBack,
  });

  final String title;
  final String message;
  final VoidCallback? onRetry;
  final VoidCallback? onBack;

  @override
  Widget build(BuildContext context) {
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(AppSizes.lg),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Container(
              width: 64,
              height: 64,
              decoration: BoxDecoration(
                color: AppColors.bgTertiary,
                borderRadius: BorderRadius.circular(AppSizes.radiusLg),
              ),
              child: const Icon(
                Icons.play_disabled_rounded,
                color: AppColors.accentPrimary,
              ),
            ),
            const SizedBox(height: AppSizes.md),
            Text(
              title,
              textAlign: TextAlign.center,
              style: AppTextStyles.headlineSmall,
            ),
            const SizedBox(height: AppSizes.sm),
            Text(
              message,
              textAlign: TextAlign.center,
              style: AppTextStyles.bodyMedium.copyWith(
                color: AppColors.textSecondary,
              ),
            ),
            const SizedBox(height: AppSizes.lg),
            if (onRetry != null)
              AppButton(
                label: 'Retry',
                onPressed: onRetry,
                width: 180,
                height: AppSizes.buttonHeightSm,
              ),
            if (onBack != null) ...[
              const SizedBox(height: AppSizes.sm),
              AppButton(
                label: 'Back',
                onPressed: onBack,
                width: 180,
                height: AppSizes.buttonHeightSm,
                variant: AppButtonVariant.ghost,
              ),
            ],
          ],
        ),
      ),
    );
  }
}

import 'package:flutter/material.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../../../shared/widgets/app_button.dart';

class CourseEmptyState extends StatelessWidget {
  const CourseEmptyState({
    super.key,
    required this.icon,
    required this.title,
    required this.body,
    this.buttonLabel,
    this.onPressed,
  });

  final IconData icon;
  final String title;
  final String body;
  final String? buttonLabel;
  final VoidCallback? onPressed;

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
                border: Border.all(color: AppColors.border),
              ),
              child: Icon(icon, color: AppColors.accentPrimary),
            ),
            const SizedBox(height: AppSizes.md),
            Text(
              title,
              textAlign: TextAlign.center,
              style: AppTextStyles.headlineSmall,
            ),
            const SizedBox(height: AppSizes.sm),
            Text(
              body,
              textAlign: TextAlign.center,
              style: AppTextStyles.bodyMedium.copyWith(
                color: AppColors.textSecondary,
              ),
            ),
            if (buttonLabel != null && onPressed != null) ...[
              const SizedBox(height: AppSizes.lg),
              AppButton(
                label: buttonLabel!,
                onPressed: onPressed,
                width: 180,
                height: AppSizes.buttonHeightSm,
              ),
            ],
          ],
        ),
      ),
    );
  }
}

import 'package:flutter/material.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';

class ExamSaveStatus extends StatelessWidget {
  const ExamSaveStatus({
    super.key,
    required this.isSaving,
    required this.saveError,
    required this.lastSavedAt,
  });

  final bool isSaving;
  final String? saveError;
  final DateTime? lastSavedAt;

  @override
  Widget build(BuildContext context) {
    final (icon, text, color) = _status();
    return Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        Icon(icon, size: AppSizes.iconXs, color: color),
        const SizedBox(width: AppSizes.xs),
        Text(text, style: AppTextStyles.caption.copyWith(color: color)),
      ],
    );
  }

  (IconData, String, Color) _status() {
    if (isSaving) {
      return (
        Icons.cloud_sync_outlined,
        'Saving...',
        AppColors.accentPrimary,
      );
    }
    if (saveError != null) {
      return (
        Icons.cloud_off_outlined,
        'Save failed - will retry',
        AppColors.warning,
      );
    }
    if (lastSavedAt != null) {
      return (
        Icons.cloud_done_outlined,
        'Saved',
        AppColors.textSecondary,
      );
    }
    return (
      Icons.cloud_outlined,
      'Autosave ready',
      AppColors.textSecondary,
    );
  }
}

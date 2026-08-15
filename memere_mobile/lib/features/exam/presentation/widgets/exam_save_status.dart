import 'package:flutter/material.dart';

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
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
      decoration: BoxDecoration(
        color: const Color(0xFF111116),
        borderRadius: BorderRadius.circular(8),
        border: Border.all(color: const Color(0xFF1E1E28)),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(icon, size: 12, color: color),
          const SizedBox(width: 5),
          Text(
            text,
            style: AppTextStyles.labelSmall.copyWith(
              color: color,
              fontSize: 11,
              fontWeight: FontWeight.w500,
            ),
          ),
        ],
      ),
    );
  }

  (IconData, String, Color) _status() {
    if (isSaving) {
      return (
        Icons.sync_rounded,
        'Saving...',
        const Color(0xFF10B981),
      );
    }
    if (saveError != null) {
      return (
        Icons.cloud_off_rounded,
        'Offline (will retry)',
        const Color(0xFFF59E0B),
      );
    }
    if (lastSavedAt != null) {
      return (
        Icons.check_circle_outline_rounded,
        'Progress saved',
        const Color(0xFF71717A),
      );
    }
    return (
      Icons.cloud_queue_rounded,
      'Autosave active',
      const Color(0xFF71717A),
    );
  }
}

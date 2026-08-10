import 'package:flutter/material.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import 'download_button.dart';

class VideoActionBar extends StatelessWidget {
  const VideoActionBar({
    super.key,
    required this.videoId,
    required this.lessonId,
    required this.courseId,
    required this.title,
    required this.onRefresh,
    required this.onMarkComplete,
    this.isSavingProgress = false,
    this.isCompleted = false,
  });

  final String videoId;
  final String lessonId;
  final String courseId;
  final String title;
  final VoidCallback onRefresh;
  final VoidCallback onMarkComplete;
  final bool isSavingProgress;
  final bool isCompleted;

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        Row(
          children: [
            Expanded(
              child: ElevatedButton.icon(
                onPressed: lessonId.isEmpty ? null : onMarkComplete,
                style: ElevatedButton.styleFrom(
                  backgroundColor: isCompleted ? AppColors.success : AppColors.textPrimary,
                  foregroundColor: isCompleted ? Colors.white : AppColors.bgPrimary,
                  minimumSize: const Size.fromHeight(46),
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(AppSizes.radiusMd),
                  ),
                ),
                icon: Icon(
                  isCompleted ? Icons.check_circle_rounded : Icons.check_circle_outline_rounded,
                  size: 20,
                ),
                label: Text(
                  isCompleted ? 'Completed ✓' : 'Mark as Complete',
                  style: const TextStyle(fontWeight: FontWeight.bold),
                ),
              ),
            ),
            const SizedBox(width: AppSizes.sm),
            OutlinedButton.icon(
              onPressed: onRefresh,
              style: OutlinedButton.styleFrom(
                minimumSize: const Size(46, 46),
                padding: EdgeInsets.zero,
                side: const BorderSide(color: AppColors.borderStrong),
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(AppSizes.radiusMd),
                ),
              ),
              icon: const Icon(Icons.refresh_rounded, color: AppColors.textPrimary, size: 20),
              label: const SizedBox.shrink(),
            ),
          ],
        ),
        const SizedBox(height: AppSizes.sm),
        Row(
          children: [
            Expanded(
              child: DownloadButton(
                videoId: videoId,
                lessonId: lessonId,
                courseId: courseId,
                title: title,
              ),
            ),
            if (isSavingProgress) ...[
              const SizedBox(width: AppSizes.sm),
              Container(
                padding: const EdgeInsets.all(12),
                decoration: BoxDecoration(
                  color: AppColors.bgSecondary,
                  borderRadius: BorderRadius.circular(AppSizes.radiusMd),
                  border: Border.all(color: AppColors.border),
                ),
                child: const SizedBox(
                  width: 18,
                  height: 18,
                  child: CircularProgressIndicator(
                    strokeWidth: 2,
                    color: Color(0xFFFF5252),
                  ),
                ),
              ),
            ],
          ],
        ),
      ],
    );
  }
}

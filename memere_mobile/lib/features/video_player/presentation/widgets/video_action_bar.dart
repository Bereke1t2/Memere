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
  });

  final String videoId;
  final String lessonId;
  final String courseId;
  final String title;
  final VoidCallback onRefresh;
  final VoidCallback onMarkComplete;
  final bool isSavingProgress;

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        IconButton(
          tooltip: 'Refresh stream',
          onPressed: onRefresh,
          icon: const Icon(Icons.refresh_rounded),
        ),
        DownloadButton(
          videoId: videoId,
          lessonId: lessonId,
          courseId: courseId,
          title: title,
        ),
        IconButton(
          tooltip: 'Mark complete',
          onPressed: lessonId.isEmpty ? null : onMarkComplete,
          icon: const Icon(Icons.check_circle_outline_rounded),
        ),
        const Spacer(),
        AnimatedSwitcher(
          duration: const Duration(milliseconds: 200),
          child: isSavingProgress
              ? const SizedBox(
                  key: ValueKey('saving'),
                  width: AppSizes.iconSm,
                  height: AppSizes.iconSm,
                  child: CircularProgressIndicator(
                    strokeWidth: 2,
                    color: AppColors.accentPrimary,
                  ),
                )
              : const Icon(
                  Icons.cloud_done_outlined,
                  key: ValueKey('saved'),
                  size: AppSizes.iconSm,
                  color: AppColors.textSecondary,
                ),
        ),
      ],
    );
  }
}

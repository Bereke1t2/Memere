import 'package:flutter/material.dart';

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
                  backgroundColor: isCompleted ? const Color(0xFF10B981) : const Color(0xFF1E293B),
                  foregroundColor: Colors.white,
                  minimumSize: const Size.fromHeight(44),
                  elevation: 0,
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(10),
                    side: BorderSide(
                      color: isCompleted ? const Color(0xFF10B981) : const Color(0xFF334155),
                    ),
                  ),
                ),
                icon: Icon(
                  isCompleted ? Icons.check_circle_outline_rounded : Icons.radio_button_unchecked_rounded,
                  size: 18,
                ),
                label: Text(
                  isCompleted ? 'Completed' : 'Mark as Complete',
                  style: const TextStyle(fontWeight: FontWeight.w600, fontSize: 13),
                ),
              ),
            ),
            const SizedBox(width: 8),
            OutlinedButton(
              onPressed: onRefresh,
              style: OutlinedButton.styleFrom(
                minimumSize: const Size(44, 44),
                padding: EdgeInsets.zero,
                side: const BorderSide(color: Color(0xFF334155)),
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(10),
                ),
              ),
              child: const Icon(Icons.refresh_rounded, color: Colors.white, size: 18),
            ),
          ],
        ),
        const SizedBox(height: 8),
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
              const SizedBox(width: 8),
              Container(
                padding: const EdgeInsets.all(12),
                decoration: BoxDecoration(
                  color: const Color(0xFF1E293B),
                  borderRadius: BorderRadius.circular(10),
                  border: Border.all(color: const Color(0xFF334155)),
                ),
                child: const SizedBox(
                  width: 16,
                  height: 16,
                  child: CircularProgressIndicator(
                    strokeWidth: 2,
                    color: Colors.white,
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

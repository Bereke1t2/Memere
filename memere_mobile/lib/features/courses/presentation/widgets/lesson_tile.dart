import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/router/app_router.dart';
import '../../domain/entities/lesson_entity.dart';

/// Clean, numbered Lesson Tile matching image copy 5.png (Screen 3)
class LessonTile extends StatefulWidget {
  const LessonTile({
    super.key,
    required this.lesson,
    required this.lessonNumber,
    this.canOpen = false,
  });

  final LessonEntity lesson;
  final int lessonNumber;
  final bool canOpen;

  @override
  State<LessonTile> createState() => _LessonTileState();
}

class _LessonTileState extends State<LessonTile> {
  LessonEntity get lesson => widget.lesson;
  int get lessonNumber => widget.lessonNumber;
  bool get canOpen => widget.canOpen;

  @override
  Widget build(BuildContext context) {
    final isOpenable = canOpen || lesson.isFreePreview;
    final numStr = lessonNumber.toString().padLeft(2, '0');

    return InkWell(
      onTap: () => _handleTap(context),
      borderRadius: BorderRadius.circular(16),
      child: Container(
        padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 12),
        decoration: BoxDecoration(
          color: AppColors.bgSecondary,
          borderRadius: BorderRadius.circular(16),
          border: Border.all(color: AppColors.borderStrong.withAlpha(90)),
        ),
        child: Row(
          children: [
            // Left Number: e.g. "01."
            Text(
              '$numStr.',
              style: const TextStyle(
                fontSize: 15,
                fontWeight: FontWeight.w800,
                color: Color(0xFF94A3B8),
              ),
            ),
            const SizedBox(width: 12),

            // Middle Column: Title & Duration / Type Subtitle
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    lesson.title,
                    maxLines: 2,
                    overflow: TextOverflow.ellipsis,
                    style: const TextStyle(
                      fontSize: 14.5,
                      fontWeight: FontWeight.w700,
                      color: Colors.white,
                      height: 1.25,
                    ),
                  ),
                  const SizedBox(height: 3),
                  Text(
                    _buildSubtitleText(lesson),
                    style: const TextStyle(
                      fontSize: 12,
                      color: Color(0xFF64748B),
                      fontWeight: FontWeight.w500,
                    ),
                  ),
                ],
              ),
            ),
            const SizedBox(width: 10),

            // Right Action Icon Pill (Play / Document / Quiz / Lock)
            Container(
              width: 36,
              height: 36,
              alignment: Alignment.center,
              decoration: BoxDecoration(
                color: isOpenable
                    ? (lesson.isFreePreview && !canOpen
                        ? const Color(0x2238BDF8)
                        : const Color(0x2210B981))
                    : const Color(0xFF1E2433),
                borderRadius: BorderRadius.circular(10),
                border: Border.all(
                  color: isOpenable
                      ? (lesson.isFreePreview && !canOpen
                          ? const Color(0x5538BDF8)
                          : const Color(0x5510B981))
                      : const Color(0xFF2A3449),
                ),
              ),
              child: Icon(
                isOpenable ? _actionIcon(lesson) : Icons.lock_outline_rounded,
                size: 18,
                color: isOpenable
                    ? (lesson.isFreePreview && !canOpen
                        ? const Color(0xFF38BDF8)
                        : const Color(0xFF10B981))
                    : const Color(0xFF64748B),
              ),
            ),
          ],
        ),
      ),
    );
  }

  IconData _actionIcon(LessonEntity lesson) {
    if (lesson.hasVideo || lesson.type == LessonType.video) {
      return Icons.play_arrow_rounded;
    }
    if (lesson.hasQuiz || lesson.type == LessonType.quiz) {
      return Icons.quiz_outlined;
    }
    if (lesson.hasPdf) {
      return Icons.picture_as_pdf_rounded;
    }
    if (lesson.hasContent || lesson.type == LessonType.note) {
      return Icons.article_rounded;
    }
    return Icons.description_rounded;
  }

  void _handleTap(BuildContext context) {
    final isOpenable = canOpen || lesson.isFreePreview;

    if (!isOpenable) {
      _showMessage(context, 'Enroll in this course to unlock this lesson.');
      return;
    }

    // 1. Playable Video content
    if (lesson.hasVideo || lesson.type == LessonType.video || (lesson.videoId != null && lesson.videoId!.isNotEmpty)) {
      final effectiveVideoId = (lesson.videoId != null && lesson.videoId!.isNotEmpty)
          ? lesson.videoId!
          : lesson.id;
      context.push(
        AppRoutes.videoPlayerPath(
          videoId: effectiveVideoId,
          lessonId: lesson.id,
          courseId: lesson.courseId,
          title: lesson.title,
        ),
      );
      return;
    }

    // 2. Quiz content
    if (lesson.hasQuiz || lesson.type == LessonType.quiz) {
      if (lesson.quizId != null && lesson.quizId!.isNotEmpty) {
        context.push(AppRoutes.quizDetailPath(lesson.quizId!));
        return;
      }
      _showMessage(context, 'This quiz is coming soon.');
      return;
    }

    // 3. Open Study Document / Notes Reader for all lessons
    final pdfName = lesson.pdfUrl ?? '';
    context.push(
      AppRoutes.pdfReaderPath(
        title: lesson.title,
        pdfUrl: pdfName,
        lessonId: lesson.id,
        content: lesson.content,
      ),
      extra: <String, dynamic>{
        'title': lesson.title,
        'pdfUrl': pdfName,
        'lessonId': lesson.id,
        'content': lesson.content,
      },
    );
  }

  String _buildSubtitleText(LessonEntity lesson) {
    if (lesson.hasPdf && lesson.hasContent) return '${lesson.durationLabel} • Notes & PDF';
    if (lesson.hasPdf) return '${lesson.durationLabel} • PDF Document';
    if (lesson.hasContent || lesson.type == LessonType.note) return '${lesson.durationLabel} • Study Notes';
    if (lesson.hasQuiz || lesson.type == LessonType.quiz) return 'Practice Quiz';
    return '${lesson.durationLabel} • Video Lesson';
  }

  void _showMessage(BuildContext context, String message) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Text(message),
        behavior: SnackBarBehavior.floating,
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(AppSizes.radiusMd)),
      ),
    );
  }
}

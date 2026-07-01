import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../../../core/router/app_router.dart';
import '../../domain/entities/lesson_entity.dart';

class LessonTile extends StatelessWidget {
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
  Widget build(BuildContext context) {
    return InkWell(
      onTap: () => _handleTap(context),
      borderRadius: BorderRadius.circular(AppSizes.radiusSm),
      child: Padding(
        padding: const EdgeInsets.symmetric(
          horizontal: AppSizes.sm,
          vertical: AppSizes.sm,
        ),
        child: Row(
          children: [
            Container(
              width: 34,
              height: 34,
              alignment: Alignment.center,
              decoration: BoxDecoration(
                color: lesson.isFreePreview
                    ? AppColors.successSurface
                    : AppColors.bgTertiary,
                borderRadius: BorderRadius.circular(AppSizes.radiusSm),
              ),
              child: Text(
                lessonNumber.toString(),
                style: AppTextStyles.labelMedium.copyWith(
                  color: lesson.isFreePreview
                      ? AppColors.success
                      : AppColors.textSecondary,
                ),
              ),
            ),
            const SizedBox(width: AppSizes.sm),
            Icon(
              _lessonIcon(lesson.type),
              size: AppSizes.iconSm,
              color: AppColors.textSecondary,
            ),
            const SizedBox(width: AppSizes.sm),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    lesson.title,
                    maxLines: 2,
                    overflow: TextOverflow.ellipsis,
                    style: AppTextStyles.bodyMedium,
                  ),
                  const SizedBox(height: 2),
                  Text(
                    '${lesson.typeLabel} • ${lesson.durationLabel}',
                    style: AppTextStyles.caption,
                  ),
                ],
              ),
            ),
            const SizedBox(width: AppSizes.sm),
            if (canOpen && lesson.hasVideo)
              const Icon(
                Icons.play_arrow_rounded,
                size: AppSizes.iconSm,
                color: AppColors.accentPrimary,
              )
            else if (lesson.isFreePreview)
              _PreviewBadge()
            else
              const Icon(
                Icons.lock_outline_rounded,
                size: AppSizes.iconSm,
                color: AppColors.textDisabled,
              ),
          ],
        ),
      ),
    );
  }

  void _handleTap(BuildContext context) {
    if (!canOpen) {
      _showMessage(context, 'Enrollment starts in Phase 6.');
      return;
    }

    if (lesson.type == LessonType.quiz) {
      if (!lesson.hasQuiz) {
        _showMessage(context, 'This lesson does not have a quiz attached yet.');
        return;
      }
      context.push(AppRoutes.quizDetailPath(lesson.quizId!));
      return;
    }

    if (lesson.type != LessonType.video) {
      _showMessage(context, 'This lesson type opens in a later phase.');
      return;
    }

    if (!lesson.hasVideo) {
      _showMessage(context, 'This lesson does not have a video attached yet.');
      return;
    }

    context.push(
      AppRoutes.videoPlayerPath(
        videoId: lesson.videoId!,
        lessonId: lesson.id,
        courseId: lesson.courseId,
        title: lesson.title,
      ),
    );
  }

  void _showMessage(BuildContext context, String message) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text(message)),
    );
  }
}

class _PreviewBadge extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(
        horizontal: AppSizes.sm,
        vertical: AppSizes.xs,
      ),
      decoration: BoxDecoration(
        color: AppColors.successSurface,
        borderRadius: BorderRadius.circular(AppSizes.radiusFull),
      ),
      child: Text(
        'Preview',
        style: AppTextStyles.labelSmall.copyWith(color: AppColors.success),
      ),
    );
  }
}

IconData _lessonIcon(LessonType type) {
  switch (type) {
    case LessonType.note:
      return Icons.article_outlined;
    case LessonType.quiz:
      return Icons.quiz_outlined;
    case LessonType.mixed:
      return Icons.widgets_outlined;
    case LessonType.video:
      return Icons.play_circle_outline_rounded;
  }
}

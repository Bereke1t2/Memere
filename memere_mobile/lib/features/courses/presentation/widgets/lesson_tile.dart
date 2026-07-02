import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_motion.dart';
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
    final completed = canOpen && lesson.isFreePreview;
    final playable = canOpen && (lesson.hasVideo || lesson.hasQuiz);

    return InkWell(
      onTap: () => _handleTap(context),
      borderRadius: BorderRadius.circular(AppSizes.radiusLg),
      child: AnimatedContainer(
        duration: AppMotion.base,
        curve: AppMotion.standard,
        padding: const EdgeInsets.symmetric(
          horizontal: AppSizes.sm,
          vertical: AppSizes.sm,
        ),
        decoration: BoxDecoration(
          color: completed ? AppColors.successSurface : AppColors.bgSecondary,
          borderRadius: BorderRadius.circular(AppSizes.radiusLg),
          border: Border.all(
            color:
                completed ? AppColors.success.withAlpha(70) : AppColors.border,
          ),
        ),
        child: Row(
          children: [
            _LessonNumber(number: lessonNumber, completed: completed),
            const SizedBox(width: AppSizes.sm),
            Container(
              width: 36,
              height: 36,
              alignment: Alignment.center,
              decoration: BoxDecoration(
                color: AppColors.bgTertiary,
                borderRadius: BorderRadius.circular(AppSizes.radiusMd),
              ),
              child: Icon(
                _lessonIcon(lesson.type),
                size: AppSizes.iconSm,
                color: AppColors.accentPrimary,
              ),
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
                    style: AppTextStyles.bodyMedium.copyWith(
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                  const SizedBox(height: 2),
                  Text(
                    '${lesson.typeLabel} • ${lesson.durationLabel}',
                    style: AppTextStyles.caption.copyWith(
                      color: AppColors.textMuted,
                    ),
                  ),
                ],
              ),
            ),
            const SizedBox(width: AppSizes.sm),
            AnimatedSwitcher(
              duration: AppMotion.base,
              switchInCurve: AppMotion.emphasized,
              child: completed
                  ? const _CompletedMark()
                  : playable
                      ? const _PlayMark()
                      : lesson.isFreePreview
                          ? _PreviewBadge()
                          : const Icon(
                              Icons.lock_outline_rounded,
                              size: AppSizes.iconSm,
                              color: AppColors.textDisabled,
                            ),
            ),
          ],
        ),
      ),
    );
  }

  void _handleTap(BuildContext context) {
    if (!canOpen) {
      _showMessage(context, 'Enroll to unlock this lesson.');
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

class _LessonNumber extends StatelessWidget {
  const _LessonNumber({
    required this.number,
    required this.completed,
  });

  final int number;
  final bool completed;

  @override
  Widget build(BuildContext context) {
    return AnimatedContainer(
      duration: AppMotion.base,
      curve: AppMotion.standard,
      width: 34,
      height: 34,
      alignment: Alignment.center,
      decoration: BoxDecoration(
        color: completed ? AppColors.success : AppColors.bgTertiary,
        shape: BoxShape.circle,
      ),
      child: AnimatedSwitcher(
        duration: AppMotion.base,
        child: completed
            ? const Icon(
                Icons.check_rounded,
                key: ValueKey('done'),
                size: AppSizes.iconSm,
                color: AppColors.textInverse,
              )
            : Text(
                number.toString(),
                key: ValueKey(number),
                style: AppTextStyles.labelMedium.copyWith(
                  color: AppColors.textSecondary,
                ),
              ),
      ),
    );
  }
}

class _CompletedMark extends StatelessWidget {
  const _CompletedMark();

  @override
  Widget build(BuildContext context) {
    return TweenAnimationBuilder<double>(
      key: const ValueKey('completed'),
      tween: Tween(begin: 0.82, end: 1),
      duration: AppMotion.base,
      curve: AppMotion.emphasized,
      builder: (context, value, child) {
        return Transform.scale(scale: value, child: child);
      },
      child: const Icon(
        Icons.check_circle_rounded,
        color: AppColors.success,
        size: AppSizes.iconMd,
      ),
    );
  }
}

class _PlayMark extends StatelessWidget {
  const _PlayMark();

  @override
  Widget build(BuildContext context) {
    return Container(
      key: const ValueKey('play'),
      width: 30,
      height: 30,
      alignment: Alignment.center,
      decoration: const BoxDecoration(
        color: AppColors.bgTertiary,
        shape: BoxShape.circle,
      ),
      child: const Icon(
        Icons.play_arrow_rounded,
        color: AppColors.accentPrimary,
        size: AppSizes.iconSm,
      ),
    );
  }
}

class _PreviewBadge extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Container(
      key: const ValueKey('preview'),
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

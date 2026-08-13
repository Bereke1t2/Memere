import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_motion.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../../../core/router/app_router.dart';
import '../../domain/entities/lesson_entity.dart';

/// Professional, Coursera/MasterClass-grade Lesson Tile with Checkbox indicator
/// and direct launcher for HLS Video Player and In-App PDF Reader Screen.
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
  bool _isCompleted = false;

  @override
  void initState() {
    super.initState();
    _isCompleted = widget.canOpen && widget.lesson.isFreePreview;
  }

  LessonEntity get lesson => widget.lesson;
  int get lessonNumber => widget.lessonNumber;
  bool get canOpen => widget.canOpen;

  @override
  Widget build(BuildContext context) {
    final isOpenable = canOpen || lesson.isFreePreview;
    final hasActiveMedia = lesson.hasVideo ||
        lesson.hasQuiz ||
        lesson.hasContent ||
        lesson.hasPdf ||
        lesson.type == LessonType.note ||
        lesson.type == LessonType.mixed;

    return InkWell(
      onTap: () => _handleTap(context),
      borderRadius: BorderRadius.circular(AppSizes.radiusLg),
      child: AnimatedContainer(
        duration: AppMotion.base,
        curve: AppMotion.standard,
        padding: const EdgeInsets.symmetric(
          horizontal: AppSizes.md,
          vertical: AppSizes.sm + 2,
        ),
        decoration: BoxDecoration(
          color: _isCompleted
              ? AppColors.bgSecondary.withAlpha(200)
              : AppColors.bgSecondary,
          borderRadius: BorderRadius.circular(AppSizes.radiusLg),
          border: Border.all(
            color: _isCompleted
                ? AppColors.border
                : AppColors.borderStrong.withAlpha(120),
          ),
        ),
        child: Row(
          children: [
            // Professional Checkbox / Progress Indicator
            GestureDetector(
              onTap: isOpenable
                  ? () {
                      setState(() => _isCompleted = !_isCompleted);
                    }
                  : null,
              child: AnimatedContainer(
                duration: AppMotion.base,
                width: 28,
                height: 28,
                alignment: Alignment.center,
                decoration: BoxDecoration(
                  color: _isCompleted
                      ? AppColors.textPrimary
                      : isOpenable
                          ? AppColors.bgTertiary
                          : AppColors.bgTertiary.withAlpha(80),
                  borderRadius: BorderRadius.circular(AppSizes.radiusSm + 2),
                  border: Border.all(
                    color: _isCompleted
                        ? AppColors.textPrimary
                        : isOpenable
                            ? AppColors.borderStrong
                            : AppColors.border,
                  ),
                ),
                child: _isCompleted
                    ? const Icon(
                        Icons.check_rounded,
                        size: 18,
                        color: AppColors.bgPrimary,
                      )
                    : isOpenable
                        ? Text(
                            '$lessonNumber',
                            style: AppTextStyles.labelMedium.copyWith(
                              color: AppColors.textSecondary,
                              fontWeight: FontWeight.bold,
                            ),
                          )
                        : const Icon(
                            Icons.lock_outline_rounded,
                            size: 14,
                            color: AppColors.textDisabled,
                          ),
              ),
            ),
            const SizedBox(width: AppSizes.md),

            // Format type icon box
            Container(
              width: 38,
              height: 38,
              alignment: Alignment.center,
              decoration: BoxDecoration(
                color: _typeBgColor(lesson),
                borderRadius: BorderRadius.circular(AppSizes.radiusMd),
                border: Border.all(color: _typeBorderColor(lesson)),
              ),
              child: Icon(
                _lessonIcon(lesson),
                size: 20,
                color: _typeIconColor(lesson),
              ),
            ),
            const SizedBox(width: AppSizes.md),

            // Lesson Title & Subtitle Metadata
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
                      color: AppColors.textPrimary,
                      letterSpacing: -0.2,
                    ),
                  ),
                  const SizedBox(height: 4),
                  Wrap(
                    crossAxisAlignment: WrapCrossAlignment.center,
                    spacing: 6,
                    runSpacing: 4,
                    children: [
                      Text(
                        _buildSubtitleText(lesson),
                        style: AppTextStyles.caption.copyWith(
                          color: AppColors.textMuted,
                          fontSize: 12,
                        ),
                      ),
                      if (lesson.hasPdf)
                        Container(
                          padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 1.5),
                          decoration: BoxDecoration(
                            color: const Color(0x22FF5252),
                            borderRadius: BorderRadius.circular(4),
                            border: Border.all(color: const Color(0x55FF5252)),
                          ),
                          child: const Text(
                            'PDF DOC',
                            style: TextStyle(
                              fontSize: 9,
                              fontWeight: FontWeight.w800,
                              color: Color(0xFFFF5252),
                              letterSpacing: 0.5,
                            ),
                          ),
                        ),
                      if (lesson.hasContent && !lesson.hasPdf)
                        Container(
                          padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 1.5),
                          decoration: BoxDecoration(
                            color: const Color(0x22448AFF),
                            borderRadius: BorderRadius.circular(4),
                            border: Border.all(color: const Color(0x55448AFF)),
                          ),
                          child: const Text(
                            'STUDY NOTE',
                            style: TextStyle(
                              fontSize: 9,
                              fontWeight: FontWeight.w800,
                              color: Color(0xFF448AFF),
                              letterSpacing: 0.5,
                            ),
                          ),
                        ),
                    ],
                  ),
                ],
              ),
            ),
            const SizedBox(width: AppSizes.sm),

            // Right Action Status Indicator
            AnimatedSwitcher(
              duration: AppMotion.base,
              switchInCurve: AppMotion.emphasized,
              child: _isCompleted
                  ? const _CompletedMark()
                  : isOpenable && hasActiveMedia
                      ? _PlayMark(isNote: lesson.hasPdf || lesson.type == LessonType.note)
                      : lesson.isFreePreview
                          ? _PreviewBadge()
                          : const Icon(
                              Icons.lock_outline_rounded,
                              size: 16,
                              color: AppColors.textDisabled,
                            ),
            ),
          ],
        ),
      ),
    );
  }

  void _handleTap(BuildContext context) {
    if (!canOpen && !lesson.isFreePreview) {
      _showMessage(context, 'Enroll in this course to unlock this lesson.');
      return;
    }

    // 1. Video playback
    if (lesson.hasVideo) {
      context.push(
        AppRoutes.videoPlayerPath(
          videoId: lesson.videoId!,
          lessonId: lesson.id,
          courseId: lesson.courseId,
          title: lesson.title,
        ),
      );
      return;
    }

    // 2. Quiz content
    if (lesson.hasQuiz || lesson.type == LessonType.quiz) {
      if (!lesson.hasQuiz) {
        _showMessage(context, 'This quiz is coming soon.');
        return;
      }
      context.push(AppRoutes.quizDetailPath(lesson.quizId!));
      return;
    }

    // 3. Note, PDF, or text content -> Directly open Full-Screen In-App PDF Reader
    if (lesson.type == LessonType.note ||
        lesson.type == LessonType.mixed ||
        lesson.hasContent ||
        lesson.hasPdf) {
      final pdfName = lesson.pdfUrl ?? '';
      context.push(
        AppRoutes.pdfReaderPath(
          title: lesson.title,
          pdfUrl: pdfName,
          content: lesson.content,
        ),
      );
      return;
    }

    // 4. Video lesson without video ID attached yet
    if (lesson.type == LessonType.video) {
      _showMessage(context, 'The video for this lesson is being processed by your instructor.');
      return;
    }

    _showMessage(context, 'Lesson material will be available soon.');
  }

  String _buildSubtitleText(LessonEntity lesson) {
    if (lesson.hasPdf) return 'PDF Document • ${lesson.durationLabel}';
    if (lesson.hasContent) return 'Study Note • ${lesson.durationLabel}';
    return '${lesson.typeLabel} • ${lesson.durationLabel}';
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
        color: AppColors.textPrimary,
        size: 20,
      ),
    );
  }
}

class _PlayMark extends StatelessWidget {
  const _PlayMark({this.isNote = false});

  final bool isNote;

  @override
  Widget build(BuildContext context) {
    return Container(
      key: ValueKey(isNote ? 'note' : 'play'),
      width: 28,
      height: 28,
      alignment: Alignment.center,
      decoration: BoxDecoration(
        color: AppColors.bgTertiary,
        shape: BoxShape.circle,
        border: Border.all(color: AppColors.borderStrong),
      ),
      child: Icon(
        isNote ? Icons.description_rounded : Icons.play_arrow_rounded,
        color: AppColors.textPrimary,
        size: 16,
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
        vertical: 2,
      ),
      decoration: BoxDecoration(
        color: AppColors.bgTertiary,
        borderRadius: BorderRadius.circular(AppSizes.radiusFull),
        border: Border.all(color: AppColors.borderStrong),
      ),
      child: Text(
        'Preview',
        style: AppTextStyles.labelSmall.copyWith(color: AppColors.textSecondary, fontSize: 10),
      ),
    );
  }
}

IconData _lessonIcon(LessonEntity lesson) {
  if (lesson.hasPdf) return Icons.picture_as_pdf_rounded;
  if (lesson.hasContent) return Icons.article_rounded;
  switch (lesson.type) {
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

Color _typeBgColor(LessonEntity lesson) {
  if (lesson.hasPdf) return const Color(0x22FF5252);
  if (lesson.hasContent) return const Color(0x22448AFF);
  return AppColors.bgTertiary;
}

Color _typeBorderColor(LessonEntity lesson) {
  if (lesson.hasPdf) return const Color(0x44FF5252);
  if (lesson.hasContent) return const Color(0x44448AFF);
  return AppColors.border;
}

Color _typeIconColor(LessonEntity lesson) {
  if (lesson.hasPdf) return const Color(0xFFFF5252);
  if (lesson.hasContent) return const Color(0xFF448AFF);
  return AppColors.textPrimary;
}

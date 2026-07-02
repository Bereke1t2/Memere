import 'package:cached_network_image/cached_network_image.dart';
import 'package:flutter/material.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_motion.dart';
import '../../../../core/constants/app_shadows.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../../../shared/widgets/app_surface.dart';
import '../../domain/entities/course_entity.dart';
import 'course_card.dart';

class CourseDetailHeader extends StatelessWidget {
  const CourseDetailHeader({
    super.key,
    required this.course,
  });

  final CourseEntity course;

  @override
  Widget build(BuildContext context) {
    final subjectColor = _subjectColor(course.subject);

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Wrap(
          spacing: AppSizes.sm,
          runSpacing: AppSizes.sm,
          children: [
            AppBadge(label: course.subject, color: subjectColor),
            AppBadge(label: 'Grade ${course.grade}', color: AppColors.info),
            AppBadge(
                label: course.levelLabel, color: AppColors.accentSecondary),
          ],
        ),
        const SizedBox(height: AppSizes.sm),
        Text(
          course.title,
          maxLines: 4,
          overflow: TextOverflow.ellipsis,
          style: AppTextStyles.headlineLarge.copyWith(height: 1.25),
        ),
        const SizedBox(height: AppSizes.sm),
        Row(
          children: [
            const CircleAvatar(
              radius: 16,
              backgroundColor: AppColors.bgTertiary,
              child: Icon(
                Icons.school_rounded,
                color: AppColors.accentPrimary,
                size: AppSizes.iconSm,
              ),
            ),
            const SizedBox(width: AppSizes.sm),
            Expanded(
              child: Text(
                'Memere instructor',
                maxLines: 1,
                overflow: TextOverflow.ellipsis,
                style: AppTextStyles.bodyMedium.copyWith(
                  color: AppColors.textSecondary,
                ),
              ),
            ),
            const Icon(
              Icons.star_rounded,
              color: AppColors.warning,
              size: AppSizes.iconSm,
            ),
            const SizedBox(width: AppSizes.xs),
            Text(
              course.ratingAvg > 0
                  ? course.ratingAvg.toStringAsFixed(1)
                  : 'New',
              style: AppTextStyles.labelMedium.copyWith(
                color: AppColors.textPrimary,
              ),
            ),
          ],
        ),
        const SizedBox(height: AppSizes.md),
        _VideoPreview(course: course, subjectColor: subjectColor),
      ],
    );
  }
}

class _VideoPreview extends StatefulWidget {
  const _VideoPreview({
    required this.course,
    required this.subjectColor,
  });

  final CourseEntity course;
  final Color subjectColor;

  @override
  State<_VideoPreview> createState() => _VideoPreviewState();
}

class _VideoPreviewState extends State<_VideoPreview> {
  bool _controlsVisible = true;

  @override
  Widget build(BuildContext context) {
    return AppSurface(
      onTap: () => setState(() => _controlsVisible = !_controlsVisible),
      padding: EdgeInsets.zero,
      radius: AppSizes.radiusXl,
      shadows: AppShadows.lg,
      child: AspectRatio(
        aspectRatio: 16 / 9,
        child: Stack(
          fit: StackFit.expand,
          children: [
            Hero(
              tag: courseHeroTag(widget.course.id),
              child: _PreviewImage(
                course: widget.course,
                color: widget.subjectColor,
              ),
            ),
            DecoratedBox(
              decoration: BoxDecoration(
                gradient: LinearGradient(
                  begin: Alignment.topCenter,
                  end: Alignment.bottomCenter,
                  colors: [
                    Colors.transparent,
                    AppColors.scrim.withAlpha(140),
                  ],
                ),
              ),
            ),
            AnimatedOpacity(
              opacity: _controlsVisible ? 1 : 0,
              duration: AppMotion.base,
              curve: AppMotion.standard,
              child: const Center(
                child: Row(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    _PreviewControl(icon: Icons.skip_previous_rounded),
                    SizedBox(width: AppSizes.sm),
                    _PreviewControl(
                      icon: Icons.play_arrow_rounded,
                      size: 54,
                      iconSize: 30,
                    ),
                    SizedBox(width: AppSizes.sm),
                    _PreviewControl(icon: Icons.skip_next_rounded),
                  ],
                ),
              ),
            ),
            Positioned(
              left: AppSizes.md,
              right: AppSizes.md,
              bottom: AppSizes.md,
              child: AnimatedOpacity(
                opacity: _controlsVisible ? 1 : 0,
                duration: AppMotion.base,
                child: Row(
                  children: [
                    Text(
                      widget.course.durationLabel,
                      style: AppTextStyles.labelSmall.copyWith(
                        color: AppColors.textPrimary,
                      ),
                    ),
                    const SizedBox(width: AppSizes.sm),
                    Expanded(
                      child: ClipRRect(
                        borderRadius:
                            BorderRadius.circular(AppSizes.radiusFull),
                        child: TweenAnimationBuilder<double>(
                          tween: Tween(begin: 0, end: 0.35),
                          duration: AppMotion.slow,
                          curve: AppMotion.standard,
                          builder: (context, value, _) {
                            return LinearProgressIndicator(
                              value: value,
                              minHeight: 4,
                              color: AppColors.accentPrimary,
                              backgroundColor:
                                  AppColors.textPrimary.withAlpha(52),
                            );
                          },
                        ),
                      ),
                    ),
                  ],
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _PreviewImage extends StatelessWidget {
  const _PreviewImage({
    required this.course,
    required this.color,
  });

  final CourseEntity course;
  final Color color;

  @override
  Widget build(BuildContext context) {
    if (course.thumbnailUrl == null || course.thumbnailUrl!.trim().isEmpty) {
      return _FallbackHeader(subject: course.subject, color: color);
    }

    return CachedNetworkImage(
      imageUrl: course.thumbnailUrl!,
      fit: BoxFit.cover,
      placeholder: (_, __) => ColoredBox(color: color.withAlpha(24)),
      errorWidget: (_, __, ___) => _FallbackHeader(
        subject: course.subject,
        color: color,
      ),
    );
  }
}

class _PreviewControl extends StatelessWidget {
  const _PreviewControl({
    required this.icon,
    this.size = 38,
    this.iconSize = AppSizes.iconSm,
  });

  final IconData icon;
  final double size;
  final double iconSize;

  @override
  Widget build(BuildContext context) {
    return Container(
      width: size,
      height: size,
      alignment: Alignment.center,
      decoration: BoxDecoration(
        color: AppColors.bgTertiary.withAlpha(232),
        shape: BoxShape.circle,
        border: Border.all(color: AppColors.borderStrong),
        boxShadow: AppShadows.md,
      ),
      child: Icon(icon, size: iconSize, color: AppColors.accentPrimary),
    );
  }
}

class _FallbackHeader extends StatelessWidget {
  const _FallbackHeader({
    required this.subject,
    required this.color,
  });

  final String subject;
  final Color color;

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: BoxDecoration(
        gradient: LinearGradient(
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
          colors: [
            color.withAlpha(50),
            AppColors.bgTertiary,
          ],
        ),
      ),
      child: Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(Icons.school_rounded, color: color, size: AppSizes.iconXl),
            const SizedBox(height: AppSizes.sm),
            Text(
              subject,
              style: AppTextStyles.titleLarge.copyWith(color: color),
            ),
          ],
        ),
      ),
    );
  }
}

Color _subjectColor(String subject) {
  switch (subject.toLowerCase()) {
    case 'mathematics':
    case 'math':
      return AppColors.subjectMath;
    case 'physics':
      return AppColors.subjectPhysics;
    case 'chemistry':
      return AppColors.subjectChem;
    case 'biology':
      return AppColors.subjectBio;
    case 'english':
      return AppColors.subjectEng;
    case 'civics':
      return AppColors.subjectCivics;
    case 'history':
      return AppColors.subjectHist;
    case 'geography':
      return AppColors.subjectGeo;
    case 'economics':
      return AppColors.subjectEcon;
    default:
      return AppColors.accentPrimary;
  }
}

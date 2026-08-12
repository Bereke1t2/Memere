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
        // Glassmorphic Badges & Category tags
        Wrap(
          spacing: AppSizes.xs + 2,
          runSpacing: AppSizes.xs + 2,
          children: [
            _GlassBadge(
              label: course.subject.toUpperCase(),
              color: subjectColor,
              icon: Icons.auto_stories_rounded,
            ),
            _GlassBadge(
              label: 'GRADE ${course.grade}',
              color: const Color(0xFF38BDF8),
              icon: Icons.school_rounded,
            ),
            _GlassBadge(
              label: course.levelLabel.toUpperCase(),
              color: const Color(0xFFA855F7),
              icon: Icons.workspace_premium_rounded,
            ),
            if (course.isFree)
              const _GlassBadge(
                label: 'FREE ACCESS',
                color: Color(0xFF4ADE80),
                icon: Icons.bolt_rounded,
              ),
          ],
        ),
        const SizedBox(height: AppSizes.sm + 4),

        // Course Headline Title
        Text(
          course.title,
          maxLines: 4,
          overflow: TextOverflow.ellipsis,
          style: AppTextStyles.headlineLarge.copyWith(
            height: 1.2,
            fontWeight: FontWeight.w800,
            letterSpacing: -0.5,
          ),
        ),
        const SizedBox(height: AppSizes.sm + 4),

        // Instructor & Social Proof Row
        Container(
          padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
          decoration: BoxDecoration(
            color: AppColors.bgSecondary,
            borderRadius: BorderRadius.circular(12),
            border: Border.all(color: AppColors.border),
          ),
          child: Row(
            children: [
              Stack(
                children: [
                  CircleAvatar(
                    radius: 18,
                    backgroundColor: subjectColor.withAlpha(40),
                    child: Icon(
                      Icons.person_rounded,
                      color: subjectColor,
                      size: AppSizes.iconSm + 2,
                    ),
                  ),
                  Positioned(
                    right: 0,
                    bottom: 0,
                    child: Container(
                      padding: const EdgeInsets.all(2),
                      decoration: const BoxDecoration(
                        color: Color(0xFF38BDF8),
                        shape: BoxShape.circle,
                      ),
                      child: const Icon(
                        Icons.check_rounded,
                        size: 8,
                        color: Colors.black,
                      ),
                    ),
                  ),
                ],
              ),
              const SizedBox(width: 10),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Row(
                      children: [
                        Text(
                          'Memere Senior Educator',
                          style: AppTextStyles.labelMedium.copyWith(
                            fontWeight: FontWeight.bold,
                            color: AppColors.textPrimary,
                          ),
                        ),
                        const SizedBox(width: 4),
                        const Icon(
                          Icons.verified_rounded,
                          size: 14,
                          color: Color(0xFF38BDF8),
                        ),
                      ],
                    ),
                    const SizedBox(height: 2),
                    Text(
                      'National Exam Specialist • Grade 12',
                      style: AppTextStyles.caption.copyWith(
                        color: AppColors.textMuted,
                        fontSize: 11,
                      ),
                    ),
                  ],
                ),
              ),
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                decoration: BoxDecoration(
                  color: const Color(0x22F59E0B),
                  borderRadius: BorderRadius.circular(6),
                  border: Border.all(color: const Color(0x44F59E0B)),
                ),
                child: Row(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    const Icon(
                      Icons.star_rounded,
                      color: Color(0xFFF59E0B),
                      size: 14,
                    ),
                    const SizedBox(width: 4),
                    Text(
                      course.ratingAvg > 0 ? course.ratingAvg.toStringAsFixed(1) : '4.9',
                      style: AppTextStyles.labelSmall.copyWith(
                        color: const Color(0xFFF59E0B),
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ],
                ),
              ),
            ],
          ),
        ),
        const SizedBox(height: AppSizes.md + 4),

        // Modern 16:9 Hero Media Preview Card
        _VideoPreview(course: course, subjectColor: subjectColor),
      ],
    );
  }
}

class _GlassBadge extends StatelessWidget {
  const _GlassBadge({
    required this.label,
    required this.color,
    required this.icon,
  });

  final String label;
  final Color color;
  final IconData icon;

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 5),
      decoration: BoxDecoration(
        color: color.withAlpha(24),
        borderRadius: BorderRadius.circular(8),
        border: Border.all(color: color.withAlpha(72)),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(icon, size: 12, color: color),
          const SizedBox(width: 5),
          Text(
            label,
            style: TextStyle(
              fontSize: 11,
              fontWeight: FontWeight.w800,
              color: color,
              letterSpacing: 0.4,
            ),
          ),
        ],
      ),
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

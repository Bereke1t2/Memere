import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_shadows.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../../../core/router/app_router.dart';
import '../../../../shared/widgets/app_surface.dart';
import '../../domain/entities/course_entity.dart';

String courseHeroTag(String courseId) => 'course-thumbnail-$courseId';

class CourseCard extends StatelessWidget {
  const CourseCard({
    super.key,
    required this.course,
  });

  final CourseEntity course;

  @override
  Widget build(BuildContext context) {
    return CourseRowCard(course: course);
  }
}

class CourseSpotlightCard extends StatelessWidget {
  const CourseSpotlightCard({
    super.key,
    required this.course,
  });

  final CourseEntity course;

  @override
  Widget build(BuildContext context) {
    final subjectColor = _subjectColor(course.subject);

    return SizedBox(
      width: 220,
      child: AppSurface(
        onTap: () => context.push(AppRoutes.courseDetailPath(course.id)),
        padding: const EdgeInsets.all(AppSizes.md),
        color: AppColors.bgSecondary,
        radius: AppSizes.radiusXl,
        borderColor: subjectColor.withAlpha(72),
        shadows: AppShadows.md,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                AppIconTile(
                  icon: Icons.school_outlined,
                  size: 38,
                  iconSize: AppSizes.iconSm,
                  color: subjectColor,
                ),
                const Spacer(),
                AppBadge(
                  label: _courseBadge(course),
                  color:
                      course.isFree ? AppColors.success : AppColors.textPrimary,
                ),
              ],
            ),
            const SizedBox(height: AppSizes.md),
            Text(
              course.title,
              maxLines: 2,
              overflow: TextOverflow.ellipsis,
              style: AppTextStyles.titleMedium,
            ),
            const SizedBox(height: AppSizes.xs),
            Text(
              course.shortDescription.isNotEmpty
                  ? course.shortDescription
                  : 'Focused lessons, examples, and exam practice.',
              maxLines: 2,
              overflow: TextOverflow.ellipsis,
              style: AppTextStyles.bodySmall,
            ),
            const Spacer(),
            AppBadge(
              label: _shortSubject(course.subject),
              color: subjectColor,
              icon: Icons.auto_stories_outlined,
            ),
            const SizedBox(height: AppSizes.md),
            Row(
              children: [
                Expanded(
                  child: _MetaItem(
                    icon: Icons.schedule_rounded,
                    label: course.durationLabel,
                  ),
                ),
                _MetaItem(
                  icon: Icons.star_rounded,
                  label: course.ratingAvg > 0
                      ? course.ratingAvg.toStringAsFixed(1)
                      : 'New',
                  color: AppColors.textPrimary,
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}

class CourseRowCard extends StatelessWidget {
  const CourseRowCard({
    super.key,
    required this.course,
  });

  final CourseEntity course;

  @override
  Widget build(BuildContext context) {
    final subjectColor = _subjectColor(course.subject);

    return AppSurface(
      onTap: () => context.push(AppRoutes.courseDetailPath(course.id)),
      padding: const EdgeInsets.all(AppSizes.md),
      color: AppColors.bgSecondary,
      radius: AppSizes.radiusXl,
      borderColor: subjectColor.withAlpha(72),
      shadows: AppShadows.sm,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              AppIconTile(
                icon: Icons.school_outlined,
                size: 40,
                iconSize: AppSizes.iconSm,
                color: subjectColor,
              ),
              const SizedBox(width: AppSizes.md),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Row(
                      children: [
                        Expanded(
                          child: Text(
                            course.title,
                            maxLines: 2,
                            overflow: TextOverflow.ellipsis,
                            style: AppTextStyles.titleMedium,
                          ),
                        ),
                        const SizedBox(width: AppSizes.sm),
                        AppBadge(
                          label: course.levelLabel,
                          color: subjectColor,
                        ),
                      ],
                    ),
                    const SizedBox(height: AppSizes.xs),
                    Text(
                      course.shortDescription.isNotEmpty
                          ? course.shortDescription
                          : 'Focused lessons, examples, and exam practice.',
                      maxLines: 2,
                      overflow: TextOverflow.ellipsis,
                      style: AppTextStyles.bodySmall,
                    ),
                    const SizedBox(height: AppSizes.sm),
                    Wrap(
                      spacing: AppSizes.sm,
                      runSpacing: AppSizes.xs,
                      children: [
                        AppBadge(
                          label: _shortSubject(course.subject),
                          color: subjectColor,
                          icon: Icons.auto_stories_outlined,
                        ),
                        AppBadge(
                          label: course.priceLabel,
                          color: course.isFree
                              ? AppColors.success
                              : AppColors.textPrimary,
                        ),
                      ],
                    ),
                  ],
                ),
              ),
            ],
          ),
          const SizedBox(height: AppSizes.md),
          Row(
            children: [
              _MetaItem(
                icon: Icons.menu_book_outlined,
                label: '${course.totalLessons} lessons',
              ),
              const SizedBox(width: AppSizes.md),
              Expanded(
                child: _MetaItem(
                  icon: Icons.schedule_rounded,
                  label: course.durationLabel,
                ),
              ),
              _MetaItem(
                icon: Icons.star_rounded,
                label: course.ratingAvg > 0
                    ? course.ratingAvg.toStringAsFixed(1)
                    : 'New',
                color: AppColors.textPrimary,
              ),
            ],
          ),
        ],
      ),
    );
  }
}

class _MetaItem extends StatelessWidget {
  const _MetaItem({
    required this.icon,
    required this.label,
    this.color = AppColors.textSecondary,
  });

  final IconData icon;
  final String label;
  final Color color;

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        Icon(icon, size: AppSizes.iconXs, color: color),
        const SizedBox(width: AppSizes.xs),
        Flexible(
          child: Text(
            label,
            maxLines: 1,
            overflow: TextOverflow.ellipsis,
            style: AppTextStyles.caption.copyWith(color: AppColors.textMuted),
          ),
        ),
      ],
    );
  }
}

String _courseBadge(CourseEntity course) {
  if (course.isFree) return 'Free';
  if (course.ratingAvg >= 4.5 || course.enrollmentCount >= 1000) {
    return 'Popular';
  }
  return 'New';
}

String _shortSubject(String subject) {
  final normalized = subject.trim();
  if (normalized.toLowerCase() == 'mathematics') return 'Math';
  return normalized;
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

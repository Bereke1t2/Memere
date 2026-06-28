import 'package:cached_network_image/cached_network_image.dart';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../../../core/router/app_router.dart';
import '../../../../shared/utils/formatters.dart';
import '../../domain/entities/course_entity.dart';

class CourseCard extends StatelessWidget {
  const CourseCard({
    super.key,
    required this.course,
  });

  final CourseEntity course;

  @override
  Widget build(BuildContext context) {
    final subjectColor = _subjectColor(course.subject);

    return InkWell(
      onTap: () => context.go(AppRoutes.courseDetailPath(course.id)),
      borderRadius: BorderRadius.circular(AppSizes.radiusMd),
      child: Container(
        decoration: BoxDecoration(
          color: AppColors.bgSecondary,
          borderRadius: BorderRadius.circular(AppSizes.radiusMd),
          border: Border.all(color: AppColors.border),
        ),
        clipBehavior: Clip.antiAlias,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            SizedBox(
              height: 116,
              width: double.infinity,
              child: course.thumbnailUrl == null
                  ? _CourseFallbackVisual(
                      subject: course.subject,
                      color: subjectColor,
                    )
                  : CachedNetworkImage(
                      imageUrl: course.thumbnailUrl!,
                      fit: BoxFit.cover,
                      placeholder: (_, __) => const ColoredBox(
                        color: AppColors.bgTertiary,
                      ),
                      errorWidget: (_, __, ___) => _CourseFallbackVisual(
                        subject: course.subject,
                        color: subjectColor,
                      ),
                    ),
            ),
            Padding(
              padding: const EdgeInsets.all(AppSizes.md),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    children: [
                      _Badge(
                        label: course.subject,
                        color: subjectColor,
                      ),
                      const Spacer(),
                      _Badge(
                        label: course.priceLabel,
                        color: course.isFree
                            ? AppColors.success
                            : AppColors.accentSecondary,
                      ),
                    ],
                  ),
                  const SizedBox(height: AppSizes.sm),
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
                        : course.description,
                    maxLines: 2,
                    overflow: TextOverflow.ellipsis,
                    style: AppTextStyles.bodySmall,
                  ),
                  const SizedBox(height: AppSizes.md),
                  Wrap(
                    spacing: AppSizes.md,
                    runSpacing: AppSizes.xs,
                    children: [
                      _MetaItem(
                        icon: Icons.menu_book_outlined,
                        label: course.lessonCountLabel,
                      ),
                      _MetaItem(
                        icon: Icons.schedule_rounded,
                        label: course.durationLabel,
                      ),
                      _MetaItem(
                        icon: Icons.star_rounded,
                        label: course.ratingAvg > 0
                            ? course.ratingAvg.toStringAsFixed(1)
                            : 'New',
                      ),
                      _MetaItem(
                        icon: Icons.group_outlined,
                        label: formatCompactCount(course.enrollmentCount),
                      ),
                    ],
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _CourseFallbackVisual extends StatelessWidget {
  const _CourseFallbackVisual({
    required this.subject,
    required this.color,
  });

  final String subject;
  final Color color;

  @override
  Widget build(BuildContext context) {
    return Container(
      color: color.withAlpha(35),
      alignment: Alignment.center,
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(Icons.school_rounded, color: color, size: AppSizes.iconLg),
          const SizedBox(height: AppSizes.sm),
          Text(
            subject,
            style: AppTextStyles.labelLarge.copyWith(color: color),
          ),
        ],
      ),
    );
  }
}

class _Badge extends StatelessWidget {
  const _Badge({
    required this.label,
    required this.color,
  });

  final String label;
  final Color color;

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(
        horizontal: AppSizes.sm,
        vertical: AppSizes.xs,
      ),
      decoration: BoxDecoration(
        color: color.withAlpha(28),
        borderRadius: BorderRadius.circular(AppSizes.radiusFull),
        border: Border.all(color: color.withAlpha(120)),
      ),
      child: Text(
        label,
        maxLines: 1,
        overflow: TextOverflow.ellipsis,
        style: AppTextStyles.labelSmall.copyWith(color: color),
      ),
    );
  }
}

class _MetaItem extends StatelessWidget {
  const _MetaItem({
    required this.icon,
    required this.label,
  });

  final IconData icon;
  final String label;

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        Icon(icon, size: AppSizes.iconXs, color: AppColors.textSecondary),
        const SizedBox(width: AppSizes.xs),
        Text(label, style: AppTextStyles.caption),
      ],
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

import 'package:cached_network_image/cached_network_image.dart';
import 'package:flutter/material.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_shadows.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../../../shared/widgets/app_surface.dart';
import '../../domain/entities/course_entity.dart';

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
        AppSurface(
          padding: EdgeInsets.zero,
          radius: AppSizes.radiusLg,
          shadows: AppShadows.md,
          child: Stack(
            children: [
              SizedBox(
                height: 212,
                width: double.infinity,
                child: course.thumbnailUrl == null
                    ? _FallbackHeader(
                        subject: course.subject,
                        color: subjectColor,
                      )
                    : CachedNetworkImage(
                        imageUrl: course.thumbnailUrl!,
                        fit: BoxFit.cover,
                        placeholder: (_, __) => const ColoredBox(
                          color: AppColors.bgTertiary,
                        ),
                        errorWidget: (_, __, ___) => _FallbackHeader(
                          subject: course.subject,
                          color: subjectColor,
                        ),
                      ),
              ),
              Positioned.fill(
                child: DecoratedBox(
                  decoration: BoxDecoration(
                    gradient: LinearGradient(
                      begin: Alignment.topCenter,
                      end: Alignment.bottomCenter,
                      colors: [
                        Colors.transparent,
                        AppColors.bgPrimary.withAlpha(165),
                      ],
                    ),
                  ),
                ),
              ),
              Positioned(
                left: AppSizes.md,
                right: AppSizes.md,
                bottom: AppSizes.md,
                child: Wrap(
                  spacing: AppSizes.sm,
                  runSpacing: AppSizes.sm,
                  children: [
                    AppBadge(label: course.subject, color: subjectColor),
                    AppBadge(
                      label: 'Grade ${course.grade}',
                      color: AppColors.info,
                    ),
                    AppBadge(
                      label: course.levelLabel,
                      color: AppColors.accentSecondary,
                    ),
                  ],
                ),
              ),
            ],
          ),
        ),
        const SizedBox(height: AppSizes.md),
        Text(course.title, style: AppTextStyles.headlineLarge),
        const SizedBox(height: AppSizes.sm),
        Text(
          course.shortDescription.isNotEmpty
              ? course.shortDescription
              : course.description,
          style: AppTextStyles.bodyMedium.copyWith(
            color: AppColors.textSecondary,
          ),
        ),
      ],
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
      color: color.withAlpha(35),
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

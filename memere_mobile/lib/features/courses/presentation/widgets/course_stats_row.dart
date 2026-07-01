import 'package:flutter/material.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_shadows.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../../../shared/utils/formatters.dart';
import '../../domain/entities/course_entity.dart';

class CourseStatsRow extends StatelessWidget {
  const CourseStatsRow({
    super.key,
    required this.course,
  });

  final CourseEntity course;

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        Expanded(
          child: _StatItem(
            icon: Icons.menu_book_outlined,
            label: 'Lessons',
            value: course.totalLessons.toString(),
          ),
        ),
        Expanded(
          child: _StatItem(
            icon: Icons.schedule_rounded,
            label: 'Duration',
            value: course.durationLabel,
          ),
        ),
        Expanded(
          child: _StatItem(
            icon: Icons.star_rounded,
            label: 'Rating',
            value: course.ratingAvg > 0
                ? course.ratingAvg.toStringAsFixed(1)
                : 'New',
          ),
        ),
        Expanded(
          child: _StatItem(
            icon: Icons.group_outlined,
            label: 'Students',
            value: formatCompactCount(course.enrollmentCount),
          ),
        ),
      ],
    );
  }
}

class _StatItem extends StatelessWidget {
  const _StatItem({
    required this.icon,
    required this.label,
    required this.value,
  });

  final IconData icon;
  final String label;
  final String value;

  @override
  Widget build(BuildContext context) {
    return Container(
      constraints: const BoxConstraints(minHeight: 78),
      padding: const EdgeInsets.all(AppSizes.sm),
      decoration: BoxDecoration(
        color: AppColors.bgSecondary,
        borderRadius: BorderRadius.circular(AppSizes.radiusSm),
        border: Border.all(color: AppColors.hairline),
        boxShadow: AppShadows.sm,
      ),
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(icon, size: AppSizes.iconSm, color: AppColors.accentPrimary),
          const SizedBox(height: AppSizes.xs),
          Text(
            value,
            maxLines: 1,
            overflow: TextOverflow.ellipsis,
            style: AppTextStyles.labelMedium,
          ),
          const SizedBox(height: 2),
          Text(
            label,
            maxLines: 1,
            overflow: TextOverflow.ellipsis,
            style: AppTextStyles.caption,
          ),
        ],
      ),
    );
  }
}

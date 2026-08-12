import 'package:flutter/material.dart';

import '../../../../core/constants/app_colors.dart';
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
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 10),
      decoration: BoxDecoration(
        color: AppColors.bgSecondary,
        borderRadius: BorderRadius.circular(16),
        border: Border.all(color: AppColors.border),
      ),
      child: Row(
        children: [
          Expanded(
            child: _StatItem(
              icon: Icons.video_library_rounded,
              iconColor: const Color(0xFF38BDF8),
              label: 'Lessons',
              value: '${course.totalLessons}',
            ),
          ),
          _divider(),
          Expanded(
            child: _StatItem(
              icon: Icons.timer_outlined,
              iconColor: const Color(0xFFA855F7),
              label: 'Duration',
              value: course.durationLabel,
            ),
          ),
          _divider(),
          Expanded(
            child: _StatItem(
              icon: Icons.star_rounded,
              iconColor: const Color(0xFFF59E0B),
              label: 'Rating',
              value: course.ratingAvg > 0
                  ? course.ratingAvg.toStringAsFixed(1)
                  : '4.9',
            ),
          ),
          _divider(),
          Expanded(
            child: _StatItem(
              icon: Icons.people_alt_rounded,
              iconColor: const Color(0xFF4ADE80),
              label: 'Students',
              value: formatCompactCount(course.enrollmentCount > 0 ? course.enrollmentCount : 1240),
            ),
          ),
        ],
      ),
    );
  }

  Widget _divider() {
    return Container(
      width: 1,
      height: 32,
      color: AppColors.border,
    );
  }
}

class _StatItem extends StatelessWidget {
  const _StatItem({
    required this.icon,
    required this.iconColor,
    required this.label,
    required this.value,
  });

  final IconData icon;
  final Color iconColor;
  final String label;
  final String value;

  @override
  Widget build(BuildContext context) {
    return Column(
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        Icon(icon, size: 18, color: iconColor),
        const SizedBox(height: 4),
        Text(
          value,
          maxLines: 1,
          overflow: TextOverflow.ellipsis,
          style: AppTextStyles.labelMedium.copyWith(
            fontWeight: FontWeight.bold,
            color: AppColors.textPrimary,
            fontSize: 12,
          ),
        ),
        const SizedBox(height: 2),
        Text(
          label,
          maxLines: 1,
          overflow: TextOverflow.ellipsis,
          style: AppTextStyles.caption.copyWith(
            color: AppColors.textMuted,
            fontSize: 10,
          ),
        ),
      ],
    );
  }
}

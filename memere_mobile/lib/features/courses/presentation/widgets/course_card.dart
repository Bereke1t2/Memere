import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/constants/app_colors.dart';
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
    return SizedBox(
      width: 220,
      child: Container(
        padding: const EdgeInsets.all(AppSizes.md),
        decoration: BoxDecoration(
          color: AppColors.bgSecondary,
          borderRadius: BorderRadius.circular(20),
          border: const Border(
            top: BorderSide(color: AppColors.borderStrong),
            left: BorderSide(color: AppColors.borderStrong),
            right: BorderSide(color: AppColors.borderStrong),
            bottom: BorderSide(color: AppColors.borderStrong, width: 3), // Tactile 3D Depth
          ),
        ),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Container(
                  padding: const EdgeInsets.all(8),
                  decoration: BoxDecoration(
                    color: AppColors.brandEmerald.withOpacity(0.15),
                    shape: BoxShape.circle,
                  ),
                  child: const Icon(
                    Icons.school_outlined,
                    size: 20,
                    color: AppColors.brandEmerald,
                  ),
                ),
                const Spacer(),
                AppBadge(
                  label: _courseBadge(course),
                  color: course.isFree ? AppColors.brandEmerald : AppColors.textPrimary,
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
              color: AppColors.brandEmerald,
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
                  color: AppColors.brandAmber,
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
    return InkWell(
      onTap: () => context.push(AppRoutes.courseDetailPath(course.id)),
      borderRadius: BorderRadius.circular(20),
      child: Container(
        padding: const EdgeInsets.all(16),
        decoration: BoxDecoration(
          color: AppColors.bgSecondary,
          borderRadius: BorderRadius.circular(20),
          border: const Border(
            top: BorderSide(color: AppColors.borderStrong),
            left: BorderSide(color: AppColors.borderStrong),
            right: BorderSide(color: AppColors.borderStrong),
            bottom: BorderSide(color: AppColors.borderStrong, width: 3.5), // Tactile 3D Depth
          ),
        ),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Container(
                  padding: const EdgeInsets.all(10),
                  decoration: BoxDecoration(
                    color: AppColors.brandEmerald.withOpacity(0.12),
                    borderRadius: BorderRadius.circular(14),
                    border: Border.all(color: AppColors.brandEmerald.withOpacity(0.3)),
                  ),
                  child: const Icon(
                    Icons.school_rounded,
                    size: 22,
                    color: AppColors.brandEmerald,
                  ),
                ),
                const SizedBox(width: 14),
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
                              style: const TextStyle(
                                fontSize: 16,
                                fontWeight: FontWeight.bold,
                                color: AppColors.textPrimary,
                                letterSpacing: -0.2,
                              ),
                            ),
                          ),
                          const SizedBox(width: 8),
                          Container(
                            padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                            decoration: BoxDecoration(
                              color: AppColors.bgPrimary,
                              borderRadius: BorderRadius.circular(10),
                              border: Border.all(color: AppColors.borderStrong),
                            ),
                            child: Text(
                              course.levelLabel,
                              style: const TextStyle(
                                fontSize: 11,
                                fontWeight: FontWeight.bold,
                                color: AppColors.brandEmerald,
                              ),
                            ),
                          ),
                        ],
                      ),
                      const SizedBox(height: 4),
                      Text(
                        course.shortDescription.isNotEmpty
                            ? course.shortDescription
                            : 'Focused lessons, examples, and exam practice.',
                        maxLines: 2,
                        overflow: TextOverflow.ellipsis,
                        style: const TextStyle(
                          fontSize: 13,
                          color: AppColors.textMuted,
                          height: 1.4,
                        ),
                      ),
                      const SizedBox(height: 10),
                      Row(
                        children: [
                          Container(
                            padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                            decoration: BoxDecoration(
                              color: AppColors.brandEmerald.withOpacity(0.12),
                              borderRadius: BorderRadius.circular(8),
                            ),
                            child: Row(
                              children: [
                                const Icon(Icons.auto_stories_outlined,
                                    size: 12, color: AppColors.brandEmerald),
                                const SizedBox(width: 4),
                                Text(
                                  _shortSubject(course.subject),
                                  style: const TextStyle(
                                    fontSize: 11,
                                    fontWeight: FontWeight.bold,
                                    color: AppColors.brandEmerald,
                                  ),
                                ),
                              ],
                            ),
                          ),
                          const SizedBox(width: 8),
                          Container(
                            padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                            decoration: BoxDecoration(
                              color: course.isFree
                                  ? AppColors.brandEmerald.withOpacity(0.15)
                                  : AppColors.bgPrimary,
                              borderRadius: BorderRadius.circular(8),
                              border: Border.all(
                                color: course.isFree
                                    ? AppColors.brandEmerald
                                    : AppColors.borderStrong,
                              ),
                            ),
                            child: Text(
                              course.priceLabel,
                              style: TextStyle(
                                fontSize: 11,
                                fontWeight: FontWeight.bold,
                                color: course.isFree
                                    ? AppColors.brandEmerald
                                    : AppColors.textPrimary,
                              ),
                            ),
                          ),
                        ],
                      ),
                    ],
                  ),
                ),
              ],
            ),
            const SizedBox(height: 14),
            const Divider(height: 1, color: AppColors.borderStrong),
            const SizedBox(height: 10),
            Row(
              children: [
                _MetaItem(
                  icon: Icons.menu_book_outlined,
                  label: '${course.totalLessons} lessons',
                ),
                const SizedBox(width: 16),
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
                  color: AppColors.brandAmber,
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}

class _MetaItem extends StatelessWidget {
  const _MetaItem({
    required this.icon,
    required this.label,
    this.color = AppColors.textMuted,
  });

  final IconData icon;
  final String label;
  final Color color;

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        Icon(icon, size: 14, color: color),
        const SizedBox(width: 4),
        Flexible(
          child: Text(
            label,
            maxLines: 1,
            overflow: TextOverflow.ellipsis,
            style: const TextStyle(fontSize: 12, color: AppColors.textMuted),
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

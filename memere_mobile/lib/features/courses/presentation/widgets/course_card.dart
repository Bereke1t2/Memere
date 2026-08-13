import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../../../core/router/app_router.dart';
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
            bottom: BorderSide(color: AppColors.borderStrong, width: 3),
          ),
        ),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Container(
                  padding: const EdgeInsets.all(8),
                  decoration: const BoxDecoration(
                    color: AppColors.bgQuaternary,
                    shape: BoxShape.circle,
                  ),
                  child: const Icon(
                    Icons.school_outlined,
                    size: 20,
                    color: AppColors.textSecondary,
                  ),
                ),
                const Spacer(),
                _CourseLevelBadge(levelLabel: course.levelLabel),
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
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
              decoration: BoxDecoration(
                color: AppColors.bgTertiary,
                borderRadius: BorderRadius.circular(8),
                border: Border.all(color: AppColors.borderStrong),
              ),
              child: Row(
                mainAxisSize: MainAxisSize.min,
                children: [
                  const Icon(Icons.auto_stories_outlined,
                      size: 12, color: AppColors.textMuted),
                  const SizedBox(width: 4),
                  Text(
                    _shortSubject(course.subject),
                    style: const TextStyle(
                      fontSize: 11,
                      fontWeight: FontWeight.bold,
                      color: AppColors.textSecondary,
                    ),
                  ),
                ],
              ),
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

/// Interactive Micro-Animated Course Card with Level Badges (Beginner, Intermediate, Advanced)
class CourseRowCard extends StatefulWidget {
  const CourseRowCard({
    super.key,
    required this.course,
  });

  final CourseEntity course;

  @override
  State<CourseRowCard> createState() => _CourseRowCardState();
}

class _CourseRowCardState extends State<CourseRowCard> {
  bool _isPressed = false;

  @override
  Widget build(BuildContext context) {
    final course = widget.course;

    return GestureDetector(
      onTapDown: (_) => setState(() => _isPressed = true),
      onTapUp: (_) => setState(() => _isPressed = false),
      onTapCancel: () => setState(() => _isPressed = false),
      onTap: () => context.push(AppRoutes.courseDetailPath(course.id)),
      child: AnimatedScale(
        scale: _isPressed ? 0.97 : 1.0,
        duration: const Duration(milliseconds: 120),
        curve: Curves.easeOutCubic,
        child: Container(
          padding: const EdgeInsets.all(16),
          decoration: BoxDecoration(
            color: AppColors.bgSecondary,
            borderRadius: BorderRadius.circular(20),
            border: Border(
              top: const BorderSide(color: AppColors.borderStrong),
              left: const BorderSide(color: AppColors.borderStrong),
              right: const BorderSide(color: AppColors.borderStrong),
              bottom: BorderSide(
                color: _isPressed ? AppColors.borderFocused : AppColors.borderStrong,
                width: 3.5,
              ),
            ),
          ),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  // Muted Neutral Icon Tile
                  Container(
                    padding: const EdgeInsets.all(10),
                    decoration: BoxDecoration(
                      color: AppColors.bgTertiary,
                      borderRadius: BorderRadius.circular(14),
                      border: Border.all(color: AppColors.borderStrong),
                    ),
                    child: const Icon(
                      Icons.school_rounded,
                      size: 22,
                      color: AppColors.textPrimary,
                    ),
                  ),
                  const SizedBox(width: 14),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Row(
                          crossAxisAlignment: CrossAxisAlignment.start,
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

                            // Vibrant Course Level Badge (Beginner, Intermediate, Advanced)
                            _CourseLevelBadge(levelLabel: course.levelLabel),
                          ],
                        ),
                        const SizedBox(height: 6),
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
                        Wrap(
                          spacing: 8,
                          runSpacing: 6,
                          children: [
                            // Neutral Subject Chip
                            Container(
                              padding: const EdgeInsets.symmetric(
                                  horizontal: 8, vertical: 4),
                              decoration: BoxDecoration(
                                color: AppColors.bgTertiary,
                                borderRadius: BorderRadius.circular(8),
                                border: Border.all(color: AppColors.borderStrong),
                              ),
                              child: Row(
                                mainAxisSize: MainAxisSize.min,
                                children: [
                                  const Icon(Icons.auto_stories_outlined,
                                      size: 12, color: AppColors.textMuted),
                                  const SizedBox(width: 4),
                                  Text(
                                    _shortSubject(course.subject),
                                    style: const TextStyle(
                                      fontSize: 11,
                                      fontWeight: FontWeight.bold,
                                      color: AppColors.textSecondary,
                                    ),
                                  ),
                                ],
                              ),
                            ),

                            // Price Chip
                            Container(
                              padding: const EdgeInsets.symmetric(
                                  horizontal: 8, vertical: 4),
                              decoration: BoxDecoration(
                                color: AppColors.bgTertiary,
                                borderRadius: BorderRadius.circular(8),
                                border: Border.all(color: AppColors.borderStrong),
                              ),
                              child: Text(
                                course.priceLabel,
                                style: TextStyle(
                                  fontSize: 11,
                                  fontWeight: FontWeight.bold,
                                  color: course.isFree
                                      ? AppColors.brandEmerald
                                      : AppColors.textSecondary,
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

                  // Animated Rating Star Pulse
                  TweenAnimationBuilder<double>(
                    tween: Tween<double>(begin: 0.8, end: 1.0),
                    duration: const Duration(milliseconds: 600),
                    curve: Curves.elasticOut,
                    builder: (context, scale, child) {
                      return Transform.scale(
                        scale: scale,
                        child: child,
                      );
                    },
                    child: _MetaItem(
                      icon: Icons.star_rounded,
                      label: course.ratingAvg > 0
                          ? course.ratingAvg.toStringAsFixed(1)
                          : 'New',
                      color: AppColors.brandAmber,
                    ),
                  ),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }
}

/// Highlighted Course Level Badge for Beginner, Intermediate, and Advanced
class _CourseLevelBadge extends StatelessWidget {
  const _CourseLevelBadge({required this.levelLabel});

  final String levelLabel;

  @override
  Widget build(BuildContext context) {
    final lower = levelLabel.toLowerCase();
    Color bgColor;
    Color textColor;
    String label;

    if (lower.contains('beginner') || lower.contains('basic')) {
      bgColor = AppColors.levelBeginner;
      textColor = Colors.white;
      label = 'BEGINNER';
    } else if (lower.contains('advanced') || lower.contains('expert')) {
      bgColor = AppColors.levelAdvanced;
      textColor = Colors.white;
      label = 'ADVANCED';
    } else {
      bgColor = AppColors.levelIntermediate;
      textColor = Colors.black;
      label = 'INTERMEDIATE';
    }

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 9, vertical: 4),
      decoration: BoxDecoration(
        color: bgColor,
        borderRadius: BorderRadius.circular(10),
        boxShadow: [
          BoxShadow(
            color: bgColor.withAlpha(90),
            blurRadius: 6,
            offset: const Offset(0, 2),
          ),
        ],
      ),
      child: Text(
        label,
        style: TextStyle(
          fontSize: 10,
          fontWeight: FontWeight.w900,
          color: textColor,
          letterSpacing: 0.5,
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

String _shortSubject(String subject) {
  final normalized = subject.trim();
  if (normalized.toLowerCase() == 'mathematics') return 'Math';
  return normalized;
}

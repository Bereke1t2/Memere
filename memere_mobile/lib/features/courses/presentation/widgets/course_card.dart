import 'package:cached_network_image/cached_network_image.dart';
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

/// Unique, polished Course Card with optional thumbnail image support
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
    final hasThumbnail = course.thumbnailUrl != null &&
        course.thumbnailUrl!.trim().isNotEmpty &&
        course.thumbnailUrl!.startsWith('http');

    return GestureDetector(
      onTapDown: (_) => setState(() => _isPressed = true),
      onTapUp: (_) => setState(() => _isPressed = false),
      onTapCancel: () => setState(() => _isPressed = false),
      onTap: () => context.push(AppRoutes.courseDetailPath(course.id)),
      child: AnimatedScale(
        scale: _isPressed ? 0.98 : 1.0,
        duration: const Duration(milliseconds: 120),
        curve: Curves.easeOutCubic,
        child: Container(
          decoration: BoxDecoration(
            color: AppColors.bgSecondary,
            borderRadius: BorderRadius.circular(20),
            border: Border.all(
              color: _isPressed ? AppColors.borderFocused : AppColors.borderStrong,
            ),
          ),
          clipBehavior: Clip.antiAlias,
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // Top Banner Area (Image thumbnail or Distinctive Subject Banner)
              SizedBox(
                height: 120,
                width: double.infinity,
                child: hasThumbnail
                    ? Stack(
                        fit: StackFit.expand,
                        children: [
                          CachedNetworkImage(
                            imageUrl: course.thumbnailUrl!,
                            fit: BoxFit.cover,
                            placeholder: (_, __) => _SubjectFallbackBanner(course: course),
                            errorWidget: (_, __, ___) => _SubjectFallbackBanner(course: course),
                          ),
                          // Subtle bottom shadow gradient for contrast
                          Positioned.fill(
                            child: Container(
                              decoration: BoxDecoration(
                                gradient: LinearGradient(
                                  begin: Alignment.topCenter,
                                  end: Alignment.bottomCenter,
                                  colors: [
                                    Colors.transparent,
                                    Colors.black.withAlpha(140),
                                  ],
                                ),
                              ),
                            ),
                          ),
                          // Overlaid Badges
                          _BannerOverlayBadges(course: course),
                        ],
                      )
                    : _SubjectFallbackBanner(course: course),
              ),

              // Bottom Content Area
              Padding(
                padding: const EdgeInsets.all(16),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    // Subject Track & Grade
                    Row(
                      children: [
                        Text(
                          'GRADE ${course.grade} • ${_streamLabel(course.subject).toUpperCase()}',
                          style: const TextStyle(
                            fontSize: 11,
                            fontWeight: FontWeight.w800,
                            color: Color(0xFF38BDF8),
                            letterSpacing: 0.5,
                          ),
                        ),
                        const Spacer(),
                        _CourseLevelBadge(levelLabel: course.levelLabel),
                      ],
                    ),
                    const SizedBox(height: 6),

                    // Course Title
                    Text(
                      course.title,
                      maxLines: 2,
                      overflow: TextOverflow.ellipsis,
                      style: const TextStyle(
                        fontSize: 16,
                        fontWeight: FontWeight.bold,
                        color: AppColors.textPrimary,
                        letterSpacing: -0.2,
                        height: 1.25,
                      ),
                    ),
                    if (course.shortDescription.isNotEmpty) ...[
                      const SizedBox(height: 4),
                      Text(
                        course.shortDescription,
                        maxLines: 2,
                        overflow: TextOverflow.ellipsis,
                        style: const TextStyle(
                          fontSize: 12.5,
                          color: AppColors.textMuted,
                          height: 1.35,
                        ),
                      ),
                    ],
                    const SizedBox(height: 12),
                    const Divider(height: 1, color: AppColors.borderStrong),
                    const SizedBox(height: 10),

                    // Footer Row: Metadata & Start Learning Action
                    Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        Row(
                          children: [
                            const Icon(
                              Icons.menu_book_outlined,
                              size: 13,
                              color: AppColors.textMuted,
                            ),
                            const SizedBox(width: 4),
                            Text(
                              '${course.totalLessons} lessons',
                              style: const TextStyle(
                                fontSize: 12,
                                color: AppColors.textMuted,
                                fontWeight: FontWeight.w500,
                              ),
                            ),
                            const SizedBox(width: 12),
                            const Icon(
                              Icons.schedule_rounded,
                              size: 13,
                              color: AppColors.textMuted,
                            ),
                            const SizedBox(width: 4),
                            Text(
                              course.durationLabel,
                              style: const TextStyle(
                                fontSize: 12,
                                color: AppColors.textMuted,
                                fontWeight: FontWeight.w500,
                              ),
                            ),
                          ],
                        ),
                        Container(
                          padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 5),
                          decoration: BoxDecoration(
                            color: course.isFree
                                ? const Color(0x1810B981)
                                : const Color(0x1838BDF8),
                            borderRadius: BorderRadius.circular(8),
                            border: Border.all(
                              color: course.isFree
                                  ? const Color(0x3510B981)
                                  : const Color(0x3538BDF8),
                            ),
                          ),
                          child: Row(
                            mainAxisSize: MainAxisSize.min,
                            children: [
                              Text(
                                course.isFree ? 'Start Free' : 'Enroll ${course.priceLabel}',
                                style: TextStyle(
                                  fontSize: 11.5,
                                  fontWeight: FontWeight.bold,
                                  color: course.isFree
                                      ? AppColors.brandEmerald
                                      : const Color(0xFF38BDF8),
                                ),
                              ),
                              const SizedBox(width: 4),
                              Icon(
                                Icons.arrow_forward_rounded,
                                size: 12,
                                color: course.isFree
                                    ? AppColors.brandEmerald
                                    : const Color(0xFF38BDF8),
                              ),
                            ],
                          ),
                        ),
                      ],
                    ),
                  ],
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

/// Fallback decorative banner when no image thumbnail is provided
class _SubjectFallbackBanner extends StatelessWidget {
  const _SubjectFallbackBanner({required this.course});

  final CourseEntity course;

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: BoxDecoration(
        gradient: LinearGradient(
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
          colors: [
            _subjectGradientStart(course.subject),
            _subjectGradientEnd(course.subject),
          ],
        ),
      ),
      child: Stack(
        children: [
          // Ambient background geometry
          Positioned(
            right: -10,
            bottom: -20,
            child: Icon(
              _subjectIcon(course.subject),
              size: 110,
              color: Colors.white.withAlpha(12),
            ),
          ),
          // Overlaid Badges
          _BannerOverlayBadges(course: course),
        ],
      ),
    );
  }
}

/// Overlaid tags & badges on the course card banner
class _BannerOverlayBadges extends StatelessWidget {
  const _BannerOverlayBadges({required this.course});

  final CourseEntity course;

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.all(12),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          // Subject Tag Pill
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 5),
            decoration: BoxDecoration(
              color: Colors.black.withAlpha(150),
              borderRadius: BorderRadius.circular(10),
              border: Border.all(color: Colors.white.withAlpha(35)),
            ),
            child: Row(
              mainAxisSize: MainAxisSize.min,
              children: [
                Icon(_subjectIcon(course.subject), size: 13, color: Colors.white),
                const SizedBox(width: 5),
                Text(
                  _shortSubject(course.subject),
                  style: const TextStyle(
                    fontSize: 11.5,
                    fontWeight: FontWeight.bold,
                    color: Colors.white,
                  ),
                ),
              ],
            ),
          ),

          // Price Tag Pill
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 9, vertical: 5),
            decoration: BoxDecoration(
              color: course.isFree
                  ? const Color(0xE0065F46)
                  : Colors.black.withAlpha(150),
              borderRadius: BorderRadius.circular(10),
              border: Border.all(
                color: course.isFree
                    ? const Color(0xFF10B981)
                    : Colors.white.withAlpha(35),
              ),
            ),
            child: Text(
              course.priceLabel,
              style: TextStyle(
                fontSize: 11,
                fontWeight: FontWeight.bold,
                color: course.isFree ? Colors.white : const Color(0xFFE2E8F0),
              ),
            ),
          ),
        ],
      ),
    );
  }
}

String _streamLabel(String subject) {
  final lower = subject.toLowerCase();
  if (lower.contains('history') || lower.contains('geography') || lower.contains('economics')) {
    return 'Social Science';
  }
  return 'Natural Science';
}

Color _subjectGradientStart(String subject) {
  final lower = subject.toLowerCase();
  if (lower.contains('bio')) return const Color(0xFF064E3B);
  if (lower.contains('phys')) return const Color(0xFF1E293B);
  if (lower.contains('chem')) return const Color(0xFF311B92);
  if (lower.contains('math')) return const Color(0xFF1E3A8A);
  return const Color(0xFF182234);
}

Color _subjectGradientEnd(String subject) {
  final lower = subject.toLowerCase();
  if (lower.contains('bio')) return const Color(0xFF022C22);
  if (lower.contains('phys')) return const Color(0xFF0F172A);
  if (lower.contains('chem')) return const Color(0xFF1A103C);
  if (lower.contains('math')) return const Color(0xFF172554);
  return const Color(0xFF0B0F17);
}

IconData _subjectIcon(String subject) {
  switch (subject.toLowerCase()) {
    case 'mathematics':
    case 'math':
      return Icons.calculate_outlined;
    case 'physics':
      return Icons.science_outlined;
    case 'chemistry':
      return Icons.biotech_outlined;
    case 'biology':
      return Icons.eco_outlined;
    case 'english':
      return Icons.menu_book_outlined;
    case 'history':
      return Icons.history_edu_outlined;
    case 'geography':
      return Icons.public_outlined;
    case 'economics':
      return Icons.trending_up_rounded;
    default:
      return Icons.school_outlined;
  }
}

String _shortSubject(String subject) {
  final normalized = subject.trim();
  if (normalized.toLowerCase() == 'mathematics') return 'Math';
  return normalized;
}

class _CourseLevelBadge extends StatelessWidget {
  const _CourseLevelBadge({required this.levelLabel});

  final String levelLabel;

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 7, vertical: 3),
      decoration: BoxDecoration(
        color: const Color(0xFF1E2433),
        borderRadius: BorderRadius.circular(6),
        border: Border.all(color: const Color(0xFF2A3449), width: 0.8),
      ),
      child: Text(
        levelLabel.toUpperCase(),
        style: const TextStyle(
          fontSize: 9.5,
          fontWeight: FontWeight.w700,
          color: Color(0xFF94A3B8),
          letterSpacing: 0.4,
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
        Icon(icon, size: 13, color: color),
        const SizedBox(width: 4),
        Flexible(
          child: Text(
            label,
            maxLines: 1,
            overflow: TextOverflow.ellipsis,
            style: const TextStyle(fontSize: 11.5, color: AppColors.textMuted),
          ),
        ),
      ],
    );
  }
}

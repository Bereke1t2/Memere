import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../../../core/errors/failures.dart';
import '../../../../shared/widgets/app_button.dart';
import '../providers/course_detail_provider.dart';
import '../widgets/course_detail_header.dart';
import '../widgets/course_detail_skeleton.dart';
import '../widgets/course_empty_state.dart';
import '../widgets/course_section_tile.dart';
import '../widgets/course_stats_row.dart';

class CourseDetailScreen extends ConsumerWidget {
  const CourseDetailScreen({
    super.key,
    required this.courseId,
  });

  final String courseId;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final detailAsync = ref.watch(courseDetailProvider(courseId));

    return Scaffold(
      backgroundColor: AppColors.bgPrimary,
      appBar: AppBar(
        title: const Text('Course detail'),
      ),
      bottomNavigationBar: detailAsync.maybeWhen(
        data: (detail) => SafeArea(
          top: false,
          child: Container(
            padding: const EdgeInsets.all(AppSizes.md),
            decoration: const BoxDecoration(
              color: AppColors.bgSecondary,
              border: Border(top: BorderSide(color: AppColors.border)),
            ),
            child: AppButton(
              label: detail.course.isFree
                  ? 'Start learning'
                  : 'Enroll for ${detail.course.priceLabel}',
              onPressed: () => _showEnrollmentMessage(context),
            ),
          ),
        ),
        orElse: () => null,
      ),
      body: SafeArea(
        top: false,
        child: detailAsync.when(
          loading: () => const CourseDetailSkeleton(),
          error: (error, _) {
            final failure = error is Failure ? error : null;
            final isNotFound =
                failure is ServerFailure && failure.statusCode == 404;
            return CourseEmptyState(
              icon: isNotFound
                  ? Icons.hide_source_rounded
                  : Icons.error_outline_rounded,
              title: isNotFound ? 'Course not found' : 'Could not load course',
              body: isNotFound
                  ? 'This course may have been removed or unpublished.'
                  : failure?.message ?? 'Please try again.',
              buttonLabel: 'Retry',
              onPressed: () => ref.invalidate(courseDetailProvider(courseId)),
            );
          },
          data: (detail) {
            final course = detail.course;
            return RefreshIndicator(
              color: AppColors.accentPrimary,
              backgroundColor: AppColors.bgSecondary,
              onRefresh: () async {
                ref.invalidate(courseDetailProvider(courseId));
                await ref.read(courseDetailProvider(courseId).future);
              },
              child: ListView(
                physics: const AlwaysScrollableScrollPhysics(),
                padding: const EdgeInsets.fromLTRB(
                  AppSizes.screenPaddingH,
                  AppSizes.md,
                  AppSizes.screenPaddingH,
                  AppSizes.xl,
                ),
                children: [
                  CourseDetailHeader(course: course),
                  const SizedBox(height: AppSizes.lg),
                  CourseStatsRow(course: course),
                  const SizedBox(height: AppSizes.lg),
                  _PriceBand(
                    label: course.priceLabel,
                    isFree: course.isFree,
                  ),
                  const SizedBox(height: AppSizes.lg),
                  const Text(
                    'About course',
                    style: AppTextStyles.headlineSmall,
                  ),
                  const SizedBox(height: AppSizes.sm),
                  Text(
                    course.description.isNotEmpty
                        ? course.description
                        : course.shortDescription,
                    style: AppTextStyles.bodyMedium.copyWith(
                      color: AppColors.textSecondary,
                    ),
                  ),
                  const SizedBox(height: AppSizes.lg),
                  const Text(
                    'Curriculum',
                    style: AppTextStyles.headlineSmall,
                  ),
                  const SizedBox(height: AppSizes.md),
                  if (detail.sections.isEmpty)
                    const CourseEmptyState(
                      icon: Icons.menu_book_outlined,
                      title: 'No curriculum yet',
                      body:
                          'Lessons will appear here when this course is ready.',
                    )
                  else
                    ...detail.sections.asMap().entries.map(
                          (entry) => Padding(
                            padding: const EdgeInsets.only(bottom: AppSizes.md),
                            child: CourseSectionTile(
                              section: entry.value,
                              sectionNumber: entry.key + 1,
                              initiallyExpanded: entry.key == 0,
                            ),
                          ),
                        ),
                ],
              ),
            );
          },
        ),
      ),
    );
  }

  void _showEnrollmentMessage(BuildContext context) {
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(content: Text('Enrollment starts in Phase 6.')),
    );
  }
}

class _PriceBand extends StatelessWidget {
  const _PriceBand({
    required this.label,
    required this.isFree,
  });

  final String label;
  final bool isFree;

  @override
  Widget build(BuildContext context) {
    final color = isFree ? AppColors.success : AppColors.accentSecondary;

    return Container(
      padding: const EdgeInsets.all(AppSizes.md),
      decoration: BoxDecoration(
        color: AppColors.bgSecondary,
        borderRadius: BorderRadius.circular(AppSizes.radiusMd),
        border: Border.all(color: AppColors.border),
      ),
      child: Row(
        children: [
          Container(
            width: 40,
            height: 40,
            decoration: BoxDecoration(
              color: color.withAlpha(28),
              borderRadius: BorderRadius.circular(AppSizes.radiusSm),
            ),
            child: Icon(
              isFree ? Icons.lock_open_rounded : Icons.payments_outlined,
              color: color,
            ),
          ),
          const SizedBox(width: AppSizes.md),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(label, style: AppTextStyles.titleMedium),
                const SizedBox(height: AppSizes.xs),
                Text(
                  isFree
                      ? 'You can preview this course for free.'
                      : 'Enrollment and payment are handled in Phase 6.',
                  style: AppTextStyles.bodySmall,
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

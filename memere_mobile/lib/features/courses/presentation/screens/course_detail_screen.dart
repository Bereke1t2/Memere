import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_shadows.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../../../core/errors/failures.dart';
import '../../../../core/router/app_router.dart';
import '../../../../shared/widgets/app_button.dart';
import '../../../../shared/widgets/app_surface.dart';
import '../../../payment/presentation/providers/checkout_flow_provider.dart';
import '../../../payment/presentation/providers/course_access_provider.dart';
import '../../../payment/presentation/providers/purchase_history_provider.dart';
import '../../../payment/presentation/widgets/payment_provider_sheet.dart';
import '../../domain/entities/course_detail_entity.dart';
import '../../domain/entities/course_entity.dart';
import '../../domain/entities/lesson_entity.dart';
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
    final accessAsync = ref.watch(courseAccessProvider(courseId));
    final hasAccess = accessAsync.valueOrNull?.hasAccess ?? false;

    return Scaffold(
      backgroundColor: AppColors.bgPrimary,
      appBar: AppBar(
        title: const Text('Course detail'),
      ),
      bottomNavigationBar: detailAsync.maybeWhen(
        data: (detail) => _CheckoutCtaBar(
          course: detail.course,
          courseId: courseId,
        ),
        orElse: () => null,
      ),
      body: AppPageBackground(
        child: SafeArea(
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
                title:
                    isNotFound ? 'Course not found' : 'Could not load course',
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
                  ref.invalidate(courseAccessProvider(courseId));
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
                      hasAccess: hasAccess,
                    ),
                    const SizedBox(height: AppSizes.lg),
                    const AppSectionHeader(title: 'About course'),
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
                    AppSectionHeader(
                      title: 'Curriculum',
                      subtitle: detail.sections.isEmpty
                          ? null
                          : '${detail.sections.length} sections',
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
                              padding:
                                  const EdgeInsets.only(bottom: AppSizes.md),
                              child: CourseSectionTile(
                                section: entry.value,
                                sectionNumber: entry.key + 1,
                                initiallyExpanded: entry.key == 0,
                                canOpenLessons: course.isFree || hasAccess,
                              ),
                            ),
                          ),
                  ],
                ),
              );
            },
          ),
        ),
      ),
    );
  }
}

/// Bottom CTA that reflects real access state and drives the free/paid flows.
class _CheckoutCtaBar extends ConsumerStatefulWidget {
  const _CheckoutCtaBar({
    required this.course,
    required this.courseId,
  });

  final CourseEntity course;
  final String courseId;

  @override
  ConsumerState<_CheckoutCtaBar> createState() => _CheckoutCtaBarState();
}

class _CheckoutCtaBarState extends ConsumerState<_CheckoutCtaBar> {
  String get _courseId => widget.courseId;

  CheckoutFlowNotifier get _notifier =>
      ref.read(checkoutFlowProvider(_courseId).notifier);

  @override
  Widget build(BuildContext context) {
    final accessAsync = ref.watch(courseAccessProvider(_courseId));
    final checkoutAsync = ref.watch(checkoutFlowProvider(_courseId));
    final busy = checkoutAsync.valueOrNull?.isWorking ?? false;

    return SafeArea(
      top: false,
      child: AppSurface(
        padding: const EdgeInsets.all(AppSizes.md),
        radius: AppSizes.radiusLg,
        color: AppColors.bgSecondary,
        shadows: AppShadows.lg,
        child: accessAsync.when(
          loading: () => const AppButton(
            label: 'Checking access…',
            onPressed: null,
            isLoading: true,
          ),
          error: (_, __) => AppButton(
            label: 'Retry',
            variant: AppButtonVariant.outline,
            onPressed: () => ref.invalidate(courseAccessProvider(_courseId)),
          ),
          data: (access) {
            if (access.hasAccess) {
              return AppButton(
                label: 'Continue learning',
                onPressed: _continueLearning,
                suffixIcon: Icons.play_arrow_rounded,
              );
            }
            if (widget.course.isFree) {
              return AppButton(
                label: 'Start learning',
                isLoading: busy,
                onPressed: busy ? null : _startFree,
                suffixIcon: Icons.arrow_forward_rounded,
              );
            }
            // Surface an in-flight payment so the user can resume verification.
            final latestPayment =
                ref.watch(latestCoursePaymentProvider(_courseId));
            if (latestPayment != null && latestPayment.isPending) {
              return AppButton(
                label: 'Payment pending',
                variant: AppButtonVariant.secondary,
                onPressed: () => context.push(
                  AppRoutes.paymentResultPath(
                    paymentId: latestPayment.paymentId,
                    courseId: _courseId,
                  ),
                ),
              );
            }
            return AppButton(
              label: 'Enroll for ${widget.course.priceLabel}',
              isLoading: busy,
              onPressed: busy ? null : _startPaid,
              suffixIcon: Icons.arrow_forward_rounded,
            );
          },
        ),
      ),
    );
  }

  /// Opens the first playable lesson (video, else quiz) so the CTA is never a
  /// dead end. Falls back to a hint if the curriculum has nothing playable yet.
  void _continueLearning() {
    final detail = ref.read(courseDetailProvider(_courseId)).valueOrNull;
    final lesson = _firstPlayableLesson(detail);
    if (lesson == null) {
      _showMessage('Lessons will appear here when this course is ready.');
      return;
    }
    if (lesson.hasVideo) {
      context.push(
        AppRoutes.videoPlayerPath(
          videoId: lesson.videoId!,
          lessonId: lesson.id,
          courseId: lesson.courseId,
          title: lesson.title,
        ),
      );
      return;
    }
    if (lesson.hasQuiz) {
      context.push(AppRoutes.quizDetailPath(lesson.quizId!));
    }
  }

  LessonEntity? _firstPlayableLesson(CourseDetailEntity? detail) {
    if (detail == null) return null;
    for (final section in detail.sections) {
      for (final lesson in section.lessons) {
        if (lesson.hasVideo || lesson.hasQuiz) return lesson;
      }
    }
    return null;
  }

  Future<void> _startFree() async {
    final ok = await _notifier.startFreeEnrollment();
    if (!mounted) return;
    if (ok) {
      _showMessage('You are enrolled. Start learning!');
    } else {
      final error =
          ref.read(checkoutFlowProvider(_courseId)).valueOrNull?.error;
      _showMessage(error ?? 'Could not enroll. Please try again.');
    }
  }

  Future<void> _startPaid() async {
    final provider = await PaymentProviderSheet.show(
      context,
      amountLabel: widget.course.priceLabel,
    );
    if (provider == null || !mounted) return;

    final initiation = await _notifier.startPaidCheckout(provider: provider);
    if (!mounted) return;

    if (initiation == null) {
      final state = ref.read(checkoutFlowProvider(_courseId)).valueOrNull;
      final code = state?.errorCode;
      if (code == 'COURSE_IS_FREE') {
        await _startFree();
        return;
      }
      if (code == 'ALREADY_ENROLLED') {
        ref.invalidate(courseAccessProvider(_courseId));
        return;
      }
      _showMessage(
          state?.error ?? 'Could not start checkout. Please try again.');
      return;
    }

    context.push(
      AppRoutes.paymentWebViewPath(
        paymentId: initiation.paymentId,
        courseId: _courseId,
        redirectUrl: initiation.redirectUrl,
      ),
    );
  }

  void _showMessage(String message) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text(message)),
    );
  }
}

class _PriceBand extends StatelessWidget {
  const _PriceBand({
    required this.label,
    required this.isFree,
    required this.hasAccess,
  });

  final String label;
  final bool isFree;
  final bool hasAccess;

  @override
  Widget build(BuildContext context) {
    final color = hasAccess
        ? AppColors.success
        : (isFree ? AppColors.success : AppColors.accentSecondary);

    final IconData icon;
    final String title;
    final String body;
    if (hasAccess) {
      icon = Icons.verified_rounded;
      title = 'You have access';
      body = 'Open any lesson below to continue learning.';
    } else if (isFree) {
      icon = Icons.lock_open_rounded;
      title = label;
      body = 'Enroll for free to unlock every lesson.';
    } else {
      icon = Icons.payments_outlined;
      title = label;
      body = 'One-time payment for full lifetime access.';
    }

    return AppSurface(
      padding: const EdgeInsets.all(AppSizes.md),
      shadows: AppShadows.sm,
      child: Row(
        children: [
          AppIconTile(
            icon: icon,
            color: color,
            size: 44,
            iconSize: AppSizes.iconSm,
          ),
          const SizedBox(width: AppSizes.md),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(title, style: AppTextStyles.titleMedium),
                const SizedBox(height: AppSizes.xs),
                Text(body, style: AppTextStyles.bodySmall),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

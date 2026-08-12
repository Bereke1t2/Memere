import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_motion.dart';
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

class CourseDetailScreen extends ConsumerStatefulWidget {
  const CourseDetailScreen({
    super.key,
    required this.courseId,
  });

  final String courseId;

  @override
  ConsumerState<CourseDetailScreen> createState() => _CourseDetailScreenState();
}

class _CourseDetailScreenState extends ConsumerState<CourseDetailScreen> {
  int _selectedTab = 0;

  String get courseId => widget.courseId;

  @override
  Widget build(BuildContext context) {
    final detailAsync = ref.watch(courseDetailProvider(courseId));
    final accessAsync = ref.watch(courseAccessProvider(courseId));
    final hasAccess = accessAsync.valueOrNull?.hasAccess ?? false;

    return Scaffold(
      backgroundColor: AppColors.bgPrimary,
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
                  padding: EdgeInsets.zero,
                  children: [
                    _DetailTopBand(title: detail.course.title),
                    _DetailContent(
                      detail: detail,
                      hasAccess: hasAccess,
                      selectedTab: _selectedTab,
                      onTabChanged: (index) {
                        setState(() => _selectedTab = index);
                      },
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

class _DetailTopBand extends StatelessWidget {
  const _DetailTopBand({required this.title});

  final String title;

  @override
  Widget build(BuildContext context) {
    final topPadding = MediaQuery.paddingOf(context).top;

    return Container(
      padding: EdgeInsets.fromLTRB(
        AppSizes.screenPaddingH,
        topPadding + AppSizes.sm,
        AppSizes.screenPaddingH,
        AppSizes.lg,
      ),
      decoration: const BoxDecoration(
        color: AppColors.bgSecondary,
        borderRadius: BorderRadius.vertical(
          bottom: Radius.circular(AppSizes.radiusLg),
        ),
        border: Border(
          bottom: BorderSide(color: AppColors.borderStrong),
        ),
      ),
      child: Row(
        children: [
          _TopIconButton(
            icon: Icons.arrow_back_ios_new_rounded,
            onTap: () {
              if (context.canPop()) {
                context.pop();
              } else {
                context.go(AppRoutes.home);
              }
            },
          ),
          const SizedBox(width: AppSizes.md),
          Expanded(
            child: Text(
              title,
              textAlign: TextAlign.center,
              maxLines: 1,
              overflow: TextOverflow.ellipsis,
              style: AppTextStyles.titleLarge,
            ),
          ),
          const SizedBox(width: AppSizes.md),
          _TopIconButton(
            icon: Icons.more_horiz_rounded,
            onTap: () {},
          ),
        ],
      ),
    );
  }
}

class _TopIconButton extends StatelessWidget {
  const _TopIconButton({
    required this.icon,
    required this.onTap,
  });

  final IconData icon;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    return AppPressable(
      onTap: onTap,
      borderRadius: AppSizes.radiusFull,
      child: Container(
        width: 36,
        height: 36,
        alignment: Alignment.center,
        decoration: BoxDecoration(
          color: AppColors.bgTertiary,
          shape: BoxShape.circle,
          border: Border.all(color: AppColors.borderStrong),
        ),
        child:
            Icon(icon, color: AppColors.textSecondary, size: AppSizes.iconSm),
      ),
    );
  }
}

class _DetailContent extends StatelessWidget {
  const _DetailContent({
    required this.detail,
    required this.hasAccess,
    required this.selectedTab,
    required this.onTabChanged,
  });

  final CourseDetailEntity detail;
  final bool hasAccess;
  final int selectedTab;
  final ValueChanged<int> onTabChanged;

  @override
  Widget build(BuildContext context) {
    final course = detail.course;
    final bottomInset = MediaQuery.paddingOf(context).bottom;

    return Container(
      padding: EdgeInsets.fromLTRB(
        AppSizes.screenPaddingH,
        AppSizes.md,
        AppSizes.screenPaddingH,
        AppSizes.xxxl + AppSizes.buttonHeight + bottomInset,
      ),
      decoration: const BoxDecoration(
        color: AppColors.bgPrimary,
        borderRadius: BorderRadius.zero,
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          CourseDetailHeader(course: course),
          const SizedBox(height: AppSizes.md),
          CourseStatsRow(course: course),
          const SizedBox(height: AppSizes.md),
          _PriceBand(
            label: course.priceLabel,
            isFree: course.isFree,
            hasAccess: hasAccess,
          ),
          const SizedBox(height: AppSizes.md),
          _DetailTabs(
            selectedIndex: selectedTab,
            onChanged: onTabChanged,
          ),
          const SizedBox(height: AppSizes.md),
          AnimatedSwitcher(
            duration: AppMotion.slow,
            switchInCurve: AppMotion.standard,
            switchOutCurve: AppMotion.exit,
            child: selectedTab == 0
                ? _CurriculumTab(
                    key: const ValueKey('classes'),
                    detail: detail,
                    canOpenLessons: course.isFree || hasAccess,
                  )
                : _DescriptionTab(
                    key: const ValueKey('description'),
                    course: course,
                  ),
          ),
        ],
      ),
    );
  }
}

class _DetailTabs extends StatelessWidget {
  const _DetailTabs({
    required this.selectedIndex,
    required this.onChanged,
  });

  final int selectedIndex;
  final ValueChanged<int> onChanged;

  @override
  Widget build(BuildContext context) {
    return Container(
      height: 44,
      decoration: BoxDecoration(
        color: AppColors.bgTertiary,
        borderRadius: BorderRadius.circular(AppSizes.radiusFull),
        border: Border.all(color: AppColors.borderStrong),
      ),
      child: LayoutBuilder(
        builder: (context, constraints) {
          final width = constraints.maxWidth / 2;
          return Stack(
            children: [
              AnimatedPositioned(
                duration: AppMotion.base,
                curve: AppMotion.standard,
                left: selectedIndex * width,
                top: 4,
                bottom: 4,
                width: width,
                child: Container(
                  decoration: BoxDecoration(
                    color: AppColors.textPrimary.withAlpha(28),
                    borderRadius: BorderRadius.circular(AppSizes.radiusFull),
                    boxShadow: AppShadows.sm,
                  ),
                ),
              ),
              Row(
                children: [
                  _TabButton(
                    label: 'All Classes',
                    selected: selectedIndex == 0,
                    onTap: () => onChanged(0),
                  ),
                  _TabButton(
                    label: 'Description',
                    selected: selectedIndex == 1,
                    onTap: () => onChanged(1),
                  ),
                ],
              ),
            ],
          );
        },
      ),
    );
  }
}

class _TabButton extends StatelessWidget {
  const _TabButton({
    required this.label,
    required this.selected,
    required this.onTap,
  });

  final String label;
  final bool selected;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    return Expanded(
      child: AppPressable(
        onTap: onTap,
        borderRadius: AppSizes.radiusFull,
        child: Center(
          child: AnimatedDefaultTextStyle(
            duration: AppMotion.base,
            style: AppTextStyles.labelMedium.copyWith(
              color: selected ? AppColors.textPrimary : AppColors.textSecondary,
            ),
            child: Text(label),
          ),
        ),
      ),
    );
  }
}

class _CurriculumTab extends StatelessWidget {
  const _CurriculumTab({
    super.key,
    required this.detail,
    required this.canOpenLessons,
  });

  final CourseDetailEntity detail;
  final bool canOpenLessons;

  @override
  Widget build(BuildContext context) {
    if (detail.sections.isEmpty) {
      return const CourseEmptyState(
        icon: Icons.menu_book_outlined,
        title: 'No curriculum yet',
        body: 'Lessons will appear here when this course is ready.',
      );
    }

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        AppSectionHeader(
          title: 'All Classes',
          subtitle: '${detail.sections.length} sections',
        ),
        const SizedBox(height: AppSizes.md),
        ...detail.sections.asMap().entries.map(
              (entry) => Padding(
                padding: const EdgeInsets.only(bottom: AppSizes.md),
                child: CourseSectionTile(
                  section: entry.value,
                  sectionNumber: entry.key + 1,
                  initiallyExpanded: true,
                  canOpenLessons: canOpenLessons,
                ),
              ),
            ),
      ],
    );
  }
}

class _DescriptionTab extends StatelessWidget {
  const _DescriptionTab({
    super.key,
    required this.course,
  });

  final CourseEntity course;

  @override
  Widget build(BuildContext context) {
    final description = course.description.isNotEmpty
        ? course.description
        : course.shortDescription;

    return AppSurface(
      padding: const EdgeInsets.all(AppSizes.lg),
      color: AppColors.bgSecondary,
      radius: AppSizes.radiusXl,
      shadows: AppShadows.sm,
      child: Text(
        description.isNotEmpty
            ? description
            : 'Course details will appear here when this class is ready.',
        style: AppTextStyles.bodyMedium.copyWith(
          color: AppColors.textSecondary,
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
        radius: AppSizes.radiusXl,
        color: AppColors.bgSecondary,
        shadows: AppShadows.lg,
        child: Row(
          children: [
            const _BookmarkAction(),
            const SizedBox(width: AppSizes.md),
            Expanded(
              child: accessAsync.when(
                loading: () => const AppButton(
                  label: 'Checking access...',
                  onPressed: null,
                  isLoading: true,
                ),
                error: (_, __) => AppButton(
                  label: 'Retry',
                  variant: AppButtonVariant.outline,
                  onPressed: () =>
                      ref.invalidate(courseAccessProvider(_courseId)),
                ),
                data: (access) {
                  if (access.hasAccess) {
                    return AppButton(
                      label: 'Continue Learning',
                      onPressed: _continueLearning,
                      suffixIcon: Icons.play_arrow_rounded,
                    );
                  }
                  if (widget.course.isFree) {
                    return AppButton(
                      label: 'Start Learning',
                      isLoading: busy,
                      onPressed: busy ? null : _startFree,
                      suffixIcon: Icons.arrow_forward_rounded,
                    );
                  }
                  final latestPayment =
                      ref.watch(latestCoursePaymentProvider(_courseId));
                  if (latestPayment != null && latestPayment.isPending) {
                    return AppButton(
                      label: 'Payment Pending',
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
          ],
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

class _BookmarkAction extends StatelessWidget {
  const _BookmarkAction();

  @override
  Widget build(BuildContext context) {
    return AppPressable(
      onTap: () {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Saved for later.')),
        );
      },
      borderRadius: AppSizes.radiusFull,
      child: Container(
        width: AppSizes.buttonHeight,
        height: AppSizes.buttonHeight,
        alignment: Alignment.center,
        decoration: BoxDecoration(
          color: AppColors.bgTertiary,
          shape: BoxShape.circle,
          border: Border.all(color: AppColors.border),
        ),
        child: const Icon(
          Icons.bookmark_border_rounded,
          color: AppColors.accentPrimary,
          size: AppSizes.iconMd,
        ),
      ),
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
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: AppColors.bgSecondary,
        borderRadius: BorderRadius.circular(16),
        border: Border.all(
          color: hasAccess ? const Color(0xFF4ADE80) : const Color(0x44FF5252),
        ),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              Container(
                padding: const EdgeInsets.all(10),
                decoration: BoxDecoration(
                  color: hasAccess
                      ? const Color(0x224ADE80)
                      : (isFree ? const Color(0x224ADE80) : const Color(0x22FF5252)),
                  shape: BoxShape.circle,
                ),
                child: Icon(
                  hasAccess
                      ? Icons.verified_rounded
                      : (isFree ? Icons.lock_open_rounded : Icons.workspace_premium_rounded),
                  color: hasAccess
                      ? const Color(0xFF4ADE80)
                      : (isFree ? const Color(0xFF4ADE80) : const Color(0xFFFF5252)),
                  size: 22,
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      hasAccess
                          ? 'Full Course Access Unlocked'
                          : (isFree ? '100% Free Course' : 'Lifetime Access • $label'),
                      style: AppTextStyles.titleMedium.copyWith(
                        fontWeight: FontWeight.bold,
                        fontSize: 15,
                      ),
                    ),
                    const SizedBox(height: 2),
                    Text(
                      hasAccess
                          ? 'You are enrolled. Open any lesson below to learn.'
                          : 'One-time payment • No monthly subscription',
                      style: AppTextStyles.caption.copyWith(
                        color: AppColors.textMuted,
                        fontSize: 12,
                      ),
                    ),
                  ],
                ),
              ),
            ],
          ),
          const SizedBox(height: 12),
          const Divider(height: 1, color: AppColors.border),
          const SizedBox(height: 10),
          const Wrap(
            spacing: 16,
            runSpacing: 8,
            children: [
              _FeatureCheck(label: 'Full HD Video Lessons'),
              _FeatureCheck(label: 'Downloadable PDF Notes'),
              _FeatureCheck(label: 'Interactive Quizzes'),
              _FeatureCheck(label: 'Certificate Included'),
            ],
          ),
        ],
      ),
    );
  }
}

class _FeatureCheck extends StatelessWidget {
  const _FeatureCheck({required this.label});

  final String label;

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        const Icon(Icons.check_circle_rounded, size: 14, color: Color(0xFF4ADE80)),
        const SizedBox(width: 6),
        Text(
          label,
          style: AppTextStyles.caption.copyWith(
            color: AppColors.textSecondary,
            fontSize: 11,
            fontWeight: FontWeight.w600,
          ),
        ),
      ],
    );
  }
}

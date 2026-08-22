import 'package:cached_network_image/cached_network_image.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/errors/failures.dart';
import '../../../../core/auth/account_gate.dart';
import '../../../../core/router/app_router.dart';
import '../../../../shared/widgets/app_button.dart';
import '../../../../shared/widgets/memere_mascot.dart';
import '../../../payment/presentation/providers/checkout_flow_provider.dart';
import '../../../payment/presentation/providers/course_access_provider.dart';
import '../../../payment/presentation/providers/purchase_history_provider.dart';
import '../../../payment/presentation/widgets/payment_provider_sheet.dart';
import '../../domain/entities/course_detail_entity.dart';
import '../../domain/entities/course_entity.dart';
import '../../domain/entities/lesson_entity.dart';
import '../../../exam/domain/entities/mock_exam_entity.dart';
import '../../../exam/presentation/providers/exam_providers.dart';
import '../../../quiz/presentation/providers/quiz_providers.dart';
import '../../../saved/presentation/providers/saved_courses_provider.dart';
import '../providers/course_detail_provider.dart';
import '../providers/course_download_provider.dart';
import '../widgets/course_detail_skeleton.dart';
import '../widgets/course_empty_state.dart';
import '../widgets/lesson_tile.dart';

/// Clean, refined Course Detail Screen matching image copy 5.png (Screen 3)
///
/// Features:
/// - Top App Bar with back button and centered Course Title
/// - Hero Illustration Area with character mascot or course thumbnail
/// - Curved Course Content Sheet with numbered lessons (01. Title -> [ ▶ ])
/// - Bottom Sticky Checkout/Enrollment Action Bar
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
              color: AppColors.brandEmerald,
              backgroundColor: AppColors.bgSecondary,
              onRefresh: () async {
                ref.invalidate(courseDetailProvider(courseId));
                ref.invalidate(courseAccessProvider(courseId));
                await ref.read(courseDetailProvider(courseId).future);
              },
              child: CustomScrollView(
                physics: const AlwaysScrollableScrollPhysics(),
                slivers: [
                  // 1. Top Bar with back button & course title
                  SliverToBoxAdapter(
                    child: _DetailTopBar(course: course),
                  ),

                  // 2. Hero Illustration Area (Mascot or Thumbnail)
                  SliverToBoxAdapter(
                    child: _HeroIllustrationArea(course: course),
                  ),

                  // 3. Curved Course Content Sheet
                  SliverToBoxAdapter(
                    child: _CourseContentSheet(
                      detail: detail,
                      hasAccess: hasAccess,
                      selectedTab: _selectedTab,
                      onTabChanged: (index) =>
                          setState(() => _selectedTab = index),
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
}

/// Toggles the course in the local favorites store and confirms with a
/// snackbar. Shared by the top-bar bookmark and the checkout-bar heart so both
/// controls stay in sync.
Future<void> _toggleFavorite(
  BuildContext context,
  WidgetRef ref,
  CourseEntity course,
) async {
  final nowSaved = await ref.read(savedCoursesProvider.notifier).toggle(course);
  if (!context.mounted) return;
  ScaffoldMessenger.of(context).showSnackBar(
    SnackBar(
      content: Text(nowSaved ? 'Added to My Courses.' : 'Removed from My Courses.'),
      duration: const Duration(seconds: 2),
    ),
  );
}

/// Top App Bar with back button and centered Title
class _DetailTopBar extends ConsumerWidget {
  const _DetailTopBar({required this.course});

  final CourseEntity course;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final topPadding = MediaQuery.paddingOf(context).top;
    final saved = ref.watch(savedCoursesProvider).valueOrNull ?? const [];
    final isSaved = saved.any((item) => item.id == course.id);

    return Padding(
      padding: EdgeInsets.fromLTRB(
        AppSizes.screenPaddingH,
        topPadding + 10,
        AppSizes.screenPaddingH,
        8,
      ),
      child: Row(
        children: [
          // Back Button
          InkWell(
            onTap: () {
              if (context.canPop()) {
                context.pop();
              } else {
                context.go(AppRoutes.home);
              }
            },
            borderRadius: BorderRadius.circular(20),
            child: Container(
              width: 38,
              height: 38,
              alignment: Alignment.center,
              decoration: BoxDecoration(
                color: AppColors.bgSecondary,
                shape: BoxShape.circle,
                border: Border.all(color: AppColors.borderStrong),
              ),
              child: const Icon(
                Icons.arrow_back_ios_new_rounded,
                color: AppColors.textPrimary,
                size: 16,
              ),
            ),
          ),
          const SizedBox(width: 12),

          // Course Subject / Title
          Expanded(
            child: Text(
              course.subject.isNotEmpty ? course.subject : course.title,
              textAlign: TextAlign.center,
              maxLines: 1,
              overflow: TextOverflow.ellipsis,
              style: const TextStyle(
                fontSize: 17,
                fontWeight: FontWeight.bold,
                color: AppColors.textPrimary,
                letterSpacing: -0.3,
              ),
            ),
          ),
          const SizedBox(width: 12),

          // Action / Bookmark Button — toggles the course in My Courses
          InkWell(
            onTap: () => _toggleFavorite(context, ref, course),
            borderRadius: BorderRadius.circular(20),
            child: Container(
              width: 38,
              height: 38,
              alignment: Alignment.center,
              decoration: BoxDecoration(
                color: AppColors.bgSecondary,
                shape: BoxShape.circle,
                border: Border.all(
                  color:
                      isSaved ? AppColors.brandEmerald : AppColors.borderStrong,
                ),
              ),
              child: Icon(
                isSaved
                    ? Icons.bookmark_rounded
                    : Icons.bookmark_border_rounded,
                color: isSaved ? AppColors.brandEmerald : AppColors.textPrimary,
                size: 18,
              ),
            ),
          ),
        ],
      ),
    );
  }
}

/// Hero Illustration Area (matching image copy 5.png Screen 3)
class _HeroIllustrationArea extends StatelessWidget {
  const _HeroIllustrationArea({required this.course});

  final CourseEntity course;

  @override
  Widget build(BuildContext context) {
    final hasImage = course.thumbnailUrl != null &&
        course.thumbnailUrl!.trim().isNotEmpty &&
        course.thumbnailUrl!.startsWith('http');

    return Container(
      width: double.infinity,
      padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 12),
      child: Center(
        child: hasImage
            ? ClipRRect(
                borderRadius: BorderRadius.circular(20),
                child: SizedBox(
                  height: 180,
                  width: double.infinity,
                  child: CachedNetworkImage(
                    imageUrl: course.thumbnailUrl!,
                    fit: BoxFit.cover,
                    placeholder: (_, __) => _MascotHeroCanvas(course: course),
                    errorWidget: (_, __, ___) =>
                        _MascotHeroCanvas(course: course),
                  ),
                ),
              )
            : _MascotHeroCanvas(course: course),
      ),
    );
  }
}

/// Animated Mascot Hero Canvas with Subject & Grade Pills
class _MascotHeroCanvas extends StatelessWidget {
  const _MascotHeroCanvas({required this.course});

  final CourseEntity course;

  @override
  Widget build(BuildContext context) {
    return Container(
      height: 180,
      width: double.infinity,
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: AppColors.bgSecondary,
        borderRadius: BorderRadius.circular(24),
        border: Border.all(color: AppColors.borderStrong),
      ),
      child: Stack(
        children: [
          // Background ambient rings
          Positioned(
            right: -20,
            top: -20,
            child: Container(
              width: 110,
              height: 110,
              decoration: BoxDecoration(
                shape: BoxShape.circle,
                color: AppColors.brandEmerald.withAlpha(12),
              ),
            ),
          ),

          // Mascot Character Illustration (standing on card surface)
          const Center(
            child: MemereMascot(
              size: Size(150, 138),
              showBackdrop: false,
            ),
          ),

          // Left Grade & Subject Badges
          Positioned(
            left: 0,
            top: 0,
            child: Container(
              padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 5),
              decoration: BoxDecoration(
                color: AppColors.bgTertiary,
                borderRadius: BorderRadius.circular(10),
                border: Border.all(color: AppColors.borderStrong),
              ),
              child: Text(
                'Grade ${course.grade}',
                style: const TextStyle(
                  fontSize: 11,
                  fontWeight: FontWeight.bold,
                  color: AppColors.brandEmerald,
                ),
              ),
            ),
          ),

          // Right Level / Free Badge
          Positioned(
            right: 0,
            top: 0,
            child: Container(
              padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 5),
              decoration: BoxDecoration(
                color: course.isFree
                    ? const Color(0x2210B981)
                    : const Color(0x2238BDF8),
                borderRadius: BorderRadius.circular(10),
                border: Border.all(
                  color: course.isFree
                      ? const Color(0x6610B981)
                      : const Color(0x6638BDF8),
                ),
              ),
              child: Text(
                course.isFree ? '100% Free' : course.priceLabel,
                style: TextStyle(
                  fontSize: 11,
                  fontWeight: FontWeight.bold,
                  color: course.isFree
                      ? AppColors.brandEmerald
                      : const Color(0xFF38BDF8),
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }
}

/// Curved Course Content Sheet matching image copy 5.png (Screen 3)
class _CourseContentSheet extends StatefulWidget {
  const _CourseContentSheet({
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
  State<_CourseContentSheet> createState() => _CourseContentSheetState();
}

class _CourseContentSheetState extends State<_CourseContentSheet> {
  late Set<int> _expandedSections;

  @override
  void initState() {
    super.initState();
    // Expand all sections initially so students can view content and tap to shrink
    _expandedSections = Set<int>.from(
      List.generate(widget.detail.sections.length, (index) => index),
    );
  }

  void _toggleSection(int index) {
    setState(() {
      if (_expandedSections.contains(index)) {
        _expandedSections.remove(index);
      } else {
        _expandedSections.add(index);
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    final course = widget.detail.course;
    final totalLessons = widget.detail.lessonCount;

    return Container(
      width: double.infinity,
      decoration: const BoxDecoration(
        color: Color(0xFF11141E),
        borderRadius: BorderRadius.vertical(top: Radius.circular(32)),
        border: Border(
          top: BorderSide(color: AppColors.borderStrong),
        ),
      ),
      padding: const EdgeInsets.fromLTRB(20, 14, 20, 40),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Top Sheet Handle Indicator
          Center(
            child: Container(
              width: 38,
              height: 4,
              decoration: BoxDecoration(
                color: const Color(0xFF334155),
                borderRadius: BorderRadius.circular(2),
              ),
            ),
          ),
          const SizedBox(height: 16),

          // Course Title & Metrics Header
          Text(
            course.title,
            style: const TextStyle(
              fontSize: 20,
              fontWeight: FontWeight.w900,
              color: Colors.white,
              letterSpacing: -0.4,
              height: 1.25,
            ),
          ),
          const SizedBox(height: 6),
          Text(
            '${course.durationLabel} • ${course.totalLessons} lessons • Grade ${course.grade}',
            style: const TextStyle(
              fontSize: 13,
              color: Color(0xFF94A3B8),
              fontWeight: FontWeight.w500,
            ),
          ),
          const SizedBox(height: 16),

          // Download the whole course for offline study — every video, PDF,
          // quiz and exam, all in one tap.
          _DownloadWholeCourseButton(
            detail: widget.detail,
            canDownload: course.isFree || widget.hasAccess,
          ),
          const SizedBox(height: 16),

          // Tab Switcher (Lessons / Quizzes / Exams / About)
          Container(
            height: 42,
            padding: const EdgeInsets.all(3),
            decoration: BoxDecoration(
              color: AppColors.bgSecondary,
              borderRadius: BorderRadius.circular(12),
              border: Border.all(color: AppColors.borderStrong),
            ),
            child: Row(
              children: [
                _TabButton(
                  label: 'Lessons',
                  selected: widget.selectedTab == 0,
                  onTap: () => widget.onTabChanged(0),
                ),
                _TabButton(
                  label: 'Quizzes',
                  selected: widget.selectedTab == 1,
                  onTap: () => widget.onTabChanged(1),
                ),
                _TabButton(
                  label: 'Exams',
                  selected: widget.selectedTab == 2,
                  onTap: () => widget.onTabChanged(2),
                ),
                _TabButton(
                  label: 'About',
                  selected: widget.selectedTab == 3,
                  onTap: () => widget.onTabChanged(3),
                ),
              ],
            ),
          ),
          const SizedBox(height: 20),

          // Content Tabs
          if (widget.selectedTab == 0) ...[
            if (totalLessons == 0)
              const Padding(
                padding: EdgeInsets.symmetric(vertical: 32),
                child: Center(
                  child: Text(
                    'Lessons will appear here when this course is published.',
                    textAlign: TextAlign.center,
                    style: TextStyle(fontSize: 13, color: Color(0xFF64748B)),
                  ),
                ),
              )
            else
              ..._buildSectionsList(course, widget.hasAccess),
          ] else if (widget.selectedTab == 1) ...[
            _CourseQuizzesTab(
              courseId: course.id,
              hasAccess: widget.hasAccess,
              isFree: course.isFree,
            ),
          ] else if (widget.selectedTab == 2) ...[
            _CourseExamsTab(
              courseId: course.id,
              hasAccess: widget.hasAccess,
              isFree: course.isFree,
            ),
          ] else ...[
            Container(
              width: double.infinity,
              padding: const EdgeInsets.all(16),
              decoration: BoxDecoration(
                color: AppColors.bgSecondary,
                borderRadius: BorderRadius.circular(16),
                border: Border.all(color: AppColors.borderStrong),
              ),
              child: Text(
                course.description.isNotEmpty
                    ? course.description
                    : (course.shortDescription.isNotEmpty
                        ? course.shortDescription
                        : 'Comprehensive lessons and practice materials prepared for national entrance exams.'),
                style: const TextStyle(
                  fontSize: 13.5,
                  color: AppColors.textSecondary,
                  height: 1.45,
                ),
              ),
            ),
          ],
        ],
      ),
    );
  }

  List<Widget> _buildSectionsList(CourseEntity course, bool hasAccess) {
    final widgets = <Widget>[];
    int globalLessonNumber = 0;

    for (int i = 0; i < widget.detail.sections.length; i++) {
      final section = widget.detail.sections[i];
      final isFirst = i == 0;
      final isExpanded = _expandedSections.contains(i);

      // Section Header (Tappable to expand and shrink)
      widgets.add(
        Padding(
          padding: EdgeInsets.only(
            top: isFirst ? 0 : 16,
            bottom: 6,
          ),
          child: InkWell(
            onTap: () => _toggleSection(i),
            borderRadius: BorderRadius.circular(10),
            child: Padding(
              padding: const EdgeInsets.symmetric(horizontal: 4, vertical: 6),
              child: Row(
                children: [
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          section.title.isNotEmpty
                              ? section.title
                              : 'Section ${i + 1}',
                          maxLines: 2,
                          overflow: TextOverflow.ellipsis,
                          style: const TextStyle(
                            fontSize: 14.5,
                            fontWeight: FontWeight.w700,
                            color: Colors.white,
                            letterSpacing: -0.2,
                          ),
                        ),
                        if (section.description.isNotEmpty) ...[
                          const SizedBox(height: 2),
                          Text(
                            section.description,
                            maxLines: 1,
                            overflow: TextOverflow.ellipsis,
                            style: const TextStyle(
                              fontSize: 12,
                              color: Color(0xFF64748B),
                            ),
                          ),
                        ],
                      ],
                    ),
                  ),
                  const SizedBox(width: 8),
                  Text(
                    section.lessonCountLabel,
                    maxLines: 1,
                    overflow: TextOverflow.ellipsis,
                    style: const TextStyle(
                      fontSize: 12,
                      color: Color(0xFF94A3B8),
                      fontWeight: FontWeight.w500,
                    ),
                  ),
                  const SizedBox(width: 6),
                  AnimatedRotation(
                    turns: isExpanded ? 0.5 : 0.0,
                    duration: const Duration(milliseconds: 180),
                    child: const Icon(
                      Icons.keyboard_arrow_down_rounded,
                      size: 20,
                      color: Color(0xFF94A3B8),
                    ),
                  ),
                ],
              ),
            ),
          ),
        ),
      );

      if (isExpanded) {
        if (section.lessons.isEmpty) {
          widgets.add(
            const Padding(
              padding: EdgeInsets.symmetric(vertical: 8, horizontal: 4),
              child: Text(
                'No lessons in this section yet.',
                style: TextStyle(fontSize: 12, color: Color(0xFF64748B)),
              ),
            ),
          );
        } else {
          for (int j = 0; j < section.lessons.length; j++) {
            final lesson = section.lessons[j];
            globalLessonNumber++;
            widgets.add(
              Padding(
                padding: const EdgeInsets.only(bottom: 10),
                child: LessonTile(
                  lesson: lesson,
                  lessonNumber: globalLessonNumber,
                  canOpen: course.isFree || hasAccess,
                ),
              ),
            );
          }
        }
      } else {
        // Increment global counter when collapsed to keep overall sequence intact
        globalLessonNumber += section.lessons.length;
      }
    }

    return widgets;
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
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(10),
        child: Container(
          alignment: Alignment.center,
          decoration: BoxDecoration(
            color: selected ? AppColors.brandEmerald : Colors.transparent,
            borderRadius: BorderRadius.circular(10),
          ),
          child: Text(
            label,
            style: TextStyle(
              fontSize: 13,
              fontWeight: selected ? FontWeight.bold : FontWeight.w600,
              color: selected ? Colors.white : AppColors.textSecondary,
            ),
          ),
        ),
      ),
    );
  }
}

/// One-tap "download the whole course" control — fetches every video, PDF,
/// quiz and exam for offline use, showing one aggregate progress value. Gated
/// on access (free or enrolled); otherwise it nudges the student to enroll.
class _DownloadWholeCourseButton extends ConsumerWidget {
  const _DownloadWholeCourseButton({
    required this.detail,
    required this.canDownload,
  });

  final CourseDetailEntity detail;
  final bool canDownload;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final courseId = detail.course.id;
    final progress = ref.watch(courseDownloadProvider(courseId));
    final running = progress.isRunning;
    final done = progress.status == CourseDownloadStatus.done;
    final failedAll = progress.status == CourseDownloadStatus.failed;
    final offlineReady = done && progress.failed == 0;

    final label = running
        ? 'Downloading ${progress.processed}/${progress.total} • ${progress.percent}%'
        : offlineReady
            ? 'Course downloaded • Offline ready'
            : done
                ? 'Downloaded ${progress.completed}/${progress.total} • ${progress.failed} unavailable'
                : failedAll
                    ? 'Download failed • Tap to retry'
                    : 'Download whole course';

    final accent = offlineReady
        ? AppColors.brandEmerald
        : (canDownload ? Colors.white : AppColors.textMuted);

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Material(
          color: offlineReady ? const Color(0x1810B981) : AppColors.bgSecondary,
          borderRadius: BorderRadius.circular(12),
          child: InkWell(
            borderRadius: BorderRadius.circular(12),
            onTap: (!canDownload || running)
                ? null
                : () => _run(context, ref, courseId),
            child: Container(
              width: double.infinity,
              padding:
                  const EdgeInsets.symmetric(horizontal: 14, vertical: 12),
              decoration: BoxDecoration(
                borderRadius: BorderRadius.circular(12),
                border: Border.all(
                  color: offlineReady
                      ? const Color(0x5510B981)
                      : AppColors.borderStrong,
                ),
              ),
              child: Row(
                children: [
                  running
                      ? const SizedBox(
                          width: 18,
                          height: 18,
                          child: CircularProgressIndicator(
                            strokeWidth: 2,
                            color: AppColors.brandEmerald,
                          ),
                        )
                      : Icon(
                          offlineReady
                              ? Icons.check_circle_rounded
                              : failedAll
                                  ? Icons.refresh_rounded
                                  : Icons.download_for_offline_outlined,
                          size: 20,
                          color: accent,
                        ),
                  const SizedBox(width: 10),
                  Expanded(
                    child: Text(
                      label,
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                      style: TextStyle(
                        fontSize: 13.5,
                        fontWeight: FontWeight.w700,
                        color: accent,
                      ),
                    ),
                  ),
                ],
              ),
            ),
          ),
        ),
        if (running) ...[
          const SizedBox(height: 8),
          ClipRRect(
            borderRadius: BorderRadius.circular(4),
            child: LinearProgressIndicator(
              value: progress.fraction == 0 ? null : progress.fraction,
              minHeight: 4,
              backgroundColor: AppColors.bgTertiary,
              valueColor:
                  const AlwaysStoppedAnimation(AppColors.brandEmerald),
            ),
          ),
          if (progress.currentLabel != null) ...[
            const SizedBox(height: 6),
            Text(
              progress.currentLabel!,
              maxLines: 1,
              overflow: TextOverflow.ellipsis,
              style: const TextStyle(fontSize: 11.5, color: Color(0xFF94A3B8)),
            ),
          ],
        ] else if (!canDownload) ...[
          const SizedBox(height: 6),
          const Text(
            'Enroll to download this course for offline study.',
            style: TextStyle(fontSize: 11.5, color: Color(0xFF94A3B8)),
          ),
        ],
      ],
    );
  }

  Future<void> _run(
    BuildContext context,
    WidgetRef ref,
    String courseId,
  ) async {
    final messenger = ScaffoldMessenger.of(context);
    await ref
        .read(courseDownloadProvider(courseId).notifier)
        .downloadCourse(detail);
    if (!context.mounted) return;
    final result = ref.read(courseDownloadProvider(courseId));
    final String message;
    if (result.total == 0) {
      message = 'Nothing to download for this course yet.';
    } else if (result.failed == 0) {
      message = 'Course downloaded for offline study.';
    } else {
      message =
          'Downloaded ${result.completed}/${result.total} — ${result.failed} item(s) unavailable.';
    }
    messenger.showSnackBar(SnackBar(content: Text(message)));
  }
}

/// Bottom Sticky CTA that reflects real access state and drives checkout
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
    final saved = ref.watch(savedCoursesProvider).valueOrNull ?? const [];
    final isSaved = saved.any((item) => item.id == widget.course.id);

    return SafeArea(
      top: false,
      child: Container(
        padding: const EdgeInsets.symmetric(horizontal: 18, vertical: 12),
        decoration: const BoxDecoration(
          color: Color(0xFF0F121B),
          border: Border(
            top: BorderSide(color: AppColors.borderStrong),
          ),
        ),
        child: Row(
          children: [
            // Bookmark Heart Action Button — toggles the course in My Courses
            InkWell(
              onTap: () => _toggleFavorite(context, ref, widget.course),
              borderRadius: BorderRadius.circular(14),
              child: Container(
                width: 48,
                height: 48,
                alignment: Alignment.center,
                decoration: BoxDecoration(
                  color: AppColors.bgSecondary,
                  borderRadius: BorderRadius.circular(14),
                  border: Border.all(
                    color: isSaved
                        ? AppColors.brandEmerald
                        : AppColors.borderStrong,
                  ),
                ),
                child: Icon(
                  isSaved
                      ? Icons.favorite_rounded
                      : Icons.favorite_border_rounded,
                  color:
                      isSaved ? AppColors.brandEmerald : AppColors.textPrimary,
                  size: 20,
                ),
              ),
            ),
            const SizedBox(width: 12),

            // Primary Action Button (in brand emerald #10B981)
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
                      suffixIcon: Icons.arrow_forward_rounded,
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

  void _continueLearning() {
    final detail = ref.read(courseDetailProvider(_courseId)).valueOrNull;
    final lesson = _firstPlayableLesson(detail);
    if (lesson == null) {
      _showMessage('Lessons will appear here when this course is ready.');
      return;
    }
    if (lesson.hasVideo || lesson.type == LessonType.video || (lesson.videoId != null && lesson.videoId!.isNotEmpty)) {
      final effectiveVideoId = (lesson.videoId != null && lesson.videoId!.isNotEmpty)
          ? lesson.videoId!
          : lesson.id;
      context.push(
        AppRoutes.videoPlayerPath(
          videoId: effectiveVideoId,
          lessonId: lesson.id,
          courseId: lesson.courseId,
          title: lesson.title,
        ),
      );
      return;
    }
    if (lesson.hasQuiz && lesson.quizId != null && lesson.quizId!.isNotEmpty) {
      context.push(AppRoutes.quizDetailPath(lesson.quizId!));
      return;
    }
    final pdfName = lesson.pdfUrl ?? '';
    context.push(
      AppRoutes.pdfReaderPath(
        title: lesson.title,
        pdfUrl: pdfName,
        lessonId: lesson.id,
        content: lesson.content,
      ),
      extra: <String, dynamic>{
        'title': lesson.title,
        'pdfUrl': pdfName,
        'lessonId': lesson.id,
        'content': lesson.content,
      },
    );
  }

  LessonEntity? _firstPlayableLesson(CourseDetailEntity? detail) {
    if (detail == null) return null;
    for (final section in detail.sections) {
      if (section.lessons.isNotEmpty) {
        return section.lessons.first;
      }
    }
    return null;
  }

  Future<void> _startFree() async {
    if (!await requireAccount(
      context,
      ref,
      title: 'Sign in to enroll',
      message:
          'Create a free account or sign in to enroll and track your progress. '
          'Your downloads and saved items stay on this device either way.',
    )) {
      return;
    }
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
    if (!await requireAccount(
      context,
      ref,
      title: 'Sign in to purchase',
      message:
          'Create a free account or sign in to buy this course. Your purchase '
          'unlocks it across your devices.',
    )) {
      return;
    }
    if (!mounted) return;
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

class _CourseQuizzesTab extends ConsumerWidget {
  const _CourseQuizzesTab({
    required this.courseId,
    required this.hasAccess,
    required this.isFree,
  });

  final String courseId;
  final bool hasAccess;
  final bool isFree;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final quizzesAsync = ref.watch(courseQuizzesProvider(courseId));
    final canOpen = isFree || hasAccess;

    return quizzesAsync.when(
      loading: () => const Padding(
        padding: EdgeInsets.symmetric(vertical: 24),
        child: Center(
          child: CircularProgressIndicator(
            color: AppColors.brandEmerald,
            strokeWidth: 2,
          ),
        ),
      ),
      error: (err, _) => Padding(
        padding: const EdgeInsets.symmetric(vertical: 20),
        child: Text(
          'Could not load quizzes.',
          style: TextStyle(color: Colors.red.shade300, fontSize: 13),
        ),
      ),
      data: (quizzes) {
        if (quizzes.isEmpty) {
          return Container(
            width: double.infinity,
            padding: const EdgeInsets.all(20),
            decoration: BoxDecoration(
              color: AppColors.bgSecondary,
              borderRadius: BorderRadius.circular(16),
              border: Border.all(color: AppColors.borderStrong),
            ),
            child: const Column(
              children: [
                Icon(
                  Icons.quiz_outlined,
                  size: 36,
                  color: Color(0xFF64748B),
                ),
                SizedBox(height: 10),
                Text(
                  'No standalone quizzes found for this course.',
                  textAlign: TextAlign.center,
                  style: TextStyle(
                    fontSize: 14,
                    fontWeight: FontWeight.w600,
                    color: Colors.white,
                  ),
                ),
                SizedBox(height: 4),
                Text(
                  'Practice quizzes may also be attached directly to individual video lessons in the Lessons tab.',
                  textAlign: TextAlign.center,
                  style: TextStyle(
                    fontSize: 12,
                    color: Color(0xFF94A3B8),
                  ),
                ),
              ],
            ),
          );
        }

        return Column(
          children: quizzes.map((quiz) {
            return Padding(
              padding: const EdgeInsets.only(bottom: 10),
              child: InkWell(
                onTap: () {
                  if (!canOpen) {
                    ScaffoldMessenger.of(context).showSnackBar(
                      const SnackBar(
                        content: Text('Enroll in this course to unlock quizzes.'),
                      ),
                    );
                    return;
                  }
                  context.push(AppRoutes.quizDetailPath(quiz.id));
                },
                borderRadius: BorderRadius.circular(16),
                child: Container(
                  padding:
                      const EdgeInsets.symmetric(horizontal: 14, vertical: 12),
                  decoration: BoxDecoration(
                    color: AppColors.bgSecondary,
                    borderRadius: BorderRadius.circular(16),
                    border: Border.all(
                      color: AppColors.borderStrong.withAlpha(90),
                    ),
                  ),
                  child: Row(
                    children: [
                      Container(
                        width: 38,
                        height: 38,
                        alignment: Alignment.center,
                        decoration: BoxDecoration(
                          color: const Color(0x2210B981),
                          borderRadius: BorderRadius.circular(10),
                          border: Border.all(color: const Color(0x5510B981)),
                        ),
                        child: const Icon(
                          Icons.quiz_outlined,
                          size: 18,
                          color: Color(0xFF10B981),
                        ),
                      ),
                      const SizedBox(width: 12),
                      Expanded(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(
                              quiz.title,
                              maxLines: 2,
                              overflow: TextOverflow.ellipsis,
                              style: const TextStyle(
                                fontSize: 14.5,
                                fontWeight: FontWeight.w700,
                                color: Colors.white,
                              ),
                            ),
                            const SizedBox(height: 3),
                            Text(
                              '${quiz.questionCount} Questions • ${quiz.passPercentage.toStringAsFixed(0)}% Pass Mark',
                              style: const TextStyle(
                                fontSize: 12,
                                color: Color(0xFF64748B),
                                fontWeight: FontWeight.w500,
                              ),
                            ),
                          ],
                        ),
                      ),
                      const Icon(
                        Icons.arrow_forward_ios_rounded,
                        size: 14,
                        color: Color(0xFF94A3B8),
                      ),
                    ],
                  ),
                ),
              ),
            );
          }).toList(),
        );
      },
    );
  }
}

class _CourseExamsTab extends ConsumerWidget {
  const _CourseExamsTab({
    required this.courseId,
    required this.hasAccess,
    required this.isFree,
  });

  final String courseId;
  final bool hasAccess;
  final bool isFree;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final examsAsync = ref.watch(courseExamsProvider(courseId));
    final canOpen = isFree || hasAccess;

    return examsAsync.when(
      loading: () => const Padding(
        padding: EdgeInsets.symmetric(vertical: 24),
        child: Center(
          child: CircularProgressIndicator(
            color: AppColors.brandEmerald,
            strokeWidth: 2,
          ),
        ),
      ),
      error: (err, _) => Padding(
        padding: const EdgeInsets.symmetric(vertical: 20),
        child: Text(
          'Could not load exams.',
          style: TextStyle(color: Colors.red.shade300, fontSize: 13),
        ),
      ),
      data: (exams) {
        if (exams.isEmpty) {
          return Container(
            width: double.infinity,
            padding: const EdgeInsets.all(20),
            decoration: BoxDecoration(
              color: AppColors.bgSecondary,
              borderRadius: BorderRadius.circular(16),
              border: Border.all(color: AppColors.borderStrong),
            ),
            child: const Column(
              children: [
                Icon(
                  Icons.assignment_turned_in_outlined,
                  size: 36,
                  color: Color(0xFF64748B),
                ),
                SizedBox(height: 10),
                Text(
                  'No national practice exams found for this course.',
                  textAlign: TextAlign.center,
                  style: TextStyle(
                    fontSize: 14,
                    fontWeight: FontWeight.w600,
                    color: Colors.white,
                  ),
                ),
                SizedBox(height: 4),
                Text(
                  'Explore general entrance mock exams in the Exams tab on the main screen.',
                  textAlign: TextAlign.center,
                  style: TextStyle(
                    fontSize: 12,
                    color: Color(0xFF94A3B8),
                  ),
                ),
              ],
            ),
          );
        }

        return Column(
          children: exams.map((exam) {
            return Padding(
              padding: const EdgeInsets.only(bottom: 10),
              child: InkWell(
                onTap: () {
                  if (!canOpen) {
                    ScaffoldMessenger.of(context).showSnackBar(
                      const SnackBar(
                        content: Text('Enroll in this course to unlock exams.'),
                      ),
                    );
                    return;
                  }
                  _startExam(context, ref, exam);
                },
                borderRadius: BorderRadius.circular(16),
                child: Container(
                  padding:
                      const EdgeInsets.symmetric(horizontal: 14, vertical: 12),
                  decoration: BoxDecoration(
                    color: AppColors.bgSecondary,
                    borderRadius: BorderRadius.circular(16),
                    border: Border.all(
                      color: AppColors.borderStrong.withAlpha(90),
                    ),
                  ),
                  child: Row(
                    children: [
                      Container(
                        width: 38,
                        height: 38,
                        alignment: Alignment.center,
                        decoration: BoxDecoration(
                          color: const Color(0x2238BDF8),
                          borderRadius: BorderRadius.circular(10),
                          border: Border.all(color: const Color(0x5538BDF8)),
                        ),
                        child: const Icon(
                          Icons.assignment_turned_in_rounded,
                          size: 18,
                          color: Color(0xFF38BDF8),
                        ),
                      ),
                      const SizedBox(width: 12),
                      Expanded(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(
                              exam.title,
                              maxLines: 2,
                              overflow: TextOverflow.ellipsis,
                              style: const TextStyle(
                                fontSize: 14.5,
                                fontWeight: FontWeight.w700,
                                color: Colors.white,
                              ),
                            ),
                            const SizedBox(height: 3),
                            Text(
                              'Grade ${exam.grade} • ${exam.durationMinutes} min • ${exam.totalMarks} Total Marks',
                              style: const TextStyle(
                                fontSize: 12,
                                color: Color(0xFF64748B),
                                fontWeight: FontWeight.w500,
                              ),
                            ),
                          ],
                        ),
                      ),
                      const Icon(
                        Icons.play_arrow_rounded,
                        size: 20,
                        color: Color(0xFF38BDF8),
                      ),
                    ],
                  ),
                ),
              ),
            );
          }).toList(),
        );
      },
    );
  }

  Future<void> _startExam(
      BuildContext context, WidgetRef ref, MockExamEntity exam) async {
    final useCase = ref.read(startExamUseCaseProvider);
    final result = await useCase(exam.id);
    result.fold(
      (failure) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text(failure.message)),
        );
      },
      (attempt) {
        context.push(
          AppRoutes.examAttemptPath(
            attemptId: attempt.attemptId,
            examId: exam.id,
          ),
        );
      },
    );
  }
}

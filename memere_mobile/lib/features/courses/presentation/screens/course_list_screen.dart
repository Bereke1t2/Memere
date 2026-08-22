import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../../../core/errors/failures.dart';
import '../../../../core/router/app_router.dart';
import '../../../../shared/widgets/memere_mascot.dart';
import '../../../auth/presentation/providers/auth_state_provider.dart';
import '../../../notifications/presentation/providers/announcement_provider.dart';
import '../providers/course_list_provider.dart';
import '../widgets/course_card.dart';
import '../widgets/course_empty_state.dart';
import '../widgets/course_list_skeleton.dart';

/// Grade shown as "Freshman" in the UI. There is no backend content for it yet
/// (the API serves grades 9–12), so selecting it short-circuits to a
/// coming-soon state instead of querying grade 13.
const _kAllGrade = 0;
const _kFreshmanGrade = 13;

/// Clean, professional Home Screen for Memere.
///
/// Features:
/// - Minimalist App Bar: 'Memere' on left, Profile avatar on right (navigates to Account)
/// - Overflow-Safe Hero Featured Banner with animated MemereMascot
/// - Grade-level filter chips (All, Freshman, Grade 12–9) + Subject chooser navigation underneath
/// - Unique, non-video Course Cards with optional thumbnail support
class CourseListScreen extends ConsumerStatefulWidget {
  const CourseListScreen({super.key});

  @override
  ConsumerState<CourseListScreen> createState() => _CourseListScreenState();
}

class _CourseListScreenState extends ConsumerState<CourseListScreen> {
  late final ScrollController _scrollController;
  bool _dismissedAnnouncement = false;

  /// The grade chip the student has selected. Held locally (not derived from
  /// the provider) so "Freshman" can stay highlighted while we short-circuit
  /// its fetch. Defaults to 0 (All), matching [CourseListState.selectedGrade].
  int _selectedGrade = _kAllGrade;

  @override
  void initState() {
    super.initState();
    _scrollController = ScrollController()..addListener(_onScroll);
  }

  @override
  void dispose() {
    _scrollController
      ..removeListener(_onScroll)
      ..dispose();
    super.dispose();
  }

  void _onScroll() {
    if (!_scrollController.hasClients) return;
    final position = _scrollController.position;
    if (position.pixels >= position.maxScrollExtent - 360) {
      ref.read(courseListProvider.notifier).loadMore();
    }
  }

  void _onSelectGrade(int grade) {
    setState(() => _selectedGrade = grade);
    // Freshman has no backend content yet — don't query grade 13; the content
    // area renders a coming-soon state while this chip is active.
    if (grade == _kFreshmanGrade) return;
    ref
        .read(courseListProvider.notifier)
        .setGrade(grade == _kAllGrade ? null : grade);
  }

  void _onSelectSubject(String? subject) {
    ref.read(courseListProvider.notifier).setSubject(subject);
  }

  void _resetFilters() {
    setState(() => _selectedGrade = _kAllGrade);
    ref.read(courseListProvider.notifier).clearFilters();
  }

  @override
  Widget build(BuildContext context) {
    final coursesAsync = ref.watch(courseListProvider);

    return Scaffold(
      backgroundColor: AppColors.bgPrimary,
      body: SafeArea(
        bottom: false,
        child: coursesAsync.when(
          loading: () {
            final previous = coursesAsync.valueOrNull;
            if (previous != null && previous.courses.isNotEmpty) {
              return _CourseListContent(
                controller: _scrollController,
                state: previous,
                selectedGrade: _selectedGrade,
                onSelectGrade: _onSelectGrade,
                onSelectSubject: _onSelectSubject,
                onResetFilters: _resetFilters,
                dismissedAnnouncement: _dismissedAnnouncement,
                onDismissAnnouncement: () =>
                    setState(() => _dismissedAnnouncement = true),
              );
            }
            return const Column(
              children: [
                _HomeTopBar(),
                SizedBox(height: AppSizes.md),
                Expanded(child: CourseListSkeleton()),
              ],
            );
          },
          error: (error, _) => Column(
            children: [
              const _HomeTopBar(),
              Expanded(
                child: CourseEmptyState(
                  icon: Icons.wifi_off_rounded,
                  title: 'Could not load courses',
                  body: error is Failure
                      ? error.message
                      : 'Check your connection and try again.',
                  buttonLabel: 'Retry',
                  onPressed: () => ref.invalidate(courseListProvider),
                ),
              ),
            ],
          ),
          data: (state) => _CourseListContent(
            controller: _scrollController,
            state: state,
            selectedGrade: _selectedGrade,
            onSelectGrade: _onSelectGrade,
            onSelectSubject: _onSelectSubject,
            onResetFilters: _resetFilters,
            dismissedAnnouncement: _dismissedAnnouncement,
            onDismissAnnouncement: () =>
                setState(() => _dismissedAnnouncement = true),
          ),
        ),
      ),
    );
  }
}

class _CourseListContent extends ConsumerWidget {
  const _CourseListContent({
    required this.controller,
    required this.state,
    required this.selectedGrade,
    required this.onSelectGrade,
    required this.onSelectSubject,
    required this.onResetFilters,
    required this.dismissedAnnouncement,
    required this.onDismissAnnouncement,
  });

  final ScrollController controller;
  final CourseListState state;
  final int selectedGrade;
  final void Function(int grade) onSelectGrade;
  final void Function(String? subject) onSelectSubject;
  final VoidCallback onResetFilters;
  final bool dismissedAnnouncement;
  final VoidCallback onDismissAnnouncement;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final announcementsAsync = ref.watch(announcementProvider);

    return RefreshIndicator(
      color: AppColors.brandEmerald,
      backgroundColor: AppColors.bgSecondary,
      onRefresh: () async {
        ref.read(courseListProvider.notifier).refresh();
        ref.invalidate(announcementProvider);
      },
      child: CustomScrollView(
        controller: controller,
        physics: const AlwaysScrollableScrollPhysics(),
        slivers: [
          // 1. Minimalist App Bar (Memere left, Profile right)
          const SliverToBoxAdapter(
            child: _HomeTopBar(),
          ),

          // 2. Overflow-Safe Hero Featured Banner with Character Mascot
          SliverToBoxAdapter(
            child: Padding(
              padding: const EdgeInsets.fromLTRB(
                AppSizes.screenPaddingH,
                AppSizes.md,
                AppSizes.screenPaddingH,
                0,
              ),
              child: _HeroFeatureCard(
                totalCount: state.courses.length,
              ),
            ),
          ),

          // 3. Real Backend Announcement Banner (placed after hero card with see more)
          if (!dismissedAnnouncement)
            SliverToBoxAdapter(
              child: announcementsAsync.when(
                data: (announcements) {
                  if (announcements.isEmpty) return const SizedBox.shrink();
                  final announcement = announcements.first;
                  return Padding(
                    padding: const EdgeInsets.fromLTRB(
                      AppSizes.screenPaddingH,
                      AppSizes.md,
                      AppSizes.screenPaddingH,
                      0,
                    ),
                    child: _AnnouncementBanner(
                      title: announcement.title,
                      body: announcement.body,
                      onDismiss: onDismissAnnouncement,
                    ),
                  );
                },
                loading: () => const SizedBox.shrink(),
                error: (_, __) => const SizedBox.shrink(),
              ),
            ),

          // 4. Grade-level filter chips & Subject chooser dropdown directly underneath
          SliverToBoxAdapter(
            child: Padding(
              padding: const EdgeInsets.only(top: AppSizes.md),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  _GradeFilterRow(
                    selectedGrade: selectedGrade,
                    onSelect: onSelectGrade,
                  ),
                  const SizedBox(height: AppSizes.sm),
                  _SubjectDropdownSelector(
                    selectedSubject: state.selectedSubject,
                    onSelect: onSelectSubject,
                  ),
                ],
              ),
            ),
          ),

          // 5. Section Header: Courses
          SliverToBoxAdapter(
            child: Padding(
              padding: const EdgeInsets.fromLTRB(
                AppSizes.screenPaddingH,
                AppSizes.lg,
                AppSizes.screenPaddingH,
                AppSizes.sm,
              ),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  const Text(
                    'Courses',
                    style: TextStyle(
                      fontSize: 18,
                      fontWeight: FontWeight.bold,
                      color: AppColors.textPrimary,
                      letterSpacing: -0.3,
                    ),
                  ),
                  if (selectedGrade != _kFreshmanGrade &&
                      state.filteredCourses.isNotEmpty)
                    Text(
                      '${state.filteredCourses.length} available',
                      style: const TextStyle(
                        fontSize: 12,
                        color: AppColors.textMuted,
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                ],
              ),
            ),
          ),

          // 6. Course Content List or Empty State
          if (selectedGrade == _kFreshmanGrade)
            SliverFillRemaining(
              hasScrollBody: false,
              child: CourseEmptyState(
                icon: Icons.auto_awesome_rounded,
                title: 'Freshman courses are coming soon',
                body:
                    'We’re preparing freshman-year lessons and practice. For '
                    'now, explore our Grade 9–12 courses.',
                buttonLabel: 'Browse Grade 12',
                onPressed: () => onSelectGrade(12),
              ),
            )
          else if (state.filteredCourses.isEmpty)
            SliverFillRemaining(
              hasScrollBody: false,
              child: CourseEmptyState(
                icon: Icons.search_off_rounded,
                title: 'No courses found',
                body: 'Try another grade or subject.',
                buttonLabel: state.hasActiveFilters ? 'Clear filters' : null,
                onPressed: state.hasActiveFilters ? onResetFilters : null,
              ),
            )
          else
            SliverPadding(
              padding: const EdgeInsets.fromLTRB(
                AppSizes.screenPaddingH,
                0,
                AppSizes.screenPaddingH,
                AppSizes.lg,
              ),
              sliver: SliverList.separated(
                itemCount: state.filteredCourses.length,
                separatorBuilder: (_, __) =>
                    const SizedBox(height: AppSizes.md),
                itemBuilder: (context, index) {
                  final course = state.filteredCourses[index];
                  return CourseRowCard(course: course);
                },
              ),
            ),

          SliverToBoxAdapter(
            child: _LoadMoreFooter(state: state),
          ),
        ],
      ),
    );
  }
}

/// Minimalist App Bar: 'Memere' on left, Profile avatar on right
class _HomeTopBar extends ConsumerWidget {
  const _HomeTopBar();

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final user = ref.watch(authStateProvider).valueOrNull?.user;
    final firstName = user?.firstName.trim();
    final initial = firstName == null || firstName.isEmpty
        ? 'M'
        : firstName.substring(0, 1).toUpperCase();

    return Padding(
      padding: const EdgeInsets.fromLTRB(
        AppSizes.screenPaddingH,
        AppSizes.md,
        AppSizes.screenPaddingH,
        0,
      ),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          // Memere Brand Name
          const Text(
            'Memere',
            style: TextStyle(
              fontSize: 24,
              fontWeight: FontWeight.w900,
              color: AppColors.brandEmerald,
              letterSpacing: -0.6,
            ),
          ),

          // Clickable Profile Avatar (navigates to Account page)
          InkWell(
            onTap: () => context.go(AppRoutes.profile),
            borderRadius: BorderRadius.circular(20),
            child: Container(
              width: 38,
              height: 38,
              alignment: Alignment.center,
              decoration: BoxDecoration(
                color: AppColors.brandEmerald,
                shape: BoxShape.circle,
                boxShadow: [
                  BoxShadow(
                    color: AppColors.brandEmerald.withAlpha(80),
                    blurRadius: 8,
                    offset: const Offset(0, 2),
                  ),
                ],
              ),
              child: Text(
                initial,
                style: const TextStyle(
                  fontSize: 15,
                  fontWeight: FontWeight.bold,
                  color: Colors.white,
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }
}

/// Overflow-Safe Hero Featured Banner Card with Character Mascot
class _HeroFeatureCard extends StatelessWidget {
  const _HeroFeatureCard({required this.totalCount});

  final int totalCount;

  @override
  Widget build(BuildContext context) {
    return Container(
      width: double.infinity,
      decoration: BoxDecoration(
        gradient: const LinearGradient(
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
          colors: [
            Color(0xFF1D4ED8), // Royal Blue
            Color(0xFF2563EB), // Vibrant Blue
            Color(0xFF1E40AF), // Deep Blue
          ],
        ),
        borderRadius: BorderRadius.circular(22),
        boxShadow: [
          BoxShadow(
            color: const Color(0xFF1D4ED8).withAlpha(90),
            blurRadius: 14,
            offset: const Offset(0, 5),
          ),
        ],
      ),
      child: ClipRRect(
        borderRadius: BorderRadius.circular(22),
        child: Stack(
          children: [
            // Decorative background ambient rings
            Positioned(
              right: -25,
              top: -25,
              child: Container(
                width: 120,
                height: 120,
                decoration: BoxDecoration(
                  shape: BoxShape.circle,
                  color: Colors.white.withAlpha(15),
                ),
              ),
            ),

            // Banner Content & Animated Mascot
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 16),
              child: Row(
                crossAxisAlignment: CrossAxisAlignment.center,
                children: [
                  Expanded(
                    child: Column(
                      mainAxisSize: MainAxisSize.min,
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        // Tag Badge
                        Container(
                          padding: const EdgeInsets.symmetric(
                              horizontal: 8, vertical: 3.5),
                          decoration: BoxDecoration(
                            color: Colors.black.withAlpha(45),
                            borderRadius: BorderRadius.circular(16),
                            border:
                                Border.all(color: Colors.white.withAlpha(35)),
                          ),
                          child: const Row(
                            mainAxisSize: MainAxisSize.min,
                            children: [
                              Icon(Icons.stars_rounded,
                                  size: 12, color: Color(0xFF4ADE80)),
                              SizedBox(width: 4),
                              Text(
                                'Featured Courses',
                                style: TextStyle(
                                  fontSize: 10,
                                  fontWeight: FontWeight.w700,
                                  color: Colors.white,
                                  letterSpacing: 0.2,
                                ),
                              ),
                            ],
                          ),
                        ),
                        const SizedBox(height: 8),

                        // Title
                        const Text(
                          'Master Your Courses',
                          style: TextStyle(
                            fontSize: 18,
                            fontWeight: FontWeight.w900,
                            color: Colors.white,
                            height: 1.15,
                            letterSpacing: -0.3,
                          ),
                        ),
                        const SizedBox(height: 4),

                        // Subtitle
                        Text(
                          'Video lessons, short notes, and practice quizzes across grades.',
                          style: TextStyle(
                            fontSize: 11.5,
                            color: Colors.white.withAlpha(200),
                            height: 1.3,
                          ),
                        ),
                        const SizedBox(height: 10),

                        // Action Button / Pill
                        InkWell(
                          onTap: () => context.go(AppRoutes.mockExams),
                          borderRadius: BorderRadius.circular(16),
                          child: Container(
                            padding: const EdgeInsets.symmetric(
                                horizontal: 12, vertical: 6),
                            decoration: BoxDecoration(
                              color: Colors.white.withAlpha(35),
                              borderRadius: BorderRadius.circular(16),
                              border:
                                  Border.all(color: Colors.white.withAlpha(60)),
                            ),
                            child: const Row(
                              mainAxisSize: MainAxisSize.min,
                              children: [
                                Text(
                                  'Explore Practice Exams',
                                  style: TextStyle(
                                    fontSize: 11,
                                    fontWeight: FontWeight.bold,
                                    color: Colors.white,
                                  ),
                                ),
                                SizedBox(width: 4),
                                Icon(Icons.arrow_forward_rounded,
                                    size: 12, color: Colors.white),
                              ],
                            ),
                          ),
                        ),
                      ],
                    ),
                  ),
                  const SizedBox(width: 6),

                  // Character Animation from Saved Screen (Safely Constrained)
                  const SizedBox(
                    width: 80,
                    height: 86,
                    child: MemereMascot(
                      size: Size(80, 86),
                      showBackdrop: false,
                    ),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}

const Map<String, IconData> _subjectIcons = {
  'Mathematics': Icons.calculate_outlined,
  'Physics': Icons.science_outlined,
  'Chemistry': Icons.biotech_outlined,
  'Biology': Icons.eco_outlined,
  'English': Icons.menu_book_outlined,
  'Civics': Icons.gavel_outlined,
  'History': Icons.auto_stories_outlined,
  'Geography': Icons.public_outlined,
  'Economics': Icons.trending_up_outlined,
};

/// Horizontal grade-level filter chips (All [with 4-dot icon], Freshman, Grade 12 → Grade 9).
class _GradeFilterRow extends StatelessWidget {
  const _GradeFilterRow({
    required this.selectedGrade,
    required this.onSelect,
  });

  final int selectedGrade;
  final void Function(int grade) onSelect;

  @override
  Widget build(BuildContext context) {
    const grades = <(String, int, IconData?)>[
      ('All', _kAllGrade, Icons.grid_view_rounded),
      ('Freshman', _kFreshmanGrade, null),
      ('Grade 12', 12, null),
      ('Grade 11', 11, null),
      ('Grade 10', 10, null),
      ('Grade 9', 9, null),
    ];

    return SizedBox(
      height: 38,
      child: ListView.separated(
        scrollDirection: Axis.horizontal,
        padding: const EdgeInsets.symmetric(horizontal: AppSizes.screenPaddingH),
        itemCount: grades.length,
        separatorBuilder: (_, __) => const SizedBox(width: 8),
        itemBuilder: (context, index) {
          final (label, grade, icon) = grades[index];
          final isSelected = selectedGrade == grade;

          return InkWell(
            onTap: () => onSelect(grade),
            borderRadius: BorderRadius.circular(12),
            child: Container(
              padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 7),
              alignment: Alignment.center,
              decoration: BoxDecoration(
                color:
                    isSelected ? AppColors.brandEmerald : AppColors.bgSecondary,
                borderRadius: BorderRadius.circular(12),
                border: Border.all(
                  color: isSelected
                      ? AppColors.brandEmerald
                      : AppColors.borderStrong,
                ),
              ),
              child: Row(
                mainAxisSize: MainAxisSize.min,
                children: [
                  if (icon != null) ...[
                    Icon(
                      icon,
                      size: 14,
                      color: isSelected ? Colors.white : AppColors.textSecondary,
                    ),
                    const SizedBox(width: 5),
                  ],
                  Text(
                    label,
                    style: TextStyle(
                      fontSize: 12,
                      fontWeight: isSelected ? FontWeight.bold : FontWeight.w600,
                      color: isSelected ? Colors.white : AppColors.textSecondary,
                    ),
                  ),
                ],
              ),
            ),
          );
        },
      ),
    );
  }
}

/// Subject chooser dropdown button positioned directly UNDER the grade navigation bar.
class _SubjectDropdownSelector extends StatelessWidget {
  const _SubjectDropdownSelector({
    required this.selectedSubject,
    required this.onSelect,
  });

  final String? selectedSubject;
  final void Function(String? subject) onSelect;

  @override
  Widget build(BuildContext context) {
    final hasSubject = selectedSubject != null;
    final displaySubject = selectedSubject ?? 'All Subjects';
    final icon = hasSubject
        ? (_subjectIcons[selectedSubject] ?? Icons.book_outlined)
        : Icons.grid_view_rounded;

    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: AppSizes.screenPaddingH),
      child: InkWell(
        onTap: () => _openSubjectSheet(context),
        borderRadius: BorderRadius.circular(12),
        child: Container(
          height: 40,
          padding: const EdgeInsets.symmetric(horizontal: 14),
          decoration: BoxDecoration(
            color: hasSubject
                ? AppColors.brandEmerald.withAlpha(20)
                : AppColors.bgSecondary,
            borderRadius: BorderRadius.circular(12),
            border: Border.all(
              color: hasSubject
                  ? AppColors.brandEmerald
                  : AppColors.borderStrong,
              width: hasSubject ? 1.5 : 1.0,
            ),
          ),
          child: Row(
            children: [
              Icon(
                icon,
                size: 16,
                color: hasSubject
                    ? AppColors.brandEmerald
                    : AppColors.textSecondary,
              ),
              const SizedBox(width: 8),
              const Text(
                'Subject:',
                style: TextStyle(
                  fontSize: 12,
                  fontWeight: FontWeight.w600,
                  color: AppColors.textMuted,
                ),
              ),
              const SizedBox(width: 6),
              Expanded(
                child: Text(
                  displaySubject,
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                  style: TextStyle(
                    fontSize: 13,
                    fontWeight: FontWeight.bold,
                    color: hasSubject
                        ? AppColors.brandEmerald
                        : AppColors.textPrimary,
                  ),
                ),
              ),
              Icon(
                Icons.keyboard_arrow_down_rounded,
                size: 20,
                color: hasSubject
                    ? AppColors.brandEmerald
                    : AppColors.textSecondary,
              ),
            ],
          ),
        ),
      ),
    );
  }

  Future<void> _openSubjectSheet(BuildContext context) async {
    const options = <String?>[null, ...phase2Subjects];

    await showModalBottomSheet<void>(
      context: context,
      backgroundColor: AppColors.bgSecondary,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (sheetContext) {
        return SafeArea(
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              const SizedBox(height: AppSizes.sm),
              Container(
                width: 38,
                height: 4,
                decoration: BoxDecoration(
                  color: AppColors.borderStrong,
                  borderRadius: BorderRadius.circular(2),
                ),
              ),
              const Padding(
                padding: EdgeInsets.fromLTRB(
                  AppSizes.screenPaddingH,
                  AppSizes.md,
                  AppSizes.screenPaddingH,
                  AppSizes.xs,
                ),
                child: Align(
                  alignment: Alignment.centerLeft,
                  child: Text(
                    'Filter by subject',
                    style: TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.bold,
                      color: AppColors.textPrimary,
                    ),
                  ),
                ),
              ),
              Flexible(
                child: ListView.builder(
                  shrinkWrap: true,
                  padding: const EdgeInsets.only(bottom: AppSizes.sm),
                  itemCount: options.length,
                  itemBuilder: (context, index) {
                    final subject = options[index];
                    final isAll = subject == null;
                    final label = isAll ? 'All Subjects' : subject;
                    final icon = isAll
                        ? Icons.grid_view_rounded
                        : (_subjectIcons[subject] ?? Icons.book_outlined);
                    final isSelected = isAll
                        ? selectedSubject == null
                        : selectedSubject == subject;

                    return ListTile(
                      leading: Icon(
                        icon,
                        size: 20,
                        color: isSelected
                            ? AppColors.brandEmerald
                            : AppColors.textSecondary,
                      ),
                      title: Text(
                        label,
                        style: TextStyle(
                          fontSize: 14,
                          color: isSelected
                              ? AppColors.brandEmerald
                              : AppColors.textPrimary,
                          fontWeight:
                              isSelected ? FontWeight.bold : FontWeight.w500,
                        ),
                      ),
                      trailing: isSelected
                          ? const Icon(
                              Icons.check_rounded,
                              color: AppColors.brandEmerald,
                              size: 20,
                            )
                          : null,
                      onTap: () {
                        Navigator.of(sheetContext).pop();
                        onSelect(subject);
                      },
                    );
                  },
                ),
              ),
            ],
          ),
        );
      },
    );
  }
}

/// Announcement Banner Card with collapsible "See more" support
class _AnnouncementBanner extends StatefulWidget {
  const _AnnouncementBanner({
    required this.title,
    required this.body,
    required this.onDismiss,
  });

  final String title;
  final String body;
  final VoidCallback onDismiss;

  @override
  State<_AnnouncementBanner> createState() => _AnnouncementBannerState();
}

class _AnnouncementBannerState extends State<_AnnouncementBanner> {
  bool _isExpanded = false;

  @override
  Widget build(BuildContext context) {
    final isLong = widget.body.length > 110;

    return Container(
      width: double.infinity,
      padding: const EdgeInsets.all(14),
      decoration: BoxDecoration(
        color: AppColors.bgSecondary,
        borderRadius: BorderRadius.circular(16),
        border: Border.all(color: AppColors.borderStrong),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              Container(
                padding: const EdgeInsets.all(4),
                decoration: BoxDecoration(
                  color: AppColors.bgTertiary,
                  borderRadius: BorderRadius.circular(6),
                ),
                child: const Icon(
                  Icons.campaign_rounded,
                  color: AppColors.brandAmber,
                  size: 15,
                ),
              ),
              const SizedBox(width: 8),
              const Expanded(
                child: Text(
                  'ANNOUNCEMENT',
                  style: TextStyle(
                    fontSize: 10.5,
                    fontWeight: FontWeight.w800,
                    color: AppColors.textMuted,
                    letterSpacing: 0.8,
                  ),
                ),
              ),
              IconButton(
                onPressed: widget.onDismiss,
                icon: const Icon(Icons.close_rounded,
                    size: 16, color: AppColors.textMuted),
                padding: EdgeInsets.zero,
                constraints: const BoxConstraints(),
              ),
            ],
          ),
          const SizedBox(height: 6),
          Text(
            widget.title,
            style: const TextStyle(
              fontSize: 14,
              fontWeight: FontWeight.bold,
              color: AppColors.textPrimary,
            ),
          ),
          const SizedBox(height: 4),
          AnimatedSize(
            duration: const Duration(milliseconds: 200),
            curve: Curves.easeInOut,
            alignment: Alignment.topLeft,
            child: Text(
              widget.body,
              maxLines: _isExpanded ? null : 3,
              overflow: _isExpanded ? TextOverflow.visible : TextOverflow.ellipsis,
              style: const TextStyle(
                fontSize: 12.5,
                color: AppColors.textSecondary,
                height: 1.4,
              ),
            ),
          ),
          if (isLong) ...[
            const SizedBox(height: 6),
            InkWell(
              onTap: () => setState(() => _isExpanded = !_isExpanded),
              borderRadius: BorderRadius.circular(6),
              child: Padding(
                padding: const EdgeInsets.symmetric(vertical: 2),
                child: Row(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    Text(
                      _isExpanded ? 'See less' : 'See more',
                      style: const TextStyle(
                        fontSize: 12,
                        fontWeight: FontWeight.bold,
                        color: AppColors.brandEmerald,
                      ),
                    ),
                    const SizedBox(width: 2),
                    Icon(
                      _isExpanded
                          ? Icons.keyboard_arrow_up_rounded
                          : Icons.keyboard_arrow_down_rounded,
                      size: 16,
                      color: AppColors.brandEmerald,
                    ),
                  ],
                ),
              ),
            ),
          ],
        ],
      ),
    );
  }
}

class _LoadMoreFooter extends StatelessWidget {
  const _LoadMoreFooter({required this.state});

  final CourseListState state;

  @override
  Widget build(BuildContext context) {
    if (state.isLoadingMore) {
      return const Padding(
        padding: EdgeInsets.only(bottom: AppSizes.xl),
        child: Center(
          child: CircularProgressIndicator(
            strokeWidth: 2,
            color: AppColors.brandEmerald,
          ),
        ),
      );
    }

    if (state.loadMoreError != null) {
      return Padding(
        padding: const EdgeInsets.fromLTRB(
          AppSizes.screenPaddingH,
          0,
          AppSizes.screenPaddingH,
          AppSizes.xl,
        ),
        child: Text(
          state.loadMoreError!.message,
          textAlign: TextAlign.center,
          style: AppTextStyles.bodySmall.copyWith(color: AppColors.error),
        ),
      );
    }

    return const SizedBox(height: AppSizes.xl);
  }
}

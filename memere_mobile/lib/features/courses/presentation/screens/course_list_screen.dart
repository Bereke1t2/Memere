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
import '../../domain/entities/course_entity.dart';
import '../providers/course_list_provider.dart';
import '../widgets/course_card.dart';
import '../widgets/course_empty_state.dart';
import '../widgets/course_list_skeleton.dart';
import '../widgets/subject_filter_chips.dart';

/// Professional Duolingo Obsidian Dark Home Screen for Memere.
///
/// Features animated subject boxes, animated daily goal progress, vibrant Course Level Badges
/// (Beginner, Intermediate, Advanced), and restrained color usage so green is not everywhere.
class CourseListScreen extends ConsumerStatefulWidget {
  const CourseListScreen({super.key});

  @override
  ConsumerState<CourseListScreen> createState() => _CourseListScreenState();
}

class _CourseListScreenState extends ConsumerState<CourseListScreen> {
  late final TextEditingController _searchController;
  late final ScrollController _scrollController;
  bool _dismissedAnnouncement = false;

  @override
  void initState() {
    super.initState();
    _searchController = TextEditingController();
    _scrollController = ScrollController()..addListener(_onScroll);
  }

  @override
  void dispose() {
    _scrollController
      ..removeListener(_onScroll)
      ..dispose();
    _searchController.dispose();
    super.dispose();
  }

  void _onScroll() {
    if (!_scrollController.hasClients) return;
    final position = _scrollController.position;
    if (position.pixels >= position.maxScrollExtent - 360) {
      ref.read(courseListProvider.notifier).loadMore();
    }
  }

  @override
  Widget build(BuildContext context) {
    final coursesAsync = ref.watch(courseListProvider);

    ref.listen(courseListProvider, (_, next) {
      final query = next.valueOrNull?.searchQuery ?? '';
      if (_searchController.text != query) {
        _searchController.value = TextEditingValue(
          text: query,
          selection: TextSelection.collapsed(offset: query.length),
        );
      }
    });

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
                searchController: _searchController,
                state: previous,
                dismissedAnnouncement: _dismissedAnnouncement,
                onDismissAnnouncement: () =>
                    setState(() => _dismissedAnnouncement = true),
              );
            }
            return Column(
              children: [
                _DashboardHeader(searchController: _searchController),
                const SizedBox(height: AppSizes.lg),
                const Expanded(child: CourseListSkeleton()),
              ],
            );
          },
          error: (error, _) => Column(
            children: [
              _DashboardHeader(searchController: _searchController),
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
            searchController: _searchController,
            state: state,
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
    required this.searchController,
    required this.state,
    required this.dismissedAnnouncement,
    required this.onDismissAnnouncement,
  });

  final ScrollController controller;
  final TextEditingController searchController;
  final CourseListState state;
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
          // 1. Header (Avatar + Streak + XP)
          SliverToBoxAdapter(
            child: _DashboardHeader(searchController: searchController),
          ),

          // 2. Backend Announcement Banner (if available & not dismissed)
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

          // 3. Animated Daily Target Banner
          const SliverToBoxAdapter(
            child: Padding(
              padding: EdgeInsets.fromLTRB(
                AppSizes.screenPaddingH,
                AppSizes.md,
                AppSizes.screenPaddingH,
                0,
              ),
              child: _DailyGoalCard(),
            ),
          ),

          // 4. Subject Filter Chips
          SliverToBoxAdapter(
            child: _DashboardFilters(state: state),
          ),

          // 5. Content List or Empty State
          if (state.filteredCourses.isEmpty)
            SliverFillRemaining(
              hasScrollBody: false,
              child: CourseEmptyState(
                icon: Icons.search_off_rounded,
                title: 'No courses found',
                body: 'Try another subject or search term.',
                buttonLabel: state.hasActiveFilters ? 'Clear filters' : null,
                onPressed: state.hasActiveFilters
                    ? () => ref.read(courseListProvider.notifier).clearFilters()
                    : null,
              ),
            )
          else ...[
            // Animated Subject Grid
            SliverToBoxAdapter(
              child: _TopicGrid(courses: state.filteredCourses),
            ),

            // Section Header: All Courses
            const SliverToBoxAdapter(
              child: Padding(
                padding: EdgeInsets.fromLTRB(
                  AppSizes.screenPaddingH,
                  AppSizes.lg,
                  AppSizes.screenPaddingH,
                  AppSizes.md,
                ),
                child: Text(
                  'Exam Prep Courses',
                  style: TextStyle(
                    fontSize: 18,
                    fontWeight: FontWeight.bold,
                    color: AppColors.textPrimary,
                    letterSpacing: -0.2,
                  ),
                ),
              ),
            ),

            // Course Rows List
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
          ],
          SliverToBoxAdapter(
            child: _LoadMoreFooter(state: state),
          ),
        ],
      ),
    );
  }
}

/// Duolingo Header Top Bar with Profile Avatar, Greeting, Streak 🔥, and XP ⚡ Badges
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
                                'National Entrance Prep',
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
                          'Master Your Exams',
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
                          'Grade 12 video lessons, short notes, and practice exams.',
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
                                  'Explore Mock Exams',
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

class _DashboardHeader extends ConsumerWidget {
  const _DashboardHeader({required this.searchController});

  final TextEditingController searchController;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final user = ref.watch(authStateProvider).valueOrNull?.user;
    final firstName = user?.firstName.trim();
    final greeting = firstName == null || firstName.isEmpty
        ? 'Selam, Student!'
        : 'Selam, $firstName!';
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
      child: Column(
        children: [
          Row(
            children: [
              // Avatar
              Container(
                width: 44,
                height: 44,
                alignment: Alignment.center,
                decoration: BoxDecoration(
                  color: AppColors.bgSecondary,
                  shape: BoxShape.circle,
                  border: Border.all(color: AppColors.borderStrong, width: 2),
                ),
                child: Text(
                  initial,
                  style: const TextStyle(
                    fontSize: 18,
                    fontWeight: FontWeight.bold,
                    color: AppColors.textPrimary,
                  ),
                ),
              ),
              const SizedBox(width: 12),

              // Greeting & Subtitle
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      greeting,
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                      style: const TextStyle(
                        fontSize: 17,
                        fontWeight: FontWeight.bold,
                        color: AppColors.textPrimary,
                      ),
                    ),
                    const SizedBox(height: 2),
                    const Text(
                      'Grade 12 National Exam Prep',
                      style: TextStyle(
                        fontSize: 12,
                        color: AppColors.textMuted,
                        fontWeight: FontWeight.w500,
                      ),
                    ),
                  ],
                ),
              ),

              // Duolingo Streak Badge (🔥 7 Days)
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 6),
                decoration: BoxDecoration(
                  color: AppColors.bgSecondary,
                  borderRadius: BorderRadius.circular(16),
                  border: Border.all(color: AppColors.borderStrong),
                ),
                child: const Row(
                  children: [
                    Icon(Icons.local_fire_department_rounded,
                        color: AppColors.brandAmber, size: 18),
                    SizedBox(width: 4),
                    Text(
                      '7',
                      style: TextStyle(
                        fontSize: 13,
                        fontWeight: FontWeight.bold,
                        color: AppColors.brandAmber,
                      ),
                    ),
                  ],
                ),
              ),
              const SizedBox(width: 8),

              // Duolingo XP Points Badge (⚡ 450 XP)
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 6),
                decoration: BoxDecoration(
                  color: AppColors.bgSecondary,
                  borderRadius: BorderRadius.circular(16),
                  border: Border.all(color: AppColors.borderStrong),
                ),
                child: const Row(
                  children: [
                    Icon(Icons.bolt_rounded, color: AppColors.brandAmber, size: 18),
                    SizedBox(width: 4),
                    Text(
                      '450',
                      style: TextStyle(
                        fontSize: 13,
                        fontWeight: FontWeight.bold,
                        color: AppColors.brandAmber,
                      ),
                    ),
                  ],
                ),
              ),
            ],
          ),
          const SizedBox(height: 14),

          // Floating Search Bar
          _FloatingSearchBar(searchController: searchController),
        ],
      ),
    );
  }
}

/// Backend Announcement Banner Card (Duolingo Tactile 3D Styling)
class _AnnouncementBanner extends StatelessWidget {
  const _AnnouncementBanner({
    required this.title,
    required this.body,
    required this.onDismiss,
  });

  final String title;
  final String body;
  final VoidCallback onDismiss;

  @override
  Widget build(BuildContext context) {
    return Container(
      width: double.infinity,
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: AppColors.bgSecondary,
        borderRadius: BorderRadius.circular(20),
        border: const Border(
          top: BorderSide(color: AppColors.borderStrong),
          left: BorderSide(color: AppColors.borderStrong),
          right: BorderSide(color: AppColors.borderStrong),
          bottom: BorderSide(color: AppColors.brandAmberDark, width: 4),
        ),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              Container(
                padding: const EdgeInsets.all(6),
                decoration: BoxDecoration(
                  color: AppColors.brandAmber.withAlpha(38),
                  borderRadius: BorderRadius.circular(8),
                ),
                child: const Icon(
                  Icons.campaign_rounded,
                  color: AppColors.brandAmber,
                  size: 18,
                ),
              ),
              const SizedBox(width: 10),
              const Expanded(
                child: Text(
                  'ANNOUNCEMENT',
                  style: TextStyle(
                    fontSize: 11,
                    fontWeight: FontWeight.w800,
                    color: AppColors.brandAmber,
                    letterSpacing: 1.0,
                  ),
                ),
              ),
              IconButton(
                onPressed: onDismiss,
                icon: const Icon(Icons.close_rounded, size: 18, color: AppColors.textMuted),
                padding: EdgeInsets.zero,
                constraints: const BoxConstraints(),
              ),
            ],
          ),
          const SizedBox(height: 8),
          Text(
            title,
            style: const TextStyle(
              fontSize: 15,
              fontWeight: FontWeight.bold,
              color: AppColors.textPrimary,
            ),
          ),
          const SizedBox(height: 4),
          Text(
            body,
            style: const TextStyle(
              fontSize: 13,
              color: AppColors.textSecondary,
              height: 1.4,
            ),
          ),
          const SizedBox(height: 12),
          SizedBox(
            width: double.infinity,
            height: 42,
            child: ElevatedButton.icon(
              onPressed: () => context.go(AppRoutes.mockExams),
              style: ElevatedButton.styleFrom(
                backgroundColor: AppColors.brandAmber,
                foregroundColor: Colors.black,
                elevation: 0,
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(12),
                  side: const BorderSide(color: AppColors.brandAmberDark, width: 2),
                ),
              ),
              icon: const Icon(Icons.arrow_forward_rounded, size: 18, color: Colors.black),
              label: const Text(
                'OPEN MOCK EXAMS',
                style: TextStyle(
                  fontSize: 13,
                  fontWeight: FontWeight.w800,
                  letterSpacing: 0.5,
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }
}

/// Duolingo Signature Daily Target Banner Card with Smooth Fill Animation
class _DailyGoalCard extends StatelessWidget {
  const _DailyGoalCard();

  @override
  Widget build(BuildContext context) {
    return Container(
      width: double.infinity,
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: AppColors.bgSecondary,
        borderRadius: BorderRadius.circular(20),
        border: const Border(
          top: BorderSide(color: AppColors.borderStrong),
          left: BorderSide(color: AppColors.borderStrong),
          right: BorderSide(color: AppColors.borderStrong),
          bottom: BorderSide(color: AppColors.brandEmeraldDark, width: 4),
        ),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Row(
                children: [
                  Container(
                    padding: const EdgeInsets.all(8),
                    decoration: BoxDecoration(
                      color: AppColors.bgTertiary,
                      borderRadius: BorderRadius.circular(10),
                      border: Border.all(color: AppColors.borderStrong),
                    ),
                    child: const Icon(
                      Icons.track_changes_rounded,
                      color: AppColors.textPrimary,
                      size: 20,
                    ),
                  ),
                  const SizedBox(width: 10),
                  const Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        'DAILY TARGET',
                        style: TextStyle(
                          fontSize: 11,
                          fontWeight: FontWeight.w800,
                          color: AppColors.brandEmerald,
                          letterSpacing: 1.0,
                        ),
                      ),
                      Text(
                        '3 of 5 lessons completed',
                        style: TextStyle(
                          fontSize: 14,
                          fontWeight: FontWeight.bold,
                          color: AppColors.textPrimary,
                        ),
                      ),
                    ],
                  ),
                ],
              ),
              const Text(
                '60%',
                style: TextStyle(
                  fontSize: 16,
                  fontWeight: FontWeight.bold,
                  color: AppColors.brandEmerald,
                ),
              ),
            ],
          ),
          const SizedBox(height: 12),

          // Animated Smooth Progress Bar
          TweenAnimationBuilder<double>(
            tween: Tween<double>(begin: 0.0, end: 0.6),
            duration: const Duration(milliseconds: 1000),
            curve: Curves.easeOutCubic,
            builder: (context, value, child) {
              return ClipRRect(
                borderRadius: BorderRadius.circular(10),
                child: LinearProgressIndicator(
                  value: value,
                  minHeight: 12,
                  backgroundColor: AppColors.bgPrimary,
                  color: AppColors.brandEmerald,
                ),
              );
            },
          ),
          const SizedBox(height: 14),

          // Tactile 3D Action Button
          SizedBox(
            width: double.infinity,
            height: 46,
            child: ElevatedButton.icon(
              onPressed: () => context.go(AppRoutes.learn),
              style: ElevatedButton.styleFrom(
                backgroundColor: AppColors.brandEmerald,
                foregroundColor: Colors.white,
                elevation: 0,
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(14),
                  side: const BorderSide(color: AppColors.brandEmeraldDark, width: 2),
                ),
              ),
              icon: const Icon(Icons.play_arrow_rounded, size: 20),
              label: const Text(
                'CONTINUE STUDYING',
                style: TextStyle(
                  fontSize: 14,
                  fontWeight: FontWeight.w800,
                  letterSpacing: 0.5,
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }
}

/// Neutral Slate Search Bar Container
class _FloatingSearchBar extends ConsumerWidget {
  const _FloatingSearchBar({required this.searchController});

  final TextEditingController searchController;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return Container(
      height: 48,
      decoration: BoxDecoration(
        color: AppColors.bgSecondary,
        borderRadius: BorderRadius.circular(16),
        border: Border.all(color: AppColors.borderStrong),
      ),
      child: Row(
        children: [
          const SizedBox(width: 14),
          const Icon(
            Icons.search_rounded,
            size: 20,
            color: AppColors.textMuted,
          ),
          const SizedBox(width: 10),
          Expanded(
            child: TextField(
              controller: searchController,
              textInputAction: TextInputAction.search,
              onChanged: ref.read(courseListProvider.notifier).setSearchQuery,
              style: const TextStyle(fontSize: 14, color: AppColors.textPrimary),
              decoration: const InputDecoration(
                hintText: 'Search subjects, topics, exams...',
                hintStyle: TextStyle(fontSize: 14, color: AppColors.textMuted),
                border: InputBorder.none,
                enabledBorder: InputBorder.none,
                focusedBorder: InputBorder.none,
                isCollapsed: true,
              ),
            ),
          ),
          IconButton(
            tooltip: 'Filter',
            onPressed: () {},
            icon: const Icon(
              Icons.tune_rounded,
              size: 20,
              color: AppColors.textSecondary,
            ),
          ),
          const SizedBox(width: 4),
        ],
      ),
    );
  }
}

class _DashboardFilters extends ConsumerWidget {
  const _DashboardFilters({required this.state});

  final CourseListState state;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const SizedBox(height: 14),
        SubjectFilterChips(
          subjects: phase2Subjects,
          selectedSubject: state.selectedSubject,
          onSelected: (subject) {
            ref.read(courseListProvider.notifier).setSubject(subject);
          },
        ),
      ],
    );
  }
}

/// Animated Subject Grid with Interactive Press Scale & Character Icons
class _TopicGrid extends StatelessWidget {
  const _TopicGrid({required this.courses});

  final List<CourseEntity> courses;

  @override
  Widget build(BuildContext context) {
    if (courses.isEmpty) return const SizedBox.shrink();
    final subjects = <String>[];
    for (final course in courses) {
      if (!subjects.contains(course.subject)) subjects.add(course.subject);
      if (subjects.length == 6) break;
    }

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Padding(
          padding: EdgeInsets.fromLTRB(
            AppSizes.screenPaddingH,
            AppSizes.lg,
            AppSizes.screenPaddingH,
            AppSizes.sm,
          ),
          child: Text(
            'Subject Lanes',
            style: TextStyle(
              fontSize: 18,
              fontWeight: FontWeight.bold,
              color: AppColors.textPrimary,
              letterSpacing: -0.2,
            ),
          ),
        ),
        Padding(
          padding:
              const EdgeInsets.symmetric(horizontal: AppSizes.screenPaddingH),
          child: GridView.builder(
            itemCount: subjects.length,
            shrinkWrap: true,
            physics: const NeverScrollableScrollPhysics(),
            gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
              crossAxisCount: 2,
              mainAxisSpacing: 10,
              crossAxisSpacing: 10,
              childAspectRatio: 2.5,
            ),
            itemBuilder: (context, index) {
              final subject = subjects[index];
              return _AnimatedSubjectCard(subject: subject);
            },
          ),
        ),
      ],
    );
  }
}

/// Interactive Micro-Animated Subject Card
class _AnimatedSubjectCard extends StatefulWidget {
  const _AnimatedSubjectCard({required this.subject});

  final String subject;

  @override
  State<_AnimatedSubjectCard> createState() => _AnimatedSubjectCardState();
}

class _AnimatedSubjectCardState extends State<_AnimatedSubjectCard> {
  bool _isPressed = false;

  @override
  Widget build(BuildContext context) {
    final subject = widget.subject;

    return GestureDetector(
      onTapDown: (_) => setState(() => _isPressed = true),
      onTapUp: (_) => setState(() => _isPressed = false),
      onTapCancel: () => setState(() => _isPressed = false),
      onTap: () {},
      child: AnimatedScale(
        scale: _isPressed ? 0.96 : 1.0,
        duration: const Duration(milliseconds: 120),
        curve: Curves.easeOutCubic,
        child: Container(
          padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
          decoration: BoxDecoration(
            color: AppColors.bgSecondary,
            borderRadius: BorderRadius.circular(16),
            border: Border(
              top: const BorderSide(color: AppColors.borderStrong),
              left: const BorderSide(color: AppColors.borderStrong),
              right: const BorderSide(color: AppColors.borderStrong),
              bottom: BorderSide(
                color: _isPressed ? AppColors.borderFocused : AppColors.borderStrong,
                width: 3,
              ),
            ),
          ),
          child: Row(
            children: [
              // Muted Character Icon Tile
              Container(
                padding: const EdgeInsets.all(8),
                decoration: BoxDecoration(
                  color: AppColors.bgTertiary,
                  shape: BoxShape.circle,
                  border: Border.all(color: AppColors.borderStrong),
                ),
                child: Icon(
                  _subjectIcon(subject),
                  size: 18,
                  color: AppColors.textPrimary,
                ),
              ),
              const SizedBox(width: 10),
              Expanded(
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      _shortSubject(subject),
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                      style: const TextStyle(
                        fontSize: 14,
                        fontWeight: FontWeight.bold,
                        color: AppColors.textPrimary,
                      ),
                    ),
                    const SizedBox(height: 2),
                    const Text(
                      'Prep Lane',
                      style: TextStyle(
                        fontSize: 11,
                        color: AppColors.textMuted,
                      ),
                    ),
                  ],
                ),
              ),
              const Icon(
                Icons.chevron_right_rounded,
                size: 18,
                color: AppColors.textMuted,
              ),
            ],
          ),
        ),
      ),
    );
  }
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

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/router/app_router.dart';
import '../../../../core/storage/hive/models/saved_item.dart';
import '../../../auth/presentation/providers/auth_state_provider.dart';
import '../../../courses/domain/entities/course_entity.dart';
import '../../../courses/presentation/providers/course_download_provider.dart';
import '../../../exam/presentation/widgets/exam_empty_state.dart';
import '../../../payment/domain/entities/enrollment_entity.dart';
import '../../../payment/presentation/providers/purchase_history_provider.dart';
import '../../../saved/presentation/providers/saved_courses_provider.dart';

enum LearningFilter { all, enrolled, favorites, downloaded }

enum CourseCardKind { enrolled, favorite, downloaded }

class _CourseCardData {
  const _CourseCardData({
    required this.id,
    required this.kind,
    required this.title,
    required this.subtitle,
    this.enrollment,
    this.savedItem,
    this.courseEntity,
  });

  final String id;
  final CourseCardKind kind;
  final String title;
  final String subtitle;
  final EnrollmentEntity? enrollment;
  final SavedItem? savedItem;
  final CourseEntity? courseEntity;
}

/// Refined Obsidian & Soft Emerald My Learning Screen matching Home & Exam catalog UX.
class MyLearningScreen extends ConsumerStatefulWidget {
  const MyLearningScreen({super.key});

  @override
  ConsumerState<MyLearningScreen> createState() => _MyLearningScreenState();
}

class _MyLearningScreenState extends ConsumerState<MyLearningScreen> {
  LearningFilter _selectedFilter = LearningFilter.all;

  @override
  Widget build(BuildContext context) {
    final isAuthenticated =
        ref.watch(authStateProvider).valueOrNull?.isAuthenticated ?? false;
    final favorites = ref.watch(savedCoursesProvider).valueOrNull ?? const [];
    final downloadedCourses = ref.watch(downloadedCoursesProvider);
    final enrollmentAsync =
        isAuthenticated ? ref.watch(enrollmentListProvider) : null;
    final enrollments = enrollmentAsync?.valueOrNull ?? const [];

    final cardItems = <_CourseCardData>[];
    final seenIds = <String>{};

    // 1. Enrolled Courses (account level)
    for (final e in enrollments) {
      if (seenIds.add(e.courseId)) {
        cardItems.add(_CourseCardData(
          id: e.courseId,
          kind: CourseCardKind.enrolled,
          title: e.courseTitle.isNotEmpty ? e.courseTitle : 'Course',
          subtitle: 'Grade ${e.courseGrade} • ${e.courseSubject}',
          enrollment: e,
        ));
      }
    }

    // 2. Downloaded Courses (offline access)
    for (final d in downloadedCourses) {
      if (seenIds.add(d.id)) {
        cardItems.add(_CourseCardData(
          id: d.id,
          kind: CourseCardKind.downloaded,
          title: d.title,
          subtitle: 'Grade ${d.grade} • ${d.subject}',
          courseEntity: d,
        ));
      }
    }

    // 3. Favorite Courses (local bookmarks)
    for (final f in favorites) {
      if (seenIds.add(f.courseId)) {
        cardItems.add(_CourseCardData(
          id: f.courseId,
          kind: CourseCardKind.favorite,
          title: f.title,
          subtitle: f.subtitle ?? 'Favorite Course',
          savedItem: f,
        ));
      }
    }

    final filteredItems = cardItems.where((item) {
      switch (_selectedFilter) {
        case LearningFilter.all:
          return true;
        case LearningFilter.enrolled:
          return item.kind == CourseCardKind.enrolled;
        case LearningFilter.favorites:
          return item.kind == CourseCardKind.favorite;
        case LearningFilter.downloaded:
          return item.kind == CourseCardKind.downloaded;
      }
    }).toList();

    final isLoading =
        (enrollmentAsync?.isLoading ?? false) && cardItems.isEmpty;

    return Scaffold(
      backgroundColor: AppColors.bgPrimary,
      body: SafeArea(
        child: RefreshIndicator(
          color: AppColors.brandEmerald,
          backgroundColor: AppColors.bgSecondary,
          onRefresh: () async {
            ref.invalidate(savedCoursesProvider);
            if (isAuthenticated) {
              ref.invalidate(enrollmentListProvider);
              await ref.read(enrollmentListProvider.future);
            }
          },
          child: CustomScrollView(
            physics: const AlwaysScrollableScrollPhysics(),
            slivers: [
              // 1. Top Header Row + Category Pills Filter
              SliverToBoxAdapter(
                child: _MyLearningHeader(
                  selectedFilter: _selectedFilter,
                  onSelectFilter: (f) => setState(() => _selectedFilter = f),
                  totalCount: cardItems.length,
                  enrolledCount: enrollments.length,
                  favoritesCount: favorites.length,
                  downloadedCount: downloadedCourses.length,
                ),
              ),

              // 2. Unauthenticated Banner if viewing Enrolled specifically as a Guest
              if (!isAuthenticated &&
                  _selectedFilter == LearningFilter.enrolled)
                const SliverToBoxAdapter(
                  child: Padding(
                    padding: EdgeInsets.symmetric(
                        horizontal: AppSizes.screenPaddingH, vertical: 8),
                    child: _GuestSignInPromptCard(),
                  ),
                ),

              // 3. Grid of Vibrant Obsidian Cards or Empty State
              if (isLoading)
                const SliverFillRemaining(
                  child: Center(
                    child: CircularProgressIndicator(
                        color: AppColors.brandEmerald),
                  ),
                )
              else if (filteredItems.isEmpty)
                SliverFillRemaining(
                  hasScrollBody: false,
                  child: ExamEmptyState(
                    icon: Icons.school_outlined,
                    title: cardItems.isEmpty
                        ? 'No courses in your learning hub'
                        : 'No ${_selectedFilter.name} courses found',
                    body: cardItems.isEmpty
                        ? 'Enroll in courses or save them to your favorites to track your progress here.'
                        : 'Try selecting another category above to view your courses.',
                    buttonLabel:
                        cardItems.isEmpty ? 'Browse Courses' : 'Show All',
                    onPressed: cardItems.isEmpty
                        ? () => context.go(AppRoutes.home)
                        : () => setState(
                            () => _selectedFilter = LearningFilter.all),
                  ),
                )
              else
                SliverPadding(
                  padding: const EdgeInsets.fromLTRB(
                    AppSizes.screenPaddingH,
                    0,
                    AppSizes.screenPaddingH,
                    AppSizes.md,
                  ),
                  sliver: SliverGrid(
                    gridDelegate:
                        const SliverGridDelegateWithFixedCrossAxisCount(
                      crossAxisCount: 2,
                      mainAxisSpacing: 10,
                      crossAxisSpacing: 10,
                      childAspectRatio: 0.72,
                    ),
                    delegate: SliverChildBuilderDelegate(
                      (context, index) {
                        final item = filteredItems[index];
                        return _CourseCardTile(
                          item: item,
                          cardIndex: index,
                          onOpen: () => context
                              .push(AppRoutes.courseDetailPath(item.id)),
                        );
                      },
                      childCount: filteredItems.length,
                    ),
                  ),
                ),
              const SliverToBoxAdapter(
                child: SizedBox(height: AppSizes.md),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

/// Header Row (Logo + Courses Badge) + Category Selector Badges
class _MyLearningHeader extends ConsumerWidget {
  const _MyLearningHeader({
    required this.selectedFilter,
    required this.onSelectFilter,
    required this.totalCount,
    required this.enrolledCount,
    required this.favoritesCount,
    required this.downloadedCount,
  });

  final LearningFilter selectedFilter;
  final ValueChanged<LearningFilter> onSelectFilter;
  final int totalCount;
  final int enrolledCount;
  final int favoritesCount;
  final int downloadedCount;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final filterOptions = [
      (LearningFilter.all, Icons.grid_view_rounded, 'All', totalCount),
      (LearningFilter.enrolled, Icons.verified_rounded, 'Enrolled', enrolledCount),
      (LearningFilter.favorites, Icons.bookmark_rounded, 'Favorites', favoritesCount),
      (LearningFilter.downloaded, Icons.download_done_rounded, 'Downloaded', downloadedCount),
    ];

    return Padding(
      padding: const EdgeInsets.fromLTRB(0, AppSizes.sm, 0, AppSizes.xs),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Top Bar: Branding + Enrolled / Learning Counter Badge
          Padding(
            padding: const EdgeInsets.symmetric(
                horizontal: AppSizes.screenPaddingH),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                const Row(
                  children: [
                    Icon(
                      Icons.school_rounded,
                      color: AppColors.brandEmerald,
                      size: 22,
                    ),
                    SizedBox(width: 8),
                    Text(
                      'Memere Learning Hub',
                      style: TextStyle(
                        fontSize: 18,
                        fontWeight: FontWeight.w900,
                        color: AppColors.textPrimary,
                        letterSpacing: -0.4,
                      ),
                    ),
                  ],
                ),

                // Course counter pill badge
                Container(
                  padding:
                      const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
                  decoration: BoxDecoration(
                    color: AppColors.bgSecondary,
                    borderRadius: BorderRadius.circular(16),
                    border: Border.all(color: AppColors.borderStrong),
                  ),
                  child: Row(
                    children: [
                      const Icon(Icons.local_library_rounded,
                          color: AppColors.brandEmerald, size: 16),
                      const SizedBox(width: 5),
                      Text(
                        '$totalCount ${totalCount == 1 ? 'course' : 'courses'}',
                        style: const TextStyle(
                          fontSize: 12.5,
                          fontWeight: FontWeight.bold,
                          color: AppColors.textPrimary,
                        ),
                      ),
                    ],
                  ),
                ),
              ],
            ),
          ),
          const SizedBox(height: 14),

          // Category Selector Pills Row
          SingleChildScrollView(
            scrollDirection: Axis.horizontal,
            physics: const BouncingScrollPhysics(),
            padding: const EdgeInsets.symmetric(
                horizontal: AppSizes.screenPaddingH),
            child: Row(
              children: List.generate(filterOptions.length, (idx) {
                final opt = filterOptions[idx];
                final isSelected = selectedFilter == opt.$1;

                return Padding(
                  padding: EdgeInsets.only(
                      right: idx == filterOptions.length - 1 ? 0 : 10),
                  child: GestureDetector(
                    onTap: () => onSelectFilter(opt.$1),
                    child: AnimatedContainer(
                      duration: const Duration(milliseconds: 180),
                      padding: const EdgeInsets.symmetric(
                          horizontal: 14, vertical: 8),
                      decoration: BoxDecoration(
                        color: isSelected
                            ? AppColors.brandEmerald
                            : AppColors.bgSecondary,
                        borderRadius: BorderRadius.circular(16),
                        border: Border.all(
                          color: isSelected
                              ? AppColors.brandEmerald
                              : AppColors.borderStrong,
                        ),
                        boxShadow: isSelected
                            ? [
                                BoxShadow(
                                  color: AppColors.brandEmerald.withAlpha(60),
                                  blurRadius: 8,
                                  offset: const Offset(0, 2),
                                ),
                              ]
                            : null,
                      ),
                      child: Row(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          Icon(
                            opt.$2,
                            size: 15,
                            color: isSelected
                                ? Colors.white
                                : AppColors.textSecondary,
                          ),
                          const SizedBox(width: 6),
                          Text(
                            '${opt.$3} (${opt.$4})',
                            style: TextStyle(
                              fontSize: 12,
                              fontWeight: isSelected
                                  ? FontWeight.bold
                                  : FontWeight.w600,
                              color: isSelected
                                  ? Colors.white
                                  : AppColors.textSecondary,
                            ),
                          ),
                        ],
                      ),
                    ),
                  ),
                );
              }),
            ),
          ),
          const SizedBox(height: 14),
        ],
      ),
    );
  }
}

class _CourseCardTile extends ConsumerStatefulWidget {
  const _CourseCardTile({
    required this.item,
    required this.cardIndex,
    required this.onOpen,
  });

  final _CourseCardData item;
  final int cardIndex;
  final VoidCallback onOpen;

  @override
  ConsumerState<_CourseCardTile> createState() => _CourseCardTileState();
}

class _CourseCardTileState extends ConsumerState<_CourseCardTile> {
  bool _isPressed = false;

  @override
  Widget build(BuildContext context) {
    final item = widget.item;
    final (gradient, glowColor, badgeColor, badgeLabel, iconData, buttonLabel) =
        _getKindStyle(item.kind);

    return GestureDetector(
      onTapDown: (_) => setState(() => _isPressed = true),
      onTapUp: (_) => setState(() => _isPressed = false),
      onTapCancel: () => setState(() => _isPressed = false),
      onTap: widget.onOpen,
      child: AnimatedScale(
        scale: _isPressed ? 0.96 : 1.0,
        duration: const Duration(milliseconds: 120),
        child: ClipRRect(
          borderRadius: BorderRadius.circular(16),
          child: Container(
            padding: const EdgeInsets.all(12),
            decoration: BoxDecoration(
              gradient: gradient,
              borderRadius: BorderRadius.circular(16),
              border: Border.all(
                color: _isPressed ? Colors.white : glowColor.withAlpha(90),
                width: 1.5,
              ),
              boxShadow: [
                BoxShadow(
                  color: glowColor.withAlpha(45),
                  blurRadius: 10,
                  offset: const Offset(0, 4),
                ),
              ],
            ),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                // Top Row: Category Pill Badge + Status Icon
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Container(
                      padding: const EdgeInsets.symmetric(
                          horizontal: 8, vertical: 3),
                      decoration: BoxDecoration(
                        color: badgeColor.withAlpha(60),
                        borderRadius: BorderRadius.circular(6),
                        border: Border.all(
                            color: Colors.white.withAlpha(50), width: 0.8),
                      ),
                      child: Text(
                        badgeLabel,
                        style: const TextStyle(
                          fontSize: 10.5,
                          fontWeight: FontWeight.bold,
                          color: Colors.white,
                        ),
                      ),
                    ),
                    Container(
                      padding: const EdgeInsets.all(4),
                      decoration: BoxDecoration(
                        color: Colors.black.withAlpha(40),
                        shape: BoxShape.circle,
                      ),
                      child: Icon(
                        iconData,
                        size: 14,
                        color: Colors.white70,
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 8),

                // Middle: Icon + Title + Subtitle
                Row(
                  children: [
                    Container(
                      width: 36,
                      height: 36,
                      decoration: BoxDecoration(
                        color: Colors.black.withAlpha(40),
                        shape: BoxShape.circle,
                        border: Border.all(color: Colors.white.withAlpha(40)),
                      ),
                      child: const Icon(Icons.school_rounded,
                          size: 18, color: Colors.white),
                    ),
                    const SizedBox(width: 8),
                    Expanded(
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(
                            item.title,
                            maxLines: 2,
                            overflow: TextOverflow.ellipsis,
                            style: const TextStyle(
                              fontSize: 13.5,
                              fontWeight: FontWeight.w800,
                              color: Colors.white,
                              height: 1.2,
                            ),
                          ),
                          const SizedBox(height: 2),
                          Text(
                            item.subtitle,
                            maxLines: 1,
                            overflow: TextOverflow.ellipsis,
                            style: TextStyle(
                              fontSize: 11,
                              color: Colors.white.withAlpha(180),
                            ),
                          ),
                        ],
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 10),

                // Bottom 1-Tap Action Pill Button
                Container(
                  width: double.infinity,
                  padding: const EdgeInsets.symmetric(vertical: 6.5),
                  decoration: BoxDecoration(
                    color: Colors.white.withAlpha(35),
                    borderRadius: BorderRadius.circular(10),
                    border: Border.all(
                        color: Colors.white.withAlpha(60), width: 0.8),
                  ),
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Text(
                        buttonLabel,
                        style: const TextStyle(
                          fontSize: 11.5,
                          fontWeight: FontWeight.bold,
                          color: Colors.white,
                        ),
                      ),
                      const SizedBox(width: 4),
                      const Icon(Icons.arrow_forward_rounded,
                          size: 12, color: Colors.white),
                    ],
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }

  (LinearGradient, Color, Color, String, IconData, String) _getKindStyle(
      CourseCardKind kind) {
    switch (kind) {
      case CourseCardKind.enrolled:
        return (
          const LinearGradient(
            begin: Alignment.topLeft,
            end: Alignment.bottomRight,
            colors: [Color(0xFF065F46), Color(0xFF047857)],
          ),
          const Color(0xFF10B981),
          const Color(0xFF059669),
          'ENROLLED',
          Icons.verified_rounded,
          'Continue Learning',
        );
      case CourseCardKind.favorite:
        return (
          const LinearGradient(
            begin: Alignment.topLeft,
            end: Alignment.bottomRight,
            colors: [Color(0xFFB45309), Color(0xFF92400E)],
          ),
          const Color(0xFFF59E0B),
          const Color(0xFFD97706),
          'FAVORITE',
          Icons.bookmark_rounded,
          'View Course',
        );
      case CourseCardKind.downloaded:
        return (
          const LinearGradient(
            begin: Alignment.topLeft,
            end: Alignment.bottomRight,
            colors: [Color(0xFF0284C7), Color(0xFF0369A1)],
          ),
          const Color(0xFF38BDF8),
          const Color(0xFF0284C7),
          'OFFLINE READY',
          Icons.download_done_rounded,
          'Start Offline',
        );
    }
  }
}

class _GuestSignInPromptCard extends StatelessWidget {
  const _GuestSignInPromptCard();

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(AppSizes.md),
      decoration: BoxDecoration(
        color: AppColors.bgSecondary,
        borderRadius: BorderRadius.circular(16),
        border: Border.all(color: AppColors.borderStrong),
      ),
      child: Row(
        children: [
          Container(
            padding: const EdgeInsets.all(8),
            decoration: BoxDecoration(
              color: AppColors.brandEmerald.withAlpha(30),
              shape: BoxShape.circle,
            ),
            child: const Icon(Icons.lock_outline_rounded,
                color: AppColors.brandEmerald, size: 20),
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const Text(
                  'Sign in to sync enrolled courses',
                  style: TextStyle(
                    fontSize: 13,
                    fontWeight: FontWeight.bold,
                    color: AppColors.textPrimary,
                  ),
                ),
                const SizedBox(height: 2),
                Text(
                  'Enrolled courses sync automatically when signed in.',
                  style: TextStyle(
                    fontSize: 11,
                    color: AppColors.textSecondary,
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../../../core/errors/failures.dart';
import '../../../../core/router/app_router.dart';
import '../../../auth/presentation/providers/auth_state_provider.dart';
import '../../domain/entities/course_entity.dart';
import '../providers/course_list_provider.dart';
import '../widgets/course_card.dart';
import '../widgets/course_empty_state.dart';
import '../widgets/course_list_skeleton.dart';
import '../widgets/subject_filter_chips.dart';

/// Professional Duolingo-Inspired Home Screen for Memere (ExamPrep).
///
/// Features a strict 4-color palette, tactile 3D cards/buttons, gamified daily target progress,
/// streak counter, XP tracker, and structured subject lanes.
class CourseListScreen extends ConsumerStatefulWidget {
  const CourseListScreen({super.key});

  @override
  ConsumerState<CourseListScreen> createState() => _CourseListScreenState();
}

class _CourseListScreenState extends ConsumerState<CourseListScreen> {
  late final TextEditingController _searchController;
  late final ScrollController _scrollController;

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
  });

  final ScrollController controller;
  final TextEditingController searchController;
  final CourseListState state;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return RefreshIndicator(
      color: AppColors.brandEmerald,
      backgroundColor: AppColors.bgSecondary,
      onRefresh: () => ref.read(courseListProvider.notifier).refresh(),
      child: CustomScrollView(
        controller: controller,
        physics: const AlwaysScrollableScrollPhysics(),
        slivers: [
          // 1. Duolingo Header (Avatar + Streak + XP)
          SliverToBoxAdapter(
            child: _DashboardHeader(searchController: searchController),
          ),

          // 2. Duolingo Daily Target Banner
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

          // 3. Subject Filters & Grade Switcher
          SliverToBoxAdapter(
            child: _DashboardFilters(state: state),
          ),

          // 4. Content List or Empty State
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
            // Subject Lanes Grid
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
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Text(
                      'Exam Prep Courses',
                      style: TextStyle(
                        fontSize: 18,
                        fontWeight: FontWeight.bold,
                        color: AppColors.textPrimary,
                        letterSpacing: -0.2,
                      ),
                    ),
                    Text(
                      'Grade 12',
                      style: TextStyle(
                        fontSize: 12,
                        fontWeight: FontWeight.w600,
                        color: AppColors.brandEmerald,
                      ),
                    ),
                  ],
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
                  border: Border.all(color: AppColors.brandEmerald, width: 2),
                ),
                child: Text(
                  initial,
                  style: const TextStyle(
                    fontSize: 18,
                    fontWeight: FontWeight.bold,
                    color: AppColors.brandEmerald,
                  ),
                ),
              ),
              const SizedBox(width: 12),

              // Greeting & Level
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
                      'Grade 12 National Prep',
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
                    Icon(Icons.bolt_rounded, color: AppColors.brandEmerald, size: 18),
                    SizedBox(width: 4),
                    Text(
                      '450',
                      style: TextStyle(
                        fontSize: 13,
                        fontWeight: FontWeight.bold,
                        color: AppColors.brandEmerald,
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

/// Duolingo Signature Tactile 3D Daily Goal Banner Card
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
          bottom: BorderSide(color: AppColors.brandEmeraldDark, width: 4), // 3D Tactile Border
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
                      color: AppColors.brandEmerald.withOpacity(0.15),
                      borderRadius: BorderRadius.circular(10),
                    ),
                    child: const Icon(
                      Icons.track_changes_rounded,
                      color: AppColors.brandEmerald,
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

          // Thick Duolingo Progress Bar
          ClipRRect(
            borderRadius: BorderRadius.circular(10),
            child: const LinearProgressIndicator(
              value: 0.6,
              minHeight: 12,
              backgroundColor: AppColors.bgPrimary,
              color: AppColors.brandEmerald,
            ),
          ),
          const SizedBox(height: 14),

          // Tactile Duolingo 3D Button
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

/// Tactile Search Bar Container
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
              color: AppColors.brandEmerald,
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
    final selectedGrade = state.selectedGrade ?? 12;

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
        const SizedBox(height: 12),
        Padding(
          padding: const EdgeInsets.symmetric(
            horizontal: AppSizes.screenPaddingH,
          ),
          child: Row(
            children: [
              const Text(
                'Grade Level',
                style: TextStyle(
                  fontSize: 13,
                  fontWeight: FontWeight.bold,
                  color: AppColors.textMuted,
                ),
              ),
              const SizedBox(width: 12),
              SegmentedButton<int>(
                segments: const [
                  ButtonSegment(value: 11, label: Text('Grade 11')),
                  ButtonSegment(value: 12, label: Text('Grade 12')),
                ],
                selected: {selectedGrade},
                onSelectionChanged: (selection) {
                  ref
                      .read(courseListProvider.notifier)
                      .setGrade(selection.first);
                },
              ),
            ],
          ),
        ),
      ],
    );
  }
}

/// Duolingo Path Subject Grid
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
              childAspectRatio: 2.6,
            ),
            itemBuilder: (context, index) {
              final subject = subjects[index];
              return Container(
                padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
                decoration: BoxDecoration(
                  color: AppColors.bgSecondary,
                  borderRadius: BorderRadius.circular(16),
                  border: const Border(
                    top: BorderSide(color: AppColors.borderStrong),
                    left: BorderSide(color: AppColors.borderStrong),
                    right: BorderSide(color: AppColors.borderStrong),
                    bottom: BorderSide(color: AppColors.borderStrong, width: 3), // Tactile 3D Depth
                  ),
                ),
                child: Row(
                  children: [
                    Container(
                      padding: const EdgeInsets.all(8),
                      decoration: BoxDecoration(
                        color: AppColors.brandEmerald.withOpacity(0.15),
                        shape: BoxShape.circle,
                      ),
                      child: Icon(
                        _subjectIcon(subject),
                        size: 18,
                        color: AppColors.brandEmerald,
                      ),
                    ),
                    const SizedBox(width: 10),
                    Expanded(
                      child: Text(
                        _shortSubject(subject),
                        maxLines: 1,
                        overflow: TextOverflow.ellipsis,
                        style: const TextStyle(
                          fontSize: 14,
                          fontWeight: FontWeight.bold,
                          color: AppColors.textPrimary,
                        ),
                      ),
                    ),
                    const Icon(
                      Icons.chevron_right_rounded,
                      size: 18,
                      color: AppColors.textMuted,
                    ),
                  ],
                ),
              );
            },
          ),
        ),
      ],
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

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_shadows.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../../../core/errors/failures.dart';
import '../../../../shared/widgets/app_surface.dart';
import '../../../auth/presentation/providers/auth_state_provider.dart';
import '../../domain/entities/course_entity.dart';
import '../providers/course_list_provider.dart';
import '../widgets/course_card.dart';
import '../widgets/course_empty_state.dart';
import '../widgets/course_list_skeleton.dart';
import '../widgets/subject_filter_chips.dart';

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
      body: AppPageBackground(
        child: SafeArea(
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
      color: AppColors.accentPrimary,
      backgroundColor: AppColors.bgSecondary,
      onRefresh: () => ref.read(courseListProvider.notifier).refresh(),
      child: CustomScrollView(
        controller: controller,
        physics: const AlwaysScrollableScrollPhysics(),
        slivers: [
          SliverToBoxAdapter(
            child: _DashboardHeader(searchController: searchController),
          ),
          SliverToBoxAdapter(
            child: _DashboardFilters(state: state),
          ),
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
            SliverToBoxAdapter(
              child: _TopicGrid(courses: state.filteredCourses),
            ),
            const SliverToBoxAdapter(
              child: Padding(
                padding: EdgeInsets.fromLTRB(
                  AppSizes.screenPaddingH,
                  AppSizes.xl,
                  AppSizes.screenPaddingH,
                  AppSizes.md,
                ),
                child: AppSectionHeader(
                  title: 'All Courses',
                  subtitle: 'Structured exam prep for Grade 12',
                ),
              ),
            ),
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
                  return AppStaggeredReveal(
                    index: index,
                    child: CourseRowCard(course: course),
                  );
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
          AppSurface(
            padding: const EdgeInsets.all(AppSizes.md),
            radius: AppSizes.radiusLg,
            borderColor: AppColors.borderStrong,
            shadows: AppShadows.md,
            child: Row(
              children: [
                _Avatar(initial: initial),
                const SizedBox(width: AppSizes.md),
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        greeting,
                        maxLines: 1,
                        overflow: TextOverflow.ellipsis,
                        style: AppTextStyles.headlineMedium,
                      ),
                      const SizedBox(height: AppSizes.xs),
                      const Text(
                        'Grade 12 exam prep',
                        maxLines: 1,
                        overflow: TextOverflow.ellipsis,
                        style: AppTextStyles.bodySmall,
                      ),
                    ],
                  ),
                ),
                const SizedBox(width: AppSizes.sm),
                const _HeaderIconButton(icon: Icons.notifications_none),
                const SizedBox(width: AppSizes.xs),
                const _HeaderIconButton(icon: Icons.bookmark_border),
              ],
            ),
          ),
          const SizedBox(height: AppSizes.md),
          _FloatingSearchBar(searchController: searchController),
        ],
      ),
    );
  }
}

class _FloatingSearchBar extends ConsumerWidget {
  const _FloatingSearchBar({required this.searchController});

  final TextEditingController searchController;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return AppSurface(
      padding: EdgeInsets.zero,
      radius: AppSizes.radiusMd,
      color: AppColors.bgSecondary,
      borderColor: AppColors.borderStrong,
      shadows: AppShadows.sm,
      child: SizedBox(
        height: AppSizes.inputHeight,
        child: Row(
          children: [
            const SizedBox(width: AppSizes.md),
            const Icon(
              Icons.search_rounded,
              size: AppSizes.iconMd,
              color: AppColors.textPrimary,
            ),
            const SizedBox(width: AppSizes.sm),
            Expanded(
              child: TextField(
                controller: searchController,
                textInputAction: TextInputAction.search,
                onChanged: ref.read(courseListProvider.notifier).setSearchQuery,
                style: AppTextStyles.bodyMedium,
                decoration: InputDecoration(
                  hintText: 'Search lessons, subjects, exams',
                  hintStyle: AppTextStyles.bodyMedium.copyWith(
                    color: AppColors.textMuted,
                  ),
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
                color: AppColors.accentPrimary,
              ),
            ),
            const SizedBox(width: AppSizes.xs),
          ],
        ),
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
        const SizedBox(height: AppSizes.md),
        SubjectFilterChips(
          subjects: phase2Subjects,
          selectedSubject: state.selectedSubject,
          onSelected: (subject) {
            ref.read(courseListProvider.notifier).setSubject(subject);
          },
        ),
        const SizedBox(height: AppSizes.md),
        Padding(
          padding: const EdgeInsets.symmetric(
            horizontal: AppSizes.screenPaddingH,
          ),
          child: Row(
            children: [
              Text(
                'Grade',
                style: AppTextStyles.labelMedium.copyWith(
                  color: AppColors.textSecondary,
                ),
              ),
              const SizedBox(width: AppSizes.md),
              SegmentedButton<int>(
                segments: const [
                  ButtonSegment(value: 11, label: Text('11')),
                  ButtonSegment(value: 12, label: Text('12')),
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
            AppSizes.md,
          ),
          child: AppSectionHeader(
            title: 'Browse Subjects',
            subtitle: 'Pick a focused lane for today',
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
              mainAxisSpacing: AppSizes.sm,
              crossAxisSpacing: AppSizes.sm,
              childAspectRatio: 2.8,
            ),
            itemBuilder: (context, index) {
              final subject = subjects[index];
              return AppSurface(
                padding: const EdgeInsets.symmetric(
                  horizontal: AppSizes.md,
                  vertical: AppSizes.sm,
                ),
                radius: AppSizes.radiusMd,
                borderColor: AppColors.borderStrong,
                shadows: AppShadows.sm,
                child: Row(
                  children: [
                    AppIconTile(
                      icon: _subjectIcon(subject),
                      size: 34,
                      iconSize: AppSizes.iconSm,
                      color: _subjectColor(subject),
                    ),
                    const SizedBox(width: AppSizes.sm),
                    Expanded(
                      child: Text(
                        _shortSubject(subject),
                        maxLines: 2,
                        overflow: TextOverflow.ellipsis,
                        style: AppTextStyles.labelMedium.copyWith(
                          fontWeight: FontWeight.w700,
                        ),
                      ),
                    ),
                    const Icon(
                      Icons.chevron_right_rounded,
                      size: AppSizes.iconSm,
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

class _Avatar extends StatelessWidget {
  const _Avatar({required this.initial});

  final String initial;

  @override
  Widget build(BuildContext context) {
    return Container(
      width: AppSizes.avatarMd,
      height: AppSizes.avatarMd,
      alignment: Alignment.center,
      decoration: BoxDecoration(
        color: AppColors.bgTertiary,
        shape: BoxShape.circle,
        border: Border.all(color: AppColors.borderStrong),
        boxShadow: AppShadows.md,
      ),
      child: Text(
        initial,
        style: AppTextStyles.titleMedium.copyWith(
          color: AppColors.accentPrimary,
        ),
      ),
    );
  }
}

class _HeaderIconButton extends StatelessWidget {
  const _HeaderIconButton({required this.icon});

  final IconData icon;

  @override
  Widget build(BuildContext context) {
    return AppPressable(
      onTap: () {},
      borderRadius: AppSizes.radiusFull,
      child: Container(
        width: 34,
        height: 34,
        alignment: Alignment.center,
        decoration: BoxDecoration(
          color: AppColors.bgTertiary,
          shape: BoxShape.circle,
          border: Border.all(color: AppColors.border),
        ),
        child:
            Icon(icon, color: AppColors.textSecondary, size: AppSizes.iconSm),
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

Color _subjectColor(String subject) {
  switch (subject.toLowerCase()) {
    case 'mathematics':
    case 'math':
      return AppColors.subjectMath;
    case 'physics':
      return AppColors.subjectPhysics;
    case 'chemistry':
      return AppColors.subjectChem;
    case 'biology':
      return AppColors.subjectBio;
    case 'english':
      return AppColors.subjectEng;
    case 'history':
      return AppColors.subjectHist;
    case 'geography':
      return AppColors.subjectGeo;
    case 'economics':
      return AppColors.subjectEcon;
    default:
      return AppColors.subjectNeutral;
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
            color: AppColors.accentPrimary,
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

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../../../core/errors/failures.dart';
import '../../../../core/router/app_router.dart';
import '../../../../shared/widgets/app_text_field.dart';
import '../../../auth/presentation/providers/auth_state_provider.dart';
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
      body: SafeArea(
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
                _CourseListHeader(onLogout: _logout),
                _SearchAndFilters(searchController: _searchController),
                const Expanded(child: CourseListSkeleton()),
              ],
            );
          },
          error: (error, _) => Column(
            children: [
              _CourseListHeader(onLogout: _logout),
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
            onLogout: _logout,
          ),
        ),
      ),
    );
  }

  Future<void> _logout() async {
    await ref.read(authStateProvider.notifier).logout();
    if (mounted) context.go(AppRoutes.login);
  }
}

class _CourseListContent extends ConsumerWidget {
  const _CourseListContent({
    required this.controller,
    required this.searchController,
    required this.state,
    this.onLogout,
  });

  final ScrollController controller;
  final TextEditingController searchController;
  final CourseListState state;
  final VoidCallback? onLogout;

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
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                _CourseListHeader(onLogout: onLogout),
                _SearchAndFilters(searchController: searchController),
              ],
            ),
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
          else
            SliverPadding(
              padding: const EdgeInsets.fromLTRB(
                AppSizes.screenPaddingH,
                AppSizes.md,
                AppSizes.screenPaddingH,
                AppSizes.lg,
              ),
              sliver: SliverList.separated(
                itemCount: state.filteredCourses.length,
                separatorBuilder: (_, __) =>
                    const SizedBox(height: AppSizes.md),
                itemBuilder: (context, index) {
                  return CourseCard(course: state.filteredCourses[index]);
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

class _CourseListHeader extends ConsumerWidget {
  const _CourseListHeader({this.onLogout});

  final VoidCallback? onLogout;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final user = ref.watch(authStateProvider).valueOrNull?.user;
    final firstName = user?.firstName.trim();

    return Padding(
      padding: const EdgeInsets.fromLTRB(
        AppSizes.screenPaddingH,
        AppSizes.md,
        AppSizes.screenPaddingH,
        AppSizes.lg,
      ),
      child: Row(
        children: [
          Container(
            width: 44,
            height: 44,
            decoration: BoxDecoration(
              gradient: AppColors.primaryGradient,
              borderRadius: BorderRadius.circular(AppSizes.radiusMd),
            ),
            child: const Icon(
              Icons.school_rounded,
              color: Colors.white,
              size: AppSizes.iconMd,
            ),
          ),
          const SizedBox(width: AppSizes.md),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  firstName == null || firstName.isEmpty
                      ? 'Explore courses'
                      : 'Hi, $firstName',
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                  style: AppTextStyles.headlineMedium,
                ),
                const SizedBox(height: AppSizes.xs),
                const Text(
                  'Build your Grade 12 exam plan',
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                  style: AppTextStyles.bodySmall,
                ),
              ],
            ),
          ),
          IconButton(
            tooltip: 'Sign out',
            onPressed: onLogout,
            icon: const Icon(Icons.logout_rounded),
          ),
        ],
      ),
    );
  }
}

class _SearchAndFilters extends ConsumerWidget {
  const _SearchAndFilters({required this.searchController});

  final TextEditingController searchController;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final state = ref.watch(courseListProvider).valueOrNull;
    final selectedGrade = state?.selectedGrade ?? 12;

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Padding(
          padding: const EdgeInsets.symmetric(
            horizontal: AppSizes.screenPaddingH,
          ),
          child: AppTextField(
            controller: searchController,
            hintText: 'Search courses',
            prefixIcon: Icons.search_rounded,
            textInputAction: TextInputAction.search,
            onChanged: ref.read(courseListProvider.notifier).setSearchQuery,
          ),
        ),
        const SizedBox(height: AppSizes.md),
        SubjectFilterChips(
          subjects: phase2Subjects,
          selectedSubject: state?.selectedSubject,
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
                style: ButtonStyle(
                  visualDensity: VisualDensity.compact,
                  foregroundColor: WidgetStateProperty.resolveWith((states) {
                    if (states.contains(WidgetState.selected)) {
                      return Colors.white;
                    }
                    return AppColors.textSecondary;
                  }),
                  backgroundColor: WidgetStateProperty.resolveWith((states) {
                    if (states.contains(WidgetState.selected)) {
                      return AppColors.accentPrimary;
                    }
                    return AppColors.bgTertiary;
                  }),
                  side: WidgetStateProperty.all(
                    const BorderSide(color: AppColors.border),
                  ),
                ),
              ),
            ],
          ),
        ),
      ],
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

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_shadows.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../../../core/errors/failures.dart';
import '../../../../shared/widgets/app_text_field.dart';
import '../../../../shared/widgets/app_surface.dart';
import '../../../courses/presentation/widgets/subject_filter_chips.dart';
import '../providers/mock_exam_catalog_provider.dart';
import '../widgets/exam_empty_state.dart';
import '../widgets/mock_exam_card.dart';
import '../widgets/mock_exam_catalog_skeleton.dart';

class MockExamCatalogScreen extends ConsumerStatefulWidget {
  const MockExamCatalogScreen({super.key});

  @override
  ConsumerState<MockExamCatalogScreen> createState() =>
      _MockExamCatalogScreenState();
}

class _MockExamCatalogScreenState extends ConsumerState<MockExamCatalogScreen> {
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
      ref.read(mockExamCatalogProvider.notifier).loadMore();
    }
  }

  @override
  Widget build(BuildContext context) {
    final examsAsync = ref.watch(mockExamCatalogProvider);

    ref.listen(mockExamCatalogProvider, (_, next) {
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
          child: examsAsync.when(
            loading: () {
              final previous = examsAsync.valueOrNull;
              if (previous != null && previous.exams.isNotEmpty) {
                return _CatalogContent(
                  controller: _scrollController,
                  searchController: _searchController,
                  state: previous,
                );
              }
              return Column(
                children: [
                  const _CatalogHeader(),
                  _SearchAndFilters(searchController: _searchController),
                  const Expanded(child: MockExamCatalogSkeleton()),
                ],
              );
            },
            error: (error, _) => Column(
              children: [
                const _CatalogHeader(),
                Expanded(
                  child: ExamEmptyState(
                    icon: Icons.wifi_off_rounded,
                    title: 'Could not load mock exams',
                    body: error is Failure
                        ? error.message
                        : 'Check your connection and try again.',
                    buttonLabel: 'Retry',
                    onPressed: () => ref.invalidate(mockExamCatalogProvider),
                  ),
                ),
              ],
            ),
            data: (state) => _CatalogContent(
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

class _CatalogContent extends ConsumerWidget {
  const _CatalogContent({
    required this.controller,
    required this.searchController,
    required this.state,
  });

  final ScrollController controller;
  final TextEditingController searchController;
  final MockExamCatalogState state;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return RefreshIndicator(
      color: AppColors.accentPrimary,
      backgroundColor: AppColors.bgSecondary,
      onRefresh: () => ref.read(mockExamCatalogProvider.notifier).refresh(),
      child: CustomScrollView(
        controller: controller,
        physics: const AlwaysScrollableScrollPhysics(),
        slivers: [
          SliverToBoxAdapter(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const _CatalogHeader(),
                _SearchAndFilters(searchController: searchController),
              ],
            ),
          ),
          if (state.filteredExams.isEmpty)
            SliverFillRemaining(
              hasScrollBody: false,
              child: ExamEmptyState(
                icon: Icons.search_off_rounded,
                title: 'No mock exams found',
                body: 'Try another subject or search term.',
                buttonLabel: state.hasActiveFilters ? 'Clear filters' : null,
                onPressed: state.hasActiveFilters
                    ? () => ref
                        .read(mockExamCatalogProvider.notifier)
                        .clearFilters()
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
                itemCount: state.filteredExams.length,
                separatorBuilder: (_, __) =>
                    const SizedBox(height: AppSizes.md),
                itemBuilder: (context, index) {
                  return AppStaggeredReveal(
                    index: index,
                    child: MockExamCard(exam: state.filteredExams[index]),
                  );
                },
              ),
            ),
          SliverToBoxAdapter(child: _LoadMoreFooter(state: state)),
        ],
      ),
    );
  }
}

class _CatalogHeader extends StatelessWidget {
  const _CatalogHeader();

  @override
  Widget build(BuildContext context) {
    return const Padding(
      padding: EdgeInsets.fromLTRB(
        AppSizes.screenPaddingH,
        AppSizes.md,
        AppSizes.screenPaddingH,
        AppSizes.lg,
      ),
      child: AppSurface(
        padding: EdgeInsets.all(AppSizes.md),
        gradient: AppColors.cardGradient,
        shadows: AppShadows.md,
        child: Row(
          children: [
            AppIconTile(
              icon: Icons.assignment_rounded,
              gradient: AppColors.examGradient,
            ),
            SizedBox(width: AppSizes.md),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    'Mock exams',
                    maxLines: 1,
                    overflow: TextOverflow.ellipsis,
                    style: AppTextStyles.headlineMedium,
                  ),
                  SizedBox(height: AppSizes.xs),
                  Text(
                    'Practice under real exam timing',
                    maxLines: 1,
                    overflow: TextOverflow.ellipsis,
                    style: AppTextStyles.bodySmall,
                  ),
                ],
              ),
            ),
            AppBadge(
              label: 'Timed',
              color: AppColors.info,
              icon: Icons.timer_outlined,
            ),
          ],
        ),
      ),
    );
  }
}

class _SearchAndFilters extends ConsumerWidget {
  const _SearchAndFilters({required this.searchController});

  final TextEditingController searchController;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final state = ref.watch(mockExamCatalogProvider).valueOrNull;
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
            hintText: 'Search mock exams',
            prefixIcon: Icons.search_rounded,
            textInputAction: TextInputAction.search,
            onChanged:
                ref.read(mockExamCatalogProvider.notifier).setSearchQuery,
          ),
        ),
        const SizedBox(height: AppSizes.md),
        SubjectFilterChips(
          subjects: mockExamSubjects,
          selectedSubject: state?.selectedSubject,
          onSelected: (subject) {
            ref.read(mockExamCatalogProvider.notifier).setSubject(subject);
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
              Expanded(
                child: Align(
                  alignment: Alignment.centerLeft,
                  child: SegmentedButton<int>(
                    segments: const [
                      ButtonSegment(value: 11, label: Text('11')),
                      ButtonSegment(value: 12, label: Text('12')),
                    ],
                    selected: {selectedGrade},
                    onSelectionChanged: (selection) {
                      ref
                          .read(mockExamCatalogProvider.notifier)
                          .setGrade(selection.first);
                    },
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

  final MockExamCatalogState state;

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

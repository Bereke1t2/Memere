import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../../../core/errors/failures.dart';
import '../../auth/presentation/providers/auth_state_provider.dart';
import '../../../progress/presentation/providers/progress_providers.dart';
import '../providers/exam_providers.dart';
import '../providers/mock_exam_catalog_provider.dart';
import '../widgets/exam_empty_state.dart';
import '../widgets/mock_exam_card.dart';
import '../widgets/mock_exam_catalog_skeleton.dart';

/// Refined Obsidian & Soft Emerald Exam Catalog Screen.
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
      body: SafeArea(
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
            return const Column(
              children: [
                _CatalogHeader(),
                Expanded(child: MockExamCatalogSkeleton()),
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
      color: AppColors.brandEmerald,
      backgroundColor: AppColors.bgSecondary,
      onRefresh: () async {
        ref.invalidate(myAllExamAttemptsProvider);
        ref.invalidate(studentPointsProvider);
        await ref.read(mockExamCatalogProvider.notifier).refresh();
      },
      child: CustomScrollView(
        controller: controller,
        physics: const AlwaysScrollableScrollPhysics(),
        slivers: [
          // 1. Top Header Row + Greeting + Circular Category Badges
          const SliverToBoxAdapter(
            child: _CatalogHeader(),
          ),

          // 2. 2-Column Grid of Vibrant Gradient Exam Cards
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
                0,
                AppSizes.screenPaddingH,
                AppSizes.md,
              ),
              sliver: SliverGrid(
                gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
                  crossAxisCount: 2,
                  mainAxisSpacing: 10,
                  crossAxisSpacing: 10,
                  childAspectRatio: 0.72,
                ),
                delegate: SliverChildBuilderDelegate(
                  (context, index) {
                    final exam = state.filteredExams[index];
                    return MockExamCard(exam: exam, cardIndex: index);
                  },
                  childCount: state.filteredExams.length,
                ),
              ),
            ),
          SliverToBoxAdapter(child: _LoadMoreFooter(state: state)),
        ],
      ),
    );
  }
}

/// Top Header Row (Logo + Dynamic Points Badge) + Clean Circular Subject Badges
class _CatalogHeader extends ConsumerWidget {
  const _CatalogHeader();

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final catalogState = ref.watch(mockExamCatalogProvider).valueOrNull;
    final currentSubject = catalogState?.selectedSubject;
    final authState = ref.watch(authStateProvider).valueOrNull;
    final isAuthenticated = authState?.isAuthenticated ?? false;
    final pointsAsync = ref.watch(studentPointsProvider);

    final String pointsText;
    if (!isAuthenticated) {
      pointsText = '0 pts';
    } else {
      pointsText = pointsAsync.when(
        data: (pts) => '${_formatPoints(pts.totalPoints)} pts',
        loading: () => '... pts',
        error: (_, __) => '0 pts',
      );
    }

    final categories = <(IconData, String, String?)>[
      (Icons.grid_view_rounded, 'All', null),
      (Icons.calculate_outlined, 'Math', 'Mathematics'),
      (Icons.science_outlined, 'Physics', 'Physics'),
      (Icons.biotech_outlined, 'Chemistry', 'Chemistry'),
      (Icons.eco_outlined, 'Biology', 'Biology'),
      (Icons.menu_book_outlined, 'English', 'English'),
      (Icons.history_edu_outlined, 'History', 'History'),
      (Icons.public_outlined, 'Geography', 'Geography'),
      (Icons.trending_up_rounded, 'Economics', 'Economics'),
    ];

    return Padding(
      padding: const EdgeInsets.fromLTRB(
        0,
        AppSizes.sm,
        0,
        AppSizes.xs,
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Top Bar: App Branding + Points / Score Pill
          Padding(
            padding: const EdgeInsets.symmetric(
              horizontal: AppSizes.screenPaddingH,
            ),
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
                      'Memere Exam Hub',
                      style: TextStyle(
                        fontSize: 18,
                        fontWeight: FontWeight.w900,
                        color: AppColors.textPrimary,
                        letterSpacing: -0.4,
                      ),
                    ),
                  ],
                ),

                // Dynamic Points / Score Pill Badge
                Container(
                  padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
                  decoration: BoxDecoration(
                    color: AppColors.bgSecondary,
                    borderRadius: BorderRadius.circular(16),
                    border: Border.all(color: AppColors.borderStrong),
                  ),
                  child: Row(
                    children: [
                      const Icon(Icons.stars_rounded, color: AppColors.brandAmber, size: 16),
                      const SizedBox(width: 5),
                      Text(
                        pointsText,
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

          // Horizontal Scrollable Row of Outlined Circular Subject Icon Badges
          SingleChildScrollView(
            scrollDirection: Axis.horizontal,
            physics: const BouncingScrollPhysics(),
            padding: const EdgeInsets.symmetric(
              horizontal: AppSizes.screenPaddingH,
            ),
            child: Row(
              children: List.generate(categories.length, (idx) {
                final cat = categories[idx];
                final isSelected = (cat.$3 == null && currentSubject == null) ||
                    (cat.$3 != null && cat.$3 == currentSubject);

                return Padding(
                  padding: EdgeInsets.only(
                    right: idx == categories.length - 1 ? 0 : 12,
                  ),
                  child: GestureDetector(
                    onTap: () {
                      ref
                          .read(mockExamCatalogProvider.notifier)
                          .setSubject(cat.$3);
                    },
                    child: Column(
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        AnimatedContainer(
                          duration: const Duration(milliseconds: 180),
                          width: 48,
                          height: 48,
                          alignment: Alignment.center,
                          decoration: BoxDecoration(
                            color: isSelected
                                ? AppColors.brandEmerald.withAlpha(35)
                                : AppColors.bgSecondary,
                            shape: BoxShape.circle,
                            border: Border.all(
                              color: isSelected
                                  ? AppColors.brandEmerald
                                  : AppColors.borderStrong,
                              width: isSelected ? 2.0 : 1.2,
                            ),
                            boxShadow: isSelected
                                ? [
                                    BoxShadow(
                                      color:
                                          AppColors.brandEmerald.withAlpha(60),
                                      blurRadius: 8,
                                      offset: const Offset(0, 2),
                                    ),
                                  ]
                                : null,
                          ),
                          child: Icon(
                            cat.$1,
                            size: 22,
                            color: isSelected
                                ? AppColors.brandEmerald
                                : AppColors.textSecondary,
                          ),
                        ),
                        const SizedBox(height: 6),
                        Text(
                          cat.$2,
                          style: TextStyle(
                            fontSize: 11,
                            fontWeight: isSelected
                                ? FontWeight.bold
                                : FontWeight.w500,
                            color: isSelected
                                ? AppColors.brandEmerald
                                : AppColors.textMuted,
                          ),
                        ),
                      ],
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

String _formatPoints(double value) {
  if (value == value.roundToDouble()) return value.toInt().toString();
  return value.toStringAsFixed(1);
}

class _LoadMoreFooter extends StatelessWidget {
  const _LoadMoreFooter({required this.state});

  final MockExamCatalogState state;

  @override
  Widget build(BuildContext context) {
    if (state.isLoadingMore) {
      return const Padding(
        padding: EdgeInsets.only(bottom: AppSizes.lg),
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
          AppSizes.lg,
        ),
        child: Text(
          state.loadMoreError!.message,
          textAlign: TextAlign.center,
          style: AppTextStyles.bodySmall.copyWith(color: AppColors.error),
        ),
      );
    }

    return const SizedBox(height: AppSizes.md);
  }
}

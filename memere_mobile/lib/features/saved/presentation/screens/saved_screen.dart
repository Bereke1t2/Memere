import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/offline/offline_attempt_factory.dart';
import '../../../../core/offline/offline_providers.dart';
import '../../../../core/router/app_router.dart';
import '../../../../core/storage/hive/models/downloaded_item.dart';
import '../../../../core/storage/hive/models/saved_item.dart';
import '../../../exam/presentation/providers/exam_attempt_provider.dart';
import '../../../exam/presentation/providers/exam_download_provider.dart';
import '../../../exam/presentation/widgets/exam_empty_state.dart';
import '../../../quiz/presentation/providers/quiz_download_provider.dart';
import '../providers/saved_courses_provider.dart';

enum SavedCategoryFilter { all, courses, quizzes, exams }

enum SavedKind { course, quiz, exam }

class _SavedCardData {
  const _SavedCardData({
    required this.id,
    required this.kind,
    required this.title,
    required this.subtitle,
    required this.savedAt,
    this.downloadedItem,
    this.savedCourseItem,
  });

  final String id;
  final SavedKind kind;
  final String title;
  final String subtitle;
  final DateTime savedAt;
  final DownloadedItem? downloadedItem;
  final SavedItem? savedCourseItem;
}

/// Refined Obsidian & Soft Emerald Saved Screen matching Exam & Home catalog UX.
class SavedScreen extends ConsumerStatefulWidget {
  const SavedScreen({super.key});

  @override
  ConsumerState<SavedScreen> createState() => _SavedScreenState();
}

class _SavedScreenState extends ConsumerState<SavedScreen> {
  SavedCategoryFilter _selectedFilter = SavedCategoryFilter.all;

  @override
  Widget build(BuildContext context) {
    final coursesAsync = ref.watch(savedCoursesProvider);
    final quizzesAsync = ref.watch(downloadedQuizzesProvider);
    final examsAsync = ref.watch(downloadedExamsProvider);

    final courses = coursesAsync.valueOrNull ?? const [];
    final quizzes = quizzesAsync.valueOrNull ?? const [];
    final exams = examsAsync.valueOrNull ?? const [];

    final allItems = <_SavedCardData>[
      ...courses.map((c) => _SavedCardData(
            id: c.id,
            kind: SavedKind.course,
            title: c.title,
            subtitle: c.subtitle ?? 'Saved Course',
            savedAt: c.savedAt,
            savedCourseItem: c,
          )),
      ...quizzes.map((q) => _SavedCardData(
            id: q.id,
            kind: SavedKind.quiz,
            title: q.title,
            subtitle: q.subtitle ?? 'Offline Quiz',
            savedAt: q.downloadedAt,
            downloadedItem: q,
          )),
      ...exams.map((e) => _SavedCardData(
            id: e.id,
            kind: SavedKind.exam,
            title: e.title,
            subtitle: e.subtitle ?? 'Offline Exam',
            savedAt: e.downloadedAt,
            downloadedItem: e,
          )),
    ]..sort((a, b) => b.savedAt.compareTo(a.savedAt));

    final filteredItems = allItems.where((item) {
      switch (_selectedFilter) {
        case SavedCategoryFilter.all:
          return true;
        case SavedCategoryFilter.courses:
          return item.kind == SavedKind.course;
        case SavedCategoryFilter.quizzes:
          return item.kind == SavedKind.quiz;
        case SavedCategoryFilter.exams:
          return item.kind == SavedKind.exam;
      }
    }).toList();

    final isLoading = allItems.isEmpty &&
        (coursesAsync.isLoading ||
            quizzesAsync.isLoading ||
            examsAsync.isLoading);

    return Scaffold(
      backgroundColor: AppColors.bgPrimary,
      body: SafeArea(
        child: RefreshIndicator(
          color: AppColors.brandEmerald,
          backgroundColor: AppColors.bgSecondary,
          onRefresh: () async {
            ref.invalidate(savedCoursesProvider);
            ref.invalidate(downloadedQuizzesProvider);
            ref.invalidate(downloadedExamsProvider);
          },
          child: CustomScrollView(
            physics: const AlwaysScrollableScrollPhysics(),
            slivers: [
              // 1. Top Header Row + Category Badges
              SliverToBoxAdapter(
                child: _SavedHeader(
                  selectedFilter: _selectedFilter,
                  onSelectFilter: (filter) =>
                      setState(() => _selectedFilter = filter),
                  totalCount: allItems.length,
                  coursesCount: courses.length,
                  quizzesCount: quizzes.length,
                  examsCount: exams.length,
                ),
              ),

              // 2. 2-Column Grid of Vibrant Gradient Cards or Empty State
              if (isLoading)
                const SliverFillRemaining(
                  child: Center(
                    child: CircularProgressIndicator(
                      color: AppColors.brandEmerald,
                    ),
                  ),
                )
              else if (filteredItems.isEmpty)
                SliverFillRemaining(
                  hasScrollBody: false,
                  child: ExamEmptyState(
                    icon: Icons.bookmark_outline_rounded,
                    title: allItems.isEmpty
                        ? 'Nothing saved yet'
                        : 'No ${_selectedFilter.name} found',
                    body: allItems.isEmpty
                        ? 'Save courses or download quizzes & exams to access them anytime offline.'
                        : 'Try switching your filter above to view other saved items.',
                    buttonLabel: allItems.isEmpty
                        ? 'Explore Courses'
                        : (_selectedFilter != SavedCategoryFilter.all
                            ? 'Show All'
                            : null),
                    onPressed: allItems.isEmpty
                        ? () => context.go(AppRoutes.home)
                        : (_selectedFilter != SavedCategoryFilter.all
                            ? () => setState(
                                () => _selectedFilter = SavedCategoryFilter.all)
                            : null),
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
                        return _SavedCardTile(
                          item: item,
                          cardIndex: index,
                          onOpen: () => _openItem(context, ref, item),
                          onRemove: () => _confirmRemove(context, ref, item),
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

  void _openItem(BuildContext context, WidgetRef ref, _SavedCardData item) {
    switch (item.kind) {
      case SavedKind.course:
        context.push(AppRoutes.courseDetailPath(item.id));
        break;
      case SavedKind.quiz:
        context.push(
          AppRoutes.quizAttemptPath(
            attemptId: newLocalAttemptId(),
            quizId: item.id,
          ),
        );
        break;
      case SavedKind.exam:
        final exam = ref.read(downloadStoreProvider).getOfflineExam(item.id);
        if (exam == null) {
          ScaffoldMessenger.of(context).showSnackBar(
            const SnackBar(content: Text('This exam is not downloaded.')),
          );
          return;
        }
        final localId = newLocalAttemptId();
        ref.read(pendingExamAttemptProvider.notifier).state =
            buildLocalExamAttempt(exam, attemptId: localId);
        context.push(
          AppRoutes.examAttemptPath(attemptId: localId, examId: item.id),
        );
        break;
    }
  }

  Future<void> _confirmRemove(
      BuildContext context, WidgetRef ref, _SavedCardData item) async {
    final kindName = item.kind == SavedKind.course
        ? 'saved course'
        : (item.kind == SavedKind.quiz
            ? 'downloaded quiz'
            : 'downloaded exam');

    final shouldDelete = await showDialog<bool>(
      context: context,
      builder: (ctx) => AlertDialog(
        backgroundColor: const Color(0xFF141824),
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
        title: Text(
          'Remove $kindName?',
          style: const TextStyle(
            color: Colors.white,
            fontSize: 16,
            fontWeight: FontWeight.bold,
          ),
        ),
        content: Text(
          'This removes "${item.title}" from your saved items. You can save or download it again anytime.',
          style: const TextStyle(color: Colors.white70, fontSize: 13),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(ctx, false),
            child: const Text('Cancel',
                style: TextStyle(color: AppColors.textMuted)),
          ),
          ElevatedButton(
            onPressed: () => Navigator.pop(ctx, true),
            style: ElevatedButton.styleFrom(
              backgroundColor: AppColors.error,
              shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(10)),
            ),
            child: const Text('Remove',
                style: TextStyle(
                    color: Colors.white, fontWeight: FontWeight.bold)),
          ),
        ],
      ),
    );

    if (shouldDelete != true) return;

    switch (item.kind) {
      case SavedKind.course:
        await ref.read(savedCoursesProvider.notifier).remove(item.id);
        break;
      case SavedKind.quiz:
        await ref.read(downloadedQuizzesProvider.notifier).remove(item.id);
        break;
      case SavedKind.exam:
        await ref.read(downloadedExamsProvider.notifier).remove(item.id);
        break;
    }

    if (context.mounted) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Removed "${item.title}" from saved items.')),
      );
    }
  }
}

/// Top Header Row (Branding + Items Badge) + Category Selector Badges
class _SavedHeader extends ConsumerWidget {
  const _SavedHeader({
    required this.selectedFilter,
    required this.onSelectFilter,
    required this.totalCount,
    required this.coursesCount,
    required this.quizzesCount,
    required this.examsCount,
  });

  final SavedCategoryFilter selectedFilter;
  final ValueChanged<SavedCategoryFilter> onSelectFilter;
  final int totalCount;
  final int coursesCount;
  final int quizzesCount;
  final int examsCount;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final filterOptions = [
      (SavedCategoryFilter.all, Icons.grid_view_rounded, 'All', totalCount),
      (
        SavedCategoryFilter.courses,
        Icons.menu_book_outlined,
        'Courses',
        coursesCount
      ),
      (
        SavedCategoryFilter.quizzes,
        Icons.quiz_outlined,
        'Quizzes',
        quizzesCount
      ),
      (
        SavedCategoryFilter.exams,
        Icons.assignment_outlined,
        'Exams',
        examsCount
      ),
    ];

    return Padding(
      padding: const EdgeInsets.fromLTRB(0, AppSizes.sm, 0, AppSizes.xs),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Top Bar: Branding + Saved Counter Badge
          Padding(
            padding: const EdgeInsets.symmetric(
                horizontal: AppSizes.screenPaddingH),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                const Row(
                  children: [
                    Icon(
                      Icons.bookmark_rounded,
                      color: AppColors.brandEmerald,
                      size: 22,
                    ),
                    SizedBox(width: 8),
                    Text(
                      'Memere Saved Hub',
                      style: TextStyle(
                        fontSize: 18,
                        fontWeight: FontWeight.w900,
                        color: AppColors.textPrimary,
                        letterSpacing: -0.4,
                      ),
                    ),
                  ],
                ),

                // Saved items count pill badge
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
                      const Icon(Icons.bookmark_added_rounded,
                          color: AppColors.brandEmerald, size: 16),
                      const SizedBox(width: 5),
                      Text(
                        '$totalCount ${totalCount == 1 ? 'item' : 'items'}',
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

class _SavedCardTile extends ConsumerStatefulWidget {
  const _SavedCardTile({
    required this.item,
    required this.cardIndex,
    required this.onOpen,
    required this.onRemove,
  });

  final _SavedCardData item;
  final int cardIndex;
  final VoidCallback onOpen;
  final VoidCallback onRemove;

  @override
  ConsumerState<_SavedCardTile> createState() => _SavedCardTileState();
}

class _SavedCardTileState extends ConsumerState<_SavedCardTile> {
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
                color:
                    _isPressed ? AppColors.borderFocused : glowColor.withAlpha(60),
                width: 1.2,
              ),
              boxShadow: [
                BoxShadow(
                  color: glowColor.withAlpha(20),
                  blurRadius: 8,
                  offset: const Offset(0, 3),
                ),
              ],
            ),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                // Top Row: Category Pill Badge + Delete Button
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Container(
                      padding: const EdgeInsets.symmetric(
                          horizontal: 8, vertical: 3),
                      decoration: BoxDecoration(
                        color: AppColors.bgTertiary,
                        borderRadius: BorderRadius.circular(6),
                        border: Border.all(
                            color: glowColor.withAlpha(80), width: 0.8),
                      ),
                      child: Text(
                        badgeLabel,
                        style: TextStyle(
                          fontSize: 10,
                          fontWeight: FontWeight.w800,
                          color: glowColor,
                          letterSpacing: 0.4,
                        ),
                      ),
                    ),
                    InkWell(
                      onTap: widget.onRemove,
                      borderRadius: BorderRadius.circular(12),
                      child: Container(
                        padding: const EdgeInsets.all(4),
                        decoration: BoxDecoration(
                          color: AppColors.bgTertiary,
                          shape: BoxShape.circle,
                          border: Border.all(color: AppColors.borderStrong),
                        ),
                        child: const Icon(
                          Icons.delete_outline_rounded,
                          size: 15,
                          color: AppColors.textMuted,
                        ),
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 8),

                // Middle: Icon + Title + Subtitle
                Row(
                  children: [
                    Container(
                      width: 34,
                      height: 34,
                      decoration: BoxDecoration(
                        color: glowColor.withAlpha(25),
                        shape: BoxShape.circle,
                        border: Border.all(color: glowColor.withAlpha(60)),
                      ),
                      child: Icon(iconData, size: 16, color: glowColor),
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
                              fontSize: 13,
                              fontWeight: FontWeight.bold,
                              color: AppColors.textPrimary,
                              height: 1.25,
                            ),
                          ),
                          const SizedBox(height: 2),
                          Text(
                            item.subtitle,
                            maxLines: 1,
                            overflow: TextOverflow.ellipsis,
                            style: const TextStyle(
                              fontSize: 11,
                              color: AppColors.textMuted,
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
                    color: AppColors.bgTertiary,
                    borderRadius: BorderRadius.circular(10),
                    border: Border.all(
                        color: glowColor.withAlpha(80), width: 0.8),
                  ),
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Text(
                        buttonLabel,
                        style: TextStyle(
                          fontSize: 11.5,
                          fontWeight: FontWeight.bold,
                          color: glowColor,
                        ),
                      ),
                      const SizedBox(width: 4),
                      Icon(Icons.arrow_forward_rounded,
                          size: 12, color: glowColor),
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
      SavedKind kind) {
    switch (kind) {
      case SavedKind.course:
        return (
          const LinearGradient(
            begin: Alignment.topLeft,
            end: Alignment.bottomRight,
            colors: [Color(0xFF0F261E), Color(0xFF141926)],
          ),
          const Color(0xFF34D399),
          const Color(0xFF059669),
          'COURSE',
          Icons.menu_book_rounded,
          'View Course',
        );
      case SavedKind.quiz:
        return (
          const LinearGradient(
            begin: Alignment.topLeft,
            end: Alignment.bottomRight,
            colors: [Color(0xFF1D1B36), Color(0xFF141926)],
          ),
          const Color(0xFFA78BFA),
          const Color(0xFF7C3AED),
          'OFFLINE QUIZ',
          Icons.quiz_rounded,
          'Start Quiz',
        );
      case SavedKind.exam:
        return (
          const LinearGradient(
            begin: Alignment.topLeft,
            end: Alignment.bottomRight,
            colors: [Color(0xFF0E2232), Color(0xFF141926)],
          ),
          const Color(0xFF38BDF8),
          const Color(0xFF0284C7),
          'OFFLINE EXAM',
          Icons.assignment_rounded,
          'Start Exam',
        );
    }
  }
}

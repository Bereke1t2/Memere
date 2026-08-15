import 'package:flutter/material.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';

/// Shows an interactive bottom sheet containing the full question navigation grid (60 questions).
void showExamQuestionPaletteSheet({
  required BuildContext context,
  required int count,
  required int currentIndex,
  required Set<int> answeredIndexes,
  required ValueChanged<int> onSelected,
}) {
  showModalBottomSheet(
    context: context,
    isScrollControlled: true,
    backgroundColor: const Color(0xFF0F0F14),
    shape: const RoundedRectangleBorder(
      borderRadius: BorderRadius.vertical(top: Radius.circular(24)),
      side: BorderSide(color: Color(0xFF22222C), width: 1),
    ),
    builder: (sheetContext) => DraggableScrollableSheet(
      initialChildSize: 0.65,
      minChildSize: 0.45,
      maxChildSize: 0.90,
      expand: false,
      builder: (_, scrollController) => _PaletteSheetContent(
        scrollController: scrollController,
        count: count,
        currentIndex: currentIndex,
        answeredIndexes: answeredIndexes,
        onSelected: (idx) {
          Navigator.of(sheetContext).pop();
          onSelected(idx);
        },
      ),
    ),
  );
}

class _PaletteSheetContent extends StatefulWidget {
  const _PaletteSheetContent({
    required this.scrollController,
    required this.count,
    required this.currentIndex,
    required this.answeredIndexes,
    required this.onSelected,
  });

  final ScrollController scrollController;
  final int count;
  final int currentIndex;
  final Set<int> answeredIndexes;
  final ValueChanged<int> onSelected;

  @override
  State<_PaletteSheetContent> createState() => _PaletteSheetContentState();
}

class _PaletteSheetContentState extends State<_PaletteSheetContent> {
  int _filterIndex = 0; // 0: All, 1: Answered, 2: Unanswered

  @override
  Widget build(BuildContext context) {
    final answeredCount = widget.answeredIndexes.length;
    final unansweredCount = widget.count - answeredCount;

    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 18),
      child: Column(
        children: [
          const SizedBox(height: 12),
          // Drag handle
          Center(
            child: Container(
              width: 40,
              height: 4,
              decoration: BoxDecoration(
                color: const Color(0xFF3F3F46),
                borderRadius: BorderRadius.circular(2),
              ),
            ),
          ),
          const SizedBox(height: 16),
          // Header Row
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text(
                'Question Navigator',
                style: AppTextStyles.titleMedium.copyWith(
                  fontWeight: FontWeight.w700,
                  color: Colors.white,
                ),
              ),
              IconButton(
                onPressed: () => Navigator.of(context).pop(),
                icon: const Icon(Icons.close_rounded, color: Color(0xFF71717A)),
                padding: EdgeInsets.zero,
                constraints: const BoxConstraints(),
              ),
            ],
          ),
          const SizedBox(height: 14),
          // Stats Row
          Row(
            children: [
              Expanded(
                child: Container(
                  padding: const EdgeInsets.symmetric(vertical: 10, horizontal: 12),
                  decoration: BoxDecoration(
                    color: const Color(0x1810B981),
                    borderRadius: BorderRadius.circular(12),
                    border: Border.all(color: const Color(0xFF10B981).withAlpha(80)),
                  ),
                  child: Row(
                    children: [
                      const Icon(Icons.check_circle_rounded, size: 16, color: Color(0xFF10B981)),
                      const SizedBox(width: 8),
                      Text(
                        '$answeredCount Answered',
                        style: AppTextStyles.labelMedium.copyWith(
                          color: const Color(0xFF10B981),
                          fontWeight: FontWeight.w700,
                        ),
                      ),
                    ],
                  ),
                ),
              ),
              const SizedBox(width: 10),
              Expanded(
                child: Container(
                  padding: const EdgeInsets.symmetric(vertical: 10, horizontal: 12),
                  decoration: BoxDecoration(
                    color: const Color(0xFF1A1A22),
                    borderRadius: BorderRadius.circular(12),
                    border: Border.all(color: const Color(0xFF2C2C38)),
                  ),
                  child: Row(
                    children: [
                      const Icon(Icons.radio_button_unchecked_rounded, size: 16, color: Color(0xFFA1A1AA)),
                      const SizedBox(width: 8),
                      Text(
                        '$unansweredCount Left',
                        style: AppTextStyles.labelMedium.copyWith(
                          color: const Color(0xFFA1A1AA),
                          fontWeight: FontWeight.w700,
                        ),
                      ),
                    ],
                  ),
                ),
              ),
            ],
          ),
          const SizedBox(height: 14),
          // Filter Pills
          Row(
            children: [
              _buildFilterPill('All (${widget.count})', 0),
              const SizedBox(width: 8),
              _buildFilterPill('Answered ($answeredCount)', 1),
              const SizedBox(width: 8),
              _buildFilterPill('Unanswered ($unansweredCount)', 2),
            ],
          ),
          const SizedBox(height: 16),
          // Grid
          Expanded(
            child: GridView.builder(
              controller: widget.scrollController,
              padding: const EdgeInsets.only(bottom: 24),
              gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
                crossAxisCount: 6,
                crossAxisSpacing: 10,
                mainAxisSpacing: 10,
                childAspectRatio: 1.0,
              ),
              itemCount: widget.count,
              itemBuilder: (context, index) {
                final isAnswered = widget.answeredIndexes.contains(index);
                final isCurrent = index == widget.currentIndex;

                // Filter logic
                if (_filterIndex == 1 && !isAnswered) return const SizedBox.shrink();
                if (_filterIndex == 2 && isAnswered) return const SizedBox.shrink();

                Color bgColor;
                Color borderColor;
                Color textColor;

                if (isCurrent) {
                  bgColor = const Color(0xFF0284C7); // Azure blue
                  borderColor = const Color(0xFF38BDF8);
                  textColor = Colors.white;
                } else if (isAnswered) {
                  bgColor = const Color(0xFF064E3B); // Dark emerald
                  borderColor = const Color(0xFF10B981);
                  textColor = const Color(0xFF34D399);
                } else {
                  bgColor = const Color(0xFF16161D);
                  borderColor = const Color(0xFF2A2A36);
                  textColor = const Color(0xFFA1A1AA);
                }

                return Material(
                  color: Colors.transparent,
                  child: InkWell(
                    onTap: () => widget.onSelected(index),
                    borderRadius: BorderRadius.circular(12),
                    child: Container(
                      alignment: Alignment.center,
                      decoration: BoxDecoration(
                        color: bgColor,
                        borderRadius: BorderRadius.circular(12),
                        border: Border.all(
                          color: borderColor,
                          width: isCurrent ? 2 : 1,
                        ),
                        boxShadow: isCurrent
                            ? [
                                BoxShadow(
                                  color: const Color(0xFF38BDF8).withAlpha(80),
                                  blurRadius: 8,
                                ),
                              ]
                            : null,
                      ),
                      child: Text(
                        '${index + 1}',
                        style: AppTextStyles.labelMedium.copyWith(
                          color: textColor,
                          fontWeight: isCurrent || isAnswered
                              ? FontWeight.w700
                              : FontWeight.w500,
                        ),
                      ),
                    ),
                  ),
                );
              },
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildFilterPill(String label, int index) {
    final selected = _filterIndex == index;
    return GestureDetector(
      onTap: () => setState(() => _filterIndex = index),
      child: Container(
        padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
        decoration: BoxDecoration(
          color: selected ? Colors.white : const Color(0xFF181820),
          borderRadius: BorderRadius.circular(16),
          border: Border.all(
            color: selected ? Colors.white : const Color(0xFF2A2A36),
          ),
        ),
        child: Text(
          label,
          style: AppTextStyles.labelSmall.copyWith(
            color: selected ? Colors.black : const Color(0xFFA1A1AA),
            fontWeight: selected ? FontWeight.w700 : FontWeight.w500,
          ),
        ),
      ),
    );
  }
}

/// Fallback palette widget for inline embedding if needed.
class ExamQuestionPalette extends StatelessWidget {
  const ExamQuestionPalette({
    super.key,
    required this.count,
    required this.currentIndex,
    required this.answeredIndexes,
    required this.onSelected,
  });

  final int count;
  final int currentIndex;
  final Set<int> answeredIndexes;
  final ValueChanged<int> onSelected;

  @override
  Widget build(BuildContext context) {
    return Wrap(
      spacing: AppSizes.sm,
      runSpacing: AppSizes.sm,
      children: List.generate(count, (index) {
        final current = index == currentIndex;
        final answered = answeredIndexes.contains(index);
        final fill = current
            ? AppColors.accentPrimary
            : answered
                ? AppColors.success
                : AppColors.bgTertiary;
        final borderColor = current
            ? AppColors.accentPrimary
            : answered
                ? AppColors.success
                : AppColors.border;
        return SizedBox(
          width: 36,
          height: 36,
          child: InkWell(
            onTap: () => onSelected(index),
            borderRadius: BorderRadius.circular(AppSizes.radiusFull),
            child: Container(
              alignment: Alignment.center,
              decoration: BoxDecoration(
                color: fill,
                borderRadius: BorderRadius.circular(AppSizes.radiusFull),
                border: Border.all(color: borderColor),
              ),
              child: Text(
                '${index + 1}',
                style: AppTextStyles.labelMedium.copyWith(
                  color: current || answered
                      ? AppColors.textInverse
                      : AppColors.textSecondary,
                ),
              ),
            ),
          ),
        );
      }),
    );
  }
}

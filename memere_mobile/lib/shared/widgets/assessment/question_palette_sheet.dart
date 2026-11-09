import 'package:flutter/material.dart';

import '../../../core/constants/app_colors.dart';
import '../../../core/constants/app_sizes.dart';
import '../../../core/constants/app_text_styles.dart';
import 'question_status.dart';

/// The whole paper at a glance, as a sheet rather than a strip in the page.
///
/// A 100-question exam puts a hundred cells on screen. Inline, that pushes the
/// question itself off the top and makes every scroll pass through the grid; in
/// a sheet it is one tap away and uses the full height when it is wanted.
class QuestionPaletteSheet extends StatelessWidget {
  const QuestionPaletteSheet({
    super.key,
    required this.statuses,
    required this.currentIndex,
    required this.onSelected,
  });

  /// Status per question, indexed in question order.
  final List<QuestionStatus> statuses;
  final int currentIndex;
  final ValueChanged<int> onSelected;

  static Future<void> show(
    BuildContext context, {
    required List<QuestionStatus> statuses,
    required int currentIndex,
    required ValueChanged<int> onSelected,
  }) {
    return showModalBottomSheet<void>(
      context: context,
      backgroundColor: AppColors.bgSecondary,
      isScrollControlled: true,
      showDragHandle: true,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(AppSizes.radiusLg)),
      ),
      builder: (context) => QuestionPaletteSheet(
        statuses: statuses,
        currentIndex: currentIndex,
        onSelected: (index) {
          Navigator.of(context).pop();
          onSelected(index);
        },
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final answered = statuses.where((s) => s.isAnswered).length;
    final flagged = statuses.where((s) => s.isFlagged).length;
    final blank = statuses.length - answered;

    return SafeArea(
      child: ConstrainedBox(
        constraints: BoxConstraints(
          maxHeight: MediaQuery.of(context).size.height * 0.78,
        ),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Padding(
              padding: const EdgeInsets.fromLTRB(
                AppSizes.screenPaddingH,
                0,
                AppSizes.screenPaddingH,
                AppSizes.md,
              ),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const Text('All questions', style: AppTextStyles.headlineSmall),
                  const SizedBox(height: AppSizes.sm),
                  Wrap(
                    spacing: AppSizes.md,
                    runSpacing: AppSizes.sm,
                    children: [
                      _Legend(
                        color: AppColors.accentPrimary,
                        label: '$answered answered',
                      ),
                      _Legend(
                        color: AppColors.warning,
                        label: '$flagged flagged',
                      ),
                      _Legend(
                        color: AppColors.borderStrong,
                        label: '$blank left',
                        emphasise: blank > 0,
                      ),
                    ],
                  ),
                ],
              ),
            ),
            const Divider(height: 1, color: AppColors.border),
            Flexible(
              child: SingleChildScrollView(
                padding: const EdgeInsets.all(AppSizes.screenPaddingH),
                child: Wrap(
                  spacing: AppSizes.sm,
                  runSpacing: AppSizes.sm,
                  children: List.generate(statuses.length, (index) {
                    return _PaletteCell(
                      number: index + 1,
                      status: statuses[index],
                      isCurrent: index == currentIndex,
                      onTap: () => onSelected(index),
                    );
                  }),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _PaletteCell extends StatelessWidget {
  const _PaletteCell({
    required this.number,
    required this.status,
    required this.isCurrent,
    required this.onTap,
  });

  final int number;
  final QuestionStatus status;
  final bool isCurrent;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    return Semantics(
      button: true,
      selected: isCurrent,
      label: 'Question $number, ${status.label}',
      child: SizedBox(
        // 44dp is the smallest comfortable touch target; the old 36dp grid was
        // easy to mis-tap, which on an exam means losing your place.
        width: 44,
        height: 44,
        child: Material(
          color: status.fill,
          borderRadius: BorderRadius.circular(AppSizes.radiusSm),
          child: InkWell(
            onTap: onTap,
            borderRadius: BorderRadius.circular(AppSizes.radiusSm),
            child: Container(
              decoration: BoxDecoration(
                borderRadius: BorderRadius.circular(AppSizes.radiusSm),
                border: Border.all(
                  // The question you are on is outlined in white so it is
                  // findable without stealing the accent from "answered".
                  color: isCurrent ? AppColors.textPrimary : status.border,
                  width: isCurrent ? 2 : 1,
                ),
              ),
              // The stack must fill the cell, not shrink to the number, or a
              // full-width marker would only span the width of the digits.
              child: Stack(
                fit: StackFit.expand,
                alignment: Alignment.center,
                children: [
                  Center(
                    child: Text(
                      '$number',
                      style: AppTextStyles.labelMedium.copyWith(
                        color: status.foreground,
                        fontWeight: isCurrent ? FontWeight.w800 : FontWeight.w600,
                      ),
                    ),
                  ),
                  // A bar along the top edge rather than a corner icon: a
                  // 44dp cell holding a three-digit number has no corner free,
                  // and an icon there sits on top of the digits.
                  if (status.isFlagged)
                    Positioned(
                      top: 0,
                      left: 0,
                      right: 0,
                      child: Container(
                        height: 4,
                        decoration: const BoxDecoration(
                          color: AppColors.warning,
                          borderRadius: BorderRadius.vertical(
                            top: Radius.circular(AppSizes.radiusSm - 1),
                          ),
                        ),
                      ),
                    ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}

class _Legend extends StatelessWidget {
  const _Legend({
    required this.color,
    required this.label,
    this.emphasise = false,
  });

  final Color color;
  final String label;
  final bool emphasise;

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        Container(
          width: 10,
          height: 10,
          decoration: BoxDecoration(
            color: color,
            borderRadius: BorderRadius.circular(3),
          ),
        ),
        const SizedBox(width: 6),
        Text(
          label,
          style: AppTextStyles.bodySmall.copyWith(
            color: emphasise ? AppColors.warning : AppColors.textSecondary,
            fontWeight: emphasise ? FontWeight.w700 : FontWeight.w500,
          ),
        ),
      ],
    );
  }
}

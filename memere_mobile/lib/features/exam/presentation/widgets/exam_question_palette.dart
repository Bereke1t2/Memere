import 'package:flutter/material.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';

/// Grid of question numbers. Shows current / answered / unanswered state only —
/// never correct or wrong, which is result-only.
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
                      ? Colors.white
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

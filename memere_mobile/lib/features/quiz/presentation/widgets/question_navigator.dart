import 'package:flutter/material.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';

class QuestionNavigator extends StatelessWidget {
  const QuestionNavigator({
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
        final color = current
            ? AppColors.accentPrimary
            : answered
                ? AppColors.success
                : AppColors.bgTertiary;
        return SizedBox(
          width: 36,
          height: 36,
          child: InkWell(
            onTap: () => onSelected(index),
            borderRadius: BorderRadius.circular(AppSizes.radiusFull),
            child: Container(
              alignment: Alignment.center,
              decoration: BoxDecoration(
                color: color.withAlpha(current || answered ? 255 : 255),
                borderRadius: BorderRadius.circular(AppSizes.radiusFull),
                border: Border.all(
                  color: current
                      ? AppColors.accentPrimary
                      : answered
                          ? AppColors.success
                          : AppColors.border,
                ),
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

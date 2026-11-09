import 'package:flutter/material.dart';

import '../../../core/constants/app_colors.dart';
import '../../../core/constants/app_sizes.dart';
import '../../../core/constants/app_text_styles.dart';

/// Which reviewed questions a student wants to see.
enum ReviewFilter {
  all,
  correct,
  missed;

  bool matches(bool isCorrect) => switch (this) {
        ReviewFilter.all => true,
        ReviewFilter.correct => isCorrect,
        ReviewFilter.missed => !isCorrect,
      };
}

/// Segmented filter over a review list.
///
/// After a 100-question paper the question a student actually has is "what did
/// I get wrong", and scrolling past ninety correct answers to find out is the
/// whole friction. Counts are on the tabs so the answer is visible before
/// anything is tapped.
class ReviewFilterTabs extends StatelessWidget {
  const ReviewFilterTabs({
    super.key,
    required this.selected,
    required this.total,
    required this.correct,
    required this.onChanged,
  });

  final ReviewFilter selected;
  final int total;
  final int correct;
  final ValueChanged<ReviewFilter> onChanged;

  int get missed => total - correct;

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(4),
      decoration: BoxDecoration(
        color: AppColors.bgSecondary,
        borderRadius: BorderRadius.circular(AppSizes.radiusMd),
        border: Border.all(color: AppColors.borderStrong),
      ),
      child: Row(
        children: [
          _Tab(
            label: 'All',
            count: total,
            isSelected: selected == ReviewFilter.all,
            onTap: () => onChanged(ReviewFilter.all),
          ),
          _Tab(
            label: 'Correct',
            count: correct,
            isSelected: selected == ReviewFilter.correct,
            accent: AppColors.accentPrimary,
            onTap: () => onChanged(ReviewFilter.correct),
          ),
          _Tab(
            label: 'Missed',
            count: missed,
            isSelected: selected == ReviewFilter.missed,
            accent: AppColors.error,
            onTap: () => onChanged(ReviewFilter.missed),
          ),
        ],
      ),
    );
  }
}

class _Tab extends StatelessWidget {
  const _Tab({
    required this.label,
    required this.count,
    required this.isSelected,
    required this.onTap,
    this.accent,
  });

  final String label;
  final int count;
  final bool isSelected;
  final VoidCallback onTap;
  final Color? accent;

  @override
  Widget build(BuildContext context) {
    final colour = accent ?? AppColors.textPrimary;
    return Expanded(
      child: Semantics(
        button: true,
        selected: isSelected,
        label: '$label, $count questions',
        child: InkWell(
          onTap: onTap,
          borderRadius: BorderRadius.circular(AppSizes.radiusSm),
          child: Container(
            padding: const EdgeInsets.symmetric(vertical: 9),
            alignment: Alignment.center,
            decoration: BoxDecoration(
              color: isSelected ? AppColors.bgQuaternary : Colors.transparent,
              borderRadius: BorderRadius.circular(AppSizes.radiusSm),
            ),
            child: Text(
              '$label ($count)',
              maxLines: 1,
              overflow: TextOverflow.ellipsis,
              style: AppTextStyles.labelSmall.copyWith(
                color: isSelected ? colour : AppColors.textMuted,
                fontWeight: isSelected ? FontWeight.w700 : FontWeight.w500,
              ),
            ),
          ),
        ),
      ),
    );
  }
}

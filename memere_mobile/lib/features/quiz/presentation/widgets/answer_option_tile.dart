import 'package:flutter/material.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../domain/entities/quiz_answer_entity.dart';

class AnswerOptionTile extends StatelessWidget {
  const AnswerOptionTile({
    super.key,
    required this.answer,
    required this.isSelected,
    required this.onTap,
    this.isMultiSelect = false,
  });

  final QuizAnswerEntity answer;
  final bool isSelected;
  final VoidCallback onTap;
  final bool isMultiSelect;

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.only(bottom: AppSizes.sm),
      child: Material(
        color: isSelected ? AppColors.accentGlow : AppColors.bgTertiary,
        borderRadius: BorderRadius.circular(AppSizes.radiusMd),
        child: InkWell(
          onTap: onTap,
          borderRadius: BorderRadius.circular(AppSizes.radiusMd),
          child: Container(
            padding: const EdgeInsets.all(AppSizes.md),
            decoration: BoxDecoration(
              borderRadius: BorderRadius.circular(AppSizes.radiusMd),
              border: Border.all(
                color: isSelected ? AppColors.accentPrimary : AppColors.border,
              ),
            ),
            child: Row(
              children: [
                Icon(
                  isMultiSelect
                      ? (isSelected
                          ? Icons.check_box_rounded
                          : Icons.check_box_outline_blank_rounded)
                      : (isSelected
                          ? Icons.radio_button_checked_rounded
                          : Icons.radio_button_unchecked_rounded),
                  color: isSelected
                      ? AppColors.accentPrimary
                      : AppColors.textSecondary,
                  size: AppSizes.iconSm,
                ),
                const SizedBox(width: AppSizes.md),
                Expanded(
                  child: Text(
                    answer.text,
                    style: AppTextStyles.bodyMedium,
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}

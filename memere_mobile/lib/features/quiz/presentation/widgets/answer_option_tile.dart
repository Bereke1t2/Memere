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
    this.optionIndex = 0,
    this.isMultiSelect = false,
  });

  final QuizAnswerEntity answer;
  final bool isSelected;
  final VoidCallback onTap;
  final int optionIndex;
  final bool isMultiSelect;

  String get _optionLetter {
    if (optionIndex >= 0 && optionIndex < 26) {
      return String.fromCharCode(65 + optionIndex);
    }
    return '${optionIndex + 1}';
  }

  @override
  Widget build(BuildContext context) {
    const emeraldColor = Color(0xFF10B981);
    final borderColor = isSelected ? emeraldColor : AppColors.borderStrong.withAlpha(100);
    final bgColor = isSelected ? const Color(0x1E10B981) : AppColors.bgSecondary;

    return Padding(
      padding: const EdgeInsets.only(bottom: AppSizes.sm + 2),
      child: AnimatedContainer(
        duration: const Duration(milliseconds: 180),
        curve: Curves.easeInOut,
        decoration: BoxDecoration(
          color: bgColor,
          borderRadius: BorderRadius.circular(16),
          border: Border.all(
            color: borderColor,
            width: isSelected ? 2.0 : 1.0,
          ),
          boxShadow: isSelected
              ? [
                  BoxShadow(
                    color: emeraldColor.withValues(alpha: 0.15),
                    blurRadius: 10,
                    offset: const Offset(0, 4),
                  ),
                ]
              : [],
        ),
        child: InkWell(
          onTap: onTap,
          borderRadius: BorderRadius.circular(16),
          child: Padding(
            padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
            child: Row(
              children: [
                AnimatedContainer(
                  duration: const Duration(milliseconds: 180),
                  width: 34,
                  height: 34,
                  alignment: Alignment.center,
                  decoration: BoxDecoration(
                    color: isSelected ? emeraldColor : AppColors.bgTertiary,
                    shape: BoxShape.circle,
                    border: Border.all(
                      color: isSelected ? emeraldColor : AppColors.borderStrong,
                      width: 1.5,
                    ),
                  ),
                  child: Text(
                    _optionLetter,
                    style: TextStyle(
                      fontSize: 14,
                      fontWeight: FontWeight.w700,
                      color: isSelected ? Colors.white : AppColors.textSecondary,
                    ),
                  ),
                ),
                const SizedBox(width: 14),
                Expanded(
                  child: Text(
                    answer.text,
                    style: AppTextStyles.bodyMedium.copyWith(
                      color: isSelected ? Colors.white : AppColors.textPrimary,
                      fontWeight: isSelected ? FontWeight.w600 : FontWeight.w400,
                      height: 1.35,
                    ),
                  ),
                ),
                const SizedBox(width: 10),
                Icon(
                  isMultiSelect
                      ? (isSelected
                          ? Icons.check_box_rounded
                          : Icons.check_box_outline_blank_rounded)
                      : (isSelected
                          ? Icons.check_circle_rounded
                          : Icons.radio_button_unchecked_rounded),
                  color: isSelected ? emeraldColor : const Color(0xFF64748B),
                  size: 20,
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}

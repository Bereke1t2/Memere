import 'package:flutter/material.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../domain/entities/exam_answer_entity.dart';

/// Modern, luxury answer tile for live exam attempts with A/B/C/D letter badges,
/// glowing selection state, and tactile feedback.
class ExamAnswerOptionTile extends StatelessWidget {
  const ExamAnswerOptionTile({
    super.key,
    required this.answer,
    required this.isSelected,
    required this.onTap,
    this.index = 0,
    this.isMultiSelect = false,
  });

  final ExamAnswerEntity answer;
  final bool isSelected;
  final VoidCallback onTap;
  final int index;
  final bool isMultiSelect;

  String get _letter {
    final letters = ['A', 'B', 'C', 'D', 'E', 'F', 'G', 'H'];
    if (index >= 0 && index < letters.length) {
      return letters[index];
    }
    return '${index + 1}';
  }

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 12),
      child: AnimatedContainer(
        duration: const Duration(milliseconds: 180),
        curve: Curves.easeInOut,
        decoration: BoxDecoration(
          color: isSelected
              ? const Color(0xFF062A1F) // Deep luminous emerald tint
              : const Color(0xFF131318), // Clean obsidian surface
          borderRadius: BorderRadius.circular(14),
          border: Border.all(
            color: isSelected
                ? AppColors.brandEmerald
                : const Color(0xFF22222B),
            width: isSelected ? 1.8 : 1.0,
          ),
          boxShadow: isSelected
              ? [
                  BoxShadow(
                    color: AppColors.brandEmerald.withAlpha(40),
                    blurRadius: 12,
                    offset: const Offset(0, 3),
                  ),
                ]
              : null,
        ),
        child: Material(
          color: Colors.transparent,
          child: InkWell(
            onTap: onTap,
            borderRadius: BorderRadius.circular(14),
            splashColor: AppColors.brandEmerald.withAlpha(30),
            highlightColor: AppColors.brandEmerald.withAlpha(15),
            child: Padding(
              padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 14),
              child: Row(
                crossAxisAlignment: CrossAxisAlignment.center,
                children: [
                  // Letter Badge (A, B, C, D)
                  AnimatedContainer(
                    duration: const Duration(milliseconds: 180),
                    width: 32,
                    height: 32,
                    alignment: Alignment.center,
                    decoration: BoxDecoration(
                      color: isSelected
                          ? AppColors.brandEmerald
                          : const Color(0xFF1E1E26),
                      borderRadius: BorderRadius.circular(8),
                      border: Border.all(
                        color: isSelected
                            ? AppColors.brandEmerald
                            : const Color(0xFF32323E),
                        width: 1,
                      ),
                    ),
                    child: isSelected && isMultiSelect
                        ? const Icon(
                            Icons.check_rounded,
                            size: 18,
                            color: Colors.black,
                          )
                        : Text(
                            _letter,
                            style: AppTextStyles.labelMedium.copyWith(
                              color: isSelected ? Colors.black : AppColors.textSecondary,
                              fontWeight: FontWeight.w700,
                            ),
                          ),
                  ),
                  const SizedBox(width: 14),
                  // Answer Text
                  Expanded(
                    child: Text(
                      answer.text,
                      style: AppTextStyles.bodyMedium.copyWith(
                        color: isSelected ? Colors.white : const Color(0xFFE4E4E7),
                        fontWeight: isSelected ? FontWeight.w600 : FontWeight.w400,
                        height: 1.4,
                      ),
                    ),
                  ),
                  const SizedBox(width: 8),
                  // Selection Indicator Icon
                  if (isSelected)
                    const Icon(
                      Icons.check_circle_rounded,
                      size: 20,
                      color: AppColors.brandEmerald,
                    )
                  else
                    Container(
                      width: 18,
                      height: 18,
                      decoration: BoxDecoration(
                        shape: BoxShape.circle,
                        border: Border.all(
                          color: const Color(0xFF383845),
                          width: 1.5,
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

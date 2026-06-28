import 'package:flutter/material.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';

class SubjectFilterChips extends StatelessWidget {
  const SubjectFilterChips({
    super.key,
    required this.subjects,
    required this.selectedSubject,
    required this.onSelected,
  });

  final List<String> subjects;
  final String? selectedSubject;
  final ValueChanged<String?> onSelected;

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      height: 42,
      child: ListView.separated(
        scrollDirection: Axis.horizontal,
        padding:
            const EdgeInsets.symmetric(horizontal: AppSizes.screenPaddingH),
        itemCount: subjects.length + 1,
        separatorBuilder: (_, __) => const SizedBox(width: AppSizes.sm),
        itemBuilder: (context, index) {
          final subject = index == 0 ? null : subjects[index - 1];
          final label = subject ?? 'All';
          final selected = selectedSubject == subject;
          return ChoiceChip(
            label: Text(label),
            selected: selected,
            onSelected: (_) => onSelected(subject),
            labelStyle: AppTextStyles.labelMedium.copyWith(
              color: selected ? Colors.white : AppColors.textSecondary,
            ),
            selectedColor: AppColors.accentPrimary,
            backgroundColor: AppColors.bgTertiary,
            side: BorderSide(
              color: selected ? AppColors.accentPrimary : AppColors.border,
            ),
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(AppSizes.radiusFull),
            ),
          );
        },
      ),
    );
  }
}

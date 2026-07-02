import 'package:flutter/material.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_motion.dart';
import '../../../../core/constants/app_shadows.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../../../shared/widgets/app_surface.dart';

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
      height: 38,
      child: ListView.separated(
        scrollDirection: Axis.horizontal,
        padding:
            const EdgeInsets.symmetric(horizontal: AppSizes.screenPaddingH),
        itemCount: subjects.length + 1,
        separatorBuilder: (_, __) => const SizedBox(width: AppSizes.sm),
        itemBuilder: (context, index) {
          final subject = index == 0 ? null : subjects[index - 1];
          final selected = selectedSubject == subject;
          return AppPressable(
            onTap: () => onSelected(subject),
            borderRadius: AppSizes.radiusFull,
            child: AnimatedContainer(
              duration: AppMotion.base,
              curve: AppMotion.standard,
              height: 36,
              padding: EdgeInsets.only(
                left: selected ? AppSizes.sm : AppSizes.md,
                right: AppSizes.md,
              ),
              decoration: BoxDecoration(
                color: selected ? AppColors.textPrimary : AppColors.bgSecondary,
                borderRadius: BorderRadius.circular(AppSizes.radiusFull),
                border: Border.all(
                  color: selected ? AppColors.textPrimary : AppColors.border,
                ),
                boxShadow: selected ? AppShadows.accentGlow : AppShadows.sm,
              ),
              child: Row(
                mainAxisSize: MainAxisSize.min,
                children: [
                  AnimatedContainer(
                    duration: AppMotion.base,
                    curve: AppMotion.standard,
                    width: selected ? 24 : 0,
                    height: 24,
                    margin: EdgeInsets.only(
                      right: selected ? AppSizes.xs : 0,
                    ),
                    decoration: BoxDecoration(
                      color: AppColors.textInverse.withAlpha(26),
                      shape: BoxShape.circle,
                    ),
                    child: selected
                        ? const Icon(
                            Icons.check_rounded,
                            size: AppSizes.iconXs,
                            color: AppColors.textInverse,
                          )
                        : null,
                  ),
                  Text(
                    subject ?? 'All',
                    style: AppTextStyles.labelMedium.copyWith(
                      color: selected
                          ? AppColors.textInverse
                          : AppColors.textSecondary,
                    ),
                  ),
                ],
              ),
            ),
          );
        },
      ),
    );
  }
}

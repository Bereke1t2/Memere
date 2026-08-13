import 'package:flutter/material.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_motion.dart';
import '../../../../core/constants/app_sizes.dart';
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
        separatorBuilder: (_, __) => const SizedBox(width: 8),
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
                left: selected ? 10 : 16,
                right: 16,
              ),
              decoration: BoxDecoration(
                color: selected ? AppColors.brandEmerald : AppColors.bgSecondary,
                borderRadius: BorderRadius.circular(18),
                border: Border.all(
                  color: selected ? AppColors.brandEmeraldDark : AppColors.borderStrong,
                  width: selected ? 2 : 1,
                ),
              ),
              child: Row(
                mainAxisSize: MainAxisSize.min,
                children: [
                  if (selected) ...[
                    const Icon(
                      Icons.check_rounded,
                      size: 16,
                      color: Colors.white,
                    ),
                    const SizedBox(width: 6),
                  ],
                  Text(
                    subject ?? 'All Subjects',
                    style: TextStyle(
                      fontSize: 13,
                      fontWeight: FontWeight.bold,
                      color: selected ? Colors.white : AppColors.textSecondary,
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

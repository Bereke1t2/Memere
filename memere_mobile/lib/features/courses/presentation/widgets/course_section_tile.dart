import 'package:flutter/material.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_motion.dart';
import '../../../../core/constants/app_shadows.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../domain/entities/course_section_entity.dart';
import 'lesson_tile.dart';

class CourseSectionTile extends StatelessWidget {
  const CourseSectionTile({
    super.key,
    required this.section,
    required this.sectionNumber,
    this.initiallyExpanded = false,
    this.canOpenLessons = false,
  });

  final CourseSectionEntity section;
  final int sectionNumber;
  final bool initiallyExpanded;
  final bool canOpenLessons;

  @override
  Widget build(BuildContext context) {
    final numPadded = sectionNumber.toString().padLeft(2, '0');

    return AnimatedContainer(
      duration: AppMotion.base,
      curve: AppMotion.standard,
      decoration: BoxDecoration(
        color: AppColors.bgSecondary,
        borderRadius: BorderRadius.circular(16),
        border: Border.all(color: AppColors.border),
        boxShadow: AppShadows.sm,
      ),
      clipBehavior: Clip.antiAlias,
      child: Theme(
        data: Theme.of(context).copyWith(
          dividerColor: Colors.transparent,
          splashColor: Colors.transparent,
          highlightColor: Colors.transparent,
        ),
        child: ExpansionTile(
          initiallyExpanded: initiallyExpanded,
          tilePadding: const EdgeInsets.symmetric(
            horizontal: 16,
            vertical: 8,
          ),
          childrenPadding: const EdgeInsets.fromLTRB(10, 0, 10, 10),
          title: Row(
            children: [
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 3),
                decoration: BoxDecoration(
                  color: const Color(0x2238BDF8),
                  borderRadius: BorderRadius.circular(6),
                  border: Border.all(color: const Color(0x4438BDF8)),
                ),
                child: Text(
                  'SECTION $numPadded',
                  style: const TextStyle(
                    fontSize: 10,
                    fontWeight: FontWeight.w800,
                    color: Color(0xFF38BDF8),
                    letterSpacing: 0.6,
                  ),
                ),
              ),
              const SizedBox(width: 8),
              Text(
                section.lessonCountLabel,
                style: AppTextStyles.caption.copyWith(
                  color: AppColors.textMuted,
                  fontSize: 11,
                ),
              ),
            ],
          ),
          subtitle: Padding(
            padding: const EdgeInsets.only(top: 6),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  section.title,
                  style: AppTextStyles.titleMedium.copyWith(
                    fontWeight: FontWeight.bold,
                    fontSize: 15,
                  ),
                ),
                if (section.description.isNotEmpty) ...[
                  const SizedBox(height: 4),
                  Text(
                    section.description,
                    maxLines: 2,
                    overflow: TextOverflow.ellipsis,
                    style: AppTextStyles.bodySmall.copyWith(
                      color: AppColors.textSecondary,
                      height: 1.3,
                    ),
                  ),
                ],
              ],
            ),
          ),
          children: [
            if (section.lessons.isEmpty)
              const Padding(
                padding: EdgeInsets.all(AppSizes.md),
                child: Text(
                  'No lessons in this section yet.',
                  style: AppTextStyles.bodySmall,
                ),
              )
            else
              ...section.lessons.asMap().entries.map(
                    (entry) => Padding(
                      padding: const EdgeInsets.only(top: 6),
                      child: LessonTile(
                        lesson: entry.value,
                        lessonNumber: entry.key + 1,
                        canOpen: canOpenLessons || entry.value.isFreePreview,
                      ),
                    ),
                  ),
          ],
        ),
      ),
    );
  }
}

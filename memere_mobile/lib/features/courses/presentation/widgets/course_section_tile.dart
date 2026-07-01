import 'package:flutter/material.dart';

import '../../../../core/constants/app_colors.dart';
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
    return Container(
      decoration: BoxDecoration(
        color: AppColors.bgSecondary,
        borderRadius: BorderRadius.circular(AppSizes.radiusMd),
        border: Border.all(color: AppColors.hairline),
        boxShadow: AppShadows.sm,
      ),
      clipBehavior: Clip.antiAlias,
      child: Theme(
        data: Theme.of(context).copyWith(dividerColor: Colors.transparent),
        child: ExpansionTile(
          initiallyExpanded: initiallyExpanded,
          tilePadding: const EdgeInsets.symmetric(
            horizontal: AppSizes.md,
            vertical: AppSizes.xs,
          ),
          childrenPadding: const EdgeInsets.fromLTRB(
            AppSizes.sm,
            0,
            AppSizes.sm,
            AppSizes.sm,
          ),
          title: Text(
            'Section $sectionNumber',
            style: AppTextStyles.labelSmall.copyWith(
              color: AppColors.accentPrimary,
            ),
          ),
          subtitle: Padding(
            padding: const EdgeInsets.only(top: AppSizes.xs),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(section.title, style: AppTextStyles.titleMedium),
                if (section.description.isNotEmpty) ...[
                  const SizedBox(height: AppSizes.xs),
                  Text(
                    section.description,
                    maxLines: 2,
                    overflow: TextOverflow.ellipsis,
                    style: AppTextStyles.bodySmall,
                  ),
                ],
              ],
            ),
          ),
          trailing: Text(
            section.lessonCountLabel,
            style: AppTextStyles.caption,
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
                    (entry) => LessonTile(
                      lesson: entry.value,
                      lessonNumber: entry.key + 1,
                      canOpen: canOpenLessons || entry.value.isFreePreview,
                    ),
                  ),
          ],
        ),
      ),
    );
  }
}

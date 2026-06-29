import 'package:flutter/material.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../domain/entities/question_feedback_entity.dart';

class QuestionFeedbackTile extends StatelessWidget {
  const QuestionFeedbackTile({
    super.key,
    required this.feedback,
    required this.index,
  });

  final QuestionFeedbackEntity feedback;
  final int index;

  @override
  Widget build(BuildContext context) {
    final color = feedback.correct ? AppColors.success : AppColors.error;
    return Container(
      margin: const EdgeInsets.only(bottom: AppSizes.md),
      padding: const EdgeInsets.all(AppSizes.md),
      decoration: BoxDecoration(
        color: AppColors.bgSecondary,
        borderRadius: BorderRadius.circular(AppSizes.radiusMd),
        border: Border.all(color: AppColors.border),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              Icon(
                feedback.correct
                    ? Icons.check_circle_outline_rounded
                    : Icons.cancel_outlined,
                color: color,
                size: AppSizes.iconSm,
              ),
              const SizedBox(width: AppSizes.sm),
              Expanded(
                child: Text(
                  'Question ${index + 1}',
                  style: AppTextStyles.titleMedium,
                ),
              ),
              Text(
                '${feedback.pointsAwarded}/${feedback.pointsPossible}',
                style: AppTextStyles.labelMedium.copyWith(color: color),
              ),
            ],
          ),
          const SizedBox(height: AppSizes.sm),
          _IdLine(label: 'Selected', values: feedback.selectedAnswers),
          const SizedBox(height: AppSizes.xs),
          _IdLine(label: 'Correct', values: feedback.correctAnswerIds),
          if (feedback.explanation != null) ...[
            const SizedBox(height: AppSizes.sm),
            Text(
              feedback.explanation!,
              style: AppTextStyles.bodySmall.copyWith(
                color: AppColors.textSecondary,
              ),
            ),
          ],
        ],
      ),
    );
  }
}

class _IdLine extends StatelessWidget {
  const _IdLine({
    required this.label,
    required this.values,
  });

  final String label;
  final List<String> values;

  @override
  Widget build(BuildContext context) {
    return Text(
      '$label: ${values.isEmpty ? 'None' : values.join(', ')}',
      maxLines: 2,
      overflow: TextOverflow.ellipsis,
      style: AppTextStyles.caption.copyWith(color: AppColors.textSecondary),
    );
  }
}

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
    final isCorrect = feedback.correct;
    const emeraldColor = Color(0xFF10B981);
    const redColor = Color(0xFFEF4444);
    final statusColor = isCorrect ? emeraldColor : redColor;

    return Container(
      margin: const EdgeInsets.only(bottom: AppSizes.md),
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: AppColors.bgSecondary,
        borderRadius: BorderRadius.circular(16),
        border: Border.all(
          color: statusColor.withValues(alpha: 0.3),
          width: 1.2,
        ),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              Container(
                padding:
                    const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
                decoration: BoxDecoration(
                  color: statusColor.withValues(alpha: 0.15),
                  borderRadius: BorderRadius.circular(10),
                  border: Border.all(color: statusColor.withValues(alpha: 0.4)),
                ),
                child: Row(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    Icon(
                      isCorrect
                          ? Icons.check_circle_rounded
                          : Icons.cancel_rounded,
                      color: statusColor,
                      size: 15,
                    ),
                    const SizedBox(width: 6),
                    Text(
                      isCorrect ? 'Correct' : 'Incorrect',
                      style: TextStyle(
                        fontSize: 12,
                        fontWeight: FontWeight.w700,
                        color: statusColor,
                      ),
                    ),
                  ],
                ),
              ),
              const Spacer(),
              Text(
                '${feedback.pointsAwarded}/${feedback.pointsPossible} ${feedback.pointsPossible == 1 ? 'pt' : 'pts'}',
                style: TextStyle(
                  fontSize: 13,
                  fontWeight: FontWeight.w700,
                  color: statusColor,
                ),
              ),
            ],
          ),
          const SizedBox(height: 12),
          Text(
            feedback.questionText != null && feedback.questionText!.isNotEmpty
                ? 'Q${index + 1}. ${feedback.questionText!}'
                : 'Question ${index + 1}',
            style: const TextStyle(
              fontSize: 15.5,
              fontWeight: FontWeight.w700,
              color: Colors.white,
              height: 1.35,
            ),
          ),
          const SizedBox(height: 14),
          if (feedback.answers.isNotEmpty) ...[
            ...feedback.answers.asMap().entries.map((entry) {
              final idx = entry.key;
              final ans = entry.value;
              final isSelected = feedback.selectedAnswers.contains(ans.id);
              final isCorrectAns = ans.isCorrect ||
                  feedback.correctAnswerIds.contains(ans.id);

              Color optBg = AppColors.bgTertiary;
              Color optBorder = AppColors.borderStrong.withAlpha(80);
              Color textColor = AppColors.textSecondary;
              Widget? icon;

              if (isSelected && isCorrectAns) {
                optBg = const Color(0x2210B981);
                optBorder = emeraldColor;
                textColor = Colors.white;
                icon = const Icon(Icons.check_circle_rounded,
                    color: emeraldColor, size: 18);
              } else if (isSelected && !isCorrectAns) {
                optBg = const Color(0x22EF4444);
                optBorder = redColor;
                textColor = Colors.white;
                icon = const Icon(Icons.cancel_rounded,
                    color: redColor, size: 18);
              } else if (isCorrectAns) {
                optBg = const Color(0x1810B981);
                optBorder = emeraldColor.withValues(alpha: 0.6);
                textColor = emeraldColor;
                icon = const Icon(Icons.check_circle_outline_rounded,
                    color: emeraldColor, size: 18);
              }

              final letter = String.fromCharCode(65 + idx);

              return Padding(
                padding: const EdgeInsets.only(bottom: 8),
                child: Container(
                  padding:
                      const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
                  decoration: BoxDecoration(
                    color: optBg,
                    borderRadius: BorderRadius.circular(12),
                    border: Border.all(color: optBorder),
                  ),
                  child: Row(
                    children: [
                      Container(
                        width: 26,
                        height: 26,
                        alignment: Alignment.center,
                        decoration: BoxDecoration(
                          color: isSelected || isCorrectAns
                              ? optBorder
                              : AppColors.bgSecondary,
                          shape: BoxShape.circle,
                        ),
                        child: Text(
                          letter,
                          style: TextStyle(
                            fontSize: 12,
                            fontWeight: FontWeight.w700,
                            color: isSelected || isCorrectAns
                                ? Colors.white
                                : AppColors.textSecondary,
                          ),
                        ),
                      ),
                      const SizedBox(width: 10),
                      Expanded(
                        child: Text(
                          ans.text,
                          style: TextStyle(
                            fontSize: 13.5,
                            fontWeight: isSelected || isCorrectAns
                                ? FontWeight.w600
                                : FontWeight.w400,
                            color: textColor,
                          ),
                        ),
                      ),
                      if (icon != null) ...[
                        const SizedBox(width: 8),
                        icon,
                      ],
                    ],
                  ),
                ),
              );
            }),
          ] else ...[
            _IdLine(label: 'Your Answer', values: feedback.selectedAnswers),
            const SizedBox(height: 4),
            _IdLine(label: 'Correct Answer', values: feedback.correctAnswerIds),
          ],
          if (feedback.explanation != null &&
              feedback.explanation!.isNotEmpty) ...[
            const SizedBox(height: 12),
            Container(
              padding: const EdgeInsets.all(12),
              decoration: BoxDecoration(
                color: const Color(0x1806B6D4),
                borderRadius: BorderRadius.circular(12),
                border: Border.all(color: const Color(0x4406B6D4)),
              ),
              child: Row(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const Icon(
                    Icons.lightbulb_outline_rounded,
                    size: 18,
                    color: Color(0xFF38BDF8),
                  ),
                  const SizedBox(width: 8),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        const Text(
                          'Explanation',
                          style: TextStyle(
                            fontSize: 12,
                            fontWeight: FontWeight.w700,
                            color: Color(0xFF38BDF8),
                          ),
                        ),
                        const SizedBox(height: 2),
                        Text(
                          feedback.explanation!,
                          style: const TextStyle(
                            fontSize: 12.5,
                            color: Color(0xFFCBD5E1),
                            height: 1.35,
                          ),
                        ),
                      ],
                    ),
                  ),
                ],
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

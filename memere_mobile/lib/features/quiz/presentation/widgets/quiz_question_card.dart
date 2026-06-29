import 'package:flutter/material.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../domain/entities/quiz_question_entity.dart';
import 'answer_option_tile.dart';
import 'short_answer_field.dart';

class QuizQuestionCard extends StatelessWidget {
  const QuizQuestionCard({
    super.key,
    required this.question,
    required this.selectedAnswer,
    required this.onSingleAnswerSelected,
    required this.onMultiAnswerToggled,
    required this.onShortAnswerChanged,
  });

  final QuizQuestionEntity question;
  final Object? selectedAnswer;
  final void Function(String questionId, String answerId)
      onSingleAnswerSelected;
  final void Function(String questionId, String answerId) onMultiAnswerToggled;
  final void Function(String questionId, String value) onShortAnswerChanged;

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(AppSizes.md),
      decoration: BoxDecoration(
        color: AppColors.bgSecondary,
        borderRadius: BorderRadius.circular(AppSizes.radiusMd),
        border: Border.all(color: AppColors.border),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Wrap(
            spacing: AppSizes.sm,
            runSpacing: AppSizes.sm,
            children: [
              _Chip(label: '${question.points} pt'),
              _Chip(label: _questionTypeLabel(question.type)),
              if (question.subject != null) _Chip(label: question.subject!),
              if (question.topic != null) _Chip(label: question.topic!),
            ],
          ),
          const SizedBox(height: AppSizes.md),
          Text(question.text, style: AppTextStyles.titleLarge),
          const SizedBox(height: AppSizes.lg),
          if (question.type == QuizQuestionType.shortAnswer)
            ShortAnswerField(
              initialValue:
                  selectedAnswer is String ? selectedAnswer as String : '',
              onChanged: (value) => onShortAnswerChanged(question.id, value),
            )
          else
            ...question.answers.map(
              (answer) {
                final multi = question.type == QuizQuestionType.multipleSelect;
                final selected = multi
                    ? selectedAnswer is List &&
                        (selectedAnswer as List).contains(answer.id)
                    : selectedAnswer == answer.id;
                return AnswerOptionTile(
                  answer: answer,
                  isSelected: selected,
                  isMultiSelect: multi,
                  onTap: () {
                    if (multi) {
                      onMultiAnswerToggled(question.id, answer.id);
                    } else {
                      onSingleAnswerSelected(question.id, answer.id);
                    }
                  },
                );
              },
            ),
        ],
      ),
    );
  }
}

class _Chip extends StatelessWidget {
  const _Chip({required this.label});

  final String label;

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(
        horizontal: AppSizes.sm,
        vertical: AppSizes.xs,
      ),
      decoration: BoxDecoration(
        color: AppColors.bgTertiary,
        borderRadius: BorderRadius.circular(AppSizes.radiusFull),
        border: Border.all(color: AppColors.border),
      ),
      child: Text(label, style: AppTextStyles.labelSmall),
    );
  }
}

String _questionTypeLabel(QuizQuestionType type) {
  switch (type) {
    case QuizQuestionType.multipleChoice:
      return 'Single choice';
    case QuizQuestionType.multipleSelect:
      return 'Multi select';
    case QuizQuestionType.trueFalse:
      return 'True/false';
    case QuizQuestionType.shortAnswer:
      return 'Short answer';
  }
}

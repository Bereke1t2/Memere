import 'package:flutter/material.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../domain/entities/exam_question_entity.dart';
import 'exam_answer_option_tile.dart';

class ExamQuestionCard extends StatelessWidget {
  const ExamQuestionCard({
    super.key,
    required this.question,
    required this.selectedAnswer,
    required this.onSingleAnswerSelected,
    required this.onMultiAnswerToggled,
    required this.onShortAnswerChanged,
  });

  final ExamQuestionEntity question;
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
              _Chip(label: '${question.marks} mark${question.marks == 1 ? '' : 's'}'),
              _Chip(label: _questionTypeLabel(question.type)),
              if (question.subject != null) _Chip(label: question.subject!),
              if (question.topic != null) _Chip(label: question.topic!),
            ],
          ),
          const SizedBox(height: AppSizes.md),
          Text(question.text, style: AppTextStyles.titleLarge),
          const SizedBox(height: AppSizes.lg),
          if (question.type == ExamQuestionType.shortAnswer)
            _ShortAnswerField(
              key: ValueKey('short-${question.questionId}'),
              initialValue:
                  selectedAnswer is String ? selectedAnswer as String : '',
              onChanged: (value) =>
                  onShortAnswerChanged(question.questionId, value),
            )
          else
            ...question.answers.map((answer) {
              final multi = question.type == ExamQuestionType.multipleSelect;
              final selected = multi
                  ? selectedAnswer is List &&
                      (selectedAnswer as List).contains(answer.id)
                  : selectedAnswer == answer.id;
              return ExamAnswerOptionTile(
                answer: answer,
                isSelected: selected,
                isMultiSelect: multi,
                onTap: () {
                  if (multi) {
                    onMultiAnswerToggled(question.questionId, answer.id);
                  } else {
                    onSingleAnswerSelected(question.questionId, answer.id);
                  }
                },
              );
            }),
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

class _ShortAnswerField extends StatefulWidget {
  const _ShortAnswerField({
    super.key,
    required this.initialValue,
    required this.onChanged,
  });

  final String initialValue;
  final ValueChanged<String> onChanged;

  @override
  State<_ShortAnswerField> createState() => _ShortAnswerFieldState();
}

class _ShortAnswerFieldState extends State<_ShortAnswerField> {
  late final TextEditingController _controller;

  @override
  void initState() {
    super.initState();
    _controller = TextEditingController(text: widget.initialValue);
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return TextField(
      controller: _controller,
      minLines: 3,
      maxLines: 5,
      onChanged: widget.onChanged,
      style: AppTextStyles.bodyMedium,
      decoration: InputDecoration(
        hintText: 'Type your answer',
        hintStyle: AppTextStyles.bodyMedium.copyWith(
          color: AppColors.textDisabled,
        ),
      ),
    );
  }
}

String _questionTypeLabel(ExamQuestionType type) {
  switch (type) {
    case ExamQuestionType.multipleChoice:
      return 'Single choice';
    case ExamQuestionType.multipleSelect:
      return 'Multi select';
    case ExamQuestionType.trueFalse:
      return 'True/false';
    case ExamQuestionType.shortAnswer:
      return 'Short answer';
  }
}

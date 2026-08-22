import 'package:flutter/material.dart';

import '../../../../core/constants/app_colors.dart';
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
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(
        color: AppColors.bgSecondary,
        borderRadius: BorderRadius.circular(20),
        border: Border.all(
          color: AppColors.borderStrong.withAlpha(90),
        ),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              _Chip(
                label: '${question.points} ${question.points == 1 ? 'Point' : 'Points'}',
                isHighlight: true,
              ),
              const SizedBox(width: 8),
              _Chip(label: _questionTypeLabel(question.type)),
              if (question.subject != null) ...[
                const SizedBox(width: 8),
                _Chip(label: question.subject!),
              ],
            ],
          ),
          const SizedBox(height: 16),
          Text(
            question.text,
            style: const TextStyle(
              fontSize: 17,
              fontWeight: FontWeight.w700,
              color: Colors.white,
              height: 1.4,
            ),
          ),
          const SizedBox(height: 20),
          if (question.type == QuizQuestionType.shortAnswer)
            ShortAnswerField(
              initialValue:
                  selectedAnswer is String ? selectedAnswer as String : '',
              onChanged: (value) => onShortAnswerChanged(question.id, value),
            )
          else
            ...question.answers.asMap().entries.map(
              (entry) {
                final index = entry.key;
                final answer = entry.value;
                final multi = question.type == QuizQuestionType.multipleSelect;
                // Answers are stored as a list of selected ids (single-choice is a
                // one-element list); tolerate a legacy bare string too.
                final selected = selectedAnswer is List
                    ? (selectedAnswer as List).contains(answer.id)
                    : selectedAnswer == answer.id;
                return AnswerOptionTile(
                  answer: answer,
                  isSelected: selected,
                  optionIndex: index,
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
  const _Chip({required this.label, this.isHighlight = false});

  final String label;
  final bool isHighlight;

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(
        horizontal: 10,
        vertical: 4,
      ),
      decoration: BoxDecoration(
        color: isHighlight ? const Color(0x2210B981) : AppColors.bgTertiary,
        borderRadius: BorderRadius.circular(12),
        border: Border.all(
          color: isHighlight ? const Color(0x5510B981) : AppColors.borderStrong,
        ),
      ),
      child: Text(
        label,
        style: TextStyle(
          fontSize: 12,
          fontWeight: FontWeight.w600,
          color: isHighlight ? const Color(0xFF10B981) : const Color(0xFF94A3B8),
        ),
      ),
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

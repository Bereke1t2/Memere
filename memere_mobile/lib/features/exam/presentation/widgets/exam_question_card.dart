import 'package:flutter/material.dart';

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
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        // Question Container Box
        Container(
          width: double.infinity,
          padding: const EdgeInsets.all(18),
          decoration: BoxDecoration(
            color: const Color(0xFF111116),
            borderRadius: BorderRadius.circular(18),
            border: Border.all(color: const Color(0xFF22222C)),
            boxShadow: [
              BoxShadow(
                color: Colors.black.withAlpha(60),
                blurRadius: 10,
                offset: const Offset(0, 4),
              ),
            ],
          ),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // Meta Tag Pills (Subject / Topic / Marks)
              Row(
                crossAxisAlignment: CrossAxisAlignment.center,
                children: [
                  // Marks Pill
                  Container(
                    padding: const EdgeInsets.symmetric(horizontal: 9, vertical: 4),
                    decoration: BoxDecoration(
                      color: const Color(0x2210B981),
                      borderRadius: BorderRadius.circular(8),
                      border: Border.all(
                        color: const Color(0xFF10B981).withAlpha(100),
                        width: 1,
                      ),
                    ),
                    child: Row(
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        const Icon(
                          Icons.verified_rounded,
                          size: 13,
                          color: Color(0xFF10B981),
                        ),
                        const SizedBox(width: 4),
                        Text(
                          '${question.marks} Mark${question.marks == 1 ? '' : 's'}',
                          style: AppTextStyles.labelSmall.copyWith(
                            color: const Color(0xFF10B981),
                            fontWeight: FontWeight.w700,
                          ),
                        ),
                      ],
                    ),
                  ),
                  const SizedBox(width: 8),
                  // Topic / Type Pill
                  if (question.topic != null && question.topic!.isNotEmpty)
                    Expanded(
                      child: Container(
                        padding: const EdgeInsets.symmetric(horizontal: 9, vertical: 4),
                        decoration: BoxDecoration(
                          color: const Color(0xFF1A1A22),
                          borderRadius: BorderRadius.circular(8),
                          border: Border.all(
                            color: const Color(0xFF2C2C38),
                            width: 1,
                          ),
                        ),
                        child: Text(
                          question.topic!,
                          maxLines: 1,
                          overflow: TextOverflow.ellipsis,
                          style: AppTextStyles.labelSmall.copyWith(
                            color: const Color(0xFFA1A1AA),
                            fontWeight: FontWeight.w600,
                          ),
                        ),
                      ),
                    )
                  else
                    Container(
                      padding: const EdgeInsets.symmetric(horizontal: 9, vertical: 4),
                      decoration: BoxDecoration(
                        color: const Color(0xFF1A1A22),
                        borderRadius: BorderRadius.circular(8),
                        border: Border.all(
                          color: const Color(0xFF2C2C38),
                          width: 1,
                        ),
                      ),
                      child: Text(
                        _questionTypeLabel(question.type),
                        style: AppTextStyles.labelSmall.copyWith(
                          color: const Color(0xFFA1A1AA),
                          fontWeight: FontWeight.w600,
                        ),
                      ),
                    ),
                ],
              ),
              const SizedBox(height: 16),
              // Question Prompt Text
              Text(
                question.text,
                style: AppTextStyles.titleMedium.copyWith(
                  fontSize: 17,
                  fontWeight: FontWeight.w600,
                  height: 1.5,
                  color: Colors.white,
                  letterSpacing: -0.2,
                ),
              ),
            ],
          ),
        ),
        const SizedBox(height: 20),

        // Section Title for Options
        Text(
          question.type == ExamQuestionType.multipleSelect
              ? 'Select all that apply:'
              : question.type == ExamQuestionType.shortAnswer
                  ? 'Your Answer:'
                  : 'Choose one answer:',
          style: AppTextStyles.labelMedium.copyWith(
            color: const Color(0xFF71717A),
            fontWeight: FontWeight.w600,
            letterSpacing: 0.3,
          ),
        ),
        const SizedBox(height: 10),

        // Options List / Short Answer
        if (question.type == ExamQuestionType.shortAnswer)
          _ShortAnswerField(
            key: ValueKey('short-${question.questionId}'),
            initialValue:
                selectedAnswer is String ? selectedAnswer as String : '',
            onChanged: (value) =>
                onShortAnswerChanged(question.questionId, value),
          )
        else
          ...question.answers.asMap().entries.map((entry) {
            final idx = entry.key;
            final answer = entry.value;
            final multi = question.type == ExamQuestionType.multipleSelect;
            // Answers are stored as a list of selected ids (single-choice is a
            // one-element list); tolerate a legacy bare string too.
            final selected = selectedAnswer is List
                ? (selectedAnswer as List).contains(answer.id)
                : selectedAnswer == answer.id;
            return ExamAnswerOptionTile(
              index: idx,
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
    return Container(
      decoration: BoxDecoration(
        color: const Color(0xFF131318),
        borderRadius: BorderRadius.circular(14),
        border: Border.all(color: const Color(0xFF242430)),
      ),
      child: TextField(
        controller: _controller,
        minLines: 3,
        maxLines: 5,
        onChanged: widget.onChanged,
        style: AppTextStyles.bodyMedium.copyWith(color: Colors.white),
        decoration: InputDecoration(
          hintText: 'Type your explanation or answer here...',
          hintStyle: AppTextStyles.bodyMedium.copyWith(
            color: const Color(0xFF52525B),
          ),
          border: InputBorder.none,
          contentPadding: const EdgeInsets.all(16),
        ),
      ),
    );
  }
}

String _questionTypeLabel(ExamQuestionType type) {
  switch (type) {
    case ExamQuestionType.multipleChoice:
      return 'Single Choice';
    case ExamQuestionType.multipleSelect:
      return 'Multi Select';
    case ExamQuestionType.trueFalse:
      return 'True/False';
    case ExamQuestionType.shortAnswer:
      return 'Short Answer';
  }
}

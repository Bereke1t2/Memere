import 'package:flutter/material.dart';

import '../../domain/entities/exam_feedback_answer_entity.dart';
import '../../domain/entities/exam_question_feedback_entity.dart';

class ExamQuestionFeedbackTile extends StatelessWidget {
  const ExamQuestionFeedbackTile({
    super.key,
    required this.feedback,
    required this.index,
    this.questionText,
  });

  final ExamQuestionFeedbackEntity feedback;
  final int index;
  final String? questionText;

  @override
  Widget build(BuildContext context) {
    final isCorrect = feedback.correct;
    final prompt = feedback.questionText ?? questionText ?? 'Question ${index + 1}';
    final topic = feedback.topic ?? feedback.subject;

    return Container(
      margin: const EdgeInsets.only(bottom: 16),
      decoration: BoxDecoration(
        color: const Color(0xFF111116),
        borderRadius: BorderRadius.circular(18),
        border: Border.all(
          color: isCorrect
              ? const Color(0xFF10B981).withAlpha(60)
              : const Color(0xFFEF4444).withAlpha(60),
          width: 1.2,
        ),
        boxShadow: [
          BoxShadow(
            color: (isCorrect ? const Color(0xFF10B981) : const Color(0xFFEF4444)).withAlpha(20),
            blurRadius: 10,
            offset: const Offset(0, 3),
          ),
        ],
      ),
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Top Badge Row: Question Number, Outcome Tag, Marks
            Row(
              children: [
                // Outcome Icon + Question Index
                Container(
                  padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 5),
                  decoration: BoxDecoration(
                    color: isCorrect
                        ? const Color(0x2210B981)
                        : const Color(0x22EF4444),
                    borderRadius: BorderRadius.circular(10),
                    border: Border.all(
                      color: isCorrect
                          ? const Color(0xFF10B981).withAlpha(100)
                          : const Color(0xFFEF4444).withAlpha(100),
                    ),
                  ),
                  child: Row(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      Icon(
                        isCorrect
                            ? Icons.check_circle_rounded
                            : Icons.cancel_rounded,
                        size: 15,
                        color: isCorrect
                            ? const Color(0xFF10B981)
                            : const Color(0xFFEF4444),
                      ),
                      const SizedBox(width: 5),
                      Text(
                        'Q${index + 1}',
                        style: TextStyle(
                          color: isCorrect
                              ? const Color(0xFF10B981)
                              : const Color(0xFFEF4444),
                          fontWeight: FontWeight.w800,
                          fontSize: 13,
                        ),
                      ),
                    ],
                  ),
                ),
                const SizedBox(width: 8),

                // Topic Tag (if available)
                if (topic != null && topic.isNotEmpty)
                  Expanded(
                    child: Text(
                      topic,
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                      style: const TextStyle(
                        color: Color(0xFFA1A1AA),
                        fontSize: 12,
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                  )
                else
                  const Spacer(),

                // Marks Tag
                Container(
                  padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                  decoration: BoxDecoration(
                    color: const Color(0xFF1A1A24),
                    borderRadius: BorderRadius.circular(8),
                    border: Border.all(color: const Color(0xFF2C2C3C)),
                  ),
                  child: Text(
                    '${feedback.marksAwarded}/${feedback.marksPossible} Marks',
                    style: TextStyle(
                      color: isCorrect
                          ? const Color(0xFF34D399)
                          : const Color(0xFFF87171),
                      fontWeight: FontWeight.w700,
                      fontSize: 11,
                    ),
                  ),
                ),
              ],
            ),
            const SizedBox(height: 14),

            // Question Prompt
            Text(
              prompt,
              style: const TextStyle(
                color: Colors.white,
                fontWeight: FontWeight.w600,
                fontSize: 15,
                height: 1.45,
              ),
            ),
            const SizedBox(height: 14),

            // Option Choices with Status
            if (feedback.answers.isNotEmpty) ...[
              ...feedback.answers.asMap().entries.map((entry) {
                final idx = entry.key;
                final answer = entry.value;
                return _buildAnswerFeedbackRow(idx, answer);
              }),
            ] else ...[
              // Fallback for short answers or legacy payload
              _buildFallbackSelectionDisplay(),
            ],

            // Explanation Box
            if (feedback.explanation != null && feedback.explanation!.trim().isNotEmpty) ...[
              const SizedBox(height: 12),
              Container(
                width: double.infinity,
                padding: const EdgeInsets.all(12),
                decoration: BoxDecoration(
                  color: const Color(0xFF151922),
                  borderRadius: BorderRadius.circular(12),
                  border: Border.all(
                    color: const Color(0xFF38BDF8).withAlpha(70),
                    width: 1,
                  ),
                ),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const Row(
                      children: [
                        Icon(
                          Icons.lightbulb_rounded,
                          size: 16,
                          color: Color(0xFF38BDF8),
                        ),
                        SizedBox(width: 6),
                        Text(
                          'Explanation & Solution',
                          style: TextStyle(
                            color: Color(0xFF38BDF8),
                            fontWeight: FontWeight.w700,
                            fontSize: 12,
                          ),
                        ),
                      ],
                    ),
                    const SizedBox(height: 6),
                    Text(
                      feedback.explanation!,
                      style: const TextStyle(
                        color: Color(0xFFD4D4D8),
                        fontSize: 13,
                        height: 1.4,
                      ),
                    ),
                  ],
                ),
              ),
            ],
          ],
        ),
      ),
    );
  }

  Widget _buildAnswerFeedbackRow(int idx, ExamFeedbackAnswerEntity answer) {
    final letters = ['A', 'B', 'C', 'D', 'E', 'F'];
    final letter = idx < letters.length ? letters[idx] : '${idx + 1}';
    final isSelected = feedback.selectedAnswers.contains(answer.id);
    final isKeyCorrect = answer.isCorrect || feedback.correctAnswerIds.contains(answer.id);

    Color bgColor = const Color(0xFF14141C);
    Color borderColor = const Color(0xFF22222E);
    Color letterBg = const Color(0xFF1C1C26);
    Color letterText = const Color(0xFFA1A1AA);
    Widget? statusBadge;

    if (isKeyCorrect) {
      // Correct answer
      bgColor = const Color(0xFF063A28);
      borderColor = const Color(0xFF10B981);
      letterBg = const Color(0xFF10B981);
      letterText = Colors.black;
      statusBadge = Container(
        padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
        decoration: BoxDecoration(
          color: const Color(0x3310B981),
          borderRadius: BorderRadius.circular(6),
        ),
        child: const Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(Icons.check_rounded, size: 12, color: Color(0xFF10B981)),
            SizedBox(width: 3),
            Text(
              'Correct Answer',
              style: TextStyle(
                color: Color(0xFF10B981),
                fontSize: 10,
                fontWeight: FontWeight.w700,
              ),
            ),
          ],
        ),
      );
    } else if (isSelected && !isKeyCorrect) {
      // User selected wrong answer
      bgColor = const Color(0xFF3B1212);
      borderColor = const Color(0xFFEF4444);
      letterBg = const Color(0xFFEF4444);
      letterText = Colors.white;
      statusBadge = Container(
        padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
        decoration: BoxDecoration(
          color: const Color(0x33EF4444),
          borderRadius: BorderRadius.circular(6),
        ),
        child: const Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(Icons.close_rounded, size: 12, color: Color(0xFFEF4444)),
            SizedBox(width: 3),
            Text(
              'Your Answer',
              style: TextStyle(
                color: Color(0xFFEF4444),
                fontSize: 10,
                fontWeight: FontWeight.w700,
              ),
            ),
          ],
        ),
      );
    }

    return Container(
      margin: const EdgeInsets.only(bottom: 8),
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
      decoration: BoxDecoration(
        color: bgColor,
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: borderColor, width: 1.2),
      ),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.center,
        children: [
          // Letter badge
          Container(
            width: 26,
            height: 26,
            alignment: Alignment.center,
            decoration: BoxDecoration(
              color: letterBg,
              borderRadius: BorderRadius.circular(6),
            ),
            child: Text(
              letter,
              style: TextStyle(
                color: letterText,
                fontWeight: FontWeight.w800,
                fontSize: 12,
              ),
            ),
          ),
          const SizedBox(width: 10),
          // Option text
          Expanded(
            child: Text(
              answer.text,
              style: TextStyle(
                color: isKeyCorrect
                    ? Colors.white
                    : (isSelected ? const Color(0xFFFCA5A5) : const Color(0xFFD4D4D8)),
                fontWeight: isKeyCorrect || isSelected ? FontWeight.w600 : FontWeight.w400,
                fontSize: 13,
              ),
            ),
          ),
          if (statusBadge != null) ...[
            const SizedBox(width: 8),
            statusBadge,
          ],
        ],
      ),
    );
  }

  Widget _buildFallbackSelectionDisplay() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        if (feedback.selectedAnswers.isNotEmpty)
          Padding(
            padding: const EdgeInsets.only(bottom: 4),
            child: Text(
              'Your Answer: ${feedback.selectedAnswers.join(", ")}',
              style: TextStyle(
                color: feedback.correct ? const Color(0xFF10B981) : const Color(0xFFEF4444),
                fontSize: 13,
                fontWeight: FontWeight.w600,
              ),
            ),
          ),
        if (!feedback.correct && feedback.correctAnswerIds.isNotEmpty)
          Text(
            'Correct: ${feedback.correctAnswerIds.join(", ")}',
            style: const TextStyle(
              color: Color(0xFF10B981),
              fontSize: 13,
              fontWeight: FontWeight.w600,
            ),
          ),
      ],
    );
  }
}

import 'exam_feedback_answer_entity.dart';

class ExamQuestionFeedbackEntity {
  const ExamQuestionFeedbackEntity({
    required this.questionId,
    this.questionText,
    this.type,
    this.subject,
    this.topic,
    required this.correct,
    required this.marksAwarded,
    required this.marksPossible,
    required this.selectedAnswers,
    required this.correctAnswerIds,
    this.explanation,
    this.answers = const [],
  });

  final String questionId;
  final String? questionText;
  final String? type;
  final String? subject;
  final String? topic;
  final bool correct;
  final int marksAwarded;
  final int marksPossible;
  final List<String> selectedAnswers;
  final List<String> correctAnswerIds;
  final String? explanation;
  final List<ExamFeedbackAnswerEntity> answers;
}

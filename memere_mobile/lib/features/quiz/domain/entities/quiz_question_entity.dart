import 'quiz_answer_entity.dart';

enum QuizQuestionType {
  multipleChoice,
  multipleSelect,
  trueFalse,
  shortAnswer,
}

class QuizQuestionEntity {
  const QuizQuestionEntity({
    required this.id,
    required this.text,
    required this.type,
    required this.points,
    this.subject,
    this.topic,
    required this.answers,
  });

  final String id;
  final String text;
  final QuizQuestionType type;
  final int points;
  final String? subject;
  final String? topic;
  final List<QuizAnswerEntity> answers;

  bool get isChoiceBased =>
      type == QuizQuestionType.multipleChoice ||
      type == QuizQuestionType.multipleSelect ||
      type == QuizQuestionType.trueFalse;
}

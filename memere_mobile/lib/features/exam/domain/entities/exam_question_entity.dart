import 'exam_answer_entity.dart';

enum ExamQuestionType {
  multipleChoice,
  multipleSelect,
  trueFalse,
  shortAnswer,
}

class ExamQuestionEntity {
  const ExamQuestionEntity({
    required this.questionId,
    required this.text,
    required this.type,
    required this.marks,
    this.subject,
    this.topic,
    required this.answers,
  });

  final String questionId;
  final String text;
  final ExamQuestionType type;
  final int marks;
  final String? subject;
  final String? topic;
  final List<ExamAnswerEntity> answers;

  bool get isChoiceBased =>
      type == ExamQuestionType.multipleChoice ||
      type == ExamQuestionType.multipleSelect ||
      type == ExamQuestionType.trueFalse;
}

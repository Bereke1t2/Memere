import '../../domain/entities/exam_question_entity.dart';
import 'exam_answer_model.dart';
import 'exam_model_helpers.dart';

class ExamQuestionModel extends ExamQuestionEntity {
  const ExamQuestionModel({
    required super.questionId,
    required super.text,
    required super.type,
    required super.marks,
    super.subject,
    super.topic,
    required super.answers,
  });

  factory ExamQuestionModel.fromJson(Map<String, dynamic> json) {
    final answersJson = json['answers'];
    final answers = answersJson is List
        ? answersJson
            .whereType<Map<String, dynamic>>()
            .map(ExamAnswerModel.fromJson)
            .toList()
        : <ExamAnswerModel>[];

    return ExamQuestionModel(
      questionId: examStringValue(json['question_id']),
      text: examStringValue(json['text'], fallback: 'Untitled question'),
      type: parseExamQuestionType(examStringValue(json['type'])),
      marks: examIntValue(json['marks'], fallback: 1),
      subject: examNullableString(json['subject']),
      topic: examNullableString(json['topic']),
      answers: answers,
    );
  }

  Map<String, dynamic> toJson() => {
        'question_id': questionId,
        'text': text,
        'type': examQuestionTypeValue(type),
        'marks': marks,
        'subject': subject,
        'topic': topic,
        'answers': answers
            .whereType<ExamAnswerModel>()
            .map((answer) => answer.toJson())
            .toList(),
      };
}

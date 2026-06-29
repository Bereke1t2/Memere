import '../../domain/entities/quiz_question_entity.dart';
import 'quiz_answer_model.dart';
import 'quiz_model_helpers.dart';

class QuizQuestionModel extends QuizQuestionEntity {
  const QuizQuestionModel({
    required super.id,
    required super.text,
    required super.type,
    required super.points,
    super.subject,
    super.topic,
    required super.answers,
  });

  factory QuizQuestionModel.fromJson(Map<String, dynamic> json) {
    final answersJson = json['answers'];
    final answers = answersJson is List
        ? answersJson
            .whereType<Map<String, dynamic>>()
            .map(QuizAnswerModel.fromJson)
            .toList()
        : <QuizAnswerModel>[];

    return QuizQuestionModel(
      id: quizStringValue(json['id']),
      text: quizStringValue(json['text'], fallback: 'Untitled question'),
      type: parseQuizQuestionType(quizStringValue(json['type'])),
      points: quizIntValue(json['points'], fallback: 1),
      subject: quizNullableString(json['subject']),
      topic: quizNullableString(json['topic']),
      answers: answers,
    );
  }

  Map<String, dynamic> toJson() => {
        'id': id,
        'text': text,
        'type': quizQuestionTypeValue(type),
        'points': points,
        'subject': subject,
        'topic': topic,
        'answers': answers
            .whereType<QuizAnswerModel>()
            .map((answer) => answer.toJson())
            .toList(),
      };
}

QuizQuestionType parseQuizQuestionType(String value) {
  switch (value.toLowerCase()) {
    case 'multiple_select':
      return QuizQuestionType.multipleSelect;
    case 'true_false':
      return QuizQuestionType.trueFalse;
    case 'short_answer':
      return QuizQuestionType.shortAnswer;
    default:
      return QuizQuestionType.multipleChoice;
  }
}

String quizQuestionTypeValue(QuizQuestionType type) {
  switch (type) {
    case QuizQuestionType.multipleChoice:
      return 'multiple_choice';
    case QuizQuestionType.multipleSelect:
      return 'multiple_select';
    case QuizQuestionType.trueFalse:
      return 'true_false';
    case QuizQuestionType.shortAnswer:
      return 'short_answer';
  }
}

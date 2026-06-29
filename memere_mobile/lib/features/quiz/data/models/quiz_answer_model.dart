import '../../domain/entities/quiz_answer_entity.dart';
import 'quiz_model_helpers.dart';

class QuizAnswerModel extends QuizAnswerEntity {
  const QuizAnswerModel({
    required super.id,
    required super.text,
    required super.orderIndex,
  });

  factory QuizAnswerModel.fromJson(Map<String, dynamic> json) {
    return QuizAnswerModel(
      id: quizStringValue(json['id']),
      text: quizStringValue(json['text']),
      orderIndex: quizIntValue(json['order_index']),
    );
  }

  Map<String, dynamic> toJson() => {
        'id': id,
        'text': text,
        'order_index': orderIndex,
      };
}

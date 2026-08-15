import '../../domain/entities/exam_feedback_answer_entity.dart';
import 'exam_model_helpers.dart';

class ExamFeedbackAnswerModel extends ExamFeedbackAnswerEntity {
  const ExamFeedbackAnswerModel({
    required super.id,
    required super.text,
    required super.isCorrect,
    super.orderIndex,
  });

  factory ExamFeedbackAnswerModel.fromJson(Map<String, dynamic> json) {
    return ExamFeedbackAnswerModel(
      id: examStringValue(json['id']),
      text: examStringValue(json['text']),
      isCorrect: examBoolValue(json['is_correct']),
      orderIndex: examIntValue(json['order_index']),
    );
  }

  Map<String, dynamic> toJson() => {
        'id': id,
        'text': text,
        'is_correct': isCorrect,
        'order_index': orderIndex,
      };
}

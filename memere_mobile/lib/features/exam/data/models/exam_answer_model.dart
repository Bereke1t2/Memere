import '../../domain/entities/exam_answer_entity.dart';
import 'exam_model_helpers.dart';

class ExamAnswerModel extends ExamAnswerEntity {
  const ExamAnswerModel({
    required super.id,
    required super.text,
    required super.orderIndex,
  });

  factory ExamAnswerModel.fromJson(Map<String, dynamic> json) {
    return ExamAnswerModel(
      id: examStringValue(json['id']),
      text: examStringValue(json['text']),
      orderIndex: examIntValue(json['order_index']),
    );
  }

  Map<String, dynamic> toJson() => {
        'id': id,
        'text': text,
        'order_index': orderIndex,
      };
}

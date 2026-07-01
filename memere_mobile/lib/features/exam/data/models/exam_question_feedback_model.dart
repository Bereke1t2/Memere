import '../../domain/entities/exam_question_feedback_entity.dart';
import 'exam_model_helpers.dart';

class ExamQuestionFeedbackModel extends ExamQuestionFeedbackEntity {
  const ExamQuestionFeedbackModel({
    required super.questionId,
    required super.correct,
    required super.marksAwarded,
    required super.marksPossible,
    required super.selectedAnswers,
    required super.correctAnswerIds,
    super.explanation,
  });

  factory ExamQuestionFeedbackModel.fromJson(Map<String, dynamic> json) {
    return ExamQuestionFeedbackModel(
      questionId: examStringValue(json['question_id']),
      correct: examBoolValue(json['correct']),
      marksAwarded: examIntValue(json['marks_awarded']),
      marksPossible: examIntValue(json['marks_possible']),
      selectedAnswers: examStringList(json['selected_answers']),
      correctAnswerIds: examStringList(json['correct_answer_ids']),
      explanation: examNullableString(json['explanation']),
    );
  }

  Map<String, dynamic> toJson() => {
        'question_id': questionId,
        'correct': correct,
        'marks_awarded': marksAwarded,
        'marks_possible': marksPossible,
        'selected_answers': selectedAnswers,
        'correct_answer_ids': correctAnswerIds,
        'explanation': explanation,
      };
}

import '../../domain/entities/question_feedback_entity.dart';
import 'quiz_model_helpers.dart';

class QuestionFeedbackModel extends QuestionFeedbackEntity {
  const QuestionFeedbackModel({
    required super.questionId,
    required super.correct,
    required super.pointsAwarded,
    required super.pointsPossible,
    required super.selectedAnswers,
    required super.correctAnswerIds,
    super.explanation,
  });

  factory QuestionFeedbackModel.fromJson(Map<String, dynamic> json) {
    return QuestionFeedbackModel(
      questionId: quizStringValue(json['question_id']),
      correct: quizBoolValue(json['correct']),
      pointsAwarded: quizIntValue(json['points_awarded']),
      pointsPossible: quizIntValue(json['points_possible']),
      selectedAnswers: quizStringList(json['selected_answers']),
      correctAnswerIds: quizStringList(json['correct_answer_ids']),
      explanation: quizNullableString(json['explanation']),
    );
  }

  Map<String, dynamic> toJson() => {
        'question_id': questionId,
        'correct': correct,
        'points_awarded': pointsAwarded,
        'points_possible': pointsPossible,
        'selected_answers': selectedAnswers,
        'correct_answer_ids': correctAnswerIds,
        'explanation': explanation,
      };
}

import '../../domain/entities/question_feedback_entity.dart';
import 'quiz_model_helpers.dart';

class QuizAnswerFeedbackModel extends QuizAnswerFeedbackEntity {
  const QuizAnswerFeedbackModel({
    required super.id,
    required super.text,
    required super.isCorrect,
  });

  factory QuizAnswerFeedbackModel.fromJson(Map<String, dynamic> json) {
    return QuizAnswerFeedbackModel(
      id: quizStringValue(json['id']),
      text: quizStringValue(json['text']),
      isCorrect: quizBoolValue(json['is_correct']),
    );
  }
}

class QuestionFeedbackModel extends QuestionFeedbackEntity {
  const QuestionFeedbackModel({
    required super.questionId,
    super.questionText,
    required super.correct,
    required super.pointsAwarded,
    required super.pointsPossible,
    required super.selectedAnswers,
    required super.correctAnswerIds,
    super.explanation,
    super.answers,
  });

  factory QuestionFeedbackModel.fromJson(Map<String, dynamic> json) {
    final answersList = (json['answers'] as List<dynamic>?)
            ?.map((e) => QuizAnswerFeedbackModel.fromJson(e as Map<String, dynamic>))
            .toList() ??
        [];

    return QuestionFeedbackModel(
      questionId: quizStringValue(json['question_id']),
      questionText: quizNullableString(json['question_text']),
      correct: quizBoolValue(json['correct']),
      pointsAwarded: quizIntValue(json['points_awarded']),
      pointsPossible: quizIntValue(json['points_possible']),
      selectedAnswers: quizStringList(json['selected_answers']),
      correctAnswerIds: quizStringList(json['correct_answer_ids']),
      explanation: quizNullableString(json['explanation']),
      answers: answersList,
    );
  }

  Map<String, dynamic> toJson() => {
        'question_id': questionId,
        'question_text': questionText,
        'correct': correct,
        'points_awarded': pointsAwarded,
        'points_possible': pointsPossible,
        'selected_answers': selectedAnswers,
        'correct_answer_ids': correctAnswerIds,
        'explanation': explanation,
      };
}

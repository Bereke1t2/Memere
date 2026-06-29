import '../../domain/entities/quiz_result_entity.dart';
import 'question_feedback_model.dart';
import 'quiz_attempt_model.dart';
import 'quiz_model_helpers.dart';
import 'subject_score_model.dart';

class QuizResultModel extends QuizResultEntity {
  const QuizResultModel({
    required super.attemptId,
    required super.quizId,
    required super.attemptNumber,
    required super.status,
    required super.score,
    required super.totalPoints,
    required super.percentage,
    required super.passed,
    super.submittedAt,
    required super.feedback,
    required super.subjectBreakdown,
  });

  factory QuizResultModel.fromJson(Map<String, dynamic> json) {
    final feedbackJson = json['feedback'];
    final feedback = feedbackJson is List
        ? feedbackJson
            .whereType<Map<String, dynamic>>()
            .map(QuestionFeedbackModel.fromJson)
            .toList()
        : <QuestionFeedbackModel>[];

    final breakdownJson = json['subject_breakdown'];
    final breakdown = <String, SubjectScoreModel>{};
    if (breakdownJson is Map) {
      for (final entry in breakdownJson.entries) {
        final value = entry.value;
        if (value is Map<String, dynamic>) {
          breakdown[entry.key.toString()] = SubjectScoreModel.fromJson(value);
        } else if (value is Map) {
          breakdown[entry.key.toString()] = SubjectScoreModel.fromJson(
            Map<String, dynamic>.from(value),
          );
        }
      }
    }

    return QuizResultModel(
      attemptId: quizStringValue(json['attempt_id']),
      quizId: quizStringValue(json['quiz_id']),
      attemptNumber: quizIntValue(json['attempt_number'], fallback: 1),
      status: parseQuizAttemptStatus(quizStringValue(json['status'])),
      score: quizDoubleValue(json['score']),
      totalPoints: quizIntValue(json['total_points']),
      percentage: quizDoubleValue(json['percentage']),
      passed: quizBoolValue(json['passed']),
      submittedAt: quizDateValue(json['submitted_at']),
      feedback: feedback,
      subjectBreakdown: breakdown,
    );
  }

  Map<String, dynamic> toJson() => {
        'attempt_id': attemptId,
        'quiz_id': quizId,
        'attempt_number': attemptNumber,
        'status': quizAttemptStatusValue(status),
        'score': score,
        'total_points': totalPoints,
        'percentage': percentage,
        'passed': passed,
        'submitted_at': submittedAt?.toIso8601String(),
        'feedback': feedback
            .whereType<QuestionFeedbackModel>()
            .map((item) => item.toJson())
            .toList(),
        'subject_breakdown': subjectBreakdown.map(
          (key, value) => MapEntry(
            key,
            value is SubjectScoreModel
                ? value.toJson()
                : SubjectScoreModel(
                    earned: value.earned,
                    possible: value.possible,
                  ).toJson(),
          ),
        ),
      };
}

import '../../domain/entities/exam_attempt_analytics_entity.dart';
import 'exam_model_helpers.dart';
import 'exam_subject_score_model.dart';

class ExamAttemptAnalyticsModel extends ExamAttemptAnalyticsEntity {
  const ExamAttemptAnalyticsModel({
    required super.attemptId,
    required super.examId,
    required super.score,
    required super.percentage,
    required super.subjectBreakdown,
    required super.weakAreas,
    super.percentile,
  });

  factory ExamAttemptAnalyticsModel.fromJson(Map<String, dynamic> json) {
    return ExamAttemptAnalyticsModel(
      attemptId: examStringValue(json['attempt_id']),
      examId: examStringValue(json['exam_id']),
      score: examDoubleValue(json['score']),
      percentage: examDoubleValue(json['percentage']),
      subjectBreakdown: _parseScores(json['subject_breakdown']),
      weakAreas: _parseScores(json['weak_areas']),
      percentile: examNullableDouble(json['percentile']),
    );
  }

  static List<ExamSubjectScoreModel> _parseScores(Object? value) {
    if (value is! List) return const [];
    return value
        .whereType<Map<String, dynamic>>()
        .map(ExamSubjectScoreModel.fromJson)
        .toList();
  }
}

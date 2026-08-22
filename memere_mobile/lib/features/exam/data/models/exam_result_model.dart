import '../../domain/entities/exam_result_entity.dart';
import '../../domain/entities/exam_subject_score_entity.dart';
import 'exam_model_helpers.dart';
import 'exam_question_feedback_model.dart';
import 'exam_subject_score_model.dart';

class ExamResultModel extends ExamResultEntity {
  const ExamResultModel({
    required super.attemptId,
    required super.examId,
    required super.status,
    required super.score,
    required super.totalMarks,
    required super.percentage,
    required super.passMarks,
    required super.passed,
    super.submittedAt,
    required super.feedback,
    required super.subjectBreakdown,
  });

  factory ExamResultModel.fromJson(Map<String, dynamic> json) {
    final feedbackJson = json['feedback'];
    final feedback = feedbackJson is List
        ? feedbackJson
            .whereType<Map<String, dynamic>>()
            .map(ExamQuestionFeedbackModel.fromJson)
            .toList()
        : <ExamQuestionFeedbackModel>[];

    final breakdownJson = json['subject_breakdown'];
    final breakdown = <String, ExamSubjectScoreEntity>{};
    if (breakdownJson is Map) {
      for (final entry in breakdownJson.entries) {
        final value = entry.value;
        final key = entry.key.toString();
        if (value is Map<String, dynamic>) {
          breakdown[key] = ExamSubjectScoreModel.fromMapEntry(key, value);
        } else if (value is Map) {
          breakdown[key] = ExamSubjectScoreModel.fromMapEntry(
            key,
            Map<String, dynamic>.from(value),
          );
        }
      }
    }

    return ExamResultModel(
      attemptId: examStringValue(json['attempt_id']),
      examId: examStringValue(json['exam_id']),
      status: parseExamAttemptStatus(examStringValue(json['status'])),
      score: examDoubleValue(json['score']),
      totalMarks: examIntValue(json['total_marks']),
      percentage: examDoubleValue(json['percentage']),
      passMarks: examIntValue(json['pass_marks']),
      passed: examBoolValue(json['passed']),
      submittedAt: examDateValue(json['submitted_at']),
      feedback: feedback,
      subjectBreakdown: breakdown,
    );
  }

  /// Round-trips with [ExamResultModel.fromJson] so an on-device-graded result
  /// can be persisted to the `offline_attempt_results` box and rehydrated on the
  /// existing result screen. `subject_breakdown` is written as the result-shape
  /// map (key → `{key, earned, possible}`) that `fromMapEntry` reads back.
  Map<String, dynamic> toJson() => {
        'attempt_id': attemptId,
        'exam_id': examId,
        'status': examAttemptStatusValue(status),
        'score': score,
        'total_marks': totalMarks,
        'percentage': percentage,
        'pass_marks': passMarks,
        'passed': passed,
        'submitted_at': submittedAt?.toIso8601String(),
        'feedback': feedback
            .whereType<ExamQuestionFeedbackModel>()
            .map((item) => item.toJson())
            .toList(),
        'subject_breakdown': subjectBreakdown.map(
          (key, value) => MapEntry(
            key,
            value is ExamSubjectScoreModel
                ? value.toJson()
                : ExamSubjectScoreModel(
                    key: key,
                    earned: value.earned,
                    possible: value.possible,
                  ).toJson(),
          ),
        ),
      };
}

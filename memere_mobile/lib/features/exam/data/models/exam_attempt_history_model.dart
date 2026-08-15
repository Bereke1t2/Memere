import '../../domain/entities/exam_attempt_history_entity.dart';
import 'exam_model_helpers.dart';

class ExamAttemptHistoryModel extends ExamAttemptHistoryEntity {
  const ExamAttemptHistoryModel({
    required super.attemptId,
    required super.examId,
    required super.status,
    super.score,
    super.percentage,
    required super.startedAt,
    super.submittedAt,
  });

  factory ExamAttemptHistoryModel.fromJson(Map<String, dynamic> json) {
    return ExamAttemptHistoryModel(
      attemptId: examStringValue(json['attempt_id']),
      examId: examStringValue(json['exam_id']),
      status: examStringValue(json['status']),
      score: json['score'] != null ? examDoubleValue(json['score']) : null,
      percentage: json['percentage'] != null ? examDoubleValue(json['percentage']) : null,
      startedAt: examDateValue(json['started_at']) ?? DateTime.now(),
      submittedAt: examDateValue(json['submitted_at']),
    );
  }

  Map<String, dynamic> toJson() => {
        'attempt_id': attemptId,
        'exam_id': examId,
        'status': status,
        'score': score,
        'percentage': percentage,
        'started_at': startedAt.toIso8601String(),
        'submitted_at': submittedAt?.toIso8601String(),
      };
}

import 'exam_question_entity.dart';

enum ExamAttemptStatus { inProgress, submitted, expired, graded }

class ExamAttemptEntity {
  const ExamAttemptEntity({
    required this.attemptId,
    required this.examId,
    required this.status,
    this.startedAt,
    this.expiresAt,
    this.remainingSeconds,
    required this.totalMarks,
    required this.questions,
  });

  final String attemptId;
  final String examId;
  final ExamAttemptStatus status;
  final DateTime? startedAt;
  final DateTime? expiresAt;
  final int? remainingSeconds;
  final int totalMarks;
  final List<ExamQuestionEntity> questions;

  bool get hasTimer => expiresAt != null || remainingSeconds != null;
  bool get isInProgress => status == ExamAttemptStatus.inProgress;
}

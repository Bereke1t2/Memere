class ExamAttemptHistoryEntity {
  const ExamAttemptHistoryEntity({
    required this.attemptId,
    required this.examId,
    required this.status,
    this.score,
    this.percentage,
    required this.startedAt,
    this.submittedAt,
  });

  final String attemptId;
  final String examId;
  final String status;
  final double? score;
  final double? percentage;
  final DateTime startedAt;
  final DateTime? submittedAt;

  bool get isGraded => status == 'graded' || status == 'submitted' || status == 'expired';
  bool get isPassed => (percentage ?? 0) >= 50.0;
}

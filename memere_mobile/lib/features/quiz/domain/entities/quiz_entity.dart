class QuizEntity {
  const QuizEntity({
    required this.id,
    required this.courseId,
    required this.title,
    this.timeLimitSeconds,
    required this.passPercentage,
    required this.randomizeQuestions,
    this.maxAttempts,
    required this.questionCount,
    required this.attemptsUsed,
  });

  final String id;
  final String courseId;
  final String title;
  final int? timeLimitSeconds;
  final double passPercentage;
  final bool randomizeQuestions;
  final int? maxAttempts;
  final int questionCount;
  final int attemptsUsed;

  bool get hasTimeLimit => timeLimitSeconds != null && timeLimitSeconds! > 0;
  bool get hasAttemptsRemaining =>
      maxAttempts == null || attemptsUsed < maxAttempts!;
  int? get attemptsRemaining => maxAttempts == null
      ? null
      : (maxAttempts! - attemptsUsed).clamp(0, maxAttempts!).toInt();
}

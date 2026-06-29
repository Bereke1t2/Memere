class LessonProgressEntity {
  const LessonProgressEntity({
    required this.lessonId,
    required this.isCompleted,
    this.completedAt,
    required this.videoProgressSeconds,
  });

  final String lessonId;
  final bool isCompleted;
  final DateTime? completedAt;
  final int videoProgressSeconds;
}

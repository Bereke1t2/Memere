class LessonEntity {
  const LessonEntity({
    required this.id,
    required this.sectionId,
    required this.courseId,
    required this.title,
    required this.type,
    required this.orderIndex,
    required this.isFreePreview,
    required this.durationSeconds,
    required this.isPublished,
    required this.createdAt,
    required this.updatedAt,
    this.videoId,
    this.quizId,
  });

  final String id;
  final String sectionId;
  final String courseId;
  final String title;
  final LessonType type;
  final int orderIndex;
  final bool isFreePreview;
  final int durationSeconds;
  final bool isPublished;
  final DateTime? createdAt;
  final DateTime? updatedAt;
  final String? videoId;
  final String? quizId;

  bool get hasVideo => videoId != null && videoId!.trim().isNotEmpty;
  bool get hasQuiz => quizId != null && quizId!.trim().isNotEmpty;

  String get durationLabel {
    if (durationSeconds <= 0) return '0m';
    final hours = durationSeconds ~/ 3600;
    final minutes = (durationSeconds % 3600) ~/ 60;
    if (hours > 0 && minutes > 0) return '${hours}h ${minutes}m';
    if (hours > 0) return '${hours}h';
    return '${minutes == 0 ? 1 : minutes}m';
  }

  String get typeLabel {
    switch (type) {
      case LessonType.video:
        return 'Video';
      case LessonType.note:
        return 'Note';
      case LessonType.quiz:
        return 'Quiz';
      case LessonType.mixed:
        return 'Mixed';
    }
  }
}

enum LessonType { video, note, quiz, mixed }

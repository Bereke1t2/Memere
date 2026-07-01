class MockExamEntity {
  const MockExamEntity({
    required this.id,
    this.courseId,
    required this.title,
    required this.subject,
    required this.grade,
    required this.durationMinutes,
    required this.totalMarks,
    required this.passMarks,
    this.instructions,
    required this.isPublished,
    this.createdAt,
    this.updatedAt,
  });

  final String id;
  final String? courseId;
  final String title;
  final String subject;
  final int grade;
  final int durationMinutes;
  final int totalMarks;
  final int passMarks;
  final String? instructions;
  final bool isPublished;
  final DateTime? createdAt;
  final DateTime? updatedAt;

  String get durationLabel {
    if (durationMinutes <= 0) return 'Untimed';
    final hours = durationMinutes ~/ 60;
    final minutes = durationMinutes % 60;
    if (hours > 0 && minutes > 0) return '${hours}h ${minutes}m';
    if (hours > 0) return '${hours}h';
    return '${minutes}m';
  }

  String get marksLabel => '$totalMarks marks · pass $passMarks';
}

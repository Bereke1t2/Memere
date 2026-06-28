import 'lesson_entity.dart';

class CourseSectionEntity {
  const CourseSectionEntity({
    required this.id,
    required this.courseId,
    required this.title,
    required this.description,
    required this.orderIndex,
    required this.isPublished,
    required this.createdAt,
    required this.updatedAt,
    required this.lessons,
  });

  final String id;
  final String courseId;
  final String title;
  final String description;
  final int orderIndex;
  final bool isPublished;
  final DateTime? createdAt;
  final DateTime? updatedAt;
  final List<LessonEntity> lessons;

  String get lessonCountLabel {
    if (lessons.length == 1) return '1 lesson';
    return '${lessons.length} lessons';
  }
}

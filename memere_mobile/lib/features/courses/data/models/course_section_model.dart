import '../../domain/entities/course_section_entity.dart';
import 'course_model.dart';
import 'lesson_model.dart';

class CourseSectionModel extends CourseSectionEntity {
  const CourseSectionModel({
    required super.id,
    required super.courseId,
    required super.title,
    required super.description,
    required super.orderIndex,
    required super.isPublished,
    required super.createdAt,
    required super.updatedAt,
    required super.lessons,
  });

  factory CourseSectionModel.fromJson(Map<String, dynamic> json) {
    final lessonsJson = json['lessons'];
    final lessons = lessonsJson is List
        ? lessonsJson
            .whereType<Map<String, dynamic>>()
            .map(LessonModel.fromJson)
            .toList()
        : <LessonModel>[];

    return CourseSectionModel(
      id: modelStringValue(json['id']),
      courseId: modelStringValue(json['course_id']),
      title: modelStringValue(json['title'], fallback: 'Untitled section'),
      description: modelStringValue(json['description']),
      orderIndex: modelIntValue(json['order_index']),
      isPublished: modelBoolValue(json['is_published'], fallback: true),
      createdAt: modelDateValue(json['created_at']),
      updatedAt: modelDateValue(json['updated_at']),
      lessons: lessons,
    );
  }

  Map<String, dynamic> toJson() => {
        'id': id,
        'course_id': courseId,
        'title': title,
        'description': description,
        'order_index': orderIndex,
        'is_published': isPublished,
        'created_at': createdAt?.toIso8601String(),
        'updated_at': updatedAt?.toIso8601String(),
        'lessons': lessons
            .whereType<LessonModel>()
            .map((lesson) => lesson.toJson())
            .toList(),
      };
}

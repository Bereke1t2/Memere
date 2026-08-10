import '../../domain/entities/lesson_entity.dart';
import 'course_model.dart';

class LessonModel extends LessonEntity {
  const LessonModel({
    required super.id,
    required super.sectionId,
    required super.courseId,
    required super.title,
    required super.type,
    required super.orderIndex,
    required super.isFreePreview,
    required super.durationSeconds,
    required super.isPublished,
    required super.createdAt,
    required super.updatedAt,
    super.videoId,
    super.quizId,
    super.content,
    super.pdfUrl,
  });

  factory LessonModel.fromJson(Map<String, dynamic> json) {
    return LessonModel(
      id: modelStringValue(json['id']),
      sectionId: modelStringValue(json['section_id']),
      courseId: modelStringValue(json['course_id']),
      title: modelStringValue(json['title'], fallback: 'Untitled lesson'),
      type: _parseLessonType(modelStringValue(json['type'], fallback: 'video')),
      orderIndex: modelIntValue(json['order_index']),
      isFreePreview: modelBoolValue(json['is_free_preview']),
      durationSeconds: modelIntValue(json['duration_seconds']),
      isPublished: modelBoolValue(json['is_published'], fallback: true),
      createdAt: modelDateValue(json['created_at']),
      updatedAt: modelDateValue(json['updated_at']),
      videoId: modelNullableString(json['video_id']),
      quizId: modelNullableString(json['quiz_id']),
      content: modelNullableString(json['content']),
      pdfUrl: modelNullableString(json['pdf_url']) ?? modelNullableString(json['note_url']),
    );
  }

  Map<String, dynamic> toJson() => {
        'id': id,
        'section_id': sectionId,
        'course_id': courseId,
        'title': title,
        'type': type.name,
        'order_index': orderIndex,
        'is_free_preview': isFreePreview,
        'duration_seconds': durationSeconds,
        'is_published': isPublished,
        'created_at': createdAt?.toIso8601String(),
        'updated_at': updatedAt?.toIso8601String(),
        'video_id': videoId,
        'quiz_id': quizId,
        'content': content,
        'pdf_url': pdfUrl,
      };

  static LessonType _parseLessonType(String value) {
    switch (value.toLowerCase()) {
      case 'note':
        return LessonType.note;
      case 'quiz':
        return LessonType.quiz;
      case 'mixed':
        return LessonType.mixed;
      default:
        return LessonType.video;
    }
  }
}

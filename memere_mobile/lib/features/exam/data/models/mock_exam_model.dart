import '../../domain/entities/mock_exam_entity.dart';
import 'exam_model_helpers.dart';

class MockExamModel extends MockExamEntity {
  const MockExamModel({
    required super.id,
    super.courseId,
    required super.title,
    required super.subject,
    required super.grade,
    required super.durationMinutes,
    required super.totalMarks,
    required super.passMarks,
    super.instructions,
    required super.isPublished,
    super.createdAt,
    super.updatedAt,
  });

  factory MockExamModel.fromJson(Map<String, dynamic> json) {
    return MockExamModel(
      id: examStringValue(json['id']),
      courseId: examNullableString(json['course_id']),
      title: examStringValue(json['title'], fallback: 'Untitled exam'),
      subject: examStringValue(json['subject']),
      grade: examIntValue(json['grade']),
      durationMinutes: examIntValue(json['duration_minutes']),
      totalMarks: examIntValue(json['total_marks']),
      passMarks: examIntValue(json['pass_marks']),
      instructions: examNullableString(json['instructions']),
      isPublished: examBoolValue(json['is_published']),
      createdAt: examDateValue(json['created_at']),
      updatedAt: examDateValue(json['updated_at']),
    );
  }

  Map<String, dynamic> toJson() => {
        'id': id,
        'course_id': courseId,
        'title': title,
        'subject': subject,
        'grade': grade,
        'duration_minutes': durationMinutes,
        'total_marks': totalMarks,
        'pass_marks': passMarks,
        'instructions': instructions,
        'is_published': isPublished,
        'created_at': createdAt?.toIso8601String(),
        'updated_at': updatedAt?.toIso8601String(),
      };
}

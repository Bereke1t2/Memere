import '../../domain/entities/course_entity.dart';

class CourseModel extends CourseEntity {
  const CourseModel({
    required super.id,
    required super.teacherId,
    required super.title,
    required super.slug,
    required super.description,
    required super.shortDescription,
    required super.subject,
    required super.grade,
    required super.thumbnailUrl,
    required super.price,
    required super.currency,
    required super.isFree,
    required super.isPublished,
    required super.language,
    required super.level,
    required super.totalDurationSeconds,
    required super.totalLessons,
    required super.ratingAvg,
    required super.enrollmentCount,
    required super.createdAt,
    required super.updatedAt,
  });

  factory CourseModel.fromJson(Map<String, dynamic> json) {
    return CourseModel(
      id: modelStringValue(json['id']),
      teacherId: modelStringValue(json['teacher_id']),
      title: modelStringValue(json['title'], fallback: 'Untitled course'),
      slug: modelStringValue(json['slug']),
      description: modelStringValue(json['description']),
      shortDescription: modelStringValue(json['short_description']),
      subject: modelStringValue(json['subject'], fallback: 'General'),
      grade: modelIntValue(json['grade'], fallback: 12),
      thumbnailUrl: modelNullableString(json['thumbnail_url']),
      price: modelDoubleValue(json['price']),
      currency: modelStringValue(json['currency'], fallback: 'ETB'),
      isFree: modelBoolValue(json['is_free']),
      isPublished: modelBoolValue(json['is_published'], fallback: true),
      language: modelStringValue(json['language'], fallback: 'en'),
      level: _parseLevel(modelStringValue(json['level'], fallback: 'beginner')),
      totalDurationSeconds: modelIntValue(json['total_duration_seconds']),
      totalLessons: modelIntValue(json['total_lessons']),
      ratingAvg: modelDoubleValue(json['rating_avg']),
      enrollmentCount: modelIntValue(json['enrollment_count']),
      createdAt: modelDateValue(json['created_at']),
      updatedAt: modelDateValue(json['updated_at']),
    );
  }

  Map<String, dynamic> toJson() => {
        'id': id,
        'teacher_id': teacherId,
        'title': title,
        'slug': slug,
        'description': description,
        'short_description': shortDescription,
        'subject': subject,
        'grade': grade,
        'thumbnail_url': thumbnailUrl,
        'price': price,
        'currency': currency,
        'is_free': isFree,
        'is_published': isPublished,
        'language': language,
        'level': level.name,
        'total_duration_seconds': totalDurationSeconds,
        'total_lessons': totalLessons,
        'rating_avg': ratingAvg,
        'enrollment_count': enrollmentCount,
        'created_at': createdAt?.toIso8601String(),
        'updated_at': updatedAt?.toIso8601String(),
      };

  static CourseLevel _parseLevel(String value) {
    switch (value.toLowerCase()) {
      case 'intermediate':
        return CourseLevel.intermediate;
      case 'advanced':
        return CourseLevel.advanced;
      default:
        return CourseLevel.beginner;
    }
  }
}

String modelStringValue(Object? value, {String fallback = ''}) {
  if (value is String) return value;
  if (value == null) return fallback;
  return value.toString();
}

String? modelNullableString(Object? value) {
  if (value is String && value.trim().isNotEmpty) return value;
  return null;
}

int modelIntValue(Object? value, {int fallback = 0}) {
  if (value is int) return value;
  if (value is num) return value.toInt();
  if (value is String) return int.tryParse(value) ?? fallback;
  return fallback;
}

double modelDoubleValue(Object? value, {double fallback = 0}) {
  if (value is int) return value.toDouble();
  if (value is double) return value;
  if (value is num) return value.toDouble();
  if (value is String) return double.tryParse(value) ?? fallback;
  return fallback;
}

bool modelBoolValue(Object? value, {bool fallback = false}) {
  if (value is bool) return value;
  if (value is String) return value.toLowerCase() == 'true';
  if (value is num) return value != 0;
  return fallback;
}

DateTime? modelDateValue(Object? value) {
  if (value is! String || value.isEmpty) return null;
  return DateTime.tryParse(value);
}

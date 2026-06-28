import '../../domain/entities/paginated_courses_entity.dart';
import 'course_model.dart';

class PaginatedCoursesModel extends PaginatedCoursesEntity {
  const PaginatedCoursesModel({
    required super.courses,
    required super.nextCursor,
    required super.limit,
  });

  factory PaginatedCoursesModel.fromJson(Map<String, dynamic> json) {
    final data = json['data'];
    final courses = data is List
        ? data
            .whereType<Map<String, dynamic>>()
            .map(CourseModel.fromJson)
            .toList()
        : <CourseModel>[];

    final nextCursor = json['next_cursor'];
    return PaginatedCoursesModel(
      courses: courses,
      nextCursor:
          nextCursor is String && nextCursor.isNotEmpty ? nextCursor : null,
      limit: modelIntValue(json['limit'], fallback: 20),
    );
  }
}

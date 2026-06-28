import 'course_entity.dart';

class PaginatedCoursesEntity {
  const PaginatedCoursesEntity({
    required this.courses,
    required this.nextCursor,
    required this.limit,
  });

  final List<CourseEntity> courses;
  final String? nextCursor;
  final int limit;

  bool get hasMore => nextCursor != null && nextCursor!.isNotEmpty;
}

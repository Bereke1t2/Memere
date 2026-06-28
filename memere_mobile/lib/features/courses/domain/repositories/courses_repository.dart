import 'package:fpdart/fpdart.dart';

import '../../../../core/errors/failures.dart';
import '../entities/course_detail_entity.dart';
import '../entities/paginated_courses_entity.dart';

abstract class CoursesRepository {
  Future<Either<Failure, PaginatedCoursesEntity>> listCourses({
    int limit = 20,
    String? after,
    String? subject,
    int? grade,
  });

  Future<Either<Failure, CourseDetailEntity>> getCourseDetail(String courseId);
}

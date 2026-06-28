import 'package:fpdart/fpdart.dart';

import '../../../../core/errors/failures.dart';
import '../entities/course_detail_entity.dart';
import '../repositories/courses_repository.dart';

class GetCourseDetailUseCase {
  const GetCourseDetailUseCase(this._repository);
  final CoursesRepository _repository;

  Future<Either<Failure, CourseDetailEntity>> call(String courseId) {
    final trimmedId = courseId.trim();
    if (trimmedId.isEmpty) {
      return Future.value(
        const Left(
          ValidationFailure('Course ID is required', field: 'courseId'),
        ),
      );
    }
    return _repository.getCourseDetail(trimmedId);
  }
}

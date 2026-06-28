import 'package:fpdart/fpdart.dart';

import '../../../../core/errors/failures.dart';
import '../entities/paginated_courses_entity.dart';
import '../repositories/courses_repository.dart';

class ListCoursesParams {
  const ListCoursesParams({
    this.limit = 20,
    this.after,
    this.subject,
    this.grade,
    this.searchQuery,
  });

  final int limit;
  final String? after;
  final String? subject;
  final int? grade;
  final String? searchQuery;
}

class ListCoursesUseCase {
  const ListCoursesUseCase(this._repository);
  final CoursesRepository _repository;

  Future<Either<Failure, PaginatedCoursesEntity>> call(
    ListCoursesParams params,
  ) {
    final limit = params.limit <= 0 ? 20 : params.limit;
    return _repository.listCourses(
      limit: limit,
      after: params.after,
      subject: params.subject,
      grade: params.grade,
    );
  }
}

import 'package:fpdart/fpdart.dart';

import '../../../../core/errors/failures.dart';
import '../entities/quiz_entity.dart';
import '../repositories/quiz_repository.dart';

class ListQuizzesByCourseUseCase {
  const ListQuizzesByCourseUseCase(this._repository);
  final QuizRepository _repository;

  Future<Either<Failure, List<QuizEntity>>> call(String courseId) {
    final trimmed = courseId.trim();
    if (trimmed.isEmpty) {
      return Future.value(
        const Left(
            ValidationFailure('Course ID is required', field: 'courseId')),
      );
    }
    return _repository.listQuizzesByCourse(trimmed);
  }
}

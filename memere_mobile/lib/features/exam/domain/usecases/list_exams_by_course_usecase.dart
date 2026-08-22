import 'package:fpdart/fpdart.dart';

import '../../../../core/errors/failures.dart';
import '../entities/mock_exam_entity.dart';
import '../repositories/exam_repository.dart';

class ListExamsByCourseUseCase {
  const ListExamsByCourseUseCase(this._repository);
  final ExamRepository _repository;

  Future<Either<Failure, List<MockExamEntity>>> call(String courseId) {
    final trimmed = courseId.trim();
    if (trimmed.isEmpty) {
      return Future.value(
        const Left(
            ValidationFailure('Course ID is required', field: 'courseId')),
      );
    }
    return _repository.listExamsByCourse(trimmed);
  }
}

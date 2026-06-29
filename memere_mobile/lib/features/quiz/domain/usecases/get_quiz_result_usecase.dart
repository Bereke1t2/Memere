import 'package:fpdart/fpdart.dart';

import '../../../../core/errors/failures.dart';
import '../entities/quiz_result_entity.dart';
import '../repositories/quiz_repository.dart';

class GetQuizResultUseCase {
  const GetQuizResultUseCase(this._repository);
  final QuizRepository _repository;

  Future<Either<Failure, QuizResultEntity>> call(String attemptId) {
    final trimmed = attemptId.trim();
    if (trimmed.isEmpty) {
      return Future.value(
        const Left(
          ValidationFailure('Attempt ID is required', field: 'attemptId'),
        ),
      );
    }
    return _repository.getResult(trimmed);
  }
}

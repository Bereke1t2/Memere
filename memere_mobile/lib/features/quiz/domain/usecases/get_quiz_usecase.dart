import 'package:fpdart/fpdart.dart';

import '../../../../core/errors/failures.dart';
import '../entities/quiz_entity.dart';
import '../repositories/quiz_repository.dart';

class GetQuizUseCase {
  const GetQuizUseCase(this._repository);
  final QuizRepository _repository;

  Future<Either<Failure, QuizEntity>> call(String quizId) {
    final trimmed = quizId.trim();
    if (trimmed.isEmpty) {
      return Future.value(
        const Left(ValidationFailure('Quiz ID is required', field: 'quizId')),
      );
    }
    return _repository.getQuiz(trimmed);
  }
}

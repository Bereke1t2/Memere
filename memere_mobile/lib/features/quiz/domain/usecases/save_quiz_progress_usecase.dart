import 'package:fpdart/fpdart.dart';

import '../../../../core/errors/failures.dart';
import '../entities/quiz_answer_payload.dart';
import '../repositories/quiz_repository.dart';

class SaveQuizProgressParams {
  const SaveQuizProgressParams({
    required this.attemptId,
    required this.answers,
  });

  final String attemptId;
  final QuizAnswerPayload answers;
}

class SaveQuizProgressUseCase {
  const SaveQuizProgressUseCase(this._repository);
  final QuizRepository _repository;

  Future<Either<Failure, void>> call(SaveQuizProgressParams params) {
    final trimmed = params.attemptId.trim();
    if (trimmed.isEmpty) {
      return Future.value(
        const Left(
          ValidationFailure('Attempt ID is required', field: 'attemptId'),
        ),
      );
    }
    return _repository.saveProgress(
      attemptId: trimmed,
      answers: params.answers,
    );
  }
}

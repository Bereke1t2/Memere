import 'package:fpdart/fpdart.dart';

import '../../../../core/errors/failures.dart';
import '../entities/quiz_answer_payload.dart';
import '../entities/quiz_result_entity.dart';
import '../repositories/quiz_repository.dart';

class SubmitQuizAttemptParams {
  const SubmitQuizAttemptParams({
    required this.attemptId,
    required this.answers,
  });

  final String attemptId;
  final QuizAnswerPayload answers;
}

class SubmitQuizAttemptUseCase {
  const SubmitQuizAttemptUseCase(this._repository);
  final QuizRepository _repository;

  Future<Either<Failure, QuizResultEntity>> call(
    SubmitQuizAttemptParams params,
  ) {
    final trimmed = params.attemptId.trim();
    if (trimmed.isEmpty) {
      return Future.value(
        const Left(
          ValidationFailure('Attempt ID is required', field: 'attemptId'),
        ),
      );
    }
    return _repository.submitAttempt(
      attemptId: trimmed,
      answers: params.answers,
    );
  }
}

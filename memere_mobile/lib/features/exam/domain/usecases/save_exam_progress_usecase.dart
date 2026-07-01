import 'package:fpdart/fpdart.dart';

import '../../../../core/errors/failures.dart';
import '../entities/exam_answer_payload.dart';
import '../repositories/exam_repository.dart';

class SaveExamProgressParams {
  const SaveExamProgressParams({
    required this.attemptId,
    required this.answers,
  });

  final String attemptId;
  final ExamAnswerPayload answers;
}

class SaveExamProgressUseCase {
  const SaveExamProgressUseCase(this._repository);
  final ExamRepository _repository;

  Future<Either<Failure, void>> call(SaveExamProgressParams params) {
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

import 'package:fpdart/fpdart.dart';

import '../../../../core/errors/failures.dart';
import '../entities/exam_answer_payload.dart';
import '../entities/exam_result_entity.dart';
import '../repositories/exam_repository.dart';

class SubmitExamParams {
  const SubmitExamParams({
    required this.attemptId,
    required this.answers,
  });

  final String attemptId;
  final ExamAnswerPayload answers;
}

class SubmitExamUseCase {
  const SubmitExamUseCase(this._repository);
  final ExamRepository _repository;

  Future<Either<Failure, ExamResultEntity>> call(SubmitExamParams params) {
    final trimmed = params.attemptId.trim();
    if (trimmed.isEmpty) {
      return Future.value(
        const Left(
          ValidationFailure('Attempt ID is required', field: 'attemptId'),
        ),
      );
    }
    return _repository.submitExam(
      attemptId: trimmed,
      answers: params.answers,
    );
  }
}

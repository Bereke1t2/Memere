import 'package:fpdart/fpdart.dart';

import '../../../../core/errors/failures.dart';
import '../entities/exam_attempt_entity.dart';
import '../repositories/exam_repository.dart';

class StartExamUseCase {
  const StartExamUseCase(this._repository);
  final ExamRepository _repository;

  Future<Either<Failure, ExamAttemptEntity>> call(String examId) {
    final trimmed = examId.trim();
    if (trimmed.isEmpty) {
      return Future.value(
        const Left(ValidationFailure('Exam ID is required', field: 'examId')),
      );
    }
    return _repository.startExam(trimmed);
  }
}

import 'package:fpdart/fpdart.dart';

import '../../../../core/errors/failures.dart';
import '../entities/exam_attempt_analytics_entity.dart';
import '../repositories/exam_repository.dart';

class GetExamAnalyticsUseCase {
  const GetExamAnalyticsUseCase(this._repository);
  final ExamRepository _repository;

  Future<Either<Failure, ExamAttemptAnalyticsEntity>> call(String attemptId) {
    final trimmed = attemptId.trim();
    if (trimmed.isEmpty) {
      return Future.value(
        const Left(
          ValidationFailure('Attempt ID is required', field: 'attemptId'),
        ),
      );
    }
    return _repository.getAnalytics(trimmed);
  }
}

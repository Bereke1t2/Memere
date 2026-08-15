import 'package:fpdart/fpdart.dart';

import '../../../../core/errors/failures.dart';
import '../entities/exam_attempt_history_entity.dart';
import '../repositories/exam_repository.dart';

class ListMyExamAttemptsUseCase {
  const ListMyExamAttemptsUseCase(this._repository);
  final ExamRepository _repository;

  Future<Either<Failure, List<ExamAttemptHistoryEntity>>> call({String? examId}) {
    return _repository.listMyAttempts(examId: examId);
  }
}

import 'package:fpdart/fpdart.dart';

import '../../../../core/errors/failures.dart';
import '../entities/paginated_mock_exams_entity.dart';
import '../repositories/exam_repository.dart';

class ListMockExamsParams {
  const ListMockExamsParams({
    this.limit = 20,
    this.after,
    this.subject,
    this.grade,
  });

  final int limit;
  final String? after;
  final String? subject;
  final int? grade;
}

class ListMockExamsUseCase {
  const ListMockExamsUseCase(this._repository);
  final ExamRepository _repository;

  Future<Either<Failure, PaginatedMockExamsEntity>> call(
    ListMockExamsParams params,
  ) {
    final limit = params.limit <= 0 || params.limit > 100 ? 20 : params.limit;
    return _repository.listMockExams(
      limit: limit,
      after: params.after,
      subject: params.subject,
      grade: params.grade,
    );
  }
}

import 'package:fpdart/fpdart.dart';

import '../../../../core/errors/failures.dart';
import '../entities/student_points_entity.dart';
import '../repositories/progress_repository.dart';

class GetMyPointsUseCase {
  const GetMyPointsUseCase(this._repository);
  final ProgressRepository _repository;

  Future<Either<Failure, StudentPointsEntity>> call() {
    return _repository.getMyPoints();
  }
}

import 'package:fpdart/fpdart.dart';

import '../../../../core/errors/failures.dart';
import '../entities/student_points_entity.dart';

abstract class ProgressRepository {
  /// The signed-in student's cumulative quiz + exam points. Requires auth on the
  /// server; guests receive a failure.
  Future<Either<Failure, StudentPointsEntity>> getMyPoints();
}

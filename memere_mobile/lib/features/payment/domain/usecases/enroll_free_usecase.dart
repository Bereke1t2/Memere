import 'package:fpdart/fpdart.dart';

import '../../../../core/errors/failures.dart';
import '../repositories/payment_repository.dart';

class EnrollFreeUseCase {
  const EnrollFreeUseCase(this._repository);
  final PaymentRepository _repository;

  Future<Either<Failure, void>> call(String courseId) {
    final trimmed = courseId.trim();
    if (trimmed.isEmpty) {
      return Future.value(
        const Left(
          ValidationFailure('Course ID is required', field: 'courseId'),
        ),
      );
    }
    return _repository.enrollFree(trimmed);
  }
}

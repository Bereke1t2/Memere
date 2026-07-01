import 'package:fpdart/fpdart.dart';

import '../../../../core/errors/failures.dart';
import '../entities/enrollment_entity.dart';
import '../repositories/payment_repository.dart';

class ListEnrollmentsUseCase {
  const ListEnrollmentsUseCase(this._repository);
  final PaymentRepository _repository;

  Future<Either<Failure, List<EnrollmentEntity>>> call({int limit = 50}) {
    return _repository.listEnrollments(limit: limit);
  }
}

import 'package:fpdart/fpdart.dart';

import '../../../../core/errors/failures.dart';
import '../entities/payment_entity.dart';
import '../repositories/payment_repository.dart';

class ListPaymentsUseCase {
  const ListPaymentsUseCase(this._repository);
  final PaymentRepository _repository;

  Future<Either<Failure, List<PaymentEntity>>> call({int limit = 50}) {
    return _repository.listPayments(limit: limit);
  }
}

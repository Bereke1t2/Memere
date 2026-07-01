import 'package:fpdart/fpdart.dart';

import '../../../../core/errors/failures.dart';
import '../entities/payment_entity.dart';
import '../repositories/payment_repository.dart';

class GetPaymentStatusUseCase {
  const GetPaymentStatusUseCase(this._repository);
  final PaymentRepository _repository;

  Future<Either<Failure, PaymentEntity>> call(String paymentId) {
    final trimmed = paymentId.trim();
    if (trimmed.isEmpty) {
      return Future.value(
        const Left(
          ValidationFailure('Payment ID is required', field: 'paymentId'),
        ),
      );
    }
    return _repository.getPaymentStatus(trimmed);
  }
}

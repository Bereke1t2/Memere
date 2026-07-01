import 'package:fpdart/fpdart.dart';

import '../../../../core/errors/failures.dart';
import '../repositories/payment_repository.dart';

class CancelSubscriptionUseCase {
  const CancelSubscriptionUseCase(this._repository);
  final PaymentRepository _repository;

  Future<Either<Failure, void>> call(String subscriptionId) {
    final trimmed = subscriptionId.trim();
    if (trimmed.isEmpty) {
      return Future.value(
        const Left(
          ValidationFailure(
            'Subscription ID is required',
            field: 'subscriptionId',
          ),
        ),
      );
    }
    return _repository.cancelSubscription(trimmed);
  }
}
